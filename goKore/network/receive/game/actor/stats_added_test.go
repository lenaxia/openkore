package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleStatsAdded(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.stats_added", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Test cases
	testCases := []struct {
		name         string
		args         map[string]interface{}
		expectedType int
		expectedVal  byte
		expectedStat string
	}{
		{
			name: "Strength Added",
			args: map[string]interface{}{
				"type":   int(VAR_STR),
				"val":    byte(30),
				"result": byte(1), // success
			},
			expectedType: VAR_STR,
			expectedVal:  30,
			expectedStat: "str",
		},
		{
			name: "Agility Added",
			args: map[string]interface{}{
				"type":   int(VAR_AGI),
				"val":    byte(25),
				"result": byte(1), // success
			},
			expectedType: VAR_AGI,
			expectedVal:  25,
			expectedStat: "agi",
		},
		{
			name: "Vitality Added",
			args: map[string]interface{}{
				"type":   int(VAR_VIT),
				"val":    byte(20),
				"result": byte(1), // success
			},
			expectedType: VAR_VIT,
			expectedVal:  20,
			expectedStat: "vit",
		},
		{
			name: "Intelligence Added",
			args: map[string]interface{}{
				"type":   int(VAR_INT),
				"val":    byte(35),
				"result": byte(1), // success
			},
			expectedType: VAR_INT,
			expectedVal:  35,
			expectedStat: "int",
		},
		{
			name: "Dexterity Added",
			args: map[string]interface{}{
				"type":   int(VAR_DEX),
				"val":    byte(40),
				"result": byte(1), // success
			},
			expectedType: VAR_DEX,
			expectedVal:  40,
			expectedStat: "dex",
		},
		{
			name: "Luck Added",
			args: map[string]interface{}{
				"type":   int(VAR_LUK),
				"val":    byte(15),
				"result": byte(1), // success
			},
			expectedType: VAR_LUK,
			expectedVal:  15,
			expectedStat: "luk",
		},
		{
			name: "Power Added",
			args: map[string]interface{}{
				"type":   int(VAR_SP_POW),
				"val":    byte(10),
				"result": byte(1), // success
			},
			expectedType: VAR_SP_POW,
			expectedVal:  10,
			expectedStat: "pow",
		},
		{
			name: "Stamina Added",
			args: map[string]interface{}{
				"type":   int(VAR_SP_STA),
				"val":    byte(12),
				"result": byte(1), // success
			},
			expectedType: VAR_SP_STA,
			expectedVal:  12,
			expectedStat: "sta",
		},
		{
			name: "Wisdom Added",
			args: map[string]interface{}{
				"type":   int(VAR_SP_WIS),
				"val":    byte(14),
				"result": byte(1), // success
			},
			expectedType: VAR_SP_WIS,
			expectedVal:  14,
			expectedStat: "wis",
		},
		{
			name: "Spell Added",
			args: map[string]interface{}{
				"type":   int(VAR_SP_SPL),
				"val":    byte(16),
				"result": byte(1), // success
			},
			expectedType: VAR_SP_SPL,
			expectedVal:  16,
			expectedStat: "spl",
		},
		{
			name: "Concentration Added",
			args: map[string]interface{}{
				"type":   int(VAR_SP_CON),
				"val":    byte(18),
				"result": byte(1), // success
			},
			expectedType: VAR_SP_CON,
			expectedVal:  18,
			expectedStat: "con",
		},
		{
			name: "Creative Added",
			args: map[string]interface{}{
				"type":   int(VAR_SP_CRT),
				"val":    byte(20),
				"result": byte(1), // success
			},
			expectedType: VAR_SP_CRT,
			expectedVal:  20,
			expectedStat: "crt",
		},
		{
			name: "Not Enough Points",
			args: map[string]interface{}{
				"type":   int(VAR_STR),
				"val":    byte(207), // special value for not enough points
				"result": byte(0),   // failure
			},
			expectedType: VAR_STR,
			expectedVal:  207,
			expectedStat: "error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleStatsAdded(tc.args)
			if err != nil {
				t.Errorf("HandleStatsAdded returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-resultChan

			// Verify the result
			if statType, ok := result["type"].(int); !ok || statType != tc.expectedType {
				t.Errorf("Expected type %d, got %v", tc.expectedType, result["type"])
			}

			if val, ok := result["val"].(byte); !ok || val != tc.expectedVal {
				t.Errorf("Expected val %d, got %v", tc.expectedVal, result["val"])
			}

			if stat, ok := result["stat"].(string); !ok || stat != tc.expectedStat {
				t.Errorf("Expected stat %s, got %v", tc.expectedStat, result["stat"])
			}
		})
	}
}

// Test for unhappy paths
func TestHandleStatsAddedUnhappy(t *testing.T) {
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
			name: "Missing type",
			args: map[string]interface{}{
				"val":    byte(30),
				"result": byte(1),
			},
		},
		{
			name: "Missing val",
			args: map[string]interface{}{
				"type":   int(VAR_STR),
				"result": byte(1),
			},
		},
		{
			name: "Invalid type",
			args: map[string]interface{}{
				"type":   "invalid",
				"val":    byte(30),
				"result": byte(1),
			},
		},
		{
			name: "Invalid val",
			args: map[string]interface{}{
				"type":   int(VAR_STR),
				"val":    "invalid",
				"result": byte(1),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleStatsAdded(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
