package skill

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// ComboDelayManager manages the combo_delay packet handler
type ComboDelayManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewComboDelayManager creates a new combo delay manager
func NewComboDelayManager(parser *core.CoreParser, hookManager *hooks.HookManager) *ComboDelayManager {
	return &ComboDelayManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterHandlers registers all handlers related to combo delay
func (m *ComboDelayManager) RegisterHandlers() {
	// Register combo_delay handler
	m.parser.RegisterHandlerFunc("01D2", "combo_delay", "a4 V",
		[]string{"ID", "delay"},
		m.handleComboDelay)
}

// handleComboDelay handles the combo_delay packet
// Packet format: 01D2 <account ID>.L <delay>.L
func (m *ComboDelayManager) handleComboDelay(args map[string]interface{}) error {
	// Process the packet
	result := m.processComboDelay(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("character.combo_delay", result)
	}

	return nil
}

// processComboDelay processes the combo_delay packet and returns a structured result
func (m *ComboDelayManager) processComboDelay(args map[string]interface{}) map[string]interface{} {
	// Extract combo delay information from args
	var actorID uint32
	var delay uint32

	// Extract actorID
	if val, ok := args["ID"].(uint32); ok {
		actorID = val
	}

	// Extract delay
	if val, ok := args["delay"].(uint32); ok {
		delay = val
	}

	// Calculate combo delay in seconds
	// In the original Perl code, there's a commented formula: ($args->{delay} * 15) / 100000
	// But it's not actually used, so we'll just store the raw delay value
	comboDelay := delay

	// Return structured result
	return map[string]interface{}{
		"actorID":    actorID,
		"delay":      delay,
		"comboDelay": comboDelay,
		"isOwnChar":  false, // This would be determined by comparing with accountID in a real implementation
	}
}
