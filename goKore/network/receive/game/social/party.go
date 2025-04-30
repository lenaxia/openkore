package social

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// PartyManager handles party-related packet handling
type PartyManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger
}

// NewPartyManager creates a new party manager
func NewPartyManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *PartyManager {
	return &PartyManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
	}
}

// RegisterHandlers registers all party-related packet handlers
func (pm *PartyManager) RegisterHandlers() {
	// Register party_join handlers for different packet versions
	// 0104 - DEFAULT OLD PACKET
	pm.parser.RegisterHandlerFunc("0104", "party_join", "V C v3 Z24 Z24 Z16",
		[]string{"ID", "role", "x", "y", "type", "name", "user", "map"}, pm.HandlePartyJoin)

	// 01E9 - PACKETVER >= 2015
	pm.parser.RegisterHandlerFunc("01E9", "party_join", "V C v3 Z24 Z24 Z16 v C2",
		[]string{"ID", "role", "x", "y", "type", "name", "user", "map", "lv", "item_pickup", "item_share"}, pm.HandlePartyJoin)

	// 0A43 - PACKETVER >= 2016
	pm.parser.RegisterHandlerFunc("0A43", "party_join", "V C v v3 Z24 Z24 Z16 C2",
		[]string{"ID", "role", "jobID", "lv", "x", "y", "type", "name", "user", "map", "item_pickup", "item_share"}, pm.HandlePartyJoin)

	// 0AE4 - PACKETVER >= 2017
	pm.parser.RegisterHandlerFunc("0AE4", "party_join", "V V C v v3 Z24 Z24 Z16 C2",
		[]string{"ID", "charID", "role", "jobID", "lv", "x", "y", "type", "name", "user", "map", "item_pickup", "item_share"}, pm.HandlePartyJoin)

	// Register party_allow_invite handler
	pm.parser.RegisterHandlerFunc("02C7", "party_allow_invite", "C",
		[]string{"type"}, pm.HandlePartyAllowInvite)

	// Register party_chat handler
	pm.parser.RegisterHandlerFunc("0109", "party_chat", "v Z*",
		[]string{"len", "message"}, pm.HandlePartyChat)

	// Register party_exp handler
	pm.parser.RegisterHandlerFunc("0101", "party_exp", "C",
		[]string{"type"}, pm.HandlePartyExp)

	// Register party_exp handler (extended version)
	pm.parser.RegisterHandlerFunc("07D8", "party_exp", "C C C",
		[]string{"type", "itemPickup", "itemDivision"}, pm.HandlePartyExp)

	// Register party_leader handler
	pm.parser.RegisterHandlerFunc("07FC", "party_leader", "V V",
		[]string{"old", "new"}, pm.HandlePartyLeader)

	// Register party_hp_info handler
	pm.parser.RegisterHandlerFunc("0106", "party_hp_info", "V v2",
		[]string{"ID", "hp", "hp_max"}, pm.HandlePartyHpInfo)

	// Register party_invite handlers
	pm.parser.RegisterHandlerFunc("00FE", "party_invite", "V Z24",
		[]string{"ID", "name"}, pm.HandlePartyInvite)

	pm.parser.RegisterHandlerFunc("02C6", "party_invite", "V Z24",
		[]string{"ID", "name"}, pm.HandlePartyInvite)

	// Register other party handlers
	pm.parser.RegisterHandlerFunc("00FC", "party_invite_result", "C Z24",
		[]string{"result", "name"}, pm.HandlePartyInviteResult)

	pm.parser.RegisterHandlerFunc("01EB", "party_location", "V v2",
		[]string{"ID", "x", "y"}, pm.HandlePartyLocation)

	pm.parser.RegisterHandlerFunc("00FA", "party_organize_result", "C",
		[]string{"fail"}, pm.HandlePartyOrganizeResult)

	pm.parser.RegisterHandlerFunc("02B8", "party_show_picker", "V v3 C2 a8 v C2 a16 C",
		[]string{"sourceID", "nameID", "identified", "broken", "upgrade", "cards", "itemName"}, pm.HandlePartyShowPicker)

	// Register party_users_info handlers for different packet versions
	// 00FB - DEFAULT [OLD]
	pm.parser.RegisterHandlerFunc("00FB", "party_users_info", "Z24 x2 Z*",
		[]string{"party_name", "playerInfo"}, pm.HandlePartyUsersInfo)

	// 0A44 - PACKETVER >= 20151007
	pm.parser.RegisterHandlerFunc("0A44", "party_users_info", "Z24 x2 Z*",
		[]string{"party_name", "playerInfo"}, pm.HandlePartyUsersInfo)

	// 0AE5 - PACKETVER >= 20171207
	pm.parser.RegisterHandlerFunc("0AE5", "party_users_info", "Z24 x2 Z*",
		[]string{"party_name", "playerInfo"}, pm.HandlePartyUsersInfo)

	pm.parser.RegisterHandlerFunc("0AB2", "party_dead", "V C",
		[]string{"ID", "isDead"}, pm.HandlePartyDead)

	pm.parser.RegisterHandlerFunc("0A43", "partylv_info", "V v v",
		[]string{"ID", "job", "level"}, pm.HandlePartyLvInfo)

	pm.parser.RegisterHandlerFunc("0105", "party_leave", "V C",
		[]string{"ID", "reason"}, pm.HandlePartyLeave)
}

// HandlePartyChat handles the party_chat packet (lines 8233-8256)
func (pm *PartyManager) HandlePartyChat(args map[string]interface{}) error {
	// Extract packet data
	user, ok := args["user"].(string)
	if !ok {
		return fmt.Errorf("invalid user in party_chat packet")
	}

	message, ok := args["message"].(string)
	if !ok {
		return fmt.Errorf("invalid message in party_chat packet")
	}

	// Log the party chat message
	pm.logger.Info("[Party] %s: %s", user, message)

	// Call hooks
	pm.hookManager.CallHook("packet_partyMsg", map[string]interface{}{
		"user":    user,
		"message": message,
	})

	return nil
}

// HandlePartyExp handles the party_exp packet (lines 8258-8286)
func (pm *PartyManager) HandlePartyExp(args map[string]interface{}) error {
	// Extract packet data
	expOption, ok := args["expOption"].(uint8)
	if !ok {
		return fmt.Errorf("invalid expOption in party_exp packet")
	}

	itemOption, ok := args["itemOption"].(uint8)
	if !ok {
		return fmt.Errorf("invalid itemOption in party_exp packet")
	}

	// Process exp distribution option
	var expMsg string
	switch expOption {
	case 0:
		expMsg = "Individual Take"
	case 1:
		expMsg = "Even Share"
	default:
		expMsg = fmt.Sprintf("Unknown (%d)", expOption)
	}

	// Process item distribution option
	var itemPickupMsg, itemDivisionMsg string
	switch itemOption {
	case 0:
		itemPickupMsg = "Individual Take"
		itemDivisionMsg = "Individual Take"
	case 1:
		itemPickupMsg = "Party Take"
		itemDivisionMsg = "Even Share"
	default:
		itemPickupMsg = fmt.Sprintf("Unknown (%d)", itemOption)
		itemDivisionMsg = fmt.Sprintf("Unknown (%d)", itemOption)
	}

	// Log the party options
	pm.logger.Info("Party EXP distribution: %s", expMsg)
	pm.logger.Info("Party item pickup: %s", itemPickupMsg)
	pm.logger.Info("Party item division: %s", itemDivisionMsg)

	return nil
}

// HandlePartyLeader handles the party_leader packet (lines 8288-8299)
func (pm *PartyManager) HandlePartyLeader(args map[string]interface{}) error {
	// Extract packet data
	newLeaderID, ok := args["newLeaderID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid newLeaderID in party_leader packet")
	}

	// Get leader name (would be implemented in a real system)
	leaderName := pm.getPartyMemberName(newLeaderID)

	// Log the leadership change
	pm.logger.Info("Party leader changed to: %s", leaderName)

	return nil
}

// HandlePartyHpInfo handles the party_hp_info packet (lines 8301-8309)
func (pm *PartyManager) HandlePartyHpInfo(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in party_hp_info packet")
	}

	hp, ok := args["hp"].(uint32)
	if !ok {
		return fmt.Errorf("invalid hp in party_hp_info packet")
	}

	maxHp, ok := args["maxHp"].(uint32)
	if !ok {
		return fmt.Errorf("invalid maxHp in party_hp_info packet")
	}

	// Get member name (would be implemented in a real system)
	memberName := pm.getPartyMemberName(ID)

	// Log the HP update (debug level)
	pm.logger.Debug("Party member HP update: %s (%d/%d)", memberName, hp, maxHp)

	return nil
}

// HandlePartyInvite handles the party_invite packet (lines 8311-8322)
func (pm *PartyManager) HandlePartyInvite(args map[string]interface{}) error {
	// Extract packet data
	partyName, ok := args["partyName"].(string)
	if !ok {
		return fmt.Errorf("invalid partyName in party_invite packet")
	}

	// Log the party invitation
	pm.logger.Info("Party invitation from: %s", partyName)

	// Call hooks
	pm.hookManager.CallHook("party_invite", map[string]interface{}{
		"partyName": partyName,
	})

	return nil
}

// HandlePartyInviteResult handles the party_invite_result packet (lines 8324-8346)
func (pm *PartyManager) HandlePartyInviteResult(args map[string]interface{}) error {
	// Extract packet data
	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in party_invite_result packet")
	}

	result, ok := args["result"].(uint8)
	if !ok {
		return fmt.Errorf("invalid result in party_invite_result packet")
	}

	// Process based on result code
	switch result {
	case 0: // ANSWER_ALREADY_OTHERGROUPM
		pm.logger.Warning("%s is already in a party", name)
	case 1: // ANSWER_JOIN_REFUSE
		pm.logger.Warning("%s refused your invitation", name)
	case 2: // ANSWER_JOIN_ACCEPT
		pm.logger.Info("%s accepted your invitation", name)
	case 3: // ANSWER_MEMBER_OVERSIZE
		pm.logger.Warning("Party is full")
	case 4: // ANSWER_DUPLICATE
		pm.logger.Warning("Same account already in party")
	case 5: // ANSWER_JOINMSG_REFUSE
		pm.logger.Warning("Join request denied")
	case 6: // ANSWER_UNKNOWN_ERROR
		pm.logger.Warning("Unknown error")
	case 7: // ANSWER_UNKNOWN_CHARACTER
		pm.logger.Warning("%s is not online or does not exist", name)
	case 8: // ANSWER_INVALID_MAPPROPERTY
		pm.logger.Warning("Not allowed in this map")
	default:
		pm.logger.Warning("Unknown invite result: %d", result)
	}

	return nil
}

// HandlePartyLocation handles the party_location packet (lines 8371-8382)
func (pm *PartyManager) HandlePartyLocation(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in party_location packet")
	}

	x, ok := args["x"].(uint16)
	if !ok {
		return fmt.Errorf("invalid x in party_location packet")
	}

	y, ok := args["y"].(uint16)
	if !ok {
		return fmt.Errorf("invalid y in party_location packet")
	}

	// Get member name (would be implemented in a real system)
	memberName := pm.getPartyMemberName(ID)

	// Log the location update (debug level)
	pm.logger.Debug("Party member location update: %s (%d, %d)", memberName, x, y)

	return nil
}

// HandlePartyOrganizeResult handles the party_organize_result packet (lines 8383-8397)
func (pm *PartyManager) HandlePartyOrganizeResult(args map[string]interface{}) error {
	// Extract packet data
	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in party_organize_result packet")
	}

	// Process based on fail code
	if fail == 0 {
		pm.logger.Success("Party created successfully")
	} else {
		switch fail {
		case 1:
			pm.logger.Warning("Party name already exists")
		case 2:
			pm.logger.Warning("Already in a party")
		case 3:
			pm.logger.Warning("Not allowed in current map")
		default:
			pm.logger.Warning("Failed to organize party (error: %d)", fail)
		}
	}

	return nil
}

// HandlePartyShowPicker handles the party_show_picker packet (lines 8399-8413)
func (pm *PartyManager) HandlePartyShowPicker(args map[string]interface{}) error {
	// Extract packet data
	sourceID, ok := args["sourceID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid sourceID in party_show_picker packet")
	}

	itemName, ok := args["itemName"].(string)
	if !ok {
		return fmt.Errorf("invalid itemName in party_show_picker packet")
	}

	// Get member name (would be implemented in a real system)
	memberName := pm.getPartyMemberName(sourceID)

	// Log the item pickup
	pm.logger.Info("Party member %s picked up: %s", memberName, itemName)

	return nil
}

// HandlePartyUsersInfo handles the party_users_info packet (lines 8415-8468)
func (pm *PartyManager) HandlePartyUsersInfo(args map[string]interface{}) error {
	// Extract packet data
	partyName, ok := args["partyName"].(string)
	if !ok {
		return fmt.Errorf("invalid partyName in party_users_info packet")
	}

	members, ok := args["members"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid members in party_users_info packet")
	}

	// Log the party information
	pm.logger.Info("Party: %s (%d members)", partyName, len(members))
	for _, member := range members {
		name, _ := member["name"].(string)
		map_, _ := member["map"].(string)
		online, _ := member["online"].(uint8)

		status := "Online"
		if online == 0 {
			status = "Offline"
		}

		pm.logger.Info("- %s (%s) [%s]", name, map_, status)
	}

	// Call hooks
	pm.hookManager.CallHook("party_users_info_ready", nil)

	return nil
}

// HandlePartyDead handles the party_dead packet (lines 8469-8488)
func (pm *PartyManager) HandlePartyDead(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in party_dead packet")
	}

	isDead, ok := args["isDead"].(uint8)
	if !ok {
		return fmt.Errorf("invalid isDead in party_dead packet")
	}

	// Get member name (would be implemented in a real system)
	memberName := pm.getPartyMemberName(ID)

	// Process based on isDead flag
	if isDead == 1 {
		pm.logger.Info("Party member %s is dead", memberName)
	} else {
		pm.logger.Info("Party member %s is alive", memberName)
	}

	return nil
}

// HandlePartyLvInfo handles the partylv_info packet (lines 9428-9435)
func (pm *PartyManager) HandlePartyLvInfo(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in partylv_info packet")
	}
	_ = ID // Will be used when implementing member name lookup

	job, ok := args["job"].(uint16)
	if !ok {
		return fmt.Errorf("invalid job in partylv_info packet")
	}

	level, ok := args["level"].(uint16)
	if !ok {
		return fmt.Errorf("invalid level in partylv_info packet")
	}

	// Get member name (would be implemented in a real system)
	memberName := pm.getPartyMemberName(ID)

	// Log the level update (debug level)
	pm.logger.Debug("Party member level update: %s (Job: %d, Level: %d)", memberName, job, level)

	return nil
}

// HandlePartyJoin handles the party_join packet (lines 8168-8220)
func (pm *PartyManager) HandlePartyJoin(args map[string]interface{}) error {
	// Extract packet data
	partyName, ok := args["partyName"].(string)
	if !ok {
		return fmt.Errorf("invalid partyName in party_join packet")
	}

	ID, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in party_join packet")
	}
	// Use ID for debugging
	pm.logger.Debug("Party join with ID: %d", ID)

	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in party_join packet")
	}

	// Check if this is the player joining or another member
	isPlayer := args["isPlayer"].(bool)

	// Log the party join
	if isPlayer {
		pm.logger.Info("You joined party: %s", partyName)
	} else {
		pm.logger.Info("%s joined the party", name)
	}

	// Call hooks
	if isPlayer {
		pm.hookManager.CallHook("packet_partyJoin", map[string]interface{}{
			"partyName": partyName,
		})
	}

	return nil
}

// HandlePartyAllowInvite handles the party_allow_invite packet (lines 8221-8231)
func (pm *PartyManager) HandlePartyAllowInvite(args map[string]interface{}) error {
	// Extract packet data
	type_, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in party_allow_invite packet")
	}

	// Process based on type
	if type_ == 0 {
		pm.logger.Info("Allowed other players to invite to party")
	} else {
		pm.logger.Info("Not allowed other players to invite to party")
	}

	return nil
}

// HandlePartyLeave handles the party_leave packet (lines 8348-8369)
func (pm *PartyManager) HandlePartyLeave(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in party_leave packet")
	}
	// Use ID for debugging
	pm.logger.Debug("Party leave with ID: %d", ID)

	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in party_leave packet")
	}

	reason, ok := args["reason"].(uint8)
	if !ok {
		return fmt.Errorf("invalid reason in party_leave packet")
	}

	// Check if this is the player leaving or another member
	isPlayer := args["isPlayer"].(bool)

	// Process based on reason code
	var reasonStr string
	switch reason {
	case 0: // GROUPMEMBER_DELETE_LEAVE
		reasonStr = "left"
	case 1: // GROUPMEMBER_DELETE_EXPEL
		reasonStr = "kicked"
	default:
		reasonStr = fmt.Sprintf("left (reason: %d)", reason)
	}

	// Log the party leave
	if isPlayer {
		pm.logger.Info("You %s the party", reasonStr)
	} else {
		pm.logger.Info("%s has %s the party", name, reasonStr)
	}

	return nil
}

// Helper function to get party member name from ID
func (pm *PartyManager) getPartyMemberName(ID uint32) string {
	// TODO: Implement party member name lookup
	// This would be handled by the game state manager
	return fmt.Sprintf("Member#%d", ID)
}
