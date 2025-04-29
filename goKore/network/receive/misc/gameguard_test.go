package misc

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestGameGuardRequest(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("misc.gameguard_request", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	config := &GameGuardConfig{}
	manager := NewGameGuardManager(nil, hookManager, config)

	// Test cases
	testCases := []struct {
		name        string
		args        map[string]interface{}
		gameGuard   int
		netVersion  int
		shouldQuery bool
	}{
		{
			name: "GameGuard Enabled - Should Query",
			args: map[string]interface{}{
				"RAW_MSG":      []byte{0x01, 0x02, 0x03, 0x04},
				"RAW_MSG_SIZE": 4,
			},
			gameGuard:   2,
			netVersion:  1,
			shouldQuery: true,
		},
		{
			name: "GameGuard Disabled - Should Not Query",
			args: map[string]interface{}{
				"RAW_MSG":      []byte{0x01, 0x02, 0x03, 0x04},
				"RAW_MSG_SIZE": 4,
			},
			gameGuard:   0,
			netVersion:  1,
			shouldQuery: false,
		},
		{
			name: "GameGuard Not 2 - Should Not Query",
			args: map[string]interface{}{
				"RAW_MSG":      []byte{0x01, 0x02, 0x03, 0x04},
				"RAW_MSG_SIZE": 4,
			},
			gameGuard:   1,
			netVersion:  1,
			shouldQuery: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set the gameGuard and netVersion values
			manager.config.GameGuard = tc.gameGuard
			manager.config.NetVersion = tc.netVersion

			// Call the handler
			err := manager.handleGameGuardRequest(tc.args)
			if err != nil {
				t.Errorf("handleGameGuardRequest returned an error: %v", err)
			}

			// Check if we should expect a result
			if tc.shouldQuery {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if result["raw_msg"] == nil {
					t.Errorf("Expected raw_msg to be set")
				}

				// Check that the raw message is correct
				rawMsg, ok := result["raw_msg"].([]byte)
				if !ok {
					t.Errorf("Expected raw_msg to be a byte slice")
				} else if len(rawMsg) != tc.args["RAW_MSG_SIZE"].(int) {
					t.Errorf("Expected raw_msg length to be %d, got %d", tc.args["RAW_MSG_SIZE"].(int), len(rawMsg))
				}
			} else {
				// Make sure no result was sent
				select {
				case result := <-resultChan:
					t.Errorf("Unexpected result: %v", result)
				default:
					// No result, which is expected
				}
			}
		})
	}
}

func TestGameGuardGrant(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("misc.gameguard_grant", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	config := &GameGuardConfig{}
	manager := NewGameGuardManager(nil, hookManager, config)

	// Test cases
	testCases := []struct {
		name          string
		args          map[string]interface{}
		expectedState string
	}{
		{
			name: "Server Denied Login",
			args: map[string]interface{}{
				"server": uint8(0),
			},
			expectedState: "denied",
		},
		{
			name: "Server Granted Login to Account Server",
			args: map[string]interface{}{
				"server": uint8(1),
			},
			expectedState: "account_server",
		},
		{
			name: "Server Granted Login to Char/Map Server",
			args: map[string]interface{}{
				"server": uint8(2),
			},
			expectedState: "char_map_server",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := manager.handleGameGuardGrant(tc.args)
			if err != nil {
				t.Errorf("handleGameGuardGrant returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-resultChan

			// Verify the result
			if result["state"] != tc.expectedState {
				t.Errorf("Expected state %q, got %q", tc.expectedState, result["state"])
			}
			if result["server"] != tc.args["server"] {
				t.Errorf("Expected server %v, got %v", tc.args["server"], result["server"])
			}
		})
	}
}
