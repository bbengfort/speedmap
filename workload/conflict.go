package workload

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bbengfort/speedmap"
)

// NewConflict workload with the specified probability.
func NewConflict(prob, readratio float32) *Conflict {
	return &Conflict{
		readratio: readratio,
		prob:      prob,
		keys:      MaxKeys,
	}
}

// Conflict allocates a key range to each thread, along with a probability
// that the thread will access a key being accessed by a different thread.
// A 0% probability means that the clients will access a disjoint key set.
type Conflict struct {
	readratio float32 // ratio of reads to writes
	prob      float32 // probability of conflict
	keys      int64   // the number of keys (each key identified by number)
	clients   int     // the number of clients
}

// Run the conflict workload for the specified number of clients.
func (c *Conflict) Run(store speedmap.Store, clients int) (*speedmap.Result, error) {
	c.clients = clients
	result := &speedmap.Result{Store: store, Workload: c, Concurrency: clients}

	group := &sync.WaitGroup{}
	group.Add(clients)

	start := time.Now()
	for i := 1; i <= clients; i++ {
		go c.client(i, store, group)
	}
	group.Wait()
	result.Duration = time.Since(start)
	result.Operations = uint64(clients) * uint64(OpsPerThread)

	return result, nil
}

// Runs the ith client in a go routine.
func (c *Conflict) client(i int, store speedmap.Store, group *sync.WaitGroup) {
	r := int64(i)
	for o := 0; o < OpsPerThread; o++ {
		var key string
		var val = []byte(fmt.Sprintf("the time is now %s", time.Now()))

		if rand.Float32() < c.prob {
			// We have a conflict select any key in the key group
			// TODO: make sure own keyspace isn't selected
			key = fmt.Sprintf("%X", rand.Int63n(c.keys))
		} else {
			key = fmt.Sprintf("%X", rand.Int63n(c.keys)*r)
		}

		if rand.Float32() <= c.readratio {
			// GetOrCreate a key
			store.GetOrCreate(key, val)
		} else {
			// Put a key
			store.Put(key, val)
		}
	}
	group.Done()
}

// String returns a representation of the conflict workload
func (c *Conflict) String() string {
	if c.readratio == 1.0 {
		return fmt.Sprintf("%0.2f%% conflict read-only", c.prob)
	}

	if c.readratio == 0.0 {
		return fmt.Sprintf("%0.2f%% conflict write-only", c.prob)
	}

	return fmt.Sprintf("%0.2f conflict %0.2f reads", c.prob, c.readratio)
}
