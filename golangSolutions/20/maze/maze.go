package maze

import (
	"errors"
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
	gScore        map[gridutils.Coordinate]int
	fScore        map[gridutils.Coordinate]int
}

func NewMaze(mazeStrings []string) Maze {
	maze := Maze{
		coordinateMap: hashset.New[gridutils.Coordinate](),
		mazeWidth:     len(mazeStrings[0]),
		mazeHeight:    len(mazeStrings),
		gScore:        make(map[gridutils.Coordinate]int),
		fScore:        make(map[gridutils.Coordinate]int),
	}

	for y, row := range mazeStrings {
		for x, cell := range row {
			c := gridutils.Coordinate{X: x, Y: y}
			switch cell {
			case START_RUNE:
				maze.startPosition = c
				maze.coordinateMap.Add(c)
				slog.Debug("found start position", "coordinate", c)
			case END_RUNE:
				maze.endPosition = c
				maze.coordinateMap.Add(c)
				slog.Debug("found start position", "coordinate", c)
			case EMPTY_RUNE:
				maze.coordinateMap.Add(c)
			case WALL_RUNE:

			default:
				slog.Debug("found unexpected rune", "rune", cell, "coordinate", c)
			}
		}
	}

	return maze
}

// --------------------------------------------------------------------------------
// Print methods

func (maze Maze) String() string {
	mazeString := make([]rune, maze.mazeHeight*(maze.mazeWidth+1))
	mazeStringIndex := 0
	for y := 0; y < maze.mazeHeight; y += 1 {
		for x := 0; x < maze.mazeWidth; x += 1 {
			c := gridutils.Coordinate{X: x, Y: y}
			if c == maze.endPosition {
				mazeString[mazeStringIndex] = END_RUNE
			} else if c == maze.startPosition {
				mazeString[mazeStringIndex] = START_RUNE
			} else if !maze.coordinateMap.Contains(c) {
				mazeString[mazeStringIndex] = WALL_RUNE
			} else {
				mazeString[mazeStringIndex] = '.'
			}
			mazeStringIndex += 1
		}
		mazeString[mazeStringIndex] = '\n'
		mazeStringIndex += 1
	}
	return string(mazeString)
}

func (maze Maze) StringWithPath(path []gridutils.Coordinate) string {
	mazeString := make([]rune, maze.mazeHeight*(maze.mazeWidth+1))
	mazeStringIndex := 0
	for y := 0; y < maze.mazeHeight; y += 1 {
		for x := 0; x < maze.mazeWidth; x += 1 {
			c := gridutils.Coordinate{X: x, Y: y}
			if c == maze.endPosition {
				mazeString[mazeStringIndex] = END_RUNE
			} else if c == maze.startPosition {
				mazeString[mazeStringIndex] = START_RUNE
			} else if !maze.coordinateMap.Contains(c) {
				mazeString[mazeStringIndex] = WALL_RUNE
			} else {
				mazeString[mazeStringIndex] = '.'
			}
			mazeStringIndex += 1
		}
		mazeString[mazeStringIndex] = '\n'
		mazeStringIndex += 1
	}
	for _, pathStep := range path {
		linearCoordinate := (maze.mazeWidth+1)*pathStep.Y + pathStep.X
		mazeString[linearCoordinate] = 'O'
	}

	return string(mazeString)
}

// --------------------------------------------------------------------------------
// Pathfinding methods

func (maze Maze) heuristic(step gridutils.Coordinate) int {
	deltaX := step.X - maze.endPosition.X
	if deltaX < 0 {
		deltaX *= -1
	}

	deltaY := step.Y - maze.endPosition.Y
	if deltaY < 0 {
		deltaY *= -1
	}

	return deltaX + deltaY
}

func (maze Maze) getGScore(step gridutils.Coordinate) int {
	if g, ok := maze.gScore[step]; ok {
		return g
	}
	return math.MaxInt
}

func (maze Maze) getFScore(step gridutils.Coordinate) int {
	if g, ok := maze.fScore[step]; ok {
		return g
	}
	return math.MaxInt
}

func (maze Maze) expandStepSingleOptimalPath(
	step gridutils.Coordinate,
	openset *priorityqueue.PriorityQueue[gridutils.Coordinate],
	cameFrom map[gridutils.Coordinate]gridutils.Coordinate,
) {
	stepGScore := maze.getGScore(step)

	for _, direction := range []gridutils.Direction{gridutils.DIRECTION_UP, gridutils.DIRECTION_RIGHT, gridutils.DIRECTION_DOWN, gridutils.DIRECTION_LEFT} {
		nextStep := step.Step(direction)
		if maze.coordinateMap.Contains(nextStep) {
			forwardGScoreViaCurrent := stepGScore + 1
			forwardGScorePrior := maze.getGScore(nextStep)
			if forwardGScoreViaCurrent < forwardGScorePrior {
				// slog.Debug("expanded better path to next neighbor", "current step", step, "next step", nextStep, "next g", forwardGScoreViaCurrent)
				cameFrom[nextStep] = step
				maze.gScore[nextStep] = forwardGScoreViaCurrent
				maze.fScore[nextStep] = forwardGScoreViaCurrent + maze.heuristic(nextStep)
				if _, err := openset.Find(func(item gridutils.Coordinate) bool {
					return item.Equal(nextStep)
				}); err != nil {
					openset.Add(nextStep)
				}
			}
		}
	}
}

// Find the optimal path using A* pathfinding
func (maze Maze) ComputeOptimalPath() ([]gridutils.Coordinate, error) {
	pathfindStepComparator := func(a, b gridutils.Coordinate) int {
		return maze.getFScore(a) - maze.getFScore(b)
	}
	openset := priorityqueue.New(pathfindStepComparator)
	cameFrom := make(map[gridutils.Coordinate]gridutils.Coordinate)

	initialPosition := maze.startPosition
	maze.gScore[initialPosition] = 0
	maze.fScore[initialPosition] = maze.heuristic(initialPosition)
	openset.Add(initialPosition)

	for openset.Size() > 0 {
		currentStep, _ := openset.Remove()
		// slog.Debug("expanding node", "current step", currentStep)

		if currentStep.Equal(maze.endPosition) {
			reconstructedPath := make([]gridutils.Coordinate, 0)
			reconstructedStep := currentStep
			for reconstructedStep != maze.startPosition {
				reconstructedPath = append(reconstructedPath, reconstructedStep)
				reconstructedStep = cameFrom[reconstructedStep]
			}
			reconstructedPath = append(reconstructedPath, maze.startPosition)
			slices.Reverse(reconstructedPath)
			return reconstructedPath, nil
		}

		maze.expandStepSingleOptimalPath(currentStep, openset, cameFrom)
	}

	return nil, errors.New("could not find path to end")
}

// --------------------------------------------------------------------------------
// Problem specific pathfinding

func (maze Maze) StringTwoStepCheat(cheatOrigin gridutils.Coordinate, cheatDirection gridutils.Direction) string {
	mazeString := make([]rune, maze.mazeHeight*(maze.mazeWidth+1))
	mazeStringIndex := 0
	for y := 0; y < maze.mazeHeight; y += 1 {
		for x := 0; x < maze.mazeWidth; x += 1 {
			c := gridutils.Coordinate{X: x, Y: y}
			if c == maze.endPosition {
				mazeString[mazeStringIndex] = END_RUNE
			} else if c == maze.startPosition {
				mazeString[mazeStringIndex] = START_RUNE
			} else if !maze.coordinateMap.Contains(c) {
				mazeString[mazeStringIndex] = WALL_RUNE
			} else {
				mazeString[mazeStringIndex] = '.'
			}
			mazeStringIndex += 1
		}
		mazeString[mazeStringIndex] = '\n'
		mazeStringIndex += 1
	}

	cheatedStep := cheatOrigin
	for cheatStepIndex := range 2 {
		cheatedStep = cheatedStep.Step(cheatDirection)
		linearCoordinate := (maze.mazeWidth+1)*cheatedStep.Y + cheatedStep.X
		mazeString[linearCoordinate] = rune(cheatStepIndex + 1 + '0')
	}

	return string(mazeString)
}
