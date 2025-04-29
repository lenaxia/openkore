package field

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestNavigateTo(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("field.navigate_to", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewFieldManager(nil, hookManager)

	// Test cases
	testCases := []struct {
		name           string
		args           map[string]interface{}
		expectedStatus string
		hasMobID       bool
	}{
		{
			name: "Navigate to Location",
			args: map[string]interface{}{
				"type":   byte(0),
				"flag":   byte(0),
				"hide":   byte(0),
				"map":    []byte("prontera"),
				"x":      uint16(150),
				"y":      uint16(200),
				"mob_id": uint16(0),
			},
			expectedStatus: "Server asked us to navigate to prontera (150,200)",
			hasMobID:       false,
		},
		{
			name: "Navigate to Monster",
			args: map[string]interface{}{
				"type":   byte(0),
				"flag":   byte(0),
				"hide":   byte(0),
				"map":    []byte("prontera"),
				"x":      uint16(0),
				"y":      uint16(0),
				"mob_id": uint16(1002),
			},
			expectedStatus: "Server asked us to navigate to prontera map and look for monster with ID 1002",
			hasMobID:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := manager.handleNavigateTo(tc.args)
			if err != nil {
				t.Errorf("handleNavigateTo returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-resultChan

			// Verify the result
			if result["status"] != tc.expectedStatus {
				t.Errorf("Expected status %q, got %q", tc.expectedStatus, result["status"])
			}

			// Verify map name
			if mapName, ok := result["map"].(string); !ok || mapName != "prontera" {
				t.Errorf("Expected map 'prontera', got %v", result["map"])
			}

			// Verify coordinates
			if x, ok := result["x"].(uint16); !ok || x != tc.args["x"] {
				t.Errorf("Expected x %v, got %v", tc.args["x"], result["x"])
			}
			if y, ok := result["y"].(uint16); !ok || y != tc.args["y"] {
				t.Errorf("Expected y %v, got %v", tc.args["y"], result["y"])
			}

			// Verify mob_id based on test case
			if tc.hasMobID {
				if mobID, ok := result["mob_id"].(uint16); !ok || mobID != tc.args["mob_id"] {
					t.Errorf("Expected mob_id %v, got %v", tc.args["mob_id"], result["mob_id"])
				}
			}
		})
	}
}
