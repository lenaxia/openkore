package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
)

// TestNewSakraySend tests the NewSakraySend function
func TestNewSakraySend(t *testing.T) {
	// Create a base send
	hookManager := hooks.NewHookManager()
	baseSend := NewBaseSend(hookManager)

	// Configure the base send
	packetConstructions := make(map[string]common.PacketConstruction)
	baseSend.Configure("Sakray", packetConstructions)

	// Set encryption keys
	baseSend.SetEncryptionKeys(0x12345678, 0x87654321, 0x11223344)

	// Create a Sakray send
	sakraySend := NewSakraySend(baseSend)

	if sakraySend == nil {
		t.Fatal("NewSakraySend() returned nil")
	}

	if sakraySend.baseSend == nil {
		t.Error("sakraySend.baseSend was not set correctly")
	}
}

// TestSakraySendToServer tests the SendToServer method
func TestSakraySendToServer(t *testing.T) {
	// Create a mock base send
	mockBaseSend := &MockBaseSend{
		sentPackets: make([][]byte, 0),
	}

	// Create a Sakray send
	sakraySend := NewSakraySend(mockBaseSend)

	// Test sending a packet
	packet := []byte{0x01, 0x02, 0x03, 0x04}
	err := sakraySend.SendToServer(packet)
	if err != nil {
		t.Fatalf("SendToServer() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockBaseSend.sentPackets) != 1 {
		t.Fatalf("len(mockBaseSend.sentPackets) = %v, want 1", len(mockBaseSend.sentPackets))
	}

	// Check that the packet was not modified
	if len(mockBaseSend.sentPackets[0]) != len(packet) {
		t.Fatalf("len(mockBaseSend.sentPackets[0]) = %v, want %v", len(mockBaseSend.sentPackets[0]), len(packet))
	}

	for i := 0; i < len(packet); i++ {
		if mockBaseSend.sentPackets[0][i] != packet[i] {
			t.Errorf("mockBaseSend.sentPackets[0][%d] = %v, want %v", i, mockBaseSend.sentPackets[0][i], packet[i])
		}
	}
}

// TestSakrayEncryptMessageID tests the EncryptMessageID method
func TestSakrayEncryptMessageID(t *testing.T) {
	// Create a mock base send
	mockBaseSend := &MockBaseSend{
		sentPackets: make([][]byte, 0),
	}

	// Create a Sakray send
	sakraySend := NewSakraySend(mockBaseSend)

	// Test encrypting a packet
	packet := []byte{0x01, 0x02, 0x03, 0x04}
	packetCopy := make([]byte, len(packet))
	copy(packetCopy, packet)

	err := sakraySend.EncryptMessageID(&packetCopy)
	if err != nil {
		t.Fatalf("EncryptMessageID() returned error: %v", err)
	}

	// Check that the packet was encrypted (should be delegated to base send)
	if mockBaseSend.encryptCalled != 1 {
		t.Errorf("mockBaseSend.encryptCalled = %v, want 1", mockBaseSend.encryptCalled)
	}
}

// MockBaseSend is a mock implementation of the Send interface for testing
type MockBaseSend struct {
	sentPackets   [][]byte
	encryptCalled int
}

func (mbs *MockBaseSend) SendToServer(msg []byte) error {
	mbs.sentPackets = append(mbs.sentPackets, msg)
	return nil
}

func (mbs *MockBaseSend) EncryptMessageID(msg *[]byte) error {
	mbs.encryptCalled++
	return nil
}

func (mbs *MockBaseSend) CryptKeys(key1, key2, key3 uint32) {
}

func (mbs *MockBaseSend) PinEncode(seed, pin int) string {
	return "1234"
}

func (mbs *MockBaseSend) InjectMessage(message string) error {
	return nil
}

func (mbs *MockBaseSend) InjectAdminMessage(message string) error {
	return nil
}

func (mbs *MockBaseSend) SendRaw(raw string) error {
	return nil
}

func (mbs *MockBaseSend) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	return []byte{0x00, 0x00}, nil
}

func (mbs *MockBaseSend) GetPacketID(name string) (string, bool) {
	return "", false
}

func (mbs *MockBaseSend) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
}

func (mbs *MockBaseSend) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (mbs *MockBaseSend) SetConnection(conn interface{}) {
}

func (mbs *MockBaseSend) GetConnection() interface{} {
	return nil
}

func (mbs *MockBaseSend) GetTime() uint32 {
	return 12345
}
