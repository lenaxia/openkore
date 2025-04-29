package card

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/factory"
	"github.com/lenaxia/goKore/network/receive/types"
)

// RegisterAllHandlers registers all handlers in the card package with the given receive component
func RegisterAllHandlers(receive types.Receive) {
	// Register the card_merge_list handler
	receive.RegisterHandler("card_merge_list", func(args map[string]interface{}) error {
		// Create a card manager for this specific call
		manager := NewCardManager(nil, nil)
		return manager.handleCardMergeList(args)
	})

	// Register the card_merge_status handler
	receive.RegisterHandler("card_merge_status", func(args map[string]interface{}) error {
		// Create a card manager for this specific call
		manager := NewCardManager(nil, nil)
		return manager.handleCardMergeStatus(args)
	})
}

// RegisterWithFactory registers all handlers in the card package with the given factory
func RegisterWithFactory(receiveFactory *factory.ReceiveFactory) {
	// Currently, the factory doesn't have a direct method to register packet handlers
	// This would typically be done through server-specific packet definitions
}

// RegisterWithCoreParser registers all handlers in the card package with the given core parser
func RegisterWithCoreParser(parser *core.CoreParser, hookManager *hooks.HookManager) {
	// Create the card manager
	manager := NewCardManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()
}

// GetPacketDefinitions returns the packet definitions for the card package
func GetPacketDefinitions() map[string]common.PacketDef {
	return map[string]common.PacketDef{
		"0A10": {
			ID:         "0A10",
			Name:       "card_merge_list",
			Format:     "v a*",
			FieldNames: []string{"len", "item_list"},
		},
		"0A11": {
			ID:         "0A11",
			Name:       "card_merge_status",
			Format:     "B B B",
			FieldNames: []string{"fail", "item_index", "card_index"},
		},
	}
}
