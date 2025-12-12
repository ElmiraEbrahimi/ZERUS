// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/VotingBatchVerifier.sol";

contract DeployVotingBatchVerifier is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        Verifier verifier = new Verifier();

        console2.log("VotingBatchVerifier deployed at", address(verifier));

        vm.stopBroadcast();
    }
}
