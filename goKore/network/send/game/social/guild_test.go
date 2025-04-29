package social

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForGuild is a mock implementation of the Send interface for testing guild functionality
type MockSendForGuild struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForGuild() *MockSendForGuild {
	return &MockSendForGuild{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForGuild) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForGuild) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForGuild) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForGuild) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForGuild) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForGuild) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForGuild) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForGuild) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForGuild) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForGuild) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForGuild) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForGuild) SetConnection(conn interface{}) {
}

func (ms *MockSendForGuild) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForGuild) GetTime() uint32 {
	return 12345
}

// TestNewGuildManager tests the NewGuildManager function
func TestNewGuildManager(t *testing.T) {
	mockSend := NewMockSendForGuild()
	gm := NewGuildManager(mockSend)

	if gm == nil {
		t.Fatal("NewGuildManager() returned nil")
	}

	if gm.baseSend == nil {
		t.Error("gm.baseSend was not set correctly")
	}
}

// TestCreateGuild tests the CreateGuild method
func TestCreateGuild(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["create_guild"] = "0165"

	gm := NewGuildManager(mockSend)

	// Test creating a guild
	name := "My Guild"
	err := gm.CreateGuild(name)
	if err != nil {
		t.Fatalf("CreateGuild() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0165"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["name"] != name {
		t.Errorf("args[\"name\"] = %v, want %v", args["name"], name)
	}
}

// TestJoinGuild tests the JoinGuild method
func TestJoinGuild(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["join_guild"] = "0168"

	gm := NewGuildManager(mockSend)

	// Test joining a guild
	guildID := uint32(12345)
	err := gm.JoinGuild(guildID)
	if err != nil {
		t.Fatalf("JoinGuild() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0168"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["guild_id"] != guildID {
		t.Errorf("args[\"guild_id\"] = %v, want %v", args["guild_id"], guildID)
	}
}

// TestLeaveGuild tests the LeaveGuild method
func TestLeaveGuild(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["leave_guild"] = "0159"

	gm := NewGuildManager(mockSend)

	// Test leaving a guild
	guildID := uint32(12345)
	accountID := uint32(67890)
	charID := uint32(54321)
	reason := "I'm leaving"
	err := gm.LeaveGuild(guildID, accountID, charID, reason)
	if err != nil {
		t.Fatalf("LeaveGuild() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0159"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["guild_id"] != guildID {
		t.Errorf("args[\"guild_id\"] = %v, want %v", args["guild_id"], guildID)
	}

	if args["account_id"] != accountID {
		t.Errorf("args[\"account_id\"] = %v, want %v", args["account_id"], accountID)
	}

	if args["char_id"] != charID {
		t.Errorf("args[\"char_id\"] = %v, want %v", args["char_id"], charID)
	}

	if args["reason"] != reason {
		t.Errorf("args[\"reason\"] = %v, want %v", args["reason"], reason)
	}
}

// TestInviteToGuild tests the InviteToGuild method
func TestInviteToGuild(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["invite_to_guild"] = "0168"

	gm := NewGuildManager(mockSend)

	// Test inviting to a guild
	accountID := uint32(67890)
	err := gm.InviteToGuild(accountID)
	if err != nil {
		t.Fatalf("InviteToGuild() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0168"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["account_id"] != accountID {
		t.Errorf("args[\"account_id\"] = %v, want %v", args["account_id"], accountID)
	}
}

// TestGuildChat tests the GuildChat method
func TestGuildChat(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["guild_chat"] = "017E"

	gm := NewGuildManager(mockSend)

	// Test sending a guild chat message
	message := "Hello, guild!"
	err := gm.GuildChat(message)
	if err != nil {
		t.Fatalf("GuildChat() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["017E"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["message"] != message {
		t.Errorf("args[\"message\"] = %v, want %v", args["message"], message)
	}
}

// TestChangeGuildPositionInfo tests the ChangeGuildPositionInfo method
func TestChangeGuildPositionInfo(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["change_guild_position_info"] = "0161"

	gm := NewGuildManager(mockSend)

	// Test changing guild position info
	positionID := uint32(1)
	name := "Officer"
	mode := uint32(0x01)
	ranking := uint32(2)
	err := gm.ChangeGuildPositionInfo(positionID, name, mode, ranking)
	if err != nil {
		t.Fatalf("ChangeGuildPositionInfo() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0161"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["position_id"] != positionID {
		t.Errorf("args[\"position_id\"] = %v, want %v", args["position_id"], positionID)
	}

	if args["name"] != name {
		t.Errorf("args[\"name\"] = %v, want %v", args["name"], name)
	}

	if args["mode"] != mode {
		t.Errorf("args[\"mode\"] = %v, want %v", args["mode"], mode)
	}

	if args["ranking"] != ranking {
		t.Errorf("args[\"ranking\"] = %v, want %v", args["ranking"], ranking)
	}
}
