package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleStatInfo(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.stat_info", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Test cases
	testCases := []struct {
		name          string
		args          map[string]interface{}
		expectedType  int
		expectedValue int32
		expectedActor string
	}{
		{
			name: "Character Walk Speed",
			args: map[string]interface{}{
				"switch": "00B0", // Character stat info
				"type":   int(0), // VAR_SPEED
				"val":    int32(150),
			},
			expectedType:  0,
			expectedValue: 150,
			expectedActor: "character",
		},
		{
			name: "Character HP",
			args: map[string]interface{}{
				"switch": "00B0", // Character stat info
				"type":   int(5), // VAR_HP
				"val":    int32(1000),
			},
			expectedType:  5,
			expectedValue: 1000,
			expectedActor: "character",
		},
		{
			name: "Character SP",
			args: map[string]interface{}{
				"switch": "00B0", // Character stat info
				"type":   int(7), // VAR_SP
				"val":    int32(500),
			},
			expectedType:  7,
			expectedValue: 500,
			expectedActor: "character",
		},
		{
			name: "Character Level",
			args: map[string]interface{}{
				"switch": "00B0",  // Character stat info
				"type":   int(11), // VAR_CLEVEL
				"val":    int32(99),
			},
			expectedType:  11,
			expectedValue: 99,
			expectedActor: "character",
		},
		{
			name: "Character Money",
			args: map[string]interface{}{
				"switch": "00B0",  // Character stat info
				"type":   int(20), // VAR_MONEY
				"val":    int32(10000),
			},
			expectedType:  20,
			expectedValue: 10000,
			expectedActor: "character",
		},
		{
			name: "Other Player Stat",
			args: map[string]interface{}{
				"switch": "01AB",             // Other player stat info
				"ID":     []byte{1, 2, 3, 4}, // Player ID
				"type":   int(5),             // VAR_HP
				"val":    int32(2000),
			},
			expectedType:  5,
			expectedValue: 2000,
			expectedActor: "other",
		},
		{
			name: "Homunculus Stat",
			args: map[string]interface{}{
				"switch": "07DB", // Homunculus stat info
				"type":   int(5), // VAR_HP
				"val":    int32(3000),
			},
			expectedType:  5,
			expectedValue: 3000,
			expectedActor: "homunculus",
		},
		{
			name: "Mercenary Stat",
			args: map[string]interface{}{
				"switch": "02A2", // Mercenary stat info
				"type":   int(5), // VAR_HP
				"val":    int32(4000),
			},
			expectedType:  5,
			expectedValue: 4000,
			expectedActor: "mercenary",
		},
		{
			name: "Elemental Stat",
			args: map[string]interface{}{
				"switch": "081E", // Elemental stat info
				"type":   int(5), // VAR_HP
				"val":    int32(5000),
			},
			expectedType:  5,
			expectedValue: 5000,
			expectedActor: "elemental",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleStatInfo(tc.args)
			if err != nil {
				t.Errorf("HandleStatInfo returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-resultChan

			// Verify the result
			if statType, ok := result["type"].(int); !ok || statType != tc.expectedType {
				t.Errorf("Expected type %d, got %v", tc.expectedType, result["type"])
			}

			if value, ok := result["value"].(int32); !ok || value != tc.expectedValue {
				t.Errorf("Expected value %d, got %v", tc.expectedValue, result["value"])
			}

			if actor, ok := result["actor"].(string); !ok || actor != tc.expectedActor {
				t.Errorf("Expected actor %s, got %v", tc.expectedActor, result["actor"])
			}
		})
	}
}

// Test for unhappy paths
func TestHandleStatInfoUnhappy(t *testing.T) {
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
			name: "Missing switch",
			args: map[string]interface{}{
				"type": int(0),
				"val":  int32(150),
			},
		},
		{
			name: "Missing type",
			args: map[string]interface{}{
				"switch": "00B0",
				"val":    int32(150),
			},
		},
		{
			name: "Missing val",
			args: map[string]interface{}{
				"switch": "00B0",
				"type":   int(0),
			},
		},
		{
			name: "Unknown actor type",
			args: map[string]interface{}{
				"switch": "FFFF", // Unknown switch
				"type":   int(0),
				"val":    int32(150),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleStatInfo(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
