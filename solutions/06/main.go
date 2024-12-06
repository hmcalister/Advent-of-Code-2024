package main

import (
	"bufio"
	"errors"
	"flag"
	"log/slog"
	"os"
	"time"
)

const (
	EMPTY_RUNE       rune = '.'
	OBSTACLE_RUNE    rune = '#'
	GUARD_UP_RUNE    rune = '^'
	GUARD_RIGHT_RUNE rune = '>'
	GUARD_DOWN_RUNE  rune = 'v'
	GUARD_LEFT_RUNE  rune = '<'
)

func main() {
	debugFlag := flag.Bool("debug", false, "Debug Flag")
	inputFilePath := flag.String("inputFile", "puzzleInput", "Path to input file.")
	selectedPart := flag.Int("part", 0, "Part to execute. Must be 1 or 2.")
	flag.Parse()
	logFileHandler := SetLogging(*debugFlag)
	defer logFileHandler.Close()

	inputFile, err := os.Open(*inputFilePath)
	if err != nil {
		slog.Error("error opening input file", "error", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	fileScanner := bufio.NewScanner(inputFile)
	if err != nil {
		slog.Error("error processing file input", "error", err)
		os.Exit(1)
	}

	var result int
	computationStartTime := time.Now()
	switch *selectedPart {
	case 1:
		result, err = Part01(fileScanner)
	case 2:
		result, err = Part02(fileScanner)
	default:
		slog.Error("invalid part selected, part must be one of 1 or 2", "part selected", *selectedPart)
		os.Exit(1)
	}
	computationEndTime := time.Now()
	if err != nil {
		slog.Error("error encountered during computation", "error", err, "part selected", *selectedPart)
		os.Exit(1)
	}

	slog.Info("computation completed", "result", result, "computation time elapsed (ns)", computationEndTime.Sub(computationStartTime).Nanoseconds())
}

type MapData struct {
	Width       int
	Height      int
	ObstacleMap map[Coordinate]interface{}
}

// returns the width, height of the grid, a map of coordinates to obstacles, the guard position and the guard direction
func parseInput(inputLines []string) (MapData, GuardState) {
	guardState := GuardState{
		Coordinate: Coordinate{-1, -1},
		Direction:  DIRECTION_UP,
	}

	obstacleMap := make(map[Coordinate]interface{})
	for y, line := range inputLines {
		slog.Debug("read line", "line", line)
		for x, repRune := range line {
			c := Coordinate{x, y}
			switch repRune {
			case OBSTACLE_RUNE:
				slog.Debug("found obstacle", "coordinate", c)
				obstacleMap[c] = struct{}{}
			case GUARD_UP_RUNE:
				slog.Debug("found guard up", "coordinate", c)
				guardState.Coordinate = c
				guardState.Direction = DIRECTION_UP
			case GUARD_RIGHT_RUNE:
				slog.Debug("found guard right", "coordinate", c)
				guardState.Coordinate = c
				guardState.Direction = DIRECTION_RIGHT
			case GUARD_DOWN_RUNE:
				slog.Debug("found guard down", "coordinate", c)
				guardState.Coordinate = c
				guardState.Direction = DIRECTION_DOWN
			case GUARD_LEFT_RUNE:
				slog.Debug("found guard left", "coordinate", c)
				guardState.Coordinate = c
				guardState.Direction = DIRECTION_LEFT
			case EMPTY_RUNE:
				// slog.Debug("empty coordinate")
			default:
				slog.Error("unexpected rune found", "rune", repRune, "location x", x, "location y", y)
			}
		}
	}

	if guardState.Coordinate.X == -1 && guardState.Coordinate.Y == -1 {
		slog.Error("no guard found")
		os.Exit(1)
	}

	return MapData{len(inputLines[0]), len(inputLines), obstacleMap}, guardState
}

// Count the number of visited cells (not states, direction is irrelevant) in path set by the initial guard state
// If the path loops at any point, an error is returned
func (m MapData) CheckVisitedCells(guardState GuardState) (int, error) {
	visitedCells := make(map[Coordinate]interface{})
	visitedStates := make(map[GuardState]interface{})

	for guardState.InBounds(m.Width, m.Height) {
		if _, ok := visitedStates[guardState]; ok {
			// We have seen this state before, therefore we are in a loop
			return -1, errors.New("a loop has occurred in the path")
		}
		visitedStates[guardState] = struct{}{}
		visitedCells[guardState.Coordinate] = struct{}{}

		nextState := guardState.Step()
		if _, ok := m.ObstacleMap[nextState.Coordinate]; ok {
			// slog.Debug("found obstacle", "current state", guardState, "numVisitedCells", len(visitedCells))
			guardState = guardState.EncounterObstacle()
			continue
		}

		// slog.Debug("making step", "current state", guardState, "numVisitedCells", len(visitedCells))
		guardState = nextState
	}

	return len(visitedCells), nil
}

func Part01(fileScanner *bufio.Scanner) (int, error) {
	inputLines := make([]string, 0)
	for fileScanner.Scan() {
		inputLines = append(inputLines, fileScanner.Text())
	}
	mapData, guardState := parseInput(inputLines)
	slog.Debug("parsed input", "mapWidth", mapData.Width, "mapHeight", mapData.Height, "obstacleMap", mapData.ObstacleMap, "guardState", guardState)

	numVisitedCells, err := mapData.CheckVisitedCells(guardState)
	if err != nil {
		slog.Error("error occurred when checking path", "error", err)
		os.Exit(1)
	}

	return numVisitedCells, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	inputLines := make([]string, 0)
	for fileScanner.Scan() {
		inputLines = append(inputLines, fileScanner.Text())
	}
	mapWidth, mapHeight, obstacleMap, guardCoordinate, guardDirection := parseInput(inputLines)
	slog.Debug("parsed input", "mapWidth", mapWidth, "mapHeight", mapHeight, "obstacleMap", obstacleMap, "guardCoordinate", guardCoordinate, "guardDirection", guardDirection)

	visitedCells := make(map[Coordinate]interface{})
	pathCorners := make(map[Coordinate]interface{})
	for guardCoordinate.X >= 0 &&
		guardCoordinate.X < mapWidth &&
		guardCoordinate.Y >= 0 &&
		guardCoordinate.Y < mapHeight {
		visitedCells[guardCoordinate] = struct{}{}
		nextCoordinate := guardCoordinate.Step(guardDirection)
		if _, ok := obstacleMap[nextCoordinate]; ok {
			slog.Debug("found obstacle", "current coordinate", guardCoordinate, "direction", guardDirection)
			pathCorners[guardCoordinate] = struct{}{}
			guardDirection = guardDirection.RotateRight()
			continue
		}

		slog.Debug("making step", "current coordinate", guardCoordinate, "direction", guardDirection)
		guardCoordinate = guardCoordinate.Step(guardDirection)
	}
	slog.Debug("path complete", "path corners", pathCorners)

	return 0, nil
}
