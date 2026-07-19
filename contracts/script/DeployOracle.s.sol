// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/oracle.sol";

contract DeployOracle is Script {
    function run() external returns (address oracleAddr) {
        uint256 deployerPrivateKey = vm.envUint("ZKSYNC_PRIVATE_KEY");

        uint256 levels = vm.envUint("SPARSE_TREE_DEPTH");
        uint256 seedX = vm.envUint("ORACLE_SEED_X");
        uint256 seedY = vm.envUint("ORACLE_SEED_Y");
        uint256 batchSize = vm.envUint("BATCH_SIZE");
        uint256 aggregatorTimeout = vm.envOr(
            "AGGREGATOR_TIMEOUT",
            uint256(300)
        );

        address votingVerifierAddress = vm.envAddress(
            "VOTING_BATCH_VERIFIER_ADDRESS"
        );

        vm.startBroadcast(deployerPrivateKey);

        Oracle oracle = new Oracle(
            levels,
            seedX,
            seedY,
            votingVerifierAddress,
            batchSize,
            aggregatorTimeout,
            vm.envUint("DESTINATION_ID")
        );

        oracleAddr = address(oracle);
        console2.log("Oracle deployed at", oracleAddr);

        vm.stopBroadcast();
    }
}
