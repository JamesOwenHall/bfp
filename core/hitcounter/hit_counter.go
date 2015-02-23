package hitcounter

import (
	"github.com/JamesOwenHall/BruteForceProtection/core/message-server"
	"time"
)

type HitCounter struct {
	Clock *Clock
	*server.Server
}

func NewHitCounter(directions []Direction) *HitCounter {
	result := new(HitCounter)
	result.Clock = NewClock()
	result.Server = server.New()

	for i := range directions {
		// Add the route
		result.Routes[directions[i].Name] = makeRoute(result, &directions[i])

		// Schedule the cleanup
		go func(dir *Direction) {
			for {
				dir.Store.CleanUp(result.Clock.GetTime())
				time.Sleep(time.Duration(dir.CleanUpTime) * time.Second)
			}
		}(&directions[i])
	}

	return result
}

func makeRoute(hitCounter *HitCounter, dir *Direction) func(interface{}) bool {
	return func(val interface{}) bool {
		return dir.Hit(hitCounter.Clock.GetTime(), val)
	}
}
