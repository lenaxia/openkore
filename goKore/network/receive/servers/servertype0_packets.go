// Package servers provides server-specific packet definitions for different server types.
package servers

import (
	"github.com/lenaxia/goKore/network/common"
)

// ServerType0PacketDefs provides packet definitions for ServerType0
func ServerType0PacketDefs() map[string]common.PacketDef {
	defs := map[string]common.PacketDef{
		"0069": {
			Name:       "account_server_info",
			Format:     "v a4 a4 a4 a4 a26 C a*",
			FieldNames: []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"},
		},
		"006A": {
			Name:       "login_error",
			Format:     "C Z20",
			FieldNames: []string{"type", "date"},
		},
		"006B": {
			Name:       "received_characters_info",
			Format:     "v C3 x20 a*",
			FieldNames: []string{"len", "total_slot", "premium_start_slot", "premium_end_slot", "charInfo"},
		},
		"006C": {
			Name:       "login_error_game_login_server",
			Format:     "",
			FieldNames: []string{},
		},
		"006D": {
			Name:       "character_creation_successful",
			Format:     "a*",
			FieldNames: []string{"charInfo"},
		},
		"006F": {
			Name:       "character_deletion_successful",
			Format:     "",
			FieldNames: []string{},
		},
		"0070": {
			Name:       "character_deletion_failed",
			Format:     "C",
			FieldNames: []string{"error"},
		},
	}

	return defs
}
