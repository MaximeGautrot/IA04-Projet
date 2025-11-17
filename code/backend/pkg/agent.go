package simulation

type Agent struct {
	name    string
	heatlth int
	alive   bool
	pos     Position
	sprite  Sprite
}

func (a *Agent) isAlive() bool {
	return a.alive
}
