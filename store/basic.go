package store

import (
	"fmt"
	"sync"
)

// Basic implements a simple key/value data structure that is synchronized
// using a RWMutex. It is ready to go on the first access (e.g. it doesn't
// have to be initialized).
type Basic struct {
	sync.RWMutex
	data map[string][]byte
}

// Init the internal map.
func (s *Basic) Init() (err error) {
	s.data = make(map[string][]byte)
	return nil
}

// Get a value by read locking the internal map and fetching it. If the key
// is not in the map, returns an error.
func (s *Basic) Get(key string) (value []byte, err error) {
	s.RLock()
	defer s.RUnlock()

	val, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("no value found for key '%s'", key)
	}
	return val, nil
}

// Put a value by locking the internal map and storing it. No error returned.
func (s *Basic) Put(key string, value []byte) (err error) {
	s.Lock()
	defer s.Unlock()

	s.data[key] = value
	return nil
}

// Delete a key by locking the internal map. No error returned even if the
// key isn't in the map to begin with.
func (s *Basic) Delete(key string) (err error) {
	s.Lock()
	defer s.Unlock()

	delete(s.data, key)
	return nil
}

// GetOrCreate returns the value stored or stores the supplied default value.
func (s *Basic) GetOrCreate(key string, value []byte) (actual []byte, created bool) {
	s.Lock()
	defer s.Unlock()

	var found bool
	actual, found = s.data[key]

	if !found {
		s.data[key] = value
		return value, true
	}

	return actual, false
}

// String returns a string representation of the Store
func (s *Basic) String() string {
	return "basic"
}
