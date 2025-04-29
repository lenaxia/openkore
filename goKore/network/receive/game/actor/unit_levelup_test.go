package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleUnitLevelup(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create channels to capture hook calls
	baseLevelChan := make(chan map[string]interface{}, 1)
	jobLevelChan := make(chan map[string]interface{}, 1)
	refineFailChan := make(chan map[string]interface{}, 1)
	refineSuccessChan := make(chan map[string]interface{}, 1)
	gameOverChan := make(chan map[string]interface{}, 1)
	potionSuccessChan := make(chan map[string]interface{}, 1)
	potionFailChan := make(chan map[string]interface{}, 1)

	// Register hooks to capture the results
	hookManager.AddHook("actor.base_level", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		baseLevelChan <- result
	}, nil)

	hookManager.AddHook("actor.job_level", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		jobLevelChan <- result
	}, nil)

	hookManager.AddHook("actor.refine_fail", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		refineFailChan <- result
	}, nil)

	hookManager.AddHook("actor.refine_success", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		refineSuccessChan <- result
	}, nil)

	hookManager.AddHook("actor.game_over", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		gameOverChan <- result
	}, nil)

	hookManager.AddHook("actor.potion_success", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		potionSuccessChan <- result
	}, nil)

	hookManager.AddHook("actor.potion_fail", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		potionFailChan <- result
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
		name       string
		args       map[string]interface{}
		effectType byte
		hookName   string
		hookChan   chan map[string]interface{}
	}{
		{
			name: "Base Level Up",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": byte(0), // LEVELUP_EFFECT
			},
			effectType: 0,
			hookName:   "actor.base_level",
			hookChan:   baseLevelChan,
		},
		{
			name: "Base Level Up (Super Novice)",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": byte(7), // LEVELUP_EFFECT2
			},
			effectType: 7,
			hookName:   "actor.base_level",
			hookChan:   baseLevelChan,
		},
		{
			name: "Base Level Up (Taekwon)",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": byte(9), // LEVELUP_EFFECT3
			},
			effectType: 9,
			hookName:   "actor.base_level",
			hookChan:   baseLevelChan,
		},
		{
			name: "Job Level Up",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": byte(1), // JOBLEVELUP_EFFECT
			},
			effectType: 1,
			hookName:   "actor.job_level",
			hookChan:   jobLevelChan,
		},
		{
			name: "Job Level Up (Super Novice)",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": byte(8), // JOBLEVELUP_EFFECT2
			},
			effectType: 8,
			hookName:   "actor.job_level",
			hookChan:   jobLevelChan,
		},
		{
			name: "Refine Failure",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": byte(2), // REFINING_FAIL_EFFECT
			},
			effectType: 2,
			hookName:   "actor.refine_fail",
			hookChan:   refineFailChan,
		},
		{
			name: "Refine Success",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": byte(3), // REFINING_SUCCESS_EFFECT
			},
			effectType: 3,
			hookName:   "actor.refine_success",
			hookChan:   refineSuccessChan,
		},
		{
			name: "Game Over",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": byte(4), // GAME_OVER_EFFECT
			},
			effectType: 4,
			hookName:   "actor.game_over",
			hookChan:   gameOverChan,
		},
		{
			name: "Potion Success",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": byte(5), // MAKEITEM_AM_SUCCESS_EFFECT
			},
			effectType: 5,
			hookName:   "actor.potion_success",
			hookChan:   potionSuccessChan,
		},
		{
			name: "Potion Failure",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": byte(6), // MAKEITEM_AM_FAIL_EFFECT
			},
			effectType: 6,
			hookName:   "actor.potion_fail",
			hookChan:   potionFailChan,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleUnitLevelup(tc.args)
			if err != nil {
				t.Errorf("HandleUnitLevelup returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-tc.hookChan

			// Verify the result
			if id, ok := result["ID"].([]byte); !ok || string(id) != string(tc.args["ID"].([]byte)) {
				t.Errorf("Expected ID %v, got %v", tc.args["ID"], result["ID"])
			}

			if effectType, ok := result["type"].(byte); !ok || effectType != tc.effectType {
				t.Errorf("Expected effect type %d, got %v", tc.effectType, result["type"])
			}

			if actor, ok := result["actor"].(Actor); !ok || actor.Name() != player.Name() {
				t.Errorf("Expected actor %s, got %v", player.Name(), result["actor"])
			}
		})
	}
}

// Test for unknown effect type
func TestHandleUnitLevelupUnknownEffect(t *testing.T) {
	// Create a handler for testing
	handler := NewHandler()

	// Create a player for testing
	player := NewPlayer([]byte{1, 2, 3, 4})
	player.SetName("TestPlayer")
	handler.playersList.Add(player)

	// Test case for unknown effect type
	args := map[string]interface{}{
		"ID":   []byte{1, 2, 3, 4},
		"type": byte(10), // Unknown effect type
	}

	// Call the handler - it should not panic
	err := handler.HandleUnitLevelup(args)
	if err != nil {
		t.Errorf("HandleUnitLevelup returned an error: %v", err)
	}
}

// Test for unhappy paths
func TestHandleUnitLevelupUnhappy(t *testing.T) {
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
				"type": byte(0),
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
				"type": byte(0),
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
			err := handler.HandleUnitLevelup(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
