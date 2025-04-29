package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForShop is a mock implementation of the Send interface for testing shop functionality
type MockSendForShop struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForShop() *MockSendForShop {
	return &MockSendForShop{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForShop) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForShop) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForShop) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForShop) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForShop) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForShop) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForShop) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForShop) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForShop) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForShop) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForShop) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForShop) SetConnection(conn interface{}) {
}

func (ms *MockSendForShop) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForShop) GetTime() uint32 {
	return 12345
}

// TestNewShopManager tests the NewShopManager function
func TestNewShopManager(t *testing.T) {
	mockSend := NewMockSendForShop()
	sm := NewShopManager(mockSend)

	if sm == nil {
		t.Fatal("NewShopManager() returned nil")
	}

	if sm.baseSend == nil {
		t.Error("sm.baseSend was not set correctly")
	}
}

// TestOpenNpcShop tests the OpenNpcShop method
func TestOpenNpcShop(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["open_npc_shop"] = "00C5"

	sm := NewShopManager(mockSend)

	// Test opening an NPC shop
	npcID := uint32(12345)
	err := sm.OpenNpcShop(npcID)
	if err != nil {
		t.Fatalf("OpenNpcShop() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00C5"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["npc_id"] != npcID {
		t.Errorf("args[\"npc_id\"] = %v, want %v", args["npc_id"], npcID)
	}
}

// TestBuyItem tests the BuyItem method
func TestBuyItem(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["buy_item"] = "00C8"

	sm := NewShopManager(mockSend)

	// Test buying an item
	itemID := uint16(501)
	amount := uint16(10)
	err := sm.BuyItem(itemID, amount)
	if err != nil {
		t.Fatalf("BuyItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00C8"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["item_id"] != itemID {
		t.Errorf("args[\"item_id\"] = %v, want %v", args["item_id"], itemID)
	}

	if args["amount"] != amount {
		t.Errorf("args[\"amount\"] = %v, want %v", args["amount"], amount)
	}
}

// TestSellItem tests the SellItem method
func TestSellItem(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["sell_item"] = "00C9"

	sm := NewShopManager(mockSend)

	// Test selling an item
	index := uint16(1)
	amount := uint16(5)
	err := sm.SellItem(index, amount)
	if err != nil {
		t.Fatalf("SellItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00C9"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["index"] != index {
		t.Errorf("args[\"index\"] = %v, want %v", args["index"], index)
	}

	if args["amount"] != amount {
		t.Errorf("args[\"amount\"] = %v, want %v", args["amount"], amount)
	}
}

// TestCloseShop tests the CloseShop method
func TestCloseShop(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["close_shop"] = "00CC"

	sm := NewShopManager(mockSend)

	// Test closing a shop
	err := sm.CloseShop()
	if err != nil {
		t.Fatalf("CloseShop() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	_, exists := mockSend.reconstructArgs["00CC"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestOpenVendingShop tests the OpenVendingShop method
func TestOpenVendingShop(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["open_vending_shop"] = "01B2"

	sm := NewShopManager(mockSend)

	// Test opening a vending shop
	title := "My Shop"
	items := []VendingItem{
		{Index: 1, Amount: 10, Price: 1000},
		{Index: 2, Amount: 5, Price: 2000},
	}
	err := sm.OpenVendingShop(title, items)
	if err != nil {
		t.Fatalf("OpenVendingShop() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["01B2"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["title"] != title {
		t.Errorf("args[\"title\"] = %v, want %v", args["title"], title)
	}

	if len(args["items"].([]VendingItem)) != len(items) {
		t.Errorf("len(args[\"items\"]) = %v, want %v", len(args["items"].([]VendingItem)), len(items))
	}
}

// TestBuyVendingItem tests the BuyVendingItem method
func TestBuyVendingItem(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["buy_vending_item"] = "0134"

	sm := NewShopManager(mockSend)

	// Test buying an item from a vending shop
	vendorID := uint32(12345)
	index := uint16(1)
	amount := uint16(5)
	err := sm.BuyVendingItem(vendorID, index, amount)
	if err != nil {
		t.Fatalf("BuyVendingItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0134"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["vendor_id"] != vendorID {
		t.Errorf("args[\"vendor_id\"] = %v, want %v", args["vendor_id"], vendorID)
	}

	if args["index"] != index {
		t.Errorf("args[\"index\"] = %v, want %v", args["index"], index)
	}

	if args["amount"] != amount {
		t.Errorf("args[\"amount\"] = %v, want %v", args["amount"], amount)
	}
}
