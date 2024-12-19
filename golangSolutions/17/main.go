package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"hmcalister/AdventOfCode/tribitemulator"
	"log/slog"
	"os"
	"runtime/pprof"
	"slices"
	"strconv"
	"strings"
	"time"

	arrayqueue "github.com/hmcalister/Go-DSA/queue/ArrayQueue"
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

func parseInputToProgramAndRegisters(fileScanner *bufio.Scanner) ([]int, int, int, int) {
	// read register A
	if !fileScanner.Scan() {
		slog.Error("could not find input line for register A")
		os.Exit(1)
	}
	registerALine := strings.Split(fileScanner.Text(), ": ")
	if len(registerALine) != 2 {
		slog.Error("register A input does not match expected format", "register A line", registerALine)
		os.Exit(1)
	}
	registerA, err := strconv.Atoi(registerALine[1])
	if err != nil {
		slog.Error("could not parse register A value to integer", "register A line", registerALine)
	}

	// read register B
	if !fileScanner.Scan() {
		slog.Error("could not find input line for register B")
		os.Exit(1)
	}
	registerBLine := strings.Split(fileScanner.Text(), ": ")
	if len(registerBLine) != 2 {
		slog.Error("register B input does not match expected format", "register B line", registerBLine)
		os.Exit(1)
	}
	registerB, err := strconv.Atoi(registerBLine[1])
	if err != nil {
		slog.Error("could not parse register B value to integer", "register B line", registerBLine)
	}

	// read register C
	if !fileScanner.Scan() {
		slog.Error("could not find input line for register C")
		os.Exit(1)
	}
	registerCLine := strings.Split(fileScanner.Text(), ": ")
	if len(registerCLine) != 2 {
		slog.Error("register C input does not match expected format", "register C line", registerCLine)
		os.Exit(1)
	}
	registerC, err := strconv.Atoi(registerCLine[1])
	if err != nil {
		slog.Error("could not parse register C value to integer", "register C line", registerCLine)
	}

	// read program
	fileScanner.Scan()
	if !fileScanner.Scan() {
		slog.Error("could not find input line for program")
		os.Exit(1)
	}
	programLine := strings.Split(fileScanner.Text(), ": ")
	if len(programLine) != 2 {
		slog.Error("program input does not match expected format", "program line", programLine)
		os.Exit(1)
	}
	programValueStrs := strings.Split(programLine[1], ",")
	program := make([]int, len(programValueStrs))
	for index, valueStr := range programValueStrs {
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			slog.Error("could not parse program value to integer", "program index", index, "program value", valueStr)
		}
		program[index] = value
	}

	return program, registerA, registerB, registerC
}

func Part01(fileScanner *bufio.Scanner) (int, error) {
	program, registerA, registerB, registerC := parseInputToProgramAndRegisters(fileScanner)
	slog.Debug("parsed input", "program", program, "registerA", registerA, "registerB", registerB, "registerC", registerC)
	emulator := tribitemulator.NewTribitEmulator(registerA, registerB, registerC)
	output := emulator.ExecuteProgram(program)
	slog.Info("program output", "output", output)

	if len(output) > 0 {
		for index := 0; index < len(output)-1; index += 1 {
			fmt.Print(output[index])
			fmt.Print(",")
		}
		fmt.Println(output[len(output)-1])
	}

	return 0, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	program, registerA, registerB, registerC := parseInputToProgramAndRegisters(fileScanner)
	slog.Debug("parsed input", "program", program, "registerA", registerA, "registerB", registerB, "registerC", registerC)

	type registerSearchData struct {
		initialAValue         int
		nextSuffixMatchLength int
	}
	registerSearchQueue := arrayqueue.New[registerSearchData]()
	registerSearchQueue.Add(registerSearchData{0, 1})

	for registerSearchQueue.Size() > 0 {
		currentRegisterSearch, _ := registerSearchQueue.Remove()
		for initialAValue := currentRegisterSearch.initialAValue; initialAValue < currentRegisterSearch.initialAValue+8; initialAValue += 1 {
			emulator := tribitemulator.NewTribitEmulator(initialAValue, registerB, registerC)
			output := emulator.ExecuteProgram(program)
			slog.Debug("register search loop", "initial a value", initialAValue, "output", output, "program prefix", program[len(program)-currentRegisterSearch.nextSuffixMatchLength:])
			if slices.Equal(output, program[len(program)-currentRegisterSearch.nextSuffixMatchLength:]) {
				if currentRegisterSearch.nextSuffixMatchLength == len(program) {
					return initialAValue, nil
				}
				registerSearchQueue.Add(registerSearchData{initialAValue * 8, currentRegisterSearch.nextSuffixMatchLength + 1})
			}
		}
	}

	return -1, errors.New("did not find any matching prefix value")
}
