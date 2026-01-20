// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Script.sol";
import "../src/L1Messenger.sol";

/// @notice Deploys L1Messenger to L1.
contract DeployL1Messenger is Script {
    function run() external {
        uint256 deployerKey = vm.envUint("ZKSYNC_PRIVATE_KEY");
        address mailbox = vm.envAddress("L1_MAILBOX_ADDRESS");
        address l2Messenger = vm.envOr("L2_MESSENGER_CONTRACT_ADDRESS", address(0));
        bool useBridgehub = vm.envBool("L1_USE_DIRECT_MESSAGING");
        uint256 l2ChainId = vm.envUint("ZKSYNC_CHAIN_ID");

        vm.startBroadcast(deployerKey);
        L1Messenger messenger = new L1Messenger(mailbox, l2Messenger, l2ChainId, useBridgehub);
        vm.stopBroadcast();

        console2.log("L1Messenger deployed at", address(messenger));
        if (l2Messenger != address(0)) {
            console2.log("L1Messenger linked to L2Messenger", l2Messenger);
        }
    }
}
