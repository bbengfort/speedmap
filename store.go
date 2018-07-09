package speedmap

// Store represents the interface for all in-memory key/value data structures
// that are being benchmarked by the Speed Map package.
type Store interface {
	Init() (err error)
	Get(key string) (value []byte, err error)
	Put(key string, value []byte) (err error)
	Delete(key string) (err error)
	GetOrCreate(key string, value []byte) (actual []byte, created bool)
}
