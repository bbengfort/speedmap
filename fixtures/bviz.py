#!/usr/bin/env python3
# Visualize the blast throughput

import json
import argparse
import numpy as np
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt


def load_data(path):
    with open(path, 'r') as f:
        for line in f:
            if not line.startswith("{"):
                continue
            yield json.loads(line)


def p95(v):
    return np.percentile(v, 95)


def plot_throughput(path, failure=False, title=None):
    if failure:
        _, axes = plt.subplots(nrows=2, figsize=(10,6), sharex=True)
    else:
        _, ax = plt.subplots(figsize=(10,6))
        axes = [ax]

    data = pd.DataFrame(load_data(path))
    data["ops"] = data["requests"] + data["failures"]

    # Plot throughput

    g = sns.barplot(x="ops", y="throughput", data=data, ci=None, estimator=p95, color="#2980b9", ax=axes[0])
    g.set_ylabel("throughput (ops/sec)")
    if failure:
        g.set_xlabel("")
    else:
        g.set_xlabel("number of requests")

    title = title or "Blast Throughput: Sync Map, Macbook Pro Local"
    g.set_title(title)

    if failure:
        g = sns.barplot(x="ops", y="failures", data=data, color="#e74c3c", ax=axes[1])
        g.set_xlabel("number of requests")
        g.set_ylabel("failures")

    return axes


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('data', help='json lines throughput generated by blast.sh')
    parser.add_argument('-t', '--title', default=None, help='title to give to the figure')
    parser.add_argument('-f', '--failure', action='store_true', help='plot failures as well')
    parser.add_argument('-s', '--savefig', default=None, help='path to save the figure')

    args = parser.parse_args()
    g = plot_throughput(args.data, args.failure, args.title)
    plt.tight_layout()

    if args.savefig:
        plt.savefig(args.savefig)
    else:
        plt.show()