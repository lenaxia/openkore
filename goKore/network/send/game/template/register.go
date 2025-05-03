// Package template provides a template for packet handler registration
package template

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
)

// RegisterWithSend registers all handlers with the send component
func RegisterWithSend(send *core.BaseSend, hookManager *hooks.HookManager, logger core.Logger) {
	// Create the manager
	manager := NewManager(send, hookManager, logger)

	// Register handlers
	manager.RegisterHandlers()

	// Log registration if logger is provided
	if logger != nil {
		logger.Debug("Registered template handlers with send")
	}
}

// GetPacketDefinitions returns the packet definitions for this package
func GetPacketDefinitions() map[string]common.PacketConstruction {
	return map[string]common.PacketConstruction{
		"0000": {
			Name:       "example_handler",
			Format:     "v",
			FieldNames: []string{"example_field"},
		},
		// Add more packet definitions as needed
	}
}
