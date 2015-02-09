package hitcounter

type StringDirection struct {
	hits       *StringMap
	name       string
	windowSize int32
	incAmount  int32
}

func NewStringDirection(name string, windowSize, maxHits int32) *StringDirection {
	return &StringDirection{
		hits:       NewStringMap(),
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
		return false
	}

	// We need to use the lock to access the hits map.
	status := s.hits.Lock(value)
	defer s.hits.Unlock(value)

	if status.FrontTile < clock {
		// No recent hits
		status.FrontTile = clock + s.incAmount
		status.IsBlocked = false
		return true
	} else if status.FrontTile < clock+s.windowSize {
		// Recent hits, but not over the threshold
		status.FrontTile += s.incAmount
		return !status.IsBlocked
	} else {
		// Over the threshold
		status.IsBlocked = true
		return false
	}
}

func (s *StringDirection) CleanUp(clock int32) {
	for i, m := range s.hits.Mutexes {
		m.Lock()

		for k := range s.hits.Shards[i] {
			if s.hits.Shards[i][k].FrontTile < clock {
				delete(s.hits.Shards[i], k)
			}
		}

		m.Unlock()
	}
}
