package acceptance

import (
	"github.com/masslessparticle/goq"
	"github.com/masslessparticle/goq/pubsub"
	"github.com/masslessparticle/goq/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Notification", func() {
	DescribeTable("Notifies the client with a publisher",
		func(publisher goq.PubSub) {
			client := testhelpers.NewTestClient("Subscription - 1")
			publisher.Subscribe(client)

			queue := goq.NewGoQ(25, publisher)

			queue.StartPublishing()

			queue.Enqueue(goq.Message{ID: "Message - 1"})

			message := goq.Message{}
			Eventually(client.Notifications).Should(Receive(&message))
			Expect(message.ID).To(Equal("Message - 1"))
		},
		Entry("round robin", pubsub.NewRoundRobinPublisher()),
		Entry("all", pubsub.NewAllPublisher()),
		Entry("least used", pubsub.NewLeastUsedPublisher()),
	)
})
