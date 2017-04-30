package pubsub_test

import (
	"github.com/masslessparticle/goq"
	"github.com/masslessparticle/goq/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/masslessparticle/goq/pubsub"
)

var _ = Describe("RoundRobin", func() {
	var roundRobin *pubsub.RoundRobinPublisher

	BeforeEach(func () {
		roundRobin = pubsub.NewRoundRobinPublisher()
	})

	Context("Publish", func() {
		It("doesn't deliver the message when there aren't subscribers", func() {
			delivered := roundRobin.Publish(goq.Message{Id: "Message - 1"})
			Expect(delivered).To(BeFalse())
		})

		It("delivers messages to all subscribers using round robin", func() {
			subscriber := testhelpers.NewTestClient("Subscriber - 1")
			subscriber2 := testhelpers.NewTestClient("Subscriber - 2")
			subscriber3 := testhelpers.NewTestClient("Subscriber - 3")

			roundRobin.Subscribe(subscriber)
			roundRobin.Subscribe(subscriber2)
			roundRobin.Subscribe(subscriber3)

			delivered := roundRobin.Publish(goq.Message{Id: "Message - 1"})
			Expect(delivered).To(BeTrue())

			delivered = roundRobin.Publish(goq.Message{Id: "Message - 2"})
			Expect(delivered).To(BeTrue())

			delivered = roundRobin.Publish(goq.Message{Id: "Message - 3"})
			Expect(delivered).To(BeTrue())

			delivered = roundRobin.Publish(goq.Message{Id: "Message - 4"})
			Expect(delivered).To(BeTrue())

			message := goq.Message{}
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))

			Eventually(subscriber2.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))

			Eventually(subscriber3.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 3"))

			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 4"))
		})

		It("doesn't send messages to removed subscribers", func() {
			subscriber := testhelpers.NewTestClient("Subscriber - 1")
			subscriber2 := testhelpers.NewTestClient("Subscriber - 2")
			subscriber3 := testhelpers.NewTestClient("Subscriber - 3")

			roundRobin.Subscribe(subscriber)
			roundRobin.Subscribe(subscriber2)
			roundRobin.Subscribe(subscriber3)

			delivered := roundRobin.Publish(goq.Message{Id: "Message - 1"})
			Expect(delivered).To(BeTrue())

			roundRobin.Unsubscribe(subscriber2)

			delivered = roundRobin.Publish(goq.Message{Id: "Message - 2"})
			Expect(delivered).To(BeTrue())

			message := goq.Message{}
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))

			Eventually(subscriber3.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))

			Consistently(subscriber2.Notifications).ShouldNot(Receive())
		})

		It("sends messages to new subscribers", func() {
			subscriber := testhelpers.NewTestClient("Subscriber - 1")
			subscriber2 := testhelpers.NewTestClient("Subscriber - 2")
			subscriber3 := testhelpers.NewTestClient("Subscriber - 3")

			roundRobin.Subscribe(subscriber)
			roundRobin.Subscribe(subscriber2)

			delivered := roundRobin.Publish(goq.Message{Id: "Message - 1"})
			Expect(delivered).To(BeTrue())

			delivered = roundRobin.Publish(goq.Message{Id: "Message - 2"})
			Expect(delivered).To(BeTrue())

			message := goq.Message{}
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))

			Eventually(subscriber2.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))

			roundRobin.Subscribe(subscriber3)

			delivered = roundRobin.Publish(goq.Message{Id: "Message - 3"})
			Expect(delivered).To(BeTrue())

			Eventually(subscriber3.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 3"))
		})
	})
})
