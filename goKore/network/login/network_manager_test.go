package login

import (
	"testing"
)

// TestNetworkManagerCompatibility tests that the MockNetworkManager can be used with the login.NetworkManager interface
func TestNetworkManagerCompatibility(t *testing.T) {
	// Create a mock network manager
	mockNetworkManager := NewMockNetworkManager()

	// Create a login config
	config := NewLoginConfig("testuser", "testpass", "testserver")

	// Create a login manager with the mock network manager
	loginManager := NewLoginManager(mockNetworkManager, config)

	// Verify that the login manager was created successfully
	if loginManager == nil {
		t.Error("Expected login manager to be created successfully")
	}
}

// TestSetStateChangeCallback tests that the SetStateChangeCallback method works correctly
func TestSetStateChangeCallback(t *testing.T) {
	// Create a mock network manager
	mockNetworkManager := NewMockNetworkManager()

	// Create variables to track state changes
	var oldState, newState int

	// Set the state change callback
	mockNetworkManager.SetStateChangeCallback(func(old, new int) {
		oldState = old
		newState = new
	})

	// Change the state
	mockNetworkManager.SetState(1) // ConnectedToMasterServer

	// Verify that the callback was called with the correct values
	if oldState != 0 || newState != 1 {
		t.Errorf("Expected oldState=0, newState=1, got oldState=%d, newState=%d", oldState, newState)
	}
}
