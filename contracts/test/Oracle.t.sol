// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Test.sol";
import "../src/oracle.sol";

contract MockVotingVerifier {
    function verifyProof(uint256[8] memory, uint256[12] memory) external pure {}
}

/// @notice Gateway conformance tests (F-30): cover the paper-mandated
/// rules that are checkable without a Groth16 proof — claim/mint gating
/// (F-05), sequential batch identifiers (F-07), aggregator gating (F-13),
/// replacement stake rule (F-03), exit-before-withdraw (F-16), timed
/// aggregator rotation (F-23), and authenticated commitment-root
/// publication (F-22). Proof-dependent paths (mint on threshold bitmask,
/// nullifier marking) are exercised by the Go integration flow.
contract OracleTest is Test {
    Oracle internal oracle;
    MockVotingVerifier internal mockVerifier;

    address internal validator0 = address(0xA0);
    address internal validator1 = address(0xA1);
    address internal user = address(0xB0);
    address internal stranger = address(0xC0);
    address internal recipient = address(0xD0);

    uint256 internal constant BATCH_SIZE = 4;
    uint256 internal constant AGG_TIMEOUT = 300;
    uint256 internal constant DESTINATION_ID = 7;

    function setUp() public {
        mockVerifier = new MockVotingVerifier();
        oracle = new Oracle(
            3, // levels
            1, // seedX
            2, // seedY
            address(mockVerifier),
            BATCH_SIZE,
            AGG_TIMEOUT,
            DESTINATION_ID
        );

        vm.deal(validator0, 100 ether);
        vm.deal(validator1, 100 ether);

        vm.prank(validator0);
        oracle.registerValidator{value: 1000}(
            0,
            Oracle.PublicKey({x: 11, y: 12})
        );
        vm.prank(validator1);
        oracle.registerValidator{value: 1000}(
            1,
            Oracle.PublicKey({x: 21, y: 22})
        );

        vm.prank(user);
        oracle.registerUser(Oracle.PublicKey({x: 31, y: 32}));
        vm.prank(recipient);
        oracle.registerUser(Oracle.PublicKey({x: 41, y: 42}));
    }

    // region registration

    function testRegisterValidatorDuplicateReverts() public {
        vm.prank(stranger);
        vm.expectRevert(bytes("validator already registered"));
        oracle.registerValidator(0, Oracle.PublicKey({x: 1, y: 2}));
    }

    function testRegisterUserDuplicateReverts() public {
        vm.prank(user);
        vm.expectRevert(bytes("user already registered"));
        oracle.registerUser(Oracle.PublicKey({x: 1, y: 2}));
    }

    function testInitialAggregatorIsValidatorZero() public view {
        assertEq(oracle.getAggregator(), 0);
    }

    // endregion

    // region burn + claim (F-05, F-07)

    function testBurnDecrementsBurnBalance() public {
        vm.prank(user);
        oracle.burn(bytes32(uint256(1)));

        vm.prank(user);
        (uint256 burnBalance, uint256 claimBalance) = oracle.viewBalance();
        assertEq(burnBalance, 1_000_000 - 100);
        assertEq(claimBalance, 0);
    }

    function testBurnUnregisteredReverts() public {
        vm.prank(stranger);
        vm.expectRevert(bytes("address not registered"));
        oracle.burn(bytes32(uint256(1)));
    }

    /// F-07: the Gateway assigns sequential identifiers so round r covers
    /// exactly [r*b, (r+1)*b - 1].
    function testClaimAssignsSequentialIdentifiers() public {
        for (uint256 i = 0; i < 2; i++) {
            vm.prank(user);
            vm.expectEmit(false, false, false, true, address(oracle));
            emit Oracle.ClaimSubmitted(
                i,
                hex"aa",
                hex"bb",
                bytes32(uint256(0x100 + i)),
                DESTINATION_ID,
                user
            );
            oracle.claim(hex"aa", hex"bb", bytes32(uint256(0x100 + i)), DESTINATION_ID, user);
        }
    }

    /// F-05: claim() records the request as pending; nothing is minted
    /// before the round's Aggregating proof verifies.
    function testClaimDoesNotMint() public {
        vm.prank(user);
        oracle.claim(hex"aa", hex"bb", bytes32(uint256(0x100)), DESTINATION_ID, user);

        vm.prank(user);
        (, uint256 claimBalance) = oracle.viewBalance();
        assertEq(claimBalance, 0);
    }

    function testClaimUnregisteredReverts() public {
        vm.prank(stranger);
        vm.expectRevert(bytes("address not registered"));
        oracle.claim(hex"aa", hex"bb", bytes32(uint256(0x100)), DESTINATION_ID, user);
    }

    function testClaimWrongDestinationReverts() public {
        vm.prank(user);
        vm.expectRevert(bytes("wrong destination"));
        oracle.claim(hex"aa", hex"bb", bytes32(uint256(0x100)), DESTINATION_ID + 1, user);
    }

    function testClaimRequiresRecipient() public {
        vm.prank(user);
        vm.expectRevert(bytes("recipient required"));
        oracle.claim(hex"aa", hex"bb", bytes32(uint256(0x100)), DESTINATION_ID, address(0));
    }

    function testAcceptedClaimMintsToProofBoundRecipientAndBlocksReplay() public {
        bytes32 nullifierHash = bytes32(uint256(0x100));
        vm.prank(user);
        oracle.claim(hex"aa", hex"bb", nullifierHash, DESTINATION_ID, recipient);

        uint256[8] memory proof;
        uint256 currentRoot = oracle.getRoot();
        vm.prank(validator0);
        oracle.submitWiVote(0, 0, 0, 0, 0, 1, currentRoot, 3, 4, proof);

        vm.prank(user);
        (, uint256 submitterBalance) = oracle.viewBalance();
        vm.prank(recipient);
        (, uint256 recipientBalance) = oracle.viewBalance();
        assertEq(submitterBalance, 0);
        assertEq(recipientBalance, 100);
        assertTrue(oracle.spentNullifiers(nullifierHash));

        vm.prank(user);
        vm.expectRevert(bytes("nullifier already spent"));
        oracle.claim(hex"aa", hex"bb", nullifierHash, DESTINATION_ID, user);
    }

    function testSubmitWiVoteBindsProofToCurrentPreStateRoot() public {
        uint256[8] memory proof;
        uint256 preStateRoot = oracle.getRoot();
        uint256 postStateRoot = preStateRoot + 123;

        uint256[12] memory expectedInput = [
            preStateRoot,
            postStateRoot,
            uint256(0),
            uint256(0),
            uint256(1),
            uint256(0),
            uint256(0),
            uint256(0),
            uint256(1),
            uint256(2),
            uint256(3),
            uint256(4)
        ];

        vm.expectCall(
            address(mockVerifier),
            abi.encodeCall(MockVotingVerifier.verifyProof, (proof, expectedInput))
        );

        vm.prank(validator0);
        oracle.submitWiVote(0, 0, 0, 0, 0, 1, postStateRoot, 3, 4, proof);
    }

    // endregion

    // region submitWiVote gating (F-13)

    function testSubmitWiVoteWrongOwnerReverts() public {
        uint256[8] memory proof;
        vm.prank(stranger);
        vm.expectRevert(bytes("invalid index"));
        oracle.submitWiVote(0, 0, 0, 0, 0, 0, 0, 0, 0, proof);
    }

    /// F-13: only the designated round aggregator may finalize.
    function testSubmitWiVoteNonAggregatorReverts() public {
        uint256[8] memory proof;
        vm.prank(validator1);
        vm.expectRevert(bytes("not current aggregator"));
        oracle.submitWiVote(1, 0, 0, 0, 0, 0, 0, 0, 0, proof);
    }

    // endregion

    // region replacement stake rule (F-03)

    function testReplaceWithoutHigherStakeReverts() public {
        Oracle.Account memory target = Oracle.Account({
            index: 0,
            pubKey: Oracle.PublicKey({x: 11, y: 12}),
            balance: 1000
        });
        uint256[] memory path = new uint256[](4);

        vm.prank(stranger);
        vm.deal(stranger, 1 ether);
        vm.expectRevert(
            bytes("replacement stake must exceed incumbent stake")
        );
        oracle.replace{value: 1000}(
            Oracle.PublicKey({x: 41, y: 42}),
            target,
            path,
            0,
            3
        );
    }

    // endregion

    // region exit / withdraw (F-16)

    function testWithdrawWithoutExitReverts() public {
        Oracle.Account memory account = Oracle.Account({
            index: 0,
            pubKey: Oracle.PublicKey({x: 11, y: 12}),
            balance: 1000
        });
        uint256[] memory path = new uint256[](4);

        vm.prank(validator0);
        vm.expectRevert(bytes("exit not requested"));
        oracle.withdraw(account, path, 0, 3);
    }

    function testExitWrongSenderReverts() public {
        Oracle.Account memory account = Oracle.Account({
            index: 0,
            pubKey: Oracle.PublicKey({x: 11, y: 12}),
            balance: 1000
        });
        uint256[] memory path = new uint256[](4);

        vm.prank(stranger);
        vm.expectRevert(bytes("wrong sender address"));
        oracle.exit(account, path, 0, 3);
    }

    // endregion

    // region aggregator rotation (F-23)

    function testRotationByNonValidatorReverts() public {
        vm.warp(block.timestamp + AGG_TIMEOUT + 1);
        vm.prank(stranger);
        vm.expectRevert(bytes("not a registered validator"));
        oracle.chooseNewAggregator();
    }

    function testRotationBeforeTimeoutReverts() public {
        vm.warp(block.timestamp + AGG_TIMEOUT - 2);
        vm.prank(validator0);
        vm.expectRevert(bytes("aggregator timeout not reached"));
        oracle.chooseNewAggregator();
    }

    function testRotationAfterTimeoutSucceeds() public {
        vm.warp(block.timestamp + AGG_TIMEOUT + 1);
        vm.prank(validator0);
        oracle.chooseNewAggregator();
        assertEq(oracle.getAggregator(), 0); // first rotation selects index 0
    }

    // endregion

    // region commitment-root publication (F-22)

    function testPublishCommitmentRootByNonAggregatorReverts() public {
        vm.prank(stranger);
        vm.expectRevert(bytes("not current aggregator"));
        oracle.publishCommitmentRoot(123, "Qm...");
    }

    function testPublishCommitmentRootRecordsRoot() public {
        assertFalse(oracle.isPublishedCommitmentRoot(123));

        vm.prank(validator0); // validators[aggregator=0]
        oracle.publishCommitmentRoot(123, "QmRef");

        assertTrue(oracle.isPublishedCommitmentRoot(123));
        (uint256 epoch, uint256 root) = oracle.viewLatestCommitmentRoot();
        assertEq(epoch, 1);
        assertEq(root, 123);
        assertEq(oracle.viewLatestIPFSHash(), "QmRef");
    }

    // endregion

    // region constructor guards

    function testConstructorRejectsZeroBatchSize() public {
        vm.expectRevert(bytes("invalid batch size"));
        new Oracle(3, 1, 2, address(0xBEEF), 0, AGG_TIMEOUT, DESTINATION_ID);
    }

    // endregion
}
