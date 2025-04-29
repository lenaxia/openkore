// Package servers provides server-specific packet constructions for different server types.
package servers

import (
	"github.com/lenaxia/goKore/network/common"
)

// SakrayPacketConstructions provides packet constructions for Sakray server type
func SakrayPacketConstructions() map[string]common.PacketConstruction {
	// Start with base constructions
	constructions := ServerType0PacketConstructions()

	// Override or add Sakray-specific packet constructions
	// For example:
	/*
		constructions["0064"] = common.PacketConstruction{
			ID:         "0064",
			Name:       "login_request",
			Format:     "v a24 a24 C C C",
			FieldNames: []string{"version", "username", "password", "clienttype", "clientinfo", "sakray_auth"},
		}
	*/

	// Add more Sakray-specific packet constructions here

	return constructions
}
