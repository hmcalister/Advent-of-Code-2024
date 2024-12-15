package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"hmcalister/AdventOfCode/gridutils"
	"hmcalister/AdventOfCode/warehouse"
	"log/slog"
	"os"
	"runtime/pprof"
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
	warehouseMapStrs := make([]string, 0)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if len(line) == 0 {
			break
		}
		warehouseMapStrs = append(warehouseMapStrs, line)
	}
	warehouseMap := warehouse.NewSingleWidthWarehouseMap(warehouseMapStrs)
	fmt.Println(warehouseMap)

	var robotStepDirection gridutils.Direction
	for fileScanner.Scan() {
		line := fileScanner.Text()
		for _, robotStepDirectionRune := range line {
			switch robotStepDirectionRune {
			case '^':
				robotStepDirection = gridutils.DIRECTION_UP
			case '>':
				robotStepDirection = gridutils.DIRECTION_RIGHT
			case 'v':
				robotStepDirection = gridutils.DIRECTION_DOWN
			case '<':
				robotStepDirection = gridutils.DIRECTION_LEFT
			default:
				slog.Error("unexpected robot direction encountered", "rune found", robotStepDirectionRune)
				return 0, errors.New("could not parse robot direction input")
			}
			slog.Debug("robot moving", "rune found", robotStepDirectionRune, "robot direction", robotStepDirection)
			warehouseMap.RobotStep(robotStepDirection)
			// fmt.Println(warehouseMap)
		}
	}

	return warehouseMap.ComputeGPS(), nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	return 0, nil
}
