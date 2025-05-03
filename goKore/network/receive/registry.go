// Package receive provides functionality for receiving and processing network packets.
package receive

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/game/achievement"
	"github.com/lenaxia/goKore/network/receive/game/auction"
	"github.com/lenaxia/goKore/network/receive/game/banking"
	"github.com/lenaxia/goKore/network/receive/game/battle"
	"github.com/lenaxia/goKore/network/receive/game/buyingstore"
	"github.com/lenaxia/goKore/network/receive/game/card"
	"github.com/lenaxia/goKore/network/receive/game/character_ban_list"
	"github.com/lenaxia/goKore/network/receive/game/crafting"
	"github.com/lenaxia/goKore/network/receive/game/deal"
	"github.com/lenaxia/goKore/network/receive/game/effects"
	"github.com/lenaxia/goKore/network/receive/game/field"
	"github.com/lenaxia/goKore/network/receive/game/gm"
	"github.com/lenaxia/goKore/network/receive/game/homunculus"
	"github.com/lenaxia/goKore/network/receive/game/item"
	"github.com/lenaxia/goKore/network/receive/game/mail"
	"github.com/lenaxia/goKore/network/receive/game/market"
	"github.com/lenaxia/goKore/network/receive/game/marriage"
	"github.com/lenaxia/goKore/network/receive/game/mercenary"
	"github.com/lenaxia/goKore/network/receive/game/minimap"
	"github.com/lenaxia/goKore/network/receive/game/mvp"
	"github.com/lenaxia/goKore/network/receive/game/npc"
	"github.com/lenaxia/goKore/network/receive/game/pet"
	"github.com/lenaxia/goKore/network/receive/game/quest"
	"github.com/lenaxia/goKore/network/receive/game/ranking"
	"github.com/lenaxia/goKore/network/receive/game/refining"
	"github.com/lenaxia/goKore/network/receive/game/rental"
	"github.com/lenaxia/goKore/network/receive/game/shop"
	"github.com/lenaxia/goKore/network/receive/game/skill"
	"github.com/lenaxia/goKore/network/receive/game/social"
	"github.com/lenaxia/goKore/network/receive/game/storage"
	"github.com/lenaxia/goKore/network/receive/game/ui"
	"github.com/lenaxia/goKore/network/receive/handlers/game"
	"github.com/lenaxia/goKore/network/receive/handlers/login"
	"github.com/lenaxia/goKore/network/receive/misc"
	"github.com/lenaxia/goKore/network/receive/security"
	"github.com/lenaxia/goKore/network/receive/types"
)

// Manager is an interface for all manager types
type Manager interface{}

// HandlerRegistry manages the registration of all packet handlers
type HandlerRegistry struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	receive     types.Receive
	logger      core.Logger

	// Map of managers
	managers map[string]Manager
}

// NewHandlerRegistry creates a new handler registry
func NewHandlerRegistry(parser *core.CoreParser, hookManager *hooks.HookManager, receive types.Receive, logger core.Logger) *HandlerRegistry {
	return &HandlerRegistry{
		parser:      parser,
		hookManager: hookManager,
		receive:     receive,
		logger:      logger,
		managers:    make(map[string]Manager),
	}
}

// RegisterAllHandlers registers all handlers with the receive component
func (hr *HandlerRegistry) RegisterAllHandlers() {
	// Register core handlers
	hr.registerCoreHandlers()

	// Register game handlers
	hr.registerGameHandlers()

	// Register security handlers
	hr.registerSecurityHandlers()

	// Register misc handlers
	hr.registerMiscHandlers()
}

// registerCoreHandlers registers all core handlers
func (hr *HandlerRegistry) registerCoreHandlers() {
	// Register core handlers with the parser
	core.RegisterAllHandlers(hr.parser, hr.hookManager)
}

// registerGameHandlers registers all game-related handlers
func (hr *HandlerRegistry) registerGameHandlers() {
	// Register handlers that use the receive interface
	game.RegisterActorHandlers(hr.receive)

	// Register handlers with both parser and receive interface
	hr.registerGameModuleHandlers()

	// Log registration
	hr.logger.Info("Registered game handlers")
}

// registerGameModuleHandlers registers all game module handlers
func (hr *HandlerRegistry) registerGameModuleHandlers() {
	// Register handlers for packages that have been updated to the unified pattern
	card.RegisterWithParser(hr.parser, hr.hookManager, hr.logger)
	card.RegisterWithReceive(hr.receive)

	character_ban_list.RegisterWithParser(hr.parser, hr.hookManager, hr.logger)
	character_ban_list.RegisterWithReceive(hr.receive)

	// Register handlers for packages that still use the old pattern
	// These should be updated to the unified pattern over time

	// Achievement
	achievement.RegisterWithCoreParser(hr.parser, hr.hookManager)
	achievement.RegisterAllHandlers(hr.receive)

	// Field
	field.RegisterWithCoreParser(hr.parser, hr.hookManager)
	field.RegisterAllHandlers(hr.receive)

	// GM
	gm.RegisterWithCoreParser(hr.parser, hr.hookManager)
	gm.RegisterAllHandlers(hr.receive)

	// Ranking
	ranking.RegisterWithCoreParser(hr.parser, hr.hookManager)
	ranking.RegisterAllHandlers(hr.receive)

	// Shop
	shop.RegisterWithCoreParser(hr.parser, hr.hookManager, hr.logger)
	shop.RegisterAllHandlers(hr.receive)

	// Skill
	skill.RegisterWithParser(hr.parser, hr.hookManager)

	// Storage
	storage.RegisterWithCoreParser(hr.parser, hr.hookManager)
	storage.RegisterAllHandlers(hr.receive)

	// Register handlers for other game modules
	hr.registerLegacyHandlers()
}

// registerLegacyHandlers registers handlers for modules that haven't been updated to the unified pattern
func (hr *HandlerRegistry) registerLegacyHandlers() {
	// Define manager creators
	managerCreators := map[string]func() Manager{
		"auction": func() Manager {
			manager := auction.NewAuctionManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"banking": func() Manager {
			manager := banking.NewBankingManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"battle": func() Manager {
			manager := battle.NewBattleManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"buyingstore": func() Manager {
			manager := buyingstore.NewBuyingStoreManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"crafting": func() Manager {
			manager := crafting.NewCraftingManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"deal": func() Manager {
			manager := deal.NewDealManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"effects": func() Manager {
			manager := effects.NewEffectsManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"homunculus": func() Manager {
			manager := homunculus.NewHomManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"item": func() Manager {
			handler := item.NewHandler(hr.parser, hr.hookManager, hr.logger)
			handler.RegisterHandlers()
			return handler
		},
		"mail": func() Manager {
			manager := mail.NewRodexManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"market": func() Manager {
			manager := market.NewMarketManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"marriage": func() Manager {
			manager := marriage.NewMarriageManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"mercenary": func() Manager {
			manager := mercenary.NewMercManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"minimap": func() Manager {
			manager := minimap.NewMinimapManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"mvp": func() Manager {
			manager := mvp.NewMVPManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"npc": func() Manager {
			manager := npc.NewInteractionManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"pet": func() Manager {
			manager := pet.NewPetManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"quest": func() Manager {
			manager := quest.NewQuestManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"refining": func() Manager {
			manager := refining.NewRefiningManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"rental": func() Manager {
			manager := rental.NewRentalManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
		"social": func() Manager {
			handler := social.NewHandler(hr.parser, hr.hookManager, hr.logger)
			handler.RegisterHandlers()
			return handler
		},
		"ui": func() Manager {
			manager := ui.NewUIManager(hr.parser, hr.hookManager, hr.logger)
			manager.RegisterHandlers()
			return manager
		},
	}

	// Create managers
	for name, creator := range managerCreators {
		hr.managers[name] = creator()
		hr.logger.Debug("Created and registered %s manager", name)
	}
}

// GetManager returns a manager by name
func (hr *HandlerRegistry) GetManager(name string) (Manager, bool) {
	manager, exists := hr.managers[name]
	return manager, exists
}

// registerSecurityHandlers registers all security-related handlers
func (hr *HandlerRegistry) registerSecurityHandlers() {
	// Register security handlers with the receive interface
	security.RegisterAllHandlers(hr.receive)
}

// registerMiscHandlers registers all miscellaneous handlers
func (hr *HandlerRegistry) registerMiscHandlers() {
	// Register misc handlers with the receive interface
	misc.RegisterAllHandlers(hr.receive)
}

// RegisterWithReceive registers all handlers that use the Receive interface
func (hr *HandlerRegistry) RegisterWithReceive() {
	// Register login handlers
	login.RegisterHandlers(hr.receive)

	// Register game handlers
	game.RegisterActorHandlers(hr.receive)

	// Register handlers for packages that have been updated to the unified pattern
	card.RegisterWithReceive(hr.receive)

	// Register handlers for packages that still use the old pattern
	achievement.RegisterAllHandlers(hr.receive)
	field.RegisterAllHandlers(hr.receive)
	gm.RegisterAllHandlers(hr.receive)
	ranking.RegisterAllHandlers(hr.receive)
	shop.RegisterAllHandlers(hr.receive)
	storage.RegisterAllHandlers(hr.receive)

	// Register security handlers
	security.RegisterAllHandlers(hr.receive)

	// Register misc handlers
	misc.RegisterAllHandlers(hr.receive)

	hr.logger.Info("Registered handlers with Receive interface")
}

// RegisterWithParser registers all handlers that use the Parser
func (hr *HandlerRegistry) RegisterWithParser() {
	// Register core handlers
	core.RegisterAllHandlers(hr.parser, hr.hookManager)

	// Register handlers for packages that have been updated to the unified pattern
	card.RegisterWithParser(hr.parser, hr.hookManager, hr.logger)

	// Register handlers for packages that still use the old pattern
	achievement.RegisterWithCoreParser(hr.parser, hr.hookManager)
	field.RegisterWithCoreParser(hr.parser, hr.hookManager)
	gm.RegisterWithCoreParser(hr.parser, hr.hookManager)
	ranking.RegisterWithCoreParser(hr.parser, hr.hookManager)
	shop.RegisterWithCoreParser(hr.parser, hr.hookManager, hr.logger)
	skill.RegisterWithParser(hr.parser, hr.hookManager)
	storage.RegisterWithCoreParser(hr.parser, hr.hookManager)

	// Register handlers for other game modules
	hr.registerLegacyHandlers()

	hr.logger.Info("Registered handlers with Parser")
}
