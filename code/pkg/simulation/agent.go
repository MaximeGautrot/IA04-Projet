package simulation

// Agent définit les comportements obligatoires
type Agent interface {
	IsAlive() bool
	GetSprite() Sprite
	GetID() uint
	Move(dx, dy float64, env *Environment)
	GetName() string
	GetHealth() int
	Kill()
	IsAttacked(damage int)
	SetID(id uint)
	GetEnergy() uint

	// Méthodes IA
	Percept(env *Environment)
	Deliberate()
	Act(env *Environment)

	// Méthodes Concurrence (Goroutines)
	Start(env *Environment)
	Sync() chan bool
	Done() chan bool
	Stop()
}

// AgentParams contient les données communes
type AgentParams struct {
	id     uint
	name   string
	health int
	alive  bool
	sprite Sprite

	// Channels pour la synchronisation
	syncChan chan bool
	doneChan chan bool
	stopChan chan bool
}

// NewAgentParams initialise les channels (Indispensable !)
func NewAgentParams(id uint, name string, health int, sprite Sprite) AgentParams {
	return AgentParams{
		id:       id,
		name:     name,
		health:   health,
		alive:    true,
		sprite:   sprite,
		syncChan: make(chan bool),
		doneChan: make(chan bool),
		stopChan: make(chan bool),
	}
}

// --- Getters et Setters ---

func (ap *AgentParams) GetName() string { 
	return ap.name 
}

func (ap *AgentParams) GetID() uint { 
	return ap.id 
}

func (ap *AgentParams) GetHealth() int { 
	return ap.health 
}

func (ap *AgentParams) GetSprite() Sprite { 
	return ap.sprite 
}

func (ap *AgentParams) GetEnergy() uint { 
	return 0 
}

func (ap *AgentParams) IsAlive() bool { 
	return ap.alive 
}

func (ap *AgentParams) Sync() chan bool { 
	return ap.syncChan 
}

func (ap *AgentParams) Done() chan bool { 
	return ap.doneChan 
}

func (ap *AgentParams) Kill() {
	ap.alive = false
}

func (ap *AgentParams) Stop() {
	// On ferme le channel pour arrêter proprement la goroutine si besoin
	select {
	case <-ap.stopChan:
		// déjà fermé
	default:
		close(ap.stopChan)
	}
}

func (ap *AgentParams) SetID(id uint) {
	ap.id = id
}

func (ap *AgentParams) IsAttacked(damage int) {
	ap.health -= damage
	if ap.health <= 0 {
		ap.alive = false
	}
}

func (ap *AgentParams) Move(dx, dy float64, env *Environment) {
	if !env.IsPositionInside(ap.sprite, dx, dy) {
		return
	}

	futureSprite := ap.sprite
	futureSprite.MovePosition(Vector{dx, dy})

	// Vérification simple de collision (sans mutex ici car géré en amont ou séquentiel)
	// Dans la version multithread, l'Environment gère les verrous, donc on peut lire env.agents
	// Mais attention aux lectures concurrentes. Idéalement env expose une méthode "IsColliding" thread-safe.
	// Pour simplifier ici, on suppose que le mouvement est validé.
	
	ap.sprite.MovePosition(Vector{dx, dy})
}

func (ap *AgentParams) Percept(env *Environment) {}
func (ap *AgentParams) Deliberate() {}
func (ap *AgentParams) Act(env *Environment) {}