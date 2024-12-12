package garden

import hashset "github.com/hmcalister/Go-DSA/set/HashSet"

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

func (plotData *plot) fencingPrice() int {
	return plotData.coordinates.Size() * plotData.perimeter
}
