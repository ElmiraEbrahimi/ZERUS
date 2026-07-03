#!/bin/bash

set -euo pipefail

server_pid=""
server_pgid=""
logger_pid=""
server_tmp=""
server_port=18001
sleep_base=3
script_pgid="$(ps -o pgid= $$ | tr -d ' ')"

calc_sleep_with_nodes() {
	awk -v base="$sleep_base" -v nodes="${node_count:-0}" 'BEGIN { printf "%.1f", base + (nodes * 0.5) }'
}

kill_tree() {
	local pid="${1:-}"
	local sig="${2:-TERM}"
	if [ -z "$pid" ]; then
		return 0
	fi
	local children
	if command -v pgrep >/dev/null 2>&1; then
		children="$(pgrep -P "$pid" 2>/dev/null || true)"
	else
		children="$(ps -o pid= -P "$pid" 2>/dev/null | tr -d ' ' || true)"
	fi
	for child in $children; do
		kill_tree "$child" "$sig"
	done
	kill "-$sig" "$pid" 2>/dev/null || true
}

wait_for_port_free() {
	if ! command -v lsof >/dev/null 2>&1; then
		sleep 0.5
		return 0
	fi
	local i
	for ((i = 0; i < 50; i++)); do
		if ! lsof -tiTCP:"$server_port" -sTCP:LISTEN >/dev/null 2>&1; then
			return 0
		fi
		sleep 0.1
	done
	return 1
}

force_kill_port_listener() {
	if ! command -v lsof >/dev/null 2>&1; then
		return 0
	fi
	local pids
	pids="$(lsof -tiTCP:"$server_port" -sTCP:LISTEN 2>/dev/null || true)"
	if [ -z "$pids" ]; then
		return 0
	fi
	kill -TERM $pids 2>/dev/null || true
	sleep 0.2
	pids="$(lsof -tiTCP:"$server_port" -sTCP:LISTEN 2>/dev/null || true)"
	if [ -n "$pids" ]; then
		kill -KILL $pids 2>/dev/null || true
	fi
	wait_for_port_free || true
}

cleanup_server() {
	local had_server=""
	if [ -n "${server_pid:-}" ]; then
		had_server=1
	fi
	if [ -n "${server_pgid:-}" ] && [ "$server_pgid" != "$script_pgid" ]; then
		kill -TERM "-$server_pgid" 2>/dev/null || true
	elif [ -n "${server_pid:-}" ]; then
		kill_tree "$server_pid" TERM
	fi
	if [ -n "${server_pid:-}" ]; then
		wait "$server_pid" 2>/dev/null || true
	fi
	if [ -n "${server_pid:-}" ] && kill -0 "$server_pid" 2>/dev/null; then
		if [ -n "${server_pgid:-}" ] && [ "$server_pgid" != "$script_pgid" ]; then
			kill -KILL "-$server_pgid" 2>/dev/null || true
		else
			kill_tree "$server_pid" KILL
		fi
		wait "$server_pid" 2>/dev/null || true
	fi
	if [ -n "${logger_pid:-}" ]; then
		kill "$logger_pid" 2>/dev/null || true
		# Reap the logger job quietly so bash does not print a
		# "Terminated" job-status notice for it.
		wait "$logger_pid" 2>/dev/null || true
	fi
	if [ -n "${server_tmp:-}" ]; then
		rm -rf "$server_tmp"
	fi
	if ! wait_for_port_free; then
		force_kill_port_listener
		echo "Port $server_port still in use after cleanup"
	fi
	server_pid=""
	server_pgid=""
	logger_pid=""
	server_tmp=""
	if [ -n "$had_server" ]; then
		echo "✅ Server stopped cleanly."
	fi
}

trap cleanup_server EXIT

usage() {
	echo "Usage: $0 [-r <runs>] [-n <node-count>] [-b <batch-size>] [-s <sleep-seconds>]"
	echo "       $0 [--runs <runs>]"
	echo "       If --runs is not provided, it defaults to 50."
}

update_env_var() {
	local key="$1"
	local value="$2"
	local file="$3"
	local tmp
	tmp="$(mktemp)"
	awk -v key="$key" -v value="$value" '
		BEGIN { found = 0 }
		$0 ~ "^" key "=" {
			print key "=" value
			found = 1
			next
		}
		{ print }
		END {
			if (found == 0) {
				print key "=" value
			}
		}
	' "$file" >"$tmp"
	mv "$tmp" "$file"
}

get_env_var() {
	local key="$1"
	local file="$2"
	awk -F= -v key="$key" '$0 ~ "^" key "=" {print $2}' "$file" | tail -n 1
}

reset_run_index() {
	local memtime_dir
	if [ -d "$(pwd)/go-server" ]; then
		memtime_dir="$(pwd)/go-server/logs/memtime"
	else
		memtime_dir="$(pwd)/logs/memtime"
	fi
	local run_index_file
	run_index_file="$(dirname "$memtime_dir")/run_index"
	mkdir -p "$(dirname "$run_index_file")"
	: >"$run_index_file"
}

runs="50"
runs_set=0
node_count_arg=""
batch_size_arg=""
node_count_set=0
batch_size_set=0
sleep_base_set=0

OPTIND=1
while getopts ":r:n:b:s:-:" opt; do
	case "$opt" in
	r)
		runs="$OPTARG"
		runs_set=1
		;;
	n)
		node_count_arg="$OPTARG"
		node_count_set=1
		;;
	b)
		batch_size_arg="$OPTARG"
		batch_size_set=1
		;;
	s)
		sleep_base="$OPTARG"
		sleep_base_set=1
		;;
	-)
		case "$OPTARG" in
		runs)
			runs="${!OPTIND}"
			OPTIND=$((OPTIND + 1))
			runs_set=1
			;;
		*)
			echo "Unknown option --$OPTARG"
			usage
			exit 1
			;;
		esac
		;;
	\?)
		echo "Unknown option -$OPTARG"
		usage
		exit 1
		;;
	:)
		echo "Option -$OPTARG requires an argument"
		usage
		exit 1
		;;
	esac
done

if ! [[ "$runs" =~ ^[0-9]+$ ]]; then
	echo "Runs must be an integer: $runs"
	exit 1
fi
if [ "$runs" -le 0 ]; then
	echo "Runs must be greater than zero: $runs"
	exit 1
fi
runs_limit="$runs"
readonly runs_limit

if [ "$runs_set" -eq 1 ]; then
	echo "Runs: $runs"
fi

if [ "$sleep_base_set" -eq 1 ]; then
	if ! [[ "$sleep_base" =~ ^[0-9]+$ ]]; then
		echo "Sleep seconds must be an integer: $sleep_base"
		exit 1
	fi
	echo "Sleep base: $sleep_base"
fi

if [ "$node_count_set" -eq 1 ]; then
	if ! [[ "$node_count_arg" =~ ^[0-9]+$ ]]; then
		echo "Node count must be an integer: $node_count_arg"
		exit 1
	fi
	if [ "$node_count_arg" -le 0 ]; then
		echo "Node count must be greater than zero: $node_count_arg"
		exit 1
	fi
	tmp_node="$node_count_arg"
	sparse_tree_depth=0
	while [ "$tmp_node" -gt 1 ]; do
		if ((tmp_node % 2 != 0)); then
			echo "Node count must be a power of two to derive SPARSE_TREE_DEPTH: $node_count_arg"
			exit 1
		fi
		tmp_node=$((tmp_node / 2))
		sparse_tree_depth=$((sparse_tree_depth + 1))
	done
fi

if [ "$batch_size_set" -eq 1 ]; then
	if ! [[ "$batch_size_arg" =~ ^[0-9]+$ ]]; then
		echo "Batch size must be an integer: $batch_size_arg"
		exit 1
	fi
fi

env_file=".env"
if [ ! -f "$env_file" ]; then
	echo "Missing .env in $(pwd)"
	exit 1
fi

updates=()
if [ "$node_count_set" -eq 1 ]; then
	current_node_count="$(get_env_var "NODE_COUNT" "$env_file")"
	current_sparse_tree_depth="$(get_env_var "SPARSE_TREE_DEPTH" "$env_file")"
	if [ "$current_node_count" != "$node_count_arg" ]; then
		update_env_var "NODE_COUNT" "$node_count_arg" "$env_file"
		updates+=("NODE_COUNT=$node_count_arg")
	fi
	if [ "$current_sparse_tree_depth" != "$sparse_tree_depth" ]; then
		update_env_var "SPARSE_TREE_DEPTH" "$sparse_tree_depth" "$env_file"
		updates+=("SPARSE_TREE_DEPTH=$sparse_tree_depth")
	fi
fi
if [ "$batch_size_set" -eq 1 ]; then
	current_batch_size="$(get_env_var "BATCH_SIZE" "$env_file")"
	if [ "$current_batch_size" != "$batch_size_arg" ]; then
		update_env_var "BATCH_SIZE" "$batch_size_arg" "$env_file"
		updates+=("BATCH_SIZE=$batch_size_arg")
	fi
fi

if [ "${#updates[@]}" -gt 0 ]; then
	echo ".env updated: ${updates[*]}"
else
	echo ".env not updated"
fi

reset_run_index

run() {
	local counter="${1:-}"
	if [ -z "$counter" ]; then
		echo "Missing run counter"
		return 1
	fi
	force_kill_port_listener
	if ! wait_for_port_free; then
		echo "Port $server_port is still in use; aborting run"
		return 1
	fi
	make deploy
	server_tmp="$(mktemp -d)"
	local server_log="$server_tmp/server.log"
	local server_ready="$server_tmp/server.ready"
	local server_fifo="$server_tmp/server.fifo"
	mkfifo "$server_fifo"
	if command -v setsid >/dev/null 2>&1; then
		setsid make server >"$server_fifo" 2>&1 &
	else
		make server >"$server_fifo" 2>&1 &
	fi
	server_pid=$!
	server_pgid="$(ps -o pgid= "$server_pid" | tr -d ' ')"
	if [ -z "$server_pgid" ]; then
		server_pgid="$server_pid"
	fi
	{
		while IFS= read -r line; do
			# Expected shutdown noise (SIGTERM during cleanup) stays in the
			# log file but is kept off the console.
			case "$line" in
			"make: *** [server]"*[Tt]erminated*) ;;
			*"signal: terminated"*) ;;
			*) printf '%s\n' "$line" ;;
			esac
			printf '%s\n' "$line" >>"$server_log"
			if [[ "$line" == *"starting server on"* ]]; then
				: >"$server_ready"
			fi
		done <"$server_fifo"
	} &
	logger_pid=$!
	while [ ! -f "$server_ready" ]; do
		if ! kill -0 "$server_pid" 2>/dev/null; then
			echo "make server exited before startup"
			cat "$server_log"
			cleanup_server
			return 1
		fi
		sleep 0.2
	done
	echo "Running $counter: NODE_COUNT=$node_count BATCH_SIZE=$batch_size"

	local sleep_with_nodes
	sleep_with_nodes="$(calc_sleep_with_nodes)"
	curl -X POST "http://localhost:${server_port}/validators/register" -H "Content-Type: application/json"
	sleep "$sleep_with_nodes"
	curl -X POST "http://localhost:${server_port}/users/default/register" -H "Content-Type: application/json"
	sleep 1
	for ((bs = 0; bs < batch_size; bs++)); do
		curl -X POST "http://localhost:${server_port}/users/default/burn" -H "Content-Type: application/json"
		sleep "$sleep_base"
	done
	sleep "$sleep_with_nodes"
	for ((bs = 0; bs < batch_size; bs++)); do
		curl -X POST "http://localhost:${server_port}/users/default/withdraw" -H "Content-Type: application/json"
		sleep "$sleep_base"
	done
	sleep "$sleep_with_nodes"
	curl -X POST "http://localhost:${server_port}/validators/replace" -H "Content-Type: application/json" -d '{"node_id":0,"replace_with_account_id":1}'
	curl -X POST "http://localhost:${server_port}/validators/exit" -H "Content-Type: application/json" -d '{"node_id":0}'
	curl -X POST "http://localhost:${server_port}/validators/withdraw" -H "Content-Type: application/json" -d '{"node_id":0}'
	echo "\n"

	cleanup_server
}

env_file=".env"

http_bind_addr=$(awk -F= '/^HTTP_BIND_ADDR=/{print $2}' "$env_file" | tail -n 1)
http_bind_addr="${http_bind_addr%\"}"
http_bind_addr="${http_bind_addr#\"}"
http_bind_port="${http_bind_addr##*:}"
if [[ "$http_bind_port" =~ ^[0-9]+$ ]]; then
	server_port="$http_bind_port"
fi

node_count=$(awk -F= '/^NODE_COUNT=/{print $2}' "$env_file" | tail -n 1)
batch_size=$(awk -F= '/^BATCH_SIZE=/{print $2}' "$env_file" | tail -n 1)

for ((i = 1; i <= runs_limit; i++)); do
	run "$i"
done
