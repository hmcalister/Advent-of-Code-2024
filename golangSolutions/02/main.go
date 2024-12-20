package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const (
	INPUT_FILE_PATH = "puzzleInput"
	// INPUT_FILE_PATH = "test"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Debug Flag")
	flag.Parse()
	logFileHandler := SetLogging(*debugFlag)
	defer logFileHandler.Close()

	inputFile, err := os.Open(INPUT_FILE_PATH)
	if err != nil {
		slog.Error("error opening input file", "error", err)
	}
	defer inputFile.Close()

	fileScanner := bufio.NewScanner(inputFile)
	result, err := Part02(fileScanner)
	if err != nil {
		slog.Error("error processing file input", "error", err)
	}

	slog.Info("computation completed", "result", result)
}

func isSafe(levels []string) bool {
	previousValue, _ := strconv.Atoi(levels[0])
	currentValue, _ := strconv.Atoi(levels[1])
	isIncreasing := (currentValue - previousValue) > 0
	for _, level := range levels[1:] {
		currentValue, _ = strconv.Atoi(level)
		difference := currentValue - previousValue
		if !isIncreasing {
			difference *= -1
		}
		if difference < 1 || difference > 3 {
			return false
		}

		previousValue = currentValue
	}

	return true
}

func Part01(fileScanner *bufio.Scanner) (int, error) {
	totalSafe := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		levels := strings.Split(line, " ")
		levelSafe := isSafe(levels)
		slog.Info("safety determined", "levels", levels, "safe", levelSafe)
		if levelSafe {
			totalSafe += 1
		}

	}
	return totalSafe, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	totalSafe := 0
reportLoop:
	for fileScanner.Scan() {
		line := fileScanner.Text()
		levels := strings.Split(line, " ")
		levelSafe := isSafe(levels)
		slog.Info("undamped safety determined", "levels", levels, "safe", levelSafe)
		if levelSafe {
			totalSafe += 1
			continue reportLoop
		}

		tempLevels := make([]string, len(levels)-1)
		for i := range len(levels) {
			copy(tempLevels, levels[:i])
			copy(tempLevels[i:], levels[i+1:])
			levelSafe = isSafe(tempLevels)
			slog.Debug("damped safety determined", "dampedLevelIndex", i, "dampedSlice", tempLevels, "safe", levelSafe)
			if levelSafe {
				totalSafe += 1
				continue reportLoop
			}
		}

	}
	return totalSafe, nil
}
