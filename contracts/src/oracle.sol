// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "./mimc.sol";
import {Verifier as VotingVerifier} from "./VotingBatchVerifier.sol";
import "./merkle_tree.sol";

/// @notice zkSync system messenger used for L2 -> L1 messages (paper SIV-B
/// Step 5; validator_functions design note).
interface IL2ToL1SystemMessenger {
    function sendToL1(bytes calldata message) external returns (bytes32);
}

/// @notice L1 -> L2 address aliasing (zkSync convention): an L1 contract
/// calling an L2 contract arrives as alias(l1Address).
library L1AliasHelper {
    uint160 internal constant OFFSET =
        uint160(0x1111000000000000000000000000000000001111);

    function applyL1ToL2Alias(address l1Address)
        internal
        pure
        returns (address)
    {
        unchecked {
            return address(uint160(l1Address) + OFFSET);
        }
    }
}

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
        address recipient;
        bytes proof;
        bytes publicWitness;
        bytes32 nullifierHash;
        bool isApproved;
        bool isClaimed;
    }

    // L1-anchored validator lifecycle (paper SIV-B/SIV-C;
    // validator_functions design note). The ABI of these structs must match
    // IL2Oracle in L1Hub.sol.

    struct ValidatorInput {
        uint256 validatorID;
        address validatorAddr;
        uint256 stake;
        PublicKey pubKey;
    }

    struct ReplacementRequest {
        uint256 requestId;
        uint256 targetValidatorID;
        uint256 candidateValidatorID;
        address candidateAddr;
        uint256 candidateStake;
        PublicKey candidatePubKey;
        Account targetAccount;
        uint256[] path;
        uint256 leafIndex;
        uint256 depth;
    }

    enum L2ToL1MsgType {
        CHECKPOINT,
        EXIT_REQUEST,
        WITHDRAW_REQUEST,
        REPLACEMENT_RESULT,
        VALIDATOR_IMPORT_RESULT
    }

    struct L2ToL1Message {
        L2ToL1MsgType msgType;
        bytes payload;
    }

    VotingVerifier votingVerifier;

    uint256 seedX;
    uint256 seedY;

    // Commitment-tree anchoring (paper SIV-E Steps 6-7): after threshold
    // agreement the aggregator publishes the appended leaves to the DFS and
    // reports the object identifier together with the finalized commitment
    // root, so any mismatch is publicly detectable from the recorded values.
    string latestIPFSHash;
    uint256 private latestCommitmentRoot;
    uint256 private commitmentEpoch;
    mapping(uint256 => uint256) public commitmentRootByEpoch;
    mapping(uint256 => bool) private publishedCommitmentRoots;

    mapping(address => bool) private users;
    mapping(address => uint) private tokenBurnBalances;
    mapping(address => uint) private tokenClaimBalances;

    mapping(uint256 => ClaimRequest) private claimRequests;

    // Defensive double-spend guard (paper Theorem 2): the committee keeps the
    // authoritative spent-list in the DFS, but the Gateway also refuses to
    // accept or mint an already-spent nullifier, so replay is impossible even
    // against a Byzantine committee.
    mapping(bytes32 => bool) public spentNullifiers;

    mapping(uint => address) private validators;
    uint256[] private validatorsList;

    uint256 private aggregator;
    uint256 private aggregatorIdx = type(uint256).max; // sentinel: forces first selection to index 0

    // Aggregator rotation (paper SIV-C): handover is allowed only after the
    // on-chain timeout elapses without a finalized round, or immediately
    // after a successful round finalization (normal round-robin rotation).
    uint256 public immutable aggregatorTimeout;
    uint256 private roundStartedAt;
    bool private rotationReady;

    mapping(address => bool) private isValidator;
    mapping(uint256 => bool) private roundFinalized;

    // L1 Hub anchoring (paper SIV-B/SIV-C). When l1Hub is unset the L2->L1
    // messaging paths are no-ops, which keeps single-chain runs (and the
    // L1-anchored baseline) working without a Hub.
    address public owner;
    address public l1Hub;
    address private constant L2_TO_L1_SYSTEM_MESSENGER =
        address(0x0000000000000000000000000000000000008008);

    // Validators that requested an exit (by account leaf index); withdrawal
    // is only possible after the exit request (paper SIV-C).
    mapping(uint256 => bool) private exitRequested;

    // Deterministic batching (paper SIV-D): claims receive sequential
    // identifiers so that round r covers exactly [r*b, (r+1)*b - 1].
    uint256 public immutable batchSize;
    uint256 public immutable destinationID;
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
        bytes32 nullifierHash,
        uint256 destinationID,
        address recipient
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

    event CommitmentRootPublished(
        uint256 indexed epoch,
        uint256 root,
        string dfsRef
    );

    /// @notice Emitted when a validator leaf is replaced. Carries the full
    /// new leaf contents (index, key, stake) so every off-chain validator can
    /// mirror the update in its local state tree.
    event Replaced(
        address indexed sender,
        address indexed replaced,
        uint256 index,
        PublicKey pubkey,
        uint256 stake
    );
    event Exiting(address indexed sender);
    event Withdrawn(address indexed sender);

    event L1HubUpdated(address indexed l1Hub);
    event ValidatorsImportedFromL1(uint256 count, uint256 newRoot);
    event ReplacementFromL1Processed(
        uint256 indexed requestId,
        bool success,
        uint256 newRoot
    );
    event L2ToL1MessageSent(L2ToL1MsgType msgType, bytes payload);

    // endregion

    constructor(
        uint256 _levels,
        uint256 _seedX,
        uint256 _seedY,
        address votingVerifierAddress,
        uint256 _batchSize,
        uint256 _aggregatorTimeout,
        uint256 _destinationID
    ) MerkleTree(_levels) {
        require(_batchSize > 0 && _batchSize <= 256, "invalid batch size");
        levels = _levels;
        seedX = _seedX;
        seedY = _seedY;
        votingVerifier = VotingVerifier(votingVerifierAddress);
        batchSize = _batchSize;
        destinationID = _destinationID;
        aggregatorTimeout = _aggregatorTimeout;
        roundStartedAt = block.timestamp;
        owner = msg.sender;
    }

    // region l1 anchoring

    function setL1Hub(address l1HubAddress) external {
        require(msg.sender == owner, "only owner");
        require(l1HubAddress != address(0), "l1 hub required");
        l1Hub = l1HubAddress;
        emit L1HubUpdated(l1HubAddress);
    }

    modifier onlyL1Hub() {
        require(l1Hub != address(0), "l1 hub not set");
        require(
            msg.sender == L1AliasHelper.applyL1ToL2Alias(l1Hub),
            "unauthorized L1 sender"
        );
        _;
    }

    /// @notice Import a batch of L1-registered validators into the L2
    /// validator-state tree (paper SIV-B/SIV-C: registration starts on the
    /// L1 Hub, which batches N registrations into one L1->L2 transaction).
    function importValidatorsFromL1(
        ValidatorInput[] calldata inputs
    ) external onlyL1Hub {
        for (uint256 i = 0; i < inputs.length; i++) {
            ValidatorInput calldata input = inputs[i];
            require(
                validators[input.validatorID] == address(0),
                "validator already registered"
            );
            validators[input.validatorID] = input.validatorAddr;
            validatorsList.push(input.validatorID);
            isValidator[input.validatorAddr] = true;

            Account memory account = Account(
                getNextLeafIndex(),
                input.pubKey,
                input.stake
            );
            accounts[account.index] = input.validatorAddr;

            uint[] memory leaf = new uint[](1);
            leaf[0] = hashAccount(account);
            insert(MiMC.hash(leaf));

            emit ValidatorRegistered(
                input.validatorAddr,
                input.validatorID,
                account.index,
                account.pubKey,
                account.balance
            );
        }

        uint256 newRoot = getRoot();
        emit ValidatorsImportedFromL1(inputs.length, newRoot);
    }

    /// @notice Apply an L1-initiated validator replacement (paper SIV-C:
    /// the candidate's higher stake is locked on the L1 Hub before the L2
    /// replacement is applied) and report the result back to L1.
    function replaceValidatorFromL1(
        ReplacementRequest calldata request
    ) external onlyL1Hub {
        bool success = request.candidateStake >
            request.targetAccount.balance &&
            request.path.length > 0 &&
            request.path[0] == hashAccount(request.targetAccount) &&
            verify(request.path, request.leafIndex, request.depth);

        if (success) {
            Account memory replaced = Account(
                request.targetAccount.index,
                request.candidatePubKey,
                request.candidateStake
            );
            update(
                hashAccount(replaced),
                request.path,
                request.leafIndex,
                request.depth
            );

            address replacedAddr = accounts[request.targetAccount.index];
            accounts[request.targetAccount.index] = request.candidateAddr;
            validators[request.targetValidatorID] = request.candidateAddr;
            isValidator[request.candidateAddr] = true;
            emit Replaced(
                request.candidateAddr,
                replacedAddr,
                request.targetAccount.index,
                request.candidatePubKey,
                request.candidateStake
            );
        }

        uint256 newRoot = getRoot();
        emit ReplacementFromL1Processed(request.requestId, success, newRoot);
        _sendToL1(
            L2ToL1MsgType.REPLACEMENT_RESULT,
            abi.encode(
                request.requestId,
                success,
                request.targetValidatorID,
                request.candidateValidatorID,
                request.candidateAddr,
                request.candidateStake,
                newRoot
            )
        );
    }

    /// @notice Send a typed message to the L1 Hub via the zkSync system
    /// messenger. The oracle itself calls the system contract so the L2->L1
    /// log key equals the oracle address, matching the Hub's proof check.
    /// No-op when no Hub is configured (single-chain runs, baseline).
    function _sendToL1(L2ToL1MsgType msgType, bytes memory payload) internal {
        if (l1Hub == address(0)) {
            return;
        }
        bytes memory message = abi.encode(
            l1Hub,
            abi.encode(L2ToL1Message({msgType: msgType, payload: payload}))
        );
        IL2ToL1SystemMessenger(L2_TO_L1_SYSTEM_MESSENGER).sendToL1(message);
        emit L2ToL1MessageSent(msgType, payload);
    }

    // endregion

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
        isValidator[msg.sender] = true;

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
        bytes32 nullifierHash,
        uint256 claimDestinationID,
        address recipient
    ) external {
        require(users[msg.sender] == true, "address not registered");
        require(claimDestinationID == destinationID, "wrong destination");
        require(recipient != address(0), "recipient required");
        require(!spentNullifiers[nullifierHash], "nullifier already spent");

        uint256 uniqueID = nextClaimID;
        nextClaimID++;

        claimRequests[uniqueID] = ClaimRequest(
            uniqueID,
            msg.sender,
            recipient,
            proof,
            publicWitness,
            nullifierHash,
            false,
            false
        );

        emit ClaimSubmitted(
            uniqueID,
            proof,
            publicWitness,
            nullifierHash,
            claimDestinationID,
            recipient
        );
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
        // Only the designated round aggregator may finalize (paper SIV-D).
        require(index == aggregator, "not current aggregator");
        // Dedicated sentinel: an all-zero majority bitmask must also
        // finalize the round exactly once.
        require(!roundFinalized[roundId], "round already finalized");
        roundFinalized[roundId] = true;

        wiVotes[roundId] = vote;

        uint256 preStateRoot = getRoot();
        uint[12] memory input = [
            preStateRoot,
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
            if (
                request.from == address(0) ||
                request.isClaimed ||
                spentNullifiers[request.nullifierHash]
            ) {
                continue;
            }
            request.isApproved = true;
            request.isClaimed = true;
            spentNullifiers[request.nullifierHash] = true;
            tokenClaimBalances[request.recipient] += BURN_AMOUNT;
            emit ClaimMinted(
                baseID + i,
                request.recipient,
                request.nullifierHash
            );
        }

        // Round finalized: allow the normal round-robin handover and restart
        // the timeout clock (paper SIV-C).
        rotationReady = true;
        roundStartedAt = block.timestamp;

        // Checkpoint the finalized validator-state root to L1 (paper SIV-B
        // Step 5).
        _sendToL1(
            L2ToL1MsgType.CHECKPOINT,
            abi.encode(roundId, postStateRoot)
        );

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
        isValidator[msg.sender] = true;
        // Interaction last
        (bool ok, ) = replacedAddr.call{value: toReplace.balance}("");
        require(ok, "ETH_TRANSFER_FAILED");

        emit Replaced(
            msg.sender,
            replacedAddr,
            toReplace.index,
            publicKey,
            msg.value
        );
    }

    /// @notice Request an exit from the validator set (paper SIV-C): the
    /// exit is recorded as a validator-state transition on L2 and forwarded
    /// to the L1 Hub, which starts the withdrawal waiting period. The stake
    /// itself is locked (and later released) on L1.
    function exit(
        Account memory account,
        uint256[] memory path,
        uint256 leafIndex,
        uint256 depth
    ) public {
        require(accounts[account.index] == msg.sender, "wrong sender address");
        require(!exitRequested[account.index], "exit already requested");

        require(
            path[0] == hashAccount(account),
            "leaf does not match account"
        );
        require(verify(path, leafIndex, depth), "invalid merkle proof");

        exitRequested[account.index] = true;

        // validatorID == account leaf index by the registration convention
        // used throughout this codebase.
        _sendToL1(
            L2ToL1MsgType.EXIT_REQUEST,
            abi.encode(account.index, msg.sender, getRoot())
        );

        emit Exiting(msg.sender);
    }

    /// @notice Withdraw a validator's balance (paper SIV-C): only possible
    /// after a recorded exit request. The leaf is zeroed on L2 and the
    /// withdrawal is forwarded to the L1 Hub, which releases the stake after
    /// the waiting period.
    function withdraw(
        Account memory account,
        uint256[] memory path,
        uint256 leafIndex,
        uint256 depth
    ) public {
        require(accounts[account.index] == msg.sender, "wrong sender address");
        require(exitRequested[account.index], "exit not requested");

        require(
            path[0] == hashAccount(account),
            "leaf does not match account"
        );
        require(verify(path, leafIndex, depth), "invalid merkle proof");

        // Root under which this withdrawal was authorized (matches the last
        // checkpoint; the leaf update below intentionally happens after).
        uint256 rootAtWithdraw = getRoot();

        delete accounts[account.index];
        delete exitRequested[account.index];

        Account memory empty = Account(account.index, account.pubKey, 0);
        update(hashAccount(empty), path, leafIndex, depth);

        _sendToL1(
            L2ToL1MsgType.WITHDRAW_REQUEST,
            abi.encode(account.index, msg.sender, account.balance, rootAtWithdraw)
        );

        emit Withdrawn(msg.sender);
    }

    // endregion

    // region aggregator

    function getAggregator() public view returns (uint) {
        return aggregator;
    }

    /// @notice Advance to the next round-robin aggregator (paper SIV-C).
    /// Validators may trigger the handover only after a finalized round
    /// (normal rotation) or once the on-chain timeout has elapsed without
    /// finalization.
    function chooseNewAggregator() external {
        require(isValidator[msg.sender], "not a registered validator");
        require(
            rotationReady ||
                block.timestamp >= roundStartedAt + aggregatorTimeout,
            "aggregator timeout not reached"
        );
        require(validatorsList.length > 0, "no validators registered");

        rotationReady = false;
        roundStartedAt = block.timestamp;

        if (aggregatorIdx >= validatorsList.length - 1) {
            aggregatorIdx = 0;
        } else {
            aggregatorIdx++;
        }

        aggregator = validatorsList[aggregatorIdx];
        emit NewAggregator(aggregator);
    }

    // endregion

    // region commitment root / dfs

    /// @notice Record the finalized commitment-tree root together with the
    /// DFS object identifier (paper SIV-E Steps 6-7). Only the current round
    /// aggregator may publish; claims are only accepted against recorded
    /// roots.
    function publishCommitmentRoot(
        uint256 root,
        string memory dfsRef
    ) external {
        require(
            validators[aggregator] == msg.sender,
            "not current aggregator"
        );
        commitmentEpoch++;
        commitmentRootByEpoch[commitmentEpoch] = root;
        publishedCommitmentRoots[root] = true;
        latestCommitmentRoot = root;
        latestIPFSHash = dfsRef;
        emit CommitmentRootPublished(commitmentEpoch, root, dfsRef);
    }

    function viewLatestIPFSHash() public view returns (string memory) {
        return latestIPFSHash;
    }

    function viewLatestCommitmentRoot()
        public
        view
        returns (uint256, uint256)
    {
        return (commitmentEpoch, latestCommitmentRoot);
    }

    function isPublishedCommitmentRoot(
        uint256 root
    ) public view returns (bool) {
        return publishedCommitmentRoots[root];
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
