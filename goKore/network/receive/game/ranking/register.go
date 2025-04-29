package ranking

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/factory"
	"github.com/lenaxia/goKore/network/receive/types"
)

// RegisterAllHandlers registers all handlers in the ranking package with the given receive component
func RegisterAllHandlers(receive types.Receive) {
	// Register the pvp_rank handler
	receive.RegisterHandler("pvp_rank", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handlePvpRank(args)
	})

	// Register the taekwon_rank handler
	receive.RegisterHandler("taekwon_rank", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTaekwonRank(args)
	})

	// Register the taekwon_packets handler
	receive.RegisterHandler("taekwon_packets", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTaekwonPackets(args)
	})

	// Register the top10_taekwon_rank handler
	receive.RegisterHandler("top10_taekwon_rank", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTop10TaekwonRank(args)
	})

	// Register the top10_pk_rank handler
	receive.RegisterHandler("top10_pk_rank", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTop10PkRank(args)
	})

	// Register the top10_blacksmith_rank handler
	receive.RegisterHandler("top10_blacksmith_rank", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTop10BlacksmithRank(args)
	})

	// Register the top10_alchemist_rank handler
	receive.RegisterHandler("top10_alchemist_rank", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTop10AlchemistRank(args)
	})

	// Register the top10 handler
	receive.RegisterHandler("top10", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTop10(args)
	})
}

// RegisterWithFactory registers all handlers in the ranking package with the given factory
func RegisterWithFactory(receiveFactory *factory.ReceiveFactory) {
	// Currently, the factory doesn't have a direct method to register packet handlers
	// This would typically be done through server-specific packet definitions
}

// RegisterWithCoreParser registers all handlers in the ranking package with the given core parser
func RegisterWithCoreParser(parser *core.CoreParser, hookManager *hooks.HookManager) {
	// Create the ranking manager
	manager := NewRankingManager(parser, hookManager)

	// Register handlers
	manager.RegisterAllHandlers()
}

// RegisterWithBaseReceive registers the ranking manager with the base receive
// This function should be called after the BaseReceive is configured
func RegisterWithBaseReceive(baseReceive *core.BaseReceive) {
	// Register the pvp_rank handler
	baseReceive.RegisterHandler("pvp_rank", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handlePvpRank(args)
	})

	// Register the taekwon_rank handler
	baseReceive.RegisterHandler("taekwon_rank", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTaekwonRank(args)
	})

	// Register the taekwon_packets handler
	baseReceive.RegisterHandler("taekwon_packets", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTaekwonPackets(args)
	})

	// Register the top10_taekwon_rank handler
	baseReceive.RegisterHandler("top10_taekwon_rank", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTop10TaekwonRank(args)
	})

	// Register the top10_pk_rank handler
	baseReceive.RegisterHandler("top10_pk_rank", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTop10PkRank(args)
	})

	// Register the top10_blacksmith_rank handler
	baseReceive.RegisterHandler("top10_blacksmith_rank", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTop10BlacksmithRank(args)
	})

	// Register the top10_alchemist_rank handler
	baseReceive.RegisterHandler("top10_alchemist_rank", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTop10AlchemistRank(args)
	})

	// Register the top10 handler
	baseReceive.RegisterHandler("top10", func(args map[string]interface{}) error {
		// Create a ranking manager for this specific call
		manager := NewRankingManager(nil, nil)
		return manager.handleTop10(args)
	})
}

// GetPacketDefinitions returns the packet definitions for the ranking package
func GetPacketDefinitions() map[string]common.PacketDef {
	return map[string]common.PacketDef{
		"09A1": {
			ID:         "09A1",
			Name:       "pvp_rank",
			Format:     "V v2",
			FieldNames: []string{"ID", "rank", "num"},
		},
		"0224": {
			ID:         "0224",
			Name:       "taekwon_rank",
			Format:     "V",
			FieldNames: []string{"rank"},
		},
		"0226": {
			ID:         "0226",
			Name:       "taekwon_packets",
			Format:     "C C Z24",
			FieldNames: []string{"flag", "value", "name"},
		},
		"0223": {
			ID:         "0223",
			Name:       "top10_taekwon_rank",
			Format:     "a*",
			FieldNames: []string{"RAW_MSG"},
		},
		"0238": {
			ID:         "0238",
			Name:       "top10_pk_rank",
			Format:     "a*",
			FieldNames: []string{"RAW_MSG"},
		},
		"0219": {
			ID:         "0219",
			Name:       "top10_blacksmith_rank",
			Format:     "a*",
			FieldNames: []string{"RAW_MSG"},
		},
		"021A": {
			ID:         "021A",
			Name:       "top10_alchemist_rank",
			Format:     "a*",
			FieldNames: []string{"RAW_MSG"},
		},
		"097D": {
			ID:         "097D",
			Name:       "top10",
			Format:     "W a*",
			FieldNames: []string{"type", "RAW_MSG"},
		},
	}
}
