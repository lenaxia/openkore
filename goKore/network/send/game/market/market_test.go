package market

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForMarket is a mock implementation of the Send interface for testing market functionality
type MockSendForMarket struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForMarket() *MockSendForMarket {
	return &MockSendForMarket{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForMarket) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForMarket) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForMarket) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForMarket) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForMarket) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForMarket) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForMarket) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForMarket) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForMarket) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForMarket) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForMarket) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForMarket) SetConnection(conn interface{}) {
}

func (ms *MockSendForMarket) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForMarket) GetTime() uint32 {
	return 12345
}

// TestNewMarketManager tests the NewMarketManager function
func TestNewMarketManager(t *testing.T) {
	mockSend := NewMockSendForMarket()
	mm := NewMarketManager(mockSend)

	if mm == nil {
		t.Fatal("NewMarketManager() returned nil")
	}

	if mm.baseSend == nil {
		t.Error("mm.baseSend was not set correctly")
	}
}

// TestReconstructBuyBulkMarket tests the ReconstructBuyBulkMarket function
func TestReconstructBuyBulkMarket(t *testing.T) {
	// Create test data
	items := []map[string]interface{}{
		{"itemID": uint16(1), "amount": uint32(2)},
		{"itemID": uint16(3), "amount": uint32(4)},
	}

	// Create args map
	args := map[string]interface{}{
		"items": items,
	}

	// Reconstruct the data
	ReconstructBuyBulkMarket(args)

	// Check that the buyInfo was reconstructed correctly
	buyInfo, ok := args["buyInfo"].([]byte)
	if !ok {
		t.Fatal("ReconstructBuyBulkMarket did not set buyInfo field correctly")
	}

	// Expected buyInfo: itemID=1, amount=2, itemID=3, amount=4
	expected := []byte{
		0x01, 0x00, 0x02, 0x00, 0x00, 0x00, // item 1: itemID=1, amount=2
		0x03, 0x00, 0x04, 0x00, 0x00, 0x00, // item 2: itemID=3, amount=4
	}

	if len(buyInfo) != len(expected) {
		t.Fatalf("len(buyInfo) = %v, want %v", len(buyInfo), len(expected))
	}

	for i := 0; i < len(buyInfo); i++ {
		if buyInfo[i] != expected[i] {
			t.Errorf("buyInfo[%d] = %v, want %v", i, buyInfo[i], expected[i])
		}
	}
}

// TestSendBuyBulkMarket tests the SendBuyBulkMarket method
func TestSendBuyBulkMarket(t *testing.T) {
	mockSend := NewMockSendForMarket()
	mockSend.packetLUT["buy_bulk_market"] = "09D6"

	mm := NewMarketManager(mockSend)

	// Test buying items from the market
	items := []map[string]interface{}{
		{"itemID": uint16(1), "amount": uint32(2)},
		{"itemID": uint16(3), "amount": uint32(4)},
	}
	err := mm.SendBuyBulkMarket(items)
	if err != nil {
		t.Fatalf("SendBuyBulkMarket() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["09D6"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	// Check that the items were passed correctly
	argsItems, ok := args["items"].([]map[string]interface{})
	if !ok {
		t.Fatal("args[\"items\"] is not a slice of maps")
	}

	if len(argsItems) != len(items) {
		t.Fatalf("len(argsItems) = %v, want %v", len(argsItems), len(items))
	}

	for i := 0; i < len(items); i++ {
		if argsItems[i]["itemID"] != items[i]["itemID"] {
			t.Errorf("argsItems[%d][\"itemID\"] = %v, want %v", i, argsItems[i]["itemID"], items[i]["itemID"])
		}
		if argsItems[i]["amount"] != items[i]["amount"] {
			t.Errorf("argsItems[%d][\"amount\"] = %v, want %v", i, argsItems[i]["amount"], items[i]["amount"])
		}
	}
}

// TestSendMarketClose tests the SendMarketClose method
func TestSendMarketClose(t *testing.T) {
	mockSend := NewMockSendForMarket()
	mockSend.packetLUT["market_close"] = "09D8"

	mm := NewMarketManager(mockSend)

	// Test closing the market
	err := mm.SendMarketClose()
	if err != nil {
		t.Fatalf("SendMarketClose() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["09D8"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestSendSellBuyComplete tests the SendSellBuyComplete method
func TestSendSellBuyComplete(t *testing.T) {
	mockSend := NewMockSendForMarket()
	mockSend.packetLUT["sell_buy_complete"] = "09D4"

	mm := NewMarketManager(mockSend)

	// Test sending sell/buy complete
	err := mm.SendSellBuyComplete()
	if err != nil {
		t.Fatalf("SendSellBuyComplete() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["09D4"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}
