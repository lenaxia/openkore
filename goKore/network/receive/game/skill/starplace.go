package skill

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// StarplaceManager manages the starplace packet handler
type StarplaceManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewStarplaceManager creates a new starplace manager
func NewStarplaceManager(parser *core.CoreParser, hookManager *hooks.HookManager) *StarplaceManager {
	return &StarplaceManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to starplace
func (m *StarplaceManager) RegisterHandlers() {
	// Register starplace handler (0253)
	m.parser.RegisterHandlerFunc("0253", "starplace", "C",
		[]string{"which"},
		m.handleStarplace)
}

// handleStarplace handles the starplace packet
// Packet format: 0253 <which>.B
// Star Gladiator's Feeling map confirmation prompt
func (m *StarplaceManager) handleStarplace(args map[string]interface{}) error {
	// Process the packet
	result := m.processStarplace(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.starplace", result)
	}

	return nil
}

// processStarplace processes the starplace packet and returns a structured result
func (m *StarplaceManager) processStarplace(args map[string]interface{}) map[string]interface{} {
	// Extract starplace information from args
	var which uint8

	// Extract which
	if val, ok := args["which"].(uint8); ok {
		which = val
	}

	// Generate message
	message := m.generateStarplaceMessage(which)

	// Return structured result
	return map[string]interface{}{
		"which":   which,
		"message": message,
	}
}

// generateStarplaceMessage generates a message for the starplace
func (m *StarplaceManager) generateStarplaceMessage(which uint8) string {
	return fmt.Sprintf("Star Gladiator's Feeling map confirmation prompt: %d", which)
}
