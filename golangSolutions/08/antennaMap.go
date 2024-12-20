package main

import (
	"bufio"
	"log/slog"
	"sync"
	"unicode"

	hashset "github.com/hmcalister/Go-DSA/set/HashSet"
)

type AntennaMap struct {
	rawMap                    [][]rune
	width                     int
	height                    int
	antennaFrequencyLocations map[rune][]Coordinate
}

func ParseInputToAntennaMap(fileScanner *bufio.Scanner) *AntennaMap {
	antennaMap := AntennaMap{
		rawMap:                    make([][]rune, 0),
		antennaFrequencyLocations: make(map[rune][]Coordinate),
	}

	for y := 0; fileScanner.Scan(); y += 1 {
		line := fileScanner.Text()
		slog.Debug("read line", "line", line)
		antennaMap.rawMap = append(antennaMap.rawMap, []rune(line))

		for x, r := range line {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				c := Coordinate{x, y}
				slog.Debug("found antenna", "coordinate", c, "frequency", r)
				// If this is the first antenna of this frequency, make a new list for it
				if _, ok := antennaMap.antennaFrequencyLocations[r]; !ok {
					antennaMap.antennaFrequencyLocations[r] = make([]Coordinate, 0)
				}
				antennaMap.antennaFrequencyLocations[r] = append(antennaMap.antennaFrequencyLocations[r], c)
			}
		}
	}

	antennaMap.height = len(antennaMap.rawMap)
	antennaMap.width = len(antennaMap.rawMap[0])
	slog.Debug("map parsed", "antenna map", antennaMap)

	return &antennaMap
}

func (antennaMap *AntennaMap) CountAntinodesPart01() int {
	validAntinodes := hashset.New[Coordinate]()
	for frequency := range antennaMap.antennaFrequencyLocations {
		frequencyAntinodes := antennaMap.countFirstOrderAntinodesOfFrequency(frequency)
		slog.Debug("found antinodes of frequency", "frequency", frequency, "antinodes", frequencyAntinodes)
		for _, antinode := range frequencyAntinodes {
			validAntinodes.Add(antinode)
		}
	}
	return validAntinodes.Size()
}

func (antennaMap *AntennaMap) CountAntinodesPart02() int {
	var validAntinodesMutex sync.Mutex
	var workerWaitGroup sync.WaitGroup
	validAntinodes := hashset.New[Coordinate]()
	for frequency := range antennaMap.antennaFrequencyLocations {
		workerWaitGroup.Add(1)
		go func() {
			defer workerWaitGroup.Done()
			frequencyAntinodes := antennaMap.countAllAntinodesOfFrequency(frequency)
			slog.Debug("found antinodes of frequency", "frequency", frequency, "antinodes", frequencyAntinodes)
			validAntinodesMutex.Lock()
			for _, antinode := range frequencyAntinodes {
				validAntinodes.Add(antinode)
			}
			validAntinodesMutex.Unlock()
		}()
	}

	workerWaitGroup.Wait()

	// DEBUG PRINT LOOP
	// for y, row := range antennaMap.rawMap {
	// 	for x, cell := range row {
	// 		if validAntinodes.Contains(Coordinate{x, y}) {
	// 			fmt.Print("#")
	// 		} else {
	// 			fmt.Print(string(cell))
	// 		}
	// 	}
	// 	fmt.Println()
	// }

	return validAntinodes.Size()
}

// Count the first order antinodes of a given frequency, returning the valid (inbound) coordinates
// This function does not mutate any attributes of the AntennaMap and is hence concurrency safe
func (antennaMap *AntennaMap) countFirstOrderAntinodesOfFrequency(frequency rune) []Coordinate {
	frequencyCoordinates, ok := antennaMap.antennaFrequencyLocations[frequency]
	validAntinodes := make([]Coordinate, 0)
	if !ok {
		slog.Debug("requested frequency not found", "frequency", frequency)
		return validAntinodes
	}

	for coordOneIndex := 0; coordOneIndex < len(frequencyCoordinates); coordOneIndex += 1 {
		coordOne := frequencyCoordinates[coordOneIndex]
		for coordTwoIndex := coordOneIndex + 1; coordTwoIndex < len(frequencyCoordinates); coordTwoIndex += 1 {
			coordTwo := frequencyCoordinates[coordTwoIndex]
			antinodeOne := determineFirstOrderAntinode(coordOne, coordTwo)
			if antinodeOne.InBounds(antennaMap.width, antennaMap.height) {
				validAntinodes = append(validAntinodes, antinodeOne)
			}

			antinodeTwo := determineFirstOrderAntinode(frequencyCoordinates[coordTwoIndex], frequencyCoordinates[coordOneIndex])
			if antinodeTwo.InBounds(antennaMap.width, antennaMap.height) {
				validAntinodes = append(validAntinodes, antinodeTwo)
			}
		}
	}

	return validAntinodes
}

// Count all antinodes of a given frequency, returning the valid (inbound) coordinates
// This function does not mutate any attributes of the AntennaMap and is hence concurrency safe
func (antennaMap *AntennaMap) countAllAntinodesOfFrequency(frequency rune) []Coordinate {
	frequencyCoordinates, ok := antennaMap.antennaFrequencyLocations[frequency]
	validAntinodes := make([]Coordinate, 0)
	if !ok {
		slog.Debug("requested frequency not found", "frequency", frequency)
		return validAntinodes
	}

	for coordOneIndex := 0; coordOneIndex < len(frequencyCoordinates); coordOneIndex += 1 {
		// Don't forget the antinode at the current position, i.e. dx=dy=0
		coordOne := frequencyCoordinates[coordOneIndex]
		validAntinodes = append(validAntinodes, coordOne)
		for coordTwoIndex := coordOneIndex + 1; coordTwoIndex < len(frequencyCoordinates); coordTwoIndex += 1 {
			coordTwo := frequencyCoordinates[coordTwoIndex]

			antinodeOneStep := coordOne.Subtract(coordTwo)
			antinodeOne := coordOne.Add(antinodeOneStep)
			for antinodeOne.InBounds(antennaMap.width, antennaMap.height) {
				validAntinodes = append(validAntinodes, antinodeOne)
				antinodeOne = antinodeOne.Add(antinodeOneStep)
			}

			antinodeTwoStep := coordTwo.Subtract(coordOne)
			antinodeTwo := coordTwo.Add(antinodeTwoStep)
			for antinodeTwo.InBounds(antennaMap.width, antennaMap.height) {
				validAntinodes = append(validAntinodes, antinodeTwo)
				antinodeTwo = antinodeTwo.Add(antinodeTwoStep)
			}
		}
	}

	return validAntinodes
}
