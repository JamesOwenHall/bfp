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
	Mutexes []sync.Mutex
	Shards  []map[string]*BlockStatus
}

func NewStringMap() *StringMap {
	result := new(StringMap)

	result.Mutexes = make([]sync.Mutex, NumShards)
	result.Shards = make([]map[string]*BlockStatus, NumShards)
	for i := 0; i < NumShards; i++ {
		result.Shards[i] = make(map[string]*BlockStatus)
	}

	return result
}

func (s *StringMap) Lock(key string) *BlockStatus {
	hash := fnv.New32()
	hash.Write([]byte(key))
	index := int(hash.Sum32() % NumShards)

	s.Mutexes[index].Lock()
	shard := s.Shards[index]
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
	Mutexes []sync.Mutex
	Shards  []map[int32]*BlockStatus
}

func NewInt32Map() *Int32Map {
	result := new(Int32Map)

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
