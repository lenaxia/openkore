package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestAttackRange(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the attack range manager
	manager := NewAttackRangeManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test case for attack range
	t.Run("Basic Attack Range", func(t *testing.T) {
		// Create a channel to receive hook events
		hookCalled := false
		var hookResult map[string]interface{}

		// Register a hook to capture the event
		hookManager.AddHook("character.attack_range", func(hookName string, arg interface{}, userData interface{}) {
			hookCalled = true
			if result, ok := arg.(map[string]interface{}); ok {
				hookResult = result
			}
		}, nil)

		// Create packet data
		args := map[string]interface{}{
			"switch": "013A",
			"type":   uint16(5), // Attack range of 5
		}

		// Call the handler
		err := manager.handleAttackRange(args)
		if err != nil {
			t.Errorf("handleAttackRange() returned error: %v", err)
		}

		// Check that the hook was called
		if !hookCalled {
			t.Error("Hook was not called")
		}

		// Check the hook result
		if hookResult == nil {
			t.Fatal("Hook result is nil")
		}

		// Check the attack range
		if attackRange, ok := hookResult["attackRange"].(uint16); !ok || attackRange != 5 {
			t.Errorf("Expected attack range 5, got %v", attackRange)
		}

		// Check the config updated values
		if configUpdated, ok := hookResult["configUpdated"].(map[string]interface{}); ok {
			// Check attackDistance
			if attackDistance, ok := configUpdated["attackDistance"].(uint16); !ok || attackDistance != 5 {
				t.Errorf("Expected attackDistance 5, got %v", attackDistance)
			}

			// Check attackMaxDistance
			if attackMaxDistance, ok := configUpdated["attackMaxDistance"].(uint16); !ok || attackMaxDistance != 5 {
				t.Errorf("Expected attackMaxDistance 5, got %v", attackMaxDistance)
			}
		} else {
			t.Error("configUpdated not found in hook result or not a map")
		}
	})
}

// Test unhappy paths
func TestAttackRangeUnhappy(t *testing.T) {
	// Create an attack range manager
	manager := NewAttackRangeManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "013A",
			// Missing type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleAttackRange(args)
		if err != nil {
			t.Errorf("handleAttackRange() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processAttackRange(args)

		// Check that the attack range is zero
		if attackRange, ok := result["attackRange"].(uint16); !ok || attackRange != 0 {
			t.Errorf("Expected attack range 0, got %v", attackRange)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "013A",
			"type":   "not a uint16", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleAttackRange(args)
		if err != nil {
			t.Errorf("handleAttackRange() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processAttackRange(args)

		// Check that the attack range is zero
		if attackRange, ok := result["attackRange"].(uint16); !ok || attackRange != 0 {
			t.Errorf("Expected attack range 0, got %v", attackRange)
		}
	})
}
