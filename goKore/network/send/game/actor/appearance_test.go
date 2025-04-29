package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForAppearance is a mock implementation of the Send interface for testing appearance functionality
type MockSendForAppearance struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForAppearance() *MockSendForAppearance {
	return &MockSendForAppearance{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForAppearance) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForAppearance) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForAppearance) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForAppearance) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForAppearance) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForAppearance) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForAppearance) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForAppearance) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForAppearance) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForAppearance) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForAppearance) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForAppearance) SetConnection(conn interface{}) {
}

func (ms *MockSendForAppearance) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForAppearance) GetTime() uint32 {
	return 12345
}

// TestNewAppearanceManager tests the NewAppearanceManager function
func TestNewAppearanceManager(t *testing.T) {
	mockSend := NewMockSendForAppearance()
	am := NewAppearanceManager(mockSend)

	if am == nil {
		t.Fatal("NewAppearanceManager() returned nil")
	}

	if am.baseSend == nil {
		t.Error("am.baseSend was not set correctly")
	}
}

// TestSendChangeClothes tests the SendChangeClothes method
func TestSendChangeClothes(t *testing.T) {
	mockSend := NewMockSendForAppearance()
	mockSend.packetLUT["change_clothes"] = "00A9"

	am := NewAppearanceManager(mockSend)

	// Test sending change clothes command
	headTop := 123    // Head top equipment
	headMid := 456    // Head mid equipment
	headBottom := 789 // Head bottom equipment
	err := am.SendChangeClothes(headTop, headMid, headBottom)
	if err != nil {
		t.Fatalf("SendChangeClothes() returned error: %v", err)
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

	if args["head_top"] != headTop {
		t.Errorf("args[\"head_top\"] = %v, want %v", args["head_top"], headTop)
	}

	if args["head_mid"] != headMid {
		t.Errorf("args[\"head_mid\"] = %v, want %v", args["head_mid"], headMid)
	}

	if args["head_bottom"] != headBottom {
		t.Errorf("args[\"head_bottom\"] = %v, want %v", args["head_bottom"], headBottom)
	}
}

// TestSendChangeHair tests the SendChangeHair method
func TestSendChangeHair(t *testing.T) {
	mockSend := NewMockSendForAppearance()
	mockSend.packetLUT["change_hair"] = "01B0"

	am := NewAppearanceManager(mockSend)

	// Test sending change hair command
	hairStyle := 5  // Hair style ID
	hairColor := 10 // Hair color ID
	err := am.SendChangeHair(hairStyle, hairColor)
	if err != nil {
		t.Fatalf("SendChangeHair() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["01B0"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["hair_style"] != hairStyle {
		t.Errorf("args[\"hair_style\"] = %v, want %v", args["hair_style"], hairStyle)
	}

	if args["hair_color"] != hairColor {
		t.Errorf("args[\"hair_color\"] = %v, want %v", args["hair_color"], hairColor)
	}
}

// TestSendChangeStat tests the SendChangeStat method
func TestSendChangeStat(t *testing.T) {
	mockSend := NewMockSendForAppearance()
	mockSend.packetLUT["change_stat"] = "00BB"

	am := NewAppearanceManager(mockSend)

	// Test sending change stat command
	statType := 13 // Stat type (e.g., STR, AGI, VIT, etc.)
	amount := 1    // Amount to increase
	err := am.SendChangeStat(statType, amount)
	if err != nil {
		t.Fatalf("SendChangeStat() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00BB"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["stat_type"] != statType {
		t.Errorf("args[\"stat_type\"] = %v, want %v", args["stat_type"], statType)
	}

	if args["amount"] != amount {
		t.Errorf("args[\"amount\"] = %v, want %v", args["amount"], amount)
	}
}
