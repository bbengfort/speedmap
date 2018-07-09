# Speed Maps

**Benchmarks for concurrent access to key/value data structures.**

Here's the scenario: multiple clients are reading and writing values to different keys in a shared data structure. Each client is in their own go routine, unfortunately this means that a lock or a channel is required to synchronize accesses to the store. In this repository we explore different synchronization methods for concurrent access to the store to find out what is the most performant for various workloads.

There are two parts to each benchmark:

1. A data structure that implements the `Store` interface.
2. A workload that defines access patterns for increasing threads.

The store interface is pretty straight forward:

```go
type Store interface {
    Init() (err error)
    Get(key string) (value []byte, err error)
    Put(key string, value []byte) (err error)
    Delete(key string) (err error)
    GetOrCreate(key string, value []byte) (actual []byte, created bool)
}
```

The following stores have been implemented:

1. Basic: wraps a `map[string][]byte` with a `sync.RWMutex` (baseline) and treats `GetDefault` as a write operation.
2. Misframe: optimizes `GetDefault` as described in [Optimizing Concurrent Map Access in Go](https://misfra.me/optimizing-concurrent-map-access-in-go/).
3. [`sync.Map`](https://golang.org/pkg/sync/#Map): the official concurrent map object in the sync package.
