package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	total := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		slog.Debug("line read", "line", line)
		exp := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)

		foundMatches := exp.FindAllStringSubmatch(line, -1)
		slog.Debug("regex applied", "matches", foundMatches)

		for _, match := range foundMatches {
			slog.Debug("match loop", "match", match)
			v1, _ := strconv.Atoi(match[1])
			v2, _ := strconv.Atoi(match[2])
			total += v1 * v2
			slog.Debug("match calculation", "v1", v1, "v2", v2, "mul", v1*v2, "new total", total)
		}

	}

	return total, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	return 0, nil
}
