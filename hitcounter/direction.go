package hitcounter

type Direction interface {
	Hit(clock int32, val interface{}) bool
	Name() string
}
