package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleStatsInfo(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.stats_info", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Test case
	args := map[string]interface{}{
		"points_free":      uint16(10),
		"str":              byte(30),
		"points_str":       byte(5),
		"agi":              byte(25),
		"points_agi":       byte(4),
		"vit":              byte(20),
		"points_vit":       byte(3),
		"int":              byte(35),
		"points_int":       byte(6),
		"dex":              byte(40),
		"points_dex":       byte(7),
		"luk":              byte(15),
		"points_luk":       byte(2),
		"attack":           uint16(100),
		"attack_bonus":     uint16(20),
		"attack_magic_min": uint16(80),
		"attack_magic_max": uint16(120),
		"def":              uint16(50),
		"def_bonus":        uint16(10),
		"def_magic":        uint16(40),
		"def_magic_bonus":  uint16(5),
		"hit":              uint16(200),
		"flee":             uint16(150),
		"flee_bonus":       uint16(30),
		"critical":         uint16(15),
	}

	// Call the handler
	err := handler.HandleStatsInfo(args)
	if err != nil {
		t.Errorf("HandleStatsInfo returned an error: %v", err)
	}

	// Get the result from the channel
	result := <-resultChan

	// Verify the result
	expectedFields := []string{
		"points_free", "str", "points_str", "agi", "points_agi", "vit", "points_vit",
		"int", "points_int", "dex", "points_dex", "luk", "points_luk", "attack",
		"attack_bonus", "attack_magic_min", "attack_magic_max", "def", "def_bonus",
		"def_magic", "def_magic_bonus", "hit", "flee", "flee_bonus", "critical",
	}

	for _, field := range expectedFields {
		if _, ok := result[field]; !ok {
			t.Errorf("Expected field %s in result, but it was not found", field)
		}
	}

	// Verify specific values
	if points, ok := result["points_free"].(uint16); !ok || points != 10 {
		t.Errorf("Expected points_free 10, got %v", result["points_free"])
	}

	if str, ok := result["str"].(byte); !ok || str != 30 {
		t.Errorf("Expected str 30, got %v", result["str"])
	}

	if attack, ok := result["attack"].(uint16); !ok || attack != 100 {
		t.Errorf("Expected attack 100, got %v", result["attack"])
	}

	if critical, ok := result["critical"].(uint16); !ok || critical != 15 {
		t.Errorf("Expected critical 15, got %v", result["critical"])
	}
}

// Test for unhappy paths
func TestHandleStatsInfoUnhappy(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Test cases for unhappy paths
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Missing points_free",
			args: map[string]interface{}{
				"str":        byte(30),
				"points_str": byte(5),
			},
		},
		{
			name: "Invalid points_free type",
			args: map[string]interface{}{
				"points_free": "invalid",
				"str":         byte(30),
				"points_str":  byte(5),
			},
		},
		{
			name: "Missing str",
			args: map[string]interface{}{
				"points_free": uint16(10),
				"points_str":  byte(5),
			},
		},
		{
			name: "Invalid str type",
			args: map[string]interface{}{
				"points_free": uint16(10),
				"str":         "invalid",
				"points_str":  byte(5),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleStatsInfo(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
