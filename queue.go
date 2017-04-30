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
	Publish(msg Message) bool
    Subscribe(client QClient) error
    Unsubscribe(qClient QClient)
    IsSubscribed(qClient QClient) bool
}

type Message struct {
	Id      string
	Payload string
}

type GoQ struct {
	Publisher Publisher
	maxDepth  int
	queue     chan Message
	doneChan  chan bool
	lock      sync.Mutex
}

func NewGoQ(queueDepth int, publisher Publisher) *GoQ {
	return &GoQ{
		maxDepth:    queueDepth,
		queue:       make(chan Message, queueDepth),
		doneChan:    make(chan bool, 1),
		Publisher:         publisher,
	}
}

func (q *GoQ) Enqueue(message Message) error {
	if len(q.queue) < q.maxDepth {
		q.queue <- message
		return nil
	}

	return errors.New("Message rejected, max queue depth reached")
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
	delivered := q.Publisher.Publish(msg)
	if !delivered {
		q.queue <- msg
	}
}

func (q *GoQ) StopPublishing() {
	q.doneChan <- true
}
