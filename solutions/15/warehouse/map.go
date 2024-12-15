package warehouse

import (
	"hmcalister/AdventOfCode/gridutils"
	"log/slog"

	hashset "github.com/hmcalister/Go-DSA/set/HashSet"
)

const (
	WALL_RUNE  rune = '#'
	BOX_RUNE   rune = 'O'
	ROBOT_RUNE rune = '@'
	EMPTY_RUNE rune = '.'
)

type WarehouseMap struct {
	wallMap       *hashset.HashSet[gridutils.Coordinate]
	boxMap        *hashset.HashSet[gridutils.Coordinate]
	robotPosition gridutils.Coordinate
	mapWidth      int
	mapHeight     int
}

