package speedmap_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bbengfort/speedmap"
)

var _ = Describe("Basic", func() {

	var store Store

	BeforeEach(func() {
		store = &Basic{}
		Ω(store.Init()).Should(Succeed())
	})

	It("should be a store", func() {
		Ω(&Basic{}).Should(BeAssignableToTypeOf(store))
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
