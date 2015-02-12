package hitcounter

import (
	"hash/fnv"
	"sync"
)

const NumShards = 256

type BlockStatus struct {
	IsBlocked bool
	FrontTile float64
}

// StringMap

type StringMap struct {
	Mutexes    []sync.Mutex
	Shards     []map[string]*BlockStatus
	MaxTracked int64
}

func NewStringMap(maxTracked int64) *StringMap {
	result := new(StringMap)
	result.MaxTracked = maxTracked
	result.Mutexes = make([]sync.Mutex, NumShards)
	result.Shards = make([]map[string]*BlockStatus, NumShards)
	for i := 0; i < NumShards; i++ {
		result.Shards[i] = make(map[string]*BlockStatus)
	}

	return result
}

func (s *StringMap) Lock(key string) *BlockStatus {
	// Find the index of the shard
	hash := fnv.New32()
	hash.Write([]byte(key))
	index := int(hash.Sum32() % NumShards)

	s.Mutexes[index].Lock()
	shard := s.Shards[index]

	// Make sure we don't track too many values
	if s.MaxTracked != 0 {
		totalApprox := NumShards * int64(len(shard))
		if totalApprox > s.MaxTracked {
			return nil
		}
	}

	result, ok := shard[key]
	if !ok {
		result = new(BlockStatus)
		shard[key] = result
	}

	return result
}

func (s *StringMap) Unlock(key string) {
	hash := fnv.New32()
	hash.Write([]byte(key))
	index := int(hash.Sum32() % NumShards)

	s.Mutexes[index].Unlock()
}

// Int32Map

type Int32Map struct {
	Mutexes    []sync.Mutex
	Shards     []map[int32]*BlockStatus
	MaxTracked int64
}

func NewInt32Map(maxTracked int64) *Int32Map {
	result := new(Int32Map)
	result.MaxTracked = maxTracked
	result.Mutexes = make([]sync.Mutex, NumShards)
	result.Shards = make([]map[int32]*BlockStatus, NumShards)
	for i := 0; i < NumShards; i++ {
		result.Shards[i] = make(map[int32]*BlockStatus)
	}

	return result
}

func (i *Int32Map) Lock(key int32) *BlockStatus {
	index := key % NumShards
	i.Mutexes[index].Lock()
	shard := i.Shards[index]

	// Make sure we don't track too many values
	if i.MaxTracked != 0 {
		totalApprox := NumShards * int64(len(shard))
		if totalApprox > i.MaxTracked {
			return nil
		}
	}

	result, ok := shard[key]
	if !ok {
		result = new(BlockStatus)
		shard[key] = result
	}

	return result
}

func (i *Int32Map) Unlock(key int32) {
	index := key % NumShards
	i.Mutexes[index].Unlock()
}
