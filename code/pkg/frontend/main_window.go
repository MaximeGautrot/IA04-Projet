package frontend

import (
	"fmt"
	"ia04project/pkg/simulation"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	SidebarWidth = 250
	GameWidth    = 800
	GameHeight   = 600
	WindowWidth  = SidebarWidth + GameWidth
)

var (
	imgRooster    *ebiten.Image
	imgCow        *ebiten.Image
	imgBull       *ebiten.Image
	imgHumanIdle  *ebiten.Image
	imgHumanWalk  *ebiten.Image
	imgHumanRun   *ebiten.Image
	imgHumanHurt  *ebiten.Image
	imgHumanDeath *ebiten.Image
)

func loadSpriteSheets() {
	var err error
	imgRooster, _, err = ebitenutil.NewImageFromFile("images/Rooster_animation_with_shadow.png")
	if err != nil {
		log.Fatalf("Erreur chargement Rooster: %v", err)
	}
	imgCow, _, err = ebitenutil.NewImageFromFile("images/Cow_animation_with_shadow.png")
	if err != nil {
		log.Fatalf("Erreur chargement Cow: %v", err)
	}
	imgBull, _, err = ebitenutil.NewImageFromFile("images/Bull_animation_with_shadow.png")
	if err != nil {
		log.Fatalf("Erreur chargement Bull: %v", err)
	}
	imgHumanIdle, _, err = ebitenutil.NewImageFromFile("images/human/Idle.png")
	if err != nil {
		log.Fatalf("Erreur chargement Human Idle: %v", err)
	}
	imgHumanWalk, _, err = ebitenutil.NewImageFromFile("images/human/Walk.png")
	if err != nil {
		log.Fatalf("Erreur chargement Human Walk: %v", err)
	}
	imgHumanRun, _, err = ebitenutil.NewImageFromFile("images/human/Run.png")
	if err != nil {
		log.Fatalf("Erreur chargement Human Run: %v", err)
	}
	imgHumanHurt, _, err = ebitenutil.NewImageFromFile("images/human/Hurt.png")
	if err != nil {
		log.Fatalf("Erreur chargement Human Hurt: %v", err)
	}
	imgHumanDeath, _, err = ebitenutil.NewImageFromFile("images/human/Death.png")
	if err != nil {
		log.Fatalf("Erreur chargement Human Death: %v", err)
	}
}

type Slider struct {
	X, Y, W, H int
	Min, Max   int
	Current    int
	IsDragging bool
}

func (s *Slider) Update() {
	mx, my := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if mx >= s.X && mx <= s.X+s.W && my >= s.Y-10 && my <= s.Y+s.H+10 {
			s.IsDragging = true
		}
	} else {
		s.IsDragging = false
	}

	if s.IsDragging {
		relX := float64(mx - s.X)
		ratio := relX / float64(s.W)
		if ratio < 0 {
			ratio = 0
		}
		if ratio > 1 {
			ratio = 1
		}

		s.Current = s.Min + int(ratio*float64(s.Max-s.Min))
	}
}

func (s *Slider) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, float64(s.X), float64(s.Y+s.H/2)-2, float64(s.W), 4, color.White)

	ratio := float64(s.Current-s.Min) / float64(s.Max-s.Min)
	cursorX := float64(s.X) + ratio*float64(s.W) - 5 // -5 pour centrer le carré de 10px
	ebitenutil.DrawRect(screen, cursorX, float64(s.Y-5), 10, 20, color.RGBA{255, 100, 100, 255})

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%dx", s.Current), int(cursorX)-5, s.Y-20)
}

type MainWindow struct {
	Sim           *simulation.Simulation
	SpriteMap     map[uint]Sprite
	LastPositions map[uint]simulation.Position

	SelectedAgent  simulation.Agent
	SelectedObject simulation.Object

	StopButton  Button
	SpeedSlider Slider

	IsFinished bool
	GameView   *ebiten.Image
}

func NewMainWindow(sim *simulation.Simulation) *MainWindow {
	loadSpriteSheets()
	mw := &MainWindow{
		Sim:           sim,
		SpriteMap:     make(map[uint]Sprite),
		LastPositions: make(map[uint]simulation.Position),
		IsFinished:    false,
		GameView:      ebiten.NewImage(GameWidth, GameHeight),
	}

	mw.StopButton = Button{
		X: 20, Y: GameHeight - 60, W: 210, H: 40,
		Label: "ARRETER SIMULATION",
		OnClick: func() {
			mw.Sim.Stop()
			mw.IsFinished = true
		},
	}

	mw.SpeedSlider = Slider{
		X: 20, Y: 500, W: 200, H: 10,
		Min: 1, Max: 100, Current: 1,
	}

	return mw
}

func (mw *MainWindow) Update() error {
	mw.SpeedSlider.Update()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if !mw.SpeedSlider.IsDragging {
			if mx < SidebarWidth {
				mw.StopButton.CheckClick(mx, my)
			} else {
				mw.handleGameClick(float64(mx-SidebarWidth), float64(my))
			}
		}
	}

	if mw.IsFinished {
		return nil
	}

	steps := mw.SpeedSlider.Current
	for i := 0; i < steps; i++ {
		mw.Sim.Step()
	}

	agents := mw.Sim.GetAllAgents()
	objects := mw.Sim.GetAllObjects()
	aliveIDs := make(map[uint]bool)

	for _, agent := range agents {
		if !agent.IsAlive() {
			continue
		}
		id := agent.GetID()
		aliveIDs[id] = true
		currentPos := agent.GetSprite().Position

		if _, exists := mw.SpriteMap[id]; !exists {
			mw.createSpriteForAgent(agent)
			mw.LastPositions[id] = currentPos
		}
		visualSprite := mw.SpriteMap[id]
		lastPos := mw.LastPositions[id]
		dx := currentPos.X - lastPos.X
		dy := currentPos.Y - lastPos.Y

		visualSprite.SetPosition(currentPos.X, currentPos.Y)
		mw.updateAgentAnimation(agent, visualSprite, dx, dy)
		visualSprite.Update()
		mw.LastPositions[id] = currentPos
	}

	for _, obj := range objects {
		if !obj.IsAlive() {
			continue
		}
		id := obj.GetID()
		aliveIDs[id] = true
		if _, exists := mw.SpriteMap[id]; !exists {
			mw.createVegetableSprite(obj)
		}
		if s, ok := mw.SpriteMap[id]; ok {
			pos := obj.GetSprite().Position
			s.SetPosition(pos.X, pos.Y)
		}
	}

	for id := range mw.SpriteMap {
		if !aliveIDs[id] {
			delete(mw.SpriteMap, id)
			delete(mw.LastPositions, id)
			if mw.SelectedAgent != nil && mw.SelectedAgent.GetID() == id {
				mw.SelectedAgent = nil
			}
		}
	}
	return nil
}

func (mw *MainWindow) handleGameClick(x, y float64) {
	mw.SelectedAgent = nil
	mw.SelectedObject = nil

	minDist := 3600.0
	var closestHuman simulation.Agent

	agents := mw.Sim.GetAllAgents()
	for _, a := range agents {
		if !a.IsAlive() {
			continue
		}
		if _, isHuman := a.(*simulation.Human); !isHuman {
			continue
		}

		pos := a.GetSprite().Position
		distSq := (pos.X-x)*(pos.X-x) + (pos.Y-y)*(pos.Y-y)

		if distSq < minDist {
			minDist = distSq
			closestHuman = a
		}
	}

	if closestHuman != nil {
		mw.SelectedAgent = closestHuman
	}
}

func (mw *MainWindow) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, SidebarWidth, float64(GameHeight), color.RGBA{50, 50, 70, 255})
	mw.drawSidebarInfo(screen)
	mw.StopButton.Draw(screen)

	mw.SpeedSlider.Draw(screen)

	mw.GameView.Fill(color.RGBA{34, 139, 34, 255})
	for _, s := range mw.SpriteMap {
		s.Draw(mw.GameView)
	}

	if mw.SelectedAgent != nil {
		pos := mw.SelectedAgent.GetSprite().Position
		ebitenutil.DrawRect(mw.GameView, pos.X+27, pos.Y+10, 10, 10, color.White)
	}

	opView := &ebiten.DrawImageOptions{}
	opView.GeoM.Translate(SidebarWidth, 0)
	screen.DrawImage(mw.GameView, opView)
}

func (mw *MainWindow) drawSidebarInfo(screen *ebiten.Image) {
	// Stats Globales
	stats := mw.Sim.GetHistory()
	if len(stats) > 0 {
		last := stats[len(stats)-1]
		y := 20
		line := 20
		ebitenutil.DebugPrintAt(screen, "--- STATISTIQUES ---", 10, y)
		y += line
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Temps: %d", last.Tick), 10, y)
		y += line
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Vitesse: x%d", mw.SpeedSlider.Current), 10, y)
		y += line
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Humains: %d", last.HumansAlive), 10, y)
		y += line
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Animaux: %d", last.AnimalsAlive), 10, y)
		y += line
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Plantes: %d", last.VegetablesAlive), 10, y)
		y += line
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pragmatiques: %d", last.CountPragmatic), 10, y)
		y += line
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Prudents: %d", last.CountCautious), 10, y)
		y += line
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Egoïstes: %d", last.CountSelfish), 10, y)
		y += line
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Collectivistes: %d", last.CountCollectivist), 10, y)
	}

	// Infos Agent Sélectionné
	y := 250
	line := 20
	ebitenutil.DebugPrintAt(screen, "--- INSPECTION ---", 10, y)
	y += line
	if mw.SelectedAgent != nil {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ID: %d", mw.SelectedAgent.GetID()), 10, y)
		y += line
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Nom: %s", mw.SelectedAgent.GetName()), 10, y)
		y += line
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Vie: %d", mw.SelectedAgent.GetHealth()), 10, y)
		y += line

		if h, ok := mw.SelectedAgent.(*simulation.Human); ok {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Faim: %d", h.GetHunger()), 10, y)
			y += line
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Energie: %d", h.GetEnergy()), 10, y)
			y += line

			prof := "Inconnu"
			switch h.GetProfile() {
			case simulation.Pragmatic:
				prof = "Pragmatique"
			case simulation.Cautious:
				prof = "Prudent"
			case simulation.Selfish:
				prof = "Egoïste"
			case simulation.Collectivist:
				prof = "Collectiviste"
			}
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Type: %s", prof), 10, y)
			y += line

			action := "Rien"
			if h.GetCurrentAction() != nil {
				switch h.GetCurrentAction().(type) {
				case *simulation.RestAction:
					action = "Dort"
				case *simulation.GatherAction:
					action = "Cueille"
				case *simulation.HuntAction:
					action = "Chasse"
				case *simulation.ReproduceAction:
					action = "Reproduction"
				}
			}
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Action: %s", action), 10, y)
			y += line
		}
	} else {
		ebitenutil.DebugPrintAt(screen, "Cliquez sur un agent", 10, y+20)
	}

	ebitenutil.DebugPrintAt(screen, "--- VITESSE SIMULATION ---", 10, 480)
}

func (mw *MainWindow) Layout(w, h int) (int, int) {
	return WindowWidth, GameHeight
}

func (mw *MainWindow) createVegetableSprite(obj simulation.Object) {
	if veg, ok := obj.(*simulation.Vegetable); ok {
		var s Sprite
		switch veg.GetType() {
		case simulation.Carrot:
			s = NewVegetableSprite(10, 10, 255, 165, 0)
		case simulation.Lettuce:
			s = NewVegetableSprite(12, 12, 50, 205, 50)
		case simulation.Berry:
			s = NewVegetableSprite(8, 8, 148, 0, 211)
		}
		if s != nil {
			pos := veg.GetSprite().Position
			s.SetPosition(pos.X, pos.Y)
			mw.SpriteMap[veg.GetID()] = s
		}
	}
}

func (mw *MainWindow) updateAgentAnimation(agent simulation.Agent, s Sprite, dx, dy float64) {
	const moveThreshold = 0.1
	isMoving := math.Abs(dx) > moveThreshold || math.Abs(dy) > moveThreshold

	switch sTyped := s.(type) {
	case *HumanSprite:
		if !agent.IsAlive() {
			sTyped.SetAnimationRow(16)
			return
		}
		if isMoving {
			baseRow := 4
			if math.Abs(dx) > math.Abs(dy) {
				if dx > 0 {
					sTyped.SetAnimationRow(baseRow + 2)
				} else {
					sTyped.SetAnimationRow(baseRow + 1)
				}
			} else {
				if dy > 0 {
					sTyped.SetAnimationRow(baseRow + 0)
				} else {
					sTyped.SetAnimationRow(baseRow + 3)
				}
			}
		} else {
			sTyped.SetAnimationRow(0)
		}
	case *CowSprite:
		if isMoving {
			if dx > 0 {
				sTyped.SetAnimationRow(1)
			} else {
				sTyped.SetAnimationRow(0)
			}
		}
	case *RoosterSprite:
		if isMoving {
			if dx > 0 {
				sTyped.SetAnimationRow(1)
			} else {
				sTyped.SetAnimationRow(0)
			}
		}
	case *BullSprite:
		if isMoving {
			if dx > 0 {
				sTyped.SetAnimationRow(1)
			} else {
				sTyped.SetAnimationRow(0)
			}
		}
	}
}

func (mw *MainWindow) createSpriteForAgent(agent simulation.Agent) {
	var s Sprite
	switch v := agent.(type) {
	case *simulation.Human:
		s = NewHumanSprite()
	case *simulation.Animal:
		switch v.GetType() {
		case simulation.Chicken:
			s = NewRoosterSprite()
		case simulation.Cow:
			s = NewCowSprite()
		case simulation.Bull:
			s = NewBullSprite()
		}
	}
	if s != nil {
		pos := agent.GetSprite().Position
		s.SetPosition(pos.X, pos.Y)
		mw.SpriteMap[agent.GetID()] = s
	}
}
