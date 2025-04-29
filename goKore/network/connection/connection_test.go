package connection

import (
	"testing"
	"time"
)

// TestConnectionStateConstants verifies that the connection state constants are defined correctly
func TestConnectionStateConstants(t *testing.T) {
	// Test that all connection states are defined with the correct values
	states := map[string]ConnectionState{
		"NOT_CONNECTED":              NOT_CONNECTED,
		"CONNECTED_TO_MASTER_SERVER": CONNECTED_TO_MASTER_SERVER,
		"CONNECTED_TO_LOGIN_SERVER":  CONNECTED_TO_LOGIN_SERVER,
		"CONNECTED_TO_CHAR_SERVER":   CONNECTED_TO_CHAR_SERVER,
		"IN_GAME":                    IN_GAME,
		"IN_GAME_BUT_UNINITIALIZED":  IN_GAME_BUT_UNINITIALIZED,
	}

	expectedValues := map[string]int{
		"NOT_CONNECTED":              1,
		"CONNECTED_TO_MASTER_SERVER": 2,
		"CONNECTED_TO_LOGIN_SERVER":  3,
		"CONNECTED_TO_CHAR_SERVER":   4,
		"IN_GAME":                    5,
		"IN_GAME_BUT_UNINITIALIZED":  -1,
	}

	for name, value := range states {
		expected, exists := expectedValues[name]
		if !exists {
			t.Errorf("Unexpected state: %s", name)
			continue
		}
		if int(value) != expected {
			t.Errorf("State %s has value %d, expected %d", name, value, expected)
		}
	}
}

// TestConnectionStateString verifies that the String method of ConnectionState works correctly
func TestConnectionStateString(t *testing.T) {
	tests := []struct {
		state ConnectionState
		want  string
	}{
		{NOT_CONNECTED, "NOT_CONNECTED"},
		{CONNECTED_TO_MASTER_SERVER, "CONNECTED_TO_MASTER_SERVER"},
		{CONNECTED_TO_LOGIN_SERVER, "CONNECTED_TO_LOGIN_SERVER"},
		{CONNECTED_TO_CHAR_SERVER, "CONNECTED_TO_CHAR_SERVER"},
		{IN_GAME, "IN_GAME"},
		{IN_GAME_BUT_UNINITIALIZED, "IN_GAME_BUT_UNINITIALIZED"},
		{ConnectionState(100), "UNKNOWN_STATE(100)"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("ConnectionState.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConnectionEventConstants verifies that the connection event constants are defined correctly
func TestConnectionEventConstants(t *testing.T) {
	// Test that all connection events are defined with the correct values
	events := map[string]ConnectionEvent{
		"EventConnecting":    EventConnecting,
		"EventConnected":     EventConnected,
		"EventDisconnecting": EventDisconnecting,
		"EventDisconnected":  EventDisconnected,
		"EventDataSent":      EventDataSent,
		"EventDataReceived":  EventDataReceived,
		"EventError":         EventError,
		"EventStateChange":   EventStateChange,
	}

	expectedValues := map[string]int{
		"EventConnecting":    0,
		"EventConnected":     1,
		"EventDisconnecting": 2,
		"EventDisconnected":  3,
		"EventDataSent":      4,
		"EventDataReceived":  5,
		"EventError":         6,
		"EventStateChange":   7,
	}

	for name, value := range events {
		expected, exists := expectedValues[name]
		if !exists {
			t.Errorf("Unexpected event: %s", name)
			continue
		}
		if int(value) != expected {
			t.Errorf("Event %s has value %d, expected %d", name, value, expected)
		}
	}
}

// TestConnectionEventString verifies that the String method of ConnectionEvent works correctly
func TestConnectionEventString(t *testing.T) {
	tests := []struct {
		event ConnectionEvent
		want  string
	}{
		{EventConnecting, "EventConnecting"},
		{EventConnected, "EventConnected"},
		{EventDisconnecting, "EventDisconnecting"},
		{EventDisconnected, "EventDisconnected"},
		{EventDataSent, "EventDataSent"},
		{EventDataReceived, "EventDataReceived"},
		{EventError, "EventError"},
		{EventStateChange, "EventStateChange"},
		{ConnectionEvent(100), "UnknownEvent"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.event.String(); got != tt.want {
				t.Errorf("ConnectionEvent.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBaseConnection tests the BaseConnection methods
func TestBaseConnection(t *testing.T) {
	// Create a new BaseConnection
	config := &ConnectionConfig{
		Host:    "example.com",
		Port:    8000,
		Timeout: 30 * time.Second,
	}
	conn := NewBaseConnection(config)

	// Test initial state
	if conn.GetState() != NOT_CONNECTED {
		t.Errorf("Initial state = %v, want %v", conn.GetState(), NOT_CONNECTED)
	}

	// Test GetConfig
	if conn.GetConfig() != config {
		t.Errorf("GetConfig() = %v, want %v", conn.GetConfig(), config)
	}

	// Test SetConfig
	newConfig := &ConnectionConfig{
		Host:    "newhost.com",
		Port:    9000,
		Timeout: 60 * time.Second,
	}
	conn.SetConfig(newConfig)
	if conn.GetConfig() != newConfig {
		t.Errorf("GetConfig() after SetConfig = %v, want %v", conn.GetConfig(), newConfig)
	}

	// Test GetLastError (should be nil initially)
	if conn.GetLastError() != nil {
		t.Errorf("Initial GetLastError() = %v, want nil", conn.GetLastError())
	}

	// Test IsIdle
	if !conn.IsIdle(1 * time.Millisecond) {
		t.Errorf("IsIdle() = false, want true for new connection")
	}
}

// TestBaseConnectionCallbacks tests the callback registration and triggering
func TestBaseConnectionCallbacks(t *testing.T) {
	// Create a new BaseConnection
	conn := NewBaseConnection(&ConnectionConfig{})

	// Create a channel to receive events
	eventCh := make(chan ConnectionEvent, 10)
	dataCh := make(chan interface{}, 10)

	// Register a callback for state change events
	conn.RegisterCallback(EventStateChange, func(event ConnectionEvent, data interface{}) {
		eventCh <- event
		dataCh <- data
	})

	// Trigger a state change
	conn.SetState(CONNECTED_TO_MASTER_SERVER)

	// Check that the callback was called
	select {
	case event := <-eventCh:
		if event != EventStateChange {
			t.Errorf("Event = %v, want %v", event, EventStateChange)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timeout waiting for event")
	}

	// Check the data
	select {
	case data := <-dataCh:
		stateData, ok := data.(map[string]interface{})
		if !ok {
			t.Errorf("Data type = %T, want map[string]interface{}", data)
		} else {
			oldState, ok := stateData["oldState"].(ConnectionState)
			if !ok || oldState != NOT_CONNECTED {
				t.Errorf("oldState = %v, want %v", oldState, NOT_CONNECTED)
			}
			newState, ok := stateData["newState"].(ConnectionState)
			if !ok || newState != CONNECTED_TO_MASTER_SERVER {
				t.Errorf("newState = %v, want %v", newState, CONNECTED_TO_MASTER_SERVER)
			}
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timeout waiting for data")
	}

	// Test unregistering a callback
	callback := func(event ConnectionEvent, data interface{}) {
		t.Error("Unregistered callback was called")
	}
	conn.RegisterCallback(EventError, callback)
	conn.UnregisterCallback(EventError, callback)

	// Trigger the error event (should not call the unregistered callback)
	conn.setLastError(nil) // This is a bit of a hack, but it triggers the error event
}
