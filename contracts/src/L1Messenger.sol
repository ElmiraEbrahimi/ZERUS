// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice Minimal L1 <-> L2 messaging demo contract for zkSync.
/// It can request an L2 transaction (L1 -> L2) and verify an L2 -> L1 message.
interface IZkSyncMailbox {
    struct L2Log {
        uint8 l2ShardId;
        bool isService;
        uint16 txNumberInBlock;
        address sender;
        bytes32 key;
        bytes32 value;
    }

    struct L2Message {
        uint16 txNumberInBlock;
        address sender;
        bytes data;
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

    function proveL2MessageInclusion(
        uint256 blockNumber,
        uint256 messageIndex,
        L2Message calldata message,
        bytes32[] calldata proof
    ) external view returns (bool);

    function proveL2LogInclusion(
        uint256 blockNumber,
        uint256 logIndex,
        L2Log calldata log,
        bytes32[] calldata proof
    ) external view returns (bool);
}

interface IBridgehubMailbox {
    function proveL2MessageInclusion(
        uint256 chainId,
        uint256 batchNumber,
        uint256 messageIndex,
        IZkSyncMailbox.L2Message calldata message,
        bytes32[] calldata proof
    ) external view returns (bool);

    function proveL2LogInclusion(
        uint256 chainId,
        uint256 batchNumber,
        uint256 logIndex,
        IZkSyncMailbox.L2Log calldata log,
        bytes32[] calldata proof
    ) external view returns (bool);
}

interface IL2MessengerReceiver {
    function receiveFromL1(string calldata message) external;
}

/// @title L1Messenger
/// @notice Sends messages to L2 and verifies L2 -> L1 messages.
contract L1Messenger {
    address public constant L2_MESSENGER_SYSTEM_CONTRACT = address(0x0000000000000000000000000000000000008008);
    address public owner;
    address public l2Messenger;
    IZkSyncMailbox public immutable mailbox;
    uint256 public l2ChainId;
    bool public useBridgehub;

    string public lastMessageToL2;
    string public lastMessageFromL2;

    mapping(bytes32 => bool) public consumedMessages;

    event L1ToL2MessageRequested(bytes32 indexed txHash, string message);
    event L2ToL1MessageReceived(bytes32 indexed messageHash, string message);
    event L2MessengerUpdated(address indexed newL2Messenger);
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    modifier onlyOwner() {
        require(msg.sender == owner, "only owner");
        _;
    }

    constructor(address mailboxAddress, address l2MessengerAddress, uint256 l2ChainId_, bool useBridgehub_) {
        require(mailboxAddress != address(0), "mailbox required");
        if (useBridgehub_) {
            require(l2ChainId_ != 0, "chain id required");
        }
        owner = msg.sender;
        mailbox = IZkSyncMailbox(mailboxAddress);
        l2Messenger = l2MessengerAddress;
        l2ChainId = l2ChainId_;
        useBridgehub = useBridgehub_;
    }

    function setL2Messenger(address l2MessengerAddress) external onlyOwner {
        l2Messenger = l2MessengerAddress;
        emit L2MessengerUpdated(l2MessengerAddress);
    }

    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "owner required");
        emit OwnershipTransferred(owner, newOwner);
        owner = newOwner;
    }

    /// @notice Send a message to L2 using requestL2Transaction (single-chain mode).
    function sendToL2(
        string calldata message,
        uint256 l2GasLimit,
        uint256 l2GasPerPubdataByteLimit
    ) external payable returns (bytes32 txHash) {
        require(l2Messenger != address(0), "l2 messenger not set");
        bytes memory data = abi.encodeWithSelector(IL2MessengerReceiver.receiveFromL1.selector, message);

        txHash = mailbox.requestL2Transaction{value: msg.value}(
            l2Messenger,
            0,
            data,
            l2GasLimit,
            l2GasPerPubdataByteLimit,
            new bytes[](0),
            msg.sender
        );

        lastMessageToL2 = message;
        emit L1ToL2MessageRequested(txHash, message);
    }

    /// @notice Send a message to L2 using requestL2TransactionDirect (zk-chains mode).
    function sendToL2Direct(
        uint256 l2ChainIdParam,
        string calldata message,
        uint256 l2GasLimit,
        uint256 l2GasPerPubdataByteLimit
    ) external payable returns (bytes32 txHash) {
        require(l2Messenger != address(0), "l2 messenger not set");
        bytes memory data = abi.encodeWithSelector(IL2MessengerReceiver.receiveFromL1.selector, message);

        // Bridgehub requires mintValue to match msg.value.
        IZkSyncMailbox.L2TransactionRequest memory request = IZkSyncMailbox.L2TransactionRequest({
            chainId: l2ChainIdParam,
            mintValue: msg.value,
            l2Contract: l2Messenger,
            l2Value: 0,
            l2Calldata: data,
            l2GasLimit: l2GasLimit,
            l2GasPerPubdataByteLimit: l2GasPerPubdataByteLimit,
            factoryDeps: new bytes[](0),
            refundRecipient: msg.sender
        });

        txHash = mailbox.requestL2TransactionDirect{value: msg.value}(request);

        lastMessageToL2 = message;
        emit L1ToL2MessageRequested(txHash, message);
    }

    /// @notice Verify and store a message sent from L2.
    function receiveFromL2(
        string calldata message,
        uint256 l2BlockNumber,
        uint256 l2MessageIndex,
        uint16 l2TxNumberInBlock,
        bytes32[] calldata proof
    ) external {
        require(l2Messenger != address(0), "l2 messenger not set");
        bytes32 logKey = bytes32(uint256(uint160(l2Messenger)));
        bytes32 logValue = keccak256(abi.encode(address(this), message));
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
            ok = IBridgehubMailbox(address(mailbox)).proveL2LogInclusion(
                l2ChainId,
                l2BlockNumber,
                l2MessageIndex,
                l2Log,
                proof
            );
        } else {
            ok = mailbox.proveL2LogInclusion(l2BlockNumber, l2MessageIndex, l2Log, proof);
        }
        require(ok, "invalid log proof");

        bytes32 messageHash = keccak256(
            abi.encode(l2BlockNumber, l2MessageIndex, l2TxNumberInBlock, l2Log.sender, l2Log.key, l2Log.value)
        );
        require(!consumedMessages[messageHash], "message already consumed");
        consumedMessages[messageHash] = true;

        lastMessageFromL2 = message;
        emit L2ToL1MessageReceived(messageHash, message);
    }
}
