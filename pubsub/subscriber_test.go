package pubsub_test

import (
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/masslessparticle/goq"
	"github.com/masslessparticle/goq/testhelpers"
	"github.com/masslessparticle/goq/pubsub"
)

var _ = Describe("PublisherSubscriber", func() {
	DescribeTable("Publishers can subscribe clients",
		func(publisher goq.PubSub) {
			client := testhelpers.NewTestClient("subscriber-1")

			err := publisher.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(publisher.SubscriberCount()).To(Equal(1))
		},
		Entry("round robin", pubsub.NewRoundRobinPublisher()),
		Entry("all the things", pubsub.NewAllPublisher()),
		Entry("least used", pubsub.NewLeastUsedPublisher()),
	)

	DescribeTable("Publishers can unsubscribe clients",
		func(publisher goq.PubSub) {
			client := testhelpers.NewTestClient("subscriber-1")

			err := publisher.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(publisher.SubscriberCount()).To(Equal(1))

			publisher.Unsubscribe(client)
			Expect(publisher.SubscriberCount()).To(Equal(0))
		},
		Entry("round robin", pubsub.NewRoundRobinPublisher()),
		Entry("all the things", pubsub.NewAllPublisher()),
		Entry("Least Used", pubsub.NewLeastUsedPublisher()),
	)
})
