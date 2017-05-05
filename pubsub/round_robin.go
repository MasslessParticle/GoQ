package pubsub

import (
	"github.com/masslessparticle/goq"
)

type RoundRobinPublisher struct {
	SubscriberList
	nextNotified int
}

func NewRoundRobinPublisher() *RoundRobinPublisher {
	publisher := RoundRobinPublisher{
		nextNotified: -1,
	}

	publisher.items = make([]goq.QClient, 0)

	return &publisher
}

func (rr *RoundRobinPublisher) Publish(msg goq.Message) bool {
	numSubScribers := rr.SubscriberCount()
	if numSubScribers > 0 {
		rr.nextNotified = (rr.nextNotified + 1) % numSubScribers
		rr.Get(rr.nextNotified).Notify(msg)
	} else {
		return false
	}

	return true
}

func (rr *RoundRobinPublisher) Done() {}