# L2-Alchemy

Local zkSync stack + Foundry contracts + a Go HTTP API that drives the Counter
contract and the on-chain Oracle flow.

## Prerequisites

Pinned toolchain (paper §V):

- Docker + Docker Compose
- Foundry `forge` 1.3.5 with the foundry-zksync fork (v0.1.5); Solidity is
  pinned to 0.8.30 in `contracts/foundry.toml`
- Go 1.24.x (gnark v0.14.0, gnark-crypto v0.19.0, go-ethereum v1.16.7 are
  pinned in `go-server/go.mod`)
- abigen (only needed for deploy targets that regenerate bindings)

## Repository Layout

```
contracts/    # Foundry contracts + deploy scripts
go-server/    # Go HTTP API + oracle runtime
local-setup/  # zkSync local stack scripts + compose files
Makefile      # handy wrapper targets
```

## Quickstart (single L2)

1. Create a local env file and update values as needed:

```sh
cp .env.example .env
```

For the default local stack, set `ZKSYNC_RPC_URL=http://localhost:3050` and
`ZKSYNC_CHAIN_ID=270`. You can grab a funded dev key from
`local-setup/rich-wallets.json`.

2. Start the local zkSync stack:

```sh
make up
```

3. Deploy contracts and regenerate Go bindings:

```sh
make deploy
```

`make deploy` requires `abigen` to be available in your PATH.

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
| L1_MAILBOX_ADDRESS            | zkSync L1 mailbox / bridgehub address                     |
| L1_USE_DIRECT_MESSAGING       | Use bridgehub direct L1->L2 flow by default               |
| ZKSYNC_CHAIN_ID               | L2 chain ID                                               |
| ZKSYNC_PRIVATE_KEY            | deployer/signer key                                       |
| HTTP_BIND_ADDR                | API bind address                                          |
| COUNTER_CONTRACT_ADDRESS      | Counter contract address                                  |
| L1_MESSENGER_CONTRACT_ADDRESS | L1 messenger demo contract address                        |
| L2_MESSENGER_CONTRACT_ADDRESS | L2 messenger demo contract address                        |
| ORACLE_CONTRACT_ADDRESS       | Oracle contract address                                   |

The deploy targets update contract address fields inside `.env` automatically.

## Makefile Shortcuts

```sh
make up          # start local-setup
make down        # stop and clear local-setup
make up-deploy   # start local-setup + deploy contracts
make deploy      # deploy core L2 contracts + L1/L2 messenger demo (requires L1_MAILBOX_ADDRESS)
make deploy-messengers # deploy L2 + L1 messenger demo and link them
make configure-l2-messenger # set L1 messenger on L2 (if needed)
make server      # run the Go API server
```

## Tests and Evaluation

```sh
cd contracts && forge test        # Gateway conformance tests
cd go-server && go test ./...     # Go unit tests (thresholds, batching, ordering, crypto vectors)
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

## L1/L2 Messaging Demo

Set `L1_MAILBOX_ADDRESS` to your zkSync L1 mailbox/bridgehub contract before deploying.
`make deploy`/`make deploy-messengers` will auto-populate it from `ZKSYNC_RPC_URL` if it is empty.
Set `L1_USE_DIRECT_MESSAGING=true` for bridgehub-based networks. The deploy script auto-sets it based on RPC capabilities.
Then run `make deploy-messengers` and use `sendToL2` / `sendToL1` on the deployed contracts.

### Messaging API (quick curl)

Set a base URL once:

```sh
BASE_URL=http://localhost:18000
```

L1 -> L2 (send + read on L2):

```sh
curl -X POST "$BASE_URL/messaging/l1/send" \
  -H "Content-Type: application/json" \
  -d '{"message":"ping from L1","l2_gas_limit":2000000,"l2_gas_per_pubdata":800,"value_wei":"0"}'

curl "$BASE_URL/messaging/l2/last-from-l1"
```

L2 -> L1 (send + read on L1):

```sh
curl -X POST "$BASE_URL/messaging/l2/send" \
  -H "Content-Type: application/json" \
  -d '{"message":"ping from L2"}'

# Wait a few seconds for the relay to finalize on L1.
curl "$BASE_URL/messaging/l1/last-from-l2"
```

### Messaging API (curl examples)

L1 -> L2 (single-chain mode):
If `value_wei` is empty or `0`, the server estimates the L1 base fee automatically.
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

L2 -> L1 (requires zkSync proof from L2 RPC):
```sh
curl -X POST http://localhost:18000/messaging/l2/send \
  -H "Content-Type: application/json" \
  -d '{"message":"hello from L2"}'

# After you fetch the proof (e.g. via zks_getL2ToL1MsgProof), call:
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
