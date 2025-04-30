package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestSkillDelete(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the skill delete manager
	manager := NewSkillDeleteManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test case for skill delete
	t.Run("Basic Skill Delete", func(t *testing.T) {
		// Create a channel to receive hook events
		hookCalled := false
		var hookResult map[string]interface{}

		// Register a hook to capture the event
		hookManager.AddHook("character.skill_delete", func(hookName string, arg interface{}, userData interface{}) {
			hookCalled = true
			if result, ok := arg.(map[string]interface{}); ok {
				hookResult = result
			}
		}, nil)

		// Create packet data
		args := map[string]interface{}{
			"switch":  "0441",
			"skillID": uint16(10),
		}

		// Call the handler
		err := manager.handleSkillDelete(args)
		if err != nil {
			t.Errorf("handleSkillDelete() returned error: %v", err)
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
	})
}

// Test unhappy paths
func TestSkillDeleteUnhappy(t *testing.T) {
	// Create a skill delete manager
	manager := NewSkillDeleteManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0441",
			// Missing skillID
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillDelete(args)
		if err != nil {
			t.Errorf("handleSkillDelete() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillDelete(args)

		// Check that the skill ID is zero
		if skillID, ok := result["skillID"].(uint16); !ok || skillID != 0 {
			t.Errorf("Expected skill ID 0, got %v", result["skillID"])
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":  "0441",
			"skillID": "not a uint16", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillDelete(args)
		if err != nil {
			t.Errorf("handleSkillDelete() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillDelete(args)

		// Check that the skill ID is zero
		if skillID, ok := result["skillID"].(uint16); !ok || skillID != 0 {
			t.Errorf("Expected skill ID 0, got %v", result["skillID"])
		}
	})
}
