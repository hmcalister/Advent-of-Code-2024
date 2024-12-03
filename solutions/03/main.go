package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	total := 0

	mulInstructionExpression := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	doInstructionExpression := regexp.MustCompile(`do\(\)`)
	dontInstructionExpression := regexp.MustCompile(`don't\(\)`)

	// Read all lines into a single string to avoid resetting enabled logic between lines
	allLinesArray := make([]string, 0)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		slog.Debug("line read", "line", line)
		allLinesArray = append(allLinesArray, line)
	}
	allLines := strings.Join(allLinesArray, " ")
	slog.Debug("all lines combined", "all lines", allLines)

	// Track the index of the next `don't()` occurrence
	var nextDisableIndex int

	// Track the remaining input (trim off the prefix that is either disabled or already matched to avoid additional computation)
	remainingInput := allLines
	for {
		slog.Debug("matching loop", "remaining input", remainingInput)

		// Find the next `don't()` occurrence
		disableIndexMatch := dontInstructionExpression.FindStringIndex(remainingInput)
		slog.Debug("disable match applied", "match results", disableIndexMatch)
		if disableIndexMatch == nil {
			// No disable found, carry on to end of slice
			nextDisableIndex = len(remainingInput)
		} else {
			nextDisableIndex = disableIndexMatch[1]
		}

		// Pick out the prefix up to the next `don't()` call, which gives the current enabled section
		enabledFragment := remainingInput[:nextDisableIndex]
		slog.Debug("enabled fragment found", "fragment", enabledFragment)

		// Apply same logic as part 1 to find matches
		foundMulMatches := mulInstructionExpression.FindAllStringSubmatch(enabledFragment, -1)
		slog.Debug("regex applied", "matches", foundMulMatches)

		for _, match := range foundMulMatches {
			v1, _ := strconv.Atoi(match[1])
			v2, _ := strconv.Atoi(match[2])
			total += v1 * v2
			slog.Debug("match calculation", "v1", v1, "v2", v2, "mul", v1*v2, "new total", total)
		}

		// Trim off the currently enabled section which has already been processed
		remainingInput = remainingInput[nextDisableIndex:]

		// Find the next `do()` occurrence, which enables the section again
		enableIndexMatch := doInstructionExpression.FindStringIndex(remainingInput)
		slog.Debug("enable match applied", "match results", enableIndexMatch)
		if enableIndexMatch == nil {
			// No enable found, we are done
			break
		}

		// Trim off the disabled section
		remainingInput = remainingInput[enableIndexMatch[0]:]

	}

	return total, nil
}
