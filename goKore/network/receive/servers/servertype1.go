// Package servers provides server-specific packet definitions for different server types.
package servers

import (
	"github.com/lenaxia/goKore/network/common"
)

// ServerType1PacketDefs provides packet definitions for ServerType1
func ServerType1PacketDefs() map[string]common.PacketDef {
	// Start with base definitions
	defs := ServerType0PacketDefs()

	// Override specific packet definitions
	defs["0069"] = common.PacketDef{
		ID:         "0069",
		Name:       "account_server_info",
		Format:     "v a4 a4 a4 a4 a26 C a* v",
		FieldNames: []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo", "additionalField"},
	}

	// Add more ServerType1-specific packet definitions here

	return defs
}
