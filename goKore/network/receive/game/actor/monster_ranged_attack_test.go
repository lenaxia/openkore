package actor

import (
	"testing"
	"time"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleMonsterRangedAttack(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("monster.ranged_attack", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Create a monster for testing
	monster := NewMonster([]byte{1, 2, 3, 4})
	monster.SetName("Poring")
	handler.monstersList.Add(monster)

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
		sourceX     int
		sourceY     int
		targetX     int
		targetY     int
		attackRange int
		expectHook  bool
	}{
		{
			name: "Monster Ranged Attack",
			args: map[string]interface{}{
				"ID":      []byte{1, 2, 3, 4},
				"sourceX": 100,
				"sourceY": 200,
				"targetX": 150,
				"targetY": 250,
				"range":   5,
			},
			sourceX:     100,
			sourceY:     200,
			targetX:     150,
			targetY:     250,
			attackRange: 5,
			expectHook:  true,
		},
		{
			name: "Unknown Monster",
			args: map[string]interface{}{
				"ID":      []byte{9, 8, 7, 6},
				"sourceX": 100,
				"sourceY": 200,
				"targetX": 150,
				"targetY": 250,
				"range":   5,
			},
			sourceX:     100,
			sourceY:     200,
			targetX:     150,
			targetY:     250,
			attackRange: 5,
			expectHook:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Record the current time for comparison
			beforeTime := time.Now()

			// Call the handler
			err := handler.HandleMonsterRangedAttack(tc.args)
			if err != nil {
				t.Errorf("HandleMonsterRangedAttack returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if id, ok := result["ID"].([]byte); !ok || string(id) != string(tc.args["ID"].([]byte)) {
					t.Errorf("Expected ID %v, got %v", tc.args["ID"], result["ID"])
				}

				if sourceX, ok := result["sourceX"].(int); !ok || sourceX != tc.sourceX {
					t.Errorf("Expected sourceX %d, got %v", tc.sourceX, result["sourceX"])
				}

				if sourceY, ok := result["sourceY"].(int); !ok || sourceY != tc.sourceY {
					t.Errorf("Expected sourceY %d, got %v", tc.sourceY, result["sourceY"])
				}

				if targetX, ok := result["targetX"].(int); !ok || targetX != tc.targetX {
					t.Errorf("Expected targetX %d, got %v", tc.targetX, result["targetX"])
				}

				if targetY, ok := result["targetY"].(int); !ok || targetY != tc.targetY {
					t.Errorf("Expected targetY %d, got %v", tc.targetY, result["targetY"])
				}

				if rangeVal, ok := result["range"].(int); !ok || rangeVal != tc.attackRange {
					t.Errorf("Expected range %d, got %v", tc.attackRange, result["range"])
				}

				// Check if monster's movetoattack_pos was updated (if monster exists)
				if tc.name == "Monster Ranged Attack" {
					if monster.MoveToAttackPos() == nil {
						t.Errorf("Expected monster movetoattack_pos to be set, but it's nil")
					} else {
						if monster.MoveToAttackPos().X != tc.sourceX {
							t.Errorf("Expected monster movetoattack_pos.X %d, got %d", tc.sourceX, monster.MoveToAttackPos().X)
						}
						if monster.MoveToAttackPos().Y != tc.sourceY {
							t.Errorf("Expected monster movetoattack_pos.Y %d, got %d", tc.sourceY, monster.MoveToAttackPos().Y)
						}
					}

					// Check if monster's movetoattack_time was updated
					if monster.MoveToAttackTime().Before(beforeTime) {
						t.Errorf("Expected monster movetoattack_time to be updated")
					}
				}

				// Check if character's movetoattack_pos was updated
				if player.MoveToAttackPos() == nil {
					t.Errorf("Expected character movetoattack_pos to be set, but it's nil")
				} else {
					if player.MoveToAttackPos().X != tc.targetX {
						t.Errorf("Expected character movetoattack_pos.X %d, got %d", tc.targetX, player.MoveToAttackPos().X)
					}
					if player.MoveToAttackPos().Y != tc.targetY {
						t.Errorf("Expected character movetoattack_pos.Y %d, got %d", tc.targetY, player.MoveToAttackPos().Y)
					}
				}

				// Check if character's movetoattack_time was updated
				if player.MoveToAttackTime().Before(beforeTime) {
					t.Errorf("Expected character movetoattack_time to be updated")
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleMonsterRangedAttackUnhappy(t *testing.T) {
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
				"sourceX": 100,
				"sourceY": 200,
				"targetX": 150,
				"targetY": 250,
				"range":   5,
			},
		},
		{
			name: "Missing sourceX",
			args: map[string]interface{}{
				"ID":      []byte{1, 2, 3, 4},
				"sourceY": 200,
				"targetX": 150,
				"targetY": 250,
				"range":   5,
			},
		},
		{
			name: "Missing sourceY",
			args: map[string]interface{}{
				"ID":      []byte{1, 2, 3, 4},
				"sourceX": 100,
				"targetX": 150,
				"targetY": 250,
				"range":   5,
			},
		},
		{
			name: "Missing targetX",
			args: map[string]interface{}{
				"ID":      []byte{1, 2, 3, 4},
				"sourceX": 100,
				"sourceY": 200,
				"targetY": 250,
				"range":   5,
			},
		},
		{
			name: "Missing targetY",
			args: map[string]interface{}{
				"ID":      []byte{1, 2, 3, 4},
				"sourceX": 100,
				"sourceY": 200,
				"targetX": 150,
				"range":   5,
			},
		},
		{
			name: "Missing range",
			args: map[string]interface{}{
				"ID":      []byte{1, 2, 3, 4},
				"sourceX": 100,
				"sourceY": 200,
				"targetX": 150,
				"targetY": 250,
			},
		},
		{
			name: "Invalid ID type",
			args: map[string]interface{}{
				"ID":      "invalid",
				"sourceX": 100,
				"sourceY": 200,
				"targetX": 150,
				"targetY": 250,
				"range":   5,
			},
		},
		{
			name: "Invalid sourceX type",
			args: map[string]interface{}{
				"ID":      []byte{1, 2, 3, 4},
				"sourceX": "invalid",
				"sourceY": 200,
				"targetX": 150,
				"targetY": 250,
				"range":   5,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleMonsterRangedAttack(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
