package social

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

// TestRequestFriendList tests the RequestFriendList method
func TestRequestFriendList(t *testing.T) {
	mockSend := NewMockSendForFriend()
	mockSend.packetLUT["request_friend_list"] = "0202"

	fm := NewFriendManager(mockSend)

	// Test requesting friend list
	err := fm.RequestFriendList()
	if err != nil {
		t.Fatalf("RequestFriendList() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	_, exists := mockSend.reconstructArgs["0202"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestAddFriend tests the AddFriend method
func TestAddFriend(t *testing.T) {
	mockSend := NewMockSendForFriend()
	mockSend.packetLUT["add_friend"] = "0206"

	fm := NewFriendManager(mockSend)

	// Test adding a friend
	name := "Friend1"
	err := fm.AddFriend(name)
	if err != nil {
		t.Fatalf("AddFriend() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0206"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["name"] != name {
		t.Errorf("args[\"name\"] = %v, want %v", args["name"], name)
	}
}

// TestRemoveFriend tests the RemoveFriend method
func TestRemoveFriend(t *testing.T) {
	mockSend := NewMockSendForFriend()
	mockSend.packetLUT["remove_friend"] = "0203"

	fm := NewFriendManager(mockSend)

	// Test removing a friend
	accountID := uint32(12345)
	err := fm.RemoveFriend(accountID)
	if err != nil {
		t.Fatalf("RemoveFriend() returned error: %v", err)
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

	if args["account_id"] != accountID {
		t.Errorf("args[\"account_id\"] = %v, want %v", args["account_id"], accountID)
	}
}

// TestAcceptFriendRequest tests the AcceptFriendRequest method
func TestAcceptFriendRequest(t *testing.T) {
	mockSend := NewMockSendForFriend()
	mockSend.packetLUT["accept_friend_request"] = "0208"

	fm := NewFriendManager(mockSend)

	// Test accepting a friend request
	accountID := uint32(12345)
	charID := uint32(67890)
	err := fm.AcceptFriendRequest(accountID, charID)
	if err != nil {
		t.Fatalf("AcceptFriendRequest() returned error: %v", err)
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

	if args["account_id"] != accountID {
		t.Errorf("args[\"account_id\"] = %v, want %v", args["account_id"], accountID)
	}

	if args["char_id"] != charID {
		t.Errorf("args[\"char_id\"] = %v, want %v", args["char_id"], charID)
	}
}

// TestRejectFriendRequest tests the RejectFriendRequest method
func TestRejectFriendRequest(t *testing.T) {
	mockSend := NewMockSendForFriend()
	mockSend.packetLUT["reject_friend_request"] = "0207"

	fm := NewFriendManager(mockSend)

	// Test rejecting a friend request
	accountID := uint32(12345)
	charID := uint32(67890)
	err := fm.RejectFriendRequest(accountID, charID)
	if err != nil {
		t.Fatalf("RejectFriendRequest() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0207"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["account_id"] != accountID {
		t.Errorf("args[\"account_id\"] = %v, want %v", args["account_id"], accountID)
	}

	if args["char_id"] != charID {
		t.Errorf("args[\"char_id\"] = %v, want %v", args["char_id"], charID)
	}
}
