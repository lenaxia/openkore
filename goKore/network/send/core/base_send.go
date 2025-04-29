// Package core provides the core functionality for sending packets to the server.
package core

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
)

// Error variables
var (
	// ErrNoConnection is returned when trying to send a packet without a connection.
	ErrNoConnection = errors.New("no connection available")

	// ErrInvalidPacketID is returned when trying to send a packet with an invalid ID.
	ErrInvalidPacketID = errors.New("invalid packet ID")

	// ErrPacketNotRegistered is returned when trying to send a packet that is not registered.
	ErrPacketNotRegistered = errors.New("packet not registered")
)

// SendHandler is a function that constructs a packet
type SendHandler func(args map[string]interface{}) ([]byte, error)

// BaseSend implements the Send interface with the aligned architecture.
type BaseSend struct {
	// Connection to the server.
	conn interface{}

	// Hook manager.
	hookManager *hooks.HookManager

	// Server type.
	serverType string

	// Packet handlers.
	handlers map[string]SendHandler

	// Packet constructions.
	packetConstructions map[string]common.PacketConstruction

	// Packet name to ID lookup table.
	packetLUT map[string]string

	// Packet builder for packet construction
	packetBuilder *PacketBuilder

	// Encryption settings.
	encryption struct {
		cryptKey   uint32
		cryptKey1  uint32
		cryptKey2  uint32
		cryptKey3  uint32
		useEncrypt bool
	}

	// Debug mode.
	debugMode bool
}

// NewBaseSend creates a new BaseSend instance.
func NewBaseSend(hookManager *hooks.HookManager) *BaseSend {
	bs := &BaseSend{
		hookManager:         hookManager,
		handlers:            make(map[string]SendHandler),
		packetConstructions: make(map[string]common.PacketConstruction),
		packetLUT:           make(map[string]string),
		packetBuilder:       NewPacketBuilder(),
	}
	return bs
}

// Configure configures the send component with server-specific packet constructions.
func (bs *BaseSend) Configure(serverType string, packetConstructions map[string]common.PacketConstruction) error {
	bs.serverType = serverType
	bs.packetConstructions = packetConstructions

	// Build the lookup table and register formats with the packet builder
	for id, construction := range packetConstructions {
		bs.packetLUT[construction.Name] = id
		bs.packetBuilder.RegisterPacket(id, construction.Name, construction.Format, construction.FieldNames)
	}

	return nil
}

// RegisterHandler registers a handler for a specific packet.
func (bs *BaseSend) RegisterHandler(packetName string, handler SendHandler) {
	bs.handlers[packetName] = handler
}

// ConstructPacket constructs a packet from a packet name and arguments.
func (bs *BaseSend) ConstructPacket(packetName string, args map[string]interface{}) ([]byte, error) {
	// Look up the packet ID
	packetID, exists := bs.packetLUT[packetName]
	if !exists {
		return nil, fmt.Errorf("packet not registered: %s", packetName)
	}

	// Call the handler if it exists
	if handler, exists := bs.handlers[packetName]; exists {
		return handler(args)
	}

	// Use the packet builder to construct the packet
	return bs.packetBuilder.BuildPacket(packetID, args)
}

// SendPacket constructs and sends a packet.
func (bs *BaseSend) SendPacket(packetName string, args map[string]interface{}) error {
	// Construct the packet
	packet, err := bs.ConstructPacket(packetName, args)
	if err != nil {
		return err
	}

	// Send the packet
	return bs.SendToServer(packet)
}

// SendToServer sends a raw packet to the server.
func (bs *BaseSend) SendToServer(packet []byte) error {
	if bs.conn == nil {
		return ErrNoConnection
	}

	// Get the message ID
	if len(packet) < 2 {
		return ErrInvalidPacketID
	}
	messageID := fmt.Sprintf("%02X%02X", packet[1], packet[0])

	// Call the hook for this packet
	hookName := "packet_send/" + messageID
	if bs.hookManager != nil {
		args := map[string]interface{}{
			"switch": messageID,
			"data":   packet,
		}
		// Call the hook
		bs.hookManager.CallHook(hookName, args)
		// In the future, we might want to check if the hook modified the packet
		// or indicated that we should not send it
	}

	// Encrypt the message ID if encryption is enabled
	if bs.encryption.useEncrypt {
		if err := bs.encryptMessageID(&packet); err != nil {
			return err
		}
	}

	// Send the packet
	// Use type assertion to check if the connection has a Send method
	if sender, ok := bs.conn.(interface{ Send([]byte) error }); ok {
		if err := sender.Send(packet); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("connection does not implement Send method")
	}

	// Debug output
	if bs.debugMode {
		fmt.Printf("Sent packet: %s [%d bytes]\n", messageID, len(packet))
	}

	return nil
}

// encryptMessageID encrypts the message ID of a packet.
func (bs *BaseSend) encryptMessageID(msg *[]byte) error {
	if len(*msg) < 2 {
		return ErrInvalidPacketID
	}

	// Extract the message ID (little-endian)
	messageID := uint16((*msg)[0]) | (uint16((*msg)[1]) << 8)

	// Check if encryption is enabled
	if bs.encryption.cryptKey3 > 0 {
		// Save the old message ID and key for debugging
		oldMID := messageID
		oldKey := (bs.encryption.cryptKey >> 16) & 0x7FFF

		// Calculate the new encryption key
		bs.encryption.cryptKey = (bs.encryption.cryptKey*bs.encryption.cryptKey3 + bs.encryption.cryptKey2) & 0xFFFFFFFF

		// XOR the message ID with the new key
		messageID = (messageID ^ uint16((bs.encryption.cryptKey>>16)&0x7FFF)) & 0xFFFF

		// Update the message ID in the packet (little-endian)
		(*msg)[0] = byte(messageID & 0xFF)
		(*msg)[1] = byte((messageID >> 8) & 0xFF)

		// Debug output
		if bs.debugMode {
			fmt.Printf("Encrypted MID: [%04X]->[%04X] / KEY: [0x%04X]->[0x%04X]\n",
				oldMID, messageID, oldKey, (bs.encryption.cryptKey>>16)&0x7FFF)
		}
	}

	return nil
}

// SetConnection sets the connection to use for sending packets.
func (bs *BaseSend) SetConnection(conn interface{}) {
	bs.conn = conn
}

// GetConnection returns the current connection.
func (bs *BaseSend) GetConnection() interface{} {
	return bs.conn
}

// SetEncryptionKeys sets the encryption keys for message ID encryption.
func (bs *BaseSend) SetEncryptionKeys(key1, key2, key3 uint32) {
	bs.encryption.cryptKey1 = key1
	bs.encryption.cryptKey2 = key2
	bs.encryption.cryptKey3 = key3
	bs.encryption.useEncrypt = true
}

// SetDebugMode sets the debug mode.
func (bs *BaseSend) SetDebugMode(debug bool) {
	bs.debugMode = debug
}

// GetServerType returns the server type.
func (bs *BaseSend) GetServerType() string {
	return bs.serverType
}

// GetPacketID returns the packet ID for a given packet name.
func (bs *BaseSend) GetPacketID(name string) (string, bool) {
	id, exists := bs.packetLUT[name]
	return id, exists
}

// RegisterHook registers a hook for a specific event.
func (bs *BaseSend) RegisterHook(hookName string, callback hooks.HookCallback) {
	if bs.hookManager != nil {
		bs.hookManager.AddHook(hookName, callback, nil)
	}
}

// EncryptMessageID encrypts the message ID of a packet.
func (bs *BaseSend) EncryptMessageID(msg *[]byte) error {
	return bs.encryptMessageID(msg)
}

// CryptKeys sets the encryption keys for message ID encryption.
func (bs *BaseSend) CryptKeys(key1, key2, key3 uint32) {
	bs.SetEncryptionKeys(key1, key2, key3)
}

// SendRaw sends a raw packet to the server from a space-delimited list of hex byte values.
func (bs *BaseSend) SendRaw(raw string) error {
	// Split the raw string into hex values and convert to bytes
	packet := make([]byte, 0)
	for _, hex := range strings.Split(raw, " ") {
		val, err := strconv.ParseUint(hex, 16, 8)
		if err != nil {
			return err
		}
		packet = append(packet, byte(val))
	}

	// Send the packet
	return bs.SendToServer(packet)
}

// PinEncode encodes a PIN code using the given seed.
func (bs *BaseSend) PinEncode(seed, pin int) string {
	// Use the PINEncryptor to encode the PIN
	encryptor := NewPINEncryptor()
	return encryptor.Encode(seed, pin)
}

// InjectMessage sends a text message to the connected client's party chat.
func (bs *BaseSend) InjectMessage(message string) error {
	if bs.conn == nil {
		return ErrNoConnection
	}

	// Create the message with the format from Send.pm
	name := []byte("|")
	msgContent := []byte(" : " + message)

	// Create the packet
	packet := []byte{0x09, 0x01}

	// Add the length (name length + message length + 12)
	length := uint16(len(name) + len(msgContent) + 12)
	lengthBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(lengthBytes, length)
	packet = append(packet, lengthBytes...)

	// Add padding (4 bytes of zeros)
	packet = append(packet, 0, 0, 0, 0)

	// Add the name and message content
	packet = append(packet, name...)
	packet = append(packet, msgContent...)

	// Add null terminator
	packet = append(packet, 0)

	// Send the packet to the client
	if sender, ok := bs.conn.(interface{ ClientSend([]byte) error }); ok {
		if err := sender.ClientSend(packet); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("connection does not implement ClientSend method")
	}

	return nil
}

// InjectAdminMessage sends a text message to the connected client's system chat.
func (bs *BaseSend) InjectAdminMessage(message string) error {
	if bs.conn == nil {
		return ErrNoConnection
	}

	// Convert message to bytes
	msgBytes := []byte(message)

	// Create the packet
	packet := []byte{0x9A, 0x00}

	// Add the length (message length + 5)
	length := uint16(len(msgBytes) + 5)
	lengthBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(lengthBytes, length)
	packet = append(packet, lengthBytes...)

	// Add the message content
	packet = append(packet, msgBytes...)

	// Add null terminator
	packet = append(packet, 0)

	// Send the packet to the client
	if sender, ok := bs.conn.(interface{ ClientSend([]byte) error }); ok {
		if err := sender.ClientSend(packet); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("connection does not implement ClientSend method")
	}

	return nil
}

// Reconstruct constructs a packet from a packet ID and arguments.
func (bs *BaseSend) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	return bs.packetBuilder.BuildPacket(packetID, args)
}

// RegisterPacketHandler registers a handler for a packet.
func (bs *BaseSend) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	// Register the packet with the packet builder
	bs.packetBuilder.RegisterPacket(packetID, name, format, keys)

	// Register the packet ID in the lookup table
	bs.packetLUT[name] = packetID

	// Register the handler if provided
	if handler != nil {
		bs.handlers[name] = func(args map[string]interface{}) ([]byte, error) {
			// Call the handler to process the arguments
			if err := handler(args); err != nil {
				return nil, err
			}

			// Construct the packet using the processed arguments
			return bs.packetBuilder.BuildPacket(packetID, args)
		}
	}
}

// GetTime returns the current time in milliseconds.
func (bs *BaseSend) GetTime() uint32 {
	return uint32(time.Now().UnixNano() / int64(time.Millisecond))
}
