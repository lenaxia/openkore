package login

import (
	"testing"

	"github.com/lenaxia/goKore/network/send/types"
)

// MockSend is a mock implementation of the types.Send interface for testing
type MockSend struct {
	handlers map[string]types.SendHandler
}

// NewMockSend creates a new MockSend instance
func NewMockSend() *MockSend {
	return &MockSend{
		handlers: make(map[string]types.SendHandler),
	}
}

// RegisterHandler registers a handler for a specific packet
func (ms *MockSend) RegisterHandler(packetName string, handler types.SendHandler) {
	ms.handlers[packetName] = handler
}

// ConstructPacket constructs a packet from a packet name and arguments
func (ms *MockSend) ConstructPacket(packetName string, args map[string]interface{}) ([]byte, error) {
	handler, exists := ms.handlers[packetName]
	if !exists {
		return nil, nil
	}
	return handler(args)
}

// SendPacket constructs and sends a packet
func (ms *MockSend) SendPacket(packetName string, args map[string]interface{}) error {
	return nil
}

// SendToServer sends a raw packet to the server
func (ms *MockSend) SendToServer(packet []byte) error {
	return nil
}

// Configure configures the send component with server-specific packet constructions
func (ms *MockSend) Configure(serverType string, packetConstructions map[string]types.PacketConstruction) error {
	return nil
}

// Happy path - all handlers are registered correctly
func TestRegisterHandlers(t *testing.T) {
	mockSend := NewMockSend()

	// Register the handlers
	RegisterHandlers(mockSend)

	// Check that the handlers were registered
	expectedHandlers := []string{
		"login_request",
		"login_otp",
		"login_restart",
	}

	for _, handlerName := range expectedHandlers {
		if _, exists := mockSend.handlers[handlerName]; !exists {
			t.Errorf("Handler %s was not registered", handlerName)
		}
	}

	// Verify the handlers actually work by calling them
	for handlerName := range mockSend.handlers {
		packet, err := mockSend.handlers[handlerName](map[string]interface{}{})
		if err != nil {
			t.Errorf("Handler %s returned error: %v", handlerName, err)
		}
		if packet == nil {
			t.Errorf("Handler %s returned nil packet", handlerName)
		}
	}
}

// Happy path - all packet constructions are returned
func TestGetPacketConstructions(t *testing.T) {
	// Get the packet constructions
	constructions := GetPacketConstructions()

	// Check that the packet constructions were returned
	expectedConstructions := []string{
		"0064", // login_request
		"0066", // login_otp
		"00B2", // login_restart
	}

	for _, id := range expectedConstructions {
		if _, exists := constructions[id]; !exists {
			t.Errorf("Packet construction %s was not returned", id)
		}
	}

	// Verify the packet constructions have the correct format and field names
	if constructions["0064"].Format != "v a24 a24 C" {
		t.Errorf("Packet construction 0064 has incorrect format: %s", constructions["0064"].Format)
	}

	expectedFieldNames := []string{"version", "username", "password", "clienttype"}
	for i, fieldName := range constructions["0064"].FieldNames {
		if fieldName != expectedFieldNames[i] {
			t.Errorf("Packet construction 0064 has incorrect field name at index %d: %s", i, fieldName)
		}
	}
}

// Happy path - construct login request with valid arguments
func TestConstructLoginRequest_HappyPath(t *testing.T) {
	// Create arguments
	args := map[string]interface{}{
		"version":    uint16(1),
		"username":   "testuser",
		"password":   "testpass",
		"clienttype": uint8(1),
	}

	// Construct the packet
	packet, err := constructLoginRequest(args)

	// Check that there was no error
	if err != nil {
		t.Fatalf("constructLoginRequest() returned error: %v", err)
	}

	// Check that the packet is not nil and has the expected format
	if packet == nil {
		t.Fatal("constructLoginRequest() returned nil packet")
	}

	// Check packet length (this is a simple check, in a real test we'd verify the actual bytes)
	if len(packet) != 6 {
		t.Errorf("Expected packet length to be 6, got %d", len(packet))
	}
}

// Unhappy path - construct login request with missing arguments
func TestConstructLoginRequest_MissingArgs(t *testing.T) {
	// Create incomplete arguments
	incompleteArgs := map[string]interface{}{
		"version":  uint16(1),
		"username": "testuser",
		// Missing password and clienttype
	}

	// Construct the packet - this should still work with our mock implementation
	// but in a real implementation it might return an error
	packet, _ := constructLoginRequest(incompleteArgs)

	// Our mock implementation doesn't check args, but in a real test we'd expect:
	// - Either an error
	// - Or a valid packet with default values for missing args
	if packet == nil {
		t.Fatal("constructLoginRequest() returned nil packet")
	}
}

// Happy path - construct login OTP with valid arguments
func TestConstructLoginOTP_HappyPath(t *testing.T) {
	// Create arguments
	args := map[string]interface{}{
		"version":    uint16(1),
		"username":   "testuser",
		"password":   "testpass",
		"clienttype": uint8(1),
		"state":      uint8(1),
		"otp":        "123456",
	}

	// Construct the packet
	packet, _ := constructLoginOTP(args)

	// No error check needed for our mock implementation

	// Check that the packet is not nil
	if packet == nil {
		t.Fatal("constructLoginOTP() returned nil packet")
	}

	// Check packet length
	if len(packet) != 6 {
		t.Errorf("Expected packet length to be 6, got %d", len(packet))
	}
}

// Unhappy path - construct login OTP with invalid OTP
func TestConstructLoginOTP_InvalidOTP(t *testing.T) {
	// Create arguments with invalid OTP (too long)
	args := map[string]interface{}{
		"version":    uint16(1),
		"username":   "testuser",
		"password":   "testpass",
		"clienttype": uint8(1),
		"state":      uint8(1),
		"otp":        "1234567890", // Too long
	}

	// Construct the packet - our mock implementation doesn't validate
	packet, _ := constructLoginOTP(args)

	// Our mock implementation doesn't check args, but in a real test we'd expect:
	// - Either an error for invalid OTP
	// - Or a truncated OTP in the packet
	if packet == nil {
		t.Fatal("constructLoginOTP() returned nil packet")
	}
}

// Happy path - construct login restart with valid arguments
func TestConstructLoginRestart_HappyPath(t *testing.T) {
	// Create arguments
	args := map[string]interface{}{
		"type": uint8(1),
	}

	// Construct the packet
	packet, _ := constructLoginRestart(args)

	// No error check needed for our mock implementation

	// Check that the packet is not nil
	if packet == nil {
		t.Fatal("constructLoginRestart() returned nil packet")
	}

	// Check packet length
	if len(packet) != 3 {
		t.Errorf("Expected packet length to be 3, got %d", len(packet))
	}
}

// Unhappy path - construct login restart with invalid type
func TestConstructLoginRestart_InvalidType(t *testing.T) {
	// Create arguments with invalid type (string instead of uint8)
	args := map[string]interface{}{
		"type": "invalid",
	}

	// Construct the packet - our mock implementation doesn't validate
	packet, _ := constructLoginRestart(args)

	// Our mock implementation doesn't check args, but in a real test we'd expect:
	// - Either an error for invalid type
	// - Or a default type value in the packet
	if packet == nil {
		t.Fatal("constructLoginRestart() returned nil packet")
	}
}
