package pubsub_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/masslessparticle/goq/testhelpers"
	"github.com/masslessparticle/goq/pubsub"
)

var _ = Describe("SubscribersList", func() {
	It("doesn't allow duplicate ids", func() {
		subscribers := pubsub.NewSubscribersList()
		err := subscribers.Subscribe(testhelpers.NewTestClient("Subscriber - 1"))
		Expect(err).ToNot(HaveOccurred())

		err = subscribers.Subscribe(testhelpers.NewTestClient("Subscriber - 1"))
		Expect(err).To(HaveOccurred())
	})
})
