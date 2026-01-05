package simulation

type vegetableType int

const (
	Carrot vegetableType = iota
	Lettuce
	Berry
)

const HungerValue = 5

type Vegetable struct {
	ObjectParams
	typ vegetableType
}

func CreateVegetable(id uint, name string, sprite Sprite, typ vegetableType) *Vegetable {
	return &Vegetable{
		ObjectParams: ObjectParams{
			id:     id,
			name:   name,
			alive:  true,
			sprite: sprite,
		},
		typ: typ,
	}
}

func (v *Vegetable) GetType() vegetableType {
	return v.typ
}

func (v *Vegetable) GetSprite() Sprite {
	return v.sprite
}

func (v *Vegetable) GetHungerValue() uint {
	switch v.typ {
	case Carrot:
		return 60
	case Lettuce:
		return 40
	case Berry:
		return 25
	default:
		return 0
	}
}

func (v *Vegetable) Consume() {
	v.alive = false
}
