package garden

type Coordinate struct {
	X int
	Y int
}

func (c Coordinate) GetOrthogonalNeighbors() []Coordinate {
	return []Coordinate{
		{c.X - 1, c.Y},
		{c.X + 1, c.Y},
		{c.X, c.Y - 1},
		{c.X, c.Y + 1},
	}
}
