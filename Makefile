.PHONY: server
server:
	cd go-server && go run ./cmd/server

.PHONY: zksync
zksync:
	@echo "Starting zkSync local stack via start.sh..."
	@cd local-setup && chmod +x ./start.sh && ./start.sh

.PHONY: down
down:
	cd local-setup && docker compose -f zk-chains-docker-compose.yml down

# Deploy the Counter contract using forge-zksync and automatically
# update go-server/.env with the deployed COUNTER_CONTRACT_ADDRESS,
# then regenerate Go bindings via abigen.
#
# Requirements:
#   - contracts/.env defines:
#       DEPLOYER_PRIVATE_KEY
#       ZKSYNC_RPC_URL
#   - forge-zksync installed
#   - abigen installed
#
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

	@echo "Generating Go bindings via abigen..."
	@cd contracts && forge build --extra-output-files abi >/dev/null
	@abigen \
	  --abi contracts/out/Counter.sol/Counter.abi.json \
	  --pkg eth \
	  --type Counter \
	  --out go-server/internal/eth/counter_abigen.go
	@echo "Regenerated go-server/internal/eth/counter_abigen.go"
	@echo "Deployment + ABI regeneration complete."

.PHONY: gnark-merkle-verifier deploy-merkle-verifier \
	gnark-votingbatch-verifier deploy-votingbatch-verifier

# ------------------------------------------------------------------------------
# MerkleProof (redeeming_circuit) verifier
# ------------------------------------------------------------------------------

gnark-merkle-verifier:
	@echo "Generating MerkleProof Solidity verifier via gnark..."
	@cd go-server && go run ./circuits/merkle_proof/cmd/gen_verifier/main.go

deploy-merkle-verifier: gnark-merkle-verifier
	@echo "Deploying MerkleProofVerifier via forge-zksync..."
	@cd contracts && set -a && . ./.env && set +a && \
	  forge script script/DeployMerkleProofVerifier.s.sol --zksync --rpc-url "$$ZKSYNC_RPC_URL" --broadcast 2>&1 | tee /tmp/forge_merkle_deploy.log

	@addr=$$(awk '/MerkleProofVerifier deployed at/ {print $$NF}' /tmp/forge_merkle_deploy.log | tail -n1); \
	if [ -z "$$addr" ]; then \
		echo "ERROR: deployment succeeded but could not extract MerkleProofVerifier address from forge output" >&2; \
		exit 1; \
	fi; \
	  echo "Detected MerkleProofVerifier contract address: $$addr"; \
	if [ -f go-server/.env ]; then \
		awk -v addr="$$addr" 'BEGIN{updated=0} \
		/^MERKLE_VERIFIER_ADDRESS=/{print "MERKLE_VERIFIER_ADDRESS=" addr; updated=1; next} \
		{print} \
		END{if(!updated) print "MERKLE_VERIFIER_ADDRESS=" addr}' go-server/.env > go-server/.env.tmp && mv go-server/.env.tmp go-server/.env; \
	else \
		echo "MERKLE_VERIFIER_ADDRESS=$$addr" > go-server/.env; \
	fi; \
	  echo "Updated go-server/.env with MERKLE_VERIFIER_ADDRESS=$$addr"

	@echo "Generating Go bindings for MerkleProofVerifier via abigen..."
	@cd contracts && forge build --extra-output-files abi >/dev/null
	@abigen \
	  --abi contracts/out/MerkleProofVerifier.sol/Verifier.abi.json \
	  --pkg eth \
	  --type MerkleProofVerifier \
	  --out go-server/internal/eth/merkle_verifier_abigen.go
	@echo "Regenerated go-server/internal/eth/merkle_verifier_abigen.go"
	@echo "MerkleProof verifier deployment + ABI regeneration complete."

# ------------------------------------------------------------------------------
# Voting batch circuit (BatchingVotingCircuit) verifier
# ------------------------------------------------------------------------------

gnark-votingbatch-verifier:
	@echo "Generating VotingBatch Solidity verifier via gnark..."
	@cd go-server && go run ./circuits/voting_batch/cmd/gen_votingbatch_verifier/main.go

deploy-votingbatch-verifier: gnark-votingbatch-verifier
	@echo "Deploying VotingBatchVerifier via forge-zksync..."
	@cd contracts && set -a && . ./.env && set +a && \
	  forge script script/DeployVotingBatchVerifier.s.sol \
	    --zksync --rpc-url "$$ZKSYNC_RPC_URL" --broadcast 2>&1 | tee /tmp/forge_votingbatch_deploy.log

	@addr=$$(awk '/VotingBatchVerifier deployed at/ {print $$NF}' /tmp/forge_votingbatch_deploy.log | tail -n1); \
	  if [ -z "$$addr" ]; then \
	    echo "ERROR: deployment succeeded but could not extract VotingBatchVerifier address from forge output" >&2; \
	    exit 1; \
	  fi; \
	  echo "Detected VotingBatchVerifier contract address: $$addr"; \
	  if [ -f go-server/.env ]; then \
	    awk -v addr="$$addr" 'BEGIN{updated=0} \
	      /^VOTING_BATCH_VERIFIER_ADDRESS=/{print "VOTING_BATCH_VERIFIER_ADDRESS=" addr; updated=1; next} \
	      {print} \
	      END{if(!updated) print "VOTING_BATCH_VERIFIER_ADDRESS=" addr}' go-server/.env > go-server/.env.tmp && mv go-server/.env.tmp go-server/.env; \
	  else \
	    echo "VOTING_BATCH_VERIFIER_ADDRESS=$$addr" > go-server/.env; \
	  fi; \
	  echo "Updated go-server/.env with VOTING_BATCH_VERIFIER_ADDRESS=$$addr"

	@echo "Generating Go bindings for VotingBatchVerifier via abigen..."
	@cd contracts && forge build --extra-output-files abi >/dev/null
	@abigen \
	  --abi contracts/out/VotingBatchVerifier.sol/Verifier.abi.json \
	  --pkg eth \
	  --type VotingBatchVerifier \
	  --out go-server/internal/eth/votingbatch_verifier_abigen.go
	@echo "Regenerated go-server/internal/eth/votingbatch_verifier_abigen.go"
	@echo "Voting batch verifier deployment + ABI regeneration complete."