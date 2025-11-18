// package frontend (dans window.go)
package frontend

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	imgRooster *ebiten.Image
	imgCow    *ebiten.Image
	imgBull    *ebiten.Image
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
}


type MainWindow struct {
	sprites []Sprite
}

func (mw *MainWindow) Update() error {
	for _, s := range mw.sprites {
		s.Update()
	}
	return nil
}

func (mw *MainWindow) Draw(screen *ebiten.Image) {
	for _, s := range mw.sprites {
		s.Draw(screen)
	}
}

func (mw *MainWindow) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 240
}

func (mw *MainWindow) Run() {
	
	loadSpriteSheets()

	rooster1 := NewRoosterSprite()
	rooster1.SetPosition(50, 100)
	rooster1.SetAnimationRow(2)

	cow1 := NewCowSprite()
	cow1.SetPosition(150, 50)
	cow1.SetAnimationRow(0)
	
	bull1 := NewBullSprite()
	bull1.SetPosition(200, 150)
	bull1.SetAnimationRow(4) 

	mw.sprites = []Sprite{rooster1, cow1, bull1}

	ebiten.SetWindowSize(640, 480) 
	ebiten.SetWindowTitle("Simulation")

	if err := ebiten.RunGame(mw); err != nil {
		log.Fatal(err)
	}
}