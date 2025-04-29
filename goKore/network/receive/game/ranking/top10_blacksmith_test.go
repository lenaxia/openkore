package ranking

import (
	"strings"
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestTop10BlacksmithRank(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("ranking.top10_blacksmith_rank", func(hookName string, arg interface{}, userData interface{}) {
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
		"Smith1",
		"Smith2",
		"Smith3",
		"Smith4",
		"Smith5",
		"Smith6",
		"Smith7",
		"Smith8",
		"Smith9",
		"Smith10",
	}

	for _, name := range names {
		nameBytes := make([]byte, 24)
		copy(nameBytes, []byte(name))
		rawMsg = append(rawMsg, nameBytes...)
	}

	// Add 10 points (4 bytes each)
	points := []uint32{
		10000,
		9000,
		8000,
		7000,
		6000,
		5000,
		4000,
		3000,
		2000,
		1000,
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
	err := manager.handleTop10BlacksmithRank(args)
	if err != nil {
		t.Errorf("handleTop10BlacksmithRank returned an error: %v", err)
	}

	// Get the result from the channel
	result := <-resultChan

	// Verify the result
	rankings := result["rankings"].([]map[string]interface{})
	if len(rankings) != 10 {
		t.Errorf("Expected 10 rankings, got %d", len(rankings))
	}

	// Check the first ranking
	if rankings[0]["name"] != "Smith1" {
		t.Errorf("Expected name 'Smith1', got '%s'", rankings[0]["name"])
	}
	if rankings[0]["points"] != uint32(10000) {
		t.Errorf("Expected points 10000, got %d", rankings[0]["points"])
	}

	// Check the last ranking
	if rankings[9]["name"] != "Smith10" {
		t.Errorf("Expected name 'Smith10', got '%s'", rankings[9]["name"])
	}
	if rankings[9]["points"] != uint32(1000) {
		t.Errorf("Expected points 1000, got %d", rankings[9]["points"])
	}

	// Check the status message format (not the exact content)
	status := result["status"].(string)
	if len(status) == 0 {
		t.Errorf("Status message is empty")
	}

	// Verify the header contains "BLACKSMITH RANK"
	if !strings.Contains(status, "BLACKSMITH RANK") {
		t.Errorf("Expected status header to contain 'BLACKSMITH RANK', got '%s'", status[:30])
	}
}
