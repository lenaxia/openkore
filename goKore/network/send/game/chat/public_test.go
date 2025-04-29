package chat

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForPublicChat is a mock implementation of the Send interface for testing public chat functionality
type MockSendForPublicChat struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForPublicChat() *MockSendForPublicChat {
	return &MockSendForPublicChat{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForPublicChat) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForPublicChat) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForPublicChat) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForPublicChat) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForPublicChat) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForPublicChat) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForPublicChat) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForPublicChat) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForPublicChat) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForPublicChat) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForPublicChat) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForPublicChat) SetConnection(conn interface{}) {
}

func (ms *MockSendForPublicChat) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForPublicChat) GetTime() uint32 {
	return 12345
}

// TestNewPublicChatManager tests the NewPublicChatManager function
func TestNewPublicChatManager(t *testing.T) {
	mockSend := NewMockSendForPublicChat()
	pcm := NewPublicChatManager(mockSend)

	if pcm == nil {
		t.Fatal("NewPublicChatManager() returned nil")
	}

	if pcm.baseSend == nil {
		t.Error("pcm.baseSend was not set correctly")
	}
}

// TestSendChat tests the SendChat method
func TestSendChat(t *testing.T) {
	mockSend := NewMockSendForPublicChat()
	mockSend.packetLUT["public_chat"] = "008C"

	pcm := NewPublicChatManager(mockSend)

	// Test sending public chat message
	message := "Hello, world!"
	err := pcm.SendChat(message)
	if err != nil {
		t.Fatalf("SendChat() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["008C"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["message"] != message {
		t.Errorf("args[\"message\"] = %v, want %v", args["message"], message)
	}
}

// TestSendGMBroadcast tests the SendGMBroadcast method
func TestSendGMBroadcast(t *testing.T) {
	mockSend := NewMockSendForPublicChat()
	mockSend.packetLUT["gm_broadcast"] = "00F3"

	pcm := NewPublicChatManager(mockSend)

	// Test sending GM broadcast message
	message := "Server maintenance in 10 minutes!"
	err := pcm.SendGMBroadcast(message)
	if err != nil {
		t.Fatalf("SendGMBroadcast() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00F3"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["message"] != message {
		t.Errorf("args[\"message\"] = %v, want %v", args["message"], message)
	}
}

// TestSendLocalBroadcast tests the SendLocalBroadcast method
func TestSendLocalBroadcast(t *testing.T) {
	mockSend := NewMockSendForPublicChat()
	mockSend.packetLUT["local_broadcast"] = "009A"

	pcm := NewPublicChatManager(mockSend)

	// Test sending local broadcast message
	message := "Welcome to Prontera!"
	err := pcm.SendLocalBroadcast(message)
	if err != nil {
		t.Fatalf("SendLocalBroadcast() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["009A"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["message"] != message {
		t.Errorf("args[\"message\"] = %v, want %v", args["message"], message)
	}
}
