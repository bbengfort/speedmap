package workload

const (
	// MaxKeys being accessed by all threads
	MaxKeys = 5000

	// OpsPerThread being executed against the store
	OpsPerThread = 5000

	// DataSize is the number of bytes written to each key
	DataSize = 32
)
