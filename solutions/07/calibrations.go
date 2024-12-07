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

