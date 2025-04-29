package misc

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestEACKey(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("misc.eac_key", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Test cases
	testCases := []struct {
		name                   string
		ignoreAntiCheatWarning bool
		expectedIgnored        bool
		expectedQuit           bool
	}{
		{
			name:                   "Ignore Anti-Cheat Warning",
			ignoreAntiCheatWarning: true,
			expectedIgnored:        true,
			expectedQuit:           false,
		},
		{
			name:                   "Do Not Ignore Anti-Cheat Warning",
			ignoreAntiCheatWarning: false,
			expectedIgnored:        false,
			expectedQuit:           true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a config for testing
			config := &AntiCheatConfig{
				IgnoreAntiCheatWarning: tc.ignoreAntiCheatWarning,
			}

			// Create a manager for testing
			manager := NewAntiCheatManager(nil, hookManager, config)

			// Call the handler
			err := manager.handleEACKey(map[string]interface{}{})
			if err != nil {
				t.Errorf("handleEACKey returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-resultChan

			// Verify the result
			if result["ignored"] != tc.expectedIgnored {
				t.Errorf("Expected ignored %v, got %v", tc.expectedIgnored, result["ignored"])
			}
			if result["quit"] != tc.expectedQuit {
				t.Errorf("Expected quit %v, got %v", tc.expectedQuit, result["quit"])
			}
		})
	}
}
