package skill

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// GospelBuffManager manages the gospel_buff_aligned packet handler
type GospelBuffManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewGospelBuffManager creates a new gospel buff manager
func NewGospelBuffManager(parser *core.CoreParser, hookManager *hooks.HookManager) *GospelBuffManager {
	return &GospelBuffManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to gospel buffs
func (m *GospelBuffManager) RegisterHandlers() {
	// Register gospel_buff_aligned handler
	m.parser.RegisterHandlerFunc("0215", "gospel_buff_aligned", "V",
		[]string{"ID"},
		m.handleGospelBuffAligned)
}

// handleGospelBuffAligned handles the gospel_buff_aligned packet
// Packet format: 0215 <msg id>.L
func (m *GospelBuffManager) handleGospelBuffAligned(args map[string]interface{}) error {
	// Process the packet
	result := m.processGospelBuffAligned(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.gospel_buff_aligned", result)
	}

	return nil
}

// processGospelBuffAligned processes the gospel_buff_aligned packet and returns a structured result
func (m *GospelBuffManager) processGospelBuffAligned(args map[string]interface{}) map[string]interface{} {
	// Extract status ID from args
	var statusID uint32

	// Extract ID
	if val, ok := args["ID"].(uint32); ok {
		statusID = val
	}

	// Get the message based on the status ID
	message := m.getGospelMessage(statusID)

	// Return structured result
	return map[string]interface{}{
		"statusID": statusID,
		"message":  message,
	}
}

// getGospelMessage returns the message for a given gospel status ID
func (m *GospelBuffManager) getGospelMessage(statusID uint32) string {
	// Map of status IDs to messages
	messages := map[uint32]string{
		21: "All abnormal status effects have been removed.",
		22: "You will be immune to abnormal status effects for the next minute.",
		23: "Your Max HP will stay increased for the next minute.",
		24: "Your Max SP will stay increased for the next minute.",
		25: "All of your Stats will stay increased for the next minute.",
		28: "Your weapon will remain blessed with Holy power for the next minute.",
		29: "Your armor will remain blessed with Holy power for the next minute.",
		30: "Your Defense will stay increased for the next 10 seconds.",
		31: "Your Attack strength will stay increased for the next minute.",
		32: "Your Accuracy and Flee Rate will stay increased for the next minute.",
		40: "Full strip failed because of coating.",
	}

	// Return the message for the status ID, or an empty string if not found
	if message, ok := messages[statusID]; ok {
		return message
	}
	return "Unknown gospel buff effect."
}
