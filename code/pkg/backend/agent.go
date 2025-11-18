package simulation

type Agent interface {
	IsAlive() bool
	GetSprite() Sprite
	GetID() uint
	Move(dx, dy float64, env *Environment)
	GetName() string
	GetHealth() int
	Kill()
	IsAttacked(damage int)
}

type AgentParams struct {
	id     uint
	name   string
	health int
	alive  bool
	sprite Sprite
}

func (ap *AgentParams) GetName() string {
	return ap.name
}

func (ap *AgentParams) GetID() uint {
	return ap.id
}

func (ap *AgentParams) GetHealth() int {
	return ap.health
}

func (ap *AgentParams) GetSprite() Sprite {
	return ap.sprite
}

func (ap *AgentParams) IsAlive() bool {
	return ap.alive
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

func (ap *AgentParams) Move(dx, dy float64, env *Environment) {
	if !env.IsPositionInside(ap.sprite, dx, dy) {
		return
	}

	futureSprite := ap.sprite
	futureSprite.MovePosition(Vector{dx, dy})

	for _, agent := range env.agents {
		if agent.GetID() != ap.id {
			if agent.GetSprite().IsColliding(&futureSprite) {
				return
			}
		}
	}

	ap.sprite.MovePosition(Vector{dx, dy})
}
