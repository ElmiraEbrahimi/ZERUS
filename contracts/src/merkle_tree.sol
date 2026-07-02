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

    /// @dev Hash of a single leaf value, matching the LeafSum convention of
    /// the gnark circuits and the gnark-crypto Merkle tree.
    function leafSum(uint256 leaf) public pure returns (uint256) {
        uint[] memory input = new uint[](1);
        input[0] = leaf;
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

    /// @dev Proof paths follow the gnark-crypto / circuit convention: path[0]
    /// is the raw leaf value (e.g. the account hash) and the tree leaf is
    /// MiMC(path[0]). Sibling ordering is decided by the bits of leafIndex.
    function verify(
        uint256[] memory path,
        uint256 leafIndex,
        uint256 depth
    ) public view returns (bool) {
        uint256 computedHash = leafSum(path[0]);

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
        uint256 computedHash = leafSum(path[0]);

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

    /// @dev Precomputed empty-subtree hashes: zeros(0) = ZERO_VALUE and
    /// zeros(i+1) = MiMC(zeros(i) || zeros(i)) under the 110-round MiMC that
    /// matches gnark-crypto's MIMC_BN254.
    function zeros(uint256 i) public pure returns (uint256) {
        if (i == 0)
            return
                4555114089170143013007615382799372902997177870479602349537353593038812875418;
        else if (i == 1)
            return
                19836274635794509466014568838483484624600013041370895881939528370985142278954;
        else if (i == 2)
            return
                13617893609837248081008283892047260108796513789343133982990357780506196299573;
        else if (i == 3)
            return
                17019866949954776883112551581972189814819563148075409968100822283266078590629;
        else if (i == 4)
            return
                2138726945387356468363054067854480562066796241025489545311417198275461384100;
        else if (i == 5)
            return
                9530699369832533880639149369063437013450242560813587771525635668574190709561;
        else if (i == 6)
            return
                3103391909757162857233288457092921903629638236252641536170476897663767868208;
        else if (i == 7)
            return
                5180474051210496608687766129670424206209071725010745029464643986723420860574;
        else revert("index out of bounds");
    }
}
