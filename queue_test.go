package goq_test

import (
	. "github.com/masslessparticle/goq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/masslessparticle/goq/testhelpers"
)

var _ = Describe("Queue", func() {
	Context("Enqueue", func() {
		It("is empty when it's new", func() {
			queue := NewGoQ(25)
			enqueuedItems := queue.QueuedMessages()
			Expect(enqueuedItems).To(Equal(0))
		})

		It("can have messages enqueued", func() {
			queue := NewGoQ(25)
			err := queue.Enqueue(Message{
				Id: "1",
				Payload: "The Message",
			})
			Expect(err).ToNot(HaveOccurred())

			err = queue.Enqueue(Message{
				Id: "2",
				Payload: "Another Message",
			})
			Expect(err).ToNot(HaveOccurred())

			Expect(queue.QueuedMessages()).To(Equal(2))
		})

		It("throws an error when the queuedepth is exceeded", func() {
			queue := NewGoQ(1)
			err := queue.Enqueue(Message{
				Id: "1",
				Payload: "The Message",
			})
			Expect(err).ToNot(HaveOccurred())

			err = queue.Enqueue(Message{
				Id: "2",
				Payload: "Another Message",
			})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Subscribe", func() {
		It("can store client subscriptions", func() {
			queue := NewGoQ(25)
			client := testhelpers.NewTestClient("subscriber-1")

			err := queue.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(queue.IsSubscribed(client)).To(BeTrue())
		})

		It("returns an error when a client with the same Id attempts to subscribe", func() {
			queue := NewGoQ(25)

			client := testhelpers.NewTestClient("subscriber-1")

			err := queue.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())

			client2 := testhelpers.NewTestClient("subscriber-1")

			err = queue.Subscribe(client2)
			Expect(err).To(HaveOccurred())

			Expect(queue.IsSubscribed(client)).To(BeTrue())
		})

		It("can unsubscribe clients", func() {
			queue := NewGoQ(25)
			client := testhelpers.NewTestClient("subscriber-1")

			err := queue.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(queue.IsSubscribed(client)).To(BeTrue())

			queue.Unsubscribe(client)
			Expect(queue.IsSubscribed(client)).To(BeFalse())
		})
	})

	Context("Notifications", func() {
		It("notifies the only subscriber when a message is recieved", func() {
			queue := NewGoQ(25)
			queue.StartPublishing()

			client := testhelpers.NewTestClient("subscriber-1")

			err := queue.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())

			queue.Enqueue(Message{
				Id: "MessageId - 1",
				Payload: "This is the message",
			})

			recievedMessage := Message{}
			Eventually(client.Notifications).Should(Receive(&recievedMessage))
			Expect(recievedMessage.Payload).To(Equal("This is the message"))
		})

		It("notifies subscribers with a round robin strategy", func() {
			queue := NewGoQ(25)

			client := testhelpers.NewTestClient("subscriber-1")
			client2 := testhelpers.NewTestClient("subscriber-2")
			client3 := testhelpers.NewTestClient("subscriber-3")

			err := queue.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())

			err = queue.Subscribe(client2)
			Expect(err).ToNot(HaveOccurred())

			err = queue.Subscribe(client3)
			Expect(err).ToNot(HaveOccurred())

			queue.Enqueue(Message{Id: "MessageId - 1"})
			queue.Enqueue(Message{Id: "MessageId - 2"})
			queue.Enqueue(Message{Id: "MessageId - 3"})
			queue.Enqueue(Message{Id: "MessageId - 4"})

			queue.StartPublishing()

			recievedMessage := Message{}
			Eventually(client.Notifications).Should(Receive(&recievedMessage))
			Expect(recievedMessage.Id).To(Equal("MessageId - 1"))

			Eventually(client2.Notifications).Should(Receive(&recievedMessage))
			Expect(recievedMessage.Id).To(Equal("MessageId - 2"))

			Eventually(client3.Notifications).Should(Receive(&recievedMessage))
			Expect(recievedMessage.Id).To(Equal("MessageId - 3"))

			Eventually(client.Notifications).Should(Receive(&recievedMessage))
			Expect(recievedMessage.Id).To(Equal("MessageId - 4"))
		})

		It("doesn't send notifications after stopping publishing", func() {
			queue := NewGoQ(25)
			queue.StartPublishing()

			client := testhelpers.NewTestClient("subscriber-1")

			err := queue.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())

			queue.Enqueue(Message{Id: "MessageId - 1"})

			recievedMessage := Message{}
			Eventually(client.Notifications).Should(Receive(&recievedMessage))
			Expect(recievedMessage.Id).To(Equal("MessageId - 1"))

			queue.StopPublishing()

			queue.Enqueue(Message{Id: "MessageId - 2"})
			Consistently(client.Notifications).ShouldNot(Receive())
		})

		It("doesn't send notifications to unsubscribed clients", func() {
			queue := NewGoQ(25)
			queue.StartPublishing()

			client := testhelpers.NewTestClient("subscriber-1")

			err := queue.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())

			queue.Enqueue(Message{Id: "MessageId - 1"})

			recievedMessage := Message{}
			Eventually(client.Notifications).Should(Receive(&recievedMessage))
			Expect(recievedMessage.Id).To(Equal("MessageId - 1"))

			queue.Unsubscribe(client)

			queue.Enqueue(Message{Id: "MessageId - 2"})
			Consistently(client.Notifications).ShouldNot(Receive())
		})

		It("sends messages to new subscribers after old ones have unsubscribed", func() {
			queue := NewGoQ(25)
			queue.StartPublishing()

			client := testhelpers.NewTestClient("subscriber-1")

			err := queue.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())

			queue.Enqueue(Message{Id: "MessageId - 1"})

			recievedMessage := Message{}
			Eventually(client.Notifications).Should(Receive(&recievedMessage))
			Expect(recievedMessage.Id).To(Equal("MessageId - 1"))

			queue.Unsubscribe(client)

			queue.Enqueue(Message{Id: "MessageId - 2"})
			Consistently(client.Notifications).ShouldNot(Receive())

			err = queue.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())
			Eventually(client.Notifications).Should(Receive(&recievedMessage))
			Expect(recievedMessage.Id).To(Equal("MessageId - 2"))
		})
	})
})
