package store

type BlockStatus struct {
	IsBlocked bool
	FrontTile float64
	Since     float64
}

type BlockedValue struct {
	Since float64
	Value interface{}
}
