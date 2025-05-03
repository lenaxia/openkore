// Package send provides functionality for sending packets to the server.
package send

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
	"github.com/lenaxia/goKore/network/send/game/actor"
	"github.com/lenaxia/goKore/network/send/game/auction"
	"github.com/lenaxia/goKore/network/send/game/banking"
	"github.com/lenaxia/goKore/network/send/game/battle"
	"github.com/lenaxia/goKore/network/send/game/buyingstore"
	"github.com/lenaxia/goKore/network/send/game/captcha"
	"github.com/lenaxia/goKore/network/send/game/card"
	"github.com/lenaxia/goKore/network/send/game/cashshop"
	"github.com/lenaxia/goKore/network/send/game/craft"
	"github.com/lenaxia/goKore/network/send/game/deal"
	"github.com/lenaxia/goKore/network/send/game/friend"
	"github.com/lenaxia/goKore/network/send/game/gm"
	"github.com/lenaxia/goKore/network/send/game/guild"
	"github.com/lenaxia/goKore/network/send/game/homunculus"
	"github.com/lenaxia/goKore/network/send/game/login"
	"github.com/lenaxia/goKore/network/send/game/macro"
	"github.com/lenaxia/goKore/network/send/game/mail"
	"github.com/lenaxia/goKore/network/send/game/market"
	"github.com/lenaxia/goKore/network/send/game/marriage"
	"github.com/lenaxia/goKore/network/send/game/mercenary"
	"github.com/lenaxia/goKore/network/send/game/misc"
	"github.com/lenaxia/goKore/network/send/game/npc"
	"github.com/lenaxia/goKore/network/send/game/party"
	"github.com/lenaxia/goKore/network/send/game/pet"
	"github.com/lenaxia/goKore/network/send/game/ranking"
	"github.com/lenaxia/goKore/network/send/game/rodex"
	"github.com/lenaxia/goKore/network/send/game/shop"
	"github.com/lenaxia/goKore/network/send/game/skill"
	"github.com/lenaxia/goKore/network/send/game/ui"
	"github.com/lenaxia/goKore/network/send/handlers/game"
	loginHandlers "github.com/lenaxia/goKore/network/send/handlers/login"
	"github.com/lenaxia/goKore/network/send/handlers/servers"
)

// Manager is an interface for all manager types
type Manager interface{}

// HandlerRegistry manages the registration of all send packet handlers
type HandlerRegistry struct {
	baseSend    *core.BaseSend
	hookManager *hooks.HookManager
	logger      core.Logger

	// Map of managers
	managers map[string]Manager
}

// NewHandlerRegistry creates a new send handler registry
func NewHandlerRegistry(baseSend *core.BaseSend, hookManager *hooks.HookManager, logger core.Logger) *HandlerRegistry {
	return &HandlerRegistry{
		baseSend:    baseSend,
		hookManager: hookManager,
		logger:      logger,
		managers:    make(map[string]Manager),
	}
}

// RegisterAllHandlers registers all handlers with the send component
func (hr *HandlerRegistry) RegisterAllHandlers() {
	// Register login handlers
	hr.registerLoginHandlers()

	// Register game handlers
	hr.registerGameHandlers()

	// Register server-specific handlers
	hr.registerServerHandlers()

	// Register manager-based handlers
	hr.registerManagerHandlers()

	// Log registration
	hr.logger.Info("Registered all send handlers")
}

// registerLoginHandlers registers all login-related handlers
func (hr *HandlerRegistry) registerLoginHandlers() {
	// Register login handlers
	loginHandlers.RegisterHandlers(hr.baseSend)

	// Log registration
	hr.logger.Debug("Registered login send handlers")
}

// registerGameHandlers registers all game-related handlers
func (hr *HandlerRegistry) registerGameHandlers() {
	// Register handlers from the handlers/game package
	game.RegisterHandlers(hr.baseSend)

	// Register handlers for packages that have been updated to the unified pattern
	actor.RegisterWithSend(hr.baseSend, hr.hookManager, hr.logger)
	card.RegisterWithSend(hr.baseSend, hr.hookManager, hr.logger)
	login.RegisterWithSend(hr.baseSend, hr.hookManager, hr.logger)

	// Log registration
	hr.logger.Debug("Registered game send handlers")
}

// registerServerHandlers registers all server-specific handlers
func (hr *HandlerRegistry) registerServerHandlers() {
	// Register server-specific handlers based on server type
	switch hr.baseSend.GetServerType() {
	case "ServerType0":
		servers.RegisterServerType0Handlers(hr.baseSend)
	// Add more server types as needed
	default:
		hr.logger.Warning("Unknown server type: %s", hr.baseSend.GetServerType())
	}

	// Log registration
	hr.logger.Debug("Registered server-specific send handlers")
}

// registerManagerHandlers registers all manager-based handlers
func (hr *HandlerRegistry) registerManagerHandlers() {
	// Create and register managers for each game package
	hr.createManagers()

	// Log registration
	hr.logger.Debug("Registered manager-based handlers")
}

// createManagers creates instances of all manager structs
func (hr *HandlerRegistry) createManagers() {
	// Define manager creators
	managerCreators := map[string]func() Manager{
		"auction": func() Manager {
			return auction.NewAuctionManager(hr.baseSend)
		},
		"banking": func() Manager {
			return banking.NewBankingManager(hr.baseSend)
		},
		"battle": func() Manager {
			return battle.NewBattleManager(hr.baseSend)
		},
		"buyingstore": func() Manager {
			return buyingstore.NewBuyingStoreManager(hr.baseSend)
		},
		"captcha": func() Manager {
			return captcha.NewCaptchaManager(hr.baseSend)
		},
		"cashshop": func() Manager {
			return cashshop.NewCashShopManager(hr.baseSend)
		},
		"craft": func() Manager {
			return craft.NewCraftManager(hr.baseSend)
		},
		"deal": func() Manager {
			return deal.NewDealManager(hr.baseSend)
		},
		"friend": func() Manager {
			return friend.NewFriendManager(hr.baseSend)
		},
		"gm": func() Manager {
			return gm.NewGMManager(hr.baseSend)
		},
		"guild": func() Manager {
			return guild.NewGuildManager(hr.baseSend)
		},
		"homunculus": func() Manager {
			return homunculus.NewHomunculusManager(hr.baseSend)
		},
		"macro": func() Manager {
			return macro.NewMacroManager(hr.baseSend)
		},
		"mail": func() Manager {
			return mail.NewMailManager(hr.baseSend)
		},
		"market": func() Manager {
			return market.NewMarketManager(hr.baseSend)
		},
		"marriage": func() Manager {
			return marriage.NewMarriageManager(hr.baseSend)
		},
		"mercenary": func() Manager {
			return mercenary.NewMercenaryManager(hr.baseSend)
		},
		"misc": func() Manager {
			return misc.NewMiscManager(hr.baseSend)
		},
		"npc": func() Manager {
			return npc.NewNPCManager(hr.baseSend)
		},
		"party": func() Manager {
			return party.NewPartyManager(hr.baseSend)
		},
		"pet": func() Manager {
			return pet.NewPetManager(hr.baseSend)
		},
		"ranking": func() Manager {
			return ranking.NewRankingManager(hr.baseSend)
		},
		"rodex": func() Manager {
			return rodex.NewRodexManager(hr.baseSend)
		},
		"shop": func() Manager {
			return shop.NewShopManager(hr.baseSend)
		},
		"skill": func() Manager {
			return skill.NewSkillManager(hr.baseSend)
		},
		"ui": func() Manager {
			return ui.NewUIManager(hr.baseSend)
		},
	}

	// Create managers
	for name, creator := range managerCreators {
		hr.managers[name] = creator()
		hr.logger.Debug("Created %s manager", name)
	}

	// Note: The following packages are now registered through the unified pattern:
	// - Actor
	// - Card
	// - Login (game/login)
	hr.logger.Debug("Some packages are registered through the unified pattern")
}

// GetManager returns a manager by name
func (hr *HandlerRegistry) GetManager(name string) (Manager, bool) {
	manager, exists := hr.managers[name]
	return manager, exists
}

// ConfigureServerType configures the send component for a specific server type
func (hr *HandlerRegistry) ConfigureServerType(serverType string, packetConstructions map[string]common.PacketConstruction) error {
	// Configure the base send component
	err := hr.baseSend.Configure(serverType, packetConstructions)
	if err != nil {
		return err
	}

	// Register handlers for the configured server type
	hr.registerServerHandlers()

	return nil
}

// GetPacketDefinitions returns all packet definitions from all packages
func (hr *HandlerRegistry) GetPacketDefinitions() map[string]common.PacketConstruction {
	// Create a map to store all packet definitions
	packetDefs := make(map[string]common.PacketConstruction)

	// Add packet definitions from packages that have been updated to the unified pattern
	for id, def := range actor.GetPacketDefinitions() {
		packetDefs[id] = def
	}

	for id, def := range card.GetPacketDefinitions() {
		packetDefs[id] = def
	}

	for id, def := range login.GetPacketDefinitions() {
		packetDefs[id] = def
	}

	// Add more packet definitions from other packages as they are updated

	return packetDefs
}
