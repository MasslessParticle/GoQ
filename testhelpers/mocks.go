package testhelpers

import (
	"github.com/masslessparticle/goq"
	"github.com/masslessparticle/goq/pubsub"
)

type TestClient struct {
	ClientID      string
	Notifications chan goq.Message
}

func NewTestClient(id string) TestClient {
	return TestClient{
		ClientID:      id,
		Notifications: make(chan goq.Message, 1000),
	}
}

func (qc TestClient) ID() string {
	return qc.ClientID
}

func (qc TestClient) Notify(message goq.Message) {
	qc.Notifications <- message
}

type TestPublisher struct {
	Responses   chan bool
	Messages    chan goq.Message
	DoneCalls   chan bool
	Subscribers chan *pubsub.SubscriberList
}

func NewTestPublisher() *TestPublisher {
	return &TestPublisher{
		Responses: make(chan bool, 1000),
		Messages:  make(chan goq.Message, 1000),
		DoneCalls: make(chan bool, 1000),
	}
}

func (tp *TestPublisher) Publish(msg goq.Message) bool {
	tp.Messages <- msg

	if len(tp.Responses) == 0 {
		panic("responses must be set on the test publisher")
	}

	return <-tp.Responses
}

func (tp *TestPublisher) Done() {
	tp.DoneCalls <- true
}

func (tp *TestPublisher) Subscribe(client goq.QClient) error {
	return nil
}

func (tp *TestPublisher) Unsubscribe(qClient goq.QClient) {}

func (tp *TestPublisher) SubscriberCount() int {
	return 0
}
