package pubsub

import (
	"errors"
	"github.com/masslessparticle/goq"
	"math"
	"sync"
)

type PQEntry struct {
	MessagesSent int
	Client       goq.QClient
}

type SubscriberPriorityQueue struct {
	lock              sync.RWMutex
	items             []PQEntry
	subscribedClients map[string]bool
}

func NewSubscriberPriorityQueue() *SubscriberPriorityQueue {
	return &SubscriberPriorityQueue{
		items:             make([]PQEntry, 1),
		subscribedClients: make(map[string]bool, 0),
	}
}

func (s *SubscriberPriorityQueue) Push(entry PQEntry) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, hasKey := s.subscribedClients[entry.Client.Id()]
	if hasKey {
		return errors.New("Duplicate client ids aren't allowed")
	}

	s.items = append(s.items, entry)
	s.bubbleUp(len(s.items) - 1)
	s.subscribedClients[entry.Client.Id()] = true

	return nil
}

func (s *SubscriberPriorityQueue) Peek() PQEntry {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.items[1]
}

func (s *SubscriberPriorityQueue) Pop() PQEntry {
	s.lock.Lock()
	defer s.lock.Unlock()

	item := s.items[1]

	s.swap(1, len(s.items)-1)
	s.items = s.items[:len(s.items)-1]

	s.bubbleDown(1)

	delete(s.subscribedClients, item.Client.Id())

	return item
}

func (s *SubscriberPriorityQueue) Subscribe(client goq.QClient) error {
	return s.Push(PQEntry{Client: client})
}

func (s *SubscriberPriorityQueue) Unsubscribe(client goq.QClient) {
	s.lock.Lock()
	defer s.lock.Unlock()

	subIndex := s.indexOf(client.Id())
	if subIndex >= 0 {
		item := s.items[subIndex]
		s.swap(subIndex, len(s.items)-1)
		s.items = s.items[:len(s.items)-1]

		if subIndex != len(s.items) {
			s.bubbleDown(subIndex)
			s.bubbleUp(subIndex)
		}

		delete(s.subscribedClients, item.Client.Id())
	}
}

func (s *SubscriberPriorityQueue) SubscriberCount() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.items) - 1
}

func (s *SubscriberPriorityQueue) indexOf(qClientId string) int {
	for i := 1; i < len(s.items); i++ {
		if s.items[i].Client.Id() == qClientId {
			return i
		}
	}
	return -1
}

func (s *SubscriberPriorityQueue) bubbleUp(clientIndex int) {
	for i := clientIndex; i > 0; i = i / 2 {
		s.swapIfSmallerParent(i/2, i)
	}
}

func (s *SubscriberPriorityQueue) swapIfSmallerParent(parent, current int) int {
	if current == 1 {
		return current
	}

	if s.items[parent].MessagesSent > s.items[current].MessagesSent {
		s.swap(parent, current)
		return parent
	}

	return current
}

func (s *SubscriberPriorityQueue) swap(toIndex, fromIndex int) {
	temp := s.items[toIndex]
	s.items[toIndex] = s.items[fromIndex]
	s.items[fromIndex] = temp
}

func (s *SubscriberPriorityQueue) bubbleDown(clientIndex int) {
	for i := clientIndex; i < len(s.items); {
		lastIndex := i
		i = s.swapIfLargerChild(i)

		if i == lastIndex {
			return
		}
	}
}

func (s *SubscriberPriorityQueue) swapIfLargerChild(nodeIndex int) int {
	lChild := nodeIndex * 2
	rChild := nodeIndex*2 + 1

	if lChild >= len(s.items) && rChild >= len(s.items) {
		return nodeIndex
	}

	lMessages := math.MaxInt32
	rMessages := math.MaxInt32

	if lChild < len(s.items) {
		lMessages = s.items[lChild].MessagesSent
	}

	if rChild < len(s.items) {
		rMessages = s.items[rChild].MessagesSent
	}

	minMessages := s.minInt(lMessages, rMessages)

	if s.items[nodeIndex].MessagesSent > minMessages {
		if minMessages == lMessages {
			s.swap(lChild, nodeIndex)
			return lChild
		} else {
			s.swap(rChild, nodeIndex)
			return rChild
		}
	}

	return nodeIndex
}

func (s *SubscriberPriorityQueue) minInt(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
