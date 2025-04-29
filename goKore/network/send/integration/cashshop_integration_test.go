package integration

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
	"github.com/lenaxia/goKore/network/send/game/cashshop"
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
			"cash_shop_request_points": "0A6A",
			"cash_shop_open":           "0844",
			"cash_shop_close":          "0845",
			"cash_shop_buy":            "0288",
			"cash_shop_list":           "0848",
			"cash_buy_bulk":            "0972",
			"cash_shop_check_coupon":   "0974",
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

// TestCashShopIntegration tests the integration of the cash shop system
func TestCashShopIntegration(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	cashShopManager := cashshop.NewCashShopManager(mockSend)

	// Test a sequence of cash shop operations

	// 1. Open the cash shop
	err := cashShopManager.OpenShop()
	if err != nil {
		t.Fatalf("OpenShop() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "cash_shop_open" {
		t.Errorf("Expected packet ID 'cash_shop_open', got '%s'", packetID)
	}

	// 2. Request cash items list
	err = cashShopManager.RequestPoints()
	if err != nil {
		t.Fatalf("RequestPoints() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "cash_shop_request_points" {
		t.Errorf("Expected packet ID 'cash_shop_request_points', got '%s'", packetID)
	}

	// 3. Buy items from the cash shop
	itemID := 12345
	amount := 1
	kafraPoints := 1000

	err = cashShopManager.Buy(itemID, amount, kafraPoints)
	if err != nil {
		t.Fatalf("Buy() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "cash_shop_buy" {
		t.Errorf("Expected packet ID 'cash_shop_buy', got '%s'", packetID)
	}

	// 4. Close the cash shop
	err = cashShopManager.CloseShop()
	if err != nil {
		t.Fatalf("CloseShop() returned error: %v", err)
	}

	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "cash_shop_close" {
		t.Errorf("Expected packet ID 'cash_shop_close', got '%s'", packetID)
	}
}
