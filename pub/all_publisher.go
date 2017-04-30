package pub

import "github.com/masslessparticle/goq"

type AllPublisher struct {
	SubscriberList
}

func NewAllPublisher() *AllPublisher {
	publisher := AllPublisher{}
	publisher.items = make([]goq.QClient, 0)

	return &publisher
}

func (rr *AllPublisher) Publish(msg goq.Message) bool {
	numSubScribers := rr.SubscriberCount()
	delivered := false

	for i:= 0; i < numSubScribers; i++ {
		rr.Get(i).Notify(msg)
		delivered = true
	}

	return delivered
}