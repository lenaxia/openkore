package login

import (
	"sync"
)

// LoginState represents the current state of the login process
type LoginState int

const (
	StateNotConnected LoginState = iota
	StateConnectedToMasterServer
	StateConnectedToCharServer
	StateConnectedToMapServer
	StateInGame
)

// String returns a string representation of the login state
func (s LoginState) String() string {
	switch s {
	case StateNotConnected:
		return "NotConnected"
	case StateConnectedToMasterServer:
		return "ConnectedToMasterServer"
	case StateConnectedToCharServer:
		return "ConnectedToCharServer"
	case StateConnectedToMapServer:
		return "ConnectedToMapServer"
	case StateInGame:
		return "InGame"
	default:
		return "Unknown"
	}
}

// StateObserver is notified when the state changes
type StateObserver interface {
	OnStateChange(oldState, newState LoginState)
}

// LoginStateManager manages the state transitions of the login process
type LoginStateManager struct {
	currentState LoginState
	mu           sync.RWMutex
	observers    []StateObserver
}

// NewLoginStateManager creates a new login state manager
func NewLoginStateManager() *LoginStateManager {
	return &LoginStateManager{
		currentState: StateNotConnected,
		observers:    make([]StateObserver, 0),
	}
}

// GetState returns the current login state
func (lsm *LoginStateManager) GetState() LoginState {
	lsm.mu.RLock()
	defer lsm.mu.RUnlock()
	return lsm.currentState
}

// SetState changes the current state and notifies observers
func (lsm *LoginStateManager) SetState(newState LoginState) {
	lsm.mu.Lock()
	oldState := lsm.currentState
	lsm.currentState = newState
	observers := make([]StateObserver, len(lsm.observers)) // Copy to avoid holding lock during callbacks
	copy(observers, lsm.observers)
	lsm.mu.Unlock()

	// Notify observers
	for _, observer := range observers {
		observer.OnStateChange(oldState, newState)
	}
}

// AddObserver registers a new state observer
func (lsm *LoginStateManager) AddObserver(observer StateObserver) {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()
	lsm.observers = append(lsm.observers, observer)
}

// RemoveObserver unregisters a state observer
func (lsm *LoginStateManager) RemoveObserver(observer StateObserver) {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()

	for i, obs := range lsm.observers {
		if obs == observer {
			// Remove observer by replacing it with the last element and truncating
			lsm.observers[i] = lsm.observers[len(lsm.observers)-1]
			lsm.observers = lsm.observers[:len(lsm.observers)-1]
			break
		}
	}
}
