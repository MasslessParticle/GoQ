package acceptance

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/masslessparticle/goq/pub"
	"github.com/masslessparticle/goq"
	"github.com/masslessparticle/goq/testhelpers"
)

var _ = Describe("Notification", func() {
	It("Notifies a subscribed client", func() {
		publisher := pub.NewRoundRobinPublisher()
		queue := goq.NewGoQ(25, publisher)

		client := testhelpers.NewTestClient("Subscription - 1")
		queue.Subscribe(client)

		queue.StartPublishing()

		queue.Enqueue(goq.Message{Id: "Message - 1"})

		message := goq.Message{}
		Eventually(client.Notifications).Should(Receive(&message))
		Expect(message.Id).To(Equal("Message - 1"))
	})
})
