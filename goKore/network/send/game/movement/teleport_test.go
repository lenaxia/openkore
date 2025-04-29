package movement

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSend is a mock implementation of the Send interface for testing
type MockSend struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSend() *MockSend {
	return &MockSend{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSend) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSend) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSend) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSend) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSend) InjectMessage(message string) error {
	return nil
}

func (ms *MockSend) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSend) SendRaw(raw string) error {
	return nil
}

func (ms *MockSend) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSend) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSend) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSend) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSend) SetConnection(conn interface{}) {
}

func (ms *MockSend) GetConnection() interface{} {
	return nil
}

func (ms *MockSend) GetTime() uint32 {
	return 12345
}

// TestNewTeleportManager tests the NewTeleportManager function
func TestNewTeleportManager(t *testing.T) {
	mockSend := NewMockSend()
	tm := NewTeleportManager(mockSend)

	if tm == nil {
		t.Fatal("NewTeleportManager() returned nil")
	}

	if tm.baseSend == nil {
		t.Error("tm.baseSend was not set correctly")
	}
}

// TestSendWarpTele tests the SendWarpTele function
func TestSendWarpTele(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetLUT["warp_select"] = "011B"
	tm := NewTeleportManager(mockSend)

	// Test sending warp teleport command
	skillID := 26 // Teleport
	mapName := "prontera"
	err := tm.SendWarpTele(skillID, mapName)
	if err != nil {
		t.Fatalf("SendWarpTele() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["011B"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["skillID"] != skillID {
		t.Errorf("args[\"skillID\"] = %v, want %v", args["skillID"], skillID)
	}

	if args["mapName"] != mapName {
		t.Errorf("args[\"mapName\"] = %v, want %v", args["mapName"], mapName)
	}
}

// TestSendPrivateAirshipRequest tests the SendPrivateAirshipRequest function
func TestSendPrivateAirshipRequest(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.packetLUT["private_airship_request"] = "0A18"
	tm := NewTeleportManager(mockSend)

	// Test sending private airship request
	mapName := "izlude"
	nameID := 123
	err := tm.SendPrivateAirshipRequest(mapName, nameID)
	if err != nil {
		t.Fatalf("SendPrivateAirshipRequest() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A18"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["map_name"] != mapName {
		t.Errorf("args[\"map_name\"] = %v, want %v", args["map_name"], mapName)
	}

	if args["nameID"] != nameID {
		t.Errorf("args[\"nameID\"] = %v, want %v", args["nameID"], nameID)
	}
}
