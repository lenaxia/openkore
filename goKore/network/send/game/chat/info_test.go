package chat

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
)

// MockSend is a mock implementation of the core.Send interface for testing
type MockInfoSend struct {
	packetIDs      map[string]string
	reconstructed  []byte
	sent           []byte
	time           uint32
	lastPacketName string
	lastArgs       map[string]interface{}
}

// NewMockInfoSend creates a new MockInfoSend instance with default values
func NewMockInfoSend() *MockInfoSend {
	return &MockInfoSend{
		packetIDs: map[string]string{
			"actor_info_request": "0A90", // Changed from 0094 to avoid conflicts
			"actor_name_request": "0095",
			"request_user_count": "01C1",
			"battleground_chat":  "02DB",
			"clan_chat":          "0B01",
		},
		time:     12345,
		lastArgs: make(map[string]interface{}),
	}
}

// SendToServer mocks sending a packet to the server
func (ms *MockInfoSend) SendToServer(msg []byte) error {
	ms.sent = msg
	return nil
}

// EncryptMessageID mocks encrypting a message ID
func (ms *MockInfoSend) EncryptMessageID(msg *[]byte) error {
	return nil
}

// CryptKeys mocks setting encryption keys
func (ms *MockInfoSend) CryptKeys(key1, key2, key3 uint32) {}

// PinEncode mocks encoding a PIN
func (ms *MockInfoSend) PinEncode(seed, pin int) string {
	return ""
}

// InjectMessage mocks injecting a message
func (ms *MockInfoSend) InjectMessage(message string) error {
	return nil
}

// InjectAdminMessage mocks injecting an admin message
func (ms *MockInfoSend) InjectAdminMessage(message string) error {
	return nil
}

// SendRaw mocks sending a raw packet
func (ms *MockInfoSend) SendRaw(raw string) error {
	return nil
}

// Reconstruct mocks reconstructing a packet
func (ms *MockInfoSend) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
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
func (ms *MockInfoSend) GetPacketID(name string) (string, bool) {
	id, ok := ms.packetIDs[name]
	if ok {
		ms.lastPacketName = name
	}
	return id, ok
}

// RegisterPacketHandler mocks registering a packet handler
func (ms *MockInfoSend) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
}

// RegisterHook mocks registering a hook
func (ms *MockInfoSend) RegisterHook(hookName string, callback hooks.HookCallback) {}

// SetConnection mocks setting a connection
func (ms *MockInfoSend) SetConnection(conn interface{}) {}

// GetConnection mocks getting a connection
func (ms *MockInfoSend) GetConnection() interface{} {
	return nil
}

// GetTime mocks getting the current time
func (ms *MockInfoSend) GetTime() uint32 {
	return ms.time
}

// LastPacketID returns the name of the last packet that was requested
func (ms *MockInfoSend) LastPacketID() (string, bool) {
	if ms.lastPacketName == "" {
		return "", false
	}
	return ms.lastPacketName, true
}

// LastArgs returns the arguments of the last packet that was reconstructed
func (ms *MockInfoSend) LastArgs() map[string]interface{} {
	return ms.lastArgs
}

// TestNewInfoChatManager tests the NewInfoChatManager function
func TestNewInfoChatManager(t *testing.T) {
	// Verify that MockInfoSend implements core.Send
	var _ core.Send = &MockInfoSend{}
	mockSend := NewMockInfoSend()
	infoChatManager := NewInfoChatManager(mockSend)

	if infoChatManager == nil {
		t.Fatal("NewInfoChatManager() returned nil")
	}

	if infoChatManager.baseSend == nil {
		t.Error("infoChatManager.baseSend was not set correctly")
	}
}

// TestSendGetPlayerInfo tests the SendGetPlayerInfo function
func TestSendGetPlayerInfo(t *testing.T) {
	mockSend := NewMockInfoSend()
	infoChatManager := NewInfoChatManager(mockSend)

	// Test sending a get player info request
	ID := uint32(12345)
	err := infoChatManager.SendGetPlayerInfo(ID)
	if err != nil {
		t.Fatalf("SendGetPlayerInfo() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "actor_info_request" {
		t.Errorf("Expected packet ID 'actor_info_request', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != ID {
		t.Errorf("Expected ID=%d, got %v", ID, mockSend.LastArgs()["ID"])
	}
}

// TestSendGetCharacterName tests the SendGetCharacterName function
func TestSendGetCharacterName(t *testing.T) {
	mockSend := NewMockInfoSend()
	infoChatManager := NewInfoChatManager(mockSend)

	// Test sending a get character name request
	ID := uint32(12345)
	err := infoChatManager.SendGetCharacterName(ID)
	if err != nil {
		t.Fatalf("SendGetCharacterName() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "actor_name_request" {
		t.Errorf("Expected packet ID 'actor_name_request', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != ID {
		t.Errorf("Expected ID=%d, got %v", ID, mockSend.LastArgs()["ID"])
	}
}

// TestSendWho tests the SendWho function
func TestSendWho(t *testing.T) {
	mockSend := NewMockInfoSend()
	infoChatManager := NewInfoChatManager(mockSend)

	// Test sending a who request
	err := infoChatManager.SendWho()
	if err != nil {
		t.Fatalf("SendWho() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "request_user_count" {
		t.Errorf("Expected packet ID 'request_user_count', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendBattlegroundChat tests the SendBattlegroundChat function
func TestSendBattlegroundChat(t *testing.T) {
	mockSend := NewMockInfoSend()
	infoChatManager := NewInfoChatManager(mockSend)

	// Test sending a battleground chat message
	message := "Hello, battleground!"
	err := infoChatManager.SendBattlegroundChat(message)
	if err != nil {
		t.Fatalf("SendBattlegroundChat() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "battleground_chat" {
		t.Errorf("Expected packet ID 'battleground_chat', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if messageVal, ok := mockSend.LastArgs()["message"].([]byte); !ok {
		t.Errorf("Expected message to be []byte, got %T", mockSend.LastArgs()["message"])
	} else {
		messageStr := string(messageVal)
		if messageStr != message {
			t.Errorf("Expected message=%s, got %s", message, messageStr)
		}
	}
}

// TestSendClanChat tests the SendClanChat function
func TestSendClanChat(t *testing.T) {
	mockSend := NewMockInfoSend()
	infoChatManager := NewInfoChatManager(mockSend)

	// Test sending a clan chat message
	message := "Hello, clan!"
	charName := "TestChar"
	err := infoChatManager.SendClanChat(message, charName)
	if err != nil {
		t.Fatalf("SendClanChat() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "clan_chat" {
		t.Errorf("Expected packet ID 'clan_chat', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	expectedMessage := charName + " : " + message
	if messageVal, ok := mockSend.LastArgs()["message"].(string); !ok || messageVal != expectedMessage {
		t.Errorf("Expected message=%s, got %v", expectedMessage, mockSend.LastArgs()["message"])
	}

	if lenVal, ok := mockSend.LastArgs()["len"].(uint16); !ok || lenVal != uint16(len(expectedMessage)+4) {
		t.Errorf("Expected len=%d, got %v", len(expectedMessage)+4, mockSend.LastArgs()["len"])
	}
}
