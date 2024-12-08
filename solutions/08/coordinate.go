package main

type Coordinate struct {
	X int
	Y int
}

func (c Coordinate) InBounds(mapWidth, mapHeight int) bool {
	return c.X >= 0 &&
		c.X < mapWidth &&
		c.Y >= 0 &&
		c.Y < mapHeight
}

func determineAntinode(c1, c2 Coordinate) Coordinate {
	return Coordinate{2*c1.X - c2.X, 2*c1.Y - c2.Y}
}
