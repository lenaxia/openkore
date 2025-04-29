package buyingstore

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
			"buy_bulk":                       "00C8",
			"sell_bulk":                      "00C9",
			"search_store_close":             "0835",
			"search_store_info":              "0838",
			"search_store_request_next_page": "0839",
			"search_store_select":            "083B",
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

// TestNewBuyingStoreManager tests the NewBuyingStoreManager function
func TestNewBuyingStoreManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	buyingStoreManager := NewBuyingStoreManager(mockSend)

	if buyingStoreManager == nil {
		t.Fatal("NewBuyingStoreManager() returned nil")
	}

	if buyingStoreManager.baseSend == nil {
		t.Error("buyingStoreManager.baseSend was not set correctly")
	}
}

// TestSendBuyBulk tests the SendBuyBulk function
func TestSendBuyBulk(t *testing.T) {
	mockSend := NewMockSend()
	buyingStoreManager := NewBuyingStoreManager(mockSend)

	// Test sending a buy bulk request
	items := []map[string]interface{}{
		{
			"amount": uint16(10),
			"itemID": uint16(501),
		},
		{
			"amount": uint16(5),
			"itemID": uint16(502),
		},
	}
	err := buyingStoreManager.SendBuyBulk(items)
	if err != nil {
		t.Fatalf("SendBuyBulk() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "buy_bulk" {
		t.Errorf("Expected packet ID 'buy_bulk', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if items, ok := mockSend.LastArgs()["items"].([]map[string]interface{}); !ok {
		t.Errorf("Expected items to be []map[string]interface{}, got %T", mockSend.LastArgs()["items"])
	} else {
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}
	}
}

// TestReconstructBuyBulk tests the ReconstructBuyBulk function
func TestReconstructBuyBulk(t *testing.T) {
	mockSend := NewMockSend()
	buyingStoreManager := NewBuyingStoreManager(mockSend)

	// Test reconstructing a buy bulk packet
	items := []map[string]interface{}{
		{
			"amount": uint16(10),
			"itemID": uint16(501),
		},
		{
			"amount": uint16(5),
			"itemID": uint16(502),
		},
	}
	args := map[string]interface{}{
		"items": items,
	}
	err := buyingStoreManager.ReconstructBuyBulk(args)
	if err != nil {
		t.Fatalf("ReconstructBuyBulk() returned error: %v", err)
	}

	// Check that the buyInfo field was added
	if _, ok := args["buyInfo"]; !ok {
		t.Error("buyInfo field was not added to args")
	}
}

// TestSendSellBulk tests the SendSellBulk function
func TestSendSellBulk(t *testing.T) {
	mockSend := NewMockSend()
	buyingStoreManager := NewBuyingStoreManager(mockSend)

	// Test sending a sell bulk request
	items := []map[string]interface{}{
		{
			"ID":     []byte{0x01, 0x02},
			"amount": uint16(10),
		},
		{
			"ID":     []byte{0x03, 0x04},
			"amount": uint16(5),
		},
	}
	err := buyingStoreManager.SendSellBulk(items)
	if err != nil {
		t.Fatalf("SendSellBulk() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "sell_bulk" {
		t.Errorf("Expected packet ID 'sell_bulk', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if items, ok := mockSend.LastArgs()["items"].([]map[string]interface{}); !ok {
		t.Errorf("Expected items to be []map[string]interface{}, got %T", mockSend.LastArgs()["items"])
	} else {
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}
	}
}

// TestReconstructSellBulk tests the ReconstructSellBulk function
func TestReconstructSellBulk(t *testing.T) {
	mockSend := NewMockSend()
	buyingStoreManager := NewBuyingStoreManager(mockSend)

	// Test reconstructing a sell bulk packet
	items := []map[string]interface{}{
		{
			"ID":     []byte{0x01, 0x02},
			"amount": uint16(10),
		},
		{
			"ID":     []byte{0x03, 0x04},
			"amount": uint16(5),
		},
	}
	args := map[string]interface{}{
		"items": items,
	}
	err := buyingStoreManager.ReconstructSellBulk(args)
	if err != nil {
		t.Fatalf("ReconstructSellBulk() returned error: %v", err)
	}

	// Check that the sellInfo field was added
	if _, ok := args["sellInfo"]; !ok {
		t.Error("sellInfo field was not added to args")
	}
}

// TestSendSearchStoreClose tests the SendSearchStoreClose function
func TestSendSearchStoreClose(t *testing.T) {
	mockSend := NewMockSend()
	buyingStoreManager := NewBuyingStoreManager(mockSend)

	// Test sending a search store close request
	err := buyingStoreManager.SendSearchStoreClose()
	if err != nil {
		t.Fatalf("SendSearchStoreClose() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "search_store_close" {
		t.Errorf("Expected packet ID 'search_store_close', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendSearchStoreSearch tests the SendSearchStoreSearch function
func TestSendSearchStoreSearch(t *testing.T) {
	mockSend := NewMockSend()
	buyingStoreManager := NewBuyingStoreManager(mockSend)

	// Test sending a search store search request
	searchType := uint8(0)
	maxPrice := uint32(1000)
	minPrice := uint32(500)
	itemList := []uint16{501, 502}
	cardList := []uint16{4001, 4002}
	err := buyingStoreManager.SendSearchStoreSearch(searchType, maxPrice, minPrice, itemList, cardList)
	if err != nil {
		t.Fatalf("SendSearchStoreSearch() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "search_store_info" {
		t.Errorf("Expected packet ID 'search_store_info', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || typeVal != searchType {
		t.Errorf("Expected type=%d, got %v", searchType, mockSend.LastArgs()["type"])
	}

	if priceVal, ok := mockSend.LastArgs()["max_price"].(uint32); !ok || priceVal != maxPrice {
		t.Errorf("Expected max_price=%d, got %v", maxPrice, mockSend.LastArgs()["max_price"])
	}

	if priceVal, ok := mockSend.LastArgs()["min_price"].(uint32); !ok || priceVal != minPrice {
		t.Errorf("Expected min_price=%d, got %v", minPrice, mockSend.LastArgs()["min_price"])
	}

	if items, ok := mockSend.LastArgs()["item_list"].([]uint16); !ok {
		t.Errorf("Expected item_list to be []uint16, got %T", mockSend.LastArgs()["item_list"])
	} else {
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}
	}

	if cards, ok := mockSend.LastArgs()["card_list"].([]uint16); !ok {
		t.Errorf("Expected card_list to be []uint16, got %T", mockSend.LastArgs()["card_list"])
	} else {
		if len(cards) != 2 {
			t.Errorf("Expected 2 cards, got %d", len(cards))
		}
	}
}

// TestReconstructSearchStoreInfo tests the ReconstructSearchStoreInfo function
func TestReconstructSearchStoreInfo(t *testing.T) {
	mockSend := NewMockSend()
	buyingStoreManager := NewBuyingStoreManager(mockSend)

	// Test reconstructing a search store info packet
	itemList := []uint16{501, 502}
	cardList := []uint16{4001, 4002}
	args := map[string]interface{}{
		"item_list": itemList,
		"card_list": cardList,
	}
	err := buyingStoreManager.ReconstructSearchStoreInfo(args)
	if err != nil {
		t.Fatalf("ReconstructSearchStoreInfo() returned error: %v", err)
	}

	// Check that the item_count and card_count fields were added
	if itemCount, ok := args["item_count"].(int); !ok || itemCount != 2 {
		t.Errorf("Expected item_count=2, got %v", args["item_count"])
	}

	if cardCount, ok := args["card_count"].(int); !ok || cardCount != 2 {
		t.Errorf("Expected card_count=2, got %v", args["card_count"])
	}

	// Check that the item_card_list field was added
	if _, ok := args["item_card_list"]; !ok {
		t.Error("item_card_list field was not added to args")
	}
}

// TestSendSearchStoreRequestNextPage tests the SendSearchStoreRequestNextPage function
func TestSendSearchStoreRequestNextPage(t *testing.T) {
	mockSend := NewMockSend()
	buyingStoreManager := NewBuyingStoreManager(mockSend)

	// Test sending a search store request next page request
	err := buyingStoreManager.SendSearchStoreRequestNextPage()
	if err != nil {
		t.Fatalf("SendSearchStoreRequestNextPage() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "search_store_request_next_page" {
		t.Errorf("Expected packet ID 'search_store_request_next_page', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendSearchStoreSelect tests the SendSearchStoreSelect function
func TestSendSearchStoreSelect(t *testing.T) {
	mockSend := NewMockSend()
	buyingStoreManager := NewBuyingStoreManager(mockSend)

	// Test sending a search store select request
	accountID := uint32(12345)
	storeID := uint32(67890)
	nameID := uint16(501)
	err := buyingStoreManager.SendSearchStoreSelect(accountID, storeID, nameID)
	if err != nil {
		t.Fatalf("SendSearchStoreSelect() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "search_store_select" {
		t.Errorf("Expected packet ID 'search_store_select', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["accountID"].(uint32); !ok || idVal != accountID {
		t.Errorf("Expected accountID=%d, got %v", accountID, mockSend.LastArgs()["accountID"])
	}

	if idVal, ok := mockSend.LastArgs()["storeID"].(uint32); !ok || idVal != storeID {
		t.Errorf("Expected storeID=%d, got %v", storeID, mockSend.LastArgs()["storeID"])
	}

	if idVal, ok := mockSend.LastArgs()["nameID"].(uint16); !ok || idVal != nameID {
		t.Errorf("Expected nameID=%d, got %v", nameID, mockSend.LastArgs()["nameID"])
	}
}
