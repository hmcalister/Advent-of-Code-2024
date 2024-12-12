package garden

import (
	"hmcalister/AdventOfCode/hashset"
)

type plot struct {
	coordinates *hashset.HashSet[Coordinate]
	perimeter   int
}

func newPlot() *plot {
	return &plot{
		coordinates: hashset.New[Coordinate](),
		perimeter:   0,
	}
}

func (plotData *plot) Add(c Coordinate) {
	plotData.coordinates.Add(c)

	// Adding a new totally disjoint cell would give four new perimeter
	plotData.perimeter += 4

	// For each neighbor that is present, two perimeter
	// (one from the new cell and one from the old) is removed
	//
	// If zero neighbors, the full four is left
	// If four neighbors, the new perimeter is removed as well as the previous void (four walls)
	for _, neighbor := range c.GetOrthogonalNeighbors() {
		if plotData.coordinates.Contains(neighbor) {
			plotData.perimeter -= 2
		}
	}
}

func (plotData *plot) isInternalCoordinate(c Coordinate) bool {
	for _, neighbor := range c.GetOrthogonalNeighbors() {
		if !plotData.coordinates.Contains(neighbor) {
			return false
		}
	}
	return true
}

// Idea for determining edge
//
// For left edge --- if coordinate above is NOT in plot
//
//	XXXX
//	XXOO	<--- Left Boundary
//	XOOO	<--- Left Boundary
//	XOOO
//
// OR coordinate above AND coordinate to right BOTH in plot
//
//	XXXX
//	XOOO	<--- Left Boundary (from above)
//	XXOO	<--- Left Boundary (of note)
//	XXOO
//
// And similar for other edges
func (plotData *plot) countEdges() int {
	numEdges := 0
	for _, c := range plotData.coordinates.Items() {
		containsUpperLeft := plotData.coordinates.Contains(Coordinate{c.X - 1, c.Y - 1})
		containsUpperMiddle := plotData.coordinates.Contains(Coordinate{c.X, c.Y - 1})
		containsUpperRight := plotData.coordinates.Contains(Coordinate{c.X + 1, c.Y - 1})
		containsMiddleLeft := plotData.coordinates.Contains(Coordinate{c.X - 1, c.Y})
		containsMiddleRight := plotData.coordinates.Contains(Coordinate{c.X + 1, c.Y})
		containsLowerLeft := plotData.coordinates.Contains(Coordinate{c.X - 1, c.Y + 1})
		containsLowerMiddle := plotData.coordinates.Contains(Coordinate{c.X, c.Y + 1})
		containsLowerRight := plotData.coordinates.Contains(Coordinate{c.X + 1, c.Y + 1})

		// Left edge
		if !containsMiddleLeft && (!containsUpperMiddle || (containsUpperMiddle && containsUpperLeft)) {
			numEdges += 1
		}

		// Top Edge
		if !containsUpperMiddle && (!containsMiddleRight || (containsMiddleRight && containsUpperRight)) {
			numEdges += 1
		}

		// Right Edge
		if !containsMiddleRight && (!containsLowerMiddle || (containsLowerMiddle && containsLowerRight)) {
			numEdges += 1
		}

		// Bottom Edge
		if !containsLowerMiddle && (!containsMiddleLeft || (containsMiddleLeft && containsLowerLeft)) {
			numEdges += 1
		}

	}
	return numEdges
}

func (plotData *plot) fencingPrice() int {
	return plotData.coordinates.Size() * plotData.perimeter
}

func (plotData *plot) discountFencingPrice() int {
	return plotData.coordinates.Size() * plotData.countEdges()
}
