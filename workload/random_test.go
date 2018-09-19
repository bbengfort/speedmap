package workload_test

import (
	"crypto/md5"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bbengfort/speedmap/workload"
)

var _ = Describe("Random", func() {

	var (
		buf []byte
		err error
	)

	It("should generate random bytes", func() {
		var b []byte

		buf, err = GenerateRandomBytes(DataSize)
		Ω(err).ShouldNot(HaveOccurred())

		for i := 0; i < 100; i++ {
			b, err = GenerateRandomBytes(DataSize)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(b).ShouldNot(Equal(buf))
			Ω(b).Should(HaveLen(len(buf)))

			buf = b
		}

	})

	It("should generate varying length byte arrays", func() {
		for i := 1; i < 17; i++ {
			size := i * 512
			buf, err = GenerateRandomBytes(size)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(buf).Should(HaveLen(size))
		}
	})

	It("should generate random keys", func() {
		cases := []struct {
			keyspace, minkey, maxkey int64
		}{
			{0, 0, 1000}, {1, 1000, 2000}, {2, 2000, 3000}, {3, 3000, 4000},
		}

		for _, test := range cases {
			dups := 0
			prev := RandomKey(test.keyspace, 1000)

			for i := 0; i < 100; i++ {
				next := RandomKey(test.keyspace, 1000)
				if next == prev {
					dups++
				}

				num, err := strconv.ParseInt(next, 16, 64)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(num).Should(BeNumerically(">=", test.minkey))
				Ω(num).Should(BeNumerically("<", test.maxkey))

				prev = next
			}

			Ω(dups).Should(BeNumerically("<", 3))
		}

	})

	It("should be able to randomly mutate a buffer", func() {

		data, err := GenerateRandomBytes(512)
		Ω(err).ShouldNot(HaveOccurred())

		checksum := md5.Sum(data)

		for i := 0; i < 1024; i++ {
			RandomMutation(data, 32)
			newsum := md5.Sum(data)
			Ω(newsum).ShouldNot(Equal(checksum))
			checksum = newsum
		}
	})

})

func BenchmarkRandomKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomKey(int64(i%1000), 1000)
	}
}

func BenchmarkGenerateRandomBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateRandomBytes(2048)
	}
}

func BenchmarkRandomMutation(b *testing.B) {
	data, _ := GenerateRandomBytes(2048)
	for i := 0; i < b.N; i++ {
		RandomMutation(data, 128)
	}
}
