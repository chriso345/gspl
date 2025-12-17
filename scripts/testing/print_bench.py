#!/usr/bin/env python3
from __future__ import annotations
import csv
import sys
from collections import defaultdict
from pathlib import Path

BenchmarkData = dict[str, dict[str, dict[str, str]]]
MetricOrder = list[str]


def read_csv(file_path: Path) -> list[list[str]]:
    """Read CSV file and strip whitespace from each cell."""
    rows: list[list[str]] = []
    with file_path.open(newline="") as f:
        reader = csv.reader(f)
        for row in reader:
            rows.append([cell.strip() for cell in row])
    return rows


def split_blocks(rows: list[list[str]]) -> list[list[list[str]]]:
    """Split rows into blocks separated by empty rows."""
    blocks: list[list[list[str]]] = []
    current_block: list[list[str]] = []
    for row in rows:
        if not any(row):
            if current_block:
                blocks.append(current_block)
                current_block = []
        else:
            current_block.append(row)
    if current_block:
        blocks.append(current_block)
    return blocks


def collect_benchmarks(
    blocks: list[list[list[str]]],
    name1: str,
    name2: str,
) -> tuple[BenchmarkData, MetricOrder]:
    """Extract benchmark data and maintain metric order."""
    packages: BenchmarkData = defaultdict(lambda: defaultdict(dict))
    metrics_order: MetricOrder = []
    current_pkg: str | None = None

    for block in blocks:
        # Find package metadata
        for r in block:
            if r and r[0].startswith("pkg:"):
                current_pkg = r[0][len("pkg:") :].strip()
                break
        if not current_pkg:
            continue

        # Find metric row
        metric_row: list[str] | None = None
        for r in block:
            if len(r) > 1 and r[1] in ["sec/op", "B/op", "allocs/op"]:
                metric_row = r
                break
        if not metric_row:
            continue

        metric_name = metric_row[1]
        if metric_name not in metrics_order:
            metrics_order.append(metric_name)

        # Extract benchmark rows
        for r in block[2:]:
            if (
                not r
                or r[0] == "geomean"
                or r[0].startswith(("pkg:", "goos:", "goarch:", "cpu:"))
                or r[0] == ""
            ):
                continue
            testname = r[0]
            val1 = r[1] if len(r) > 1 else ""
            val2 = r[3] if len(r) > 3 else ""
            packages[current_pkg][testname][metric_name] = {  # pyright: ignore[reportArgumentType]
                name1: val1,
                name2: val2,
            }

    return packages, metrics_order


def calc_improvement(val1: str, val2: str) -> str:
    """Compute improvement percentage from val1 to val2."""
    try:
        v1 = float(val1)
        v2 = float(val2)
        if v1 == 0:
            return ""
        return f"{((v1 - v2) / v1 * 100):+.2f}%"
    except ValueError:
        return ""


def print_table(
    pkg_path: str,
    tests: dict[str, dict[str, dict[str, str]]],
    metrics_order: MetricOrder,
    name1: str,
    name2: str,
) -> None:
    """Print an ASCII table for a given package, with padded names."""
    if not tests:
        return

    # Determine max length of the name labels
    name_width = max(len(name1), len(name2))

    # Compute first column width, including test names
    first_col_width = max(
        [len(f"{name1.ljust(name_width)} {t}") for t in tests]
        + [len(f"{name2.ljust(name_width)} {t}") for t in tests]
        + [len(pkg_path)]
    )
    col_widths = [first_col_width]

    # Compute widths for metric columns
    for m in metrics_order:
        w = max(
            len(m),
            *(len(tests[t][m][name1]) for t in tests),
            *(len(tests[t][m][name2]) for t in tests),
            *(
                len(calc_improvement(tests[t][m][name1], tests[t][m][name2]))
                for t in tests
            ),
        )
        col_widths.append(w)

    total_width = sum(col_widths) + 3 * len(col_widths) + 1
    hline = "+" + "-" * (total_width - 2) + "+"
    print(hline)

    # Package header
    pkg_header = "| " + pkg_path.ljust(col_widths[0])
    for i, m in enumerate(metrics_order):
        pkg_header += " | " + m.ljust(col_widths[i + 1])
    pkg_header += " |"
    print(pkg_header)
    print(hline)

    # Print rows per test
    for t in tests.keys():
        row1 = [f"{name1.ljust(name_width)} {t}"] + [
            tests[t][m][name1] for m in metrics_order
        ]
        row2 = [f"{name2.ljust(name_width)} {t}"] + [
            tests[t][m][name2] for m in metrics_order
        ]
        improvement_row = ["".ljust(first_col_width)] + [
            calc_improvement(tests[t][m][name1], tests[t][m][name2])
            for m in metrics_order
        ]

        for row in [row1, row2, improvement_row]:
            line = (
                "| "
                + " | ".join(str(row[i]).ljust(col_widths[i]) for i in range(len(row)))
                + " |"
            )
            print(line)
        print(hline)
    print()


def main() -> None:
    if len(sys.argv) < 4:
        print("Usage: script.py <csv_file> <name1> <name2>")
        sys.exit(1)

    csv_file = Path(sys.argv[1])
    name1 = sys.argv[2]
    name2 = sys.argv[3]

    if not csv_file.exists():
        print(f"File not found: {csv_file}")
        sys.exit(1)

    rows = read_csv(csv_file)
    if not rows:
        print("No data found in CSV")
        sys.exit(1)

    blocks = split_blocks(rows)
    packages, metrics_order = collect_benchmarks(blocks, name1, name2)

    for pkg_path, tests in packages.items():
        print_table(pkg_path, tests, metrics_order, name1, name2)  # pyright: ignore[reportArgumentType]


if __name__ == "__main__":
    main()
