package speedmap

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Version of the speedmap package
const Version = "1.0"

// Store represents the interface for all in-memory key/value data structures
// that are being benchmarked by the Speed Map package.
type Store interface {
	Init() (err error)
	Get(key string) (value []byte, err error)
	Put(key string, value []byte) (err error)
	Delete(key string) (err error)
	GetOrCreate(key string, value []byte) (actual []byte, created bool)
	String() string
}

// Workload is an interface for creating a benchmark that runs operations
// against a data store for the specified number of clients and returns a
// result object with the number of successfully completed operations and
// the duration of the benchmark.
type Workload interface {
	Run(store Store, clients int) (*Result, error)
	String() string
}

// Result holds a record for a run of the workload execution.
type Result struct {
	Store       Store         // The store that the result is for
	Workload    Workload      // The workload of the result
	Concurrency int           // The number of concurrent clients executed on the store
	Operations  uint64        // The number of operations successfully executed
	Duration    time.Duration // The length of time the workload run took
}

// Throughput returns the number of operations per second achieved.
func (r *Result) Throughput() float64 {
	if r.Duration == 0 || r.Operations == 0 {
		return 0.0
	}

	return float64(r.Operations) / r.Duration.Seconds()
}

// String returns a CSV value for writing the record to disk:
// store,workload,concurrency,operations,duration (ns),throughput
func (r *Result) String() string {
	return fmt.Sprintf(
		"%s,%s,%d,%d,%d,%0.3f\n",
		r.Store,
		r.Workload,
		r.Concurrency,
		r.Operations,
		r.Duration,
		r.Throughput(),
	)
}
