package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSend is a mock implementation of the Send interface for testing
type MockSend struct {
	sentPackets      [][]byte
	encryptedPackets [][]byte
	cryptKey1        uint32
	cryptKey2        uint32
	cryptKey3        uint32
	conn             interface{}
	hookManager      *hooks.HookManager
	packetLUT        map[string]string
}

func NewMockSend() *MockSend {
	return &MockSend{
		sentPackets:      make([][]byte, 0),
		encryptedPackets: make([][]byte, 0),
		hookManager:      hooks.NewHookManager(),
		packetLUT:        make(map[string]string),
	}
}

func (ms *MockSend) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSend) EncryptMessageID(msg *[]byte) error {
	ms.encryptedPackets = append(ms.encryptedPackets, *msg)
	return nil
}

func (ms *MockSend) CryptKeys(key1, key2, key3 uint32) {
	ms.cryptKey1 = key1
	ms.cryptKey2 = key2
	ms.cryptKey3 = key3
}

func (ms *MockSend) PinEncode(seed, pin int) string {
	// Simple mock implementation for testing
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
	// Simple mock implementation that returns a packet with the ID
	if len(packetID) != 4 {
		return nil, ErrInvalidPacketID
	}

	// Create a simple packet with the ID
	packet := []byte{0, 0}
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
	ms.hookManager.AddHook(hookName, callback, nil)
}

func (ms *MockSend) SetConnection(conn interface{}) {
	ms.conn = conn
}

func (ms *MockSend) GetConnection() interface{} {
	return ms.conn
}

func (ms *MockSend) GetTime() uint32 {
	return 12345
}

// TestSendConfig tests the SendConfig struct
func TestSendConfig(t *testing.T) {
	config := &SendConfig{
		ServerType:    "ServerType0",
		PacketVersion: 23,
		UseEncryption: true,
		CryptKey1:     0x12345678,
		CryptKey2:     0x87654321,
		CryptKey3:     0xABCDEF01,
	}

	if config.ServerType != "ServerType0" {
		t.Errorf("config.ServerType = %v, want ServerType0", config.ServerType)
	}

	if config.PacketVersion != 23 {
		t.Errorf("config.PacketVersion = %v, want 23", config.PacketVersion)
	}

	if !config.UseEncryption {
		t.Errorf("config.UseEncryption = %v, want true", config.UseEncryption)
	}

	if config.CryptKey1 != 0x12345678 {
		t.Errorf("config.CryptKey1 = %v, want 0x12345678", config.CryptKey1)
	}

	if config.CryptKey2 != 0x87654321 {
		t.Errorf("config.CryptKey2 = %v, want 0x87654321", config.CryptKey2)
	}

	if config.CryptKey3 != 0xABCDEF01 {
		t.Errorf("config.CryptKey3 = %v, want 0xABCDEF01", config.CryptKey3)
	}
}

// TestPacketDefinition tests the PacketDefinition struct
func TestPacketDefinition(t *testing.T) {
	handler := func(args map[string]interface{}) error {
		return nil
	}

	def := PacketDefinition{
		ID:      "0064",
		Name:    "master_login",
		Format:  "V Z24 Z24 C",
		Keys:    []string{"version", "username", "password", "master_version"},
		Handler: handler,
	}

	if def.ID != "0064" {
		t.Errorf("def.ID = %v, want 0064", def.ID)
	}

	if def.Name != "master_login" {
		t.Errorf("def.Name = %v, want master_login", def.Name)
	}

	if def.Format != "V Z24 Z24 C" {
		t.Errorf("def.Format = %v, want V Z24 Z24 C", def.Format)
	}

	if len(def.Keys) != 4 {
		t.Errorf("len(def.Keys) = %v, want 4", len(def.Keys))
	}

	if def.Keys[0] != "version" {
		t.Errorf("def.Keys[0] = %v, want version", def.Keys[0])
	}

	if def.Keys[1] != "username" {
		t.Errorf("def.Keys[1] = %v, want username", def.Keys[1])
	}

	if def.Keys[2] != "password" {
		t.Errorf("def.Keys[2] = %v, want password", def.Keys[2])
	}

	if def.Keys[3] != "master_version" {
		t.Errorf("def.Keys[3] = %v, want master_version", def.Keys[3])
	}

	if def.Handler == nil {
		t.Errorf("def.Handler is nil")
	}
}

// TestMockSend tests the MockSend implementation
func TestMockSend(t *testing.T) {
	mock := NewMockSend()

	// Test SendToServer
	err := mock.SendToServer([]byte{1, 2, 3, 4})
	if err != nil {
		t.Errorf("SendToServer() returned error: %v", err)
	}
	if len(mock.sentPackets) != 1 {
		t.Errorf("len(mock.sentPackets) = %v, want 1", len(mock.sentPackets))
	}
	if len(mock.sentPackets[0]) != 4 {
		t.Errorf("len(mock.sentPackets[0]) = %v, want 4", len(mock.sentPackets[0]))
	}

	// Test EncryptMessageID
	msg := []byte{5, 6, 7, 8}
	err = mock.EncryptMessageID(&msg)
	if err != nil {
		t.Errorf("EncryptMessageID() returned error: %v", err)
	}
	if len(mock.encryptedPackets) != 1 {
		t.Errorf("len(mock.encryptedPackets) = %v, want 1", len(mock.encryptedPackets))
	}
	if len(mock.encryptedPackets[0]) != 4 {
		t.Errorf("len(mock.encryptedPackets[0]) = %v, want 4", len(mock.encryptedPackets[0]))
	}

	// Test CryptKeys
	mock.CryptKeys(1, 2, 3)
	if mock.cryptKey1 != 1 {
		t.Errorf("mock.cryptKey1 = %v, want 1", mock.cryptKey1)
	}
	if mock.cryptKey2 != 2 {
		t.Errorf("mock.cryptKey2 = %v, want 2", mock.cryptKey2)
	}
	if mock.cryptKey3 != 3 {
		t.Errorf("mock.cryptKey3 = %v, want 3", mock.cryptKey3)
	}

	// Test PinEncode
	pin := mock.PinEncode(1234, 5678)
	if pin != "1234" {
		t.Errorf("PinEncode() = %v, want 1234", pin)
	}

	// Test RegisterPacketHandler and GetPacketID
	mock.RegisterPacketHandler("0064", "master_login", "V Z24 Z24 C", []string{"version", "username", "password", "master_version"}, nil)
	id, exists := mock.GetPacketID("master_login")
	if !exists {
		t.Errorf("GetPacketID() returned exists = false, want true")
	}
	if id != "0064" {
		t.Errorf("GetPacketID() = %v, want 0064", id)
	}

	// Test SetConnection and GetConnection
	conn := "test_connection"
	mock.SetConnection(conn)
	gotConn := mock.GetConnection()
	if gotConn != conn {
		t.Errorf("GetConnection() = %v, want %v", gotConn, conn)
	}
}
