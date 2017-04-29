package goq

import (
	"sync"
	"errors"
)

type QClient interface {
	Id() string
	Notify(message Message)
}

type Message struct {
	Id      string
	Payload string
}

type GoQ struct {
	maxDepth    int
	queue       chan Message
	subscribers subscribers
	doneChan    chan bool
	lock        sync.Mutex
}

func NewGoQ(queueDepth int) *GoQ {
	return &GoQ{
		maxDepth: queueDepth,
		queue: make(chan Message, queueDepth),
		subscribers: newSubscribersList(),
		doneChan: make(chan bool, 1),
	}
}

func (q *GoQ) QueuedMessages() int {
	return len(q.queue)
}

func (q *GoQ) Enqueue(message Message) error {
	if len(q.queue) < q.maxDepth {
		q.queue <- message
		return nil
	}

	return errors.New("Message rejected, max queue depth reached")
}

func (q *GoQ) Subscribe(client QClient) error {
	if q.subscribers.contains(client) {
		return errors.New("Duplicate Clinet Id")
	}

	q.subscribers.append(client)
	return nil
}

func (q *GoQ) Unsubscribe(qClient QClient) {
	q.subscribers.remove(qClient)
}

func (q *GoQ) IsSubscribed(qClient QClient) bool {
	return q.subscribers.contains(qClient)
}

func (q *GoQ) StartPublishing() {
	//TODO: Here is where to implement notification strategies
	go func() {
		nextNotified := -1
		for msg := range q.queue {
			select {
			case <-q.doneChan:
				return
			default:
				nextNotified = q.notifyNext(nextNotified, msg)
			}
		}
	}()
}

func (q *GoQ) StopPublishing() {
	q.doneChan <- true
}

func (q *GoQ) notifyNext(nextIndex int, msg Message) int {
	numSubScribers := q.subscribers.size()

	if numSubScribers > 0 {
		nextIndex = (nextIndex + 1) % numSubScribers
		q.subscribers.get(nextIndex).Notify(msg)
	} else {
		q.queue <- msg
	}

	return nextIndex
}
