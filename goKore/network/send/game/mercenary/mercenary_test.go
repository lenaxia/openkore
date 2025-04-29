package mercenary

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
			"mercenary_command": "0234",
			"companion_release": "02A5",
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

// TestNewMercenaryManager tests the NewMercenaryManager function
func TestNewMercenaryManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	mercenaryManager := NewMercenaryManager(mockSend)

	if mercenaryManager == nil {
		t.Fatal("NewMercenaryManager() returned nil")
	}

	if mercenaryManager.baseSend == nil {
		t.Error("mercenaryManager.baseSend was not set correctly")
	}
}

// TestSendMercenaryCommand tests the SendMercenaryCommand function
func TestSendMercenaryCommand(t *testing.T) {
	mockSend := NewMockSend()
	mercenaryManager := NewMercenaryManager(mockSend)

	// Test cases for different command flags
	testCases := []struct {
		command int
		name    string
	}{
		{0, "COMMAND_REQ_NONE"},
		{1, "COMMAND_REQ_PROPERTY"},
		{2, "COMMAND_REQ_DELETE"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := mercenaryManager.SendMercenaryCommand(tc.command)
			if err != nil {
				t.Fatalf("SendMercenaryCommand(%d) returned error: %v", tc.command, err)
			}

			if mockSend.sent == nil {
				t.Fatal("No packet was sent")
			}

			// Check that the correct packet ID was used
			if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "mercenary_command" {
				t.Errorf("Expected packet ID 'mercenary_command', got '%s'", packetID)
			}

			// Check that the correct arguments were used
			if flagVal, ok := mockSend.LastArgs()["flag"].(uint8); !ok || int(flagVal) != tc.command {
				t.Errorf("Expected flag=%d, got %v", tc.command, mockSend.LastArgs()["flag"])
			}
		})
	}

	// Test with invalid command
	err := mercenaryManager.SendMercenaryCommand(3)
	if err == nil {
		t.Error("Expected error for invalid command, got nil")
	}
}

// TestSendCompanionRelease tests the SendCompanionRelease function
func TestSendCompanionRelease(t *testing.T) {
	mockSend := NewMockSend()
	mercenaryManager := NewMercenaryManager(mockSend)

	// Test sending a companion release command
	err := mercenaryManager.SendCompanionRelease()
	if err != nil {
		t.Fatalf("SendCompanionRelease() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "companion_release" {
		t.Errorf("Expected packet ID 'companion_release', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}
