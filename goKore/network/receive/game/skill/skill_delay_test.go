package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestSkillPostDelay(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the skill delay manager
	manager := NewSkillDelayManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test case for skill post delay
	t.Run("Basic Skill Post Delay", func(t *testing.T) {
		// Create a channel to receive hook events
		hookCalled := false
		var hookResult map[string]interface{}

		// Register a hook to capture the event
		hookManager.AddHook("character.skill_post_delay", func(hookName string, arg interface{}, userData interface{}) {
			hookCalled = true
			if result, ok := arg.(map[string]interface{}); ok {
				hookResult = result
			}
		}, nil)

		// Create packet data
		args := map[string]interface{}{
			"switch": "043D",
			"ID":     uint16(10),   // Skill ID
			"time":   uint32(5000), // 5 seconds
		}

		// Call the handler
		err := manager.handleSkillPostDelay(args)
		if err != nil {
			t.Errorf("handleSkillPostDelay() returned error: %v", err)
		}

		// Check that the hook was called
		if !hookCalled {
			t.Error("Hook was not called")
		}

		// Check the hook result
		if hookResult == nil {
			t.Fatal("Hook result is nil")
		}

		// Check the skill delay info
		if skillDelay, ok := hookResult["skillDelay"].(SkillDelayInfo); ok {
			if skillDelay.ID != 10 {
				t.Errorf("Expected skill ID 10, got %d", skillDelay.ID)
			}
			if skillDelay.RemainTime != 5000 {
				t.Errorf("Expected remain time 5000, got %d", skillDelay.RemainTime)
			}
		} else {
			t.Error("Skill delay info not found in hook result")
		}
	})
}

func TestSkillPostDelayList(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the skill delay manager
	manager := NewSkillDelayManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for different packet types
	testCases := []struct {
		name           string
		packetSwitch   string
		skillList      []byte
		expectedSkills int
		hasTotal       bool
	}{
		{
			name:         "Standard Skill Post Delay List (043E)",
			packetSwitch: "043E",
			skillList: []byte{
				10, 0, // Skill ID 10
				100, 0, 0, 0, // Remain time 100ms
				20, 0, // Skill ID 20
				200, 0, 0, 0, // Remain time 200ms
			},
			expectedSkills: 2,
			hasTotal:       false,
		},
		{
			name:         "Expanded Skill Post Delay List (0985)",
			packetSwitch: "0985",
			skillList: []byte{
				10, 0, // Skill ID 10
				150, 0, 0, 0, // Total time 150ms
				100, 0, 0, 0, // Remain time 100ms
				20, 0, // Skill ID 20
				250, 0, 0, 0, // Total time 250ms
				200, 0, 0, 0, // Remain time 200ms
			},
			expectedSkills: 2,
			hasTotal:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.skill_post_delaylist", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Create packet data
			args := map[string]interface{}{
				"switch":     tc.packetSwitch,
				"skill_list": tc.skillList,
			}

			// Call the handler
			err := manager.handleSkillPostDelayList(args)
			if err != nil {
				t.Errorf("handleSkillPostDelayList() returned error: %v", err)
			}

			// Check that the hook was called
			if !hookCalled {
				t.Error("Hook was not called")
			}

			// Check the hook result
			if hookResult == nil {
				t.Fatal("Hook result is nil")
			}

			// Check the skill delays
			if skillDelays, ok := hookResult["skillDelays"].([]SkillDelayInfo); ok {
				if len(skillDelays) != tc.expectedSkills {
					t.Errorf("Expected %d skill delays, got %d", tc.expectedSkills, len(skillDelays))
				}

				// Check the first skill delay
				if len(skillDelays) > 0 {
					if skillDelays[0].ID != 10 {
						t.Errorf("Expected skill ID 10, got %d", skillDelays[0].ID)
					}
					if skillDelays[0].RemainTime != 100 {
						t.Errorf("Expected remain time 100, got %d", skillDelays[0].RemainTime)
					}
					if tc.hasTotal && skillDelays[0].TotalTime != 150 {
						t.Errorf("Expected total time 150, got %d", skillDelays[0].TotalTime)
					}
				}

				// Check the second skill delay
				if len(skillDelays) > 1 {
					if skillDelays[1].ID != 20 {
						t.Errorf("Expected skill ID 20, got %d", skillDelays[1].ID)
					}
					if skillDelays[1].RemainTime != 200 {
						t.Errorf("Expected remain time 200, got %d", skillDelays[1].RemainTime)
					}
					if tc.hasTotal && skillDelays[1].TotalTime != 250 {
						t.Errorf("Expected total time 250, got %d", skillDelays[1].TotalTime)
					}
				}
			} else {
				t.Error("Skill delays not found in hook result")
			}
		})
	}
}

// Test unhappy paths
func TestSkillDelayUnhappy(t *testing.T) {
	// Create a skill delay manager
	manager := NewSkillDelayManager(nil, nil)

	// Test with missing fields for skill_post_delay
	t.Run("MissingFieldsPostDelay", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "043D",
			// Missing ID and time
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillPostDelay(args)
		if err != nil {
			t.Errorf("handleSkillPostDelay() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillPostDelay(args)

		// Check that the skill delay info has zero values
		if skillDelay, ok := result["skillDelay"].(SkillDelayInfo); ok {
			if skillDelay.ID != 0 {
				t.Errorf("Expected skill ID 0, got %d", skillDelay.ID)
			}
			if skillDelay.RemainTime != 0 {
				t.Errorf("Expected remain time 0, got %d", skillDelay.RemainTime)
			}
		} else {
			t.Error("Skill delay info not found in result")
		}
	})

	// Test with wrong field types for skill_post_delay
	t.Run("WrongFieldTypesPostDelay", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "043D",
			"ID":     "not a uint16", // Wrong type
			"time":   "not a uint32", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillPostDelay(args)
		if err != nil {
			t.Errorf("handleSkillPostDelay() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillPostDelay(args)

		// Check that the skill delay info has zero values
		if skillDelay, ok := result["skillDelay"].(SkillDelayInfo); ok {
			if skillDelay.ID != 0 {
				t.Errorf("Expected skill ID 0, got %d", skillDelay.ID)
			}
			if skillDelay.RemainTime != 0 {
				t.Errorf("Expected remain time 0, got %d", skillDelay.RemainTime)
			}
		} else {
			t.Error("Skill delay info not found in result")
		}
	})

	// Test with missing fields for skill_post_delaylist
	t.Run("MissingFieldsPostDelayList", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "043E",
			// Missing skill_list
		}

		// This should not return an error, but the result should be an empty list
		err := manager.handleSkillPostDelayList(args)
		if err != nil {
			t.Errorf("handleSkillPostDelayList() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillPostDelayList(args)

		// Check that the skill delays is an empty list
		if skillDelays, ok := result["skillDelays"].([]SkillDelayInfo); ok {
			if len(skillDelays) != 0 {
				t.Errorf("Expected 0 skill delays, got %d", len(skillDelays))
			}
		} else {
			t.Error("Skill delays not found in result")
		}
	})

	// Test with wrong field types for skill_post_delaylist
	t.Run("WrongFieldTypesPostDelayList", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":     "043E",
			"skill_list": "not a []byte", // Wrong type
		}

		// This should not return an error, but the result should be an empty list
		err := manager.handleSkillPostDelayList(args)
		if err != nil {
			t.Errorf("handleSkillPostDelayList() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillPostDelayList(args)

		// Check that the skill delays is an empty list
		if skillDelays, ok := result["skillDelays"].([]SkillDelayInfo); ok {
			if len(skillDelays) != 0 {
				t.Errorf("Expected 0 skill delays, got %d", len(skillDelays))
			}
		} else {
			t.Error("Skill delays not found in result")
		}
	})
}
