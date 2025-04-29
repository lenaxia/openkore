package marriage

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
)

// MockSend is a mock implementation of the core.Send interface for testing
type MockSend struct {
	packetIDs      map[string]string
	reconstructed  []byte
	sent           []byte
	time           uint32
	lastPacketName string
	lastArgs       map[string]interface{}
}

// NewMockSend creates a new MockSend instance with default values
func NewMockSend() *MockSend {
	return &MockSend{
		packetIDs: map[string]string{
			"adopt_request":       "01F7",
			"adopt_reply_request": "01F9",
		},
		time:     12345,
		lastArgs: make(map[string]interface{}),
	}
}

// SendToServer mocks sending a packet to the server
func (ms *MockSend) SendToServer(msg []byte) error {
	ms.sent = msg
	return nil
}

// EncryptMessageID mocks encrypting a message ID
func (ms *MockSend) EncryptMessageID(msg *[]byte) error {
	return nil
}

// CryptKeys mocks setting encryption keys
func (ms *MockSend) CryptKeys(key1, key2, key3 uint32) {}

// PinEncode mocks encoding a PIN
func (ms *MockSend) PinEncode(seed, pin int) string {
	return ""
}

// InjectMessage mocks injecting a message
func (ms *MockSend) InjectMessage(message string) error {
	return nil
}

// InjectAdminMessage mocks injecting an admin message
func (ms *MockSend) InjectAdminMessage(message string) error {
	return nil
}

// SendRaw mocks sending a raw packet
func (ms *MockSend) SendRaw(raw string) error {
	return nil
}

// Reconstruct mocks reconstructing a packet
func (ms *MockSend) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the last packet name and arguments for testing
	for name, id := range ms.packetIDs {
		if id == packetID {
			ms.lastPacketName = name
			break
		}
	}

	// Store the arguments for testing
	ms.lastArgs = args

	// Simple mock implementation that just returns the packet ID as bytes
	ms.reconstructed = []byte{0x00, 0x00}
	return ms.reconstructed, nil
}

// GetPacketID mocks getting a packet ID by name
func (ms *MockSend) GetPacketID(name string) (string, bool) {
	id, ok := ms.packetIDs[name]
	if ok {
		ms.lastPacketName = name
	}
	return id, ok
}

// RegisterPacketHandler mocks registering a packet handler
func (ms *MockSend) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
}

// RegisterHook mocks registering a hook
func (ms *MockSend) RegisterHook(hookName string, callback hooks.HookCallback) {}

// SetConnection mocks setting a connection
func (ms *MockSend) SetConnection(conn interface{}) {}

// GetConnection mocks getting a connection
func (ms *MockSend) GetConnection() interface{} {
	return nil
}

// GetTime mocks getting the current time
func (ms *MockSend) GetTime() uint32 {
	return ms.time
}

// LastPacketID returns the name of the last packet that was requested
func (ms *MockSend) LastPacketID() (string, bool) {
	if ms.lastPacketName == "" {
		return "", false
	}
	return ms.lastPacketName, true
}

// LastArgs returns the arguments of the last packet that was reconstructed
func (ms *MockSend) LastArgs() map[string]interface{} {
	return ms.lastArgs
}

// TestNewMarriageManager tests the NewMarriageManager function
func TestNewMarriageManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	marriageManager := NewMarriageManager(mockSend)

	if marriageManager == nil {
		t.Fatal("NewMarriageManager() returned nil")
	}

	if marriageManager.baseSend == nil {
		t.Error("marriageManager.baseSend was not set correctly")
	}
}

// TestSendAdoptRequest tests the SendAdoptRequest function
func TestSendAdoptRequest(t *testing.T) {
	mockSend := NewMockSend()
	marriageManager := NewMarriageManager(mockSend)

	// Test sending an adopt request
	playerID := uint32(12345)
	err := marriageManager.SendAdoptRequest(playerID)
	if err != nil {
		t.Fatalf("SendAdoptRequest() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "adopt_request" {
		t.Errorf("Expected packet ID 'adopt_request', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != playerID {
		t.Errorf("Expected ID=%d, got %v", playerID, mockSend.LastArgs()["ID"])
	}
}

// TestSendAdoptReply tests the SendAdoptReply function
func TestSendAdoptReply(t *testing.T) {
	mockSend := NewMockSend()
	marriageManager := NewMarriageManager(mockSend)

	// Test sending an adopt reply
	parentID1 := uint32(12345)
	parentID2 := uint32(67890)
	result := uint8(1) // 1 = accept, 0 = reject
	err := marriageManager.SendAdoptReply(parentID1, parentID2, result)
	if err != nil {
		t.Fatalf("SendAdoptReply() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "adopt_reply_request" {
		t.Errorf("Expected packet ID 'adopt_reply_request', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["parentID1"].(uint32); !ok || idVal != parentID1 {
		t.Errorf("Expected parentID1=%d, got %v", parentID1, mockSend.LastArgs()["parentID1"])
	}

	if idVal, ok := mockSend.LastArgs()["parentID2"].(uint32); !ok || idVal != parentID2 {
		t.Errorf("Expected parentID2=%d, got %v", parentID2, mockSend.LastArgs()["parentID2"])
	}

	if resultVal, ok := mockSend.LastArgs()["result"].(uint8); !ok || resultVal != result {
		t.Errorf("Expected result=%d, got %v", result, mockSend.LastArgs()["result"])
	}

	// Test with invalid result
	err = marriageManager.SendAdoptReply(parentID1, parentID2, 2)
	if err == nil {
		t.Error("Expected error for invalid result, got nil")
	}
}
