package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
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
	result, err := Part02(fileScanner)
	if err != nil {
		slog.Error("error processing file input", "error", err)
	}

	slog.Info("computation completed", "result", result)
}

func Part01(fileScanner *bufio.Scanner) (int, error) {
	list1 := make([]int, 0)
	list2 := make([]int, 0)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		lineIntegers := strings.Split(line, "   ")
		i1, _ := strconv.Atoi(lineIntegers[0])
		i2, _ := strconv.Atoi(lineIntegers[1])
		list1 = append(list1, i1)
		list2 = append(list2, i2)
	}

	slices.Sort(list1)
	slices.Sort(list2)

	difference := 0
	for i := 0; i < len(list1); i++ {
		d := list1[i] - list2[i]
		if d < 0 {
			difference -= d
		} else {
			difference += d
		}
	}

	return difference, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	list1 := make([]int, 0)
	list2 := make(map[int]int)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		lineIntegers := strings.Split(line, "   ")
		i1, _ := strconv.Atoi(lineIntegers[0])
		i2, _ := strconv.Atoi(lineIntegers[1])
		list1 = append(list1, i1)
		list2[i2] += 1
		slog.Debug("list two count", "integer", i2, "count", list2[i2])
	}

	score := 0
	for _, i := range list1 {
		score += list2[i] * i
	}

	return score, nil
}
