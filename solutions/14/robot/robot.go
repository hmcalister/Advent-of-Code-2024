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

