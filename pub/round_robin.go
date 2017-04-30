package pub

import (
	"github.com/masslessparticle/goq"
)

type RoundRobinPublisher struct {
	nextNotified   int
}

func NewRoundRobinPublisher() *RoundRobinPublisher {
	return &RoundRobinPublisher{
		nextNotified: -1,
	}
}

func (rr *RoundRobinPublisher) Publish(msg goq.Message, subscribers *goq.Subscribers) bool {
	if subscribers == nil {
		return false
	}

	numSubScribers := subscribers.Size()
	if numSubScribers > 0 {
		rr.nextNotified = (rr.nextNotified + 1) % numSubScribers
		subscribers.Get(rr.nextNotified).Notify(msg)
	} else {
		return false
	}

	return true
}
