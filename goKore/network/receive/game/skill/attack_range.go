package skill

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// AttackRangeManager manages the attack_range packet handler
type AttackRangeManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewAttackRangeManager creates a new attack range manager
func NewAttackRangeManager(parser *core.CoreParser, hookManager *hooks.HookManager) *AttackRangeManager {
	return &AttackRangeManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to attack range
func (m *AttackRangeManager) RegisterHandlers() {
	// Register attack_range handler
	m.parser.RegisterHandlerFunc("013A", "attack_range", "v",
		[]string{"type"},
		m.handleAttackRange)
}

// handleAttackRange handles the attack_range packet
// Packet format: 013A <atk range>.W
func (m *AttackRangeManager) handleAttackRange(args map[string]interface{}) error {
	// Process the packet
	result := m.processAttackRange(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.attack_range", result)
	}

	return nil
}

// processAttackRange processes the attack_range packet and returns a structured result
func (m *AttackRangeManager) processAttackRange(args map[string]interface{}) map[string]interface{} {
	// Extract attack range information from args
	var attackRange uint16

	// Extract attack range
	if val, ok := args["type"].(uint16); ok {
		attackRange = val
	}

	// Return structured result
	return map[string]interface{}{
		"attackRange": attackRange,
		"configUpdated": map[string]interface{}{
			"attackDistance":    attackRange,
			"attackMaxDistance": attackRange,
		},
	}
}
