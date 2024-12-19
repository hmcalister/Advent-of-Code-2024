package main

import "log/slog"

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

type GuardState struct {
	Coordinate Coordinate
	Direction  Direction
}

func (state GuardState) Step() GuardState {
	return GuardState{
		Coordinate: state.Coordinate.Step(state.Direction),
		Direction:  state.Direction,
	}
}

func (state GuardState) EncounterObstacle() GuardState {
	return GuardState{
		Coordinate: state.Coordinate,
		Direction:  state.Direction.RotateRight(),
	}
}

func (state GuardState) InBounds(mapWidth, mapHeight int) bool {
	return state.Coordinate.X >= 0 &&
		state.Coordinate.X < mapWidth &&
		state.Coordinate.Y >= 0 &&
		state.Coordinate.Y < mapHeight
}

type Coordinate struct {
	X int
	Y int
}

func (c Coordinate) Step(d Direction) Coordinate {
	directionCoord := directionMap[d]
	return Coordinate{
		c.X + directionCoord.X,
		c.Y + directionCoord.Y,
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
