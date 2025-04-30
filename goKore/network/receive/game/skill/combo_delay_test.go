package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestComboDelay(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the combo delay manager
	manager := NewComboDelayManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test case for combo delay
	t.Run("Basic Combo Delay", func(t *testing.T) {
		// Create a channel to receive hook events
		hookCalled := false
		var hookResult map[string]interface{}

		// Register a hook to capture the event
		hookManager.AddHook("character.combo_delay", func(hookName string, arg interface{}, userData interface{}) {
			hookCalled = true
			if result, ok := arg.(map[string]interface{}); ok {
				hookResult = result
			}
		}, nil)

		// Create packet data
		args := map[string]interface{}{
			"switch": "01D2",
			"ID":     uint32(100),  // Actor ID
			"delay":  uint32(5000), // Delay
		}

		// Call the handler
		err := manager.handleComboDelay(args)
		if err != nil {
			t.Errorf("handleComboDelay() returned error: %v", err)
		}

		// Check that the hook was called
		if !hookCalled {
			t.Error("Hook was not called")
		}

		// Check the hook result
		if hookResult == nil {
			t.Fatal("Hook result is nil")
		}

		// Check the actor ID
		if actorID, ok := hookResult["actorID"].(uint32); !ok || actorID != 100 {
			t.Errorf("Expected actor ID 100, got %v", actorID)
		}

		// Check the delay
		if delay, ok := hookResult["delay"].(uint32); !ok || delay != 5000 {
			t.Errorf("Expected delay 5000, got %v", delay)
		}

		// Check the combo delay
		if comboDelay, ok := hookResult["comboDelay"].(uint32); !ok || comboDelay != 5000 {
			t.Errorf("Expected combo delay 5000, got %v", comboDelay)
		}

		// Check the isOwnChar flag
		if isOwnChar, ok := hookResult["isOwnChar"].(bool); !ok {
			t.Error("isOwnChar not found in hook result or not a bool")
		} else if isOwnChar != false {
			t.Errorf("Expected isOwnChar to be false, got %v", isOwnChar)
		}
	})
}

// Test unhappy paths
func TestComboDelayUnhappy(t *testing.T) {
	// Create a combo delay manager
	manager := NewComboDelayManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "01D2",
			// Missing ID and delay
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleComboDelay(args)
		if err != nil {
			t.Errorf("handleComboDelay() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processComboDelay(args)

		// Check that the actor ID and delay are zero
		if actorID, ok := result["actorID"].(uint32); !ok || actorID != 0 {
			t.Errorf("Expected actor ID 0, got %v", actorID)
		}
		if delay, ok := result["delay"].(uint32); !ok || delay != 0 {
			t.Errorf("Expected delay 0, got %v", delay)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "01D2",
			"ID":     "not a uint32", // Wrong type
			"delay":  "not a uint32", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleComboDelay(args)
		if err != nil {
			t.Errorf("handleComboDelay() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processComboDelay(args)

		// Check that the actor ID and delay are zero
		if actorID, ok := result["actorID"].(uint32); !ok || actorID != 0 {
			t.Errorf("Expected actor ID 0, got %v", actorID)
		}
		if delay, ok := result["delay"].(uint32); !ok || delay != 0 {
			t.Errorf("Expected delay 0, got %v", delay)
		}
	})
}
