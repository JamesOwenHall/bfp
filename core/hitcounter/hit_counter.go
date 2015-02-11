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
		result.Routes[dir.Name()] = makeRoute(result, dir)

		// Schedule the cleanup
		go func(dir Direction) {
			for {
				dir.CleanUp(result.clock.GetTime())
				time.Sleep(time.Duration(dir.CleanUpTime()) * time.Second)
			}
		}(dir)
	}

	return result
}

func makeRoute(hitCounter *HitCounter, dir Direction) func(interface{}) bool {
	return func(val interface{}) bool {
		return dir.Hit(hitCounter.clock.GetTime(), val)
	}
}
