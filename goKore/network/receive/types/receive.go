// Package types provides type definitions for the receive component.
package types

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
)

// ReceiveHandler is a function that processes a parsed packet
type ReceiveHandler = common.PacketHandler

// Receive defines the interface for packet receiving and handling
type Receive interface {
	// RegisterHandler registers a handler for a specific packet
	RegisterHandler(packetName string, handler ReceiveHandler)

	// Process processes a packet, calling the appropriate handler and hooks
	Process(packet []byte) error

	// Configure configures the receive component with server-specific packet definitions
	Configure(serverType string, packetDefs map[string]common.PacketDef) error

	// GetPacketID returns the packet ID for a given packet name
	GetPacketID(name string) (string, bool)

	// RegisterHook registers a hook for a specific event
	RegisterHook(hookName string, callback hooks.HookCallback)

	// GetServerType returns the server type
	GetServerType() string

	// SetDebugMode sets the debug mode
	SetDebugMode(debug bool)

	// ParsePacket parses a packet and returns the parsed arguments
	ParsePacket(packet []byte) (map[string]interface{}, error)
}
