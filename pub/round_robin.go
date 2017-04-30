package pub

import (
	"github.com/masslessparticle/goq"
)

type RoundRobinPublisher struct {
	nextNotified   int
	subscribers *SubscriberList
}

func NewRoundRobinPublisher() *RoundRobinPublisher {
	return &RoundRobinPublisher{
		nextNotified: -1,
		subscribers: NewSubscribersList(),
	}
}

func (rr *RoundRobinPublisher) Publish(msg goq.Message) bool {
	if rr.subscribers == nil {
		return false
	}

	numSubScribers := rr.subscribers.Size()
	if numSubScribers > 0 {
		rr.nextNotified = (rr.nextNotified + 1) % numSubScribers
		rr.subscribers.Get(rr.nextNotified).Notify(msg)
	} else {
		return false
	}

	return true
}

func (rr *RoundRobinPublisher) Subscribe(client goq.QClient) error {
	return rr.subscribers.Append(client)
}

func (rr *RoundRobinPublisher) Unsubscribe(qClient goq.QClient) {
	rr.subscribers.Remove(qClient)
}

func (rr *RoundRobinPublisher) IsSubscribed(qClient goq.QClient) bool {
	return rr.subscribers.Contains(qClient)
}
