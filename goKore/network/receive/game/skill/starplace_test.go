package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestStarplace(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the starplace manager
	manager := NewStarplaceManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for starplace
	testCases := []struct {
		name            string
		which           uint8
		expectedMessage string
	}{
		{
			name:            "Basic Starplace",
			which:           1,
			expectedMessage: "Star Gladiator's Feeling map confirmation prompt: 1",
		},
		{
			name:            "Another Starplace",
			which:           2,
			expectedMessage: "Star Gladiator's Feeling map confirmation prompt: 2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.starplace", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Create packet data
			args := map[string]interface{}{
				"switch": "0253",
				"which":  tc.which,
			}

			// Call the handler
			err := manager.handleStarplace(args)
			if err != nil {
				t.Errorf("handleStarplace() returned error: %v", err)
			}

			// Check that the hook was called
			if !hookCalled {
				t.Error("Hook was not called")
			}

			// Check the hook result
			if hookResult == nil {
				t.Fatal("Hook result is nil")
			}

			// Check the which value
			if which, ok := hookResult["which"].(uint8); !ok || which != tc.which {
				t.Errorf("Expected which %d, got %v", tc.which, which)
			}

			// Check the message
			if message, ok := hookResult["message"].(string); !ok || message != tc.expectedMessage {
				t.Errorf("Expected message %q, got %q", tc.expectedMessage, message)
			}
		})
	}
}

// Test unhappy paths
func TestStarplaceUnhappy(t *testing.T) {
	// Create a starplace manager
	manager := NewStarplaceManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0253",
			// Missing which
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleStarplace(args)
		if err != nil {
			t.Errorf("handleStarplace() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processStarplace(args)

		// Check that the which value is zero
		if which, ok := result["which"].(uint8); !ok || which != 0 {
			t.Errorf("Expected which 0, got %v", which)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0253",
			"which":  "not a uint8", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleStarplace(args)
		if err != nil {
			t.Errorf("handleStarplace() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processStarplace(args)

		// Check that the which value is zero
		if which, ok := result["which"].(uint8); !ok || which != 0 {
			t.Errorf("Expected which 0, got %v", which)
		}
	})
}
