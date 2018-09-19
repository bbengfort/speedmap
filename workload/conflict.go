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
		size:      DataSize,
	}
}

// Conflict allocates a key range to each thread, along with a probability
// that the thread will access a key being accessed by a different thread.
// A 0% probability means that the clients will access a disjoint key set.
type Conflict struct {
	readratio float32 // ratio of reads to writes
	prob      float32 // probability of conflict
	keys      int64   // the number of keys (each key identified by number)
	size      int     // the size of the value to write
}

// Run the conflict workload for the specified number of clients.
func (c *Conflict) Run(store speedmap.Store, clients int) (*speedmap.Result, error) {
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

// Runs the ith client in a go routine generating values as the byte string of
// the client number - operation number. Note that i must be 1-index to ensure
// that keyspace 0 is the conflict space.
func (c *Conflict) client(i int, store speedmap.Store, group *sync.WaitGroup) {
	r := int64(i)

	for o := 0; o < OpsPerThread; o++ {
		var key string
		val := []byte(fmt.Sprintf("%X-%X", r, o))

		if rand.Float32() < c.prob {
			// We have a conflict select any key in the key group
			// TODO: make sure own keyspace isn't selected
			key = RandomKey(0, c.keys)
		} else {
			key = RandomKey(r, c.keys)
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

// Runs the ith client in a go routine generating random values and mutating
// them according to the size of the writes. Note that i must be 1-index to
// ensure that keyspace 0 is the conflict space.
func (c *Conflict) complexClient(i int, store speedmap.Store, group *sync.WaitGroup) {
	r := int64(i)
	val, _ := GenerateRandomBytes(c.size)

	for o := 0; o < OpsPerThread; o++ {
		var key string
		RandomMutation(val, 8)

		if rand.Float32() < c.prob {
			// We have a conflict select any key in the shared key group
			key = RandomKey(0, c.keys)
		} else {
			// Select a key in the clients own keyspace
			key = RandomKey(r, c.keys)
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
	prob := c.prob * 100

	if c.readratio == 1.0 {
		return fmt.Sprintf("%0.0f%% conflict read-only", prob)
	}

	if c.readratio == 0.0 {
		return fmt.Sprintf("%0.0f%% conflict write-only", prob)
	}

	return fmt.Sprintf("%0.0f%% conflict %0.0f%% reads", prob, c.readratio*100)
}
