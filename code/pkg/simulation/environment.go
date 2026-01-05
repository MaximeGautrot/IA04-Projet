package simulation

import (
	"sync"
)

type Environment struct {
	width   int
	height  int
	agents  []Agent
	objects []Object
	mutex   sync.RWMutex
}

func CreateEnvironment(width int, height int) Environment {
	return Environment{
		width:   width,
		height:  height,
		agents:  []Agent{},
		objects: []Object{},
	}
}

func (e *Environment) AddAgent(agent Agent) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.agents = append(e.agents, agent)
}

func (e *Environment) AddObject(obj Object) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.objects = append(e.objects, obj)
}

func (e *Environment) RemoveObject(id uint) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	newObjects := []Object{}
	for _, o := range e.objects {
		if o.GetID() != id {
			newObjects = append(newObjects, o)
		}
	}
	e.objects = newObjects
}

func (e *Environment) RemoveDeadAgents() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	newAgents := []Agent{}
	for _, a := range e.agents {
		if a.IsAlive() {
			newAgents = append(newAgents, a)
		}
	}
	e.agents = newAgents
}

func (e *Environment) RemoveDeadObjects() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	newObjects := []Object{}
	for _, o := range e.objects {
		if o.IsAlive() {
			newObjects = append(newObjects, o)
		}
	}
	e.objects = newObjects
}

// IsPositionInside reste inchangé...
func (e *Environment) IsPositionInside(s Sprite, dx, dy float64) bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	newX := s.Position.X + dx
	newY := s.Position.Y + dy

	return !(newX < 0 ||
		newY < 0 ||
		newX+float64(s.width) > float64(e.width) ||
		newY+float64(s.height) > float64(e.height))
}

func (e *Environment) IsLocationFree(x, y, minDist float64) bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	for _, agent := range e.agents {
		if agent.IsAlive() {
			pos := agent.GetSprite().Position
			dist := (pos.X-x)*(pos.X-x) + (pos.Y-y)*(pos.Y-y)
			if dist < minDist*minDist {
				return false
			}
		}
	}

	for _, obj := range e.objects {
		if obj.IsAlive() {
			pos := obj.GetSprite().Position
			dist := (pos.X-x)*(pos.X-x) + (pos.Y-y)*(pos.Y-y)
			if dist < minDist*minDist {
				return false
			}
		}
	}

	return true
}