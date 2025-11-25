package simulation

type Object interface {
	IsAlive() bool
	GetSprite() Sprite
	GetID() uint
	GetName() string
	Spawn()
}

type ObjectParams struct {
	id     uint
	name   string
	alive  bool
	sprite Sprite
}

func (o *ObjectParams) IsAlive() bool {
	return o.alive
}

func (o *ObjectParams) GetSprite() Sprite {
	return o.sprite
}

func (o *ObjectParams) GetID() uint {
	return o.id
}

func (o *ObjectParams) GetName() string {
	return o.name
}

func (o *ObjectParams) Spawn(x, y float64) {
	o.alive = true
	o.sprite.SetPosition(x, y)
}
