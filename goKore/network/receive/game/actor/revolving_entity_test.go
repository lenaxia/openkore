package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleRevolvingEntity(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.revolving_entity", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Create a player for testing
	player := NewPlayer([]byte{1, 2, 3, 4})
	player.SetName("TestPlayer")
	handler.playersList.Add(player)

	// Test cases
	testCases := []struct {
		name          string
		args          map[string]interface{}
		switchID      string
		entityType    string
		entityNum     uint16
		entityElement string
		expectHook    bool
	}{
		{
			name: "Monk Spirits",
			args: map[string]interface{}{
				"sourceID": []byte{1, 2, 3, 4},
				"entity":   uint16(3),
				"switch":   "01D0",
			},
			switchID:      "01D0",
			entityType:    "spirit",
			entityNum:     3,
			entityElement: "",
			expectHook:    true,
		},
		{
			name: "Gunslinger Coins",
			args: map[string]interface{}{
				"sourceID": []byte{1, 2, 3, 4},
				"entity":   uint16(5),
				"switch":   "01E1",
			},
			switchID:      "01E1",
			entityType:    "coin",
			entityNum:     5,
			entityElement: "",
			expectHook:    true,
		},
		{
			name: "Ninja Amulet",
			args: map[string]interface{}{
				"sourceID": []byte{1, 2, 3, 4},
				"entity":   uint16(4),
				"type":     uint16(1), // Water element
				"switch":   "08CF",
			},
			switchID:      "08CF",
			entityType:    "amulet",
			entityNum:     4,
			entityElement: "Water",
			expectHook:    true,
		},
		{
			name: "Soul Energy",
			args: map[string]interface{}{
				"sourceID": []byte{1, 2, 3, 4},
				"entity":   uint16(2),
				"switch":   "0B73",
			},
			switchID:      "0B73",
			entityType:    "soul energy",
			entityNum:     2,
			entityElement: "",
			expectHook:    true,
		},
		{
			name: "Unknown Entity Type",
			args: map[string]interface{}{
				"sourceID": []byte{1, 2, 3, 4},
				"entity":   uint16(1),
				"switch":   "FFFF", // Unknown switch
			},
			switchID:      "FFFF",
			entityType:    "entity unknown",
			entityNum:     1,
			entityElement: "",
			expectHook:    true,
		},
		{
			name: "Unknown Actor",
			args: map[string]interface{}{
				"sourceID": []byte{5, 6, 7, 8},
				"entity":   uint16(3),
				"switch":   "01D0",
			},
			switchID:      "01D0",
			entityType:    "spirit",
			entityNum:     3,
			entityElement: "",
			expectHook:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleRevolvingEntity(tc.args)
			if err != nil {
				t.Errorf("HandleRevolvingEntity returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if sourceID, ok := result["sourceID"].([]byte); !ok || string(sourceID) != string(tc.args["sourceID"].([]byte)) {
					t.Errorf("Expected sourceID %v, got %v", tc.args["sourceID"], result["sourceID"])
				}

				if entityType, ok := result["entityType"].(string); !ok || entityType != tc.entityType {
					t.Errorf("Expected entityType %s, got %v", tc.entityType, result["entityType"])
				}

				if entityNum, ok := result["entityNum"].(uint16); !ok || entityNum != tc.entityNum {
					t.Errorf("Expected entityNum %d, got %v", tc.entityNum, result["entityNum"])
				}

				if tc.entityElement != "" {
					if entityElement, ok := result["entityElement"].(string); !ok || entityElement != tc.entityElement {
						t.Errorf("Expected entityElement %s, got %v", tc.entityElement, result["entityElement"])
					}
				}

				// Check if player's spirits were updated
				if tc.args["sourceID"].([]byte)[0] == 1 && tc.args["sourceID"].([]byte)[1] == 2 &&
					tc.args["sourceID"].([]byte)[2] == 3 && tc.args["sourceID"].([]byte)[3] == 4 {
					if uint16(player.Spirits()) != tc.entityNum {
						t.Errorf("Expected player spirits %d, got %d", tc.entityNum, player.Spirits())
					}

					if tc.entityElement != "" && player.AmuletType() != tc.entityElement {
						t.Errorf("Expected player amulet type %s, got %s", tc.entityElement, player.AmuletType())
					}

					if player.SpiritsType() != tc.entityType {
						t.Errorf("Expected player spirits type %s, got %s", tc.entityType, player.SpiritsType())
					}
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleRevolvingEntityUnhappy(t *testing.T) {
	// Create a handler for testing
	handler := NewHandler()

	// Test cases for unhappy paths
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Missing sourceID",
			args: map[string]interface{}{
				"entity": uint16(3),
				"switch": "01D0",
			},
		},
		{
			name: "Missing entity",
			args: map[string]interface{}{
				"sourceID": []byte{1, 2, 3, 4},
				"switch":   "01D0",
			},
		},
		{
			name: "Missing switch",
			args: map[string]interface{}{
				"sourceID": []byte{1, 2, 3, 4},
				"entity":   uint16(3),
			},
		},
		{
			name: "Invalid sourceID type",
			args: map[string]interface{}{
				"sourceID": "invalid",
				"entity":   uint16(3),
				"switch":   "01D0",
			},
		},
		{
			name: "Invalid entity type",
			args: map[string]interface{}{
				"sourceID": []byte{1, 2, 3, 4},
				"entity":   "invalid",
				"switch":   "01D0",
			},
		},
		{
			name: "Invalid switch type",
			args: map[string]interface{}{
				"sourceID": []byte{1, 2, 3, 4},
				"entity":   uint16(3),
				"switch":   123,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleRevolvingEntity(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
