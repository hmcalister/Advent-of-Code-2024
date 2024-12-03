package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Debug Flag")
	inputFilePath := flag.String("inputFile", "puzzleInput", "Path to input file.")
	flag.Parse()
	logFileHandler := SetLogging(*debugFlag)
	defer logFileHandler.Close()

	inputFile, err := os.Open(*inputFilePath)
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
