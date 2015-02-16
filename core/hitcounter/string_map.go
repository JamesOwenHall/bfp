package hitcounter

import (
	"hash/fnv"
	"sync"
)

type StringMap struct {
	mutexes    []sync.Mutex
	shards     []map[string]*BlockStatus
	maxTracked int64
}

func NewStringMap(maxTracked int64) *StringMap {
	result := new(StringMap)
	result.maxTracked = maxTracked
	result.mutexes = make([]sync.Mutex, NumShards)
	result.shards = make([]map[string]*BlockStatus, NumShards)
	for i := 0; i < NumShards; i++ {
		result.shards[i] = make(map[string]*BlockStatus)
	}

	return result
}

func (s *StringMap) Get(key interface{}) (*BlockStatus, *sync.Mutex) {
	// Type check
	keyString, ok := key.(string)
	if !ok {
		return nil, nil
	}

	// Find the index of the shard
	hash := fnv.New32()
	hash.Write([]byte(keyString))
	index := int(hash.Sum32() % NumShards)

	mutex := &s.mutexes[index]
	mutex.Lock()
	shard := s.shards[index]

	// Make sure we don't track too many values
	if s.maxTracked != 0 {
		totalApprox := NumShards * int64(len(shard))
		if totalApprox > s.maxTracked {
			return nil, nil
		}
	}

	// Create a status if it doesn't exist
	status, exists := shard[keyString]
	if !exists {
		status = new(BlockStatus)
		shard[keyString] = status
	}

	return status, mutex
}

func (s *StringMap) CleanUp(clock int32) {
	fClock := float64(clock)
	for i := range s.mutexes {
		m := &s.mutexes[i]
		m.Lock()

		for k := range s.shards[i] {
			if s.shards[i][k].FrontTile < fClock {
				delete(s.shards[i], k)
			}
		}

		m.Unlock()
	}
}

func (s *StringMap) Type() string {
	return "string"
}

func (s *StringMap) BlockedValues() []interface{} {
	result := make([]interface{}, 0)

	for i := range s.mutexes {
		mutex := &s.mutexes[i]
		shard := &s.shards[i]
		mutex.Lock()

		for key, status := range *shard {
			if status.IsBlocked {
				result = append(result, key)
			}
		}

		mutex.Unlock()
	}

	return result
}
