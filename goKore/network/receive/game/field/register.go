package field

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/factory"
	"github.com/lenaxia/goKore/network/receive/types"
)

// RegisterAllHandlers registers all handlers in the field package with the given receive component
func RegisterAllHandlers(receive types.Receive) {
	// Register the private_airship_type handler
	receive.RegisterHandler("private_airship_type", func(args map[string]interface{}) error {
		// Create a field manager for this specific call
		manager := NewFieldManager(nil, nil)
		return manager.handlePrivateAirshipType(args)
	})
}

// RegisterWithFactory registers all handlers in the field package with the given factory
func RegisterWithFactory(receiveFactory *factory.ReceiveFactory) {
	// Currently, the factory doesn't have a direct method to register packet handlers
	// This would typically be done through server-specific packet definitions
}

// RegisterWithCoreParser registers all handlers in the field package with the given core parser
func RegisterWithCoreParser(parser *core.CoreParser, hookManager *hooks.HookManager) {
	// Create the field manager
	manager := NewFieldManager(parser, hookManager)

	// Register handlers
	manager.RegisterTransportHandlers()
}

// GetPacketDefinitions returns the packet definitions for the field package
func GetPacketDefinitions() map[string]common.PacketDef {
	return map[string]common.PacketDef{
		"0A4B": {
			ID:         "0A4B",
			Name:       "private_airship_type",
			Format:     "B",
			FieldNames: []string{"fail"},
		},
		"0293": {
			ID:         "0293",
			Name:       "boss_map_info",
			Format:     "B Z24 v2 B2",
			FieldNames: []string{"flag", "name", "x", "y", "hours", "minutes"},
		},
		"08E2": {
			ID:         "08E2",
			Name:       "navigate_to",
			Format:     "B3 Z16 v3",
			FieldNames: []string{"type", "flag", "hide", "map", "x", "y", "mob_id"},
		},
		"0B1B": {
			ID:         "0B1B",
			Name:       "warp_portal_list",
			Format:     "B Z16 Z16 Z16 Z16",
			FieldNames: []string{"type", "memo1", "memo2", "memo3", "memo4"},
		},
	}
}
