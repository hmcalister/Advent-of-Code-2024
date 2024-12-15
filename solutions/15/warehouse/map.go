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

func NewWarehouseMap(warehouseMapStrs []string) *WarehouseMap {
	warehouse := &WarehouseMap{
		wallMap:   hashset.New[gridutils.Coordinate](),
		boxMap:    hashset.New[gridutils.Coordinate](),
		mapWidth:  len(warehouseMapStrs[0]),
		mapHeight: len(warehouseMapStrs),
	}

	for y, row := range warehouseMapStrs {
		slog.Debug("parsing row", "row", row)
		for x, cell := range row {
			currentCoordinate := gridutils.Coordinate{X: x, Y: y}
			switch cell {
			case WALL_RUNE:
				slog.Debug("found wall", "coordinate", currentCoordinate)
				warehouse.wallMap.Add(currentCoordinate)
			case BOX_RUNE:
				slog.Debug("found box", "coordinate", currentCoordinate)
				warehouse.boxMap.Add(currentCoordinate)
			case ROBOT_RUNE:
				slog.Debug("found robot", "coordinate", currentCoordinate)
				warehouse.robotPosition = currentCoordinate
			}
		}
	}

	return warehouse
}

func (warehouse *WarehouseMap) String() string {
	s := make([]rune, warehouse.mapHeight*(warehouse.mapWidth+1))

	currentStrIndex := 0
	for y := 0; y < warehouse.mapHeight; y += 1 {
		for x := 0; x < warehouse.mapWidth; x += 1 {
			coord := gridutils.Coordinate{X: x, Y: y}

			if warehouse.robotPosition.Equal(coord) {
				s[currentStrIndex] = ROBOT_RUNE
			} else if warehouse.wallMap.Contains(coord) {
				s[currentStrIndex] = WALL_RUNE
			} else if warehouse.boxMap.Contains(coord) {
				s[currentStrIndex] = BOX_RUNE
			} else {
				s[currentStrIndex] = EMPTY_RUNE
			}
			currentStrIndex += 1
		}
		s[currentStrIndex] = '\n'
		currentStrIndex += 1
	}

	return string(s)
}

