package main

import "strconv"

type StoneCounter struct {
	// For a given stone value (key) store how many stones with that value exist (value)
	stoneCountMap map[int]int
}

func NewStoneCounter() *StoneCounter {
	return &StoneCounter{
		stoneCountMap: make(map[int]int),
	}
}

// Return the total number of stones
func (stoneCounter *StoneCounter) NumStones() int {
	totalStones := 0
	for _, count := range stoneCounter.stoneCountMap {
		totalStones += count
	}
	return totalStones
}

func (stoneCounter *StoneCounter) AddStone(stoneValue int, count int) {
	value, ok := stoneCounter.stoneCountMap[stoneValue]
	if !ok {
		stoneCounter.stoneCountMap[stoneValue] = count
		return
	}
	stoneCounter.stoneCountMap[stoneValue] = value + count
}

func (stoneCounter *StoneCounter) Blink() {
	newStoneCounter := NewStoneCounter()

	for currentStoneValue, count := range stoneCounter.stoneCountMap {
		if currentStoneValue == 0 {
			newStoneCounter.AddStone(1, count)
			continue
		}

		currentStoneValueString := strconv.FormatInt(int64(currentStoneValue), 10)
		if len(currentStoneValueString)%2 == 0 {
			leftStoneValueString := currentStoneValueString[:len(currentStoneValueString)/2]
			rightStoneValueString := currentStoneValueString[len(currentStoneValueString)/2:]
			leftStoneValue, _ := strconv.Atoi(leftStoneValueString)
			rightStoneValue, _ := strconv.Atoi(rightStoneValueString)
			newStoneCounter.AddStone(leftStoneValue, count)
			newStoneCounter.AddStone(rightStoneValue, count)
		} else {
			newStoneCounter.AddStone(currentStoneValue*2024, count)
		}
	}
	stoneCounter.stoneCountMap = newStoneCounter.stoneCountMap
}
