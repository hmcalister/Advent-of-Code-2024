package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"hmcalister/AdventOfCode/gridutils"
	"hmcalister/AdventOfCode/warehouse"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"log/slog"
	"os"
	"runtime/pprof"
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

	frames := make([]string, 0)
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
			frames = append(frames, warehouseMap.String())
		}
	}

	createGIF(frames, "part01.gif", 12)
	return warehouseMap.ComputeGPS(), nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	warehouseMapStrs := make([]string, 0)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if len(line) == 0 {
			break
		}
		warehouseMapStrs = append(warehouseMapStrs, line)
	}
	warehouseMap := warehouse.NewDoubleWidthWarehouseMap(warehouseMapStrs)
	fmt.Println(warehouseMap)

	frames := make([]string, 0)
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
			frames = append(frames, warehouseMap.String())
		}
	}

	createGIF(frames, "part02.gif", 12)
	return warehouseMap.ComputeGPS(), nil
}

// Turn the warehouse map string into an image
func stringToImage(s string) *image.Paletted {
	cellSize := 32
	borderWidth := 4
	backgroundColor := color.RGBA{255, 255, 255, 255}
	wallColor := color.RGBA{0, 0, 0, 255}
	boxBorderColor := color.RGBA{96, 96, 96, 255}
	boxColor := color.RGBA{128, 128, 128, 255}
	robotColor := color.RGBA{128, 196, 128, 255}

	lines := strings.Split(s, "\n")

	img := image.NewRGBA(
		image.Rect(
			0,
			0,
			cellSize*len(lines[0]),
			cellSize*(len(lines)-1),
		),
	)
	draw.Draw(img, img.Bounds(), &image.Uniform{backgroundColor}, image.Point{}, draw.Src)

	for y, line := range lines {
		for x, cell := range line {
			var currentCellColor color.RGBA
			switch cell {
			case warehouse.ROBOT_RUNE:
				currentCellColor = robotColor
			case warehouse.WALL_RUNE:
				currentCellColor = wallColor
			case warehouse.BOX_RUNE:
				currentCellColor = boxColor
			case '[':
				// DUAL BOX!!!
				currentCellColor = boxColor
				borderRect := image.Rect(
					cellSize*x,
					cellSize*y,
					cellSize*(x+2),
					cellSize*(y+1),
				)
				draw.Draw(img, borderRect, &image.Uniform{boxBorderColor}, image.Point{}, draw.Src)
				cellRect := image.Rect(
					cellSize*x+borderWidth,
					cellSize*y+borderWidth,
					cellSize*(x+2)-borderWidth,
					cellSize*(y+1)-borderWidth,
				)
				draw.Draw(img, cellRect, &image.Uniform{boxColor}, image.Point{}, draw.Src)
				continue
			case ']':
				continue
			default:
				currentCellColor = backgroundColor
			}
			cellRect := image.Rect(cellSize*x, cellSize*y, cellSize*(x+1), cellSize*(y+1))
			draw.Draw(img, cellRect, &image.Uniform{currentCellColor}, image.Point{}, draw.Src)
		}
	}

	pal := []color.Color{backgroundColor, wallColor, boxColor, boxBorderColor, robotColor}
	palettedImg := image.NewPaletted(img.Bounds(), pal)
	draw.FloydSteinberg.Draw(palettedImg, img.Bounds(), img, image.Point{})

	return palettedImg
}

func createGIF(frameStrings []string, outputFilePath string, delay int) error {
	frames := make([]*image.Paletted, 0)
	delays := make([]int, 0)

	for frameIndex, frameString := range frameStrings {
		slog.Debug("converting frame", "frame index", frameIndex, "total frames", len(frameStrings))
		img := stringToImage(frameString)
		frames = append(frames, img)
		delays = append(delays, delay)
	}

	f, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return gif.EncodeAll(f, &gif.GIF{
		Image: frames,
		Delay: delays,
	})
}
