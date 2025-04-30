package skill

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// HighJumpManager manages the high_jump packet handler
type HighJumpManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewHighJumpManager creates a new high jump manager
func NewHighJumpManager(parser *core.CoreParser, hookManager *hooks.HookManager) *HighJumpManager {
	return &HighJumpManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to high jump
func (m *HighJumpManager) RegisterHandlers() {
	// Register high_jump handler
	m.parser.RegisterHandlerFunc("01FF", "high_jump", "a4 v2",
		[]string{"ID", "x", "y"},
		m.handleHighJump)
}

// handleHighJump handles the high_jump packet
// Packet format: 01FF <ID>.L <x>.W <y>.W
func (m *HighJumpManager) handleHighJump(args map[string]interface{}) error {
	// Process the packet
	result := m.processHighJump(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.high_jump", result)
	}

	return nil
}

// processHighJump processes the high_jump packet and returns a structured result
func (m *HighJumpManager) processHighJump(args map[string]interface{}) map[string]interface{} {
	// Extract high jump information from args
	var actorID uint32
	var x, y uint16

	// Extract actorID
	if val, ok := args["ID"].(uint32); ok {
		actorID = val
	}

	// Extract x coordinate
	if val, ok := args["x"].(uint16); ok {
		x = val
	}

	// Extract y coordinate
	if val, ok := args["y"].(uint16); ok {
		y = val
	}

	// Check if the actor exists and if the move was successful
	// In a real implementation, this would check if the actor exists and if the move was successful
	// For now, we'll assume the move was successful
	moveSuccessful := true

	// Generate message
	message := m.generateHighJumpMessage(actorID, x, y, moveSuccessful)

	// Return structured result
	return map[string]interface{}{
		"actorID":        actorID,
		"x":              x,
		"y":              y,
		"moveSuccessful": moveSuccessful,
		"message":        message,
	}
}

// generateHighJumpMessage generates a message for the high jump skill
func (m *HighJumpManager) generateHighJumpMessage(actorID uint32, x, y uint16, moveSuccessful bool) string {
	// Get actor name
	actorName := m.getActorName(actorID)

	// Generate message based on move success
	if moveSuccessful {
		return fmt.Sprintf("%s instantly moved to %d, %d", actorName, x, y)
	} else {
		return fmt.Sprintf("%s failed to instantly move", actorName)
	}
}

// getActorName returns the name of an actor by ID
// In a real implementation, this would look up the actor name from a database or map
func (m *HighJumpManager) getActorName(actorID uint32) string {
	// This is a simplified implementation
	// In a real implementation, this would look up the actor name from a database or map
	return fmt.Sprintf("Actor_%d", actorID)
}
