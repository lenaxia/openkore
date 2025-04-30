package login

import (
	"sync"
	"testing"
)

// MockStateObserver implements the StateObserver interface for testing
type MockStateObserver struct {
	oldState     LoginState
	newState     LoginState
	callCount    int
	notifyLock   sync.Mutex
	notifyCalled chan struct{}
}

func NewMockStateObserver() *MockStateObserver {
	return &MockStateObserver{
		notifyCalled: make(chan struct{}, 10), // Buffer to avoid blocking
	}
}

func (m *MockStateObserver) OnStateChange(oldState, newState LoginState) {
	m.notifyLock.Lock()
	defer m.notifyLock.Unlock()

	m.oldState = oldState
	m.newState = newState
	m.callCount++
	m.notifyCalled <- struct{}{}
}

func (m *MockStateObserver) GetCallCount() int {
	m.notifyLock.Lock()
	defer m.notifyLock.Unlock()
	return m.callCount
}

func (m *MockStateObserver) GetLastStates() (LoginState, LoginState) {
	m.notifyLock.Lock()
	defer m.notifyLock.Unlock()
	return m.oldState, m.newState
}

func TestLoginStateManager_InitialState(t *testing.T) {
	lsm := NewLoginStateManager()

	if lsm.GetState() != StateNotConnected {
		t.Errorf("Expected initial state to be StateNotConnected, got %v", lsm.GetState())
	}
}

func TestLoginStateManager_SetState(t *testing.T) {
	lsm := NewLoginStateManager()

	// Change state
	lsm.SetState(StateConnectedToMasterServer)

	if lsm.GetState() != StateConnectedToMasterServer {
		t.Errorf("Expected state to be StateConnectedToMasterServer, got %v", lsm.GetState())
	}
}

func TestLoginStateManager_ObserverNotification(t *testing.T) {
	lsm := NewLoginStateManager()
	observer := NewMockStateObserver()

	// Register observer
	lsm.AddObserver(observer)

	// Change state
	lsm.SetState(StateConnectedToMasterServer)

	// Wait for notification
	<-observer.notifyCalled

	// Check if observer was notified
	oldState, newState := observer.GetLastStates()
	if oldState != StateNotConnected || newState != StateConnectedToMasterServer {
		t.Errorf("Expected state transition from StateNotConnected to StateConnectedToMasterServer, got %v to %v", oldState, newState)
	}

	if observer.GetCallCount() != 1 {
		t.Errorf("Expected observer to be called once, got %d", observer.GetCallCount())
	}
}

func TestLoginStateManager_MultipleObservers(t *testing.T) {
	lsm := NewLoginStateManager()
	observer1 := NewMockStateObserver()
	observer2 := NewMockStateObserver()

	// Register observers
	lsm.AddObserver(observer1)
	lsm.AddObserver(observer2)

	// Change state
	lsm.SetState(StateConnectedToMasterServer)

	// Wait for notifications
	<-observer1.notifyCalled
	<-observer2.notifyCalled

	// Check if both observers were notified
	if observer1.GetCallCount() != 1 {
		t.Errorf("Expected observer1 to be called once, got %d", observer1.GetCallCount())
	}

	if observer2.GetCallCount() != 1 {
		t.Errorf("Expected observer2 to be called once, got %d", observer2.GetCallCount())
	}
}

func TestLoginStateManager_RemoveObserver(t *testing.T) {
	lsm := NewLoginStateManager()
	observer := NewMockStateObserver()

	// Register observer
	lsm.AddObserver(observer)

	// Change state
	lsm.SetState(StateConnectedToMasterServer)

	// Wait for notification
	<-observer.notifyCalled

	// Remove observer
	lsm.RemoveObserver(observer)

	// Change state again
	lsm.SetState(StateConnectedToCharServer)

	// Check if observer was not notified again
	if observer.GetCallCount() > 1 {
		t.Errorf("Expected observer to be called only once, got %d", observer.GetCallCount())
	}
}

func TestLoginStateManager_StateTransitionSequence(t *testing.T) {
	lsm := NewLoginStateManager()
	observer := NewMockStateObserver()

	// Register observer
	lsm.AddObserver(observer)

	// Define expected state sequence
	stateSequence := []LoginState{
		StateConnectedToMasterServer,
		StateConnectedToCharServer,
		StateConnectedToMapServer,
		StateInGame,
	}

	// Transition through states
	for _, state := range stateSequence {
		lsm.SetState(state)
		<-observer.notifyCalled
	}

	// Check final state
	if lsm.GetState() != StateInGame {
		t.Errorf("Expected final state to be StateInGame, got %v", lsm.GetState())
	}

	// Check call count
	if observer.GetCallCount() != len(stateSequence) {
		t.Errorf("Expected observer to be called %d times, got %d", len(stateSequence), observer.GetCallCount())
	}
}

func TestLoginStateManager_ConcurrentAccess(t *testing.T) {
	lsm := NewLoginStateManager()
	observer := NewMockStateObserver()

	// Register observer
	lsm.AddObserver(observer)

	// Number of concurrent goroutines
	numGoroutines := 10

	// Use WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch goroutines that change state
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Each goroutine sets a different state
			state := LoginState((id % 5) + 1) // Ensure state is valid (1-5)
			lsm.SetState(state)
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check that observer was called the correct number of times
	// Note: Due to race conditions, we can't guarantee exactly how many times
	// the observer will be called, but it should be at least once
	if observer.GetCallCount() < 1 {
		t.Errorf("Expected observer to be called at least once, got %d", observer.GetCallCount())
	}
}
