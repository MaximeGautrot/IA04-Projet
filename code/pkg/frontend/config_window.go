package frontend

import (
	"fmt"
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ConfigParams struct {
	LambdaAnimals float64
	LambdaPlants  float64
	InitHumans    int
	InitAnimals   int
	InitPlants    int
	
	MaxAnimals    int
	MaxPlants     int
	MaxSteps      int
	
	WeightPragmatic    float64
	WeightCautious     float64
	WeightSelfish      float64
	WeightCollectivist float64
}
type Button struct {
	X, Y, W, H int
	Label      string
	OnClick    func()
}

func (b *Button) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, float64(b.X), float64(b.Y), float64(b.W), float64(b.H), color.RGBA{100, 100, 100, 255})
	ebitenutil.DebugPrintAt(screen, b.Label, b.X+10, b.Y+5)
}

func (b *Button) CheckClick(mx, my int) {
	if mx >= b.X && mx <= b.X+b.W && my >= b.Y && my <= b.Y+b.H {
		b.OnClick()
	}
}

type ConfigScreen struct {
	Params  ConfigParams
	IsDone  bool
	Buttons []Button
}

func NewConfigScreen() *ConfigScreen {
	cs := &ConfigScreen{
		Params: ConfigParams{
			LambdaAnimals: 0.35,
			LambdaPlants:  0.35,
			InitHumans:    20,
			InitAnimals:   20,
			InitPlants:    20,

			MaxAnimals:    300,
			MaxPlants:     300,
			MaxSteps:      150000,
			
			WeightPragmatic:    1.0,
			WeightCautious:     1.0,
			WeightSelfish:      1.0,
			WeightCollectivist: 1.0,
		},
		IsDone: false,
	}

	yBase := 50
	step := 35

	cs.Buttons = []Button{
		// 0. Taux Animaux
		{250, yBase, 30, 20, "-", func() { cs.Params.LambdaAnimals -= 0.005; if cs.Params.LambdaAnimals < 0 { cs.Params.LambdaAnimals = 0 } }},
		{300, yBase, 30, 20, "+", func() { cs.Params.LambdaAnimals += 0.005 }},

		// 1. Taux Plantes
		{250, yBase + step, 30, 20, "-", func() { cs.Params.LambdaPlants -= 0.005; if cs.Params.LambdaPlants < 0 { cs.Params.LambdaPlants = 0 } }},
		{300, yBase + step, 30, 20, "+", func() { cs.Params.LambdaPlants += 0.005 }},

		// 2. Init Humains
		{250, yBase + step*2, 30, 20, "-", func() { cs.Params.InitHumans -= 5; if cs.Params.InitHumans < 0 { cs.Params.InitHumans = 0 } }},
		{300, yBase + step*2, 30, 20, "+", func() { cs.Params.InitHumans += 5 }},

		// 3. Init Animaux
		{250, yBase + step*3, 30, 20, "-", func() { cs.Params.InitAnimals -= 5; if cs.Params.InitAnimals < 0 { cs.Params.InitAnimals = 0 } }},
		{300, yBase + step*3, 30, 20, "+", func() { cs.Params.InitAnimals += 5 }},

		// 4. Init Plantes
		{250, yBase + step*4, 30, 20, "-", func() { cs.Params.InitPlants -= 5; if cs.Params.InitPlants < 0 { cs.Params.InitPlants = 0 } }},
		{300, yBase + step*4, 30, 20, "+", func() { cs.Params.InitPlants += 5 }},

		// 5. MAX Animaux (NOUVEAU)
		{250, yBase + step*5, 30, 20, "-", func() { cs.Params.MaxAnimals -= 10; if cs.Params.MaxAnimals < 10 { cs.Params.MaxAnimals = 10 } }},
		{300, yBase + step*5, 30, 20, "+", func() { cs.Params.MaxAnimals += 10 }},

		// 6. MAX Plantes (NOUVEAU)
		{250, yBase + step*6, 30, 20, "-", func() { cs.Params.MaxPlants -= 10; if cs.Params.MaxPlants < 10 { cs.Params.MaxPlants = 10 } }},
		{300, yBase + step*6, 30, 20, "+", func() { cs.Params.MaxPlants += 10 }},

		// 7. Max Steps
		{250, yBase + step*7, 30, 20, "-", func() { cs.Params.MaxSteps -= 500; if cs.Params.MaxSteps < 100 { cs.Params.MaxSteps = 100 } }},
		{300, yBase + step*7, 30, 20, "+", func() { cs.Params.MaxSteps += 500 }},

		// Profils (8, 9, 10, 11)
		{250, yBase + step*8, 30, 20, "-", func() { cs.Params.WeightPragmatic -= 1; if cs.Params.WeightPragmatic < 0 { cs.Params.WeightPragmatic = 0 } }},
		{300, yBase + step*8, 30, 20, "+", func() { cs.Params.WeightPragmatic += 1 }},

		{250, yBase + step*9, 30, 20, "-", func() { cs.Params.WeightCautious -= 1; if cs.Params.WeightCautious < 0 { cs.Params.WeightCautious = 0 } }},
		{300, yBase + step*9, 30, 20, "+", func() { cs.Params.WeightCautious += 1 }},

		{250, yBase + step*10, 30, 20, "-", func() { cs.Params.WeightSelfish -= 1; if cs.Params.WeightSelfish < 0 { cs.Params.WeightSelfish = 0 } }},
		{300, yBase + step*10, 30, 20, "+", func() { cs.Params.WeightSelfish += 1 }},

		{250, yBase + step*11, 30, 20, "-", func() { cs.Params.WeightCollectivist -= 1; if cs.Params.WeightCollectivist < 0 { cs.Params.WeightCollectivist = 0 } }},
		{300, yBase + step*11, 30, 20, "+", func() { cs.Params.WeightCollectivist += 1 }},

		{250, 520, 125, 40, "LANCER SIMULATION", func() { cs.IsDone = true }},
	}
	return cs
}

func (c *ConfigScreen) Update() error {
	// Gestion des clics répétitifs 
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		for i := range c.Buttons {
			c.Buttons[i].CheckClick(mx, my)
		}
	}
	return nil
}

func (c *ConfigScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 50, 255})
	
	ebitenutil.DebugPrintAt(screen, "CONFIGURATION PARAMETRES", 20, 20)
	
	y := 55
	step := 35

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Taux Apparition A. : %.3f", c.Params.LambdaAnimals), 20, y); y+=step
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Taux Apparition P. : %.3f", c.Params.LambdaPlants), 20, y); y+=step
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Humains Depart     : %d", c.Params.InitHumans), 20, y); y+=step
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Animaux Depart     : %d", c.Params.InitAnimals), 20, y); y+=step
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Plantes Depart     : %d", c.Params.InitPlants), 20, y); y+=step
	
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Limite Max Animaux : %d", c.Params.MaxAnimals), 20, y); y+=step
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Limite Max Plantes : %d", c.Params.MaxPlants), 20, y); y+=step
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Max Steps          : %d", c.Params.MaxSteps), 20, y); y+=step

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Poids Pragmatique  : %.0f", c.Params.WeightPragmatic), 20, y); y+=step
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Poids Prudent      : %.0f", c.Params.WeightCautious), 20, y); y+=step
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Poids Egoiste      : %.0f", c.Params.WeightSelfish), 20, y); y+=step
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Poids Collectiviste: %.0f", c.Params.WeightCollectivist), 20, y); y+=step

	for _, b := range c.Buttons {
		b.Draw(screen)
	}
}