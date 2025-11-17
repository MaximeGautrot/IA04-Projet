package simulation

type Simulation struct {
	maxSteps    int
	currentStep int
	isRunning   bool
	agents      []Agent
	environment Environment
}
