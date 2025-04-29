package macro

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
			"macro_start":             "0871",
			"macro_stop":              "0872",
			"macro_detector_download": "0A5A",
			"macro_detector_answer":   "0A5C",
			"req_cash_tabcode":        "0A69", // Changed from 0A68 to avoid conflict with open_ui_request
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

// TestNewMacroManager tests the NewMacroManager function
func TestNewMacroManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	macroManager := NewMacroManager(mockSend)

	if macroManager == nil {
		t.Fatal("NewMacroManager() returned nil")
	}

	if macroManager.baseSend == nil {
		t.Error("macroManager.baseSend was not set correctly")
	}
}

// TestSendMacroStart tests the SendMacroStart function
func TestSendMacroStart(t *testing.T) {
	mockSend := NewMockSend()
	macroManager := NewMacroManager(mockSend)

	// Test sending a macro start request
	err := macroManager.SendMacroStart()
	if err != nil {
		t.Fatalf("SendMacroStart() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "macro_start" {
		t.Errorf("Expected packet ID 'macro_start', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendMacroStop tests the SendMacroStop function
func TestSendMacroStop(t *testing.T) {
	mockSend := NewMockSend()
	macroManager := NewMacroManager(mockSend)

	// Test sending a macro stop request
	err := macroManager.SendMacroStop()
	if err != nil {
		t.Fatalf("SendMacroStop() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "macro_stop" {
		t.Errorf("Expected packet ID 'macro_stop', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendMacroDetectorDownload tests the SendMacroDetectorDownload function
func TestSendMacroDetectorDownload(t *testing.T) {
	mockSend := NewMockSend()
	macroManager := NewMacroManager(mockSend)

	// Test sending a macro detector download request
	err := macroManager.SendMacroDetectorDownload()
	if err != nil {
		t.Fatalf("SendMacroDetectorDownload() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "macro_detector_download" {
		t.Errorf("Expected packet ID 'macro_detector_download', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendMacroDetectorAnswer tests the SendMacroDetectorAnswer function
func TestSendMacroDetectorAnswer(t *testing.T) {
	mockSend := NewMockSend()
	macroManager := NewMacroManager(mockSend)

	// Test sending a macro detector answer
	answer := "test answer"
	err := macroManager.SendMacroDetectorAnswer(answer)
	if err != nil {
		t.Fatalf("SendMacroDetectorAnswer() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "macro_detector_answer" {
		t.Errorf("Expected packet ID 'macro_detector_answer', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if answerVal, ok := mockSend.LastArgs()["answer"].([]byte); !ok {
		t.Errorf("Expected answer to be []byte, got %T", mockSend.LastArgs()["answer"])
	} else {
		answerStr := string(answerVal)
		if answerStr != answer {
			t.Errorf("Expected answer=%s, got %s", answer, answerStr)
		}
	}
}

// TestSendReqCashTabCode tests the SendReqCashTabCode function
func TestSendReqCashTabCode(t *testing.T) {
	mockSend := NewMockSend()
	macroManager := NewMacroManager(mockSend)

	// Test sending a request for cash tab code
	tabID := uint16(123)
	err := macroManager.SendReqCashTabCode(tabID)
	if err != nil {
		t.Fatalf("SendReqCashTabCode() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "req_cash_tabcode" {
		t.Errorf("Expected packet ID 'req_cash_tabcode', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint16); !ok || idVal != tabID {
		t.Errorf("Expected ID=%d, got %v", tabID, mockSend.LastArgs()["ID"])
	}
}
