package cashshop

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
			"cash_shop_request_points": "0A6A", // Changed from 0A68 to avoid conflict with open_ui_request
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

// TestNewCashShopManager tests the NewCashShopManager function
func TestNewCashShopManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	cashShopManager := NewCashShopManager(mockSend)

	if cashShopManager == nil {
		t.Fatal("NewCashShopManager() returned nil")
	}

	if cashShopManager.baseSend == nil {
		t.Error("cashShopManager.baseSend was not set correctly")
	}

	// Test GetManagerName
	if name := cashShopManager.GetManagerName(); name != "CashShopManager" {
		t.Errorf("Expected manager name 'CashShopManager', got '%s'", name)
	}
}

// TestRequestPoints tests the RequestPoints function
func TestRequestPoints(t *testing.T) {
	mockSend := NewMockSend()
	cashShopManager := NewCashShopManager(mockSend)

	// Test sending a request for cash items list
	err := cashShopManager.RequestPoints()
	if err != nil {
		t.Fatalf("RequestPoints() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "cash_shop_request_points" {
		t.Errorf("Expected packet ID 'cash_shop_request_points', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestOpenShop tests the OpenShop function
func TestOpenShop(t *testing.T) {
	mockSend := NewMockSend()
	cashShopManager := NewCashShopManager(mockSend)

	// Test sending a cash shop open request
	err := cashShopManager.OpenShop()
	if err != nil {
		t.Fatalf("OpenShop() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "cash_shop_open" {
		t.Errorf("Expected packet ID 'cash_shop_open', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestCloseShop tests the CloseShop function
func TestCloseShop(t *testing.T) {
	mockSend := NewMockSend()
	cashShopManager := NewCashShopManager(mockSend)

	// Test sending a cash shop close request
	err := cashShopManager.CloseShop()
	if err != nil {
		t.Fatalf("CloseShop() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "cash_shop_close" {
		t.Errorf("Expected packet ID 'cash_shop_close', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestBuyBulk tests the BuyBulk function
func TestBuyBulk(t *testing.T) {
	mockSend := NewMockSend()
	cashShopManager := NewCashShopManager(mockSend)

	// Test sending a cash buy bulk request
	kafraPoints := 1000
	items := []int{1001, 1002, 1003}

	err := cashShopManager.BuyBulk(kafraPoints, items)
	if err != nil {
		t.Fatalf("BuyBulk() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "cash_buy_bulk" {
		t.Errorf("Expected packet ID 'cash_buy_bulk', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if args["kafra_points"] != kafraPoints {
		t.Errorf("Expected kafra_points %d, got %v", kafraPoints, args["kafra_points"])
	}

	if items, ok := args["items"].([]int); !ok {
		t.Errorf("Expected items to be []int, got %T", args["items"])
	} else if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}
}

// TestBuy tests the Buy function
func TestBuy(t *testing.T) {
	mockSend := NewMockSend()
	cashShopManager := NewCashShopManager(mockSend)

	// Test sending a cash shop buy request
	itemID := 1001
	amount := 5
	kafraPoints := 1000

	err := cashShopManager.Buy(itemID, amount, kafraPoints)
	if err != nil {
		t.Fatalf("Buy() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "cash_shop_buy" {
		t.Errorf("Expected packet ID 'cash_shop_buy', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if args["item_id"] != itemID {
		t.Errorf("Expected item_id %d, got %v", itemID, args["item_id"])
	}
	if args["amount"] != amount {
		t.Errorf("Expected amount %d, got %v", amount, args["amount"])
	}
	if args["kafra_points"] != kafraPoints {
		t.Errorf("Expected kafra_points %d, got %v", kafraPoints, args["kafra_points"])
	}
}

// TestList tests the List function
func TestList(t *testing.T) {
	mockSend := NewMockSend()
	cashShopManager := NewCashShopManager(mockSend)

	// Test sending a cash shop list request
	tab := 2

	err := cashShopManager.List(tab)
	if err != nil {
		t.Fatalf("List() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "cash_shop_list" {
		t.Errorf("Expected packet ID 'cash_shop_list', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if args["tab"] != tab {
		t.Errorf("Expected tab %d, got %v", tab, args["tab"])
	}
}

// TestCheckCoupon tests the CheckCoupon function
func TestCheckCoupon(t *testing.T) {
	mockSend := NewMockSend()
	cashShopManager := NewCashShopManager(mockSend)

	// Test sending a cash shop check coupon request
	couponCode := "ABC123"

	err := cashShopManager.CheckCoupon(couponCode)
	if err != nil {
		t.Fatalf("CheckCoupon() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "cash_shop_check_coupon" {
		t.Errorf("Expected packet ID 'cash_shop_check_coupon', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	args := mockSend.LastArgs()
	if args["coupon_code"] != couponCode {
		t.Errorf("Expected coupon_code %s, got %v", couponCode, args["coupon_code"])
	}
}
