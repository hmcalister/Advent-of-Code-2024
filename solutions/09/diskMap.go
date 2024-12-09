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

