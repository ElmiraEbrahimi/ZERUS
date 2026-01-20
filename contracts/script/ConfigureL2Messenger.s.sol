// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Script.sol";
import "../src/L2Messenger.sol";

/// @notice Updates the L2Messenger with the L1 messenger address.
contract ConfigureL2Messenger is Script {
    function run() external {
        uint256 deployerKey = vm.envUint("ZKSYNC_PRIVATE_KEY");
        address l2Messenger = vm.envAddress("L2_MESSENGER_CONTRACT_ADDRESS");
        address l1Messenger = vm.envAddress("L1_MESSENGER_CONTRACT_ADDRESS");

        vm.startBroadcast(deployerKey);
        L2Messenger(l2Messenger).setL1Messenger(l1Messenger);
        vm.stopBroadcast();

        console2.log("L2Messenger updated with L1Messenger", l1Messenger);
    }
}
