package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestSkillUseFailed(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the skill use failed manager
	manager := NewSkillUseFailedManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for different failure types
	testCases := []struct {
		name            string
		packetSwitch    string
		args            map[string]interface{}
		expectedSkillID uint16
		expectedCause   uint8
		expectedMessage string
		isHomunculus    bool
	}{
		{
			name:         "Basic Skill Failure",
			packetSwitch: "0110",
			args: map[string]interface{}{
				"skillID": uint16(10),
				"btype":   uint16(0),
				"itemId":  uint32(0),
				"flag":    uint32(0),
				"cause":   uint8(0),
				"unknown": uint8(0),
			},
			expectedSkillID: 10,
			expectedCause:   0,
			expectedMessage: "Basic",
			isHomunculus:    false,
		},
		{
			name:         "Insufficient SP",
			packetSwitch: "0110",
			args: map[string]interface{}{
				"skillID": uint16(20),
				"btype":   uint16(0),
				"itemId":  uint32(0),
				"flag":    uint32(0),
				"cause":   uint8(1),
				"unknown": uint8(0),
			},
			expectedSkillID: 20,
			expectedCause:   1,
			expectedMessage: "Insufficient SP",
			isHomunculus:    false,
		},
		{
			name:         "Missing Required Item",
			packetSwitch: "0110",
			args: map[string]interface{}{
				"skillID": uint16(30),
				"btype":   uint16(0),
				"itemId":  uint32(501),
				"flag":    uint32(0),
				"cause":   uint8(71),
				"unknown": uint8(0),
			},
			expectedSkillID: 30,
			expectedCause:   71,
			expectedMessage: "Missing Required Item - item 501",
			isHomunculus:    false,
		},
		{
			name:         "Base Fail Type",
			packetSwitch: "0110",
			args: map[string]interface{}{
				"skillID": uint16(1),
				"btype":   uint16(3),
				"itemId":  uint32(0),
				"flag":    uint32(0),
				"cause":   uint8(0),
				"unknown": uint8(0),
			},
			expectedSkillID: 1,
			expectedCause:   0,
			expectedMessage: "No chat",
			isHomunculus:    false,
		},
		{
			name:         "Resurrect Homunculus Failed",
			packetSwitch: "0110",
			args: map[string]interface{}{
				"skillID": uint16(247),
				"btype":   uint16(0),
				"itemId":  uint32(0),
				"flag":    uint32(0),
				"cause":   uint8(0),
				"unknown": uint8(0),
			},
			expectedSkillID: 247,
			expectedCause:   0,
			expectedMessage: "Basic",
			isHomunculus:    true,
		},
		{
			name:         "Call Homunculus Failed - No Vaporized",
			packetSwitch: "0110",
			args: map[string]interface{}{
				"skillID": uint16(243),
				"btype":   uint16(0),
				"itemId":  uint32(0),
				"flag":    uint32(0),
				"cause":   uint8(0),
				"unknown": uint8(0),
			},
			expectedSkillID: 243,
			expectedCause:   0,
			expectedMessage: "Basic",
			isHomunculus:    true,
		},
		{
			name:         "Call Homunculus Failed - Missing Item",
			packetSwitch: "0110",
			args: map[string]interface{}{
				"skillID": uint16(243),
				"btype":   uint16(0),
				"itemId":  uint32(0),
				"flag":    uint32(0),
				"cause":   uint8(71),
				"unknown": uint8(0),
			},
			expectedSkillID: 243,
			expectedCause:   71,
			expectedMessage: "Missing Required Item - item 0",
			isHomunculus:    true,
		},
		{
			name:         "Unknown Error",
			packetSwitch: "0110",
			args: map[string]interface{}{
				"skillID": uint16(50),
				"btype":   uint16(0),
				"itemId":  uint32(0),
				"flag":    uint32(0),
				"cause":   uint8(255), // Unknown cause
				"unknown": uint8(0),
			},
			expectedSkillID: 50,
			expectedCause:   255,
			expectedMessage: "Unknown error",
			isHomunculus:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.skill_use_failed", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Add the switch to the args
			tc.args["switch"] = tc.packetSwitch

			// Call the handler
			err := manager.handleSkillUseFailed(tc.args)
			if err != nil {
				t.Errorf("handleSkillUseFailed() returned error: %v", err)
			}

			// Check that the hook was called
			if !hookCalled {
				t.Error("Hook was not called")
			}

			// Check the hook result
			if hookResult == nil {
				t.Fatal("Hook result is nil")
			}

			// Check the skill ID
			if skillID, ok := hookResult["skillID"].(uint16); !ok || skillID != tc.expectedSkillID {
				t.Errorf("Expected skill ID %d, got %v", tc.expectedSkillID, skillID)
			}

			// Check the cause
			if cause, ok := hookResult["cause"].(uint8); !ok || cause != tc.expectedCause {
				t.Errorf("Expected cause %d, got %v", tc.expectedCause, cause)
			}

			// Check the error message
			if errorMessage, ok := hookResult["errorMessage"].(string); !ok || errorMessage != tc.expectedMessage {
				t.Errorf("Expected error message %q, got %q", tc.expectedMessage, errorMessage)
			}

			// Check the isHomunculus flag
			if isHomunculus, ok := hookResult["isHomunculus"].(bool); !ok || isHomunculus != tc.isHomunculus {
				t.Errorf("Expected isHomunculus %v, got %v", tc.isHomunculus, isHomunculus)
			}
		})
	}
}

// Test unhappy paths
func TestSkillUseFailedUnhappy(t *testing.T) {
	// Create a skill use failed manager
	manager := NewSkillUseFailedManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0110",
			// Missing skillID, btype, etc.
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillUseFailed(args)
		if err != nil {
			t.Errorf("handleSkillUseFailed() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillUseFailed(args)

		// Check that the skill ID and cause are zero
		if skillID, ok := result["skillID"].(uint16); !ok || skillID != 0 {
			t.Errorf("Expected skill ID 0, got %v", skillID)
		}
		if cause, ok := result["cause"].(uint8); !ok || cause != 0 {
			t.Errorf("Expected cause 0, got %v", cause)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":  "0110",
			"skillID": "not a uint16", // Wrong type
			"btype":   "not a uint16", // Wrong type
			"itemId":  "not a uint32", // Wrong type
			"flag":    "not a uint32", // Wrong type
			"cause":   "not a uint8",  // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillUseFailed(args)
		if err != nil {
			t.Errorf("handleSkillUseFailed() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillUseFailed(args)

		// Check that the skill ID and cause are zero
		if skillID, ok := result["skillID"].(uint16); !ok || skillID != 0 {
			t.Errorf("Expected skill ID 0, got %v", skillID)
		}
		if cause, ok := result["cause"].(uint8); !ok || cause != 0 {
			t.Errorf("Expected cause 0, got %v", cause)
		}
	})
}
