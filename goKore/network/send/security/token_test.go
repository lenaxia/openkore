package security

import (
	"bytes"
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForToken is a mock implementation of the Send interface for testing token functionality
type MockSendForToken struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForToken() *MockSendForToken {
	return &MockSendForToken{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForToken) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForToken) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForToken) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForToken) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForToken) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForToken) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForToken) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForToken) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForToken) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForToken) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForToken) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForToken) SetConnection(conn interface{}) {
}

func (ms *MockSendForToken) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForToken) GetTime() uint32 {
	return 12345
}

// TestNewTokenManager tests the NewTokenManager function
func TestNewTokenManager(t *testing.T) {
	mockSend := NewMockSendForToken()
	tm := NewTokenManager(mockSend)

	if tm == nil {
		t.Fatal("NewTokenManager() returned nil")
	}

	if tm.baseSend == nil {
		t.Error("tm.baseSend was not set correctly")
	}
}

// TestSendTokenLogin tests the SendTokenLogin method
func TestSendTokenLogin(t *testing.T) {
	mockSend := NewMockSendForToken()
	mockSend.packetLUT["token_login"] = "0825"

	tm := NewTokenManager(mockSend)

	// Test sending token login
	username := "testuser"
	token := []byte{1, 2, 3, 4, 5}
	mac := "AABBCCDDEEFF"
	ip := "192.168.1.1"
	version := 23
	masterVersion := 1
	err := tm.SendTokenLogin(username, token, mac, ip, version, masterVersion)
	if err != nil {
		t.Fatalf("SendTokenLogin() returned error: %v", err)
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

	if args["username"] != username {
		t.Errorf("args[\"username\"] = %v, want %v", args["username"], username)
	}

	if !bytes.Equal(args["token"].([]byte), token) {
		t.Errorf("args[\"token\"] = %v, want %v", args["token"], token)
	}

	if args["mac"] != mac {
		t.Errorf("args[\"mac\"] = %v, want %v", args["mac"], mac)
	}

	if args["ip"] != ip {
		t.Errorf("args[\"ip\"] = %v, want %v", args["ip"], ip)
	}

	if args["version"] != version {
		t.Errorf("args[\"version\"] = %v, want %v", args["version"], version)
	}

	if args["master_version"] != masterVersion {
		t.Errorf("args[\"master_version\"] = %v, want %v", args["master_version"], masterVersion)
	}
}

// TestSendSecureLogin tests the SendSecureLogin method
func TestSendSecureLogin(t *testing.T) {
	mockSend := NewMockSendForToken()
	mockSend.packetLUT["secure_login"] = "0277"

	tm := NewTokenManager(mockSend)

	// Test sending secure login
	username := "testuser"
	password := "testpass"
	salt := []byte{1, 2, 3, 4}
	loginType := 1
	err := tm.SendSecureLogin(username, password, salt, loginType)
	if err != nil {
		t.Fatalf("SendSecureLogin() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0277"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["username"] != username {
		t.Errorf("args[\"username\"] = %v, want %v", args["username"], username)
	}

	// The password should be hashed with the salt
	if args["password_hash"] == nil {
		t.Error("args[\"password_hash\"] is nil")
	}
}
