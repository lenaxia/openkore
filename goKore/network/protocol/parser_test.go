package protocol

import (
	"encoding/binary"
	"reflect"
	"testing"
)

// TestNewPacketParser tests the creation of a new packet parser
func TestNewPacketParser(t *testing.T) {
	parser := NewPacketParser()
	if parser == nil {
		t.Fatal("NewPacketParser() returned nil")
	}

	// Check that the packet lists are initialized
	if parser.PacketList == nil {
		t.Error("PacketList was not initialized")
	}
	if parser.PacketLUT == nil {
		t.Error("PacketLUT was not initialized")
	}
}

// TestRegisterHandler tests the registration of packet handlers
func TestRegisterHandler(t *testing.T) {
	parser := NewPacketParser()

	// Register the handler for a specific packet ID
	packetID := "0102"
	parser.RegisterHandler(packetID, "test_handler", "v1 C1 v1", []string{"param1", "param2", "param3"}, nil)

	// Verify the handler was registered
	if _, exists := parser.PacketList[packetID]; !exists {
		t.Errorf("Handler was not registered for packet ID %s", packetID)
	}

	// Verify the handler function was stored
	handlerInfo := parser.PacketList[packetID]
	if handlerInfo.Name != "test_handler" {
		t.Errorf("Handler name mismatch: got %s, want test_handler", handlerInfo.Name)
	}

	// Verify the format string was stored
	if handlerInfo.Format != "v1 C1 v1" {
		t.Errorf("Format string mismatch: got %s, want v1 C1 v1", handlerInfo.Format)
	}

	// Verify the parameter names were stored
	expectedParams := []string{"param1", "param2", "param3"}
	if !reflect.DeepEqual(handlerInfo.ParamNames, expectedParams) {
		t.Errorf("Parameter names mismatch: got %v, want %v", handlerInfo.ParamNames, expectedParams)
	}
}

// TestLookupHandlerByName tests looking up a packet ID by handler name
func TestLookupHandlerByName(t *testing.T) {
	parser := NewPacketParser()

	// Register a handler
	packetID := "0102"
	handlerName := "test_handler"
	parser.RegisterHandler(packetID, handlerName, "v1", []string{"param1"}, nil)

	// Look up the packet ID by handler name
	lookupID, exists := parser.LookupPacketID(handlerName)
	if !exists {
		t.Errorf("Failed to look up packet ID for handler %s", handlerName)
	}
	if lookupID != packetID {
		t.Errorf("Packet ID mismatch: got %s, want %s", lookupID, packetID)
	}

	// Test looking up a non-existent handler
	_, exists = parser.LookupPacketID("non_existent_handler")
	if exists {
		t.Error("LookupPacketID returned true for non-existent handler")
	}
}

// TestParsePacket tests parsing a packet into a map of arguments
func TestParsePacket(t *testing.T) {
	parser := NewPacketParser()

	// Register a handler
	packetID := "0102"
	parser.RegisterHandler(packetID, "test_handler", "v1 C1 v1", []string{"param1", "param2", "param3"}, nil)

	// Create a test packet
	packet := make([]byte, 7) // 2 bytes for switch + 5 bytes for data
	packet[0] = 0x02          // Little endian switch
	packet[1] = 0x01

	// Set values for parameters
	binary.LittleEndian.PutUint16(packet[2:4], 12345) // param1 (v1 = 2 bytes)
	packet[4] = 67                                    // param2 (C1 = 1 byte)
	binary.LittleEndian.PutUint16(packet[5:7], 9876)  // param3 (v1 = 2 bytes)

	// Parse the packet
	args, err := parser.Parse(packet)
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	// Verify the parsed arguments
	if args["switch"] != packetID {
		t.Errorf("Switch mismatch: got %s, want %s", args["switch"], packetID)
	}

	if args["param1"] != uint16(12345) {
		t.Errorf("param1 mismatch: got %v, want %v", args["param1"], uint16(12345))
	}

	if args["param2"] != uint8(67) {
		t.Errorf("param2 mismatch: got %v, want %v", args["param2"], uint8(67))
	}

	if args["param3"] != uint16(9876) {
		t.Errorf("param3 mismatch: got %v, want %v", args["param3"], uint16(9876))
	}
}

// TestParseUnknownPacket tests parsing a packet with an unknown ID
func TestParseUnknownPacket(t *testing.T) {
	parser := NewPacketParser()

	// Create a test packet with an unknown ID
	packet := []byte{0xFF, 0xFF, 0x01, 0x02}

	// Parse the packet - should return nil args but no error
	args, err := parser.Parse(packet)
	if err != nil {
		t.Fatalf("Parse() returned error for unknown packet: %v", err)
	}
	if args != nil {
		t.Errorf("Parse() returned non-nil args for unknown packet: %v", args)
	}
}

// TestReconstructPacket tests reconstructing a packet from arguments
func TestReconstructPacket(t *testing.T) {
	parser := NewPacketParser()

	// Register a handler
	packetID := "0102"
	parser.RegisterHandler(packetID, "test_handler", "v1 C1 v1", []string{"param1", "param2", "param3"}, nil)

	// Create arguments
	args := map[string]interface{}{
		"switch": packetID,
		"param1": uint16(12345),
		"param2": uint8(67),
		"param3": uint16(9876),
	}

	// Reconstruct the packet
	packet, err := parser.Reconstruct(args)
	if err != nil {
		t.Fatalf("Reconstruct() returned error: %v", err)
	}

	// Verify the packet length
	expectedLength := 7 // 2 bytes for switch + 5 bytes for data
	if len(packet) != expectedLength {
		t.Errorf("Packet length mismatch: got %d, want %d", len(packet), expectedLength)
	}

	// Verify the packet switch
	if packet[0] != 0x02 || packet[1] != 0x01 {
		t.Errorf("Packet switch mismatch: got %02X%02X, want %s", packet[1], packet[0], packetID)
	}

	// Verify the parameter values
	if binary.LittleEndian.Uint16(packet[2:4]) != 12345 {
		t.Errorf("param1 mismatch: got %d, want %d", binary.LittleEndian.Uint16(packet[2:4]), 12345)
	}

	if packet[4] != 67 {
		t.Errorf("param2 mismatch: got %d, want %d", packet[4], 67)
	}

	if binary.LittleEndian.Uint16(packet[5:7]) != 9876 {
		t.Errorf("param3 mismatch: got %d, want %d", binary.LittleEndian.Uint16(packet[5:7]), 9876)
	}
}

// TestReconstructByHandlerName tests reconstructing a packet using a handler name
func TestReconstructByHandlerName(t *testing.T) {
	parser := NewPacketParser()

	// Register a handler
	packetID := "0102"
	handlerName := "test_handler"
	parser.RegisterHandler(packetID, handlerName, "v1 C1 v1", []string{"param1", "param2", "param3"}, nil)

	// Create arguments using handler name instead of switch
	args := map[string]interface{}{
		"switch": handlerName,
		"param1": uint16(12345),
		"param2": uint8(67),
		"param3": uint16(9876),
	}

	// Reconstruct the packet
	packet, err := parser.Reconstruct(args)
	if err != nil {
		t.Fatalf("Reconstruct() returned error: %v", err)
	}

	// Verify the packet switch
	if packet[0] != 0x02 || packet[1] != 0x01 {
		t.Errorf("Packet switch mismatch: got %02X%02X, want %s", packet[1], packet[0], packetID)
	}
}

// TestProcessPackets tests processing multiple packets from a tokenizer
func TestProcessPackets(t *testing.T) {
	parser := NewPacketParser()
	tokenizer := NewTokenizer(map[string]PacketDef{
		"0102": {Length: 7, HasLength: false},
		"0304": {Length: 5, HasLength: false},
	})

	// Register handlers
	handlerCalled1 := false
	handler1 := func(packet []byte) error {
		handlerCalled1 = true
		return nil
	}

	handlerCalled2 := false
	handler2 := func(packet []byte) error {
		handlerCalled2 = true
		return nil
	}

	parser.RegisterHandler("0102", "test_handler1", "v1 C1 v1", []string{"param1", "param2", "param3"}, handler1)
	parser.RegisterHandler("0304", "test_handler2", "v1 C1", []string{"param1", "param2"}, handler2)

	// Create test packets
	packet1 := []byte{0x02, 0x01, 0x39, 0x30, 0x43, 0xD4, 0x26}
	packet2 := []byte{0x04, 0x03, 0x01, 0x02, 0x03}

	// Add packets to tokenizer
	tokenizer.Add(append(packet1, packet2...))

	// Process packets
	err := parser.Process(tokenizer)
	if err != nil {
		t.Fatalf("Process() returned error: %v", err)
	}

	// Verify handlers were called
	if !handlerCalled1 {
		t.Error("Handler 1 was not called")
	}
	if !handlerCalled2 {
		t.Error("Handler 2 was not called")
	}
}
