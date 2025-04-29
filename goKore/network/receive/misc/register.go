package misc

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/factory"
	"github.com/lenaxia/goKore/network/receive/types"
)

// RegisterAllHandlers registers all handlers in the misc package with the given receive component
func RegisterAllHandlers(receive types.Receive) {
	// Register the dynamicnpc_create_result handler
	receive.RegisterHandler("dynamicnpc_create_result", func(args map[string]interface{}) error {
		// Create a misc manager for this specific call
		manager := NewMiscManager(nil, nil)
		return manager.handleDynamicNPCCreateResult(args)
	})

	// Register GameGuard handlers with default configuration
	RegisterGameGuardHandlers(receive, DefaultGameGuardConfig())

	// Register AntiCheat handlers with default configuration
	RegisterAntiCheatHandlers(receive, DefaultAntiCheatConfig())
}

// RegisterWithFactory registers all handlers in the misc package with the given factory
func RegisterWithFactory(receiveFactory *factory.ReceiveFactory) {
	// Currently, the factory doesn't have a direct method to register packet handlers
	// This would typically be done through server-specific packet definitions
}

// RegisterWithCoreParser registers all handlers in the misc package with the given core parser
func RegisterWithCoreParser(parser *core.CoreParser, hookManager *hooks.HookManager) {
	// Create the misc manager
	manager := NewMiscManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()
}

// GetPacketDefinitions returns the packet definitions for the misc package
func GetPacketDefinitions() map[string]common.PacketDef {
	return map[string]common.PacketDef{
		"0A17": {
			ID:         "0A17",
			Name:       "dynamicnpc_create_result",
			Format:     "B",
			FieldNames: []string{"result"},
		},
		"0277": {
			ID:         "0277",
			Name:       "gameguard_request",
			Format:     "",
			FieldNames: []string{},
		},
		"02DC": {
			ID:         "02DC",
			Name:       "gameguard_grant",
			Format:     "C",
			FieldNames: []string{"server"},
		},
		"0A7F": {
			ID:         "0A7F",
			Name:       "EAC_key",
			Format:     "",
			FieldNames: []string{},
		},
	}
}

// RegisterGameGuardHandlers registers all GameGuard-related handlers with the given receive component
func RegisterGameGuardHandlers(receive types.Receive, config *GameGuardConfig) {
	// Register the gameguard_request handler
	receive.RegisterHandler("gameguard_request", func(args map[string]interface{}) error {
		// Create a GameGuard manager for this specific call
		manager := NewGameGuardManager(nil, nil, config)
		return manager.handleGameGuardRequest(args)
	})

	// Register the gameguard_grant handler
	receive.RegisterHandler("gameguard_grant", func(args map[string]interface{}) error {
		// Create a GameGuard manager for this specific call
		manager := NewGameGuardManager(nil, nil, config)
		return manager.handleGameGuardGrant(args)
	})
}

// RegisterAntiCheatHandlers registers all AntiCheat-related handlers with the given receive component
func RegisterAntiCheatHandlers(receive types.Receive, config *AntiCheatConfig) {
	// Register the EAC_key handler
	receive.RegisterHandler("EAC_key", func(args map[string]interface{}) error {
		// Create an AntiCheat manager for this specific call
		manager := NewAntiCheatManager(nil, nil, config)
		return manager.handleEACKey(args)
	})
}
