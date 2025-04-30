package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestSkillMsg(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the skill message manager
	manager := NewSkillMsgManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test case for skill message
	t.Run("Basic Skill Message", func(t *testing.T) {
		// Create a channel to receive hook events
		hookCalled := false
		var hookResult map[string]interface{}

		// Register a hook to capture the event
		hookManager.AddHook("character.skill_msg", func(hookName string, arg interface{}, userData interface{}) {
			hookCalled = true
			if result, ok := arg.(map[string]interface{}); ok {
				hookResult = result
			}
		}, nil)

		// Create packet data
		args := map[string]interface{}{
			"switch": "0215",
			"id":     uint16(10), // Skill ID
			"msgid":  uint16(0),  // Message ID (will be incremented to 1 in the handler)
		}

		// Call the handler
		err := manager.handleSkillMsg(args)
		if err != nil {
			t.Errorf("handleSkillMsg() returned error: %v", err)
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
		if skillID, ok := hookResult["skillID"].(uint16); !ok || skillID != 10 {
			t.Errorf("Expected skill ID 10, got %v", hookResult["skillID"])
		}

		// Check the message ID
		if msgID, ok := hookResult["msgID"].(uint16); !ok || msgID != 0 {
			t.Errorf("Expected message ID 0, got %v", hookResult["msgID"])
		}

		// Check that a message was returned
		if message, ok := hookResult["message"].(string); !ok || message == "" {
			t.Errorf("Expected non-empty message, got %v", hookResult["message"])
		}
	})

	// Test case for unknown message ID
	t.Run("Unknown Message ID", func(t *testing.T) {
		// Create a channel to receive hook events
		hookCalled := false
		var hookResult map[string]interface{}

		// Register a hook to capture the event
		hookManager.AddHook("character.skill_msg", func(hookName string, arg interface{}, userData interface{}) {
			hookCalled = true
			if result, ok := arg.(map[string]interface{}); ok {
				hookResult = result
			}
		}, nil)

		// Create packet data with an unknown message ID
		args := map[string]interface{}{
			"switch": "0215",
			"id":     uint16(10),
			"msgid":  uint16(999), // Unknown message ID
		}

		// Call the handler
		err := manager.handleSkillMsg(args)
		if err != nil {
			t.Errorf("handleSkillMsg() returned error: %v", err)
		}

		// Check that the hook was called
		if !hookCalled {
			t.Error("Hook was not called")
		}

		// Check the hook result
		if hookResult == nil {
			t.Fatal("Hook result is nil")
		}

		// Check that an "unknown message" was returned
		if message, ok := hookResult["message"].(string); !ok || message == "" {
			t.Errorf("Expected unknown message text, got %v", hookResult["message"])
		} else if message[:7] != "Unknown" {
			t.Errorf("Expected message to start with 'Unknown', got %v", message)
		}
	})
}

// Test unhappy paths
func TestSkillMsgUnhappy(t *testing.T) {
	// Create a skill message manager
	manager := NewSkillMsgManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0215",
			// Missing id and msgid
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillMsg(args)
		if err != nil {
			t.Errorf("handleSkillMsg() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillMsg(args)

		// Check that the skill ID and message ID are zero
		if skillID, ok := result["skillID"].(uint16); !ok || skillID != 0 {
			t.Errorf("Expected skill ID 0, got %v", result["skillID"])
		}
		if msgID, ok := result["msgID"].(uint16); !ok || msgID != 0 {
			t.Errorf("Expected message ID 0, got %v", result["msgID"])
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0215",
			"id":     "not a uint16", // Wrong type
			"msgid":  "not a uint16", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillMsg(args)
		if err != nil {
			t.Errorf("handleSkillMsg() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillMsg(args)

		// Check that the skill ID and message ID are zero
		if skillID, ok := result["skillID"].(uint16); !ok || skillID != 0 {
			t.Errorf("Expected skill ID 0, got %v", result["skillID"])
		}
		if msgID, ok := result["msgID"].(uint16); !ok || msgID != 0 {
			t.Errorf("Expected message ID 0, got %v", result["msgID"])
		}
	})
}
