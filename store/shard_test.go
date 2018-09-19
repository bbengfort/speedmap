package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bbengfort/speedmap"
	. "github.com/bbengfort/speedmap/store"
)

var _ = Describe("Shard", func() {

	var (
		err   error
		store speedmap.Store
	)

	BeforeEach(func() {
		store, err = NewShard()
		Ω(err).ShouldNot(HaveOccurred())
	})

	It("should be a store", func() {
		Ω(make(Shard, 0)).Should(BeAssignableToTypeOf(store))
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

})
