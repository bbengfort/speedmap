package store_test

import (
	"fmt"
	"testing"

	"github.com/bbengfort/speedmap"
	. "github.com/bbengfort/speedmap/store"
)

// Create a list of stores for benchmarking
func makeStores(t testing.TB) []speedmap.Store {
	var (
		e error
		s speedmap.Store
	)

	// Create the stores array
	stores := make([]speedmap.Store, 0, 4)

	// Add the basic store
	if s, e = NewBasic(); e != nil {
		t.Fatalf("could not create basic store: %s", e)
	}
	stores = append(stores, s)

	// Add the misframe store
	if s, e = NewMisframe(); e != nil {
		t.Fatalf("could not create misframe store: %s", e)
	}
	stores = append(stores, s)

	// Add the sync map store
	if s, e = NewSyncMap(); e != nil {
		t.Fatalf("could not create sync map store: %s", e)
	}
	stores = append(stores, s)

	// Add the shard store
	if s, e = NewShard(); e != nil {
		t.Fatalf("could not create shard store: %s", e)
	}
	stores = append(stores, s)

	return stores
}

func makeKey(i int) string {
	return fmt.Sprintf("%X", i)
}

// Benchmark the Get operation on the available stores; this measures the
// number of nanoseconds per operation, not necessarily the concurrent access
// throughput as defined by the Measurement tests with each store.
func BenchmarkGet(b *testing.B) {
	for _, store := range makeStores(b) {
		b.Run(store.String(), func(b *testing.B) {
			// Allocate a bunch of keys into the store
			for k := 0; k < 256; k++ {
				key := makeKey(k)
				store.Put(key, []byte(key))
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := makeKey(i % 256)
				if _, err := store.Get(key); err != nil {
					b.Fatalf("could not get key %s: %s", key, err)
				}
			}
		})
	}
}

// Benchmark the Put operation on the available stores; this measures the
// number of nanoseconds per operation, not necessarily the concurrent access
// throughput as defined by the Measurement tests with each store.
func BenchmarkPut(b *testing.B) {
	for _, store := range makeStores(b) {
		b.Run(store.String(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				key := makeKey(i % 256)
				if err := store.Put(key, []byte(key)); err != nil {
					b.Fatalf("could not put key %s: %s", key, err)
				}
			}
		})
	}
}

// Benchmark the GetorCreate operation on the available stores when the store
// is empty (e.g. the default value will be inserted); this measures the
// number of nanoseconds per operation, not necessarily the concurrent access
// throughput as defined by the Measurement tests with each store.
func BenchmarkGetOrCreateEmpty(b *testing.B) {
	for _, store := range makeStores(b) {
		b.Run(store.String(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				key := makeKey(i)
				store.GetOrCreate(key, []byte(key))
			}
		})
	}
}

// Benchmark the GetorCreate operation on the available stores when the store
// is not empty (e.g. the actual value will be fetched); this measures the
// number of nanoseconds per operation, not necessarily the concurrent access
// throughput as defined by the Measurement tests with each store.
func BenchmarkGetOrCreateFull(b *testing.B) {
	for _, store := range makeStores(b) {
		b.Run(store.String(), func(b *testing.B) {
			// Allocate a bunch of keys into the store
			for k := 0; k < 256; k++ {
				key := makeKey(k)
				store.Put(key, []byte(key))
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := makeKey(i % 256)
				if _, created := store.GetOrCreate(key, []byte(key)); created {
					b.Fatalf("could not get or create key %s, not created", key)
				}
			}
		})
	}
}

// Benchmark the Delete operation on the available stores; this measures the
// number of nanoseconds per operation, not necessarily the concurrent access
// throughput as defined by the Measurement tests with each store.
func BenchmarkDelete(b *testing.B) {
	for _, store := range makeStores(b) {
		b.Run(store.String(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// NOTE: cannot protect the put key in stop and start timer because it hangs.
				// See: https://stackoverflow.com/questions/37620251/golang-benchmarking-b-stoptimer-hangs-is-it-me
				// b.StopTimer()
				// Make sure the key is in the database
				key := makeKey(i % 256)
				if err := store.Put(key, []byte(key)); err != nil {
					b.Fatalf("could not prepare store for delete operation for key %s: %s", key, err)
				}

				// b.StartTimer()
				if err := store.Delete(key); err != nil {
					b.Fatalf("could not delete key %s: %s", key, err)
				}
			}
		})
	}
}
