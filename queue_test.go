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

		It("Returns an error when a message is sent to a done queue", func() {
			publisher := testhelpers.NewTestPublisher()
			publisher.Responses <- true

			queue := NewGoQ(25, publisher)
			queue.StartPublishing()

			queue.Enqueue(Message{Id: "MessageId - 1"})

			recievedMessage := Message{}
			Eventually(publisher.Messages).Should(Receive(&recievedMessage))
			Expect(recievedMessage.Id).To(Equal("MessageId - 1"))

			queue.StopPublishing()

			err := queue.Enqueue(Message{Id: "MessageId - 2"})
			Expect(err).To(HaveOccurred())
		})


	})

	Context("Notifications", func() {
		It("passes to the pubsub whether or not the channel is done", func () {
			publisher := testhelpers.NewTestPublisher()
			publisher.Responses <- true

			queue := NewGoQ(25, publisher)
			queue.StartPublishing()

			queue.Enqueue(Message{Id: "MessageId - 1"})

			recievedMessage := Message{}
			Eventually(publisher.Messages).Should(Receive(&recievedMessage))
			Expect(recievedMessage.Id).To(Equal("MessageId - 1"))

			queue.StopPublishing()

			Eventually(publisher.DoneCalls).Should(Receive(Equal(true)))
		})

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

			queue.PausePublishing()

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
