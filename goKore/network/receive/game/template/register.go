// Package template provides a template for packet handler registration
package template

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/types"
)

// RegisterWithParser registers all handlers with the parser
func RegisterWithParser(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) {
	// Create the manager
	manager := NewManager(parser, hookManager, logger)

	// Register handlers
	manager.RegisterHandlers()
}

// RegisterWithReceive registers all handlers with the receive interface
func RegisterWithReceive(receive types.Receive) {
	// Register the example_handler
	receive.RegisterHandler("example_handler", func(args map[string]interface{}) error {
		// Create a manager for this specific call
		manager := NewManager(nil, nil, nil)
		return manager.handleExample(args)
	})

	// Register additional handlers as needed
}

// GetPacketDefinitions returns the packet definitions for this package
func GetPacketDefinitions() map[string]common.PacketConstruction {
	return map[string]common.PacketConstruction{
		"0000": {
			ID:         "0000",
			Name:       "example_handler",
			Format:     "v",
			FieldNames: []string{"example_field"},
		},
		// Add more packet definitions as needed
	}
}
