package simulation

import (
	"math"
	"math/rand"
)

const (
	AnimalVisionRadius = 100.0
	AnimalSpeed        = 0.5
	WanderDuration     = 100
)

type AnimalType int

const (
	Chicken AnimalType = iota
	Cow
	Bull
)

type AnimalState int

const (
	AnimalStateWander AnimalState = iota
	AnimalStateFlee
	AnimalStateStay
)

type Animal struct {
	AgentParams
	typ             AnimalType
	peopleNeeded    int
	state           AnimalState
	targetPos       Position
	stepsInState    int
	detectedThreats []Agent
}

func CreateAnimal(name string, sprite Sprite, typ AnimalType) *Animal {
	var health int
	var peopleNeeded int

	switch typ {
	case Chicken:
		health = 40
		peopleNeeded = 1
	case Cow:
		health = 120
		peopleNeeded = 2
	case Bull:
		health = 160
		peopleNeeded = 3
	}

	baseParams := NewAgentParams(0, name, health, sprite)

	return &Animal{
		AgentParams:  baseParams,
		typ:          typ,
		peopleNeeded: peopleNeeded,
		state:        AnimalStateWander,
	}
}

// Start lance la boucle d'intelligence de l'animal
func (a *Animal) Start(env *Environment) {
	go func() {
		for {
			select {
			case <-a.syncChan:
				if a.alive {
					a.Percept(env)
					a.Deliberate()
					a.Act(env)
				}
				a.doneChan <- true
			case <-a.stopChan:
				return
			}
		}
	}()
}

func (a *Animal) GetType() AnimalType       { return a.typ }
func (a *Animal) GetPeopleNeeded() int      { return a.peopleNeeded }

func (a *Animal) GetHungerValue() uint {
	switch a.typ {
	case Chicken: 
		return 100
	case Cow:     
		return 260
	case Bull:    
		return 490
	default:      
		return 0
	}
}

// --- IA ANIMAL (Rappel du code précédent) ---

func (a *Animal) Percept(env *Environment) {
	a.detectedThreats = []Agent{}
	// Utilisation de RLock si on veut être strict, mais ici on simplifie
	// env.mutex.RLock(); defer env.mutex.RUnlock() (si ajouté dans env)
	
	for _, agent := range env.agents {
		if _, ok := agent.(*Human); ok && agent.IsAlive() {
			dist := a.GetSprite().Position.DistanceTo(agent.GetSprite().Position)
			if dist < AnimalVisionRadius {
				a.detectedThreats = append(a.detectedThreats, agent)
			}
		}
	}
}

func (a *Animal) Deliberate() {
	if len(a.detectedThreats) > 0 {
		a.state = AnimalStateFlee
	} else {
		a.state = AnimalStateWander
	}
}

func (a *Animal) Act(env *Environment) {
	a.stepsInState++
	currentPos := a.GetSprite().Position

	switch a.state {
	case AnimalStateFlee:
		var fleeX, fleeY float64
		for _, threat := range a.detectedThreats {
			tPos := threat.GetSprite().Position
			fleeX += (currentPos.X - tPos.X)
			fleeY += (currentPos.Y - tPos.Y)
		}
		length := math.Sqrt(fleeX*fleeX + fleeY*fleeY)
		if length > 0 {
			dx := (fleeX / length) * AnimalSpeed * 1.5
			dy := (fleeY / length) * AnimalSpeed * 1.5
			a.Move(dx, dy, env)
		}

	case AnimalStateWander:
		dist := currentPos.DistanceTo(a.targetPos)
		if dist < 5.0 || a.stepsInState > WanderDuration {
			a.stepsInState = 0
			a.targetPos = Position{
				X: rand.Float64() * float64(env.width),
				Y: rand.Float64() * float64(env.height),
			}
		}
		dx := a.targetPos.X - currentPos.X
		dy := a.targetPos.Y - currentPos.Y
		length := math.Sqrt(dx*dx + dy*dy)
		if length > 0 {
			dx = (dx / length) * AnimalSpeed
			dy = (dy / length) * AnimalSpeed
			a.Move(dx, dy, env)
		}
	case AnimalStateStay:
	}
}