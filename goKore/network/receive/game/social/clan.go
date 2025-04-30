package social

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// ClanManager handles clan-related packet handling
type ClanManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger
}

// NewClanManager creates a new clan manager
func NewClanManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *ClanManager {
	return &ClanManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
	}
}

// RegisterHandlers registers all clan-related packet handlers
func (cm *ClanManager) RegisterHandlers() {
	// Register clan_user handler
	cm.parser.RegisterHandlerFunc("0988", "clan_user", "v2",
		[]string{"onlineuser", "totalmembers"}, cm.HandleClanUser)

	// Register clan_info handler
	cm.parser.RegisterHandlerFunc("0989", "clan_info", "V Z24 Z24 Z16 v2 Z*",
		[]string{"clan_ID", "clan_name", "clan_master", "clan_map", "alliance_count", "antagonist_count", "ally_antagonist_names"}, cm.HandleClanInfo)

	// Register clan_chat handler
	cm.parser.RegisterHandlerFunc("098A", "clan_chat", "v Z24 Z*",
		[]string{"len", "charname", "message"}, cm.HandleClanChat)

	// Register clan_leave handler
	cm.parser.RegisterHandlerFunc("098D", "clan_leave", "",
		[]string{}, cm.HandleClanLeave)
}

// HandleClanUser handles the clan_user packet (lines 8904-8911)
func (cm *ClanManager) HandleClanUser(args map[string]interface{}) error {
	// Extract packet data
	onlineUser, ok := args["onlineuser"].(uint16)
	if !ok {
		return fmt.Errorf("invalid onlineuser in clan_user packet")
	}

	totalMembers, ok := args["totalmembers"].(uint16)
	if !ok {
		return fmt.Errorf("invalid totalmembers in clan_user packet")
	}

	// Log clan user info
	cm.logger.Info("Clan online users: %d, Total members: %d", onlineUser, totalMembers)

	return nil
}

// HandleClanInfo handles the clan_info packet (lines 8913-8942)
func (cm *ClanManager) HandleClanInfo(args map[string]interface{}) error {
	// Extract packet data
	clanID, ok := args["clan_ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid clan_ID in clan_info packet")
	}

	clanName, ok := args["clan_name"].(string)
	if !ok {
		return fmt.Errorf("invalid clan_name in clan_info packet")
	}

	clanMaster, ok := args["clan_master"].(string)
	if !ok {
		return fmt.Errorf("invalid clan_master in clan_info packet")
	}

	clanMap, ok := args["clan_map"].(string)
	if !ok {
		return fmt.Errorf("invalid clan_map in clan_info packet")
	}

	allianceCount, ok := args["alliance_count"].(uint16)
	if !ok {
		return fmt.Errorf("invalid alliance_count in clan_info packet")
	}

	antagonistCount, ok := args["antagonist_count"].(uint16)
	if !ok {
		return fmt.Errorf("invalid antagonist_count in clan_info packet")
	}

	// Log clan info
	cm.logger.Info("Clan Info: %s (ID: %d, Master: %s, Map: %s)", clanName, clanID, clanMaster, clanMap)
	cm.logger.Info("Clan has %d allies and %d antagonists", allianceCount, antagonistCount)

	return nil
}

// HandleClanChat handles the clan_chat packet (lines 8944-8965)
func (cm *ClanManager) HandleClanChat(args map[string]interface{}) error {
	// Extract packet data
	charName, ok := args["charname"].(string)
	if !ok {
		return fmt.Errorf("invalid charname in clan_chat packet")
	}

	message, ok := args["message"].(string)
	if !ok {
		return fmt.Errorf("invalid message in clan_chat packet")
	}

	// Log clan chat message
	cm.logger.Info("[Clan] %s: %s", charName, message)

	// Call hook for clan message
	cm.hookManager.CallHook("packet_clanMsg", map[string]interface{}{
		"MsgUser": charName,
		"Msg":     message,
		"RawMsg":  message,
	})

	return nil
}

// HandleClanLeave handles the clan_leave packet (lines 8967-8974)
func (cm *ClanManager) HandleClanLeave(args map[string]interface{}) error {
	// Log clan leave
	cm.logger.Info("You left the clan")

	return nil
}
