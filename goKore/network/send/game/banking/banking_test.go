package banking

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForBanking is a mock implementation of the Send interface for testing banking functionality
type MockSendForBanking struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForBanking() *MockSendForBanking {
	return &MockSendForBanking{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForBanking) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForBanking) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForBanking) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForBanking) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForBanking) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForBanking) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForBanking) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForBanking) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForBanking) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForBanking) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForBanking) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForBanking) SetConnection(conn interface{}) {
}

func (ms *MockSendForBanking) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForBanking) GetTime() uint32 {
	return 12345
}

// TestNewBankingManager tests the NewBankingManager function
func TestNewBankingManager(t *testing.T) {
	mockSend := NewMockSendForBanking()
	bm := NewBankingManager(mockSend)

	if bm == nil {
		t.Fatal("NewBankingManager() returned nil")
	}

	if bm.baseSend == nil {
		t.Error("bm.baseSend was not set correctly")
	}
}

// TestSendBankingCheck tests the SendBankingCheck method
func TestSendBankingCheck(t *testing.T) {
	mockSend := NewMockSendForBanking()
	mockSend.packetLUT["banking_check_request"] = "09AB"

	bm := NewBankingManager(mockSend)

	// Test checking banking data
	accountID := uint32(12345)
	err := bm.SendBankingCheck(accountID)
	if err != nil {
		t.Fatalf("SendBankingCheck() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["09AB"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["accountID"] != accountID {
		t.Errorf("args[\"accountID\"] = %v, want %v", args["accountID"], accountID)
	}
}

// TestSendBankingWithdraw tests the SendBankingWithdraw method
func TestSendBankingWithdraw(t *testing.T) {
	mockSend := NewMockSendForBanking()
	mockSend.packetLUT["banking_withdraw_request"] = "09A9"

	bm := NewBankingManager(mockSend)

	// Test withdrawing from bank
	accountID := uint32(12345)
	zeny := uint32(1000)
	err := bm.SendBankingWithdraw(accountID, zeny)
	if err != nil {
		t.Fatalf("SendBankingWithdraw() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["09A9"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["accountID"] != accountID {
		t.Errorf("args[\"accountID\"] = %v, want %v", args["accountID"], accountID)
	}

	if args["zeny"] != zeny {
		t.Errorf("args[\"zeny\"] = %v, want %v", args["zeny"], zeny)
	}
}

// TestSendBankingDeposit tests the SendBankingDeposit method
func TestSendBankingDeposit(t *testing.T) {
	mockSend := NewMockSendForBanking()
	mockSend.packetLUT["banking_deposit_request"] = "09A7"

	bm := NewBankingManager(mockSend)

	// Test depositing to bank
	accountID := uint32(12345)
	zeny := uint32(1000)
	err := bm.SendBankingDeposit(accountID, zeny)
	if err != nil {
		t.Fatalf("SendBankingDeposit() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["09A7"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["accountID"] != accountID {
		t.Errorf("args[\"accountID\"] = %v, want %v", args["accountID"], accountID)
	}

	if args["zeny"] != zeny {
		t.Errorf("args[\"zeny\"] = %v, want %v", args["zeny"], zeny)
	}
}
