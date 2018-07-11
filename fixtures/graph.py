#!/usr/bin/env python3

import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt

sns.set_style("whitegrid")
sns.set_context("notebook")


def plot(path="fixtures/results.csv"):
    df = pd.read_csv(path)

    g = sns.barplot(x="concurrency", y="throughput", hue="store", data=df)
    plt.show()


if __name__ == '__main__':
    plot()
