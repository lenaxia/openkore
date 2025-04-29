package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleMonsterTypechange(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("monster.typechange", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Create a monster for testing
	monster := NewMonster([]byte{1, 2, 3, 4})
	monster.SetName("Poring")
	monster.SetNameID(1002) // Poring ID
	handler.monstersList.Add(monster)

	// Test cases
	testCases := []struct {
		name       string
		args       map[string]interface{}
		oldName    string
		newNameID  uint16
		newName    string
		expectHook bool
	}{
		{
			name: "Monster Type Change",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": uint16(1001), // Drops ID
			},
			oldName:    "Poring",
			newNameID:  1001,
			newName:    "Drops", // This would be looked up from a monster name table
			expectHook: true,
		},
		{
			name: "Unknown Monster Type",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": uint16(9999), // Unknown monster ID
			},
			oldName:    "Drops", // From previous test
			newNameID:  9999,
			newName:    "Unknown #67305985", // Unknown monster name with nameID
			expectHook: true,
		},
		{
			name: "Unknown Monster",
			args: map[string]interface{}{
				"ID":   []byte{5, 6, 7, 8},
				"type": uint16(1001),
			},
			expectHook: false,
		},
	}

	// Create a simple monster name lookup table for testing
	monsterNames := map[uint16]string{
		1001: "Drops",
		1002: "Poring",
		1003: "Poporing",
	}
	handler.SetMonsterNames(monsterNames)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleMonsterTypechange(tc.args)
			if err != nil {
				t.Errorf("HandleMonsterTypechange returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if id, ok := result["ID"].([]byte); !ok || string(id) != string(tc.args["ID"].([]byte)) {
					t.Errorf("Expected ID %v, got %v", tc.args["ID"], result["ID"])
				}

				if monsterType, ok := result["type"].(uint16); !ok || monsterType != tc.args["type"].(uint16) {
					t.Errorf("Expected type %d, got %v", tc.args["type"].(uint16), result["type"])
				}

				if oldName, ok := result["oldName"].(string); !ok || oldName != tc.oldName {
					t.Errorf("Expected oldName %s, got %v", tc.oldName, result["oldName"])
				}

				if newName, ok := result["newName"].(string); !ok || newName != tc.newName {
					t.Errorf("Expected newName %s, got %v", tc.newName, result["newName"])
				}

				// Check if monster's name and nameID were updated
				if uint16(monster.NameID()) != tc.newNameID {
					t.Errorf("Expected monster nameID %d, got %d", tc.newNameID, monster.NameID())
				}

				if monster.Name() != tc.newName {
					t.Errorf("Expected monster name %s, got %s", tc.newName, monster.Name())
				}

				// Check if damage counters were reset
				if monster.DmgToParty() != 0 {
					t.Errorf("Expected dmgToParty to be reset to 0, got %d", monster.DmgToParty())
				}

				if monster.DmgFromParty() != 0 {
					t.Errorf("Expected dmgFromParty to be reset to 0, got %d", monster.DmgFromParty())
				}

				if monster.MissedToParty() != 0 {
					t.Errorf("Expected missedToParty to be reset to 0, got %d", monster.MissedToParty())
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleMonsterTypechangeUnhappy(t *testing.T) {
	// Create a handler for testing
	handler := NewHandler()

	// Test cases for unhappy paths
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Missing ID",
			args: map[string]interface{}{
				"type": uint16(1001),
			},
		},
		{
			name: "Missing type",
			args: map[string]interface{}{
				"ID": []byte{1, 2, 3, 4},
			},
		},
		{
			name: "Invalid ID type",
			args: map[string]interface{}{
				"ID":   "invalid",
				"type": uint16(1001),
			},
		},
		{
			name: "Invalid type type",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleMonsterTypechange(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
