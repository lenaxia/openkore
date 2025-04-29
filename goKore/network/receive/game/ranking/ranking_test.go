package ranking

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestPvpRank(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("ranking.pvp_rank", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewRankingManager(nil, hookManager)

	// Test cases
	testCases := []struct {
		name           string
		args           map[string]interface{}
		initialRank    uint16
		initialNum     uint16
		expectedStatus string
		expectedRank   uint16
		expectedNum    uint16
		pvpEnabled     bool
		expectMessage  bool // Flag to indicate if we expect a message
	}{
		{
			name: "PVP Enabled - Rank Update",
			args: map[string]interface{}{
				"ID":   uint32(12345),
				"rank": uint16(10),
				"num":  uint16(100),
			},
			initialRank:    0, // Different from the new rank
			initialNum:     0, // Different from the new num
			expectedStatus: "Your PvP rank is: 10/100",
			expectedRank:   uint16(10),
			expectedNum:    uint16(100),
			pvpEnabled:     true,
			expectMessage:  true,
		},
		{
			name: "PVP Disabled - No Message",
			args: map[string]interface{}{
				"ID":   uint32(12345),
				"rank": uint16(10),
				"num":  uint16(100),
			},
			initialRank:    0, // Different from the new rank
			initialNum:     0, // Different from the new num
			expectedStatus: "",
			expectedRank:   uint16(10),
			expectedNum:    uint16(100),
			pvpEnabled:     false,
			expectMessage:  false,
		},
		{
			name: "Same Rank - No Update",
			args: map[string]interface{}{
				"ID":   uint32(12345),
				"rank": uint16(10),
				"num":  uint16(100),
			},
			initialRank:    10,  // Same as the new rank
			initialNum:     100, // Same as the new num
			expectedStatus: "",
			expectedRank:   uint16(10),
			expectedNum:    uint16(100),
			pvpEnabled:     true,
			expectMessage:  false,
		},
		{
			name: "Different Rank - Update",
			args: map[string]interface{}{
				"ID":   uint32(12345),
				"rank": uint16(5),
				"num":  uint16(100),
			},
			initialRank:    10,  // Different from the new rank
			initialNum:     100, // Same as the new num
			expectedStatus: "Your PvP rank is: 5/100",
			expectedRank:   uint16(5),
			expectedNum:    uint16(100),
			pvpEnabled:     true,
			expectMessage:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set initial state for this test
			manager.pvpRank = tc.initialRank
			manager.pvpNum = tc.initialNum
			manager.pvpEnabled = tc.pvpEnabled

			// Call the handler
			err := manager.handlePvpRank(tc.args)
			if err != nil {
				t.Errorf("handlePvpRank returned an error: %v", err)
			}

			// For tests that expect a message, check the hook result
			if tc.expectMessage {
				// Get the result from the channel
				select {
				case result := <-resultChan:
					// Verify the result
					if result["status"] != tc.expectedStatus {
						t.Errorf("Expected status %q, got %q", tc.expectedStatus, result["status"])
					}
				default:
					t.Errorf("Expected a message but none was received")
				}
			} else {
				// Make sure no message was sent
				select {
				case result := <-resultChan:
					t.Errorf("Unexpected message received: %v", result)
				default:
					// This is the expected case - no message
				}
			}

			// Verify the internal state was updated correctly
			if manager.pvpRank != tc.expectedRank {
				t.Errorf("Expected pvpRank %v, got %v", tc.expectedRank, manager.pvpRank)
			}
			if manager.pvpNum != tc.expectedNum {
				t.Errorf("Expected pvpNum %v, got %v", tc.expectedNum, manager.pvpNum)
			}
		})
	}
}
