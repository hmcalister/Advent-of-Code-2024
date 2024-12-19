package gridutils

import "log/slog"

type Coordinate struct {
	X int
	Y int
}

func (c Coordinate) Equal(otherCoord Coordinate) bool {
	return c.X == otherCoord.X && c.Y == otherCoord.Y
}

func (c Coordinate) Step(d Direction) Coordinate {
	directionCoord := directionMap[d]
	return Coordinate{
		c.X + directionCoord.X,
		c.Y + directionCoord.Y,
	}
}

func (c Coordinate) GetOrthogonalNeighbors() []Coordinate {
	return []Coordinate{
		{c.X - 1, c.Y},
		{c.X + 1, c.Y},
		{c.X, c.Y - 1},
		{c.X, c.Y + 1},
	}
}

func (d Direction) RotateLeft() Direction {
	switch d {
	case DIRECTION_UP:
		return DIRECTION_LEFT
	case DIRECTION_RIGHT:
		return DIRECTION_UP
	case DIRECTION_DOWN:
		return DIRECTION_RIGHT
	case DIRECTION_LEFT:
		return DIRECTION_DOWN
	default:
		slog.Error("unexpected direction", "direction", d)
		return DIRECTION_UP
	}
}

func (d Direction) RotateRight() Direction {
	switch d {
	case DIRECTION_UP:
		return DIRECTION_RIGHT
	case DIRECTION_RIGHT:
		return DIRECTION_DOWN
	case DIRECTION_DOWN:
		return DIRECTION_LEFT
	case DIRECTION_LEFT:
		return DIRECTION_UP
	default:
		slog.Error("unexpected direction", "direction", d)
		return DIRECTION_UP
	}
}
