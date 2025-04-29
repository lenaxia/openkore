// Package servers provides server-specific packet constructions for different server types.
package servers

import (
	"github.com/lenaxia/goKore/network/common"
)

// ServerType0PacketConstructions provides packet constructions for ServerType0
func ServerType0PacketConstructions() map[string]common.PacketConstruction {
	return map[string]common.PacketConstruction{
		"0064": {
			ID:         "0064",
			Name:       "login_request",
			Format:     "v a24 a24 C",
			FieldNames: []string{"version", "username", "password", "clienttype"},
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
		"0234": {
			ID:         "0234",
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
		"00CF": {
			ID:         "00CF",
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
		"0090": {
			ID:         "0090",
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
		"0288": {
			ID:         "0288",
			Name:       "cash_shop_buy",
			Format:     "V v a*",
			FieldNames: []string{"kafra_points", "count", "buy_info"},
		},
		"0848": {
			ID:         "0848",
			Name:       "cash_dealer_buy",
			Format:     "V v V",
			FieldNames: []string{"itemid", "amount", "kafra_points"},
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
		"0447": {
			ID:         "0447",
			Name:       "blocking_play_cancel",
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
}
