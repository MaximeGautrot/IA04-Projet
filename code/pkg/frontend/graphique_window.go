package frontend

import (
	"fmt"
	"ia04project/pkg/simulation"
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GraphScreen struct {
	History []simulation.TurnData
}

func NewGraphScreen(history []simulation.TurnData) *GraphScreen {
	return &GraphScreen{History: history}
}

func (g *GraphScreen) Update() error {
	return nil
}

func (g *GraphScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	
	if len(g.History) < 2 {
		ebitenutil.DebugPrintAt(screen, "Pas assez de données...", 10, 10)
		return
	}

	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	
	rectTop := Rect{X: 50, Y: 50, W: float64(w) - 100, H: float64(h)/2 - 80}
	g.drawGlobalGraph(screen, rectTop)

	rectBot := Rect{X: 50, Y: float64(h)/2 + 20, W: float64(w) - 100, H: float64(h)/2 - 80}
	g.drawProfilesGraph(screen, rectBot)
}

type Rect struct {
	X, Y, W, H float64
}

func (g *GraphScreen) drawGlobalGraph(screen *ebiten.Image, r Rect) {
	ebitenutil.DebugPrintAt(screen, "POPULATION TOTALE (Bleu=Humains, Rouge=Animaux, Vert=Plantes)", int(r.X), int(r.Y)-20)
	
	ebitenutil.DrawRect(screen, r.X, r.Y, r.W, r.H, color.RGBA{240, 240, 240, 255})

	maxPop := 10.0
	for _, d := range g.History {
		if float64(d.HumansAlive) > maxPop { maxPop = float64(d.HumansAlive) }
		if float64(d.AnimalsAlive) > maxPop { maxPop = float64(d.AnimalsAlive) }
		if float64(d.VegetablesAlive) > maxPop { maxPop = float64(d.VegetablesAlive) }
	}

	ebitenutil.DrawLine(screen, r.X, r.Y+r.H, r.X+r.W, r.Y+r.H, color.Black) // Axe X
	ebitenutil.DrawLine(screen, r.X, r.Y, r.X, r.Y+r.H, color.Black)         // Axe Y
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f", maxPop), int(r.X)-30, int(r.Y))

	stepX := r.W / float64(len(g.History))

	for i := 0; i < len(g.History)-1; i++ {
		d1 := g.History[i]
		d2 := g.History[i+1]

		x1 := r.X + float64(i)*stepX
		x2 := r.X + float64(i+1)*stepX

		y1 := r.Y + r.H - (float64(d1.HumansAlive)/maxPop)*r.H
		y2 := r.Y + r.H - (float64(d2.HumansAlive)/maxPop)*r.H
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.RGBA{0, 0, 255, 255})

		y1 = r.Y + r.H - (float64(d1.AnimalsAlive)/maxPop)*r.H
		y2 = r.Y + r.H - (float64(d2.AnimalsAlive)/maxPop)*r.H
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.RGBA{255, 0, 0, 255})
		
		y1 = r.Y + r.H - (float64(d1.VegetablesAlive)/maxPop)*r.H
		y2 = r.Y + r.H - (float64(d2.VegetablesAlive)/maxPop)*r.H
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.RGBA{0, 255, 0, 255})
	}
}

func (g *GraphScreen) drawProfilesGraph(screen *ebiten.Image, r Rect) {
	ebitenutil.DebugPrintAt(screen, "PROFILS HUMAINS", int(r.X), int(r.Y)-20)
	ebitenutil.DebugPrintAt(screen, "Cyan: Pragm | Jaune: Prudent | Violet: Egoiste | Orange: Collectif", int(r.X)+150, int(r.Y)-20)

	ebitenutil.DrawRect(screen, r.X, r.Y, r.W, r.H, color.RGBA{240, 240, 240, 255})

	maxVal := 5.0
	for _, d := range g.History {
		if float64(d.CountPragmatic) > maxVal { maxVal = float64(d.CountPragmatic) }
		if float64(d.CountCautious) > maxVal { maxVal = float64(d.CountCautious) }
		if float64(d.CountSelfish) > maxVal { maxVal = float64(d.CountSelfish) }
		if float64(d.CountCollectivist) > maxVal { maxVal = float64(d.CountCollectivist) }
	}

	// Grille
	ebitenutil.DrawLine(screen, r.X, r.Y+r.H, r.X+r.W, r.Y+r.H, color.Black)
	ebitenutil.DrawLine(screen, r.X, r.Y, r.X, r.Y+r.H, color.Black)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f", maxVal), int(r.X)-30, int(r.Y))

	stepX := r.W / float64(len(g.History))

	for i := 0; i < len(g.History)-1; i++ {
		d1 := g.History[i]
		d2 := g.History[i+1]

		x1 := r.X + float64(i)*stepX
		x2 := r.X + float64(i+1)*stepX

		// Pragmatique (Cyan)
		y1 := r.Y + r.H - (float64(d1.CountPragmatic)/maxVal)*r.H
		y2 := r.Y + r.H - (float64(d2.CountPragmatic)/maxVal)*r.H
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.RGBA{0, 255, 255, 255})

		// Prudent (Jaune / Or)
		y1 = r.Y + r.H - (float64(d1.CountCautious)/maxVal)*r.H
		y2 = r.Y + r.H - (float64(d2.CountCautious)/maxVal)*r.H
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.RGBA{218, 165, 32, 255})

		// Egoiste (Violet)
		y1 = r.Y + r.H - (float64(d1.CountSelfish)/maxVal)*r.H
		y2 = r.Y + r.H - (float64(d2.CountSelfish)/maxVal)*r.H
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.RGBA{138, 43, 226, 255})

		// Collectiviste (Orange)
		y1 = r.Y + r.H - (float64(d1.CountCollectivist)/maxVal)*r.H
		y2 = r.Y + r.H - (float64(d2.CountCollectivist)/maxVal)*r.H
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.RGBA{255, 140, 0, 255})
	}
}