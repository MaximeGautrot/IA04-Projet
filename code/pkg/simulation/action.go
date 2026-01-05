package simulation

type Action interface {
	Execute(a Agent, env *Environment)
	evaluateUtility(a Agent, env *Environment) float64
}
