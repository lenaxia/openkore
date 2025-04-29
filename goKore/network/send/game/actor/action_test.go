package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForAction is a mock implementation of the Send interface for testing action functionality
type MockSendForAction struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForAction() *MockSendForAction {
	return &MockSendForAction{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForAction) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForAction) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForAction) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForAction) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForAction) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForAction) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForAction) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForAction) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForAction) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForAction) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForAction) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForAction) SetConnection(conn interface{}) {
}

func (ms *MockSendForAction) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForAction) GetTime() uint32 {
	return 12345
}

// TestNewActionManager tests the NewActionManager function
func TestNewActionManager(t *testing.T) {
	mockSend := NewMockSendForAction()
	am := NewActionManager(mockSend)

	if am == nil {
		t.Fatal("NewActionManager() returned nil")
	}

	if am.baseSend == nil {
		t.Error("am.baseSend was not set correctly")
	}
}

// TestSendAction tests the SendAction method
func TestSendAction(t *testing.T) {
	mockSend := NewMockSendForAction()
	mockSend.packetLUT["send_action"] = "0089"

	am := NewActionManager(mockSend)

	// Test sending action command
	targetID := uint32(12345)
	action := 1 // 1 = attack
	err := am.SendAction(targetID, action)
	if err != nil {
		t.Fatalf("SendAction() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0089"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["target_id"] != targetID {
		t.Errorf("args[\"target_id\"] = %v, want %v", args["target_id"], targetID)
	}

	if args["action"] != action {
		t.Errorf("args[\"action\"] = %v, want %v", args["action"], action)
	}
}

// TestSendLook tests the SendLook method
func TestSendLook(t *testing.T) {
	mockSend := NewMockSendForAction()
	mockSend.packetLUT["send_look"] = "009B"

	am := NewActionManager(mockSend)

	// Test sending look command
	body := 1 // Body direction (0-7)
	head := 2 // Head direction (0-2)
	err := am.SendLook(body, head)
	if err != nil {
		t.Fatalf("SendLook() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["009B"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["body"] != body {
		t.Errorf("args[\"body\"] = %v, want %v", args["body"], body)
	}

	if args["head"] != head {
		t.Errorf("args[\"head\"] = %v, want %v", args["head"], head)
	}
}

// TestSendEmotion tests the SendEmotion method
func TestSendEmotion(t *testing.T) {
	mockSend := NewMockSendForAction()
	mockSend.packetLUT["send_emotion"] = "00BF"

	am := NewActionManager(mockSend)

	// Test sending emotion command
	emotion := 10 // Emotion ID (0-100)
	err := am.SendEmotion(emotion)
	if err != nil {
		t.Fatalf("SendEmotion() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00BF"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["emotion"] != emotion {
		t.Errorf("args[\"emotion\"] = %v, want %v", args["emotion"], emotion)
	}
}

// TestSendSit tests the SendSit method
func TestSendSit(t *testing.T) {
	mockSend := NewMockSendForAction()
	mockSend.packetLUT["send_sit"] = "0089"

	am := NewActionManager(mockSend)

	// Test sending sit command
	err := am.SendSit()
	if err != nil {
		t.Fatalf("SendSit() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0089"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["action"] != 2 { // 2 = sit
		t.Errorf("args[\"action\"] = %v, want 2", args["action"])
	}
}

// TestSendStand tests the SendStand method
func TestSendStand(t *testing.T) {
	mockSend := NewMockSendForAction()
	mockSend.packetLUT["send_stand"] = "0089"

	am := NewActionManager(mockSend)

	// Test sending stand command
	err := am.SendStand()
	if err != nil {
		t.Fatalf("SendStand() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0089"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["action"] != 3 { // 3 = stand
		t.Errorf("args[\"action\"] = %v, want 3", args["action"])
	}
}
