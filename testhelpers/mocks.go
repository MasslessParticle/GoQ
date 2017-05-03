package testhelpers

import (
	"github.com/masslessparticle/goq"
	"github.com/masslessparticle/goq/pubsub"
)

type TestClient struct {
	ClientId      string
	Notifications chan goq.Message
}

func NewTestClient(id string) TestClient {
	return TestClient{
		ClientId:      id,
		Notifications: make(chan goq.Message, 1000),
	}
}

func (qc TestClient) Id() string {
	return qc.ClientId
}

func (qc TestClient) Notify(message goq.Message) {
	qc.Notifications <- message
}

type TestPublisher struct {
	Responses   chan bool
	Messages    chan goq.Message
	Subscribers chan *pubsub.SubscriberList
}

func NewTestPublisher() *TestPublisher {
	return &TestPublisher{
		Responses: make(chan bool, 1000),
		Messages:  make(chan goq.Message, 1000),
	}
}

func (tp *TestPublisher) Publish(msg goq.Message) bool {
	tp.Messages <- msg

	if len(tp.Responses) == 0 {
		panic("responses must be set on the test publisher")
	}

	return <-tp.Responses
}

func (tp *TestPublisher) Subscribe(client goq.QClient) error {
	return nil
}

func (tp *TestPublisher) Unsubscribe(qClient goq.QClient) {}

func (tp *TestPublisher) SubscriberCount() int {
	return 0
}
