package store

import (
	"fmt"
	"sync"
)

// ShardCount specifies the number of shards the store contains.
const ShardCount = 32

// Shard implements a high performance solution to concurrent reads and
// writes by sharding the keyspace across multiple maps. This allows many
// concurrent operations to be lock-free in their operation. The original
// implementation is https://github.com/orcaman/concurrent-map.
//
// TODO: do we use a sync.Map or another store and make Shard generic?
type Shard []*shard

// A thread safe map that implements the portion of the shard's keyspace.
type shard struct {
	sync.RWMutex
	data map[string][]byte
}

// NewShard creates the shard store and initializes the keyspace and the
// internal shards of the keyspace.
func NewShard() (store Shard, err error) {
	store = make(Shard, ShardCount)
	for i := 0; i < ShardCount; i++ {
		store[i] = &shard{data: make(map[string][]byte)}
	}
	return store, err
}

// GetShard returns the shard the key is assigned to.
func (s Shard) GetShard(key string) *shard {
	return s[uint(fnv32(key))%uint(ShardCount)]
}

// Get a value by finding the shard the key belongs to, fetching and locking
// that shard before returning. Unlike the Basic store, does not use defer unlock.
func (s Shard) Get(key string) (value []byte, err error) {
	shard := s.GetShard(key)
	shard.RLock()
	val, ok := shard.data[key]
	shard.RUnlock()

	if !ok {
		return nil, fmt.Errorf("no value found for key '%s'", key)
	}
	return val, nil
}

// Put a value by finding the shard the key belongs to, fetching and locking
// the shard and assigning the value to the key. Unlike the Basic store,
// does not use defer unlock.
func (s Shard) Put(key string, value []byte) (err error) {
	shard := s.GetShard(key)
	shard.Lock()
	shard.data[key] = value
	shard.Unlock()
	return nil
}

// Delete a key by finding the shard the key belongs to, fetching and locking
// the shard and deleting the key. Unlike the Basic store, does not use defer unlock.
func (s Shard) Delete(key string) (err error) {
	shard := s.GetShard(key)
	shard.Lock()
	delete(shard.data, key)
	shard.Unlock()
	return nil
}

// GetOrCreate returns the value stored or stores the default value by finding
// the shard the key belongs to, fetching and locking it then checking if the
// key exists, setting if necessary. Unlike the Basic store, does not use defer unlock.
func (s Shard) GetOrCreate(key string, value []byte) (actual []byte, created bool) {
	shard := s.GetShard(key)
	shard.Lock()

	var found bool
	actual, found = shard.data[key]

	if !found {
		shard.data[key] = value
		shard.Unlock()
		return value, true
	}

	shard.Unlock()
	return actual, false
}

// String returns the string representation of the sharded store.
func (s Shard) String() string {
	return "shard"
}

// Computes the uint32 fingerprint of the specified key.
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
