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
	num := s.hits.Lock(value)
	defer s.hits.Unlock(value)

	if *num < clock {
		// No recent hits
		*num = clock + s.incAmount
		return true
	} else if *num < clock+s.windowSize {
		// Recent hits, but not over the threshold
		*num += s.incAmount
		return true
	} else {
		// Over the threshold
		return false
	}
}

func (s *StringDirection) CleanUp(clock int32) {
	for i, m := range s.hits.Mutexes {
		m.Lock()

		for k := range s.hits.Shards[i] {
			if *s.hits.Shards[i][k] < clock {
				delete(s.hits.Shards[i], k)
			}
		}

		m.Unlock()
	}
}
