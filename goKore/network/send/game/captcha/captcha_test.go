package captcha

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
			"captcha_answer":          "07E7",
			"captcha_preview_request": "07E5",
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

// TestNewCaptchaManager tests the NewCaptchaManager function
func TestNewCaptchaManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	captchaManager := NewCaptchaManager(mockSend)

	if captchaManager == nil {
		t.Fatal("NewCaptchaManager() returned nil")
	}

	if captchaManager.baseSend == nil {
		t.Error("captchaManager.baseSend was not set correctly")
	}
}

// TestSendCaptchaAnswer tests the SendCaptchaAnswer function
func TestSendCaptchaAnswer(t *testing.T) {
	mockSend := NewMockSend()
	captchaManager := NewCaptchaManager(mockSend)

	// Test sending a captcha answer
	answer := "test answer"
	err := captchaManager.SendCaptchaAnswer(answer)
	if err != nil {
		t.Fatalf("SendCaptchaAnswer() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "captcha_answer" {
		t.Errorf("Expected packet ID 'captcha_answer', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if accountIDVal, ok := mockSend.LastArgs()["accountID"].(uint32); !ok || accountIDVal != uint32(12345) {
		t.Errorf("Expected accountID=%d, got %v", uint32(12345), mockSend.LastArgs()["accountID"])
	}

	if answerVal, ok := mockSend.LastArgs()["answer"].(string); !ok || answerVal != answer {
		t.Errorf("Expected answer=%s, got %v", answer, mockSend.LastArgs()["answer"])
	}

	// Check that the len field was set
	if _, ok := mockSend.LastArgs()["len"].(int); !ok {
		t.Errorf("Expected len to be set, got %v", mockSend.LastArgs()["len"])
	}
}

// TestSendCaptchaPreviewRequest tests the SendCaptchaPreviewRequest function
func TestSendCaptchaPreviewRequest(t *testing.T) {
	mockSend := NewMockSend()
	captchaManager := NewCaptchaManager(mockSend)

	// Test sending a captcha preview request
	captchaKey := uint32(12345)
	err := captchaManager.SendCaptchaPreviewRequest(captchaKey)
	if err != nil {
		t.Fatalf("SendCaptchaPreviewRequest() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "captcha_preview_request" {
		t.Errorf("Expected packet ID 'captcha_preview_request', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if keyVal, ok := mockSend.LastArgs()["captcha_key"].(uint32); !ok || keyVal != captchaKey {
		t.Errorf("Expected captcha_key=%d, got %v", captchaKey, mockSend.LastArgs()["captcha_key"])
	}
}
