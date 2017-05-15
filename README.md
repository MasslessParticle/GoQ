# GoQ   [![Build Status](https://travis-ci.org/MasslessParticle/GoQ.svg)](https://travis-ci.org/MasslessParticle/GoQ) [![Go Report Card](https://goreportcard.com/badge/github.com/masslessparticle/goq)](https://goreportcard.com/report/github.com/masslessparticle/goq)

Package Goq provides a lightweight, extensible, in-memory message broker.

`go get github.com/masslessparticle/goq`

Running tests:

```
go get github.com/onsi/ginkgo
go get github.com/onsi/gomega
go install github.com/onsi/ginkgo/ginkgo

ginkgo -r
```

## QuickStart
```go
client := testhelpers.NewTestClient("Subscription - 1")

publisher := pubsub.NewRoundRobinPublisher()
publisher.Subscribe(client)

queue := goq.NewGoQ(25, publisher)
queue.StartPublishing()
queue.Enqueue(goq.Message{Id: "Message - 1"})
queue.StopPublishing()
```

## Creating a Queue

The goq.NewGoQ() method takes a max queue-depth and a PubSub
```go
queue := goq.NewGoQ(25, pubsub)
```

## QClient

A subscriber to the message broker. QClients are called according to the strategy provided by the PubSub. A QClient is anything implementing `goq.QClient`:

```go
type QClient interface {
	Id() string
	Notify(message Message) error
}
```

## PubSub

PubSub is the component that handles client subscription and message delivery. GoQ provides three message delivery strategies:
- Deliver to all clients.
- Round Robin.
- Least Used.

Anything implementing the `goq.PubSub` interface can be a PubSub:

```go
type PubSub interface {
	Done()
	Publish(msg Message) bool
	Subscribe(client QClient) error
	Unsubscribe(qClient QClient)
	SubscriberCount() int
}
```

## Message

The type received and emitted by the queue.

## Contributing

Pull requests, bug fixes and issue reports are welcome and appreciated.

## Licence

MIT
