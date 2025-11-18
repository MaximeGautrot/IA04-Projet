package frontend

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// ======================================================================
// 1. L'INTERFACE SPRITE
// ======================================================================

type Sprite interface {
	Update()
	Draw(screen *ebiten.Image)
	SetPosition(x, y float64)
	SetAnimationRow(row int)
}

// ======================================================================
// 2. LE STRUCT DE BASE (Modifié)
// ======================================================================

type BaseSprite struct {
	sheet *ebiten.Image
	x, y  float64

	animTick     int
	currentFrame int
	currentAnim  int

	frameWidth  int
	frameHeight int
	frameCount  int
	animSpeed   int
}

func (s *BaseSprite) GetFrameCount() int {
	return s.frameCount
}

func (s *BaseSprite) Update() {
	s.animTick++
	if s.animTick >= s.animSpeed {
		s.animTick = 0
		s.currentFrame = (s.currentFrame + 1) % s.GetFrameCount()
	}
}

func (s *BaseSprite) Draw(screen *ebiten.Image) {
	sx := s.currentFrame * s.frameWidth
	sy := s.currentAnim * s.frameHeight
	rect := image.Rect(sx, sy, sx+s.frameWidth, sy+s.frameHeight)
	subImg := s.sheet.SubImage(rect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(s.x, s.y)

	screen.DrawImage(subImg, op)
}

func (s *BaseSprite) SetPosition(x, y float64) {
	s.x = x
	s.y = y
}

// ======================================================================
// 3. SPRITES SPÉCIFIQUES (Modifié)
// ======================================================================

type RoosterSprite struct {
	BaseSprite
}

func NewRoosterSprite() *RoosterSprite {
	return &RoosterSprite{
		BaseSprite: BaseSprite{
			sheet:       imgRooster,
			currentAnim: 0,
			frameWidth:  32,
			frameHeight: 32,
			frameCount:  6,
			animSpeed:   20,
		},
	}
}

func (s *RoosterSprite) SetAnimationRow(row int) {
	if s.currentAnim != row {
		s.currentAnim = row
		s.currentFrame = 0
		s.animTick = 0
	}
}

type CowSprite struct {
	BaseSprite
}

func NewCowSprite() *CowSprite {
	return &CowSprite{
		BaseSprite: BaseSprite{
			sheet:       imgCow,
			currentAnim: 0,
			frameWidth:  64,
			frameHeight: 64,
			frameCount:  6,
			animSpeed:   20,
		},
	}
}

func (s *CowSprite) SetAnimationRow(row int) {
	if s.currentAnim != row {
		s.currentAnim = row
		s.currentFrame = 0
		s.animTick = 0

		if row >= 0 && row < 4 {
			s.frameCount = 6
		} else {
			s.frameCount = 4
		}
	}
}

type BullSprite struct {
	BaseSprite
}

func NewBullSprite() *BullSprite {
	return &BullSprite{
		BaseSprite: BaseSprite{
			sheet:       imgBull,
			currentAnim: 0,
			frameWidth:  64,
			frameHeight: 64,
			frameCount:  6,
			animSpeed:   20,
		},
	}
}

func (s *BullSprite) SetAnimationRow(row int) {
	if s.currentAnim != row {
		s.currentAnim = row
		s.currentFrame = 0
		s.animTick = 0

		if row >= 0 && row < 4 {
			s.frameCount = 6
		} else {
			s.frameCount = 4
		}
	}
}