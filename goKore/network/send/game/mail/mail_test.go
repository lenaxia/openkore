package mail

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForMail is a mock implementation of the Send interface for testing mail functionality
type MockSendForMail struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForMail() *MockSendForMail {
	return &MockSendForMail{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForMail) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForMail) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForMail) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForMail) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForMail) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForMail) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForMail) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForMail) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForMail) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForMail) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForMail) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForMail) SetConnection(conn interface{}) {
}

func (ms *MockSendForMail) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForMail) GetTime() uint32 {
	return 12345
}

// TestNewMailManager tests the NewMailManager function
func TestNewMailManager(t *testing.T) {
	mockSend := NewMockSendForMail()
	mm := NewMailManager(mockSend)

	if mm == nil {
		t.Fatal("NewMailManager() returned nil")
	}

	if mm.baseSend == nil {
		t.Error("mm.baseSend was not set correctly")
	}
}

// TestSendMailboxOpen tests the SendMailboxOpen method
func TestSendMailboxOpen(t *testing.T) {
	mockSend := NewMockSendForMail()
	mockSend.packetLUT["mailbox_open"] = "0260"

	mm := NewMailManager(mockSend)

	// Test opening the mailbox
	err := mm.SendMailboxOpen()
	if err != nil {
		t.Fatalf("SendMailboxOpen() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["0260"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestSendMailRead tests the SendMailRead method
func TestSendMailRead(t *testing.T) {
	mockSend := NewMockSendForMail()
	mockSend.packetLUT["mail_read"] = "0241"

	mm := NewMailManager(mockSend)

	// Test reading a mail
	mailID := uint32(12345)
	err := mm.SendMailRead(mailID)
	if err != nil {
		t.Fatalf("SendMailRead() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0241"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["mailID"] != mailID {
		t.Errorf("args[\"mailID\"] = %v, want %v", args["mailID"], mailID)
	}
}

// TestSendMailDelete tests the SendMailDelete method
func TestSendMailDelete(t *testing.T) {
	mockSend := NewMockSendForMail()
	mockSend.packetLUT["mail_delete"] = "0243"

	mm := NewMailManager(mockSend)

	// Test deleting a mail
	mailID := uint32(12345)
	err := mm.SendMailDelete(mailID)
	if err != nil {
		t.Fatalf("SendMailDelete() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0243"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["mailID"] != mailID {
		t.Errorf("args[\"mailID\"] = %v, want %v", args["mailID"], mailID)
	}
}

// TestSendMailGetAttach tests the SendMailGetAttach method
func TestSendMailGetAttach(t *testing.T) {
	mockSend := NewMockSendForMail()
	mockSend.packetLUT["mail_attachment_get"] = "0244"

	mm := NewMailManager(mockSend)

	// Test getting a mail attachment
	mailID := uint32(12345)
	err := mm.SendMailGetAttach(mailID)
	if err != nil {
		t.Fatalf("SendMailGetAttach() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0244"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["mailID"] != mailID {
		t.Errorf("args[\"mailID\"] = %v, want %v", args["mailID"], mailID)
	}
}

// TestSendMailOperateWindow tests the SendMailOperateWindow method
func TestSendMailOperateWindow(t *testing.T) {
	mockSend := NewMockSendForMail()
	mockSend.packetLUT["mail_remove"] = "0246"

	mm := NewMailManager(mockSend)

	// Test operating the mail window
	flag := uint8(1)
	err := mm.SendMailOperateWindow(flag)
	if err != nil {
		t.Fatalf("SendMailOperateWindow() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0246"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["flag"] != flag {
		t.Errorf("args[\"flag\"] = %v, want %v", args["flag"], flag)
	}
}

// TestSendMailSetAttach tests the SendMailSetAttach method
func TestSendMailSetAttach(t *testing.T) {
	mockSend := NewMockSendForMail()
	mockSend.packetLUT["mail_attachment_set"] = "0247"
	mockSend.packetLUT["mail_remove"] = "0246"

	mm := NewMailManager(mockSend)

	// Test setting a mail attachment
	amount := uint32(1000)
	ID := uint16(12345)
	err := mm.SendMailSetAttach(amount, ID)
	if err != nil {
		t.Fatalf("SendMailSetAttach() returned error: %v", err)
	}

	// Check that the packets were sent
	if len(mockSend.sentPackets) != 2 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 2", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct for the mail_remove packet
	args1, exists := mockSend.reconstructArgs["0246"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct for mail_remove")
	}

	if args1["flag"] != uint8(1) {
		t.Errorf("args1[\"flag\"] = %v, want 1", args1["flag"])
	}

	// Check that the arguments were correct for the mail_attachment_set packet
	args2, exists := mockSend.reconstructArgs["0247"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct for mail_attachment_set")
	}

	if args2["ID"] != ID {
		t.Errorf("args2[\"ID\"] = %v, want %v", args2["ID"], ID)
	}

	if args2["amount"] != amount {
		t.Errorf("args2[\"amount\"] = %v, want %v", args2["amount"], amount)
	}
}

// TestSendMailSetAttachWithZeny tests the SendMailSetAttach method with zeny
func TestSendMailSetAttachWithZeny(t *testing.T) {
	mockSend := NewMockSendForMail()
	mockSend.packetLUT["mail_attachment_set"] = "0247"
	mockSend.packetLUT["mail_remove"] = "0246"

	mm := NewMailManager(mockSend)

	// Test setting a mail attachment with zeny
	amount := uint32(1000)
	ID := uint16(0) // 0 means zeny
	err := mm.SendMailSetAttach(amount, ID)
	if err != nil {
		t.Fatalf("SendMailSetAttach() returned error: %v", err)
	}

	// Check that the packets were sent
	if len(mockSend.sentPackets) != 2 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 2", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct for the mail_remove packet
	args1, exists := mockSend.reconstructArgs["0246"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct for mail_remove")
	}

	if args1["flag"] != uint8(2) {
		t.Errorf("args1[\"flag\"] = %v, want 2", args1["flag"])
	}

	// Check that the arguments were correct for the mail_attachment_set packet
	args2, exists := mockSend.reconstructArgs["0247"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct for mail_attachment_set")
	}

	if args2["ID"] != ID {
		t.Errorf("args2[\"ID\"] = %v, want %v", args2["ID"], ID)
	}

	if args2["amount"] != amount {
		t.Errorf("args2[\"amount\"] = %v, want %v", args2["amount"], amount)
	}
}

// TestReconstructMailSend tests the ReconstructMailSend function
func TestReconstructMailSend(t *testing.T) {
	// Create test data
	body := "Test body"
	bodyLen := len(body)

	// Create args map
	args := map[string]interface{}{
		"body":     body,
		"body_len": bodyLen,
	}

	// Reconstruct the data
	ReconstructMailSend(args)

	// Check that the body was reconstructed correctly
	bodyBytes, ok := args["body"].([]byte)
	if !ok {
		t.Fatal("ReconstructMailSend did not set body field correctly")
	}

	// Check that the body was converted to bytes with null terminator
	if string(bodyBytes[:bodyLen]) != body {
		t.Errorf("args[\"body\"] = %v, want %v", string(bodyBytes[:bodyLen]), body)
	}

	if len(bodyBytes) != bodyLen+1 {
		t.Errorf("len(args[\"body\"]) = %v, want %v", len(bodyBytes), bodyLen+1)
	}

	if bodyBytes[bodyLen] != 0 {
		t.Errorf("args[\"body\"] does not end with null terminator")
	}
}

// TestSendMailSend tests the SendMailSend method
func TestSendMailSend(t *testing.T) {
	mockSend := NewMockSendForMail()
	mockSend.packetLUT["mail_send"] = "0248"

	mm := NewMailManager(mockSend)

	// Test sending a mail
	receiver := "TestReceiver"
	title := "Test Title"
	message := "Test Message"
	err := mm.SendMailSend(receiver, title, message)
	if err != nil {
		t.Fatalf("SendMailSend() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0248"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	// Check that receiver was converted to bytes
	receiverBytes, ok := args["recipient"].([]byte)
	if !ok {
		t.Errorf("args[\"recipient\"] is not a byte slice")
	} else if string(receiverBytes) != receiver {
		t.Errorf("args[\"recipient\"] = %v, want %v", string(receiverBytes), receiver)
	}

	// Check that title was converted to bytes
	titleBytes, ok := args["title"].([]byte)
	if !ok {
		t.Errorf("args[\"title\"] is not a byte slice")
	} else if string(titleBytes) != title {
		t.Errorf("args[\"title\"] = %v, want %v", string(titleBytes), title)
	}

	if args["body_len"] != len(message) {
		t.Errorf("args[\"body_len\"] = %v, want %v", args["body_len"], len(message))
	}

	// Check that body was converted to bytes with null terminator
	bodyBytes, ok := args["body"].([]byte)
	if !ok {
		t.Errorf("args[\"body\"] is not a byte slice")
	} else if string(bodyBytes[:len(message)]) != message {
		t.Errorf("args[\"body\"] = %v, want %v", string(bodyBytes[:len(message)]), message)
	} else if bodyBytes[len(bodyBytes)-1] != 0 {
		t.Errorf("args[\"body\"] does not end with null terminator")
	}
}

// TestSendMailReturn tests the SendMailReturn method
func TestSendMailReturn(t *testing.T) {
	mockSend := NewMockSendForMail()
	mockSend.packetLUT["mail_return"] = "0273"

	mm := NewMailManager(mockSend)

	// Test returning a mail
	mailID := uint32(12345)
	sender := uint32(67890)
	err := mm.SendMailReturn(mailID, sender)
	if err != nil {
		t.Fatalf("SendMailReturn() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0273"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["mailID"] != mailID {
		t.Errorf("args[\"mailID\"] = %v, want %v", args["mailID"], mailID)
	}

	if args["sender"] != sender {
		t.Errorf("args[\"sender\"] = %v, want %v", args["sender"], sender)
	}
}
