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

func (s *Sprite) Move(v Vector) {
	s.x += v.dx
	s.y += v.dy
}

func (s *Sprite) SetPosition(x, y float64) {
	s.x = x
	s.y = y
}

func (s *Sprite) GetPosition() Position {
	return CreatePosition(s.x, s.y)
}

func (s1, s2 *Sprite) IsColliding() bool {
	return !(s1.x+s1.width < s2.x ||
		s1.x > s2.x+s2.width ||
		s1.y+s1.height < s2.y ||
		s1.y > s2.y+s2.height)
}

func (p *Position) DistanceTo(other Position) float64 {
	dx := other.x - p.x
	dy := other.y - p.y
	return math.Sqrt(float64(dx*dx + dy*dy))
}
