package user

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"l2alchemy/circuits/merkle_proof"
	"l2alchemy/internal/config"
	bc "l2alchemy/internal/eth"
	"l2alchemy/internal/memtime"
	"l2alchemy/internal/oracle-repo"
	"l2alchemy/internal/oracle-repo/db"
	"l2alchemy/internal/oracle-repo/gnark"
	"l2alchemy/internal/oracle-repo/merkle"
	"l2alchemy/internal/oracle-repo/util"
	"log"
	"math/big"
	"math/bits"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	_ "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/backend/groth16"
	gnarkwitness "github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	// This is used for key generation, signing, verifying (outside circuit)
	eddsa "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	// Used for assigning EdDSA witness variables (inside circuit)
)

type User struct {
	cfg             *config.Config
	wg              *sync.WaitGroup
	ecdsaPrivateKey *ecdsa.PrivateKey
	ethClient       *ethclient.Client
	IPFSClient      *db.IPFSClient
	R1CS            constraint.ConstraintSystem
	PK              groth16.ProvingKey
	VK              groth16.VerifyingKey

	Name string `json:"name"`

	circuit             *merkleproof.MerkleProofCircuit
	commitmentHashBytes []byte
	nullifierBytes      []byte
	nullifierHashBytes  []byte
	secretBytes         []byte
	destinationIDBytes  []byte

	Account *gnark.Account
}

type persistedUserState struct {
	CommitmentHashBytes []byte
	NullifierBytes      []byte
	NullifierHashBytes  []byte
	SecretBytes         []byte
	DestinationIDBytes  []byte
}

type CircuitMemTime struct {
	IncProvingTime        int
	IncProvingMemoryUsage int
	CompileMemory         int
	CompileTime           int
	VerifyMemory          int
}

var issuerPrivateKey *eddsa.PrivateKey
var issuerPublicKey *eddsa.PublicKey

func init() {
	var err error
	issuerPrivateKey, err = eddsa.GenerateKey(rand.Reader)
	if err != nil {
		panic("failed to generate issuer key")
	}
	issuerPublicKey = &issuerPrivateKey.PublicKey
}

func NewUser(cfg *config.Config, ethClient *ethclient.Client, ipfsClient *db.IPFSClient, circuit *merkleproof.MerkleProofCircuit, r1cs constraint.ConstraintSystem, pk groth16.ProvingKey, vk groth16.VerifyingKey, name string, privateKey *eddsa.PrivateKey, account *gnark.Account) User {
	ecdsaPrivateKey, err := crypto.HexToECDSA(cfg.UserPK)
	if err != nil {
		panic(fmt.Errorf("failed to parse ecdsaPrivateKey: %v", err))
	}

	usr := User{
		cfg:             cfg,
		ecdsaPrivateKey: ecdsaPrivateKey,
		ethClient:       ethClient,
		IPFSClient:      ipfsClient,
		R1CS:            r1cs,
		PK:              pk,
		VK:              vk,
		Name:            name,
		circuit:         circuit,
		Account:         account,
	}
	if err := usr.loadState(); err != nil {
		log.Printf("user=%s failed to load persisted state: %v", usr.Name, err)
	}

	return usr
}

// func (u *User) ListenToBc(quit chan struct{}) {
// 	contractAddr := common.HexToAddress(u.cfg.OracleContractAddress)

// 	// filterers:
// 	oracleFilterer, err := bc.NewOracleFilterer(contractAddr, u.ethClient)
// 	if err != nil {
// 		log.Fatalf("create oracle filterer: %v", err)
// 	}
// 	registeredEventChan := make(chan *bc.OracleUserRegistered)

// 	// event subs:
// 	ctx := context.Background()
// 	registeredSub, err := oracleFilterer.WatchUserRegistered(&bind.WatchOpts{Context: ctx}, registeredEventChan)
// 	if err != nil {
// 		log.Fatalf("set up \"registered\" subscription: %v", err)
// 	}

// 	u.wg.Add(1)
// 	defer u.wg.Done()
// 	go func() {
// 		for {
// 			select {
// 			// new aggregator event:
// 			case registeredEvent := <-registeredEventChan:
// 				u.Account.Index = registeredEvent.Index
// 				u.Account.Balance = registeredEvent.Balance
// 				fmt.Printf("user=%v registered: %v\n", u.Name, u.Account)
// 			// sub errors:
// 			case err := <-registeredSub.Err():
// 				log.Fatalf("registered sub error: %v\n", err)
// 			// quit:
// 			case <-quit:
// 				fmt.Printf("blockchain listener of user %v quitting...\n", u.Name)
// 				registeredSub.Unsubscribe()
// 				close(registeredEventChan)
// 				return
// 			}
// 		}
// 	}()
// }

func (u *User) GetBalance() (uint, uint, error) {
	fmt.Printf("getting balance for user=%v ...\n", u.Name)
	oracleContractAddr := common.HexToAddress(u.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, u.ethClient)
	if err != nil {
		return 0, 0, fmt.Errorf("create contract client instance: %w", err)
	}
	chainID := big.NewInt(u.cfg.ChainID)
	trxOpts, err := bind.NewKeyedTransactorWithChainID(u.ecdsaPrivateKey, chainID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create keyed transactor: %w", err)
	}
	callOpts := &bind.CallOpts{
		Context: context.Background(),
		From:    trxOpts.From,
	}
	burnBalance, claimBalance, err := bcClient.ViewBalance(callOpts)
	if err != nil {
		return 0, 0, fmt.Errorf("call ViewBalance() function: %w", err)
	}

	if burnBalance.BitLen() > bits.UintSize {
		return 0, 0, fmt.Errorf("burn balance overflows uint")
	}
	if claimBalance.BitLen() > bits.UintSize {
		return 0, 0, fmt.Errorf("claim balance overflows uint")
	}

	return uint(burnBalance.Uint64()), uint(claimBalance.Uint64()), nil
}

func (u *User) statePath() string {
	if u == nil || u.cfg == nil {
		return ""
	}
	return u.cfg.UserStatePath
}

func (u *User) loadState() error {
	path := u.statePath()
	if path == "" {
		return nil
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("user state path is a directory")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := gob.NewDecoder(file)
	var state persistedUserState
	if err := dec.Decode(&state); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	u.commitmentHashBytes = state.CommitmentHashBytes
	u.nullifierBytes = state.NullifierBytes
	u.nullifierHashBytes = state.NullifierHashBytes
	u.secretBytes = state.SecretBytes
	u.destinationIDBytes = state.DestinationIDBytes

	return nil
}

func (u *User) saveState() error {
	path := u.statePath()
	if path == "" {
		return nil
	}
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	state := persistedUserState{
		CommitmentHashBytes: u.commitmentHashBytes,
		NullifierBytes:      u.nullifierBytes,
		NullifierHashBytes:  u.nullifierHashBytes,
		SecretBytes:         u.secretBytes,
		DestinationIDBytes:  u.destinationIDBytes,
	}
	if err := enc.Encode(state); err != nil {
		return err
	}

	return nil
}

func (u *User) newTransactOpts() (*bind.TransactOpts, error) {
	chainID := big.NewInt(u.cfg.ChainID)
	trxOpts, err := bind.NewKeyedTransactorWithChainID(u.ecdsaPrivateKey, chainID)
	if err != nil {
		return nil, err
	}
	if u.cfg.TxGasLimit > 0 {
		// Avoid RPC gas estimation when the endpoint doesn't support eth_estimateGas.
		trxOpts.GasLimit = uint64(u.cfg.TxGasLimit)
	}
	if u.cfg.TxGasPriceWei > 0 {
		// Avoid RPC gas price discovery when eth_gasPrice isn't available.
		trxOpts.GasPrice = big.NewInt(u.cfg.TxGasPriceWei)
	} else if u.cfg.TxGasFeeCapWei > 0 || u.cfg.TxGasTipCapWei > 0 {
		// Allow explicit EIP-1559 values without RPC lookups.
		if u.cfg.TxGasFeeCapWei > 0 {
			trxOpts.GasFeeCap = big.NewInt(u.cfg.TxGasFeeCapWei)
		}
		if u.cfg.TxGasTipCapWei > 0 {
			trxOpts.GasTipCap = big.NewInt(u.cfg.TxGasTipCapWei)
		}
	}
	return trxOpts, nil
}

// RegisterUserTx registers the user on-chain and waits for the receipt.
// It returns the transaction hash on success.
func (u *User) RegisterUserTx() (string, error) {
	fmt.Printf("registering user=%v ...\n", u.Name)
	oracleContractAddr := common.HexToAddress(u.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, u.ethClient)
	if err != nil {
		return "", fmt.Errorf("create oracle contract client: %w", err)
	}

	trxOpts, err := u.newTransactOpts()
	if err != nil {
		return "", fmt.Errorf("create keyed transactor: %w", err)
	}

	pendingNonce, err := u.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		return "", fmt.Errorf("fetch pending nonce: %w", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))
	if trxOpts.GasPrice == nil && trxOpts.GasFeeCap == nil {
		gasPrice, err := u.ethClient.SuggestGasPrice(context.Background())
		if err != nil {
			return "", fmt.Errorf("suggest gas price: %w", err)
		}
		trxOpts.GasPrice = gasPrice
	}

	pk := gnark.PublicKeyToOraclePublicKey(u.Account.PublicKey)
	tx, err := bcClient.RegisterUser(trxOpts, *pk)
	if err != nil {
		return "", fmt.Errorf("call RegisterUser(): %w", err)
	}

	receipt, err := bind.WaitMined(context.Background(), u.ethClient, tx)
	if err != nil {
		return "", fmt.Errorf("wait for tx mined: %w", err)
	}
	bc.LogTxReceipt(fmt.Sprintf("register user name=%s", u.Name), tx, receipt)
	if receipt.Status != 1 {
		return "", fmt.Errorf("transaction reverted (tx=%s)", tx.Hash().Hex())
	}

	log.Printf("successfully registered user=%v (tx=%s)", u.Name, tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func (u *User) getLatestIPFSHashView() (string, error) {
	// fmt.Printf("getting latest ipfs hash for user=%v ...\n", u.Name)
	oracleContractAddr := common.HexToAddress(u.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, u.ethClient)
	if err != nil {
		log.Fatalf("create contract client instance: %v", err)
	}
	chainID := big.NewInt(u.cfg.ChainID)
	trxOpts, err := bind.NewKeyedTransactorWithChainID(u.ecdsaPrivateKey, chainID)
	if err != nil {
		log.Fatalf("failed to create keyed transactor: %v", err)
	}
	callOpts := &bind.CallOpts{
		Context: context.Background(),
		From:    trxOpts.From,
	}
	latestIPFSHash, err := bcClient.ViewLatestIPFSHash(callOpts)
	if err != nil {
		log.Fatalf("call ViewLatestIPFSHash() function: %v", err)
	}

	return latestIPFSHash, nil
}

// destinationID returns the per-deployment destination-rollup identifier
// d_dst (paper SIV-E; F-24): configured per Gateway instance, used when
// forming the burn commitment C = H(n_rd || s_rd || d_dst).
func (u *User) destinationID() (*big.Int, error) {
	id, ok := new(big.Int).SetString(u.cfg.DestinationID, 10)
	if !ok {
		return nil, fmt.Errorf("invalid DESTINATION_ID %q", u.cfg.DestinationID)
	}
	return id, nil
}

// isPublishedCommitmentRoot checks the Gateway's record of commitment roots
// (paper SIV-E Steps 6-7).
func (u *User) isPublishedCommitmentRoot(root *big.Int) (bool, error) {
	oracleContractAddr := common.HexToAddress(u.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, u.ethClient)
	if err != nil {
		return false, fmt.Errorf("create contract client instance: %w", err)
	}
	callOpts := &bind.CallOpts{Context: context.Background()}
	return bcClient.IsPublishedCommitmentRoot(callOpts, root)
}

func (u *User) BurnTx() (string, string, error) {
	log.Printf("user=%v is burning...", u.Name)
	preBurnBalance, preClaimBalance, err := u.GetBalance()
	if err != nil {
		return "", "", fmt.Errorf("fetch pre-burn balance (user=%v): %w", u.Name, err)
	}
	oracleContractAddr := common.HexToAddress(u.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, u.ethClient)
	if err != nil {
		return "", "", fmt.Errorf("create contract client instance: %w", err)
	}
	trxOpts, err := u.newTransactOpts()
	if err != nil {
		return "", "", fmt.Errorf("create keyed transactor: %w", err)
	}

	pendingNonce, err := u.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch pending nonce: %w", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))
	if trxOpts.GasPrice == nil && trxOpts.GasFeeCap == nil {
		gasPrice, err := u.ethClient.SuggestGasPrice(context.Background())
		if err != nil {
			return "", "", fmt.Errorf("suggest gas price: %w", err)
		}
		trxOpts.GasPrice = gasPrice
	}

	// calculate the commitment hash:
	destinationID, err := u.destinationID()
	if err != nil {
		return "", "", err
	}
	commitmentHashBytes, nullifierBytes, nullifierHashBytes, secretBytes, destinationIDBytes := CalculateCommitmentHash(destinationID)

	// burn the amount:
	var commitmentHash [32]byte
	copy(commitmentHash[:], commitmentHashBytes)
	tx, err := bcClient.Burn(trxOpts, commitmentHash)
	if err != nil {
		return "", "", fmt.Errorf("call Burn() function: %w", err)
	}

	receipt, err := bind.WaitMined(context.Background(), u.ethClient, tx)
	if err != nil {
		return "", "", fmt.Errorf("failed to wait for transaction mining: %w", err)
	}
	bc.LogTxReceipt(fmt.Sprintf("burn user=%s", u.Name), tx, receipt)
	if receipt.Status != 1 {
		return "", "", fmt.Errorf("transaction reverted (tx=%s)", tx.Hash().Hex())
	}
	log.Printf("successfully burned from user=%v", u.Name)

	postBurnBalance, postClaimBalance, err := u.GetBalance()
	if err != nil {
		return "", "", fmt.Errorf("fetch post-burn balance (user=%v): %w", u.Name, err)
	}
	if postBurnBalance >= preBurnBalance {
		return "", "", fmt.Errorf("burn did not decrease burn balance (user=%v before=%d after=%d)", u.Name, preBurnBalance, postBurnBalance)
	}
	if postClaimBalance != preClaimBalance {
		return "", "", fmt.Errorf("burn unexpectedly changed claim balance (user=%v before=%d after=%d)", u.Name, preClaimBalance, postClaimBalance)
	}
	log.Printf("burn balances user=%v burn=%d->%d claim=%d", u.Name, preBurnBalance, postBurnBalance, postClaimBalance)

	// save the calculated values:
	u.commitmentHashBytes = commitmentHashBytes[:]
	u.nullifierBytes = nullifierBytes[:]
	u.nullifierHashBytes = nullifierHashBytes[:]
	u.secretBytes = secretBytes[:]
	u.destinationIDBytes = destinationIDBytes[:]
	if err := u.saveState(); err != nil {
		return "", "", fmt.Errorf("failed to persist user state (user=%v): %w", u.Name, err)
	}

	return tx.Hash().Hex(), common.BytesToHash(commitmentHashBytes).Hex(), nil
}

func (u *User) WithdrawTx() (string, error) {
	preBurnBalance, preClaimBalance, err := u.GetBalance()
	if err != nil {
		return "", fmt.Errorf("fetch pre-withdraw balance (user=%v): %w", u.Name, err)
	}
	if len(u.commitmentHashBytes) == 0 {
		if err := u.loadState(); err != nil {
			return "", fmt.Errorf("failed to load persisted user state (user=%v): %w", u.Name, err)
		}
	}
	if len(u.commitmentHashBytes) == 0 {
		return "", fmt.Errorf("no commitment hash available for withdrawal (user=%v)", u.Name)
	}
	// fetch from ipfs:
	latestIPFSHash, err := u.getLatestIPFSHashView()
	if err != nil {
		return "", fmt.Errorf("failed to get the latest ipfs hash (user=%v): %v", u.Name, err)
	}
	content, err := u.IPFSClient.Download(latestIPFSHash)
	if err != nil {
		return "", fmt.Errorf("failed to download from ipfs (user=%v): %v", u.Name, err)
	}
	ipfsContent, err := oracle.DeserializeIPFSContent(content)
	if err != nil {
		return "", fmt.Errorf("failed to deserialize ipfs content (user=%v): %v", u.Name, err)
	}
	// search ipfs for the commitment hash:
	incVote, ok := ipfsContent.CommitmentHashIncVote[string(u.commitmentHashBytes)]
	if !ok {
		return "", fmt.Errorf("commitment hash not found in ipfs (user=%v hash=%s)", u.Name, common.BytesToHash(u.commitmentHashBytes).Hex())
	}
	tree := ipfsContent.IncMerkleTree
	if reflect.DeepEqual(tree, merkle.IncrementalMerkleTree{}) {
		return "", fmt.Errorf("empty inc merkle tree in ipfs (user=%v)", u.Name)
	}

	// setup to generate:

	proofIndex := incVote.IncTreeIndex
	depth := ipfsContent.IncMerkleTree.Depth
	merkleRoot, proofPath, err := tree.GetProofPath(proofIndex)
	if err != nil {
		return "", fmt.Errorf("failed to get proof path from inc merkle tree (user=%v): %v", u.Name, err)
	}

	// The Gateway must have recorded this commitment root (paper SIV-E
	// Steps 6-7); otherwise the fetched DFS content cannot be trusted as a
	// basis for the Redeeming proof.
	rootPublished, err := u.isPublishedCommitmentRoot(new(big.Int).SetBytes(merkleRoot))
	if err != nil {
		return "", fmt.Errorf("failed to validate commitment root (user=%v): %w", u.Name, err)
	}
	if !rootPublished {
		return "", fmt.Errorf("DFS commitment root not recorded by the gateway (user=%v root=%s)", u.Name, common.BytesToHash(merkleRoot).Hex())
	}

	var witness merkleproof.MerkleProofCircuit

	witness.Nullifier = new(big.Int).SetBytes(u.nullifierBytes)
	witness.Secret = new(big.Int).SetBytes(u.secretBytes)
	witness.DestinationID = new(big.Int).SetBytes(u.destinationIDBytes)

	witness.Leaf = proofIndex
	// witness.CommitmentHash = u.commitmentHashBytes
	witness.NullifierHash = u.nullifierHashBytes
	witness.M.RootHash = merkleRoot

	witness.M.Path = make([]frontend.Variable, depth+1)
	for i := 0; i < depth+1; i++ {
		witness.M.Path[i] = frontend.Variable(new(big.Int).SetBytes(proofPath[i]))
	}

	fmt.Println("Debugging Witness Before Proof Generation:")
	fmt.Printf("Nullifier: %x\n", witness.Nullifier)
	fmt.Printf("Secret: %x\n", witness.Secret)
	fmt.Printf("DestinationID: %x\n", witness.DestinationID)
	fmt.Printf("Leaf Index: %d\n", witness.Leaf)
	// fmt.Printf("Commitment Hash: %x\n", witness.CommitmentHash)
	fmt.Printf("Nullifier Hash: %x\n", witness.NullifierHash)

	fmt.Printf("Merkle Root: %x\n", witness.M.RootHash)
	for i, path := range witness.M.Path {
		fmt.Printf("Merkle Path[%d]: %x\n", i, path)
	}

	var (
		fullWitness gnarkwitness.Witness
		proof       groth16.Proof
	)

	_, err = memtime.MeasurePeak(
		"groth16.Prove merkle_proof",
		memtime.PeakSampleInterval,
		func() error {
			var e error

			fullWitness, e = frontend.NewWitness(&witness, ecc.BN254.ScalarField())
			if e != nil {
				return e
			}

			proof, e = groth16.Prove(u.R1CS, u.PK, fullWitness)
			return e
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to prove (user=%v): %v", u.Name, err)
	}

	publicWitness, err := frontend.NewWitness(&witness, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		return "", fmt.Errorf("failed to create public witness (user=%v): %v", u.Name, err)
	}

	_, err = memtime.MeasurePeak(
		"groth16.Verify merkle_proof",
		memtime.PeakSampleInterval,
		func() error {
			return groth16.Verify(proof, u.VK, publicWitness)
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to verify proof (user=%v): %v", u.Name, err)
	}

	// send trx to claim:
	log.Printf("claiming (user=%v)...", u.Name)
	oracleContractAddr := common.HexToAddress(u.cfg.OracleContractAddress)
	bcClient, err := bc.NewOracle(oracleContractAddr, u.ethClient)
	if err != nil {
		return "", fmt.Errorf("create contract client instance: %w", err)
	}
	trxOpts, err := u.newTransactOpts()
	if err != nil {
		return "", fmt.Errorf("failed to create keyed transactor: %w", err)
	}

	pendingNonce, err := u.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		return "", fmt.Errorf("failed to fetch pending nonce: %w", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))
	if trxOpts.GasPrice == nil && trxOpts.GasFeeCap == nil {
		gasPrice, err := u.ethClient.SuggestGasPrice(context.Background())
		if err != nil {
			return "", fmt.Errorf("suggest gas price: %w", err)
		}
		trxOpts.GasPrice = gasPrice
	}

	var buf bytes.Buffer
	_, err = proof.WriteTo(&buf)
	if err != nil {
		return "", fmt.Errorf("failed to serialize proof: %w", err)
	}
	proofBytes := buf.Bytes()

	publicWitnessBytes, err := publicWitness.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("failed to marshal public witness (user=%v): %v", u.Name, err)
	}

	tx, err := bcClient.Claim(trxOpts, proofBytes, publicWitnessBytes, [32]byte(u.nullifierHashBytes))
	if err != nil {
		return "", fmt.Errorf("call Claim() function: %w", err)
	}

	receipt, err := bind.WaitMined(context.Background(), u.ethClient, tx)
	if err != nil {
		return "", fmt.Errorf("failed to wait for transaction mining: %w", err)
	}
	bc.LogTxReceipt(fmt.Sprintf("claim user=%s", u.Name), tx, receipt)
	if receipt.Status == 1 {
		log.Printf("successfully sent claim tx (user=%v)", u.Name)
	} else {
		return "", fmt.Errorf("transaction reverted (tx=%s)", tx.Hash().Hex())
	}

	postBurnBalance, postClaimBalance, err := u.GetBalance()
	if err != nil {
		return "", fmt.Errorf("fetch post-withdraw balance (user=%v): %w", u.Name, err)
	}
	if postBurnBalance != preBurnBalance {
		return "", fmt.Errorf("withdraw unexpectedly changed burn balance (user=%v before=%d after=%d)", u.Name, preBurnBalance, postBurnBalance)
	}
	// Claims are pending until the committee finalizes the round: tokens are
	// minted by submitWiVote after the Aggregating proof verifies (paper
	// SIV-E Step 17), so the claim balance must be unchanged here.
	if postClaimBalance != preClaimBalance {
		return "", fmt.Errorf("claim credited before batch finalization (user=%v before=%d after=%d)", u.Name, preClaimBalance, postClaimBalance)
	}
	log.Printf("claim submitted, pending batch verification: user=%v burn=%d claim=%d", u.Name, postBurnBalance, postClaimBalance)

	return tx.Hash().Hex(), nil
}

// AwaitClaimCredit polls the user's claim balance until it exceeds the given
// pre-claim balance (i.e. the round containing the claim was finalized and
// minted by submitWiVote) or the timeout elapses.
func (u *User) AwaitClaimCredit(preClaimBalance uint, timeout time.Duration) (uint, error) {
	deadline := time.Now().Add(timeout)
	for {
		_, claimBalance, err := u.GetBalance()
		if err != nil {
			return 0, fmt.Errorf("fetch claim balance (user=%v): %w", u.Name, err)
		}
		if claimBalance > preClaimBalance {
			return claimBalance, nil
		}
		if time.Now().After(deadline) {
			return claimBalance, fmt.Errorf("claim not credited within %s (user=%v balance=%d)", timeout, u.Name, claimBalance)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// NullifierHashBytes computes H(n_rd) over the canonical 32-byte field
// encoding of the nullifier, matching the in-circuit MiMC hash of the field
// element. Hashing the unpadded nullifier bytes diverges for nullifiers
// with leading zero bytes (~1/256 of runs; F-25).
func NullifierHashBytes(nullifier *big.Int) []byte {
	h := hash.MIMC_BN254.New()
	h.Reset()
	h.Write(util.PadTo32Bytes(nullifier))
	return h.Sum(nil)
}

// CalculateCommitmentHash draws a fresh (nullifier, secret) pair and forms
// the burn commitment C = H(n_rd || s_rd || d_dst) (paper SIV-E). The
// destination-rollup identifier d_dst is a per-deployment configuration
// value (F-24), one per Gateway instance.
func CalculateCommitmentHash(destinationID *big.Int) (commitmentHashBytes, nullifierBytes, nullifierHashBytes, secretBytes, destinationIDBytes []byte) {
	mod := ecc.BN254.ScalarField()
	nullifier, _ := util.GenerateRandomBigInt32Bytes(mod)
	secret, _ := util.GenerateRandomBigInt32Bytes(mod)

	mimcHash := hash.MIMC_BN254.New()
	mimcHash.Reset()
	mimcHash.Write(util.PadTo32Bytes(nullifier))
	mimcHash.Write(util.PadTo32Bytes(secret))
	mimcHash.Write(util.PadTo32Bytes(destinationID))
	commitmentHash := mimcHash.Sum(nil)

	nullifierHash := NullifierHashBytes(nullifier)

	commitmentHashBytes = []byte(commitmentHash)
	nullifierBytes = []byte(util.PadTo32Bytes(nullifier))
	nullifierHashBytes = []byte(nullifierHash)
	secretBytes = []byte(util.PadTo32Bytes(secret))
	destinationIDBytes = []byte(util.PadTo32Bytes(destinationID))

	return commitmentHash, nullifierBytes, nullifierHashBytes, secretBytes, destinationIDBytes
}
