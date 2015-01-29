package hitcounter

import (
	"sync"
)

type StringDirection struct {
	hits       map[string]int32
	name       string
	lock       sync.Mutex
	windowSize int32
	incAmount  int32
}

func NewStringDirection(name string, windowSize, maxHits int32) *StringDirection {
	return &StringDirection{
		hits:       make(map[string]int32),
		name:       name,
		windowSize: windowSize,
		incAmount:  windowSize / maxHits,
	}
}

func (s *StringDirection) Name() string {
	return s.name
}

func (s *StringDirection) Hit(clock int32, val interface{}) bool {
	// Make sure we have a string
	value, ok := val.(string)
	if !ok {
		return true
	}

	// We need to use the lock to access the hits map.
	s.lock.Lock()
	defer s.lock.Unlock()

	num, ok := s.hits[value]
	if !ok || num < clock {
		// No recent hits
		s.hits[value] = clock + s.incAmount
		return true
	} else if num < clock+s.windowSize {
		// Recent hits, but not over the threshold
		s.hits[value] = num + s.incAmount
		return true
	} else {
		// Over the threshold
		return false
	}
}
