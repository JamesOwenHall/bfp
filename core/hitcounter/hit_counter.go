package hitcounter

import (
	"github.com/JamesOwenHall/BruteForceProtection/core/server"
	"time"
)

type HitCounter struct {
	clock *Clock
	*server.Server
}

func NewHitCounter(directions []Direction) *HitCounter {
	result := new(HitCounter)
	result.clock = NewClock()
	result.Server = server.New()

	for _, dir := range directions {
		// Add the route
		result.Routes[dir.Name()] = func(val interface{}) bool {
			return dir.Hit(result.clock.GetTime(), val)
		}

		// Schedule the cleanup
		go func(dir Direction) {
			for {
				dir.CleanUp(result.clock.GetTime())
				time.Sleep(5 * time.Second)
			}
		}(dir)
	}

	return result
}
