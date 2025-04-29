package chat

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForPrivateChat is a mock implementation of the Send interface for testing private chat functionality
type MockSendForPrivateChat struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForPrivateChat() *MockSendForPrivateChat {
	return &MockSendForPrivateChat{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForPrivateChat) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForPrivateChat) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForPrivateChat) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForPrivateChat) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForPrivateChat) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForPrivateChat) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForPrivateChat) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForPrivateChat) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForPrivateChat) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForPrivateChat) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForPrivateChat) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForPrivateChat) SetConnection(conn interface{}) {
}

func (ms *MockSendForPrivateChat) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForPrivateChat) GetTime() uint32 {
	return 12345
}

// TestNewPrivateChatManager tests the NewPrivateChatManager function
func TestNewPrivateChatManager(t *testing.T) {
	mockSend := NewMockSendForPrivateChat()
	pcm := NewPrivateChatManager(mockSend)

	if pcm == nil {
		t.Fatal("NewPrivateChatManager() returned nil")
	}

	if pcm.baseSend == nil {
		t.Error("pcm.baseSend was not set correctly")
	}
}

// TestSendPrivateMessage tests the SendPrivateMessage method
func TestSendPrivateMessage(t *testing.T) {
	mockSend := NewMockSendForPrivateChat()
	mockSend.packetLUT["private_message"] = "0096"

	pcm := NewPrivateChatManager(mockSend)

	// Test sending private message
	target := "PlayerName"
	message := "Hello, how are you?"
	err := pcm.SendPrivateMessage(target, message)
	if err != nil {
		t.Fatalf("SendPrivateMessage() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0096"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["target"] != target {
		t.Errorf("args[\"target\"] = %v, want %v", args["target"], target)
	}

	if args["message"] != message {
		t.Errorf("args[\"message\"] = %v, want %v", args["message"], message)
	}
}

// TestSendWhisperResponse tests the SendWhisperResponse method
func TestSendWhisperResponse(t *testing.T) {
	mockSend := NewMockSendForPrivateChat()
	mockSend.packetLUT["whisper_response"] = "0097"

	pcm := NewPrivateChatManager(mockSend)

	// Test sending whisper response
	target := "PlayerName"
	response := 1 // 1 = success, 2 = target offline, 3 = target is ignoring you
	err := pcm.SendWhisperResponse(target, response)
	if err != nil {
		t.Fatalf("SendWhisperResponse() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0097"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["target"] != target {
		t.Errorf("args[\"target\"] = %v, want %v", args["target"], target)
	}

	if args["response"] != response {
		t.Errorf("args[\"response\"] = %v, want %v", args["response"], response)
	}
}

// TestSendIgnorePlayer tests the SendIgnorePlayer method
func TestSendIgnorePlayer(t *testing.T) {
	mockSend := NewMockSendForPrivateChat()
	mockSend.packetLUT["ignore_player"] = "00CF"

	pcm := NewPrivateChatManager(mockSend)

	// Test sending ignore player
	target := "PlayerName"
	flag := 1 // 0 = unignore, 1 = ignore
	err := pcm.SendIgnorePlayer(target, flag)
	if err != nil {
		t.Fatalf("SendIgnorePlayer() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00CF"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["target"] != target {
		t.Errorf("args[\"target\"] = %v, want %v", args["target"], target)
	}

	if args["flag"] != flag {
		t.Errorf("args[\"flag\"] = %v, want %v", args["flag"], flag)
	}
}

// TestSendTalk tests the SendTalk method
func TestSendTalk(t *testing.T) {
	mockSend := NewMockSendForPrivateChat()
	mockSend.packetLUT["talk"] = "0090"

	pcm := NewPrivateChatManager(mockSend)

	// Test sending talk
	npcID := uint32(12345)
	err := pcm.SendTalk(npcID)
	if err != nil {
		t.Fatalf("SendTalk() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0090"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["npc_id"] != npcID {
		t.Errorf("args[\"npc_id\"] = %v, want %v", args["npc_id"], npcID)
	}
}

// TestSendTalkResponse tests the SendTalkResponse method
func TestSendTalkResponse(t *testing.T) {
	mockSend := NewMockSendForPrivateChat()
	mockSend.packetLUT["talk_response"] = "00B8"

	pcm := NewPrivateChatManager(mockSend)

	// Test sending talk response
	npcID := uint32(12345)
	response := 1
	err := pcm.SendTalkResponse(npcID, response)
	if err != nil {
		t.Fatalf("SendTalkResponse() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00B8"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["npc_id"] != npcID {
		t.Errorf("args[\"npc_id\"] = %v, want %v", args["npc_id"], npcID)
	}

	if args["response"] != uint8(response) {
		t.Errorf("args[\"response\"] = %v, want %v", args["response"], uint8(response))
	}
}
