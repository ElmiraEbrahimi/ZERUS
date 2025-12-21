// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Script.sol";
import "../src/MerkleProofVerifier.sol";

contract DeployMerkleProofVerifier is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("ZKSYNC_PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        // Adjust the type name if ExportSolidity generated a different contract name.
        Verifier verifier = new Verifier();

        console2.log("MerkleProofVerifier deployed at", address(verifier));

        vm.stopBroadcast();
    }
}
