package main

import (
	"bufio"
	"flag"
	"fmt"
	"hmcalister/AdventOfCode/gridutils"
	"hmcalister/AdventOfCode/maze"
	"log/slog"
	"os"
	"runtime/pprof"
	"slices"
	"time"
)

const (
	CPU_PROFILE_FILEPATH string = "profile"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Debug Flag")
	inputFilePath := flag.String("inputFile", "puzzleInput", "Path to input file.")
	selectedPart := flag.Int("part", 0, "Part to execute. Must be 1 or 2.")
	profile := flag.Bool("profile", false, "Flag to profile program")
	flag.Parse()
	if *profile {
		f, err := os.Create(CPU_PROFILE_FILEPATH)
		if err != nil {
			slog.Error("could not create cpu profile file", "file", CPU_PROFILE_FILEPATH)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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

func getAllCheatsUpToLength(initialPosition gridutils.Coordinate, maxCheatLength int) []gridutils.Coordinate {
	allCheats := make([]gridutils.Coordinate, 0)
	var remainingCheatLength int
	for x := -maxCheatLength; x <= maxCheatLength; x += 1 {
		if x >= 0 {
			remainingCheatLength = maxCheatLength - x
		} else {
			remainingCheatLength = maxCheatLength + x
		}
		for y := -remainingCheatLength; y <= remainingCheatLength; y += 1 {
			allCheats = append(allCheats, gridutils.Coordinate{
				X: initialPosition.X + x,
				Y: initialPosition.Y + y,
			})
		}
	}

	return allCheats
}

func getCheatLength(initialPosition, cheatPosition gridutils.Coordinate) int {
	delX := initialPosition.X - cheatPosition.X
	if delX < 0 {
		delX *= -1
	}
	delY := initialPosition.Y - cheatPosition.Y
	if delY < 0 {
		delY *= -1
	}

	return delX + delY
}

func Part01(fileScanner *bufio.Scanner) (int, error) {
	mazeStrs := make([]string, 0)
	for fileScanner.Scan() {
		mazeStrs = append(mazeStrs, fileScanner.Text())
	}
	mazeData := maze.NewMaze(mazeStrs)
	fmt.Println(mazeData)

	honestPath, err := mazeData.ComputeOptimalPath()
	if err != nil {
		slog.Error("error when computing honest optimal path", "error", err)
		return -1, err
	}

	cheatedPathSavingCounts := make(map[int]int)
	for honestPathIndex, honestPathStep := range honestPath {
		for _, d := range gridutils.AllDirections {
			// Step (twice) in specific direction
			cheatedStep := honestPathStep.Step(d).Step(d)
			cheatedStepIndex := slices.Index(honestPath, cheatedStep)

			// If the cheated step is somewhere further along the path we can save time
			// This also handles the case of the cheated step not being found (-1)
			if cheatedStepIndex > honestPathIndex+2 {
				cheatSaving := cheatedStepIndex - honestPathIndex - 2
				// if cheatSaving > 60 {
				// 	fmt.Printf("%vCheat Saving: %v\n", mazeData.StringWithPathAndCheat(honestPath, honestPathStep, d), cheatSaving)
				// }
				cheatedPathSavingCounts[cheatSaving] += 1
			}
		}
	}

	numCheatsAbove100 := 0
	for cheatLength, cheatCount := range cheatedPathSavingCounts {
		slog.Info("cheated path saving count", "cheated path saving", cheatLength, "number of cheats", cheatCount)
		if cheatLength >= 100 {
			numCheatsAbove100 += cheatCount
		}
	}

	return numCheatsAbove100, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	return 0, nil
}
