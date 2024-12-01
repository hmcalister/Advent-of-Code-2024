package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
)

const (
	INPUT_FILE_PATH = "puzzleInput"
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
	result, err := Part01(fileScanner)
	if err != nil {
		slog.Error("error processing file input", "error", err)
	}

	slog.Info("computation completed", "result", result)
}

func Part01(fileScanner *bufio.Scanner) (int, error) {
	return 0, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	return 0, nil
}
