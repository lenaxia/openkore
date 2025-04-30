package skill

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// SkillUpdateManager manages the skill_update packet handler
type SkillUpdateManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewSkillUpdateManager creates a new skill update manager
func NewSkillUpdateManager(parser *core.CoreParser, hookManager *hooks.HookManager) *SkillUpdateManager {
	return &SkillUpdateManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to skill update
func (m *SkillUpdateManager) RegisterHandlers() {
	// Register skill_update handler
	m.parser.RegisterHandlerFunc("0110", "skill_update", "v v V v3 C V",
		[]string{"skillID", "lv", "sp", "range", "up", "lv2"},
		m.handleSkillUpdate)
}

// handleSkillUpdate handles the skill_update packet
// Packet format: 0110 <skill id>.W <level>.W <sp cost>.W <attack range>.W <upgradable>.B <level2>.W
func (m *SkillUpdateManager) handleSkillUpdate(args map[string]interface{}) error {
	// Process the packet
	result := m.processSkillUpdate(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.skill_update", result)
	}

	return nil
}

// processSkillUpdate processes the skill_update packet and returns a structured result
func (m *SkillUpdateManager) processSkillUpdate(args map[string]interface{}) map[string]interface{} {
	// Extract skill information from args
	var skillID uint16
	var level, sp, attackRange uint16
	var upgradable uint8
	var level2 uint32

	// Extract skillID
	if val, ok := args["skillID"].(uint16); ok {
		skillID = val
	}

	// Extract level
	if val, ok := args["lv"].(uint16); ok {
		level = val
	}

	// Extract sp cost
	if val, ok := args["sp"].(uint16); ok {
		sp = val
	}

	// Extract attack range
	if val, ok := args["range"].(uint16); ok {
		attackRange = val
	}

	// Extract upgradable flag
	if val, ok := args["up"].(uint8); ok {
		upgradable = val
	}

	// Extract level2
	if val, ok := args["lv2"].(uint32); ok {
		level2 = val
	}

	// Create skill info
	skill := SkillInfo{
		ID:     skillID,
		Level:  level,
		SP:     sp,
		Range:  attackRange,
		Up:     upgradable,
		Level2: uint16(level2),
		// Note: We don't have the handle (name) here, but in a real implementation
		// we would look it up from a skill database using the ID
	}

	// Return structured result
	return map[string]interface{}{
		"skill":     skill,
		"ownerType": OwnerChar, // Skill update is always for character
	}
}
