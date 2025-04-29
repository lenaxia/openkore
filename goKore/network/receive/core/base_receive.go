// Package core provides core functionality for receiving and processing network packets.
package core

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
)

// BaseReceive implements the Receive interface using the CoreParser
type BaseReceive struct {
	serverType  string
	packetDefs  map[string]common.PacketDef
	coreParser  *CoreParser
	hookManager *hooks.HookManager
}

// NewBaseReceive creates a new BaseReceive instance
func NewBaseReceive(hookManager *hooks.HookManager) *BaseReceive {
	return &BaseReceive{
		packetDefs:  make(map[string]common.PacketDef),
		coreParser:  NewCoreParser("", hookManager),
		hookManager: hookManager,
	}
}

// RegisterHandler registers a handler for a specific packet
func (br *BaseReceive) RegisterHandler(packetName string, handler PacketHandler) {
	// Find the packet ID for the given name
	for id, def := range br.packetDefs {
		if def.Name == packetName {
			br.coreParser.RegisterHandler(id, packetName, def.Format, def.FieldNames, handler)
			return
		}
	}
}

// Process processes a packet, calling the appropriate handler and hooks
func (br *BaseReceive) Process(packet []byte) error {
	return br.coreParser.Process(packet)
}

// Configure configures the receive component with server-specific packet definitions
func (br *BaseReceive) Configure(serverType string, packetDefs map[string]common.PacketDef) error {
	br.serverType = serverType
	br.packetDefs = packetDefs
	br.coreParser.serverType = serverType

	// Register packet formats with the parser
	for id, def := range packetDefs {
		br.coreParser.RegisterHandler(id, def.Name, def.Format, def.FieldNames, nil)
	}

	return nil
}

// GetServerType returns the server type
func (br *BaseReceive) GetServerType() string {
	return br.serverType
}
