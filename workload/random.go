package workload

import (
	crypto "crypto/rand"
	"fmt"
	"math/rand"
)

// GenerateRandomBytes returns securely generated random bytes. It will
// return an error if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue.
func GenerateRandomBytes(n int) (b []byte, err error) {
	b = make([]byte, n)

	if _, err = crypto.Read(b); err != nil {
		// Note that err == nil only if we read len(b) bytes.
		return nil, err
	}

	return b, nil
}

// RandomKey returns a random key in the specified key range. For example if
// the keylen is 100, then range 0 is 0-99, range 1 is 100-199, etc. This
// function is used to generate random keys where keys in Range 0 are
// conflicting and keys in another range are not.
func RandomKey(keyspace, keylen int64) string {

	kmin := keyspace * keylen
	kmax := kmin + keylen

	key := rand.Int63n(kmax-kmin) + kmin
	return fmt.Sprintf("%X", key)
}

// RandomMutation modifies the byte array by replacing bytes with random bytes
// this method should be much quicker at modifying an existing random array
// than generating an entirely new random array of equal size.
func RandomMutation(b []byte, m int) (err error) {
	mutations := make([]byte, m)

	if _, err = crypto.Read(mutations); err != nil {
		return err
	}

	blen := len(b)
	for _, mutation := range mutations {
		idx := rand.Intn(blen)
		b[idx] = mutation
	}

	return nil
}
