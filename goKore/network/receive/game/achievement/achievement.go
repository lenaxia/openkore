// Package achievement provides handlers for achievement-related packets.
package achievement

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// AchievementManager manages achievement-related packet handlers
type AchievementManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager

	// Internal state
	achievementList map[uint32]map[string]interface{}
}

// NewAchievementManager creates a new achievement manager
func NewAchievementManager(parser *core.CoreParser, hookManager *hooks.HookManager) *AchievementManager {
	return &AchievementManager{
		parser:          parser,
		hookManager:     hookManager,
		achievementList: make(map[uint32]map[string]interface{}),
	}
}

// RegisterAchievementHandlers registers all handlers related to achievements
func (m *AchievementManager) RegisterAchievementHandlers() {
	// Register achievement handlers
	if m.parser != nil {
		// Register achievement_list handler
		m.parser.RegisterHandlerFunc("0A23", "achievement_list", "a*",
			[]string{"RAW_MSG"},
			m.HandleAchievementList)

		// Register achievement_update handler
		m.parser.RegisterHandlerFunc("0A24", "achievement_update", "V C V10 V C",
			[]string{"achievementID", "completed", "objective1", "objective2", "objective3", "objective4", "objective5", "objective6", "objective7", "objective8", "objective9", "objective10", "completed_at", "reward"},
			m.HandleAchievementUpdate)

		// Register achievement_reward_ack handler
		m.parser.RegisterHandlerFunc("0A25", "achievement_reward_ack", "V",
			[]string{"achievementID"},
			m.HandleAchievementRewardAck)
	}
}

// RegisterAllHandlers registers all achievement-related handlers
func (m *AchievementManager) RegisterAllHandlers() {
	// Register achievement handlers
	m.RegisterAchievementHandlers()
}

// HandleAchievementList handles the achievement_list packet
// Packet format: 0A23 <packet len>.W <ACH count>.L <ACH ID>.L <completed>.B <objective[]>.L*10 <completed at>.L <reward>.B
func (m *AchievementManager) HandleAchievementList(args map[string]interface{}) error {
	// Process the packet
	result := m.processAchievementList(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("achievement.list", result)
	}

	return nil
}

// processAchievementList processes the achievement_list packet and returns a structured result
func (m *AchievementManager) processAchievementList(args map[string]interface{}) map[string]interface{} {
	var rawMsg []byte
	var msgSize int

	// Reset achievement list
	m.achievementList = make(map[uint32]map[string]interface{})

	// Extract raw message from args
	if rawMsgVal, ok := args["RAW_MSG"].([]byte); ok {
		rawMsg = rawMsgVal
	}

	// Extract message size from args
	if msgSizeVal, ok := args["RAW_MSG_SIZE"].(int); ok {
		msgSize = msgSizeVal
	} else {
		msgSize = len(rawMsg)
	}

	// Define constants
	headerLen := 22
	achieveLen := 49 // 4 (ID) + 1 (completed) + 10*4 (objectives) + 4 (completed_at) + 1 (reward)

	// Process each achievement
	for i := headerLen; i < msgSize; i += achieveLen {
		if i+achieveLen > msgSize {
			break
		}

		achieve := make(map[string]interface{})

		// Extract achievement ID (4 bytes)
		achieve["achievementID"] = uint32(rawMsg[i]) |
			(uint32(rawMsg[i+1]) << 8) |
			(uint32(rawMsg[i+2]) << 16) |
			(uint32(rawMsg[i+3]) << 24)

		// Extract completed status (1 byte)
		achieve["completed"] = rawMsg[i+4]

		// Extract 10 objectives (4 bytes each)
		for j := 0; j < 10; j++ {
			objectiveOffset := i + 5 + (j * 4)
			objectiveName := fmt.Sprintf("objective%d", j+1)
			achieve[objectiveName] = uint32(rawMsg[objectiveOffset]) |
				(uint32(rawMsg[objectiveOffset+1]) << 8) |
				(uint32(rawMsg[objectiveOffset+2]) << 16) |
				(uint32(rawMsg[objectiveOffset+3]) << 24)
		}

		// Extract completed_at timestamp (4 bytes)
		completedAtOffset := i + 5 + (10 * 4)
		achieve["completed_at"] = uint32(rawMsg[completedAtOffset]) |
			(uint32(rawMsg[completedAtOffset+1]) << 8) |
			(uint32(rawMsg[completedAtOffset+2]) << 16) |
			(uint32(rawMsg[completedAtOffset+3]) << 24)

		// Extract reward status (1 byte)
		achieve["reward"] = rawMsg[completedAtOffset+4]

		// Add to achievement list
		achievementID := achieve["achievementID"].(uint32)
		m.achievementList[achievementID] = achieve

		// Log message
		// Note: In the original Perl code, there's a reference to $achievements which seems to be
		// a global variable containing achievement titles. We'll just use the ID for now.
		fmt.Printf("Achievement '%d' added.\n", achievementID)
	}

	// Create the result
	result := map[string]interface{}{
		"achievements": m.achievementList,
		"status":       fmt.Sprintf("Received %d achievements.", len(m.achievementList)),
	}

	return result
}

// HandleAchievementUpdate handles the achievement_update packet
// Packet format: 0A24 <ACH ID>.L <completed>.B <objective[]>.L*10 <completed at>.L <reward>.B
func (m *AchievementManager) HandleAchievementUpdate(args map[string]interface{}) error {
	// Process the packet
	result := m.processAchievementUpdate(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("achievement.update", result)
	}

	return nil
}

// processAchievementUpdate processes the achievement_update packet and returns a structured result
func (m *AchievementManager) processAchievementUpdate(args map[string]interface{}) map[string]interface{} {
	// Create a new achievement entry
	achieve := make(map[string]interface{})

	// Copy all fields from args to achieve
	for key, value := range args {
		achieve[key] = value
	}

	// Add to achievement list
	achievementID := achieve["achievementID"].(uint32)
	m.achievementList[achievementID] = achieve

	// Log message
	// Note: In the original Perl code, there's a reference to $achievements which seems to be
	// a global variable containing achievement titles. We'll just use the ID for now.
	fmt.Printf("Achievement '%d' added or updated.\n", achievementID)

	// Create the result
	result := map[string]interface{}{
		"achievement": achieve,
		"status":      fmt.Sprintf("Achievement '%d' added or updated.", achievementID),
	}

	return result
}

// HandleAchievementRewardAck handles the achievement_reward_ack packet
// Packet format: 0A25 <ACH ID>.L
func (m *AchievementManager) HandleAchievementRewardAck(args map[string]interface{}) error {
	// Process the packet
	result := m.processAchievementRewardAck(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("achievement.reward_ack", result)
	}

	return nil
}

// processAchievementRewardAck processes the achievement_reward_ack packet and returns a structured result
func (m *AchievementManager) processAchievementRewardAck(args map[string]interface{}) map[string]interface{} {
	var achievementID uint32

	// Extract achievement ID from args
	if idVal, ok := args["achievementID"].(uint32); ok {
		achievementID = idVal
	}

	// Log message
	// Note: In the original Perl code, there's a reference to $achievements which seems to be
	// a global variable containing achievement titles. We'll just use the ID for now.
	fmt.Printf("Received reward for achievement '%d'.\n", achievementID)

	// Create the result
	result := map[string]interface{}{
		"achievementID": achievementID,
		"status":        fmt.Sprintf("Received reward for achievement '%d'.", achievementID),
	}

	return result
}
