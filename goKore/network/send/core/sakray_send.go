// Package core provides the core functionality for sending packets to the server.
package core

import "github.com/lenaxia/goKore/network/hooks"

// SakraySend implements the Send interface for Sakray servers.
type SakraySend struct {
	// Base send implementation
	baseSend Send
}

// NewSakraySend creates a new SakraySend instance.
func NewSakraySend(baseSend Send) *SakraySend {
	return &SakraySend{
		baseSend: baseSend,
	}
}

// SendToServer sends a raw packet to the server.
func (ss *SakraySend) SendToServer(msg []byte) error {
	// Delegate to base send
	return ss.baseSend.SendToServer(msg)
}

// EncryptMessageID encrypts the message ID of a packet.
func (ss *SakraySend) EncryptMessageID(msg *[]byte) error {
	// Delegate to base send
	return ss.baseSend.EncryptMessageID(msg)
}

// CryptKeys sets the encryption keys for message ID encryption.
func (ss *SakraySend) CryptKeys(key1, key2, key3 uint32) {
	// Delegate to base send
	ss.baseSend.CryptKeys(key1, key2, key3)
}

// PinEncode encodes a PIN code using the given seed.
func (ss *SakraySend) PinEncode(seed, pin int) string {
	// Delegate to base send
	return ss.baseSend.PinEncode(seed, pin)
}

// InjectMessage sends a text message to the connected client's party chat.
func (ss *SakraySend) InjectMessage(message string) error {
	// Delegate to base send
	return ss.baseSend.InjectMessage(message)
}

// InjectAdminMessage sends a text message to the connected client's system chat.
func (ss *SakraySend) InjectAdminMessage(message string) error {
	// Delegate to base send
	return ss.baseSend.InjectAdminMessage(message)
}

// SendRaw sends a raw packet to the server from a space-delimited list of hex byte values.
func (ss *SakraySend) SendRaw(raw string) error {
	// Delegate to base send
	return ss.baseSend.SendRaw(raw)
}

// Reconstruct constructs a packet from a packet ID and arguments.
func (ss *SakraySend) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Check if this is a Sakray-specific packet
	if sakrayPacket, ok := ss.getSakrayPacket(packetID, args); ok {
		return sakrayPacket, nil
	}

	// Delegate to base send
	return ss.baseSend.Reconstruct(packetID, args)
}

// GetPacketID returns the packet ID for a given packet name.
func (ss *SakraySend) GetPacketID(name string) (string, bool) {
	// Check if this is a Sakray-specific packet name
	if packetID, ok := ss.getSakrayPacketID(name); ok {
		return packetID, true
	}

	// Delegate to base send
	return ss.baseSend.GetPacketID(name)
}

// RegisterPacketHandler registers a handler for a packet.
func (ss *SakraySend) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	// Delegate to base send
	ss.baseSend.RegisterPacketHandler(packetID, name, format, keys, handler)
}

// RegisterHook registers a hook for a specific event.
func (ss *SakraySend) RegisterHook(hookName string, callback hooks.HookCallback) {
	// Delegate to base send
	ss.baseSend.RegisterHook(hookName, callback)
}

// SetConnection sets the connection to use for sending packets.
func (ss *SakraySend) SetConnection(conn interface{}) {
	// Delegate to base send
	ss.baseSend.SetConnection(conn)
}

// GetConnection returns the current connection.
func (ss *SakraySend) GetConnection() interface{} {
	// Delegate to base send
	return ss.baseSend.GetConnection()
}

// GetTime returns the current time in milliseconds.
func (ss *SakraySend) GetTime() uint32 {
	// Delegate to base send
	return ss.baseSend.GetTime()
}

// getSakrayPacket returns a Sakray-specific packet.
func (ss *SakraySend) getSakrayPacket(packetID string, args map[string]interface{}) ([]byte, bool) {
	// TODO: Implement Sakray-specific packet construction
	return nil, false
}

// getSakrayPacketID returns a Sakray-specific packet ID.
func (ss *SakraySend) getSakrayPacketID(name string) (string, bool) {
	// TODO: Implement Sakray-specific packet ID lookup
	return "", false
}

// RegisterSakrayPackets registers Sakray-specific packets.
func (ss *SakraySend) RegisterSakrayPackets() {
	// TODO: Register Sakray-specific packets
	// Example:
	// ss.baseSend.RegisterPacketHandler("0123", "sakray_packet", "v", []string{"value"}, func(args map[string]interface{}) error {
	//     // Handle packet
	//     return nil
	// })
}
