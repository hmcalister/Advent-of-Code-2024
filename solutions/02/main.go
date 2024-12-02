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

