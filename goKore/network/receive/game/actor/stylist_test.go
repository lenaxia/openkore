package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleStylistRes(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("game.actor.stylist_res", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a stylist manager for testing
	manager := NewStylistManager(hookManager)

	// Test cases
	testCases := []struct {
		name        string
		args        map[string]interface{}
		expectHook  bool
		expectedMsg string
		success     bool
	}{
		{
			name: "Success Result",
			args: map[string]interface{}{
				"result": byte(1),
			},
			expectHook:  true,
			expectedMsg: "[Stylist UI] Success.",
			success:     true,
		},
		{
			name: "Fail Result",
			args: map[string]interface{}{
				"result": byte(0),
			},
			expectHook:  true,
			expectedMsg: "[Stylist UI] Fail.",
			success:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := manager.HandleStylistRes(tc.args)
			if err != nil {
				t.Errorf("HandleStylistRes returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if result["result"] != tc.args["result"] {
					t.Errorf("Expected result %v, got %v", tc.args["result"], result["result"])
				}

				// Check that the success flag was set correctly
				if success, ok := result["success"].(bool); !ok || success != tc.success {
					t.Errorf("Expected success %v, got %v", tc.success, result["success"])
				}

				// Check that the message was formatted correctly
				message, ok := result["message"].(string)
				if !ok {
					t.Errorf("Expected message to be a string")
				} else if message != tc.expectedMsg {
					t.Errorf("Expected message %q, got %q", tc.expectedMsg, message)
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleStylistResUnhappy(t *testing.T) {
	// Create a stylist manager for testing
	manager := NewStylistManager(nil)

	// Test cases for unhappy paths
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Missing Result",
			args: map[string]interface{}{},
		},
		{
			name: "Invalid Result Type",
			args: map[string]interface{}{
				"result": "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := manager.HandleStylistRes(tc.args)
			if err != nil {
				t.Errorf("Expected no error for unhappy path, but got: %v", err)
			}
		})
	}
}
