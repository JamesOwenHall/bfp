package hitcounter

type Int32Direction struct {
	hits       *Int32Map
	name       string
	windowSize int32
	incAmount  int32
}

func NewInt32Direction(name string, windowSize, maxHits int32) *Int32Direction {
	return &Int32Direction{
		hits:       NewInt32Map(),
		name:       name,
		windowSize: windowSize,
		incAmount:  windowSize / maxHits,
	}
}

func (i *Int32Direction) Name() string {
	return i.name
}

func (i *Int32Direction) Hit(clock int32, val interface{}) bool {
	// JSON encodes numbers as float64, so we need to type assert, then cast.
	valueFloat, ok := val.(float64)
	if !ok {
		return true
	}
	value := int32(valueFloat)

	// We need to use the lock to access the hits map.
	num := i.hits.Lock(value)
	defer i.hits.Unlock(value)

	if *num < clock {
		// No recent hits
		*num = clock + i.incAmount
		return true
	} else if *num < clock+i.windowSize {
		// Recent hits, but not over the threshold
		*num += i.incAmount
		return true
	} else {
		// Over the threshold
		return false
	}
}

func (i *Int32Direction) CleanUp(clock int32) {
	for j, m := range i.hits.Mutexes {
		m.Lock()

		for k := range i.hits.Shards[j] {
			if *i.hits.Shards[j][k] < clock {
				delete(i.hits.Shards[j], k)
			}
		}

		m.Unlock()
	}
}
