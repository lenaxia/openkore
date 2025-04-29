package social

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// Handler is the main handler for social-related packets
type Handler struct {
	chatManager   *ChatManager
	partyManager  *PartyManager
	friendManager *FriendManager
	guildManager  *GuildManager
	logger        core.Logger
}

// NewHandler creates a new social handler
func NewHandler(baseParse core.Parser, hooks *hooks.HookManager, logger core.Logger) *Handler {
	return &Handler{
		chatManager:   NewChatManager(baseParse, hooks, logger),
		partyManager:  NewPartyManager(baseParse, hooks, logger),
		friendManager: NewFriendManager(baseParse, hooks, logger),
		guildManager:  NewGuildManager(baseParse, hooks, logger),
		logger:        logger,
	}
}

// RegisterHandlers registers all social-related packet handlers
func (h *Handler) RegisterHandlers() {
	// Register chat handlers
	h.chatManager.RegisterHandlers()

	// Register party handlers
	h.partyManager.RegisterHandlers()

	// Register friend handlers
	h.friendManager.RegisterHandlers()

	// Register guild handlers
	h.guildManager.RegisterHandlers()
}

// GetChatManager returns the chat manager
func (h *Handler) GetChatManager() *ChatManager {
	return h.chatManager
}

// GetPartyManager returns the party manager
func (h *Handler) GetPartyManager() *PartyManager {
	return h.partyManager
}

// GetFriendManager returns the friend manager
func (h *Handler) GetFriendManager() *FriendManager {
	return h.friendManager
}

// GetGuildManager returns the guild manager
func (h *Handler) GetGuildManager() *GuildManager {
	return h.guildManager
}
