package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestSkillExchangeItem(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the skill exchange item manager
	manager := NewSkillExchangeItemManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for different exchange types
	testCases := []struct {
		name             string
		exchangeType     uint16
		val              uint32
		expectedTypeName string
		expectedCommand  string
		expectedMessage  string
	}{
		{
			name:             "Change Material",
			exchangeType:     0,
			val:              123,
			expectedTypeName: "Change Material",
			expectedCommand:  "cm",
			expectedMessage:  "Change Material is ready. Use command 'cm' to continue.",
		},
		{
			name:             "Elemental Analysis Lv 1",
			exchangeType:     1,
			val:              456,
			expectedTypeName: "Elemental Analysis Lv 1",
			expectedCommand:  "analysis",
			expectedMessage:  "Four Spirit Analysis is ready. Use command 'analysis' to continue.",
		},
		{
			name:             "Elemental Analysis Lv 2",
			exchangeType:     2,
			val:              789,
			expectedTypeName: "Elemental Analysis Lv 2",
			expectedCommand:  "analysis",
			expectedMessage:  "Four Spirit Analysis is ready. Use command 'analysis' to continue.",
		},
		{
			name:             "Unknown Type",
			exchangeType:     99,
			val:              999,
			expectedTypeName: "Unknown",
			expectedCommand:  "",
			expectedMessage:  "Unknown skill exchange type.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.skill_exchange_item", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Create packet data
			args := map[string]interface{}{
				"switch": "0917",
				"type":   tc.exchangeType,
				"val":    tc.val,
			}

			// Call the handler
			err := manager.handleSkillExchangeItem(args)
			if err != nil {
				t.Errorf("handleSkillExchangeItem() returned error: %v", err)
			}

			// Check that the hook was called
			if !hookCalled {
				t.Error("Hook was not called")
			}

			// Check the hook result
			if hookResult == nil {
				t.Fatal("Hook result is nil")
			}

			// Check the exchange type
			if exchangeType, ok := hookResult["exchangeType"].(uint16); !ok || exchangeType != tc.exchangeType {
				t.Errorf("Expected exchange type %d, got %v", tc.exchangeType, exchangeType)
			}

			// Check the val
			if val, ok := hookResult["val"].(uint32); !ok || val != tc.val {
				t.Errorf("Expected val %d, got %v", tc.val, val)
			}

			// Check the exchange type name
			if typeName, ok := hookResult["exchangeTypeName"].(string); !ok || typeName != tc.expectedTypeName {
				t.Errorf("Expected type name %q, got %q", tc.expectedTypeName, typeName)
			}

			// Check the command
			if command, ok := hookResult["command"].(string); !ok || command != tc.expectedCommand {
				t.Errorf("Expected command %q, got %q", tc.expectedCommand, command)
			}

			// Check the message
			if message, ok := hookResult["message"].(string); !ok || message != tc.expectedMessage {
				t.Errorf("Expected message %q, got %q", tc.expectedMessage, message)
			}
		})
	}
}

// Test unhappy paths
func TestSkillExchangeItemUnhappy(t *testing.T) {
	// Create a skill exchange item manager
	manager := NewSkillExchangeItemManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0917",
			// Missing type and val
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillExchangeItem(args)
		if err != nil {
			t.Errorf("handleSkillExchangeItem() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillExchangeItem(args)

		// Check that the exchange type is zero
		if exchangeType, ok := result["exchangeType"].(uint16); !ok || exchangeType != 0 {
			t.Errorf("Expected exchange type 0, got %v", exchangeType)
		}

		// Check that the val is zero
		if val, ok := result["val"].(uint32); !ok || val != 0 {
			t.Errorf("Expected val 0, got %v", val)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0917",
			"type":   "not a uint16", // Wrong type
			"val":    "not a uint32", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillExchangeItem(args)
		if err != nil {
			t.Errorf("handleSkillExchangeItem() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillExchangeItem(args)

		// Check that the exchange type is zero
		if exchangeType, ok := result["exchangeType"].(uint16); !ok || exchangeType != 0 {
			t.Errorf("Expected exchange type 0, got %v", exchangeType)
		}

		// Check that the val is zero
		if val, ok := result["val"].(uint32); !ok || val != 0 {
			t.Errorf("Expected val 0, got %v", val)
		}
	})
}
