package simulation

type Agent interface {
	isAlive() bool
}

type AgentParams struct {
	name    string
	health  int
	alive   bool
	sprite  Sprite
}

func (ap *AgentParams) GetName() string {
	return ap.name
}

func (ap *AgentParams) GetHealth() int {
	return ap.health
}

func (ap *AgentParams) GetSprite() Sprite {
	return ap.sprite
}

func (ap *AgentParams) Kill() {
	ap.alive = false
}

func (ap *AgentParams) IsAttacked(damage int) {
	ap.health -= damage
	if ap.health <= 0 {
		ap.alive = false
	}
}