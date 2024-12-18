package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"hmcalister/AdventOfCode/gridutils"
	"hmcalister/AdventOfCode/maze"
	"log/slog"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
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

func parseInput(fileScanner *bufio.Scanner) (int, int, []gridutils.Coordinate) {
	if !fileScanner.Scan() {
		slog.Error("input file is empty")
		os.Exit(1)
	}
	mazeDimensionStr := strings.Split(fileScanner.Text(), ",")
	if len(mazeDimensionStr) != 2 {
		slog.Error("maze dimension string does not match expected format", "maze dimension string", mazeDimensionStr)
		os.Exit(1)
	}
	mazeWidth, mazeWidthErr := strconv.Atoi(mazeDimensionStr[0])
	mazeHeight, mazeHeightErr := strconv.Atoi(mazeDimensionStr[1])
	if errors.Join(mazeWidthErr, mazeHeightErr) != nil {
		slog.Error("could not parse maze dimension string to integer", "maze dimension string", mazeDimensionStr)
		os.Exit(1)
	}

	fallingByteCoords := make([]gridutils.Coordinate, 0)
	for fileScanner.Scan() {
		fallingByteCoordStr := strings.Split(fileScanner.Text(), ",")
		if len(fallingByteCoordStr) != 2 {
			slog.Error("falling byte string does not match expected format", "falling byte string", mazeDimensionStr)
			os.Exit(1)
		}
		byteX, byteXErr := strconv.Atoi(fallingByteCoordStr[0])
		byteY, byteYErr := strconv.Atoi(fallingByteCoordStr[1])
		if errors.Join(byteXErr, byteYErr) != nil {
			slog.Error("could not parse falling byte string to integer", "falling byte string", fallingByteCoordStr)
			os.Exit(1)
		}
		fallingByteCoords = append(fallingByteCoords, gridutils.Coordinate{X: byteX, Y: byteY})
	}

	return mazeWidth, mazeHeight, fallingByteCoords
}

func Part01(fileScanner *bufio.Scanner) (int, error) {
	mazeWidth, mazeHeight, fallingByteCoords := parseInput(fileScanner)
	slog.Debug("parsed input", "maze width", mazeWidth, "maze height", mazeHeight, "num falling bytes", len(fallingByteCoords))
	maze := maze.NewMaze(mazeWidth, mazeHeight, fallingByteCoords[:1024])
	fmt.Println(maze)
	optimalPath, err := maze.ComputeOptimalPath()
	if err != nil {
		return -1, err
	}
	fmt.Println(maze.StringWithPath(optimalPath))

	return len(optimalPath) - 1, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	mazeWidth, mazeHeight, fallingByteCoords := parseInput(fileScanner)
	slog.Debug("parsed input", "maze width", mazeWidth, "maze height", mazeHeight, "num falling bytes", len(fallingByteCoords))

	for byteIndex := 0; byteIndex < len(fallingByteCoords); byteIndex += 1 {
		slog.Info("attempting to block maze", "byte index", byteIndex, "total bytes", len(fallingByteCoords))
		maze := maze.NewMaze(mazeWidth, mazeHeight, fallingByteCoords[:byteIndex+1])
		_, err := maze.ComputeOptimalPath()
		if err != nil {
			fmt.Println(maze)
			fmt.Println(byteIndex, fallingByteCoords[byteIndex])
			return byteIndex, nil
		}
	}

	return -1, errors.New("path always available")
}
