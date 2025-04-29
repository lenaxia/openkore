// Package servers provides server-specific packet constructions for different server types.
package servers

import (
	"github.com/lenaxia/goKore/network/common"
)

// ServerType1PacketConstructions provides packet constructions for ServerType1
func ServerType1PacketConstructions() map[string]common.PacketConstruction {
	// Start with base constructions
	constructions := ServerType0PacketConstructions()

	// Override specific packet constructions
	constructions["0064"] = common.PacketConstruction{
		ID:         "0064",
		Name:       "login_request",
		Format:     "v a24 a24 C C",
		FieldNames: []string{"version", "username", "password", "clienttype", "clientinfo"},
	}

	// Add more ServerType1-specific packet constructions here

	return constructions
}
