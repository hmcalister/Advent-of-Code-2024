package main

import (
	"bufio"
	"log/slog"
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

func (antennaMap *AntennaMap) CountAntinodes() int {
	validAntinodes := hashset.New[Coordinate]()
	for frequency := range antennaMap.antennaFrequencyLocations {
		frequencyAntinodes := antennaMap.countAntinodesOfFrequency(frequency)
		slog.Debug("found antinodes of frequency", "frequency", frequency, "antinodes", frequencyAntinodes)
		for _, antinode := range frequencyAntinodes {
			validAntinodes.Add(antinode)
		}
	}
	return validAntinodes.Size()
}

// Count the antinodes of a given frequency, returning the valid (inbound) coordinates
// This function does not mutate any attributes of the AntennaMap and is hence concurrency safe
func (antennaMap *AntennaMap) countAntinodesOfFrequency(frequency rune) []Coordinate {
	frequencyCoordinates, ok := antennaMap.antennaFrequencyLocations[frequency]
	validAntinodes := make([]Coordinate, 0)
	if !ok {
		slog.Debug("requested frequency not found", "frequency", frequency)
		return validAntinodes
	}

	for c1 := 0; c1 < len(frequencyCoordinates); c1 += 1 {
		for c2 := c1 + 1; c2 < len(frequencyCoordinates); c2 += 1 {
			coord1 := frequencyCoordinates[c1]
			coord2 := frequencyCoordinates[c2]
			a1 := determineAntinode(coord1, coord2)
			if a1.InBounds(antennaMap.width, antennaMap.height) {
				// slog.Debug("found valid antinode", "frequency", frequency, "coord1", coord1, "coord2", coord2, "antinode", a1)
				validAntinodes = append(validAntinodes, a1)
			}

			a2 := determineAntinode(frequencyCoordinates[c2], frequencyCoordinates[c1])
			if a2.InBounds(antennaMap.width, antennaMap.height) {
				// slog.Debug("found valid antinode", "frequency", frequency, "coord1", coord2, "coord2", coord1, "antinode", a2)
				validAntinodes = append(validAntinodes, a2)
			}
		}
	}

	return validAntinodes
}
