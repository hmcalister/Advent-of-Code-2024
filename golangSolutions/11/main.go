package main

import (
	"bufio"
	"flag"
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

func Part01(fileScanner *bufio.Scanner) (int, error) {
	fileScanner.Scan()
	line := fileScanner.Text()
	stoneValueStrings := strings.Split(line, " ")

	stoneCountMap := NewStoneCounter()
	for _, stoneValueString := range stoneValueStrings {
		stoneValue, err := strconv.Atoi(stoneValueString)
		if err != nil {
			slog.Error("could not parse stone value string", "stone value string", stoneValueString, "error", err)
			continue
		}

		slog.Debug("added next stone", "stone value", stoneValue)
		stoneCountMap.AddStone(stoneValue, 1)
	}

	for i := 1; i <= 25; i += 1 {
		stoneCountMap.Blink()
		slog.Debug("finished blink", "blink index", i, "num stones", stoneCountMap.NumStones())
	}

	return stoneCountMap.NumStones(), nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	fileScanner.Scan()
	line := fileScanner.Text()
	stoneValueStrings := strings.Split(line, " ")

	stoneCountMap := NewStoneCounter()
	for _, stoneValueString := range stoneValueStrings {
		stoneValue, err := strconv.Atoi(stoneValueString)
		if err != nil {
			slog.Error("could not parse stone value string", "stone value string", stoneValueString, "error", err)
			continue
		}

		slog.Debug("added next stone", "stone value", stoneValue)
		stoneCountMap.AddStone(stoneValue, 1)
	}

	for i := 1; i <= 75; i += 1 {
		stoneCountMap.Blink()
		slog.Debug("finished blink", "blink index", i, "num stones", stoneCountMap.NumStones())
	}

	return stoneCountMap.NumStones(), nil
}
