package main

import "log/slog"

type Direction Coordinate

var (
	DIRECTION_UP    Direction = Direction{0, -1}
	DIRECTION_RIGHT Direction = Direction{1, 0}
	DIRECTION_DOWN  Direction = Direction{0, 1}
	DIRECTION_LEFT  Direction = Direction{-1, 0}
)

type GuardState struct {
	C Coordinate
	D Direction
}

func (state GuardState) Step() GuardState {
	return GuardState{
		C: state.C.Step(state.D),
		D: state.D,
	}
}

func (state GuardState) RotateRight() GuardState {
	return GuardState{
		C: state.C,
		D: state.D.RotateRight(),
	}
}

type Coordinate struct {
	X int
	Y int
}

func (c Coordinate) Step(d Direction) Coordinate {
	return Coordinate{
		c.X + d.X,
		c.Y + d.Y,
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
