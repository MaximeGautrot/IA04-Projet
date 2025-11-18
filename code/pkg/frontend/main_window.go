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
	imgHumanIdle *ebiten.Image
	imgHumanWalk *ebiten.Image
	imgHumanRun  *ebiten.Image
	imgHumanHurt *ebiten.Image
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

	human1 := NewHumanSprite()
	human1.SetPosition(100, 120)
	human1.SetAnimationRow(4)

	human2 := NewHumanSprite()
	human2.SetPosition(250, 80)
	human2.SetAnimationRow(1)

	human3 := NewHumanSprite()
	human3.SetPosition(250, 180)
	human3.SetAnimationRow(10)

	mw.sprites = []Sprite{rooster1, cow1, bull1, human1, human2, human3}

	ebiten.SetWindowSize(640, 480) 
	ebiten.SetWindowTitle("Simulation")

	if err := ebiten.RunGame(mw); err != nil {
		log.Fatal(err)
	}
}