package skill

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// SkillAddManager manages the skill_add packet handler
type SkillAddManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewSkillAddManager creates a new skill add manager
func NewSkillAddManager(parser *core.CoreParser, hookManager *hooks.HookManager) *SkillAddManager {
	return &SkillAddManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to skill add
func (m *SkillAddManager) RegisterHandlers() {
	// Register skill_add handler
	m.parser.RegisterHandlerFunc("0111", "skill_add", "v V v3 C V",
		[]string{"skillID", "target", "lv", "sp", "range", "upgradable", "lv2"},
		m.handleSkillAdd)

	// Register skill_add handler (alternative packet)
	m.parser.RegisterHandlerFunc("09FE", "skill_add_new", "v V v3 C V Z24",
		[]string{"skillID", "target", "lv", "sp", "range", "upgradable", "lv2", "name"},
		m.handleSkillAdd)
}

// handleSkillAdd handles the skill_add packet
// Packet formats:
// 0111: <skill id>.W <target type>.L <level>.W <sp cost>.W <attack range>.W <upgradable>.B <level2>.W
// 09FE: <skill id>.W <target type>.L <level>.W <sp cost>.W <attack range>.W <upgradable>.B <level2>.W <skill name>.24B
func (m *SkillAddManager) handleSkillAdd(args map[string]interface{}) error {
	// Process the packet
	result := m.processSkillAdd(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.skill_add", result)
	}

	return nil
}

// processSkillAdd processes the skill_add packet and returns a structured result
func (m *SkillAddManager) processSkillAdd(args map[string]interface{}) map[string]interface{} {
	// Extract skill information from args
	var skillID uint16
	var targetType uint32
	var level, sp, attackRange uint16
	var upgradable uint8
	var level2 uint32
	var name string

	// Extract skillID
	if val, ok := args["skillID"].(uint16); ok {
		skillID = val
	}

	// Extract targetType
	if val, ok := args["target"].(uint32); ok {
		targetType = val
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
	if val, ok := args["upgradable"].(uint8); ok {
		upgradable = val
	}

	// Extract level2
	if val, ok := args["lv2"].(uint32); ok {
		level2 = val
	}

	// Extract name if available
	if val, ok := args["name"].(string); ok {
		name = val
	}

	// Create skill info
	skill := SkillInfo{
		ID:         skillID,
		TargetType: targetType,
		Level:      level,
		SP:         sp,
		Range:      attackRange,
		Handle:     name, // Will be empty for packet 0111
		Up:         upgradable,
		Level2:     uint16(level2),
	}

	// Return structured result
	return map[string]interface{}{
		"skill":     skill,
		"ownerType": OwnerChar, // Skill add is always for character
	}
}
