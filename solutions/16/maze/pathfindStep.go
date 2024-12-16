package maze

import "hmcalister/AdventOfCode/gridutils"

type pathfindStepData struct {
	position          gridutils.Coordinate
	incomingDirection gridutils.Direction
}
