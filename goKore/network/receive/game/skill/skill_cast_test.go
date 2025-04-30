package skill

import (
	"testing"
	"time"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestSkillCast(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the skill cast manager
	manager := NewSkillCastManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for different packet types
	testCases := []struct {
		name               string
		packetSwitch       string
		args               map[string]interface{}
		expectedSourceID   uint32
		expectedTargetID   uint32
		expectedSkillID    uint16
		expectedX          uint16
		expectedY          uint16
		expectedWait       uint32
		isLocationTargeted bool
	}{
		{
			name:         "Basic Skill Cast (013E)",
			packetSwitch: "013E",
			args: map[string]interface{}{
				"sourceID": uint32(100),
				"targetID": uint32(200),
				"x":        uint16(0),
				"y":        uint16(0),
				"skillID":  uint16(10),
				"type":     uint8(1),
				"wait":     uint32(1000),
			},
			expectedSourceID:   100,
			expectedTargetID:   200,
			expectedSkillID:    10,
			expectedX:          0,
			expectedY:          0,
			expectedWait:       1000,
			isLocationTargeted: false,
		},
		{
			name:         "Location-Targeted Skill Cast (013E)",
			packetSwitch: "013E",
			args: map[string]interface{}{
				"sourceID": uint32(100),
				"targetID": uint32(0),
				"x":        uint16(150),
				"y":        uint16(200),
				"skillID":  uint16(20),
				"type":     uint8(1),
				"wait":     uint32(2000),
			},
			expectedSourceID:   100,
			expectedTargetID:   0,
			expectedSkillID:    20,
			expectedX:          150,
			expectedY:          200,
			expectedWait:       2000,
			isLocationTargeted: true,
		},
		{
			name:         "Expanded Skill Cast (07FB)",
			packetSwitch: "07FB",
			args: map[string]interface{}{
				"sourceID": uint32(100),
				"targetID": uint32(200),
				"x":        uint16(0),
				"y":        uint16(0),
				"skillID":  uint16(30),
				"type":     uint8(1),
				"wait":     uint32(1500),
				"unknown":  uint32(0),
			},
			expectedSourceID:   100,
			expectedTargetID:   200,
			expectedSkillID:    30,
			expectedX:          0,
			expectedY:          0,
			expectedWait:       1500,
			isLocationTargeted: false,
		},
		{
			name:         "No Damage Skill Cast (0A1C)",
			packetSwitch: "0A1C",
			args: map[string]interface{}{
				"sourceID": uint32(100),
				"targetID": uint32(200),
				"skillID":  uint16(40),
				"unknown":  uint16(0),
				"type":     uint8(1),
				"wait":     uint32(500),
				"unknown2": uint16(0),
			},
			expectedSourceID:   100,
			expectedTargetID:   200,
			expectedSkillID:    40,
			expectedX:          0,
			expectedY:          0,
			expectedWait:       500,
			isLocationTargeted: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.skill_cast", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Add the switch to the args
			tc.args["switch"] = tc.packetSwitch

			// Call the appropriate handler based on the packet switch
			var err error
			if tc.packetSwitch == "0A1C" {
				err = manager.handleSkillCastNoDamage(tc.args)
			} else {
				err = manager.handleSkillCast(tc.args)
			}

			if err != nil {
				t.Errorf("handleSkillCast() returned error: %v", err)
			}

			// Check that the hook was called
			if !hookCalled {
				t.Error("Hook was not called")
			}

			// Check the hook result
			if hookResult == nil {
				t.Fatal("Hook result is nil")
			}

			// Check the isLocationTargeted flag
			if isLocationTargeted, ok := hookResult["isLocationTargeted"].(bool); !ok || isLocationTargeted != tc.isLocationTargeted {
				t.Errorf("Expected isLocationTargeted %v, got %v", tc.isLocationTargeted, isLocationTargeted)
			}

			// Check the waitTimeSeconds
			if waitTimeSeconds, ok := hookResult["waitTimeSeconds"].(float64); !ok || waitTimeSeconds != float64(tc.expectedWait)/1000.0 {
				t.Errorf("Expected waitTimeSeconds %v, got %v", float64(tc.expectedWait)/1000.0, waitTimeSeconds)
			}

			// Check the cast info
			if castInfo, ok := hookResult["castInfo"].(CastInfo); ok {
				if castInfo.SourceID != tc.expectedSourceID {
					t.Errorf("Expected source ID %d, got %d", tc.expectedSourceID, castInfo.SourceID)
				}
				if castInfo.TargetID != tc.expectedTargetID {
					t.Errorf("Expected target ID %d, got %d", tc.expectedTargetID, castInfo.TargetID)
				}
				if castInfo.SkillID != tc.expectedSkillID {
					t.Errorf("Expected skill ID %d, got %d", tc.expectedSkillID, castInfo.SkillID)
				}
				if castInfo.X != tc.expectedX {
					t.Errorf("Expected X %d, got %d", tc.expectedX, castInfo.X)
				}
				if castInfo.Y != tc.expectedY {
					t.Errorf("Expected Y %d, got %d", tc.expectedY, castInfo.Y)
				}
				if castInfo.Wait != tc.expectedWait {
					t.Errorf("Expected wait %d, got %d", tc.expectedWait, castInfo.Wait)
				}
				// Check that StartTime is set to a reasonable value (within the last second)
				if time.Since(castInfo.StartTime) > time.Second {
					t.Errorf("StartTime is too old: %v", castInfo.StartTime)
				}
			} else {
				t.Error("Cast info not found in hook result")
			}
		})
	}
}

// Test unhappy paths
func TestSkillCastUnhappy(t *testing.T) {
	// Create a skill cast manager
	manager := NewSkillCastManager(nil, nil)

	// Test with missing fields
	t.Run("MissingFields", func(t *testing.T) {
		args := map[string]interface{}{
			"switch": "013E",
			// Missing sourceID, targetID, etc.
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillCast(args)
		if err != nil {
			t.Errorf("handleSkillCast() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillCast(args)

		// Check that the cast info has zero values
		if castInfo, ok := result["castInfo"].(CastInfo); ok {
			if castInfo.SourceID != 0 {
				t.Errorf("Expected source ID 0, got %d", castInfo.SourceID)
			}
			if castInfo.TargetID != 0 {
				t.Errorf("Expected target ID 0, got %d", castInfo.TargetID)
			}
			if castInfo.SkillID != 0 {
				t.Errorf("Expected skill ID 0, got %d", castInfo.SkillID)
			}
		} else {
			t.Error("Cast info not found in result")
		}
	})

	// Test with wrong field types
	t.Run("WrongFieldTypes", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":   "013E",
			"sourceID": "not a uint32", // Wrong type
			"targetID": "not a uint32", // Wrong type
			"x":        "not a uint16", // Wrong type
			"y":        "not a uint16", // Wrong type
			"skillID":  "not a uint16", // Wrong type
			"type":     "not a uint8",  // Wrong type
			"wait":     "not a uint32", // Wrong type
		}

		// This should not return an error, but the fields should be zero values
		err := manager.handleSkillCast(args)
		if err != nil {
			t.Errorf("handleSkillCast() returned error: %v", err)
		}

		// Process the args directly to check the result
		result := manager.processSkillCast(args)

		// Check that the cast info has zero values
		if castInfo, ok := result["castInfo"].(CastInfo); ok {
			if castInfo.SourceID != 0 {
				t.Errorf("Expected source ID 0, got %d", castInfo.SourceID)
			}
			if castInfo.TargetID != 0 {
				t.Errorf("Expected target ID 0, got %d", castInfo.TargetID)
			}
			if castInfo.SkillID != 0 {
				t.Errorf("Expected skill ID 0, got %d", castInfo.SkillID)
			}
		} else {
			t.Error("Cast info not found in result")
		}
	})
}
