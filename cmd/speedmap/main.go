package main

import (
	"github.com/bbengfort/speedmap"
	"github.com/bbengfort/speedmap/store"
	"github.com/bbengfort/speedmap/workload"
)

func main() {
	workload := workload.NewConflict(0.0)
	bench := speedmap.New(workload, 24)

	basic := &store.Basic{}
	basic.Init()
	bench.Run(basic)

	msfr := &store.Misframe{}
	msfr.Init()
	bench.Run(msfr)

	smap := &store.SyncMap{}
	smap.Init()
	bench.Run(smap)

	bench.Save("fixtures/results.csv")

}
