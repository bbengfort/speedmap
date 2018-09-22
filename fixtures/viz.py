#!/usr/bin/env python3
# Plot the results of a benchmark file

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
RESULTS = os.path.join(FIXTURES, "data", "results.csv")


def lineplot(x, y, hue, data, ax=None):
    if ax is None:
        ax = plt.gca()

    xr = np.arange(data[x].min(), data[x].max()+1)

    for h in data[hue].unique():
        means = data[data[hue]==h].groupby(x)[y].mean()
        std = data[data[hue]==h].groupby(x)[y].std()

        ax.plot(means, label=h)
        ax.fill_between(xr, means+std, means-std, alpha=0.25)

    ax.set_xlim(data[x].min(), data[x].max())
    ax.legend(frameon=True)

    return ax


def plot(path=RESULTS, outpath=None, bar=True):
    _, ax = plt.subplots(figsize=(9,6))

    df = pd.read_csv(path)
    workloads = df['workload'].unique()
    if len(workloads) > 1:
        raise ValueError("results set needs to be filtered by workload")

    if bar:
        sns.barplot(x="concurrency", y="throughput", hue="store", data=df, ax=ax)
    else:
        lineplot(x="concurrency", y="throughput", hue="store", data=df, ax=ax)

    ax.set_xlabel("number of concurrent clients")
    ax.set_ylabel("throughput (ops/sec)")
    ax.set_title("Concurrent Map Access for {} Workload".format(workloads[0].title()))

    if outpath is not None:
        plt.savefig(outpath)
    else:
        plt.show()


if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        description="plot the results of a benchmark file"
    )

    parser.add_argument(
        "-o", "--outpath", default=None, metavar="PATH",
        help="path to save the figure out to",
    )
    parser.add_argument(
        "-l", "--line", action="store_true", default=False,
        help="make a line plot instead of a bar plot",
    )
    parser.add_argument(
        "data", nargs="?", default=RESULTS,
        help="path to the results CSV file"
    )

    args = parser.parse_args()
    plot(args.data, args.outpath, bar=not args.line)
