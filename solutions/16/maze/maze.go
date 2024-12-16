package maze

import (
	"errors"
	"fmt"
	"hmcalister/AdventOfCode/gridutils"
	"log/slog"
	"math"

	arrayqueue "github.com/hmcalister/Go-DSA/queue/ArrayQueue"
	priorityqueue "github.com/hmcalister/Go-DSA/queue/PriorityQueue"
	hashset "github.com/hmcalister/Go-DSA/set/HashSet"
)

type Maze struct {
	startPosition gridutils.Coordinate
	endPosition   gridutils.Coordinate
	coordinateMap *hashset.HashSet[gridutils.Coordinate]
	mazeWidth     int
	mazeHeight    int
	gScore        map[pathfindStepData]int
	fScore        map[pathfindStepData]int
}

func NewMaze(mazeStrs []string) Maze {
	maze := Maze{
		coordinateMap: hashset.New[gridutils.Coordinate](),
		mazeWidth:     len(mazeStrs[0]),
		mazeHeight:    len(mazeStrs),
		gScore:        make(map[pathfindStepData]int),
		fScore:        make(map[pathfindStepData]int),
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

func (maze Maze) heuristic(step pathfindStepData) int {
	deltaX := step.position.X - maze.endPosition.X
	if deltaX < 0 {
		deltaX *= -1
	}

	deltaY := step.position.Y - maze.endPosition.Y
	if deltaY < 0 {
		deltaY *= -1
	}

	return deltaX + deltaY
}

func (maze Maze) getGScore(step pathfindStepData) int {
	if g, ok := maze.gScore[step]; ok {
		return g
	}
	return math.MaxInt
}

func (maze Maze) getFScore(step pathfindStepData) int {
	if g, ok := maze.fScore[step]; ok {
		return g
	}
	return math.MaxInt
}

func (maze Maze) expandStepSingleOptimalPath(
	step pathfindStepData,
	openset *priorityqueue.PriorityQueue[pathfindStepData],
	cameFrom map[pathfindStepData]pathfindStepData,
) {
	stepGScore := maze.getGScore(step)

	forwardCoord := step.position.Step(step.incomingDirection)
	if maze.coordinateMap.Contains(forwardCoord) {
		forwardStep := pathfindStepData{
			position:          forwardCoord,
			incomingDirection: step.incomingDirection,
		}
		forwardGScoreViaCurrent := stepGScore + 1
		forwardGScorePrior := maze.getGScore(forwardStep)
		if forwardGScoreViaCurrent < forwardGScorePrior {
			slog.Debug("expanded better path to forward neighbor", "current step", step, "forward step", forwardStep, "forward g", forwardGScoreViaCurrent)
			cameFrom[forwardStep] = step
			maze.gScore[forwardStep] = forwardGScoreViaCurrent
			maze.fScore[forwardStep] = forwardGScoreViaCurrent + maze.heuristic(forwardStep)
			if _, err := openset.Find(func(item pathfindStepData) bool {
				return item == forwardStep
			}); err != nil {
				openset.Add(forwardStep)
			}
		}
	}

	leftDirection := step.incomingDirection.RotateLeft()
	leftCoord := step.position.Step(leftDirection)
	if maze.coordinateMap.Contains(leftCoord) {
		leftStep := pathfindStepData{
			position:          leftCoord,
			incomingDirection: leftDirection,
		}
		leftGScoreViaCurrent := stepGScore + 1001
		leftGScorePrior := maze.getGScore(leftStep)
		if leftGScoreViaCurrent < leftGScorePrior {
			slog.Debug("expanded better path to left neighbor", "current step", step, "left step", leftStep, "left g", leftGScoreViaCurrent)
			cameFrom[leftStep] = step
			maze.gScore[leftStep] = leftGScoreViaCurrent
			maze.fScore[leftStep] = leftGScoreViaCurrent + maze.heuristic(leftStep)
			if _, err := openset.Find(func(item pathfindStepData) bool {
				return item == leftStep
			}); err != nil {
				openset.Add(leftStep)
			}
		}
	}

	rightDirection := step.incomingDirection.RotateRight()
	rightCoord := step.position.Step(rightDirection)
	if maze.coordinateMap.Contains(rightCoord) {
		rightStep := pathfindStepData{
			position:          rightCoord,
			incomingDirection: rightDirection,
		}
		rightGScoreViaCurrent := stepGScore + 1001
		rightGScorePrior := maze.getGScore(rightStep)
		if rightGScoreViaCurrent < rightGScorePrior {
			slog.Debug("expanded better path to right neighbor", "current step", step, "right step", rightStep, "right g", rightGScoreViaCurrent)
			cameFrom[rightStep] = step
			maze.gScore[rightStep] = rightGScoreViaCurrent
			maze.fScore[rightStep] = rightGScoreViaCurrent + maze.heuristic(rightStep)
			if _, err := openset.Find(func(item pathfindStepData) bool {
				return item == rightStep
			}); err != nil {
				openset.Add(rightStep)
			}
		}
	}
}

func (maze Maze) expandStepAllOptimalPaths(
	step pathfindStepData,
	openset *priorityqueue.PriorityQueue[pathfindStepData],
	cameFrom map[pathfindStepData][]pathfindStepData,
) {
	stepGScore := maze.getGScore(step)

	forwardCoord := step.position.Step(step.incomingDirection)
	if maze.coordinateMap.Contains(forwardCoord) {
		forwardStep := pathfindStepData{
			position:          forwardCoord,
			incomingDirection: step.incomingDirection,
		}
		forwardGScoreViaCurrent := stepGScore + 1
		forwardGScorePrior := maze.getGScore(forwardStep)
		if forwardGScoreViaCurrent == forwardGScorePrior {
			// append the paths if we are JUST AS GOOD
			cameFrom[forwardStep] = append(cameFrom[forwardStep], step)
		}
		if forwardGScoreViaCurrent < forwardGScorePrior {
			slog.Debug("expanded better path to forward neighbor", "current step", step, "forward step", forwardStep, "forward g", forwardGScoreViaCurrent)
			// overwrite the paths if we have done STRICTLY better
			cameFrom[forwardStep] = []pathfindStepData{step}
			maze.gScore[forwardStep] = forwardGScoreViaCurrent
			maze.fScore[forwardStep] = forwardGScoreViaCurrent + maze.heuristic(forwardStep)
			if _, err := openset.Find(func(item pathfindStepData) bool {
				return item == forwardStep
			}); err != nil {
				openset.Add(forwardStep)
			}
		}
	}

	leftDirection := step.incomingDirection.RotateLeft()
	leftCoord := step.position.Step(leftDirection)
	if maze.coordinateMap.Contains(leftCoord) {
		leftStep := pathfindStepData{
			position:          leftCoord,
			incomingDirection: leftDirection,
		}
		leftGScoreViaCurrent := stepGScore + 1001
		leftGScorePrior := maze.getGScore(leftStep)
		if leftGScoreViaCurrent == leftGScorePrior {
			// append the paths if we are JUST AS GOOD
			cameFrom[leftStep] = append(cameFrom[leftStep], step)
		}
		if leftGScoreViaCurrent < leftGScorePrior {
			slog.Debug("expanded better path to left neighbor", "current step", step, "left step", leftStep, "left g", leftGScoreViaCurrent)
			// overwrite the paths if we have done STRICTLY better
			cameFrom[leftStep] = []pathfindStepData{step}
			maze.gScore[leftStep] = leftGScoreViaCurrent
			maze.fScore[leftStep] = leftGScoreViaCurrent + maze.heuristic(leftStep)
			if _, err := openset.Find(func(item pathfindStepData) bool {
				return item == leftStep
			}); err != nil {
				openset.Add(leftStep)
			}
		}
	}

	rightDirection := step.incomingDirection.RotateRight()
	rightCoord := step.position.Step(rightDirection)
	if maze.coordinateMap.Contains(rightCoord) {
		rightStep := pathfindStepData{
			position:          rightCoord,
			incomingDirection: rightDirection,
		}
		rightGScoreViaCurrent := stepGScore + 1001
		rightGScorePrior := maze.getGScore(rightStep)
		if rightGScoreViaCurrent == rightGScorePrior {
			// append the paths if we are JUST AS GOOD
			cameFrom[rightStep] = append(cameFrom[rightStep], step)
		}
		if rightGScoreViaCurrent < rightGScorePrior {
			slog.Debug("expanded better path to right neighbor", "current step", step, "right step", rightStep, "right g", rightGScoreViaCurrent)
			// overwrite the paths if we have done STRICTLY better
			cameFrom[rightStep] = []pathfindStepData{step}
			maze.gScore[rightStep] = rightGScoreViaCurrent
			maze.fScore[rightStep] = rightGScoreViaCurrent + maze.heuristic(rightStep)
			if _, err := openset.Find(func(item pathfindStepData) bool {
				return item == rightStep
			}); err != nil {
				openset.Add(rightStep)
			}
		}
	}
}

	completePathSteps := make(map[gridutils.Coordinate]pathfindStepData)
	reconstructedStep := finalStep
	for reconstructedStep.position != maze.startPosition {
		fmt.Printf("reconstructing path: %+v\n", reconstructedStep)
		completePathSteps[reconstructedStep.position] = reconstructedStep
		reconstructedStep = cameFrom[reconstructedStep]
	}

	mazeString := make([]rune, maze.mazeHeight*(maze.mazeWidth+1))
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
}

// Find the optimal path using A* pathfinding
func (maze Maze) ComputeOptimalPath() (int, error) {
	pathfindStepComparator := func(a, b pathfindStepData) int {
		return maze.getFScore(a) - maze.getFScore(b)
	}
	openset := priorityqueue.New(pathfindStepComparator)
	cameFrom := make(map[pathfindStepData]pathfindStepData)

	initialPosition := pathfindStepData{
		position:          maze.startPosition,
		incomingDirection: gridutils.DIRECTION_RIGHT,
	}
	maze.gScore[initialPosition] = 0
	maze.fScore[initialPosition] = maze.heuristic(initialPosition)
	openset.Add(initialPosition)

	for openset.Size() > 0 {
		currentStep, _ := openset.Remove()
		slog.Debug("expanding node", "current step", currentStep)
		currentGScore := maze.getGScore(currentStep)

		if currentStep.position.Equal(maze.endPosition) {
			maze.reconstructPath(currentStep, cameFrom)
			return currentGScore, nil
		}

		maze.expandStep(currentStep, openset, cameFrom)
	}

	return -1, errors.New("could not find path to end")
}
