package security

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForPin is a mock implementation of the core.Send interface for testing PIN functionality
type MockSendForPin struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
	pinEncodeResult string
}

func NewMockSendForPin() *MockSendForPin {
	return &MockSendForPin{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
		pinEncodeResult: "1234",
	}
}

func (ms *MockSendForPin) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForPin) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForPin) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForPin) PinEncode(seed, pin int) string {
	return ms.pinEncodeResult
}

func (ms *MockSendForPin) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForPin) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForPin) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForPin) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForPin) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForPin) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForPin) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForPin) SetConnection(conn interface{}) {
}

func (ms *MockSendForPin) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForPin) GetTime() uint32 {
	return 12345
}

// TestNewPINManager tests the NewPINManager function
func TestNewPINManager(t *testing.T) {
	mockSend := NewMockSendForPin()
	pm := NewPINManager(mockSend)

	if pm == nil {
		t.Fatal("NewPINManager() returned nil")
	}

	if pm.baseSend == nil {
		t.Error("pm.baseSend was not set correctly")
	}
}

// TestSendPINCode tests the SendPINCode method
func TestSendPINCode(t *testing.T) {
	mockSend := NewMockSendForPin()
	mockSend.packetLUT["send_pin_code"] = "0292"

	pm := NewPINManager(mockSend)

	// Test sending PIN code
	pin := 1234
	seed := 5678
	err := pm.SendPINCode(pin, seed)
	if err != nil {
		t.Fatalf("SendPINCode() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0292"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	// The PIN should be encoded
	if args["pin"] != "1234" {
		t.Errorf("args[\"pin\"] = %v, want 1234", args["pin"])
	}
}

// TestSendPINCodeState tests the SendPINCodeState method
func TestSendPINCodeState(t *testing.T) {
	mockSend := NewMockSendForPin()
	mockSend.packetLUT["send_pin_state"] = "08B8"

	pm := NewPINManager(mockSend)

	// Test sending PIN code state
	state := 0 // 0 = disabled, 1 = enabled, 2 = requested
	err := pm.SendPINCodeState(state)
	if err != nil {
		t.Fatalf("SendPINCodeState() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["08B8"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["state"] != state {
		t.Errorf("args[\"state\"] = %v, want %v", args["state"], state)
	}
}

// TestSendPINCodeChange tests the SendPINCodeChange method
func TestSendPINCodeChange(t *testing.T) {
	mockSend := NewMockSendForPin()
	mockSend.packetLUT["send_pin_change"] = "0362"

	pm := NewPINManager(mockSend)

	// Test sending PIN code change
	oldPin := 1234
	newPin := 5678
	seed := 9012
	err := pm.SendPINCodeChange(oldPin, newPin, seed)
	if err != nil {
		t.Fatalf("SendPINCodeChange() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0362"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	// Both PINs should be encoded
	if args["old_pin"] != "1234" {
		t.Errorf("args[\"old_pin\"] = %v, want 1234", args["old_pin"])
	}

	if args["new_pin"] != "1234" {
		t.Errorf("args[\"new_pin\"] = %v, want 1234", args["new_pin"])
	}
}

// TestSendLoginPinCode tests the SendLoginPinCode method
func TestSendLoginPinCode(t *testing.T) {
	mockSend := NewMockSendForPin()
	mockSend.packetLUT["send_pin_password"] = "0825"
	mockSend.packetLUT["new_pin_password"] = "0826"

	pm := NewPINManager(mockSend)

	// Test sending login PIN code (type 0)
	seed := 5678
	err := pm.SendLoginPinCode(seed, 0)
	if err != nil {
		t.Fatalf("SendLoginPinCode() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0825"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	// The PIN should be encoded
	if args["pin"] != "1234" {
		t.Errorf("args[\"pin\"] = %v, want 1234", args["pin"])
	}

	// Test sending new PIN code (type 1)
	err = pm.SendLoginPinCode(seed, 1)
	if err != nil {
		t.Fatalf("SendLoginPinCode() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 2 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 2", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists = mockSend.reconstructArgs["0826"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	// The PIN should be encoded
	if args["pin"] != "1234" {
		t.Errorf("args[\"pin\"] = %v, want 1234", args["pin"])
	}

	// Test sending invalid type
	err = pm.SendLoginPinCode(seed, 2)
	if err == nil {
		t.Fatal("SendLoginPinCode() did not return an error for invalid type")
	}
}
