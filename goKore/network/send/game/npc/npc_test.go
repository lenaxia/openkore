package npc

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForNPC is a mock implementation of the Send interface for testing NPC functionality
type MockSendForNPC struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForNPC() *MockSendForNPC {
	return &MockSendForNPC{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForNPC) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForNPC) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForNPC) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForNPC) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForNPC) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForNPC) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForNPC) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForNPC) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForNPC) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForNPC) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForNPC) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForNPC) SetConnection(conn interface{}) {
}

func (ms *MockSendForNPC) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForNPC) GetTime() uint32 {
	return 12345
}

// TestNewNPCManager tests the NewNPCManager function
func TestNewNPCManager(t *testing.T) {
	mockSend := NewMockSendForNPC()
	nm := NewNPCManager(mockSend)

	if nm == nil {
		t.Fatal("NewNPCManager() returned nil")
	}

	if nm.baseSend == nil {
		t.Error("nm.baseSend was not set correctly")
	}
}

// TestSendNPCBuySellList tests the SendNPCBuySellList method
func TestSendNPCBuySellList(t *testing.T) {
	mockSend := NewMockSendForNPC()
	mockSend.packetLUT["request_buy_sell_list"] = "00C5"

	nm := NewNPCManager(mockSend)

	// Test requesting buy/sell list
	ID := uint32(12345)
	type_ := 0 // 0 = buy list, 1 = sell list
	err := nm.SendNPCBuySellList(ID, type_)
	if err != nil {
		t.Fatalf("SendNPCBuySellList() returned error: %v", err)
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

	if args["ID"] != ID {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], ID)
	}

	if args["type"] != type_ {
		t.Errorf("args[\"type\"] = %v, want %v", args["type"], type_)
	}
}

// TestSendNPCCreateRequest tests the SendNPCCreateRequest method
func TestSendNPCCreateRequest(t *testing.T) {
	mockSend := NewMockSendForNPC()
	mockSend.packetLUT["dynamicnpc_create_request"] = "0A16"

	nm := NewNPCManager(mockSend)

	// Test requesting to create an NPC
	name := "TestNPC"
	err := nm.SendNPCCreateRequest(name)
	if err != nil {
		t.Fatalf("SendNPCCreateRequest() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A16"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != name {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], name)
	}
}
