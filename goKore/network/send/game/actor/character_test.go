package actor

import (
	"reflect"
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
)

// MockBaseSend implements the core.BaseSend interface for testing
type MockBaseSend struct {
	registeredHandlers map[string]core.SendHandler
	constructedPackets map[string][]byte
	serverType         string
}

func NewMockBaseSend() *MockBaseSend {
	return &MockBaseSend{
		registeredHandlers: make(map[string]core.SendHandler),
		constructedPackets: make(map[string][]byte),
		serverType:         "ServerType0",
	}
}

func (m *MockBaseSend) RegisterHandler(packetName string, handler core.SendHandler) {
	m.registeredHandlers[packetName] = handler
}

func (m *MockBaseSend) ConstructPacket(packetName string, args map[string]interface{}) ([]byte, error) {
	// For char_delete, we expect charID and email
	if packetName == "char_delete" {
		// Create a simple mock packet for testing
		// In a real implementation, this would construct the packet based on the format
		packet := []byte{0x00, 0x68} // Packet ID 0068

		// Add charID (4 bytes)
		if charID, ok := args["charID"].([]byte); ok && len(charID) == 4 {
			packet = append(packet, charID...)
		} else {
			// Default test value
			packet = append(packet, []byte{0x01, 0x02, 0x03, 0x04}...)
		}

		// Add email (40 bytes)
		if email, ok := args["email"].([]byte); ok && len(email) <= 40 {
			// Pad email to 40 bytes
			paddedEmail := make([]byte, 40)
			copy(paddedEmail, email)
			packet = append(packet, paddedEmail...)
		} else {
			// Default test value
			defaultEmail := make([]byte, 40)
			copy(defaultEmail, []byte("test@example.com"))
			packet = append(packet, defaultEmail...)
		}

		m.constructedPackets[packetName] = packet
		return packet, nil
	}

	// Default empty packet
	return []byte{}, nil
}

func (m *MockBaseSend) GetServerType() string {
	return m.serverType
}

// MockLogger implements the core.Logger interface for testing
type MockLogger struct {
	debugMessages []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		debugMessages: make([]string, 0),
	}
}

func (m *MockLogger) Debug(format string, args ...interface{}) {
	m.debugMessages = append(m.debugMessages, format)
}

func (m *MockLogger) Info(format string, args ...interface{})    {}
func (m *MockLogger) Warning(format string, args ...interface{}) {}
func (m *MockLogger) Error(format string, args ...interface{})   {}

func TestHandleCharDelete(t *testing.T) {
	// Create mock objects
	mockSend := NewMockBaseSend()
	mockHookManager := hooks.NewHookManager()
	mockLogger := NewMockLogger()

	// Create the manager
	manager := NewManager(mockSend, mockHookManager, mockLogger)

	// Test data
	charID := []byte{0x01, 0x02, 0x03, 0x04}
	email := []byte("test@example.com")

	// Call the handler
	args := map[string]interface{}{
		"charID": charID,
		"email":  email,
	}
	packet, err := manager.HandleCharDelete(args)

	// Check for errors
	if err != nil {
		t.Errorf("HandleCharDelete returned an error: %v", err)
	}

	// Check that the packet was constructed correctly
	expectedPacket := []byte{0x00, 0x68, 0x01, 0x02, 0x03, 0x04}
	emailBytes := make([]byte, 40)
	copy(emailBytes, email)
	expectedPacket = append(expectedPacket, emailBytes...)

	if !reflect.DeepEqual(packet, expectedPacket) {
		t.Errorf("HandleCharDelete constructed incorrect packet.\nGot: %v\nExpected: %v", packet, expectedPacket)
	}

	// Check that the logger was called
	if len(mockLogger.debugMessages) == 0 || mockLogger.debugMessages[0] != "Handling char_delete packet" {
		t.Errorf("Logger was not called correctly")
	}

	// Check that the hook was called
	hookCalled := false
	mockHookManager.RegisterHook("send/actor/char_delete", func(data interface{}) {
		hookCalled = true
		// Check that the hook arguments match the original arguments
		if hookArgs, ok := data.(map[string]interface{}); ok {
			if !reflect.DeepEqual(hookArgs, args) {
				t.Errorf("Hook arguments do not match original arguments.\nGot: %v\nExpected: %v", hookArgs, args)
			}
		} else {
			t.Errorf("Hook data is not of expected type map[string]interface{}")
		}
	})

	// Call the handler again to trigger the hook
	manager.HandleCharDelete(args)

	if !hookCalled {
		t.Errorf("Hook was not called")
	}
}
