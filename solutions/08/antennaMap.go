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

