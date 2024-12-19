package main

import (
	"bufio"
	"errors"
	"flag"
	"hmcalister/AdventOfCode/clawmachine"
	"log/slog"
	"os"
	"regexp"
	"runtime/pprof"
	"strconv"
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

func parseInputToClawMachines(fileScanner *bufio.Scanner) []*clawmachine.ClawMachine {
	machines := make([]*clawmachine.ClawMachine, 0)

	buttonARegex := regexp.MustCompile(`Button A: X([+-]\d*), Y([+-]\d*)`)
	buttonBRegex := regexp.MustCompile(`Button B: X([+-]\d*), Y([+-]\d*)`)
	prizeRegex := regexp.MustCompile(`Prize: X=([+-]?\d*), Y=([+-]?\d*)`)

	// We still have some input
	for fileScanner.Scan() {
		buttonALine := fileScanner.Text()
		if !fileScanner.Scan() {
			slog.Debug("unexpected EOF when reading button B line")
			return machines
		}
		buttonBLine := fileScanner.Text()
		if !fileScanner.Scan() {
			slog.Debug("unexpected EOF when reading prize line")
			return machines
		}
		prizeLine := fileScanner.Text()
		fileScanner.Scan()

		slog.Debug("next claw machine input read", "button A line", buttonALine, "button B line", buttonBLine, "prize line", prizeLine)

		buttonAMatches := buttonARegex.FindStringSubmatch(buttonALine)
		if buttonAMatches == nil || len(buttonAMatches) != 3 {
			slog.Error("button A line not of expected format", "button A line", buttonALine, "button A regex matches", buttonAMatches)
			continue
		}
		buttonAX, errA := strconv.ParseFloat(buttonAMatches[1], 64)
		buttonAY, errB := strconv.ParseFloat(buttonAMatches[2], 64)
		if errors.Join(errA, errB) != nil {
			slog.Error("could not parse button A line", "button A line", buttonALine, "button A regex matches", buttonAMatches)
			continue
		}

		buttonBMatches := buttonBRegex.FindStringSubmatch(buttonBLine)
		if buttonBMatches == nil || len(buttonBMatches) != 3 {
			slog.Error("button B line not of expected format", "button B line", buttonBLine, "button B regex matches", buttonBMatches)
			continue
		}
		buttonBX, errA := strconv.ParseFloat(buttonBMatches[1], 64)
		buttonBY, errB := strconv.ParseFloat(buttonBMatches[2], 64)
		if errors.Join(errA, errB) != nil {
			slog.Error("could not parse button B line", "button B line", buttonBLine, "button B regex matches", buttonBMatches)
			continue
		}

		prizeMatches := prizeRegex.FindStringSubmatch(prizeLine)
		if prizeMatches == nil || len(prizeMatches) != 3 {
			slog.Error("prize line not of expected format", "prize line", prizeLine, "prize regex matches", prizeMatches)
			continue
		}
		prizeX, errA := strconv.ParseFloat(prizeMatches[1], 64)
		prizeY, errB := strconv.ParseFloat(prizeMatches[2], 64)
		if errors.Join(errA, errB) != nil {
			slog.Error("could not parse prize line", "prize line", prizeLine, "prize regex matches", prizeMatches)
		}

		nextClawMachine := clawmachine.NewClawMachine(
			buttonAX,
			buttonAY,
			buttonBX,
			buttonBY,
			prizeX,
			prizeY,
		)
		slog.Debug("input parsed to claw machine", "claw machine", nextClawMachine)
		machines = append(machines, nextClawMachine)
	}

	return machines
}

func Part01(fileScanner *bufio.Scanner) (int, error) {
	clawMachines := parseInputToClawMachines(fileScanner)
	totalCost := 0
	validMachines := 0
	for _, machine := range clawMachines {
		cost, err := machine.ComputeLowestTokenCost()
		if err != nil {
			slog.Error("error when computing token cost", "error", err)
			continue
		}
		totalCost += cost
		validMachines += 1
		slog.Debug("found lowest token cost", "current machine token cost", cost, "updated total token cost", totalCost)

	}
	slog.Debug("finished computing claw results", "total claw machines", len(clawMachines), "valid claw machines", validMachines)
	return totalCost, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	clawMachines := parseInputToClawMachines(fileScanner)
	totalCost := 0
	validMachines := 0
	for _, machine := range clawMachines {
		machine.FixUnitConversion()
		cost, err := machine.ComputeLowestTokenCost()
		if err != nil {
			slog.Error("error when computing token cost", "error", err)
			continue
		}
		totalCost += cost
		validMachines += 1
		slog.Debug("found lowest token cost", "current machine token cost", cost, "updated total token cost", totalCost)

	}
	slog.Debug("finished computing claw results", "total claw machines", len(clawMachines), "valid claw machines", validMachines)
	return totalCost, nil
}
