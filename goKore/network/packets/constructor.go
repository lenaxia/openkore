// Package packets provides functionality for handling Ragnarok Online network packets.
// This file implements packet construction for outgoing packets.
package packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"

	"github.com/lenaxia/goKore/network/connection"
	"github.com/lenaxia/goKore/network/protocol"
)

// PacketConstructor is responsible for constructing outgoing packets
type PacketConstructor struct {
	// Database of packet definitions
	packetDB *PacketDatabase

	// Encryption keys for packet IDs
	cryptKey1 *big.Int
	cryptKey2 *big.Int
	cryptKey3 *big.Int

	// Current encryption key
	currentCryptKey *big.Int

	// Padded packets handler
	paddedPackets *protocol.PaddedPackets

	// Network state
	networkState int

	// Debug mode
	debugMode bool
}

// NewPacketConstructor creates a new packet constructor
func NewPacketConstructor(packetDB *PacketDatabase) *PacketConstructor {
	return &PacketConstructor{
		packetDB:      packetDB,
		paddedPackets: protocol.NewPaddedPackets(),
		debugMode:     false,
	}
}

// SetCryptKeys sets the encryption keys for packet IDs
func (c *PacketConstructor) SetCryptKeys(key1, key2, key3 uint32) {
	c.cryptKey1 = new(big.Int).SetUint64(uint64(key1))
	c.cryptKey2 = new(big.Int).SetUint64(uint64(key2))
	c.cryptKey3 = new(big.Int).SetUint64(uint64(key3))
}

// SetNetworkState sets the current network state
func (c *PacketConstructor) SetNetworkState(state int) {
	c.networkState = state
}

// SetDebugMode enables or disables debug mode
func (c *PacketConstructor) SetDebugMode(debug bool) {
	c.debugMode = debug
}

// SetPaddedPacketsEnabled enables or disables padded packets
func (c *PacketConstructor) SetPaddedPacketsEnabled(enabled bool) {
	c.paddedPackets.SetEnabled(enabled)
}

// SetPaddedPacketsData sets the data needed for padded packets
func (c *PacketConstructor) SetPaddedPacketsData(accountID, mapSync, sync uint32) {
	c.paddedPackets.SetHashData(accountID, mapSync, sync)
}

// EncryptMessageID encrypts the message ID of a packet
func (c *PacketConstructor) EncryptMessageID(packet []byte) []byte {
	if len(packet) < 2 {
		return packet
	}

	// Extract the message ID
	messageID := binary.LittleEndian.Uint16(packet[0:2])

	// Check if we're in game
	if c.networkState != int(connection.IN_GAME) {
		// Turn off encryption
		c.currentCryptKey = big.NewInt(0)
		return packet
	}

	// Check if encryption is active
	if c.currentCryptKey != nil && c.currentCryptKey.Uint64() > 0 && c.cryptKey3 != nil {
		// Save old values for debugging
		oldMID := messageID
		oldKey := (c.currentCryptKey.Uint64() >> 16) & 0x7FFF

		// Calculate the new encryption key
		c.currentCryptKey = new(big.Int).Mul(c.currentCryptKey, c.cryptKey3)
		c.currentCryptKey = new(big.Int).Add(c.currentCryptKey, c.cryptKey2)
		c.currentCryptKey = new(big.Int).And(c.currentCryptKey, big.NewInt(0xFFFFFFFF))

		// XOR the message ID
		messageID = (messageID ^ uint16((c.currentCryptKey.Uint64()>>16)&0x7FFF)) & 0xFFFF

		// Update the packet with the encrypted message ID
		binary.LittleEndian.PutUint16(packet[0:2], messageID)

		// Debug log
		if c.debugMode {
			fmt.Printf("Encrypted MID : [%04X]->[%04X] / KEY : [0x%04X]->[0x%04X]\n",
				oldMID, messageID, oldKey, (c.currentCryptKey.Uint64()>>16)&0x7FFF)
		}
	}

	return packet
}

// ConstructPacket constructs a packet using the given packet name and arguments
func (c *PacketConstructor) ConstructPacket(name string, args map[string]interface{}) ([]byte, error) {
	// Look up the packet definition
	def, exists := c.packetDB.GetPacketByName(name)
	if !exists {
		return nil, fmt.Errorf("packet not found: %s", name)
	}

	// Create a buffer for the packet
	buf := new(bytes.Buffer)

	// Write the packet ID
	binary.Write(buf, binary.LittleEndian, def.ID)

	// Parse the format string and write the data
	format := def.Format
	if format == "" {
		// If no format is specified, just return the packet ID
		return buf.Bytes(), nil
	}

	// Split the format string into fields
	fields := strings.Fields(format)
	for i, field := range fields {
		if i >= len(def.ParamNames) {
			break
		}

		key := def.ParamNames[i]
		value, ok := args[key]
		if !ok {
			return nil, fmt.Errorf("missing argument: %s", key)
		}

		switch field[0] {
		case 'C': // unsigned char (1 byte)
			if v, ok := value.(uint8); ok {
				binary.Write(buf, binary.LittleEndian, v)
			} else if v, ok := value.(int); ok {
				binary.Write(buf, binary.LittleEndian, uint8(v))
			} else {
				return nil, fmt.Errorf("invalid type for field %s: expected uint8, got %T", key, value)
			}
		case 'v': // unsigned short (2 bytes)
			if v, ok := value.(uint16); ok {
				binary.Write(buf, binary.LittleEndian, v)
			} else if v, ok := value.(int); ok {
				binary.Write(buf, binary.LittleEndian, uint16(v))
			} else {
				return nil, fmt.Errorf("invalid type for field %s: expected uint16, got %T", key, value)
			}
		case 'V': // unsigned int (4 bytes)
			if v, ok := value.(uint32); ok {
				binary.Write(buf, binary.LittleEndian, v)
			} else if v, ok := value.(int); ok {
				binary.Write(buf, binary.LittleEndian, uint32(v))
			} else {
				return nil, fmt.Errorf("invalid type for field %s: expected uint32, got %T", key, value)
			}
		case 'a': // fixed-length string
			length, _ := parseLength(field)
			if v, ok := value.(string); ok {
				writeFixedString(buf, v, length)
			} else if v, ok := value.([]byte); ok {
				writeFixedBytes(buf, v, length)
			} else {
				return nil, fmt.Errorf("invalid type for field %s: expected string or []byte, got %T", key, value)
			}
		case 'Z': // null-terminated string
			if v, ok := value.(string); ok {
				buf.WriteString(v)
				buf.WriteByte(0)
			} else {
				return nil, fmt.Errorf("invalid type for field %s: expected string, got %T", key, value)
			}
		case 'f': // float (4 bytes)
			if v, ok := value.(float32); ok {
				binary.Write(buf, binary.LittleEndian, v)
			} else if v, ok := value.(float64); ok {
				binary.Write(buf, binary.LittleEndian, float32(v))
			} else {
				return nil, fmt.Errorf("invalid type for field %s: expected float32, got %T", key, value)
			}
		case 'd': // double (8 bytes)
			if v, ok := value.(float64); ok {
				binary.Write(buf, binary.LittleEndian, v)
			} else {
				return nil, fmt.Errorf("invalid type for field %s: expected float64, got %T", key, value)
			}
		default:
			return nil, fmt.Errorf("unknown format specifier: %s", field)
		}
	}

	// Encrypt the message ID if needed
	packet := buf.Bytes()
	packet = c.EncryptMessageID(packet)

	return packet, nil
}

// Helper function to parse length from format specifier
func parseLength(format string) (int, error) {
	if len(format) < 2 {
		return 0, fmt.Errorf("invalid format specifier: %s", format)
	}
	if format[1] == '*' {
		return 0, nil // variable length
	}
	var length int
	_, err := fmt.Sscanf(format[1:], "%d", &length)
	if err != nil {
		return 0, fmt.Errorf("invalid length in format specifier: %s", format)
	}
	return length, nil
}

// Helper function to write a fixed-length string
func writeFixedString(buf *bytes.Buffer, s string, length int) {
	if len(s) > length {
		s = s[:length]
	}
	buf.WriteString(s)
	for i := len(s); i < length; i++ {
		buf.WriteByte(0)
	}
}

// Helper function to write fixed-length bytes
func writeFixedBytes(buf *bytes.Buffer, b []byte, length int) {
	if len(b) > length {
		b = b[:length]
	}
	buf.Write(b)
	for i := len(b); i < length; i++ {
		buf.WriteByte(0)
	}
}

// PinEncode encodes a PIN using the given seed
func (c *PacketConstructor) PinEncode(seed, pin int) string {
	// Constants for the algorithm
	mulFactor := 0x3498
	addFactor := 0x881234

	// Create a big.Int for the seed
	seedInt := big.NewInt(int64(seed))

	// Calculate keys order (they are randomized based on seed value)
	keypadKeysOrder := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	if len(keypadKeysOrder) >= 1 {
		k := 2
		for pos := 1; pos < len(keypadKeysOrder); pos++ {
			// Calculate next seed value
			seedInt = new(big.Int).Mul(seedInt, big.NewInt(int64(mulFactor)))
			seedInt = new(big.Int).Add(seedInt, big.NewInt(int64(addFactor)))
			seedInt = new(big.Int).And(seedInt, big.NewInt(0xFFFFFFFF))

			replacePos := int(seedInt.Int64()) % k
			if pos != replacePos {
				// Swap values
				keypadKeysOrder[pos], keypadKeysOrder[replacePos] = keypadKeysOrder[replacePos], keypadKeysOrder[pos]
			}
			k++
		}
	}

	// Associate keys values with their position using a map
	keypad := make(map[byte]int)
	for pos, key := range keypadKeysOrder {
		keypad[key] = pos
	}

	// Encode the PIN
	pinReply := ""
	pinStr := fmt.Sprintf("%d", pin)
	for _, digit := range pinStr {
		pinReply += fmt.Sprintf("%d", keypad[byte(digit)])
	}

	return pinReply
}

// SendToServer sends a packet to the server
func (c *PacketConstructor) SendToServer(packet []byte) []byte {
	// Encrypt the message ID if needed
	packet = c.EncryptMessageID(packet)
	return packet
}

// SendRaw sends a raw packet to the server
func (c *PacketConstructor) SendRaw(raw string) ([]byte, error) {
	// Parse the raw string
	fields := strings.Fields(raw)
	packet := make([]byte, len(fields))
	for i, field := range fields {
		b, err := parseHexByte(field)
		if err != nil {
			return nil, err
		}
		packet[i] = b
	}
	return c.SendToServer(packet), nil
}

// Helper function to parse a hex byte
func parseHexByte(s string) (byte, error) {
	var b byte
	_, err := fmt.Sscanf(s, "%x", &b)
	return b, err
}

// GenerateSitStand generates a sit/stand packet
func (c *PacketConstructor) GenerateSitStand(sit bool) []byte {
	return c.paddedPackets.GenerateSitStand(sit)
}

// GenerateAttack generates an attack packet
func (c *PacketConstructor) GenerateAttack(targetID uint32, flag uint32) []byte {
	return c.paddedPackets.GenerateAttack(targetID, flag)
}

// GenerateSkillUse generates a skill use packet
func (c *PacketConstructor) GenerateSkillUse(skillID, skillLv, targetID uint32) []byte {
	return c.paddedPackets.GenerateSkillUse(skillID, skillLv, targetID)
}

// InjectMessage injects a message into the client's party chat
func (c *PacketConstructor) InjectMessage(message string) []byte {
	name := "|"
	msg := name + " : " + message + "\x00"

	// Create the packet
	packet := make([]byte, 2+2+4+len(name)+len(message)+3)
	packet[0] = 0x09
	packet[1] = 0x01

	// Set the length
	binary.LittleEndian.PutUint16(packet[2:4], uint16(len(name)+len(message)+12))

	// Copy the message
	copy(packet[8:], []byte(msg))

	return packet
}

// InjectAdminMessage injects a message into the client's system chat
func (c *PacketConstructor) InjectAdminMessage(message string) []byte {
	// Create the packet
	packet := make([]byte, 2+2+len(message)+1)
	packet[0] = 0x9A
	packet[1] = 0x00

	// Set the length
	binary.LittleEndian.PutUint16(packet[2:4], uint16(len(message)+5))

	// Copy the message
	copy(packet[4:], []byte(message))
	packet[len(packet)-1] = 0 // Null terminator

	return packet
}

// String returns a string representation of the packet constructor
func (c *PacketConstructor) String() string {
	return fmt.Sprintf("PacketConstructor{debugMode=%v, networkState=%d, paddedPackets=%v}",
		c.debugMode, c.networkState, c.paddedPackets)
}
