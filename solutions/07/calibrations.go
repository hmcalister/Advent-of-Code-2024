package main

import (
	"errors"
	"log/slog"
	"math"
	"strconv"
	"strings"
)

var (
	ErrorMalformedLine = errors.New("line not parsable to equation data")
)

type CalibrationData struct {
	TargetNumber int
	EquationData []int
}

// Line is expected to be of the form "(targetNumber): (equationData[0]) (equationData[1]) (equationData[2]) ..."
// where each group is an integer, e.g.
//
// 10: 2 4 8 16
func ParseLineToCalibrationData(line string) (*CalibrationData, error) {
	lineParts := strings.Split(line, ":")
	if len(lineParts) != 2 {
		return nil, ErrorMalformedLine
	}

	targetNumberStr := lineParts[0]
	targetNumber, err := strconv.Atoi(targetNumberStr)
	if err != nil {
		slog.Error("target number not parsable to integer", "line", line, "target number string", targetNumberStr)
		return nil, err
	}

	equationDataStrs := strings.Split(strings.TrimSpace(lineParts[1]), " ")
	equationData := make([]int, len(equationDataStrs))
	for i, s := range equationDataStrs {
		num, err := strconv.Atoi(s)
		if err != nil {
			slog.Error("equation data number not parsable to integer", "line", line, "equation data number string", s)
			return nil, err
		}
		equationData[i] = num
	}

	return &CalibrationData{
		TargetNumber: targetNumber,
		EquationData: equationData,
	}, nil
}

