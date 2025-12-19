package simulation

type Environment struct {
	width   int
	height  int
	agents  []Agent
	objects []Object
}

func CreateEnvironment(width int, height int) Environment {
	return Environment{
		width:   width,
		height:  height,
		agents:  []Agent{},
		objects: []Object{},
	}
}

func (e *Environment) AddAgent(agent Agent) {
	e.agents = append(e.agents, agent)
}

func (e *Environment) IsPositionInside(s Sprite, dx, dy float64) bool {
	newX := s.x + dx
	newY := s.y + dy

	return !(newX < 0 ||
		newY < 0 ||
		newX+float64(s.width) > float64(e.width) ||
		newY+float64(s.height) > float64(e.height))
}
