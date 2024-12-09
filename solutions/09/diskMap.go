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

