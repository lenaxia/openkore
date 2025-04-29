package achievement

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/factory"
	"github.com/lenaxia/goKore/network/receive/types"
)

// RegisterAllHandlers registers all handlers in the achievement package with the given receive component
func RegisterAllHandlers(receive types.Receive) {
	// Register the achievement_list handler
	receive.RegisterHandler("achievement_list", func(args map[string]interface{}) error {
		// Create an achievement manager for this specific call
		manager := NewAchievementManager(nil, nil)
		return manager.HandleAchievementList(args)
	})

	// Register the achievement_update handler
	receive.RegisterHandler("achievement_update", func(args map[string]interface{}) error {
		// Create an achievement manager for this specific call
		manager := NewAchievementManager(nil, nil)
		return manager.HandleAchievementUpdate(args)
	})

	// Register the achievement_reward_ack handler
	receive.RegisterHandler("achievement_reward_ack", func(args map[string]interface{}) error {
		// Create an achievement manager for this specific call
		manager := NewAchievementManager(nil, nil)
		return manager.HandleAchievementRewardAck(args)
	})
}

// RegisterWithFactory registers all handlers in the achievement package with the given factory
func RegisterWithFactory(receiveFactory *factory.ReceiveFactory) {
	// Currently, the factory doesn't have a direct method to register packet handlers
	// This would typically be done through server-specific packet definitions
}

// RegisterWithCoreParser registers all handlers in the achievement package with the given core parser
func RegisterWithCoreParser(parser *core.CoreParser, hookManager *hooks.HookManager) {
	// Create the achievement manager
	manager := NewAchievementManager(parser, hookManager)

	// Register handlers
	manager.RegisterAllHandlers()
}

// RegisterWithBaseReceive registers the achievement manager with the base receive
// This function should be called after the BaseReceive is configured
func RegisterWithBaseReceive(baseReceive *core.BaseReceive) {
	// Register the achievement_list handler
	baseReceive.RegisterHandler("achievement_list", func(args map[string]interface{}) error {
		// Create an achievement manager for this specific call
		manager := NewAchievementManager(nil, nil)
		return manager.HandleAchievementList(args)
	})

	// Register the achievement_update handler
	baseReceive.RegisterHandler("achievement_update", func(args map[string]interface{}) error {
		// Create an achievement manager for this specific call
		manager := NewAchievementManager(nil, nil)
		return manager.HandleAchievementUpdate(args)
	})

	// Register the achievement_reward_ack handler
	baseReceive.RegisterHandler("achievement_reward_ack", func(args map[string]interface{}) error {
		// Create an achievement manager for this specific call
		manager := NewAchievementManager(nil, nil)
		return manager.HandleAchievementRewardAck(args)
	})
}

// GetPacketDefinitions returns the packet definitions for the achievement package
func GetPacketDefinitions() map[string]common.PacketDef {
	return map[string]common.PacketDef{
		"0A23": {
			ID:         "0A23",
			Name:       "achievement_list",
			Format:     "a*",
			FieldNames: []string{"RAW_MSG"},
		},
		"0A24": {
			ID:         "0A24",
			Name:       "achievement_update",
			Format:     "V C V10 V C",
			FieldNames: []string{"achievementID", "completed", "objective1", "objective2", "objective3", "objective4", "objective5", "objective6", "objective7", "objective8", "objective9", "objective10", "completed_at", "reward"},
		},
		"0A25": {
			ID:         "0A25",
			Name:       "achievement_reward_ack",
			Format:     "V",
			FieldNames: []string{"achievementID"},
		},
	}
}
