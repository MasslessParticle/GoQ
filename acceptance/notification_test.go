package acceptance

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/masslessparticle/goq"
	"github.com/masslessparticle/goq/testhelpers"
	"github.com/masslessparticle/goq/pubsub"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("Notification", func() {
	DescribeTable("Notifies the client with a publisher",
		func(publisher goq.PubSub) {
			client := testhelpers.NewTestClient("Subscription - 1")
			publisher.Subscribe(client)

			queue := goq.NewGoQ(25, publisher)

			queue.StartPublishing()

			queue.Enqueue(goq.Message{Id: "Message - 1"})

			message := goq.Message{}
			Eventually(client.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))
		},
		Entry("round robin", pubsub.NewRoundRobinPublisher()),
		Entry("all", pubsub.NewAllPublisher()),
		Entry("least used", pubsub.NewLeastUsedPublisher()),
	)
})
