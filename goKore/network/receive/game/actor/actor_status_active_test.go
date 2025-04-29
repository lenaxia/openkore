package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleActorStatusActive(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.status_active", func(hookName string, arg interface{}, userData interface{}) {
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
		statusType uint16
		statusName string
		flag       byte
		tick       uint32
		expectHook bool
	}{
		{
			name: "Player Status Active",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": uint16(1), // Assuming 1 is a valid status type
				"tick": uint32(5000),
				"flag": byte(1),
			},
			statusType: 1,
			statusName: "UNKNOWN_STATUS_1", // This would be mapped to a real status name in production
			flag:       1,
			tick:       5000,
			expectHook: true,
		},
		{
			name: "Player Status Inactive",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": uint16(2), // Assuming 2 is a valid status type
				"tick": uint32(0),
				"flag": byte(0),
			},
			statusType: 2,
			statusName: "UNKNOWN_STATUS_2", // This would be mapped to a real status name in production
			flag:       0,
			tick:       0,
			expectHook: true,
		},
		{
			name: "Unknown Actor",
			args: map[string]interface{}{
				"ID":   []byte{5, 6, 7, 8},
				"type": uint16(1),
				"tick": uint32(5000),
				"flag": byte(1),
			},
			expectHook: false,
		},
		{
			name: "Cart Active",
			args: map[string]interface{}{
				"ID":       []byte{1, 2, 3, 4},
				"type":     uint16(673), // Cart active status
				"tick":     uint32(5000),
				"flag":     byte(1),
				"unknown1": uint32(1), // Cart type
			},
			statusType: 673,
			statusName: "CART_ACTIVE", // This would be mapped to a real status name in production
			flag:       1,
			tick:       5000,
			expectHook: true,
		},
		{
			name: "Rolling Cutter",
			args: map[string]interface{}{
				"ID":       []byte{1, 2, 3, 4},
				"type":     uint16(0x153), // Rolling Cutter status
				"tick":     uint32(5000),
				"flag":     byte(1),
				"unknown1": uint32(3), // Number of counters
			},
			statusType: 0x153,
			statusName: "ROLLING_CUTTER", // This would be mapped to a real status name in production
			flag:       1,
			tick:       5000,
			expectHook: true,
		},
		{
			name: "Infinite Duration",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": uint16(3),
				"tick": uint32(9999), // Special value for infinite duration
				"flag": byte(1),
			},
			statusType: 3,
			statusName: "UNKNOWN_STATUS_3", // This would be mapped to a real status name in production
			flag:       1,
			tick:       0, // Should be converted to 0 (infinite)
			expectHook: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleActorStatusActive(tc.args)
			if err != nil {
				t.Errorf("HandleActorStatusActive returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if id, ok := result["ID"].([]byte); !ok || string(id) != string(tc.args["ID"].([]byte)) {
					t.Errorf("Expected ID %v, got %v", tc.args["ID"], result["ID"])
				}

				if statusType, ok := result["type"].(uint16); !ok || statusType != tc.statusType {
					t.Errorf("Expected status type %d, got %v", tc.statusType, result["type"])
				}

				if statusName, ok := result["statusName"].(string); !ok || statusName != tc.statusName {
					t.Errorf("Expected status name %s, got %v", tc.statusName, result["statusName"])
				}

				if flag, ok := result["flag"].(byte); !ok || flag != tc.flag {
					t.Errorf("Expected flag %d, got %v", tc.flag, result["flag"])
				}

				// Special case for cart active
				if tc.statusType == 673 {
					// Check if cart type was updated
					if player.CartType() != tc.args["unknown1"].(uint32) {
						t.Errorf("Expected cart type %d, got %d", tc.args["unknown1"].(uint32), player.CartType())
					}
				}

				// Special case for rolling cutter
				if tc.statusType == 0x153 {
					// Check if spirits count was updated
					if player.Spirits() != tc.args["unknown1"].(uint32) {
						t.Errorf("Expected spirits %d, got %d", tc.args["unknown1"].(uint32), player.Spirits())
					}
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleActorStatusActiveUnhappy(t *testing.T) {
	// Create a handler for testing
	handler := NewHandler()

	// Test cases for unhappy paths
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Missing ID",
			args: map[string]interface{}{
				"type": uint16(1),
				"tick": uint32(5000),
				"flag": byte(1),
			},
		},
		{
			name: "Missing type",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"tick": uint32(5000),
				"flag": byte(1),
			},
		},
		{
			name: "Invalid ID type",
			args: map[string]interface{}{
				"ID":   "invalid",
				"type": uint16(1),
				"tick": uint32(5000),
				"flag": byte(1),
			},
		},
		{
			name: "Invalid type type",
			args: map[string]interface{}{
				"ID":   []byte{1, 2, 3, 4},
				"type": "invalid",
				"tick": uint32(5000),
				"flag": byte(1),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleActorStatusActive(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}
