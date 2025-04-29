package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleHpSpChanged(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.hp_sp_changed", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Test cases
	testCases := []struct {
		name           string
		args           map[string]interface{}
		expectedType   int
		expectedAmount int32
		expectedStat   string
	}{
		{
			name: "HP Increase",
			args: map[string]interface{}{
				"type":   int(VAR_HP), // 5
				"amount": int32(100),
			},
			expectedType:   VAR_HP,
			expectedAmount: 100,
			expectedStat:   "hp",
		},
		{
			name: "SP Increase",
			args: map[string]interface{}{
				"type":   int(VAR_SP), // 7
				"amount": int32(50),
			},
			expectedType:   VAR_SP,
			expectedAmount: 50,
			expectedStat:   "sp",
		},
		{
			name: "HP Decrease",
			args: map[string]interface{}{
				"type":   int(VAR_HP), // 5
				"amount": int32(-30),
			},
			expectedType:   VAR_HP,
			expectedAmount: -30,
			expectedStat:   "hp",
		},
		{
			name: "SP Decrease",
			args: map[string]interface{}{
				"type":   int(VAR_SP), // 7
				"amount": int32(-20),
			},
			expectedType:   VAR_SP,
			expectedAmount: -20,
			expectedStat:   "sp",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleHpSpChanged(tc.args)
			if err != nil {
				t.Errorf("HandleHpSpChanged returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-resultChan

			// Verify the result
			if statType, ok := result["type"].(int); !ok || statType != tc.expectedType {
				t.Errorf("Expected type %d, got %v", tc.expectedType, result["type"])
			}

			if amount, ok := result["amount"].(int32); !ok || amount != tc.expectedAmount {
				t.Errorf("Expected amount %d, got %v", tc.expectedAmount, result["amount"])
			}

			if stat, ok := result["stat"].(string); !ok || stat != tc.expectedStat {
				t.Errorf("Expected stat %s, got %v", tc.expectedStat, result["stat"])
			}
		})
	}
}

// Test for unhappy paths
func TestHandleHpSpChangedUnhappy(t *testing.T) {
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
				"amount": int32(100),
			},
		},
		{
			name: "Missing amount",
			args: map[string]interface{}{
				"type": int(VAR_HP),
			},
		},
		{
			name: "Invalid type",
			args: map[string]interface{}{
				"type":   "invalid",
				"amount": int32(100),
			},
		},
		{
			name: "Invalid amount",
			args: map[string]interface{}{
				"type":   int(VAR_HP),
				"amount": "invalid",
			},
		},
		{
			name: "Unknown type",
			args: map[string]interface{}{
				"type":   int(999), // Unknown type
				"amount": int32(100),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleHpSpChanged(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
