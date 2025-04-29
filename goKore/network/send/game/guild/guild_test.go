package guild

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

// TestSendGuildMasterMemberCheck tests the SendGuildMasterMemberCheck method
func TestSendGuildMasterMemberCheck(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["guild_check"] = "014D"

	gm := NewGuildManager(mockSend)

	// Test checking guild master/member status
	err := gm.SendGuildMasterMemberCheck()
	if err != nil {
		t.Fatalf("SendGuildMasterMemberCheck() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["014D"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestSendGuildRequestInfo tests the SendGuildRequestInfo method
func TestSendGuildRequestInfo(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["guild_info_request"] = "014E"

	gm := NewGuildManager(mockSend)

	// Test requesting guild info
	page := 1 // 0-4
	err := gm.SendGuildRequestInfo(page)
	if err != nil {
		t.Fatalf("SendGuildRequestInfo() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["014E"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["type"] != page {
		t.Errorf("args[\"type\"] = %v, want %v", args["type"], page)
	}
}

// TestSendGuildAlly tests the SendGuildAlly method
func TestSendGuildAlly(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["guild_alliance_reply"] = "0173"

	gm := NewGuildManager(mockSend)

	// Test responding to guild alliance request
	ID := uint32(12345)
	flag := 1
	err := gm.SendGuildAlly(ID, flag)
	if err != nil {
		t.Fatalf("SendGuildAlly() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0173"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != ID {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], ID)
	}

	if args["flag"] != flag {
		t.Errorf("args[\"flag\"] = %v, want %v", args["flag"], flag)
	}
}

// TestSendGuildRequestEmblem tests the SendGuildRequestEmblem method
func TestSendGuildRequestEmblem(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["guild_emblem_request"] = "0151"

	gm := NewGuildManager(mockSend)

	// Test requesting guild emblem
	guildID := uint32(12345)
	err := gm.SendGuildRequestEmblem(guildID)
	if err != nil {
		t.Fatalf("SendGuildRequestEmblem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0151"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["guildID"] != guildID {
		t.Errorf("args[\"guildID\"] = %v, want %v", args["guildID"], guildID)
	}
}

// TestSendGuildBreak tests the SendGuildBreak method
func TestSendGuildBreak(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["guild_break"] = "015D"

	gm := NewGuildManager(mockSend)

	// Test breaking a guild
	guildName := "TestGuild"
	err := gm.SendGuildBreak(guildName)
	if err != nil {
		t.Fatalf("SendGuildBreak() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["015D"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	// Check that name was converted to bytes
	nameBytes, ok := args["guildName"].([]byte)
	if !ok {
		t.Errorf("args[\"guildName\"] is not a byte slice")
	} else if string(nameBytes) != guildName {
		t.Errorf("args[\"guildName\"] = %v, want %v", string(nameBytes), guildName)
	}
}

// TestSendGuildLeave tests the SendGuildLeave method
func TestSendGuildLeave(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["guild_leave"] = "0159"

	gm := NewGuildManager(mockSend)

	// Test leaving a guild
	reason := "Test reason"
	guildID := uint32(12345)
	accountID := uint32(67890)
	charID := uint32(54321)
	err := gm.SendGuildLeave(reason, guildID, accountID, charID)
	if err != nil {
		t.Fatalf("SendGuildLeave() returned error: %v", err)
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

	if args["guildID"] != guildID {
		t.Errorf("args[\"guildID\"] = %v, want %v", args["guildID"], guildID)
	}

	if args["accountID"] != accountID {
		t.Errorf("args[\"accountID\"] = %v, want %v", args["accountID"], accountID)
	}

	if args["charID"] != charID {
		t.Errorf("args[\"charID\"] = %v, want %v", args["charID"], charID)
	}

	// Check that reason was converted to bytes
	reasonBytes, ok := args["reason"].([]byte)
	if !ok {
		t.Errorf("args[\"reason\"] is not a byte slice")
	} else if string(reasonBytes) != reason {
		t.Errorf("args[\"reason\"] = %v, want %v", string(reasonBytes), reason)
	}
}

// TestSendGuildMemberKick tests the SendGuildMemberKick method
func TestSendGuildMemberKick(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["guild_kick"] = "015B"

	gm := NewGuildManager(mockSend)

	// Test kicking a guild member
	guildID := uint32(12345)
	accountID := uint32(67890)
	charID := uint32(54321)
	reason := "Test reason"
	err := gm.SendGuildMemberKick(guildID, accountID, charID, reason)
	if err != nil {
		t.Fatalf("SendGuildMemberKick() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["015B"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["guildID"] != guildID {
		t.Errorf("args[\"guildID\"] = %v, want %v", args["guildID"], guildID)
	}

	if args["accountID"] != accountID {
		t.Errorf("args[\"accountID\"] = %v, want %v", args["accountID"], accountID)
	}

	if args["charID"] != charID {
		t.Errorf("args[\"charID\"] = %v, want %v", args["charID"], charID)
	}

	// Check that reason was converted to bytes
	reasonBytes, ok := args["reason"].([]byte)
	if !ok {
		t.Errorf("args[\"reason\"] is not a byte slice")
	} else if string(reasonBytes) != reason {
		t.Errorf("args[\"reason\"] = %v, want %v", string(reasonBytes), reason)
	}
}

// TestSendGuildCreate tests the SendGuildCreate method
func TestSendGuildCreate(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["guild_create"] = "0165"

	gm := NewGuildManager(mockSend)

	// Test creating a guild
	name := "TestGuild"
	charID := uint32(12345)
	err := gm.SendGuildCreate(name, charID)
	if err != nil {
		t.Fatalf("SendGuildCreate() returned error: %v", err)
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

	if args["charID"] != charID {
		t.Errorf("args[\"charID\"] = %v, want %v", args["charID"], charID)
	}

	// Check that name was converted to bytes
	nameBytes, ok := args["guildName"].([]byte)
	if !ok {
		t.Errorf("args[\"guildName\"] is not a byte slice")
	} else if string(nameBytes) != name {
		t.Errorf("args[\"guildName\"] = %v, want %v", string(nameBytes), name)
	}
}

// TestSendGuildJoin tests the SendGuildJoin method
func TestSendGuildJoin(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["guild_join"] = "0169"

	gm := NewGuildManager(mockSend)

	// Test joining a guild
	ID := uint32(12345)
	flag := 1
	err := gm.SendGuildJoin(ID, flag)
	if err != nil {
		t.Fatalf("SendGuildJoin() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0169"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != ID {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], ID)
	}

	if args["flag"] != flag {
		t.Errorf("args[\"flag\"] = %v, want %v", args["flag"], flag)
	}
}

// TestSendGuildJoinRequest tests the SendGuildJoinRequest method
func TestSendGuildJoinRequest(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["guild_join_request"] = "016E"

	gm := NewGuildManager(mockSend)

	// Test requesting to join a guild
	ID := uint32(12345)
	accountID := uint32(67890)
	charID := uint32(54321)
	err := gm.SendGuildJoinRequest(ID, accountID, charID)
	if err != nil {
		t.Fatalf("SendGuildJoinRequest() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["016E"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != ID {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], ID)
	}

	if args["accountID"] != accountID {
		t.Errorf("args[\"accountID\"] = %v, want %v", args["accountID"], accountID)
	}

	if args["charID"] != charID {
		t.Errorf("args[\"charID\"] = %v, want %v", args["charID"], charID)
	}
}

// TestSendGuildSetAlly tests the SendGuildSetAlly method
func TestSendGuildSetAlly(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["guild_alliance_request"] = "0170"

	gm := NewGuildManager(mockSend)

	// Test requesting a guild alliance
	targetAID := uint32(12345)
	myAID := uint32(67890)
	charID := uint32(54321)
	err := gm.SendGuildSetAlly(targetAID, myAID, charID)
	if err != nil {
		t.Fatalf("SendGuildSetAlly() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0170"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["targetAccountID"] != targetAID {
		t.Errorf("args[\"targetAccountID\"] = %v, want %v", args["targetAccountID"], targetAID)
	}

	if args["accountID"] != myAID {
		t.Errorf("args[\"accountID\"] = %v, want %v", args["accountID"], myAID)
	}

	if args["charID"] != charID {
		t.Errorf("args[\"charID\"] = %v, want %v", args["charID"], charID)
	}
}

// TestSendGuildNotice tests the SendGuildNotice method
func TestSendGuildNotice(t *testing.T) {
	mockSend := NewMockSendForGuild()
	mockSend.packetLUT["guild_notice"] = "016E"

	gm := NewGuildManager(mockSend)

	// Test setting guild notice
	guildID := uint32(12345)
	name := "TestGuild"
	notice := "Test notice"
	err := gm.SendGuildNotice(guildID, name, notice)
	if err != nil {
		t.Fatalf("SendGuildNotice() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["016E"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["guildID"] != guildID {
		t.Errorf("args[\"guildID\"] = %v, want %v", args["guildID"], guildID)
	}

	// Check that name was converted to bytes
	nameBytes, ok := args["name"].([]byte)
	if !ok {
		t.Errorf("args[\"name\"] is not a byte slice")
	} else if string(nameBytes) != name {
		t.Errorf("args[\"name\"] = %v, want %v", string(nameBytes), name)
	}

	// Check that notice was converted to bytes
	noticeBytes, ok := args["notice"].([]byte)
	if !ok {
		t.Errorf("args[\"notice\"] is not a byte slice")
	} else if string(noticeBytes) != notice {
		t.Errorf("args[\"notice\"] = %v, want %v", string(noticeBytes), notice)
	}
}
