package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleHotkeys(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("core.hotkeys", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a hotkey manager for testing
	manager := NewHotkeyManager(hookManager)

	// Test cases
	testCases := []struct {
		name          string
		args          map[string]interface{}
		expectedCount int
		expectHook    bool
	}{
		{
			name: "Valid Hotkeys",
			args: map[string]interface{}{
				"hotkeys": []byte{
					// Hotkey 1: Skill, ID 1234, Level 5
					1, 210, 4, 0, 0, 5, 0,
					// Hotkey 2: Item, ID 5678, Level 0
					0, 46, 22, 0, 0, 0, 0,
				},
			},
			expectedCount: 2,
			expectHook:    true,
		},
		{
			name: "Empty Hotkeys",
			args: map[string]interface{}{
				"hotkeys": []byte{},
			},
			expectedCount: 0,
			expectHook:    true,
		},
		{
			name: "Partial Hotkey Data",
			args: map[string]interface{}{
				"hotkeys": []byte{
					// Complete hotkey
					1, 210, 4, 0, 0, 5, 0,
					// Partial hotkey (should be ignored)
					0, 46, 22,
				},
			},
			expectedCount: 1,
			expectHook:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := manager.HandleHotkeys(tc.args)
			if err != nil && tc.expectHook {
				t.Errorf("HandleHotkeys returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				hotkeys, ok := result["hotkeys"].([]Hotkey)
				if !ok {
					t.Errorf("Expected hotkeys to be a []Hotkey")
				} else if len(hotkeys) != tc.expectedCount {
					t.Errorf("Expected %d hotkeys, got %d", tc.expectedCount, len(hotkeys))
				}

				// Check specific hotkey values for the first test case
				if tc.name == "Valid Hotkeys" && len(hotkeys) >= 2 {
					// Check first hotkey
					if hotkeys[0].Type != HotkeyTypeSkill {
						t.Errorf("Expected first hotkey type to be skill, got %v", hotkeys[0].Type)
					}
					if hotkeys[0].ID != 1234 {
						t.Errorf("Expected first hotkey ID to be 1234, got %d", hotkeys[0].ID)
					}
					if hotkeys[0].Lv != 5 {
						t.Errorf("Expected first hotkey level to be 5, got %d", hotkeys[0].Lv)
					}

					// Check second hotkey
					if hotkeys[1].Type != HotkeyTypeItem {
						t.Errorf("Expected second hotkey type to be item, got %v", hotkeys[1].Type)
					}
					if hotkeys[1].ID != 5678 {
						t.Errorf("Expected second hotkey ID to be 5678, got %d", hotkeys[1].ID)
					}
					if hotkeys[1].Lv != 0 {
						t.Errorf("Expected second hotkey level to be 0, got %d", hotkeys[1].Lv)
					}
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleHotkeysUnhappy(t *testing.T) {
	// Create a hotkey manager for testing
	manager := NewHotkeyManager(nil)

	// Test cases for unhappy paths
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Missing Hotkeys",
			args: map[string]interface{}{},
		},
		{
			name: "Invalid Hotkeys Type",
			args: map[string]interface{}{
				"hotkeys": "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should return an error
			err := manager.HandleHotkeys(tc.args)
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
