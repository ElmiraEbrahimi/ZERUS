// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Script.sol";
import "../src/L1Hub.sol";

/// @notice Deploys the L1 Hub (validator staking + L1<->L2 anchoring) to L1
/// (paper SIV-B/SIV-C).
contract DeployL1Hub is Script {
    function run() external {
        uint256 deployerKey = vm.envUint("ZKSYNC_PRIVATE_KEY");
        address mailbox = vm.envAddress("L1_MAILBOX_ADDRESS");
        address l2Oracle = vm.envOr("ORACLE_CONTRACT_ADDRESS", address(0));
        bool useBridgehub = vm.envBool("L1_USE_DIRECT_MESSAGING");
        uint256 l2ChainId = vm.envUint("ZKSYNC_CHAIN_ID");
        uint256 withdrawDelay = vm.envOr(
            "HUB_WITHDRAW_DELAY",
            uint256(86400)
        );

        vm.startBroadcast(deployerKey);
        L1Hub hub = new L1Hub(
            mailbox,
            l2Oracle,
            l2ChainId,
            useBridgehub,
            withdrawDelay
        );
        vm.stopBroadcast();

        console2.log("L1Hub deployed at", address(hub));
        if (l2Oracle != address(0)) {
            console2.log("L1Hub linked to Oracle", l2Oracle);
        }
    }
}
