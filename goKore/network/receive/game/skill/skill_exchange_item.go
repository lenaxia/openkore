package skill

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// SkillExchangeItemManager manages the skill_exchange_item packet handler
type SkillExchangeItemManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewSkillExchangeItemManager creates a new skill exchange item manager
func NewSkillExchangeItemManager(parser *core.CoreParser, hookManager *hooks.HookManager) *SkillExchangeItemManager {
	return &SkillExchangeItemManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// ExchangeType represents the type of skill exchange
type ExchangeType int

const (
	// ChangeMaterial represents the Change Material skill
	ChangeMaterial ExchangeType = iota
	// ElementalAnalysisLv1 represents the Elemental Analysis Lv 1 skill
	ElementalAnalysisLv1
	// ElementalAnalysisLv2 represents the Elemental Analysis Lv 2 skill
	ElementalAnalysisLv2
)

// RegisterHandlers registers all handlers related to skill exchange item
func (m *SkillExchangeItemManager) RegisterHandlers() {
	// Register skill_exchange_item handler
	m.parser.RegisterHandlerFunc("0917", "skill_exchange_item", "v V",
		[]string{"type", "val"},
		m.handleSkillExchangeItem)
}

// handleSkillExchangeItem handles the skill_exchange_item packet
// Packet format: 0917 <type>.W <val>.L
func (m *SkillExchangeItemManager) handleSkillExchangeItem(args map[string]interface{}) error {
	// Process the packet
	result := m.processSkillExchangeItem(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.skill_exchange_item", result)
	}

	return nil
}

// processSkillExchangeItem processes the skill_exchange_item packet and returns a structured result
func (m *SkillExchangeItemManager) processSkillExchangeItem(args map[string]interface{}) map[string]interface{} {
	// Extract skill exchange item information from args
	var exchangeType uint16
	var val uint32

	// Extract type
	if val, ok := args["type"].(uint16); ok {
		exchangeType = val
	}

	// Extract val
	if v, ok := args["val"].(uint32); ok {
		val = v
	}

	// Get the exchange type name and command
	typeName, command := m.getExchangeTypeInfo(exchangeType)

	// Generate message
	message := m.generateMessage(exchangeType)

	// Return structured result
	return map[string]interface{}{
		"exchangeType":     exchangeType,
		"exchangeTypeName": typeName,
		"command":          command,
		"val":              val,
		"message":          message,
	}
}

// getExchangeTypeInfo returns the exchange type name and command based on the type
func (m *SkillExchangeItemManager) getExchangeTypeInfo(exchangeType uint16) (string, string) {
	switch exchangeType {
	case uint16(ChangeMaterial):
		return "Change Material", "cm"
	case uint16(ElementalAnalysisLv1):
		return "Elemental Analysis Lv 1", "analysis"
	case uint16(ElementalAnalysisLv2):
		return "Elemental Analysis Lv 2", "analysis"
	default:
		return "Unknown", ""
	}
}

// generateMessage generates a message based on the exchange type
func (m *SkillExchangeItemManager) generateMessage(exchangeType uint16) string {
	switch exchangeType {
	case uint16(ChangeMaterial):
		return "Change Material is ready. Use command 'cm' to continue."
	case uint16(ElementalAnalysisLv1), uint16(ElementalAnalysisLv2):
		return "Four Spirit Analysis is ready. Use command 'analysis' to continue."
	default:
		return "Unknown skill exchange type."
	}
}
