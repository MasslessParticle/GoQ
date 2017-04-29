package goq

import "sync"

type subscribers struct {
	sync.RWMutex
	items []QClient
}

func newSubscribersList() subscribers {
	return subscribers{
		items: make([]QClient, 0),
	}
}

func (s *subscribers) append(client QClient) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	s.items = append(s.items, client)
}

func (s *subscribers) remove(client QClient) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	subIndex := s.indexOf(client.Id())
	if subIndex >= 0 {
		s.items = append(s.items[:subIndex], s.items[subIndex + 1:]...)
	}
}

func (s *subscribers) contains(client QClient) bool {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.indexOf(client.Id()) >= 0
}

func (s *subscribers) get(index int) QClient {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.items[index]
}

func (s *subscribers) size() int {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return len(s.items)
}

func (s *subscribers) indexOf(qClientId string) int {
	for i, item := range s.items {
		if item.Id() == qClientId {
			return i
		}
	}
	return -1
}
