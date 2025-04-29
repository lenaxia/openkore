package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleStatInfo2(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.stat_info2", func(hookName string, arg interface{}, userData interface{}) {
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
		expectedVal  int32
		expectedVal2 int32
		expectedStat string
	}{
		{
			name: "Strength",
			args: map[string]interface{}{
				"type": int(VAR_STR),
				"val":  int32(30),
				"val2": int32(5),
			},
			expectedType: VAR_STR,
			expectedVal:  30,
			expectedVal2: 5,
			expectedStat: "str",
		},
		{
			name: "Agility",
			args: map[string]interface{}{
				"type": int(VAR_AGI),
				"val":  int32(25),
				"val2": int32(3),
			},
			expectedType: VAR_AGI,
			expectedVal:  25,
			expectedVal2: 3,
			expectedStat: "agi",
		},
		{
			name: "Vitality",
			args: map[string]interface{}{
				"type": int(VAR_VIT),
				"val":  int32(20),
				"val2": int32(2),
			},
			expectedType: VAR_VIT,
			expectedVal:  20,
			expectedVal2: 2,
			expectedStat: "vit",
		},
		{
			name: "Intelligence",
			args: map[string]interface{}{
				"type": int(VAR_INT),
				"val":  int32(35),
				"val2": int32(4),
			},
			expectedType: VAR_INT,
			expectedVal:  35,
			expectedVal2: 4,
			expectedStat: "int",
		},
		{
			name: "Dexterity",
			args: map[string]interface{}{
				"type": int(VAR_DEX),
				"val":  int32(40),
				"val2": int32(6),
			},
			expectedType: VAR_DEX,
			expectedVal:  40,
			expectedVal2: 6,
			expectedStat: "dex",
		},
		{
			name: "Luck",
			args: map[string]interface{}{
				"type": int(VAR_LUK),
				"val":  int32(15),
				"val2": int32(1),
			},
			expectedType: VAR_LUK,
			expectedVal:  15,
			expectedVal2: 1,
			expectedStat: "luk",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleStatInfo2(tc.args)
			if err != nil {
				t.Errorf("HandleStatInfo2 returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-resultChan

			// Verify the result
			if statType, ok := result["type"].(int); !ok || statType != tc.expectedType {
				t.Errorf("Expected type %d, got %v", tc.expectedType, result["type"])
			}

			if val, ok := result["val"].(int32); !ok || val != tc.expectedVal {
				t.Errorf("Expected val %d, got %v", tc.expectedVal, result["val"])
			}

			if val2, ok := result["val2"].(int32); !ok || val2 != tc.expectedVal2 {
				t.Errorf("Expected val2 %d, got %v", tc.expectedVal2, result["val2"])
			}

			if stat, ok := result["stat"].(string); !ok || stat != tc.expectedStat {
				t.Errorf("Expected stat %s, got %v", tc.expectedStat, result["stat"])
			}
		})
	}
}

// Test for unhappy paths
func TestHandleStatInfo2Unhappy(t *testing.T) {
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
				"val":  int32(30),
				"val2": int32(5),
			},
		},
		{
			name: "Missing val",
			args: map[string]interface{}{
				"type": int(VAR_STR),
				"val2": int32(5),
			},
		},
		{
			name: "Missing val2",
			args: map[string]interface{}{
				"type": int(VAR_STR),
				"val":  int32(30),
			},
		},
		{
			name: "Invalid type",
			args: map[string]interface{}{
				"type": "invalid",
				"val":  int32(30),
				"val2": int32(5),
			},
		},
		{
			name: "Invalid val",
			args: map[string]interface{}{
				"type": int(VAR_STR),
				"val":  "invalid",
				"val2": int32(5),
			},
		},
		{
			name: "Invalid val2",
			args: map[string]interface{}{
				"type": int(VAR_STR),
				"val":  int32(30),
				"val2": "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleStatInfo2(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
