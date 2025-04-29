package homunculus

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForHomunculus is a mock implementation of the Send interface for testing homunculus functionality
type MockSendForHomunculus struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForHomunculus() *MockSendForHomunculus {
	return &MockSendForHomunculus{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForHomunculus) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForHomunculus) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForHomunculus) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForHomunculus) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForHomunculus) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForHomunculus) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForHomunculus) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForHomunculus) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForHomunculus) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForHomunculus) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForHomunculus) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForHomunculus) SetConnection(conn interface{}) {
}

func (ms *MockSendForHomunculus) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForHomunculus) GetTime() uint32 {
	return 12345
}

// TestNewHomunculusManager tests the NewHomunculusManager function
func TestNewHomunculusManager(t *testing.T) {
	mockSend := NewMockSendForHomunculus()
	hm := NewHomunculusManager(mockSend)

	if hm == nil {
		t.Fatal("NewHomunculusManager() returned nil")
	}

	if hm.baseSend == nil {
		t.Error("hm.baseSend was not set correctly")
	}
}

// TestSendHomunculusName tests the SendHomunculusName method
func TestSendHomunculusName(t *testing.T) {
	mockSend := NewMockSendForHomunculus()
	mockSend.packetLUT["homunculus_name"] = "0231"

	hm := NewHomunculusManager(mockSend)

	// Test renaming a homunculus
	name := "TestHomunculus"
	err := hm.SendHomunculusName(name)
	if err != nil {
		t.Fatalf("SendHomunculusName() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0231"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	// Check that name was converted to bytes
	nameBytes, ok := args["name"].([]byte)
	if !ok {
		t.Errorf("args[\"name\"] is not a byte slice")
	} else if string(nameBytes) != name {
		t.Errorf("args[\"name\"] = %v, want %v", string(nameBytes), name)
	}
}

// TestSendHomunculusCommand tests the SendHomunculusCommand method
func TestSendHomunculusCommand(t *testing.T) {
	mockSend := NewMockSendForHomunculus()
	mockSend.packetLUT["homunculus_command"] = "0232"

	hm := NewHomunculusManager(mockSend)

	// Test sending a homunculus command
	command := uint8(1) // 0:get stats, 1:feed or 2:fire
	commandType := uint8(0)
	err := hm.SendHomunculusCommand(command, commandType)
	if err != nil {
		t.Fatalf("SendHomunculusCommand() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0232"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["commandID"] != command {
		t.Errorf("args[\"commandID\"] = %v, want %v", args["commandID"], command)
	}

	if args["commandType"] != commandType {
		t.Errorf("args[\"commandType\"] = %v, want %v", args["commandType"], commandType)
	}
}
