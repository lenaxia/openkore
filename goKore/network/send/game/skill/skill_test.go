package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForSkill is a mock implementation of the Send interface for testing skill functionality
type MockSendForSkill struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
	hookCalled      bool
	hookArgs        map[string]interface{}
}

func NewMockSendForSkill() *MockSendForSkill {
	return &MockSendForSkill{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
		hookArgs:        make(map[string]interface{}),
	}
}

func (ms *MockSendForSkill) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForSkill) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForSkill) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForSkill) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForSkill) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForSkill) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForSkill) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForSkill) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForSkill) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForSkill) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForSkill) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForSkill) SetConnection(conn interface{}) {
}

func (ms *MockSendForSkill) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForSkill) GetTime() uint32 {
	return 12345
}

// TestNewSkillManager tests the NewSkillManager function
func TestNewSkillManager(t *testing.T) {
	mockSend := NewMockSendForSkill()
	sm := NewSkillManager(mockSend)

	if sm == nil {
		t.Fatal("NewSkillManager() returned nil")
	}

	if sm.baseSend == nil {
		t.Error("sm.baseSend was not set correctly")
	}
}

// TestSendSkillUse tests the SendSkillUse method
func TestSendSkillUse(t *testing.T) {
	mockSend := NewMockSendForSkill()
	mockSend.packetLUT["skill_use"] = "0113"

	sm := NewSkillManager(mockSend)

	// Test using a skill
	skillID := uint16(1)
	level := uint16(10)
	targetID := uint32(12345)
	err := sm.SendSkillUse(skillID, level, targetID)
	if err != nil {
		t.Fatalf("SendSkillUse() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0113"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["skillID"] != skillID {
		t.Errorf("args[\"skillID\"] = %v, want %v", args["skillID"], skillID)
	}

	if args["lv"] != level {
		t.Errorf("args[\"lv\"] = %v, want %v", args["lv"], level)
	}

	if args["targetID"] != targetID {
		t.Errorf("args[\"targetID\"] = %v, want %v", args["targetID"], targetID)
	}
}

// TestSendSkillUseLoc tests the SendSkillUseLoc method
func TestSendSkillUseLoc(t *testing.T) {
	mockSend := NewMockSendForSkill()
	mockSend.packetLUT["skill_use_location"] = "0116"

	sm := NewSkillManager(mockSend)

	// Test using a skill on a location
	skillID := uint16(1)
	level := uint16(10)
	x := uint16(100)
	y := uint16(200)
	err := sm.SendSkillUseLoc(skillID, level, x, y)
	if err != nil {
		t.Fatalf("SendSkillUseLoc() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0116"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["skillID"] != skillID {
		t.Errorf("args[\"skillID\"] = %v, want %v", args["skillID"], skillID)
	}

	if args["lv"] != level {
		t.Errorf("args[\"lv\"] = %v, want %v", args["lv"], level)
	}

	if args["x"] != x {
		t.Errorf("args[\"x\"] = %v, want %v", args["x"], x)
	}

	if args["y"] != y {
		t.Errorf("args[\"y\"] = %v, want %v", args["y"], y)
	}
}

// TestSendSkillUseLocInfo tests the SendSkillUseLocInfo method
func TestSendSkillUseLocInfo(t *testing.T) {
	mockSend := NewMockSendForSkill()
	mockSend.packetLUT["skill_use_location_text"] = "0190"

	sm := NewSkillManager(mockSend)

	// Test using a skill on a location with additional info
	skillID := uint16(1)
	level := uint16(10)
	x := uint16(100)
	y := uint16(200)
	info := "Additional info"
	err := sm.SendSkillUseLocInfo(skillID, level, x, y, info)
	if err != nil {
		t.Fatalf("SendSkillUseLocInfo() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0190"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != skillID {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], skillID)
	}

	if args["lvl"] != level {
		t.Errorf("args[\"lvl\"] = %v, want %v", args["lvl"], level)
	}

	if args["x"] != x {
		t.Errorf("args[\"x\"] = %v, want %v", args["x"], x)
	}

	if args["y"] != y {
		t.Errorf("args[\"y\"] = %v, want %v", args["y"], y)
	}

	if args["info"] != info {
		t.Errorf("args[\"info\"] = %v, want %v", args["info"], info)
	}
}

// TestSendSkillSelect tests the SendSkillSelect method
func TestSendSkillSelect(t *testing.T) {
	mockSend := NewMockSendForSkill()
	mockSend.packetLUT["skill_select"] = "0442"

	sm := NewSkillManager(mockSend)

	// Test selecting a skill
	skillID := uint16(1)
	why := uint16(2)
	err := sm.SendSkillSelect(skillID, why)
	if err != nil {
		t.Fatalf("SendSkillSelect() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0442"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["skillID"] != skillID {
		t.Errorf("args[\"skillID\"] = %v, want %v", args["skillID"], skillID)
	}

	if args["why"] != why {
		t.Errorf("args[\"why\"] = %v, want %v", args["why"], why)
	}
}

// TestSendStartSkillUse tests the SendStartSkillUse method
func TestSendStartSkillUse(t *testing.T) {
	mockSend := NewMockSendForSkill()
	mockSend.packetLUT["start_skill_use"] = "0439"

	sm := NewSkillManager(mockSend)

	// Test starting a continuous skill
	skillID := uint16(1)
	level := uint16(10)
	targetID := uint32(12345)
	err := sm.SendStartSkillUse(skillID, level, targetID)
	if err != nil {
		t.Fatalf("SendStartSkillUse() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0439"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["skillID"] != skillID {
		t.Errorf("args[\"skillID\"] = %v, want %v", args["skillID"], skillID)
	}

	if args["lv"] != level {
		t.Errorf("args[\"lv\"] = %v, want %v", args["lv"], level)
	}

	if args["targetID"] != targetID {
		t.Errorf("args[\"targetID\"] = %v, want %v", args["targetID"], targetID)
	}
}

// TestSendStopSkillUse tests the SendStopSkillUse method
func TestSendStopSkillUse(t *testing.T) {
	mockSend := NewMockSendForSkill()
	mockSend.packetLUT["stop_skill_use"] = "0440"

	sm := NewSkillManager(mockSend)

	// Test stopping a continuous skill
	skillID := uint16(1)
	err := sm.SendStopSkillUse(skillID)
	if err != nil {
		t.Fatalf("SendStopSkillUse() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0440"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["skillID"] != skillID {
		t.Errorf("args[\"skillID\"] = %v, want %v", args["skillID"], skillID)
	}
}

// TestSendAutoSpell tests the SendAutoSpell method
func TestSendAutoSpell(t *testing.T) {
	mockSend := NewMockSendForSkill()
	mockSend.packetLUT["auto_spell"] = "0443"

	sm := NewSkillManager(mockSend)

	// Test setting auto-spell
	skillID := uint16(1)
	err := sm.SendAutoSpell(skillID)
	if err != nil {
		t.Fatalf("SendAutoSpell() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0443"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != skillID {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], skillID)
	}
}
