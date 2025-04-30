package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestHighJump(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the high jump manager
	manager := NewHighJumpManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for high jump
	testCases := []struct {
		name            string
		actorID         uint32
		x               uint16
		y               uint16
		moveSuccessful  bool
		expectedMessage string
	}{
		{
			name:            "Successful High Jump",
			actorID:         1001,
			x:               150,
			y:               100,
			moveSuccessful:  true,
			expectedMessage: "Actor_1001 instantly moved to 150, 100",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.high_jump", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Create packet data
			args := map[string]interface{}{
				"switch": "01FF",
				"ID":     tc.actorID,
				"x":      tc.x,
				"y":      tc.y,
			}

			// Call the handler
			err := manager.handleHighJump(args)
			if err != nil {
				t.Errorf("handleHighJump() returned error: %v", err)
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
			if actorID, ok := hookResult["actorID"].(uint32); !ok || actorID != tc.actorID {
				t.Errorf("Expected actor ID %d, got %v", tc.actorID, actorID)
			}

			// Check the x coordinate
			if x, ok := hookResult["x"].(uint16); !ok || x != tc.x {
				t.Errorf("Expected x %d, got %v", tc.x, x)
			}

			// Check the y coordinate
			if y, ok := hookResult["y"].(uint16); !ok || y != tc.y {
				t.Errorf("Expected y %d, got %v", tc.y, y)
			}

			// Check the move successful flag
			if moveSuccessful, ok := hookResult["moveSuccessful"].(bool); !ok || moveSuccessful != tc.moveSuccessful {
				t.Errorf("Expected moveSuccessful %v, got %v", tc.moveSuccessful, moveSuccessful)
			}

			// Check the message
			if message, ok := hookResult["message"].(string); !ok || message != tc.expectedMessage {
				t.Errorf("Expected message %q, got %q", tc.expectedMessage, message)
			}
		})
	}
}

// Test unhappy paths
func TestHighJumpUnhappy(t *testing.T) {
	// Create a high jump manager
	manager := NewHighJumpManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "01FF",
			// Missing ID, x, and y
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleHighJump(args)
		if err != nil {
			t.Errorf("handleHighJump() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processHighJump(args)

		// Check that the actor ID is zero
		if actorID, ok := result["actorID"].(uint32); !ok || actorID != 0 {
			t.Errorf("Expected actor ID 0, got %v", actorID)
		}

		// Check that the x coordinate is zero
		if x, ok := result["x"].(uint16); !ok || x != 0 {
			t.Errorf("Expected x 0, got %v", x)
		}

		// Check that the y coordinate is zero
		if y, ok := result["y"].(uint16); !ok || y != 0 {
			t.Errorf("Expected y 0, got %v", y)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "01FF",
			"ID":     "not a uint32", // Wrong type
			"x":      "not a uint16", // Wrong type
			"y":      "not a uint16", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleHighJump(args)
		if err != nil {
			t.Errorf("handleHighJump() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processHighJump(args)

		// Check that the actor ID is zero
		if actorID, ok := result["actorID"].(uint32); !ok || actorID != 0 {
			t.Errorf("Expected actor ID 0, got %v", actorID)
		}

		// Check that the x coordinate is zero
		if x, ok := result["x"].(uint16); !ok || x != 0 {
			t.Errorf("Expected x 0, got %v", x)
		}

		// Check that the y coordinate is zero
		if y, ok := result["y"].(uint16); !ok || y != 0 {
			t.Errorf("Expected y 0, got %v", y)
		}
	})
}
