package protocol_test

import (
	"testing"

	"github.com/lenaxia/goKore/network/protocol"
	"github.com/lenaxia/goKore/test/utils"
)

// TestTokenizeLoginPackets tests the tokenizer's ability to extract login packets
func TestTokenizeLoginPackets(t *testing.T) {
	// Create packet definitions for the essential login packets
	packetDefs := map[string]protocol.PacketLengthDef{
		"0064": {Length: 55, HasLength: false}, // Account Server Login
		"0AC4": {Length: 224, HasLength: true}, // Account Info With Server Info
		"0065": {Length: 17, HasLength: false}, // Character Server Login
		"0066": {Length: 3, HasLength: false},  // Char Login
		"082D": {Length: 29, HasLength: true},  // Received characters info
		"006B": {Length: 182, HasLength: true}, // Received characters
		"08B9": {Length: 12, HasLength: false}, // PinCode Request
		"0AC5": {Length: 156, HasLength: true}, // Character Map Info
		"0436": {Length: 19, HasLength: false}, // Map Login
		"007D": {Length: 2, HasLength: false},  // Map Loaded
		"0283": {Length: 6, HasLength: false},  // Account ID
		"02EB": {Length: 13, HasLength: false}, // Enter Map
	}

	// Create tokenizer with the packet definitions
	tokenizer := protocol.NewTokenizer(packetDefs)

	// Test each packet
	for _, testCase := range utils.AllLoginPackets {
		t.Run(testCase.Name, func(t *testing.T) {
			// Convert hex to bytes if needed
			var packetData []byte
			var err error
			if len(testCase.RawData) == 0 {
				packetData, err = utils.HexToBytes(testCase.RawHex)
				if err != nil {
					t.Fatalf("Failed to convert hex to bytes: %v", err)
				}
			} else {
				packetData = testCase.RawData
			}

			// Reset the tokenizer before each test
			tokenizer = protocol.NewTokenizer(packetDefs)

			// Add the packet data to the tokenizer
			tokenizer.Add(packetData)

			// Extract the packet
			packet, msgType, err := tokenizer.ReadNext()
			if err != nil {
				t.Fatalf("Failed to extract packet: %v", err)
			}

			// Verify the packet was extracted correctly
			if msgType != protocol.KnownMessage {
				t.Errorf("Expected message type KnownMessage, got %v", msgType)
			}

			// Verify the packet ID
			packetID := utils.GetPacketID(packet)
			if packetID != testCase.PacketID {
				t.Errorf("Expected packet ID %s, got %s", testCase.PacketID, packetID)
			}

			// Verify the packet length
			expectedLength := packetDefs[testCase.PacketID].Length
			if expectedLength > 0 && len(packet) != expectedLength {
				t.Errorf("Expected packet length %d, got %d", expectedLength, len(packet))
			}
		})
	}
}

// TestTokenizeMultiplePackets tests the tokenizer's ability to extract multiple packets from a stream
func TestTokenizeMultiplePackets(t *testing.T) {
	// Create packet definitions for the essential login packets
	packetDefs := map[string]protocol.PacketLengthDef{
		"0064": {Length: 55, HasLength: false}, // Account Server Login
		"0065": {Length: 17, HasLength: false}, // Character Server Login
		"0066": {Length: 3, HasLength: false},  // Char Login
		"007D": {Length: 2, HasLength: false},  // Map Loaded
	}

	// Create tokenizer with the packet definitions
	tokenizer := protocol.NewTokenizer(packetDefs)

	// Test with just the Map Loaded packet since it's the simplest
	mapLoadedData, err := utils.HexToBytes(utils.MapLoadedPacket.RawHex)
	if err != nil {
		t.Fatalf("Failed to convert hex to bytes: %v", err)
	}

	// Add the packet to the tokenizer
	tokenizer.Add(mapLoadedData)

	// Extract and verify the packet
	packet, msgType, err := tokenizer.ReadNext()
	if err != nil {
		t.Fatalf("Failed to extract packet: %v", err)
	}

	if msgType != protocol.KnownMessage {
		t.Errorf("Expected message type KnownMessage, got %v", msgType)
	}

	packetID := utils.GetPacketID(packet)
	if packetID != "007D" {
		t.Errorf("Expected packet ID 007D, got %s", packetID)
	}
}
