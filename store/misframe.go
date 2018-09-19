package store

import (
	"fmt"
	"sync"
)

// Misframe extends the Basic synchronized map structure with an optimized
// GetOrCreate method that doesn't use a Lock but rather uses a two phase
// check and lock cycle with both read and write locks, as described at:
// https://misfra.me/optimizing-concurrent-map-access-in-go/
type Misframe struct {
	sync.RWMutex
	data map[string][]byte
}

// NewMisframe creates the Misframe store and initializes the internal map.
func NewMisframe() (store *Misframe, err error) {
	store = new(Misframe)
	store.data = make(map[string][]byte)
	return store, nil
}

// Get a value by read locking the internal map and fetching it. If the key
// is not in the map, returns an error.
func (s *Misframe) Get(key string) (value []byte, err error) {
	s.RLock()
	defer s.RUnlock()

	val, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("no value found for key '%s'", key)
	}
	return val, nil
}

// Put a value by locking the internal map and storing it. No error returned.
func (s *Misframe) Put(key string, value []byte) (err error) {
	s.Lock()
	defer s.Unlock()

	s.data[key] = value
	return nil
}

// Delete a key by locking the internal map. No error returned even if the
// key isn't in the map to begin with.
func (s *Misframe) Delete(key string) (err error) {
	s.Lock()
	defer s.Unlock()

	delete(s.data, key)
	return nil
}

// GetOrCreate returns the value stored or stores the supplied default value.
func (s *Misframe) GetOrCreate(key string, value []byte) (actual []byte, created bool) {
	var present bool

	s.RLock()
	if actual, present = s.data[key]; !present {
		// The source wasn't found, so we'll create it.
		s.RUnlock()
		s.Lock()
		if actual, present = s.data[key]; !present {
			// Insert the value.
			s.data[key] = value
			actual = value
		}
		s.Unlock()
		return actual, true
	}
	s.RUnlock()
	return actual, false
}

// String returns a string representation of the Store
func (s *Misframe) String() string {
	return "misframe"
}
