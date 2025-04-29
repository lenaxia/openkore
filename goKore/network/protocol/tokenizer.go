// Package protocol provides functionality for handling the Ragnarok Online network protocol.
// This file implements the message tokenizer which extracts discrete packets from byte streams.
package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// MessageType represents the type of message received
type MessageType int

const (
	// KnownMessage indicates a recognized packet type
	KnownMessage MessageType = iota
	// UnknownMessage indicates an unrecognized packet type
	UnknownMessage
	// AccountID indicates a special message containing the account ID
	AccountID
)

// Errors
var (
	ErrIncompletePacket = errors.New("incomplete packet")
	ErrInvalidPacket    = errors.New("invalid packet")
)

// PacketDef contains information about packet structure
type PacketDef struct {
	Length    int  // Fixed length or -1 for variable length
	HasLength bool // Whether packet has length field
}

// Tokenizer handles breaking byte streams into discrete packets
type Tokenizer struct {
	buffer               []byte
	packetDefs           map[string]PacketDef
	nextMightBeAccountID bool
}

// NewTokenizer creates a new message tokenizer
func NewTokenizer(packetDefs map[string]PacketDef) *Tokenizer {
	return &Tokenizer{
		buffer:     make([]byte, 0),
		packetDefs: packetDefs,
	}
}

// Add appends data to the buffer
func (t *Tokenizer) Add(data []byte) {
	t.buffer = append(t.buffer, data...)
}

// Clear removes data from the buffer
func (t *Tokenizer) Clear(size int) {
	if size <= 0 || size > len(t.buffer) {
		t.buffer = make([]byte, 0)
	} else {
		t.buffer = t.buffer[size:]
	}
}

// GetBuffer returns the current buffer contents
func (t *Tokenizer) GetBuffer() []byte {
	return t.buffer
}

// NextMessageMightBeAccountID marks that the next message might be an account ID
func (t *Tokenizer) NextMessageMightBeAccountID() {
	t.nextMightBeAccountID = true
}

// GetMessageID extracts the message ID from a packet
func GetMessageID(packet []byte) string {
	if len(packet) < 2 {
		return ""
	}
	return fmt.Sprintf("%02X%02X", packet[1], packet[0])
}

// ReadNext extracts the next complete packet from the buffer
func (t *Tokenizer) ReadNext() ([]byte, MessageType, error) {
	if len(t.buffer) < 2 {
		return nil, UnknownMessage, ErrIncompletePacket
	}

	// Check if this might be an account ID
	if t.nextMightBeAccountID && len(t.buffer) >= 4 {
		t.nextMightBeAccountID = false
		// In a real implementation, we would compare with the global accountID
		// For now, just return the first 4 bytes as an account ID
		accountIDBytes := make([]byte, 4)
		copy(accountIDBytes, t.buffer[:4])
		t.buffer = t.buffer[4:]
		return accountIDBytes, AccountID, nil
	}

	// Get packet ID and look up definition
	packetID := GetMessageID(t.buffer)
	packetDef, exists := t.packetDefs[packetID]

	if !exists {
		// Unknown packet type
		// In a real implementation, we might want to handle this differently
		// For now, just return the entire buffer as an unknown message
		result := make([]byte, len(t.buffer))
		copy(result, t.buffer)
		t.buffer = make([]byte, 0)
		return result, UnknownMessage, nil
	}

	if packetDef.Length > 0 {
		// Fixed length packet
		if len(t.buffer) < packetDef.Length {
			return nil, UnknownMessage, ErrIncompletePacket
		}

		result := make([]byte, packetDef.Length)
		copy(result, t.buffer[:packetDef.Length])
		t.buffer = t.buffer[packetDef.Length:]
		return result, KnownMessage, nil
	} else if packetDef.HasLength {
		// Variable length packet
		if len(t.buffer) < 4 {
			return nil, UnknownMessage, ErrIncompletePacket
		}

		length := int(binary.LittleEndian.Uint16(t.buffer[2:4]))
		if len(t.buffer) < length {
			return nil, UnknownMessage, ErrIncompletePacket
		}

		result := make([]byte, length)
		copy(result, t.buffer[:length])
		t.buffer = t.buffer[length:]
		return result, KnownMessage, nil
	}

	// Should never reach here if packet definitions are correct
	return nil, UnknownMessage, ErrInvalidPacket
}
