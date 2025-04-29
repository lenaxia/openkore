// Package misc provides handlers for miscellaneous packets that don't fit into other categories.
package misc

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// AntiCheatConfig holds configuration for anti-cheat systems
type AntiCheatConfig struct {
	// Whether to ignore anti-cheat warnings
	IgnoreAntiCheatWarning bool
}

// DefaultAntiCheatConfig returns a default anti-cheat configuration
func DefaultAntiCheatConfig() *AntiCheatConfig {
	return &AntiCheatConfig{
		IgnoreAntiCheatWarning: false,
	}
}

// AntiCheatManager manages anti-cheat related packet handlers
type AntiCheatManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	config      *AntiCheatConfig
}

// NewAntiCheatManager creates a new anti-cheat manager
func NewAntiCheatManager(parser *core.CoreParser, hookManager *hooks.HookManager, config *AntiCheatConfig) *AntiCheatManager {
	if config == nil {
		config = DefaultAntiCheatConfig()
	}

	return &AntiCheatManager{
		parser:      parser,
		hookManager: hookManager,
		config:      config,
	}
}

// RegisterHandlers registers all handlers related to anti-cheat
func (m *AntiCheatManager) RegisterHandlers() {
	// Register EAC_key handler
	// Note: The actual packet ID may need to be updated based on the server implementation
	m.parser.RegisterHandlerFunc("0A7F", "EAC_key", "",
		[]string{},
		m.handleEACKey)
}

// handleEACKey handles the EAC_key packet
// Packet format: 0A7F
func (m *AntiCheatManager) handleEACKey(args map[string]interface{}) error {
	// Process the packet
	result := m.processEACKey(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("misc.eac_key", result)
	}

	return nil
}

// processEACKey processes the EAC_key packet and returns a structured result
func (m *AntiCheatManager) processEACKey(args map[string]interface{}) map[string]interface{} {
	// Check if we should ignore anti-cheat warnings
	if m.config.IgnoreAntiCheatWarning {
		return map[string]interface{}{
			"ignored": true,
			"message": "Easy Anti-Cheat warning ignored",
			"quit":    false,
		}
	}

	// In the original implementation, this would log a message and quit
	// We'll use the hook system to handle this

	// Return structured result
	return map[string]interface{}{
		"ignored": false,
		"message": "*** Easy Anti-Cheat Detected ***\nOpenKore doesn't have support for servers with Easy Anti-Cheat Shield, please read the FAQ (github).",
		"quit":    true,
	}
}

// UpdateConfig updates the anti-cheat configuration
func (m *AntiCheatManager) UpdateConfig(config *AntiCheatConfig) {
	if config != nil {
		m.config = config
	}
}

// GetConfig returns the current anti-cheat configuration
func (m *AntiCheatManager) GetConfig() *AntiCheatConfig {
	return m.config
}
