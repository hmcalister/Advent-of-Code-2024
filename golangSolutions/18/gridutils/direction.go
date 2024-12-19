package gridutils

//go:generate stringer --type Direction
type Direction int

const (
	DIRECTION_UP    Direction = 0
	DIRECTION_RIGHT Direction = 1
	DIRECTION_DOWN  Direction = 2
	DIRECTION_LEFT  Direction = 3
)

var (
	directionMap = []Coordinate{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
)
