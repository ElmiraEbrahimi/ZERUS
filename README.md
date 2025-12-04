# ZERUS: A Zero-Knowledge Embedded Rollup Architecture for Cross-Rollup Interoperability


## Features

* **Modular Go architecture** – code is organised into small packages
  under `internal/` and `cmd/` with clear responsibilities. The
  architecture is intentionally simple so you can layer in more
  complexity as needed.
* **Ethereum and zkSync clients** – thin wrappers around the
  `go‑ethereum` RPC client and a generic JSON‑RPC client for zkSync
  operations. Both clients are configured via environment variables
  (RPC URL, chain ID and timeouts).
* **Bridge abstraction** – a service in `internal/bridge` that exposes
  high‑level functions such as `SendL1ToL2Message` and `TrackStatus`. At
  this stage the implementation simply records messages in memory and
  marks them as sent. It’s designed so that you can later plug in
  contract calls and event monitoring without refactoring callers.
* **HTTP API** – a small server in `internal/api` exposes the bridge
  functionality via JSON endpoints: `POST /bridge/l1-to-l2`, `POST
  /bridge/l2-to-l1` and `GET /bridge/status/{id}`.
* **Containerised development environment** – a `docker‑compose.yml`
  brings up Postgres, a local Ethereum devnet using Reth and the Go
  service. The Reth container is configured according to the devnet
  example which runs a node in development mode with one second block
  times and exposes JSON‑RPC/WebSocket interfaces on ports 8545 and
  8546【155669985012267†L111-L124】. If you wish to run a local zkSync
  node as well, see the optional section below.

## Getting Started

### Prerequisites

* Docker and docker‑compose installed on your machine.
* Go ≥1.20 if you intend to run the service outside of Docker.

### Clone and bootstrap

```sh
git clone <this‑repo>
cd bridge-service

# Copy the environment template and adjust values if needed
cp .env.example .env

# Build dependencies
go mod tidy

# Build and run tests
make build
make test
```

### Running locally with Docker

Start the full stack:

```sh
make up
```

This command starts Postgres, the Reth Ethereum node and the Go API. The API
listens on `http://localhost:8080` by default (configurable via
`APP_PORT`). The Reth devnet exposes JSON‑RPC on `http://localhost:8545` and
WebSocket on `ws://localhost:8546` using the configuration borrowed from
the Reth devnet example【155669985012267†L111-L124】. The default zkSync
variables in `.env.example` point to `http://localhost:3050` and `ws://localhost:3051`,
which correspond to the L2 RPC and WebSocket endpoints defined in the
official local setup guide【344146198151240†L365-L372】. If you are not
running a local zkSync node, point `ZKSYNC_RPC_URL` to a testnet or
mainnet RPC provider.

To stop and remove the containers:

```sh
make down
```

### Optional: Running a local zkSync node

Matter Labs publishes a [local setup](https://github.com/matter-labs/local-setup)
repository that runs a full zkSync Era stack (Postgres, local L1 via Reth and
the zkSync server). The guide lists the default RPC endpoints: L1 RPC
`http://localhost:8545` and L2 RPC `http://localhost:3050`【344146198151240†L365-L372】.
To integrate such a node into this project you have two options:

1. Clone the `local-setup` repo and run its `start.sh` script. Update your
   `.env` file so `ZKSYNC_RPC_URL` points to `http://localhost:3050` and
   restart the API container.
2. Uncomment the `zksync` service in `docker-compose.yml` and provide a
   compatible image. Matter Labs publishes a `local-node` image on Docker
   Hub, but note that it is large and subject to change. Ensure the
   `DATABASE_URL` and `ETH_CLIENT_WEB3_URL` variables are set appropriately
   (see comments in `docker-compose.yml`).

### Using the HTTP API

With the stack running you can exercise the endpoints with `curl` or any
HTTP client. The examples below assume the API is running on port 8080.

**Send a message from L1 to L2**

```sh
curl -X POST -H "Content-Type: application/json" \
  -d '{"from":"0xSender","to":"0xRecipient","data":"0x"}' \
  http://localhost:8080/bridge/l1-to-l2

# Response:
{"id":"<uuid>"}
```

**Send a message from L2 to L1**

```sh
curl -X POST -H "Content-Type: application/json" \
  -d '{"from":"0xSender","to":"0xRecipient","data":"0x"}' \
  http://localhost:8080/bridge/l2-to-l1
```

**Check message status**

```sh
curl http://localhost:8080/bridge/status/<uuid>
```

Initially the status will be `pending`; after a short delay it becomes
`sent`. In a real implementation you would update this status based on
transaction receipts and contract events.

### Project Structure

```
bridge-service/
├── cmd/               # entrypoints
│   └── api/main.go    # starts the HTTP server
├── internal/          # private application code
│   ├── api/           # HTTP handlers
│   ├── bridge/        # cross-layer messaging logic
│   │   └── contracts/ # interfaces for future contract bindings
│   ├── config/        # environment configuration loader
│   ├── db/            # database connection setup
│   ├── eth/           # Ethereum L1 client wrapper
│   ├── log/           # logger initialisation
│   └── zksync/        # zkSync L2 client wrapper
├── contracts/         # Solidity sources and deployment scripts
│   ├── src/
│   ├── deploy/
│   └── README.md
├── scripts/           # helper scripts (migrations, setup)
├── docker-compose.yml # orchestrates local services
├── Dockerfile         # builds the Go binary
├── .env.example       # sample environment configuration
├── Makefile           # common development commands
├── go.mod / go.sum    # Go module definitions
└── README.md          # this document
```
