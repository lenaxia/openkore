package chat

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForChannelChat is a mock implementation of the Send interface for testing channel chat functionality
type MockSendForChannelChat struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForChannelChat() *MockSendForChannelChat {
	return &MockSendForChannelChat{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForChannelChat) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForChannelChat) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForChannelChat) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForChannelChat) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForChannelChat) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForChannelChat) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForChannelChat) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForChannelChat) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForChannelChat) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForChannelChat) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForChannelChat) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForChannelChat) SetConnection(conn interface{}) {
}

func (ms *MockSendForChannelChat) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForChannelChat) GetTime() uint32 {
	return 12345
}

// TestNewChannelChatManager tests the NewChannelChatManager function
func TestNewChannelChatManager(t *testing.T) {
	mockSend := NewMockSendForChannelChat()
	ccm := NewChannelChatManager(mockSend)

	if ccm == nil {
		t.Fatal("NewChannelChatManager() returned nil")
	}

	if ccm.baseSend == nil {
		t.Error("ccm.baseSend was not set correctly")
	}
}

// TestJoinChannel tests the JoinChannel method
func TestJoinChannel(t *testing.T) {
	mockSend := NewMockSendForChannelChat()
	mockSend.packetLUT["join_channel"] = "00B2"

	ccm := NewChannelChatManager(mockSend)

	// Test joining a channel
	channelName := "Global"
	password := "secret123"
	err := ccm.JoinChannel(channelName, password)
	if err != nil {
		t.Fatalf("JoinChannel() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00B2"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["channel_name"] != channelName {
		t.Errorf("args[\"channel_name\"] = %v, want %v", args["channel_name"], channelName)
	}

	if args["password"] != password {
		t.Errorf("args[\"password\"] = %v, want %v", args["password"], password)
	}
}

// TestLeaveChannel tests the LeaveChannel method
func TestLeaveChannel(t *testing.T) {
	mockSend := NewMockSendForChannelChat()
	mockSend.packetLUT["leave_channel"] = "00B3"

	ccm := NewChannelChatManager(mockSend)

	// Test leaving a channel
	channelID := uint32(1)
	err := ccm.LeaveChannel(channelID)
	if err != nil {
		t.Fatalf("LeaveChannel() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00B3"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["channel_id"] != channelID {
		t.Errorf("args[\"channel_id\"] = %v, want %v", args["channel_id"], channelID)
	}
}

// TestSendChannelMessage tests the SendChannelMessage method
func TestSendChannelMessage(t *testing.T) {
	mockSend := NewMockSendForChannelChat()
	mockSend.packetLUT["channel_message"] = "00B5"

	ccm := NewChannelChatManager(mockSend)

	// Test sending a channel message
	channelID := uint32(1)
	message := "Hello, channel!"
	err := ccm.SendChannelMessage(channelID, message)
	if err != nil {
		t.Fatalf("SendChannelMessage() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00B5"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["channel_id"] != channelID {
		t.Errorf("args[\"channel_id\"] = %v, want %v", args["channel_id"], channelID)
	}

	if args["message"] != message {
		t.Errorf("args[\"message\"] = %v, want %v", args["message"], message)
	}
}

// TestListChannels tests the ListChannels method
func TestListChannels(t *testing.T) {
	mockSend := NewMockSendForChannelChat()
	mockSend.packetLUT["list_channels"] = "00B4"

	ccm := NewChannelChatManager(mockSend)

	// Test listing channels
	err := ccm.ListChannels()
	if err != nil {
		t.Fatalf("ListChannels() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	_, exists := mockSend.reconstructArgs["00B4"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}
