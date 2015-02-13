package hitcounter

type Direction struct {
	Store       ShardMap
	Name        string
	CleanUpTime float64
	MaxHits     float64
	WindowSize  float64
}

func (d *Direction) Hit(clock int32, val interface{}) bool {
	status, mutex := d.Store.Get(val)
	if status == nil {
		return false
	}
	defer mutex.Unlock()

	incAmount := (d.WindowSize / d.MaxHits)

	fClock := float64(clock)
	if status.FrontTile < fClock {
		// No recent hits
		status.FrontTile = fClock + incAmount
		status.IsBlocked = false
		return true
	} else if status.FrontTile > fClock+incAmount {
		// We've crossed the threshold, start blocking
		status.IsBlocked = true
		return false
	} else {
		// We haven't crossed the threshold yet, let's increment
		status.FrontTile += incAmount

		// Now that we've incremented, we may have crossed the threshold
		if status.FrontTile > fClock+d.WindowSize {
			// We crossed the threshold, start blocking
			status.IsBlocked = true
			return false
		} else {
			// We're not over the threshold even after incrementing.  But it's
			// possible that we crossed it earlier, so let's make sure we're
			// not already blocking.
			return !status.IsBlocked
		}
	}
}
