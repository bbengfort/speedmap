package main

import (
	"fmt"
	"os"

	"github.com/bbengfort/speedmap"
	"github.com/bbengfort/speedmap/store"
	"github.com/bbengfort/speedmap/workload"
	"github.com/urfave/cli"
)

func main() {

	// Instantiate the command line application
	app := cli.NewApp()
	app.Name = "speedmap"
	app.Version = speedmap.Version
	app.Usage = "benchmarks various concurrent map implementations"
	app.Action = bench
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "n, rounds",
			Usage: "number of benchmarking rounds",
			Value: 10,
		},
		cli.IntFlag{
			Name:  "t, threads",
			Usage: "maximum number of concurrent threads",
			Value: 10,
		},
		cli.StringFlag{
			Name:  "o, outpath",
			Usage: "path to write the results csv to",
		},
		cli.Float64Flag{
			Name:  "p, prob",
			Usage: "conflict probability in workload",
			Value: 0.5,
		},
		cli.Float64Flag{
			Name:  "r, readratio",
			Usage: "percent of reads in workload (0 for all writes)",
			Value: 0.5,
		},
		cli.BoolFlag{
			Name:  "B, no-basic",
			Usage: "exclude the basic store from evaluation",
		},
		cli.BoolFlag{
			Name:  "M, no-misframe",
			Usage: "exclude the misframe store from evaluation",
		},
		cli.BoolFlag{
			Name:  "S, no-sync",
			Usage: "exclude the sync map store from evaluation",
		},
	}

	// Run the CLI program
	app.Run(os.Args)

}

func bench(c *cli.Context) error {

	N := c.Int("rounds")
	T := c.Int("threads")

	stores := make([]speedmap.Store, 0, 4)
	workload := workload.NewConflict(float32(c.Float64("prob")), float32(c.Float64("readratio")))
	bench := speedmap.New(workload, T)

	if !c.Bool("no-basic") {
		basic := new(store.Basic)
		if err := basic.Init(); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		stores = append(stores, basic)
	}

	if !c.Bool("no-misframe") {
		msfr := new(store.Misframe)
		if err := msfr.Init(); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		stores = append(stores, msfr)
	}

	if !c.Bool("no-sync") {
		smap := new(store.SyncMap)
		if err := smap.Init(); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		stores = append(stores, smap)
	}

	rounds := N * T * len(stores)
	fmt.Printf("%s workload commencing for %d stores in %d rounds\n", workload, len(stores), rounds)

	for n := 0; n < N; n++ {
		for _, s := range stores {
			if err := bench.Run(s); err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			fmt.Print(".")
		}
	}

	if err := bench.Save(c.String("outpath")); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	fmt.Print("\n")
	return nil
}
