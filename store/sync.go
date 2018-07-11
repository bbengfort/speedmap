package store

import (
	"fmt"
	"sync"
)

// SyncMap is just a wrapper to sync.Map to provide the specified interface.
type SyncMap struct {
	data *sync.Map
}

// Init creates the sync.Map
func (s *SyncMap) Init() (err error) {
	s.data = new(sync.Map)
	return nil
}

// Get is an alias for sync.Map.Load. Returns an error if key not found or
// if the value cannot be cast to bytes.
func (s *SyncMap) Get(key string) (value []byte, err error) {
	data, ok := s.data.Load(key)
	if !ok {
		return nil, fmt.Errorf("no value found for key '%s'", key)
	}

	value, ok = data.([]byte)
	if !ok {
		return nil, fmt.Errorf("could not cast value to bytes")
	}

	return value, nil
}

// Put is an alias for sync.Map.Store. Does not return an error.
func (s *SyncMap) Put(key string, value []byte) (err error) {
	s.data.Store(key, value)
	return nil
}

// Delete is an alias for sync.Map.Delete. Does not return an error.
func (s *SyncMap) Delete(key string) (err error) {
	s.data.Delete(key)
	return nil
}

// GetOrCreate is an alias for sync.Map.LoadOrStore.
func (s *SyncMap) GetOrCreate(key string, value []byte) (actual []byte, created bool) {
	data, loaded := s.data.LoadOrStore(key, value)
	return data.([]byte), !loaded
}
