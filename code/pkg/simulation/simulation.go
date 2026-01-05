package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type TurnData struct {
	Tick              int
	HumansAlive       int
	AnimalsAlive      int
	VegetablesAlive   int
	CountPragmatic    int
	CountCautious     int
	CountSelfish      int
	CountCollectivist int
}

type Simulation struct {
	maxSteps    int
	MaxAnimals    int
	MaxPlants     int

	currentStep int
	isRunning   bool
	agents      []Agent
	environment Environment
	
	// Paramètres de spawn
	lambdaAnimals float64
	lambdaPlants  float64
	
	// Populations Initiales
	InitHumans    int
	InitAnimals   int
	InitPlants    int

	nextAnimalTime float64
	nextPlantTime  float64

	// Probabilités cumulatives pour les profils
	distPragmatic    float64
	distCautious     float64
	distSelfish      float64
	distCollectivist float64

	History         []TurnData
	globalIDCounter uint
}

func CreateSimulation(width, height int) *Simulation {
	rand.Seed(time.Now().UnixNano())
	return &Simulation{
		maxSteps:        5000,
		MaxAnimals:      100,
		MaxPlants:       100,
		currentStep:     0,
		isRunning:       false,
		agents:          []Agent{},
		environment:     CreateEnvironment(width, height),
		History:         []TurnData{},
		globalIDCounter: 0,
		nextAnimalTime:  0,
		nextPlantTime:   0,
	}
}

// Mise à jour pour inclure MaxSteps et les 4 Poids
func (s *Simulation) SetParameters(maxSteps, maxAnimals, maxPlants int, lambdaAnimals, lambdaPlants float64, initHumans, initAnimals, initPlants int, wPrag, wCaut, wSelf, wColl float64) {
	s.maxSteps = maxSteps
	s.lambdaAnimals = lambdaAnimals
	s.lambdaPlants = lambdaPlants
	s.InitHumans = initHumans
	s.InitAnimals = initAnimals
	s.InitPlants = initPlants
	
	// Enregistrement des limites
	s.MaxAnimals = maxAnimals
	s.MaxPlants = maxPlants

	// Calcul des probabilités cumulatives
	totalWeight := wPrag + wCaut + wSelf + wColl
	if totalWeight <= 0 { 
		totalWeight = 1; wPrag = 1 
	}
	
	s.distPragmatic = wPrag / totalWeight
	s.distCautious = s.distPragmatic + (wCaut / totalWeight)
	s.distSelfish = s.distCautious + (wSelf / totalWeight)
	s.distCollectivist = 1.0 

	s.nextAnimalTime = s.getExponentialTime(s.lambdaAnimals)
	s.nextPlantTime = s.getExponentialTime(s.lambdaPlants)
}

func (s *Simulation) getExponentialTime(lambda float64) float64 {
	if lambda <= 0 { return math.Inf(1) }
	u := rand.Float64()
	if u == 0 { u = 0.00000001 }
	return -math.Log(u) / lambda
}

func (s *Simulation) pickRandomAnimalType() AnimalType {
	roll := rand.Float64() * 100
	if roll < 60 { return Chicken }
	if roll < 90 { return Cow }
	return Bull
}

func (s *Simulation) pickRandomVegetableType() vegetableType {
	roll := rand.Float64() * 100
	if roll < 50 { return Berry }
	if roll < 80 { return Lettuce }
	return Carrot
}

func (s *Simulation) PickRandomProfile() Profile {
	r := rand.Float64()
	
	if r < s.distPragmatic { 
		return Pragmatic 
	} else if r < s.distCautious { 
		return Cautious 
	} else if r < s.distSelfish { 
		return Selfish 
	} else { 
		return Collectivist 
	}
}

func (s *Simulation) Start() {
	s.isRunning = true
	
	// Spawn Humains
	safeW := float64(s.environment.width - 16)
	safeH := float64(s.environment.height - 16)

	for i := 0; i < s.InitHumans; i++ {
		profile := s.PickRandomProfile()
		
		x := rand.Float64() * safeW
		y := rand.Float64() * safeH
		
		pos := CreatePosition(x, y)
		sprite := CreateSprite(pos.X, pos.Y, 16, 16)
		name := fmt.Sprintf("H-%d", i)
		h := CreateHuman(name, 100, sprite, 50, 100, profile, "Base")
		s.AddAgent(h)
	}
	
	for i := 0; i < s.InitAnimals; i++ {
		s.spawnAnimal()
	}

	for i := 0; i < s.InitPlants; i++ {
		s.spawnVegetable()
	}

	for _, a := range s.environment.agents {
		a.Start(&s.environment)
	}
}

func (s *Simulation) Stop() {
	s.isRunning = false
	for _, a := range s.environment.agents {
		a.Stop()
	}
}

func (s *Simulation) AddAgent(agent Agent) {
	s.globalIDCounter++
	agent.SetID(s.globalIDCounter)
	s.environment.AddAgent(agent)
	if s.isRunning {
		agent.Start(&s.environment)
	}
}

func (s *Simulation) Step() {
	if !s.isRunning { return }

	s.currentStep++
	if s.maxSteps > 0 && s.currentStep >= s.maxSteps {
		s.Stop()
		fmt.Println("Simulation terminée (Temps).")
		return
	}

	activeAgents := make([]Agent, 0)
	for _, a := range s.environment.agents {
		if a.IsAlive() {
			activeAgents = append(activeAgents, a)
		}
	}

	for _, a := range activeAgents { a.Sync() <- true }
	for _, a := range activeAgents { <-a.Done() }

	s.ManageSpawns()
	s.environment.RemoveDeadAgents()
	s.environment.RemoveDeadObjects()
	s.RecordStats()
}

func (s *Simulation) ManageSpawns() {
	maxSpawnsPerTick := 5
	s.nextAnimalTime -= 1.0
	c := 0
	for s.nextAnimalTime <= 0 {
		if c >= maxSpawnsPerTick { break }
		s.spawnAnimal()
		s.nextAnimalTime += s.getExponentialTime(s.lambdaAnimals)
		c++
	}
	s.nextPlantTime -= 1.0
	c = 0
	for s.nextPlantTime <= 0 {
		if c >= maxSpawnsPerTick { break }
		s.spawnVegetable()
		s.nextPlantTime += s.getExponentialTime(s.lambdaPlants)
		c++
	}
}

func (s *Simulation) spawnAnimal() {
	currentAnimals := 0
	for _, a := range s.environment.agents {
		if a.IsAlive() {
			if _, ok := a.(*Animal); ok {
				currentAnimals++
			}
		}
	}

	if currentAnimals >= s.MaxAnimals {
		return 
	}

	for i := 0; i < 10; i++ {
		typ := s.pickRandomAnimalType()
		size := 32
		if typ == Chicken { size = 20 }
		if typ == Bull { size = 48 }

		safeW := float64(s.environment.width - size)
		safeH := float64(s.environment.height - size)
		x := rand.Float64() * safeW
		y := rand.Float64() * safeH

		if s.environment.IsLocationFree(x, y, 40.0) {
			sprite := CreateSprite(x, y, size, size)
			animal := CreateAnimal("Wild", sprite, typ)
			s.AddAgent(animal)
			return
		}
	}
}

func (s *Simulation) spawnVegetable() {
	currentVegetables := 0
	for _, o := range s.environment.objects {
		if o.IsAlive() {
			currentVegetables++
		}
	}

	if currentVegetables >= s.MaxPlants {
		return
	}

	size := 16
	safeW := float64(s.environment.width - size)
	safeH := float64(s.environment.height - size)

	for i := 0; i < 10; i++ {
		x := rand.Float64() * safeW
		y := rand.Float64() * safeH
		if s.environment.IsLocationFree(x, y, 20.0) {
			s.globalIDCounter++
			typ := s.pickRandomVegetableType()
			sprite := CreateSprite(x, y, size, size)
			veg := CreateVegetable(s.globalIDCounter, "Plant", sprite, typ)
			s.environment.AddObject(veg)
			return
		}
	}
}

func (s *Simulation) RecordStats() {
	humans, animals, veg := 0, 0, 0
	cPrag, cCaut, cSelf, cColl := 0, 0, 0, 0
	
	for _, a := range s.environment.agents {
		if a.IsAlive() {
			if h, ok := a.(*Human); ok {
				humans++
				switch h.GetProfile() {
				case Pragmatic: cPrag++
				case Cautious: cCaut++
				case Selfish: cSelf++
				case Collectivist: cColl++
				}
			} else if _, ok := a.(*Animal); ok { animals++ }
		}
	}
	for _, o := range s.environment.objects { if o.IsAlive() { veg++ } }

	s.History = append(s.History, TurnData{
		Tick: s.currentStep,
		HumansAlive: humans, AnimalsAlive: animals, VegetablesAlive: veg,
		CountPragmatic: cPrag, CountCautious: cCaut, CountSelfish: cSelf, CountCollectivist: cColl,
	})
}

func (s *Simulation) GetAllAgents() []Agent { 
	return s.environment.agents 
}

func (s *Simulation) GetAllObjects() []Object { 
	return s.environment.objects 
}

func (s *Simulation) GetHistory() []TurnData { 
	return s.History 
}