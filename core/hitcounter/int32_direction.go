package hitcounter

import (
	"log"
	"sync"
)

type Int32Direction struct {
	hits       map[int32]int32
	name       string
	lock       sync.Mutex
	windowSize int32
	incAmount  int32
}

func NewInt32Direction(name string, windowSize, maxHits int32) *Int32Direction {
	return &Int32Direction{
		hits:       make(map[int32]int32),
		name:       name,
		windowSize: windowSize,
		incAmount:  windowSize / maxHits,
	}
}

func (i *Int32Direction) Name() string {
	return i.name
}

func (i *Int32Direction) Hit(clock int32, val interface{}) bool {
	// JSON encodes numbers as float64, so we need to type assert, then cast.
	valueFloat, ok := val.(float64)
	if !ok {
		return true
	}
	value := int32(valueFloat)

	// We need to use the lock to access the hits map.
	i.lock.Lock()
	defer i.lock.Unlock()

	num, ok := i.hits[value]
	log.Println("num =", num, "clock =", clock)

	if !ok || num < clock {
		// No recent hits
		i.hits[value] = clock + i.incAmount
		return true
	} else if num < clock+i.windowSize {
		// Recent hits, but not over the threshold
		i.hits[value] = num + i.incAmount
		return true
	} else {
		// Over the threshold
		return false
	}
}
