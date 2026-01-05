package simulation

import "math"

type Position struct {
	X float64
	Y float64
}

type Sprite struct {
	Position
	width  int
	height int
}

type Vector struct {
	dx float64
	dy float64
}

func CreatePosition(x, y float64) Position {
	return Position{X: x, Y: y}
}

func CreateSprite(x, y float64, width, height int) Sprite {
	return Sprite{
		Position: CreatePosition(x, y),
		width:    width,
		height:   height,
	}
}

func CreateVector(dx, dy float64) Vector {
	return Vector{dx: dx, dy: dy}
}

func (p *Position) MovePosition(v Vector) {
	p.X += v.dx
	p.Y += v.dy
}

func (p *Position) SetPosition(x, y float64) {
	p.X = x
	p.Y = y
}

func (p Position) GetPosition() Position {
	return p
}

func (p Position) DistanceTo(other Position) float64 {
	dx := other.X - p.X
	dy := other.Y - p.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (s1 Sprite) IsColliding(s2 *Sprite) bool {
	return !(s1.X+float64(s1.width) < s2.X ||
		s1.X > s2.X+float64(s2.width) ||
		s1.Y+float64(s1.height) < s2.Y ||
		s1.Y > s2.Y+float64(s2.height))
}