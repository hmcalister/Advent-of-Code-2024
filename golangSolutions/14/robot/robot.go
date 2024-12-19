package robot

type Robot struct {
	initialPosition Vector2
	velocity        Vector2
}

func NewRobot(initialPosition, velocity Vector2) *Robot {
	return &Robot{
		initialPosition: initialPosition,
		velocity:        velocity,
	}
}

func (robot *Robot) ComputePosition(gridX, gridY int, numSteps int) Vector2 {
	nextX := (robot.initialPosition.X + numSteps*robot.velocity.X) % gridX
	if nextX < 0 {
		nextX += gridX
	}

	nextY := (robot.initialPosition.Y + numSteps*robot.velocity.Y) % gridY
	if nextY < 0 {
		nextY += gridY
	}
	return Vector2{
		X: nextX,
		Y: nextY,
	}
}
