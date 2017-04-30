package goq_test

import (
	. "github.com/masslessparticle/goq"

	"github.com/masslessparticle/goq/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SubscribersList", func() {
	It("doesn't allow duplicate ids", func() {
		subscribers := NewSubscribersList()
		err := subscribers.Append(testhelpers.NewTestClient("Subscriber - 1"))
		Expect(err).ToNot(HaveOccurred())

		err = subscribers.Append(testhelpers.NewTestClient("Subscriber - 1"))
		Expect(err).To(HaveOccurred())
	})
})
