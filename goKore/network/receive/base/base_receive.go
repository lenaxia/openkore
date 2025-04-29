// Package base provides the base implementation of the Receive interface.
package base

import (
	"fmt"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/types"
)

// BaseReceive implements the Receive interface
type BaseReceive struct {
	// Server type
	serverType string

	// Packet definitions
	packetDefs map[string]common.PacketDef

	// Core parser for packet parsing
	coreParser *core.CoreParser

	// Hook manager
	hookManager *hooks.HookManager

	// Packet handlers
	handlers map[string]types.ReceiveHandler

	// Packet ID lookup table
	packetLUT map[string]string

	// Debug mode
	debugMode bool
}

// NewBaseReceive creates a new BaseReceive instance
func NewBaseReceive(hookManager *hooks.HookManager) *BaseReceive {
	return &BaseReceive{
		packetDefs:  make(map[string]common.PacketDef),
		handlers:    make(map[string]types.ReceiveHandler),
		packetLUT:   make(map[string]string),
		coreParser:  core.NewCoreParser("", hookManager),
		hookManager: hookManager,
	}
}

// RegisterHandler registers a handler for a specific packet
func (br *BaseReceive) RegisterHandler(packetName string, handler types.ReceiveHandler) {
	// Store the handler
	br.handlers[packetName] = handler

	// Find the packet ID for the given name
	packetID, exists := br.packetLUT[packetName]
	if exists {
		// Register with the core parser
		def := br.packetDefs[packetID]
		br.coreParser.RegisterHandler(packetID, packetName, def.Format, def.FieldNames, core.PacketHandler(handler))
	}
}

// Process processes a packet, calling the appropriate handler and hooks
func (br *BaseReceive) Process(packet []byte) error {
	if br.debugMode {
		fmt.Printf("Processing packet: %X [%d bytes]\n", packet[:2], len(packet))
	}

	return br.coreParser.Process(packet)
}

// Configure configures the receive component with server-specific packet definitions
func (br *BaseReceive) Configure(serverType string, packetDefs map[string]common.PacketDef) error {
	br.serverType = serverType
	br.packetDefs = packetDefs
	br.coreParser.SetDefaultState(0) // Default state

	// Build the lookup table and register packet formats with the parser
	for id, def := range packetDefs {
		br.packetLUT[def.Name] = id
		br.coreParser.RegisterHandler(id, def.Name, def.Format, def.FieldNames, nil)
	}

	// Register any handlers that were added before configuration
	for packetName, handler := range br.handlers {
		packetID, exists := br.packetLUT[packetName]
		if exists {
			def := br.packetDefs[packetID]
			br.coreParser.RegisterHandler(packetID, packetName, def.Format, def.FieldNames, core.PacketHandler(handler))
		}
	}

	return nil
}

// GetPacketID returns the packet ID for a given packet name
func (br *BaseReceive) GetPacketID(name string) (string, bool) {
	id, exists := br.packetLUT[name]
	return id, exists
}

// RegisterHook registers a hook for a specific event
func (br *BaseReceive) RegisterHook(hookName string, callback hooks.HookCallback) {
	if br.hookManager != nil {
		br.hookManager.AddHook(hookName, callback, nil)
	}
}

// GetServerType returns the server type
func (br *BaseReceive) GetServerType() string {
	return br.serverType
}

// SetDebugMode sets the debug mode
func (br *BaseReceive) SetDebugMode(debug bool) {
	br.debugMode = debug
}

// ParsePacket parses a packet and returns the parsed arguments
func (br *BaseReceive) ParsePacket(packet []byte) (map[string]interface{}, error) {
	return br.coreParser.Parse(packet)
}
