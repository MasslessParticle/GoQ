package pub

import (
	"errors"
	"sync"
	"github.com/masslessparticle/goq"
)

type SubscriberList struct {
	sync.RWMutex
	items []goq.QClient
}

func NewSubscribersList() *SubscriberList {
	return &SubscriberList{
		items: make([]goq.QClient, 0),
	}
}

func (s *SubscriberList) Subscribe(client goq.QClient) error {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	if s.indexOf(client.Id()) >= 0 {
		return errors.New("Duplicate client ids aren't allowed")
	}

	s.items = append(s.items, client)
	return nil
}

func (s *SubscriberList) Unsubscribe(client goq.QClient) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	subIndex := s.indexOf(client.Id())
	if subIndex >= 0 {
		s.items = append(s.items[:subIndex], s.items[subIndex+1:]...)
	}
}

func (s *SubscriberList) IsSubscribed(client goq.QClient) bool {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.indexOf(client.Id()) >= 0
}

func (s *SubscriberList) Get(index int) goq.QClient {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.items[index]
}

func (s *SubscriberList) Size() int {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return len(s.items)
}

func (s *SubscriberList) indexOf(qClientId string) int {
	for i, item := range s.items {
		if item.Id() == qClientId {
			return i
		}
	}
	return -1
}
