package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
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
	totalCalibrationResult := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		calibrationData, err := ParseLineToCalibrationData(line)
		if err != nil {
			slog.Error("found error when parsing line", "line", line, "error", err)
			continue
		}

		if calibrationData.IsValidPart01() {
			totalCalibrationResult += calibrationData.TargetNumber
			slog.Debug("found valid calibration data", "calibration data", calibrationData, "updated total", totalCalibrationResult)
		}
	}
	return totalCalibrationResult, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	totalCalibrationResult := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		calibrationData, err := ParseLineToCalibrationData(line)
		if err != nil {
			slog.Error("found error when parsing line", "line", line, "error", err)
			continue
		}

		if calibrationData.IsValidPart02() {
			totalCalibrationResult += calibrationData.TargetNumber
			slog.Debug("found valid calibration data", "calibration data", calibrationData, "updated total", totalCalibrationResult)
		}
	}
	return totalCalibrationResult, nil
}
