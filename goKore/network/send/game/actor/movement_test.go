package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForMovement is a mock implementation of the Send interface for testing movement functionality
type MockSendForMovement struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForMovement() *MockSendForMovement {
	return &MockSendForMovement{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForMovement) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForMovement) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForMovement) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForMovement) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForMovement) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForMovement) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForMovement) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForMovement) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForMovement) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForMovement) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForMovement) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForMovement) SetConnection(conn interface{}) {
}

func (ms *MockSendForMovement) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForMovement) GetTime() uint32 {
	return 12345
}

// TestNewMovementManager tests the NewMovementManager function
func TestNewMovementManager(t *testing.T) {
	mockSend := NewMockSendForMovement()
	mm := NewMovementManager(mockSend)

	if mm == nil {
		t.Fatal("NewMovementManager() returned nil")
	}

	if mm.baseSend == nil {
		t.Error("mm.baseSend was not set correctly")
	}
}

// TestSendMove tests the SendMove method
func TestSendMove(t *testing.T) {
	mockSend := NewMockSendForMovement()
	mockSend.packetLUT["move_to"] = "0085"

	mm := NewMovementManager(mockSend)

	// Test sending move command
	x := 150
	y := 100
	err := mm.SendMove(x, y)
	if err != nil {
		t.Fatalf("SendMove() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0085"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["x"] != x {
		t.Errorf("args[\"x\"] = %v, want %v", args["x"], x)
	}

	if args["y"] != y {
		t.Errorf("args[\"y\"] = %v, want %v", args["y"], y)
	}
}

// TestSendSlaveMove tests the SendSlaveMove method
func TestSendSlaveMove(t *testing.T) {
	mockSend := NewMockSendForMovement()
	mockSend.packetLUT["slave_move_to"] = "0232"

	mm := NewMovementManager(mockSend)

	// Test sending slave move command
	slaveID := uint32(12345)
	x := 150
	y := 100
	err := mm.SendSlaveMove(slaveID, x, y)
	if err != nil {
		t.Fatalf("SendSlaveMove() returned error: %v", err)
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

	if args["slave_id"] != slaveID {
		t.Errorf("args[\"slave_id\"] = %v, want %v", args["slave_id"], slaveID)
	}

	if args["x"] != x {
		t.Errorf("args[\"x\"] = %v, want %v", args["x"], x)
	}

	if args["y"] != y {
		t.Errorf("args[\"y\"] = %v, want %v", args["y"], y)
	}
}

// TestSendActorMove tests the SendActorMove method
func TestSendActorMove(t *testing.T) {
	mockSend := NewMockSendForMovement()
	mockSend.packetLUT["actor_move_to"] = "0437"

	mm := NewMovementManager(mockSend)

	// Test sending actor move command
	actorID := uint32(67890)
	x := 150
	y := 100
	err := mm.SendActorMove(actorID, x, y)
	if err != nil {
		t.Fatalf("SendActorMove() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0437"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["actor_id"] != actorID {
		t.Errorf("args[\"actor_id\"] = %v, want %v", args["actor_id"], actorID)
	}

	if args["x"] != x {
		t.Errorf("args[\"x\"] = %v, want %v", args["x"], x)
	}

	if args["y"] != y {
		t.Errorf("args[\"y\"] = %v, want %v", args["y"], y)
	}
}
