ANNOUNCE_TARGET = @echo "==> $@"

.PHONY: server
server:
	$(ANNOUNCE_TARGET)
	cd go-server && go run ./cmd/server

.PHONY: run-index
run-index:
	$(ANNOUNCE_TARGET)
	cd go-server && go run ./cmd/runindex

.PHONY: zksync
zksync:
	$(ANNOUNCE_TARGET)
	@echo "Starting zkSync local stack via start.sh..."
	@cd local-setup && chmod +x ./start.sh && ./start.sh

.PHONY: l1
l1:
	$(ANNOUNCE_TARGET)
	@echo "Starting local L1 node via start.sh..."
	@cd local-setup && chmod +x ./start.sh && ./start.sh

.PHONY: up-deploy
up-deploy:
	$(ANNOUNCE_TARGET)
	cd local-setup && ./start.sh
	$(MAKE) deploy

.PHONY: up
up:
	$(ANNOUNCE_TARGET)
	cd local-setup && ./start.sh

.PHONY: down
down:
	$(ANNOUNCE_TARGET)
	cd local-setup && ./clear.sh

.PHONY: deploy
deploy: run-index deploy-counter deploy-merkle-verifier deploy-votingbatch-verifier deploy-merkle-tree deploy-oracle deploy-messengers deploy-l1hub-if-configured

.PHONY: deploy-counter deploy-l2-messenger deploy-l1-messenger configure-l2-messenger deploy-messengers deploy-messengers-if-configured abigen-messengers ensure-l1-mailbox
deploy-counter:
	$(ANNOUNCE_TARGET)
	@echo "Deploying Counter via forge-zksync..."
	@cd contracts && set -a && . ../.env && set +a && \
	  forge script script/DeployCounter.s.sol --zksync --rpc-url "$$ZKSYNC_RPC_URL" \
	    --private-key "$$ZKSYNC_PRIVATE_KEY" --suppress-warnings assemblycreate --broadcast \
	    2>&1 | tee /tmp/forge_deploy.log

	@addr=$$(scripts/deployed-address.sh DeployCounter.s.sol); \
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

deploy-l2-messenger:
	$(ANNOUNCE_TARGET)
	@echo "Deploying L2Messenger via forge-zksync..."
	@cd contracts && set -a && . ../.env && set +a && \
	  forge script script/DeployL2Messenger.s.sol --zksync --rpc-url "$$ZKSYNC_RPC_URL" \
	    --private-key "$$ZKSYNC_PRIVATE_KEY" --suppress-warnings assemblycreate --broadcast \
	    2>&1 | tee /tmp/forge_l2_messenger_deploy.log

	@addr=$$(scripts/deployed-address.sh DeployL2Messenger.s.sol); \
	if [ -z "$$addr" ]; then \
		echo "ERROR: deployment succeeded but could not extract L2Messenger address from forge output" >&2; \
		exit 1; \
	fi; \
	  echo "Detected L2Messenger contract address: $$addr"; \
	if [ -f .env ]; then \
		awk -v addr="$$addr" 'BEGIN{updated=0} \
		/^L2_MESSENGER_CONTRACT_ADDRESS=/{print "L2_MESSENGER_CONTRACT_ADDRESS=" addr; updated=1; next} \
		{print} \
		END{if(!updated) print "L2_MESSENGER_CONTRACT_ADDRESS=" addr}' .env > .env.tmp && mv .env.tmp .env; \
	else \
		echo "L2_MESSENGER_CONTRACT_ADDRESS=$$addr" > .env; \
	fi; \
	  echo "Updated .env with L2_MESSENGER_CONTRACT_ADDRESS=$$addr"

deploy-l1-messenger: ensure-l1-mailbox
	$(ANNOUNCE_TARGET)
	@echo "Deploying L1Messenger to L1 via forge..."
	@cd contracts && set -a && . ../.env && set +a && \
	  gas_price="$${L1_GAS_PRICE_WEI:-1000000000}"; \
	  forge script script/DeployL1Messenger.s.sol --rpc-url "$$L1_RPC_URL" \
	    --private-key "$$ZKSYNC_PRIVATE_KEY" --legacy --gas-price "$$gas_price" --broadcast \
	    2>&1 | tee /tmp/forge_l1_messenger_deploy.log

	@addr=$$(scripts/deployed-address.sh DeployL1Messenger.s.sol); \
	if [ -z "$$addr" ]; then \
		echo "ERROR: deployment succeeded but could not extract L1Messenger address from forge output" >&2; \
		exit 1; \
	fi; \
	  echo "Detected L1Messenger contract address: $$addr"; \
	if [ -f .env ]; then \
		awk -v addr="$$addr" 'BEGIN{updated=0} \
		/^L1_MESSENGER_CONTRACT_ADDRESS=/{print "L1_MESSENGER_CONTRACT_ADDRESS=" addr; updated=1; next} \
		{print} \
		END{if(!updated) print "L1_MESSENGER_CONTRACT_ADDRESS=" addr}' .env > .env.tmp && mv .env.tmp .env; \
	else \
		echo "L1_MESSENGER_CONTRACT_ADDRESS=$$addr" > .env; \
	fi; \
	  echo "Updated .env with L1_MESSENGER_CONTRACT_ADDRESS=$$addr"

configure-l2-messenger:
	$(ANNOUNCE_TARGET)
	@echo "Configuring L2Messenger with L1 address..."
	@cd contracts && set -a && . ../.env && set +a && \
	if [ -z "$$L1_MESSENGER_CONTRACT_ADDRESS" ] || [ -z "$$L2_MESSENGER_CONTRACT_ADDRESS" ]; then \
		echo "ERROR: L1_MESSENGER_CONTRACT_ADDRESS and L2_MESSENGER_CONTRACT_ADDRESS must be set before configuring L2Messenger." >&2; \
		exit 1; \
	fi; \
	  forge script script/ConfigureL2Messenger.s.sol --zksync --rpc-url "$$ZKSYNC_RPC_URL" \
	    --private-key "$$ZKSYNC_PRIVATE_KEY" --suppress-warnings assemblycreate --broadcast \
	    2>&1 | tee /tmp/forge_l2_messenger_config.log

abigen-messengers:
	$(ANNOUNCE_TARGET)
	@echo "Generating Go bindings for L1Messenger + L2Messenger via abigen..."
	@cd contracts && forge build --extra-output-files abi >/dev/null
	@abigen \
	  --abi contracts/out/L1Messenger.sol/L1Messenger.abi.json \
	  --pkg eth \
	  --type L1Messenger \
	  --out go-server/internal/eth/l1_messenger_abigen.go
	@abigen \
	  --abi contracts/out/L2Messenger.sol/L2Messenger.abi.json \
	  --pkg eth \
	  --type L2Messenger \
	  --out go-server/internal/eth/l2_messenger_abigen.go
	@echo "Regenerated go-server/internal/eth/l1_messenger_abigen.go"
	@echo "Regenerated go-server/internal/eth/l2_messenger_abigen.go"

deploy-messengers-if-configured:
	$(ANNOUNCE_TARGET)
	@if [ ! -f .env ]; then \
		echo "Skipping L1 messenger deploy (.env not found)"; \
		exit 0; \
	fi
	@set -a && . ./.env && set +a && \
	if [ -z "$$L1_MAILBOX_ADDRESS" ]; then \
		echo "Skipping L1 messenger deploy (L1_MAILBOX_ADDRESS not set)"; \
		exit 0; \
	fi; \
	  $(MAKE) deploy-l1-messenger configure-l2-messenger

deploy-messengers: deploy-l2-messenger deploy-l1-messenger configure-l2-messenger abigen-messengers

.PHONY: deploy-l1hub deploy-l1hub-if-configured
# Deploy the L1 Hub (validator staking + L1<->L2 anchoring, paper SIV-B/SIV-C)
# to L1, link the L2 Oracle to it (setL1Hub), and regenerate Go bindings.
deploy-l1hub: ensure-l1-mailbox
	$(ANNOUNCE_TARGET)
	@echo "Deploying L1Hub to L1 via forge..."
	@cd contracts && set -a && . ../.env && set +a && \
	  gas_price="$${L1_GAS_PRICE_WEI:-1000000000}"; \
	  forge script script/DeployL1Hub.s.sol --rpc-url "$$L1_RPC_URL" \
	    --private-key "$$ZKSYNC_PRIVATE_KEY" --legacy --gas-price "$$gas_price" --broadcast \
	    2>&1 | tee /tmp/forge_l1hub_deploy.log

	@addr=$$(scripts/deployed-address.sh DeployL1Hub.s.sol); \
	if [ -z "$$addr" ]; then \
		echo "ERROR: deployment succeeded but could not extract L1Hub address from forge output" >&2; \
		exit 1; \
	fi; \
	  echo "Detected L1Hub contract address: $$addr"; \
	if [ -f .env ]; then \
		awk -v addr="$$addr" 'BEGIN{updated=0} \
		/^L1_HUB_CONTRACT_ADDRESS=/{print "L1_HUB_CONTRACT_ADDRESS=" addr; updated=1; next} \
		{print} \
		END{if(!updated) print "L1_HUB_CONTRACT_ADDRESS=" addr}' .env > .env.tmp && mv .env.tmp .env; \
	else \
		echo "L1_HUB_CONTRACT_ADDRESS=$$addr" > .env; \
	fi; \
	  echo "Updated .env with L1_HUB_CONTRACT_ADDRESS=$$addr"

	@echo "Linking Oracle to L1Hub (setL1Hub)..."
	@set -a && . ./.env && set +a && \
	if [ -z "$$ORACLE_CONTRACT_ADDRESS" ]; then \
		echo "WARNING: ORACLE_CONTRACT_ADDRESS not set; run 'make deploy-oracle' first, then re-run 'make deploy-l1hub'." >&2; \
		exit 0; \
	fi; \
	  cast send "$$ORACLE_CONTRACT_ADDRESS" "setL1Hub(address)" "$$L1_HUB_CONTRACT_ADDRESS" \
	    --rpc-url "$$ZKSYNC_RPC_URL" --private-key "$$ZKSYNC_PRIVATE_KEY" >/dev/null && \
	  echo "Oracle.setL1Hub($$L1_HUB_CONTRACT_ADDRESS) done"

	@echo "Generating Go bindings via abigen..."
	@cd contracts && forge build --extra-output-files abi >/dev/null
	@abigen \
	  --abi contracts/out/L1Hub.sol/L1Hub.abi.json \
	  --pkg eth \
	  --type L1Hub \
	  --out go-server/internal/eth/l1hub_abigen.go
	@echo "Regenerated go-server/internal/eth/l1hub_abigen.go"

deploy-l1hub-if-configured:
	$(ANNOUNCE_TARGET)
	@if [ ! -f .env ]; then \
		echo "Skipping L1Hub deploy (.env not found)"; \
		exit 0; \
	fi
	@set -a && . ./.env && set +a && \
	if [ -z "$$L1_MAILBOX_ADDRESS" ]; then \
		echo "Skipping L1Hub deploy (L1_MAILBOX_ADDRESS not set)"; \
		exit 0; \
	fi; \
	  $(MAKE) deploy-l1hub

ensure-l1-mailbox:
	$(ANNOUNCE_TARGET)
	@bash -x ./scripts/ensure-l1-mailbox.sh
.PHONY: gnark-merkle-verifier deploy-merkle-verifier \
	gnark-votingbatch-verifier deploy-votingbatch-verifier
FORCE_ZK_KEYGEN ?= 1
gnark-merkle-verifier:
	$(ANNOUNCE_TARGET)
	@echo "Generating MerkleProof Solidity verifier via gnark..."
	@cd go-server && FORCE_ZK_KEYGEN=$(FORCE_ZK_KEYGEN) go run ./circuits/merkle_proof/cmd/gen_verifier/main.go

deploy-merkle-verifier: gnark-merkle-verifier
	$(ANNOUNCE_TARGET)
	@echo "Deploying MerkleProofVerifier via forge-zksync..."
	@cd contracts && set -a && . ../.env && set +a && \
	  forge script script/DeployMerkleProofVerifier.s.sol --zksync --rpc-url "$$ZKSYNC_RPC_URL" \
	    --private-key "$$ZKSYNC_PRIVATE_KEY" --suppress-warnings assemblycreate --broadcast \
	    2>&1 | tee /tmp/forge_merkle_deploy.log

	@addr=$$(scripts/deployed-address.sh DeployMerkleProofVerifier.s.sol); \
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
	$(ANNOUNCE_TARGET)
	@echo "Generating VotingBatch Solidity verifier via gnark..."
	@cd go-server && FORCE_ZK_KEYGEN=$(FORCE_ZK_KEYGEN) go run ./circuits/voting_batch/cmd/gen_votingbatch_verifier/main.go

deploy-votingbatch-verifier: gnark-votingbatch-verifier
	$(ANNOUNCE_TARGET)
	@echo "Deploying VotingBatchVerifier via forge-zksync..."
	@cd contracts && set -a && . ../.env && set +a && \
	  forge script script/DeployVotingBatchVerifier.s.sol \
	    --zksync --rpc-url "$$ZKSYNC_RPC_URL" --private-key "$$ZKSYNC_PRIVATE_KEY" \
	    --suppress-warnings assemblycreate --broadcast 2>&1 | tee /tmp/forge_votingbatch_deploy.log

	@addr=$$(scripts/deployed-address.sh DeployVotingBatchVerifier.s.sol); \
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
	$(ANNOUNCE_TARGET)
	@echo "Deploying Oracle via forge-zksync..."
	@cd contracts && set -a && . ../.env && set +a && \
	  forge script script/DeployOracle.s.sol --zksync --rpc-url "$$ZKSYNC_RPC_URL" \
	    --private-key "$$ZKSYNC_PRIVATE_KEY" --suppress-warnings assemblycreate --broadcast \
	    2>&1 | tee /tmp/forge_oracle_deploy.log

	@addr=$$(scripts/deployed-address.sh DeployOracle.s.sol); \
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
	$(ANNOUNCE_TARGET)
	@echo "Deploying MerkleTree via forge-zksync..."
	@cd contracts && set -a && . ../.env && set +a && \
	  forge script script/DeployMerkleTree.s.sol --zksync --rpc-url "$$ZKSYNC_RPC_URL" \
	    --private-key "$$ZKSYNC_PRIVATE_KEY" --suppress-warnings assemblycreate --broadcast \
	    2>&1 | tee /tmp/forge_merkle_tree_deploy.log

	@addr=$$(scripts/deployed-address.sh DeployMerkleTree.s.sol); \
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
	$(ANNOUNCE_TARGET)
	@TIMESTAMP=$$(date +"%Y%m%d_%H%M%S"); \
	zip -r archive_$$TIMESTAMP.zip . \
		-x ".git/*" \
		-x "go-server/circuits/build/*" \
		-x "*.csv" \
		-x "*.gob" \
		-x "*.log" \
		-x "*.zip"
