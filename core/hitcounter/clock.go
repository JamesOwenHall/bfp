package hitcounter

import (
	"sync"
	"time"
)

// Clock represents an integer clock that starts at 0.
type Clock struct {
	lock   sync.RWMutex
	ticks  int32
	ticker <-chan time.Time
}

// NewClock returns an initialized and started *Clock.
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

// GetTime returns the number of seconds that have passed since intiialization.
func (c *Clock) GetTime() int32 {
	var result int32
	c.lock.RLock()
	result = c.ticks
	c.lock.RUnlock()
	return result
}
