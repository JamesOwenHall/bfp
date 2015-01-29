package hitcounter

import (
	"sync"
	"time"
)

type Clock struct {
	lock   sync.Mutex
	ticks  int32
	ticker <-chan time.Time
}

func NewClock() *Clock {
	result := &Clock{
		ticker: time.Tick(time.Second),
	}

	// Start a goroutine to capture the ticks
	go func() {
		for {
			<-result.ticker
			result.lock.Lock()
			result.ticks++
			result.lock.Unlock()
		}
	}()

	return result
}

func (c *Clock) GetTime() int32 {
	var result int32
	c.lock.Lock()
	result = c.ticks
	c.lock.Unlock()
	return result
}
