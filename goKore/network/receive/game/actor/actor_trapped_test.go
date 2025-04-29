package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleActorTrapped(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.trapped", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Create a player for testing
	player := NewPlayer([]byte{1, 2, 3, 4})
	player.SetName("TestPlayer")
	handler.playersList.Add(player)

	// Test cases
	testCases := []struct {
		name       string
		args       map[string]interface{}
		expectHook bool
	}{
		{
			name: "Player Trapped",
			args: map[string]interface{}{
				"ID": []byte{1, 2, 3, 4},
			},
			expectHook: true,
		},
		{
			name: "Unknown Actor",
			args: map[string]interface{}{
				"ID": []byte{5, 6, 7, 8},
			},
			expectHook: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleActorTrapped(tc.args)
			if err != nil {
				t.Errorf("HandleActorTrapped returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if id, ok := result["ID"].([]byte); !ok || string(id) != string(tc.args["ID"].([]byte)) {
					t.Errorf("Expected ID %v, got %v", tc.args["ID"], result["ID"])
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleActorTrappedUnhappy(t *testing.T) {
	// Create a handler for testing
	handler := NewHandler()

	// Test cases for unhappy paths
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Missing ID",
			args: map[string]interface{}{},
		},
		{
			name: "Invalid ID type",
			args: map[string]interface{}{
				"ID": "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleActorTrapped(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
