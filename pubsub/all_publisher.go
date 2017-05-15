package pubsub

import "github.com/masslessparticle/goq"

type AllPublisher struct {
	SubscriberList
	EmptyDone
}

func NewAllPublisher() *AllPublisher {
	publisher := AllPublisher{}
	publisher.items = make([]goq.QClient, 0)

	return &publisher
}

func (ap *AllPublisher) Publish(msg goq.Message) bool {
	numSubScribers := ap.SubscriberCount()
	delivered := false

	for i := 0; i < numSubScribers; i++ {
		ap.Get(i).Notify(msg)
		delivered = true
	}

	return delivered
}
