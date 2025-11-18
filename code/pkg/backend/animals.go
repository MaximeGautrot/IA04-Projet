package simulation

type animalType int

const (
	Chicken animalType = iota
	Cow
	Bull
)

type Animal struct {
	AgentParams
	typ          animalType
	peopleNeeded int
}

func CreateAnimal(name string, sprite Sprite, typ animalType) *Animal {
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
		typ: typ,
	}
}

func (a *Animal) GetType() animalType {
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
