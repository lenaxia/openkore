// Package servers provides server-specific packet definitions for different server types.
package servers

import (
	"github.com/lenaxia/goKore/network/common"
)

// SakrayPacketDefs provides packet definitions for Sakray server type
func SakrayPacketDefs() map[string]common.PacketConstruction {
	// Start with base definitions
	defs := ServerType0PacketConstructions()

	// Override or add Sakray-specific packet definitions
	// For example:
	/*
		defs["0069"] = common.PacketConstruction{
			ID:         "0069",
			Name:       "account_server_info",
			Format:     "v a4 a4 a4 a4 a26 C a* v C",
			FieldNames: []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo", "additionalField", "sakrayField"},
		}
	*/

	// Add more Sakray-specific packet definitions here

	return defs
}
