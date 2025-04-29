package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForEquipment is a mock implementation of the Send interface for testing equipment functionality
type MockSendForEquipment struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForEquipment() *MockSendForEquipment {
	return &MockSendForEquipment{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForEquipment) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForEquipment) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForEquipment) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForEquipment) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForEquipment) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForEquipment) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForEquipment) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForEquipment) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForEquipment) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForEquipment) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForEquipment) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForEquipment) SetConnection(conn interface{}) {
}

func (ms *MockSendForEquipment) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForEquipment) GetTime() uint32 {
	return 12345
}

// TestNewEquipmentManager tests the NewEquipmentManager function
func TestNewEquipmentManager(t *testing.T) {
	mockSend := NewMockSendForEquipment()
	em := NewEquipmentManager(mockSend)

	if em == nil {
		t.Fatal("NewEquipmentManager() returned nil")
	}

	if em.baseSend == nil {
		t.Error("em.baseSend was not set correctly")
	}
}

// TestEquipItem tests the EquipItem method
func TestEquipItem(t *testing.T) {
	mockSend := NewMockSendForEquipment()
	mockSend.packetLUT["equip_item"] = "00A9"

	em := NewEquipmentManager(mockSend)

	// Test equipping an item
	index := uint16(1)
	position := uint16(2) // 2 = armor
	err := em.EquipItem(index, position)
	if err != nil {
		t.Fatalf("EquipItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00A9"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["index"] != index {
		t.Errorf("args[\"index\"] = %v, want %v", args["index"], index)
	}

	if args["position"] != position {
		t.Errorf("args[\"position\"] = %v, want %v", args["position"], position)
	}
}

// TestUnequipItem tests the UnequipItem method
func TestUnequipItem(t *testing.T) {
	mockSend := NewMockSendForEquipment()
	mockSend.packetLUT["unequip_item"] = "00AB"

	em := NewEquipmentManager(mockSend)

	// Test unequipping an item
	index := uint16(1)
	err := em.UnequipItem(index)
	if err != nil {
		t.Fatalf("UnequipItem() returned error: %v", err)
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

	if args["index"] != index {
		t.Errorf("args[\"index\"] = %v, want %v", args["index"], index)
	}
}

// TestUpgradeItem tests the UpgradeItem method
func TestUpgradeItem(t *testing.T) {
	mockSend := NewMockSendForEquipment()
	mockSend.packetLUT["upgrade_item"] = "0181"

	em := NewEquipmentManager(mockSend)

	// Test upgrading an item
	index := uint16(1)
	err := em.UpgradeItem(index)
	if err != nil {
		t.Fatalf("UpgradeItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0181"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["index"] != index {
		t.Errorf("args[\"index\"] = %v, want %v", args["index"], index)
	}
}

// TestRefineItem tests the RefineItem method
func TestRefineItem(t *testing.T) {
	mockSend := NewMockSendForEquipment()
	mockSend.packetLUT["refine_item"] = "0222"

	em := NewEquipmentManager(mockSend)

	// Test refining an item
	index := uint16(1)
	err := em.RefineItem(index)
	if err != nil {
		t.Fatalf("RefineItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0222"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["index"] != index {
		t.Errorf("args[\"index\"] = %v, want %v", args["index"], index)
	}
}

// TestRepairItem tests the RepairItem method
func TestRepairItem(t *testing.T) {
	mockSend := NewMockSendForEquipment()
	mockSend.packetLUT["repair_item"] = "01FD"

	em := NewEquipmentManager(mockSend)

	// Test repairing an item
	index := uint16(1)
	itemID := uint16(1001)
	err := em.RepairItem(index, itemID)
	if err != nil {
		t.Fatalf("RepairItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["01FD"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["index"] != index {
		t.Errorf("args[\"index\"] = %v, want %v", args["index"], index)
	}

	if args["item_id"] != itemID {
		t.Errorf("args[\"item_id\"] = %v, want %v", args["item_id"], itemID)
	}
}
