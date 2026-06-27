#!/usr/bin/env bash
# docker compose up -d

# usage: ./start.sh INSTANCE_TYPE 
# Instance type is specifying the docker image to take:
# see https://hub.docker.com/r/matterlabs/local-node/tags for full list.
# latest2.0 - is the 'main' one.

INSTANCE_TYPE=${1:-latest2.0}

export INSTANCE_TYPE=$INSTANCE_TYPE

detect_platform() {
  case "$(uname -m)" in
    x86_64|amd64)
      echo "linux/amd64"
      ;;
    *)
      echo ""
      ;;
  esac
}

if [ -z "${L2A_PLATFORM:-}" ]; then
  L2A_PLATFORM=$(detect_platform)
  if [ -n "$L2A_PLATFORM" ]; then
    export L2A_PLATFORM
  fi
fi

docker compose up -d


check_all_services_healthy() {
  docker compose ps zksync 2>/dev/null | grep -q "(healthy)"
}

# True if the zksync container has actually failed (exited / crash-looping /
# marked unhealthy) rather than simply still starting up.
zksync_has_failed() {
  docker compose ps zksync 2>/dev/null | grep -qE "Exited|Restarting|unhealthy"
}

# NOTE: on Apple Silicon the zksync image only ships for linux/amd64, so it runs
# under emulation. First "healthy" state typically takes ~2-4 min (genesis +
# contract deploy), not seconds. The wait below is normal, not a hang.
echo "Waiting for zksync to become healthy (runs under amd64 emulation; expect ~2-4 min)..."

SECONDS=0
# Loop until all services are healthy
while ! check_all_services_healthy; do
  if zksync_has_failed; then
    echo "ERROR: zksync container exited / is restarting / is unhealthy after ${SECONDS}s." >&2
    echo "Recent logs:" >&2
    docker compose logs --tail=40 zksync >&2
    exit 1
  fi
  echo "Services are not yet healthy, waiting... (${SECONDS}s elapsed)"
  sleep 10  # Check every 10 seconds
done

echo "All services are healthy! (took ${SECONDS}s)"
