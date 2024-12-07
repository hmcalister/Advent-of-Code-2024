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

func (c *CalibrationData) IsValidPart01() bool {
	return c.isValidPart01Recursive(0, 0)
}

func (c *CalibrationData) isValidPart01Recursive(currentValue int, currentIndex int) bool {
	// If we are at the end, then evaluate that current value is exactly equal to the target
	if currentIndex == len(c.EquationData) {
		return currentValue == c.TargetNumber
	}

	// Since part 01 is only concerned with addition and multiplication, the current value can never decrease
	if currentValue > c.TargetNumber {
		return false
	}

	currentEquationData := c.EquationData[currentIndex]
	return c.isValidPart01Recursive(currentValue+currentEquationData, currentIndex+1) ||
		c.isValidPart01Recursive(currentValue*currentEquationData, currentIndex+1)
}

func (c *CalibrationData) IsValidPart02() bool {
	return c.isValidPart02Recursive(0, 0)
}

func (c *CalibrationData) isValidPart02Recursive(currentValue int, currentIndex int) bool {
	// If we are at the end, then evaluate that current value is exactly equal to the target
	if currentIndex == len(c.EquationData) {
		return currentValue == c.TargetNumber
	}

	// Since addition, multiplication, and concatenation are monotonically increasing, the current value can never decrease
	if currentValue > c.TargetNumber {
		return false
	}

	currentEquationData := c.EquationData[currentIndex]
	return c.isValidPart02Recursive(currentValue+currentEquationData, currentIndex+1) ||
		c.isValidPart02Recursive(currentValue*currentEquationData, currentIndex+1) ||
		c.isValidPart02Recursive(concatIntegers(currentValue, currentEquationData), currentIndex+1)
}

func concatIntegers(a, b int) int {
	bLen := len(strconv.FormatInt(int64(b), 10))
	concatValue := int(math.Pow10(bLen))*a + b
	return concatValue
}
