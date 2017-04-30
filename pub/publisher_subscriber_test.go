package pub_test

import (
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/masslessparticle/goq"
	"github.com/masslessparticle/goq/testhelpers"
	"github.com/masslessparticle/goq/pub"
)

var _ = Describe("PublisherSubscriber", func() {
	DescribeTable("Publishers can subscribe clients",
		func(publisher goq.Publisher) {
			client := testhelpers.NewTestClient("subscriber-1")

			err := publisher.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(publisher.SubscriberCount()).To(Equal(1))
		},
		Entry("round robin", pub.NewRoundRobinPublisher()),
		Entry("all the things", pub.NewAllPublisher()),
	)

	DescribeTable("Publishers can unsubscribe clients",
		func(publisher goq.Publisher) {
			client := testhelpers.NewTestClient("subscriber-1")

			err := publisher.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(publisher.SubscriberCount()).To(Equal(1))

			publisher.Unsubscribe(client)
			Expect(publisher.SubscriberCount()).To(Equal(0))
		},
		Entry("round robin", pub.NewRoundRobinPublisher()),
		Entry("all the things", pub.NewAllPublisher()),
	)
})
