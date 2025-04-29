package social

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForParty is a mock implementation of the Send interface for testing party functionality
type MockSendForParty struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForParty() *MockSendForParty {
	return &MockSendForParty{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForParty) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForParty) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForParty) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForParty) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForParty) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForParty) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForParty) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForParty) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForParty) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForParty) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForParty) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForParty) SetConnection(conn interface{}) {
}

func (ms *MockSendForParty) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForParty) GetTime() uint32 {
	return 12345
}

// TestNewPartyManager tests the NewPartyManager function
func TestNewPartyManager(t *testing.T) {
	mockSend := NewMockSendForParty()
	pm := NewPartyManager(mockSend)

	if pm == nil {
		t.Fatal("NewPartyManager() returned nil")
	}

	if pm.baseSend == nil {
		t.Error("pm.baseSend was not set correctly")
	}
}

// TestCreateParty tests the CreateParty method
func TestCreateParty(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["create_party"] = "00F9"

	pm := NewPartyManager(mockSend)

	// Test creating a party
	name := "My Party"
	shareExp := true
	shareItems := false
	err := pm.CreateParty(name, shareExp, shareItems)
	if err != nil {
		t.Fatalf("CreateParty() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00F9"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["name"] != name {
		t.Errorf("args[\"name\"] = %v, want %v", args["name"], name)
	}

	if args["share_exp"] != shareExp {
		t.Errorf("args[\"share_exp\"] = %v, want %v", args["share_exp"], shareExp)
	}

	if args["share_items"] != shareItems {
		t.Errorf("args[\"share_items\"] = %v, want %v", args["share_items"], shareItems)
	}
}

// TestJoinParty tests the JoinParty method
func TestJoinParty(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["join_party"] = "00FC"

	pm := NewPartyManager(mockSend)

	// Test joining a party
	partyID := uint32(12345)
	err := pm.JoinParty(partyID)
	if err != nil {
		t.Fatalf("JoinParty() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00FC"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["party_id"] != partyID {
		t.Errorf("args[\"party_id\"] = %v, want %v", args["party_id"], partyID)
	}
}

// TestLeaveParty tests the LeaveParty method
func TestLeaveParty(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["leave_party"] = "0100"

	pm := NewPartyManager(mockSend)

	// Test leaving a party
	err := pm.LeaveParty()
	if err != nil {
		t.Fatalf("LeaveParty() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	_, exists := mockSend.reconstructArgs["0100"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestInviteToParty tests the InviteToParty method
func TestInviteToParty(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["invite_to_party"] = "00FC"

	pm := NewPartyManager(mockSend)

	// Test inviting to a party
	playerName := "Player1"
	err := pm.InviteToParty(playerName)
	if err != nil {
		t.Fatalf("InviteToParty() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00FC"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["player_name"] != playerName {
		t.Errorf("args[\"player_name\"] = %v, want %v", args["player_name"], playerName)
	}
}

// TestKickFromParty tests the KickFromParty method
func TestKickFromParty(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["kick_from_party"] = "0103"

	pm := NewPartyManager(mockSend)

	// Test kicking from a party
	playerName := "Player1"
	err := pm.KickFromParty(playerName)
	if err != nil {
		t.Fatalf("KickFromParty() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0103"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["player_name"] != playerName {
		t.Errorf("args[\"player_name\"] = %v, want %v", args["player_name"], playerName)
	}
}

// TestPartyChat tests the PartyChat method
func TestPartyChat(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["party_chat"] = "0108"

	pm := NewPartyManager(mockSend)

	// Test sending a party chat message
	message := "Hello, party!"
	err := pm.PartyChat(message)
	if err != nil {
		t.Fatalf("PartyChat() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0108"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["message"] != message {
		t.Errorf("args[\"message\"] = %v, want %v", args["message"], message)
	}
}
