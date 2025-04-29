package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleErrors(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)
	disconnectedCalled := make(chan bool, 1)

	// Register hooks to capture results
	hookManager.AddHook("core.errors", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	hookManager.AddHook("disconnected", func(hookName string, arg interface{}, userData interface{}) {
		disconnectedCalled <- true
	}, nil)

	// Create a config for testing
	config := map[string]interface{}{
		"dcOnDisconnect":     1,
		"dcOnServerShutDown": 1,
		"dcOnServerClose":    1,
		"dcOnDualLogin":      1,
	}

	// Test cases
	testCases := []struct {
		name             string
		args             map[string]interface{}
		netState         int
		expectHook       bool
		expectDisconnect bool
		errorType        byte
	}{
		{
			name: "Server Shutdown",
			args: map[string]interface{}{
				"type": byte(ErrorServerShutdown),
			},
			netState:         NetworkStateInGame,
			expectHook:       true,
			expectDisconnect: true,
			errorType:        ErrorServerShutdown,
		},
		{
			name: "Server Closed",
			args: map[string]interface{}{
				"type": byte(ErrorServerClosed),
			},
			netState:         NetworkStateInGame,
			expectHook:       true,
			expectDisconnect: true,
			errorType:        ErrorServerClosed,
		},
		{
			name: "Dual Login",
			args: map[string]interface{}{
				"type": byte(ErrorDualLogin),
			},
			netState:         NetworkStateInGame,
			expectHook:       true,
			expectDisconnect: true,
			errorType:        ErrorDualLogin,
		},
		{
			name: "Out of Sync",
			args: map[string]interface{}{
				"type": byte(ErrorOutOfSync),
			},
			netState:         NetworkStateInGame,
			expectHook:       true,
			expectDisconnect: false, // This error type is excluded from auto-disconnect
			errorType:        ErrorOutOfSync,
		},
		{
			name: "Unknown Error",
			args: map[string]interface{}{
				"type": byte(255),
			},
			netState:         NetworkStateInGame,
			expectHook:       true,
			expectDisconnect: true,
			errorType:        255,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create an errors manager for testing
			manager := NewErrorsManager(hookManager, tc.netState, config, "TestApp")

			// Call the handler
			err := manager.HandleErrors(tc.args)
			if err != nil {
				t.Errorf("HandleErrors returned an error: %v", err)
			}

			// Check if disconnected hook was called
			if tc.netState == NetworkStateInGame {
				select {
				case <-disconnectedCalled:
					// This is expected
				default:
					t.Errorf("Expected disconnected hook to be called")
				}
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if errorType, ok := result["type"].(byte); !ok || errorType != tc.errorType {
					t.Errorf("Expected error type %v, got %v", tc.errorType, result["type"])
				}

				// Check auto disconnect flag
				if autoDisconnect, ok := result["autoDisconnect"].(bool); !ok || autoDisconnect != tc.expectDisconnect {
					t.Errorf("Expected autoDisconnect %v, got %v", tc.expectDisconnect, autoDisconnect)
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleErrorsUnhappy(t *testing.T) {
	// Create an errors manager with nil hook manager
	manager := NewErrorsManager(nil, NetworkStateDisconnected, nil, "TestApp")

	// Test cases for unhappy paths
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Missing Error Type",
			args: map[string]interface{}{},
		},
		{
			name: "Invalid Error Type",
			args: map[string]interface{}{
				"type": "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := manager.HandleErrors(tc.args)
			if err != nil {
				t.Errorf("Expected no error for unhappy path, but got: %v", err)
			}
		})
	}
}
