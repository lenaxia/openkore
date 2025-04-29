package auction

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
			"auction_add_item_cancel": "0366",
			"auction_add_item":        "0367",
			"auction_create":          "0368",
			"auction_cancel":          "0369",
			"auction_buy":             "036A",
			"auction_search":          "036B",
			"auction_info_self":       "036C",
			"auction_sell_stop":       "036D",
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

// TestNewAuctionManager tests the NewAuctionManager function
func TestNewAuctionManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	auctionManager := NewAuctionManager(mockSend)

	if auctionManager == nil {
		t.Fatal("NewAuctionManager() returned nil")
	}

	if auctionManager.baseSend == nil {
		t.Error("auctionManager.baseSend was not set correctly")
	}
}

// TestSendAuctionAddItemCancel tests the SendAuctionAddItemCancel function
func TestSendAuctionAddItemCancel(t *testing.T) {
	mockSend := NewMockSend()
	auctionManager := NewAuctionManager(mockSend)

	// Test sending an auction add item cancel request
	flag := uint8(1)
	err := auctionManager.SendAuctionAddItemCancel(flag)
	if err != nil {
		t.Fatalf("SendAuctionAddItemCancel() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "auction_add_item_cancel" {
		t.Errorf("Expected packet ID 'auction_add_item_cancel', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if flagVal, ok := mockSend.LastArgs()["flag"].(uint8); !ok || flagVal != flag {
		t.Errorf("Expected flag=%d, got %v", flag, mockSend.LastArgs()["flag"])
	}
}

// TestSendAuctionAddItem tests the SendAuctionAddItem function
func TestSendAuctionAddItem(t *testing.T) {
	mockSend := NewMockSend()
	auctionManager := NewAuctionManager(mockSend)

	// Test sending an auction add item request
	itemID := uint32(12345)
	amount := uint16(1)
	err := auctionManager.SendAuctionAddItem(itemID, amount)
	if err != nil {
		t.Fatalf("SendAuctionAddItem() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "auction_add_item" {
		t.Errorf("Expected packet ID 'auction_add_item', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != itemID {
		t.Errorf("Expected ID=%d, got %v", itemID, mockSend.LastArgs()["ID"])
	}

	if amountVal, ok := mockSend.LastArgs()["amount"].(uint16); !ok || amountVal != amount {
		t.Errorf("Expected amount=%d, got %v", amount, mockSend.LastArgs()["amount"])
	}
}

// TestSendAuctionCreate tests the SendAuctionCreate function
func TestSendAuctionCreate(t *testing.T) {
	mockSend := NewMockSend()
	auctionManager := NewAuctionManager(mockSend)

	// Test sending an auction create request
	nowPrice := uint32(1000)
	maxPrice := uint32(2000)
	deleteTime := uint32(24) // hours
	err := auctionManager.SendAuctionCreate(nowPrice, maxPrice, deleteTime)
	if err != nil {
		t.Fatalf("SendAuctionCreate() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "auction_create" {
		t.Errorf("Expected packet ID 'auction_create', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if priceVal, ok := mockSend.LastArgs()["now_price"].(uint32); !ok || priceVal != nowPrice {
		t.Errorf("Expected now_price=%d, got %v", nowPrice, mockSend.LastArgs()["now_price"])
	}

	if priceVal, ok := mockSend.LastArgs()["max_price"].(uint32); !ok || priceVal != maxPrice {
		t.Errorf("Expected max_price=%d, got %v", maxPrice, mockSend.LastArgs()["max_price"])
	}

	if timeVal, ok := mockSend.LastArgs()["delete_time"].(uint32); !ok || timeVal != deleteTime {
		t.Errorf("Expected delete_time=%d, got %v", deleteTime, mockSend.LastArgs()["delete_time"])
	}
}

// TestSendAuctionCancel tests the SendAuctionCancel function
func TestSendAuctionCancel(t *testing.T) {
	mockSend := NewMockSend()
	auctionManager := NewAuctionManager(mockSend)

	// Test sending an auction cancel request
	auctionID := uint32(12345)
	err := auctionManager.SendAuctionCancel(auctionID)
	if err != nil {
		t.Fatalf("SendAuctionCancel() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "auction_cancel" {
		t.Errorf("Expected packet ID 'auction_cancel', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != auctionID {
		t.Errorf("Expected ID=%d, got %v", auctionID, mockSend.LastArgs()["ID"])
	}
}

// TestSendAuctionBuy tests the SendAuctionBuy function
func TestSendAuctionBuy(t *testing.T) {
	mockSend := NewMockSend()
	auctionManager := NewAuctionManager(mockSend)

	// Test sending an auction buy request
	auctionID := uint32(12345)
	price := uint32(1000)
	err := auctionManager.SendAuctionBuy(auctionID, price)
	if err != nil {
		t.Fatalf("SendAuctionBuy() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "auction_buy" {
		t.Errorf("Expected packet ID 'auction_buy', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != auctionID {
		t.Errorf("Expected ID=%d, got %v", auctionID, mockSend.LastArgs()["ID"])
	}

	if priceVal, ok := mockSend.LastArgs()["price"].(uint32); !ok || priceVal != price {
		t.Errorf("Expected price=%d, got %v", price, mockSend.LastArgs()["price"])
	}
}

// TestSendAuctionItemSearch tests the SendAuctionItemSearch function
func TestSendAuctionItemSearch(t *testing.T) {
	mockSend := NewMockSend()
	auctionManager := NewAuctionManager(mockSend)

	// Test cases for different search types
	testCases := []struct {
		searchType uint8
		name       string
	}{
		{0, "armor"},
		{1, "weapon"},
		{2, "card"},
		{3, "misc"},
		{4, "name"},
		{5, "auction id"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			price := uint32(1000)
			searchString := "test search"
			page := uint16(1)
			err := auctionManager.SendAuctionItemSearch(tc.searchType, price, searchString, page)
			if err != nil {
				t.Fatalf("SendAuctionItemSearch(%d) returned error: %v", tc.searchType, err)
			}

			if mockSend.sent == nil {
				t.Fatal("No packet was sent")
			}

			// Check that the correct packet ID was used
			if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "auction_search" {
				t.Errorf("Expected packet ID 'auction_search', got '%s'", packetID)
			}

			// Check that the correct arguments were used
			if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || typeVal != tc.searchType {
				t.Errorf("Expected type=%d, got %v", tc.searchType, mockSend.LastArgs()["type"])
			}

			if priceVal, ok := mockSend.LastArgs()["price"].(uint32); !ok || priceVal != price {
				t.Errorf("Expected price=%d, got %v", price, mockSend.LastArgs()["price"])
			}

			if strVal, ok := mockSend.LastArgs()["search_string"].(string); !ok || strVal != searchString {
				t.Errorf("Expected search_string=%s, got %v", searchString, mockSend.LastArgs()["search_string"])
			}

			if pageVal, ok := mockSend.LastArgs()["page"].(uint16); !ok || pageVal != page {
				t.Errorf("Expected page=%d, got %v", page, mockSend.LastArgs()["page"])
			}
		})
	}

	// Test with invalid search type
	err := auctionManager.SendAuctionItemSearch(6, 1000, "test", 1)
	if err == nil {
		t.Error("Expected error for invalid search type, got nil")
	}
}

// TestSendAuctionReqMyInfo tests the SendAuctionReqMyInfo function
func TestSendAuctionReqMyInfo(t *testing.T) {
	mockSend := NewMockSend()
	auctionManager := NewAuctionManager(mockSend)

	// Test sending an auction request my info request
	infoType := uint8(1)
	err := auctionManager.SendAuctionReqMyInfo(infoType)
	if err != nil {
		t.Fatalf("SendAuctionReqMyInfo() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "auction_info_self" {
		t.Errorf("Expected packet ID 'auction_info_self', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || typeVal != infoType {
		t.Errorf("Expected type=%d, got %v", infoType, mockSend.LastArgs()["type"])
	}
}

// TestSendAuctionMySellStop tests the SendAuctionMySellStop function
func TestSendAuctionMySellStop(t *testing.T) {
	mockSend := NewMockSend()
	auctionManager := NewAuctionManager(mockSend)

	// Test sending an auction my sell stop request
	auctionID := uint32(12345)
	err := auctionManager.SendAuctionMySellStop(auctionID)
	if err != nil {
		t.Fatalf("SendAuctionMySellStop() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "auction_sell_stop" {
		t.Errorf("Expected packet ID 'auction_sell_stop', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != auctionID {
		t.Errorf("Expected ID=%d, got %v", auctionID, mockSend.LastArgs()["ID"])
	}
}
