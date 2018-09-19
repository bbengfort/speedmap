package store_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bbengfort/speedmap"
	. "github.com/bbengfort/speedmap/store"
)

var _ = Describe("Sync", func() {

	var (
		err   error
		store speedmap.Store
	)

	BeforeEach(func() {
		store, err = NewSyncMap()
		Ω(err).ShouldNot(HaveOccurred())
	})

	It("should be a store", func() {
		Ω(&SyncMap{}).Should(BeAssignableToTypeOf(store))
	})

	It("should be able to perform store operations", func() {
		Ω(store.Put("foo", []byte("bar"))).Should(Succeed())

		val, err := store.Get("foo")
		Ω(err).ShouldNot(HaveOccurred())
		Ω(val).Should(Equal([]byte("bar")))

		Ω(store.Delete("foo")).Should(Succeed())

		val, err = store.Get("foo")
		Ω(err).Should(HaveOccurred())
		Ω(val).Should(BeNil())
	})

	It("should be able to get or create a value", func() {
		actual, created := store.GetOrCreate("foo", []byte("bar"))
		Ω(actual).Should(Equal([]byte("bar")))
		Ω(created).Should(BeTrue())

		actual, created = store.GetOrCreate("foo", []byte("red"))
		Ω(actual).Should(Equal([]byte("bar")))
		Ω(created).Should(BeFalse())
	})

	Measure("get throughput", func(b Benchmarker) {
		// Populate the store
		for i := 0; i < 5000; i++ {
			key := fmt.Sprintf("%X", i)
			store.Put(key, []byte(key))
		}

		results, err := Blast(store, 5000, "Get")
		Ω(err).ShouldNot(HaveOccurred())
		b.RecordValue("throughput", results.Throughput)
	}, 10)

	Measure("put throughput", func(b Benchmarker) {
		results, err := Blast(store, 5000, "Put")
		Ω(err).ShouldNot(HaveOccurred())
		b.RecordValue("throughput", results.Throughput)
	}, 10)

	Measure("delete throughput", func(b Benchmarker) {
		// Populate the store
		for i := 0; i < 5000; i++ {
			key := fmt.Sprintf("%X", i)
			store.Put(key, []byte(key))
		}

		results, err := Blast(store, 5000, "Delete")
		Ω(err).ShouldNot(HaveOccurred())
		b.RecordValue("throughput", results.Throughput)
	}, 10)

	Measure("get or create throughput", func(b Benchmarker) {
		results, err := Blast(store, 5000, "GetOrCreate")
		Ω(err).ShouldNot(HaveOccurred())
		b.RecordValue("throughput", results.Throughput)
	}, 10)
})
