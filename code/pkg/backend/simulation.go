package simulation

type Simulation struct {
	maxSteps    int
	currentStep int
	isRunning   bool
	agents      []Agent
	environment Environment
}

func CreateSimulation(maxSteps int, environment Environment) *Simulation {
	return &Simulation{
		maxSteps:    maxSteps,
		currentStep: 0,
		isRunning:   false,
		agents:      []Agent{},
		environment: environment,
	}
}

func (s *Simulation) AddAgent(agent Agent) {
	s.agents = append(s.agents, agent)
	s.environment.AddAgent(agent)
}
