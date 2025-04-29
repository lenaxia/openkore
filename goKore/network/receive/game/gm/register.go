package gm

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/factory"
	"github.com/lenaxia/goKore/network/receive/types"
)

// RegisterAllHandlers registers all handlers in the gm package with the given receive component
func RegisterAllHandlers(receive types.Receive) {
	// Register the GM_silence handler
	receive.RegisterHandler("GM_silence", func(args map[string]interface{}) error {
		// Create a GM manager for this specific call
		manager := NewGMManager(nil, nil)
		return manager.HandleGMSilence(args)
	})

	// Register the GM_req_acc_name handler
	receive.RegisterHandler("GM_req_acc_name", func(args map[string]interface{}) error {
		// Create a GM manager for this specific call
		manager := NewGMManager(nil, nil)
		return manager.HandleGMReqAccName(args)
	})
}

// RegisterWithFactory registers all handlers in the gm package with the given factory
func RegisterWithFactory(receiveFactory *factory.ReceiveFactory) {
	// Currently, the factory doesn't have a direct method to register packet handlers
	// This would typically be done through server-specific packet definitions
}

// RegisterWithCoreParser registers all handlers in the gm package with the given core parser
func RegisterWithCoreParser(parser *core.CoreParser, hookManager *hooks.HookManager) {
	// Create the GM manager
	manager := NewGMManager(parser, hookManager)

	// Register handlers
	manager.RegisterAllHandlers()
}

// RegisterWithBaseReceive registers the GM manager with the base receive
// This function should be called after the BaseReceive is configured
func RegisterWithBaseReceive(baseReceive *core.BaseReceive) {
	// Register the GM_silence handler
	baseReceive.RegisterHandler("GM_silence", func(args map[string]interface{}) error {
		// Create a GM manager for this specific call
		manager := NewGMManager(nil, nil)
		return manager.HandleGMSilence(args)
	})

	// Register the GM_req_acc_name handler
	baseReceive.RegisterHandler("GM_req_acc_name", func(args map[string]interface{}) error {
		// Create a GM manager for this specific call
		manager := NewGMManager(nil, nil)
		return manager.HandleGMReqAccName(args)
	})
}

// GetPacketDefinitions returns the packet definitions for the gm package
func GetPacketDefinitions() map[string]common.PacketDef {
	return map[string]common.PacketDef{
		"0149": {
			ID:         "0149",
			Name:       "GM_silence",
			Format:     "C Z24",
			FieldNames: []string{"flag", "name"},
		},
		"01B3": {
			ID:         "01B3",
			Name:       "GM_req_acc_name",
			Format:     "V Z24",
			FieldNames: []string{"targetID", "accountName"},
		},
	}
}
