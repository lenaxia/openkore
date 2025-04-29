package rodex

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForRodex is a mock implementation of the Send interface for testing rodex functionality
type MockSendForRodex struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForRodex() *MockSendForRodex {
	return &MockSendForRodex{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForRodex) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForRodex) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForRodex) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForRodex) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForRodex) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForRodex) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForRodex) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForRodex) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForRodex) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForRodex) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForRodex) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForRodex) SetConnection(conn interface{}) {
}

func (ms *MockSendForRodex) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForRodex) GetTime() uint32 {
	return 12345
}

// TestNewRodexManager tests the NewRodexManager function
func TestNewRodexManager(t *testing.T) {
	mockSend := NewMockSendForRodex()
	rm := NewRodexManager(mockSend)

	if rm == nil {
		t.Fatal("NewRodexManager() returned nil")
	}

	if rm.baseSend == nil {
		t.Error("rm.baseSend was not set correctly")
	}
}

// TestRodexDeleteMail tests the RodexDeleteMail method
func TestRodexDeleteMail(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_delete_mail"] = "0A51"

	rm := NewRodexManager(mockSend)

	// Test deleting a mail
	type_ := uint8(1)
	mailID1 := uint32(12345)
	mailID2 := uint32(67890)
	err := rm.RodexDeleteMail(type_, mailID1, mailID2)
	if err != nil {
		t.Fatalf("RodexDeleteMail() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A51"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["type"] != type_ {
		t.Errorf("args[\"type\"] = %v, want %v", args["type"], type_)
	}

	if args["mailID1"] != mailID1 {
		t.Errorf("args[\"mailID1\"] = %v, want %v", args["mailID1"], mailID1)
	}

	if args["mailID2"] != mailID2 {
		t.Errorf("args[\"mailID2\"] = %v, want %v", args["mailID2"], mailID2)
	}
}

// TestRodexRequestZeny tests the RodexRequestZeny method
func TestRodexRequestZeny(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_request_zeny"] = "0A55"

	rm := NewRodexManager(mockSend)

	// Test requesting zeny from a mail
	mailID1 := uint32(12345)
	mailID2 := uint32(67890)
	type_ := uint8(1)
	err := rm.RodexRequestZeny(mailID1, mailID2, type_)
	if err != nil {
		t.Fatalf("RodexRequestZeny() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A55"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["mailID1"] != mailID1 {
		t.Errorf("args[\"mailID1\"] = %v, want %v", args["mailID1"], mailID1)
	}

	if args["mailID2"] != mailID2 {
		t.Errorf("args[\"mailID2\"] = %v, want %v", args["mailID2"], mailID2)
	}

	if args["type"] != type_ {
		t.Errorf("args[\"type\"] = %v, want %v", args["type"], type_)
	}
}

// TestRodexRequestItems tests the RodexRequestItems method
func TestRodexRequestItems(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_request_items"] = "0A4B"

	rm := NewRodexManager(mockSend)

	// Test requesting items from a mail
	mailID1 := uint32(12345)
	mailID2 := uint32(67890)
	type_ := uint8(1)
	err := rm.RodexRequestItems(mailID1, mailID2, type_)
	if err != nil {
		t.Fatalf("RodexRequestItems() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A4B"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["mailID1"] != mailID1 {
		t.Errorf("args[\"mailID1\"] = %v, want %v", args["mailID1"], mailID1)
	}

	if args["mailID2"] != mailID2 {
		t.Errorf("args[\"mailID2\"] = %v, want %v", args["mailID2"], mailID2)
	}

	if args["type"] != type_ {
		t.Errorf("args[\"type\"] = %v, want %v", args["type"], type_)
	}
}

// TestRodexCancelWriteMail tests the RodexCancelWriteMail method
func TestRodexCancelWriteMail(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_cancel_write_mail"] = "0A52"

	rm := NewRodexManager(mockSend)

	// Test canceling writing a mail
	err := rm.RodexCancelWriteMail()
	if err != nil {
		t.Fatalf("RodexCancelWriteMail() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["0A52"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestRodexAddItem tests the RodexAddItem method
func TestRodexAddItem(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_add_item"] = "0A4C"

	rm := NewRodexManager(mockSend)

	// Test adding an item to a mail
	ID := uint16(12345)
	amount := uint16(10)
	err := rm.RodexAddItem(ID, amount)
	if err != nil {
		t.Fatalf("RodexAddItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A4C"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != ID {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], ID)
	}

	if args["amount"] != amount {
		t.Errorf("args[\"amount\"] = %v, want %v", args["amount"], amount)
	}
}

// TestRodexRemoveItem tests the RodexRemoveItem method
func TestRodexRemoveItem(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_remove_item"] = "0A4D"

	rm := NewRodexManager(mockSend)

	// Test removing an item from a mail
	ID := uint16(12345)
	amount := uint16(10)
	err := rm.RodexRemoveItem(ID, amount)
	if err != nil {
		t.Fatalf("RodexRemoveItem() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A4D"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["ID"] != ID {
		t.Errorf("args[\"ID\"] = %v, want %v", args["ID"], ID)
	}

	if args["amount"] != amount {
		t.Errorf("args[\"amount\"] = %v, want %v", args["amount"], amount)
	}
}

// TestRodexOpenWriteMail tests the RodexOpenWriteMail method
func TestRodexOpenWriteMail(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_open_write_mail"] = "0A41"

	rm := NewRodexManager(mockSend)

	// Test opening a mail to write
	name := "TestPlayer"
	err := rm.RodexOpenWriteMail(name)
	if err != nil {
		t.Fatalf("RodexOpenWriteMail() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A41"]
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
}

// TestRodexCheckname tests the RodexCheckname method
func TestRodexCheckname(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_checkname"] = "0A4E"

	rm := NewRodexManager(mockSend)

	// Test checking a name
	name := "TestPlayer"
	err := rm.RodexCheckname(name)
	if err != nil {
		t.Fatalf("RodexCheckname() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A4E"]
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
}

// TestRodexSendMail tests the RodexSendMail method
func TestRodexSendMail(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_send_mail"] = "0A43"

	rm := NewRodexManager(mockSend)

	// Test sending a mail
	receiver := "TestReceiver"
	sender := "TestSender"
	zeny := uint32(1000)
	title := "Test Title"
	body := "Test Body"
	charID := uint32(12345)
	err := rm.RodexSendMail(receiver, sender, zeny, title, body, charID)
	if err != nil {
		t.Fatalf("RodexSendMail() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A43"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["receiver"] != receiver {
		t.Errorf("args[\"receiver\"] = %v, want %v", args["receiver"], receiver)
	}

	// Check that sender was converted to bytes
	senderBytes, ok := args["sender"].([]byte)
	if !ok {
		t.Errorf("args[\"sender\"] is not a byte slice")
	} else if string(senderBytes) != sender {
		t.Errorf("args[\"sender\"] = %v, want %v", string(senderBytes), sender)
	}

	if args["zeny1"] != zeny {
		t.Errorf("args[\"zeny1\"] = %v, want %v", args["zeny1"], zeny)
	}

	if args["zeny2"] != uint32(0) {
		t.Errorf("args[\"zeny2\"] = %v, want 0", args["zeny2"])
	}

	if args["char_id"] != charID {
		t.Errorf("args[\"char_id\"] = %v, want %v", args["char_id"], charID)
	}

	// Check that title was converted to bytes with null terminator
	titleBytes, ok := args["title"].([]byte)
	if !ok {
		t.Errorf("args[\"title\"] is not a byte slice")
	} else if string(titleBytes[:len(title)]) != title {
		t.Errorf("args[\"title\"] = %v, want %v", string(titleBytes[:len(title)]), title)
	} else if titleBytes[len(titleBytes)-1] != 0 {
		t.Errorf("args[\"title\"] does not end with null terminator")
	}

	// Check that body was converted to bytes with null terminator
	bodyBytes, ok := args["body"].([]byte)
	if !ok {
		t.Errorf("args[\"body\"] is not a byte slice")
	} else if string(bodyBytes[:len(body)]) != body {
		t.Errorf("args[\"body\"] = %v, want %v", string(bodyBytes[:len(body)]), body)
	} else if bodyBytes[len(bodyBytes)-1] != 0 {
		t.Errorf("args[\"body\"] does not end with null terminator")
	}

	if args["title_len"] != uint16(len(titleBytes)) {
		t.Errorf("args[\"title_len\"] = %v, want %v", args["title_len"], len(titleBytes))
	}

	if args["body_len"] != uint16(len(bodyBytes)) {
		t.Errorf("args[\"body_len\"] = %v, want %v", args["body_len"], len(bodyBytes))
	}
}

// TestRodexRefreshMaillist tests the RodexRefreshMaillist method
func TestRodexRefreshMaillist(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_refresh_maillist"] = "0A49"

	rm := NewRodexManager(mockSend)

	// Test refreshing the mail list
	type_ := uint8(1)
	mailID1 := uint32(12345)
	mailID2 := uint32(67890)
	err := rm.RodexRefreshMaillist(type_, mailID1, mailID2)
	if err != nil {
		t.Fatalf("RodexRefreshMaillist() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A49"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["type"] != type_ {
		t.Errorf("args[\"type\"] = %v, want %v", args["type"], type_)
	}

	if args["mailID1"] != mailID1 {
		t.Errorf("args[\"mailID1\"] = %v, want %v", args["mailID1"], mailID1)
	}

	if args["mailID2"] != mailID2 {
		t.Errorf("args[\"mailID2\"] = %v, want %v", args["mailID2"], mailID2)
	}

	if args["mailReturnID1"] != uint32(0) {
		t.Errorf("args[\"mailReturnID1\"] = %v, want 0", args["mailReturnID1"])
	}

	if args["mailReturnID2"] != uint32(0) {
		t.Errorf("args[\"mailReturnID2\"] = %v, want 0", args["mailReturnID2"])
	}

	if args["mailAccountID1"] != uint32(0) {
		t.Errorf("args[\"mailAccountID1\"] = %v, want 0", args["mailAccountID1"])
	}

	if args["mailAccountID2"] != uint32(0) {
		t.Errorf("args[\"mailAccountID2\"] = %v, want 0", args["mailAccountID2"])
	}
}

// TestRodexReadMail tests the RodexReadMail method
func TestRodexReadMail(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_read_mail"] = "0A4A"

	rm := NewRodexManager(mockSend)

	// Test reading a mail
	type_ := uint8(1)
	mailID1 := uint32(12345)
	mailID2 := uint32(67890)
	err := rm.RodexReadMail(type_, mailID1, mailID2)
	if err != nil {
		t.Fatalf("RodexReadMail() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A4A"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["type"] != type_ {
		t.Errorf("args[\"type\"] = %v, want %v", args["type"], type_)
	}

	if args["mailID1"] != mailID1 {
		t.Errorf("args[\"mailID1\"] = %v, want %v", args["mailID1"], mailID1)
	}

	if args["mailID2"] != mailID2 {
		t.Errorf("args[\"mailID2\"] = %v, want %v", args["mailID2"], mailID2)
	}
}

// TestRodexNextMaillist tests the RodexNextMaillist method
func TestRodexNextMaillist(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_next_maillist"] = "0A4F"

	rm := NewRodexManager(mockSend)

	// Test getting the next mail list
	type_ := uint8(1)
	mailID1 := uint32(12345)
	mailID2 := uint32(67890)
	err := rm.RodexNextMaillist(type_, mailID1, mailID2)
	if err != nil {
		t.Fatalf("RodexNextMaillist() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A4F"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["type"] != type_ {
		t.Errorf("args[\"type\"] = %v, want %v", args["type"], type_)
	}

	if args["mailID1"] != mailID1 {
		t.Errorf("args[\"mailID1\"] = %v, want %v", args["mailID1"], mailID1)
	}

	if args["mailID2"] != mailID2 {
		t.Errorf("args[\"mailID2\"] = %v, want %v", args["mailID2"], mailID2)
	}
}

// TestRodexOpenMailbox tests the RodexOpenMailbox method
func TestRodexOpenMailbox(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_open_mailbox"] = "0A40"

	rm := NewRodexManager(mockSend)

	// Test opening the mailbox
	type_ := uint8(1)
	mailID1 := uint32(12345)
	mailID2 := uint32(67890)
	err := rm.RodexOpenMailbox(type_, mailID1, mailID2)
	if err != nil {
		t.Fatalf("RodexOpenMailbox() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0A40"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["type"] != type_ {
		t.Errorf("args[\"type\"] = %v, want %v", args["type"], type_)
	}

	if args["mailID1"] != mailID1 {
		t.Errorf("args[\"mailID1\"] = %v, want %v", args["mailID1"], mailID1)
	}

	if args["mailID2"] != mailID2 {
		t.Errorf("args[\"mailID2\"] = %v, want %v", args["mailID2"], mailID2)
	}

	if args["mailReturnID1"] != uint32(0) {
		t.Errorf("args[\"mailReturnID1\"] = %v, want 0", args["mailReturnID1"])
	}

	if args["mailReturnID2"] != uint32(0) {
		t.Errorf("args[\"mailReturnID2\"] = %v, want 0", args["mailReturnID2"])
	}

	if args["mailAccountID1"] != uint32(0) {
		t.Errorf("args[\"mailAccountID1\"] = %v, want 0", args["mailAccountID1"])
	}

	if args["mailAccountID2"] != uint32(0) {
		t.Errorf("args[\"mailAccountID2\"] = %v, want 0", args["mailAccountID2"])
	}
}

// TestRodexCloseMailbox tests the RodexCloseMailbox method
func TestRodexCloseMailbox(t *testing.T) {
	mockSend := NewMockSendForRodex()
	mockSend.packetLUT["rodex_close_mailbox"] = "0A42"

	rm := NewRodexManager(mockSend)

	// Test closing the mailbox
	err := rm.RodexCloseMailbox()
	if err != nil {
		t.Fatalf("RodexCloseMailbox() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["0A42"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}
