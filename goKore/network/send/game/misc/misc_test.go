package misc

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
			"token_login":          "0825",
			"request_remain_time":  "0A37",
			"blocking_play_cancel": "0447",
			"recall_sso":           "0842",
			"remove_aid_sso":       "0843",
			"starplace_agree":      "0B0D",
			"sync_request_ex":      "09F1",
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
	return "encoded_pin"
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

// TestNewMiscManager tests the NewMiscManager function
func TestNewMiscManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	miscManager := NewMiscManager(mockSend)

	if miscManager == nil {
		t.Fatal("NewMiscManager() returned nil")
	}

	if miscManager.baseSend == nil {
		t.Error("miscManager.baseSend was not set correctly")
	}
}

// TestSendTokenToServer tests the SendTokenToServer function
func TestSendTokenToServer(t *testing.T) {
	mockSend := NewMockSend()
	miscManager := NewMiscManager(mockSend)

	// Test sending a token to server
	username := "testuser"
	password := "testpass"
	masterVersion := uint32(1)
	version := uint32(2)
	token := "testtoken"
	length := uint16(10)
	otpIP := "127.0.0.1"
	otpPort := uint16(6900)
	err := miscManager.SendTokenToServer(username, password, masterVersion, version, token, length, otpIP, otpPort)
	if err != nil {
		t.Fatalf("SendTokenToServer() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "token_login" {
		t.Errorf("Expected packet ID 'token_login', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if usernameVal, ok := mockSend.LastArgs()["username"].(string); !ok || usernameVal != username {
		t.Errorf("Expected username=%s, got %v", username, mockSend.LastArgs()["username"])
	}

	if passwordVal, ok := mockSend.LastArgs()["password"].(string); !ok || passwordVal != password {
		t.Errorf("Expected password=%s, got %v", password, mockSend.LastArgs()["password"])
	}

	if tokenVal, ok := mockSend.LastArgs()["token"].(string); !ok || tokenVal != token {
		t.Errorf("Expected token=%s, got %v", token, mockSend.LastArgs()["token"])
	}
}

// TestEncryptPassword tests the EncryptPassword function
func TestEncryptPassword(t *testing.T) {
	mockSend := NewMockSend()
	miscManager := NewMiscManager(mockSend)

	// Test encrypting a password
	password := "testpass"
	encryptedPassword := miscManager.EncryptPassword(password)

	// Since we can't easily test the actual encryption, just check that it returns something
	if encryptedPassword == "" {
		t.Error("EncryptPassword() returned empty string")
	}
}

// TestSendReqRemainTime tests the SendReqRemainTime function
func TestSendReqRemainTime(t *testing.T) {
	mockSend := NewMockSend()
	miscManager := NewMiscManager(mockSend)

	// Test sending a request for remain time
	err := miscManager.SendReqRemainTime()
	if err != nil {
		t.Fatalf("SendReqRemainTime() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "request_remain_time" {
		t.Errorf("Expected packet ID 'request_remain_time', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendBlockingPlayerCancel tests the SendBlockingPlayerCancel function
func TestSendBlockingPlayerCancel(t *testing.T) {
	mockSend := NewMockSend()
	miscManager := NewMiscManager(mockSend)

	// Test sending a blocking player cancel request
	err := miscManager.SendBlockingPlayerCancel()
	if err != nil {
		t.Fatalf("SendBlockingPlayerCancel() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "blocking_play_cancel" {
		t.Errorf("Expected packet ID 'blocking_play_cancel', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendRecallSso tests the SendRecallSso function
func TestSendRecallSso(t *testing.T) {
	mockSend := NewMockSend()
	miscManager := NewMiscManager(mockSend)

	// Test sending a recall SSO request
	accountID := uint32(12345)
	err := miscManager.SendRecallSso(accountID)
	if err != nil {
		t.Fatalf("SendRecallSso() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "recall_sso" {
		t.Errorf("Expected packet ID 'recall_sso', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != accountID {
		t.Errorf("Expected ID=%d, got %v", accountID, mockSend.LastArgs()["ID"])
	}
}

// TestSendRemoveAidSso tests the SendRemoveAidSso function
func TestSendRemoveAidSso(t *testing.T) {
	mockSend := NewMockSend()
	miscManager := NewMiscManager(mockSend)

	// Test sending a remove AID SSO request
	accountID := uint32(12345)
	err := miscManager.SendRemoveAidSso(accountID)
	if err != nil {
		t.Fatalf("SendRemoveAidSso() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "remove_aid_sso" {
		t.Errorf("Expected packet ID 'remove_aid_sso', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != accountID {
		t.Errorf("Expected ID=%d, got %v", accountID, mockSend.LastArgs()["ID"])
	}
}

// TestSendFeelSaveOk tests the SendFeelSaveOk function
func TestSendFeelSaveOk(t *testing.T) {
	mockSend := NewMockSend()
	miscManager := NewMiscManager(mockSend)

	// Test sending a feel save ok request
	flag := uint8(1)
	err := miscManager.SendFeelSaveOk(flag)
	if err != nil {
		t.Fatalf("SendFeelSaveOk() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "starplace_agree" {
		t.Errorf("Expected packet ID 'starplace_agree', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if flagVal, ok := mockSend.LastArgs()["flag"].(uint8); !ok || flagVal != flag {
		t.Errorf("Expected flag=%d, got %v", flag, mockSend.LastArgs()["flag"])
	}
}

// TestSendReplySyncRequestEx tests the SendReplySyncRequestEx function
func TestSendReplySyncRequestEx(t *testing.T) {
	mockSend := NewMockSend()
	miscManager := NewMiscManager(mockSend)

	// Test sending a reply sync request ex
	syncID := uint16(12345)
	err := miscManager.SendReplySyncRequestEx(syncID)
	if err != nil {
		t.Fatalf("SendReplySyncRequestEx() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "sync_request_ex" {
		t.Errorf("Expected packet ID 'sync_request_ex', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if syncIDVal, ok := mockSend.LastArgs()["syncID"].(uint16); !ok || syncIDVal != syncID {
		t.Errorf("Expected syncID=%d, got %v", syncID, mockSend.LastArgs()["syncID"])
	}
}
