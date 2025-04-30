package skill

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// SkillMsgManager manages the skill_msg packet handler
type SkillMsgManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	msgTable    map[int]string // Message table for skill messages
}

// NewSkillMsgManager creates a new skill message manager
func NewSkillMsgManager(parser *core.CoreParser, hookManager *hooks.HookManager) *SkillMsgManager {
	return &SkillMsgManager{
		parser:      parser,
		hookManager: hookManager,
		msgTable:    loadMsgTable(), // In a real implementation, this would load from msgstringtable.txt
	}
}

// loadMsgTable loads the message table from msgstringtable.txt
// This is a simplified implementation for demonstration purposes
func loadMsgTable() map[int]string {
	// In a real implementation, this would load from msgstringtable.txt
	// For now, we'll just return a few sample messages
	return map[int]string{
		1: "You have learned a new skill!",
		2: "You don't have enough SP to use this skill.",
		3: "You can't use this skill in this area.",
		4: "This skill is currently in cooldown.",
		5: "You need a weapon to use this skill.",
	}
}

// RegisterHandlers registers all handlers related to skill messages
func (m *SkillMsgManager) RegisterHandlers() {
	// Register skill_msg handler
	m.parser.RegisterHandlerFunc("0215", "skill_msg", "v v",
		[]string{"id", "msgid"},
		m.handleSkillMsg)
}

// handleSkillMsg handles the skill_msg packet
// Packet format: 0215 <skill id>.W <msg id>.W
func (m *SkillMsgManager) handleSkillMsg(args map[string]interface{}) error {
	// Process the packet
	result := m.processSkillMsg(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.skill_msg", result)
	}

	return nil
}

// processSkillMsg processes the skill_msg packet and returns a structured result
func (m *SkillMsgManager) processSkillMsg(args map[string]interface{}) map[string]interface{} {
	// Extract skill ID and message ID from args
	var skillID uint16
	var msgID uint16
	var message string

	// Extract skillID
	if val, ok := args["id"].(uint16); ok {
		skillID = val
	}

	// Extract msgID
	if val, ok := args["msgid"].(uint16); ok {
		msgID = val
	}

	// Get the message from the message table
	// In the original Perl code, the msgID is incremented before lookup
	if msg, ok := m.msgTable[int(msgID)+1]; ok {
		message = msg
	} else {
		message = fmt.Sprintf("Unknown skill message (msgid: %d, skill: %d)", msgID, skillID)
	}

	// Return structured result
	return map[string]interface{}{
		"skillID": skillID,
		"msgID":   msgID,
		"message": message,
	}
}
