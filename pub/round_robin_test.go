package pub_test

import (
	. "github.com/masslessparticle/goq/pub"

	"github.com/masslessparticle/goq"
	"github.com/masslessparticle/goq/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RoundRobin", func() {
	Context("Publish", func() {
		var subscribers goq.Subscribers

		BeforeEach(func() {
			subscribers = goq.NewSubscribersList()
		})

		It("doesn't deliver the message when there aren't subscribers", func() {
			roundRobin := NewRoundRobinPublisher()

			delivered := roundRobin.Publish(goq.Message{Id: "Message - 1"}, &subscribers)
			Expect(delivered).To(BeFalse())
		})

		It("doesn't deliver messages when subscribers haven't been set", func() {
			roundRobin := NewRoundRobinPublisher()

			delivered := roundRobin.Publish(goq.Message{Id: "Message - 1"}, &subscribers)
			Expect(delivered).To(BeFalse())
		})

		It("delivers messages to all subscribers using round robin", func() {
			subscriber := testhelpers.NewTestClient("Subscriber - 1")
			subscriber2 := testhelpers.NewTestClient("Subscriber - 2")
			subscriber3 := testhelpers.NewTestClient("Subscriber - 3")

			subscribers.Append(subscriber)
			subscribers.Append(subscriber2)
			subscribers.Append(subscriber3)

			roundRobin := NewRoundRobinPublisher()

			delivered := roundRobin.Publish(goq.Message{Id: "Message - 1"}, &subscribers)
			Expect(delivered).To(BeTrue())

			delivered = roundRobin.Publish(goq.Message{Id: "Message - 2"}, &subscribers)
			Expect(delivered).To(BeTrue())

			delivered = roundRobin.Publish(goq.Message{Id: "Message - 3"}, &subscribers)
			Expect(delivered).To(BeTrue())

			delivered = roundRobin.Publish(goq.Message{Id: "Message - 4"}, &subscribers)
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

			subscribers.Append(subscriber)
			subscribers.Append(subscriber2)
			subscribers.Append(subscriber3)

			roundRobin := NewRoundRobinPublisher()

			delivered := roundRobin.Publish(goq.Message{Id: "Message - 1"}, &subscribers)
			Expect(delivered).To(BeTrue())

			subscribers.Remove(subscriber2)

			delivered = roundRobin.Publish(goq.Message{Id: "Message - 2"}, &subscribers)
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

			subscribers.Append(subscriber)
			subscribers.Append(subscriber2)

			roundRobin := NewRoundRobinPublisher()

			delivered := roundRobin.Publish(goq.Message{Id: "Message - 1"}, &subscribers)
			Expect(delivered).To(BeTrue())

			delivered = roundRobin.Publish(goq.Message{Id: "Message - 2"}, &subscribers)
			Expect(delivered).To(BeTrue())

			message := goq.Message{}
			Eventually(subscriber.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 1"))

			Eventually(subscriber2.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 2"))

			subscribers.Append(subscriber3)

			delivered = roundRobin.Publish(goq.Message{Id: "Message - 3"}, &subscribers)
			Expect(delivered).To(BeTrue())

			Eventually(subscriber3.Notifications).Should(Receive(&message))
			Expect(message.Id).To(Equal("Message - 3"))
		})
	})
})
