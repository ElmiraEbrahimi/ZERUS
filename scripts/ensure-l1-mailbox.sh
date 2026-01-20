#!/usr/bin/env bash
set -euo pipefail

ENV_FILE="${1:-.env}"

if [[ ! -f "$ENV_FILE" ]]; then
  echo "ERROR: $ENV_FILE not found (copy .env.example to .env)" >&2
  exit 1
fi

if ! command -v curl >/dev/null 2>&1; then
  echo "ERROR: curl is required to resolve L1 mailbox address" >&2
  exit 1
fi

# shellcheck disable=SC1090
set -a
. "$ENV_FILE"
set +a

if [[ -z "${ZKSYNC_RPC_URL:-}" ]]; then
  echo "ERROR: ZKSYNC_RPC_URL is not set in $ENV_FILE" >&2
  exit 1
fi

rpc_call() {
  local method="$1"
  local url="$2"
  curl -sS --max-time 5 -H 'Content-Type: application/json' \
    -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"$method\",\"params\":[]}" \
    "$url"
}

extract_result() {
  grep -oE '0x[0-9a-fA-F]{40}' || true
}

extract_error() {
  sed -n 's/.*"message"[[:space:]]*:[[:space:]]*"\([^"]\+\)".*/\1/p'
}

extract_call_address() {
  local hex
  hex="$(sed -n 's/.*"result"[[:space:]]*:[[:space:]]*"\(0x[0-9a-fA-F]*\)".*/\1/p' | head -n1)"
  if [[ -z "$hex" || "$hex" == "0x" ]]; then
    return 1
  fi
  hex="${hex#0x}"
  if (( ${#hex} < 40 )); then
    return 1
  fi
  printf '0x%s\n' "${hex: -40}"
}

l1_eth_call() {
  local to="$1"
  local data="$2"
  curl -sS --max-time 5 -H 'Content-Type: application/json' \
    -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_call\",\"params\":[{\"to\":\"$to\",\"data\":\"$data\"},\"latest\"]}" \
    "$L1_RPC_URL"
}

add_candidate() {
  local url="$1"
  if [[ -z "$url" ]]; then
    return
  fi
  for existing in "${candidates[@]:-}"; do
    if [[ "$existing" == "$url" ]]; then
      return
    fi
  done
  candidates+=("$url")
}

candidates=()
add_candidate "$ZKSYNC_RPC_URL"

case "$ZKSYNC_RPC_URL" in
  ws://*)
    add_candidate "${ZKSYNC_RPC_URL/ws:\/\//http://}"
    ;;
  wss://*)
    add_candidate "${ZKSYNC_RPC_URL/wss:\/\//https://}"
    ;;
esac

if [[ "$ZKSYNC_RPC_URL" == *localhost:3051* ]]; then
  add_candidate "${ZKSYNC_RPC_URL/localhost:3051/localhost:3050}"
fi
if [[ "$ZKSYNC_RPC_URL" == *127.0.0.1:3051* ]]; then
  add_candidate "${ZKSYNC_RPC_URL/127.0.0.1:3051/127.0.0.1:3050}"
fi
if [[ "$ZKSYNC_RPC_URL" == *0.0.0.0:3051* ]]; then
  add_candidate "${ZKSYNC_RPC_URL/0.0.0.0:3051/0.0.0.0:3050}"
fi
if [[ "$ZKSYNC_RPC_URL" == *localhost:15101* ]]; then
  add_candidate "${ZKSYNC_RPC_URL/localhost:15101/localhost:15100}"
fi
if [[ "$ZKSYNC_RPC_URL" == *127.0.0.1:15101* ]]; then
  add_candidate "${ZKSYNC_RPC_URL/127.0.0.1:15101/127.0.0.1:15100}"
fi
if [[ "$ZKSYNC_RPC_URL" == *0.0.0.0:15101* ]]; then
  add_candidate "${ZKSYNC_RPC_URL/0.0.0.0:15101/0.0.0.0:15100}"
fi

address=""
use_direct=""
last_error=""
last_url=""
last_method=""
last_resolved_method=""
for url in "${candidates[@]}"; do
  for method in zks_getBridgehubAddress zks_getBridgehub zks_getMainContract; do
    resp="$(rpc_call "$method" "$url" || true)"
    addr="$(printf '%s' "$resp" | extract_result | head -n1)"
    if [[ -n "$addr" ]]; then
      address="$addr"
      last_resolved_method="$method"
      break 2
    fi
    err="$(printf '%s' "$resp" | extract_error | head -n1)"
    if [[ -n "$err" ]]; then
      last_error="$err"
      last_url="$url"
      last_method="$method"
    fi
  done
done

if [[ -z "$address" ]]; then
  if [[ -n "$last_error" ]]; then
    echo "ERROR: failed to resolve L1 mailbox via RPC ($last_method @ $last_url): $last_error" >&2
  else
    echo "ERROR: failed to resolve L1 mailbox via RPC" >&2
  fi
  exit 1
fi

if [[ "$last_resolved_method" == "zks_getBridgehubAddress" || "$last_resolved_method" == "zks_getBridgehub" ]]; then
  use_direct="true"
elif [[ "$last_resolved_method" == "zks_getMainContract" ]]; then
  use_direct="false"
fi

if [[ "$last_resolved_method" == "zks_getMainContract" && -n "${L1_RPC_URL:-}" ]]; then
  bridgehub_resp="$(l1_eth_call "$address" "0x3591c1a0" || true)"
  bridgehub_addr="$(printf '%s' "$bridgehub_resp" | extract_call_address)"
  if [[ -n "$bridgehub_addr" && "$bridgehub_addr" != "0x0000000000000000000000000000000000000000" ]]; then
    address="$bridgehub_addr"
    use_direct="true"
  fi
fi

# If the node only supports `zks_getMainContract`, we usually land on an Era-style
# single-chain devnet. In that case the correct mailbox to use is the main
# contract returned by `zks_getMainContract`, and the correct API is
# `proveL2LogInclusion(blockNumber, logIndex, ...)` (no chainId).
#
# Some zk-chains deployments do not expose bridgehub RPC helpers, but the main
# contract still exposes `getBridgehub()`. When available we treat it as the
# Bridgehub address and enable direct messaging.

tmp="$(mktemp "${ENV_FILE}.tmp.XXXXXX")"
awk -v addr="$address" -v direct="$use_direct" 'BEGIN{updated_addr=0; updated_direct=0} \
/^L1_MAILBOX_ADDRESS=/{print "L1_MAILBOX_ADDRESS=" addr; updated_addr=1; next} \
/^L1_USE_DIRECT_MESSAGING=/{if(direct != ""){print "L1_USE_DIRECT_MESSAGING=" direct; updated_direct=1; next}} \
{print} \
END{if(!updated_addr) print "L1_MAILBOX_ADDRESS=" addr; if(direct != "" && !updated_direct) print "L1_USE_DIRECT_MESSAGING=" direct}' "$ENV_FILE" > "$tmp"
mv "$tmp" "$ENV_FILE"
echo "Updated $ENV_FILE with L1_MAILBOX_ADDRESS=$address"
if [[ -n "$use_direct" ]]; then
  echo "Updated $ENV_FILE with L1_USE_DIRECT_MESSAGING=$use_direct"
fi
