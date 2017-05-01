package pubsub_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/masslessparticle/goq/pubsub"
	"github.com/masslessparticle/goq/testhelpers"
)

var _ = Describe("PriorityQueue", func() {
	var pq *pubsub.SubscriberPriorityQueue

	BeforeEach(func() {
		pq = pubsub.NewSubscriberPriorityQueue()
	})

	Context("Subscribe", func () {
		It("thows an error with duplicate client ids", func() {
			err := pq.Subscribe(testhelpers.NewTestClient("subscriber - 1"))
			Expect(err).ToNot(HaveOccurred())

			Expect(pq.SubscriberCount()).To(Equal(1))

			err = pq.Subscribe(testhelpers.NewTestClient("subscriber - 1"))
			Expect(err).To(HaveOccurred())

			Expect(pq.SubscriberCount()).To(Equal(1))
		})

		It("moves new subscribers to the top of the queue", func () {
			entry := pubsub.PQEntry{
				MessagesSent: 3,
				Client: testhelpers.NewTestClient("subscriber - 1"),
			}

			err := pq.Push(entry)
			Expect(err).ToNot(HaveOccurred())

			err = pq.Subscribe(testhelpers.NewTestClient("subscriber - 2"))
			Expect(err).ToNot(HaveOccurred())

			Expect(pq.Peek().Client.Id()).To(Equal("subscriber - 2"))
		})
	})

	Context("Push", func() {
		It("can't push duplicate clients", func() {
			entry := pubsub.PQEntry{
				MessagesSent: 2,
				Client: testhelpers.NewTestClient("subscriber - 1"),
			}

			err := pq.Push(entry)
			Expect(err).ToNot(HaveOccurred())

			Expect(pq.SubscriberCount()).To(Equal(1))

			err = pq.Push(entry)
			Expect(err).To(HaveOccurred())

			Expect(pq.SubscriberCount()).To(Equal(1))
		})

		It("allows adding elements with message counts", func() {
			entry := pubsub.PQEntry{
				MessagesSent: 2,
				Client: testhelpers.NewTestClient("subscriber - 1"),
			}

			pq.Push(entry)
			Expect(pq.SubscriberCount()).To(Equal(1))
		})
	})

	Context("Peek", func() {
		It("Returns the element with the lowest message calls", func() {
			entry := pubsub.PQEntry{
				MessagesSent: 4,
				Client: testhelpers.NewTestClient("subscriber - 1"),
			}

			entry2 := pubsub.PQEntry{
				MessagesSent: 3,
				Client: testhelpers.NewTestClient("subscriber - 2"),
			}

			entry3 := pubsub.PQEntry{
				MessagesSent: 2,
				Client: testhelpers.NewTestClient("subscriber - 3"),
			}


			err := pq.Push(entry)
			Expect(err).ToNot(HaveOccurred())
			err = pq.Push(entry2)
			Expect(err).ToNot(HaveOccurred())
			err = pq.Push(entry3)
			Expect(err).ToNot(HaveOccurred())

			Expect(pq.Peek().Client.Id()).To(Equal("subscriber - 3"))
		})
	})

	Context("Pop" , func() {
		It("returns the entry with lowest message calls and deletes it the queue remains prioritized", func() {
			entry := pubsub.PQEntry{
				MessagesSent: 4,
				Client: testhelpers.NewTestClient("subscriber - 1"),
			}

			entry2 := pubsub.PQEntry{
				MessagesSent: 3,
				Client: testhelpers.NewTestClient("subscriber - 2"),
			}

			entry3 := pubsub.PQEntry{
				MessagesSent: 2,
				Client: testhelpers.NewTestClient("subscriber - 3"),
			}

			entry4 := pubsub.PQEntry{
				MessagesSent: 6,
				Client: testhelpers.NewTestClient("subscriber - 4"),
			}


			err := pq.Push(entry)
			Expect(err).ToNot(HaveOccurred())
			err = pq.Push(entry2)
			Expect(err).ToNot(HaveOccurred())
			err = pq.Push(entry3)
			Expect(err).ToNot(HaveOccurred())
			err = pq.Push(entry4)
			Expect(err).ToNot(HaveOccurred())

			Expect(pq.Pop().Client.Id()).To(Equal("subscriber - 3"))
			Expect(pq.Peek().Client.Id()).To(Equal("subscriber - 2"))
		})
	})

	Context("Unsubscribe", func() {
		It("maintains the queue when something is removed from the middle", func() {
			deleteClient := testhelpers.NewTestClient("subscriber - 1")

			entry := pubsub.PQEntry{
				MessagesSent: 4,
				Client: deleteClient,
			}

			entry2 := pubsub.PQEntry{
				MessagesSent: 3,
				Client: testhelpers.NewTestClient("subscriber - 2"),
			}

			entry3 := pubsub.PQEntry{
				MessagesSent: 2,
				Client: testhelpers.NewTestClient("subscriber - 3"),
			}

			entry4 := pubsub.PQEntry{
				MessagesSent: 6,
				Client: testhelpers.NewTestClient("subscriber - 4"),
			}


			err := pq.Push(entry)
			Expect(err).ToNot(HaveOccurred())
			err = pq.Push(entry2)
			Expect(err).ToNot(HaveOccurred())
			err = pq.Push(entry3)
			Expect(err).ToNot(HaveOccurred())
			err = pq.Push(entry4)
			Expect(err).ToNot(HaveOccurred())

			pq.Unsubscribe(deleteClient)

			Expect(pq.Pop().Client.Id()).To(Equal("subscriber - 3"))
			Expect(pq.Pop().Client.Id()).To(Equal("subscriber - 2"))
			Expect(pq.Pop().Client.Id()).To(Equal("subscriber - 4"))
		})

		It("maintains the queue when something is removed from the middle", func() {
			deleteClient := testhelpers.NewTestClient("subscriber - 1")

			entry := pubsub.PQEntry{
				MessagesSent: 4,
				Client: deleteClient,
			}

			entry2 := pubsub.PQEntry{
				MessagesSent: 3,
				Client: testhelpers.NewTestClient("subscriber - 2"),
			}

			err := pq.Push(entry)
			Expect(err).ToNot(HaveOccurred())
			err = pq.Push(entry2)
			Expect(err).ToNot(HaveOccurred())

			pq.Unsubscribe(deleteClient)

			Expect(pq.Pop().Client.Id()).To(Equal("subscriber - 2"))
		})
	})
})
