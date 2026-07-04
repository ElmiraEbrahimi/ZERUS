# ZERUS (L2 implementation)

Local zkSync stack + Foundry contracts + a Go HTTP API that runs the ZERUS
protocol: the L2 Gateway contract (named `Oracle` in code), the L1 Hub, and
the off-chain validator committee that batches and proves cross-rollup
claims. Paper-to-code naming is documented in `MAPPING.md`.

## Prerequisites

Pinned toolchain (paper §V):

- Docker + Docker Compose
- Foundry `forge` 1.3.5 with the foundry-zksync fork (v0.1.5); Solidity is
  pinned to 0.8.30 in `contracts/foundry.toml`
- Go 1.24.x (gnark v0.14.0, gnark-crypto v0.19.0, go-ethereum v1.16.7 are
  pinned in `go-server/go.mod`)
- abigen (only needed for deploy targets that regenerate bindings)
- Python 3 (used by `scripts/deployed-address.sh` during `make deploy` and
  by the evaluation aggregation script)

## Repository Layout

```
contracts/    # Foundry contracts + deploy scripts
go-server/    # Go HTTP API + oracle runtime
local-setup/  # zkSync local stack scripts + compose files
scripts/      # deploy-address resolution + evaluation sweep scripts
Makefile      # handy wrapper targets
```

## Quickstart (single L2)

1. Create a local env file and update values as needed:

```sh
cp .env.example .env
```

The defaults in `.env.example` already target the local stack
(`ZKSYNC_RPC_URL=http://localhost:3051`, `ZKSYNC_CHAIN_ID=270`). You can grab
a funded dev key from `local-setup/rich-wallets.json`.

2. Start the local zkSync stack:

```sh
make up
```

3. Deploy contracts and regenerate Go bindings:

```sh
make deploy
```

`make deploy` requires `abigen` and `python3` in your PATH. It clears stale
local protocol state, deploys the verifiers, the MiMC library (pinned into
Oracle/MerkleTree via `--libraries`), the Gateway (`Oracle`), the L1/L2
messengers, and the L1 Hub, and rewrites the address fields in `.env` from
the forge broadcast journals.

4. Start the Go API server:

```sh
make server
```

The API binds to `HTTP_BIND_ADDR` (default `:18000` from `.env.example`).

## Local zkSync Ports (single L2)

These are the defaults from `local-setup/docker-compose.yml`:

| Component   | URL                   |
| ----------- | --------------------- |
| L2 JSON-RPC | http://localhost:3050 |
| L2 WS       | ws://localhost:3051   |

The local node also answers HTTP JSON-RPC on port 3051; `.env.example`
targets `http://localhost:3051` so the Go server and forge share one
endpoint. Either port works for `ZKSYNC_RPC_URL`.

## ZK Chains (multi L2, optional)

If you want the multi-chain stack with explorers, run:

```sh
cd local-setup
./start-zk-chains.sh
```

Defaults from `local-setup/zk-chains-docker-compose.yml`:

| Component          | URL                    |
| ------------------ | ---------------------- |
| L2 JSON-RPC (L2-0) | http://localhost:15100 |
| L2 WS (L2-0)       | ws://localhost:15101   |
| Explorer           | http://localhost:15005 |
| Hyperexplorer      | http://localhost:15000 |

Set `ZKSYNC_RPC_URL` and `ZKSYNC_CHAIN_ID` accordingly if you use this mode
(the Foundry config targets chain ID 271 for zk-chains).

### Cross-rollup deployment (source + destination)

The paper (§V) evaluates ZERUS on two independent zkSync Era instances:
burns execute on the source rollup (L2A) and claims on the destination
rollup (L2B). To reproduce that topology:

1. Start the multi-chain stack (`./start-zk-chains.sh`) so two independent
   L2 instances are available.
2. Deploy the contracts once per instance (`make deploy` against each
   instance's RPC), giving each Gateway its own `DESTINATION_ID`.
3. On the destination server, point the source profile at the source
   instance's Gateway:

```sh
SOURCE_RPC_URL=http://localhost:15100     # source rollup RPC
SOURCE_ORACLE_CONTRACT_ADDRESS=0x...      # source Gateway
SOURCE_CHAIN_ID=271                       # optional; auto-detected if 0
BURN_CONFIRMATION_DEPTH=1                 # source finality rule (blocks)
```

With the source profile set, the committee observes `BurnSubmitted` events
on the source chain, buffers them until the confirmation depth elapses,
inserts them into the commitment tree in deterministic
(source, block, log index, commitment) order (§IV-E), and verifies claims
against the local (destination) Gateway. Without the profile the runtime
falls back to the single-rollup loop for local testing.

## Environment Variables

The Go server requires a full set of environment variables. Use
`.env.example` as the source of truth. The most commonly edited values are:

| Variable                      | Purpose                                                   |
| ----------------------------- | --------------------------------------------------------- |
| ZKSYNC_RPC_URL                | L2 HTTP RPC endpoint                                      |
| L1_RPC_URL                    | L1 HTTP RPC endpoint (local-setup: http://localhost:8545) |
| L1_CHAIN_ID                   | L1 chain ID (0 = auto-detect from RPC)                    |
| L1_GAS_PRICE_WEI              | L1 gas price override for deployments and L1->L2 base cost |
| L1_MAILBOX_ADDRESS            | zkSync L1 mailbox / bridgehub address (auto-populated)    |
| L1_USE_DIRECT_MESSAGING       | Use bridgehub direct L1->L2 flow by default               |
| ZKSYNC_CHAIN_ID               | L2 chain ID                                               |
| ZKSYNC_PRIVATE_KEY            | deployer/signer key                                       |
| HTTP_BIND_ADDR                | API bind address                                          |
| NODE_COUNT                    | committee size n; must equal 2^SPARSE_TREE_DEPTH          |
| SPARSE_TREE_DEPTH             | validator-state tree depth (max 8)                        |
| BATCH_SIZE                    | claims per verification round b (paper §IV-D); validated against the deployed Gateway's immutable batch size at server startup |
| INC_TREE_DEPTH                | commitment (burn) tree depth                              |
| DESTINATION_ID                | destination rollup identifier d_dst bound into commitments |
| INITIAL_VALIDATOR_STAKE       | stake locked per validator at registration                |
| BURN_CONFIRMATION_DEPTH       | blocks before a burn counts as final (source rollup)      |
| AGGREGATOR_TIMEOUT            | seconds without a finalized round before validators may rotate the aggregator |
| HUB_WITHDRAW_DELAY            | Hub waiting period (seconds) between exit finalization and stake release |
| IPFS_API_URL                  | real IPFS daemon API; empty = simulated DFS               |
| COUNTER_CONTRACT_ADDRESS      | Counter contract address                                  |
| L1_MESSENGER_CONTRACT_ADDRESS | L1 messenger demo contract address                        |
| L2_MESSENGER_CONTRACT_ADDRESS | L2 messenger demo contract address                        |
| ORACLE_CONTRACT_ADDRESS       | Gateway (Oracle) contract address                         |
| MIMC_LIBRARY_ADDRESS          | deployed MiMC library linked into Oracle/MerkleTree       |
| L1_HUB_CONTRACT_ADDRESS       | L1 Hub contract address                                   |

The deploy targets update the contract address fields (including the MiMC
library and the L1 mailbox/Hub) inside `.env` automatically, reading each
address from the forge broadcast journal rather than console output.

## Makefile Shortcuts

```sh
make up          # start local-setup
make down        # stop and clear local-setup
make up-deploy   # start local-setup + deploy contracts
make deploy      # reset local state + deploy all contracts (verifiers, MiMC, Gateway, messengers, L1 Hub)
make deploy-mimc # deploy the MiMC library alone (pinned into Oracle/MerkleTree)
make reset-local-state # clear persisted user burn notes + simulated DFS
make deploy-messengers # deploy L2 + L1 messenger demo and link them
make configure-l2-messenger # set L1 messenger on L2 (if needed)
make server      # run the Go API server
```

## Tests and Evaluation

```sh
cd contracts && forge test        # Gateway conformance tests
cd go-server && go test ./...     # Go unit tests (thresholds, batching, ordering, crypto vectors)
```

For an automated end-to-end run (deploy, server, register → burn → claim →
round finalization → validator replace/exit/withdraw) against a running
stack, use the simulation driver:

```sh
./simulate.sh -r 1 -n 8 -b 4 -s 1   # 1 run, 8 validators, batch size 4
```

To regenerate the paper's §VI evaluation data (Tables I–II, Figures 5–8),
run the sweep against a running local stack (`make up-deploy`):

```sh
./scripts/run-evaluation.sh                 # 50 trials, n ∈ {4..128}, b ∈ {1,5,10,15}
./scripts/run-evaluation.sh -r 10 -n "4 8" -b "1 5"   # smaller sweep
```

Each configuration's raw CSV logs are copied to
`eval-results/<timestamp>/n{N}_b{B}/` and aggregated into
`eval-results/<timestamp>/summary.csv` (per-metric count/mean/std by
configuration and measurement source).

The transaction CSV separates L1, L2, and cross-layer message costs. In
particular, validator lifecycle measurements are logged as distinct rows for
L1 staking/requests, L2 execution, and L1 finalization:

- `L1 validator stake registration node=<id>`: stake registration on the L1 Hub.
- `L1->L2 validator import request chunk=<i>/<k> count=<m>`: L1 request that imports a bounded chunk of validators to L2. Large committees are split into chunks instead of one oversized L1->L2 message.
- `L2 validator import execution from L1 count=<m>`: L2 execution of one validator-import chunk; this is the step that inserts those validators into the L2 validator tree.
- `L2->L1 validator import result finalization on L1`: L1 finalization of the import result.
- `L2 submitWiVote batch finalization`: L2 batch-finalization transaction.
- `L2->L1 checkpoint finalization on L1`: L1 finalization of the submitted checkpoint.
- `L1->L2 validator replacement request index=<id>`: L1 replacement request.
- `L2 validator replacement execution from L1`: L2 execution of replacement.
- `L2->L1 validator replacement result finalization on L1`: L1 finalization of replacement result.
- `L2 validator exit request index=<id>` and `L2 validator withdraw request index=<id>`: L2 validator lifecycle requests.
- `L2->L1 validator exit finalization on L1` and `L2->L1 validator withdraw finalization on L1`: L1 finalization of the corresponding L2 messages.

## L1/L2 Messaging Demo

`make deploy`/`make deploy-messengers` auto-populate `L1_MAILBOX_ADDRESS`
(zkSync L1 mailbox/bridgehub) from `ZKSYNC_RPC_URL` and auto-set
`L1_USE_DIRECT_MESSAGING` based on RPC capabilities; set them manually only
for networks the probe cannot reach. Then use `sendToL2` / `sendToL1` on the
deployed demo contracts (or the `/messaging/*` endpoints below).

### Messaging API (curl examples)

L1 -> L2 (send + read on L2). If `value_wei` is empty or `0`, the server
estimates the L1 base fee automatically:

```sh
curl -X POST http://localhost:18000/messaging/l1/send \
  -H "Content-Type: application/json" \
  -d '{"message":"hello from L1","l2_gas_limit":2000000,"l2_gas_per_pubdata":800,"value_wei":"0"}'

curl http://localhost:18000/messaging/l2/last-from-l1
```

<!-- L1 -> L2 (zk-chains mode, explicit L2 chain id):
```sh
curl -X POST http://localhost:18000/messaging/l1/send-direct \
  -H "Content-Type: application/json" \
  -d '{"l2_chain_id":271,"message":"hello from L1","l2_gas_limit":2000000,"l2_gas_per_pubdata":800,"value_wei":"0"}'
``` -->

L2 -> L1 (send + read on L1; the server relays the zkSync inclusion proof
automatically — allow a few seconds before reading):

```sh
curl -X POST http://localhost:18000/messaging/l2/send \
  -H "Content-Type: application/json" \
  -d '{"message":"hello from L2"}'

# To relay a message manually instead (proof via zks_getL2ToL1LogProof):
# curl -X POST http://localhost:18000/messaging/l1/receive \
#   -H "Content-Type: application/json" \
#   -d '{"message":"hello from L2","l2_block_number":123,"l2_log_index":0,"l2_tx_number_in_block":1,"proof":["0x..."]}'

curl http://localhost:18000/messaging/l1/last-from-l2
```

## Workflow

Terminal 1 (blockchain and server)
```sh
make down  # optional
make up-deploy
make server
```

Terminal 2 (requests)

> A verification round finalizes (and mints accepted claims) only after
> `BATCH_SIZE` claims have been submitted (paper SIV-D: round `r` covers claim
> identifiers `[r*b, (r+1)*b-1]`). For this single burn/withdraw walkthrough
> set `BATCH_SIZE=1` in `.env` before `make up-deploy`; with a larger batch
> size, repeat the burn and withdraw calls `BATCH_SIZE` times — the balance
> credits when the round's Aggregating proof is verified on-chain.

```sh
curl -X POST http://localhost:18000/users/default/register -H "Content-Type: application/json"
curl -X POST http://localhost:18000/validators/register -H "Content-Type: application/json"
curl -X POST http://localhost:18000/users/default/burn -H "Content-Type: application/json"
curl -X POST http://localhost:18000/users/default/withdraw -H "Content-Type: application/json"
curl http://localhost:18000/users/default/balance
curl -X POST http://localhost:18000/validators/replace -H "Content-Type: application/json" -d '{"node_id":0,"replace_with_account_id":1}'
curl -X POST http://localhost:18000/validators/exit -H "Content-Type: application/json" -d '{"node_id":0}'
curl -X POST http://localhost:18000/validators/withdraw -H "Content-Type: application/json" -d '{"node_id":0}'
```

> `validators/register` must run before the first burn: the aggregator can
> only publish commitment roots once the validator set is registered
> on-chain.

With `L1_HUB_CONTRACT_ADDRESS` configured (set automatically by
`make deploy`), every finalized round's `CHECKPOINT` is relayed to the L1
Hub and consumed by `finalizeFromL2` with a zkSync inclusion proof. The
validator import, replacement, exit, and withdrawal flows are also relayed
and logged as separate L1/L2 CSV rows so lifecycle costs can be reported
without mixing them into the base L2 operation rows.

### Troubleshooting

- **Server exits with `could not read the Gateway's batch size`** — either
  `ORACLE_CONTRACT_ADDRESS` in `.env` is stale (redeploy with
  `make deploy`), or the local zkSync node lost recently deployed state
  (it can roll back fresh miniblocks when its L1 view is inconsistent);
  reset the stack with `make down && make up-deploy`. The server retries
  the read for ~30 s before giving up, so a briefly lagging node recovers
  on its own.
- **Server exits with `BATCH_SIZE=x does not match the deployed Gateway's
  batch size y`** — the Gateway's batch size is immutable; after changing
  `BATCH_SIZE` in `.env`, redeploy (`make deploy` or `make deploy-oracle`).
- **Withdrawals rejected with `commitment root not recorded` after a
  redeploy** — stale local state from a previous deployment; run
  `make reset-local-state` (included in `make deploy`) before
  `make server`.

## API Endpoints

| Method | Path                    | Example                                                                                                                                     |
| ------ | ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| GET    | /counter                | `curl http://localhost:18000/counter`                                                                                                       |
| POST   | /counter/increment      | `curl -X POST http://localhost:18000/counter/increment`                                                                                     |
| GET    | /health                 | `curl http://localhost:18000/health`                                                                                                        |
| POST   | /circuits/keygen        | `curl -X POST http://localhost:18000/circuits/keygen -H "Content-Type: application/json"`                                                   |
| GET    | /circuits/keys          | `curl http://localhost:18000/circuits/keys`                                                                                                 |
| POST   | /users/default/register | `curl -X POST http://localhost:18000/users/default/register -H "Content-Type: application/json"`                                            |
| GET    | /users/default/balance  | `curl http://localhost:18000/users/default/balance`                                                                                         |
| POST   | /users/default/burn     | `curl -X POST http://localhost:18000/users/default/burn -H "Content-Type: application/json"`                                                |
| POST   | /users/default/withdraw | `curl -X POST http://localhost:18000/users/default/withdraw -H "Content-Type: application/json"`                                            |
| POST   | /validators/register    | `curl -X POST http://localhost:18000/validators/register -H "Content-Type: application/json"`                                               |
| POST   | /validators/replace     | `curl -X POST http://localhost:18000/validators/replace -H "Content-Type: application/json" -d '{"node_id":0,"replace_with_account_id":1}'` |
| POST   | /validators/exit        | `curl -X POST http://localhost:18000/validators/exit -H "Content-Type: application/json" -d '{"node_id":0}'`                                |
| POST   | /validators/withdraw    | `curl -X POST http://localhost:18000/validators/withdraw -H "Content-Type: application/json" -d '{"node_id":0}'`                            |

Notes:
- `GET /users/default/balance` returns `token_one` = burn balance, `token_two` = claim balance.
