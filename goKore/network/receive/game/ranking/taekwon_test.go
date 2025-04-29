package ranking

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestTaekwonRank(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("ranking.taekwon_rank", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewRankingManager(nil, hookManager)

	// Test cases
	testCases := []struct {
		name           string
		args           map[string]interface{}
		expectedStatus string
		expectedRank   uint32
	}{
		{
			name: "Rank 1",
			args: map[string]interface{}{
				"rank": uint32(1),
			},
			expectedStatus: "TaeKwon Mission Rank : 1",
			expectedRank:   uint32(1),
		},
		{
			name: "Rank 10",
			args: map[string]interface{}{
				"rank": uint32(10),
			},
			expectedStatus: "TaeKwon Mission Rank : 10",
			expectedRank:   uint32(10),
		},
		{
			name: "Rank 100",
			args: map[string]interface{}{
				"rank": uint32(100),
			},
			expectedStatus: "TaeKwon Mission Rank : 100",
			expectedRank:   uint32(100),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := manager.handleTaekwonRank(tc.args)
			if err != nil {
				t.Errorf("handleTaekwonRank returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-resultChan

			// Verify the result
			if result["status"] != tc.expectedStatus {
				t.Errorf("Expected status %q, got %q", tc.expectedStatus, result["status"])
			}
			if result["rank"] != tc.args["rank"] {
				t.Errorf("Expected rank %v, got %v", tc.args["rank"], result["rank"])
			}
		})
	}
}
