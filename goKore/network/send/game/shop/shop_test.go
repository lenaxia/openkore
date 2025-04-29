package shop

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

// TestParseBuyBulkVender tests the ParseBuyBulkVender function
func TestParseBuyBulkVender(t *testing.T) {
	// Create test data
	itemInfo := []byte{
		0x01, 0x00, 0x02, 0x00, // item 1: amount=1, itemIndex=2
		0x03, 0x00, 0x04, 0x00, // item 2: amount=3, itemIndex=4
	}

	// Create args map
	args := map[string]interface{}{
		"itemInfo": itemInfo,
	}

	// Parse the data
	ParseBuyBulkVender(args)

	// Check that the items were parsed correctly
	items, ok := args["items"].([]map[string]interface{})
	if !ok {
		t.Fatal("ParseBuyBulkVender did not set items field correctly")
	}

	if len(items) != 2 {
		t.Fatalf("len(items) = %v, want 2", len(items))
	}

	// Check item 1
	if items[0]["amount"] != uint16(1) {
		t.Errorf("items[0][\"amount\"] = %v, want 1", items[0]["amount"])
	}
	if items[0]["itemIndex"] != uint16(2) {
		t.Errorf("items[0][\"itemIndex\"] = %v, want 2", items[0]["itemIndex"])
	}

	// Check item 2
	if items[1]["amount"] != uint16(3) {
		t.Errorf("items[1][\"amount\"] = %v, want 3", items[1]["amount"])
	}
	if items[1]["itemIndex"] != uint16(4) {
		t.Errorf("items[1][\"itemIndex\"] = %v, want 4", items[1]["itemIndex"])
	}
}

// TestReconstructBuyBulkVender tests the ReconstructBuyBulkVender function
func TestReconstructBuyBulkVender(t *testing.T) {
	// Create test data
	items := []map[string]interface{}{
		{"amount": uint16(1), "itemIndex": uint16(2)},
		{"amount": uint16(3), "itemIndex": uint16(4)},
	}

	// Create args map
	args := map[string]interface{}{
		"items": items,
	}

	// Reconstruct the data
	ReconstructBuyBulkVender(args)

	// Check that the itemInfo was reconstructed correctly
	itemInfo, ok := args["itemInfo"].([]byte)
	if !ok {
		t.Fatal("ReconstructBuyBulkVender did not set itemInfo field correctly")
	}

	// Expected itemInfo: amount=1, itemIndex=2, amount=3, itemIndex=4
	expected := []byte{
		0x01, 0x00, 0x02, 0x00, // item 1: amount=1, itemIndex=2
		0x03, 0x00, 0x04, 0x00, // item 2: amount=3, itemIndex=4
	}

	if len(itemInfo) != len(expected) {
		t.Fatalf("len(itemInfo) = %v, want %v", len(itemInfo), len(expected))
	}

	for i := 0; i < len(itemInfo); i++ {
		if itemInfo[i] != expected[i] {
			t.Errorf("itemInfo[%d] = %v, want %v", i, itemInfo[i], expected[i])
		}
	}
}

// TestSendBuyBulkVender tests the SendBuyBulkVender method
func TestSendBuyBulkVender(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["buy_bulk_vender"] = "0134"

	sm := NewShopManager(mockSend)

	// Test buying items from a vender
	venderID := uint32(12345)
	venderCID := uint32(67890)
	items := []map[string]interface{}{
		{"amount": uint16(1), "itemIndex": uint16(2)},
		{"amount": uint16(3), "itemIndex": uint16(4)},
	}
	err := sm.SendBuyBulkVender(venderID, items, venderCID)
	if err != nil {
		t.Fatalf("SendBuyBulkVender() returned error: %v", err)
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

	if args["venderID"] != venderID {
		t.Errorf("args[\"venderID\"] = %v, want %v", args["venderID"], venderID)
	}

	if args["venderCID"] != venderCID {
		t.Errorf("args[\"venderCID\"] = %v, want %v", args["venderCID"], venderCID)
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
		if argsItems[i]["amount"] != items[i]["amount"] {
			t.Errorf("argsItems[%d][\"amount\"] = %v, want %v", i, argsItems[i]["amount"], items[i]["amount"])
		}
		if argsItems[i]["itemIndex"] != items[i]["itemIndex"] {
			t.Errorf("argsItems[%d][\"itemIndex\"] = %v, want %v", i, argsItems[i]["itemIndex"], items[i]["itemIndex"])
		}
	}
}

// TestReconstructBuyBulkBuyer tests the ReconstructBuyBulkBuyer function
func TestReconstructBuyBulkBuyer(t *testing.T) {
	// Create test data
	items := []map[string]interface{}{
		{"ID": []byte{0x01, 0x02}, "itemID": uint16(3), "amount": uint16(4)},
		{"ID": []byte{0x05, 0x06}, "itemID": uint16(7), "amount": uint16(8)},
	}

	// Create args map
	args := map[string]interface{}{
		"items": items,
	}

	// Reconstruct the data
	ReconstructBuyBulkBuyer(args)

	// Check that the itemInfo was reconstructed correctly
	itemInfo, ok := args["itemInfo"].([]byte)
	if !ok {
		t.Fatal("ReconstructBuyBulkBuyer did not set itemInfo field correctly")
	}

	// Expected itemInfo: ID=0x0102, itemID=3, amount=4, ID=0x0506, itemID=7, amount=8
	expected := []byte{
		0x01, 0x02, 0x03, 0x00, 0x04, 0x00, // item 1: ID=0x0102, itemID=3, amount=4
		0x05, 0x06, 0x07, 0x00, 0x08, 0x00, // item 2: ID=0x0506, itemID=7, amount=8
	}

	if len(itemInfo) != len(expected) {
		t.Fatalf("len(itemInfo) = %v, want %v", len(itemInfo), len(expected))
	}

	for i := 0; i < len(itemInfo); i++ {
		if itemInfo[i] != expected[i] {
			t.Errorf("itemInfo[%d] = %v, want %v", i, itemInfo[i], expected[i])
		}
	}
}

// TestSendBuyBulkBuyer tests the SendBuyBulkBuyer method
func TestSendBuyBulkBuyer(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["buy_bulk_buyer"] = "0819"

	sm := NewShopManager(mockSend)

	// Test buying items from a buyer
	buyerID := uint32(12345)
	buyingStoreID := uint32(67890)
	items := []map[string]interface{}{
		{"ID": []byte{0x01, 0x02}, "itemID": uint16(3), "amount": uint16(4)},
		{"ID": []byte{0x05, 0x06}, "itemID": uint16(7), "amount": uint16(8)},
	}
	err := sm.SendBuyBulkBuyer(buyerID, items, buyingStoreID)
	if err != nil {
		t.Fatalf("SendBuyBulkBuyer() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0819"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["buyerID"] != buyerID {
		t.Errorf("args[\"buyerID\"] = %v, want %v", args["buyerID"], buyerID)
	}

	if args["buyingStoreID"] != buyingStoreID {
		t.Errorf("args[\"buyingStoreID\"] = %v, want %v", args["buyingStoreID"], buyingStoreID)
	}

	// Check that the length was calculated correctly
	expectedLen := 12 + (len(items) * 8)
	if args["len"] != expectedLen {
		t.Errorf("args[\"len\"] = %v, want %v", args["len"], expectedLen)
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
		if string(argsItems[i]["ID"].([]byte)) != string(items[i]["ID"].([]byte)) {
			t.Errorf("argsItems[%d][\"ID\"] = %v, want %v", i, argsItems[i]["ID"], items[i]["ID"])
		}
		if argsItems[i]["itemID"] != items[i]["itemID"] {
			t.Errorf("argsItems[%d][\"itemID\"] = %v, want %v", i, argsItems[i]["itemID"], items[i]["itemID"])
		}
		if argsItems[i]["amount"] != items[i]["amount"] {
			t.Errorf("argsItems[%d][\"amount\"] = %v, want %v", i, argsItems[i]["amount"], items[i]["amount"])
		}
	}
}

// TestReconstructBuyBulkOpenShop tests the ReconstructBuyBulkOpenShop function
func TestReconstructBuyBulkOpenShop(t *testing.T) {
	// Create test data
	items := []map[string]interface{}{
		{"nameID": uint16(1), "amount": uint16(2), "price": uint32(3)},
		{"nameID": uint16(4), "amount": uint16(5), "price": uint32(6)},
	}

	// Create args map
	args := map[string]interface{}{
		"items": items,
	}

	// Reconstruct the data
	ReconstructBuyBulkOpenShop(args)

	// Check that the itemInfo was reconstructed correctly
	itemInfo, ok := args["itemInfo"].([]byte)
	if !ok {
		t.Fatal("ReconstructBuyBulkOpenShop did not set itemInfo field correctly")
	}

	// Expected itemInfo: nameID=1, amount=2, price=3, nameID=4, amount=5, price=6
	expected := []byte{
		0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x00, 0x00, // item 1: nameID=1, amount=2, price=3
		0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x00, 0x00, // item 2: nameID=4, amount=5, price=6
	}

	if len(itemInfo) != len(expected) {
		t.Fatalf("len(itemInfo) = %v, want %v", len(itemInfo), len(expected))
	}

	for i := 0; i < len(itemInfo); i++ {
		if itemInfo[i] != expected[i] {
			t.Errorf("itemInfo[%d] = %v, want %v", i, itemInfo[i], expected[i])
		}
	}
}

// TestSendBuyBulkOpenShop tests the SendBuyBulkOpenShop method
func TestSendBuyBulkOpenShop(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["buy_bulk_openShop"] = "0815"

	sm := NewShopManager(mockSend)

	// Test opening a buying shop
	limitZeny := uint32(1000)
	result := uint8(1)
	storeName := "Test Shop"
	items := []map[string]interface{}{
		{"nameID": uint16(1), "amount": uint16(2), "price": uint32(3)},
		{"nameID": uint16(4), "amount": uint16(5), "price": uint32(6)},
	}
	err := sm.SendBuyBulkOpenShop(limitZeny, result, storeName, items)
	if err != nil {
		t.Fatalf("SendBuyBulkOpenShop() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0815"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["limitZeny"] != limitZeny {
		t.Errorf("args[\"limitZeny\"] = %v, want %v", args["limitZeny"], limitZeny)
	}

	if args["result"] != result {
		t.Errorf("args[\"result\"] = %v, want %v", args["result"], result)
	}

	if args["storeName"] != storeName {
		t.Errorf("args[\"storeName\"] = %v, want %v", args["storeName"], storeName)
	}

	// Check that the length was calculated correctly
	expectedLen := 89 + (len(items) * 8)
	if args["len"] != expectedLen {
		t.Errorf("args[\"len\"] = %v, want %v", args["len"], expectedLen)
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
		if argsItems[i]["nameID"] != items[i]["nameID"] {
			t.Errorf("argsItems[%d][\"nameID\"] = %v, want %v", i, argsItems[i]["nameID"], items[i]["nameID"])
		}
		if argsItems[i]["amount"] != items[i]["amount"] {
			t.Errorf("argsItems[%d][\"amount\"] = %v, want %v", i, argsItems[i]["amount"], items[i]["amount"])
		}
		if argsItems[i]["price"] != items[i]["price"] {
			t.Errorf("argsItems[%d][\"price\"] = %v, want %v", i, argsItems[i]["price"], items[i]["price"])
		}
	}
}

// TestSendEnteringBuyer tests the SendEnteringBuyer method
func TestSendEnteringBuyer(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["buy_bulk_request"] = "0817"

	sm := NewShopManager(mockSend)

	// Test entering a buyer's shop
	ID := uint32(12345)
	err := sm.SendEnteringBuyer(ID)
	if err != nil {
		t.Fatalf("SendEnteringBuyer() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0817"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != ID {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], ID)
	}
}

// TestSendCloseShop tests the SendCloseShop method
func TestSendCloseShop(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["shop_close"] = "012E"

	sm := NewShopManager(mockSend)

	// Test closing a shop
	err := sm.SendCloseShop()
	if err != nil {
		t.Fatalf("SendCloseShop() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["012E"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestSendCloseBuyShop tests the SendCloseBuyShop method
func TestSendCloseBuyShop(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["buy_bulk_closeShop"] = "0816"

	sm := NewShopManager(mockSend)

	// Test closing a buying shop
	err := sm.SendCloseBuyShop()
	if err != nil {
		t.Fatalf("SendCloseBuyShop() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["0816"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestSendEnteringVender tests the SendEnteringVender method
func TestSendEnteringVender(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["send_entering_vending"] = "0130"

	sm := NewShopManager(mockSend)

	// Test entering a vender's shop
	accountID := uint32(12345)
	err := sm.SendEnteringVender(accountID)
	if err != nil {
		t.Fatalf("SendEnteringVender() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0130"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["accountID"] != accountID {
		t.Errorf("args[\"accountID\"] = %v, want %v", args["accountID"], accountID)
	}
}

// TestReconstructShopOpen tests the ReconstructShopOpen function
func TestReconstructShopOpen(t *testing.T) {
	// Create test data
	items := []map[string]interface{}{
		{"ID": []byte{0x01, 0x02}, "amount": uint16(3), "price": uint32(4)},
		{"ID": []byte{0x05, 0x06}, "amount": uint16(7), "price": uint32(8)},
	}

	// Create args map
	args := map[string]interface{}{
		"items": items,
	}

	// Reconstruct the data
	ReconstructShopOpen(args)

	// Check that the vendingInfo was reconstructed correctly
	vendingInfo, ok := args["vendingInfo"].([]byte)
	if !ok {
		t.Fatal("ReconstructShopOpen did not set vendingInfo field correctly")
	}

	// Expected vendingInfo: ID=0x0102, amount=3, price=4, ID=0x0506, amount=7, price=8
	expected := []byte{
		0x01, 0x02, 0x03, 0x00, 0x04, 0x00, 0x00, 0x00, // item 1: ID=0x0102, amount=3, price=4
		0x05, 0x06, 0x07, 0x00, 0x08, 0x00, 0x00, 0x00, // item 2: ID=0x0506, amount=7, price=8
	}

	if len(vendingInfo) != len(expected) {
		t.Fatalf("len(vendingInfo) = %v, want %v", len(vendingInfo), len(expected))
	}

	for i := 0; i < len(vendingInfo); i++ {
		if vendingInfo[i] != expected[i] {
			t.Errorf("vendingInfo[%d] = %v, want %v", i, vendingInfo[i], expected[i])
		}
	}
}

// TestSendOpenShop tests the SendOpenShop method
func TestSendOpenShop(t *testing.T) {
	mockSend := NewMockSendForShop()
	mockSend.packetLUT["shop_open"] = "012D"

	sm := NewShopManager(mockSend)

	// Test opening a shop
	title := "Test Shop"
	items := []map[string]interface{}{
		{"ID": []byte{0x01, 0x02}, "amount": uint16(3), "price": uint32(4)},
		{"ID": []byte{0x05, 0x06}, "amount": uint16(7), "price": uint32(8)},
	}
	err := sm.SendOpenShop(title, items)
	if err != nil {
		t.Fatalf("SendOpenShop() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["012D"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	// Check that title was converted to bytes
	titleBytes, ok := args["title"].([]byte)
	if !ok {
		t.Errorf("args[\"title\"] is not a byte slice")
	} else if string(titleBytes) != title {
		t.Errorf("args[\"title\"] = %v, want %v", string(titleBytes), title)
	}

	if args["result"] != 1 {
		t.Errorf("args[\"result\"] = %v, want 1", args["result"])
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
		if string(argsItems[i]["ID"].([]byte)) != string(items[i]["ID"].([]byte)) {
			t.Errorf("argsItems[%d][\"ID\"] = %v, want %v", i, argsItems[i]["ID"], items[i]["ID"])
		}
		if argsItems[i]["amount"] != items[i]["amount"] {
			t.Errorf("argsItems[%d][\"amount\"] = %v, want %v", i, argsItems[i]["amount"], items[i]["amount"])
		}
		if argsItems[i]["price"] != items[i]["price"] {
			t.Errorf("argsItems[%d][\"price\"] = %v, want %v", i, argsItems[i]["price"], items[i]["price"])
		}
	}
}
