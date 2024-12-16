package maze

import "hmcalister/AdventOfCode/gridutils"

type pathfindStepData struct {
	position          gridutils.Coordinate
	incomingDirection gridutils.Direction
	g                 int
	h                 int
}

type coordDirectionTuple struct {
	position          gridutils.Coordinate
	incomingDirection gridutils.Direction
}

func pathfindStepComparator(a, b pathfindStepData) int {
	return a.computeFScore() - b.computeFScore()
}

func (pathfindStep pathfindStepData) getPositionDirection() coordDirectionTuple {
	return coordDirectionTuple{
		position:          pathfindStep.position,
		incomingDirection: pathfindStep.incomingDirection,
	}
}

func (pathfindStep pathfindStepData) computeFScore() int {
	return pathfindStep.g + pathfindStep.h
}
