package topographicmap

import (
	"hmcalister/AdventOfCode/hashset"
	"log/slog"
)

const (
	MAX_HEIGHT = 9
)

type TopographicMap struct {
	heightMap                        map[Coordinate]int
	trailheads                       []Coordinate
	memoizedReachableEndpoints       map[Coordinate]*hashset.HashSet[Coordinate]
	memoizedDistinctPathsToEndpoints map[Coordinate]int
}

func ParseInputToTopographicMap(lines []string) *TopographicMap {
	topographicMap := &TopographicMap{
		heightMap:                        make(map[Coordinate]int),
		trailheads:                       make([]Coordinate, 0),
		memoizedReachableEndpoints:       make(map[Coordinate]*hashset.HashSet[Coordinate]),
		memoizedDistinctPathsToEndpoints: make(map[Coordinate]int),
	}
	for y, line := range lines {
		for x, heightRune := range line {
			coord := Coordinate{x, y}
			height := (int(heightRune) - '0')
			slog.Debug("found next height", "coordinate", coord, "height", height)
			topographicMap.heightMap[coord] = height
			if height == 0 {
				topographicMap.trailheads = append(topographicMap.trailheads, coord)
			}
		}
	}

	return topographicMap
}

// --------------------------------------------------------------------------------
// Part 01

// Given the current coordinate of a trail, find all the neighbors that are slightly uphill (height increases by 1)
// and call this function on them. Base case is height=MAX_HEIGHT where trails stop.
//
// Returns the hashset of all endpoints reachable by the current coordinate.
//
// Also uses a memoized map to avoid recomputation of points.
func (topographicMap *TopographicMap) calculateReachableEndpointsRecursive(currentCoordinate Coordinate) *hashset.HashSet[Coordinate] {
	if memoizedValue, ok := topographicMap.memoizedReachableEndpoints[currentCoordinate]; ok {
		return memoizedValue
	}

	currentHeight := topographicMap.heightMap[currentCoordinate]
	reachableCoordinates := hashset.New[Coordinate]()

	// Compute base case
	if currentHeight == MAX_HEIGHT {
		reachableCoordinates.Add(currentCoordinate)
		return reachableCoordinates
	}

	for _, neighbor := range currentCoordinate.GetOrthogonalNeighbors() {
		if neighborHeight, ok := topographicMap.heightMap[neighbor]; ok && neighborHeight == currentHeight+1 {
			neighborReachableCoordinates := topographicMap.calculateReachableEndpointsRecursive(neighbor)
			hashset.CombineHashSets(reachableCoordinates, neighborReachableCoordinates)
		}
	}
	topographicMap.memoizedReachableEndpoints[currentCoordinate] = reachableCoordinates

	return reachableCoordinates
}

func (topographicMap *TopographicMap) calculateTrailheadOrthogonalScore(startCoordinate Coordinate) int {
	initialHeight, ok := topographicMap.heightMap[startCoordinate]
	if !ok {
		slog.Error("given coordinate is not in topographic map", "coordinate", startCoordinate)
		return 0
	}
	if initialHeight != 0 {
		slog.Error("given coordinate is not a trailhead", "coordinate", startCoordinate, "initial height", initialHeight)
		return 0
	}

	reachableEndpoints := topographicMap.calculateReachableEndpointsRecursive(startCoordinate)
	slog.Debug("found reachable endpoints", "trailhead coordinate", startCoordinate, "number of reachable endpoints", reachableEndpoints.Size(), "reachable endpoints", reachableEndpoints)
	return reachableEndpoints.Size()
}

func (topographicMap *TopographicMap) CalculateAllTrailheadOrthogonalScores() int {
	totalScore := 0
	for _, trailhead := range topographicMap.trailheads {
		slog.Debug("computing score of trailhead", "trailhead coordinate", trailhead)
		totalScore += topographicMap.calculateTrailheadOrthogonalScore(trailhead)
	}

	return totalScore
}

// --------------------------------------------------------------------------------
// Part 02

func (topographicMap *TopographicMap) calculateDistinctPathsToEndpointsRecursive(currentCoordinate Coordinate) int {
	if memoizedValue, ok := topographicMap.memoizedDistinctPathsToEndpoints[currentCoordinate]; ok {
		return memoizedValue
	}

	currentHeight := topographicMap.heightMap[currentCoordinate]
	distinctPathCount := 0

	// Compute base case
	if currentHeight == MAX_HEIGHT {
		distinctPathCount += 1
		return distinctPathCount
	}

	for _, neighbor := range currentCoordinate.GetOrthogonalNeighbors() {
		if neighborHeight, ok := topographicMap.heightMap[neighbor]; ok && neighborHeight == currentHeight+1 {
			neighborDistinctPathCount := topographicMap.calculateDistinctPathsToEndpointsRecursive(neighbor)
			distinctPathCount += neighborDistinctPathCount
		}
	}
	topographicMap.memoizedDistinctPathsToEndpoints[currentCoordinate] = distinctPathCount

	return distinctPathCount
}

func (topographicMap *TopographicMap) calculateTrailheadOrthogonalRating(startCoordinate Coordinate) int {
	initialHeight, ok := topographicMap.heightMap[startCoordinate]
	if !ok {
		slog.Error("given coordinate is not in topographic map", "coordinate", startCoordinate)
		return 0
	}
	if initialHeight != 0 {
		slog.Error("given coordinate is not a trailhead", "coordinate", startCoordinate, "initial height", initialHeight)
		return 0
	}

	trailheadRating := topographicMap.calculateDistinctPathsToEndpointsRecursive(startCoordinate)
	slog.Debug("found reachable endpoints", "trailhead coordinate", startCoordinate, "trailhead rating", trailheadRating)
	return trailheadRating
}

func (topographicMap *TopographicMap) CalculateAllTrailheadOrthogonalRatings() int {
	totalRating := 0
	for _, trailhead := range topographicMap.trailheads {
		slog.Debug("computing rating of trailhead", "trailhead coordinate", trailhead)
		totalRating += topographicMap.calculateTrailheadOrthogonalRating(trailhead)
	}

	return totalRating
}
