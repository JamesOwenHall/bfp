package store

import (
	"sync"
)

type Int32Map struct {
	mutexes    []sync.Mutex
	shards     []map[int32]*BlockStatus
	maxTracked int64
}

func NewInt32Map(maxTracked int64) *Int32Map {
	result := new(Int32Map)
	result.maxTracked = maxTracked
	result.mutexes = make([]sync.Mutex, NumShards)
	result.shards = make([]map[int32]*BlockStatus, NumShards)
	for i := 0; i < NumShards; i++ {
		result.shards[i] = make(map[int32]*BlockStatus)
	}

	return result
}

func (i *Int32Map) Get(key interface{}) (*BlockStatus, *sync.Mutex) {
	// Type check
	keyFloat64, ok := key.(float64)
	if !ok {
		return nil, nil
	}
	keyInt32 := int32(keyFloat64)

	index := keyInt32 % NumShards

	// Locking the mutex
	mutex := &i.mutexes[index]
	mutex.Lock()
	shard := i.shards[index]

	// Make sure we don't track too many values
	if i.maxTracked != 0 {
		totalApprox := NumShards * int64(len(shard))
		if totalApprox > i.maxTracked {
			mutex.Unlock()
			return nil, nil
		}
	}

	// Create the status if it doesn't exist
	status, exists := shard[keyInt32]
	if !exists {
		status = new(BlockStatus)
		shard[keyInt32] = status
	}

	return status, mutex
}

func (i *Int32Map) CleanUp(clock int32) {
	fClock := float64(clock)
	for j := range i.mutexes {
		i.mutexes[j].Lock()

		for k := range i.shards[j] {
			if i.shards[j][k].FrontTile < fClock {
				delete(i.shards[j], k)
			}
		}

		i.mutexes[j].Unlock()
	}
}

func (i *Int32Map) Type() string {
	return "int32"
}

func (i *Int32Map) BlockedValues() []BlockedValue {
	result := make([]BlockedValue, 0)

	for j := range i.mutexes {
		shard := &i.shards[j]
		i.mutexes[j].Lock()

		for key, status := range *shard {
			if status.IsBlocked {
				result = append(
					result,
					BlockedValue{Since: status.Since, Value: key},
				)
			}
		}

		i.mutexes[j].Unlock()
	}

	return result
}
