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

	hashset "github.com/hmcalister/Go-DSA/set/HashSet"
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
	if !fileScanner.Scan() {
		slog.Error("no input found")
		return 0, errors.New("no input found")
	}
	gridSizeLine := fileScanner.Text()
	gridSizeStrs := strings.Split(gridSizeLine, ",")
	if len(gridSizeStrs) != 2 {
		slog.Error("did not find an initial line with format x,y for grid size", "line", gridSizeLine, "split line", gridSizeStrs)
		return 0, errors.New("grid size line does not match expected format x,y")
	}
	gridX, gridXErr := strconv.Atoi(gridSizeStrs[0])
	gridY, gridYErr := strconv.Atoi(gridSizeStrs[1])
	if errors.Join(gridXErr, gridYErr) != nil {
		slog.Error("could not parse grid line to integers", "line", gridSizeLine, "grid size strs", gridSizeStrs)
		return 0, errors.New("grid size line does not contain integers")
	}

	robots := parseInputToRobots(fileScanner)
	// slog.Debug("parsed input", "gridX", gridX, "gridY", gridY, "robots", robots)

	quadrantCounts := []int{0, 0, 0, 0}
	coordinateCounts := make(map[robot.Vector2]int)
	for _, robot := range robots {
		nextPosition := robot.ComputePosition(gridX, gridY, 100)

		if nextPosition.X < gridX/2 && nextPosition.Y < gridY/2 {
			quadrantCounts[0] += 1
		} else if nextPosition.X > gridX/2 && nextPosition.Y < gridY/2 {
			quadrantCounts[1] += 1
		} else if nextPosition.X < gridX/2 && nextPosition.Y > gridY/2 {
			quadrantCounts[2] += 1
		} else if nextPosition.X > gridX/2 && nextPosition.Y > gridY/2 {
			quadrantCounts[3] += 1
		}
		coordinateCounts[nextPosition] += 1
		slog.Debug("robot stepped", "robot", *robot, "next position", nextPosition, "updated quadrant counts", quadrantCounts)
	}

	for y := 0; y < gridY; y += 1 {
		for x := 0; x < gridX; x += 1 {
			coordinate := robot.Vector2{X: x, Y: y}
			coordinateCount, ok := coordinateCounts[coordinate]
			if !ok {
				fmt.Print(".")
			} else {
				fmt.Print(coordinateCount)
			}
		}
		fmt.Println()
	}

	quadrantCountProduct := 1
	for _, count := range quadrantCounts {
		quadrantCountProduct *= count
	}

	return quadrantCountProduct, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	if !fileScanner.Scan() {
		slog.Error("no input found")
		return 0, errors.New("no input found")
	}
	gridSizeLine := fileScanner.Text()
	gridSizeStrs := strings.Split(gridSizeLine, ",")
	if len(gridSizeStrs) != 2 {
		slog.Error("did not find an initial line with format x,y for grid size", "line", gridSizeLine, "split line", gridSizeStrs)
		return 0, errors.New("grid size line does not match expected format x,y")
	}
	gridX, gridXErr := strconv.Atoi(gridSizeStrs[0])
	gridY, gridYErr := strconv.Atoi(gridSizeStrs[1])
	if errors.Join(gridXErr, gridYErr) != nil {
		slog.Error("could not parse grid line to integers", "line", gridSizeLine, "grid size strs", gridSizeStrs)
		return 0, errors.New("grid size line does not contain integers")
	}

	robots := parseInputToRobots(fileScanner)
	// slog.Debug("parsed input", "gridX", gridX, "gridY", gridY, "robots", robots)

	keyboardScanner := bufio.NewScanner(os.Stdin)

	for stepIndex := 0; stepIndex < 10000; stepIndex += 1 {
		robotInCoordinate := hashset.New[robot.Vector2]()
		toPrint := true
		for _, robot := range robots {
			nextPosition := robot.ComputePosition(gridX, gridY, stepIndex)
			if robotInCoordinate.Contains(nextPosition) {
				toPrint = false
			}
			robotInCoordinate.Add(nextPosition)
			// slog.Debug("robot stepped", "robot", *robot, "next position", nextPosition, "updated quadrant counts", quadrantCounts)
		}

		if toPrint {
			fmt.Printf("\n\nStep Index: %v\n", stepIndex)
			for y := 0; y < gridY; y += 1 {
				for x := 0; x < gridX; x += 1 {
					coordinate := robot.Vector2{X: x, Y: y}
					if robotInCoordinate.Contains(coordinate) {
						fmt.Print("#")
					} else {
						fmt.Print(".")
					}
				}
				fmt.Println()
			}
			keyboardScanner.Scan()
		}
	}

	return 0, nil
}
