package pub_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/masslessparticle/goq/testhelpers"
	"github.com/masslessparticle/goq/pub"
)

var _ = Describe("SubscribersList", func() {
	It("doesn't allow duplicate ids", func() {
		subscribers := pub.NewSubscribersList()
		err := subscribers.Append(testhelpers.NewTestClient("Subscriber - 1"))
		Expect(err).ToNot(HaveOccurred())

		err = subscribers.Append(testhelpers.NewTestClient("Subscriber - 1"))
		Expect(err).To(HaveOccurred())
	})
})
