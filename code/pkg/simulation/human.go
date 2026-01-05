package simulation

type Profile int

const (
	Selfish      Profile = iota
	Collectivist
	Pragmatic
	Cautious
)

type Human struct {
	AgentParams
	hunger         uint
	profile        Profile
	strategyType   string
	energy         uint
	currentAction  Action
	visibleAgents  []Agent
	visibleObjects []Object
	tickCounter    int 
	actionDuration int
}

// CreateHuman initialise un humain
func CreateHuman(name string, health int, sprite Sprite, hunger, energy uint, profile Profile, strategyType string) *Human {
	baseParams := NewAgentParams(0, name, health, sprite)

	return &Human{
		AgentParams:    baseParams,
		hunger:         hunger,
		profile:        profile,
		strategyType:   strategyType,
		energy:         energy,
		visibleAgents:  []Agent{},
		visibleObjects: []Object{},
		tickCounter:    0,
		actionDuration: 0,
	}
}

// Start lance la boucle de vie
func (h *Human) Start(env *Environment) {
	go func() {
		for {
			select {
			case <-h.syncChan:
				if h.alive {
					h.Percept(env)
					h.Deliberate()
					h.Act(env)
				}
				h.doneChan <- true

			case <-h.stopChan:
				return
			}
		}
	}()
}

func (h *Human) GetHunger() uint { 
	return h.hunger 
}

func (h *Human) GetProfile() Profile { 
	return h.profile 
}

func (h *Human) GetStrategyType() string { 
	return h.strategyType 
}

func (h *Human) GetEnergy() uint { 
	return h.energy 
}

func (h *Human) GetCurrentAction() Action { 
	return h.currentAction 
}

func (h *Human) Percept(env *Environment) {
	h.visibleAgents = []Agent{}
	h.visibleObjects = []Object{}

	for _, a := range env.agents {
		if a.GetID() == h.GetID() || !a.IsAlive() { continue }
		if h.GetSprite().Position.DistanceTo(a.GetSprite().Position) <= VisionRadius {
			h.visibleAgents = append(h.visibleAgents, a)
		}
	}

	for _, o := range env.objects {
		if !o.IsAlive() { continue }
		if h.GetSprite().Position.DistanceTo(o.GetSprite().Position) <= VisionRadius {
			h.visibleObjects = append(h.visibleObjects, o)
		}
	}
}

func (h *Human) Deliberate() {
	if h.currentAction != nil && h.actionDuration < 90 {
		return
	}

	h.currentAction = nil
	
	candidates := []Action{
		&RestAction{},
		&GatherAction{},
		&HuntAction{},
		&ReproduceAction{},
	}

	var bestAction Action
	bestUtility := 0.0

	for _, action := range candidates {
		utility := action.evaluateUtility(h, nil)
		if utility > bestUtility {
			bestUtility = utility
			bestAction = action
		}
	}

	if bestAction != nil && bestUtility > 0 {
		h.currentAction = bestAction
		h.actionDuration = 0
	}
}

func (h *Human) Act(env *Environment) {
	h.tickCounter++

	if h.tickCounter >= 30 {
		h.tickCounter = 0
		
		if h.energy > 0 {
			h.energy--
		} else {
			h.IsAttacked(1)
		}

		if h.hunger < MaxHunger {
			h.hunger++
		} else {
			h.IsAttacked(1)
		}
	}

	if !h.IsAlive() {
		return
	}

	if h.currentAction != nil {
		h.currentAction.Execute(h, env)
		h.actionDuration++
	} else {
		h.actionDuration = 0
	}
}