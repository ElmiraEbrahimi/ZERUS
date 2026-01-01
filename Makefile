.PHONY: server
server:
	cd go-server && go run ./cmd/server

.PHONY: zksync
zksync:
	@echo "Starting zkSync local stack via start.sh..."
	@cd local-setup && chmod +x ./start.sh && ./start.sh

.PHONY: up-deploy
up-deploy:
	cd local-setup && ./start.sh
	$(MAKE) deploy

.PHONY: up
up:
	cd local-setup && ./start.sh

.PHONY: down
down:
	cd local-setup && ./clear.sh

.PHONY: deploy
deploy: deploy-counter deploy-merkle-verifier deploy-votingbatch-verifier deploy-merkle-tree deploy-oracle

.PHONY: deploy-counter
deploy-counter:
	@echo "Deploying Counter via forge-zksync..."
	@cd contracts && set -a && . ../.env && set +a && \
	  forge script script/DeployCounter.s.sol --zksync --rpc-url "$$ZKSYNC_RPC_URL" \
	    --private-key "$$ZKSYNC_PRIVATE_KEY" --suppress-warnings assemblycreate --broadcast \
	    2>&1 | tee /tmp/forge_deploy.log

	@addr=$$(awk '/Counter deployed at/ {print $$4}' /tmp/forge_deploy.log | tail -n1); \
	if [ -z "$$addr" ]; then \
		echo "ERROR: deployment succeeded but could not extract Counter address from forge output" >&2; \
		exit 1; \
	fi; \
	  echo "Detected Counter contract address: $$addr"; \
	if [ -f .env ]; then \
		awk -v addr="$$addr" 'BEGIN{updated=0} \
		/^COUNTER_CONTRACT_ADDRESS=/{print "COUNTER_CONTRACT_ADDRESS=" addr; updated=1; next} \
		{print} \
		END{if(!updated) print "COUNTER_CONTRACT_ADDRESS=" addr}' .env > .env.tmp && mv .env.tmp .env; \
	else \
		echo "COUNTER_CONTRACT_ADDRESS=$$addr" > .env; \
	fi; \
	  echo "Updated .env with COUNTER_CONTRACT_ADDRESS=$$addr"

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
FORCE_ZK_KEYGEN ?= 1
gnark-merkle-verifier:
	@echo "Generating MerkleProof Solidity verifier via gnark..."
	@cd go-server && FORCE_ZK_KEYGEN=$(FORCE_ZK_KEYGEN) go run ./circuits/merkle_proof/cmd/gen_verifier/main.go

deploy-merkle-verifier: gnark-merkle-verifier
	@echo "Deploying MerkleProofVerifier via forge-zksync..."
	@cd contracts && set -a && . ../.env && set +a && \
	  forge script script/DeployMerkleProofVerifier.s.sol --zksync --rpc-url "$$ZKSYNC_RPC_URL" \
	    --private-key "$$ZKSYNC_PRIVATE_KEY" --suppress-warnings assemblycreate --broadcast \
	    2>&1 | tee /tmp/forge_merkle_deploy.log

	@addr=$$(awk '/MerkleProofVerifier deployed at/ {print $$NF}' /tmp/forge_merkle_deploy.log | tail -n1); \
	if [ -z "$$addr" ]; then \
		echo "ERROR: deployment succeeded but could not extract MerkleProofVerifier address from forge output" >&2; \
		exit 1; \
	fi; \
	  echo "Detected MerkleProofVerifier contract address: $$addr"; \
	if [ -f .env ]; then \
		awk -v addr="$$addr" 'BEGIN{updated=0} \
		/^MERKLE_VERIFIER_ADDRESS=/{print "MERKLE_VERIFIER_ADDRESS=" addr; updated=1; next} \
		{print} \
		END{if(!updated) print "MERKLE_VERIFIER_ADDRESS=" addr}' .env > .env.tmp && mv .env.tmp .env; \
	else \
		echo "MERKLE_VERIFIER_ADDRESS=$$addr" > .env; \
	fi; \
	  echo "Updated .env with MERKLE_VERIFIER_ADDRESS=$$addr"

	@echo "Generating Go bindings for MerkleProofVerifier via abigen..."
	@cd contracts && forge build --extra-output-files abi >/dev/null
	@abigen \
	  --abi contracts/out/MerkleProofVerifier.sol/Verifier.abi.json \
	  --pkg eth \
	  --type MerkleProofVerifier \
	  --out go-server/internal/eth/merkle_verifier_abigen.go
	@echo "Regenerated go-server/internal/eth/merkle_verifier_abigen.go"
	@echo "MerkleProof verifier deployment + ABI regeneration complete."

gnark-votingbatch-verifier:
	@echo "Generating VotingBatch Solidity verifier via gnark..."
	@cd go-server && FORCE_ZK_KEYGEN=$(FORCE_ZK_KEYGEN) go run ./circuits/voting_batch/cmd/gen_votingbatch_verifier/main.go

deploy-votingbatch-verifier: gnark-votingbatch-verifier
	@echo "Deploying VotingBatchVerifier via forge-zksync..."
	@cd contracts && set -a && . ../.env && set +a && \
	  forge script script/DeployVotingBatchVerifier.s.sol \
	    --zksync --rpc-url "$$ZKSYNC_RPC_URL" --private-key "$$ZKSYNC_PRIVATE_KEY" \
	    --suppress-warnings assemblycreate --broadcast 2>&1 | tee /tmp/forge_votingbatch_deploy.log

	@addr=$$(awk '/VotingBatchVerifier deployed at/ {print $$NF}' /tmp/forge_votingbatch_deploy.log | tail -n1); \
	if [ -z "$$addr" ]; then \
		echo "ERROR: deployment succeeded but could not extract VotingBatchVerifier address from forge output" >&2; \
		exit 1; \
	fi; \
	  echo "Detected VotingBatchVerifier contract address: $$addr"; \
	if [ -f .env ]; then \
		awk -v addr="$$addr" 'BEGIN{updated=0} \
		/^VOTING_BATCH_VERIFIER_ADDRESS=/{print "VOTING_BATCH_VERIFIER_ADDRESS=" addr; updated=1; next} \
		{print} \
		END{if(!updated) print "VOTING_BATCH_VERIFIER_ADDRESS=" addr}' .env > .env.tmp && mv .env.tmp .env; \
	else \
		echo "VOTING_BATCH_VERIFIER_ADDRESS=$$addr" > .env; \
	fi; \
	  echo "Updated .env with VOTING_BATCH_VERIFIER_ADDRESS=$$addr"

	@echo "Generating Go bindings for VotingBatchVerifier via abigen..."
	@cd contracts && forge build --extra-output-files abi >/dev/null
	@abigen \
	  --abi contracts/out/VotingBatchVerifier.sol/Verifier.abi.json \
	  --pkg eth \
	  --type VotingBatchVerifier \
	  --out go-server/internal/eth/votingbatch_verifier_abigen.go
	@echo "Regenerated go-server/internal/eth/votingbatch_verifier_abigen.go"
	@echo "Voting batch verifier deployment + ABI regeneration complete."

.PHONY: deploy-oracle
deploy-oracle:
	@echo "Deploying Oracle via forge-zksync..."
	@cd contracts && set -a && . ../.env && set +a && \
	  forge script script/DeployOracle.s.sol --zksync --rpc-url "$$ZKSYNC_RPC_URL" \
	    --private-key "$$ZKSYNC_PRIVATE_KEY" --suppress-warnings assemblycreate --broadcast \
	    2>&1 | tee /tmp/forge_oracle_deploy.log

	@addr=$$(awk '/Oracle deployed at/ {print $$4}' /tmp/forge_oracle_deploy.log | tail -n1); \
	if [ -z "$$addr" ]; then \
		echo "ERROR: deployment succeeded but could not extract Oracle address from forge output" >&2; \
		exit 1; \
	fi; \
	echo "Detected Oracle contract address: $$addr"; \
	if [ -f .env ]; then \
		awk -v addr="$$addr" 'BEGIN{updated=0} \
		/^ORACLE_CONTRACT_ADDRESS=/{print "ORACLE_CONTRACT_ADDRESS=" addr; updated=1; next} \
		{print} \
		END{if(!updated) print "ORACLE_CONTRACT_ADDRESS=" addr}' .env > .env.tmp && mv .env.tmp .env; \
	else \
		echo "ORACLE_CONTRACT_ADDRESS=$$addr" > .env; \
	fi; \
	echo "Updated .env with ORACLE_CONTRACT_ADDRESS=$$addr"

	@echo "Generating Go bindings via abigen..."
	@cd contracts && forge build --extra-output-files abi >/dev/null
	@abigen \
	  --abi contracts/out/oracle.sol/Oracle.abi.json \
	  --pkg eth \
	  --type Oracle \
	  --out go-server/internal/eth/oracle_abigen.go
	@echo "Regenerated go-server/internal/eth/oracle_abigen.go"
	@echo "Deployment + ABI regeneration complete."
.PHONY: deploy-merkle-tree
deploy-merkle-tree:
	@echo "Deploying MerkleTree via forge-zksync..."
	@cd contracts && set -a && . ../.env && set +a && \
	  forge script script/DeployMerkleTree.s.sol --zksync --rpc-url "$$ZKSYNC_RPC_URL" \
	    --private-key "$$ZKSYNC_PRIVATE_KEY" --suppress-warnings assemblycreate --broadcast \
	    2>&1 | tee /tmp/forge_merkle_tree_deploy.log

	@addr=$$(awk '/MerkleTree deployed at/ {print $$NF}' /tmp/forge_merkle_tree_deploy.log | tail -n1); \
	if [ -z "$$addr" ]; then \
		echo "ERROR: deployment succeeded but could not extract MerkleTree address from forge output" >&2; \
		exit 1; \
	fi; \
	echo "Detected MerkleTree contract address: $$addr"; \
 \
	if [ -f .env ]; then \
		awk -v addr="$$addr" 'BEGIN{updated=0} \
		/^MERKLE_TREE_CONTRACT_ADDRESS=/{print "MERKLE_TREE_CONTRACT_ADDRESS=" addr; updated=1; next} \
		{print} \
		END{if(!updated) print "MERKLE_TREE_CONTRACT_ADDRESS=" addr}' .env > .env.tmp && mv .env.tmp .env; \
	else \
		echo "MERKLE_TREE_CONTRACT_ADDRESS=$$addr" > .env; \
	fi; \
	echo "Updated .env with MERKLE_TREE_CONTRACT_ADDRESS=$$addr"

	@echo "Generating Go bindings via abigen..."
	@cd contracts && forge build --extra-output-files abi >/dev/null
	@abigen \
	  --abi contracts/out/merkle_tree.sol/MerkleTree.abi.json \
	  --pkg eth \
	  --type MerkleTree \
	  --out go-server/internal/eth/merkle_tree_abigen.go
	@echo "Regenerated go-server/internal/eth/merkle_tree_abigen.go"
	@echo "MerkleTree deployment + ABI regeneration complete."

# Utils:

.PHONY: zip
zip:
	@TIMESTAMP=$$(date +"%Y%m%d_%H%M%S"); \
	zip -r archive_$$TIMESTAMP.zip . \
		-x ".git/*" \
		-x "go-server/circuits/build/*" \
		-x "*.csv" \
		-x "*.gob" \
		-x "*.log" \
		-x "*.zip"
