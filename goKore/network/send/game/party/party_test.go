package party

import (
	"fmt"
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

// TestSendPartyOption tests the SendPartyOption method
func TestSendPartyOption(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["party_setting"] = "07D7"

	pm := NewPartyManager(mockSend)

	// Test setting party options
	exp := 1
	itemPickup := 1
	itemDivision := 1
	err := pm.SendPartyOption(exp, itemPickup, itemDivision)
	if err != nil {
		t.Fatalf("SendPartyOption() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["07D7"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["exp"] != exp {
		t.Errorf("args[\"exp\"] = %v, want %v", args["exp"], exp)
	}

	if args["itemPickup"] != itemPickup {
		t.Errorf("args[\"itemPickup\"] = %v, want %v", args["itemPickup"], itemPickup)
	}

	if args["itemDivision"] != itemDivision {
		t.Errorf("args[\"itemDivision\"] = %v, want %v", args["itemDivision"], itemDivision)
	}
}

// TestSendPartyLeader tests the SendPartyLeader method
func TestSendPartyLeader(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["party_leader"] = "07DA"

	pm := NewPartyManager(mockSend)

	// Test changing party leader
	accountID := uint32(12345)
	err := pm.SendPartyLeader(accountID)
	if err != nil {
		t.Fatalf("SendPartyLeader() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["07DA"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["accountID"] != accountID {
		t.Errorf("args[\"accountID\"] = %v, want %v", args["accountID"], accountID)
	}
}

// TestSendPartyJoinRequest tests the SendPartyJoinRequest method
func TestSendPartyJoinRequest(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["party_join_request"] = "00FC"

	pm := NewPartyManager(mockSend)

	// Test requesting to join a party
	ID := uint32(12345)
	err := pm.SendPartyJoinRequest(ID)
	if err != nil {
		t.Fatalf("SendPartyJoinRequest() returned error: %v", err)
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

	if args["ID"] != ID {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], ID)
	}
}

// TestSendPartyJoin tests the SendPartyJoin method
func TestSendPartyJoin(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["party_join"] = "00FF"

	pm := NewPartyManager(mockSend)

	// Test joining a party
	ID := uint32(12345)
	flag := 1
	err := pm.SendPartyJoin(ID, flag)
	if err != nil {
		t.Fatalf("SendPartyJoin() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00FF"]
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

// TestSendPartyKick tests the SendPartyKick method
func TestSendPartyKick(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["party_kick"] = "0103"

	pm := NewPartyManager(mockSend)

	// Test kicking a party member
	ID := uint32(12345)
	name := "TestPlayer"
	err := pm.SendPartyKick(ID, name)
	if err != nil {
		t.Fatalf("SendPartyKick() returned error: %v", err)
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

	if args["ID"] != ID {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], ID)
	}

	// Check that name was converted to bytes
	nameBytes, ok := args["name"].([]byte)
	if !ok {
		t.Errorf("args[\"name\"] is not a byte slice")
	} else if string(nameBytes) != name {
		t.Errorf("args[\"name\"] = %v, want %v", string(nameBytes), name)
	}
}

// TestSendPartyJoinRequestByName tests the SendPartyJoinRequestByName method
func TestSendPartyJoinRequestByName(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["party_join_request_by_name"] = "02C4"

	pm := NewPartyManager(mockSend)

	// Test requesting to join a party by name
	name := "TestParty"
	err := pm.SendPartyJoinRequestByName(name)
	if err != nil {
		t.Fatalf("SendPartyJoinRequestByName() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["02C4"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	// Check that name was converted to bytes
	nameBytes, ok := args["partyName"].([]byte)
	if !ok {
		t.Errorf("args[\"partyName\"] is not a byte slice")
	} else if string(nameBytes) != name {
		t.Errorf("args[\"partyName\"] = %v, want %v", string(nameBytes), name)
	}
}

// TestSendPartyJoinRequestByNameReply tests the SendPartyJoinRequestByNameReply method
func TestSendPartyJoinRequestByNameReply(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["party_join_request_by_name_reply"] = "02C7"

	pm := NewPartyManager(mockSend)

	// Test replying to a party join request by name
	accountID := uint32(12345)
	flag := 1
	err := pm.SendPartyJoinRequestByNameReply(accountID, flag)
	if err != nil {
		t.Fatalf("SendPartyJoinRequestByNameReply() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["02C7"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["accountID"] != accountID {
		t.Errorf("args[\"accountID\"] = %v, want %v", args["accountID"], accountID)
	}

	if args["flag"] != flag {
		t.Errorf("args[\"flag\"] = %v, want %v", args["flag"], flag)
	}
}

// TestSendPartyBookingRegister tests the SendPartyBookingRegister method
func TestSendPartyBookingRegister(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["booking_register"] = "0802"

	pm := NewPartyManager(mockSend)

	// Test registering a party booking
	level := 100
	mapID := 1
	jobList := []int{1, 2, 3, 4, 5, 6}
	err := pm.SendPartyBookingRegister(level, mapID, jobList)
	if err != nil {
		t.Fatalf("SendPartyBookingRegister() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0802"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["level"] != level {
		t.Errorf("args[\"level\"] = %v, want %v", args["level"], level)
	}

	if args["MapID"] != mapID {
		t.Errorf("args[\"MapID\"] = %v, want %v", args["MapID"], mapID)
	}

	for i := 0; i < 6; i++ {
		jobKey := fmt.Sprintf("job%d", i)
		if args[jobKey] != jobList[i] {
			t.Errorf("args[\"%s\"] = %v, want %v", jobKey, args[jobKey], jobList[i])
		}
	}
}

// TestSendPartyBookingReqSearch tests the SendPartyBookingReqSearch method
func TestSendPartyBookingReqSearch(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["booking_search"] = "0804"

	pm := NewPartyManager(mockSend)

	// Test searching for party bookings
	level := 100
	mapID := 1
	job := 0 // Should be converted to 65535
	lastIndex := 0
	resultCount := 0 // Should be converted to 10
	err := pm.SendPartyBookingReqSearch(level, mapID, job, lastIndex, resultCount)
	if err != nil {
		t.Fatalf("SendPartyBookingReqSearch() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0804"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["level"] != level {
		t.Errorf("args[\"level\"] = %v, want %v", args["level"], level)
	}

	if args["MapID"] != mapID {
		t.Errorf("args[\"MapID\"] = %v, want %v", args["MapID"], mapID)
	}

	if args["job"] != 65535 {
		t.Errorf("args[\"job\"] = %v, want 65535", args["job"])
	}

	if args["LastIndex"] != lastIndex {
		t.Errorf("args[\"LastIndex\"] = %v, want %v", args["LastIndex"], lastIndex)
	}

	if args["ResultCount"] != 10 {
		t.Errorf("args[\"ResultCount\"] = %v, want 10", args["ResultCount"])
	}
}

// TestSendPartyBookingDelete tests the SendPartyBookingDelete method
func TestSendPartyBookingDelete(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["booking_delete"] = "0806"

	pm := NewPartyManager(mockSend)

	// Test deleting a party booking
	err := pm.SendPartyBookingDelete()
	if err != nil {
		t.Fatalf("SendPartyBookingDelete() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["0806"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestSendPartyBookingUpdate tests the SendPartyBookingUpdate method
func TestSendPartyBookingUpdate(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["booking_update"] = "0808"

	pm := NewPartyManager(mockSend)

	// Test updating a party booking
	jobList := []int{1, 2, 3, 4, 5, 6}
	err := pm.SendPartyBookingUpdate(jobList)
	if err != nil {
		t.Fatalf("SendPartyBookingUpdate() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0808"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	for i := 0; i < 6; i++ {
		jobKey := fmt.Sprintf("job%d", i)
		if args[jobKey] != jobList[i] {
			t.Errorf("args[\"%s\"] = %v, want %v", jobKey, args[jobKey], jobList[i])
		}
	}
}

// TestSendPartyLeave tests the SendPartyLeave method
func TestSendPartyLeave(t *testing.T) {
	mockSend := NewMockSendForParty()
	mockSend.packetLUT["party_leave"] = "0100"

	pm := NewPartyManager(mockSend)

	// Test leaving a party
	err := pm.SendPartyLeave()
	if err != nil {
		t.Fatalf("SendPartyLeave() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["0100"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}
