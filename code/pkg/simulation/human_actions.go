package simulation

import (
	"math"
	"math/rand"
)

const (
	ActionRange    = 10.0
	VisionRadius   = 250.0
	MoveSpeed      = 2.0
	MaxEnergy      = 500
	MaxHunger      = 500
	MaxHealth      = 100
	EnergyRestRate = 1
	HungerCost     = 1
	HealthRestRate = 1
	EnergyMoveCost = 0
)

func moveTowards(a Agent, target Position, env *Environment) bool {
	pos := a.GetSprite().Position
	dist := pos.DistanceTo(target)

	if dist <= ActionRange {
		return true
	}

	dx := target.X - pos.X
	dy := target.Y - pos.Y
	
	length := math.Sqrt(dx*dx + dy*dy)
	if length > 0 {
		dx = (dx / length) * MoveSpeed
		dy = (dy / length) * MoveSpeed
	}

	a.Move(dx, dy, env)
	return false
}

type RestAction struct{}

func (r *RestAction) Execute(a Agent, env *Environment) {
	h := a.(*Human)
	
	// Modulo 2 : Récupère de l'énergie tous les 2 ticks (environ 30 fois par seconde)
	if h.actionDuration % 2 == 0 {
		h.energy += EnergyRestRate
		if h.energy > MaxEnergy {
			h.energy = MaxEnergy
		}
		
	}

	if h.actionDuration % 4 == 0 {
		h.hunger += HungerCost 
		if h.hunger > MaxHunger {
			h.hunger = MaxHunger
		}
	}

	// Modulo 5 : Soigne tous les 5 ticks (environ 6 fois par seconde)
	if h.actionDuration % 5 == 0 {
		if h.health < MaxHealth {
			h.health += 1 // +1 PV
			if h.health > MaxHealth {
				h.health = MaxHealth
			}
		}
	}
}

func (r *RestAction) evaluateUtility(a Agent, env *Environment) float64 {
	h := a.(*Human)
	
	utilityEnergy := float64(MaxEnergy - h.energy) / 2 - float64(h.hunger)
	
	utilityHealth := float64(MaxHealth - h.health) * 2.0 

	utility := utilityEnergy + utilityHealth

	switch h.profile {
	case Cautious:
		utility *= 1.25
	case Pragmatic:
		utility *= 1.0
	default:
		utility *= 0.8
	}
	
	return utility
}

type GatherAction struct {
	TargetID  uint
	TargetPos Position
}

func (g *GatherAction) Execute(a Agent, env *Environment) {
	h := a.(*Human)

	var target *Vegetable
	for _, obj := range env.objects {
		if obj.GetID() == g.TargetID {
			if veg, ok := obj.(*Vegetable); ok && veg.IsAlive() {
				target = veg
			}
			break
		}
	}

	if target == nil {
		h.currentAction = nil
		return
	}

	arrived := moveTowards(a, target.GetSprite().Position, env)
	
	if h.energy >= EnergyMoveCost { h.energy -= EnergyMoveCost } else { h.energy = 0 }
	h.hunger += HungerCost
	if h.hunger > MaxHunger { h.hunger = MaxHunger }

	if arrived {
		target.Consume()
		val := target.GetHungerValue()
		if h.hunger < val {
			h.hunger = 0
		} else {
			h.hunger -= val
		}
		h.currentAction = nil
	}
}

func (g *GatherAction) evaluateUtility(a Agent, env *Environment) float64 {
	h := a.(*Human)
	var closest *Vegetable
	minDist := 99999.0

	for _, obj := range h.visibleObjects {
		if veg, ok := obj.(*Vegetable); ok && veg.IsAlive() {
			
			alreadyTargeted := false
			
			for _, neighbor := range h.visibleAgents {
				if neighborHuman, ok := neighbor.(*Human); ok {
					if act, isGathering := neighborHuman.GetCurrentAction().(*GatherAction); isGathering && act != nil {
						if act.TargetID == veg.GetID() {
							alreadyTargeted = true
							break
						}
					}
				}
			}

			if alreadyTargeted {
				continue
			}

			d := h.GetSprite().Position.DistanceTo(veg.GetSprite().Position)
			if d < minDist {
				minDist = d
				closest = veg
			}
		}
	}

	if closest == nil {
		return 0.0
	}

	g.TargetID = closest.GetID()
	g.TargetPos = closest.GetSprite().Position

	utility := float64(h.hunger) - (minDist * 0.1)

	switch h.profile {
	case Pragmatic:
		utility *= 1.25
	case Collectivist:
		utility *= 0.8
	case Cautious:
		utility *= 1.5
	}

	return utility
}

type HuntAction struct {
	TargetID uint
}

func (hu *HuntAction) Execute(a Agent, env *Environment) {
	h := a.(*Human)

	var target *Animal
	for _, agent := range env.agents {
		if agent.GetID() == hu.TargetID {
			if ani, ok := agent.(*Animal); ok && ani.IsAlive() {
				target = ani
			}
			break
		}
	}

	if target == nil {
		h.currentAction = nil
		return
	}

	arrived := moveTowards(a, target.GetSprite().Position, env)
	
	if h.energy >= EnergyMoveCost { h.energy -= EnergyMoveCost }
	h.hunger += HungerCost
	if h.hunger > MaxHunger { h.hunger = MaxHunger }

	if arrived {
		hunters := 0
		participatingHunters := []*Human{}

		for _, other := range env.agents {
			if hum, ok := other.(*Human); ok && hum.IsAlive() {
				isMe := (hum.GetID() == h.GetID())
				isClose := (hum.GetSprite().Position.DistanceTo(target.GetSprite().Position) <= ActionRange*1.5)

				if isMe || isClose {
					hunters++
					participatingHunters = append(participatingHunters, hum)
				}
			}
		}

		if hunters >= target.GetPeopleNeeded() {
			target.IsAttacked(20)
			h.IsAttacked(4)
		} else {
			target.IsAttacked(10)
		}

		if target.GetHealth() <= 0 {
			target.Kill()
			val := target.GetHungerValue()/uint(target.GetPeopleNeeded())

			for _, hunter := range participatingHunters {
				if hunter.hunger < val {
					hunter.hunger = 0
				} else {
					hunter.hunger -= val
				}
				hunter.currentAction = nil
			}
		}
		
		if !target.IsAlive() {
			h.currentAction = nil
		}
	}
}

func (hu *HuntAction) evaluateUtility(a Agent, env *Environment) float64 {
	h := a.(*Human)
	var closest *Animal
	minDist := 99999.0

	var currentTargetID uint = 0
	isAlreadyHunting := false
	currentHunters := 0
	
	if currentAct, ok := h.GetCurrentAction().(*HuntAction); ok && currentAct != nil {
		currentTargetID = currentAct.TargetID
		isAlreadyHunting = true
	}

	visibleAlliesCount := 0
	for _, neighbor := range h.visibleAgents {
		if _, ok := neighbor.(*Human); ok {
			visibleAlliesCount++
		}
	}

	for _, ag := range h.visibleAgents {
		if ani, ok := ag.(*Animal); ok && ani.IsAlive() {
			
			if !isAlreadyHunting {
				if (1 + visibleAlliesCount) < ani.GetPeopleNeeded() {
					continue
				}
			}

			for _, neighbor := range h.visibleAgents {
				if neighborHuman, ok := neighbor.(*Human); ok {
					if act, isHunting := neighborHuman.GetCurrentAction().(*HuntAction); isHunting && act != nil {
						if act.TargetID == ani.GetID() {
							currentHunters++
						}
					}
				}
			}

			if currentHunters >= ani.GetPeopleNeeded() {
				if !(isAlreadyHunting && currentTargetID == ani.GetID()) {
					continue
				}
			}

			d := h.GetSprite().Position.DistanceTo(ani.GetSprite().Position)
			if d < minDist {
				minDist = d
				closest = ani
			}
		}
	}

	if closest == nil {
		return 0.0
	}

	hu.TargetID = closest.GetID()
	
	utility := (float64(h.hunger) * 1.5) - (minDist * 0.1)
	risk := float64(closest.GetPeopleNeeded()) * 10

	switch h.profile {
	case Cautious:
		utility -= risk * 1.25
	case Selfish:
		utility -= (float64(closest.GetPeopleNeeded()) - 1) * 15
	case Collectivist:
		utility += 50
	case Pragmatic:
		if h.energy > 200 {
			utility += 50
		}
		utility += float64(currentHunters) * 15
	}

	return math.Max(0, utility)
}

type ReproduceAction struct {
	MateID uint
}

func isPhysicallyReady(h *Human) bool {
	return h.energy >= 400 && h.hunger <= 150
}

func (r *ReproduceAction) Execute(a Agent, env *Environment) {
	h := a.(*Human)
	var mate *Human

	for _, ag := range env.agents {
		if ag.GetID() == r.MateID {
			if m, ok := ag.(*Human); ok && m.IsAlive() {
				mate = m
			}
			break
		}
	}

	if mate == nil {
		h.currentAction = nil
		return
	}

	arrived := moveTowards(a, mate.GetSprite().Position, env)

	if arrived {
		_, mateIsReproducing := mate.currentAction.(*ReproduceAction)
		
		if !mateIsReproducing {
			return 
		}


		if h.energy >= 250 && mate.energy >= 250 {
			h.energy -= 250
			mate.energy -= 250
		} else {
			h.currentAction = nil
			mate.currentAction = nil
			return
		}

		childProfile := determineChildProfile(h.profile, mate.profile)
		childName := h.GetName() + "-Jr"
		offsetX := (rand.Float64() * 20) - 10
		offsetY := (rand.Float64() * 20) - 10
		
		newSprite := CreateSprite(h.GetSprite().Position.X+offsetX, h.GetSprite().Position.Y+offsetY, 16, 16)
		
		child := CreateHuman(childName, 100, newSprite, 20, 80, childProfile, "Child")
		child.SetID(uint(rand.Uint32())) // ID Temporaire
		
		env.AddAgent(child)
		child.Start(env)

		h.currentAction = nil
		mate.currentAction = nil
	}
}

func determineChildProfile(p1, p2 Profile) Profile {
	roll := rand.Float64()
	if roll < 0.4 { return p1 }
	if roll < 0.8 { return p2 }
	profiles := []Profile{Selfish, Collectivist, Pragmatic, Cautious}
	return profiles[rand.Intn(len(profiles))]
}

func (r *ReproduceAction) evaluateUtility(a Agent, env *Environment) float64 {
	h := a.(*Human)

	if !isPhysicallyReady(h) {
		return 0.0
	}

	visibleFood := 0
	nearbyHumans := 0

	for _, obj := range h.visibleObjects {
		if veg, ok := obj.(*Vegetable); ok && veg.IsAlive() {
			visibleFood++
		}
	}
	
	for _, ag := range h.visibleAgents {
		if ani, ok := ag.(*Animal); ok && ani.IsAlive() {
			visibleFood++
		}
		if _, ok := ag.(*Human); ok {
			nearbyHumans++
		}
	}

	if visibleFood < 2 {
		return 0.0
	}

	if nearbyHumans > 5 {
		return 0.0
	}

	var closestMate *Human
	minDist := 99999.0

	for _, ag := range h.visibleAgents {
		if mate, ok := ag.(*Human); ok && mate.GetID() != h.GetID() {
			if isPhysicallyReady(mate) {
				d := h.GetSprite().Position.DistanceTo(mate.GetSprite().Position)
				if d < minDist {
					minDist = d
					closestMate = mate
				}
			}
		}
	}

	if closestMate == nil {
		return 0.0
	}
	
	r.MateID = closestMate.GetID()

	utility := float64(h.GetEnergy())

	utility -= (minDist * 0.5)

	utility += float64(visibleFood) * 10.0

	return math.Max(0.0, utility)
}