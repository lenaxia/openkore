package ranking

import (
	"strings"
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestTop10PkRank(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("ranking.top10_pk_rank", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewRankingManager(nil, hookManager)

	// Create test data
	// 10 names (24 bytes each) followed by 10 points (4 bytes each)
	rawMsg := make([]byte, 0)

	// Add header bytes (2 bytes)
	rawMsg = append(rawMsg, 0x00, 0x00)

	// Add 10 names (24 bytes each)
	names := []string{
		"PKPlayer1",
		"PKPlayer2",
		"PKPlayer3",
		"PKPlayer4",
		"PKPlayer5",
		"PKPlayer6",
		"PKPlayer7",
		"PKPlayer8",
		"PKPlayer9",
		"PKPlayer10",
	}

	for _, name := range names {
		nameBytes := make([]byte, 24)
		copy(nameBytes, []byte(name))
		rawMsg = append(rawMsg, nameBytes...)
	}

	// Add 10 points (4 bytes each)
	points := []uint32{
		5000,
		4500,
		4000,
		3500,
		3000,
		2500,
		2000,
		1500,
		1000,
		500,
	}

	for _, point := range points {
		pointBytes := []byte{
			byte(point & 0xFF),
			byte((point >> 8) & 0xFF),
			byte((point >> 16) & 0xFF),
			byte((point >> 24) & 0xFF),
		}
		rawMsg = append(rawMsg, pointBytes...)
	}

	// Test case
	args := map[string]interface{}{
		"RAW_MSG": rawMsg,
	}

	// Call the handler
	err := manager.handleTop10PkRank(args)
	if err != nil {
		t.Errorf("handleTop10PkRank returned an error: %v", err)
	}

	// Get the result from the channel
	result := <-resultChan

	// Verify the result
	rankings := result["rankings"].([]map[string]interface{})
	if len(rankings) != 10 {
		t.Errorf("Expected 10 rankings, got %d", len(rankings))
	}

	// Check the first ranking
	if rankings[0]["name"] != "PKPlayer1" {
		t.Errorf("Expected name 'PKPlayer1', got '%s'", rankings[0]["name"])
	}
	if rankings[0]["points"] != uint32(5000) {
		t.Errorf("Expected points 5000, got %d", rankings[0]["points"])
	}

	// Check the last ranking
	if rankings[9]["name"] != "PKPlayer10" {
		t.Errorf("Expected name 'PKPlayer10', got '%s'", rankings[9]["name"])
	}
	if rankings[9]["points"] != uint32(500) {
		t.Errorf("Expected points 500, got %d", rankings[9]["points"])
	}

	// Check the status message format (not the exact content)
	status := result["status"].(string)
	if len(status) == 0 {
		t.Errorf("Status message is empty")
	}

	// Verify the header contains "PVP RANK"
	if !strings.Contains(status, "PVP RANK") {
		t.Errorf("Expected status header to contain 'PVP RANK', got '%s'", status[:30])
	}
}
