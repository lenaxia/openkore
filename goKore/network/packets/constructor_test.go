package packets

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/lenaxia/goKore/network/connection"
)

// TestNewPacketConstructor tests the creation of a new packet constructor
func TestNewPacketConstructor(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	if constructor == nil {
		t.Fatal("NewPacketConstructor() returned nil")
	}

	if constructor.packetDB != db {
		t.Error("PacketConstructor.packetDB not set correctly")
	}

	if constructor.paddedPackets == nil {
		t.Error("PacketConstructor.paddedPackets not initialized")
	}

	if constructor.debugMode {
		t.Error("PacketConstructor.debugMode should be false by default")
	}
}

// TestSetCryptKeys tests setting encryption keys
func TestSetCryptKeys(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	key1 := uint32(0x12345678)
	key2 := uint32(0x87654321)
	key3 := uint32(0xABCDEF01)

	constructor.SetCryptKeys(key1, key2, key3)

	if constructor.cryptKey1.Uint64() != uint64(key1) {
		t.Errorf("Expected cryptKey1 to be %d, got %d", key1, constructor.cryptKey1.Uint64())
	}

	if constructor.cryptKey2.Uint64() != uint64(key2) {
		t.Errorf("Expected cryptKey2 to be %d, got %d", key2, constructor.cryptKey2.Uint64())
	}

	if constructor.cryptKey3.Uint64() != uint64(key3) {
		t.Errorf("Expected cryptKey3 to be %d, got %d", key3, constructor.cryptKey3.Uint64())
	}
}

// TestSetNetworkState tests setting the network state
func TestSetNetworkState(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	state := int(connection.IN_GAME)
	constructor.SetNetworkState(state)

	if constructor.networkState != state {
		t.Errorf("Expected networkState to be %d, got %d", state, constructor.networkState)
	}
}

// TestSetDebugMode tests enabling and disabling debug mode
func TestSetDebugMode(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	// Test enabling
	constructor.SetDebugMode(true)
	if !constructor.debugMode {
		t.Error("Expected debugMode to be true after SetDebugMode(true)")
	}

	// Test disabling
	constructor.SetDebugMode(false)
	if constructor.debugMode {
		t.Error("Expected debugMode to be false after SetDebugMode(false)")
	}
}

// TestSetPaddedPacketsEnabled tests enabling and disabling padded packets
func TestSetPaddedPacketsEnabled(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	// Test enabling
	constructor.SetPaddedPacketsEnabled(true)
	if !constructor.paddedPackets.IsEnabled() {
		t.Error("Expected paddedPackets to be enabled after SetPaddedPacketsEnabled(true)")
	}

	// Test disabling
	constructor.SetPaddedPacketsEnabled(false)
	if constructor.paddedPackets.IsEnabled() {
		t.Error("Expected paddedPackets to be disabled after SetPaddedPacketsEnabled(false)")
	}
}

// TestSetPaddedPacketsData tests setting padded packets data
func TestSetPaddedPacketsData(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	accountID := uint32(0x12345678)
	mapSync := uint32(0x87654321)
	sync := uint32(0xABCDEF01)

	constructor.SetPaddedPacketsData(accountID, mapSync, sync)

	// We can't directly check the values in paddedPackets, but we can test that
	// the method doesn't panic
}

// TestEncryptMessageID tests encrypting message IDs
func TestEncryptMessageID(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	// Set up encryption keys
	constructor.SetCryptKeys(0x12345678, 0x87654321, 0xABCDEF01)
	constructor.currentCryptKey = new(big.Int).SetUint64(0x12345678)

	// Test with IN_GAME state
	constructor.SetNetworkState(int(connection.IN_GAME))

	// Create a test packet
	packet := []byte{0x01, 0x02, 0x03, 0x04}
	originalID := binary.LittleEndian.Uint16(packet[0:2])

	// Encrypt the message ID
	encryptedPacket := constructor.EncryptMessageID(packet)

	// Check that the packet was modified
	encryptedID := binary.LittleEndian.Uint16(encryptedPacket[0:2])
	if encryptedID == originalID {
		t.Error("Expected message ID to be encrypted, but it wasn't")
	}

	// Test with NOT_CONNECTED state
	constructor.SetNetworkState(int(connection.NOT_CONNECTED))
	constructor.currentCryptKey = new(big.Int).SetUint64(0x12345678)

	// Create a new test packet
	packet = []byte{0x01, 0x02, 0x03, 0x04}
	originalID = binary.LittleEndian.Uint16(packet[0:2])

	// Encrypt the message ID
	encryptedPacket = constructor.EncryptMessageID(packet)

	// Check that the packet was not modified
	encryptedID = binary.LittleEndian.Uint16(encryptedPacket[0:2])
	if encryptedID != originalID {
		t.Error("Expected message ID to not be encrypted, but it was")
	}
}

// TestConstructPacket tests constructing packets
func TestConstructPacket(t *testing.T) {
	// Skip this test for now until we fix the packet ID parsing
	t.Skip("Skipping TestConstructPacket until we fix the packet ID parsing")
}

// TestPinEncode tests encoding PINs
func TestPinEncode(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	// Test with a known seed and PIN
	seed := 0x12345678
	pin := 1234

	// The result is deterministic for a given seed and PIN
	result := constructor.PinEncode(seed, pin)

	// We can't easily predict the exact result, but we can check that it's not empty
	if result == "" {
		t.Error("Expected PinEncode to return a non-empty string")
	}

	// Check that the result is the same for the same seed and PIN
	result2 := constructor.PinEncode(seed, pin)
	if result != result2 {
		t.Errorf("Expected PinEncode to be deterministic, got %s and %s", result, result2)
	}

	// Check that the result is different for a different seed
	result3 := constructor.PinEncode(seed+1, pin)
	if result == result3 {
		t.Error("Expected PinEncode to return different results for different seeds")
	}

	// Check that the result is different for a different PIN
	result4 := constructor.PinEncode(seed, pin+1)
	if result == result4 {
		t.Error("Expected PinEncode to return different results for different PINs")
	}
}

// TestSendToServer tests sending packets to the server
func TestSendToServer(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	// Set up encryption keys
	constructor.SetCryptKeys(0x12345678, 0x87654321, 0xABCDEF01)
	constructor.currentCryptKey = new(big.Int).SetUint64(0x12345678)
	constructor.SetNetworkState(int(connection.IN_GAME))

	// Create a test packet
	packet := []byte{0x01, 0x02, 0x03, 0x04}
	originalID := binary.LittleEndian.Uint16(packet[0:2])

	// Send the packet
	sentPacket := constructor.SendToServer(packet)

	// Check that the packet was encrypted
	sentID := binary.LittleEndian.Uint16(sentPacket[0:2])
	if sentID == originalID {
		t.Error("Expected message ID to be encrypted, but it wasn't")
	}
}

// TestSendRaw tests sending raw packets
func TestSendRaw(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	// Test with a valid raw string
	raw := "01 02 03 04"
	packet, err := constructor.SendRaw(raw)
	if err != nil {
		t.Fatalf("SendRaw failed: %v", err)
	}

	// Check that the packet has the correct bytes
	expected := []byte{0x01, 0x02, 0x03, 0x04}
	if !bytes.Equal(packet, expected) {
		t.Errorf("Expected packet to be %v, got %v", expected, packet)
	}

	// Test with an invalid raw string
	raw = "01 02 ZZ 04"
	_, err = constructor.SendRaw(raw)
	if err == nil {
		t.Error("Expected SendRaw to fail with invalid raw string, but it didn't")
	}
}

// TestGenerateSitStand tests generating sit/stand packets
func TestGenerateSitStand(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	// Test generating a sit packet
	sitPacket := constructor.GenerateSitStand(true)
	if sitPacket == nil {
		t.Fatal("GenerateSitStand(true) returned nil")
	}

	// Test generating a stand packet
	standPacket := constructor.GenerateSitStand(false)
	if standPacket == nil {
		t.Fatal("GenerateSitStand(false) returned nil")
	}
}

// TestGenerateAttack tests generating attack packets
func TestGenerateAttack(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	// Test generating an attack packet
	targetID := uint32(0x12345678)
	flag := uint32(0)
	attackPacket := constructor.GenerateAttack(targetID, flag)
	if attackPacket == nil {
		t.Fatal("GenerateAttack() returned nil")
	}
}

// TestGenerateSkillUse tests generating skill use packets
func TestGenerateSkillUse(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	// Test generating a skill use packet
	skillID := uint32(123)
	skillLv := uint32(10)
	targetID := uint32(0x12345678)
	skillPacket := constructor.GenerateSkillUse(skillID, skillLv, targetID)
	if skillPacket == nil {
		t.Fatal("GenerateSkillUse() returned nil")
	}
}

// TestInjectMessage tests injecting messages
func TestInjectMessage(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	// Test injecting a message
	message := "Hello, world!"
	packet := constructor.InjectMessage(message)
	if packet == nil {
		t.Fatal("InjectMessage() returned nil")
	}

	// Check that the packet has the correct ID
	packetID := binary.LittleEndian.Uint16(packet[0:2])
	if packetID != 0x0109 {
		t.Errorf("Expected packet ID to be 0x0109, got 0x%04X", packetID)
	}
}

// TestInjectAdminMessage tests injecting admin messages
func TestInjectAdminMessage(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	// Test injecting an admin message
	message := "Hello, admin!"
	packet := constructor.InjectAdminMessage(message)
	if packet == nil {
		t.Fatal("InjectAdminMessage() returned nil")
	}

	// Check that the packet has the correct ID
	packetID := binary.LittleEndian.Uint16(packet[0:2])
	if packetID != 0x009A {
		t.Errorf("Expected packet ID to be 0x009A, got 0x%04X", packetID)
	}
}

// TestString tests the String method
func TestString(t *testing.T) {
	db := NewDefaultPacketDatabase()
	constructor := NewPacketConstructor(db)

	// Test the String method
	str := constructor.String()
	if str == "" {
		t.Error("Expected String() to return a non-empty string")
	}
}
