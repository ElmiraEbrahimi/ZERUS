// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Script.sol";
import "../src/Counter.sol";

/// @title DeployCounter
/// @notice A Foundry script that deploys the Counter contract to the configured chain (L1 or L2).
/// It reads the deployer's private key from the `ZKSYNC_PRIVATE_KEY` environment
/// variable and broadcasts the transaction to the network specified by the Foundry
/// configuration. After deployment the script logs the address of the new contract.
contract DeployCounter is Script {
    function run() external {
        // The deployer's private key is expected to be provided via environment
        // variables. See foundry.toml for the network configuration.
        uint256 deployerKey = vm.envUint("ZKSYNC_PRIVATE_KEY");
        vm.startBroadcast(deployerKey);
        Counter counter = new Counter();
        vm.stopBroadcast();

        // Log the deployed address so it can be captured by caller tooling.
        console2.log("Counter deployed at", address(counter));
    }
}
