package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleUsersOnline(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("core.users_online", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a users online manager for testing
	manager := NewUsersOnlineManager(hookManager)

	// Test cases
	testCases := []struct {
		name        string
		args        map[string]interface{}
		expectHook  bool
		expectedMsg string
	}{
		{
			name: "Valid Users Online",
			args: map[string]interface{}{
				"users": uint32(1234),
			},
			expectHook:  true,
			expectedMsg: "There are currently 1234 users online",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := manager.HandleUsersOnline(tc.args)
			if err != nil {
				t.Errorf("HandleUsersOnline returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if users, ok := result["users"].(uint32); !ok || users != tc.args["users"].(uint32) {
					t.Errorf("Expected users %v, got %v", tc.args["users"], result["users"])
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
func TestHandleUsersOnlineUnhappy(t *testing.T) {
	// Create a users online manager for testing
	manager := NewUsersOnlineManager(nil)

	// Test cases for unhappy paths
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Missing Users",
			args: map[string]interface{}{},
		},
		{
			name: "Invalid Users Type",
			args: map[string]interface{}{
				"users": "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := manager.HandleUsersOnline(tc.args)
			if err != nil {
				t.Errorf("Expected no error for unhappy path, but got: %v", err)
			}
		})
	}
}
