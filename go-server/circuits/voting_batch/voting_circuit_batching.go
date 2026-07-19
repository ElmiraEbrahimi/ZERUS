package votingbatch

import (
	"fmt"

	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/math/bits"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const (
	RewardAggregator = 500000000000000
	RewardValidator  = 20000000000
	PenaltyValidator = 10000000000

	cost1 = 1
)

type BatchingVotingCircuit struct {
	PreStateRoot       frontend.Variable `gnark:",public"`
	ResultingStateRoot frontend.Variable `gnark:",public"`
	//*** RoundID, BatchCommitment, WithdrawalReqIDs are newly added ***
	RoundID          frontend.Variable   `gnark:",public"` // Public round index
	BatchCommitment  frontend.Variable   `gnark:",public"` // Commitment hash of private WithdrawalReqIDs
	MajorityVote     frontend.Variable   `gnark:",public"` // Decimal value of a batch-sized vote bitmask
	ValidatorBits    frontend.Variable   `gnark:",public"`
	HonestBits       frontend.Variable   `gnark:",public"`
	WithdrawalReqIDs []frontend.Variable // Slice
	Aggregator       BatchingAggregatorConstraints
	Validators       []BatchingValidatorConstraints
}

type BatchingAggregatorConstraints struct {
	Index       frontend.Variable    `gnark:",public"`
	PreSeed     twistededwards.Point `gnark:",public"`
	PostSeed    twistededwards.Point `gnark:",public"`
	SecretKey   frontend.Variable
	Balance     frontend.Variable
	MerkleProof MerkleProofW
}

type BatchingValidatorConstraints struct {
	Index       frontend.Variable
	PublicKey   eddsa.PublicKey
	Balance     frontend.Variable
	MerkleProof MerkleProofW
	Signature   eddsa.Signature
	Vote        frontend.Variable // Validator's decision as decimal of a bitmask for b=5 vote 31 means 11111

}

func powbatching(api frontend.API, x frontend.Variable, y frontend.Variable) frontend.Variable {
	output := frontend.Variable(1)
	b := bits.ToBinary(api, y, bits.WithNbDigits(256))
	for i := 0; i < len(b); i++ {
		if i != 0 {
			output = api.Mul(output, output)
		}
		multiply := api.Mul(output, x)
		output = api.Select(b[len(b)-1-i], multiply, output)
	}
	return output
}

func (c *BatchingVotingCircuit) Define(api frontend.API) error {

	api.Println(" DUMPING VOTING CIRCUIT BATCHING INPUTS FORM OUTSIDE ")

	api.Println("[batching_out] PreStateRoot:", c.PreStateRoot)
	api.Println("[batching_out] RoundID:", c.RoundID)
	api.Println("[batching_out]  MajorityVote:", c.MajorityVote)
	api.Println("[batching_out_circuit] ***** Aggregator Index:", c.Aggregator.Index)
	api.Println("[batching_out]  Aggregator PreSeed:", c.Aggregator.PreSeed)
	api.Println("[batching_out]  Aggregator Balance:", c.Aggregator.Balance)
	api.Println("[batching_out]  Aggregator MerkleProof Path:", c.Aggregator.MerkleProof.Path)

	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return fmt.Errorf("curve initialization: %w", err)
	}

	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return fmt.Errorf("hash function initialization: %w", err)
	}

	// Ensure unique validator IDs
	for i := 0; i < len(c.Validators); i++ {
		for j := 0; j < len(c.Validators); j++ {
			if i == j {
				continue
			}
			api.AssertIsDifferent(c.Validators[i].Index, c.Validators[j].Index)
		}
	}

	// *** Recompute withdrawal request hash commitment ***
	hFunc.Reset()
	for i := 0; i < len(c.WithdrawalReqIDs); i++ {
		api.Println("[batching_circuit_ids] WithdrawalReqIDs[", i, "]:", c.WithdrawalReqIDs[i])
		hFunc.Write(c.WithdrawalReqIDs[i])
	}

	api.Println(" [batching_out] ***BatchCommitment :", c.BatchCommitment)
	api.Println("[batching_circuit] ***BatchCommitment:", hFunc.Sum())
	api.AssertIsEqual(c.BatchCommitment, hFunc.Sum()) // Bind proof to specific batch (public commitment)

	// Compute next seed for the aggregator
	seedAfter := curve.ScalarMul(c.Aggregator.PreSeed, c.Aggregator.SecretKey)

	api.Println("[batching_out] Aggregator PostSeed:", c.Aggregator.PostSeed)
	api.Println("[batching_circuit]  Aggregator expected PostSeed:", seedAfter)

	api.AssertIsEqual(c.Aggregator.PostSeed.X, seedAfter.X)
	api.AssertIsEqual(c.Aggregator.PostSeed.Y, seedAfter.Y)

	// Compute aggregator public key
	base := curve.Params().Base
	basePoint := twistededwards.Point{X: base[0], Y: base[1]}
	aggregatorPubKey := curve.ScalarMul(basePoint, c.Aggregator.SecretKey)

	api.Println("[batching_circuit] Aggregator expected pubkey:", aggregatorPubKey)
	curve.AssertIsOnCurve(aggregatorPubKey)

	// Verify the aggregator Merkle proof
	hFunc.Reset()
	hFunc.Write(c.Aggregator.Index)
	hFunc.Write(aggregatorPubKey.X)
	hFunc.Write(aggregatorPubKey.Y)
	hFunc.Write(c.Aggregator.Balance)

	api.Println("[batching_out]  Aggregator MerkleProof Root:", c.Aggregator.MerkleProof.RootHash)
	api.Println("[batching_out]  Aggregator MerkleProof Path[0]:", c.Aggregator.MerkleProof.Path[0])
	api.Println("[batching_circuit] Aggregator computed leaf hash:", hFunc.Sum())

	// Keep the witness root consistent with the public pre-state root. The
	// actual membership check below is against PreStateRoot, which the Gateway
	// supplies from its current stored validator-state root.
	api.AssertIsEqual(c.Aggregator.MerkleProof.RootHash, c.PreStateRoot)
	api.AssertIsEqual(hFunc.Sum(), c.Aggregator.MerkleProof.Path[0])
	hFunc.Reset()
	// The aggregator's membership seeds the R_int chain from the pre-round
	// state root (paper Alg. 2 line 8).
	c.Aggregator.MerkleProof.VerifyProof(api, hFunc, c.Aggregator.Index, c.PreStateRoot)
	api.Println("[batching_circuit] Aggregator Verified Merkle proof...")
	// Reward the aggregator for the batch
	hFunc.Reset()
	hFunc.Write(c.Aggregator.Index)
	hFunc.Write(aggregatorPubKey.X)
	hFunc.Write(aggregatorPubKey.Y)
	hFunc.Write(api.Add(c.Aggregator.Balance, RewardAggregator))
	c.Aggregator.MerkleProof.Path[0] = hFunc.Sum()

	hFunc.Reset()
	intermediateRoot := c.Aggregator.MerkleProof.ComputeRootFromPath(api, hFunc, c.Aggregator.Index)
	api.Println("[circuit] Aggregator computed root hash:", intermediateRoot)
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// for checking the aggregator majorityvote correctness
	majorityCount := frontend.Variable(0)

	validatorBits := frontend.Variable(0)
	honestBits := frontend.Variable(0)

	// Process validators
	for _, validator := range c.Validators {
		api.Println("*****Starting validator loop******")
		api.Println("  *****[batching_out_Circuit] Validator Index:", validator.Index)
		api.Println("  [batching_out] Validator PublicKey:", validator.PublicKey)
		api.Println("  [batching_out] Validator Balance:", validator.Balance)
		api.Println("  [batching_out] Validator Vote:", validator.Vote)
		api.Println("  [batching_out] Validator MerkleProof Path:", validator.MerkleProof.Path)
		api.Println("  [batching_out] Validator MerkleProof root hash:", validator.MerkleProof.RootHash)
		api.Println("  [batching_out] Validator Signature:", validator.Signature)

		hFunc.Reset()
		api.Println("***starting valiator leaf hash computation***")
		hFunc.Write(validator.Index)
		hFunc.Write(validator.PublicKey.A.X)
		hFunc.Write(validator.PublicKey.A.Y)
		hFunc.Write(validator.Balance)

		api.Println("[batching_circuit] Validator LeafHash  Index", validator.Index)
		api.Println("[batching_circuit] Validator LeafHash  PubKey.X", validator.PublicKey.A.X)
		api.Println("[batching_circuit] Validator LeafHash  PubKey.Y", validator.PublicKey.A.Y)
		api.Println("[batching_circuit] Validator LeafHash  Balance", validator.Balance)

		api.Println(validator.Index, "[batching_out]  Validator MerkleProof Root:", validator.MerkleProof.RootHash)
		api.Println(validator.Index, "***[batching_circuit] Validator computed leaf hash:", hFunc.Sum())
		api.Println(validator.Index, "***[batching_out]  Validator MerkleProof Path[0]:", validator.MerkleProof.Path[0])
		api.AssertIsEqual(hFunc.Sum(), validator.MerkleProof.Path[0])
		api.Println(validator.Index, "Assertion passed")

		hFunc.Reset()
		// Thread the intermediate state root (paper Alg. 2 lines 11/17):
		// every validator's membership is verified against the running
		// R_int, not a free per-validator witness root.
		validator.MerkleProof.VerifyProof(api, hFunc, validator.Index, intermediateRoot)
		api.Println(validator.Index, "[batching_circuit] Validator Verified Merkle proof succssed...")

		////////////////////////////////////////////////////////////////////////////////////////////////////
		// *** Signature includes all WithdrawalReqIDs[] + vote bitmask, binding vote to batch ***
		hFunc.Reset()
		hFunc.Write(validator.Index)
		hFunc.Write(c.BatchCommitment)
		hFunc.Write(validator.Vote)
		hFunc.Write(c.RoundID) // bind to roundID:so even if someone tried to reuse their vote in a different round, it would fail verification.
		msg := hFunc.Sum()

		hFunc.Reset()
		if err := eddsa.Verify(curve, validator.Signature, msg, validator.PublicKey, &hFunc); err != nil {
			return fmt.Errorf("signature verification failed for batching circuit: %w", err)
		}
		/////////////////////////////////////////////////////////////////////////////////////////////////////////////
		// 	// All validators must match the same majority vote value
		// NEW PART here we can add if the assertion fails punish the validator
		// api.AssertIsEqual(c.MajorityVote, validator.Vote)

		/////////////////////////THIS PART IS CORRECT/////////////////////////////
		// diff = validator.Vote - MajorityVote
		diff := api.Sub(validator.Vote, c.MajorityVote)

		// isHonest = 1 if diff == 0, else 0
		isHonest := api.IsZero(diff)  // 1 if diff==0 else 0
		api.AssertIsBoolean(isHonest) // safety
		// honest gets +RewardValidator
		rewardedBalance := api.Add(validator.Balance, RewardValidator)

		// dishonest gets max(balance - PenaltyValidator, 0)
		cmpBal := api.Cmp(validator.Balance, PenaltyValidator) // -1 if bal<penalty, 0 if ==, 1 if >
		isLess := api.IsZero(api.Add(cmpBal, 1))               // 1 if cmpBal == -1 else 0
		canPay := api.Sub(1, isLess)                           // 1 if bal>=penalty else 0

		penalized := api.Sub(validator.Balance, PenaltyValidator) // field subtraction OK; clamped by Select
		penalizedOrZero := api.Select(canPay, penalized, 0)

		// newBalance = isHonest ? (bal+reward) : max(bal-penalty, 0)
		newBalance := api.Select(isHonest, rewardedBalance, penalizedOrZero)

		//////////////////////////////////////////////////////
		// Majority counting
		// isHonest = 1 if vote == majority, else 0

		// count how many validators voted for MajorityVote
		majorityCount = api.Add(majorityCount, isHonest)

		//////////////////////////////////////////////////////
		// Reward the validator
		hFunc.Reset()
		hFunc.Write(validator.Index)
		hFunc.Write(validator.PublicKey.A.X)
		hFunc.Write(validator.PublicKey.A.Y)
		hFunc.Write(newBalance)
		validator.MerkleProof.Path[0] = hFunc.Sum()

		// 	//****sort and just need the last validator ComputeRootFromPath
		hFunc.Reset()
		intermediateRoot = validator.MerkleProof.ComputeRootFromPath(api, hFunc, validator.Index)

		bitMask := powbatching(api, 2, validator.Index)
		validatorBits = api.Add(validatorBits, bitMask)

		// set honest bit if isHonest==1
		honestBits = api.Add(honestBits, api.Mul(isHonest, bitMask))
	}
	// The witness contains the selected quorum, not the full committee. Under
	// n=3f+1, the quorum size is 2f+1, so f=(quorum-1)/2 and at least f+1
	// validators in this quorum must agree with MajorityVote.
	quorumSize := len(c.Validators)
	f := (quorumSize - 1) / 2
	threshold := frontend.Variable(f + 1)

	// require majorityCount >= threshold
	cmp := api.Cmp(majorityCount, threshold) // -1 if <, 0 if ==, 1 if >
	api.AssertIsDifferent(cmp, -1)

	api.Println("[batching_out] ValidatorBits:", c.ValidatorBits)
	api.Println(" [batching_circuit] ValidatorBits (bitmask):", validatorBits)

	api.Println("[batching_out] ResultingStateRoot:", c.ResultingStateRoot)
	api.Println(" [batching_circuit] ResultingStateRoot:", intermediateRoot)
	api.Println("[batching_out] HonestBits:", c.HonestBits)
	api.Println(" [batching_circuit] HonestBits (bitmask):", honestBits)
	api.AssertIsEqual(c.ValidatorBits, validatorBits)
	api.AssertIsEqual(c.HonestBits, honestBits)
	api.AssertIsEqual(c.ResultingStateRoot, intermediateRoot)
	api.Println("Successfully verified the batching circuit...")

	return nil
}
