package goq

import (
	"errors"
	"sync"
)

type QClient interface {
	Id() string
	Notify(message Message)
}

type Publisher interface {
	Publish(msg Message, subscribers *Subscribers) bool
}

type Message struct {
	Id      string
	Payload string
}

type GoQ struct {
	maxDepth    int
	queue       chan Message
	subscribers Subscribers
	doneChan    chan bool
	lock        sync.Mutex
	pub         Publisher
}

func NewGoQ(queueDepth int, publisher Publisher) *GoQ {
	subscribers := NewSubscribersList()

	return &GoQ{
		maxDepth:    queueDepth,
		queue:       make(chan Message, queueDepth),
		subscribers: subscribers,
		doneChan:    make(chan bool, 1),
		pub:         publisher,
	}
}

func (q *GoQ) Enqueue(message Message) error {
	if len(q.queue) < q.maxDepth {
		q.queue <- message
		return nil
	}

	return errors.New("Message rejected, max queue depth reached")
}

func (q *GoQ) Subscribe(client QClient) error {
	return q.subscribers.Append(client)
}

func (q *GoQ) Unsubscribe(qClient QClient) {
	q.subscribers.Remove(qClient)
}

func (q *GoQ) IsSubscribed(qClient QClient) bool {
	return q.subscribers.Contains(qClient)
}

func (q *GoQ) StartPublishing() {
	go func() {
		for msg := range q.queue {
			select {
			case <-q.doneChan:
				return
			default:
				q.publishMessage(msg)
			}
		}
	}()
}

func (q *GoQ) publishMessage(msg Message) {
	delivered := q.pub.Publish(msg, &q.subscribers)
	if !delivered {
		q.queue <- msg
	}
}

func (q *GoQ) StopPublishing() {
	q.doneChan <- true
}
