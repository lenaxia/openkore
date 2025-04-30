package skill

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// ResurrectionManager manages the resurrection packet handler
type ResurrectionManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewResurrectionManager creates a new resurrection manager
func NewResurrectionManager(parser *core.CoreParser, hookManager *hooks.HookManager) *ResurrectionManager {
	return &ResurrectionManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to resurrection
func (m *ResurrectionManager) RegisterHandlers() {
	// Register resurrection handler
	m.parser.RegisterHandlerFunc("0148", "resurrection", "a4 C",
		[]string{"targetID", "type"},
		m.handleResurrection)
}

// handleResurrection handles the resurrection packet
// Packet format: 0148 <ID>.L <type>.B
func (m *ResurrectionManager) handleResurrection(args map[string]interface{}) error {
	// Process the packet
	result := m.processResurrection(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.resurrection", result)
	}

	return nil
}

// processResurrection processes the resurrection packet and returns a structured result
func (m *ResurrectionManager) processResurrection(args map[string]interface{}) map[string]interface{} {
	// Extract resurrection information from args
	var targetID uint32
	var resType uint8

	// Extract targetID
	if val, ok := args["targetID"].(uint32); ok {
		targetID = val
	}

	// Extract type
	if val, ok := args["type"].(uint8); ok {
		resType = val
	}

	// Determine if this is the player's character or another entity
	isOwnCharacter := m.isOwnCharacter(targetID)
	isHomunculus := m.isHomunculus(targetID)

	// Generate message
	message := m.generateResurrectionMessage(targetID, isOwnCharacter, isHomunculus)

	// Return structured result
	return map[string]interface{}{
		"targetID":       targetID,
		"type":           resType,
		"isOwnCharacter": isOwnCharacter,
		"isHomunculus":   isHomunculus,
		"message":        message,
	}
}

// isOwnCharacter checks if the target ID is the player's character
// In a real implementation, this would check against the player's account ID
func (m *ResurrectionManager) isOwnCharacter(targetID uint32) bool {
	// This is a simplified implementation
	// In a real implementation, this would check against the player's account ID
	// For now, we'll use a placeholder value
	accountID := uint32(1000) // Placeholder
	return targetID == accountID
}

// isHomunculus checks if the target ID is a homunculus
// In a real implementation, this would check if the target is a homunculus
func (m *ResurrectionManager) isHomunculus(targetID uint32) bool {
	// This is a simplified implementation
	// In a real implementation, this would check if the target is a homunculus
	// For now, we'll return false
	return false
}

// generateResurrectionMessage generates a message for the resurrection
func (m *ResurrectionManager) generateResurrectionMessage(targetID uint32, isOwnCharacter bool, isHomunculus bool) string {
	if isOwnCharacter {
		return "You have been resurrected"
	} else if isHomunculus {
		return fmt.Sprintf("Slave Resurrected: %s", m.getActorName(targetID))
	} else {
		return fmt.Sprintf("%s has been resurrected", m.getActorName(targetID))
	}
}

// getActorName returns the name of an actor by ID
// In a real implementation, this would look up the actor name from a database or map
func (m *ResurrectionManager) getActorName(actorID uint32) string {
	// This is a simplified implementation
	// In a real implementation, this would look up the actor name from a database or map
	return fmt.Sprintf("Actor_%d", actorID)
}
