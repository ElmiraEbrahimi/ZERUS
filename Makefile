### Convenience Makefile for the zkSync Counter project

.PHONY: zksync
zksync:
	cd local-setup && docker compose -f zk-chains-docker-compose.yml up -d

.PHONY: down
down:
	cd local-setup && docker compose -f zk-chains-docker-compose.yml down

# Deploy the Counter contract using forge-zksync and automatically
# update go-server/.env with the deployed COUNTER_CONTRACT_ADDRESS.
#
# Requirements:
#   - contracts/.env defines:
#       DEPLOYER_PRIVATE_KEY
#       ZKSYNC_RPC_URL (e.g. http://localhost:3050)
#   - forge-zksync is installed and accessible as `forge`
#
# This target:
#   1. Runs the DeployCounter.s.sol script with --zksync and --broadcast.
#   2. Prints the full forge output.
#   3. On success, extracts the "Contract Address: 0x..." line.
#   4. Updates go-server/.env:
#        - if COUNTER_CONTRACT_ADDRESS exists, replaces its value
#        - otherwise, appends COUNTER_CONTRACT_ADDRESS=0x...
.PHONY: deploy
deploy:
	@echo "Deploying Counter via forge-zksync..."
	@cd contracts && set -a && . ./.env && set +a && \
	  forge script script/DeployCounter.s.sol --zksync --rpc-url "$$ZKSYNC_RPC_URL" --broadcast 2>&1 | tee /tmp/forge_deploy.log
	@addr=$$(awk '/Counter deployed at/ {print $$4}' /tmp/forge_deploy.log | tail -n1); \
	  if [ -z "$$addr" ]; then \
	    echo "ERROR: deployment succeeded but could not extract Counter address from forge output" >&2; \
	    exit 1; \
	  fi; \
	  echo "Detected Counter contract address: $$addr"; \
	  if [ -f go-server/.env ]; then \
	    awk -v addr="$$addr" 'BEGIN{updated=0} \
	      /^COUNTER_CONTRACT_ADDRESS=/{print "COUNTER_CONTRACT_ADDRESS=" addr; updated=1; next} \
	      {print} \
	      END{if(!updated) print "COUNTER_CONTRACT_ADDRESS=" addr}' go-server/.env > go-server/.env.tmp && mv go-server/.env.tmp go-server/.env; \
	  else \
	    echo "COUNTER_CONTRACT_ADDRESS=$$addr" > go-server/.env; \
	  fi; \
	  echo "Updated go-server/.env with COUNTER_CONTRACT_ADDRESS=$$addr"

# Run the Go API server locally.
#
# Reads configuration from go-server/.env (see go-server/.env.example). The
# server uses github.com/joho/godotenv to load this file automatically and then
# relies on internal/config.Load for validation.
.PHONY: server
server:
	cd go-server && go run ./cmd/server
