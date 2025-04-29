package storage

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/factory"
	"github.com/lenaxia/goKore/network/receive/types"
)

// RegisterAllHandlers registers all handlers in the storage package with the given receive component
func RegisterAllHandlers(receive types.Receive) {
	// Register the guild_storage_log handler
	receive.RegisterHandler("guild_storage_log", func(args map[string]interface{}) error {
		// Create a storage manager for this specific call
		manager := NewStorageManager(nil, nil)
		return manager.HandleGuildStorageLog(args)
	})
}

// RegisterWithFactory registers all handlers in the storage package with the given factory
func RegisterWithFactory(receiveFactory *factory.ReceiveFactory) {
	// Currently, the factory doesn't have a direct method to register packet handlers
	// This would typically be done through server-specific packet definitions
}

// RegisterWithCoreParser registers all handlers in the storage package with the given core parser
func RegisterWithCoreParser(parser *core.CoreParser, hookManager *hooks.HookManager) {
	// Create the storage manager
	manager := NewStorageManager(parser, hookManager)

	// Register handlers
	manager.RegisterAllHandlers()
}

// RegisterWithBaseReceive registers the storage manager with the base receive
// This function should be called after the BaseReceive is configured
func RegisterWithBaseReceive(baseReceive *core.BaseReceive) {
	// Register the guild_storage_log handler
	baseReceive.RegisterHandler("guild_storage_log", func(args map[string]interface{}) error {
		// Create a storage manager for this specific call
		manager := NewStorageManager(nil, nil)
		return manager.HandleGuildStorageLog(args)
	})
}

// GetPacketDefinitions returns the packet definitions for the storage package
func GetPacketDefinitions() map[string]common.PacketDef {
	return map[string]common.PacketDef{
		"09A6": {
			ID:         "09A6",
			Name:       "guild_storage_log",
			Format:     "v a*",
			FieldNames: []string{"result", "log"},
		},
	}
}
