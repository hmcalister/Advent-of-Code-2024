package main

import "log/slog"

type Interval struct {
	start  int
	length int
}

type MapInterval struct {
	Interval

	// A value of -1 indicates an empty interval
	value int
}

func (interval MapInterval) computeChecksumContribution(startDiskIndex, endDiskIndex int) int {
	// slog.Debug("checksum computation", "interval", interval, "start disk index", startDiskIndex, "end disk index", endDiskIndex)
	return interval.value * ((endDiskIndex * (endDiskIndex - 1) / 2) - (startDiskIndex * (startDiskIndex - 1) / 2))
}

type DiskMap struct {
	mapValues       []MapInterval
	totalDiskLength int
}

func ParseLineToDiskMap(line string) *DiskMap {
	diskMap := DiskMap{
		mapValues: make([]MapInterval, len(line)),
	}

	currentMapIndex := 0
	for runeIndex, runeValue := range line {
		intervalLength := (int(runeValue) - '0')
		interval := Interval{
			start:  currentMapIndex,
			length: intervalLength,
		}
		currentMapIndex += intervalLength
		if runeIndex%2 == 0 {
			diskMap.mapValues[runeIndex] = MapInterval{
				Interval: interval,
				value:    runeIndex / 2,
			}
		} else {
			diskMap.mapValues[runeIndex] = MapInterval{
				Interval: interval,
				value:    -1,
			}
		}
		slog.Debug("parsed next interval", "new interval", diskMap.mapValues[runeIndex])
	}

	finalMapInterval := diskMap.mapValues[len(diskMap.mapValues)-1]
	diskMap.totalDiskLength = finalMapInterval.start + finalMapInterval.length

	return &diskMap
}

func (diskMap *DiskMap) ComputeBlockMoveChecksum() int {
	checksum := 0
	currentDiskIndex := 0
	forwardIntervalIndex := 0
	currentForwardInterval := diskMap.mapValues[forwardIntervalIndex]
	reverseIntervalIndex := len(diskMap.mapValues) - 1
	currentReverseInterval := diskMap.mapValues[reverseIntervalIndex]

	// While we have not seen all the intervals
intervalWalkLoop:
	for forwardIntervalIndex < reverseIntervalIndex {
		if currentForwardInterval.value != -1 {
			// We have a non-empty interval
			checksumContribution := currentForwardInterval.computeChecksumContribution(currentDiskIndex, currentDiskIndex+currentForwardInterval.length)
			checksum += checksumContribution
			slog.Debug(
				"found non-empty forward interval",
				"forward interval", currentForwardInterval,
				"current disk index", currentDiskIndex,
				"checksum contribution", checksumContribution,
				"updated checksum", checksum,
			)
			currentDiskIndex += currentForwardInterval.length
			forwardIntervalIndex += 1
			currentForwardInterval = diskMap.mapValues[forwardIntervalIndex]
		} else {
			// We have an empty interval, fill with non-empty intervals in reverse
			// Three cases:
			// 	- we can fit the reverse interval into this space with room to spare (currentForwardInterval.length > currentReverseInterval.length)
			// 	- we cannot fit the reverse interval into this space (currentForwardInterval.length < currentReverseInterval.length)
			// 	- we can fit the reverse interval into this space *exactly* (currentForwardInterval.length == currentReverseInterval.length)

			if currentForwardInterval.length > currentReverseInterval.length {
				// Case one:
				// 	put the entire reverse interval in,
				//	add the checksum,
				// 	increment the disk index,
				//  decrement the empty space length,
				//  and get the next (non-empty) reverse interval

				checksumContribution := currentReverseInterval.computeChecksumContribution(currentDiskIndex, currentDiskIndex+currentReverseInterval.length)
				checksum += checksumContribution
				slog.Debug(
					// "found empty forward interval with length greater than current reverse interval",
					"found empty forward interval case one",
					"forward interval", currentForwardInterval,
					"reverse interval", currentReverseInterval,
					"current disk index", currentDiskIndex,
					"checksum contribution", checksumContribution,
					"updated checksum", checksum,
				)

				currentDiskIndex += currentReverseInterval.length

				currentForwardInterval.length -= currentReverseInterval.length

				for ok := true; ok; ok = currentReverseInterval.value == -1 {
					reverseIntervalIndex -= 1
					currentReverseInterval = diskMap.mapValues[reverseIntervalIndex]
					if forwardIntervalIndex == reverseIntervalIndex {
						break intervalWalkLoop
					}
				}
			} else if currentForwardInterval.length < currentReverseInterval.length {
				// Case two:
				// 	put as much of the reverse interval in as possible,
				//	add the checksum,
				// 	increment the disk index,
				// 	decrement the amount of the reverse interval remaining,
				//  and move to the next forward interval,

				checksumContribution := currentReverseInterval.computeChecksumContribution(currentDiskIndex, currentDiskIndex+currentForwardInterval.length)
				checksum += checksumContribution
				slog.Debug(
					// "found empty forward interval with length less than current reverse interval",
					"found empty forward interval case two",
					"forward interval", currentForwardInterval,
					"reverse interval", currentReverseInterval,
					"current disk index", currentDiskIndex,
					"checksum contribution", checksumContribution,
					"updated checksum", checksum,
				)

				currentDiskIndex += currentForwardInterval.length

				currentReverseInterval.length -= currentForwardInterval.length

				forwardIntervalIndex += 1
				currentForwardInterval = diskMap.mapValues[forwardIntervalIndex]
			} else {
				// Case three:
				// 	put all of the reverse interval in,
				//	add the checksum,
				// 	increment the disk index,
				//  move to the next forward interval,
				//	and get the next (non-empty) reverse interval

				checksumContribution := currentReverseInterval.computeChecksumContribution(currentDiskIndex, currentDiskIndex+currentReverseInterval.length)
				checksum += checksumContribution
				slog.Debug(
					// "found empty forward interval with length greater than current reverse interval",
					"found empty forward interval case three",
					"forward interval", currentForwardInterval,
					"reverse interval", currentReverseInterval,
					"current disk index", currentDiskIndex,
					"checksum contribution", checksumContribution,
					"updated checksum", checksum,
				)

				currentDiskIndex += currentReverseInterval.length

				forwardIntervalIndex += 1
				currentForwardInterval = diskMap.mapValues[forwardIntervalIndex]

				for ok := true; ok; ok = currentReverseInterval.value == -1 {
					reverseIntervalIndex -= 1
					currentReverseInterval = diskMap.mapValues[reverseIntervalIndex]
					if forwardIntervalIndex == reverseIntervalIndex {
						break intervalWalkLoop
					}
				}
			}
		}
	}

	// If there is anything remaining in the reverse interval, that must be accounted for
	if currentReverseInterval.value != -1 && currentReverseInterval.length > 0 {
		checksumContribution := currentReverseInterval.computeChecksumContribution(currentDiskIndex, currentDiskIndex+currentReverseInterval.length)
		checksum += checksumContribution
		slog.Debug(
			"remaining reverse interval",
			"forward interval", currentForwardInterval,
			"reverse interval", currentReverseInterval,
			"current disk index", currentDiskIndex,
			"checksum contribution", checksumContribution,
			"updated checksum", checksum,
		)
	}

	return checksum
}
