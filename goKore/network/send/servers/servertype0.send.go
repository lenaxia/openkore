// Package servers provides server-specific packet constructions for different server types.
package servers

import (
	"github.com/lenaxia/goKore/network/common"
)

// ServerType0PacketConstructions provides packet constructions for ServerType0
func ServerType0PacketConstructions() map[string]common.PacketConstruction {
	packets := map[string]common.PacketConstruction{}

	packets["0064"] = common.PacketConstruction{
		ID:         "0064",
		Name:       "master_login",
		Format:     "V Z24 Z24 C",
		FieldNames: []string{"version", "username", "password", "master_version"},
	}
	packets["0065"] = common.PacketConstruction{
		ID:         "0065",
		Name:       "game_login",
		Format:     "a4 a4 a4 v C",
		FieldNames: []string{"accountID", "sessionID", "sessionID2", "userLevel", "accountSex"},
	}
	packets["0066"] = common.PacketConstruction{
		ID:         "0066",
		Name:       "char_login",
		Format:     "C",
		FieldNames: []string{"slot"},
	}
	packets["0067"] = common.PacketConstruction{
		ID:         "0067",
		Name:       "char_create",
		Format:     "a24 C7 v2",
		FieldNames: []string{"name", "str", "agi", "vit", "int", "dex", "luk", "slot", "hair_color", "hair_style"},
	}
	packets["0068"] = common.PacketConstruction{
		ID:         "0068",
		Name:       "char_delete",
		Format:     "a4 a40",
		FieldNames: []string{"charID", "email"},
	}
	packets["0072"] = common.PacketConstruction{
		ID:         "0072",
		Name:       "map_login",
		Format:     "a4 a4 a4 V C",
		FieldNames: []string{"accountID", "charID", "sessionID", "tick", "sex"},
	}
	packets["007D"] = common.PacketConstruction{
		ID:         "007D",
		Name:       "map_loaded",
		Format:     "",
		FieldNames: []string{},
	}
	packets["007E"] = common.PacketConstruction{
		ID:         "007E",
		Name:       "sync",
		Format:     "V",
		FieldNames: []string{"time"},
	}
	packets["0085"] = common.PacketConstruction{
		ID:         "0085",
		Name:       "character_move",
		Format:     "a3",
		FieldNames: []string{"coords"},
	}
	packets["0089"] = common.PacketConstruction{
		ID:         "0089",
		Name:       "actor_action",
		Format:     "a4 C",
		FieldNames: []string{"targetID", "type"},
	}
	packets["008C"] = common.PacketConstruction{
		ID:         "008C",
		Name:       "public_chat",
		Format:     "v Z*",
		FieldNames: []string{"len", "message"},
	}
	packets["0090"] = common.PacketConstruction{
		ID:         "0090",
		Name:       "npc_talk",
		Format:     "a4 C",
		FieldNames: []string{"ID", "type"},
	}
	packets["0094"] = common.PacketConstruction{
		ID:         "0094",
		Name:       "actor_info_request",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["0096"] = common.PacketConstruction{
		ID:         "0096",
		Name:       "private_message",
		Format:     "x2 Z24 Z*",
		FieldNames: []string{"privMsgUser", "privMsg"},
	}
	packets["0099"] = common.PacketConstruction{
		ID:         "0099",
		Name:       "gm_broadcast",
		Format:     "v Z*",
		FieldNames: []string{"len", "message"},
	}
	packets["009B"] = common.PacketConstruction{
		ID:         "009B",
		Name:       "actor_look_at",
		Format:     "v C",
		FieldNames: []string{"head", "body"},
	}
	packets["009F"] = common.PacketConstruction{
		ID:         "009F",
		Name:       "item_take",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["00A2"] = common.PacketConstruction{
		ID:         "00A2",
		Name:       "item_drop",
		Format:     "a2 v",
		FieldNames: []string{"ID", "amount"},
	}
	packets["00A7"] = common.PacketConstruction{
		ID:         "00A7",
		Name:       "item_use",
		Format:     "a2 a4",
		FieldNames: []string{"ID", "targetID"},
	}
	packets["00A9"] = common.PacketConstruction{
		ID:         "00A9",
		Name:       "send_equip",
		Format:     "a2 v",
		FieldNames: []string{"ID", "type"},
	}
	packets["00AB"] = common.PacketConstruction{
		ID:         "00AB",
		Name:       "send_unequip_item",
		Format:     "a2",
		FieldNames: []string{"ID"},
	}
	packets["00B2"] = common.PacketConstruction{
		ID:         "00B2",
		Name:       "restart",
		Format:     "C",
		FieldNames: []string{"type"},
	}
	packets["00B8"] = common.PacketConstruction{
		ID:         "00B8",
		Name:       "npc_talk_response",
		Format:     "a4 C",
		FieldNames: []string{"ID", "response"},
	}
	packets["00B9"] = common.PacketConstruction{
		ID:         "00B9",
		Name:       "npc_talk_continue",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["00BB"] = common.PacketConstruction{
		ID:         "00BB",
		Name:       "send_add_status_point",
		Format:     "v2",
		FieldNames: []string{"statusID", "Amount"},
	}
	packets["00BF"] = common.PacketConstruction{
		ID:         "00BF",
		Name:       "send_emotion",
		Format:     "C",
		FieldNames: []string{"ID"},
	}
	packets["00C1"] = common.PacketConstruction{
		ID:         "00C1",
		Name:       "request_user_count",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00C5"] = common.PacketConstruction{
		ID:         "00C5",
		Name:       "request_buy_sell_list",
		Format:     "a4 C",
		FieldNames: []string{"ID", "type"},
	}
	packets["00C8"] = common.PacketConstruction{
		ID:         "00C8",
		Name:       "buy_bulk",
		Format:     "v a*",
		FieldNames: []string{"len", "buyInfo"},
	}
	packets["00C9"] = common.PacketConstruction{
		ID:         "00C9",
		Name:       "sell_bulk",
		Format:     "v a*",
		FieldNames: []string{"len", "sellInfo"},
	}
	packets["00CC"] = common.PacketConstruction{
		ID:         "00CC",
		Name:       "gm_kick",
		Format:     "a4",
		FieldNames: []string{"targetAccountID"},
	}
	packets["00CE"] = common.PacketConstruction{
		ID:         "00CE",
		Name:       "gm_kick_all",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00CF"] = common.PacketConstruction{
		ID:         "00CF",
		Name:       "ignore_player",
		Format:     "Z24 C",
		FieldNames: []string{"name", "flag"},
	}
	packets["00D0"] = common.PacketConstruction{
		ID:         "00D0",
		Name:       "ignore_all",
		Format:     "C",
		FieldNames: []string{"flag"},
	}
	packets["00D3"] = common.PacketConstruction{
		ID:         "00D3",
		Name:       "get_ignore_list",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00D5"] = common.PacketConstruction{
		ID:         "00D5",
		Name:       "chat_room_create",
		Format:     "v2 C Z8 a*",
		FieldNames: []string{"len", "limit", "public", "password", "title"},
	}
	packets["00D9"] = common.PacketConstruction{
		ID:         "00D9",
		Name:       "chat_room_join",
		Format:     "a4 Z8",
		FieldNames: []string{"ID", "password"},
	}
	packets["00DE"] = common.PacketConstruction{
		ID:         "00DE",
		Name:       "chat_room_change",
		Format:     "v2 C Z8 a*",
		FieldNames: []string{"len", "limit", "public", "password", "title"},
	}
	packets["00E0"] = common.PacketConstruction{
		ID:         "00E0",
		Name:       "chat_room_bestow",
		Format:     "V Z24",
		FieldNames: []string{"role", "name"},
	}
	packets["00E2"] = common.PacketConstruction{
		ID:         "00E2",
		Name:       "chat_room_kick",
		Format:     "Z24",
		FieldNames: []string{"name"},
	}
	packets["00E3"] = common.PacketConstruction{
		ID:         "00E3",
		Name:       "chat_room_leave",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00E4"] = common.PacketConstruction{
		ID:         "00E4",
		Name:       "deal_initiate",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["00E6"] = common.PacketConstruction{
		ID:         "00E6",
		Name:       "deal_reply",
		Format:     "C",
		FieldNames: []string{"action"},
	}
	packets["00E8"] = common.PacketConstruction{
		ID:         "00E8",
		Name:       "deal_item_add",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["00EB"] = common.PacketConstruction{
		ID:         "00EB",
		Name:       "deal_finalize",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00ED"] = common.PacketConstruction{
		ID:         "00ED",
		Name:       "deal_cancel",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00EF"] = common.PacketConstruction{
		ID:         "00EF",
		Name:       "deal_trade",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00F3"] = common.PacketConstruction{
		ID:         "00F3",
		Name:       "storage_item_add",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["00F5"] = common.PacketConstruction{
		ID:         "00F5",
		Name:       "storage_item_remove",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["00F7"] = common.PacketConstruction{
		ID:         "00F7",
		Name:       "storage_close",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00FC"] = common.PacketConstruction{
		ID:         "00FC",
		Name:       "party_join_request",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["00FF"] = common.PacketConstruction{
		ID:         "00FF",
		Name:       "party_join",
		Format:     "a4 V",
		FieldNames: []string{"ID", "flag"},
	}
	packets["0100"] = common.PacketConstruction{
		ID:         "0100",
		Name:       "party_leave",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0102"] = common.PacketConstruction{
		ID:         "0102",
		Name:       "party_setting",
		Format:     "V",
		FieldNames: []string{"exp"},
	}
	packets["0103"] = common.PacketConstruction{
		ID:         "0103",
		Name:       "party_kick",
		Format:     "a4 Z24",
		FieldNames: []string{"ID", "name"},
	}
	packets["0108"] = common.PacketConstruction{
		ID:         "0108",
		Name:       "party_chat",
		Format:     "x2 Z*",
		FieldNames: []string{"message"},
	}
	packets["0112"] = common.PacketConstruction{
		ID:         "0112",
		Name:       "send_add_skill_point",
		Format:     "v",
		FieldNames: []string{"skillID"},
	}
	packets["0113"] = common.PacketConstruction{
		ID:         "0113",
		Name:       "skill_use",
		Format:     "v2 a4",
		FieldNames: []string{"lv", "skillID", "targetID"},
	}
	packets["0116"] = common.PacketConstruction{
		ID:         "0116",
		Name:       "skill_use_location",
		Format:     "v4",
		FieldNames: []string{"lv", "skillID", "x", "y"},
	}
	packets["011B"] = common.PacketConstruction{
		ID:         "011B",
		Name:       "warp_select",
		Format:     "v Z16",
		FieldNames: []string{"skillID", "mapName"},
	}
	packets["011D"] = common.PacketConstruction{
		ID:         "011D",
		Name:       "memo_request",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0126"] = common.PacketConstruction{
		ID:         "0126",
		Name:       "cart_add",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["0127"] = common.PacketConstruction{
		ID:         "0127",
		Name:       "cart_get",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["0128"] = common.PacketConstruction{
		ID:         "0128",
		Name:       "storage_to_cart",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["0129"] = common.PacketConstruction{
		ID:         "0129",
		Name:       "cart_to_storage",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["012A"] = common.PacketConstruction{
		ID:         "012A",
		Name:       "companion_release",
		Format:     "",
		FieldNames: []string{},
	}
	packets["012E"] = common.PacketConstruction{
		ID:         "012E",
		Name:       "shop_close",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0130"] = common.PacketConstruction{
		ID:         "0130",
		Name:       "send_entering_vending",
		Format:     "a4",
		FieldNames: []string{"accountID"},
	}
	packets["0134"] = common.PacketConstruction{
		ID:         "0134",
		Name:       "buy_bulk_vender",
		Format:     "x2 a4 a*",
		FieldNames: []string{"venderID", "itemInfo"},
	}
	packets["013F"] = common.PacketConstruction{
		ID:         "013F",
		Name:       "gm_item_mob_create",
		Format:     "a24",
		FieldNames: []string{"name"},
	}
	packets["0140"] = common.PacketConstruction{
		ID:         "0140",
		Name:       "gm_move_to_map",
		Format:     "Z16 v v",
		FieldNames: []string{"mapName", "x", "y"},
	}
	packets["0143"] = common.PacketConstruction{
		ID:         "0143",
		Name:       "npc_talk_number",
		Format:     "a4 V",
		FieldNames: []string{"ID", "value"},
	}
	packets["0146"] = common.PacketConstruction{
		ID:         "0146",
		Name:       "npc_talk_cancel",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["0149"] = common.PacketConstruction{
		ID:         "0149",
		Name:       "alignment",
		Format:     "a4 C v",
		FieldNames: []string{"targetID", "type", "point"},
	}
	packets["014D"] = common.PacketConstruction{
		ID:         "014D",
		Name:       "guild_check",
		Format:     "",
		FieldNames: []string{},
	}
	packets["014F"] = common.PacketConstruction{
		ID:         "014F",
		Name:       "guild_info_request",
		Format:     "V",
		FieldNames: []string{"type"},
	}
	packets["0151"] = common.PacketConstruction{
		ID:         "0151",
		Name:       "guild_emblem_request",
		Format:     "a4",
		FieldNames: []string{"guildID"},
	}
	packets["0159"] = common.PacketConstruction{
		ID:         "0159",
		Name:       "guild_leave",
		Format:     "a4 a4 a4 Z40",
		FieldNames: []string{"guildID", "accountID", "charID", "reason"},
	}
	packets["015B"] = common.PacketConstruction{
		ID:         "015B",
		Name:       "guild_kick",
		Format:     "a4 a4 a4 Z40",
		FieldNames: []string{"guildID", "accountID", "charID", "reason"},
	}
	packets["015D"] = common.PacketConstruction{
		ID:         "015D",
		Name:       "guild_break",
		Format:     "a4",
		FieldNames: []string{"guildName"},
	}
	packets["0165"] = common.PacketConstruction{
		ID:         "0165",
		Name:       "guild_create",
		Format:     "a4 Z24",
		FieldNames: []string{"charID", "guildName"},
	}
	packets["0168"] = common.PacketConstruction{
		ID:         "0168",
		Name:       "guild_join_request",
		Format:     "a4 a4 a4",
		FieldNames: []string{"ID", "accountID", "charID"},
	}
	packets["016B"] = common.PacketConstruction{
		ID:         "016B",
		Name:       "guild_join",
		Format:     "a4 V",
		FieldNames: []string{"ID", "flag"},
	}
	packets["016E"] = common.PacketConstruction{
		ID:         "016E",
		Name:       "guild_notice",
		Format:     "a4 Z60 Z120",
		FieldNames: []string{"guildID", "name", "notice"},
	}
	packets["0170"] = common.PacketConstruction{
		ID:         "0170",
		Name:       "guild_alliance_request",
		Format:     "a4 a4 a4",
		FieldNames: []string{"targetAccountID", "accountID", "charID"},
	}
	packets["0172"] = common.PacketConstruction{
		ID:         "0172",
		Name:       "guild_alliance_reply",
		Format:     "a4 V",
		FieldNames: []string{"ID", "flag"},
	}
	packets["0178"] = common.PacketConstruction{
		ID:         "0178",
		Name:       "identify",
		Format:     "a2",
		FieldNames: []string{"ID"},
	}
	packets["017A"] = common.PacketConstruction{
		ID:         "017A",
		Name:       "card_merge_request",
		Format:     "a2",
		FieldNames: []string{"cardID"},
	}
	packets["017C"] = common.PacketConstruction{
		ID:         "017C",
		Name:       "card_merge",
		Format:     "a2 a2",
		FieldNames: []string{"cardID", "itemID"},
	}
	packets["017E"] = common.PacketConstruction{
		ID:         "017E",
		Name:       "guild_chat",
		Format:     "x2 Z*",
		FieldNames: []string{"message"},
	}
	packets["0187"] = common.PacketConstruction{
		ID:         "0187",
		Name:       "ban_check",
		Format:     "a4",
		FieldNames: []string{"accountID"},
	}
	packets["018A"] = common.PacketConstruction{
		ID:         "018A",
		Name:       "quit_request",
		Format:     "v",
		FieldNames: []string{"type"},
	}
	packets["018E"] = common.PacketConstruction{
		ID:         "018E",
		Name:       "make_item_request",
		Format:     "v4",
		FieldNames: []string{"nameID", "material_nameID1", "material_nameID2", "material_nameID3"},
	}
	packets["0190"] = common.PacketConstruction{
		ID:         "0190",
		Name:       "skill_use_location_text",
		Format:     "v5 Z80",
		FieldNames: []string{"lvl", "ID", "x", "y", "info"},
	}
	packets["0193"] = common.PacketConstruction{
		ID:         "0193",
		Name:       "actor_name_request",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["0197"] = common.PacketConstruction{
		ID:         "0197",
		Name:       "gm_reset_state_skill",
		Format:     "v",
		FieldNames: []string{"type"},
	}
	packets["0198"] = common.PacketConstruction{
		ID:         "0198",
		Name:       "gm_change_cell_type",
		Format:     "v v v",
		FieldNames: []string{"x", "y", "type"},
	}
	packets["019C"] = common.PacketConstruction{
		ID:         "019C",
		Name:       "gm_broadcast_local",
		Format:     "v Z*",
		FieldNames: []string{"len", "message"},
	}
	packets["019D"] = common.PacketConstruction{
		ID:         "019D",
		Name:       "gm_change_effect_state",
		Format:     "V",
		FieldNames: []string{"effect_state"},
	}
	packets["019F"] = common.PacketConstruction{
		ID:         "019F",
		Name:       "pet_capture",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["01A1"] = common.PacketConstruction{
		ID:         "01A1",
		Name:       "pet_menu",
		Format:     "C",
		FieldNames: []string{"action"},
	}
	packets["01A5"] = common.PacketConstruction{
		ID:         "01A5",
		Name:       "pet_name",
		Format:     "a24",
		FieldNames: []string{"name"},
	}
	packets["01A7"] = common.PacketConstruction{
		ID:         "01A7",
		Name:       "pet_hatch",
		Format:     "a2",
		FieldNames: []string{"ID"},
	}
	packets["01A9"] = common.PacketConstruction{
		ID:         "01A9",
		Name:       "pet_emotion",
		Format:     "V",
		FieldNames: []string{"ID"},
	}
	packets["01AE"] = common.PacketConstruction{
		ID:         "01AE",
		Name:       "make_arrow",
		Format:     "v",
		FieldNames: []string{"nameID"},
	}
	packets["01AF"] = common.PacketConstruction{
		ID:         "01AF",
		Name:       "change_cart",
		Format:     "v",
		FieldNames: []string{"lvl"},
	}
	packets["01B2"] = common.PacketConstruction{
		ID:         "01B2",
		Name:       "shop_open",
		Format:     "v a80 C a*",
		FieldNames: []string{"len", "title", "result", "vendingInfo"},
	}
	packets["01BA"] = common.PacketConstruction{
		ID:         "01BA",
		Name:       "gm_remove",
		Format:     "a24",
		FieldNames: []string{"playerName"},
	}
	packets["01BB"] = common.PacketConstruction{
		ID:         "01BB",
		Name:       "gm_shift",
		Format:     "a24",
		FieldNames: []string{"playerName"},
	}
	packets["01BC"] = common.PacketConstruction{
		ID:         "01BC",
		Name:       "gm_recall",
		Format:     "a24",
		FieldNames: []string{"playerName"},
	}
	packets["01BD"] = common.PacketConstruction{
		ID:         "01BD",
		Name:       "gm_summon_player",
		Format:     "a24",
		FieldNames: []string{"playerName"},
	}
	packets["01C0"] = common.PacketConstruction{
		ID:         "01C0",
		Name:       "request_remain_time",
		Format:     "",
		FieldNames: []string{},
	}
	packets["01CE"] = common.PacketConstruction{
		ID:         "01CE",
		Name:       "auto_spell",
		Format:     "V",
		FieldNames: []string{"ID"},
	}
	packets["01D5"] = common.PacketConstruction{
		ID:         "01D5",
		Name:       "npc_talk_text",
		Format:     "v a4 Z*",
		FieldNames: []string{"len", "ID", "text"},
	}
	packets["01DB"] = common.PacketConstruction{
		ID:         "01DB",
		Name:       "secure_login_key_request",
		Format:     "",
		FieldNames: []string{},
	}
	packets["01DD"] = common.PacketConstruction{
		ID:         "01DD",
		Name:       "master_login",
		Format:     "V Z24 a16 C",
		FieldNames: []string{"version", "username", "password_salted_md5", "master_version"},
	}
	packets["01DF"] = common.PacketConstruction{
		ID:         "01DF",
		Name:       "gm_request_account_name",
		Format:     "V",
		FieldNames: []string{"targetID"},
	}
	packets["01E7"] = common.PacketConstruction{
		ID:         "01E7",
		Name:       "novice_dori_dori",
		Format:     "",
		FieldNames: []string{},
	}
	packets["01ED"] = common.PacketConstruction{
		ID:         "01ED",
		Name:       "novice_explosion_spirits",
		Format:     "",
		FieldNames: []string{},
	}
	packets["01F7"] = common.PacketConstruction{
		ID:         "01F7",
		Name:       "adopt_reply_request",
		Format:     "V3",
		FieldNames: []string{"parentID1", "parentID2", "result"},
	}
	packets["01F9"] = common.PacketConstruction{
		ID:         "01F9",
		Name:       "adopt_request",
		Format:     "V",
		FieldNames: []string{"ID"},
	}
	packets["01FA"] = common.PacketConstruction{
		ID:         "01FA",
		Name:       "master_login",
		Format:     "V Z24 a16 C C",
		FieldNames: []string{"version", "username", "password_salted_md5", "master_version", "clientInfo"},
	}
	packets["01FD"] = common.PacketConstruction{
		ID:         "01FD",
		Name:       "repair_item",
		Format:     "v2 C a8",
		FieldNames: []string{"index", "nameID", "upgrade", "cards"},
	}
	packets["0202"] = common.PacketConstruction{
		ID:         "0202",
		Name:       "friend_request",
		Format:     "a*",
		FieldNames: []string{"username"},
	}
	packets["0203"] = common.PacketConstruction{
		ID:         "0203",
		Name:       "friend_remove",
		Format:     "a4 a4",
		FieldNames: []string{"accountID", "charID"},
	}
	packets["0204"] = common.PacketConstruction{
		ID:         "0204",
		Name:       "client_hash",
		Format:     "a16",
		FieldNames: []string{"hash"},
	}
	packets["0208"] = common.PacketConstruction{
		ID:         "0208",
		Name:       "friend_response",
		Format:     "a4 a4 V",
		FieldNames: []string{"friendAccountID", "friendCharID", "type"},
	}
	packets["0212"] = common.PacketConstruction{
		ID:         "0212",
		Name:       "manner_by_name",
		Format:     "Z24",
		FieldNames: []string{"playerName"},
	}
	packets["0213"] = common.PacketConstruction{
		ID:         "0213",
		Name:       "gm_request_status",
		Format:     "Z24",
		FieldNames: []string{"playerName"},
	}
	packets["0217"] = common.PacketConstruction{
		ID:         "0217",
		Name:       "rank_blacksmith",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0218"] = common.PacketConstruction{
		ID:         "0218",
		Name:       "rank_alchemist",
		Format:     "",
		FieldNames: []string{},
	}
	packets["021D"] = common.PacketConstruction{
		ID:         "021D",
		Name:       "less_effect",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0222"] = common.PacketConstruction{
		ID:         "0222",
		Name:       "refine_item",
		Format:     "V",
		FieldNames: []string{"ID"},
	}
	packets["0225"] = common.PacketConstruction{
		ID:         "0225",
		Name:       "rank_taekwon",
		Format:     "",
		FieldNames: []string{},
	}
	packets["022D"] = common.PacketConstruction{
		ID:         "022D",
		Name:       "homunculus_command",
		Format:     "v C",
		FieldNames: []string{"commandType", "commandID"},
	}
	packets["0231"] = common.PacketConstruction{
		ID:         "0231",
		Name:       "homunculus_name",
		Format:     "a24",
		FieldNames: []string{"name"},
	}
	packets["0232"] = common.PacketConstruction{
		ID:         "0232",
		Name:       "actor_move",
		Format:     "a4 a3",
		FieldNames: []string{"ID", "coords"},
	}
	packets["0233"] = common.PacketConstruction{
		ID:         "0233",
		Name:       "slave_attack",
		Format:     "a4 a4 C",
		FieldNames: []string{"slaveID", "targetID", "flag"},
	}
	packets["0234"] = common.PacketConstruction{
		ID:         "0234",
		Name:       "slave_move_to_master",
		Format:     "a4",
		FieldNames: []string{"slaveID"},
	}
	packets["0237"] = common.PacketConstruction{
		ID:         "0237",
		Name:       "rank_killer",
		Format:     "",
		FieldNames: []string{},
	}
	packets["023B"] = common.PacketConstruction{
		ID:         "023B",
		Name:       "storage_password",
		Format:     "v a*",
		FieldNames: []string{"type", "data"},
	}
	packets["023F"] = common.PacketConstruction{
		ID:         "023F",
		Name:       "mailbox_open",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0241"] = common.PacketConstruction{
		ID:         "0241",
		Name:       "mail_read",
		Format:     "V",
		FieldNames: []string{"mailID"},
	}
	packets["0243"] = common.PacketConstruction{
		ID:         "0243",
		Name:       "mail_delete",
		Format:     "V",
		FieldNames: []string{"mailID"},
	}
	packets["0244"] = common.PacketConstruction{
		ID:         "0244",
		Name:       "mail_attachment_get",
		Format:     "V",
		FieldNames: []string{"mailID"},
	}
	packets["0246"] = common.PacketConstruction{
		ID:         "0246",
		Name:       "mail_remove",
		Format:     "v",
		FieldNames: []string{"flag"},
	}
	packets["0247"] = common.PacketConstruction{
		ID:         "0247",
		Name:       "mail_attachment_set",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["0248"] = common.PacketConstruction{
		ID:         "0248",
		Name:       "mail_send",
		Format:     "v Z24 a40 C a*",
		FieldNames: []string{"len", "recipient", "title", "body_len", "body"},
	}
	packets["024B"] = common.PacketConstruction{
		ID:         "024B",
		Name:       "auction_add_item_cancel",
		Format:     "v",
		FieldNames: []string{"flag"},
	}
	packets["024C"] = common.PacketConstruction{
		ID:         "024C",
		Name:       "auction_add_item",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["024D"] = common.PacketConstruction{
		ID:         "024D",
		Name:       "auction_create",
		Format:     "V V v",
		FieldNames: []string{"now_price", "max_price", "delete_time"},
	}
	packets["024E"] = common.PacketConstruction{
		ID:         "024E",
		Name:       "auction_cancel",
		Format:     "V",
		FieldNames: []string{"ID"},
	}
	packets["024F"] = common.PacketConstruction{
		ID:         "024F",
		Name:       "auction_buy",
		Format:     "V V",
		FieldNames: []string{"ID", "price"},
	}
	packets["0251"] = common.PacketConstruction{
		ID:         "0251",
		Name:       "auction_search",
		Format:     "v V Z24 v",
		FieldNames: []string{"type", "price", "search_string", "page"},
	}
	packets["0254"] = common.PacketConstruction{
		ID:         "0254",
		Name:       "starplace_agree",
		Format:     "C",
		FieldNames: []string{"flag"},
	}
	packets["025B"] = common.PacketConstruction{
		ID:         "025B",
		Name:       "cook_request",
		Format:     "v2",
		FieldNames: []string{"type", "nameID"},
	}
	packets["025C"] = common.PacketConstruction{
		ID:         "025C",
		Name:       "auction_info_self",
		Format:     "v",
		FieldNames: []string{"type"},
	}
	packets["025D"] = common.PacketConstruction{
		ID:         "025D",
		Name:       "auction_sell_stop",
		Format:     "V",
		FieldNames: []string{"ID"},
	}
	packets["0273"] = common.PacketConstruction{
		ID:         "0273",
		Name:       "mail_return",
		Format:     "V Z24",
		FieldNames: []string{"mailID", "sender"},
	}
	packets["0275"] = common.PacketConstruction{
		ID:         "0275",
		Name:       "game_login",
		Format:     "a4 a4 a4 v C Z16 V",
		FieldNames: []string{"accountID", "sessionID", "sessionID2", "userLevel", "accountSex", "mac", "iAccountSID"},
	}
	packets["0288"] = common.PacketConstruction{
		ID:         "0288",
		Name:       "cash_dealer_buy",
		Format:     "v2 V",
		FieldNames: []string{"itemid", "amount", "kafra_points"},
	}
	packets["0292"] = common.PacketConstruction{
		ID:         "0292",
		Name:       "auto_revive",
		Format:     "",
		FieldNames: []string{},
	}
	packets["029F"] = common.PacketConstruction{
		ID:         "029F",
		Name:       "mercenary_command",
		Format:     "C",
		FieldNames: []string{"flag"},
	}
	packets["02B0"] = common.PacketConstruction{
		ID:         "02B0",
		Name:       "master_login",
		Format:     "V Z24 a24 C Z16 Z14 C",
		FieldNames: []string{"version", "username", "password_rijndael", "master_version", "ip", "mac", "isGravityID"},
	}
	packets["02B6"] = common.PacketConstruction{
		ID:         "02B6",
		Name:       "send_quest_state",
		Format:     "V C",
		FieldNames: []string{"questID", "state"},
	}
	packets["02BA"] = common.PacketConstruction{
		ID:         "02BA",
		Name:       "hotkey_change",
		Format:     "v C V v",
		FieldNames: []string{"idx", "type", "id", "lvl"},
	}
	packets["02C4"] = common.PacketConstruction{
		ID:         "02C4",
		Name:       "party_join_request_by_name",
		Format:     "Z24",
		FieldNames: []string{"partyName"},
	}
	packets["02C7"] = common.PacketConstruction{
		ID:         "02C7",
		Name:       "party_join_request_by_name_reply",
		Format:     "a4 C",
		FieldNames: []string{"accountID", "flag"},
	}
	packets["02CF"] = common.PacketConstruction{
		ID:         "02CF",
		Name:       "memorial_dungeon_command",
		Format:     "V",
		FieldNames: []string{"command"},
	}
	packets["02D6"] = common.PacketConstruction{
		ID:         "02D6",
		Name:       "view_player_equip_request",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["02D8"] = common.PacketConstruction{
		ID:         "02D8",
		Name:       "misc_config_set",
		Format:     "V2",
		FieldNames: []string{"type", "flag"},
	}
	packets["02DB"] = common.PacketConstruction{
		ID:         "02DB",
		Name:       "battleground_chat",
		Format:     "v Z*",
		FieldNames: []string{"len", "message"},
	}
	packets["02F1"] = common.PacketConstruction{
		ID:         "02F1",
		Name:       "notify_progress_bar_complete",
		Format:     "",
		FieldNames: []string{},
	}
	packets["035F"] = common.PacketConstruction{
		ID:         "035F",
		Name:       "character_move",
		Format:     "a3",
		FieldNames: []string{"coords"},
	}
	packets["0360"] = common.PacketConstruction{
		ID:         "0360",
		Name:       "sync",
		Format:     "V",
		FieldNames: []string{"time"},
	}
	packets["0361"] = common.PacketConstruction{
		ID:         "0361",
		Name:       "actor_look_at",
		Format:     "v C",
		FieldNames: []string{"head", "body"},
	}
	packets["0362"] = common.PacketConstruction{
		ID:         "0362",
		Name:       "item_take",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["0363"] = common.PacketConstruction{
		ID:         "0363",
		Name:       "item_drop",
		Format:     "a2 v",
		FieldNames: []string{"ID", "amount"},
	}
	packets["0364"] = common.PacketConstruction{
		ID:         "0364",
		Name:       "storage_item_add",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["0365"] = common.PacketConstruction{
		ID:         "0365",
		Name:       "storage_item_remove",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["0366"] = common.PacketConstruction{
		ID:         "0366",
		Name:       "skill_use_location",
		Format:     "v4",
		FieldNames: []string{"lv", "skillID", "x", "y"},
	}
	packets["0367"] = common.PacketConstruction{
		ID:         "0367",
		Name:       "skill_use_location_text",
		Format:     "v5 Z80",
		FieldNames: []string{"lvl", "ID", "x", "y", "info"},
	}
	packets["0368"] = common.PacketConstruction{
		ID:         "0368",
		Name:       "actor_info_request",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["0369"] = common.PacketConstruction{
		ID:         "0369",
		Name:       "actor_name_request",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["0436"] = common.PacketConstruction{
		ID:         "0436",
		Name:       "map_login",
		Format:     "a4 a4 a4 V C",
		FieldNames: []string{"accountID", "charID", "sessionID", "tick", "sex"},
	}
	packets["0437"] = common.PacketConstruction{
		ID:         "0437",
		Name:       "actor_action",
		Format:     "a4 C",
		FieldNames: []string{"targetID", "type"},
	}
	packets["0438"] = common.PacketConstruction{
		ID:         "0438",
		Name:       "skill_use",
		Format:     "v2 a4",
		FieldNames: []string{"lv", "skillID", "targetID"},
	}
	packets["0439"] = common.PacketConstruction{
		ID:         "0439",
		Name:       "item_use",
		Format:     "a2 a4",
		FieldNames: []string{"ID", "targetID"},
	}
	packets["0443"] = common.PacketConstruction{
		ID:         "0443",
		Name:       "skill_select",
		Format:     "V v",
		FieldNames: []string{"why", "skillID"},
	}
	packets["0447"] = common.PacketConstruction{
		ID:         "0447",
		Name:       "blocking_play_cancel",
		Format:     "",
		FieldNames: []string{},
	}
	packets["044A"] = common.PacketConstruction{
		ID:         "044A",
		Name:       "client_version",
		Format:     "V",
		FieldNames: []string{"clientVersion"},
	}
	packets["07D7"] = common.PacketConstruction{
		ID:         "07D7",
		Name:       "party_setting",
		Format:     "V C2",
		FieldNames: []string{"exp", "itemPickup", "itemDivision"},
	}
	packets["07DA"] = common.PacketConstruction{
		ID:         "07DA",
		Name:       "party_leader",
		Format:     "a4",
		FieldNames: []string{"accountID"},
	}
	packets["07E4"] = common.PacketConstruction{
		ID:         "07E4",
		Name:       "item_list_window_selected",
		Format:     "v V V a*",
		FieldNames: []string{"len", "type", "act", "itemInfo"},
	}
	packets["07E7"] = common.PacketConstruction{
		ID:         "07E7",
		Name:       "captcha_answer",
		Format:     "v a4 a24",
		FieldNames: []string{"len", "accountID", "answer"},
	}
	packets["0801"] = common.PacketConstruction{
		ID:         "0801",
		Name:       "buy_bulk_vender",
		Format:     "v a4 a4 a*",
		FieldNames: []string{"len", "venderID", "venderCID", "itemInfo"},
	}
	packets["0802"] = common.PacketConstruction{
		ID:         "0802",
		Name:       "booking_register",
		Format:     "v8",
		FieldNames: []string{"level", "MapID", "job0", "job1", "job2", "job3", "job4", "job5"},
	}
	packets["0804"] = common.PacketConstruction{
		ID:         "0804",
		Name:       "booking_search",
		Format:     "v3 V s",
		FieldNames: []string{"level", "MapID", "job", "LastIndex", "ResultCount"},
	}
	packets["0806"] = common.PacketConstruction{
		ID:         "0806",
		Name:       "booking_delete",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0808"] = common.PacketConstruction{
		ID:         "0808",
		Name:       "booking_update",
		Format:     "v6",
		FieldNames: []string{"job0", "job1", "job2", "job3", "job4", "job5"},
	}
	packets["0811"] = common.PacketConstruction{
		ID:         "0811",
		Name:       "buy_bulk_openShop",
		Format:     "v V C Z80 a*",
		FieldNames: []string{"len", "limitZeny", "result", "storeName", "itemInfo"},
	}
	packets["0815"] = common.PacketConstruction{
		ID:         "0815",
		Name:       "buy_bulk_closeShop",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0817"] = common.PacketConstruction{
		ID:         "0817",
		Name:       "buy_bulk_request",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["0819"] = common.PacketConstruction{
		ID:         "0819",
		Name:       "buy_bulk_buyer",
		Format:     "v a4 a4 a*",
		FieldNames: []string{"len", "buyerID", "buyingStoreID", "itemInfo"},
	}
	packets["0825"] = common.PacketConstruction{
		ID:         "0825",
		Name:       "token_login",
		Format:     "v v x v Z24 a27 Z17 Z15 a*",
		FieldNames: []string{"len", "version", "master_version", "username", "password_rijndael", "mac", "ip", "token"},
	}
	packets["0827"] = common.PacketConstruction{
		ID:         "0827",
		Name:       "char_delete2",
		Format:     "a4",
		FieldNames: []string{"charID"},
	}
	packets["0829"] = common.PacketConstruction{
		ID:         "0829",
		Name:       "char_delete2_accept",
		Format:     "a4 a6",
		FieldNames: []string{"charID", "code"},
	}
	packets["082B"] = common.PacketConstruction{
		ID:         "082B",
		Name:       "char_delete2_cancel",
		Format:     "a4",
		FieldNames: []string{"charID"},
	}
	packets["0835"] = common.PacketConstruction{
		ID:         "0835",
		Name:       "search_store_info",
		Format:     "v C V2 C2 a*",
		FieldNames: []string{"len", "type", "max_price", "min_price", "item_count", "card_count", "item_card_list"},
	}
	packets["0838"] = common.PacketConstruction{
		ID:         "0838",
		Name:       "search_store_request_next_page",
		Format:     "",
		FieldNames: []string{},
	}
	packets["083B"] = common.PacketConstruction{
		ID:         "083B",
		Name:       "search_store_close",
		Format:     "",
		FieldNames: []string{},
	}
	packets["083C"] = common.PacketConstruction{
		ID:         "083C",
		Name:       "search_store_select",
		Format:     "a4 a4 v",
		FieldNames: []string{"accountID", "storeID", "nameID"},
	}
	packets["0842"] = common.PacketConstruction{
		ID:         "0842",
		Name:       "recall_sso",
		Format:     "V",
		FieldNames: []string{"ID"},
	}
	packets["0843"] = common.PacketConstruction{
		ID:         "0843",
		Name:       "remove_aid_sso",
		Format:     "V",
		FieldNames: []string{"ID"},
	}
	packets["0844"] = common.PacketConstruction{
		ID:         "0844",
		Name:       "cash_shop_open",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0846"] = common.PacketConstruction{
		ID:         "0846",
		Name:       "req_cash_tabcode",
		Format:     "v",
		FieldNames: []string{"ID"},
	}
	packets["084A"] = common.PacketConstruction{
		ID:         "084A",
		Name:       "cash_shop_close",
		Format:     "",
		FieldNames: []string{},
	}
	packets["08B5"] = common.PacketConstruction{
		ID:         "08B5",
		Name:       "pet_capture",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["08B8"] = common.PacketConstruction{
		ID:         "08B8",
		Name:       "send_pin_password",
		Format:     "a4 Z*",
		FieldNames: []string{"accountID", "pin"},
	}
	packets["08BA"] = common.PacketConstruction{
		ID:         "08BA",
		Name:       "new_pin_password",
		Format:     "a4 Z*",
		FieldNames: []string{"accountID", "pin"},
	}
	packets["08BE"] = common.PacketConstruction{
		ID:         "08BE",
		Name:       "change_pin_password",
		Format:     "a*",
		FieldNames: []string{"accountID", "oldPin", "newPin"},
	}
	packets["08C1"] = common.PacketConstruction{
		ID:         "08C1",
		Name:       "macro_start",
		Format:     "",
		FieldNames: []string{},
	}
	packets["08C2"] = common.PacketConstruction{
		ID:         "08C2",
		Name:       "macro_stop",
		Format:     "",
		FieldNames: []string{},
	}
	packets["08C9"] = common.PacketConstruction{
		ID:         "08C9",
		Name:       "request_cashitems",
		Format:     "",
		FieldNames: []string{},
	}
	packets["096E"] = common.PacketConstruction{
		ID:         "096E",
		Name:       "merge_item_request",
		Format:     "v a*",
		FieldNames: []string{"length", "itemList"},
	}
	packets["0970"] = common.PacketConstruction{
		ID:         "0970",
		Name:       "char_create",
		Format:     "a24 C v2",
		FieldNames: []string{"name", "slot", "hair_style", "hair_color"},
	}
	packets["0974"] = common.PacketConstruction{
		ID:         "0974",
		Name:       "merge_item_cancel",
		Format:     "",
		FieldNames: []string{},
	}
	packets["097C"] = common.PacketConstruction{
		ID:         "097C",
		Name:       "rank_general",
		Format:     "v",
		FieldNames: []string{"type"},
	}
	packets["0987"] = common.PacketConstruction{
		ID:         "0987",
		Name:       "master_login",
		Format:     "V Z24 a32 C",
		FieldNames: []string{"version", "username", "password_md5_hex", "master_version"},
	}
	packets["098D"] = common.PacketConstruction{
		ID:         "098D",
		Name:       "clan_chat",
		Format:     "v Z*",
		FieldNames: []string{"len", "message"},
	}
	packets["098F"] = common.PacketConstruction{
		ID:         "098F",
		Name:       "char_delete2_accept",
		Format:     "v a4 a*",
		FieldNames: []string{"len", "charID", "code"},
	}
	packets["0998"] = common.PacketConstruction{
		ID:         "0998",
		Name:       "send_equip",
		Format:     "a2 V",
		FieldNames: []string{"ID", "type"},
	}
	packets["09A1"] = common.PacketConstruction{
		ID:         "09A1",
		Name:       "sync_received_characters",
		Format:     "",
		FieldNames: []string{},
	}
	packets["09A7"] = common.PacketConstruction{
		ID:         "09A7",
		Name:       "banking_deposit_request",
		Format:     "a4 V",
		FieldNames: []string{"accountID", "zeny"},
	}
	packets["09A9"] = common.PacketConstruction{
		ID:         "09A9",
		Name:       "banking_withdraw_request",
		Format:     "a4 V",
		FieldNames: []string{"accountID", "zeny"},
	}
	packets["09AB"] = common.PacketConstruction{
		ID:         "09AB",
		Name:       "banking_check_request",
		Format:     "a4",
		FieldNames: []string{"accountID"},
	}
	packets["09D0"] = common.PacketConstruction{
		ID:         "09D0",
		Name:       "gameguard_reply",
		Format:     "",
		FieldNames: []string{},
	}
	packets["09D4"] = common.PacketConstruction{
		ID:         "09D4",
		Name:       "sell_buy_complete",
		Format:     "",
		FieldNames: []string{},
	}
	packets["09D6"] = common.PacketConstruction{
		ID:         "09D6",
		Name:       "buy_bulk_market",
		Format:     "v a*",
		FieldNames: []string{"len", "buyInfo"},
	}
	packets["09D8"] = common.PacketConstruction{
		ID:         "09D8",
		Name:       "market_close",
		Format:     "",
		FieldNames: []string{},
	}
	packets["09E1"] = common.PacketConstruction{
		ID:         "09E1",
		Name:       "guild_storage_item_add",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["09E2"] = common.PacketConstruction{
		ID:         "09E2",
		Name:       "guild_storage_item_remove",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["09E3"] = common.PacketConstruction{
		ID:         "09E3",
		Name:       "cart_to_guild_storage",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["09E4"] = common.PacketConstruction{
		ID:         "09E4",
		Name:       "guild_storage_to_cart",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["09E8"] = common.PacketConstruction{
		ID:         "09E8",
		Name:       "rodex_open_mailbox",
		Format:     "C V2",
		FieldNames: []string{"type", "mailID1", "mailID2"},
	}
	packets["09E9"] = common.PacketConstruction{
		ID:         "09E9",
		Name:       "rodex_close_mailbox",
		Format:     "",
		FieldNames: []string{},
	}
	packets["09EA"] = common.PacketConstruction{
		ID:         "09EA",
		Name:       "rodex_read_mail",
		Format:     "C V2",
		FieldNames: []string{"type", "mailID1", "mailID2"},
	}
	packets["09EE"] = common.PacketConstruction{
		ID:         "09EE",
		Name:       "rodex_next_maillist",
		Format:     "C V2",
		FieldNames: []string{"type", "mailID1", "mailID2"},
	}
	packets["09EF"] = common.PacketConstruction{
		ID:         "09EF",
		Name:       "rodex_refresh_maillist",
		Format:     "C V2",
		FieldNames: []string{"type", "mailID1", "mailID2"},
	}
	packets["09F1"] = common.PacketConstruction{
		ID:         "09F1",
		Name:       "rodex_request_zeny",
		Format:     "V2 C",
		FieldNames: []string{"mailID1", "mailID2", "type"},
	}
	packets["09F3"] = common.PacketConstruction{
		ID:         "09F3",
		Name:       "rodex_request_items",
		Format:     "V2 C",
		FieldNames: []string{"mailID1", "mailID2", "type"},
	}
	packets["09F5"] = common.PacketConstruction{
		ID:         "09F5",
		Name:       "rodex_delete_mail",
		Format:     "C V2",
		FieldNames: []string{"type", "mailID1", "mailID2"},
	}
	packets["09FB"] = common.PacketConstruction{
		ID:         "09FB",
		Name:       "pet_evolution",
		Format:     "a4 a*",
		FieldNames: []string{"ID", "itemInfo"},
	}
	packets["0A03"] = common.PacketConstruction{
		ID:         "0A03",
		Name:       "rodex_cancel_write_mail",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0A04"] = common.PacketConstruction{
		ID:         "0A04",
		Name:       "rodex_add_item",
		Format:     "a2 v",
		FieldNames: []string{"ID", "amount"},
	}
	packets["0A06"] = common.PacketConstruction{
		ID:         "0A06",
		Name:       "rodex_remove_item",
		Format:     "a2 v",
		FieldNames: []string{"ID", "amount"},
	}
	packets["0A08"] = common.PacketConstruction{
		ID:         "0A08",
		Name:       "rodex_open_write_mail",
		Format:     "Z24",
		FieldNames: []string{"name"},
	}
	packets["0A13"] = common.PacketConstruction{
		ID:         "0A13",
		Name:       "rodex_checkname",
		Format:     "Z24",
		FieldNames: []string{"name"},
	}
	packets["0A16"] = common.PacketConstruction{
		ID:         "0A16",
		Name:       "dynamicnpc_create_request",
		Format:     "Z24",
		FieldNames: []string{"name"},
	}
	packets["0A19"] = common.PacketConstruction{
		ID:         "0A19",
		Name:       "roulette_window_open",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0A1B"] = common.PacketConstruction{
		ID:         "0A1B",
		Name:       "roulette_info_request",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0A1D"] = common.PacketConstruction{
		ID:         "0A1D",
		Name:       "roulette_close",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0A1F"] = common.PacketConstruction{
		ID:         "0A1F",
		Name:       "roulette_start",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0A21"] = common.PacketConstruction{
		ID:         "0A21",
		Name:       "roulette_claim_prize",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0A25"] = common.PacketConstruction{
		ID:         "0A25",
		Name:       "achievement_get_reward",
		Format:     "V",
		FieldNames: []string{"achievementID"},
	}
	packets["0A2E"] = common.PacketConstruction{
		ID:         "0A2E",
		Name:       "send_change_title",
		Format:     "V",
		FieldNames: []string{"ID"},
	}
	packets["0A39"] = common.PacketConstruction{
		ID:         "0A39",
		Name:       "char_create",
		Format:     "a24 C v4 C",
		FieldNames: []string{"name", "slot", "hair_color", "hair_style", "job_id", "unknown", "sex"},
	}
	packets["0A52"] = common.PacketConstruction{
		ID:         "0A52",
		Name:       "captcha_register",
		Format:     "Z16 v",
		FieldNames: []string{"answer", "image_size"},
	}
	packets["0A54"] = common.PacketConstruction{
		ID:         "0A54",
		Name:       "captcha_upload_request_ack",
		Format:     "v Z4 a*",
		FieldNames: []string{"len", "captcha_key", "captcha_image"},
	}
	packets["0A56"] = common.PacketConstruction{
		ID:         "0A56",
		Name:       "macro_reporter_ack",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["0A5A"] = common.PacketConstruction{
		ID:         "0A5A",
		Name:       "macro_detector_download",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0A5C"] = common.PacketConstruction{
		ID:         "0A5C",
		Name:       "macro_detector_answer",
		Format:     "Z16",
		FieldNames: []string{"answer"},
	}
	packets["0A68"] = common.PacketConstruction{
		ID:         "0A68",
		Name:       "open_ui_request",
		Format:     "C",
		FieldNames: []string{"UIType"},
	}
	packets["0A69"] = common.PacketConstruction{
		ID:         "0A69",
		Name:       "captcha_preview_request",
		Format:     "V",
		FieldNames: []string{"captcha_key"},
	}
	packets["0A6C"] = common.PacketConstruction{
		ID:         "0A6C",
		Name:       "macro_reporter_select",
		Format:     "v2 C",
		FieldNames: []string{"x", "y", "range"},
	}
	packets["0A6E"] = common.PacketConstruction{
		ID:         "0A6E",
		Name:       "rodex_send_mail",
		Format:     "v Z24 Z24 V2 v v V a* a*",
		FieldNames: []string{"len", "receiver", "sender", "zeny1", "zeny2", "title_len", "body_len", "char_id", "title", "body"},
	}
	packets["0A76"] = common.PacketConstruction{
		ID:         "0A76",
		Name:       "master_login",
		Format:     "V Z40 a32 v",
		FieldNames: []string{"version", "username", "password_rijndael", "master_version"},
	}
	packets["0A97"] = common.PacketConstruction{
		ID:         "0A97",
		Name:       "equip_switch_add",
		Format:     "a2 V",
		FieldNames: []string{"ID", "position"},
	}
	packets["0A99"] = common.PacketConstruction{
		ID:         "0A99",
		Name:       "equip_switch_remove",
		Format:     "a2",
		FieldNames: []string{"ID"},
	}
	packets["0A9C"] = common.PacketConstruction{
		ID:         "0A9C",
		Name:       "equip_switch_run",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0AAC"] = common.PacketConstruction{
		ID:         "0AAC",
		Name:       "master_login",
		Format:     "V Z30 a32 C",
		FieldNames: []string{"version", "username", "password_hex", "master_version"},
	}
	packets["0AC0"] = common.PacketConstruction{
		ID:         "0AC0",
		Name:       "rodex_open_mailbox",
		Format:     "C V6",
		FieldNames: []string{"type", "mailID1", "mailID2", "mailReturnID1", "mailReturnID2", "mailAccountID1", "mailAccountID2"},
	}
	packets["0AC1"] = common.PacketConstruction{
		ID:         "0AC1",
		Name:       "rodex_refresh_maillist",
		Format:     "C V6",
		FieldNames: []string{"type", "mailID1", "mailID2", "mailReturnID1", "mailReturnID2", "mailAccountID1", "mailAccountID2"},
	}
	packets["0ACE"] = common.PacketConstruction{
		ID:         "0ACE",
		Name:       "equip_switch_single",
		Format:     "a2",
		FieldNames: []string{"ID"},
	}
	packets["0ACF"] = common.PacketConstruction{
		ID:         "0ACF",
		Name:       "master_login",
		Format:     "a4 Z25 a32 a5",
		FieldNames: []string{"game_code", "username", "password_rijndael", "flag"},
	}
	packets["0AE8"] = common.PacketConstruction{
		ID:         "0AE8",
		Name:       "change_dress",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0AEF"] = common.PacketConstruction{
		ID:         "0AEF",
		Name:       "attendance_reward_request",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0AF4"] = common.PacketConstruction{
		ID:         "0AF4",
		Name:       "skill_use_location",
		Format:     "v4 C",
		FieldNames: []string{"lv", "skillID", "x", "y", "unknown"},
	}
	packets["0B10"] = common.PacketConstruction{
		ID:         "0B10",
		Name:       "start_skill_use",
		Format:     "v2 a4",
		FieldNames: []string{"skillID", "lv", "targetID"},
	}
	packets["0B11"] = common.PacketConstruction{
		ID:         "0B11",
		Name:       "stop_skill_use",
		Format:     "v",
		FieldNames: []string{"skillID"},
	}
	packets["0B14"] = common.PacketConstruction{
		ID:         "0B14",
		Name:       "inventory_expansion_request",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0B19"] = common.PacketConstruction{
		ID:         "0B19",
		Name:       "inventory_expansion_rejected",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0B1C"] = common.PacketConstruction{
		ID:         "0B1C",
		Name:       "ping",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0B21"] = common.PacketConstruction{
		ID:         "0B21",
		Name:       "hotkey_change",
		Format:     "v2 C V v",
		FieldNames: []string{"tab", "idx", "type", "id", "lvl"},
	}

	return packets
}