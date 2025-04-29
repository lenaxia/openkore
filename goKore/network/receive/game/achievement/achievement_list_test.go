package achievement

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestAchievementList(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("achievement.list", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewAchievementManager(nil, hookManager)

	// Create test data
	// Header (22 bytes) + 3 achievements (49 bytes each)
	rawMsg := make([]byte, 0)

	// Add header bytes (22 bytes)
	for i := 0; i < 22; i++ {
		rawMsg = append(rawMsg, 0x00)
	}

	// Add 3 achievements
	achievements := []struct {
		id          uint32
		completed   byte
		objectives  []uint32
		completedAt uint32
		reward      byte
	}{
		{
			id:          1001,
			completed:   1,
			objectives:  []uint32{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000},
			completedAt: 12345678,
			reward:      1,
		},
		{
			id:          1002,
			completed:   0,
			objectives:  []uint32{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			completedAt: 0,
			reward:      0,
		},
		{
			id:          1003,
			completed:   2,
			objectives:  []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			completedAt: 87654321,
			reward:      2,
		},
	}

	for _, achievement := range achievements {
		// Add achievement ID (4 bytes)
		rawMsg = append(rawMsg,
			byte(achievement.id&0xFF),
			byte((achievement.id>>8)&0xFF),
			byte((achievement.id>>16)&0xFF),
			byte((achievement.id>>24)&0xFF),
		)

		// Add completed status (1 byte)
		rawMsg = append(rawMsg, achievement.completed)

		// Add 10 objectives (4 bytes each)
		for _, objective := range achievement.objectives {
			rawMsg = append(rawMsg,
				byte(objective&0xFF),
				byte((objective>>8)&0xFF),
				byte((objective>>16)&0xFF),
				byte((objective>>24)&0xFF),
			)
		}

		// Add completed_at timestamp (4 bytes)
		rawMsg = append(rawMsg,
			byte(achievement.completedAt&0xFF),
			byte((achievement.completedAt>>8)&0xFF),
			byte((achievement.completedAt>>16)&0xFF),
			byte((achievement.completedAt>>24)&0xFF),
		)

		// Add reward status (1 byte)
		rawMsg = append(rawMsg, achievement.reward)
	}

	// Test case
	args := map[string]interface{}{
		"RAW_MSG":      rawMsg,
		"RAW_MSG_SIZE": len(rawMsg),
	}

	// Call the handler
	err := manager.HandleAchievementList(args)
	if err != nil {
		t.Errorf("HandleAchievementList returned an error: %v", err)
	}

	// Get the result from the channel
	result := <-resultChan

	// Verify the result
	achievementList := result["achievements"].(map[uint32]map[string]interface{})
	if len(achievementList) != 3 {
		t.Errorf("Expected 3 achievements, got %d", len(achievementList))
	}

	// Check the first achievement
	achievement := achievementList[1001]
	if achievement["achievementID"] != uint32(1001) {
		t.Errorf("Expected achievementID 1001, got %d", achievement["achievementID"])
	}
	if achievement["completed"] != byte(1) {
		t.Errorf("Expected completed 1, got %d", achievement["completed"])
	}
	if achievement["objective1"] != uint32(100) {
		t.Errorf("Expected objective1 100, got %d", achievement["objective1"])
	}
	if achievement["completed_at"] != uint32(12345678) {
		t.Errorf("Expected completed_at 12345678, got %d", achievement["completed_at"])
	}
	if achievement["reward"] != byte(1) {
		t.Errorf("Expected reward 1, got %d", achievement["reward"])
	}

	// Check the status message
	status := result["status"].(string)
	if len(status) == 0 {
		t.Errorf("Status message is empty")
	}
}
