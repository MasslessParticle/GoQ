package pubsub_test

import (
	. "github.com/masslessparticle/goq/pubsub"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/masslessparticle/goq"
	"github.com/masslessparticle/goq/testhelpers"
)

var _ = Describe("LastUsed", func() {
	var leastUsed *LeastUsedPublisher

	BeforeEach(func() {
		leastUsed = NewLeastUsedPublisher()
	})

	Context("Publish", func() {
		It("doesn't deliver the message when there aren't subscribers", func() {
			delivered := leastUsed.Publish(goq.Message{Id: "Message - 1"})
			Expect(delivered).To(BeFalse())
		})

		It("delivers messages to the least used subscriber", func() {
			subscriber := testhelpers.NewTestClient("Subscriber - 1")
			subscriber2 := testhelpers.NewTestClient("Subscriber - 2")

			leastUsed.Subscribe(subscriber)

			delivered := leastUsed.Publish(goq.Message{Id: "Message - 1"})
			Expect(delivered).To(BeTrue())

			delivered = leastUsed.Publish(goq.Message{Id: "Message - 2"})
			Expect(delivered).To(BeTrue())

			leastUsed.Subscribe(subscriber2)

			delivered = leastUsed.Publish(goq.Message{Id: "Message - 3"})
			Expect(delivered).To(BeTrue())

			message := goq.Message{}
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))

			Eventually(subscriber2.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 3"))
		})

		It("doesn't send messages to removed subscribers", func() {
			subscriber := testhelpers.NewTestClient("Subscriber - 1")
			subscriber2 := testhelpers.NewTestClient("Subscriber - 2")

			leastUsed.Subscribe(subscriber)
			leastUsed.Subscribe(subscriber2)

			delivered := leastUsed.Publish(goq.Message{Id: "Message - 1"})
			Expect(delivered).To(BeTrue())

			delivered = leastUsed.Publish(goq.Message{Id: "Message - 2"})
			Expect(delivered).To(BeTrue())

			leastUsed.Unsubscribe(subscriber2)

			delivered = leastUsed.Publish(goq.Message{Id: "Message - 3"})
			Expect(delivered).To(BeTrue())

			message := goq.Message{}
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))
			Eventually(subscriber2.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))

			Consistently(subscriber2.Notifications).ShouldNot(Receive())
		})

		It("sends messages to new subscribers", func() {
			subscriber := testhelpers.NewTestClient("Subscriber - 1")
			subscriber2 := testhelpers.NewTestClient("Subscriber - 2")
			subscriber3 := testhelpers.NewTestClient("Subscriber - 3")

			leastUsed.Subscribe(subscriber)
			leastUsed.Subscribe(subscriber2)

			delivered := leastUsed.Publish(goq.Message{Id: "Message - 1"})
			Expect(delivered).To(BeTrue())

			delivered = leastUsed.Publish(goq.Message{Id: "Message - 2"})
			Expect(delivered).To(BeTrue())

			message := goq.Message{}
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))
			Eventually(subscriber2.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))

			leastUsed.Subscribe(subscriber3)

			delivered = leastUsed.Publish(goq.Message{Id: "Message - 3"})
			Expect(delivered).To(BeTrue())

			Eventually(subscriber3.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 3"))
		})
	})
})
