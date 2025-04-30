package skill

import (
	"testing"
	"time"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestCastCancelled(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the cast cancelled manager
	manager := NewCastCancelledManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for different packet types
	testCases := []struct {
		name            string
		packetSwitch    string
		args            map[string]interface{}
		expectedActorID uint32
	}{
		{
			name:         "Basic Cast Cancelled (01B9)",
			packetSwitch: "01B9",
			args: map[string]interface{}{
				"ID": uint32(100),
			},
			expectedActorID: 100,
		},
		{
			name:         "Expanded Cast Cancelled (08CD)",
			packetSwitch: "08CD",
			args: map[string]interface{}{
				"ID": uint32(200),
			},
			expectedActorID: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.cast_cancelled", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Add the switch to the args
			tc.args["switch"] = tc.packetSwitch

			// Call the handler
			err := manager.handleCastCancelled(tc.args)
			if err != nil {
				t.Errorf("handleCastCancelled() returned error: %v", err)
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
			if actorID, ok := hookResult["actorID"].(uint32); !ok || actorID != tc.expectedActorID {
				t.Errorf("Expected actor ID %d, got %v", tc.expectedActorID, actorID)
			}

			// Check that cancelledTime is set to a reasonable value (within the last second)
			if cancelledTime, ok := hookResult["cancelledTime"].(time.Time); !ok {
				t.Error("cancelledTime not found in hook result or not a time.Time")
			} else if time.Since(cancelledTime) > time.Second {
				t.Errorf("cancelledTime is too old: %v", cancelledTime)
			}

			// Check the isOwnCharacter flag
			if _, ok := hookResult["isOwnCharacter"].(bool); !ok {
				t.Error("isOwnCharacter not found in hook result or not a bool")
			}
		})
	}
}

// Test unhappy paths
func TestCastCancelledUnhappy(t *testing.T) {
	// Create a cast cancelled manager
	manager := NewCastCancelledManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "01B9",
			// Missing ID
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleCastCancelled(args)
		if err != nil {
			t.Errorf("handleCastCancelled() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processCastCancelled(args)

		// Check that the actor ID is zero
		if actorID, ok := result["actorID"].(uint32); !ok || actorID != 0 {
			t.Errorf("Expected actor ID 0, got %v", actorID)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "01B9",
			"ID":     "not a uint32", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleCastCancelled(args)
		if err != nil {
			t.Errorf("handleCastCancelled() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processCastCancelled(args)

		// Check that the actor ID is zero
		if actorID, ok := result["actorID"].(uint32); !ok || actorID != 0 {
			t.Errorf("Expected actor ID 0, got %v", actorID)
		}
	})
}
