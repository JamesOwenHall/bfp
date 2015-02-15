package hitcounter

import (
	"sync"
	"time"
)

type Counter struct {
	counts    []uint64
	numCounts int
	scale     int
	lock      sync.Mutex
	index     int
}

func NewCounter(n, scale int) *Counter {
	return &Counter{
		counts:    make([]uint64, n),
		numCounts: n,
		scale:     scale,
	}
}

func (c *Counter) Hit() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.counts[c.index]++
}

func (c *Counter) start() {
	go func() {
		ticker := time.Tick(time.Duration(c.scale) * time.Second)
		for range ticker {
			c.lock.Lock()
			c.index = (c.index + 1) % c.numCounts
			c.counts[c.index] = 0
			c.lock.Unlock()
		}
	}()
}

func (c *Counter) Read() []uint64 {
	result := make([]uint64, c.numCounts-1)

	c.lock.Lock()
	defer c.lock.Unlock()

	for iResult := range result {
		iCounts := (1 + c.index + iResult) % c.numCounts
		result[iResult] = c.counts[iCounts]
	}

	return result
}

type History struct {
	Counters []*Counter
}

func DefaultHistory() History {
	result := History{
		Counters: []*Counter{
			NewCounter(60, 1),
			NewCounter(144, 300),
		},
	}

	for i := range result.Counters {
		result.Counters[i].start()
	}

	return result
}

func (h *History) Hit() {
	for i := range h.Counters {
		h.Counters[i].Hit()
	}
}
