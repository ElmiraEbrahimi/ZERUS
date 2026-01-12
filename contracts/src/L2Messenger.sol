// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IL1MessengerSystem {
    function sendToL1(bytes calldata message) external returns (bytes32);
}

library AddressAliasHelper {
    uint160 internal constant OFFSET = uint160(0x1111000000000000000000000000000000001111);

    function applyL1ToL2Alias(address l1Address) internal pure returns (address) {
        unchecked {
            return address(uint160(l1Address) + OFFSET);
        }
    }

    function undoL1ToL2Alias(address l2Address) internal pure returns (address) {
        unchecked {
            return address(uint160(l2Address) - OFFSET);
        }
    }
}

/// @title L2Messenger
/// @notice Sends messages to L1 and receives messages from L1.
contract L2Messenger {
    address public owner;
    address public l1Messenger;

    string public lastMessageToL1;
    string public lastMessageFromL1;

    address private constant L1_MESSENGER_SYSTEM_CONTRACT = address(0x0000000000000000000000000000000000008008);

    event L2ToL1MessageSent(bytes32 indexed messageHash, string message);
    event L1ToL2MessageReceived(string message);
    event L1MessengerUpdated(address indexed newL1Messenger);
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    modifier onlyOwner() {
        require(msg.sender == owner, "only owner");
        _;
    }

    constructor(address l1MessengerAddress) {
        owner = msg.sender;
        l1Messenger = l1MessengerAddress;
    }

    function setL1Messenger(address l1MessengerAddress) external onlyOwner {
        l1Messenger = l1MessengerAddress;
        emit L1MessengerUpdated(l1MessengerAddress);
    }

    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "owner required");
        emit OwnershipTransferred(owner, newOwner);
        owner = newOwner;
    }

    /// @notice Send a message to L1 via the zkSync system messenger.
    function sendToL1(string calldata message) external returns (bytes32 messageHash) {
        require(l1Messenger != address(0), "l1 messenger not set");
        bytes memory payload = abi.encode(l1Messenger, message);
        messageHash = IL1MessengerSystem(L1_MESSENGER_SYSTEM_CONTRACT).sendToL1(payload);
        lastMessageToL1 = message;
        emit L2ToL1MessageSent(messageHash, message);
    }

    /// @notice Receive a message from L1 (called by zkSync when L1 -> L2 tx executes).
    function receiveFromL1(string calldata message) external {
        require(l1Messenger != address(0), "l1 messenger not set");
        require(
            msg.sender == AddressAliasHelper.applyL1ToL2Alias(l1Messenger),
            "unauthorized L1 sender"
        );

        lastMessageFromL1 = message;
        emit L1ToL2MessageReceived(message);
    }
}
