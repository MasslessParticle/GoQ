package goq

import (
	"errors"
	"sync"
)

type Subscribers struct {
	sync.RWMutex
	items []QClient
}

func NewSubscribersList() Subscribers {
	return Subscribers{
		items: make([]QClient, 0),
	}
}

func (s *Subscribers) Append(client QClient) error {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	if s.indexOf(client.Id()) >= 0 {
		return errors.New("Duplicate client ids aren't allowed")
	}

	s.items = append(s.items, client)
	return nil
}

func (s *Subscribers) Remove(client QClient) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	subIndex := s.indexOf(client.Id())
	if subIndex >= 0 {
		s.items = append(s.items[:subIndex], s.items[subIndex+1:]...)
	}
}

func (s *Subscribers) Contains(client QClient) bool {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.indexOf(client.Id()) >= 0
}

func (s *Subscribers) Get(index int) QClient {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.items[index]
}

func (s *Subscribers) Size() int {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return len(s.items)
}

func (s *Subscribers) indexOf(qClientId string) int {
	for i, item := range s.items {
		if item.Id() == qClientId {
			return i
		}
	}
	return -1
}
