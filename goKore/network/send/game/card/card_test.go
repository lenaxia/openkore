package card

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
			"card_merge_request": "017A",
			"card_merge":         "017C",
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

// TestNewCardManager tests the NewCardManager function
func TestNewCardManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	cardManager := NewCardManager(mockSend)

	if cardManager == nil {
		t.Fatal("NewCardManager() returned nil")
	}

	if cardManager.baseSend == nil {
		t.Error("cardManager.baseSend was not set correctly")
	}
}

// TestSendCardMergeRequest tests the SendCardMergeRequest function
func TestSendCardMergeRequest(t *testing.T) {
	mockSend := NewMockSend()
	cardManager := NewCardManager(mockSend)

	// Test sending a card merge request
	cardID := uint32(12345)
	err := cardManager.SendCardMergeRequest(cardID)
	if err != nil {
		t.Fatalf("SendCardMergeRequest() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "card_merge_request" {
		t.Errorf("Expected packet ID 'card_merge_request', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["cardID"].(uint32); !ok || idVal != cardID {
		t.Errorf("Expected cardID=%d, got %v", cardID, mockSend.LastArgs()["cardID"])
	}
}

// TestSendCardMerge tests the SendCardMerge function
func TestSendCardMerge(t *testing.T) {
	mockSend := NewMockSend()
	cardManager := NewCardManager(mockSend)

	// Test sending a card merge
	cardID := uint32(12345)
	itemID := uint32(67890)
	err := cardManager.SendCardMerge(cardID, itemID)
	if err != nil {
		t.Fatalf("SendCardMerge() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "card_merge" {
		t.Errorf("Expected packet ID 'card_merge', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if cardIDVal, ok := mockSend.LastArgs()["cardID"].(uint32); !ok || cardIDVal != cardID {
		t.Errorf("Expected cardID=%d, got %v", cardID, mockSend.LastArgs()["cardID"])
	}

	if itemIDVal, ok := mockSend.LastArgs()["itemID"].(uint32); !ok || itemIDVal != itemID {
		t.Errorf("Expected itemID=%d, got %v", itemID, mockSend.LastArgs()["itemID"])
	}
}
