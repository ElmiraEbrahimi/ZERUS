// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Script.sol";
import "../src/L2Messenger.sol";

/// @notice Deploys L2Messenger to zkSync L2.
contract DeployL2Messenger is Script {
    function run() external {
        uint256 deployerKey = vm.envUint("ZKSYNC_PRIVATE_KEY");
        address l1Messenger = vm.envOr("L1_MESSENGER_CONTRACT_ADDRESS", address(0));

        vm.startBroadcast(deployerKey);
        L2Messenger messenger = new L2Messenger(l1Messenger);
        vm.stopBroadcast();

        console2.log("L2Messenger deployed at", address(messenger));
        if (l1Messenger != address(0)) {
            console2.log("L2Messenger linked to L1Messenger", l1Messenger);
        }
    }
}
