package field

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestBossMapInfo(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("field.boss_map_info", func(hookName string, arg interface{}, userData interface{}) {
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
		expectedFlag   byte
	}{
		{
			name: "No Boss Found",
			args: map[string]interface{}{
				"flag": byte(0),
				"name": []byte("Baphomet"),
			},
			expectedStatus: "You cannot find any trace of a Boss Monster in this area.",
			expectedFlag:   0,
		},
		{
			name: "Boss Location Known",
			args: map[string]interface{}{
				"flag": byte(1),
				"name": []byte("Baphomet"),
				"x":    uint16(150),
				"y":    uint16(200),
			},
			expectedStatus: "MVP Boss Baphomet is now on location: (150, 200)",
			expectedFlag:   1,
		},
		{
			name: "Boss Detected on Map",
			args: map[string]interface{}{
				"flag": byte(2),
				"name": []byte("Baphomet"),
			},
			expectedStatus: "MVP Boss Baphomet has been detected on this map!",
			expectedFlag:   2,
		},
		{
			name: "Boss Dead with Respawn Time",
			args: map[string]interface{}{
				"flag":    byte(3),
				"name":    []byte("Baphomet"),
				"hours":   byte(2),
				"minutes": byte(30),
			},
			expectedStatus: "MVP Boss Baphomet is dead, but will spawn again in 2 hour(s) and 30 minutes(s).",
			expectedFlag:   3,
		},
		{
			name: "Unknown Flag",
			args: map[string]interface{}{
				"flag": byte(99),
				"name": []byte("Baphomet"),
			},
			expectedStatus: "Unknown boss_map_info result (flag: 99)",
			expectedFlag:   99,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := manager.handleBossMapInfo(tc.args)
			if err != nil {
				t.Errorf("handleBossMapInfo returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-resultChan

			// Verify the result
			if result["status"] != tc.expectedStatus {
				t.Errorf("Expected status %q, got %q", tc.expectedStatus, result["status"])
			}
			if result["flag"] != tc.args["flag"] {
				t.Errorf("Expected flag %v, got %v", tc.args["flag"], result["flag"])
			}

			// Verify boss name
			if bossName, ok := result["bossName"].(string); !ok || bossName != "Baphomet" {
				t.Errorf("Expected bossName 'Baphomet', got %v", result["bossName"])
			}

			// Verify additional fields based on flag
			switch tc.expectedFlag {
			case 1:
				// Check coordinates for flag 1
				if result["x"] != tc.args["x"] {
					t.Errorf("Expected x %v, got %v", tc.args["x"], result["x"])
				}
				if result["y"] != tc.args["y"] {
					t.Errorf("Expected y %v, got %v", tc.args["y"], result["y"])
				}
			case 3:
				// Check respawn time for flag 3
				if result["hours"] != tc.args["hours"] {
					t.Errorf("Expected hours %v, got %v", tc.args["hours"], result["hours"])
				}
				if result["minutes"] != tc.args["minutes"] {
					t.Errorf("Expected minutes %v, got %v", tc.args["minutes"], result["minutes"])
				}
			}
		})
	}
}
