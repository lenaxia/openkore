// Package servers provides server-specific packet definitions for different server types.
package servers

import (
	"github.com/lenaxia/goKore/network/common"
)

// ServerType0PacketDefs provides packet definitions for ServerType0
func ServerType0PacketDefs() map[string]common.PacketDef {
	return map[string]common.PacketDef{
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
		"006B": {
			ID:         "006B",
			Name:       "received_characters_info",
			Format:     "v C3 x20 a*",
			FieldNames: []string{"len", "total_slot", "premium_start_slot", "premium_end_slot", "charInfo"},
		},
		"006C": {
			ID:         "006C",
			Name:       "login_error_game_login_server",
			Format:     "",
			FieldNames: []string{},
		},
		"006D": {
			ID:         "006D",
			Name:       "character_creation_successful",
			Format:     "a*",
			FieldNames: []string{"charInfo"},
		},
		"0071": {
			ID:         "0071",
			Name:       "received_character_ID_and_Map",
			Format:     "a4 Z16 a4 v",
			FieldNames: []string{"charID", "mapName", "mapIP", "mapPort"},
		},
		"0073": {
			ID:         "0073",
			Name:       "map_loaded",
			Format:     "V a3 C2",
			FieldNames: []string{"syncMapSync", "coords", "xSize", "ySize"},
		},
		"0074": {
			ID:         "0074",
			Name:       "map_load_error",
			Format:     "C",
			FieldNames: []string{"error"},
		},
		"007E": {
			ID:         "007E",
			Name:       "map_change",
			Format:     "Z16 v2",
			FieldNames: []string{"map", "x", "y"},
		},
		"007F": {
			ID:         "007F",
			Name:       "connection_refused",
			Format:     "C",
			FieldNames: []string{"error"},
		},
		"0088": {
			ID:         "0088",
			Name:       "actor_movement_interrupted",
			Format:     "a4 v2",
			FieldNames: []string{"ID", "posX", "posY"},
		},
		"0091": {
			ID:         "0091",
			Name:       "map_changed",
			Format:     "Z16 v2 V2",
			FieldNames: []string{"map", "x", "y", "ip", "port"},
		},
		"00AE": {
			ID:         "00AE",
			Name:       "received_sync",
			Format:     "V",
			FieldNames: []string{"time"},
		},
		"07E5": {
			ID:         "07E5",
			Name:       "sync_request_ex",
			Format:     "",
			FieldNames: []string{"switch"},
		},
		"02F1": {
			ID:         "02F1",
			Name:       "initialize_message_id_encryption",
			Format:     "V V",
			FieldNames: []string{"param1", "param2"},
		},
		"0B1D": {
			ID:         "0B1D",
			Name:       "ping",
			Format:     "",
			FieldNames: []string{},
		},
		"0AB8": {
			ID:         "0AB8",
			Name:       "move_interrupt",
			Format:     "",
			FieldNames: []string{},
		},
		// More packet definitions would be added here
		// This file will grow as more packet definitions are defined
	}
}
