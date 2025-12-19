package simulation

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
	targetPos       Position // Pour Wander ou Flee
	stepsInState    int      // Pour changer d'état après un temps
	detectedThreats []Agent  // Agents perçus comme des menaces
}

func CreateAnimal(name string, sprite Sprite, typ AnimalType) *Animal {
	var health int
	var peopleNeeded int

	switch typ {
	case Chicken:
		health = 20
		peopleNeeded = 1
	case Cow:
		health = 50
		peopleNeeded = 2
	case Bull:
		health = 80
		peopleNeeded = 3
	}

	return &Animal{
		AgentParams: AgentParams{
			name:   name,
			health: health,
			alive:  true,
			sprite: sprite,
		},
		typ:          typ,
		peopleNeeded: peopleNeeded,
	}
}

func (a *Animal) GetType() AnimalType {
	return a.typ
}

func (a *Animal) GetPeopleNeeded() int {
	return a.peopleNeeded
}

func (a *Animal) GetEnergyValue() uint {
	switch a.typ {
	case Chicken:
		return 20
	case Cow:
		return 50
	case Bull:
		return 80
	default:
		return 0
	}
}

func (a *Animal) Percept(env *Environment) {
}

func (a *Animal) Deliberate() {
}

func (a *Animal) Act(env *Environment) {
}
