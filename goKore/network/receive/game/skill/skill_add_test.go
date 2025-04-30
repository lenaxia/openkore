package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestSkillAdd(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the skill add manager
	manager := NewSkillAddManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for different packet types
	testCases := []struct {
		name          string
		packetSwitch  string
		args          map[string]interface{}
		expectedID    uint16
		expectedLevel uint16
		expectedName  string
	}{
		{
			name:         "Basic Skill Add (0111)",
			packetSwitch: "0111",
			args: map[string]interface{}{
				"skillID":    uint16(10),
				"target":     uint32(1),
				"lv":         uint16(5),
				"sp":         uint16(20),
				"range":      uint16(3),
				"upgradable": uint8(1),
				"lv2":        uint32(2),
			},
			expectedID:    10,
			expectedLevel: 5,
			expectedName:  "",
		},
		{
			name:         "Skill Add with Name (09FE)",
			packetSwitch: "09FE",
			args: map[string]interface{}{
				"skillID":    uint16(20),
				"target":     uint32(2),
				"lv":         uint16(3),
				"sp":         uint16(15),
				"range":      uint16(4),
				"upgradable": uint8(0),
				"lv2":        uint32(1),
				"name":       "SM_SWORD",
			},
			expectedID:    20,
			expectedLevel: 3,
			expectedName:  "SM_SWORD",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.skill_add", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Add the switch to the args
			tc.args["switch"] = tc.packetSwitch

			// Call the handler
			err := manager.handleSkillAdd(tc.args)
			if err != nil {
				t.Errorf("handleSkillAdd() returned error: %v", err)
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
				if skill.ID != tc.expectedID {
					t.Errorf("Expected skill ID %d, got %d", tc.expectedID, skill.ID)
				}
				if skill.Level != tc.expectedLevel {
					t.Errorf("Expected skill level %d, got %d", tc.expectedLevel, skill.Level)
				}
				if skill.Handle != tc.expectedName {
					t.Errorf("Expected skill name %q, got %q", tc.expectedName, skill.Handle)
				}
			} else {
				t.Error("Skill info not found in hook result")
			}
		})
	}
}

// Test unhappy paths
func TestSkillAddUnhappy(t *testing.T) {
	// Create a skill add manager
	manager := NewSkillAddManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "0111",
			// Missing skillID, target, etc.
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillAdd(args)
		if err != nil {
			t.Errorf("handleSkillAdd() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillAdd(args)

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
			"switch":     "0111",
			"skillID":    "not a uint16", // Wrong type
			"target":     "not a uint32", // Wrong type
			"lv":         "not a uint16", // Wrong type
			"sp":         "not a uint16", // Wrong type
			"range":      "not a uint16", // Wrong type
			"upgradable": "not a uint8",  // Wrong type
			"lv2":        "not a uint32", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillAdd(args)
		if err != nil {
			t.Errorf("handleSkillAdd() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillAdd(args)

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
