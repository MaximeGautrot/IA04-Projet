package simulation

import (
	"math"
	"math/rand"
)

const (
	ActionRange    = 10.0 // Distance pour interagir
	VisionRadius   = 150.0
	MoveSpeed      = 2.0
	MaxEnergy      = 100
	MaxHunger      = 100
	EnergyRestRate = 5
	HungerCost     = 1
	EnergyMoveCost = 1
)

// --- Helper pour le déplacement progressif ---
func moveTowards(a Agent, target Position, env *Environment) bool {
	pos := a.GetSprite().Position
	dist := pos.DistanceTo(target)

	if dist <= ActionRange {
		return true // Arrivé
	}

	// Calcul vecteur
	dx := target.x - pos.x
	dy := target.y - pos.y
	// Normalisation
	length := math.Sqrt(dx*dx + dy*dy)
	if length > 0 {
		dx = (dx / length) * MoveSpeed
		dy = (dy / length) * MoveSpeed
	}

	a.Move(dx, dy, env)
	return false
}

// ================== REST ACTION ==================
type RestAction struct{}

func (r *RestAction) Execute(a Agent, env *Environment) {
	h := a.(*Human)
	// Restaure énergie, augmente faim
	h.energy += EnergyRestRate
	if h.energy > MaxEnergy {
		h.energy = MaxEnergy
	}
	h.hunger += HungerCost
}

func (r *RestAction) evaluateUtility(a Agent, env *Environment) float64 {
	h := a.(*Human)
	utility := float64(MaxEnergy - h.energy)

	// Profil influence
	switch h.profile {
	case Cautious:
		utility *= 1.5 // Très prudent sur son énergie
	case Pragmatic:
		utility *= 1.0
	default:
		utility *= 0.8
	}
	return utility
}

// ================== GATHER ACTION ==================
type GatherAction struct {
	TargetID  uint
	TargetPos Position
}

func (g *GatherAction) Execute(a Agent, env *Environment) {
	h := a.(*Human)

	// Vérifier si la cible existe encore
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
		h.currentAction = nil // Cible perdue
		return
	}

	// Déplacement progressif
	arrived := moveTowards(a, target.GetSprite().Position, env)
	h.energy -= EnergyMoveCost
	h.hunger += HungerCost

	if arrived {
		target.Consume()
		h.hunger -= target.GetEnergyValue()
		if h.hunger < 0 {
			h.hunger = 0
		} // uint wrap fix logic needed ideally, assuming hunger doesn't go negative or is handled
		h.currentAction = nil // Action finie
	}
}

func (g *GatherAction) evaluateUtility(a Agent, env *Environment) float64 {
	h := a.(*Human)
	// Trouver le légume le plus proche
	var closest *Vegetable
	minDist := 99999.0

	for _, obj := range h.visibleObjects {
		if veg, ok := obj.(*Vegetable); ok && veg.IsAlive() {
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

	// Sauvegarde pour l'exécution
	g.TargetID = closest.GetID()
	g.TargetPos = closest.GetSprite().Position

	utility := float64(h.hunger) - (minDist * 0.1) // Utilité = Faim - Coût distance

	switch h.profile {
	case Pragmatic:
		utility *= 1.2 // Efficace
	case Collectivist:
		utility *= 0.8 // Préfère chasser en groupe souvent
	}

	return utility
}

// ================== HUNT ACTION ==================
type HuntAction struct {
	TargetID uint
}

func (hu *HuntAction) Execute(a Agent, env *Environment) {
	h := a.(*Human)

	// Retrouver la cible
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
	h.energy -= EnergyMoveCost
	h.hunger += HungerCost

	if arrived {
		// Logique de chasse de groupe
		hunters := 0
		for _, other := range env.agents {
			if hum, ok := other.(*Human); ok && hum.IsAlive() {
				if hum.GetID() == h.GetID() {
					hunters++
					continue
				}
				// Si un autre humain est proche de la cible (supposons qu'il chasse aussi ou aide)
				if hum.GetSprite().Position.DistanceTo(target.GetSprite().Position) <= ActionRange*1.5 {
					hunters++
				}
			}
		}

		if hunters >= target.GetPeopleNeeded() {
			// Succès
			target.Kill()
			h.hunger -= target.GetEnergyValue() // Partage simplifié : celui qui tue mange
			if h.hunger < 0 {
				h.hunger = 0
			}
		} else {
			// Échec et dégâts
			target.IsAttacked(5) // Dégâts mineurs
			h.IsAttacked(20)     // Gros dégâts au chasseur
		}
		h.currentAction = nil
	}
}

func (hu *HuntAction) evaluateUtility(a Agent, env *Environment) float64 {
	h := a.(*Human)
	var closest *Animal
	minDist := 99999.0

	for _, ag := range h.visibleAgents {
		if ani, ok := ag.(*Animal); ok && ani.IsAlive() {
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

	// Utilité de base
	utility := (float64(h.hunger) * 1.5) - (minDist * 0.1)

	// Risque perçu
	risk := float64(closest.GetPeopleNeeded()) * 10.0

	switch h.profile {
	case Cautious:
		utility -= risk * 2 // Évite le risque
	case Selfish:
		utility -= risk // N'aime pas trop le risque solo
	case Collectivist:
		utility += 20 // Aime la coopération (suppose que les autres viendront)
	case Pragmatic:
		if h.energy > 50 {
			utility += 10
		} // Chasse si en forme
	}

	return math.Max(0, utility)
}

// ================== REPRODUCE ACTION ==================
type ReproduceAction struct {
	MateID uint
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
		// Coût
		h.energy -= 30
		mate.energy -= 30 // Simplification: on touche à l'énergie du partenaire directement ici pour l'exemple

		// Création des enfants
		for i := 0; i < 2; i++ {
			childProfile := determineChildProfile(h.profile, mate.profile)
			childName := h.GetName() + "-Jr"

			// Position proche
			offsetX := (rand.Float64() * 20) - 10
			offsetY := (rand.Float64() * 20) - 10

			newSprite := CreateSprite(h.GetSprite().Position.x+offsetX, h.GetSprite().Position.y+offsetY, 16, 16)

			child := CreateHuman(childName, 100, newSprite, 20, 80, childProfile, "Child")
			env.AddAgent(child)
		}
		h.currentAction = nil
	}
}

func determineChildProfile(p1, p2 Profile) Profile {
	roll := rand.Float64()
	if roll < 0.4 {
		return p1
	} else if roll < 0.8 {
		return p2
	}
	// 20% random
	profiles := []Profile{Selfish, Collectivist, Pragmatic, Cautious}
	return profiles[rand.Intn(len(profiles))]
}

func (r *ReproduceAction) evaluateUtility(a Agent, env *Environment) float64 {
	h := a.(*Human)

	// Conditions minimales
	if h.energy < 60 || h.hunger > 50 {
		return 0.0
	}

	// Trouver un partenaire
	var closestMate *Human
	minDist := 99999.0

	for _, ag := range h.visibleAgents {
		if mate, ok := ag.(*Human); ok && mate.GetID() != h.GetID() {
			d := h.GetSprite().Position.DistanceTo(mate.GetSprite().Position)
			if d < minDist {
				minDist = d
				closestMate = mate
			}
		}
	}

	if closestMate == nil {
		return 0.0
	}
	r.MateID = closestMate.GetID()

	utility := 50.0 - (minDist * 0.2)

	switch h.profile {
	case Selfish:
		utility -= 20 // Préfère garder son énergie
	case Collectivist:
		utility += 30 // Priorité à la survie du groupe
	}

	return utility
}
