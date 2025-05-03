// Package actor provides functionality for character-related packets
package actor

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
)

// RegisterWithSend registers all handlers with the send component
func RegisterWithSend(send core.BaseSend, hookManager *hooks.HookManager, logger core.Logger) {
	// Create the manager
	manager := NewManager(send, hookManager, logger)

	// Register handlers
	manager.RegisterHandlers()

	// Log registration if logger is provided
	if logger != nil {
		logger.Debug("Registered actor handlers with send")
	}
}

// GetPacketDefinitions returns the packet definitions for this package
func GetPacketDefinitions() map[string]common.PacketConstruction {
	return map[string]common.PacketConstruction{
		"0067": {
			Name:       "char_create",
			Format:     "a24 C7 v2",
			FieldNames: []string{"name", "str", "agi", "vit", "int", "dex", "luk", "slot", "hair_color", "hair_style"},
		},
		"0068": {
			Name:       "char_delete",
			Format:     "a4 a40",
			FieldNames: []string{"charID", "email"},
		},
		// Add more packet definitions as needed
	}
}
