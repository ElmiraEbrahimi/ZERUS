# Paper ↔ Code Terminology Mapping

The paper (ZERUS) uses protocol-level names; the implementation predates
them in places. This table maps every paper term to its code entry point so
artifact reviewers can cross-reference (F-29).

| Paper term | Code term / entry point |
|---|---|
| Gateway contract (L2) | `contracts/src/oracle.sol` (`Oracle`, extends `MerkleTree`) |
| Hub contract (L1) | `contracts/src/L1Hub.sol` (`L1Hub`; deployed via `make deploy-l1hub`) |
| Burn (source rollup) | `Oracle.burn(commitmentHash)`; client `user.BurnTx()` binds `a_dst` to the burn sender |
| Claim (mint on destination) | user "withdraw": `POST /users/default/withdraw` → `user.WithdrawTx()` → `Oracle.claim(...)` |
| Claim vote (per request) | `WiVote` ("withdrawal vote") |
| Batch vote (b-bit bitmask) | `SignedBatchVote` (validator-signed, paper §IV-D Step 14) |
| Aggregating-proof submission (Step 17) | `Oracle.submitWiVote(...)` |
| Commitment tree (append-only) | incremental tree (`merkle.IncrementalMerkleTree`, `IncVote` flow) |
| Commitment-root anchoring (Steps 6–7) | `Oracle.publishCommitmentRoot(root, dfsRef)` |
| Nullifier spent-list | DFS `IPFSContent.SpentNullifiers` + on-chain `Oracle.spentNullifiers` |
| Validator-state tree | sparse tree (`gnark/state.go`) + on-chain `MerkleTree` in `oracle.sol` |
| Validator-state checkpoint (Step 5) | `L2ToL1MsgType.CHECKPOINT` → `L1Hub.finalizeFromL2` (Go relayer `OracleL2ToL1Relayer`) |
| Redeeming circuit (Alg. 1) | `go-server/circuits/merkle_proof/redeeming_circuit.go` |
| Aggregating circuit (Alg. 2) | `go-server/circuits/voting_batch/voting_circuit_batching.go` + `merkle.go` |
| Finalization threshold `f+1` | `bftThreshold(n)` in `node.go`; same rule in the Aggregating circuit |
| Round aggregator rotation (§IV-C) | `Oracle.chooseNewAggregator()` (timeout-gated) |
| Validator registration on L1 (§IV-C) | `L1Hub.registerValidatorL1` → `batchImportValidatorsToL2` → `Oracle.importValidatorsFromL1` |
| Validator replacement (§IV-C) | `L1Hub.requestReplacementL1` → `Oracle.replaceValidatorFromL1`; direct L2 path `Oracle.replace` |
| Exit / stake withdrawal (§IV-C) | `Oracle.exit` / `Oracle.withdraw` → L2→L1 messages → `L1Hub.finalizeFromL2` (timelocked release) |
| Redeeming public inputs | `d_dst, R_comm, h_n, a_dst` |
| Destination-rollup id `d_dst` | `DESTINATION_ID` config, used in `C = H(n_rd ∥ s_rd ∥ d_dst ∥ a_dst)` |
| DFS | IPFS via HTTP API (`IPFS_API_URL`); local simulation fallback (`internal/oracle-repo/db/ipfs.go`) |
| Source rollup L2A / destination L2B (§V) | `SOURCE_RPC_URL` + `SOURCE_ORACLE_CONTRACT_ADDRESS` profile (source) vs. local Gateway (destination) |
| Go module | `l2alchemy` (ZERUS); the L1-anchored baseline uses `l1alchemy` |

Historical note: the repository predates the paper's naming ("L2-Alchemy",
`Oracle`); the contract and module names are kept to avoid invalidating
deployed artifacts and generated bindings, with this file as the canonical
cross-reference.
