package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/bbengfort/speedmap"
)

// Blast is a throughput measurement that launches n go routines performing
// the specified operation against the store and timing how long it takes for
// the n routines to complete, computing the throughput (ops/second).
func Blast(store speedmap.Store, n int, operation string) (*BlastResult, error) {
	if operation != "Get" && operation != "Put" && operation != "Delete" && operation != "GetOrCreate" {
		return nil, fmt.Errorf("unknown operation '%s'", operation)
	}

	wait := new(sync.WaitGroup)
	wait.Add(n)
	result := &BlastResult{Operations: uint64(n)}

	start := time.Now()
	for i := 0; i < n; i++ {
		go func(k int) {
			key := fmt.Sprintf("%X", k)
			switch operation {
			case "Get":
				store.Get(key)
			case "Put":
				store.Put(key, []byte(key))
			case "Delete":
				store.Delete(key)
			case "GetOrCreate":
				store.GetOrCreate(key, []byte(key))
			}

			wait.Done()
		}(i)
	}

	wait.Wait()
	result.Duration = time.Since(start)
	result.Throughput = float64(result.Operations) / result.Duration.Seconds()
	return result, nil
}

// BlastResult contains the throughput, number of operations, and the duration
// of the Blast for recording measurements.
type BlastResult struct {
	Operations uint64
	Duration   time.Duration
	Throughput float64
}
