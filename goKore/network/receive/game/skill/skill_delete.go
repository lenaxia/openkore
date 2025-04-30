package skill

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// SkillDeleteManager manages the skill_delete packet handler
type SkillDeleteManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewSkillDeleteManager creates a new skill delete manager
func NewSkillDeleteManager(parser *core.CoreParser, hookManager *hooks.HookManager) *SkillDeleteManager {
	return &SkillDeleteManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to skill delete
func (m *SkillDeleteManager) RegisterHandlers() {
	// Register skill_delete handler
	m.parser.RegisterHandlerFunc("0441", "skill_delete", "v",
		[]string{"skillID"},
		m.handleSkillDelete)
}

// handleSkillDelete handles the skill_delete packet
// Packet format: 0441 <skill id>.W
func (m *SkillDeleteManager) handleSkillDelete(args map[string]interface{}) error {
	// Process the packet
	result := m.processSkillDelete(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.skill_delete", result)
	}

	return nil
}

// processSkillDelete processes the skill_delete packet and returns a structured result
func (m *SkillDeleteManager) processSkillDelete(args map[string]interface{}) map[string]interface{} {
	// Extract skill ID from args
	var skillID uint16

	// Extract skillID
	if val, ok := args["skillID"].(uint16); ok {
		skillID = val
	}

	// Return structured result
	return map[string]interface{}{
		"skillID": skillID,
		// Note: In a real implementation, we would look up the skill name
		// from a skill database using the ID and include it in the result
	}
}
