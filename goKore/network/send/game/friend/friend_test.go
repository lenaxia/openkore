package friend

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForFriend is a mock implementation of the Send interface for testing friend functionality
type MockSendForFriend struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForFriend() *MockSendForFriend {
	return &MockSendForFriend{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForFriend) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForFriend) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForFriend) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForFriend) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForFriend) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForFriend) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForFriend) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForFriend) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForFriend) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForFriend) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForFriend) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForFriend) SetConnection(conn interface{}) {
}

func (ms *MockSendForFriend) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForFriend) GetTime() uint32 {
	return 12345
}

// TestNewFriendManager tests the NewFriendManager function
func TestNewFriendManager(t *testing.T) {
	mockSend := NewMockSendForFriend()
	fm := NewFriendManager(mockSend)

	if fm == nil {
		t.Fatal("NewFriendManager() returned nil")
	}

	if fm.baseSend == nil {
		t.Error("fm.baseSend was not set correctly")
	}
}

// TestSendFriendListReply tests the SendFriendListReply method
func TestSendFriendListReply(t *testing.T) {
	mockSend := NewMockSendForFriend()
	mockSend.packetLUT["friend_response"] = "0208"

	fm := NewFriendManager(mockSend)

	// Test replying to a friend request
	accountID := uint32(12345)
	charID := uint32(54321)
	flag := 1
	err := fm.SendFriendListReply(accountID, charID, flag)
	if err != nil {
		t.Fatalf("SendFriendListReply() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0208"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["friendAccountID"] != accountID {
		t.Errorf("args[\"friendAccountID\"] = %v, want %v", args["friendAccountID"], accountID)
	}

	if args["friendCharID"] != charID {
		t.Errorf("args[\"friendCharID\"] = %v, want %v", args["friendCharID"], charID)
	}

	if args["type"] != flag {
		t.Errorf("args[\"type\"] = %v, want %v", args["type"], flag)
	}
}

// TestSendFriendRequest tests the SendFriendRequest method
func TestSendFriendRequest(t *testing.T) {
	mockSend := NewMockSendForFriend()
	mockSend.packetLUT["friend_request"] = "0202"

	fm := NewFriendManager(mockSend)

	// Test sending a friend request
	name := "TestPlayer"
	err := fm.SendFriendRequest(name)
	if err != nil {
		t.Fatalf("SendFriendRequest() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0202"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	// Check that name was converted to bytes and padded to 24 bytes
	nameBytes, ok := args["username"].([]byte)
	if !ok {
		t.Errorf("args[\"username\"] is not a byte slice")
	} else {
		// Check that the name is correct (up to the length of the original name)
		if string(nameBytes[:len(name)]) != name {
			t.Errorf("args[\"username\"] = %v, want %v", string(nameBytes[:len(name)]), name)
		}
		// Check that the length is 24 bytes
		if len(nameBytes) != 24 {
			t.Errorf("len(args[\"username\"]) = %v, want 24", len(nameBytes))
		}
	}
}

// TestSendFriendRemove tests the SendFriendRemove method
func TestSendFriendRemove(t *testing.T) {
	mockSend := NewMockSendForFriend()
	mockSend.packetLUT["friend_remove"] = "0203"

	fm := NewFriendManager(mockSend)

	// Test removing a friend
	accountID := uint32(12345)
	charID := uint32(54321)
	err := fm.SendFriendRemove(accountID, charID)
	if err != nil {
		t.Fatalf("SendFriendRemove() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0203"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["accountID"] != accountID {
		t.Errorf("args[\"accountID\"] = %v, want %v", args["accountID"], accountID)
	}

	if args["charID"] != charID {
		t.Errorf("args[\"charID\"] = %v, want %v", args["charID"], charID)
	}
}

// TestSendIgnore tests the SendIgnore method
func TestSendIgnore(t *testing.T) {
	mockSend := NewMockSendForFriend()
	mockSend.packetLUT["ignore_player"] = "00CF"

	fm := NewFriendManager(mockSend)

	// Test ignoring a player
	name := "TestPlayer"
	flag := 1
	err := fm.SendIgnore(name, flag)
	if err != nil {
		t.Fatalf("SendIgnore() returned error: %v", err)
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

	// Check that name was converted to bytes
	nameBytes, ok := args["name"].([]byte)
	if !ok {
		t.Errorf("args[\"name\"] is not a byte slice")
	} else if string(nameBytes) != name {
		t.Errorf("args[\"name\"] = %v, want %v", string(nameBytes), name)
	}

	if args["flag"] != flag {
		t.Errorf("args[\"flag\"] = %v, want %v", args["flag"], flag)
	}
}

// TestSendIgnoreAll tests the SendIgnoreAll method
func TestSendIgnoreAll(t *testing.T) {
	mockSend := NewMockSendForFriend()
	mockSend.packetLUT["ignore_all"] = "00D0"

	fm := NewFriendManager(mockSend)

	// Test ignoring all players
	flag := 1
	err := fm.SendIgnoreAll(flag)
	if err != nil {
		t.Fatalf("SendIgnoreAll() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00D0"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["flag"] != flag {
		t.Errorf("args[\"flag\"] = %v, want %v", args["flag"], flag)
	}
}

// TestSendGetIgnoreList tests the SendGetIgnoreList method
func TestSendGetIgnoreList(t *testing.T) {
	mockSend := NewMockSendForFriend()
	mockSend.packetLUT["get_ignore_list"] = "00D3"

	fm := NewFriendManager(mockSend)

	// Test getting the ignore list
	err := fm.SendGetIgnoreList()
	if err != nil {
		t.Fatalf("SendGetIgnoreList() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["00D3"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}
