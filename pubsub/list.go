package pubsub

import (
	"errors"
	"sync"

	"github.com/masslessparticle/goq"
)

type SubscriberList struct {
	lock sync.Mutex
	items []goq.QClient
}

func NewSubscribersList() *SubscriberList {
	return &SubscriberList{
		items: make([]goq.QClient, 0),
	}
}

func (s *SubscriberList) Subscribe(client goq.QClient) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.indexOf(client.ID()) >= 0 {
		return errors.New("Duplicate client ids aren't allowed")
	}

	s.items = append(s.items, client)
	return nil
}

func (s *SubscriberList) Unsubscribe(client goq.QClient) {
	s.lock.Lock()
	defer s.lock.Unlock()

	subIndex := s.indexOf(client.ID())
	if subIndex >= 0 {
		s.items = append(s.items[:subIndex], s.items[subIndex+1:]...)
	}
}

func (s *SubscriberList) Get(index int) goq.QClient {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.items[index]
}

func (s *SubscriberList) SubscriberCount() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.items)
}

func (s *SubscriberList) indexOf(qClientID string) int {
	for i, item := range s.items {
		if item.ID() == qClientID {
			return i
		}
	}
	return -1
}
