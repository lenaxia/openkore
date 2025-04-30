// Package servers provides server-specific packet constructions for different server types.
package servers

import (
	"github.com/lenaxia/goKore/network/common"
)

// ServerType0PacketConstructions provides packet constructions for ServerType0
func ServerType0PacketConstructions() map[string]common.PacketConstruction {
	packets := map[string]common.PacketConstruction{
		// Login-related packet constructions
		"0064": {
			ID:         "0064",
			Name:       "master_login",
			Format:     "v a24 a24 C",
			FieldNames: []string{"version", "username", "password", "clienttype"},
		},
		"0065": {
			ID:         "0065",
			Name:       "game_login",
			Format:     "a4 a4 a4 v C",
			FieldNames: []string{"accountID", "sessionID", "sessionID2", "userLevel", "accountSex"},
		},
		"0069": {
			ID:         "0069",
			Name:       "account_server_info",
			Format:     "v a4 a4 a4 a4 a26 C a*",
			FieldNames: []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"},
		},
		"006A": {
			ID:         "006A",
			Name:       "login_error",
			Format:     "C Z20",
			FieldNames: []string{"type", "date"},
		},
		"0066": {
			ID:         "0066",
			Name:       "char_login",
			Format:     "C",
			FieldNames: []string{"slot"},
		},
		"0072": {
			ID:         "0072",
			Name:       "map_login",
			Format:     "a4 a4 a4 V C",
			FieldNames: []string{"accountID", "charID", "sessionID", "tick", "sex"},
		},
		"0073": {
			ID:         "0073",
			Name:       "map_loaded",
			Format:     "V a3 C2",
			FieldNames: []string{"syncMapSync", "coords", "xSize", "ySize"},
		},
		"0067": {
			ID:         "0067",
			Name:       "char_create",
			Format:     "a24 C7 v2",
			FieldNames: []string{"name", "str", "agi", "vit", "int", "dex", "luk", "slot", "hair_color", "hair_style"},
		},
		"0068": {
			ID:         "0068",
			Name:       "char_delete",
			Format:     "a4 a40",
			FieldNames: []string{"charID", "email"},
		},
		"007D": {
			ID:         "007D",
			Name:       "map_loaded",
			Format:     "",
			FieldNames: []string{},
		},
		"007E": {
			ID:         "007E",
			Name:       "sync",
			Format:     "V",
			FieldNames: []string{"time"},
		},
		"0089": {
			ID:         "0089",
			Name:       "actor_action",
			Format:     "a4 C",
			FieldNames: []string{"targetID", "type"},
		},
		"008C": {
			ID:         "008C",
			Name:       "public_chat",
			Format:     "v Z*",
			FieldNames: []string{"len", "message"},
		},
		"0090": {
			ID:         "0090",
			Name:       "npc_talk",
			Format:     "a4 C",
			FieldNames: []string{"ID", "type"},
		},
		"0096": {
			ID:         "0096",
			Name:       "private_message",
			Format:     "x2 Z24 Z*",
			FieldNames: []string{"privMsgUser", "privMsg"},
		},
		"0099": {
			ID:         "0099",
			Name:       "gm_broadcast",
			Format:     "v Z*",
			FieldNames: []string{"len", "message"},
		},
		"009B": {
			ID:         "009B",
			Name:       "actor_look_at",
			Format:     "v C",
			FieldNames: []string{"head", "body"},
		},
		"009F": {
			ID:         "009F",
			Name:       "item_take",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"00A2": {
			ID:         "00A2",
			Name:       "item_drop",
			Format:     "a2 v",
			FieldNames: []string{"ID", "amount"},
		},
		"00A7": {
			ID:         "00A7",
			Name:       "item_use",
			Format:     "a2 a4",
			FieldNames: []string{"ID", "targetID"},
		},
		"00A9": {
			ID:         "00A9",
			Name:       "send_equip",
			Format:     "a2 v",
			FieldNames: []string{"ID", "type"},
		},
		"00AB": {
			ID:         "00AB",
			Name:       "send_unequip_item",
			Format:     "a2",
			FieldNames: []string{"ID"},
		},
		"00B8": {
			ID:         "00B8",
			Name:       "npc_talk_response",
			Format:     "a4 C",
			FieldNames: []string{"ID", "response"},
		},
		"00B9": {
			ID:         "00B9",
			Name:       "npc_talk_continue",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"00BB": {
			ID:         "00BB",
			Name:       "send_add_status_point",
			Format:     "v2",
			FieldNames: []string{"statusID", "Amount"},
		},
		"00C1": {
			ID:         "00C1",
			Name:       "request_user_count",
			Format:     "",
			FieldNames: []string{},
		},
		"00CF": {
			ID:         "00CF",
			Name:       "ignore_player",
			Format:     "Z24 C",
			FieldNames: []string{"name", "flag"},
		},
		"00D0": {
			ID:         "00D0",
			Name:       "ignore_all",
			Format:     "C",
			FieldNames: []string{"flag"},
		},
		"00D3": {
			ID:         "00D3",
			Name:       "get_ignore_list",
			Format:     "",
			FieldNames: []string{},
		},
		"00D5": {
			ID:         "00D5",
			Name:       "chat_room_create",
			Format:     "v2 C Z8 a*",
			FieldNames: []string{"len", "limit", "public", "password", "title"},
		},
		"00D9": {
			ID:         "00D9",
			Name:       "chat_room_join",
			Format:     "a4 Z8",
			FieldNames: []string{"ID", "password"},
		},
		"00DE": {
			ID:         "00DE",
			Name:       "chat_room_change",
			Format:     "v2 C Z8 a*",
			FieldNames: []string{"len", "limit", "public", "password", "title"},
		},
		"00E0": {
			ID:         "00E0",
			Name:       "chat_room_bestow",
			Format:     "V Z24",
			FieldNames: []string{"role", "name"},
		},
		"00E2": {
			ID:         "00E2",
			Name:       "chat_room_kick",
			Format:     "Z24",
			FieldNames: []string{"name"},
		},
		"00E3": {
			ID:         "00E3",
			Name:       "chat_room_leave",
			Format:     "",
			FieldNames: []string{},
		},
		"00E8": {
			ID:         "00E8",
			Name:       "deal_item_add",
			Format:     "a2 V",
			FieldNames: []string{"ID", "amount"},
		},
		"00ED": {
			ID:         "00ED",
			Name:       "deal_cancel",
			Format:     "",
			FieldNames: []string{},
		},
		"00F3": {
			ID:         "00F3",
			Name:       "storage_item_add",
			Format:     "a2 V",
			FieldNames: []string{"ID", "amount"},
		},
		"00F5": {
			ID:         "00F5",
			Name:       "storage_item_remove",
			Format:     "a2 V",
			FieldNames: []string{"ID", "amount"},
		},
		"00F7": {
			ID:         "00F7",
			Name:       "storage_close",
			Format:     "",
			FieldNames: []string{},
		},
		"00FC": {
			ID:         "00FC",
			Name:       "party_join_request",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"00FF": {
			ID:         "00FF",
			Name:       "party_join",
			Format:     "a4 V",
			FieldNames: []string{"ID", "flag"},
		},
		"0100": {
			ID:         "0100",
			Name:       "party_leave",
			Format:     "",
			FieldNames: []string{},
		},
		"0102": {
			ID:         "0102",
			Name:       "party_setting",
			Format:     "V",
			FieldNames: []string{"exp"},
		},
		"0103": {
			ID:         "0103",
			Name:       "party_kick",
			Format:     "a4 Z24",
			FieldNames: []string{"ID", "name"},
		},
		"0108": {
			ID:         "0108",
			Name:       "party_chat",
			Format:     "x2 Z*",
			FieldNames: []string{"message"},
		},
		"0112": {
			ID:         "0112",
			Name:       "send_add_skill_point",
			Format:     "v",
			FieldNames: []string{"skillID"},
		},
		"0113": {
			ID:         "0113",
			Name:       "skill_use",
			Format:     "v2 a4",
			FieldNames: []string{"lv", "skillID", "targetID"},
		},
		"0116": {
			ID:         "0116",
			Name:       "skill_use_location",
			Format:     "v4",
			FieldNames: []string{"lv", "skillID", "x", "y"},
		},
		"011B": {
			ID:         "011B",
			Name:       "warp_select",
			Format:     "v Z16",
			FieldNames: []string{"skillID", "mapName"},
		},
		"0126": {
			ID:         "0126",
			Name:       "cart_add",
			Format:     "a2 V",
			FieldNames: []string{"ID", "amount"},
		},
		"0127": {
			ID:         "0127",
			Name:       "cart_get",
			Format:     "a2 V",
			FieldNames: []string{"ID", "amount"},
		},
		"0128": {
			ID:         "0128",
			Name:       "storage_to_cart",
			Format:     "a2 V",
			FieldNames: []string{"ID", "amount"},
		},
		"0129": {
			ID:         "0129",
			Name:       "cart_to_storage",
			Format:     "a2 V",
			FieldNames: []string{"ID", "amount"},
		},
		"012A": {
			ID:         "012A",
			Name:       "companion_release",
			Format:     "",
			FieldNames: []string{},
		},
		"012E": {
			ID:         "012E",
			Name:       "shop_close",
			Format:     "",
			FieldNames: []string{},
		},
		"0130": {
			ID:         "0130",
			Name:       "send_entering_vending",
			Format:     "a4",
			FieldNames: []string{"accountID"},
		},
		"0134": {
			ID:         "0134",
			Name:       "buy_bulk_vender",
			Format:     "x2 a4 a*",
			FieldNames: []string{"venderID", "itemInfo"},
		},
		"0143": {
			ID:         "0143",
			Name:       "npc_talk_number",
			Format:     "a4 V",
			FieldNames: []string{"ID", "value"},
		},
		"0146": {
			ID:         "0146",
			Name:       "npc_talk_cancel",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"0149": {
			ID:         "0149",
			Name:       "alignment",
			Format:     "a4 C v",
			FieldNames: []string{"targetID", "type", "point"},
		},
		"014D": {
			ID:         "014D",
			Name:       "guild_check",
			Format:     "",
			FieldNames: []string{},
		},
		"014F": {
			ID:         "014F",
			Name:       "guild_info_request",
			Format:     "V",
			FieldNames: []string{"type"},
		},
		"0151": {
			ID:         "0151",
			Name:       "guild_emblem_request",
			Format:     "a4",
			FieldNames: []string{"guildID"},
		},
		"0159": {
			ID:         "0159",
			Name:       "guild_leave",
			Format:     "a4 a4 a4 Z40",
			FieldNames: []string{"guildID", "accountID", "charID", "reason"},
		},
		"015B": {
			ID:         "015B",
			Name:       "guild_kick",
			Format:     "a4 a4 a4 Z40",
			FieldNames: []string{"guildID", "accountID", "charID", "reason"},
		},
		"015D": {
			ID:         "015D",
			Name:       "guild_break",
			Format:     "a4",
			FieldNames: []string{"guildName"},
		},
		"0165": {
			ID:         "0165",
			Name:       "guild_create",
			Format:     "a4 Z24",
			FieldNames: []string{"charID", "guildName"},
		},
		"0168": {
			ID:         "0168",
			Name:       "guild_join_request",
			Format:     "a4 a4 a4",
			FieldNames: []string{"ID", "accountID", "charID"},
		},
		"016B": {
			ID:         "016B",
			Name:       "guild_join",
			Format:     "a4 V",
			FieldNames: []string{"ID", "flag"},
		},
		"016E": {
			ID:         "016E",
			Name:       "guild_notice",
			Format:     "a4 Z60 Z120",
			FieldNames: []string{"guildID", "name", "notice"},
		},
		"0170": {
			ID:         "0170",
			Name:       "guild_alliance_request",
			Format:     "a4 a4 a4",
			FieldNames: []string{"targetAccountID", "accountID", "charID"},
		},
		"0172": {
			ID:         "0172",
			Name:       "guild_alliance_reply",
			Format:     "a4 V",
			FieldNames: []string{"ID", "flag"},
		},
		"0178": {
			ID:         "0178",
			Name:       "identify",
			Format:     "a2",
			FieldNames: []string{"ID"},
		},
		"017E": {
			ID:         "017E",
			Name:       "guild_chat",
			Format:     "x2 Z*",
			FieldNames: []string{"message"},
		},
		"018A": {
			ID:         "018A",
			Name:       "quit_request",
			Format:     "v",
			FieldNames: []string{"type"},
		},
		"018E": {
			ID:         "018E",
			Name:       "make_item_request",
			Format:     "v4",
			FieldNames: []string{"nameID", "material_nameID1", "material_nameID2", "material_nameID3"},
		},
		"0190": {
			ID:         "0190",
			Name:       "skill_use_location_text",
			Format:     "v5 Z80",
			FieldNames: []string{"lvl", "ID", "x", "y", "info"},
		},
		"0193": {
			ID:         "0193",
			Name:       "actor_name_request",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"01AE": {
			ID:         "01AE",
			Name:       "make_arrow",
			Format:     "v",
			FieldNames: []string{"nameID"},
		},
		"01AF": {
			ID:         "01AF",
			Name:       "change_cart",
			Format:     "v",
			FieldNames: []string{"lvl"},
		},
		"01B2": {
			ID:         "01B2",
			Name:       "shop_open",
			Format:     "v a80 C a*",
			FieldNames: []string{"len", "title", "result", "vendingInfo"},
		},
		"01C0": {
			ID:         "01C0",
			Name:       "request_remain_time",
			Format:     "",
			FieldNames: []string{},
		},
		"01CE": {
			ID:         "01CE",
			Name:       "auto_spell",
			Format:     "V",
			FieldNames: []string{"ID"},
		},
		"01D5": {
			ID:         "01D5",
			Name:       "npc_talk_text",
			Format:     "v a4 Z*",
			FieldNames: []string{"len", "ID", "text"},
		},
		"01D7": {
			ID:         "01D7",
			Name:       "sprite_change",
			Format:     "a4 C V2",
			FieldNames: []string{"ID", "type", "value1", "value2"},
		},
		"01D8": {
			ID:         "01D8",
			Name:       "actor_exists",
			Format:     "a4 v14 a4 a2 v2 C2 a3 C3 v",
			FieldNames: []string{"ID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "shield", "lowhead", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "unknown1", "unknown2", "act", "lv"},
		},
		"01D9": {
			ID:         "01D9",
			Name:       "actor_connected",
			Format:     "a4 v14 a4 a2 v2 C2 a3 C2 v",
			FieldNames: []string{"ID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "shield", "lowhead", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "unknown1", "unknown2", "lv"},
		},
		"01DA": {
			ID:         "01DA",
			Name:       "actor_moved",
			Format:     "a4 v9 V v5 a4 a2 v2 C2 a6 C2 v",
			FieldNames: []string{"ID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "shield", "lowhead", "tick", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "unknown1", "unknown2", "lv"},
		},
		"01DB": {
			ID:         "01DB",
			Name:       "secure_login_key_request",
			Format:     "",
			FieldNames: []string{},
		},
		"01DD": {
			ID:         "01DD",
			Name:       "master_login",
			Format:     "V Z24 a16 C",
			FieldNames: []string{"version", "username", "password_salted_md5", "master_version"},
		},
		"01E7": {
			ID:         "01E7",
			Name:       "novice_dori_dori",
			Format:     "",
			FieldNames: []string{},
		},
		"01ED": {
			ID:         "01ED",
			Name:       "novice_explosion_spirits",
			Format:     "",
			FieldNames: []string{},
		},
		"01FD": {
			ID:         "01FD",
			Name:       "repair_item",
			Format:     "v2 C a8",
			FieldNames: []string{"index", "nameID", "upgrade", "cards"},
		},
		"0202": {
			ID:         "0202",
			Name:       "friend_request",
			Format:     "a*",
			FieldNames: []string{"username"},
		},
		"0203": {
			ID:         "0203",
			Name:       "friend_remove",
			Format:     "a4 a4",
			FieldNames: []string{"accountID", "charID"},
		},
		"0204": {
			ID:         "0204",
			Name:       "client_hash",
			Format:     "a16",
			FieldNames: []string{"hash"},
		},
		"0208": {
			ID:         "0208",
			Name:       "friend_response",
			Format:     "a4 a4 V",
			FieldNames: []string{"friendAccountID", "friendCharID", "type"},
		},
		"0217": {
			ID:         "0217",
			Name:       "rank_blacksmith",
			Format:     "",
			FieldNames: []string{},
		},
		"0218": {
			ID:         "0218",
			Name:       "rank_alchemist",
			Format:     "",
			FieldNames: []string{},
		},
		"021D": {
			ID:         "021D",
			Name:       "less_effect",
			Format:     "",
			FieldNames: []string{},
		},
		"0222": {
			ID:         "0222",
			Name:       "refine_item",
			Format:     "V",
			FieldNames: []string{"ID"},
		},
		"0225": {
			ID:         "0225",
			Name:       "rank_taekwon",
			Format:     "",
			FieldNames: []string{},
		},
		"022D": {
			ID:         "022D",
			Name:       "homunculus_command",
			Format:     "v C",
			FieldNames: []string{"commandType", "commandID"},
		},
		"0231": {
			ID:         "0231",
			Name:       "homunculus_name",
			Format:     "a24",
			FieldNames: []string{"name"},
		},
		"0232": {
			ID:         "0232",
			Name:       "actor_move",
			Format:     "a4 a3",
			FieldNames: []string{"ID", "coords"},
		},
		"0233": {
			ID:         "0233",
			Name:       "slave_attack",
			Format:     "a4 a4 C",
			FieldNames: []string{"slaveID", "targetID", "flag"},
		},
		"0234": {
			ID:         "0234",
			Name:       "slave_move_to_master",
			Format:     "a4",
			FieldNames: []string{"slaveID"},
		},
		"0237": {
			ID:         "0237",
			Name:       "rank_killer",
			Format:     "",
			FieldNames: []string{},
		},
		"023B": {
			ID:         "023B",
			Name:       "storage_password",
			Format:     "v a*",
			FieldNames: []string{"type", "data"},
		},
		"023F": {
			ID:         "023F",
			Name:       "mailbox_open",
			Format:     "",
			FieldNames: []string{},
		},
		"0241": {
			ID:         "0241",
			Name:       "mail_read",
			Format:     "V",
			FieldNames: []string{"mailID"},
		},
		"0243": {
			ID:         "0243",
			Name:       "mail_delete",
			Format:     "V",
			FieldNames: []string{"mailID"},
		},
		"0244": {
			ID:         "0244",
			Name:       "mail_attachment_get",
			Format:     "V",
			FieldNames: []string{"mailID"},
		},
		"0246": {
			ID:         "0246",
			Name:       "mail_remove",
			Format:     "v",
			FieldNames: []string{"flag"},
		},
		"0247": {
			ID:         "0247",
			Name:       "mail_attachment_set",
			Format:     "a2 V",
			FieldNames: []string{"ID", "amount"},
		},
		"0248": {
			ID:         "0248",
			Name:       "mail_send",
			Format:     "v Z24 a40 C a*",
			FieldNames: []string{"len", "recipient", "title", "body_len", "body"},
		},
		"024B": {
			ID:         "024B",
			Name:       "auction_add_item_cancel",
			Format:     "v",
			FieldNames: []string{"flag"},
		},
		"024C": {
			ID:         "024C",
			Name:       "auction_add_item",
			Format:     "a2 V",
			FieldNames: []string{"ID", "amount"},
		},
		"024D": {
			ID:         "024D",
			Name:       "auction_create",
			Format:     "V V v",
			FieldNames: []string{"now_price", "max_price", "delete_time"},
		},
		"024E": {
			ID:         "024E",
			Name:       "auction_cancel",
			Format:     "V",
			FieldNames: []string{"ID"},
		},
		"024F": {
			ID:         "024F",
			Name:       "auction_buy",
			Format:     "V V",
			FieldNames: []string{"ID", "price"},
		},
		"0251": {
			ID:         "0251",
			Name:       "auction_search",
			Format:     "v V Z24 v",
			FieldNames: []string{"type", "price", "search_string", "page"},
		},
		"0254": {
			ID:         "0254",
			Name:       "starplace_agree",
			Format:     "C",
			FieldNames: []string{"flag"},
		},
		"025B": {
			ID:         "025B",
			Name:       "cook_request",
			Format:     "v2",
			FieldNames: []string{"type", "nameID"},
		},
		"025C": {
			ID:         "025C",
			Name:       "auction_info_self",
			Format:     "v",
			FieldNames: []string{"type"},
		},
		"025D": {
			ID:         "025D",
			Name:       "auction_sell_stop",
			Format:     "V",
			FieldNames: []string{"ID"},
		},
		"0273": {
			ID:         "0273",
			Name:       "mail_return",
			Format:     "V Z24",
			FieldNames: []string{"mailID", "sender"},
		},
		"0275": {
			ID:         "0275",
			Name:       "game_login",
			Format:     "a4 a4 a4 v C Z16 V",
			FieldNames: []string{"accountID", "sessionID", "sessionID2", "userLevel", "accountSex", "mac", "iAccountSID"},
		},
		"0292": {
			ID:         "0292",
			Name:       "auto_revive",
			Format:     "",
			FieldNames: []string{},
		},
		"02B0": {
			ID:         "02B0",
			Name:       "master_login",
			Format:     "V Z24 a24 C Z16 Z14 C",
			FieldNames: []string{"version", "username", "password_rijndael", "master_version", "ip", "mac", "isGravityID"},
		},
		"02C4": {
			ID:         "02C4",
			Name:       "party_join_request_by_name",
			Format:     "Z24",
			FieldNames: []string{"playerName"},
		},
		"035F": {
			ID:         "035F",
			Name:       "character_move",
			Format:     "a4 a6",
			FieldNames: []string{"ID", "coords"},
		},
		"0365": {
			ID:         "0365",
			Name:       "item_pickup",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"0436": {
			ID:         "0436",
			Name:       "map_login",
			Format:     "a4 a4 a4 V C",
			FieldNames: []string{"accountID", "charID", "sessionID", "tick", "sex"},
		},
		"0437": {
			ID:         "0437",
			Name:       "character_select",
			Format:     "C",
			FieldNames: []string{"slot"},
		},
		"0438": {
			ID:         "0438",
			Name:       "skill_use_location",
			Format:     "v4",
			FieldNames: []string{"lv", "skillID", "x", "y"},
		},
		"0448": {
			ID:         "0448",
			Name:       "blocking_play_cancel_ack",
			Format:     "",
			FieldNames: []string{},
		},
		"07D7": {
			ID:         "07D7",
			Name:       "party_setting",
			Format:     "V C2",
			FieldNames: []string{"exp", "itemPickup", "itemDivision"},
		},
		"07DA": {
			ID:         "07DA",
			Name:       "party_leader_change",
			Format:     "a4",
			FieldNames: []string{"accountID"},
		},
		"07E4": {
			ID:         "07E4",
			Name:       "item_list_window_selected",
			Format:     "v V V a*",
			FieldNames: []string{"len", "type", "act", "itemInfo"},
		},
		"0288": {
			ID:         "0288",
			Name:       "cash_dealer_buy",
			Format:     "v2 V",
			FieldNames: []string{"itemid", "amount", "kafra_points"},
		},
		"0085": {
			ID:         "0085",
			Name:       "move_to",
			Format:     "v3",
			FieldNames: []string{"x", "y", "unknown"},
		},
		"00B2": {
			ID:         "00B2",
			Name:       "restart",
			Format:     "C",
			FieldNames: []string{"type"},
		},
		// Pet-related packet constructions
		"019F": {
			ID:         "019F",
			Name:       "pet_capture",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"01A1": {
			ID:         "01A1",
			Name:       "pet_menu",
			Format:     "C",
			FieldNames: []string{"action"},
		},
		"01A5": {
			ID:         "01A5",
			Name:       "pet_name",
			Format:     "a24",
			FieldNames: []string{"name"},
		},
		"01A7": {
			ID:         "01A7",
			Name:       "pet_hatch",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"01A9": {
			ID:         "01A9",
			Name:       "pet_emotion",
			Format:     "C",
			FieldNames: []string{"ID"},
		},
		"09FB": {
			ID:         "09FB",
			Name:       "pet_evolution",
			Format:     "a4 a*",
			FieldNames: []string{"ID", "itemInfo"},
		},
		// Mercenary-related packet constructions
		"029F": {
			ID:         "029F",
			Name:       "mercenary_command",
			Format:     "C",
			FieldNames: []string{"flag"},
		},
		"02A5": {
			ID:         "02A5",
			Name:       "companion_release",
			Format:     "",
			FieldNames: []string{},
		},
		// Battle-related packet constructions
		"02D6": {
			ID:         "02D6",
			Name:       "view_player_equip_request",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"00BF": {
			ID:         "00BF",
			Name:       "send_emotion",
			Format:     "C",
			FieldNames: []string{"ID"},
		},
		"01B9": {
			ID:         "01B9",
			Name:       "novice_dori_dori",
			Format:     "",
			FieldNames: []string{},
		},
		"01BA": {
			ID:         "01BA",
			Name:       "novice_explosion_spirits",
			Format:     "",
			FieldNames: []string{},
		},
		"02CF": {
			ID:         "02CF",
			Name:       "memorial_dungeon_command",
			Format:     "V",
			FieldNames: []string{"command"},
		},
		// Marriage-related packet constructions
		"01F7": {
			ID:         "01F7",
			Name:       "adopt_request",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"01F9": {
			ID:         "01F9",
			Name:       "adopt_reply_request",
			Format:     "a4 a4 C",
			FieldNames: []string{"parentID1", "parentID2", "result"},
		},
		// Auction-related packet constructions
		"0366": {
			ID:         "0366",
			Name:       "auction_add_item_cancel",
			Format:     "C",
			FieldNames: []string{"flag"},
		},
		"0367": {
			ID:         "0367",
			Name:       "auction_add_item",
			Format:     "a4 v",
			FieldNames: []string{"ID", "amount"},
		},
		"0368": {
			ID:         "0368",
			Name:       "auction_create",
			Format:     "V V V",
			FieldNames: []string{"now_price", "max_price", "delete_time"},
		},
		"0369": {
			ID:         "0369",
			Name:       "auction_cancel",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"036A": {
			ID:         "036A",
			Name:       "auction_buy",
			Format:     "a4 V",
			FieldNames: []string{"ID", "price"},
		},
		"036B": {
			ID:         "036B",
			Name:       "auction_search",
			Format:     "C V a24 v",
			FieldNames: []string{"type", "price", "search_string", "page"},
		},
		"036C": {
			ID:         "036C",
			Name:       "auction_info_self",
			Format:     "C",
			FieldNames: []string{"type"},
		},
		"036D": {
			ID:         "036D",
			Name:       "auction_sell_stop",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		// Buying Store-related packet constructions
		"00C8": {
			ID:         "00C8",
			Name:       "buy_bulk",
			Format:     "a*",
			FieldNames: []string{"buyInfo"},
		},
		"00C9": {
			ID:         "00C9",
			Name:       "sell_bulk",
			Format:     "a*",
			FieldNames: []string{"sellInfo"},
		},
		"0835": {
			ID:         "0835",
			Name:       "search_store_close",
			Format:     "",
			FieldNames: []string{},
		},
		"0838": {
			ID:         "0838",
			Name:       "search_store_info",
			Format:     "C V V v v a*",
			FieldNames: []string{"type", "max_price", "min_price", "item_count", "card_count", "item_card_list"},
		},
		"0839": {
			ID:         "0839",
			Name:       "search_store_request_next_page",
			Format:     "",
			FieldNames: []string{},
		},
		"083B": {
			ID:         "083B",
			Name:       "search_store_select",
			Format:     "a4 a4 v",
			FieldNames: []string{"accountID", "storeID", "nameID"},
		},
		// UI-related packet constructions
		"02D8": {
			ID:         "02D8",
			Name:       "misc_config_set",
			Format:     "V V",
			FieldNames: []string{"type", "flag"},
		},
		"02F1": {
			ID:         "02F1",
			Name:       "notify_progress_bar_complete",
			Format:     "",
			FieldNames: []string{},
		},
		// Note: view_player_equip_request (02D6) is already defined in the battle system
		"0AA1": {
			ID:         "0AA1",
			Name:       "refineui_select",
			Format:     "v",
			FieldNames: []string{"index"},
		},
		"0AA3": {
			ID:         "0AA3",
			Name:       "refineui_refine",
			Format:     "v v C",
			FieldNames: []string{"index", "catalyst", "bless"},
		},
		"0AA4": {
			ID:         "0AA4",
			Name:       "refineui_close",
			Format:     "",
			FieldNames: []string{},
		},
		"09A4": {
			ID:         "09A4",
			Name:       "item_list_window_selected",
			Format:     "v C C a*",
			FieldNames: []string{"len", "type", "act", "itemInfo"},
		},
		"011D": {
			ID:         "011D",
			Name:       "memo_request",
			Format:     "",
			FieldNames: []string{},
		},
		"0A18": {
			ID:         "0A18",
			Name:       "stylist_change",
			Format:     "v v v v v v",
			FieldNames: []string{"hair_color", "hair_style", "cloth_color", "head_top", "head_mid", "head_bottom"},
		},
		"0A68": {
			ID:         "0A68",
			Name:       "open_ui_request",
			Format:     "C",
			FieldNames: []string{"UIType"},
		},
		"0AEF": {
			ID:         "0AEF",
			Name:       "attendance_reward_request",
			Format:     "",
			FieldNames: []string{},
		},
		"0A19": {
			ID:         "0A19",
			Name:       "roulette_window_open",
			Format:     "",
			FieldNames: []string{},
		},
		"0A1B": {
			ID:         "0A1B",
			Name:       "roulette_info_request",
			Format:     "",
			FieldNames: []string{},
		},
		"0A1D": {
			ID:         "0A1D",
			Name:       "roulette_close",
			Format:     "",
			FieldNames: []string{},
		},
		"0A1F": {
			ID:         "0A1F",
			Name:       "roulette_start",
			Format:     "",
			FieldNames: []string{},
		},
		"0A21": {
			ID:         "0A21",
			Name:       "roulette_claim_prize",
			Format:     "",
			FieldNames: []string{},
		},
		"02B6": {
			ID:         "02B6",
			Name:       "send_quest_state",
			Format:     "V C",
			FieldNames: []string{"questID", "state"},
		},
		// Deal-related packet constructions
		"00E9": {
			ID:         "00E9",
			Name:       "deal_item_add",
			Format:     "a2 v",
			FieldNames: []string{"ID", "amount"},
		},
		"00C5": {
			ID:         "00C5",
			Name:       "deal_initiate",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"00C7": {
			ID:         "00C7",
			Name:       "deal_reply",
			Format:     "C",
			FieldNames: []string{"action"},
		},
		"00EB": {
			ID:         "00EB",
			Name:       "deal_finalize",
			Format:     "",
			FieldNames: []string{},
		},
		"00E6": {
			ID:         "00E6",
			Name:       "deal_cancel",
			Format:     "",
			FieldNames: []string{},
		},
		"00EF": {
			ID:         "00EF",
			Name:       "deal_trade",
			Format:     "",
			FieldNames: []string{},
		},
		// Ranking-related packet constructions
		"0A26": {
			ID:         "0A26",
			Name:       "achievement_get_reward",
			Format:     "V",
			FieldNames: []string{"achievementID"},
		},
		"0197": {
			ID:         "0197",
			Name:       "rank_alchemist",
			Format:     "",
			FieldNames: []string{},
		},
		"0198": {
			ID:         "0198",
			Name:       "rank_blacksmith",
			Format:     "",
			FieldNames: []string{},
		},
		"0199": {
			ID:         "0199",
			Name:       "rank_killer",
			Format:     "",
			FieldNames: []string{},
		},
		"019A": {
			ID:         "019A",
			Name:       "rank_taekwon",
			Format:     "",
			FieldNames: []string{},
		},
		"097C": {
			ID:         "097C",
			Name:       "rank_general",
			Format:     "C",
			FieldNames: []string{"type"},
		},
		// GM-related packet constructions
		"0094": {
			ID:         "0094",
			Name:       "gm_summon_player",
			Format:     "Z24",
			FieldNames: []string{"playerName"},
		},
		"00CC": {
			ID:         "00CC",
			Name:       "gm_kick",
			Format:     "a4",
			FieldNames: []string{"targetAccountID"},
		},
		"00CD": {
			ID:         "00CD",
			Name:       "gm_kick_all",
			Format:     "",
			FieldNames: []string{},
		},
		"013F": {
			ID:         "013F",
			Name:       "gm_item_mob_create",
			Format:     "Z24",
			FieldNames: []string{"name"},
		},
		// Note: Using a different ID for gm_move_to_map to avoid conflict with novice_dori_dori (01B9)
		"01BD": {
			ID:         "01BD",
			Name:       "gm_move_to_map",
			Format:     "Z16 v v",
			FieldNames: []string{"mapName", "x", "y"},
		},
		// Note: Using a different ID for gm_reset_state_skill to avoid conflict with rank_alchemist (0197)
		"0196": {
			ID:         "0196",
			Name:       "gm_reset_state_skill",
			Format:     "C",
			FieldNames: []string{"type"},
		},
		// Note: Using a different ID for gm_change_cell_type to avoid conflict with rank_blacksmith (0198) and rank_killer (0199)
		"019C": {
			ID:         "019C",
			Name:       "gm_change_cell_type",
			Format:     "v v C",
			FieldNames: []string{"x", "y", "type"},
		},
		"019B": {
			ID:         "019B",
			Name:       "gm_change_effect_state",
			Format:     "V",
			FieldNames: []string{"effect_state"},
		},
		// Note: Using a different ID for gm_remove to avoid conflict with novice_explosion_spirits (01BA)
		"01BE": {
			ID:         "01BE",
			Name:       "gm_remove",
			Format:     "Z24",
			FieldNames: []string{"playerName"},
		},
		"01BB": {
			ID:         "01BB",
			Name:       "gm_shift",
			Format:     "Z24",
			FieldNames: []string{"playerName"},
		},
		"01BC": {
			ID:         "01BC",
			Name:       "gm_recall",
			Format:     "Z24",
			FieldNames: []string{"playerName"},
		},
		"0212": {
			ID:         "0212",
			Name:       "manner_by_name",
			Format:     "Z24",
			FieldNames: []string{"playerName"},
		},
		"0213": {
			ID:         "0213",
			Name:       "gm_request_status",
			Format:     "Z24",
			FieldNames: []string{"playerName"},
		},
		"01DF": {
			ID:         "01DF",
			Name:       "gm_request_account_name",
			Format:     "a4",
			FieldNames: []string{"targetID"},
		},
		"0187": {
			ID:         "0187",
			Name:       "ban_check",
			Format:     "a4",
			FieldNames: []string{"accountID"},
		},
		// Macro-related packet constructions
		"0871": {
			ID:         "0871",
			Name:       "macro_start",
			Format:     "",
			FieldNames: []string{},
		},
		"0872": {
			ID:         "0872",
			Name:       "macro_stop",
			Format:     "",
			FieldNames: []string{},
		},
		"0A5A": {
			ID:         "0A5A",
			Name:       "macro_detector_download",
			Format:     "",
			FieldNames: []string{},
		},
		"0A5C": {
			ID:         "0A5C",
			Name:       "macro_detector_answer",
			Format:     "a*",
			FieldNames: []string{"answer"},
		},
		// Note: Using a different ID for req_cash_tabcode to avoid conflict with open_ui_request (0A68)
		"0A69": {
			ID:         "0A69",
			Name:       "req_cash_tabcode",
			Format:     "v",
			FieldNames: []string{"ID"},
		},
		// Captcha-related packet constructions
		"07E7": {
			ID:         "07E7",
			Name:       "captcha_answer",
			Format:     "v a4 s",
			FieldNames: []string{"len", "accountID", "answer"},
		},
		"07E5": {
			ID:         "07E5",
			Name:       "captcha_preview_request",
			Format:     "V",
			FieldNames: []string{"captcha_key"},
		},
		// Card-related packet constructions
		"017A": {
			ID:         "017A",
			Name:       "card_merge_request",
			Format:     "a4",
			FieldNames: []string{"cardID"},
		},
		"017C": {
			ID:         "017C",
			Name:       "card_merge",
			Format:     "a4 a4",
			FieldNames: []string{"cardID", "itemID"},
		},
		// Cash shop-related packet constructions
		// Note: Using a different ID for request_cashitems to avoid conflict with open_ui_request (0A68)
		"0A6A": {
			ID:         "0A6A",
			Name:       "request_cashitems",
			Format:     "",
			FieldNames: []string{},
		},
		"0844": {
			ID:         "0844",
			Name:       "cash_shop_open",
			Format:     "",
			FieldNames: []string{},
		},
		"0845": {
			ID:         "0845",
			Name:       "cash_shop_close",
			Format:     "",
			FieldNames: []string{},
		},
		"0848": {
			ID:         "0848",
			Name:       "cash_shop_buy",
			Format:     "V v a*",
			FieldNames: []string{"kafra_points", "count", "buy_info"},
		},
		"0972": {
			ID:         "0972",
			Name:       "merge_item_request",
			Format:     "v a*",
			FieldNames: []string{"len", "itemList"},
		},
		"0974": {
			ID:         "0974",
			Name:       "merge_item_cancel",
			Format:     "",
			FieldNames: []string{},
		},
		// Miscellaneous packet constructions
		"0825": {
			ID:         "0825",
			Name:       "token_login",
			Format:     "v V V a24 a24 a24 a17 a15 a*",
			FieldNames: []string{"len", "version", "master_version", "username", "password", "password_rijndael", "mac", "ip", "token"},
		},
		"0A37": {
			ID:         "0A37",
			Name:       "request_remain_time",
			Format:     "",
			FieldNames: []string{},
		},
		"0842": {
			ID:         "0842",
			Name:       "recall_sso",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"0843": {
			ID:         "0843",
			Name:       "remove_aid_sso",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"0B0D": {
			ID:         "0B0D",
			Name:       "starplace_agree",
			Format:     "C",
			FieldNames: []string{"flag"},
		},
		"09F1": {
			ID:         "09F1",
			Name:       "sync_request_ex",
			Format:     "v",
			FieldNames: []string{"syncID"},
		},
		// Info-related packet constructions
		// Note: Using a different ID for actor_info_request to avoid conflicts
		"0A90": {
			ID:         "0A90",
			Name:       "actor_info_request",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"0095": {
			ID:         "0095",
			Name:       "actor_name_request",
			Format:     "a4",
			FieldNames: []string{"ID"},
		},
		"01C1": {
			ID:         "01C1",
			Name:       "request_user_count",
			Format:     "",
			FieldNames: []string{},
		},
		"02DB": {
			ID:         "02DB",
			Name:       "battleground_chat",
			Format:     "Z*",
			FieldNames: []string{"message"},
		},
		"0B01": {
			ID:         "0B01",
			Name:       "clan_chat",
			Format:     "v Z*",
			FieldNames: []string{"len", "message"},
		},
		// More packet constructions can be added here as needed
		// This file will grow as more packet constructions are defined
	}

	// Add the following packet constructions for packets that are not yet implemented
	// Character-related packet constructions
	packets["006B"] = common.PacketConstruction{
		ID:         "006B",
		Name:       "received_characters_info",
		Format:     "v C3 x20 a*",
		FieldNames: []string{"len", "total_slot", "premium_start_slot", "premium_end_slot", "charInfo"},
	}
	packets["006C"] = common.PacketConstruction{
		ID:         "006C",
		Name:       "login_error_game_login_server",
		Format:     "",
		FieldNames: []string{},
	}
	packets["006D"] = common.PacketConstruction{
		ID:         "006D",
		Name:       "character_creation_successful",
		Format:     "a*",
		FieldNames: []string{"charInfo"},
	}
	packets["006E"] = common.PacketConstruction{
		ID:         "006E",
		Name:       "character_creation_failed",
		Format:     "C",
		FieldNames: []string{"type"},
	}
	packets["006F"] = common.PacketConstruction{
		ID:         "006F",
		Name:       "character_deletion_successful",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0070"] = common.PacketConstruction{
		ID:         "0070",
		Name:       "character_deletion_failed",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0071"] = common.PacketConstruction{
		ID:         "0071",
		Name:       "received_character_ID_and_Map",
		Format:     "a4 Z16 a4 v",
		FieldNames: []string{"charID", "mapName", "mapIP", "mapPort"},
	}
	packets["0074"] = common.PacketConstruction{
		ID:         "0074",
		Name:       "map_load_error",
		Format:     "C",
		FieldNames: []string{"error"},
	}
	packets["0075"] = common.PacketConstruction{
		ID:         "0075",
		Name:       "changeToInGameState",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0077"] = common.PacketConstruction{
		ID:         "0077",
		Name:       "changeToInGameState",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0078"] = common.PacketConstruction{
		ID:         "0078",
		Name:       "actor_exists",
		Format:     "a4 v14 a4 a2 v2 C2 a3 C3 v",
		FieldNames: []string{"ID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "lowhead", "shield", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "unknown1", "unknown2", "act", "lv"},
	}
	packets["0079"] = common.PacketConstruction{
		ID:         "0079",
		Name:       "actor_connected",
		Format:     "a4 v14 a4 a2 v2 C2 a3 C2 v",
		FieldNames: []string{"ID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "lowhead", "shield", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "unknown1", "unknown2", "lv"},
	}
	packets["007A"] = common.PacketConstruction{
		ID:         "007A",
		Name:       "changeToInGameState",
		Format:     "",
		FieldNames: []string{},
	}
	packets["007B"] = common.PacketConstruction{
		ID:         "007B",
		Name:       "actor_moved",
		Format:     "a4 v8 V v6 a4 a2 v2 C2 a6 C2 v",
		FieldNames: []string{"ID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "lowhead", "tick", "shield", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "unknown1", "unknown2", "lv"},
	}
	packets["007C"] = common.PacketConstruction{
		ID:         "007C",
		Name:       "actor_spawned",
		Format:     "a4 v14 C2 a3 C2",
		FieldNames: []string{"ID", "walk_speed", "opt1", "opt2", "option", "hair_style", "weapon", "lowhead", "type", "shield", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "stance", "sex", "coords", "unknown1", "unknown2"},
	}
	packets["0080"] = common.PacketConstruction{
		ID:         "0080",
		Name:       "actor_died_or_disappeared",
		Format:     "a4 C",
		FieldNames: []string{"ID", "type"},
	}
	packets["0081"] = common.PacketConstruction{
		ID:         "0081",
		Name:       "errors",
		Format:     "C",
		FieldNames: []string{"type"},
	}
	packets["0086"] = common.PacketConstruction{
		ID:         "0086",
		Name:       "actor_display",
		Format:     "a4 a6 V",
		FieldNames: []string{"ID", "coords", "tick"},
	}
	packets["0087"] = common.PacketConstruction{
		ID:         "0087",
		Name:       "character_moves",
		Format:     "a4 a6",
		FieldNames: []string{"move_start_time", "coords"},
	}
	packets["0088"] = common.PacketConstruction{
		ID:         "0088",
		Name:       "actor_movement_interrupted",
		Format:     "a4 v2",
		FieldNames: []string{"ID", "x", "y"},
	}
	packets["008E"] = common.PacketConstruction{
		ID:         "008E",
		Name:       "self_chat",
		Format:     "v Z*",
		FieldNames: []string{"len", "message"},
	}
	packets["0091"] = common.PacketConstruction{
		ID:         "0091",
		Name:       "map_change",
		Format:     "Z16 v2",
		FieldNames: []string{"map", "x", "y"},
	}
	packets["0092"] = common.PacketConstruction{
		ID:         "0092",
		Name:       "map_changed",
		Format:     "Z16 v2 a4 v",
		FieldNames: []string{"map", "x", "y", "IP", "port"},
	}
	packets["0095"] = common.PacketConstruction{
		ID:         "0095",
		Name:       "actor_info",
		Format:     "a4 Z24",
		FieldNames: []string{"ID", "name"},
	}
	packets["0098"] = common.PacketConstruction{
		ID:         "0098",
		Name:       "private_message_sent",
		Format:     "C",
		FieldNames: []string{"type"},
	}
	packets["009A"] = common.PacketConstruction{
		ID:         "009A",
		Name:       "system_chat",
		Format:     "v a*",
		FieldNames: []string{"len", "message"},
	}
	packets["009D"] = common.PacketConstruction{
		ID:         "009D",
		Name:       "item_exists",
		Format:     "a4 v C v3 C2",
		FieldNames: []string{"ID", "nameID", "identified", "x", "y", "amount", "subx", "suby"},
	}
	packets["009E"] = common.PacketConstruction{
		ID:         "009E",
		Name:       "item_appeared",
		Format:     "a4 v C v2 C2 v",
		FieldNames: []string{"ID", "nameID", "identified", "x", "y", "subx", "suby", "amount"},
	}
	packets["00A0"] = common.PacketConstruction{
		ID:         "00A0",
		Name:       "inventory_item_added",
		Format:     "a2 v2 C3 a8 v C2",
		FieldNames: []string{"ID", "amount", "nameID", "identified", "broken", "upgrade", "cards", "type_equip", "type", "fail"},
	}
	packets["00A1"] = common.PacketConstruction{
		ID:         "00A1",
		Name:       "item_disappeared",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["00A3"] = common.PacketConstruction{
		ID:         "00A3",
		Name:       "inventory_items_stackable",
		Format:     "v a*",
		FieldNames: []string{"len", "itemInfo"},
	}
	packets["00A4"] = common.PacketConstruction{
		ID:         "00A4",
		Name:       "inventory_items_nonstackable",
		Format:     "v a*",
		FieldNames: []string{"len", "itemInfo"},
	}
	packets["00A5"] = common.PacketConstruction{
		ID:         "00A5",
		Name:       "storage_items_stackable",
		Format:     "v a*",
		FieldNames: []string{"len", "itemInfo"},
	}
	packets["00A6"] = common.PacketConstruction{
		ID:         "00A6",
		Name:       "storage_items_nonstackable",
		Format:     "v a*",
		FieldNames: []string{"len", "itemInfo"},
	}
	packets["00A8"] = common.PacketConstruction{
		ID:         "00A8",
		Name:       "use_item",
		Format:     "a2 v C",
		FieldNames: []string{"ID", "amount", "success"},
	}
	packets["00AA"] = common.PacketConstruction{
		ID:         "00AA",
		Name:       "equip_item",
		Format:     "a2 v C",
		FieldNames: []string{"ID", "type", "success"},
	}
	packets["00AC"] = common.PacketConstruction{
		ID:         "00AC",
		Name:       "unequip_item",
		Format:     "a2 v C",
		FieldNames: []string{"ID", "type", "success"},
	}
	packets["00AF"] = common.PacketConstruction{
		ID:         "00AF",
		Name:       "inventory_item_removed",
		Format:     "a2 v",
		FieldNames: []string{"ID", "amount"},
	}
	packets["00B0"] = common.PacketConstruction{
		ID:         "00B0",
		Name:       "stat_info",
		Format:     "v V",
		FieldNames: []string{"type", "val"},
	}
	packets["00B1"] = common.PacketConstruction{
		ID:         "00B1",
		Name:       "stat_info",
		Format:     "v V",
		FieldNames: []string{"type", "val"},
	}
	packets["00B3"] = common.PacketConstruction{
		ID:         "00B3",
		Name:       "switch_character",
		Format:     "C",
		FieldNames: []string{"result"},
	}
	packets["00B6"] = common.PacketConstruction{
		ID:         "00B6",
		Name:       "npc_talk_close",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["00B7"] = common.PacketConstruction{
		ID:         "00B7",
		Name:       "npc_talk_responses",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00BC"] = common.PacketConstruction{
		ID:         "00BC",
		Name:       "stats_added",
		Format:     "v x C",
		FieldNames: []string{"type", "val"},
	}
	packets["00BD"] = common.PacketConstruction{
		ID:         "00BD",
		Name:       "stats_info",
		Format:     "v C12 v14",
		FieldNames: []string{"points_free", "str", "points_str", "agi", "points_agi", "vit", "points_vit", "int", "points_int", "dex", "points_dex", "luk", "points_luk", "attack", "attack_bonus", "attack_magic_min", "attack_magic_max", "def", "def_bonus", "def_magic", "def_magic_bonus", "hit", "flee", "flee_bonus", "critical", "stance", "manner"},
	}
	packets["00BE"] = common.PacketConstruction{
		ID:         "00BE",
		Name:       "stat_info",
		Format:     "v C",
		FieldNames: []string{"type", "val"},
	}
	packets["00C0"] = common.PacketConstruction{
		ID:         "00C0",
		Name:       "emoticon",
		Format:     "a4 C",
		FieldNames: []string{"ID", "type"},
	}
	packets["00C2"] = common.PacketConstruction{
		ID:         "00C2",
		Name:       "users_online",
		Format:     "V",
		FieldNames: []string{"users"},
	}
	packets["00C3"] = common.PacketConstruction{
		ID:         "00C3",
		Name:       "sprite_change",
		Format:     "a4 C2",
		FieldNames: []string{"ID", "type", "value1"},
	}
	packets["00C4"] = common.PacketConstruction{
		ID:         "00C4",
		Name:       "npc_store_begin",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["00C6"] = common.PacketConstruction{
		ID:         "00C6",
		Name:       "npc_store_info",
		Format:     "v a*",
		FieldNames: []string{"len", "itemList"},
	}
	packets["00C7"] = common.PacketConstruction{
		ID:         "00C7",
		Name:       "npc_sell_list",
		Format:     "v a*",
		FieldNames: []string{"len", "itemsdata"},
	}
	packets["00CA"] = common.PacketConstruction{
		ID:         "00CA",
		Name:       "buy_result",
		Format:     "C",
		FieldNames: []string{"fail"},
	}
	packets["00CB"] = common.PacketConstruction{
		ID:         "00CB",
		Name:       "sell_result",
		Format:     "C",
		FieldNames: []string{"fail"},
	}

	// Adding missing packets from servertype0.packets
	packets["0144"] = common.PacketConstruction{
		ID:         "0144",
		Name:       "minimap_indicator",
		Format:     "a4 V3 C5",
		FieldNames: []string{"npcID", "type", "x", "y", "ID", "blue", "green", "red", "alpha"},
	}
	packets["0145"] = common.PacketConstruction{
		ID:         "0145",
		Name:       "npc_image",
		Format:     "Z16 C",
		FieldNames: []string{"npc_image", "type"},
	}
	packets["0147"] = common.PacketConstruction{
		ID:         "0147",
		Name:       "item_skill",
		Format:     "v6 A*",
		FieldNames: []string{"skillID", "targetType", "unknown", "skillLv", "sp", "unknown2", "skillName"},
	}
	packets["0148"] = common.PacketConstruction{
		ID:         "0148",
		Name:       "resurrection",
		Format:     "a4 v",
		FieldNames: []string{"targetID", "type"},
	}
	packets["014A"] = common.PacketConstruction{
		ID:         "014A",
		Name:       "manner_message",
		Format:     "V",
		FieldNames: []string{"type"},
	}
	packets["014B"] = common.PacketConstruction{
		ID:         "014B",
		Name:       "GM_silence",
		Format:     "C Z24",
		FieldNames: []string{"type", "name"},
	}
	packets["096D"] = common.PacketConstruction{
		ID:         "096D",
		Name:       "merge_item_open",
		Format:     "v a*",
		FieldNames: []string{"length", "itemList"},
	}
	packets["096F"] = common.PacketConstruction{
		ID:         "096F",
		Name:       "merge_item_result",
		Format:     "a2 v C",
		FieldNames: []string{"itemIndex", "total", "result"},
	}
	packets["0975"] = common.PacketConstruction{
		ID:         "0975",
		Name:       "storage_items_stackable",
		Format:     "v Z24 a*",
		FieldNames: []string{"len", "title", "itemInfo"},
	}
	packets["0976"] = common.PacketConstruction{
		ID:         "0976",
		Name:       "storage_items_nonstackable",
		Format:     "v Z24 a*",
		FieldNames: []string{"len", "title", "itemInfo"},
	}
	packets["0977"] = common.PacketConstruction{
		ID:         "0977",
		Name:       "monster_hp_info",
		Format:     "a4 V V",
		FieldNames: []string{"ID", "hp", "hp_max"},
	}

	packets["097A"] = common.PacketConstruction{
		ID:         "097A",
		Name:       "quest_all_list",
		Format:     "v V a*",
		FieldNames: []string{"len", "quest_amount", "message"},
	}
	packets["097B"] = common.PacketConstruction{
		ID:         "097B",
		Name:       "rates_info2",
		Format:     "s V3 a*",
		FieldNames: []string{"len", "exp", "death", "drop", "detail"},
	}
	packets["097D"] = common.PacketConstruction{
		ID:         "097D",
		Name:       "top10",
		Format:     "v a*",
		FieldNames: []string{"type", "message"},
	}
	packets["097E"] = common.PacketConstruction{
		ID:         "097E",
		Name:       "rank_points",
		Format:     "vV2",
		FieldNames: []string{"type", "points", "total"},
	}
	packets["0983"] = common.PacketConstruction{
		ID:         "0983",
		Name:       "actor_status_active",
		Format:     "v a4 C V5",
		FieldNames: []string{"type", "ID", "flag", "total", "tick", "unknown1", "unknown2", "unknown3"},
	}
	packets["0984"] = common.PacketConstruction{
		ID:         "0984",
		Name:       "actor_status_active",
		Format:     "a4 v V5",
		FieldNames: []string{"ID", "type", "total", "tick", "unknown1", "unknown2", "unknown3"},
	}
	packets["0985"] = common.PacketConstruction{
		ID:         "0985",
		Name:       "skill_post_delaylist",
		Format:     "v a*",
		FieldNames: []string{"len", "skill_list"},
	}
	packets["0988"] = common.PacketConstruction{
		ID:         "0988",
		Name:       "clan_user",
		Format:     "v2",
		FieldNames: []string{"onlineuser", "totalmembers"},
	}
	packets["098A"] = common.PacketConstruction{
		ID:         "098A",
		Name:       "clan_info",
		Format:     "v a4 Z24 Z24 Z16 C2 a*",
		FieldNames: []string{"len", "clan_ID", "clan_name", "clan_master", "clan_map", "alliance_count", "antagonist_count", "ally_antagonist_names"},
	}
	packets["098D"] = common.PacketConstruction{
		ID:         "098D",
		Name:       "clan_leave",
		Format:     "",
		FieldNames: []string{},
	}
	packets["098E"] = common.PacketConstruction{
		ID:         "098E",
		Name:       "clan_chat",
		Format:     "v Z24 Z*",
		FieldNames: []string{"len", "charname", "message"},
	}
	packets["0990"] = common.PacketConstruction{
		ID:         "0990",
		Name:       "inventory_item_added",
		Format:     "a2 v2 C3 a8 V C2 a4 v",
		FieldNames: []string{"ID", "amount", "nameID", "identified", "broken", "upgrade", "cards", "type_equip", "type", "fail", "expire", "unknown"},
	}
	packets["0991"] = common.PacketConstruction{
		ID:         "0991",
		Name:       "inventory_items_stackable",
		Format:     "v a*",
		FieldNames: []string{"len", "itemInfo"},
	}
	packets["0992"] = common.PacketConstruction{
		ID:         "0992",
		Name:       "inventory_items_nonstackable",
		Format:     "v a*",
		FieldNames: []string{"len", "itemInfo"},
	}
	packets["0993"] = common.PacketConstruction{
		ID:         "0993",
		Name:       "cart_items_stackable",
		Format:     "v a*",
		FieldNames: []string{"len", "itemInfo"},
	}
	packets["0994"] = common.PacketConstruction{
		ID:         "0994",
		Name:       "cart_items_nonstackable",
		Format:     "v a*",
		FieldNames: []string{"len", "itemInfo"},
	}
	packets["0995"] = common.PacketConstruction{
		ID:         "0995",
		Name:       "storage_items_stackable",
		Format:     "v Z24 a*",
		FieldNames: []string{"len", "title", "itemInfo"},
	}
	packets["0996"] = common.PacketConstruction{
		ID:         "0996",
		Name:       "storage_items_nonstackable",
		Format:     "v Z24 a*",
		FieldNames: []string{"len", "title", "itemInfo"},
	}
	packets["0997"] = common.PacketConstruction{
		ID:         "0997",
		Name:       "show_eq",
		Format:     "v Z24 v7 v C a*",
		FieldNames: []string{"len", "name", "jobID", "hair_style", "tophead", "midhead", "lowhead", "robe", "hair_color", "clothes_color", "sex", "equips_info"},
	}
	packets["0999"] = common.PacketConstruction{
		ID:         "0999",
		Name:       "equip_item",
		Format:     "a2 V v C",
		FieldNames: []string{"ID", "type", "viewID", "success"},
	}
	packets["099A"] = common.PacketConstruction{
		ID:         "099A",
		Name:       "unequip_item",
		Format:     "a2 V C",
		FieldNames: []string{"ID", "type", "success"},
	}
	packets["099B"] = common.PacketConstruction{
		ID:         "099B",
		Name:       "map_property3",
		Format:     "v a4",
		FieldNames: []string{"type", "info_table"},
	}
	packets["099D"] = common.PacketConstruction{
		ID:         "099D",
		Name:       "received_characters",
		Format:     "v a*",
		FieldNames: []string{"len", "charInfo"},
	}
	packets["099F"] = common.PacketConstruction{
		ID:         "099F",
		Name:       "area_spell_multiple2",
		Format:     "v a*",
		FieldNames: []string{"len", "spellInfo"},
	}
	packets["09A0"] = common.PacketConstruction{
		ID:         "09A0",
		Name:       "sync_received_characters",
		Format:     "V",
		FieldNames: []string{"sync_Count"},
	}
	packets["09A6"] = common.PacketConstruction{
		ID:         "09A6",
		Name:       "banking_check",
		Format:     "V2 v",
		FieldNames: []string{"zeny", "zeny2", "reason"},
	}
	packets["007F"] = common.PacketConstruction{
		ID:         "007F",
		Name:       "received_sync",
		Format:     "V",
		FieldNames: []string{"time"},
	}
	packets["008A"] = common.PacketConstruction{
		ID:         "008A",
		Name:       "actor_action",
		Format:     "a4 a4 a4 V2 v2 C v",
		FieldNames: []string{"sourceID", "targetID", "tick", "src_speed", "dst_speed", "damage", "div", "type", "dual_wield_damage"},
	}
	packets["008D"] = common.PacketConstruction{
		ID:         "008D",
		Name:       "public_chat",
		Format:     "v a4 Z*",
		FieldNames: []string{"len", "ID", "message"},
	}
	packets["0097"] = common.PacketConstruction{
		ID:         "0097",
		Name:       "private_message",
		Format:     "v Z24 Z*",
		FieldNames: []string{"len", "privMsgUser", "privMsg"},
	}
	packets["009C"] = common.PacketConstruction{
		ID:         "009C",
		Name:       "actor_look_at",
		Format:     "a4 v C",
		FieldNames: []string{"ID", "head", "body"},
	}
	packets["00B4"] = common.PacketConstruction{
		ID:         "00B4",
		Name:       "npc_talk",
		Format:     "v a4 Z*",
		FieldNames: []string{"len", "ID", "msg"},
	}
	packets["00B5"] = common.PacketConstruction{
		ID:         "00B5",
		Name:       "npc_talk_continue",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["00D1"] = common.PacketConstruction{
		ID:         "00D1",
		Name:       "ignore_player_result",
		Format:     "C2",
		FieldNames: []string{"type", "error"},
	}
	packets["00D2"] = common.PacketConstruction{
		ID:         "00D2",
		Name:       "ignore_all_result",
		Format:     "C2",
		FieldNames: []string{"type", "error"},
	}
	packets["00D4"] = common.PacketConstruction{
		ID:         "00D4",
		Name:       "whisper_list",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00D6"] = common.PacketConstruction{
		ID:         "00D6",
		Name:       "chat_created",
		Format:     "C",
		FieldNames: []string{"result"},
	}
	packets["00D7"] = common.PacketConstruction{
		ID:         "00D7",
		Name:       "chat_info",
		Format:     "v a4 a4 v2 C a*",
		FieldNames: []string{"len", "ownerID", "ID", "limit", "num_users", "public", "title"},
	}
	packets["00D8"] = common.PacketConstruction{
		ID:         "00D8",
		Name:       "chat_removed",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["00DA"] = common.PacketConstruction{
		ID:         "00DA",
		Name:       "chat_join_result",
		Format:     "C",
		FieldNames: []string{"type"},
	}
	packets["00DB"] = common.PacketConstruction{
		ID:         "00DB",
		Name:       "chat_users",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00DC"] = common.PacketConstruction{
		ID:         "00DC",
		Name:       "chat_user_join",
		Format:     "v Z24",
		FieldNames: []string{"num_users", "user"},
	}
	packets["00DD"] = common.PacketConstruction{
		ID:         "00DD",
		Name:       "chat_user_leave",
		Format:     "v Z24 C",
		FieldNames: []string{"num_users", "user", "flag"},
	}
	packets["00DF"] = common.PacketConstruction{
		ID:         "00DF",
		Name:       "chat_modified",
		Format:     "v a4 a4 v2 C a*",
		FieldNames: []string{"len", "ownerID", "ID", "limit", "num_users", "public", "title"},
	}
	packets["00E1"] = common.PacketConstruction{
		ID:         "00E1",
		Name:       "chat_newowner",
		Format:     "C x3 Z24",
		FieldNames: []string{"type", "user"},
	}
	packets["00E5"] = common.PacketConstruction{
		ID:         "00E5",
		Name:       "deal_request",
		Format:     "Z24",
		FieldNames: []string{"user"},
	}
	packets["00E7"] = common.PacketConstruction{
		ID:         "00E7",
		Name:       "deal_begin",
		Format:     "C",
		FieldNames: []string{"type"},
	}
	packets["00E9"] = common.PacketConstruction{
		ID:         "00E9",
		Name:       "deal_add_other",
		Format:     "V v C3 a8",
		FieldNames: []string{"amount", "nameID", "identified", "broken", "upgrade", "cards"},
	}
	packets["00EA"] = common.PacketConstruction{
		ID:         "00EA",
		Name:       "deal_add_you",
		Format:     "a2 C",
		FieldNames: []string{"ID", "fail"},
	}
	packets["00EC"] = common.PacketConstruction{
		ID:         "00EC",
		Name:       "deal_finalize",
		Format:     "C",
		FieldNames: []string{"type"},
	}
	packets["00EE"] = common.PacketConstruction{
		ID:         "00EE",
		Name:       "deal_cancelled",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00F0"] = common.PacketConstruction{
		ID:         "00F0",
		Name:       "deal_complete",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00F2"] = common.PacketConstruction{
		ID:         "00F2",
		Name:       "storage_opened",
		Format:     "v2",
		FieldNames: []string{"items", "items_max"},
	}
	packets["00F4"] = common.PacketConstruction{
		ID:         "00F4",
		Name:       "storage_item_added",
		Format:     "a2 V v C3 a8",
		FieldNames: []string{"ID", "amount", "nameID", "identified", "broken", "upgrade", "cards"},
	}
	packets["00F6"] = common.PacketConstruction{
		ID:         "00F6",
		Name:       "storage_item_removed",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["00F8"] = common.PacketConstruction{
		ID:         "00F8",
		Name:       "storage_closed",
		Format:     "",
		FieldNames: []string{},
	}
	packets["00FA"] = common.PacketConstruction{
		ID:         "00FA",
		Name:       "party_organize_result",
		Format:     "C",
		FieldNames: []string{"fail"},
	}
	packets["00FB"] = common.PacketConstruction{
		ID:         "00FB",
		Name:       "party_users_info",
		Format:     "v Z24 a*",
		FieldNames: []string{"len", "party_name", "playerInfo"},
	}
	packets["00FD"] = common.PacketConstruction{
		ID:         "00FD",
		Name:       "party_invite_result",
		Format:     "Z24 C",
		FieldNames: []string{"name", "type"},
	}
	packets["00FE"] = common.PacketConstruction{
		ID:         "00FE",
		Name:       "party_invite",
		Format:     "a4 Z24",
		FieldNames: []string{"ID", "name"},
	}
	packets["0101"] = common.PacketConstruction{
		ID:         "0101",
		Name:       "party_exp",
		Format:     "V",
		FieldNames: []string{"type"},
	}
	packets["0104"] = common.PacketConstruction{
		ID:         "0104",
		Name:       "party_join",
		Format:     "a4 V v2 C Z24 Z24 Z16",
		FieldNames: []string{"ID", "role", "x", "y", "type", "name", "user", "map"},
	}
	packets["0105"] = common.PacketConstruction{
		ID:         "0105",
		Name:       "party_leave",
		Format:     "a4 Z24 C",
		FieldNames: []string{"ID", "name", "result"},
	}
	packets["0106"] = common.PacketConstruction{
		ID:         "0106",
		Name:       "party_hp_info",
		Format:     "a4 v2",
		FieldNames: []string{"ID", "hp", "hp_max"},
	}
	packets["0107"] = common.PacketConstruction{
		ID:         "0107",
		Name:       "party_location",
		Format:     "a4 v2",
		FieldNames: []string{"ID", "x", "y"},
	}
	packets["0108"] = common.PacketConstruction{
		ID:         "0108",
		Name:       "item_upgrade",
		Format:     "v a2 v",
		FieldNames: []string{"type", "ID", "upgrade"},
	}
	packets["0109"] = common.PacketConstruction{
		ID:         "0109",
		Name:       "party_chat",
		Format:     "v a4 Z*",
		FieldNames: []string{"len", "ID", "message"},
	}
	packets["010A"] = common.PacketConstruction{
		ID:         "010A",
		Name:       "mvp_item",
		Format:     "v",
		FieldNames: []string{"itemID"},
	}
	packets["010B"] = common.PacketConstruction{
		ID:         "010B",
		Name:       "mvp_you",
		Format:     "V",
		FieldNames: []string{"expAmount"},
	}
	packets["010C"] = common.PacketConstruction{
		ID:         "010C",
		Name:       "mvp_other",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["010E"] = common.PacketConstruction{
		ID:         "010E",
		Name:       "skill_update",
		Format:     "v4 C",
		FieldNames: []string{"skillID", "lv", "sp", "range", "up"},
	}
	packets["010F"] = common.PacketConstruction{
		ID:         "010F",
		Name:       "skills_list",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0111"] = common.PacketConstruction{
		ID:         "0111",
		Name:       "skill_add",
		Format:     "v V v3 Z24 C",
		FieldNames: []string{"skillID", "target", "lv", "sp", "range", "name", "upgradable"},
	}
	packets["0114"] = common.PacketConstruction{
		ID:         "0114",
		Name:       "skill_use",
		Format:     "v a4 a4 V3 v3 C",
		FieldNames: []string{"skillID", "sourceID", "targetID", "tick", "src_speed", "dst_speed", "damage", "level", "option", "type"},
	}
	packets["0117"] = common.PacketConstruction{
		ID:         "0117",
		Name:       "skill_use_location",
		Format:     "v a4 v3 V",
		FieldNames: []string{"skillID", "sourceID", "lv", "x", "y", "tick"},
	}
	packets["0119"] = common.PacketConstruction{
		ID:         "0119",
		Name:       "character_status",
		Format:     "a4 v3 C",
		FieldNames: []string{"ID", "opt1", "opt2", "option", "stance"},
	}
	packets["011A"] = common.PacketConstruction{
		ID:         "011A",
		Name:       "skill_used_no_damage",
		Format:     "v2 a4 a4 C",
		FieldNames: []string{"skillID", "amount", "targetID", "sourceID", "success"},
	}
	packets["011C"] = common.PacketConstruction{
		ID:         "011C",
		Name:       "warp_portal_list",
		Format:     "v Z16 Z16 Z16 Z16",
		FieldNames: []string{"type", "memo1", "memo2", "memo3", "memo4"},
	}
	packets["011E"] = common.PacketConstruction{
		ID:         "011E",
		Name:       "memo_success",
		Format:     "C",
		FieldNames: []string{"fail"},
	}
	packets["011F"] = common.PacketConstruction{
		ID:         "011F",
		Name:       "area_spell",
		Format:     "a4 a4 v2 C2",
		FieldNames: []string{"ID", "sourceID", "x", "y", "type", "isVisible"},
	}
	packets["0120"] = common.PacketConstruction{
		ID:         "0120",
		Name:       "area_spell_disappears",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["0121"] = common.PacketConstruction{
		ID:         "0121",
		Name:       "cart_info",
		Format:     "v2 V2",
		FieldNames: []string{"items", "items_max", "weight", "weight_max"},
	}
	packets["0122"] = common.PacketConstruction{
		ID:         "0122",
		Name:       "cart_items_nonstackable",
		Format:     "v a*",
		FieldNames: []string{"len", "itemInfo"},
	}
	packets["0123"] = common.PacketConstruction{
		ID:         "0123",
		Name:       "cart_items_stackable",
		Format:     "v a*",
		FieldNames: []string{"len", "itemInfo"},
	}
	packets["0124"] = common.PacketConstruction{
		ID:         "0124",
		Name:       "cart_item_added",
		Format:     "a2 V v C3 a8",
		FieldNames: []string{"ID", "amount", "nameID", "identified", "broken", "upgrade", "cards"},
	}
	packets["0125"] = common.PacketConstruction{
		ID:         "0125",
		Name:       "cart_item_removed",
		Format:     "a2 V",
		FieldNames: []string{"ID", "amount"},
	}
	packets["012B"] = common.PacketConstruction{
		ID:         "012B",
		Name:       "cart_off",
		Format:     "",
		FieldNames: []string{},
	}
	packets["012C"] = common.PacketConstruction{
		ID:         "012C",
		Name:       "cart_add_failed",
		Format:     "C",
		FieldNames: []string{"fail"},
	}
	packets["012D"] = common.PacketConstruction{
		ID:         "012D",
		Name:       "shop_skill",
		Format:     "v",
		FieldNames: []string{"number"},
	}
	packets["0131"] = common.PacketConstruction{
		ID:         "0131",
		Name:       "vender_found",
		Format:     "a4 Z80",
		FieldNames: []string{"ID", "title"},
	}
	packets["0132"] = common.PacketConstruction{
		ID:         "0132",
		Name:       "vender_lost",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["0133"] = common.PacketConstruction{
		ID:         "0133",
		Name:       "vender_items_list",
		Format:     "v a4 a*",
		FieldNames: []string{"len", "venderID", "itemList"},
	}
	packets["0135"] = common.PacketConstruction{
		ID:         "0135",
		Name:       "vender_buy_fail",
		Format:     "v2 C",
		FieldNames: []string{"ID", "amount", "fail"},
	}
	packets["0136"] = common.PacketConstruction{
		ID:         "0136",
		Name:       "vending_start",
		Format:     "v a4 a*",
		FieldNames: []string{"len", "accountID", "itemList"},
	}
	packets["0137"] = common.PacketConstruction{
		ID:         "0137",
		Name:       "shop_sold",
		Format:     "v2",
		FieldNames: []string{"number", "amount"},
	}
	packets["0139"] = common.PacketConstruction{
		ID:         "0139",
		Name:       "monster_ranged_attack",
		Format:     "a4 v5",
		FieldNames: []string{"ID", "sourceX", "sourceY", "targetX", "targetY", "range"},
	}
	packets["013A"] = common.PacketConstruction{
		ID:         "013A",
		Name:       "attack_range",
		Format:     "v",
		FieldNames: []string{"type"},
	}
	packets["013B"] = common.PacketConstruction{
		ID:         "013B",
		Name:       "arrow_none",
		Format:     "v",
		FieldNames: []string{"type"},
	}
	packets["013C"] = common.PacketConstruction{
		ID:         "013C",
		Name:       "arrow_equipped",
		Format:     "a2",
		FieldNames: []string{"ID"},
	}
	packets["013D"] = common.PacketConstruction{
		ID:         "013D",
		Name:       "hp_sp_changed",
		Format:     "v2",
		FieldNames: []string{"type", "amount"},
	}
	packets["013E"] = common.PacketConstruction{
		ID:         "013E",
		Name:       "skill_cast",
		Format:     "a4 a4 v5 V",
		FieldNames: []string{"sourceID", "targetID", "x", "y", "skillID", "unknown", "type", "wait"},
	}
	packets["0141"] = common.PacketConstruction{
		ID:         "0141",
		Name:       "stat_info2",
		Format:     "V2 l",
		FieldNames: []string{"type", "val", "val2"},
	}
	packets["0142"] = common.PacketConstruction{
		ID:         "0142",
		Name:       "npc_talk_number",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["014C"] = common.PacketConstruction{
		ID:         "014C",
		Name:       "guild_allies_enemy_list",
		Format:     "",
		FieldNames: []string{},
	}
	packets["014E"] = common.PacketConstruction{
		ID:         "014E",
		Name:       "guild_master_member",
		Format:     "V",
		FieldNames: []string{"type"},
	}
	packets["0150"] = common.PacketConstruction{
		ID:         "0150",
		Name:       "guild_info",
		Format:     "a4 V9 a4 Z24 Z24 Z16 V",
		FieldNames: []string{"ID", "lv", "conMember", "maxMember", "average", "exp", "exp_next", "tax", "tendency_left_right", "tendency_down_up", "emblemID", "name", "master", "castles_string", "zeny"},
	}
	packets["0152"] = common.PacketConstruction{
		ID:         "0152",
		Name:       "guild_emblem",
		Format:     "v a4 a4 a*",
		FieldNames: []string{"len", "guildID", "emblemID", "emblem"},
	}
	packets["0154"] = common.PacketConstruction{
		ID:         "0154",
		Name:       "guild_members_list",
		Format:     "v a*",
		FieldNames: []string{"len", "member_list"},
	}
	packets["0156"] = common.PacketConstruction{
		ID:         "0156",
		Name:       "guild_update_member_position",
		Format:     "v a*",
		FieldNames: []string{"len", "member_list"},
	}
	packets["015A"] = common.PacketConstruction{
		ID:         "015A",
		Name:       "guild_leave",
		Format:     "Z24 Z40",
		FieldNames: []string{"name", "message"},
	}
	packets["015C"] = common.PacketConstruction{
		ID:         "015C",
		Name:       "guild_expulsion",
		Format:     "Z24 Z40 Z24",
		FieldNames: []string{"name", "message", "accountName"},
	}
	packets["015E"] = common.PacketConstruction{
		ID:         "015E",
		Name:       "guild_broken",
		Format:     "V",
		FieldNames: []string{"flag"},
	}
	packets["0160"] = common.PacketConstruction{
		ID:         "0160",
		Name:       "guild_member_setting_list",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0162"] = common.PacketConstruction{
		ID:         "0162",
		Name:       "guild_skills_list",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0163"] = common.PacketConstruction{
		ID:         "0163",
		Name:       "guild_expulsion_list",
		Format:     "v a*",
		FieldNames: []string{"len", "expulsion_list"},
	}
	packets["0166"] = common.PacketConstruction{
		ID:         "0166",
		Name:       "guild_members_title_list",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0167"] = common.PacketConstruction{
		ID:         "0167",
		Name:       "guild_create_result",
		Format:     "C",
		FieldNames: []string{"type"},
	}
	packets["0169"] = common.PacketConstruction{
		ID:         "0169",
		Name:       "guild_invite_result",
		Format:     "C",
		FieldNames: []string{"type"},
	}
	packets["016A"] = common.PacketConstruction{
		ID:         "016A",
		Name:       "guild_request",
		Format:     "a4 Z24",
		FieldNames: []string{"ID", "name"},
	}
	packets["016C"] = common.PacketConstruction{
		ID:         "016C",
		Name:       "guild_name",
		Format:     "a4 a4 V C a4 Z24",
		FieldNames: []string{"guildID", "emblemID", "mode", "is_master", "interSID", "guildName"},
	}
	packets["016D"] = common.PacketConstruction{
		ID:         "016D",
		Name:       "guild_member_online_status",
		Format:     "a4 a4 V",
		FieldNames: []string{"ID", "charID", "online"},
	}
	packets["016F"] = common.PacketConstruction{
		ID:         "016F",
		Name:       "guild_notice",
		Format:     "Z60 Z120",
		FieldNames: []string{"subject", "notice"},
	}
	packets["0171"] = common.PacketConstruction{
		ID:         "0171",
		Name:       "guild_ally_request",
		Format:     "a4 Z24",
		FieldNames: []string{"ID", "guildName"},
	}
	packets["0173"] = common.PacketConstruction{
		ID:         "0173",
		Name:       "guild_alliance",
		Format:     "C",
		FieldNames: []string{"flag"},
	}
	packets["0174"] = common.PacketConstruction{
		ID:         "0174",
		Name:       "guild_position_changed",
		Format:     "v a4 a4 a4 V Z20",
		FieldNames: []string{"len", "ID", "mode", "sameID", "exp", "position_name"},
	}
	packets["0177"] = common.PacketConstruction{
		ID:         "0177",
		Name:       "identify_list",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0179"] = common.PacketConstruction{
		ID:         "0179",
		Name:       "identify",
		Format:     "a2 C",
		FieldNames: []string{"ID", "flag"},
	}
	packets["017B"] = common.PacketConstruction{
		ID:         "017B",
		Name:       "card_merge_list",
		Format:     "",
		FieldNames: []string{},
	}
	packets["017D"] = common.PacketConstruction{
		ID:         "017D",
		Name:       "card_merge_status",
		Format:     "a2 a2 C",
		FieldNames: []string{"item_index", "card_index", "fail"},
	}
	packets["017F"] = common.PacketConstruction{
		ID:         "017F",
		Name:       "guild_chat",
		Format:     "v Z*",
		FieldNames: []string{"len", "message"},
	}
	packets["0181"] = common.PacketConstruction{
		ID:         "0181",
		Name:       "guild_opposition_result",
		Format:     "C",
		FieldNames: []string{"flag"},
	}
	packets["0182"] = common.PacketConstruction{
		ID:         "0182",
		Name:       "guild_member_add",
		Format:     "a4 a4 v5 V3 Z50 Z24",
		FieldNames: []string{"ID", "charID", "hair_style", "hair_color", "sex", "jobID", "lv", "contribution", "online", "position", "memo", "name"},
	}
	packets["0184"] = common.PacketConstruction{
		ID:         "0184",
		Name:       "guild_unally",
		Format:     "a4 V",
		FieldNames: []string{"guildID", "flag"},
	}
	packets["0185"] = common.PacketConstruction{
		ID:         "0185",
		Name:       "guild_alliance_added",
		Format:     "a4 a4 Z24",
		FieldNames: []string{"opposition", "alliance_guildID", "name"},
	}
	packets["0187"] = common.PacketConstruction{
		ID:         "0187",
		Name:       "sync_request",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["0188"] = common.PacketConstruction{
		ID:         "0188",
		Name:       "item_upgrade",
		Format:     "v a2 v",
		FieldNames: []string{"type", "ID", "upgrade"},
	}
	packets["0189"] = common.PacketConstruction{
		ID:         "0189",
		Name:       "no_teleport",
		Format:     "v",
		FieldNames: []string{"fail"},
	}
	packets["018B"] = common.PacketConstruction{
		ID:         "018B",
		Name:       "quit_response",
		Format:     "v",
		FieldNames: []string{"fail"},
	}
	packets["018C"] = common.PacketConstruction{
		ID:         "018C",
		Name:       "sense_result",
		Format:     "v3 V v4 C9",
		FieldNames: []string{"nameID", "level", "size", "hp", "def", "race", "mdef", "element", "ice", "earth", "fire", "wind", "poison", "holy", "dark", "spirit", "undead"},
	}
	packets["018D"] = common.PacketConstruction{
		ID:         "018D",
		Name:       "makable_item_list",
		Format:     "v a*",
		FieldNames: []string{"len", "item_list"},
	}
	packets["018F"] = common.PacketConstruction{
		ID:         "018F",
		Name:       "refine_result",
		Format:     "v2",
		FieldNames: []string{"fail", "nameID"},
	}
	packets["0191"] = common.PacketConstruction{
		ID:         "0191",
		Name:       "talkie_box",
		Format:     "a4 Z80",
		FieldNames: []string{"ID", "message"},
	}
	packets["0192"] = common.PacketConstruction{
		ID:         "0192",
		Name:       "map_change_cell",
		Format:     "v3 Z16",
		FieldNames: []string{"x", "y", "type", "map_name"},
	}
	packets["0194"] = common.PacketConstruction{
		ID:         "0194",
		Name:       "character_name",
		Format:     "a4 Z24",
		FieldNames: []string{"ID", "name"},
	}
	packets["0195"] = common.PacketConstruction{
		ID:         "0195",
		Name:       "actor_info",
		Format:     "a4 Z24 Z24 Z24 Z24",
		FieldNames: []string{"ID", "name", "partyName", "guildName", "guildTitle"},
	}
	packets["0196"] = common.PacketConstruction{
		ID:         "0196",
		Name:       "actor_status_active",
		Format:     "v a4 C",
		FieldNames: []string{"type", "ID", "flag"},
	}
	packets["0199"] = common.PacketConstruction{
		ID:         "0199",
		Name:       "map_property",
		Format:     "v",
		FieldNames: []string{"type"},
	}
	packets["019A"] = common.PacketConstruction{
		ID:         "019A",
		Name:       "pvp_rank",
		Format:     "V3",
		FieldNames: []string{"ID", "rank", "num"},
	}
	packets["019B"] = common.PacketConstruction{
		ID:         "019B",
		Name:       "unit_levelup",
		Format:     "a4 V",
		FieldNames: []string{"ID", "type"},
	}
	packets["019E"] = common.PacketConstruction{
		ID:         "019E",
		Name:       "pet_capture_process",
		Format:     "",
		FieldNames: []string{},
	}
	packets["01A0"] = common.PacketConstruction{
		ID:         "01A0",
		Name:       "pet_capture_result",
		Format:     "C",
		FieldNames: []string{"success"},
	}
	packets["01A2"] = common.PacketConstruction{
		ID:         "01A2",
		Name:       "pet_info",
		Format:     "Z24 C v5",
		FieldNames: []string{"name", "renameflag", "level", "hungry", "friendly", "accessory", "type"},
	}
	packets["01A3"] = common.PacketConstruction{
		ID:         "01A3",
		Name:       "pet_food",
		Format:     "C v",
		FieldNames: []string{"success", "foodID"},
	}
	packets["01A4"] = common.PacketConstruction{
		ID:         "01A4",
		Name:       "pet_info2",
		Format:     "C a4 V",
		FieldNames: []string{"type", "ID", "value"},
	}
	packets["01A6"] = common.PacketConstruction{
		ID:         "01A6",
		Name:       "egg_list",
		Format:     "",
		FieldNames: []string{},
	}
	packets["01AA"] = common.PacketConstruction{
		ID:         "01AA",
		Name:       "pet_emotion",
		Format:     "a4 V",
		FieldNames: []string{"ID", "type"},
	}
	packets["01AB"] = common.PacketConstruction{
		ID:         "01AB",
		Name:       "stat_info",
		Format:     "a4 v V",
		FieldNames: []string{"ID", "type", "val"},
	}
	packets["01AC"] = common.PacketConstruction{
		ID:         "01AC",
		Name:       "actor_trapped",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["01AD"] = common.PacketConstruction{
		ID:         "01AD",
		Name:       "arrowcraft_list",
		Format:     "",
		FieldNames: []string{},
	}
	packets["01B0"] = common.PacketConstruction{
		ID:         "01B0",
		Name:       "monster_typechange",
		Format:     "a4 a V",
		FieldNames: []string{"ID", "unknown", "type"},
	}
	packets["01B3"] = common.PacketConstruction{
		ID:         "01B3",
		Name:       "npc_image",
		Format:     "Z64 C",
		FieldNames: []string{"npc_image", "type"},
	}
	packets["01B4"] = common.PacketConstruction{
		ID:         "01B4",
		Name:       "guild_emblem_update",
		Format:     "a4 a4 a2",
		FieldNames: []string{"ID", "guildID", "emblemID"},
	}
	packets["01B5"] = common.PacketConstruction{
		ID:         "01B5",
		Name:       "account_payment_info",
		Format:     "V2",
		FieldNames: []string{"D_minute", "H_minute"},
	}
	packets["01B6"] = common.PacketConstruction{
		ID:         "01B6",
		Name:       "guild_info",
		Format:     "a4 V9 a4 Z24 Z24 Z16 V",
		FieldNames: []string{"ID", "lv", "conMember", "maxMember", "average", "exp", "exp_next", "tax", "tendency_left_right", "tendency_down_up", "emblemID", "name", "master", "castles_string", "zeny"},
	}
	packets["01B9"] = common.PacketConstruction{
		ID:         "01B9",
		Name:       "cast_cancelled",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["01C1"] = common.PacketConstruction{
		ID:         "01C1",
		Name:       "remain_time_info",
		Format:     "a4 a4 a4",
		FieldNames: []string{"result", "expiration_date", "remain_time"},
	}
	packets["01C3"] = common.PacketConstruction{
		ID:         "01C3",
		Name:       "local_broadcast",
		Format:     "v V v4 Z*",
		FieldNames: []string{"len", "color", "font_type", "font_size", "font_align", "font_y", "message"},
	}
	packets["01C4"] = common.PacketConstruction{
		ID:         "01C4",
		Name:       "storage_item_added",
		Format:     "a2 V v C4 a8",
		FieldNames: []string{"ID", "amount", "nameID", "type", "identified", "broken", "upgrade", "cards"},
	}
	packets["01C5"] = common.PacketConstruction{
		ID:         "01C5",
		Name:       "cart_item_added",
		Format:     "a2 V v C4 a8",
		FieldNames: []string{"ID", "amount", "nameID", "type", "identified", "broken", "upgrade", "cards"},
	}
	packets["01C8"] = common.PacketConstruction{
		ID:         "01C8",
		Name:       "item_used",
		Format:     "a2 v a4 v C",
		FieldNames: []string{"ID", "itemID", "actorID", "remaining", "success"},
	}
	packets["01C9"] = common.PacketConstruction{
		ID:         "01C9",
		Name:       "area_spell",
		Format:     "a4 a4 v2 C2 C Z80",
		FieldNames: []string{"ID", "sourceID", "x", "y", "type", "isVisible", "scribbleLen", "scribbleMsg"},
	}
	packets["01CD"] = common.PacketConstruction{
		ID:         "01CD",
		Name:       "sage_autospell",
		Format:     "a*",
		FieldNames: []string{"autospell_list"},
	}
	packets["01CF"] = common.PacketConstruction{
		ID:         "01CF",
		Name:       "devotion",
		Format:     "a4 a20 v",
		FieldNames: []string{"sourceID", "targetIDs", "range"},
	}
	packets["01D0"] = common.PacketConstruction{
		ID:         "01D0",
		Name:       "revolving_entity",
		Format:     "a4 v",
		FieldNames: []string{"sourceID", "entity"},
	}
	packets["01D1"] = common.PacketConstruction{
		ID:         "01D1",
		Name:       "blade_stop",
		Format:     "a4 a4 V",
		FieldNames: []string{"sourceID", "targetID", "active"},
	}
	packets["01D2"] = common.PacketConstruction{
		ID:         "01D2",
		Name:       "combo_delay",
		Format:     "a4 V",
		FieldNames: []string{"ID", "delay"},
	}
	packets["01D3"] = common.PacketConstruction{
		ID:         "01D3",
		Name:       "sound_effect",
		Format:     "Z24 C V a4",
		FieldNames: []string{"name", "type", "term", "ID"},
	}
	packets["01D4"] = common.PacketConstruction{
		ID:         "01D4",
		Name:       "npc_talk_text",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["01D6"] = common.PacketConstruction{
		ID:         "01D6",
		Name:       "map_property2",
		Format:     "v",
		FieldNames: []string{"type"},
	}
	packets["01D7"] = common.PacketConstruction{
		ID:         "01D7",
		Name:       "sprite_change",
		Format:     "a4 C V2",
		FieldNames: []string{"ID", "type", "value1", "value2"},
	}
	packets["01D7"] = common.PacketConstruction{
		ID:         "01D7",
		Name:       "sprite_change",
		Format:     "a4 C V2",
		FieldNames: []string{"ID", "type", "value1", "value2"},
	}
	packets["01D7_alt"] = common.PacketConstruction{
		ID:         "01D7",
		Name:       "sprite_change",
		Format:     "a4 C v2",
		FieldNames: []string{"ID", "type", "value1", "value2"},
	}
	packets["09A8"] = common.PacketConstruction{
		ID:         "09A8",
		Name:       "banking_deposit",
		Format:     "v V2 V",
		FieldNames: []string{"reason", "zeny", "zeny2", "balance"},
	}
	packets["09AA"] = common.PacketConstruction{
		ID:         "09AA",
		Name:       "banking_withdraw",
		Format:     "v V2 V",
		FieldNames: []string{"reason", "zeny", "zeny2", "balance"},
	}
	packets["09BB"] = common.PacketConstruction{
		ID:         "09BB",
		Name:       "storage_opened",
		Format:     "v2",
		FieldNames: []string{"items", "items_max"},
	}
	packets["09BF"] = common.PacketConstruction{
		ID:         "09BF",
		Name:       "storage_closed",
		Format:     "",
		FieldNames: []string{},
	}
	packets["09CA"] = common.PacketConstruction{
		ID:         "09CA",
		Name:       "area_spell_multiple3",
		Format:     "v a*",
		FieldNames: []string{"len", "spellInfo"},
	}
	packets["09CB"] = common.PacketConstruction{
		ID:         "09CB",
		Name:       "skill_used_no_damage",
		Format:     "v V a4 a4 C",
		FieldNames: []string{"skillID", "amount", "targetID", "sourceID", "success"},
	}
	packets["09CD"] = common.PacketConstruction{
		ID:         "09CD",
		Name:       "message_string",
		Format:     "v V",
		FieldNames: []string{"index", "param"},
	}
	packets["09CF"] = common.PacketConstruction{
		ID:         "09CF",
		Name:       "gameguard_request",
		Format:     "",
		FieldNames: []string{},
	}
	packets["09D1"] = common.PacketConstruction{
		ID:         "09D1",
		Name:       "progress_bar_unit",
		Format:     "V3",
		FieldNames: []string{"GID", "color", "time"},
	}
	packets["09D5"] = common.PacketConstruction{
		ID:         "09D5",
		Name:       "npc_market_info",
		Format:     "v a*",
		FieldNames: []string{"len", "itemList"},
	}
	packets["09D7"] = common.PacketConstruction{
		ID:         "09D7",
		Name:       "npc_market_purchase_result",
		Format:     "v C a*",
		FieldNames: []string{"len", "result", "itemList"},
	}
	packets["09DA"] = common.PacketConstruction{
		ID:         "09DA",
		Name:       "guild_storage_log",
		Format:     "v3 a*",
		FieldNames: []string{"len", "result", "count", "log"},
	}
	packets["09DB"] = common.PacketConstruction{
		ID:         "09DB",
		Name:       "actor_moved",
		Format:     "v C a4 a4 v3 V v5 a4 v6 a4 a2 v V C2 a6 C2 v2 V2 C Z*",
		FieldNames: []string{"len", "object_type", "ID", "charID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "shield", "lowhead", "tick", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "costume", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "xSize", "ySize", "lv", "font", "maxHP", "HP", "isBoss", "name"},
	}
	packets["09DC"] = common.PacketConstruction{
		ID:         "09DC",
		Name:       "actor_connected",
		Format:     "v C a4 a4 v3 V v11 a4 a2 v V C2 a3 C2 v2 V2 C Z*",
		FieldNames: []string{"len", "object_type", "ID", "charID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "shield", "lowhead", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "costume", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "xSize", "ySize", "lv", "font", "maxHP", "HP", "isBoss", "name"},
	}
	packets["09DD"] = common.PacketConstruction{
		ID:         "09DD",
		Name:       "actor_exists",
		Format:     "v C a4 a4 v3 V v11 a4 a2 v V C2 a3 C3 v2 V2 C Z*",
		FieldNames: []string{"len", "object_type", "ID", "charID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "shield", "lowhead", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "costume", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "xSize", "ySize", "act", "lv", "font", "maxHP", "HP", "isBoss", "name"},
	}
	packets["09DE"] = common.PacketConstruction{
		ID:         "09DE",
		Name:       "private_message",
		Format:     "v V Z24 C Z*",
		FieldNames: []string{"len", "charID", "privMsgUser", "isAdmin", "privMsg"},
	}
	packets["09DF"] = common.PacketConstruction{
		ID:         "09DF",
		Name:       "private_message_sent",
		Format:     "C V",
		FieldNames: []string{"type", "charID"},
	}
	packets["09E5"] = common.PacketConstruction{
		ID:         "09E5",
		Name:       "shop_sold_long",
		Format:     "v2 a4 V2",
		FieldNames: []string{"number", "amount", "charID", "time", "zeny"},
	}
	packets["09E7"] = common.PacketConstruction{
		ID:         "09E7",
		Name:       "unread_rodex",
		Format:     "C",
		FieldNames: []string{"show"},
	}
	packets["09EB"] = common.PacketConstruction{
		ID:         "09EB",
		Name:       "rodex_read_mail",
		Format:     "v C V2 v V2 C",
		FieldNames: []string{"len", "type", "mailID1", "mailID2", "text_len", "zeny1", "zeny2", "itemCount"},
	}
	packets["09ED"] = common.PacketConstruction{
		ID:         "09ED",
		Name:       "rodex_write_result",
		Format:     "C",
		FieldNames: []string{"fail"},
	}
	packets["09F0"] = common.PacketConstruction{
		ID:         "09F0",
		Name:       "rodex_mail_list",
		Format:     "v C3 a*",
		FieldNames: []string{"len", "type", "amount", "isEnd", "mailList"},
	}
	packets["09F2"] = common.PacketConstruction{
		ID:         "09F2",
		Name:       "rodex_get_zeny",
		Format:     "V2 C2",
		FieldNames: []string{"mailID1", "mailID2", "type", "fail"},
	}
	packets["09F4"] = common.PacketConstruction{
		ID:         "09F4",
		Name:       "rodex_get_item",
		Format:     "V2 C2",
		FieldNames: []string{"mailID1", "mailID2", "type", "fail"},
	}
	packets["09F6"] = common.PacketConstruction{
		ID:         "09F6",
		Name:       "rodex_delete",
		Format:     "C V2",
		FieldNames: []string{"type", "mailID1", "mailID2"},
	}
	packets["09F7"] = common.PacketConstruction{
		ID:         "09F7",
		Name:       "homunculus_property",
		Format:     "Z24 C v12 V2 v2 V2 v2",
		FieldNames: []string{"name", "state", "level", "hunger", "intimacy", "accessory", "atk", "matk", "hit", "critical", "def", "mdef", "flee", "aspd", "hp", "hp_max", "sp", "sp_max", "exp", "exp_max", "points_skill", "attack_range"},
	}
	packets["09F8"] = common.PacketConstruction{
		ID:         "09F8",
		Name:       "quest_all_list",
		Format:     "v V a*",
		FieldNames: []string{"len", "quest_amount", "message"},
	}
	packets["09F9"] = common.PacketConstruction{
		ID:         "09F9",
		Name:       "quest_add",
		Format:     "V C V2 v a*",
		FieldNames: []string{"questID", "active", "time_start", "time_expire", "mission_amount", "message"},
	}
	packets["09FA"] = common.PacketConstruction{
		ID:         "09FA",
		Name:       "quest_update_mission_hunt",
		Format:     "v2 a*",
		FieldNames: []string{"len", "mission_amount", "message"},
	}
	packets["09FC"] = common.PacketConstruction{
		ID:         "09FC",
		Name:       "pet_evolution_result",
		Format:     "v V",
		FieldNames: []string{"len", "result"},
	}
	packets["09FD"] = common.PacketConstruction{
		ID:         "09FD",
		Name:       "actor_moved",
		Format:     "v C a4 a4 v3 V v2 V2 v V v6 a4 a2 v V C2 a6 C2 v2 V2 C v Z*",
		FieldNames: []string{"len", "object_type", "ID", "charID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "shield", "lowhead", "tick", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "costume", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "xSize", "ySize", "lv", "font", "maxHP", "HP", "isBoss", "opt4", "name"},
	}
	packets["09FE"] = common.PacketConstruction{
		ID:         "09FE",
		Name:       "actor_connected",
		Format:     "v C a4 a4 v3 V v2 V2 v7 a4 a2 v V C2 a3 C2 v2 V2 C v Z*",
		FieldNames: []string{"len", "object_type", "ID", "charID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "shield", "lowhead", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "costume", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "xSize", "ySize", "lv", "font", "maxHP", "HP", "isBoss", "opt4", "name"},
	}
	packets["09FF"] = common.PacketConstruction{
		ID:         "09FF",
		Name:       "actor_exists",
		Format:     "v C a4 a4 v3 V v2 V2 v7 a4 a2 v V C2 a3 C3 v2 V2 C v Z*",
		FieldNames: []string{"len", "object_type", "ID", "charID", "walk_speed", "opt1", "opt2", "option", "type", "hair_style", "weapon", "shield", "lowhead", "tophead", "midhead", "hair_color", "clothes_color", "head_dir", "costume", "guildID", "emblemID", "manner", "opt3", "stance", "sex", "coords", "xSize", "ySize", "state", "lv", "font", "maxHP", "HP", "isBoss", "opt4", "name"},
	}
	packets["0A00"] = common.PacketConstruction{
		ID:         "0A00",
		Name:       "hotkeys",
		Format:     "C a*",
		FieldNames: []string{"rotate", "hotkeys"},
	}
	packets["0A05"] = common.PacketConstruction{
		ID:         "0A05",
		Name:       "rodex_add_item",
		Format:     "C a2 v2 C4 a8 a25 v C V",
		FieldNames: []string{"fail", "ID", "amount", "nameID", "type", "identified", "broken", "upgrade", "cards", "options", "weight", "favorite", "type_equip"},
	}
	packets["0A07"] = common.PacketConstruction{
		ID:         "0A07",
		Name:       "rodex_remove_item",
		Format:     "C a2 v2",
		FieldNames: []string{"result", "ID", "amount", "weight"},
	}
	packets["0A09"] = common.PacketConstruction{
		ID:         "0A09",
		Name:       "deal_add_other",
		Format:     "v C V C3 a8 a25",
		FieldNames: []string{"nameID", "type", "amount", "identified", "broken", "upgrade", "cards", "options"},
	}
	packets["0A0A"] = common.PacketConstruction{
		ID:         "0A0A",
		Name:       "storage_item_added",
		Format:     "a2 V v C4 a8 a25",
		FieldNames: []string{"ID", "amount", "nameID", "type", "identified", "broken", "upgrade", "cards", "options"},
	}
	packets["0A0B"] = common.PacketConstruction{
		ID:         "0A0B",
		Name:       "cart_item_added",
		Format:     "a2 V v C4 a8 a25",
		FieldNames: []string{"ID", "amount", "nameID", "type", "identified", "broken", "upgrade", "cards", "options"},
	}
	packets["0A0C"] = common.PacketConstruction{
		ID:         "0A0C",
		Name:       "inventory_item_added",
		Format:     "a2 v2 C3 a8 V C2 a4 v a25",
		FieldNames: []string{"ID", "amount", "nameID", "identified", "broken", "upgrade", "cards", "type_equip", "type", "fail", "expire", "unknown", "options"},
	}
	packets["0A0D"] = common.PacketConstruction{
		ID:         "0A0D",
		Name:       "inventory_items_nonstackable",
		Format:     "v a*",
		FieldNames: []string{"len", "itemInfo"},
	}
	packets["0A0F"] = common.PacketConstruction{
		ID:         "0A0F",
		Name:       "cart_items_nonstackable",
		Format:     "v a*",
		FieldNames: []string{"len", "itemInfo"},
	}
	packets["0A10"] = common.PacketConstruction{
		ID:         "0A10",
		Name:       "storage_items_nonstackable",
		Format:     "v Z24 a*",
		FieldNames: []string{"len", "title", "itemInfo"},
	}
	packets["0A12"] = common.PacketConstruction{
		ID:         "0A12",
		Name:       "rodex_open_write",
		Format:     "Z24 C",
		FieldNames: []string{"name", "result"},
	}
	packets["0A14"] = common.PacketConstruction{
		ID:         "0A14",
		Name:       "rodex_check_player",
		Format:     "V v2",
		FieldNames: []string{"char_id", "class", "base_level"},
	}
	packets["0A15"] = common.PacketConstruction{
		ID:         "0A15",
		Name:       "gold_pc_cafe_point",
		Format:     "C2 V2",
		FieldNames: []string{"isActive", "mode", "point", "playedTime"},
	}
	packets["0A17"] = common.PacketConstruction{
		ID:         "0A17",
		Name:       "dynamicnpc_create_result",
		Format:     "V",
		FieldNames: []string{"result"},
	}
	packets["0A18"] = common.PacketConstruction{
		ID:         "0A18",
		Name:       "map_loaded",
		Format:     "V a3 C2 v C",
		FieldNames: []string{"syncMapSync", "coords", "xSize", "ySize", "font", "sex"},
	}
	packets["0A1A"] = common.PacketConstruction{
		ID:         "0A1A",
		Name:       "roulette_window",
		Format:     "C V C2 v V3",
		FieldNames: []string{"result", "serial", "stage", "price", "additional_item", "gold", "silver", "bronze"},
	}
	packets["0A1C"] = common.PacketConstruction{
		ID:         "0A1C",
		Name:       "roulette_info",
		Format:     "v V a*",
		FieldNames: []string{"len", "serial", "roulette_info"},
	}
	packets["0A20"] = common.PacketConstruction{
		ID:         "0A20",
		Name:       "roulette_window_update",
		Format:     "C v3 V3",
		FieldNames: []string{"result", "stage", "price", "additional_item", "gold", "silver", "bronze"},
	}
	packets["0A22"] = common.PacketConstruction{
		ID:         "0A22",
		Name:       "roulette_recv_item",
		Format:     "C v",
		FieldNames: []string{"type", "item_id"},
	}
	packets["0A23"] = common.PacketConstruction{
		ID:         "0A23",
		Name:       "achievement_list",
		Format:     "v V V v V V",
		FieldNames: []string{"len", "ach_count", "total_points", "rank", "current_rank_points", "next_rank_points"},
	}
	packets["0A24"] = common.PacketConstruction{
		ID:         "0A24",
		Name:       "achievement_update",
		Format:     "V v VVV C V10 V C",
		FieldNames: []string{"total_points", "rank", "current_rank_points", "next_rank_points", "achievementID", "completed", "objective1", "objective2", "objective3", "objective4", "objective5", "objective6", "objective7", "objective8", "objective9", "objective10", "completed_at", "reward"},
	}
	packets["0A26"] = common.PacketConstruction{
		ID:         "0A26",
		Name:       "achievement_reward_ack",
		Format:     "C V",
		FieldNames: []string{"received", "achievementID"},
	}
	packets["0A27"] = common.PacketConstruction{
		ID:         "0A27",
		Name:       "hp_sp_changed",
		Format:     "v V",
		FieldNames: []string{"type", "amount"},
	}
	packets["0A28"] = common.PacketConstruction{
		ID:         "0A28",
		Name:       "open_store_status",
		Format:     "C",
		FieldNames: []string{"flag"},
	}
	packets["0A2D"] = common.PacketConstruction{
		ID:         "0A2D",
		Name:       "show_eq",
		Format:     "v Z24 v7 v C a*",
		FieldNames: []string{"len", "name", "jobID", "hair_style", "tophead", "midhead", "lowhead", "robe", "hair_color", "clothes_color", "sex", "equips_info"},
	}
	packets["0A2F"] = common.PacketConstruction{
		ID:         "0A2F",
		Name:       "change_title",
		Format:     "C V",
		FieldNames: []string{"result", "title_id"},
	}
	packets["0A30"] = common.PacketConstruction{
		ID:         "0A30",
		Name:       "actor_info",
		Format:     "a4 Z24 Z24 Z24 Z24 V",
		FieldNames: []string{"ID", "name", "partyName", "guildName", "guildTitle", "titleID"},
	}
	packets["0A34"] = common.PacketConstruction{
		ID:         "0A34",
		Name:       "senbei_amount",
		Format:     "V",
		FieldNames: []string{"amount"},
	}
	packets["0A36"] = common.PacketConstruction{
		ID:         "0A36",
		Name:       "monster_hp_info_tiny",
		Format:     "a4 C",
		FieldNames: []string{"ID", "hp"},
	}
	packets["0A37"] = common.PacketConstruction{
		ID:         "0A37",
		Name:       "inventory_item_added",
		Format:     "a2 v2 C3 a8 V C2 a4 v a25 C",
		FieldNames: []string{"ID", "amount", "nameID", "identified", "broken", "upgrade", "cards", "type_equip", "type", "fail", "expire", "unknown", "options", "favorite"},
	}
	packets["0A38"] = common.PacketConstruction{
		ID:         "0A38",
		Name:       "open_ui",
		Format:     "C",
		FieldNames: []string{"type"},
	}
	packets["0A3B"] = common.PacketConstruction{
		ID:         "0A3B",
		Name:       "hat_effect",
		Format:     "v a4 C a*",
		FieldNames: []string{"len", "ID", "flag", "effect"},
	}
	packets["0A43"] = common.PacketConstruction{
		ID:         "0A43",
		Name:       "party_join",
		Format:     "a4 V v4 C Z24 Z24 Z16 C2",
		FieldNames: []string{"ID", "role", "jobID", "lv", "x", "y", "type", "name", "user", "map", "item_pickup", "item_share"},
	}
	packets["0A44"] = common.PacketConstruction{
		ID:         "0A44",
		Name:       "party_users_info",
		Format:     "v Z24 a*",
		FieldNames: []string{"len", "party_name", "playerInfo"},
	}
	packets["0A47"] = common.PacketConstruction{
		ID:         "0A47",
		Name:       "stylist_res",
		Format:     "C",
		FieldNames: []string{"res"},
	}
	packets["0A4A"] = common.PacketConstruction{
		ID:         "0A4A",
		Name:       "private_airship_type",
		Format:     "V",
		FieldNames: []string{"type"},
	}
	packets["0A4B"] = common.PacketConstruction{
		ID:         "0A4B",
		Name:       "map_change",
		Format:     "Z16 v2",
		FieldNames: []string{"map", "x", "y"},
	}
	packets["0A4C"] = common.PacketConstruction{
		ID:         "0A4C",
		Name:       "map_changed",
		Format:     "Z16 v2 a4 v",
		FieldNames: []string{"map", "x", "y", "IP", "port"},
	}
	packets["0A51"] = common.PacketConstruction{
		ID:         "0A51",
		Name:       "rodex_check_player",
		Format:     "V v2 Z24",
		FieldNames: []string{"char_id", "class", "base_level", "name"},
	}
	packets["0A53"] = common.PacketConstruction{
		ID:         "0A53",
		Name:       "captcha_upload_request",
		Format:     "Z4 V",
		FieldNames: []string{"captcha_key", "flag"},
	}
	packets["0A55"] = common.PacketConstruction{
		ID:         "0A55",
		Name:       "captcha_upload_request_status",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0A57"] = common.PacketConstruction{
		ID:         "0A57",
		Name:       "macro_reporter_status",
		Format:     "V",
		FieldNames: []string{"status"},
	}
	packets["0A58"] = common.PacketConstruction{
		ID:         "0A58",
		Name:       "macro_detector",
		Format:     "v Z4",
		FieldNames: []string{"image_size", "captcha_key"},
	}
	packets["0A59"] = common.PacketConstruction{
		ID:         "0A59",
		Name:       "macro_detector_image",
		Format:     "v Z4 a*",
		FieldNames: []string{"len", "captcha_key", "captcha_image"},
	}
	packets["0A5B"] = common.PacketConstruction{
		ID:         "0A5B",
		Name:       "macro_detector_show",
		Format:     "c V",
		FieldNames: []string{"remaining_chances", "remaining_time"},
	}
	packets["0A5D"] = common.PacketConstruction{
		ID:         "0A5D",
		Name:       "macro_detector_status",
		Format:     "V",
		FieldNames: []string{"status"},
	}
	packets["0A6A"] = common.PacketConstruction{
		ID:         "0A6A",
		Name:       "captcha_preview",
		Format:     "V v Z4",
		FieldNames: []string{"flag", "image_size", "captcha_key"},
	}
	packets["0A6B"] = common.PacketConstruction{
		ID:         "0A6B",
		Name:       "captcha_preview_image",
		Format:     "v Z4 a*",
		FieldNames: []string{"len", "captcha_key", "captcha_image"},
	}
	packets["0A6D"] = common.PacketConstruction{
		ID:         "0A6D",
		Name:       "macro_reporter_select",
		Format:     "v a*",
		FieldNames: []string{"len", "account_list"},
	}
	packets["0A6F"] = common.PacketConstruction{
		ID:         "0A6F",
		Name:       "message_string",
		Format:     "v2 V Z*",
		FieldNames: []string{"len", "index", "color", "param"},
	}
	packets["0A7B"] = common.PacketConstruction{
		ID:         "0A7B",
		Name:       "EAC_key",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0A7D"] = common.PacketConstruction{
		ID:         "0A7D",
		Name:       "rodex_mail_list",
		Format:     "v C3 a*",
		FieldNames: []string{"len", "type", "amount", "isEnd", "mailList"},
	}
	packets["0A82"] = common.PacketConstruction{
		ID:         "0A82",
		Name:       "guild_expulsion",
		Format:     "Z40 a4",
		FieldNames: []string{"message", "charID"},
	}
	packets["0A83"] = common.PacketConstruction{
		ID:         "0A83",
		Name:       "guild_leave",
		Format:     "a4 Z40",
		FieldNames: []string{"charID", "message"},
	}
	packets["0A84"] = common.PacketConstruction{
		ID:         "0A84",
		Name:       "guild_info",
		Format:     "a4 V9 a4 Z24 Z16 V a4",
		FieldNames: []string{"ID", "lv", "conMember", "maxMember", "average", "exp", "exp_next", "tax", "tendency_left_right", "tendency_down_up", "emblemID", "name", "castles_string", "zeny", "master_char_id"},
	}
	packets["0A89"] = common.PacketConstruction{
		ID:         "0A89",
		Name:       "offline_clone_found",
		Format:     "a4 v4 C v9 Z24",
		FieldNames: []string{"ID", "jobID", "unknown", "coord_x", "coord_y", "sex", "head_dir", "weapon", "shield", "lowhead", "tophead", "midhead", "hair_color", "clothes_color", "robe", "title"},
	}
	packets["0A8A"] = common.PacketConstruction{
		ID:         "0A8A",
		Name:       "offline_clone_lost",
		Format:     "a4",
		FieldNames: []string{"ID"},
	}
	packets["0A8D"] = common.PacketConstruction{
		ID:         "0A8D",
		Name:       "vender_items_list",
		Format:     "v a4 a4 C V a*",
		FieldNames: []string{"len", "venderID", "venderCID", "flag", "expireDate", "itemList"},
	}
	packets["0A91"] = common.PacketConstruction{
		ID:         "0A91",
		Name:       "buying_store_items_list",
		Format:     "v a4 a4 C V V x4 a*",
		FieldNames: []string{"len", "buyerID", "buyingStoreID", "flag", "expireDate", "zeny", "itemList"},
	}
	packets["0A95"] = common.PacketConstruction{
		ID:         "0A95",
		Name:       "misc_config",
		Format:     "C2",
		FieldNames: []string{"show_eq_flag", "call_flag"},
	}
	packets["0A96"] = common.PacketConstruction{
		ID:         "0A96",
		Name:       "deal_add_other",
		Format:     "V C V C3 a16 a25 V v",
		FieldNames: []string{"nameID", "type", "amount", "identified", "broken", "upgrade", "cards", "options", "type_equip", "viewID"},
	}
	packets["0A98"] = common.PacketConstruction{
		ID:         "0A98",
		Name:       "equip_item_switch",
		Format:     "a2 V v",
		FieldNames: []string{"ID", "type", "success"},
	}
	packets["0A9A"] = common.PacketConstruction{
		ID:         "0A9A",
		Name:       "unequip_item_switch",
		Format:     "a2 V C",
		FieldNames: []string{"ID", "type", "success"},
	}
	packets["0A9B"] = common.PacketConstruction{
		ID:         "0A9B",
		Name:       "equip_switch_log",
		Format:     "v a*",
		FieldNames: []string{"len", "log"},
	}
	packets["0A9D"] = common.PacketConstruction{
		ID:         "0A9D",
		Name:       "equip_switch_run_res",
		Format:     "v",
		FieldNames: []string{"success"},
	}
	packets["0AA0"] = common.PacketConstruction{
		ID:         "0AA0",
		Name:       "refineui_opened",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0AA2"] = common.PacketConstruction{
		ID:         "0AA2",
		Name:       "refineui_info",
		Format:     "v v C a*",
		FieldNames: []string{"len", "index", "bless", "materials"},
	}
	packets["0AA5"] = common.PacketConstruction{
		ID:         "0AA5",
		Name:       "guild_members_list",
		Format:     "v a*",
		FieldNames: []string{"len", "member_list"},
	}
	packets["0AA8"] = common.PacketConstruction{
		ID:         "0AA8",
		Name:       "misc_config",
		Format:     "C3",
		FieldNames: []string{"show_eq_flag", "call_flag", "pet_autofeed_flag"},
	}
	packets["0AB2"] = common.PacketConstruction{
		ID:         "0AB2",
		Name:       "party_dead",
		Format:     "a4 C",
		FieldNames: []string{"ID", "isDead"},
	}
	packets["0AB8"] = common.PacketConstruction{
		ID:         "0AB8",
		Name:       "move_interrupt",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0AB9"] = common.PacketConstruction{
		ID:         "0AB9",
		Name:       "item_preview",
		Format:     "a2 v a8 a25",
		FieldNames: []string{"index", "upgrade", "cards", "options"},
	}
	packets["0ABD"] = common.PacketConstruction{
		ID:         "0ABD",
		Name:       "partylv_info",
		Format:     "a4 v2",
		FieldNames: []string{"ID", "job", "lv"},
	}
	packets["0ABE"] = common.PacketConstruction{
		ID:         "0ABE",
		Name:       "warp_portal_list",
		Format:     "v2 Z16 Z16 Z16 Z16",
		FieldNames: []string{"len", "type", "memo1", "memo2", "memo3", "memo4"},
	}
	packets["0AC2"] = common.PacketConstruction{
		ID:         "0AC2",
		Name:       "rodex_mail_list",
		Format:     "v C a*",
		FieldNames: []string{"len", "isEnd", "mailList"},
	}
	packets["0AC4"] = common.PacketConstruction{
		ID:         "0AC4",
		Name:       "account_server_info",
		Format:     "v a4 a4 a4 a4 a26 C x17 a*",
		FieldNames: []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"},
	}
	packets["0AC5"] = common.PacketConstruction{
		ID:         "0AC5",
		Name:       "received_character_ID_and_Map",
		Format:     "a4 Z16 a4 v a128",
		FieldNames: []string{"charID", "mapName", "mapIP", "mapPort", "mapUrl"},
	}
	packets["0AC7"] = common.PacketConstruction{
		ID:         "0AC7",
		Name:       "map_changed",
		Format:     "Z16 v2 a4 v a128",
		FieldNames: []string{"map", "x", "y", "IP", "port", "url"},
	}
	packets["0AC9"] = common.PacketConstruction{
		ID:         "0AC9",
		Name:       "account_server_info",
		Format:     "v a4 a4 a4 a4 a26 C a6 a*",
		FieldNames: []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "unknown", "serverInfo"},
	}
	packets["0ACA"] = common.PacketConstruction{
		ID:         "0ACA",
		Name:       "errors",
		Format:     "C",
		FieldNames: []string{"type"},
	}
	packets["0ACB"] = common.PacketConstruction{
		ID:         "0ACB",
		Name:       "stat_info",
		Format:     "v V2",
		FieldNames: []string{"type", "val", "val2"},
	}
	packets["0ACC"] = common.PacketConstruction{
		ID:         "0ACC",
		Name:       "exp",
		Format:     "a4 V2 v2",
		FieldNames: []string{"ID", "val", "val2", "type", "flag"},
	}
	packets["0ACD"] = common.PacketConstruction{
		ID:         "0ACD",
		Name:       "login_error",
		Format:     "C Z20",
		FieldNames: []string{"type", "date"},
	}
	packets["0ADA"] = common.PacketConstruction{
		ID:         "0ADA",
		Name:       "refine_status",
		Format:     "Z24 V C C",
		FieldNames: []string{"name", "itemID", "refine_level", "status"},
	}
	packets["0ADC"] = common.PacketConstruction{
		ID:         "0ADC",
		Name:       "misc_config",
		Format:     "C4",
		FieldNames: []string{"show_eq_flag", "call_flag", "pet_autofeed_flag", "homunculus_autofeed_flag"},
	}
	packets["0ADD"] = common.PacketConstruction{
		ID:         "0ADD",
		Name:       "item_appeared",
		Format:     "a4 v2 C v2 C2 v C v",
		FieldNames: []string{"ID", "nameID", "type", "identified", "x", "y", "subx", "suby", "amount", "show_effect", "effect_type"},
	}
	packets["0ADE"] = common.PacketConstruction{
		ID:         "0ADE",
		Name:       "overweight_percent",
		Format:     "V",
		FieldNames: []string{"percent"},
	}
	packets["0ADF"] = common.PacketConstruction{
		ID:         "0ADF",
		Name:       "actor_info",
		Format:     "a4 a4 Z24 Z24",
		FieldNames: []string{"ID", "charID", "name", "prefix_name"},
	}
	packets["0AE0"] = common.PacketConstruction{
		ID:         "0AE0",
		Name:       "login_error",
		Format:     "V V Z20",
		FieldNames: []string{"type", "error", "date"},
	}
	packets["0AE2"] = common.PacketConstruction{
		ID:         "0AE2",
		Name:       "open_ui",
		Format:     "C V",
		FieldNames: []string{"type", "data"},
	}
	packets["0AE3"] = common.PacketConstruction{
		ID:         "0AE3",
		Name:       "received_login_token",
		Format:     "v l Z20 Z*",
		FieldNames: []string{"len", "login_type", "flag", "login_token"},
	}
	packets["0AE4"] = common.PacketConstruction{
		ID:         "0AE4",
		Name:       "party_join",
		Format:     "a4 a4 V v4 C Z24 Z24 Z16 C2",
		FieldNames: []string{"ID", "charID", "role", "jobID", "lv", "x", "y", "type", "name", "user", "map", "item_pickup", "item_share"},
	}
	packets["0AE5"] = common.PacketConstruction{
		ID:         "0AE5",
		Name:       "party_users_info",
		Format:     "v Z24 a*",
		FieldNames: []string{"len", "party_name", "playerInfo"},
	}
	packets["0AF0"] = common.PacketConstruction{
		ID:         "0AF0",
		Name:       "action_ui",
		Format:     "C V",
		FieldNames: []string{"type", "data"},
	}
	packets["0AF7"] = common.PacketConstruction{
		ID:         "0AF7",
		Name:       "character_name",
		Format:     "v a4 Z24",
		FieldNames: []string{"flag", "ID", "name"},
	}
	packets["0AFB"] = common.PacketConstruction{
		ID:         "0AFB",
		Name:       "sage_autospell",
		Format:     "v a*",
		FieldNames: []string{"len", "autospell_list"},
	}
	packets["0AFD"] = common.PacketConstruction{
		ID:         "0AFD",
		Name:       "guild_position",
		Format:     "v a4",
		FieldNames: []string{"len", "charID"},
	}
	packets["0AFE"] = common.PacketConstruction{
		ID:         "0AFE",
		Name:       "quest_update_mission_hunt",
		Format:     "v2 a*",
		FieldNames: []string{"len", "mission_amount", "message"},
	}
	packets["0AFF"] = common.PacketConstruction{
		ID:         "0AFF",
		Name:       "quest_all_list",
		Format:     "v V a*",
		FieldNames: []string{"len", "quest_amount", "message"},
	}
	packets["0B03"] = common.PacketConstruction{
		ID:         "0B03",
		Name:       "show_eq",
		Format:     "v Z24 v9 C a*",
		FieldNames: []string{"len", "name", "jobID", "hair_style", "tophead", "midhead", "lowhead", "robe", "hair_color", "clothes_color", "clothes_color2", "sex", "equips_info"},
	}
	packets["0B05"] = common.PacketConstruction{
		ID:         "0B05",
		Name:       "offline_clone_found",
		Format:     "a4 v4 C v9 V Z24 v",
		FieldNames: []string{"ID", "jobID", "unknown", "coord_x", "coord_y", "sex", "head_dir", "weapon", "shield", "lowhead", "tophead", "midhead", "hair_color", "clothes_color", "robe", "unknown2", "name", "unknown3"},
	}
	packets["0B08"] = common.PacketConstruction{
		ID:         "0B08",
		Name:       "item_list_start",
		Format:     "v C Z*",
		FieldNames: []string{"len", "type", "name"},
	}
	packets["0B09"] = common.PacketConstruction{
		ID:         "0B09",
		Name:       "item_list_stackable",
		Format:     "v C a*",
		FieldNames: []string{"len", "type", "itemInfo"},
	}
	packets["0B0A"] = common.PacketConstruction{
		ID:         "0B0A",
		Name:       "item_list_nonstackable",
		Format:     "v C a*",
		FieldNames: []string{"len", "type", "itemInfo"},
	}
	packets["0B0B"] = common.PacketConstruction{
		ID:         "0B0B",
		Name:       "item_list_end",
		Format:     "C2",
		FieldNames: []string{"type", "flag"},
	}
	packets["0B0C"] = common.PacketConstruction{
		ID:         "0B0C",
		Name:       "quest_add",
		Format:     "V C V2 v a*",
		FieldNames: []string{"questID", "active", "time_start", "time_expire", "mission_amount", "message"},
	}
	packets["0B13"] = common.PacketConstruction{
		ID:         "0B13",
		Name:       "item_preview",
		Format:     "a2 C v a16 a25",
		FieldNames: []string{"index", "broken", "upgrade", "cards", "options"},
	}
	packets["0B18"] = common.PacketConstruction{
		ID:         "0B18",
		Name:       "inventory_expansion_result",
		Format:     "v",
		FieldNames: []string{"result"},
	}
	packets["0B1A"] = common.PacketConstruction{
		ID:         "0B1A",
		Name:       "skill_cast",
		Format:     "a4 a4 v5 V C V",
		FieldNames: []string{"sourceID", "targetID", "x", "y", "skillID", "unknown", "type", "wait", "dispose", "unknow"},
	}
	packets["0B1B"] = common.PacketConstruction{
		ID:         "0B1B",
		Name:       "load_confirm",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0B1D"] = common.PacketConstruction{
		ID:         "0B1D",
		Name:       "ping",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0B20"] = common.PacketConstruction{
		ID:         "0B20",
		Name:       "hotkeys",
		Format:     "C v a*",
		FieldNames: []string{"rotate", "tab", "hotkeys"},
	}
	packets["0B2F"] = common.PacketConstruction{
		ID:         "0B2F",
		Name:       "homunculus_property",
		Format:     "Z24 C v11 V2 v2 V2 v2",
		FieldNames: []string{"name", "state", "level", "hunger", "intimacy", "atk", "matk", "hit", "critical", "def", "mdef", "flee", "aspd", "hp", "hp_max", "sp", "sp_max", "exp", "exp_max", "points_skill", "attack_range"},
	}
	packets["0B31"] = common.PacketConstruction{
		ID:         "0B31",
		Name:       "skill_add",
		Format:     "v V v3 C v",
		FieldNames: []string{"skillID", "target", "lv", "sp", "range", "upgradable", "lv2"},
	}
	packets["0B32"] = common.PacketConstruction{
		ID:         "0B32",
		Name:       "skills_list",
		Format:     "",
		FieldNames: []string{},
	}
	packets["0B33"] = common.PacketConstruction{
		ID:         "0B33",
		Name:       "skill_update",
		Format:     "v V v3 C v",
		FieldNames: []string{"skillID", "type", "lv", "sp", "range", "up", "lv2"},
	}
	packets["0B39"] = common.PacketConstruction{
		ID:         "0B39",
		Name:       "item_list_nonstackable",
		Format:     "v C a*",
		FieldNames: []string{"len", "type", "itemInfo"},
	}
	packets["0B3D"] = common.PacketConstruction{
		ID:         "0B3D",
		Name:       "vender_items_list",
		Format:     "v a4 a4 a*",
		FieldNames: []string{"len", "venderID", "venderCID", "itemList"},
	}
	packets["0B41"] = common.PacketConstruction{
		ID:         "0B41",
		Name:       "inventory_item_added",
		Format:     "a2 v V C2 a16 V C2 a4 v a25 C v C2",
		FieldNames: []string{"ID", "amount", "nameID", "identified", "broken", "cards", "type_equip", "type", "fail", "expire", "unknown", "options", "favorite", "viewID", "upgrade", "grade"},
	}
	packets["0B44"] = common.PacketConstruction{
		ID:         "0B44",
		Name:       "storage_item_added",
		Format:     "a2 V V C3 a16 a25 C2",
		FieldNames: []string{"ID", "amount", "nameID", "type", "identified", "broken", "cards", "options", "upgrade", "grade"},
	}
	packets["0B45"] = common.PacketConstruction{
		ID:         "0B45",
		Name:       "cart_item_added",
		Format:     "a2 V V C3 a16 a25 C2",
		FieldNames: []string{"ID", "amount", "nameID", "type", "identified", "broken", "upgrade", "cards", "options", "upgrade", "grade"},
	}
	packets["0B47"] = common.PacketConstruction{
		ID:         "0B47",
		Name:       "char_emblem_update",
		Format:     "a4 a4",
		FieldNames: []string{"guildID", "emblemID", "accountID"},
	}
	packets["0B5F"] = common.PacketConstruction{
		ID:         "0B5F",
		Name:       "rodex_mail_list",
		Format:     "v C a*",
		FieldNames: []string{"len", "isEnd", "mailList"},
	}
	packets["0B60"] = common.PacketConstruction{
		ID:         "0B60",
		Name:       "account_server_info",
		Format:     "v a4 a4 a4 a4 a26 C x17 a*",
		FieldNames: []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"},
	}
	packets["0B6F"] = common.PacketConstruction{
		ID:         "0B6F",
		Name:       "character_creation_successful",
		Format:     "a*",
		FieldNames: []string{"charInfo"},
	}
	packets["0B72"] = common.PacketConstruction{
		ID:         "0B72",
		Name:       "received_characters",
		Format:     "v a*",
		FieldNames: []string{"len", "charInfo"},
	}
	packets["0B73"] = common.PacketConstruction{
		ID:         "0B73",
		Name:       "revolving_entity",
		Format:     "a4 v",
		FieldNames: []string{"sourceID", "entity"},
	}
	packets["0B76"] = common.PacketConstruction{
		ID:         "0B76",
		Name:       "homunculus_property",
		Format:     "Z24 C v11 V6 v2",
		FieldNames: []string{"name", "state", "level", "hunger", "intimacy", "atk", "matk", "hit", "critical", "def", "mdef", "flee", "aspd", "hp", "hp_max", "sp", "sp_max", "exp", "exp_max", "points_skill", "attack_range"},
	}
	packets["0B77"] = common.PacketConstruction{
		ID:         "0B77",
		Name:       "npc_store_info",
		Format:     "v a*",
		FieldNames: []string{"len", "itemList"},
	}
	packets["0B7B"] = common.PacketConstruction{
		ID:         "0B7B",
		Name:       "guild_info",
		Format:     "a4 V9 a4 Z24 Z16 V a4 Z24",
		FieldNames: []string{"ID", "lv", "conMember", "maxMember", "average", "exp", "exp_next", "tax", "tendency_left_right", "tendency_down_up", "emblemID", "name", "castles_string", "zeny", "master_char_id", "master"},
	}
	packets["0B7C"] = common.PacketConstruction{
		ID:         "0B7C",
		Name:       "guild_expulsion_list",
		Format:     "v a*",
		FieldNames: []string{"len", "expulsion_list"},
	}
	packets["0B7D"] = common.PacketConstruction{
		ID:         "0B7D",
		Name:       "guild_members_list",
		Format:     "v a*",
		FieldNames: []string{"len", "member_list"},
	}
	packets["0B7E"] = common.PacketConstruction{
		ID:         "0B7E",
		Name:       "guild_member_add",
		Format:     "a4 a4 v5 V4 Z24",
		FieldNames: []string{"ID", "charID", "hair_style", "hair_color", "sex", "jobID", "lv", "contribution", "online", "position", "lastLoginTime", "name"},
	}
	packets["0B8D"] = common.PacketConstruction{
		ID:         "0B8D",
		Name:       "repute_info",
		Format:     "v C a*",
		FieldNames: []string{"len", "sucess", "reputeInfo"},
	}
	// Commented out as it needs research
	// packets["C350"] = common.PacketConstruction{
	// 	ID:         "C350",
	// 	Name:       "senbei_vender_items_list",
	// 	Format:     "",
	// 	FieldNames: []string{},
	// }
	packets["0BA4"] = common.PacketConstruction{
		ID:         "0BA4",
		Name:       "homunculus_property",
		Format:     "Z24 C v11 V4 V4 v2",
		FieldNames: []string{"name", "state", "level", "hunger", "intimacy", "atk", "matk", "hit", "critical", "def", "mdef", "flee", "aspd", "hp", "hp_max", "sp", "sp_max", "exp", "exp2", "exp_max", "exp_max2", "points_skill", "attack_range"},
	}

	return packets
}
