# L2-Alchemy

Local zkSync stack + Foundry contracts + a Go HTTP API that drives the Counter
contract and the on-chain Oracle flow.

## Prerequisites

- Docker + Docker Compose
- Foundry + foundry-zksync
- Go 1.24.x
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

## Environment Variables

The Go server requires a full set of environment variables. Use
`.env.example` as the source of truth. The most commonly edited values are:

| Variable                 | Purpose                  |
| ------------------------ | ------------------------ |
| ZKSYNC_RPC_URL           | L2 HTTP RPC endpoint     |
| ZKSYNC_CHAIN_ID          | L2 chain ID              |
| ZKSYNC_PRIVATE_KEY       | deployer/signer key      |
| HTTP_BIND_ADDR           | API bind address         |
| COUNTER_CONTRACT_ADDRESS | Counter contract address |
| ORACLE_CONTRACT_ADDRESS  | Oracle contract address  |

The deploy targets update contract address fields inside `.env` automatically.

## Makefile Shortcuts

```sh
make up          # start local-setup
make down        # stop and clear local-setup
make up-deploy   # start local-setup + deploy contracts
make deploy      # deploy all contracts and regenerate bindings
make server      # run the Go API server
```

## Workflow

Terminal 1 (blockchain and server)
```sh
make down  # optional
make up-deploy
make server
```

Terminal 2 (requests)
```sh
curl -X POST http://localhost:18000/users/default/register -H "Content-Type: application/json"
curl -X POST http://localhost:18000/validators/register -H "Content-Type: application/json"
curl -X POST http://localhost:18000/users/default/burn -H "Content-Type: application/json"
curl -X POST http://localhost:18000/users/default/withdraw -H "Content-Type: application/json"
curl http://localhost:18000/users/default/balance
```

## API Endpoints

| Method | Path                    | Example                                                                                          |
| ------ | ----------------------- | ------------------------------------------------------------------------------------------------ |
| GET    | /counter                | `curl http://localhost:18000/counter`                                                            |
| POST   | /counter/increment      | `curl -X POST http://localhost:18000/counter/increment`                                          |
| GET    | /health                 | `curl http://localhost:18000/health`                                                             |
| POST   | /circuits/keygen        | `curl -X POST http://localhost:18000/circuits/keygen -H "Content-Type: application/json"`        |
| GET    | /circuits/keys          | `curl http://localhost:18000/circuits/keys`                                                      |
| POST   | /users/default/register | `curl -X POST http://localhost:18000/users/default/register -H "Content-Type: application/json"` |
| GET    | /users/default/balance  | `curl http://localhost:18000/users/default/balance`                                              |
| POST   | /users/default/burn     | `curl -X POST http://localhost:18000/users/default/burn -H "Content-Type: application/json"`     |
| POST   | /users/default/withdraw | `curl -X POST http://localhost:18000/users/default/withdraw -H "Content-Type: application/json"` |
| POST   | /validators/register    | `curl -X POST http://localhost:18000/validators/register -H "Content-Type: application/json"`    |
| POST   | /validators/replace     | `curl -X POST http://localhost:18000/validators/replace -H "Content-Type: application/json" -d '{"node_id":0,"replace_with_account_id":1}'` |
| POST   | /validators/exit        | `curl -X POST http://localhost:18000/validators/exit -H "Content-Type: application/json" -d '{"node_id":0}'` |
| POST   | /validators/withdraw    | `curl -X POST http://localhost:18000/validators/withdraw -H "Content-Type: application/json" -d '{"node_id":0}'` |

Notes:
- `GET /users/default/balance` returns `token_one` = burn balance, `token_two` = claim balance.
