package pubsub_test

import (
	"github.com/masslessparticle/goq"
	"github.com/masslessparticle/goq/pubsub"
	"github.com/masslessparticle/goq/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AllPublisher", func() {
	var allPublisher *pubsub.AllPublisher

	BeforeEach(func() {
		allPublisher = pubsub.NewAllPublisher()
	})

	Context("Publish", func() {
		It("doesn't deliver the message when there aren't subscribers", func() {
			delivered := allPublisher.Publish(goq.Message{Id: "Message - 1"})
			Expect(delivered).To(BeFalse())
		})

		It("delivers all messages to all subscribers", func() {
			subscriber := testhelpers.NewTestClient("Subscriber - 1")
			subscriber2 := testhelpers.NewTestClient("Subscriber - 2")
			subscriber3 := testhelpers.NewTestClient("Subscriber - 3")

			allPublisher.Subscribe(subscriber)
			allPublisher.Subscribe(subscriber2)
			allPublisher.Subscribe(subscriber3)

			delivered := allPublisher.Publish(goq.Message{Id: "Message - 1"})
			Expect(delivered).To(BeTrue())

			delivered = allPublisher.Publish(goq.Message{Id: "Message - 2"})
			Expect(delivered).To(BeTrue())

			delivered = allPublisher.Publish(goq.Message{Id: "Message - 3"})
			Expect(delivered).To(BeTrue())

			message := goq.Message{}
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 3"))

			Eventually(subscriber2.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))
			Eventually(subscriber2.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))
			Eventually(subscriber2.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 3"))

			Eventually(subscriber3.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))
			Eventually(subscriber3.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))
			Eventually(subscriber3.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 3"))
		})

		It("doesn't send messages to removed subscribers", func() {
			subscriber := testhelpers.NewTestClient("Subscriber - 1")
			subscriber2 := testhelpers.NewTestClient("Subscriber - 2")

			allPublisher.Subscribe(subscriber)
			allPublisher.Subscribe(subscriber2)

			delivered := allPublisher.Publish(goq.Message{Id: "Message - 1"})
			Expect(delivered).To(BeTrue())

			allPublisher.Unsubscribe(subscriber2)

			delivered = allPublisher.Publish(goq.Message{Id: "Message - 2"})
			Expect(delivered).To(BeTrue())

			message := goq.Message{}
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))

			Eventually(subscriber2.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))

			Consistently(subscriber2.Notifications).ShouldNot(Receive())
		})

		It("sends messages to new subscribers", func() {
			subscriber := testhelpers.NewTestClient("Subscriber - 1")
			subscriber2 := testhelpers.NewTestClient("Subscriber - 2")
			subscriber3 := testhelpers.NewTestClient("Subscriber - 3")

			allPublisher.Subscribe(subscriber)
			allPublisher.Subscribe(subscriber2)

			delivered := allPublisher.Publish(goq.Message{Id: "Message - 1"})
			Expect(delivered).To(BeTrue())

			message := goq.Message{}
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))
			Eventually(subscriber2.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))

			allPublisher.Subscribe(subscriber3)

			delivered = allPublisher.Publish(goq.Message{Id: "Message - 2"})
			Expect(delivered).To(BeTrue())

			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))
			Eventually(subscriber2.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))
			Eventually(subscriber3.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))
		})
	})
})
