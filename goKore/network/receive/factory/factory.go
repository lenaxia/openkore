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

// PacketDefProvider provides packet definitions for a specific server type
type PacketDefProvider func() map[string]common.PacketDef

// ReceiveFactory creates and configures Receive implementations
type ReceiveFactory struct {
	// Map of server type to packet definition provider
	packetDefProviders map[string]PacketDefProvider
}

// NewReceiveFactory creates a new receive factory
func NewReceiveFactory() *ReceiveFactory {
	return &ReceiveFactory{
		packetDefProviders: make(map[string]PacketDefProvider),
	}
}

// RegisterServerType registers packet definitions for a server type
func (rf *ReceiveFactory) RegisterServerType(serverType string, provider PacketDefProvider) {
	rf.packetDefProviders[serverType] = provider
}

// CreateReceive creates and configures a Receive implementation for a server type
func (rf *ReceiveFactory) CreateReceive(serverType string, hookManager *hooks.HookManager) (types.Receive, error) {
	provider, exists := rf.packetDefProviders[serverType]
	if !exists {
		return nil, fmt.Errorf("no packet definitions registered for server type: %s", serverType)
	}

	// Create a new BaseReceive instance
	receive := base.NewBaseReceive(hookManager)
	err := receive.Configure(serverType, provider())
	if err != nil {
		return nil, err
	}

	return receive, nil
}

// RegisterDefaultServerTypes registers the default server types
func (rf *ReceiveFactory) RegisterDefaultServerTypes() {
	// Register ServerType0
	rf.RegisterServerType("ServerType0", servers.ServerType0PacketDefs)

	// Register other server types
	rf.RegisterServerType("ServerTypeSakray", servers.SakrayPacketDefs)
}

// Packet definition functions have been moved to the servers package
