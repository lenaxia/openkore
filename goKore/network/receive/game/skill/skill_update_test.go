package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestSkillUpdate(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the skill update manager
	manager := NewSkillUpdateManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test case for skill update
	t.Run("Basic Skill Update", func(t *testing.T) {
		// Create a channel to receive hook events
		hookCalled := false
		var hookResult map[string]interface{}

		// Register a hook to capture the event
		hookManager.AddHook("character.skill_update", func(hookName string, arg interface{}, userData interface{}) {
			hookCalled = true
			if result, ok := arg.(map[string]interface{}); ok {
				hookResult = result
			}
		}, nil)

		// Create packet data
		args := map[string]interface{}{
			"switch":  "0110",
			"skillID": uint16(10),
			"lv":      uint16(5),
			"sp":      uint16(20),
			"range":   uint16(3),
			"up":      uint8(1),
			"lv2":     uint32(2),
		}

		// Call the handler
		err := manager.handleSkillUpdate(args)
		if err != nil {
			t.Errorf("handleSkillUpdate() returned error: %v", err)
		}

		// Check that the hook was called
		if !hookCalled {
			t.Error("Hook was not called")
		}

		// Check the hook result
		if hookResult == nil {
			t.Fatal("Hook result is nil")
		}

		// Check the owner type
		if ownerType, ok := hookResult["ownerType"].(SkillOwnerType); !ok || ownerType != OwnerChar {
			t.Errorf("Expected owner type %v, got %v", OwnerChar, hookResult["ownerType"])
		}

		// Check the skill info
		if skill, ok := hookResult["skill"].(SkillInfo); ok {
			if skill.ID != 10 {
				t.Errorf("Expected skill ID 10, got %d", skill.ID)
			}
			if skill.Level != 5 {
				t.Errorf("Expected skill level 5, got %d", skill.Level)
			}
			if skill.SP != 20 {
				t.Errorf("Expected skill SP 20, got %d", skill.SP)
			}
			if skill.Range != 3 {
				t.Errorf("Expected skill range 3, got %d", skill.Range)
			}
			if skill.Up != 1 {
				t.Errorf("Expected skill up 1, got %d", skill.Up)
			}
			if skill.Level2 != 2 {
				t.Errorf("Expected skill level2 2, got %d", skill.Level2)
			}
		} else {
			t.Error("Skill info not found in hook result")
		}
	})
}

// Test unhappy paths
func TestSkillUpdateUnhappy(t *testing.T) {
	// Create a skill update manager
	manager := NewSkillUpdateManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0110",
			// Missing skillID, lv, etc.
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillUpdate(args)
		if err != nil {
			t.Errorf("handleSkillUpdate() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillUpdate(args)

		// Check that the skill info has zero values
		if skill, ok := result["skill"].(SkillInfo); ok {
			if skill.ID != 0 {
				t.Errorf("Expected skill ID 0, got %d", skill.ID)
			}
			if skill.Level != 0 {
				t.Errorf("Expected skill level 0, got %d", skill.Level)
			}
		} else {
			t.Error("Skill info not found in result")
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":  "0110",
			"skillID": "not a uint16", // Wrong type
			"lv":      "not a uint16", // Wrong type
			"sp":      "not a uint16", // Wrong type
			"range":   "not a uint16", // Wrong type
			"up":      "not a uint8",  // Wrong type
			"lv2":     "not a uint32", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillUpdate(args)
		if err != nil {
			t.Errorf("handleSkillUpdate() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillUpdate(args)

		// Check that the skill info has zero values
		if skill, ok := result["skill"].(SkillInfo); ok {
			if skill.ID != 0 {
				t.Errorf("Expected skill ID 0, got %d", skill.ID)
			}
			if skill.Level != 0 {
				t.Errorf("Expected skill level 0, got %d", skill.Level)
			}
		} else {
			t.Error("Skill info not found in result")
		}
	})
}
