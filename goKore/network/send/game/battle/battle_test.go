package battle

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
)

// MockSend is a mock implementation of the core.Send interface for testing
type MockSend struct {
	packetIDs      map[string]string
	reconstructed  []byte
	sent           []byte
	time           uint32
	lastPacketName string
	lastArgs       map[string]interface{}
}

// NewMockSend creates a new MockSend instance with default values
func NewMockSend() *MockSend {
	return &MockSend{
		packetIDs: map[string]string{
			"view_player_equip_request": "02D6",
			"send_emotion":              "00BF",
			"novice_dori_dori":          "01B9",
			"novice_explosion_spirits":  "01BA",
			"memorial_dungeon_command":  "02CF",
		},
		time:     12345,
		lastArgs: make(map[string]interface{}),
	}
}

// SendToServer mocks sending a packet to the server
func (ms *MockSend) SendToServer(msg []byte) error {
	ms.sent = msg
	return nil
}

// EncryptMessageID mocks encrypting a message ID
func (ms *MockSend) EncryptMessageID(msg *[]byte) error {
	return nil
}

// CryptKeys mocks setting encryption keys
func (ms *MockSend) CryptKeys(key1, key2, key3 uint32) {}

// PinEncode mocks encoding a PIN
func (ms *MockSend) PinEncode(seed, pin int) string {
	return ""
}

// InjectMessage mocks injecting a message
func (ms *MockSend) InjectMessage(message string) error {
	return nil
}

// InjectAdminMessage mocks injecting an admin message
func (ms *MockSend) InjectAdminMessage(message string) error {
	return nil
}

// SendRaw mocks sending a raw packet
func (ms *MockSend) SendRaw(raw string) error {
	return nil
}

// Reconstruct mocks reconstructing a packet
func (ms *MockSend) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the last packet name and arguments for testing
	for name, id := range ms.packetIDs {
		if id == packetID {
			ms.lastPacketName = name
			break
		}
	}

	// Store the arguments for testing
	ms.lastArgs = args

	// Simple mock implementation that just returns the packet ID as bytes
	ms.reconstructed = []byte{0x00, 0x00}
	return ms.reconstructed, nil
}

// GetPacketID mocks getting a packet ID by name
func (ms *MockSend) GetPacketID(name string) (string, bool) {
	id, ok := ms.packetIDs[name]
	if ok {
		ms.lastPacketName = name
	}
	return id, ok
}

// RegisterPacketHandler mocks registering a packet handler
func (ms *MockSend) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
}

// RegisterHook mocks registering a hook
func (ms *MockSend) RegisterHook(hookName string, callback hooks.HookCallback) {}

// SetConnection mocks setting a connection
func (ms *MockSend) SetConnection(conn interface{}) {}

// GetConnection mocks getting a connection
func (ms *MockSend) GetConnection() interface{} {
	return nil
}

// GetTime mocks getting the current time
func (ms *MockSend) GetTime() uint32 {
	return ms.time
}

// LastPacketID returns the name of the last packet that was requested
func (ms *MockSend) LastPacketID() (string, bool) {
	if ms.lastPacketName == "" {
		return "", false
	}
	return ms.lastPacketName, true
}

// LastArgs returns the arguments of the last packet that was reconstructed
func (ms *MockSend) LastArgs() map[string]interface{} {
	return ms.lastArgs
}

// TestNewBattleManager tests the NewBattleManager function
func TestNewBattleManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	battleManager := NewBattleManager(mockSend)

	if battleManager == nil {
		t.Fatal("NewBattleManager() returned nil")
	}

	if battleManager.baseSend == nil {
		t.Error("battleManager.baseSend was not set correctly")
	}
}

// TestSendShowEquipPlayer tests the SendShowEquipPlayer function
func TestSendShowEquipPlayer(t *testing.T) {
	mockSend := NewMockSend()
	battleManager := NewBattleManager(mockSend)

	// Test sending a show equip player request
	playerID := uint32(12345)
	err := battleManager.SendShowEquipPlayer(playerID)
	if err != nil {
		t.Fatalf("SendShowEquipPlayer() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "view_player_equip_request" {
		t.Errorf("Expected packet ID 'view_player_equip_request', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != playerID {
		t.Errorf("Expected ID=%d, got %v", playerID, mockSend.LastArgs()["ID"])
	}
}

// TestSendEmotion tests the SendEmotion function
func TestSendEmotion(t *testing.T) {
	mockSend := NewMockSend()
	battleManager := NewBattleManager(mockSend)

	// Test sending an emotion
	emotionID := uint8(7) // Example emotion ID
	err := battleManager.SendEmotion(emotionID)
	if err != nil {
		t.Fatalf("SendEmotion() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "send_emotion" {
		t.Errorf("Expected packet ID 'send_emotion', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint8); !ok || idVal != emotionID {
		t.Errorf("Expected ID=%d, got %v", emotionID, mockSend.LastArgs()["ID"])
	}
}

// TestSendNoviceDoriDori tests the SendNoviceDoriDori function
func TestSendNoviceDoriDori(t *testing.T) {
	mockSend := NewMockSend()
	battleManager := NewBattleManager(mockSend)

	// Test sending a novice dori dori request
	err := battleManager.SendNoviceDoriDori()
	if err != nil {
		t.Fatalf("SendNoviceDoriDori() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "novice_dori_dori" {
		t.Errorf("Expected packet ID 'novice_dori_dori', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendNoviceExplosionSpirits tests the SendNoviceExplosionSpirits function
func TestSendNoviceExplosionSpirits(t *testing.T) {
	mockSend := NewMockSend()
	battleManager := NewBattleManager(mockSend)

	// Test sending a novice explosion spirits request
	err := battleManager.SendNoviceExplosionSpirits()
	if err != nil {
		t.Fatalf("SendNoviceExplosionSpirits() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "novice_explosion_spirits" {
		t.Errorf("Expected packet ID 'novice_explosion_spirits', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendMemorialDungeonCommand tests the SendMemorialDungeonCommand function
func TestSendMemorialDungeonCommand(t *testing.T) {
	mockSend := NewMockSend()
	battleManager := NewBattleManager(mockSend)

	// Test sending a memorial dungeon command
	command := uint32(1) // Example command
	err := battleManager.SendMemorialDungeonCommand(command)
	if err != nil {
		t.Fatalf("SendMemorialDungeonCommand() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "memorial_dungeon_command" {
		t.Errorf("Expected packet ID 'memorial_dungeon_command', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if cmdVal, ok := mockSend.LastArgs()["command"].(uint32); !ok || cmdVal != command {
		t.Errorf("Expected command=%d, got %v", command, mockSend.LastArgs()["command"])
	}
}
