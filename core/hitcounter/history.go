package hitcounter

import (
	"sync"
	"time"
)

type Counter struct {
	counts    []uint64
	NumCounts int
	Scale     int
	lock      sync.Mutex
	index     int
}

func NewCounter(n, scale int) *Counter {
	return &Counter{
		counts:    make([]uint64, n),
		NumCounts: n,
		Scale:     scale,
	}
}

func (c *Counter) Hit() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.counts[c.index]++
}

func (c *Counter) start() {
	go func() {
		ticker := time.Tick(time.Duration(c.Scale) * time.Second)
		for range ticker {
			c.lock.Lock()
			c.index = (c.index + 1) % c.NumCounts
			c.counts[c.index] = 0
			c.lock.Unlock()
		}
	}()
}

func (c *Counter) Read() []uint64 {
	result := make([]uint64, c.NumCounts)

	c.lock.Lock()
	defer c.lock.Unlock()

	for iResult := range result {
		iCounts := (c.index + iResult) % c.NumCounts
		result[iResult] = c.counts[iCounts]
	}

	return result
}

func (c *Counter) TimeRange() string {
	return (time.Duration(c.NumCounts*c.Scale) * time.Second).String()
}

type History struct {
	Short *Counter
	Long  *Counter
}

func DefaultHistory() History {
	result := History{
		Short: NewCounter(60, 1),
		Long:  NewCounter(144, 600),
	}

	result.Short.start()
	result.Long.start()

	return result
}

func (h *History) Hit() {
	h.Short.Hit()
	h.Long.Hit()
}
