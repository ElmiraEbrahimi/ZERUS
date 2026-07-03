#!/usr/bin/env python3
"""Aggregate evaluation CSV logs into summary statistics (F-30).

Walks the sweep results directory produced by run-evaluation.sh, reads
every memtime n{N}_b{B}.csv, and emits per-(configuration, source, metric)
count/mean/std rows on stdout as CSV — the inputs from which the paper's
Tables I-II and Figures 5-8 are produced.
"""

import csv
import math
import re
import sys
from pathlib import Path

COMBO_RE = re.compile(r"n(\d+)_b(\d+)\.csv$")


def is_number(value: str) -> bool:
    try:
        float(value)
        return True
    except (TypeError, ValueError):
        return False


def aggregate(root: Path):
    rows = []
    for csv_path in sorted(root.rglob("n*_b*.csv")):
        match = COMBO_RE.search(csv_path.name)
        if not match:
            continue
        nodes, batch = int(match.group(1)), int(match.group(2))

        with csv_path.open(newline="") as fh:
            reader = csv.reader(fh)
            records = [r for r in reader if r]
        if len(records) < 2:
            continue
        header, body = records[0], records[1:]

        # Group samples by the `source` column when present.
        try:
            source_idx = [h.strip().lower() for h in header].index("source")
        except ValueError:
            source_idx = None

        groups = {}
        for record in body:
            key = record[source_idx] if source_idx is not None and source_idx < len(record) else "all"
            groups.setdefault(key, []).append(record)

        for source, records_for_source in sorted(groups.items()):
            for col, name in enumerate(header):
                values = [
                    float(r[col])
                    for r in records_for_source
                    if col < len(r) and is_number(r[col])
                ]
                # Skip non-numeric and index-like columns.
                if not values or name.strip().lower() in {"", "index", "run", "timestamp"}:
                    continue
                count = len(values)
                mean = sum(values) / count
                var = sum((v - mean) ** 2 for v in values) / count if count > 1 else 0.0
                rows.append(
                    (nodes, batch, source, name.strip(), count, mean, math.sqrt(var))
                )
    return rows


def main() -> int:
    if len(sys.argv) != 2:
        print(f"usage: {sys.argv[0]} <results-dir>", file=sys.stderr)
        return 1
    root = Path(sys.argv[1])
    if not root.is_dir():
        print(f"not a directory: {root}", file=sys.stderr)
        return 1

    writer = csv.writer(sys.stdout)
    writer.writerow(["node_count", "batch_size", "source", "metric", "samples", "mean", "std"])
    for row in aggregate(root):
        writer.writerow(row)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
