package craft

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForCraft is a mock implementation of the Send interface for testing craft functionality
type MockSendForCraft struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForCraft() *MockSendForCraft {
	return &MockSendForCraft{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForCraft) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForCraft) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForCraft) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForCraft) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForCraft) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForCraft) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForCraft) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForCraft) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForCraft) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForCraft) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForCraft) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForCraft) SetConnection(conn interface{}) {
}

func (ms *MockSendForCraft) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForCraft) GetTime() uint32 {
	return 12345
}

// TestNewCraftManager tests the NewCraftManager function
func TestNewCraftManager(t *testing.T) {
	mockSend := NewMockSendForCraft()
	cm := NewCraftManager(mockSend)

	if cm == nil {
		t.Fatal("NewCraftManager() returned nil")
	}

	if cm.baseSend == nil {
		t.Error("cm.baseSend was not set correctly")
	}
}

// TestSendIdentify tests the SendIdentify method
func TestSendIdentify(t *testing.T) {
	mockSend := NewMockSendForCraft()
	mockSend.packetLUT["identify"] = "0178"

	cm := NewCraftManager(mockSend)

	// Test identifying an item
	ID := uint16(12345)
	err := cm.SendIdentify(ID)
	if err != nil {
		t.Fatalf("SendIdentify() returned error: %v", err)
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

	if args["ID"] != ID {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], ID)
	}
}

// TestSendWeaponRefine tests the SendWeaponRefine method
func TestSendWeaponRefine(t *testing.T) {
	mockSend := NewMockSendForCraft()
	mockSend.packetLUT["refine_item"] = "0222"

	cm := NewCraftManager(mockSend)

	// Test refining a weapon
	ID := uint16(12345)
	err := cm.SendWeaponRefine(ID)
	if err != nil {
		t.Fatalf("SendWeaponRefine() returned error: %v", err)
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

	if args["ID"] != ID {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], ID)
	}
}

// TestSendCooking tests the SendCooking method
func TestSendCooking(t *testing.T) {
	mockSend := NewMockSendForCraft()
	mockSend.packetLUT["cook_request"] = "025B"

	cm := NewCraftManager(mockSend)

	// Test cooking an item
	type_ := uint16(1)
	nameID := uint16(12345)
	err := cm.SendCooking(type_, nameID)
	if err != nil {
		t.Fatalf("SendCooking() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["025B"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["type"] != type_ {
		t.Errorf("args[\"type\"] = %v, want %v", args["type"], type_)
	}

	if args["nameID"] != nameID {
		t.Errorf("args[\"nameID\"] = %v, want %v", args["nameID"], nameID)
	}
}

// TestSendRepairItem tests the SendRepairItem method
func TestSendRepairItem(t *testing.T) {
	mockSend := NewMockSendForCraft()
	mockSend.packetLUT["repair_item"] = "01FD"

	cm := NewCraftManager(mockSend)

	// Test repairing an item
	index := uint16(1)
	nameID := uint16(12345)
	upgrade := uint8(7)
	cards := []uint32{1, 2, 3, 4}
	err := cm.SendRepairItem(index, nameID, upgrade, cards)
	if err != nil {
		t.Fatalf("SendRepairItem() returned error: %v", err)
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

	if args["nameID"] != nameID {
		t.Errorf("args[\"nameID\"] = %v, want %v", args["nameID"], nameID)
	}

	if args["upgrade"] != upgrade {
		t.Errorf("args[\"upgrade\"] = %v, want %v", args["upgrade"], upgrade)
	}

	// Check that the cards were passed correctly
	argsCards, ok := args["cards"].([]uint32)
	if !ok {
		t.Fatal("args[\"cards\"] is not a slice of uint32")
	}

	if len(argsCards) != len(cards) {
		t.Fatalf("len(argsCards) = %v, want %v", len(argsCards), len(cards))
	}

	for i := 0; i < len(cards); i++ {
		if argsCards[i] != cards[i] {
			t.Errorf("argsCards[%d] = %v, want %v", i, argsCards[i], cards[i])
		}
	}
}

// TestSendArrowCraft tests the SendArrowCraft method
func TestSendArrowCraft(t *testing.T) {
	mockSend := NewMockSendForCraft()
	mockSend.packetLUT["make_arrow"] = "01AE"

	cm := NewCraftManager(mockSend)

	// Test crafting arrows
	nameID := uint16(12345)
	err := cm.SendArrowCraft(nameID)
	if err != nil {
		t.Fatalf("SendArrowCraft() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["01AE"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["nameID"] != nameID {
		t.Errorf("args[\"nameID\"] = %v, want %v", args["nameID"], nameID)
	}
}

// TestSendCardMergeRequest tests the SendCardMergeRequest method
func TestSendCardMergeRequest(t *testing.T) {
	mockSend := NewMockSendForCraft()
	mockSend.packetLUT["card_merge_request"] = "017C"

	cm := NewCraftManager(mockSend)

	// Test requesting to merge a card
	cardID := uint16(12345)
	err := cm.SendCardMergeRequest(cardID)
	if err != nil {
		t.Fatalf("SendCardMergeRequest() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["017C"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["cardID"] != cardID {
		t.Errorf("args[\"cardID\"] = %v, want %v", args["cardID"], cardID)
	}
}

// TestSendCardMerge tests the SendCardMerge method
func TestSendCardMerge(t *testing.T) {
	mockSend := NewMockSendForCraft()
	mockSend.packetLUT["card_merge"] = "017D"

	cm := NewCraftManager(mockSend)

	// Test merging a card with an item
	cardID := uint16(12345)
	itemID := uint16(32890)
	err := cm.SendCardMerge(cardID, itemID)
	if err != nil {
		t.Fatalf("SendCardMerge() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["017D"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["cardID"] != cardID {
		t.Errorf("args[\"cardID\"] = %v, want %v", args["cardID"], cardID)
	}

	if args["itemID"] != itemID {
		t.Errorf("args[\"itemID\"] = %v, want %v", args["itemID"], itemID)
	}
}

// TestSendMakeItemRequest tests the SendMakeItemRequest method
func TestSendMakeItemRequest(t *testing.T) {
	mockSend := NewMockSendForCraft()
	mockSend.packetLUT["make_item_request"] = "018E"

	cm := NewCraftManager(mockSend)

	// Test making an item
	nameID := uint16(12345)
	material1 := uint16(1)
	material2 := uint16(2)
	material3 := uint16(3)
	err := cm.SendMakeItemRequest(nameID, material1, material2, material3)
	if err != nil {
		t.Fatalf("SendMakeItemRequest() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["018E"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["nameID"] != nameID {
		t.Errorf("args[\"nameID\"] = %v, want %v", args["nameID"], nameID)
	}

	if args["material_nameID1"] != material1 {
		t.Errorf("args[\"material_nameID1\"] = %v, want %v", args["material_nameID1"], material1)
	}

	if args["material_nameID2"] != material2 {
		t.Errorf("args[\"material_nameID2\"] = %v, want %v", args["material_nameID2"], material2)
	}

	if args["material_nameID3"] != material3 {
		t.Errorf("args[\"material_nameID3\"] = %v, want %v", args["material_nameID3"], material3)
	}
}
