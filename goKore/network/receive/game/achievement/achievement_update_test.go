package achievement

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestAchievementUpdate(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("achievement.update", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewAchievementManager(nil, hookManager)

	// Test case
	args := map[string]interface{}{
		"achievementID": uint32(1001),
		"completed":     byte(1),
		"objective1":    uint32(100),
		"objective2":    uint32(200),
		"objective3":    uint32(300),
		"objective4":    uint32(400),
		"objective5":    uint32(500),
		"objective6":    uint32(600),
		"objective7":    uint32(700),
		"objective8":    uint32(800),
		"objective9":    uint32(900),
		"objective10":   uint32(1000),
		"completed_at":  uint32(12345678),
		"reward":        byte(1),
	}

	// Call the handler
	err := manager.HandleAchievementUpdate(args)
	if err != nil {
		t.Errorf("HandleAchievementUpdate returned an error: %v", err)
	}

	// Get the result from the channel
	result := <-resultChan

	// Verify the result
	achievement := result["achievement"].(map[string]interface{})
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

func TestAchievementRewardAck(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("achievement.reward_ack", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewAchievementManager(nil, hookManager)

	// Test case
	args := map[string]interface{}{
		"achievementID": uint32(1001),
	}

	// Call the handler
	err := manager.HandleAchievementRewardAck(args)
	if err != nil {
		t.Errorf("HandleAchievementRewardAck returned an error: %v", err)
	}

	// Get the result from the channel
	result := <-resultChan

	// Verify the result
	if result["achievementID"] != uint32(1001) {
		t.Errorf("Expected achievementID 1001, got %d", result["achievementID"])
	}

	// Check the status message
	status := result["status"].(string)
	if len(status) == 0 {
		t.Errorf("Status message is empty")
	}
}
