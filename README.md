
# zkSync Local Stack + Counter Contract + Go HTTP API

This repository extends the Matter Labs **dockerized L2 stack** with:

- A **Foundry project** containing a simple `Counter` contract  
- A **deterministic deployment script** for zkSync local chain (chain ID `271`)  
- A **Go HTTP API server** (Go 1.24, geth 1.16.7) exposing REST endpoints that interact with the Counter contract  
- A **Dockerized Go service** included in the zk-chains Docker Compose stack  

## Table of Contents

1. Prerequisites
2. Repository Structure
3. Start the zkSync Local Environment
4. Deploy the Counter Contract
5. Generate Go Bindings (Optional)
6. Start the Go HTTP API
7. API Usage Examples
8. Environment Variables
9. Makefile Shortcuts

## Prerequisites

Install:

- Docker + Docker Compose
- Foundry (forge/cast)
- Go 1.24.x
- abigen (optional)

## Repository Structure

```
dockerized_l2/
├── local-setup/
├── contracts/
│   ├── src/Counter.sol
│   ├── script/DeployCounter.s.sol
│   └── foundry.toml
├── go-server/
│   ├── cmd/server/main.go
│   ├── internal/{config,eth,server}
│   ├── Dockerfile
│   └── go.mod
└── Makefile
```

## Start the zkSync Local Environment

```
make up
```

Endpoints:

| Component     | URL                    |
| ------------- | ---------------------- |
| L2 JSON-RPC   | http://localhost:15100 |
| L2 WS         | ws://localhost:15101   |
| Explorer      | http://localhost:15005 |
| Hyperexplorer | http://localhost:15000 |

## Deploy the Counter Contract

```
export ZKSYNC_PRIVATE_KEY=<private key>
export ZKSYNC_RPC_URL=http://localhost:15100

cd contracts
forge script script/DeployCounter.s.sol --rpc-url $ZKSYNC_RPC_URL --broadcast --private-key $ZKSYNC_PRIVATE_KEY
```

Capture the deployed contract address and export:

```
export COUNTER_CONTRACT_ADDRESS=0x...
```

## Generate Go Bindings (Optional)

```
cd contracts
abigen   --abi out/Counter.sol/Counter.abi.json   --bin out/Counter.sol/Counter.bin   --pkg eth   --type Counter   --out ../go-server/internal/eth/counter_binding.go
```

## Start the Go HTTP API

```
export ZKSYNC_RPC_URL=http://localhost:15100
export ZKSYNC_CHAIN_ID=271
export ZKSYNC_PRIVATE_KEY=0x...
export COUNTER_CONTRACT_ADDRESS=0x...
export HTTP_BIND_ADDR=:18000

cd go-server
go run ./cmd/server
```

## API Usage

### Read counter value

```
curl http://localhost:18000/counter
```

### Increment counter

```
curl -X POST http://localhost:18000/counter/increment
```

### Health check

```
curl http://localhost:18000/health
```

## Environment Variables

| Variable                 | Purpose            |
| ------------------------ | ------------------ |
| ZKSYNC_RPC_URL           | RPC endpoint       |
| ZKSYNC_CHAIN_ID          | Must be 271        |
| ZKSYNC_PRIVATE_KEY       | Signer private key |
| COUNTER_CONTRACT_ADDRESS | Contract address   |
| HTTP_BIND_ADDR           | Bind address       |

## Makefile Shortcuts

```sh
make up
make down
make deploy
make server
make build-docker
```

## Install Forge

```sh
curl -L https://foundry.paradigm.xyz | bash
```

and then:

```sh
foundryup
```

### foundry-zksync

```sh
curl -L https://raw.githubusercontent.com/matter-labs/foundry-zksync/main/install-foundry-zksync | bash
```

```sh
foundryup-zksync
```

## Run

1. make zksync  # wait for it set up
2. make deploy  # deploys and copies the contract address to go-server
3. make server  # runs the server

### test go-server api

```sh
curl http://localhost:18000/health
curl http://localhost:18000/counter
curl -X POST http://localhost:18000/counter/increment 
curl http://localhost:18000/counter

curl -X POST http://localhost:18000/circuits/keygen -H "Content-Type: application/json"
curl http://localhost:18000/circuits/keys 
```
