package core

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
)

// RemainTimeManager manages remain time-related functionality
type RemainTimeManager struct {
	hookManager *hooks.HookManager
}

// NewRemainTimeManager creates a new remain time manager
func NewRemainTimeManager(hookManager *hooks.HookManager) *RemainTimeManager {
	return &RemainTimeManager{
		hookManager: hookManager,
	}
}

// HandleRemainTimeInfo handles the remain_time_info packet
// Packet format: 0C42 <result>.W <expiration_date>.L <remain_time>.L
func (m *RemainTimeManager) HandleRemainTimeInfo(args map[string]interface{}) error {
	// Extract result with safety check
	var result uint16
	if resultVal, ok := args["result"].(uint16); ok {
		result = resultVal
	}

	// Extract expiration date with safety check
	var expirationDate uint32
	if expirationDateVal, ok := args["expiration_date"].(uint32); ok {
		expirationDate = expirationDateVal
	}

	// Extract remain time with safety check
	var remainTime uint32
	if remainTimeVal, ok := args["remain_time"].(uint32); ok {
		remainTime = remainTimeVal
	}

	// Format the message
	message := fmt.Sprintf("Remain Time - Result: %d - Expiration Date: %d - Time: %d",
		result, expirationDate, remainTime)

	// Log the message
	// In a real implementation, this would use a proper logger
	// logger.Debug(message)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("core.remain_time_info", map[string]interface{}{
			"result":          result,
			"expiration_date": expirationDate,
			"remain_time":     remainTime,
			"message":         message,
		})
	}

	return nil
}

// RegisterHandlers registers remain time-related packet handlers with the given parser
func (m *RemainTimeManager) RegisterHandlers(parser *CoreParser) {
	parser.RegisterHandlerFunc("0C42", "remain_time_info", "W L L",
		[]string{"result", "expiration_date", "remain_time"},
		m.HandleRemainTimeInfo)
}
