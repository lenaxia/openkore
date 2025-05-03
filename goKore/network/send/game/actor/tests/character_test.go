package tests

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
	"github.com/lenaxia/goKore/network/send/game/actor"
)

// TestData represents the structure of our test data JSON
type TestData struct {
	Method  string        `json:"method"`
	Args    []interface{} `json:"args"`
	Packets []struct {
		MessageID string `json:"messageID"`
		Hex       string `json:"hex"`
		Bytes     []byte `json:"bytes"`
	} `json:"packets"`
}

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
	// For testing, we'll just return the expected bytes from our test data
	if packetName == "char_delete" {
		// In a real implementation, this would construct the packet based on the format
		// But for testing, we'll just return a predefined value
		return []byte{0x00, 0x68, 0x41, 0x42, 0x43, 0x44, 0x31, 0x32, 0x33, 0x34, 0x74, 0x65, 0x73, 0x74, 0x40, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, nil
	}
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

func TestCharDelete(t *testing.T) {
	// Load test data
	data, err := os.ReadFile("char_delete.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	var testData TestData
	err = json.Unmarshal(data, &testData)
	if err != nil {
		t.Fatalf("Failed to parse test data: %v", err)
	}

	// Create mock objects
	mockSend := NewMockBaseSend()
	mockHookManager := hooks.NewHookManager()
	mockLogger := NewMockLogger()

	// Create the manager
	manager := actor.NewManager(mockSend, mockHookManager, mockLogger)

	// Register handlers
	manager.RegisterHandlers()

	// Prepare arguments
	charID := []byte("ABCD1234")
	email := []byte("test@example.com")
	args := map[string]interface{}{
		"charID": charID,
		"email":  email,
	}

	// Call the handler
	packet, err := manager.HandleCharDelete(args)
	if err != nil {
		t.Errorf("HandleCharDelete returned an error: %v", err)
	}

	// Check that the packet matches the expected output
	expectedPacket := testData.Packets[0].Bytes
	if !reflect.DeepEqual(packet, expectedPacket) {
		t.Errorf("HandleCharDelete constructed incorrect packet.\nGot: %v\nExpected: %v", packet, expectedPacket)
	}

	// Check that the logger was called
	if len(mockLogger.debugMessages) == 0 || mockLogger.debugMessages[0] != "Handling char_delete packet" {
		t.Errorf("Logger was not called correctly")
	}
}
