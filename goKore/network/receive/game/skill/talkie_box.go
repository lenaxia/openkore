package skill

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TalkieBoxManager manages the talkie_box packet handler
type TalkieBoxManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewTalkieBoxManager creates a new talkie box manager
func NewTalkieBoxManager(parser *core.CoreParser, hookManager *hooks.HookManager) *TalkieBoxManager {
	return &TalkieBoxManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to talkie box
func (m *TalkieBoxManager) RegisterHandlers() {
	// Register talkie_box handler (0191)
	m.parser.RegisterHandlerFunc("0191", "talkie_box", "a4 Z80",
		[]string{"ID", "message"},
		m.handleTalkieBox)
}

// handleTalkieBox handles the talkie_box packet
// Packet format: 0191 <ID>.L <message>.80B
// Displays a talkie box message from an actor
func (m *TalkieBoxManager) handleTalkieBox(args map[string]interface{}) error {
	// Process the packet
	result := m.processTalkieBox(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.talkie_box", result)
	}

	return nil
}

// processTalkieBox processes the talkie_box packet and returns a structured result
func (m *TalkieBoxManager) processTalkieBox(args map[string]interface{}) map[string]interface{} {
	// Extract talkie box information from args
	var id uint32
	var message string

	// Extract ID
	if val, ok := args["ID"].(uint32); ok {
		id = val
	}

	// Extract message
	if val, ok := args["message"].(string); ok {
		message = val
	}

	// Get actor name
	actorName := m.getActorName(id)

	// Generate display message
	displayMessage := m.generateTalkieBoxMessage(actorName, message)

	// Return structured result
	return map[string]interface{}{
		"ID":             id,
		"actorName":      actorName,
		"message":        message,
		"displayMessage": displayMessage,
	}
}

// getActorName returns the name of an actor by ID
// In a real implementation, this would look up the actor name from a database or map
func (m *TalkieBoxManager) getActorName(actorID uint32) string {
	// This is a simplified implementation
	// In a real implementation, this would look up the actor name from a database or map
	return fmt.Sprintf("Actor_%d", actorID)
}

// generateTalkieBoxMessage generates a message for the talkie box
func (m *TalkieBoxManager) generateTalkieBoxMessage(actorName, message string) string {
	return fmt.Sprintf("%s's talkie box message: %s", actorName, message)
}
