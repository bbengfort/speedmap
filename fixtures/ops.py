#!/usr/bin/env python3
# Plot the results of a benchmark or measurement file

import os
import argparse
import numpy as np
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt

sns.set_style("whitegrid")
sns.set_context("notebook")

FIXTURES = os.path.dirname(__file__)
FIGURES = os.path.join(FIXTURES, "figures")
OPS = os.path.join(FIXTURES, "data", "ops.csv")
BLAST = os.path.join(FIXTURES, "data", "blast.csv")


def plot_ops(path=None, outpath=None):
    path = path or OPS
    outpath = outpath or os.path.join(FIGURES, "benchmark_operations.png")
    _, ax = plt.subplots(figsize=(9,6))

    df = pd.read_csv(path)
    sns.barplot(x='op', y='benchmark', hue='store', ax=ax, data=df)

    ax.set_xlabel("operation")
    ax.set_ylabel("ns/op")
    ax.set_title("Sequential Access Go Benchmark")

    plt.savefig(outpath)


def plot_blast(path=None, outpath=None):
    path = path or BLAST
    outpath = outpath or os.path.join(FIGURES, "benchmark_blast_throughput.png")
    _, ax = plt.subplots(figsize=(9,6))

    df = pd.read_csv(path)
    sns.barplot(x='op', y='throughput', hue='store', ax=ax, data=df)

    ax.set_xlabel('operation')
    ax.set_ylabel('throughput (ops/sec)')
    ax.set_title("Blast of 5000 Concurrent Accesses: Throughput")

    plt.savefig(outpath)


if __name__ == '__main__':
    parser = argparse.ArgumentParser()

    parser.add_argument(
        '-v', '--viz', choices=('ops', 'blast'), default='ops',
        help='specify the plot to draw from test data',
    )
    parser.add_argument(
        'data', nargs="?", default=None, help="path to CSV data"
    )
    parser.add_argument(
        "-o", "--outpath", default=None, help="path to store the file"
    )

    args = parser.parse_args()
    if args.viz == 'ops':
        plot_ops(args.data, args.outpath)
    elif args.viz == 'blast':
        plot_blast(args.data, args.outpath)
