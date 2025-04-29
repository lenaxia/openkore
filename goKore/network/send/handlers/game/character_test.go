package game

import (
	"reflect"
	"testing"

	"github.com/lenaxia/goKore/network/send/core"
)

// MockBaseSend is a mock implementation of the core.BaseSend for testing
type MockBaseSend struct {
	handlers map[string]core.SendHandler
}

// NewMockBaseSend creates a new MockBaseSend instance
func NewMockBaseSend() *MockBaseSend {
	return &MockBaseSend{
		handlers: make(map[string]core.SendHandler),
	}
}

// RegisterHandler registers a handler for a packet name
func (m *MockBaseSend) RegisterHandler(packetName string, handler core.SendHandler) {
	m.handlers[packetName] = handler
}

// GetHandler returns a handler for a packet name
func (m *MockBaseSend) GetHandler(packetName string) (core.SendHandler, bool) {
	handler, exists := m.handlers[packetName]
	return handler, exists
}

// TestRegisterCharacterHandlers tests the RegisterCharacterHandlers function
func TestRegisterCharacterHandlers(t *testing.T) {
	mockSend := NewMockBaseSend()
	RegisterCharacterHandlers(mockSend)

	// Check that the restart handler was registered
	_, exists := mockSend.GetHandler("restart")
	if !exists {
		t.Error("restart handler was not registered")
	}
}

// TestRestartHandler tests the restart handler
func TestRestartHandler(t *testing.T) {
	mockSend := NewMockBaseSend()
	RegisterCharacterHandlers(mockSend)

	// Get the restart handler
	handler, exists := mockSend.GetHandler("restart")
	if !exists {
		t.Fatal("restart handler was not registered")
	}

	// Test with type 0 (respawn)
	args := map[string]interface{}{
		"type": uint8(0),
	}
	packet, err := handler(args)
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	// Check the packet
	expected := []byte{0xb2, 0x00, 0x00}
	if !reflect.DeepEqual(packet, expected) {
		t.Errorf("Expected packet %v, got %v", expected, packet)
	}

	// Test with type 1 (quit to char select)
	args = map[string]interface{}{
		"type": uint8(1),
	}
	packet, err = handler(args)
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	// Check the packet
	expected = []byte{0xb2, 0x00, 0x01}
	if !reflect.DeepEqual(packet, expected) {
		t.Errorf("Expected packet %v, got %v", expected, packet)
	}

	// Test with missing type (should default to 0)
	args = map[string]interface{}{}
	packet, err = handler(args)
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	// Check the packet
	expected = []byte{0xb2, 0x00, 0x00}
	if !reflect.DeepEqual(packet, expected) {
		t.Errorf("Expected packet %v, got %v", expected, packet)
	}
}
