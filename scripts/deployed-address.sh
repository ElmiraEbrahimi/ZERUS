#!/usr/bin/env bash
# Resolve the address of the contract actually deployed by a forge script run.
#
# forge's console output prints addresses computed during script SIMULATION.
# On zkSync (foundry-zksync), scripts that implicitly deploy linked libraries
# (e.g. MiMC for Oracle/MerkleTree) broadcast an extra library-deployment
# transaction, shifting the deployment nonce, so the printed address can
# diverge from the address actually created on-chain. The broadcast journal
# is the ground truth:
#   - zkSync: receipts carry ContractDeployed events emitted by the
#     ContractDeployer system contract (0x...8006); topics[3] is the created
#     address. Libraries deploy before the script's main contract, so the
#     last event is the contract we want.
#   - plain EVM: the journal records CREATE transactions with their
#     contractAddress.
#
# Usage: deployed-address.sh <ScriptName.s.sol> [broadcast-root]
set -euo pipefail

script_name="${1:?usage: deployed-address.sh <ScriptName.s.sol> [broadcast-root]}"
root="${2:-contracts/broadcast}"

python3 - "$script_name" "$root" <<'PY'
import glob
import json
import sys

script, root = sys.argv[1], sys.argv[2]
runs = glob.glob(f"{root}/{script}/*/run-latest.json")
if not runs:
    sys.exit(f"deployed-address: no broadcast journal under {root}/{script}/")


def load(path):
    with open(path) as fh:
        return json.load(fh)


# Newest run across chain-id directories.
path = max(runs, key=lambda p: load(p).get("timestamp", 0))
run = load(path)

CONTRACT_DEPLOYER = "0x0000000000000000000000000000000000008006"
# keccak256("ContractDeployed(address,bytes32,address)")
CONTRACT_DEPLOYED_TOPIC = (
    "0x290afdae231a3fc0bbae8b1af63698b0a1d79b21ad17df0342dfb952fe74f8e5"
)

deployed = []
for receipt in run.get("receipts", []):
    for log in receipt.get("logs", []):
        topics = log.get("topics") or []
        if (
            log.get("address", "").lower() == CONTRACT_DEPLOYER
            and topics
            and topics[0].lower() == CONTRACT_DEPLOYED_TOPIC
            and len(topics) >= 4
        ):
            deployed.append("0x" + topics[3][-40:])

if not deployed:
    # Plain-EVM script: use the recorded CREATE transactions.
    for tx in run.get("transactions", []):
        if tx.get("transactionType") == "CREATE" and tx.get("contractAddress"):
            deployed.append(tx["contractAddress"])

if not deployed:
    sys.exit(f"deployed-address: no deployed contract recorded in {path}")

print(deployed[-1])
PY
