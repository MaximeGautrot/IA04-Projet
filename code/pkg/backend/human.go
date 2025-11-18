package simulation

type Profile int

const (
	Selfish      Profile = iota // The Free-Rider
	Collectivist                // The Unconditional Cooperator
	Pragmatic                   // The Reciprocator
	Cautious                    // The Risk-Averse
)

type Human struct {
	AgentParams
	hunger        uint
	profile       Profile
	strategyType  string
	energy        uint
	currentAction Action
}

func CreateGHuman(name string, health uint, sprite Sprite, hunger, energy uint, profile Profile, strategyType string) *Human {
	return &Human{
		AgentParams: AgentParams{
			name:   name,
			health: health,
			alive:  true,
			sprite: sprite,
		},
		hunger:       hunger,
		profile:      profile,
		strategyType: strategyType,
		energy:       energy,
	}
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

func (h *Human) ChoseAction() {
}

func (h *Human) RunAction() {
}
