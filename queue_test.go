package goq_test

import (
	. "github.com/masslessparticle/goq"

	"github.com/masslessparticle/goq/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Queue", func() {

	Context("Enqueue", func() {
		var publisher *testhelpers.TestPublisher

		BeforeEach(func() {
			publisher = testhelpers.NewTestPublisher()
			publisher.Responses <- true
		})

		It("can enqueue a message", func() {
			queue := NewGoQ(25, publisher)
			err := queue.Enqueue(Message{
				Id:      "1",
				Payload: "The Message",
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("throws an error when the queuedepth is exceeded", func() {
			queue := NewGoQ(1, publisher)
			err := queue.Enqueue(Message{
				Id:      "1",
				Payload: "The Message",
			})
			Expect(err).ToNot(HaveOccurred())

			err = queue.Enqueue(Message{
				Id:      "2",
				Payload: "Another Message",
			})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Subscribe", func() {
		var publisher *testhelpers.TestPublisher

		BeforeEach(func() {
			publisher = testhelpers.NewTestPublisher()
			publisher.Responses <- true
		})

		It("can store client subscriptions", func() {
			queue := NewGoQ(25, publisher)
			client := testhelpers.NewTestClient("subscriber-1")

			err := queue.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(queue.IsSubscribed(client)).To(BeTrue())
		})

		It("returns an error when a client with the same Id attempts to subscribe", func() {
			queue := NewGoQ(25, publisher)

			client := testhelpers.NewTestClient("subscriber-1")

			err := queue.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())

			client2 := testhelpers.NewTestClient("subscriber-1")

			err = queue.Subscribe(client2)
			Expect(err).To(HaveOccurred())

			Expect(queue.IsSubscribed(client)).To(BeTrue())
		})

		It("can unsubscribe clients", func() {
			queue := NewGoQ(25, publisher)
			client := testhelpers.NewTestClient("subscriber-1")

			err := queue.Subscribe(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(queue.IsSubscribed(client)).To(BeTrue())

			queue.Unsubscribe(client)
			Expect(queue.IsSubscribed(client)).To(BeFalse())
		})
	})

	Context("Notifications", func() {
		It("sends the message to the publisher", func() {
			publisher := testhelpers.NewTestPublisher()
			publisher.Responses <- true

			queue := NewGoQ(25, publisher)
			queue.StartPublishing()

			queue.Enqueue(Message{
				Id:      "MessageId - 1",
				Payload: "This is the message",
			})

			recievedMessage := Message{}
			Eventually(publisher.Messages).Should(Receive(&recievedMessage))
			Expect(recievedMessage.Payload).To(Equal("This is the message"))
		})

		It("doesn't send notifications after stopping publishing", func() {
			publisher := testhelpers.NewTestPublisher()
			publisher.Responses <- true

			queue := NewGoQ(25, publisher)
			queue.StartPublishing()

			queue.Enqueue(Message{Id: "MessageId - 1"})

			recievedMessage := Message{}
			Eventually(publisher.Messages).Should(Receive(&recievedMessage))
			Expect(recievedMessage.Id).To(Equal("MessageId - 1"))

			queue.StopPublishing()

			queue.Enqueue(Message{Id: "MessageId - 2"})
			Consistently(publisher.Messages).ShouldNot(Receive())
		})

		It("retries message if delivery fails", func() {
			publisher := testhelpers.NewTestPublisher()
			publisher.Responses <- false
			publisher.Responses <- true

			queue := NewGoQ(25, publisher)
			queue.StartPublishing()

			queue.Enqueue(Message{Id: "MessageId - 1"})

			Eventually(func() int {
				return len(publisher.Messages)
			}).Should(Equal(2))
		})
	})
})
