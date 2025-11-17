package simulation

import "math"

type Position struct {
	x float64
	y float64
}

type Sprite struct {
	Position
	imagePath string
	width     float64
	height    float64
}

type Vector struct {
	dx float64
	dy float64
}

func CreatePosition(x, y float64) Position {
	return Position{x: x, y: y}
}

func CreateSprite(x, y float64, imagePath string, width, height float64) Sprite {
	return Sprite{
		Position:  CreatePosition(x, y),
		imagePath: imagePath,
		width:     width,
		height:    height,
	}
}

func CreateVector(dx, dy float64) Vector {
	return Vector{dx: dx, dy: dy}
}

func (p *Position) Move(v Vector) {
	p.x += v.dx
	p.y += v.dy
}

func (p *Position) SetPosition(x, y float64) {
	p.x = x
	p.y = y
}

func (p Position) GetPosition() Position {
	return p
}

func (p *Position) DistanceTo(other Position) float64 {
	dx := other.x - p.x
	dy := other.y - p.y
	return math.Sqrt(dx*dx + dy*dy)
}

func (s1 *Sprite) IsColliding(s2 *Sprite) bool {
	return !(s1.x+s1.width < s2.x ||
		s1.x > s2.x+s2.width ||
		s1.y+s1.height < s2.y ||
		s1.y > s2.y+s2.height)
}
