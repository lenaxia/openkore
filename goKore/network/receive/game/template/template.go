// Package template provides a template for packet handler registration
package template

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// Manager handles template-related packets
type Manager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger
}

// NewManager creates a new template manager
func NewManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *Manager {
	return &Manager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
	}
}

// RegisterHandlers registers all template-related packet handlers
func (m *Manager) RegisterHandlers() {
	// Register example handler
	if m.parser != nil {
		m.parser.RegisterHandlerFunc("0000", "example_handler", "v",
			[]string{"example_field"},
			m.handleExample)
	}

	// Register additional handlers as needed
}

// handleExample handles the example packet
func (m *Manager) handleExample(args map[string]interface{}) error {
	// Log the event
	if m.logger != nil {
		m.logger.Debug("Handling example packet")
	}

	// Process the packet
	// ...

	// Call hook if needed
	if m.hookManager != nil {
		m.hookManager.CallHook("game/template/example", map[string]interface{}{
			"example_data": "example value",
		})
	}

	// Return nil if successful
	return nil
}
