package main

import (
	"bufio"
	"errors"
	"flag"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	linkedlist "github.com/hmcalister/Go-DSA/list/LinkedList"
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

// Defines a dependency graph --- key page depends on all pages in the value list and must occur *after* these\
func parsePageDependencyGraph(pageDependencies []string) map[int][]int {
	dependencyGraph := make(map[int][]int)
	for _, dependency := range pageDependencies {
		line := strings.TrimSpace(dependency)

		dependencyPages := strings.Split(line, "|")
		if len(dependencyPages) != 2 {
			slog.Error("found dependency line that is not of form 'a|b'", "offending line", dependency)
			continue
		}
		beforePage, err1 := strconv.Atoi(dependencyPages[0])
		afterPage, err2 := strconv.Atoi(dependencyPages[1])
		if errors.Join(err1, err2) != nil {
			slog.Error("found dependency line that is not of form '[int]|[int]'", "offending line", dependency)
		}

		// If afterPage is not already in the dependency graph, add a new list
		if _, ok := dependencyGraph[afterPage]; !ok {
			dependencyGraph[afterPage] = make([]int, 0)
		}
		dependencyGraph[afterPage] = append(dependencyGraph[afterPage], beforePage)
		slog.Debug("parsed dependency", "line", line, "updated dependency list", dependencyGraph[afterPage])
	}

	return dependencyGraph
}

func parseUpdateLine(updateLine []string) []int {
	updatePagesList := make([]int, len(updateLine))
	for i := 0; i < len(updateLine); i += 1 {
		updatePagesList[i], _ = strconv.Atoi(updateLine[i])
	}
	return updatePagesList
}

func Part01(fileScanner *bufio.Scanner) (int, error) {
	return 0, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	return 0, nil
}
