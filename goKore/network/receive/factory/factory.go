// Package factory provides functionality for creating and configuring Receive implementations.
package factory

import (
	"fmt"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/base"
	"github.com/lenaxia/goKore/network/receive/servers"
	"github.com/lenaxia/goKore/network/receive/types"
)

// PacketConstructionProvider provides packet constructions for a specific server type
type PacketConstructionProvider func() map[string]common.PacketConstruction

// ReceiveFactory creates and configures Receive implementations
type ReceiveFactory struct {
	// Map of server type to packet construction provider
	packetConstructionProviders map[string]PacketConstructionProvider
}

// NewReceiveFactory creates a new receive factory
func NewReceiveFactory() *ReceiveFactory {
	return &ReceiveFactory{
		packetConstructionProviders: make(map[string]PacketConstructionProvider),
	}
}

// RegisterServerType registers packet constructions for a server type
func (rf *ReceiveFactory) RegisterServerType(serverType string, provider PacketConstructionProvider) {
	rf.packetConstructionProviders[serverType] = provider
}

// CreateReceive creates and configures a Receive implementation for a server type
func (rf *ReceiveFactory) CreateReceive(serverType string, hookManager *hooks.HookManager) (types.Receive, error) {
	provider, exists := rf.packetConstructionProviders[serverType]
	if !exists {
		return nil, fmt.Errorf("no packet constructions registered for server type: %s", serverType)
	}

	// Create a new BaseReceive instance
	receive := base.NewBaseReceive(hookManager)

	// Convert PacketConstruction to PacketDef for backward compatibility
	packetDefs := make(map[string]common.PacketConstruction)
	for id, construction := range provider() {
		packetDefs[id] = common.PacketConstruction{
			Name:       construction.Name,
			Format:     construction.Format,
			FieldNames: construction.FieldNames,
		}
	}

	err := receive.Configure(serverType, packetDefs)
	if err != nil {
		return nil, err
	}

	return receive, nil
}

// RegisterDefaultServerTypes registers the default server types
func (rf *ReceiveFactory) RegisterDefaultServerTypes() {
	// Register ServerType0
	rf.RegisterServerType("ServerType0", servers.ServerType0PacketConstructions)

	// Register other server types
	// We need to create a wrapper for SakrayPacketDefs since it's still using the old format
	rf.RegisterServerType("ServerTypeSakray", func() map[string]common.PacketConstruction {
		// Convert from PacketDef to PacketConstruction
		sakrayDefs := servers.SakrayPacketDefs()
		constructions := make(map[string]common.PacketConstruction)

		for id, def := range sakrayDefs {
			constructions[id] = common.PacketConstruction{
				ID:         id,
				Name:       def.Name,
				Format:     def.Format,
				FieldNames: def.FieldNames,
			}
		}

		return constructions
	})
}

// Packet construction functions have been moved to the servers package
