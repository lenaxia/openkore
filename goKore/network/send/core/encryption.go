package core

import (
	"encoding/binary"
	"math/big"
)

// MessageEncryptor handles message ID encryption for RO packets.
type MessageEncryptor struct {
	// Encryption keys
	cryptKey  uint32
	cryptKey1 uint32
	cryptKey2 uint32
	cryptKey3 uint32

	// Whether encryption is enabled
	enabled bool

	// Whether we're in game (affects encryption)
	inGame bool

	// Map login packet ID (used to reset encryption)
	mapLoginID string
}

// NewMessageEncryptor creates a new message encryptor.
func NewMessageEncryptor() *MessageEncryptor {
	return &MessageEncryptor{
		cryptKey:   0,
		cryptKey1:  0,
		cryptKey2:  0,
		cryptKey3:  0,
		enabled:    false,
		inGame:     false,
		mapLoginID: "0072", // Default map login packet ID
	}
}

// SetKeys sets the encryption keys.
func (me *MessageEncryptor) SetKeys(key1, key2, key3 uint32) {
	me.cryptKey1 = key1
	me.cryptKey2 = key2
	me.cryptKey3 = key3
	me.enabled = true
}

// SetMapLoginID sets the map login packet ID.
func (me *MessageEncryptor) SetMapLoginID(id string) {
	me.mapLoginID = id
}

// SetInGame sets whether we're in game.
func (me *MessageEncryptor) SetInGame(inGame bool) {
	me.inGame = inGame
}

// EncryptMessageID encrypts the message ID of a packet.
func (me *MessageEncryptor) EncryptMessageID(msg *[]byte) error {
	if len(*msg) < 2 {
		return ErrInvalidPacketID
	}

	// Extract the message ID
	messageID := binary.LittleEndian.Uint16((*msg)[:2])
	messageIDHex := getHexString(messageID)

	// Check if we need to reset the encryption key
	if me.enabled && messageIDHex == me.mapLoginID {
		me.cryptKey = me.cryptKey1
	}

	// If not in game, turn off keys and return
	if !me.inGame {
		me.cryptKey = 0
		return nil
	}

	// Check if encryption is enabled and we're in game
	if me.enabled && me.cryptKey > 0 && me.inGame {
		// Save old values for debugging if needed
		// oldMID := messageID
		// oldKey := (me.cryptKey >> 16) & 0x7FFF

		// Calculate the new encryption key
		me.cryptKey = (me.cryptKey*me.cryptKey3 + me.cryptKey2) & 0xFFFFFFFF

		// XOR the message ID with the new key
		messageID = (messageID ^ uint16((me.cryptKey>>16)&0x7FFF)) & 0xFFFF

		// Update the message ID in the packet
		binary.LittleEndian.PutUint16((*msg)[:2], messageID)
	} else if me.cryptKey3 == 0 && me.inGame && me.cryptKey > 0 {
		// Legacy encryption method (enc_val1/enc_val2)
		// Only apply if we're in game and using legacy encryption

		// Prepare encryption
		me.cryptKey = ((0x000343FD * me.cryptKey) + me.cryptKey2) & 0xFFFFFFFF

		// Encrypt message ID
		messageID = (messageID ^ uint16((me.cryptKey>>16)&0x7FFF)) & 0xFFFF
		binary.LittleEndian.PutUint16((*msg)[:2], messageID)
	}

	return nil
}

// getHexString returns the hex string representation of a uint16.
func getHexString(val uint16) string {
	return big.NewInt(int64(val)).Text(16)
}

// PINEncryptor handles PIN code encryption.
type PINEncryptor struct{}

// NewPINEncryptor creates a new PIN encryptor.
func NewPINEncryptor() *PINEncryptor {
	return &PINEncryptor{}
}

// Encode encodes a PIN code using the given seed.
// This is a Go implementation of the pinEncode function from Send.pm
func (pe *PINEncryptor) Encode(seed, pin int) string {
	// Convert seed to big.Int for precise arithmetic
	seedBig := big.NewInt(int64(seed))

	// Constants for the algorithm
	mulFactor := big.NewInt(0x3498)
	addFactor := big.NewInt(0x881234)

	// Create a slice of keypad keys
	keypadKeysOrder := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	// Calculate keys order (randomized based on seed value)
	if len(keypadKeysOrder) >= 1 {
		k := 2
		for pos := 1; pos < len(keypadKeysOrder); pos++ {
			// Calculate next seed value
			seedBig.Mul(seedBig, mulFactor)
			seedBig.Add(seedBig, addFactor)
			seedBig.And(seedBig, big.NewInt(0xFFFFFFFF))

			// Calculate replacement position
			replacePos := int(new(big.Int).Mod(seedBig, big.NewInt(int64(k))).Int64())

			// Swap values if positions are different
			if pos != replacePos {
				keypadKeysOrder[pos], keypadKeysOrder[replacePos] = keypadKeysOrder[replacePos], keypadKeysOrder[pos]
			}

			k++
		}
	}

	// Create a map of key values to positions
	keypad := make(map[string]int)
	for pos, key := range keypadKeysOrder {
		keypad[key] = pos
	}

	// Convert PIN to string and process each digit
	pinStr := big.NewInt(int64(pin)).Text(10)
	pinReply := ""
	for _, digit := range pinStr {
		pinReply += big.NewInt(int64(keypad[string(digit)])).Text(10)
	}

	return pinReply
}
