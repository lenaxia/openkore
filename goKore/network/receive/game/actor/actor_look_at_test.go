package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleActorLookAt(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.look_at", func(hookName string, arg interface{}, userData interface{}) {
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
		headDir    byte
		bodyDir    byte
		expectHook bool
	}{
		{
			name: "Player Look Direction",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"head": byte(2),
				"body": byte(5),
			},
			headDir:    2,
			bodyDir:    5,
			expectHook: true,
		},
		{
			name: "Unknown Actor",
			args: map[string]interface{}{
				"ID":   []byte{5, 6, 7, 8},
				"head": byte(1),
				"body": byte(3),
			},
			expectHook: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleActorLookAt(tc.args)
			if err != nil {
				t.Errorf("HandleActorLookAt returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if id, ok := result["ID"].([]byte); !ok || string(id) != string(tc.args["ID"].([]byte)) {
					t.Errorf("Expected ID %v, got %v", tc.args["ID"], result["ID"])
				}

				if head, ok := result["head"].(byte); !ok || head != tc.headDir {
					t.Errorf("Expected head %d, got %v", tc.headDir, result["head"])
				}

				if body, ok := result["body"].(byte); !ok || body != tc.bodyDir {
					t.Errorf("Expected body %d, got %v", tc.bodyDir, result["body"])
				}

				// Verify the player's look directions were updated
				if player.HeadDirection() != tc.headDir {
					t.Errorf("Expected player head direction %d, got %d", tc.headDir, player.HeadDirection())
				}

				if player.BodyDirection() != tc.bodyDir {
					t.Errorf("Expected player body direction %d, got %d", tc.bodyDir, player.BodyDirection())
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleActorLookAtUnhappy(t *testing.T) {
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
			name: "Missing ID",
			args: map[string]interface{}{
				"head": byte(2),
				"body": byte(5),
			},
		},
		{
			name: "Missing head",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"body": byte(5),
			},
		},
		{
			name: "Missing body",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"head": byte(2),
			},
		},
		{
			name: "Invalid ID type",
			args: map[string]interface{}{
				"ID":   "invalid",
				"head": byte(2),
				"body": byte(5),
			},
		},
		{
			name: "Invalid head type",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"head": "invalid",
				"body": byte(5),
			},
		},
		{
			name: "Invalid body type",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"head": byte(2),
				"body": "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleActorLookAt(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
