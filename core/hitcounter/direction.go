package hitcounter

type Direction interface {
	CleanUp(clock int32)
	Hit(clock int32, val interface{}) bool
	Name() string
}
