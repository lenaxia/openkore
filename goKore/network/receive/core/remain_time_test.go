package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleRemainTimeInfo(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("core.remain_time_info", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a remain time manager for testing
	manager := NewRemainTimeManager(hookManager)

	// Test cases
	testCases := []struct {
		name       string
		args       map[string]interface{}
		expectHook bool
	}{
		{
			name: "Valid Remain Time Info",
			args: map[string]interface{}{
				"result":          uint16(1),
				"expiration_date": uint32(1609459200), // 2021-01-01 00:00:00 UTC
				"remain_time":     uint32(86400),      // 1 day in seconds
			},
			expectHook: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := manager.HandleRemainTimeInfo(tc.args)
			if err != nil {
				t.Errorf("HandleRemainTimeInfo returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if result["result"] != tc.args["result"] {
					t.Errorf("Expected result %v, got %v", tc.args["result"], result["result"])
				}

				if result["expiration_date"] != tc.args["expiration_date"] {
					t.Errorf("Expected expiration_date %v, got %v", tc.args["expiration_date"], result["expiration_date"])
				}

				if result["remain_time"] != tc.args["remain_time"] {
					t.Errorf("Expected remain_time %v, got %v", tc.args["remain_time"], result["remain_time"])
				}

				// Check that the message was formatted correctly
				message, ok := result["message"].(string)
				if !ok {
					t.Errorf("Expected message to be a string")
				} else {
					expectedMessage := "Remain Time - Result: 1 - Expiration Date: 1609459200 - Time: 86400"
					if message != expectedMessage {
						t.Errorf("Expected message %q, got %q", expectedMessage, message)
					}
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleRemainTimeInfoUnhappy(t *testing.T) {
	// Create a remain time manager for testing
	manager := NewRemainTimeManager(nil)

	// Test cases for unhappy paths
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Missing Result",
			args: map[string]interface{}{
				"expiration_date": uint32(1609459200),
				"remain_time":     uint32(86400),
			},
		},
		{
			name: "Invalid Result Type",
			args: map[string]interface{}{
				"result":          "invalid",
				"expiration_date": uint32(1609459200),
				"remain_time":     uint32(86400),
			},
		},
		{
			name: "Missing Expiration Date",
			args: map[string]interface{}{
				"result":      uint16(1),
				"remain_time": uint32(86400),
			},
		},
		{
			name: "Invalid Expiration Date Type",
			args: map[string]interface{}{
				"result":          uint16(1),
				"expiration_date": "invalid",
				"remain_time":     uint32(86400),
			},
		},
		{
			name: "Missing Remain Time",
			args: map[string]interface{}{
				"result":          uint16(1),
				"expiration_date": uint32(1609459200),
			},
		},
		{
			name: "Invalid Remain Time Type",
			args: map[string]interface{}{
				"result":          uint16(1),
				"expiration_date": uint32(1609459200),
				"remain_time":     "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := manager.HandleRemainTimeInfo(tc.args)
			if err != nil {
				t.Errorf("Expected no error for unhappy path, but got: %v", err)
			}
		})
	}
}
