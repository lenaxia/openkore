// Package core provides the core functionality for sending packets to the server.
package core

import (
	"github.com/lenaxia/goKore/network/hooks"
)

// Send defines the interface for sending packets to the server.
type Send interface {
	// SendToServer sends a raw packet to the server.
	SendToServer(msg []byte) error

	// EncryptMessageID encrypts the message ID of a packet.
	EncryptMessageID(msg *[]byte) error

	// CryptKeys sets the encryption keys for message ID encryption.
	CryptKeys(key1, key2, key3 uint32)

	// PinEncode encodes a PIN code using the given seed.
	PinEncode(seed, pin int) string

	// InjectMessage sends a text message to the connected client's party chat.
	InjectMessage(message string) error

	// InjectAdminMessage sends a text message to the connected client's system chat.
	InjectAdminMessage(message string) error

	// SendRaw sends a raw packet to the server from a space-delimited list of hex byte values.
	SendRaw(raw string) error

	// Reconstruct constructs a packet from a packet ID and arguments.
	Reconstruct(packetID string, args map[string]interface{}) ([]byte, error)

	// GetPacketID returns the packet ID for a given packet name.
	GetPacketID(name string) (string, bool)

	// RegisterPacketHandler registers a handler for a packet.
	RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error)

	// RegisterHook registers a hook for a specific event.
	RegisterHook(hookName string, callback hooks.HookCallback)

	// SetConnection sets the connection to use for sending packets.
	SetConnection(conn interface{})

	// GetConnection returns the current connection.
	GetConnection() interface{}

	// GetTime returns the current time in milliseconds.
	GetTime() uint32
}

// SendConfig contains configuration for the Send implementation.
type SendConfig struct {
	// ServerType is the type of server to connect to.
	ServerType string

	// PacketVersion is the version of the packet structure to use.
	PacketVersion int

	// Encryption settings
	UseEncryption bool
	CryptKey1     uint32
	CryptKey2     uint32
	CryptKey3     uint32
}

// PacketHandler is a function that handles a packet.
type PacketHandler func(args map[string]interface{}) error

// PacketDefinition defines the structure of a packet.
type PacketDefinition struct {
	// ID is the packet ID.
	ID string

	// Name is the name of the packet.
	Name string

	// Format is the format of the packet.
	Format string

	// Keys are the names of the fields in the packet.
	Keys []string

	// Handler is the function that handles the packet.
	Handler PacketHandler
}
