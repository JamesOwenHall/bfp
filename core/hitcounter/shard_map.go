package hitcounter

import (
	"sync"
)

const NumShards = 128

type ShardMap interface {
	Get(key interface{}) (*BlockStatus, *sync.Mutex)
	CleanUp(clock int32)
	Type() string
	BlockedValues() []BlockedValue
}

type BlockedValue struct {
	Since float64
	Value interface{}
}
