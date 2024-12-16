package maze

import (
	"errors"
	"fmt"
	"hmcalister/AdventOfCode/gridutils"
	"log/slog"
	"math"
	"slices"

	priorityqueue "github.com/hmcalister/Go-DSA/queue/PriorityQueue"
	hashset "github.com/hmcalister/Go-DSA/set/HashSet"
)

type Maze struct {
	startPosition gridutils.Coordinate
	endPosition   gridutils.Coordinate
	coordinateMap *hashset.HashSet[gridutils.Coordinate]
	mazeWidth     int
	mazeHeight    int
}

func NewMaze(mazeStrs []string) Maze {
	maze := Maze{
		coordinateMap: hashset.New[gridutils.Coordinate](),
		mazeWidth:     len(mazeStrs[0]),
		mazeHeight:    len(mazeStrs),
	}

	for y, line := range mazeStrs {
		slog.Debug("creating maze", "next row", line)
		for x, cell := range line {
			c := gridutils.Coordinate{X: x, Y: y}
			switch cell {
			case EMPTY_RUNE:
				slog.Debug("found empty cell", "coordinate", c)
				maze.coordinateMap.Add(c)
			case START_RUNE:
				slog.Debug("found start", "coordinate", c)
				maze.coordinateMap.Add(c)
				maze.startPosition = c
			case END_RUNE:
				slog.Debug("found end", "coordinate", c)
				maze.coordinateMap.Add(c)
				maze.endPosition = c
			}
		}
	}

	return maze
}

func (maze Maze) heuristic(c gridutils.Coordinate) int {
	// deltaX := c.X - maze.endPosition.X
	// if deltaX < 0 {
	// 	deltaX *= -1
	// }

	// deltaY := c.Y - maze.endPosition.Y
	// if deltaY < 0 {
	// 	deltaY *= -1
	// }

	// return deltaX + deltaY

	return 0
}

func (maze Maze) getPathfindStepNeighbors(step pathfindStepData) []pathfindStepData {
	neighbors := make([]pathfindStepData, 0)
	forwardCoord := step.position.Step(step.incomingDirection)
	if maze.coordinateMap.Contains(forwardCoord) {
		neighbors = append(neighbors, pathfindStepData{
			position:          forwardCoord,
			incomingDirection: step.incomingDirection,
			g:                 1 + step.g,
			h:                 maze.heuristic(forwardCoord),
		})
	}

	leftDirection := step.incomingDirection.RotateLeft()
	leftCoord := step.position.Step(leftDirection)
	if maze.coordinateMap.Contains(leftCoord) {
		neighbors = append(neighbors, pathfindStepData{
			position:          leftCoord,
			incomingDirection: step.incomingDirection,
			g:                 1001 + step.g,
			h:                 maze.heuristic(leftCoord),
		})
	}

	rightDirection := step.incomingDirection.RotateRight()
	rightCoord := step.position.Step(rightDirection)
	if maze.coordinateMap.Contains(rightCoord) {
		neighbors = append(neighbors, pathfindStepData{
			position:          rightCoord,
			incomingDirection: step.incomingDirection,
			g:                 1001 + step.g,
			h:                 maze.heuristic(rightCoord),
		})
	}

	return neighbors
}

// Find the optimal path using A* pathfinding
//
// Heuristic is manhattan distance to end
func (maze Maze) ComputeOptimalPath() (int, error) {

	// Track the (currently known) best path  costs to each coordinate
	gScores := make(map[coordDirectionTuple]int)

	cameFrom := make(map[pathfindStepData]pathfindStepData)
	priorityQueue := priorityqueue.New[pathfindStepData](pathfindStepComparator)
	// We have to manage the first step ourselves, unfortunately, to handle the possible paths
	for _, direction := range []gridutils.Direction{gridutils.DIRECTION_UP, gridutils.DIRECTION_RIGHT, gridutils.DIRECTION_DOWN, gridutils.DIRECTION_LEFT} {
		firstStepPosition := maze.startPosition.Step(direction)
		if maze.coordinateMap.Contains(firstStepPosition) {
			firstPathfindStep := pathfindStepData{
				position:          firstStepPosition,
				incomingDirection: direction,
				g:                 1,
				h:                 maze.heuristic(firstStepPosition),
			}
			slog.Debug("found potential first step", "first step", firstPathfindStep)
			priorityQueue.Add(firstPathfindStep)
			gScores[firstPathfindStep.getPositionDirection()] = 0
			cameFrom[firstPathfindStep] = pathfindStepData{
				position:          maze.startPosition,
				incomingDirection: direction,
				g:                 0,
				h:                 maze.heuristic(maze.startPosition),
			}
		}
	}

	for priorityQueue.Size() > 0 {
		currentStep, _ := priorityQueue.Remove()

		slog.Debug("considering pathfind step", "current step", currentStep)

		// sort queue for debugging
		queueData := priorityQueue.Items()
		slices.SortFunc(queueData, pathfindStepComparator)
		fmt.Printf("Next Item: %+v\n", currentStep)
		for index, item := range queueData {
			fmt.Printf("\t%v: %+v\n", index, item)
		}

		if maze.endPosition.Equal(currentStep.position) {
			completePathSteps := make(map[gridutils.Coordinate]pathfindStepData)
			reconstructedStep := currentStep
			for reconstructedStep.position != maze.startPosition {
				fmt.Printf("reconstructing path: %+v\n", reconstructedStep)
				completePathSteps[reconstructedStep.position] = reconstructedStep
				reconstructedStep = cameFrom[reconstructedStep]
			}

			mazeString := make([]rune, maze.mazeWidth*(maze.mazeWidth+1))
			mazeStringIndex := 0
			for y := 0; y < maze.mazeHeight; y += 1 {
				for x := 0; x < maze.mazeWidth; x += 1 {
					c := gridutils.Coordinate{X: x, Y: y}
					if c == maze.endPosition {
						mazeString[mazeStringIndex] = END_RUNE
					} else if c == maze.startPosition {
						mazeString[mazeStringIndex] = START_RUNE
					} else if reconstructedPathStep, ok := completePathSteps[c]; ok {
						switch reconstructedPathStep.incomingDirection {
						case gridutils.DIRECTION_UP:
							mazeString[mazeStringIndex] = '^'
						case gridutils.DIRECTION_RIGHT:
							mazeString[mazeStringIndex] = '>'
						case gridutils.DIRECTION_DOWN:
							mazeString[mazeStringIndex] = 'v'
						case gridutils.DIRECTION_LEFT:
							mazeString[mazeStringIndex] = '<'
						}
					} else if !maze.coordinateMap.Contains(c) {
						mazeString[mazeStringIndex] = WALL_RUNE
					} else {
						mazeString[mazeStringIndex] = ' '
					}
					mazeStringIndex += 1
				}
				mazeString[mazeStringIndex] = '\n'
				mazeStringIndex += 1
			}
			fmt.Println(string(mazeString))

			// We have found the end
			return currentStep.g, nil
		}

		// Get all the neighbors of the current coordinate and add them to the queue
		// accounting for the additional costs of turning and so on
		currentNeighbors := maze.getPathfindStepNeighbors(currentStep)
		fmt.Printf("considering neighbors %+v\n", currentNeighbors)
		for _, neighbor := range currentNeighbors {
			slog.Debug("pathfind step neighbor", "neighbor", neighbor)
			neighborBestGScore, ok := gScores[neighbor.getPositionDirection()]
			// If we have never seen this coordinate before, the score is infinite
			if !ok {
				neighborBestGScore = math.MaxInt
			}

			if neighbor.g <= neighborBestGScore {
				gScores[neighbor.getPositionDirection()] = neighbor.g
				fmt.Printf("\tnew path added %+v -> %+v\n", currentStep, neighbor)
				priorityQueue.Add(neighbor)
				cameFrom[neighbor] = currentStep
			}
		}

		fmt.Println()
	}

	return -1, errors.New("could not find path to end")
}
