// Package protocol provides functionality for handling the Ragnarok Online network protocol.
// This file implements the padded packets system which handles special packet formatting for security.
package protocol

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
)

// PaddedPackets handles the generation of specially formatted packets
// for security purposes. Some RO servers use padded packets to prevent
// packet spoofing and bot detection.
type PaddedPackets struct {
	// Configuration
	enabled    bool
	attackID   uint16
	skillUseID uint16

	// Hash data
	accountID uint32
	mapSync   uint32
	sync      uint32

	// Buffer for packet generation
	buffer []byte
}

// NewPaddedPackets creates a new padded packets handler
func NewPaddedPackets() *PaddedPackets {
	return &PaddedPackets{
		enabled:    false,
		attackID:   0x89,
		skillUseID: 0x113,
		buffer:     make([]byte, 512),
	}
}

// SetEnabled enables or disables padded packets
func (p *PaddedPackets) SetEnabled(enabled bool) {
	p.enabled = enabled
}

// IsEnabled returns whether padded packets are enabled
func (p *PaddedPackets) IsEnabled() bool {
	return p.enabled
}

// SetPacketIDs sets the packet IDs for attack and skill use
func (p *PaddedPackets) SetPacketIDs(attackID, skillUseID uint16) {
	p.attackID = attackID
	p.skillUseID = skillUseID
}

// SetHashData sets the hash data used for packet generation
func (p *PaddedPackets) SetHashData(accountID, mapSync, sync uint32) {
	p.accountID = accountID
	p.mapSync = mapSync
	p.sync = sync
}

// SetAccountID sets the account ID used for hash generation
func (p *PaddedPackets) SetAccountID(accountID uint32) {
	p.accountID = accountID
}

// SetMapSync sets the map sync used for hash generation
func (p *PaddedPackets) SetMapSync(mapSync uint32) {
	p.mapSync = mapSync
}

// SetSync sets the sync used for hash generation
func (p *PaddedPackets) SetSync(sync uint32) {
	p.sync = sync
}

// Reset resets the padded packets handler to its default state
func (p *PaddedPackets) Reset() {
	p.enabled = false
	p.attackID = 0x89
	p.skillUseID = 0x113
	p.accountID = 0
	p.mapSync = 0
	p.sync = 0
}

// GenerateSitStand generates a sit/stand packet
func (p *PaddedPackets) GenerateSitStand(sit bool) []byte {
	if !p.enabled {
		// If padded packets are disabled, return a simple packet
		packet := make([]byte, 2)
		binary.LittleEndian.PutUint16(packet, p.attackID)
		if sit {
			packet = append(packet, 0x02)
		} else {
			packet = append(packet, 0x03)
		}
		return packet
	}

	// Create input keys for the padded packet
	var sitValue uint32 = 3
	if sit {
		sitValue = 2
	}

	// Generate the padded packet
	return p.generatePaddedPacket(p.attackID, []uint32{sitValue})
}

// GenerateAttack generates an attack packet
func (p *PaddedPackets) GenerateAttack(targetID uint32, flag uint32) []byte {
	if !p.enabled {
		// If padded packets are disabled, return a simple packet
		packet := make([]byte, 6)
		binary.LittleEndian.PutUint16(packet[0:2], p.attackID)
		binary.LittleEndian.PutUint32(packet[2:6], targetID)
		return packet
	}

	// Create input keys for the padded packet
	return p.generatePaddedPacket(p.attackID, []uint32{targetID, 7})
}

// GenerateSkillUse generates a skill use packet
func (p *PaddedPackets) GenerateSkillUse(skillID, skillLv, targetID uint32) []byte {
	if !p.enabled {
		// If padded packets are disabled, return a simple packet
		packet := make([]byte, 10)
		binary.LittleEndian.PutUint16(packet[0:2], p.skillUseID)
		binary.LittleEndian.PutUint16(packet[2:4], uint16(skillID))
		binary.LittleEndian.PutUint16(packet[4:6], uint16(skillLv))
		binary.LittleEndian.PutUint32(packet[6:10], targetID)
		return packet
	}

	// Create input keys for the padded packet
	return p.generatePaddedPacket(p.skillUseID, []uint32{skillLv, skillID, targetID})
}

// generatePaddedPacket generates a padded packet with the given packet ID and input keys
func (p *PaddedPackets) generatePaddedPacket(packetID uint16, inputKeys []uint32) []byte {
	// Clear the buffer
	for i := range p.buffer {
		p.buffer[i] = 0
	}

	// Set the packet ID
	binary.LittleEndian.PutUint16(p.buffer[0:2], packetID)

	// Generate the hash data
	hashData := p.createHash(p.mapSync, p.sync, p.accountID, packetID)

	// Calculate the packet length
	packetLength := (1 + len(inputKeys)) * 4

	// Offsets used in the padding algorithm
	offsets := []uint32{15, 14, 12, 9, 5, 0}

	// Apply the padding algorithm
	for iter := 0; iter <= 5; iter++ {
		packetLength = (1 + len(inputKeys)) * 4

		intCtr := uint32(5)
		writePtr := 4

		for pass := 0; pass < len(inputKeys); pass++ {
			magic := ((intCtr * uint32(pass)) + (hashData - offsets[iter])) % 0x27
			packetLength += int(magic)
			intCtr += 3

			writePtr += (4 + int(magic))
			binary.LittleEndian.PutUint32(p.buffer[writePtr-4:writePtr], inputKeys[pass]+uint32(iter)-5)
		}
	}

	// Set the packet length
	p.buffer[2] = byte(packetLength)

	// Copy the result to a new slice
	result := make([]byte, packetLength)
	copy(result, p.buffer[:packetLength])

	return result
}

// createHash generates a hash value based on the given parameters
func (p *PaddedPackets) createHash(mapSync, sync, accountID uint32, packetID uint16) uint32 {
	// Calculate which algorithm to use
	slot := (uint32(packetID)*uint32(packetID) + mapSync + sync + accountID) & 0xF

	// Calculate the key for the hash function
	key := uint32(packetID)*accountID + mapSync*sync

	// Use the appropriate hash function
	return p.hashWithAlgorithm(int(slot), key)
}

// hashWithAlgorithm applies the specified hashing algorithm to the key
func (p *PaddedPackets) hashWithAlgorithm(algorithmID int, key uint32) uint32 {
	// Convert key to bytes
	keyBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(keyBytes, key)

	var h hash.Hash

	// Select the appropriate hash function
	// In the original implementation, there are 16 different algorithms
	// For simplicity, we'll use a few common hash functions and some placeholders
	switch algorithmID {
	case 0:
		// Use MD5
		h = md5.New()
	case 1:
		// Use SHA1
		h = sha1.New()
	case 2:
		// Use SHA256
		h = sha256.New()
	case 3:
		// Use a simple hash function for algorithm 3
		return key*0x9E3779B9 + 0x12345678
	default:
		// For other algorithms, use a simple hash function
		// This is a placeholder and should be replaced with the actual algorithms
		return key ^ 0xDEADBEEF ^ uint32(algorithmID*1337)
	}

	// Calculate the hash
	h.Write(keyBytes)
	hashBytes := h.Sum(nil)

	// Return the first 4 bytes as a uint32
	return binary.LittleEndian.Uint32(hashBytes[:4])
}

// DecodePacket decodes a padded packet and extracts the keys
func (p *PaddedPackets) DecodePacket(packet []byte, keyCount int) []uint32 {
	if len(packet) < 2 {
		return nil
	}

	// Extract the packet ID
	packetID := binary.LittleEndian.Uint16(packet[0:2])

	// Generate the hash data
	hashData := p.createHash(p.mapSync, p.sync, p.accountID, packetID)

	// Extract the keys
	keys := make([]uint32, keyCount)

	// Offsets used in the padding algorithm
	offsets := []uint32{15, 14, 12, 9, 5, 0}

	// Use the last iteration's offsets (iter = 5)
	offset := offsets[5]

	intCtr := uint32(5)
	readPtr := 4

	for pass := 0; pass < keyCount; pass++ {
		magic := ((intCtr * uint32(pass)) + (hashData - offset)) % 0x27
		intCtr += 3

		readPtr += (4 + int(magic))
		if readPtr+4 <= len(packet) {
			// In the encode function, we add (iter - 5) to each key
			// In the final iteration, iter = 5, so we add 0
			// So the decoded keys should match the input keys
			keys[pass] = binary.LittleEndian.Uint32(packet[readPtr-4 : readPtr])
		}
	}

	return keys
}

// String returns a string representation of the PaddedPackets
func (p *PaddedPackets) String() string {
	return fmt.Sprintf("PaddedPackets{enabled=%v, attackID=0x%04X, skillUseID=0x%04X, accountID=%d, mapSync=%d, sync=%d}",
		p.enabled, p.attackID, p.skillUseID, p.accountID, p.mapSync, p.sync)
}
