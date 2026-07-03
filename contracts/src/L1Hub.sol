// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice Minimal L1 hub for validator staking + zkSync L1<->L2 messaging.
/// The hub locks stake on L1, sends batched validator imports to L2, and
/// finalizes L2 -> L1 messages (checkpoints, exits, withdrawals, replacements).
interface IZkSyncMailbox {
    struct L2Log {
        uint8 l2ShardId;
        bool isService;
        uint16 txNumberInBlock;
        address sender;
        bytes32 key;
        bytes32 value;
    }

    struct L2TransactionRequest {
        uint256 chainId;
        uint256 mintValue;
        address l2Contract;
        uint256 l2Value;
        bytes l2Calldata;
        uint256 l2GasLimit;
        uint256 l2GasPerPubdataByteLimit;
        bytes[] factoryDeps;
        address refundRecipient;
    }

    function requestL2Transaction(
        address contractL2,
        uint256 l2Value,
        bytes calldata calldataL2,
        uint256 l2GasLimit,
        uint256 l2GasPerPubdataByteLimit,
        bytes[] calldata factoryDeps,
        address refundRecipient
    ) external payable returns (bytes32 canonicalTxHash);

    function requestL2TransactionDirect(L2TransactionRequest calldata request) external payable returns (bytes32 canonicalTxHash);

    function proveL2LogInclusion(
        uint256 blockNumber,
        uint256 logIndex,
        L2Log calldata log,
        bytes32[] calldata proof
    ) external view returns (bool);
}

interface IBridgehubMailbox {
    function proveL2LogInclusion(
        uint256 chainId,
        uint256 batchNumber,
        uint256 logIndex,
        IZkSyncMailbox.L2Log calldata log,
        bytes32[] calldata proof
    ) external view returns (bool);
}

interface IL2Oracle {
    struct PublicKey {
        uint256 x;
        uint256 y;
    }

    struct ValidatorInput {
        uint256 validatorID;
        address validatorAddr;
        uint256 stake;
        PublicKey pubKey;
    }

    struct Account {
        uint256 index;
        PublicKey pubKey;
        uint256 balance;
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

    function importValidatorsFromL1(ValidatorInput[] calldata inputs) external;

    function replaceValidatorFromL1(ReplacementRequest calldata request) external;
}

/// @title L1Hub
/// @notice L1 hub for validator staking + L1->L2 and L2->L1 messaging.
contract L1Hub {
    address public constant L2_MESSENGER_SYSTEM_CONTRACT = address(0x0000000000000000000000000000000000008008);

    address public owner;
    address public l2Oracle;
    IZkSyncMailbox public immutable mailbox;
    uint256 public l2ChainId;
    bool public useBridgehub;

    struct PublicKey {
        uint256 x;
        uint256 y;
    }

    struct ValidatorRecord {
        address validatorAddr;
        uint256 stake;
        PublicKey pubKey;
        bool active;
        bool exitFinalized;
        bool withdrawn;
        bool importRequested;
        // Timestamp of exit finalization; withdrawals are only released
        // after the waiting period (paper SIV-C, F-04).
        uint256 exitFinalizedAt;
    }

    /// @notice Waiting period between a finalized exit and the stake
    /// release (paper SIV-C): prevents a validator from replacing another
    /// and immediately withdrawing.
    uint256 public withdrawDelay;

    mapping(uint256 => ValidatorRecord) public validators;
    uint256[] public pendingValidatorIds;

    struct ReplacementRecord {
        uint256 requestId;
        uint256 targetValidatorID;
        uint256 candidateValidatorID;
        address candidateAddr;
        uint256 candidateStake;
        PublicKey candidatePubKey;
        bool finalized;
        bool success;
    }

    struct ReplacementParams {
        uint256 targetValidatorID;
        uint256 candidateValidatorID;
        PublicKey candidatePubKey;
        uint256 candidateStake;
        uint256 targetLeafIndex;
        uint256[] path;
        uint256 depth;
    }

    mapping(uint256 => ReplacementRecord) public replacements;
    uint256 public nextReplacementRequestId = 1;

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

    mapping(bytes32 => bool) public consumedMessages;
    mapping(uint256 => uint256) public checkpointRootByRound;
    uint256 public latestRound;
    uint256 public latestRoot;

    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);
    event L2OracleUpdated(address indexed l2Oracle);
    event ValidatorRegisteredL1(uint256 indexed validatorID, address indexed validatorAddr, uint256 stake);
    event ValidatorsImportRequested(bytes32 indexed txHash, uint256 count);
    event ValidatorsImportFinalized(uint256 count, uint256 newRoot);
    event ReplacementRequested(uint256 indexed requestId, uint256 indexed targetValidatorID, address indexed candidateAddr, uint256 stake);
    event L2MessageConsumed(bytes32 indexed messageHash, L2ToL1MsgType msgType);
    event CheckpointFinalized(uint256 indexed round, uint256 root);
    event ExitFinalized(uint256 indexed validatorID, address indexed validatorAddr);
    event WithdrawalFinalized(uint256 indexed validatorID, address indexed validatorAddr, uint256 amount);
    event ReplacementFinalized(uint256 indexed requestId, bool success, uint256 newRoot);

    modifier onlyOwner() {
        require(msg.sender == owner, "only owner");
        _;
    }

    constructor(
        address mailboxAddress,
        address l2OracleAddress,
        uint256 l2ChainId_,
        bool useBridgehub_,
        uint256 withdrawDelay_
    ) {
        require(mailboxAddress != address(0), "mailbox required");
        owner = msg.sender;
        mailbox = IZkSyncMailbox(mailboxAddress);
        l2Oracle = l2OracleAddress;
        l2ChainId = l2ChainId_;
        useBridgehub = useBridgehub_;
        withdrawDelay = withdrawDelay_;
    }

    function setWithdrawDelay(uint256 withdrawDelay_) external onlyOwner {
        withdrawDelay = withdrawDelay_;
    }

    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "owner required");
        emit OwnershipTransferred(owner, newOwner);
        owner = newOwner;
    }

    function setL2Oracle(address l2OracleAddress) external onlyOwner {
        require(l2OracleAddress != address(0), "l2 oracle required");
        l2Oracle = l2OracleAddress;
        emit L2OracleUpdated(l2OracleAddress);
    }

    function setUseBridgehub(bool useBridgehub_) external onlyOwner {
        useBridgehub = useBridgehub_;
    }

    function registerValidatorL1(uint256 validatorID, PublicKey calldata pubKey) external payable {
        require(validators[validatorID].validatorAddr == address(0), "validator exists");

        validators[validatorID] = ValidatorRecord({
            validatorAddr: msg.sender,
            stake: msg.value,
            pubKey: pubKey,
            active: false,
            exitFinalized: false,
            withdrawn: false,
            importRequested: false,
            exitFinalizedAt: 0
        });

        pendingValidatorIds.push(validatorID);
        emit ValidatorRegisteredL1(validatorID, msg.sender, msg.value);
    }

    function batchImportValidatorsToL2(
        uint256[] calldata validatorIds,
        uint256 l2GasLimit,
        uint256 l2GasPerPubdataByteLimit,
        address refundRecipient
    ) external payable returns (bytes32 txHash) {
        require(l2Oracle != address(0), "l2 oracle not set");
        require(validatorIds.length > 0, "empty batch");

        IL2Oracle.ValidatorInput[] memory inputs = new IL2Oracle.ValidatorInput[](validatorIds.length);
        for (uint256 i = 0; i < validatorIds.length; i++) {
            uint256 validatorID = validatorIds[i];
            ValidatorRecord storage record = validators[validatorID];
            require(record.validatorAddr != address(0), "validator missing");
            require(!record.importRequested, "validator already imported");

            inputs[i] = IL2Oracle.ValidatorInput({
                validatorID: validatorID,
                validatorAddr: record.validatorAddr,
                stake: record.stake,
                pubKey: IL2Oracle.PublicKey({x: record.pubKey.x, y: record.pubKey.y})
            });
        }

        bytes memory calldataL2 = abi.encodeCall(IL2Oracle.importValidatorsFromL1, (inputs));
        txHash = _requestL2Transaction(calldataL2, l2GasLimit, l2GasPerPubdataByteLimit, refundRecipient, msg.value);
        emit ValidatorsImportRequested(txHash, validatorIds.length);
    }

    function requestReplacementL1(
        ReplacementParams calldata params,
        uint256 l2GasLimit,
        uint256 l2GasPerPubdataByteLimit,
        address refundRecipient
    ) external payable returns (uint256 requestId, bytes32 txHash) {
        require(l2Oracle != address(0), "l2 oracle not set");
        require(params.candidateStake > 0, "stake required");
        require(msg.value >= params.candidateStake, "insufficient value");

        ValidatorRecord storage target = validators[params.targetValidatorID];
        require(target.validatorAddr != address(0), "target missing");
        require(target.active, "target not active");

        requestId = nextReplacementRequestId++;
        replacements[requestId] = ReplacementRecord({
            requestId: requestId,
            targetValidatorID: params.targetValidatorID,
            candidateValidatorID: params.candidateValidatorID,
            candidateAddr: msg.sender,
            candidateStake: params.candidateStake,
            candidatePubKey: params.candidatePubKey,
            finalized: false,
            success: false
        });

        IL2Oracle.ReplacementRequest memory req;
        req.requestId = requestId;
        req.targetValidatorID = params.targetValidatorID;
        req.candidateValidatorID = params.candidateValidatorID;
        req.candidateAddr = msg.sender;
        req.candidateStake = params.candidateStake;
        req.candidatePubKey = IL2Oracle.PublicKey({x: params.candidatePubKey.x, y: params.candidatePubKey.y});
        req.targetAccount = IL2Oracle.Account({
            index: params.targetLeafIndex,
            pubKey: IL2Oracle.PublicKey({x: target.pubKey.x, y: target.pubKey.y}),
            balance: target.stake
        });
        req.path = params.path;
        req.leafIndex = params.targetLeafIndex;
        req.depth = params.depth;

        bytes memory calldataL2 = abi.encodeCall(IL2Oracle.replaceValidatorFromL1, (req));
        uint256 feeValue = msg.value - params.candidateStake;
        txHash = _requestL2Transaction(calldataL2, l2GasLimit, l2GasPerPubdataByteLimit, refundRecipient, feeValue);
        emit ReplacementRequested(requestId, params.targetValidatorID, msg.sender, params.candidateStake);
    }

    function finalizeFromL2(
        bytes calldata message,
        uint256 l2BlockNumber,
        uint256 l2LogIndex,
        uint16 l2TxNumberInBlock,
        bytes32[] calldata proof
    ) external {
        require(l2Oracle != address(0), "l2 oracle not set");

        bytes32 logKey = bytes32(uint256(uint160(l2Oracle)));
        bytes32 logValue = keccak256(message);
        IZkSyncMailbox.L2Log memory l2Log = IZkSyncMailbox.L2Log({
            l2ShardId: 0,
            isService: true,
            txNumberInBlock: l2TxNumberInBlock,
            sender: L2_MESSENGER_SYSTEM_CONTRACT,
            key: logKey,
            value: logValue
        });

        bool ok;
        if (useBridgehub) {
            ok = IBridgehubMailbox(address(mailbox)).proveL2LogInclusion(l2ChainId, l2BlockNumber, l2LogIndex, l2Log, proof);
        } else {
            ok = mailbox.proveL2LogInclusion(l2BlockNumber, l2LogIndex, l2Log, proof);
        }
        require(ok, "invalid log proof");

        bytes32 messageHash = keccak256(message);
        require(!consumedMessages[messageHash], "message already consumed");
        consumedMessages[messageHash] = true;

        (address receiver, bytes memory payload) = abi.decode(message, (address, bytes));
        require(receiver == address(this), "invalid receiver");

        L2ToL1Message memory decoded = abi.decode(payload, (L2ToL1Message));
        emit L2MessageConsumed(messageHash, decoded.msgType);

        if (decoded.msgType == L2ToL1MsgType.CHECKPOINT) {
            (uint256 round, uint256 root) = abi.decode(decoded.payload, (uint256, uint256));
            checkpointRootByRound[round] = root;
            if (round > latestRound) {
                latestRound = round;
            }
            latestRoot = root;
            emit CheckpointFinalized(round, root);
            return;
        }

        if (decoded.msgType == L2ToL1MsgType.EXIT_REQUEST) {
            (uint256 validatorID, address validatorAddr, uint256 rootAtExit) =
                abi.decode(decoded.payload, (uint256, address, uint256));

            require(rootAtExit == latestRoot || checkpointRootByRound[latestRound] == rootAtExit, "unknown root");
            ValidatorRecord storage record = validators[validatorID];
            require(record.validatorAddr == validatorAddr, "validator mismatch");

            record.exitFinalized = true;
            record.exitFinalizedAt = block.timestamp;
            emit ExitFinalized(validatorID, validatorAddr);
            return;
        }

        if (decoded.msgType == L2ToL1MsgType.WITHDRAW_REQUEST) {
            (uint256 validatorID, address validatorAddr, uint256 amount, uint256 rootAtWithdraw) =
                abi.decode(decoded.payload, (uint256, address, uint256, uint256));

            require(rootAtWithdraw == latestRoot || checkpointRootByRound[latestRound] == rootAtWithdraw, "unknown root");
            ValidatorRecord storage record = validators[validatorID];
            require(record.validatorAddr == validatorAddr, "validator mismatch");
            require(record.exitFinalized, "exit not finalized");
            require(!record.withdrawn, "already withdrawn");
            // Waiting period before the stake is released (paper SIV-C,
            // F-04).
            require(
                block.timestamp >= record.exitFinalizedAt + withdrawDelay,
                "withdraw delay not elapsed"
            );
            // The reported amount is the validator's L2 balance (stake plus
            // accrued rewards/penalties); the Hub releases the L1-locked
            // stake.
            require(amount >= record.stake, "amount below locked stake");

            record.withdrawn = true;
            uint256 released = record.stake;
            (bool sent, ) = payable(validatorAddr).call{value: released}("");
            require(sent, "eth transfer failed");
            emit WithdrawalFinalized(validatorID, validatorAddr, released);
            return;
        }

        if (decoded.msgType == L2ToL1MsgType.REPLACEMENT_RESULT) {
            (
                uint256 requestId,
                bool success,
                uint256 targetValidatorID,
                uint256 candidateValidatorID,
                address candidateAddr,
                uint256 candidateStake,
                uint256 newRoot
            ) = abi.decode(decoded.payload, (uint256, bool, uint256, uint256, address, uint256, uint256));

            ReplacementRecord storage req = replacements[requestId];
            require(req.requestId == requestId, "replacement missing");
            require(!req.finalized, "replacement finalized");
            require(req.targetValidatorID == targetValidatorID, "target mismatch");
            require(req.candidateValidatorID == candidateValidatorID, "candidate mismatch");
            require(req.candidateAddr == candidateAddr, "candidate addr mismatch");
            require(req.candidateStake == candidateStake, "candidate stake mismatch");

            req.finalized = true;
            req.success = success;

            if (success) {
                ValidatorRecord storage target = validators[targetValidatorID];
                address oldValidator = target.validatorAddr;
                uint256 oldStake = target.stake;

                target.validatorAddr = candidateAddr;
                target.stake = candidateStake;
                target.pubKey = req.candidatePubKey;
                target.active = true;
                target.exitFinalized = false;
                target.withdrawn = false;

                if (oldStake > 0) {
                    (bool sent, ) = payable(oldValidator).call{value: oldStake}("");
                    require(sent, "old stake transfer failed");
                }
            } else {
                if (req.candidateStake > 0) {
                    (bool refunded, ) = payable(req.candidateAddr).call{value: req.candidateStake}("");
                    require(refunded, "candidate refund failed");
                }
            }

            if (newRoot != 0) {
                latestRoot = newRoot;
            }
            emit ReplacementFinalized(requestId, success, newRoot);
            return;
        }

        if (decoded.msgType == L2ToL1MsgType.VALIDATOR_IMPORT_RESULT) {
            (uint256[] memory validatorIds, uint256 newRoot) = abi.decode(decoded.payload, (uint256[], uint256));
            for (uint256 i = 0; i < validatorIds.length; i++) {
                uint256 validatorID = validatorIds[i];
                ValidatorRecord storage record = validators[validatorID];
                require(record.validatorAddr != address(0), "validator missing");
                if (record.importRequested) {
                    continue;
                }
                record.importRequested = true;
                record.active = true;
            }
            emit ValidatorsImportFinalized(validatorIds.length, newRoot);
            return;
        }

        revert("unknown message type");
    }

    function _requestL2Transaction(
        bytes memory calldataL2,
        uint256 l2GasLimit,
        uint256 l2GasPerPubdataByteLimit,
        address refundRecipient,
        uint256 feeValue
    ) internal returns (bytes32 canonicalTxHash) {
        if (useBridgehub) {
            IZkSyncMailbox.L2TransactionRequest memory request = IZkSyncMailbox.L2TransactionRequest({
                chainId: l2ChainId,
                mintValue: feeValue,
                l2Contract: l2Oracle,
                l2Value: 0,
                l2Calldata: calldataL2,
                l2GasLimit: l2GasLimit,
                l2GasPerPubdataByteLimit: l2GasPerPubdataByteLimit,
                factoryDeps: new bytes[](0),
                refundRecipient: refundRecipient
            });
            canonicalTxHash = mailbox.requestL2TransactionDirect{value: feeValue}(request);
        } else {
            canonicalTxHash = mailbox.requestL2Transaction{value: feeValue}(
                l2Oracle,
                0,
                calldataL2,
                l2GasLimit,
                l2GasPerPubdataByteLimit,
                new bytes[](0),
                refundRecipient
            );
        }
    }
}
