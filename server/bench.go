package server

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bbengfort/speedmap"
)

// Blast implements Benchmark by sending n Put requests to the specified server
// each in its own thread. It then records the total time it takes to complete
// all n requests and uses that to compute the throughput. Additionally, each
// thread records the latency of each request, so that outlier requests can
// be removed from the blast computation.
//
// Note: this benchmark Puts a unique key and short value to the server, its
// intent is to compute pedal to the metal write throughput.
type Blast struct {
	requests  uint64          // the number of successful requests
	failures  uint64          // the number of failed requests
	started   time.Time       // the time the benchmark was started
	duration  time.Duration   // the duration of the benchmark period
	latencies []time.Duration // observed latencies in the number of requests
}

// Run the blast benchmark against the system by putting a unique key and
// small value to the server as fast as possible and measuring the duration.
func (b *Blast) Run(addr string, N, S uint) (err error) {
	// Initialize the blast latencies and results (resetting if rerun)
	b.requests = 0
	b.failures = 0
	b.latencies = make([]time.Duration, N)
	results := make([]error, N)

	// Initialize the keys and values so that it's not part of throughput.
	keys := make([]string, N)
	vals := make([][]byte, N)

	for i := uint(0); i < N; i++ {
		keys[i] = fmt.Sprintf("%X", i)
		vals[i] = make([]byte, S)
		rand.Read(vals[i])
	}

	// Create the wait group for all threads
	group := new(sync.WaitGroup)
	group.Add(int(N))

	// Initialize a client per operation in an array and connect to server.
	// NOTE: done separately from the key/value loop to ensure connections are
	// open for as little as time as possible.
	// clients := make([]*Client, N)
	// for i := uint(0); i < N; i++ {
	// 	clients[i] = NewClient(keys[i])
	// 	if err = clients[i].Connect(addr); err != nil {
	// 		return fmt.Errorf("could not create client %d: %s", i, err)
	// 	}
	// }
	clients := NewClient("")
	if err = clients.Connect(addr); err != nil {
		return err
	}

	// Execute the blast operation against the server
	b.started = time.Now()
	for i := uint(0); i < N; i++ {
		go func(k uint) {
			// Make Put request and if there is no error, store true!
			start := time.Now()
			// _, results[k] = clients[k].Put(keys[k], vals[k])
			_, results[k] = clients.Put(keys[k], vals[k])

			// Record the latency of the result, success or failure
			b.latencies[k] = time.Since(start)
			group.Done()
		}(i)
	}

	group.Wait()
	b.duration = time.Since(b.started)

	// Compute successes and failures
	errs := make(map[string]uint64)
	for _, r := range results {
		if r == nil {
			b.requests++
		} else {
			b.failures++
			errs[r.Error()]++
		}
	}

	// Print any errors that occurred
	for e, c := range errs {
		fmt.Printf("%d: %s\n", c, e)
	}

	return nil
}

// Complete returns true if requests and duration is greater than 0.
func (b *Blast) Complete() bool {
	return b.requests > 0 && b.duration > 0
}

// Throughput computes the number of requests (excluding failures) by the
// total duration of the experiment, e.g. the operations per second.
func (b *Blast) Throughput() float64 {
	if b.duration == 0 {
		return 0.0
	}

	return float64(b.requests) / b.duration.Seconds()
}

// CSV returns a results row delimited by commas as:
//     requests,failures,duration,throughput,version,benchmark
// If header is specified then string contains two rows with the header first.
func (b *Blast) CSV(header bool) (string, error) {
	if !b.Complete() {
		return "", errors.New("benchmark has not been run yet")
	}

	row := fmt.Sprintf(
		"%d,%d,%s,%0.4f,%s,blast",
		b.requests, b.failures, b.duration, b.Throughput(), speedmap.Version,
	)

	if header {
		return fmt.Sprintf("requests,failures,duration,throughput,version,benchmark\n%s", row), nil
	}

	return row, nil
}

// JSON returns a results row as a json object, formatted with or without the
// number of spaces specified by indent. Use no indent for JSON lines format.
func (b *Blast) JSON(indent int) ([]byte, error) {
	data := b.serialize()

	if indent > 0 {
		indent := strings.Repeat(" ", indent)
		return json.MarshalIndent(data, "", indent)
	}

	return json.Marshal(data)
}

// serialize converts the benchmark into a map[string]interface{} -- useful
// for dumping the benchmark as JSON and used from structs that embed benchmark
// to include more data in the results.
func (b *Blast) serialize() map[string]interface{} {
	data := make(map[string]interface{})

	data["requests"] = b.requests
	data["failures"] = b.failures
	data["duration"] = b.duration.String()
	data["throughput"] = b.Throughput()
	data["version"] = speedmap.Version
	data["benchmark"] = "multiclient"

	return data
}
