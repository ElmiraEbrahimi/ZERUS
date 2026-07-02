// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;
import "./mimc.sol";

contract MerkleTree {
    uint256 public constant ZERO_VALUE =
        4555114089170143013007615382799372902997177870479602349537353593038812875418;

    uint256 public levels;

    uint256 private root;
    uint256 private nextIndex;
    mapping(uint256 => uint256) private filledSubtrees;

    // zeros() provides precomputed subtree hashes up to index 7, which bounds
    // the supported depth to 8 levels (256 leaves).
    uint256 public constant MAX_LEVELS = 8;

    constructor(uint256 _levels) {
        require(_levels > 0, "_levels must be greater than zero");
        require(_levels <= MAX_LEVELS, "_levels must be at most 8");

        levels = _levels;

        for (uint32 i = 0; i < levels; i++) {
            filledSubtrees[i] = zeros(i);
        }

        root = zeros(levels - 1);
    }

    function hashLeftRight(
        uint256 left,
        uint256 right
    ) public pure returns (uint256) {
        uint[] memory input = new uint[](2);
        input[0] = left;
        input[1] = right;
        return MiMC.hash(input);
    }

    function insert(uint256 leaf) internal {
        require(nextIndex != 2 ** levels, "tree is full");
        uint256 left;
        uint256 right;
        uint256 currentIndex = nextIndex;
        uint256 currentHash = leaf;

        for (uint i = 0; i < levels; i++) {
            if (currentIndex % 2 == 0) {
                left = currentHash;
                right = zeros(i);
                filledSubtrees[i] = currentHash;
            } else {
                left = filledSubtrees[i];
                right = currentHash;
            }
            currentHash = hashLeftRight(left, right);
            currentIndex /= 2;
        }

        root = currentHash;

        nextIndex += 1;
    }

    function update(
        uint256 leaf,
        uint256[] memory path,
        uint256 leafIndex,
        uint256 depth
    ) internal {
        // Call verify with leafIndex and depth
        verify(path, leafIndex, depth);
        path[0] = leaf;
        // Call computeRootFromPath with leafIndex and depth
        root = computeRootFromPath(path, leafIndex, depth);
    }

    function verify(
        uint256[] memory path,
        uint256 leafIndex,
        uint256 depth
    ) public view returns (bool) {
        uint256 computedHash = path[0]; // Start with the leaf itself

        // Iterate through the Merkle proof path using binary representation of leafIndex
        for (uint256 i = 1; i <= depth; i++) {
            uint256 bit = (leafIndex >> (i - 1)) & 1; // Extract the i-th bit (LSB first)

            if (bit == 1) {
                computedHash = hashLeftRight(path[i], computedHash); // Right child
            } else {
                computedHash = hashLeftRight(computedHash, path[i]); // Left child
            }
        }

        return computedHash == root;
    }

    function computeRootFromPath(
        uint256[] memory path,
        uint256 leafIndex,
        uint256 depth
    ) public pure returns (uint256) {
        uint256 computedHash = path[0]; // Start with the leaf itself

        // Iterate through the Merkle proof path using binary representation of leafIndex
        for (uint256 i = 1; i <= depth; i++) {
            uint256 bit = (leafIndex >> (i - 1)) & 1; // Extract the i-th bit (LSB first)

            if (bit == 1) {
                computedHash = hashLeftRight(path[i], computedHash); // Right child
            } else {
                computedHash = hashLeftRight(computedHash, path[i]); // Left child
            }
        }

        return computedHash;
    }

    function getRoot() public view returns (uint256) {
        return root;
    }

    function setRoot(uint256 _root) internal {
        root = _root;
    }

    function getLevels() public view returns (uint256) {
        return levels;
    }

    function getNextLeafIndex() public view returns (uint256) {
        return nextIndex;
    }

    function zeros(uint256 i) public pure returns (uint256) {
        if (i == 0)
            return
                4555114089170143013007615382799372902997177870479602349537353593038812875418;
        else if (i == 1)
            return
                9614759978327623946452646332910600180945773348102064399025967221305784063943;
        else if (i == 2)
            return
                15762506290347708512348905356059207185046946323941989403490412292643733744343;
        else if (i == 3)
            return
                2078761282949659850987695139228042067769933906673781403014209677812047702550;
        else if (i == 4)
            return
                20395412135670005113982952294980644860334762516174975965396550838918688642133;
        else if (i == 5)
            return
                17560454953585356688949527489694249319830065182192048221544285096802159445652;
        else if (i == 6)
            return
                20019762671335178393512154978075455201849332419823879662510519485824706883752;
        else if (i == 7)
            return
                12065157948427223398688325372361960271507319753018415581972466307863230644170;
        else revert("index out of bounds");
    }
}
