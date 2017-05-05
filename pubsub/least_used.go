package pubsub

import "github.com/masslessparticle/goq"

type LeastUsedPublisher struct {
	SubscriberPriorityQueue
}

func NewLeastUsedPublisher() *LeastUsedPublisher {
	publisher := LeastUsedPublisher{}
	publisher.items = make([]PQEntry, 1)
	publisher.subscribedClients = make(map[string]bool, 0)
	return &publisher
}

func (pq *LeastUsedPublisher) Publish(msg goq.Message) bool {
	numSubscribers := pq.SubscriberCount()
	if numSubscribers > 0 {
		entry := pq.Pop()
		entry.Client.Notify(msg)
		entry.MessagesSent = entry.MessagesSent + 1
		pq.Push(entry)

		return true
	} else {
		return false
	}
}

func (pq *LeastUsedPublisher) Done() {}
