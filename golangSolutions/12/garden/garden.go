package garden

import (
	"log/slog"

	hashset "github.com/hmcalister/Go-DSA/set/HashSet"
	arraystack "github.com/hmcalister/Go-DSA/stack/ArrayStack"
)

type Garden struct {
	gardenData         [][]rune
	width              int
	height             int
	plots              []*plot
	coordinatesInPlots *hashset.HashSet[Coordinate]
}

func NewGarden(gardenData [][]rune) *Garden {
	garden := &Garden{
		gardenData:         gardenData,
		width:              len(gardenData),
		height:             len(gardenData[0]),
		plots:              make([]*plot, 0),
		coordinatesInPlots: hashset.New[Coordinate](),
	}

	garden.initialize()

	return garden
}

func (garden *Garden) isInBounds(c Coordinate) bool {
	return c.X >= 0 && c.X < garden.width && c.Y >= 0 && c.Y < garden.height
}

func (garden *Garden) initialize() {
	for y := 0; y < garden.height; y += 1 {
		for x := 0; x < garden.width; x += 1 {
			c := Coordinate{x, y}
			if !garden.coordinatesInPlots.Contains(c) {
				garden.addNewPlot(c)
			}
		}
	}
}

func (garden *Garden) addNewPlot(initialCoordinate Coordinate) {
	p := newPlot()
	garden.plots = append(garden.plots, p)
	initialRune := garden.gardenData[initialCoordinate.Y][initialCoordinate.X]

	fillCoordinateStack := arraystack.New[Coordinate]()
	fillCoordinateStack.Add(initialCoordinate)

	for fillCoordinateStack.Size() > 0 {
		currentCoordinate, _ := fillCoordinateStack.Remove()
		if garden.coordinatesInPlots.Contains(currentCoordinate) ||
			!garden.isInBounds(currentCoordinate) {
			continue
		}

		currentRune := garden.gardenData[currentCoordinate.Y][currentCoordinate.X]
		if currentRune != initialRune {
			continue
		}
		// slog.Debug("found additional plot coordinate", "initial coordinate", initialCoordinate, "current coordinate", currentCoordinate)
		garden.coordinatesInPlots.Add(currentCoordinate)
		p.Add(currentCoordinate)
		for _, neighbor := range currentCoordinate.GetOrthogonalNeighbors() {
			fillCoordinateStack.Add(neighbor)
		}
	}
	slog.Info("new plot initialized", "plot rune", initialRune, "plot area", p.coordinates.Size(), "plot perimeter", p.perimeter, "plot edges", p.countEdges())
}

func (garden *Garden) FencingPrice() int {
	totalFencingPrice := 0
	for _, p := range garden.plots {
		totalFencingPrice += p.fencingPrice()
	}
	return totalFencingPrice
}

func (garden *Garden) DiscountFencingPrice() int {
	totalFencingPrice := 0
	for _, p := range garden.plots {
		totalFencingPrice += p.discountFencingPrice()
	}
	return totalFencingPrice
}
