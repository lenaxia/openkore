package social

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// GuildManager handles guild-related packet handling
type GuildManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger
}

// NewGuildManager creates a new guild manager
func NewGuildManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *GuildManager {
	return &GuildManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
	}
}

// RegisterHandlers registers all guild-related packet handlers
func (gm *GuildManager) RegisterHandlers() {
	// Guild members list handlers for different packet versions
	// 0154 - DEFAULT OLD PACKET
	gm.parser.RegisterHandlerFunc("0154", "guild_members_list", "v Z*",
		[]string{"len", "member_list"}, gm.HandleGuildMembersList)

	// 0AA5 - PACKETVER >= 20151001
	gm.parser.RegisterHandlerFunc("0AA5", "guild_members_list", "v Z*",
		[]string{"len", "member_list"}, gm.HandleGuildMembersList)

	// 0B7D - PACKETVER >= 20170502
	gm.parser.RegisterHandlerFunc("0B7D", "guild_members_list", "v Z*",
		[]string{"len", "member_list"}, gm.HandleGuildMembersList)

	// Guild name handler
	gm.parser.RegisterHandlerFunc("016C", "guild_name", "V V V C V Z24",
		[]string{"guildID", "emblemID", "mode", "ismaster", "interSID", "guildName"}, gm.HandleGuildName)

	// Guild member online status handlers
	gm.parser.RegisterHandlerFunc("016D", "guild_member_online_status", "V V V",
		[]string{"ID", "charID", "online"}, gm.HandleGuildMemberOnlineStatus)

	// Extended version with gender, hair style, hair color
	gm.parser.RegisterHandlerFunc("01F2", "guild_member_online_status", "V V V v2",
		[]string{"ID", "charID", "online", "gender", "hair_style", "hair_color"}, gm.HandleGuildMemberOnlineStatus)

	// Guild notice handler
	gm.parser.RegisterHandlerFunc("016F", "guild_notice", "Z60 Z120",
		[]string{"subject", "notice"}, gm.HandleGuildNotice)

	// Guild allies/enemy list handler
	gm.parser.RegisterHandlerFunc("014C", "guild_allies_enemy_list", "v Z*",
		[]string{"len", "RAW_MSG"}, gm.HandleGuildAlliesEnemyList)

	// Guild ally request handler
	gm.parser.RegisterHandlerFunc("0171", "guild_ally_request", "V Z24",
		[]string{"ID", "guildName"}, gm.HandleGuildAllyRequest)

	// Guild broken handler
	gm.parser.RegisterHandlerFunc("015E", "guild_broken", "V",
		[]string{"flag"}, gm.HandleGuildBroken)

	// Guild create result handler
	gm.parser.RegisterHandlerFunc("0167", "guild_create_result", "C",
		[]string{"type"}, gm.HandleGuildCreateResult)

	// Guild info handler
	gm.parser.RegisterHandlerFunc("01B6", "guild_info", "V v V V V Z24 Z24 Z16 V",
		[]string{"guildID", "level", "connect_member", "max_member", "average_lv", "exp", "next_exp", "name", "master", "master_connect"}, gm.HandleGuildInfo)

	// Guild invite result handler
	gm.parser.RegisterHandlerFunc("0169", "guild_invite_result", "C",
		[]string{"type"}, gm.HandleGuildInviteResult)

	// Guild location handler
	gm.parser.RegisterHandlerFunc("01EB", "guild_location", "V v2",
		[]string{"ID", "x", "y"}, gm.HandleGuildLocation)

	// Guild leave handlers
	gm.parser.RegisterHandlerFunc("015A", "guild_leave", "Z24 Z40",
		[]string{"name", "message"}, gm.HandleGuildLeave)

	// Extended version
	gm.parser.RegisterHandlerFunc("0A83", "guild_leave", "V Z40",
		[]string{"charID", "message"}, gm.HandleGuildLeave)

	// Guild expulsion handlers
	gm.parser.RegisterHandlerFunc("015C", "guild_expulsion", "Z24 Z40 Z24",
		[]string{"name", "message", "acc"}, gm.HandleGuildExpulsion)

	// Simplified version
	gm.parser.RegisterHandlerFunc("0839", "guild_expulsion", "Z24 Z40",
		[]string{"name", "message"}, gm.HandleGuildExpulsion)

	// Extended version
	gm.parser.RegisterHandlerFunc("0A82", "guild_expulsion", "V Z40",
		[]string{"charID", "message"}, gm.HandleGuildExpulsion)

	// Guild update member position handler
	gm.parser.RegisterHandlerFunc("0156", "guild_update_member_position", "v Z*",
		[]string{"len", "member_list"}, gm.HandleGuildUpdateMemberPosition)

	// Guild members title list handler
	gm.parser.RegisterHandlerFunc("0166", "guild_members_title_list", "v Z*",
		[]string{"len", "RAW_MSG"}, gm.HandleGuildMembersTitleList)

	// Guild request handler
	gm.parser.RegisterHandlerFunc("016A", "guild_request", "V Z24",
		[]string{"ID", "name"}, gm.HandleGuildRequest)

	// Guild master member handler
	gm.parser.RegisterHandlerFunc("014E", "guild_master_member", "V",
		[]string{"type"}, gm.HandleGuildMasterMember)

	// Guild alliance handler
	gm.parser.RegisterHandlerFunc("0173", "guild_alliance", "C",
		[]string{"flag"}, gm.HandleGuildAlliance)

	// Guild member setting list handler
	gm.parser.RegisterHandlerFunc("0160", "guild_member_setting_list", "v Z*",
		[]string{"len", "RAW_MSG"}, gm.HandleGuildMemberSettingList)

	// Guild skills list handler
	gm.parser.RegisterHandlerFunc("0162", "guild_skills_list", "v v Z*",
		[]string{"len", "skillPoints", "RAW_MSG"}, gm.HandleGuildSkillsList)

	// Guild expulsion list handlers
	gm.parser.RegisterHandlerFunc("0163", "guild_expulsion_list", "v Z*",
		[]string{"len", "expulsion_list"}, gm.HandleGuildExpulsionList)

	// Extended version
	gm.parser.RegisterHandlerFunc("0B7C", "guild_expulsion_list", "v Z*",
		[]string{"len", "expulsion_list"}, gm.HandleGuildExpulsionList)

	// Guild member map change handler
	gm.parser.RegisterHandlerFunc("01EC", "guild_member_map_change", "V V V Z16",
		[]string{"ID", "charID", "status", "mapName"}, gm.HandleGuildMemberMapChange)

	// Guild member add handlers
	gm.parser.RegisterHandlerFunc("0182", "guild_member_add", "V V v5 V3 Z50 Z24",
		[]string{"ID", "charID", "hair_style", "hair_color", "sex", "jobID", "lv", "contribution", "online", "position", "memo", "name"}, gm.HandleGuildMemberAdd)

	// Extended version
	gm.parser.RegisterHandlerFunc("0B7E", "guild_member_add", "V V v5 V4 Z24",
		[]string{"ID", "charID", "hair_style", "hair_color", "sex", "jobID", "lv", "contribution", "online", "position", "lastLoginTime", "name"}, gm.HandleGuildMemberAdd)

	// Guild emblem handlers
	gm.parser.RegisterHandlerFunc("0152", "guild_emblem", "V V v",
		[]string{"guildID", "emblemID", "emblemVersion"}, gm.HandleGuildEmblem)

	// Guild emblem update handler
	gm.parser.RegisterHandlerFunc("01B4", "guild_emblem_update", "V V v",
		[]string{"guildID", "emblemID", "emblemVersion"}, gm.HandleGuildEmblemUpdate)

	// Guild position changed handler
	gm.parser.RegisterHandlerFunc("0174", "guild_position_changed", "v Z*",
		[]string{"len", "RAW_MSG"}, gm.HandleGuildPositionChanged)

	// Guild position handler
	gm.parser.RegisterHandlerFunc("0AFD", "guild_position", "v Z*",
		[]string{"len", "RAW_MSG"}, gm.HandleGuildPosition)

	// Guild unally handler
	gm.parser.RegisterHandlerFunc("0184", "guild_unally", "V Z24",
		[]string{"guildID", "guildName"}, gm.HandleGuildUnally)

	// Guild opposition result handler
	gm.parser.RegisterHandlerFunc("0181", "guild_opposition_result", "C",
		[]string{"flag"}, gm.HandleGuildOppositionResult)

	// Guild alliance added handler
	gm.parser.RegisterHandlerFunc("0185", "guild_alliance_added", "V V Z24",
		[]string{"opposition", "guildID", "guildName"}, gm.HandleGuildAllianceAdded)
}

// HandleGuildMembersList handles the guild_members_list packet (lines 6435-6478)
func (gm *GuildManager) HandleGuildMembersList(args map[string]interface{}) error {
	// TODO: Implement guild members list handling
	gm.logger.Info("Guild members list received")
	return nil
}

// HandleGuildName handles the guild_name packet (lines 6632-6666)
func (gm *GuildManager) HandleGuildName(args map[string]interface{}) error {
	// TODO: Implement guild name handling
	gm.logger.Info("Guild information received")
	return nil
}

// HandleGuildMemberOnlineStatus handles the guild_member_online_status packet (lines 6571-6591)
func (gm *GuildManager) HandleGuildMemberOnlineStatus(args map[string]interface{}) error {
	// TODO: Implement guild member online status handling
	gm.logger.Info("Guild member online status update received")
	return nil
}

// HandleGuildNotice handles the guild_notice packet (lines 6842-6857)
func (gm *GuildManager) HandleGuildNotice(args map[string]interface{}) error {
	// TODO: Implement guild notice handling
	gm.logger.Info("Guild notice received")
	return nil
}

// HandleGuildAlliesEnemyList handles the guild_allies_enemy_list packet (lines 6327-6357)
func (gm *GuildManager) HandleGuildAlliesEnemyList(args map[string]interface{}) error {
	// TODO: Implement guild allies/enemy list handling
	gm.logger.Info("Guild allies/enemy list received")
	return nil
}

// HandleGuildAllyRequest handles the guild_ally_request packet (lines 6359-6370)
func (gm *GuildManager) HandleGuildAllyRequest(args map[string]interface{}) error {
	// Extract packet data
	_, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in guild_ally_request packet")
	}

	guildName, ok := args["guildName"].(string)
	if !ok {
		return fmt.Errorf("invalid guildName in guild_ally_request packet")
	}

	// Log the alliance request
	gm.logger.Info("Incoming Request to Ally Guild '%s'", guildName)

	return nil
}

// HandleGuildCreateResult handles the guild_create_result packet (lines 6396-6417)
func (gm *GuildManager) HandleGuildCreateResult(args map[string]interface{}) error {
	// Extract packet data
	typeVal, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in guild_create_result packet")
	}

	// Process based on type
	switch typeVal {
	case 0:
		gm.logger.Info("Guild create successful.")
	case 2:
		gm.logger.Info("Guild create failed: Guild name already exists.")
	case 3:
		gm.logger.Info("Guild create failed: Emperium is needed.")
	default:
		gm.logger.Info("Guild create: Unknown error %d", typeVal)
	}

	return nil
}

// HandleGuildInfo handles the guild_info packet (lines 6424-6433)
func (gm *GuildManager) HandleGuildInfo(args map[string]interface{}) error {
	// Extract packet data
	guildID, ok := args["guildID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid guildID in guild_info packet")
	}

	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in guild_info packet")
	}

	master, ok := args["master"].(string)
	if !ok {
		return fmt.Errorf("invalid master in guild_info packet")
	}

	// Log guild info
	gm.logger.Info("Guild Info: %s (ID: %d, Master: %s)", name, guildID, master)

	return nil
}

// HandleGuildInviteResult handles the guild_invite_result packet (lines 6480-6503)
func (gm *GuildManager) HandleGuildInviteResult(args map[string]interface{}) error {
	// Extract packet data
	typeVal, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in guild_invite_result packet")
	}

	// Process based on type
	switch typeVal {
	case 0:
		gm.logger.Info("Target is already in a guild.")
	case 1:
		gm.logger.Info("Target has denied.")
	case 2:
		gm.logger.Info("Target has accepted.")
	case 3:
		gm.logger.Info("Your guild is full.")
	default:
		gm.logger.Info("Guild join request: Unknown %d", typeVal)
	}

	return nil
}

// HandleGuildLocation handles the guild_location packet (lines 6507-6522)
func (gm *GuildManager) HandleGuildLocation(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in guild_location packet")
	}

	x, ok := args["x"].(uint16)
	if !ok {
		return fmt.Errorf("invalid x in guild_location packet")
	}

	y, ok := args["y"].(uint16)
	if !ok {
		return fmt.Errorf("invalid y in guild_location packet")
	}

	// Log location update (debug level)
	gm.logger.Debug("Guild member location update: ID %d (%d, %d)", id, x, y)

	return nil
}

// HandleGuildLeave handles the guild_leave packet (lines 6524-6545)
func (gm *GuildManager) HandleGuildLeave(args map[string]interface{}) error {
	// Extract packet data
	var name string
	var message string

	// Extract name (could be in args or we might need to look it up by charID)
	if val, ok := args["name"].(string); ok {
		name = val
	} else if val, ok := args["charID"].(uint32); ok {
		// In a real implementation, we would look up the name by charID
		name = fmt.Sprintf("Member_%d", val)
	} else {
		return fmt.Errorf("invalid name/charID in guild_leave packet")
	}

	// Extract message
	if val, ok := args["message"].(string); ok {
		message = val
	}

	// Log the guild leave
	gm.logger.Info("%s has left the guild. Reason: %s", name, message)

	return nil
}

// HandleGuildExpulsion handles the guild_expulsion packet (lines 6547-6569)
func (gm *GuildManager) HandleGuildExpulsion(args map[string]interface{}) error {
	// Extract packet data
	var name string
	var message string

	// Extract name (could be in args or we might need to look it up by charID)
	if val, ok := args["name"].(string); ok {
		name = val
	} else if val, ok := args["charID"].(uint32); ok {
		// In a real implementation, we would look up the name by charID
		name = fmt.Sprintf("Member_%d", val)
	} else {
		return fmt.Errorf("invalid name/charID in guild_expulsion packet")
	}

	// Extract message
	if val, ok := args["message"].(string); ok {
		message = val
	}

	// Log the guild expulsion
	gm.logger.Info("%s has been removed from the guild. Reason: %s", name, message)

	return nil
}

// HandleGuildBroken handles the guild_broken packet (lines 6378-6394)
func (gm *GuildManager) HandleGuildBroken(args map[string]interface{}) error {
	// Extract packet data
	flag, ok := args["flag"].(uint32)
	if !ok {
		return fmt.Errorf("invalid flag in guild_broken packet")
	}

	// Process based on flag
	switch flag {
	case 0:
		gm.logger.Info("Guild broken.")
	case 1:
		gm.logger.Error("Guild can not be undone: invalid key")
	case 2:
		gm.logger.Error("Guild can not be undone: there are still members in the guild")
	default:
		gm.logger.Error("Guild can not be undone: unknown reason (flag: %d)", flag)
	}

	return nil
}

// HandleGuildUpdateMemberPosition handles the guild_update_member_position packet (lines 6593-6615)
func (gm *GuildManager) HandleGuildUpdateMemberPosition(args map[string]interface{}) error {
	// Extract packet data
	_, ok := args["member_list"].([]byte)
	if !ok {
		return fmt.Errorf("invalid member_list in guild_update_member_position packet")
	}

	// Log the position update
	gm.logger.Info("Guild member positions updated")

	return nil
}

// HandleGuildMembersTitleList handles the guild_members_title_list packet (lines 6617-6630)
func (gm *GuildManager) HandleGuildMembersTitleList(args map[string]interface{}) error {
	// Extract packet data
	_, ok := args["RAW_MSG"].([]byte)
	if !ok {
		return fmt.Errorf("invalid RAW_MSG in guild_members_title_list packet")
	}

	// Log the title list update
	gm.logger.Info("Guild position titles updated")

	return nil
}

// HandleGuildRequest handles the guild_request packet (lines 6670-6680)
func (gm *GuildManager) HandleGuildRequest(args map[string]interface{}) error {
	// Extract packet data
	_, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in guild_request packet")
	}

	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in guild_request packet")
	}

	// Log the guild request
	gm.logger.Info("Incoming Request to join Guild '%s'", name)

	return nil
}

// HandleGuildMasterMember handles the guild_master_member packet (lines 6692-6703)
func (gm *GuildManager) HandleGuildMasterMember(args map[string]interface{}) error {
	// Extract packet data
	typeVal, ok := args["type"].(uint32)
	if !ok {
		return fmt.Errorf("invalid type in guild_master_member packet")
	}

	// Process based on type
	if typeVal == 0xd7 {
		gm.logger.Info("You are a guildmaster.")
	} else if typeVal == 0x57 {
		gm.logger.Info("You are not a guildmaster.")
	} else {
		gm.logger.Warning("Unknown results in guild_master_member (type: %d)", typeVal)
	}

	return nil
}

// HandleGuildAlliance handles the guild_alliance packet (lines 6714-6729)
func (gm *GuildManager) HandleGuildAlliance(args map[string]interface{}) error {
	// Extract packet data
	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in guild_alliance packet")
	}

	// Process based on flag
	switch flag {
	case 0:
		gm.logger.Info("Already allied.")
	case 1:
		gm.logger.Info("You rejected the offer.")
	case 2:
		gm.logger.Info("You accepted the offer.")
	case 3:
		gm.logger.Info("They have too any alliances.")
	case 4:
		gm.logger.Info("You have too many alliances.")
	default:
		gm.logger.Warning("Unknown results in guild_alliance (flag: %d)", flag)
	}

	return nil
}

// HandleGuildMemberSettingList handles the guild_member_setting_list packet (lines 6738-6751)
func (gm *GuildManager) HandleGuildMemberSettingList(args map[string]interface{}) error {
	// Extract packet data
	_, ok := args["RAW_MSG"].([]byte)
	if !ok {
		return fmt.Errorf("invalid RAW_MSG in guild_member_setting_list packet")
	}

	// Log the setting list update
	gm.logger.Info("Guild member settings updated")

	return nil
}

// HandleGuildSkillsList handles the guild_skills_list packet (lines 6753-6773)
func (gm *GuildManager) HandleGuildSkillsList(args map[string]interface{}) error {
	// Extract packet data
	skillPoints, ok := args["skillPoints"].(uint16)
	if !ok {
		return fmt.Errorf("invalid skillPoints in guild_skills_list packet")
	}

	_, ok = args["RAW_MSG"].([]byte)
	if !ok {
		return fmt.Errorf("invalid RAW_MSG in guild_skills_list packet")
	}

	// Log the skills list update
	gm.logger.Info("Guild skills list updated (Skill Points: %d)", skillPoints)

	return nil
}

// HandleGuildExpulsionList handles the guild_expulsion_list packet (lines 6780-6807)
func (gm *GuildManager) HandleGuildExpulsionList(args map[string]interface{}) error {
	// Extract packet data
	_, ok := args["expulsion_list"].([]byte)
	if !ok {
		return fmt.Errorf("invalid expulsion_list in guild_expulsion_list packet")
	}

	// Log the expulsion list update
	gm.logger.Info("Guild expulsion list updated")

	return nil
}

// HandleGuildMemberMapChange handles the guild_member_map_change packet (lines 6809-6823)
func (gm *GuildManager) HandleGuildMemberMapChange(args map[string]interface{}) error {
	// Extract packet data
	charID, ok := args["charID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid charID in guild_member_map_change packet")
	}

	mapName, ok := args["mapName"].(string)
	if !ok {
		return fmt.Errorf("invalid mapName in guild_member_map_change packet")
	}

	// In a real implementation, we would look up the member name by charID
	memberName := fmt.Sprintf("Member_%d", charID)

	// Log the map change
	gm.logger.Info("Guild Member: %s changed map to %s", memberName, mapName)

	return nil
}

// HandleGuildMemberAdd handles the guild_member_add packet (lines 6824-6840)
func (gm *GuildManager) HandleGuildMemberAdd(args map[string]interface{}) error {
	// Extract packet data
	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in guild_member_add packet")
	}

	// Log the member addition
	gm.logger.Info("Guild member added: %s", name)

	return nil
}

// HandleGuildEmblem handles the guild_emblem packet (lines 10370-10373)
func (gm *GuildManager) HandleGuildEmblem(args map[string]interface{}) error {
	// Extract packet data
	guildID, ok := args["guildID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid guildID in guild_emblem packet")
	}

	emblemID, ok := args["emblemID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid emblemID in guild_emblem packet")
	}

	// Log the emblem update
	gm.logger.Info("Guild emblem updated for guild ID %d (Emblem ID: %d)", guildID, emblemID)

	return nil
}

// HandleGuildEmblemUpdate handles the guild_emblem_update packet (lines 10377-10380)
func (gm *GuildManager) HandleGuildEmblemUpdate(args map[string]interface{}) error {
	// Extract packet data
	guildID, ok := args["guildID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid guildID in guild_emblem_update packet")
	}

	emblemID, ok := args["emblemID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid emblemID in guild_emblem_update packet")
	}

	// Log the emblem update
	gm.logger.Info("Guild emblem updated for guild ID %d (Emblem ID: %d)", guildID, emblemID)

	return nil
}

// HandleGuildPositionChanged handles the guild_position_changed packet (lines 10391-10394)
func (gm *GuildManager) HandleGuildPositionChanged(args map[string]interface{}) error {
	// Extract packet data
	_, ok := args["RAW_MSG"].([]byte)
	if !ok {
		return fmt.Errorf("invalid RAW_MSG in guild_position_changed packet")
	}

	// Log the position change
	gm.logger.Info("Guild positions changed")

	return nil
}

// HandleGuildPosition handles the guild_position packet (lines 10398-10401)
func (gm *GuildManager) HandleGuildPosition(args map[string]interface{}) error {
	// Extract packet data
	_, ok := args["RAW_MSG"].([]byte)
	if !ok {
		return fmt.Errorf("invalid RAW_MSG in guild_position packet")
	}

	// Log the position update
	gm.logger.Info("Guild positions updated")

	return nil
}

// HandleGuildUnally handles the guild_unally packet (lines 10405-10408)
func (gm *GuildManager) HandleGuildUnally(args map[string]interface{}) error {
	// Extract packet data
	guildID, ok := args["guildID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid guildID in guild_unally packet")
	}

	guildName, ok := args["guildName"].(string)
	if !ok {
		return fmt.Errorf("invalid guildName in guild_unally packet")
	}

	// Log the unally
	gm.logger.Info("Guild alliance broken with %s (ID: %d)", guildName, guildID)

	return nil
}

// HandleGuildOppositionResult handles the guild_opposition_result packet (lines 10412-10415)
func (gm *GuildManager) HandleGuildOppositionResult(args map[string]interface{}) error {
	// Extract packet data
	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in guild_opposition_result packet")
	}

	// Log the opposition result
	gm.logger.Info("Guild opposition result: %d", flag)

	return nil
}

// HandleGuildAllianceAdded handles the guild_alliance_added packet (lines 10419-10422)
func (gm *GuildManager) HandleGuildAllianceAdded(args map[string]interface{}) error {
	// Extract packet data
	opposition, ok := args["opposition"].(uint32)
	if !ok {
		return fmt.Errorf("invalid opposition in guild_alliance_added packet")
	}

	guildID, ok := args["guildID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid guildID in guild_alliance_added packet")
	}

	guildName, ok := args["guildName"].(string)
	if !ok {
		return fmt.Errorf("invalid guildName in guild_alliance_added packet")
	}

	// Log the alliance addition
	if opposition == 0 {
		gm.logger.Info("Alliance formed with guild %s (ID: %d)", guildName, guildID)
	} else {
		gm.logger.Info("Opposition formed with guild %s (ID: %d)", guildName, guildID)
	}

	return nil
}
