package goq

import (
	"errors"
	"sync"
)

type QClient interface {
	Id() string
	Notify(message Message) error
}

type PubSub interface {
	Done()
	Publish(msg Message) bool
	Subscribe(client QClient) error
	Unsubscribe(qClient QClient)
	SubscriberCount() int
}

type Message struct {
	Id      string
	Payload interface{}
}

type GoQ struct {
	pubsub    PubSub
	maxDepth  int
	queue     chan Message
	pauseChan chan bool
	lock      sync.Mutex
	done      bool
}

func NewGoQ(queueDepth int, publisher PubSub) *GoQ {
	return &GoQ{
		maxDepth:  queueDepth,
		queue:     make(chan Message, queueDepth),
		pauseChan: make(chan bool, 1),
		pubsub:    publisher,
	}
}

func (q *GoQ) Enqueue(message Message) error {
	if q.done {
		return errors.New("Queue closed")
	}

	select {
	case q.queue <- message:
		return nil
	default:
		return errors.New("Message rejected, max queue depth reached")
	}
}

func (q *GoQ) StartPublishing() {
	go func() {
		for {
			msg, ok := <-q.queue
			if ok {
				select {
				case <-q.pauseChan:
					return
				default:
					q.publishMessage(msg)
				}
			} else {
				q.pubsub.Done()
				return
			}
		}
	}()
}

func (q *GoQ) StopPublishing() {
	q.lock.Lock()
	defer q.lock.Unlock()

	close(q.queue)
	q.done = true
}

func (q *GoQ) publishMessage(msg Message) {
	delivered := q.pubsub.Publish(msg)
	if !delivered {
		q.queue <- msg
	}
}

func (q *GoQ) PausePublishing() {
	q.pauseChan <- true
}
