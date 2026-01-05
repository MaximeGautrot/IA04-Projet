package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"ia04project/pkg/frontend"
	"ia04project/pkg/simulation"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	StateConfig = iota
	StateSimulation
	StateStats
)

type App struct {
	State        int
	ConfigScreen *frontend.ConfigScreen
	MainWindow   *frontend.MainWindow
	GraphScreen  *frontend.GraphScreen
	Sim          *simulation.Simulation
}

func NewApp() *App {
	rand.Seed(time.Now().UnixNano())
	return &App{
		State:        StateConfig,
		ConfigScreen: frontend.NewConfigScreen(),
	}
}

func (a *App) Update() error {
	switch a.State {
	case StateConfig:
		a.ConfigScreen.Update()
		if a.ConfigScreen.IsDone {
			a.StartSimulation()
			a.State = StateSimulation
		}

	case StateSimulation:
		if a.MainWindow != nil {
			a.MainWindow.Update()
			if a.MainWindow.IsFinished {
				a.GraphScreen = frontend.NewGraphScreen(a.Sim.GetHistory())
				a.State = StateStats
			}
		}

	case StateStats:
		a.GraphScreen.Update()
	}
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	switch a.State {
	case StateConfig:
		a.ConfigScreen.Draw(screen)
	case StateSimulation:
		if a.MainWindow != nil {
			a.MainWindow.Draw(screen)
		}
	case StateStats:
		if a.GraphScreen != nil {
			a.GraphScreen.Draw(screen)
		}
	}
}

func (a *App) Layout(w, h int) (int, int) {
	if a.State == StateSimulation {
		return 1050, 600 // Largeur augmentée pour Sidebar + Jeu
	}
	return 800, 600
}

func (a *App) StartSimulation() {
	params := a.ConfigScreen.Params
	
	fmt.Println("Lancement de la simulation...")
	
	envWidth, envHeight := 800, 600
	sim := simulation.CreateSimulation(envWidth, envHeight) // maxSteps passé après
	
	// Configuration avec TOUS les paramètres de la fenêtre de config
	sim.SetParameters(
		params.MaxSteps,
        params.MaxAnimals,
        params.MaxPlants,

		params.LambdaAnimals, 
		params.LambdaPlants,
		params.InitHumans,    
		params.InitAnimals,   
		params.InitPlants,    
		
		// Proportions Profils dynamiques
		params.WeightPragmatic, 
		params.WeightCautious, 
		params.WeightSelfish, 
		params.WeightCollectivist,
	)

	sim.Start()

	a.Sim = sim
	a.MainWindow = frontend.NewMainWindow(sim)
}

func main() {
	ebiten.SetWindowSize(1050, 600)
	ebiten.SetWindowTitle("IA04 - Simulation Préhistorique Avancée")
	
	app := NewApp()

	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}