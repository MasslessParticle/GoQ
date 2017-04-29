package testhelpers

import (
	"github.com/masslessparticle/goq"
)

type QClient struct {
	ClientId      string
	Notifications chan goq.Message
}

func NewTestClient(id string) QClient {
	return QClient{
		ClientId: id,
		Notifications: make(chan goq.Message, 5000000),
	}
}

func (qc QClient) Id() string {
	return qc.ClientId
}

func (qc QClient) Notify(message goq.Message) {
	qc.Notifications <- message
}



