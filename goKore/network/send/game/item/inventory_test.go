package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForInventory is a mock implementation of the Send interface for testing inventory functionality
type MockSendForInventory struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForInventory() *MockSendForInventory {
	return &MockSendForInventory{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForInventory) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForInventory) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForInventory) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForInventory) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForInventory) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForInventory) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForInventory) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForInventory) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForInventory) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForInventory) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForInventory) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForInventory) SetConnection(conn interface{}) {
}

func (ms *MockSendForInventory) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForInventory) GetTime() uint32 {
	return 12345
}

// TestNewInventoryManager tests the NewInventoryManager function
func TestNewInventoryManager(t *testing.T) {
	mockSend := NewMockSendForInventory()
	im := NewInventoryManager(mockSend)

	if im == nil {
		t.Fatal("NewInventoryManager() returned nil")
	}

	if im.baseSend == nil {
		t.Error("im.baseSend was not set correctly")
	}
}

// TestUseItem tests the UseItem method
func TestUseItem(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["use_item"] = "00A7"

	im := NewInventoryManager(mockSend)

	// Test using an item
	index := uint16(1)
	targetID := uint32(0) // 0 means self
	err := im.UseItem(index, targetID)
	if err != nil {
		t.Fatalf("UseItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00A7"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["index"] != index {
		t.Errorf("args[\"index\"] = %v, want %v", args["index"], index)
	}

	if args["target_id"] != targetID {
		t.Errorf("args[\"target_id\"] = %v, want %v", args["target_id"], targetID)
	}
}

// TestDropItem tests the DropItem method
func TestDropItem(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["drop_item"] = "00A2"

	im := NewInventoryManager(mockSend)

	// Test dropping an item
	index := uint16(1)
	amount := uint16(10)
	err := im.DropItem(index, amount)
	if err != nil {
		t.Fatalf("DropItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00A2"]
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

// TestMoveItem tests the MoveItem method
func TestMoveItem(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["move_item"] = "00AB"

	im := NewInventoryManager(mockSend)

	// Test moving an item
	fromIndex := uint16(1)
	toIndex := uint16(2)
	amount := uint16(10)
	err := im.MoveItem(fromIndex, toIndex, amount)
	if err != nil {
		t.Fatalf("MoveItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00AB"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["from_index"] != fromIndex {
		t.Errorf("args[\"from_index\"] = %v, want %v", args["from_index"], fromIndex)
	}

	if args["to_index"] != toIndex {
		t.Errorf("args[\"to_index\"] = %v, want %v", args["to_index"], toIndex)
	}

	if args["amount"] != amount {
		t.Errorf("args[\"amount\"] = %v, want %v", args["amount"], amount)
	}
}

// TestSplitItem tests the SplitItem method
func TestSplitItem(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["split_item"] = "00AC"

	im := NewInventoryManager(mockSend)

	// Test splitting an item
	index := uint16(1)
	amount := uint16(5)
	err := im.SplitItem(index, amount)
	if err != nil {
		t.Fatalf("SplitItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00AC"]
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

// TestIdentifyItem tests the IdentifyItem method
func TestIdentifyItem(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["identify_item"] = "0178"

	im := NewInventoryManager(mockSend)

	// Test identifying an item
	index := uint16(1)
	err := im.IdentifyItem(index)
	if err != nil {
		t.Fatalf("IdentifyItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0178"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["index"] != index {
		t.Errorf("args[\"index\"] = %v, want %v", args["index"], index)
	}
}

// TestSendTake tests the SendTake method
func TestSendTake(t *testing.T) {
	mockSend := NewMockSendForInventory()
	mockSend.packetLUT["take"] = "00A1"

	im := NewInventoryManager(mockSend)

	// Test taking an item from the ground
	objectID := uint32(12345)
	err := im.SendTake(objectID)
	if err != nil {
		t.Fatalf("SendTake() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00A1"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["object_id"] != objectID {
		t.Errorf("args[\"object_id\"] = %v, want %v", args["object_id"], objectID)
	}
}
