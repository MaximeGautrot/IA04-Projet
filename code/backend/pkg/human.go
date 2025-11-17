package simulation

type human struct {
	Agent
	hunger        float64
	riskProfile   string
	strategyType  string
	energy        float64
	currentAction Action
}
