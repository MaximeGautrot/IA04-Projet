package frontend

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Interface Sprite =================================

type Sprite interface {
	Update()
	Draw(screen *ebiten.Image)
	SetPosition(x, y float64)
	SetAnimationRow(row int)
}

// BaseSprite =======================================

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

// Rooster Sprite, Cow Sprite, Bull Sprite ==============================

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

// Sprite Humain =======================================

const (
	HumanFrameWidth  = 64
	HumanFrameHeight = 64
	HumanAnimSpeed   = 10

	HumanRowsPerSheet = 4

	HumanFrameCountIdle  = 9
	HumanFrameCountWalk  = 6
	HumanFrameCountRun   = 8
	HumanFrameCountHurt  = 5
	HumanFrameCountDeath = 7
)

type humanAnim struct {
	sheet      *ebiten.Image
	frameCount int
}

type HumanSprite struct {
	BaseSprite
	anims        []humanAnim 
	globalRow    int         
	rowsPerSheet int        
}

func NewHumanSprite() *HumanSprite {
	// Crée la liste ordonnée des animations : ORDRE EST IMPORTANT
	anims := []humanAnim{
		{sheet: imgHumanIdle, frameCount: HumanFrameCountIdle},   // index 0: Lignes 0-3
		{sheet: imgHumanWalk, frameCount: HumanFrameCountWalk},   // index 1: Lignes 4-7
		{sheet: imgHumanRun, frameCount: HumanFrameCountRun},    // index 2: Lignes 8-11
		{sheet: imgHumanHurt, frameCount: HumanFrameCountHurt},   // index 3: Lignes 12-15
		{sheet: imgHumanDeath, frameCount: HumanFrameCountDeath}, // index 4: Lignes 16-19
	}

	s := &HumanSprite{
		BaseSprite: BaseSprite{
			sheet:       anims[0].sheet,
			currentAnim: 0,
			frameWidth:  HumanFrameWidth,
			frameHeight: HumanFrameHeight,
			frameCount:  anims[0].frameCount,
			animSpeed:   HumanAnimSpeed,
		},
		anims:        anims,
		globalRow:    0,
		rowsPerSheet: HumanRowsPerSheet,
	}
	return s
}

func (s *HumanSprite) SetAnimationRow(globalRow int) {
	if s.globalRow == globalRow {
		return
	}
	s.globalRow = globalRow

	sheetIndex := globalRow / s.rowsPerSheet

	localRow := globalRow % s.rowsPerSheet

	if sheetIndex >= len(s.anims) {
		sheetIndex = len(s.anims) - 1
	}

	anim := s.anims[sheetIndex]

	s.sheet = anim.sheet
	s.frameCount = anim.frameCount
	s.currentAnim = localRow 

	s.currentFrame = 0
	s.animTick = 0
}