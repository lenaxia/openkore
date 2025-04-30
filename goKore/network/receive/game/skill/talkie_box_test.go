package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestTalkieBox(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the talkie box manager
	manager := NewTalkieBoxManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for talkie box
	testCases := []struct {
		name            string
		id              uint32
		message         string
		expectedMessage string
	}{
		{
			name:            "Basic Talkie Box",
			id:              1001,
			message:         "Hello, world!",
			expectedMessage: "Actor_1001's talkie box message: Hello, world!",
		},
		{
			name:            "Another Talkie Box",
			id:              2002,
			message:         "This is a test message.",
			expectedMessage: "Actor_2002's talkie box message: This is a test message.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.talkie_box", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Create packet data
			args := map[string]interface{}{
				"switch":  "0191",
				"ID":      tc.id,
				"message": tc.message,
			}

			// Call the handler
			err := manager.handleTalkieBox(args)
			if err != nil {
				t.Errorf("handleTalkieBox() returned error: %v", err)
			}

			// Check that the hook was called
			if !hookCalled {
				t.Error("Hook was not called")
			}

			// Check the hook result
			if hookResult == nil {
				t.Fatal("Hook result is nil")
			}

			// Check the ID
			if id, ok := hookResult["ID"].(uint32); !ok || id != tc.id {
				t.Errorf("Expected ID %d, got %v", tc.id, id)
			}

			// Check the actor name
			expectedActorName := manager.getActorName(tc.id)
			if actorName, ok := hookResult["actorName"].(string); !ok || actorName != expectedActorName {
				t.Errorf("Expected actor name %q, got %q", expectedActorName, actorName)
			}

			// Check the message
			if message, ok := hookResult["message"].(string); !ok || message != tc.message {
				t.Errorf("Expected message %q, got %q", tc.message, message)
			}

			// Check the display message
			if displayMessage, ok := hookResult["displayMessage"].(string); !ok || displayMessage != tc.expectedMessage {
				t.Errorf("Expected display message %q, got %q", tc.expectedMessage, displayMessage)
			}
		})
	}
}

// Test unhappy paths
func TestTalkieBoxUnhappy(t *testing.T) {
	// Create a talkie box manager
	manager := NewTalkieBoxManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0191",
			// Missing ID and message
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleTalkieBox(args)
		if err != nil {
			t.Errorf("handleTalkieBox() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processTalkieBox(args)

		// Check that the ID is zero
		if id, ok := result["ID"].(uint32); !ok || id != 0 {
			t.Errorf("Expected ID 0, got %v", id)
		}

		// Check that the message is empty
		if message, ok := result["message"].(string); !ok || message != "" {
			t.Errorf("Expected empty message, got %q", message)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":  "0191",
			"ID":      "not a uint32", // Wrong type
			"message": 123,            // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleTalkieBox(args)
		if err != nil {
			t.Errorf("handleTalkieBox() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processTalkieBox(args)

		// Check that the ID is zero
		if id, ok := result["ID"].(uint32); !ok || id != 0 {
			t.Errorf("Expected ID 0, got %v", id)
		}

		// Check that the message is empty
		if message, ok := result["message"].(string); !ok || message != "" {
			t.Errorf("Expected empty message, got %q", message)
		}
	})
}
