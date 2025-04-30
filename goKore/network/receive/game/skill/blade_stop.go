package skill

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// BladeStopManager manages the blade_stop packet handler
type BladeStopManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewBladeStopManager creates a new blade stop manager
func NewBladeStopManager(parser *core.CoreParser, hookManager *hooks.HookManager) *BladeStopManager {
	return &BladeStopManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to blade stop
func (m *BladeStopManager) RegisterHandlers() {
	// Register blade_stop handler
	m.parser.RegisterHandlerFunc("01D1", "blade_stop", "a4 a4 C",
		[]string{"sourceID", "targetID", "active"},
		m.handleBladeStop)
}

// handleBladeStop handles the blade_stop packet
// Packet format: 01D1 <sourceID>.L <targetID>.L <active>.B
func (m *BladeStopManager) handleBladeStop(args map[string]interface{}) error {
	// Process the packet
	result := m.processBladeStop(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.blade_stop", result)
	}

	return nil
}

// processBladeStop processes the blade_stop packet and returns a structured result
func (m *BladeStopManager) processBladeStop(args map[string]interface{}) map[string]interface{} {
	// Extract blade stop information from args
	var sourceID uint32
	var targetID uint32
	var active uint8

	// Extract sourceID
	if val, ok := args["sourceID"].(uint32); ok {
		sourceID = val
	}

	// Extract targetID
	if val, ok := args["targetID"].(uint32); ok {
		targetID = val
	}

	// Extract active
	if val, ok := args["active"].(uint8); ok {
		active = val
	}

	// Generate message
	message := m.generateBladeStopMessage(sourceID, targetID, active)

	// Return structured result
	return map[string]interface{}{
		"sourceID": sourceID,
		"targetID": targetID,
		"active":   active,
		"message":  message,
	}
}

// generateBladeStopMessage generates a message for the blade stop skill
func (m *BladeStopManager) generateBladeStopMessage(sourceID, targetID uint32, active uint8) string {
	// Get source and target actor names
	sourceName := m.getActorName(sourceID)
	targetName := m.getActorName(targetID)

	// Generate message based on active status
	if active == 0 {
		return fmt.Sprintf("Blade Stop by %s on %s is deactivated.", sourceName, targetName)
	} else {
		return fmt.Sprintf("Blade Stop by %s on %s is active.", sourceName, targetName)
	}
}

// getActorName returns the name of an actor by ID
// In a real implementation, this would look up the actor name from a database or map
func (m *BladeStopManager) getActorName(actorID uint32) string {
	// This is a simplified implementation
	// In a real implementation, this would look up the actor name from a database or map
	return fmt.Sprintf("Actor_%d", actorID)
}
