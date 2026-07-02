// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "./mimc.sol";
import {Verifier as VotingVerifier} from "./VotingBatchVerifier.sol";
import "./merkle_tree.sol";

contract Oracle is MerkleTree {
    uint256 constant INCENTIVIZE_AMOUNT = 1;
    uint256 constant BURN_AMOUNT = 100;
    uint256 constant INIT_TOKEN_BURN_AMOUNT = 1_000_000;
    uint256 constant INIT_TOKEN_CLAIM_AMOUNT = 0;
    uint256 public constant AGGREGATOR_REWARD = 500000000000000;
    uint256 public constant VALIDATOR_REWARD = 20000000000;

    struct Account {
        uint256 index;
        PublicKey pubKey;
        uint256 balance;
    }

    struct PublicKey {
        uint256 x;
        uint256 y;
    }

    struct ClaimRequest {
        uint256 uniqueID;
        address from;
        bytes proof;
        bytes publicWitness;
        bytes32 nullifierHash;
        bool isApproved;
        bool isClaimed;
    }

    VotingVerifier votingVerifier;

    uint256 seedX;
    uint256 seedY;

    string latestIPFSHash;

    mapping(address => bool) private users;
    mapping(address => uint) private tokenBurnBalances;
    mapping(address => uint) private tokenClaimBalances;

    mapping(uint256 => ClaimRequest) private claimRequests;

    mapping(uint => address) private validators;
    uint256[] private validatorsList;

    uint256 private aggregator;
    uint256 private aggregatorIdx = type(uint256).max; // sentinel: forces first selection to index 0

    // Deterministic batching (paper SIV-D): claims receive sequential
    // identifiers so that round r covers exactly [r*b, (r+1)*b - 1].
    uint256 public immutable batchSize;
    uint256 private nextClaimID;

    mapping(uint256 => address) accounts;
    mapping(uint256 => uint256) wiVotes;

    // region events

    event NewAggregator(uint256 validatorID);

    event UserRegistered(
        address addr,
        uint256 index,
        PublicKey pubkey,
        uint256 balance
    );

    event ValidatorRegistered(
        address addr,
        uint256 validatorID,
        uint256 index,
        PublicKey pubkey,
        uint256 balance
    );

    event Registered(
        address sender,
        uint256 index,
        PublicKey pubkey,
        uint256 value
    );

    event BurnSubmitted(bytes32 commitmentHash);

    event ClaimSubmitted(
        uint256 uniqueID,
        bytes proof,
        bytes publicWitness,
        bytes32 nullifierHash
    );

    event WiVoteSubmitted(
        uint256 submitter,
        uint256 validators,
        uint256 honestBits,
        uint256 request,
        uint256 majorityVote
    );

    event ClaimMinted(
        uint256 uniqueID,
        address recipient,
        bytes32 nullifierHash
    );

    event Replaced(address indexed sender, address indexed replaced);
    event Exiting(address indexed sender);
    event Withdrawn(address indexed sender);

    // endregion

    constructor(
        uint256 _levels,
        uint256 _seedX,
        uint256 _seedY,
        address votingVerifierAddress,
        uint256 _batchSize
    ) MerkleTree(_levels) {
        require(_batchSize > 0 && _batchSize <= 256, "invalid batch size");
        levels = _levels;
        seedX = _seedX;
        seedY = _seedY;
        votingVerifier = VotingVerifier(votingVerifierAddress);
        batchSize = _batchSize;
    }

    // region 1.register

    function registerValidator(
        uint validatorID,
        PublicKey memory publicKey
    ) external payable {
        require(
            validators[validatorID] == address(0),
            "validator already registered"
        );
        validators[validatorID] = msg.sender;
        validatorsList.push(validatorID);

        Account memory account = Account(
            getNextLeafIndex(),
            publicKey,
            msg.value
        );
        accounts[account.index] = msg.sender;
        uint256 accountHash = hashAccount(account);

        uint[] memory input = new uint[](1);
        input[0] = accountHash;
        uint256 h = MiMC.hash(input);

        insert(h);

        emit ValidatorRegistered(
            msg.sender,
            validatorID,
            account.index,
            account.pubKey,
            account.balance
        );

        emit Registered(
            msg.sender,
            account.index,
            account.pubKey,
            account.balance
        );
    }

    function registerUser(PublicKey memory publicKey) external payable {
        require(!users[msg.sender], "user already registered");
        users[msg.sender] = true;

        tokenBurnBalances[msg.sender] = INIT_TOKEN_BURN_AMOUNT;
        tokenClaimBalances[msg.sender] = INIT_TOKEN_CLAIM_AMOUNT;

        emit UserRegistered(
            msg.sender,
            0, // Index not used for users
            publicKey,
            msg.value
        );
    }

    function viewBalance() external view returns (uint, uint) {
        require(users[msg.sender] == true, "address not registered");
        return (tokenBurnBalances[msg.sender], tokenClaimBalances[msg.sender]);
    }

    // endregion

    // region 2.burn

    function burn(bytes32 commitmentHash) external {
        require(users[msg.sender] == true, "address not registered");
        if (BURN_AMOUNT > tokenBurnBalances[msg.sender]) {
            revert("insufficient balance");
        }
        tokenBurnBalances[msg.sender] -= BURN_AMOUNT;
        emit BurnSubmitted(commitmentHash);
    }

    // endregion

    // region 3.claim

    /// @notice Record a claim request (paper SIV-E Steps 10-11): the Gateway
    /// assigns a batch-policy identifier and emits a claim event. Tokens are
    /// minted only after the round's Aggregating proof verifies (Step 17).
    function claim(
        bytes memory proof,
        bytes memory publicWitness,
        bytes32 nullifierHash
    ) external {
        require(users[msg.sender] == true, "address not registered");

        uint256 uniqueID = nextClaimID;
        nextClaimID++;

        claimRequests[uniqueID] = ClaimRequest(
            uniqueID,
            msg.sender,
            proof,
            publicWitness,
            nullifierHash,
            false,
            false
        );

        emit ClaimSubmitted(uniqueID, proof, publicWitness, nullifierHash);
    }

    // endregion

    // region 4.wivote

    /// @notice Finalize a verification round (paper SIV-E Steps 16-17): the
    /// aggregator submits the batch outcome with the Aggregating proof; the
    /// Gateway verifies it, mints tokens for accepted claims, and applies the
    /// validator-state update.
    function submitWiVote(
        uint256 index,
        uint256 roundId,
        uint256 batchCommitment,
        uint256 validatorBits,
        uint256 honestBits,
        uint256 vote,
        uint256 postStateRoot,
        uint256 postSeedX,
        uint256 postSeedY,
        uint256[8] memory proof
    ) public {
        require(accounts[index] == msg.sender, "invalid index");
        require(wiVotes[roundId] == 0, "already submitted");

        wiVotes[roundId] = vote;

        uint[11] memory input = [
            postStateRoot,
            roundId,
            batchCommitment,
            vote, // threshold-supported majority vote bitmask
            validatorBits,
            honestBits,
            index,
            seedX,
            seedY,
            postSeedX,
            postSeedY
        ];

        votingVerifier.verifyProof(proof, input);

        seedX = postSeedX;
        seedY = postSeedY;

        setRoot(postStateRoot);

        // Mint accepted claims: bit i of the threshold-supported bitmask
        // corresponds to claim identifier roundId*batchSize + i (SIV-D).
        uint256 baseID = roundId * batchSize;
        for (uint256 i = 0; i < batchSize; i++) {
            if ((vote >> i) & 1 == 0) {
                continue;
            }
            ClaimRequest storage request = claimRequests[baseID + i];
            if (request.from == address(0) || request.isClaimed) {
                continue;
            }
            request.isApproved = true;
            request.isClaimed = true;
            tokenClaimBalances[request.from] += BURN_AMOUNT;
            emit ClaimMinted(baseID + i, request.from, request.nullifierHash);
        }

        emit WiVoteSubmitted(index, validatorBits, honestBits, roundId, vote);
    }

    // endregion

    // region 5.exit

    function replace(
        PublicKey memory publicKey,
        Account memory toReplace,
        uint256[] memory path,
        uint256 leafIndex,
        uint256 depth
    ) public payable {
        // Paper SIV-C: a replacement candidate must lock a strictly higher
        // stake than the validator being replaced.
        require(
            msg.value > toReplace.balance,
            "replacement stake must exceed incumbent stake"
        );

        require(
            path[0] == hashAccount(toReplace),
            "leaf does not match account"
        );
        require(verify(path, leafIndex, depth), "invalid merkle proof");

        Account memory replaced = Account(
            toReplace.index,
            publicKey,
            msg.value
        );

        update(hashAccount(replaced), path, leafIndex, depth);

        address payable replacedAddr = payable(accounts[toReplace.index]);
        // Effects first: prevent the old owner from re-entering as the current owner
        accounts[toReplace.index] = msg.sender;
        // Interaction last
        (bool ok, ) = replacedAddr.call{value: toReplace.balance}("");
        require(ok, "ETH_TRANSFER_FAILED");

        emit Replaced(msg.sender, replacedAddr);
    }

    function exit(
        Account memory account,
        uint256[] memory path,
        uint256 leafIndex,
        uint256 depth
    ) public {
        require(accounts[account.index] == msg.sender, "wrong sender address");

        require(
            path[0] == hashAccount(account),
            "leaf does not match account"
        );
        require(verify(path, leafIndex, depth), "invalid merkle proof");

        emit Exiting(msg.sender);
    }

    function withdraw(
        Account memory account,
        uint256[] memory path,
        uint256 leafIndex,
        uint256 depth
    ) public {
        require(accounts[account.index] == msg.sender, "wrong sender address");

        require(
            path[0] == hashAccount(account),
            "leaf does not match account"
        );
        require(verify(path, leafIndex, depth), "invalid merkle proof");

        // payable(msg.sender).transfer(account.balance);
        delete accounts[account.index];

        Account memory empty = Account(account.index, account.pubKey, 0);
        update(hashAccount(empty), path, leafIndex, depth);
        emit Withdrawn(msg.sender);
    }

    // endregion

    // region aggregator

    function getAggregator() public view returns (uint) {
        return aggregator;
    }

    function chooseNewAggregator() external {
        require(validatorsList.length > 0, "no validators registered");

        if (aggregatorIdx >= validatorsList.length - 1) {
            aggregatorIdx = 0;
        } else {
            aggregatorIdx++;
        }

        aggregator = validatorsList[aggregatorIdx];
        emit NewAggregator(aggregator);
    }

    // endregion

    // region ipfs

    function updateLatestIPFSHash(string memory ipfsHash) external {
        latestIPFSHash = ipfsHash;
    }

    function viewLatestIPFSHash() public view returns (string memory) {
        return latestIPFSHash;
    }

    // endregion

    // region utils

    function hashAccount(Account memory account) public pure returns (uint256) {
        uint[] memory input = new uint[](4);
        input[0] = account.index;
        input[1] = account.pubKey.x;
        input[2] = account.pubKey.y;
        input[3] = account.balance;

        return MiMC.hash(input);
    }

    function getSeed() public view returns (uint256, uint256) {
        return (seedX, seedY);
    }

    function getReward() public view returns (uint256) {
        return
            AGGREGATOR_REWARD + (getNextLeafIndex() / 2 + 1) * VALIDATOR_REWARD;
    }

    // endregion
}
