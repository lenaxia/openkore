package actor

import (
	"sync"
)

// ActorList is a thread-safe collection of actors
type ActorList struct {
	actors map[string]Actor // Map of actor ID (as string) to Actor
	mutex  sync.RWMutex     // Mutex for thread safety
}

// NewActorList creates a new empty ActorList
func NewActorList() *ActorList {
	return &ActorList{
		actors: make(map[string]Actor),
	}
}

// Add adds an actor to the list
func (l *ActorList) Add(actor Actor) {
	if actor == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.actors[string(actor.ID())] = actor
}

// Remove removes an actor from the list
func (l *ActorList) Remove(actor Actor) {
	if actor == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	delete(l.actors, string(actor.ID()))
}

// RemoveByID removes an actor from the list by ID
func (l *ActorList) RemoveByID(id []byte) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	delete(l.actors, string(id))
}

// GetByID returns an actor by ID
func (l *ActorList) GetByID(id []byte) Actor {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.actors[string(id)]
}

// GetAll returns all actors in the list
func (l *ActorList) GetAll() []Actor {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	actors := make([]Actor, 0, len(l.actors))
	for _, actor := range l.actors {
		actors = append(actors, actor)
	}

	return actors
}

// Count returns the number of actors in the list
func (l *ActorList) Count() int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return len(l.actors)
}

// Clear removes all actors from the list
func (l *ActorList) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.actors = make(map[string]Actor)
}

// ForEach executes a function for each actor in the list
func (l *ActorList) ForEach(fn func(Actor)) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	for _, actor := range l.actors {
		fn(actor)
	}
}

// Filter returns a new list containing only actors that match the filter function
func (l *ActorList) Filter(fn func(Actor) bool) *ActorList {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	result := NewActorList()
	for _, actor := range l.actors {
		if fn(actor) {
			result.Add(actor)
		}
	}

	return result
}

// Contains returns true if the list contains an actor with the given ID
func (l *ActorList) Contains(id []byte) bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	_, ok := l.actors[string(id)]
	return ok
}

// PlayersList is a specialized ActorList for players
type PlayersList struct {
	*ActorList
}

// NewPlayersList creates a new empty PlayersList
func NewPlayersList() *PlayersList {
	return &PlayersList{
		ActorList: NewActorList(),
	}
}

// Add adds a player to the list
func (l *PlayersList) Add(player *Player) {
	l.ActorList.Add(player)
}

// GetByID returns a player by ID
func (l *PlayersList) GetByID(id []byte) *Player {
	actor := l.ActorList.GetByID(id)
	if actor == nil {
		return nil
	}
	return actor.(*Player)
}

// MonstersList is a specialized ActorList for monsters
type MonstersList struct {
	*ActorList
}

// NewMonstersList creates a new empty MonstersList
func NewMonstersList() *MonstersList {
	return &MonstersList{
		ActorList: NewActorList(),
	}
}

// Add adds a monster to the list
func (l *MonstersList) Add(monster *Monster) {
	l.ActorList.Add(monster)
}

// GetByID returns a monster by ID
func (l *MonstersList) GetByID(id []byte) *Monster {
	actor := l.ActorList.GetByID(id)
	if actor == nil {
		return nil
	}
	return actor.(*Monster)
}

// NPCsList is a specialized ActorList for NPCs
type NPCsList struct {
	*ActorList
}

// NewNPCsList creates a new empty NPCsList
func NewNPCsList() *NPCsList {
	return &NPCsList{
		ActorList: NewActorList(),
	}
}

// Add adds an NPC to the list
func (l *NPCsList) Add(npc *NPC) {
	l.ActorList.Add(npc)
}

// GetByID returns an NPC by ID
func (l *NPCsList) GetByID(id []byte) *NPC {
	actor := l.ActorList.GetByID(id)
	if actor == nil {
		return nil
	}
	return actor.(*NPC)
}
