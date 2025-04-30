package types

import (
	"errors"
	"testing"
)

// Happy path - create a packet construction with all fields
func TestPacketConstruction_HappyPath(t *testing.T) {
	// Create a new PacketConstruction
	pc := PacketConstruction{
		ID:         "0064",
		Name:       "login_request",
		Format:     "v a24 a24 C",
		FieldNames: []string{"version", "username", "password", "clienttype"},
	}

	// Check that the fields were set correctly
	if pc.ID != "0064" {
		t.Errorf("Expected ID to be 0064, got %s", pc.ID)
	}

	if pc.Name != "login_request" {
		t.Errorf("Expected Name to be login_request, got %s", pc.Name)
	}

	if pc.Format != "v a24 a24 C" {
		t.Errorf("Expected Format to be v a24 a24 C, got %s", pc.Format)
	}

	if len(pc.FieldNames) != 4 {
		t.Errorf("Expected FieldNames to have 4 elements, got %d", len(pc.FieldNames))
	}

	expectedFieldNames := []string{"version", "username", "password", "clienttype"}
	for i, fieldName := range pc.FieldNames {
		if fieldName != expectedFieldNames[i] {
			t.Errorf("Expected FieldNames[%d] to be %s, got %s", i, expectedFieldNames[i], fieldName)
		}
	}
}

// Edge case - create a packet construction with empty format and field names
func TestPacketConstruction_EmptyFields(t *testing.T) {
	// Create a new PacketConstruction with empty format and field names
	pc := PacketConstruction{
		ID:         "0001",
		Name:       "empty_packet",
		Format:     "",
		FieldNames: []string{},
	}

	// Check that the fields were set correctly
	if pc.ID != "0001" {
		t.Errorf("Expected ID to be 0001, got %s", pc.ID)
	}

	if pc.Name != "empty_packet" {
		t.Errorf("Expected Name to be empty_packet, got %s", pc.Name)
	}

	if pc.Format != "" {
		t.Errorf("Expected Format to be empty, got %s", pc.Format)
	}

	if len(pc.FieldNames) != 0 {
		t.Errorf("Expected FieldNames to be empty, got %d elements", len(pc.FieldNames))
	}
}

// Edge case - create a packet construction with mismatched format and field names
func TestPacketConstruction_MismatchedFields(t *testing.T) {
	// Create a new PacketConstruction with mismatched format and field names
	// Format "v C" suggests 2 fields, but we provide 3 field names
	pc := PacketConstruction{
		ID:         "0002",
		Name:       "mismatched_packet",
		Format:     "v C",
		FieldNames: []string{"field1", "field2", "field3"},
	}

	// This is just a structural test - in a real implementation, we might want to
	// validate that the number of format specifiers matches the number of field names
	if len(pc.FieldNames) != 3 {
		t.Errorf("Expected FieldNames to have 3 elements, got %d", len(pc.FieldNames))
	}
}

// MockSend is a mock implementation of the Send interface for testing
type MockSend struct {
	handlers     map[string]SendHandler
	lastPacket   []byte
	lastPacketID string
	shouldFail   bool
}

// NewMockSend creates a new MockSend instance
func NewMockSend() *MockSend {
	return &MockSend{
		handlers:   make(map[string]SendHandler),
		shouldFail: false,
	}
}

// RegisterHandler registers a handler for a specific packet
func (ms *MockSend) RegisterHandler(packetName string, handler SendHandler) {
	ms.handlers[packetName] = handler
}

// ConstructPacket constructs a packet from a packet name and arguments
func (ms *MockSend) ConstructPacket(packetName string, args map[string]interface{}) ([]byte, error) {
	if ms.shouldFail {
		return nil, errors.New("mock failure")
	}

	handler, exists := ms.handlers[packetName]
	if !exists {
		return nil, errors.New("unknown packet: " + packetName)
	}

	ms.lastPacketID = packetName
	packet, err := handler(args)
	ms.lastPacket = packet
	return packet, err
}

// SendPacket constructs and sends a packet
func (ms *MockSend) SendPacket(packetName string, args map[string]interface{}) error {
	if ms.shouldFail {
		return errors.New("mock failure")
	}

	packet, err := ms.ConstructPacket(packetName, args)
	if err != nil {
		return err
	}
	return ms.SendToServer(packet)
}

// SendToServer sends a raw packet to the server
func (ms *MockSend) SendToServer(packet []byte) error {
	if ms.shouldFail {
		return errors.New("mock failure")
	}

	ms.lastPacket = packet
	return nil
}

// Configure configures the send component with server-specific packet constructions
func (ms *MockSend) Configure(serverType string, packetConstructions map[string]PacketConstruction) error {
	if ms.shouldFail {
		return errors.New("mock failure")
	}
	return nil
}

// SetShouldFail sets whether the mock should fail
func (ms *MockSend) SetShouldFail(shouldFail bool) {
	ms.shouldFail = shouldFail
}

// Happy path - register and use a handler
func TestMockSend_HappyPath(t *testing.T) {
	// Create a new MockSend
	ms := NewMockSend()

	// Register a handler
	ms.RegisterHandler("test_packet", func(args map[string]interface{}) ([]byte, error) {
		return []byte{0x01, 0x02, 0x03}, nil
	})

	// Construct a packet
	packet, err := ms.ConstructPacket("test_packet", nil)

	// Check that there was no error
	if err != nil {
		t.Fatalf("ConstructPacket() returned error: %v", err)
	}

	// Check that the packet is correct
	expectedPacket := []byte{0x01, 0x02, 0x03}
	if len(packet) != len(expectedPacket) {
		t.Fatalf("Expected packet length to be %d, got %d", len(expectedPacket), len(packet))
	}

	for i, b := range packet {
		if b != expectedPacket[i] {
			t.Errorf("Expected packet[%d] to be %d, got %d", i, expectedPacket[i], b)
		}
	}

	// Check that the last packet ID is correct
	if ms.lastPacketID != "test_packet" {
		t.Errorf("Expected lastPacketID to be test_packet, got %s", ms.lastPacketID)
	}

	// Send a packet
	err = ms.SendPacket("test_packet", nil)

	// Check that there was no error
	if err != nil {
		t.Fatalf("SendPacket() returned error: %v", err)
	}

	// Check that the last packet is correct
	if len(ms.lastPacket) != len(expectedPacket) {
		t.Fatalf("Expected lastPacket length to be %d, got %d", len(expectedPacket), len(ms.lastPacket))
	}

	for i, b := range ms.lastPacket {
		if b != expectedPacket[i] {
			t.Errorf("Expected lastPacket[%d] to be %d, got %d", i, expectedPacket[i], b)
		}
	}
}

// Unhappy path - unknown packet
func TestMockSend_UnknownPacket(t *testing.T) {
	// Create a new MockSend
	ms := NewMockSend()

	// Try to construct an unknown packet
	packet, err := ms.ConstructPacket("unknown_packet", nil)

	// Check that there was an error
	if err == nil {
		t.Fatal("ConstructPacket() did not return an error for unknown packet")
	}

	// Check that the packet is nil
	if packet != nil {
		t.Errorf("Expected packet to be nil, got %v", packet)
	}

	// Try to send an unknown packet
	err = ms.SendPacket("unknown_packet", nil)

	// Check that there was an error
	if err == nil {
		t.Fatal("SendPacket() did not return an error for unknown packet")
	}
}

// Unhappy path - mock failure
func TestMockSend_Failure(t *testing.T) {
	// Create a new MockSend
	ms := NewMockSend()

	// Register a handler
	ms.RegisterHandler("test_packet", func(args map[string]interface{}) ([]byte, error) {
		return []byte{0x01, 0x02, 0x03}, nil
	})

	// Set the mock to fail
	ms.SetShouldFail(true)

	// Try to construct a packet
	packet, err := ms.ConstructPacket("test_packet", nil)

	// Check that there was an error
	if err == nil {
		t.Fatal("ConstructPacket() did not return an error when mock is set to fail")
	}

	// Check that the packet is nil
	if packet != nil {
		t.Errorf("Expected packet to be nil, got %v", packet)
	}

	// Try to send a packet
	err = ms.SendPacket("test_packet", nil)

	// Check that there was an error
	if err == nil {
		t.Fatal("SendPacket() did not return an error when mock is set to fail")
	}

	// Try to send a raw packet
	err = ms.SendToServer([]byte{0x01, 0x02, 0x03})

	// Check that there was an error
	if err == nil {
		t.Fatal("SendToServer() did not return an error when mock is set to fail")
	}

	// Try to configure the mock
	err = ms.Configure("test", make(map[string]PacketConstruction))

	// Check that there was an error
	if err == nil {
		t.Fatal("Configure() did not return an error when mock is set to fail")
	}
}

// Test handler that returns an error
func TestMockSend_HandlerError(t *testing.T) {
	// Create a new MockSend
	ms := NewMockSend()

	// Register a handler that returns an error
	ms.RegisterHandler("error_packet", func(args map[string]interface{}) ([]byte, error) {
		return nil, errors.New("handler error")
	})

	// Try to construct a packet
	packet, err := ms.ConstructPacket("error_packet", nil)

	// Check that there was an error
	if err == nil {
		t.Fatal("ConstructPacket() did not return an error when handler returns an error")
	}

	// Check that the packet is nil
	if packet != nil {
		t.Errorf("Expected packet to be nil, got %v", packet)
	}

	// Try to send a packet
	err = ms.SendPacket("error_packet", nil)

	// Check that there was an error
	if err == nil {
		t.Fatal("SendPacket() did not return an error when handler returns an error")
	}
}

// Test handler with arguments
func TestMockSend_HandlerWithArgs(t *testing.T) {
	// Create a new MockSend
	ms := NewMockSend()

	// Register a handler that uses arguments
	ms.RegisterHandler("args_packet", func(args map[string]interface{}) ([]byte, error) {
		// Check that the arguments are correct
		if args["test"] != "value" {
			return nil, errors.New("incorrect arguments")
		}
		return []byte{0x01, 0x02, 0x03}, nil
	})

	// Create arguments
	args := map[string]interface{}{
		"test": "value",
	}

	// Construct a packet with arguments
	packet, err := ms.ConstructPacket("args_packet", args)

	// Check that there was no error
	if err != nil {
		t.Fatalf("ConstructPacket() returned error: %v", err)
	}

	// Check that the packet is correct
	expectedPacket := []byte{0x01, 0x02, 0x03}
	if len(packet) != len(expectedPacket) {
		t.Fatalf("Expected packet length to be %d, got %d", len(expectedPacket), len(packet))
	}

	// Try with incorrect arguments
	args = map[string]interface{}{
		"test": "wrong",
	}

	// Construct a packet with incorrect arguments
	packet, err = ms.ConstructPacket("args_packet", args)

	// Check that there was an error
	if err == nil {
		t.Fatal("ConstructPacket() did not return an error with incorrect arguments")
	}

	// Check that the packet is nil
	if packet != nil {
		t.Errorf("Expected packet to be nil, got %v", packet)
	}
}
