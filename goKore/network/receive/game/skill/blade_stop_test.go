package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestBladeStop(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the blade stop manager
	manager := NewBladeStopManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for different blade stop states
	testCases := []struct {
		name            string
		sourceID        uint32
		targetID        uint32
		active          uint8
		expectedMessage string
	}{
		{
			name:            "Blade Stop Active",
			sourceID:        1001,
			targetID:        2001,
			active:          1,
			expectedMessage: "Blade Stop by Actor_1001 on Actor_2001 is active.",
		},
		{
			name:            "Blade Stop Deactivated",
			sourceID:        1002,
			targetID:        2002,
			active:          0,
			expectedMessage: "Blade Stop by Actor_1002 on Actor_2002 is deactivated.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.blade_stop", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Create packet data
			args := map[string]interface{}{
				"switch":   "01D1",
				"sourceID": tc.sourceID,
				"targetID": tc.targetID,
				"active":   tc.active,
			}

			// Call the handler
			err := manager.handleBladeStop(args)
			if err != nil {
				t.Errorf("handleBladeStop() returned error: %v", err)
			}

			// Check that the hook was called
			if !hookCalled {
				t.Error("Hook was not called")
			}

			// Check the hook result
			if hookResult == nil {
				t.Fatal("Hook result is nil")
			}

			// Check the source ID
			if sourceID, ok := hookResult["sourceID"].(uint32); !ok || sourceID != tc.sourceID {
				t.Errorf("Expected source ID %d, got %v", tc.sourceID, sourceID)
			}

			// Check the target ID
			if targetID, ok := hookResult["targetID"].(uint32); !ok || targetID != tc.targetID {
				t.Errorf("Expected target ID %d, got %v", tc.targetID, targetID)
			}

			// Check the active status
			if active, ok := hookResult["active"].(uint8); !ok || active != tc.active {
				t.Errorf("Expected active %d, got %v", tc.active, active)
			}

			// Check the message
			if message, ok := hookResult["message"].(string); !ok || message != tc.expectedMessage {
				t.Errorf("Expected message %q, got %q", tc.expectedMessage, message)
			}
		})
	}
}

// Test unhappy paths
func TestBladeStopUnhappy(t *testing.T) {
	// Create a blade stop manager
	manager := NewBladeStopManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "01D1",
			// Missing sourceID, targetID, and active
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleBladeStop(args)
		if err != nil {
			t.Errorf("handleBladeStop() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processBladeStop(args)

		// Check that the source ID is zero
		if sourceID, ok := result["sourceID"].(uint32); !ok || sourceID != 0 {
			t.Errorf("Expected source ID 0, got %v", sourceID)
		}

		// Check that the target ID is zero
		if targetID, ok := result["targetID"].(uint32); !ok || targetID != 0 {
			t.Errorf("Expected target ID 0, got %v", targetID)
		}

		// Check that the active status is zero
		if active, ok := result["active"].(uint8); !ok || active != 0 {
			t.Errorf("Expected active 0, got %v", active)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":   "01D1",
			"sourceID": "not a uint32", // Wrong type
			"targetID": "not a uint32", // Wrong type
			"active":   "not a uint8",  // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleBladeStop(args)
		if err != nil {
			t.Errorf("handleBladeStop() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processBladeStop(args)

		// Check that the source ID is zero
		if sourceID, ok := result["sourceID"].(uint32); !ok || sourceID != 0 {
			t.Errorf("Expected source ID 0, got %v", sourceID)
		}

		// Check that the target ID is zero
		if targetID, ok := result["targetID"].(uint32); !ok || targetID != 0 {
			t.Errorf("Expected target ID 0, got %v", targetID)
		}

		// Check that the active status is zero
		if active, ok := result["active"].(uint8); !ok || active != 0 {
			t.Errorf("Expected active 0, got %v", active)
		}
	})
}
