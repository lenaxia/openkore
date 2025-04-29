package deal

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
			"deal_item_add": "00E9",
			"deal_initiate": "00C5",
			"deal_reply":    "00C7",
			"deal_finalize": "00CF",
			"deal_cancel":   "00E6",
			"deal_trade":    "00EF",
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

// TestNewDealManager tests the NewDealManager function
func TestNewDealManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	dealManager := NewDealManager(mockSend)

	if dealManager == nil {
		t.Fatal("NewDealManager() returned nil")
	}

	if dealManager.baseSend == nil {
		t.Error("dealManager.baseSend was not set correctly")
	}
}

// TestSendDealAddItem tests the SendDealAddItem function
func TestSendDealAddItem(t *testing.T) {
	mockSend := NewMockSend()
	dealManager := NewDealManager(mockSend)

	// Test sending a deal add item request
	itemID := []byte{0x01, 0x02}
	amount := uint16(10)
	err := dealManager.SendDealAddItem(itemID, amount)
	if err != nil {
		t.Fatalf("SendDealAddItem() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "deal_item_add" {
		t.Errorf("Expected packet ID 'deal_item_add', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].([]byte); !ok || len(idVal) != 2 || idVal[0] != itemID[0] || idVal[1] != itemID[1] {
		t.Errorf("Expected ID=%v, got %v", itemID, mockSend.LastArgs()["ID"])
	}

	if amountVal, ok := mockSend.LastArgs()["amount"].(uint16); !ok || amountVal != amount {
		t.Errorf("Expected amount=%d, got %v", amount, mockSend.LastArgs()["amount"])
	}
}

// TestSendDeal tests the SendDeal function
func TestSendDeal(t *testing.T) {
	mockSend := NewMockSend()
	dealManager := NewDealManager(mockSend)

	// Test sending a deal initiate request
	playerID := uint32(12345)
	err := dealManager.SendDeal(playerID)
	if err != nil {
		t.Fatalf("SendDeal() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "deal_initiate" {
		t.Errorf("Expected packet ID 'deal_initiate', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != playerID {
		t.Errorf("Expected ID=%d, got %v", playerID, mockSend.LastArgs()["ID"])
	}
}

// TestSendDealReply tests the SendDealReply function
func TestSendDealReply(t *testing.T) {
	mockSend := NewMockSend()
	dealManager := NewDealManager(mockSend)

	// Test cases for different actions
	testCases := []struct {
		action uint8
		name   string
	}{
		{3, "accept"},
		{4, "cancel"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := dealManager.SendDealReply(tc.action)
			if err != nil {
				t.Fatalf("SendDealReply(%d) returned error: %v", tc.action, err)
			}

			if mockSend.sent == nil {
				t.Fatal("No packet was sent")
			}

			// Check that the correct packet ID was used
			if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "deal_reply" {
				t.Errorf("Expected packet ID 'deal_reply', got '%s'", packetID)
			}

			// Check that the correct arguments were used
			if actionVal, ok := mockSend.LastArgs()["action"].(uint8); !ok || actionVal != tc.action {
				t.Errorf("Expected action=%d, got %v", tc.action, mockSend.LastArgs()["action"])
			}
		})
	}

	// Test with invalid action
	err := dealManager.SendDealReply(5)
	if err == nil {
		t.Error("Expected error for invalid action, got nil")
	}
}

// TestSendDealFinalize tests the SendDealFinalize function
func TestSendDealFinalize(t *testing.T) {
	mockSend := NewMockSend()
	dealManager := NewDealManager(mockSend)

	// Test sending a deal finalize request
	err := dealManager.SendDealFinalize()
	if err != nil {
		t.Fatalf("SendDealFinalize() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "deal_finalize" {
		t.Errorf("Expected packet ID 'deal_finalize', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendCurrentDealCancel tests the SendCurrentDealCancel function
func TestSendCurrentDealCancel(t *testing.T) {
	mockSend := NewMockSend()
	dealManager := NewDealManager(mockSend)

	// Test sending a current deal cancel request
	err := dealManager.SendCurrentDealCancel()
	if err != nil {
		t.Fatalf("SendCurrentDealCancel() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "deal_cancel" {
		t.Errorf("Expected packet ID 'deal_cancel', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendDealTrade tests the SendDealTrade function
func TestSendDealTrade(t *testing.T) {
	mockSend := NewMockSend()
	dealManager := NewDealManager(mockSend)

	// Test sending a deal trade request
	err := dealManager.SendDealTrade()
	if err != nil {
		t.Fatalf("SendDealTrade() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "deal_trade" {
		t.Errorf("Expected packet ID 'deal_trade', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}
