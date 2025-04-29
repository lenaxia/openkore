package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleElementalInfo(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("elemental.info", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Create a player for testing (to simulate the character)
	player := NewPlayer([]byte{5, 6, 7, 8})
	player.SetName("TestPlayer")
	handler.playersList.Add(player)

	// Set the player as the main character
	handler.SetMainCharacter(player)

	// Test cases
	testCases := []struct {
		name        string
		args        map[string]interface{}
		elementalID []byte
		hp          uint32
		maxHP       uint32
		sp          uint32
		maxSP       uint32
		level       uint16
		expectHook  bool
	}{
		{
			name: "New Elemental Info",
			args: map[string]interface{}{
				"ID":    []byte{1, 2, 3, 4},
				"hp":    uint32(500),
				"maxHP": uint32(1000),
				"sp":    uint32(200),
				"maxSP": uint32(300),
				"level": uint16(10),
				"KEYS":  []string{"hp", "maxHP", "sp", "maxSP", "level"},
			},
			elementalID: []byte{1, 2, 3, 4},
			hp:          500,
			maxHP:       1000,
			sp:          200,
			maxSP:       300,
			level:       10,
			expectHook:  true,
		},
		{
			name: "Update Existing Elemental Info",
			args: map[string]interface{}{
				"ID":    []byte{1, 2, 3, 4},
				"hp":    uint32(400),
				"maxHP": uint32(1000),
				"sp":    uint32(150),
				"maxSP": uint32(300),
				"level": uint16(10),
				"KEYS":  []string{"hp", "sp"},
			},
			elementalID: []byte{1, 2, 3, 4},
			hp:          400,
			maxHP:       1000,
			sp:          150,
			maxSP:       300,
			level:       10,
			expectHook:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleElementalInfo(tc.args)
			if err != nil {
				t.Errorf("HandleElementalInfo returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if id, ok := result["ID"].([]byte); !ok || string(id) != string(tc.elementalID) {
					t.Errorf("Expected ID %v, got %v", tc.elementalID, result["ID"])
				}

				// Check if the elemental was set on the main character
				if player.Elemental() == nil {
					t.Errorf("Expected elemental to be set on main character, but it's nil")
				} else {
					// Check elemental properties
					elemental := player.Elemental()

					if elemental.HP() != tc.hp {
						t.Errorf("Expected elemental HP %d, got %d", tc.hp, elemental.HP())
					}

					if elemental.MaxHP() != tc.maxHP {
						t.Errorf("Expected elemental MaxHP %d, got %d", tc.maxHP, elemental.MaxHP())
					}

					if elemental.SP() != tc.sp {
						t.Errorf("Expected elemental SP %d, got %d", tc.sp, elemental.SP())
					}

					if elemental.MaxSP() != tc.maxSP {
						t.Errorf("Expected elemental MaxSP %d, got %d", tc.maxSP, elemental.MaxSP())
					}

					if elemental.Level() != tc.level {
						t.Errorf("Expected elemental Level %d, got %d", tc.level, elemental.Level())
					}
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleElementalInfoUnhappy(t *testing.T) {
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
				"hp":    uint32(500),
				"maxHP": uint32(1000),
				"KEYS":  []string{"hp", "maxHP"},
			},
		},
		{
			name: "Missing KEYS",
			args: map[string]interface{}{
				"ID":    []byte{1, 2, 3, 4},
				"hp":    uint32(500),
				"maxHP": uint32(1000),
			},
		},
		{
			name: "Invalid ID type",
			args: map[string]interface{}{
				"ID":    "invalid",
				"hp":    uint32(500),
				"maxHP": uint32(1000),
				"KEYS":  []string{"hp", "maxHP"},
			},
		},
		{
			name: "Invalid KEYS type",
			args: map[string]interface{}{
				"ID":    []byte{1, 2, 3, 4},
				"hp":    uint32(500),
				"maxHP": uint32(1000),
				"KEYS":  "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleElementalInfo(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
