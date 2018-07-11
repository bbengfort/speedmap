package speedmap

import (
	"errors"
	"os"
)

// New returns a Benchmark object ready to evaluate Stores. If maxthreads is
// less than one, sets the maximum number of threads to a default of 10.
func New(workload Workload, maxthreads int) *Benchmark {
	if maxthreads < 1 {
		maxthreads = 10
	}

	bench := &Benchmark{Workload: workload, MaxConcurrency: maxthreads}
	bench.Results = make([]*Result, 0)
	return bench
}

// Benchmark runs the workload against multiple stores with multiple clients
// and then saves the results as a CSV file to disk.
type Benchmark struct {
	Workload       Workload
	MaxConcurrency int
	Results        []*Result
}

// Run the benchmark against the specified Store.
func (b *Benchmark) Run(store Store) (err error) {
	for i := 1; i <= b.MaxConcurrency; i++ {
		var result *Result
		if result, err = b.Workload.Run(store, i); err != nil {
			return err
		}
		b.Results = append(b.Results, result)
	}
	return nil
}

// Save the benchmarks to disk.
func (b *Benchmark) Save(path string) (err error) {
	if len(b.Results) < 1 {
		return errors.New("no results to save")
	}

	var file *os.File
	if file, err = os.Create(path); err != nil {
		return err
	}
	defer file.Close()

	// Write the header of the CSV file.
	header := "store,workload,concurrency,operations,duration (ns),throughput\n"
	if _, err = file.Write([]byte(header)); err != nil {
		return err
	}

	// Write each of the result rows
	for _, result := range b.Results {
		if _, err = file.Write([]byte(result.String())); err != nil {
			return err
		}
	}

	return nil
}
