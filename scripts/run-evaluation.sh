#!/bin/bash
# Evaluation runner (F-30): reproduces the paper's SVI sweeps by driving
# simulate.sh over committee sizes and batch sizes with repeated trials,
# then aggregating the per-(n, b) CSV logs into summary statistics.
#
# Usage:
#   ./scripts/run-evaluation.sh [-r <runs>] [-n "<node counts>"] [-b "<batch sizes>"]
#
# Defaults mirror the paper: 50 trials per configuration,
# n in {4 8 16 32 64 128}, b in {1 5 10 15} (Tables I-II, Figures 5-8).
# Requires the local stack up and contracts deployed (make up-deploy).
set -euo pipefail

cd "$(dirname "$0")/.."

runs=50
node_counts="4 8 16 32 64 128"
batch_sizes="1 5 10 15"

while getopts ":r:n:b:" opt; do
	case "$opt" in
	r) runs="$OPTARG" ;;
	n) node_counts="$OPTARG" ;;
	b) batch_sizes="$OPTARG" ;;
	*)
		echo "Usage: $0 [-r <runs>] [-n \"<node counts>\"] [-b \"<batch sizes>\"]" >&2
		exit 1
		;;
	esac
done

results_dir="eval-results/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$results_dir"
echo "Evaluation sweep: runs=$runs nodes={$node_counts} batches={$batch_sizes}"
echo "Results directory: $results_dir"

for n in $node_counts; do
	for b in $batch_sizes; do
		echo "=== sweep: NODE_COUNT=$n BATCH_SIZE=$b runs=$runs ==="
		./simulate.sh -n "$n" -b "$b" -r "$runs"

		combo_dir="$results_dir/n${n}_b${b}"
		mkdir -p "$combo_dir"
		if [ -d go-server/logs ]; then
			cp -R go-server/logs/. "$combo_dir/"
		fi
	done
done

echo "Aggregating summary statistics..."
python3 scripts/aggregate-results.py "$results_dir" >"$results_dir/summary.csv"
echo "Summary written to $results_dir/summary.csv"
