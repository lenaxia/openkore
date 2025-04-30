package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestSageAutospell(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the sage autospell manager
	manager := NewSageAutospellManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for different packet types
	testCases := []struct {
		name             string
		packetSwitch     string
		args             map[string]interface{}
		expectedWhy      uint16
		expectedSkillIDs []uint16
		isAutoShadow     bool
	}{
		{
			name:         "Sage's Hindsight (01CD)",
			packetSwitch: "01CD",
			args: map[string]interface{}{
				"why": uint16(1),
				"autospell_list": []byte{
					10, 0, 0, 0, // Skill ID 10
					20, 0, 0, 0, // Skill ID 20
					30, 0, 0, 0, // Skill ID 30
				},
			},
			expectedWhy:      1,
			expectedSkillIDs: []uint16{10, 20, 30},
			isAutoShadow:     false,
		},
		{
			name:         "Shadow Chaser's Auto Shadow Spell (0442)",
			packetSwitch: "0442",
			args: map[string]interface{}{
				"why": uint16(2),
				"autoshadowspell_list": []byte{
					10, 0, // Skill ID 10
					20, 0, // Skill ID 20
					30, 0, // Skill ID 30
				},
			},
			expectedWhy:      2,
			expectedSkillIDs: []uint16{10, 20, 30},
			isAutoShadow:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.sage_autospell", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Add the switch to the args
			tc.args["switch"] = tc.packetSwitch

			// Call the handler
			err := manager.handleSageAutospell(tc.args)
			if err != nil {
				t.Errorf("handleSageAutospell() returned error: %v", err)
			}

			// Check that the hook was called
			if !hookCalled {
				t.Error("Hook was not called")
			}

			// Check the hook result
			if hookResult == nil {
				t.Fatal("Hook result is nil")
			}

			// Check the why value
			if why, ok := hookResult["why"].(uint16); !ok || why != tc.expectedWhy {
				t.Errorf("Expected why %d, got %v", tc.expectedWhy, why)
			}

			// Check the skill IDs
			if skillIDs, ok := hookResult["skillIDs"].([]uint16); ok {
				if len(skillIDs) != len(tc.expectedSkillIDs) {
					t.Errorf("Expected %d skill IDs, got %d", len(tc.expectedSkillIDs), len(skillIDs))
				} else {
					for i, skillID := range skillIDs {
						if skillID != tc.expectedSkillIDs[i] {
							t.Errorf("Expected skill ID %d at index %d, got %d", tc.expectedSkillIDs[i], i, skillID)
						}
					}
				}
			} else {
				t.Error("skillIDs not found in hook result or not a []uint16")
			}

			// Check the isAutoShadow flag
			if isAutoShadow, ok := hookResult["isAutoShadow"].(bool); !ok || isAutoShadow != tc.isAutoShadow {
				t.Errorf("Expected isAutoShadow %v, got %v", tc.isAutoShadow, isAutoShadow)
			}

			// Check that the message is not empty
			if message, ok := hookResult["message"].(string); !ok || message == "" {
				t.Error("message not found in hook result or empty")
			}

			// Check that the skillInfoList is not empty
			if skillInfoList, ok := hookResult["skillInfoList"].([]map[string]interface{}); !ok {
				t.Error("skillInfoList not found in hook result or not a []map[string]interface{}")
			} else if len(skillInfoList) != len(tc.expectedSkillIDs) {
				t.Errorf("Expected %d skill infos, got %d", len(tc.expectedSkillIDs), len(skillInfoList))
			}
		})
	}
}

// Test the reconstruct function
func TestReconstructSageAutospell(t *testing.T) {
	// Create a sage autospell manager
	manager := NewSageAutospellManager(nil, nil)

	// Test case
	skillIDs := []uint16{10, 20, 30}

	// Call the reconstruct function
	autoshadowspellList, autospellList := manager.reconstructSageAutospell(skillIDs, true)

	// Check the autoshadowspell_list
	expectedAutoshadowspellList := []byte{
		10, 0, // Skill ID 10
		20, 0, // Skill ID 20
		30, 0, // Skill ID 30
	}
	if len(autoshadowspellList) != len(expectedAutoshadowspellList) {
		t.Errorf("Expected autoshadowspell_list length %d, got %d", len(expectedAutoshadowspellList), len(autoshadowspellList))
	} else {
		for i, b := range autoshadowspellList {
			if b != expectedAutoshadowspellList[i] {
				t.Errorf("Expected autoshadowspell_list byte %d at index %d, got %d", expectedAutoshadowspellList[i], i, b)
			}
		}
	}

	// Check the autospell_list
	expectedAutospellList := []byte{
		10, 0, 0, 0, // Skill ID 10
		20, 0, 0, 0, // Skill ID 20
		30, 0, 0, 0, // Skill ID 30
	}
	if len(autospellList) != len(expectedAutospellList) {
		t.Errorf("Expected autospell_list length %d, got %d", len(expectedAutospellList), len(autospellList))
	} else {
		for i, b := range autospellList {
			if b != expectedAutospellList[i] {
				t.Errorf("Expected autospell_list byte %d at index %d, got %d", expectedAutospellList[i], i, b)
			}
		}
	}
}

// Test unhappy paths
func TestSageAutospellUnhappy(t *testing.T) {
	// Create a sage autospell manager
	manager := NewSageAutospellManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "01CD",
			// Missing why and autospell_list
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSageAutospell(args)
		if err != nil {
			t.Errorf("handleSageAutospell() returned error: %v", err)
		}

		// Parse the args
		parsedArgs := manager.parseSageAutospell(args)

		// Check that the why is zero
		if why, ok := parsedArgs["why"].(uint16); !ok || why != 0 {
			t.Errorf("Expected why 0, got %v", why)
		}

		// Check that the skillIDs is empty
		if skillIDs, ok := parsedArgs["skillIDs"].([]uint16); !ok || len(skillIDs) != 0 {
			t.Errorf("Expected empty skillIDs, got %v", skillIDs)
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":         "01CD",
			"why":            "not a uint16", // Wrong type
			"autospell_list": "not a []byte", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSageAutospell(args)
		if err != nil {
			t.Errorf("handleSageAutospell() returned error: %v", err)
		}

		// Parse the args
		parsedArgs := manager.parseSageAutospell(args)

		// Check that the why is zero
		if why, ok := parsedArgs["why"].(uint16); !ok || why != 0 {
			t.Errorf("Expected why 0, got %v", why)
		}

		// Check that the skillIDs is empty
		if skillIDs, ok := parsedArgs["skillIDs"].([]uint16); !ok || len(skillIDs) != 0 {
			t.Errorf("Expected empty skillIDs, got %v", skillIDs)
		}
	})
}
