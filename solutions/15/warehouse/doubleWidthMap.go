package warehouse

import (
	"hmcalister/AdventOfCode/gridutils"
	"log/slog"

	hashset "github.com/hmcalister/Go-DSA/set/HashSet"
	arraystack "github.com/hmcalister/Go-DSA/stack/ArrayStack"
)

type DoubleWidthWarehouseMap struct {
	wallMap *hashset.HashSet[gridutils.Coordinate]
	// Boxes are only stored by their left half...
	boxMap        *hashset.HashSet[gridutils.Coordinate]
	robotPosition gridutils.Coordinate
	mapWidth      int
	mapHeight     int
}

func NewDoubleWidthWarehouseMap(warehouseMapStrs []string) *DoubleWidthWarehouseMap {
	warehouse := &DoubleWidthWarehouseMap{
		wallMap:   hashset.New[gridutils.Coordinate](),
		boxMap:    hashset.New[gridutils.Coordinate](),
		mapWidth:  2 * len(warehouseMapStrs[0]),
		mapHeight: len(warehouseMapStrs),
	}

	for y, row := range warehouseMapStrs {
		slog.Debug("parsing row", "row", row)
		for x, cell := range row {
			currentCoordinate := gridutils.Coordinate{X: 2 * x, Y: y}
			nextCoordinate := gridutils.Coordinate{X: 2*x + 1, Y: y}
			switch cell {
			case WALL_RUNE:
				slog.Debug("found wall", "coordinate", currentCoordinate)
				slog.Debug("found wall", "coordinate", nextCoordinate)
				warehouse.wallMap.Add(currentCoordinate)
				warehouse.wallMap.Add(nextCoordinate)
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

func (warehouse *DoubleWidthWarehouseMap) String() string {
	s := make([]rune, warehouse.mapHeight*(warehouse.mapWidth+1)+10)

	currentStrIndex := 0
	for y := 0; y < warehouse.mapHeight; y += 1 {
		for x := 0; x < warehouse.mapWidth; x += 1 {
			coord := gridutils.Coordinate{X: x, Y: y}

			if warehouse.robotPosition.Equal(coord) {
				copy(s[currentStrIndex:], []rune("\033[32m"))
				currentStrIndex += 5
				s[currentStrIndex] = ROBOT_RUNE
				currentStrIndex += 1
				copy(s[currentStrIndex:], []rune("\033[0m"))
				currentStrIndex += 4
			} else if warehouse.wallMap.Contains(coord) {
				s[currentStrIndex] = WALL_RUNE
			} else if warehouse.boxMap.Contains(coord) {
				s[currentStrIndex] = '['
				s[currentStrIndex+1] = ']'
				currentStrIndex += 1
				x += 1
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

func (warehouse *DoubleWidthWarehouseMap) RobotStep(stepDirection gridutils.Direction) {
	proposedRobotPosition := warehouse.robotPosition.Step(stepDirection)

	// If robot is trying to walk into a wall: don't
	if warehouse.wallMap.Contains(proposedRobotPosition) {
		return
	}

	// Boxes are only stored with respect to their left hand side
	// So if moving up/down, we must check for boxes both directly infront of
	// *and* to the left side
	//
	// 		....[]...	< Detected with immediate check above
	// 		....@....
	// 		...[]....	< Detected only with check below *and* to left

	// If robot is not moving into a box, just move the robot
	if !warehouse.boxMap.Contains(proposedRobotPosition) && !warehouse.boxMap.Contains(proposedRobotPosition.Step(gridutils.DIRECTION_LEFT)) {
		warehouse.robotPosition = proposedRobotPosition
		return
	}

	// Robot is moving into a box
	//
	// We must find *all* boxes that are being affected.
	// This may include a multitude of boxes if moving up/down
	//
	// e.g. this robot moving up
	// ##############
	// ##......##..##
	// ##..........##
	// ##...[][]...##
	// ##....[]....##
	// ##.....@....##
	// ##############
	//
	// however, as seen later:
	// ##############
	// ##......##..##
	// ##...[][]...##
	// ##....[]....##
	// ##.....@....##
	// ##..........##
	// ##############
	//
	// One box that is "stuck" (upper right) prevents all boxes from moving even if they are free (upper left)

	// Iterate through the boxes we will affect with this push
	affectedBoxesStack := arraystack.New[gridutils.Coordinate]()
	// Track all the boxes we have seen (put into touchedBoxesStack at any point)
	affectedBoxesSet := hashset.New[gridutils.Coordinate]()

	if warehouse.boxMap.Contains(proposedRobotPosition) {
		affectedBoxesStack.Add(proposedRobotPosition)
		affectedBoxesSet.Add(proposedRobotPosition)
	}
	leftOfProposed := proposedRobotPosition.Step(gridutils.DIRECTION_LEFT)
	if warehouse.boxMap.Contains(leftOfProposed) {
		affectedBoxesStack.Add(leftOfProposed)
		affectedBoxesSet.Add(leftOfProposed)
	}

	for affectedBoxesStack.Size() > 0 {
		// Get the next box to check
		nextBoxPosition, _ := affectedBoxesStack.Remove()

		// These positions are affected by the push of the current box, and hence must be checked for other boxes
		// These depend on the direction of the push!
		// If moving left, we must check two left of current box only
		// Similar for right
		// If moving up, we must check both up+left, up, and up+right for potential boxes
		// Similar for down
		var potentialNextBoxPositions []gridutils.Coordinate
		var potentialWallPositions []gridutils.Coordinate
		switch stepDirection {
		case gridutils.DIRECTION_UP:
			up := nextBoxPosition.Step(gridutils.DIRECTION_UP)
			potentialNextBoxPositions = []gridutils.Coordinate{up.Step(gridutils.DIRECTION_LEFT), up, up.Step(gridutils.DIRECTION_RIGHT)}
			potentialWallPositions = []gridutils.Coordinate{up, up.Step(gridutils.DIRECTION_RIGHT)}
		case gridutils.DIRECTION_RIGHT:
			doubleRight := nextBoxPosition.Step(gridutils.DIRECTION_RIGHT).Step(gridutils.DIRECTION_RIGHT)
			potentialNextBoxPositions = []gridutils.Coordinate{doubleRight}
			potentialWallPositions = []gridutils.Coordinate{doubleRight}
		case gridutils.DIRECTION_DOWN:
			down := nextBoxPosition.Step(gridutils.DIRECTION_DOWN)
			potentialNextBoxPositions = []gridutils.Coordinate{down.Step(gridutils.DIRECTION_LEFT), down, down.Step(gridutils.DIRECTION_RIGHT)}
			potentialWallPositions = []gridutils.Coordinate{down, down.Step(gridutils.DIRECTION_RIGHT)}
		case gridutils.DIRECTION_LEFT:
			left := nextBoxPosition.Step(gridutils.DIRECTION_LEFT)
			potentialNextBoxPositions = []gridutils.Coordinate{left.Step(gridutils.DIRECTION_LEFT)}
			potentialWallPositions = []gridutils.Coordinate{left}
		}

		// Check each potential wall position. If we have encountered a wall, we cannot move the boxes or the robot so just return
		for _, potentialWallPosition := range potentialWallPositions {
			if warehouse.wallMap.Contains(potentialWallPosition) {
				return
			}
		}
		// This box has found no walls, so we are free to move it

		// Find all boxes affected by pushing this box
		for _, potentialBoxPosition := range potentialNextBoxPositions {
			// If the potential position is indeed a box (and we have not seen it before) then add it to the stack to process
			if warehouse.boxMap.Contains(potentialBoxPosition) && !affectedBoxesSet.Contains(potentialBoxPosition) {
				affectedBoxesStack.Add(potentialBoxPosition)
				affectedBoxesSet.Add(potentialBoxPosition)
			}
		}
	}

	// Debugging --- print the state before and after updating but only when a box moves!
	// fmt.Println(warehouse.String())
	// fmt.Println(stepDirection)
	// r := bufio.NewScanner(os.Stdin)
	// defer func() {
	// 	fmt.Println(warehouse.String())
	// 	r.Scan()
	// }()

	// We have determined that all boxes are free to move, and stored those boxes in affectedBoxesStack
	// First delete all the boxes then add all boxes back in an updated position
	for updatedBoxPosition := range affectedBoxesSet.Iterator() {
		warehouse.boxMap.Remove(updatedBoxPosition)
	}
	for updatedBoxPosition := range affectedBoxesSet.Iterator() {
		warehouse.boxMap.Add(updatedBoxPosition.Step(stepDirection))
	}

	// Finally, update the robot position
	warehouse.robotPosition = proposedRobotPosition

}

func (warehouse *DoubleWidthWarehouseMap) ComputeGPS() int {
	totalGps := 0
	for boxPosition := range warehouse.boxMap.Iterator() {
		currentGps := 100*boxPosition.Y + boxPosition.X
		totalGps += currentGps
		slog.Debug("computing box gps", "box position", boxPosition, "gps", currentGps, "updated total gps", totalGps)
	}
	return totalGps
}
