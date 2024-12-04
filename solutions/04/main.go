package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
	"time"
)

var (
	TARGET_WORD_BYTES = []byte{'X', 'M', 'A', 'S'}
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

func parseInputToByteArray(fileScanner *bufio.Scanner) [][]byte {
	byteArray := make([][]byte, 0)
	for fileScanner.Scan() {
		byteArray = append(byteArray, []byte(fileScanner.Text()))
	}

	for _, row := range byteArray {
		slog.Debug("byteArray", "row", string(row))
	}
	return byteArray
}

func findTargetWords(byteArray [][]byte, rowStartIndex, colStartIndex int) int {
	// Check written along row, col, diagonal, forward and backwards

	numRows, numCols := len(byteArray), len(byteArray[0])
	bufferRequired := len(TARGET_WORD_BYTES) - 1

	totalFound := 0

	// Row Forward
	if colStartIndex+bufferRequired < numCols {
		if byteArray[rowStartIndex][colStartIndex+1] == TARGET_WORD_BYTES[1] &&
			byteArray[rowStartIndex][colStartIndex+2] == TARGET_WORD_BYTES[2] &&
			byteArray[rowStartIndex][colStartIndex+3] == TARGET_WORD_BYTES[3] {
			slog.Debug("found target word", "rowStartIndex", rowStartIndex, "colStartIndex", colStartIndex, "direction", "Row Forward")
			totalFound += 1
		}
	}

	// Row Backward
	if colStartIndex-bufferRequired >= 0 {
		if byteArray[rowStartIndex][colStartIndex-1] == TARGET_WORD_BYTES[1] &&
			byteArray[rowStartIndex][colStartIndex-2] == TARGET_WORD_BYTES[2] &&
			byteArray[rowStartIndex][colStartIndex-3] == TARGET_WORD_BYTES[3] {
			slog.Debug("found target word", "rowStartIndex", rowStartIndex, "colStartIndex", colStartIndex, "direction", "Row Backward")
			totalFound += 1
		}
	}

	// Col Forward
	if rowStartIndex+bufferRequired < numRows {
		if byteArray[rowStartIndex+1][colStartIndex] == TARGET_WORD_BYTES[1] &&
			byteArray[rowStartIndex+2][colStartIndex] == TARGET_WORD_BYTES[2] &&
			byteArray[rowStartIndex+3][colStartIndex] == TARGET_WORD_BYTES[3] {
			slog.Debug("found target word", "rowStartIndex", rowStartIndex, "colStartIndex", colStartIndex, "direction", "Col Forward")
			totalFound += 1
		}
	}

	// Col Backward
	if rowStartIndex-bufferRequired >= 0 {
		if byteArray[rowStartIndex-1][colStartIndex] == TARGET_WORD_BYTES[1] &&
			byteArray[rowStartIndex-2][colStartIndex] == TARGET_WORD_BYTES[2] &&
			byteArray[rowStartIndex-3][colStartIndex] == TARGET_WORD_BYTES[3] {
			slog.Debug("found target word", "rowStartIndex", rowStartIndex, "colStartIndex", colStartIndex, "direction", "Col Backward")
			totalFound += 1
		}
	}

	// Diagonal Down Right
	if (colStartIndex+bufferRequired < numCols) && (rowStartIndex+bufferRequired < numRows) {
		if byteArray[rowStartIndex+1][colStartIndex+1] == TARGET_WORD_BYTES[1] &&
			byteArray[rowStartIndex+2][colStartIndex+2] == TARGET_WORD_BYTES[2] &&
			byteArray[rowStartIndex+3][colStartIndex+3] == TARGET_WORD_BYTES[3] {
			slog.Debug("found target word", "rowStartIndex", rowStartIndex, "colStartIndex", colStartIndex, "direction", "Diagonal Down Right")
			totalFound += 1
		}
	}

	// Diagonal Up Right
	if (colStartIndex+bufferRequired < numCols) && (rowStartIndex-bufferRequired >= 0) {
		if byteArray[rowStartIndex-1][colStartIndex+1] == TARGET_WORD_BYTES[1] &&
			byteArray[rowStartIndex-2][colStartIndex+2] == TARGET_WORD_BYTES[2] &&
			byteArray[rowStartIndex-3][colStartIndex+3] == TARGET_WORD_BYTES[3] {
			slog.Debug("found target word", "rowStartIndex", rowStartIndex, "colStartIndex", colStartIndex, "direction", "Diagonal Up Right")
			totalFound += 1
		}
	}

	// Diagonal Down Left
	if (colStartIndex-bufferRequired >= 0) && (rowStartIndex+bufferRequired < numRows) {
		if byteArray[rowStartIndex+1][colStartIndex-1] == TARGET_WORD_BYTES[1] &&
			byteArray[rowStartIndex+2][colStartIndex-2] == TARGET_WORD_BYTES[2] &&
			byteArray[rowStartIndex+3][colStartIndex-3] == TARGET_WORD_BYTES[3] {
			slog.Debug("found target word", "rowStartIndex", rowStartIndex, "colStartIndex", colStartIndex, "direction", "Diagonal Down Left")
			totalFound += 1
		}
	}

	// Diagonal Up Left
	if (colStartIndex-bufferRequired >= 0) && (rowStartIndex-bufferRequired >= 0) {
		if byteArray[rowStartIndex-1][colStartIndex-1] == TARGET_WORD_BYTES[1] &&
			byteArray[rowStartIndex-2][colStartIndex-2] == TARGET_WORD_BYTES[2] &&
			byteArray[rowStartIndex-3][colStartIndex-3] == TARGET_WORD_BYTES[3] {
			slog.Debug("found target word", "rowStartIndex", rowStartIndex, "colStartIndex", colStartIndex, "direction", "Diagonal Up Left")
			totalFound += 1
		}
	}

	return totalFound
}

func Part01(fileScanner *bufio.Scanner) (int, error) {

	input := parseInputToByteArray(fileScanner)
	slog.Debug("parsed input", "byte array", input)

	totalTargetWords := 0
	for rowIndex, row := range input {
		for colIndex, cell := range row {
			if cell == TARGET_WORD_BYTES[0] {
				// slog.Debug("found start of target word", "rowIndex", rowIndex, "colIndex", colIndex)
				totalTargetWords += findTargetWords(input, rowIndex, colIndex)
			}
		}
	}

	return totalTargetWords, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	input := parseInputToByteArray(fileScanner)
	slog.Debug("parsed input", "byte array", input)

	totalCrossedMas := 0
	for rowIndex := 1; rowIndex < len(input)-1; rowIndex += 1 {
		for colIndex := 1; colIndex < len(input[0])-1; colIndex += 1 {
			if input[rowIndex][colIndex] == 'A' {
				slog.Debug("found center of crossed mas", "rowIndex", rowIndex, "colIndex", colIndex)

				// M.M
				// .A.
				// S.S
				if input[rowIndex-1][colIndex-1] == 'M' &&
					input[rowIndex-1][colIndex+1] == 'M' &&
					input[rowIndex+1][colIndex-1] == 'S' &&
					input[rowIndex+1][colIndex+1] == 'S' {
					slog.Debug("found cross", "rowIndex", rowIndex, "colIndex", colIndex, "orientation", "0")
					totalCrossedMas += 1
				}

				// S.M
				// .A.
				// S.M
				if input[rowIndex-1][colIndex-1] == 'S' &&
					input[rowIndex-1][colIndex+1] == 'M' &&
					input[rowIndex+1][colIndex-1] == 'S' &&
					input[rowIndex+1][colIndex+1] == 'M' {
					slog.Debug("found cross", "rowIndex", rowIndex, "colIndex", colIndex, "orientation", "0")
					totalCrossedMas += 1
				}

				// S.S
				// .A.
				// M.M
				if input[rowIndex-1][colIndex-1] == 'S' &&
					input[rowIndex-1][colIndex+1] == 'S' &&
					input[rowIndex+1][colIndex-1] == 'M' &&
					input[rowIndex+1][colIndex+1] == 'M' {
					slog.Debug("found cross", "rowIndex", rowIndex, "colIndex", colIndex, "orientation", "0")
					totalCrossedMas += 1
				}

				// M.S
				// .A.
				// M.S
				if input[rowIndex-1][colIndex-1] == 'M' &&
					input[rowIndex-1][colIndex+1] == 'S' &&
					input[rowIndex+1][colIndex-1] == 'M' &&
					input[rowIndex+1][colIndex+1] == 'S' {
					slog.Debug("found cross", "rowIndex", rowIndex, "colIndex", colIndex, "orientation", "0")
					totalCrossedMas += 1
				}
			}
		}
	}

	return totalCrossedMas, nil
}
