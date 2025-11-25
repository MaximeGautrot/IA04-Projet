package simulation

type vegetableType int

const (
	Carrot vegetableType = iota
	Lettuce
	Berry
)

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

func (v *Vegetable) GetEnergyValue() uint {
	switch v.typ {
	case Carrot:
		return 12
	case Lettuce:
		return 8
	case Berry:
		return 5
	default:
		return 0
	}
}

func (v *Vegetable) Consume() {
	v.alive = false
}
