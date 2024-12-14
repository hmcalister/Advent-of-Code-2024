package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"hmcalister/AdventOfCode/robot"
	"log/slog"
	"os"
	"regexp"
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

func parseInputToRobots(fileScanner *bufio.Scanner) []*robot.Robot {
	robotRegex := regexp.MustCompile(`p=([+-]?\d*),([+-]?\d*) v=([+-]?\d*),([+-]?\d*)`)
	robots := make([]*robot.Robot, 0)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		matches := robotRegex.FindStringSubmatch(line)
		if matches == nil || len(matches) != 5 {
			slog.Error("line does not match expected robot data format", "line", line, "matches", matches)
			continue
		}

		px, pxErr := strconv.Atoi(matches[1])
		py, pyErr := strconv.Atoi(matches[2])
		vx, vxErr := strconv.Atoi(matches[3])
		vy, vyErr := strconv.Atoi(matches[4])
		if errors.Join(pxErr, pyErr, vxErr, vyErr) != nil {
			slog.Error("could not parse matched groups to robot data",
				"line", line,
				"matches", matches,
				"position x error", pxErr,
				"position y error", pyErr,
				"velocity x error", vxErr,
				"velocity y error", vyErr,
			)
			continue
		}

		newRobot := robot.NewRobot(
			robot.Vector2{X: px, Y: py},
			robot.Vector2{X: vx, Y: vy},
		)
		slog.Debug("parsed robot", "line", line, "robot", *newRobot)
		robots = append(robots, newRobot)
	}

	return robots
}

func Part01(fileScanner *bufio.Scanner) (int, error) {
	return 0, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	return 0, nil
}
