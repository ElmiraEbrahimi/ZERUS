// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/merkle_tree.sol";

contract DeployMerkleTree is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("ZKSYNC_PRIVATE_KEY");
        uint256 levels = vm.envUint("ORACLE_LEVELS"); // reuse same env var

        vm.startBroadcast(deployerPrivateKey);

        MerkleTree tree = new MerkleTree(levels);

        console2.log("MerkleTree deployed at", address(tree));

        vm.stopBroadcast();
    }
}
