package store

import (
	"hash/fnv"
	"sync"
)

const NumShards = 128

type ShardMap struct {
	Type       string
	mutexes    []sync.Mutex
	shards     []map[interface{}]*BlockStatus
	maxTracked int64
}

func NewShardMap(maxTracked int64) *ShardMap {
	result := new(ShardMap)
	result.maxTracked = maxTracked
	result.mutexes = make([]sync.Mutex, NumShards)
	result.shards = make([]map[interface{}]*BlockStatus, NumShards)
	for i := 0; i < NumShards; i++ {
		result.shards[i] = make(map[interface{}]*BlockStatus)
	}

	return result
}

func (s *ShardMap) Get(key interface{}) (*BlockStatus, *sync.Mutex) {
	index := s.Hash(key)
	if index == -1 {
		return nil, nil
	}

	mutex := &s.mutexes[index]
	mutex.Lock()
	shard := s.shards[index]

	// Make sure we don't track too many values
	if s.maxTracked != 0 {
		totalApprox := NumShards * int64(len(shard))
		if totalApprox > s.maxTracked {
			mutex.Unlock()
			return nil, nil
		}
	}

	// Create a status if it doesn't exist
	status, exists := shard[key]
	if !exists {
		status = new(BlockStatus)
		shard[key] = status
	}

	return status, mutex
}

func (s *ShardMap) Hash(v interface{}) int {
	switch s.Type {
	case "string":
		val, ok := v.(string)
		if !ok {
			return -1
		}

		hash := fnv.New32()
		hash.Write([]byte(val))
		return int(hash.Sum32() % NumShards)
	case "int32":
		// Values that come as numbers in JSON are of type float64.
		val, ok := v.(float64)
		if !ok {
			return -1
		}

		return int(int32(val) % NumShards)
	default:
		return -1
	}
}

func (s *ShardMap) CleanUp(clock int32) {
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

func (s *ShardMap) BlockedValues() []BlockedValue {
	result := make([]BlockedValue, 0)

	for i := range s.mutexes {
		mutex := &s.mutexes[i]
		shard := &s.shards[i]
		mutex.Lock()

		for key, status := range *shard {
			if status.IsBlocked {
				result = append(
					result,
					BlockedValue{Since: status.Since, Value: key},
				)
			}
		}

		mutex.Unlock()
	}

	return result
}
