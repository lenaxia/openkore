package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleISVRDisconnect(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("core.isvr_disconnect", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create an ISVR disconnect manager for testing
	manager := NewISVRDisconnectManager(hookManager)

	// Test case
	t.Run("ISVR Disconnect", func(t *testing.T) {
		// Call the handler
		err := manager.HandleISVRDisconnect(map[string]interface{}{})
		if err != nil {
			t.Errorf("HandleISVRDisconnect returned an error: %v", err)
		}

		// Get the result from the channel
		result := <-resultChan

		// Check that the message was formatted correctly
		message, ok := result["message"].(string)
		if !ok {
			t.Errorf("Expected message to be a string")
		} else if message != "Received the package 'isvr_disconnect'" {
			t.Errorf("Expected message %q, got %q", "Received the package 'isvr_disconnect'", message)
		}
	})
}

// Test for nil hook manager
func TestHandleISVRDisconnectNilHookManager(t *testing.T) {
	// Create an ISVR disconnect manager with nil hook manager
	manager := NewISVRDisconnectManager(nil)

	// Test case
	t.Run("ISVR Disconnect with Nil Hook Manager", func(t *testing.T) {
		// Call the handler - it should not panic
		err := manager.HandleISVRDisconnect(map[string]interface{}{})
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
	})
}
