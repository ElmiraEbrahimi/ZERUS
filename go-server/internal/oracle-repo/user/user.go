package user

import (
	merkleproof "l2alchemy/circuits/merkle_proof"
	"l2alchemy/internal/config"
	bc "l2alchemy/internal/eth"
	"l2alchemy/internal/oracle-repo"
	"l2alchemy/internal/oracle-repo/db"
	"l2alchemy/internal/oracle-repo/gnark"
	"l2alchemy/internal/oracle-repo/merkle"
	"l2alchemy/internal/oracle-repo/util"

	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"reflect"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	// This is used for key generation, signing, verifying (outside circuit)
	eddsa "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	// Used for assigning EdDSA witness variables (inside circuit)
	"github.com/consensys/gnark-crypto/ecc"
	_ "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/hash"
)

type User struct {
	cfg             *config.Config
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

	GasCosts *GasCosts
}

type GasCosts struct {
	BurnCost  uint64
	ClaimCost uint64
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

	return User{
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
		GasCosts:        &GasCosts{},
	}
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

func (u *User) GetBalance() (uint, uint) {
	fmt.Printf("getting balance for user=%v ...\n", u.Name)
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
	tokenOneBalance, tokenTwoBalance, err := bcClient.ViewBalance(callOpts)
	if err != nil {
		log.Fatalf("call ViewBalance() function: %v", err)
	}

	return uint(tokenOneBalance.Uint64()), uint(tokenTwoBalance.Uint64())
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

	chainID := big.NewInt(u.cfg.ChainID)
	trxOpts, err := bind.NewKeyedTransactorWithChainID(u.ecdsaPrivateKey, chainID)
	if err != nil {
		return "", fmt.Errorf("create keyed transactor: %w", err)
	}

	pendingNonce, err := u.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		return "", fmt.Errorf("fetch pending nonce: %w", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))
	gasPrice, err := u.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("suggest gas price: %w", err)
	}
	trxOpts.GasPrice = gasPrice
	trxOpts.GasLimit = 300_000

	pk := gnark.PublicKeyToOraclePublicKey(u.Account.PublicKey)
	tx, err := bcClient.RegisterUser(trxOpts, *pk)
	if err != nil {
		return "", fmt.Errorf("call RegisterUser(): %w", err)
	}

	receipt, err := bind.WaitMined(context.Background(), u.ethClient, tx)
	if err != nil {
		return "", fmt.Errorf("wait for tx mined: %w", err)
	}
	if receipt.Status != 1 {
		return "", fmt.Errorf("transaction reverted (tx=%s)", tx.Hash().Hex())
	}

	fmt.Printf("successfully registered user=%v (tx=%s)\n", u.Name, tx.Hash().Hex())
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

func (u *User) BurnTx() error {
	fmt.Printf("user=%v is burning...\n", u.Name)
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

	pendingNonce, err := u.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))
	gasPrice, err := u.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed to suggest gas price: %v", err)
	}
	trxOpts.GasPrice = gasPrice

	// calculate the commitment hash:
	commitmentHashBytes, nullifierBytes, nullifierHashBytes, secretBytes, destinationIDBytes := CalculateCommitmentHash()

	// burn the amount:
	var commitmentHash [32]byte
	copy(commitmentHash[:], commitmentHashBytes)
	tx, err := bcClient.Burn(trxOpts, commitmentHash)
	if err != nil {
		log.Fatalf("call Burn() function: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), u.ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for transaction mining: %v", err)
	}
	if receipt.Status == 1 {
		fmt.Printf("successfully burned from user=%v\n", u.Name)
	} else {
		fmt.Printf("Transaction failed (user=%v\n)", u.Name)
	}

	// save the calculated values:
	u.commitmentHashBytes = commitmentHashBytes[:]
	u.nullifierBytes = nullifierBytes[:]
	u.nullifierHashBytes = nullifierHashBytes[:]
	u.secretBytes = secretBytes[:]
	u.destinationIDBytes = destinationIDBytes[:]

	u.GasCosts.BurnCost = receipt.GasUsed

	return nil
}

func (u *User) WithdrawTx() error {
	// fetch from ipfs:
	latestIPFSHash, err := u.getLatestIPFSHashView()
	if err != nil {
		return fmt.Errorf("failed to get the latest ipfs hash (user=%v): %v", u.Name, err)
	}
	content, err := u.IPFSClient.Download(latestIPFSHash)
	if err != nil {
		return fmt.Errorf("failed to download from ipfs (user=%v): %v", u.Name, err)
	}
	ipfsContent, err := oracle.DeserializeIPFSContent(content)
	if err != nil {
		return fmt.Errorf("failed to deserialize ipfs content (user=%v): %v", u.Name, err)
	}
	// search ipfs for the commitment hash:
	incVote, ok := ipfsContent.CommitmentHashIncVote[string(u.commitmentHashBytes)]
	if !ok {
		return fmt.Errorf("commitment hash not found in ipfs (user=%v)", u.Name)
	}
	tree := ipfsContent.IncMerkleTree
	if reflect.DeepEqual(tree, merkle.IncrementalMerkleTree{}) {
		return fmt.Errorf("empty inc merkle tree in ipfs (user=%v)", u.Name)
	}

	// setup to generate:

	proofIndex := incVote.IncTreeIndex
	depth := ipfsContent.IncMerkleTree.Depth
	merkleRoot, proofPath, err := tree.GetProofPath(proofIndex)
	if err != nil {
		return fmt.Errorf("failed to get proof path from inc merkle tree (user=%v): %v", u.Name, err)
	}

	var witness merkleproof.MerkleProofCircuit

	witness.Nullifier = new(big.Int).SetBytes(u.nullifierBytes)
	witness.Secret = new(big.Int).SetBytes(u.secretBytes)
	witness.DestinationID = new(big.Int).SetBytes(u.destinationIDBytes)

	witness.Leaf = proofIndex
	witness.NullifierHash = u.nullifierHashBytes
	witness.M.RootHash = merkleRoot

	for i := 0; i < depth+1; i++ {
		witness.M.Path[i] = frontend.Variable(new(big.Int).SetBytes(proofPath[i]))
	}

	// generate the proof
	fullWitness, err := frontend.NewWitness(&witness, ecc.BN254.ScalarField())
	if err != nil {
		return fmt.Errorf("failed to create witness (user=%v): %v", u.Name, err)
	}

	fmt.Println("Debugging Witness Before Proof Generation:")
	fmt.Printf("Nullifier: %x\n", witness.Nullifier)
	fmt.Printf("Secret: %x\n", witness.Secret)
	fmt.Printf("DestinationID: %x\n", witness.DestinationID)
	fmt.Printf("Leaf Index: %d\n", witness.Leaf)
	fmt.Printf("Nullifier Hash: %x\n", witness.NullifierHash)
	fmt.Printf("Merkle Root: %x\n", witness.M.RootHash)
	for i, path := range witness.M.Path {
		fmt.Printf("Merkle Path[%d]: %x\n", i, path)
	}

	proof, err := groth16.Prove(u.R1CS, u.PK, fullWitness)
	if err != nil {
		return fmt.Errorf("failed to generate Groth16 proof (user=%v): %v", u.Name, err)
	}

	publicWitness, err := frontend.NewWitness(&witness, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		return fmt.Errorf("failed to create public witness (user=%v): %v", u.Name, err)
	}

	_ = groth16.Verify(proof, u.VK, publicWitness)

	// send trx to claim:
	fmt.Printf("claiming (user=%v) ...\n", u.Name)
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

	pendingNonce, err := u.ethClient.PendingNonceAt(context.Background(), trxOpts.From)
	if err != nil {
		log.Fatalf("failed to fetch pending nonce: %v", err)
	}
	trxOpts.Nonce = big.NewInt(int64(pendingNonce))
	gasPrice, err := u.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed to suggest gas price: %v", err)
	}
	trxOpts.GasPrice = gasPrice

	var buf bytes.Buffer
	_, err = proof.WriteTo(&buf)
	if err != nil {
		log.Fatal("Failed to serialize proof:", err)
	}
	proofBytes := buf.Bytes()

	publicWitnessBytes, err := publicWitness.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal public witness (user=%v): %v", u.Name, err)
	}

	tx, err := bcClient.Claim(trxOpts, proofBytes, publicWitnessBytes, [32]byte(u.nullifierHashBytes))
	if err != nil {
		log.Fatalf("call Claim() function: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), u.ethClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for transaction mining: %v", err)
	}
	if receipt.Status == 1 {
		fmt.Printf("successfully sent claim trx (user=%v)\n", u.Name)
	} else {
		fmt.Printf("Transaction failed (user=%v\n)", u.Name)
	}

	tokenOneBalance, tokenTwoBalance := u.GetBalance()
	fmt.Printf("balance of user=%v: tokenOne=%v,tokenTwo=%v\n", u.Name, tokenOneBalance, tokenTwoBalance)

	u.GasCosts.ClaimCost = receipt.GasUsed

	return nil
}

func CalculateCommitmentHash() (commitmentHashBytes, nullifierBytes, nullifierHashBytes, secretBytes, destinationIDBytes []byte) {
	hGo := hash.MIMC_BN254.New()
	mod := ecc.BN254.ScalarField()
	nullifier, _ := util.GenerateRandomBigInt32Bytes(mod)
	secret, _ := util.GenerateRandomBigInt32Bytes(mod)

	mimcHash := hash.MIMC_BN254.New()
	mimcHash.Reset()
	mimcHash.Write(util.PadTo32Bytes(nullifier))
	mimcHash.Write(util.PadTo32Bytes(secret))
	destinationID := new(big.Int)
	destinationID.SetString("9636219578937187601590327046728695236698322465209974782280717458744997515735", 10)
	mimcHash.Write(util.PadTo32Bytes(destinationID))
	commitmentHash := mimcHash.Sum(nil)

	hGo.Reset()
	hGo.Write(nullifier.Bytes())
	nullifierHash := hGo.Sum(nil)

	commitmentHashBytes = []byte(commitmentHash)
	nullifierBytes = []byte(util.PadTo32Bytes(nullifier))
	nullifierHashBytes = []byte(nullifierHash)
	secretBytes = []byte(util.PadTo32Bytes(secret))
	destinationIDBytes = []byte(util.PadTo32Bytes(destinationID))

	return commitmentHash, nullifierBytes, nullifierHashBytes, secretBytes, destinationIDBytes
}
