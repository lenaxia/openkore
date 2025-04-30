package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestGospelBuffAligned(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the gospel buff manager
	manager := NewGospelBuffManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for different gospel buff effects
	testCases := []struct {
		name            string
		statusID        uint32
		expectedMessage string
	}{
		{
			name:            "Remove Abnormal Status",
			statusID:        21,
			expectedMessage: "All abnormal status effects have been removed.",
		},
		{
			name:            "Immunity to Abnormal Status",
			statusID:        22,
			expectedMessage: "You will be immune to abnormal status effects for the next minute.",
		},
		{
			name:            "Increased Max HP",
			statusID:        23,
			expectedMessage: "Your Max HP will stay increased for the next minute.",
		},
		{
			name:            "Increased Max SP",
			statusID:        24,
			expectedMessage: "Your Max SP will stay increased for the next minute.",
		},
		{
			name:            "Increased Stats",
			statusID:        25,
			expectedMessage: "All of your Stats will stay increased for the next minute.",
		},
		{
			name:            "Holy Weapon",
			statusID:        28,
			expectedMessage: "Your weapon will remain blessed with Holy power for the next minute.",
		},
		{
			name:            "Holy Armor",
			statusID:        29,
			expectedMessage: "Your armor will remain blessed with Holy power for the next minute.",
		},
		{
			name:            "Increased Defense",
			statusID:        30,
			expectedMessage: "Your Defense will stay increased for the next 10 seconds.",
		},
		{
			name:            "Increased Attack",
			statusID:        31,
			expectedMessage: "Your Attack strength will stay increased for the next minute.",
		},
		{
			name:            "Increased Accuracy and Flee",
			statusID:        32,
			expectedMessage: "Your Accuracy and Flee Rate will stay increased for the next minute.",
		},
		{
			name:            "Full Strip Failed",
			statusID:        40,
			expectedMessage: "Full strip failed because of coating.",
		},
		{
			name:            "Unknown Status",
			statusID:        999,
			expectedMessage: "Unknown gospel buff effect.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.gospel_buff_aligned", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Create packet data
			args := map[string]interface{}{
				"switch": "0215",
				"ID":     tc.statusID,
			}

			// Call the handler
			err := manager.handleGospelBuffAligned(args)
			if err != nil {
				t.Errorf("handleGospelBuffAligned() returned error: %v", err)
			}

			// Check that the hook was called
			if !hookCalled {
				t.Error("Hook was not called")
			}

			// Check the hook result
			if hookResult == nil {
				t.Fatal("Hook result is nil")
			}

			// Check the status ID
			if statusID, ok := hookResult["statusID"].(uint32); !ok || statusID != tc.statusID {
				t.Errorf("Expected status ID %d, got %v", tc.statusID, statusID)
			}

			// Check the message
			if message, ok := hookResult["message"].(string); !ok || message != tc.expectedMessage {
				t.Errorf("Expected message %q, got %q", tc.expectedMessage, message)
			}
		})
	}
}

// Test unhappy paths
func TestGospelBuffUnhappy(t *testing.T) {
	// Create a gospel buff manager
	manager := NewGospelBuffManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0215",
			// Missing ID
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleGospelBuffAligned(args)
		if err != nil {
			t.Errorf("handleGospelBuffAligned() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processGospelBuffAligned(args)

		// Check that the status ID is zero
		if statusID, ok := result["statusID"].(uint32); !ok || statusID != 0 {
			t.Errorf("Expected status ID 0, got %v", statusID)
		}

		// Check that the message is for unknown status
		if message, ok := result["message"].(string); !ok || message != "Unknown gospel buff effect." {
			t.Errorf("Expected message %q, got %q", "Unknown gospel buff effect.", message)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0215",
			"ID":     "not a uint32", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleGospelBuffAligned(args)
		if err != nil {
			t.Errorf("handleGospelBuffAligned() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processGospelBuffAligned(args)

		// Check that the status ID is zero
		if statusID, ok := result["statusID"].(uint32); !ok || statusID != 0 {
			t.Errorf("Expected status ID 0, got %v", statusID)
		}

		// Check that the message is for unknown status
		if message, ok := result["message"].(string); !ok || message != "Unknown gospel buff effect." {
			t.Errorf("Expected message %q, got %q", "Unknown gospel buff effect.", message)
		}
	})
}
