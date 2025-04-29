package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSendForStorage is a mock implementation of the Send interface for testing storage functionality
type MockSendForStorage struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
}

func NewMockSendForStorage() *MockSendForStorage {
	return &MockSendForStorage{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
	}
}

func (ms *MockSendForStorage) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSendForStorage) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSendForStorage) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSendForStorage) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSendForStorage) InjectMessage(message string) error {
	return nil
}

func (ms *MockSendForStorage) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSendForStorage) SendRaw(raw string) error {
	return nil
}

func (ms *MockSendForStorage) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSendForStorage) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSendForStorage) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSendForStorage) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSendForStorage) SetConnection(conn interface{}) {
}

func (ms *MockSendForStorage) GetConnection() interface{} {
	return nil
}

func (ms *MockSendForStorage) GetTime() uint32 {
	return 12345
}

// TestNewStorageManager tests the NewStorageManager function
func TestNewStorageManager(t *testing.T) {
	mockSend := NewMockSendForStorage()
	sm := NewStorageManager(mockSend)

	if sm == nil {
		t.Fatal("NewStorageManager() returned nil")
	}

	if sm.baseSend == nil {
		t.Error("sm.baseSend was not set correctly")
	}
}

// TestOpenStorage tests the OpenStorage method
func TestOpenStorage(t *testing.T) {
	mockSend := NewMockSendForStorage()
	mockSend.packetLUT["open_storage"] = "00F3"

	sm := NewStorageManager(mockSend)

	// Test opening storage
	err := sm.OpenStorage()
	if err != nil {
		t.Fatalf("OpenStorage() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	_, exists := mockSend.reconstructArgs["00F3"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestCloseStorage tests the CloseStorage method
func TestCloseStorage(t *testing.T) {
	mockSend := NewMockSendForStorage()
	mockSend.packetLUT["close_storage"] = "00F7"

	sm := NewStorageManager(mockSend)

	// Test closing storage
	err := sm.CloseStorage()
	if err != nil {
		t.Fatalf("CloseStorage() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	_, exists := mockSend.reconstructArgs["00F7"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}

// TestMoveToStorage tests the MoveToStorage method
func TestMoveToStorage(t *testing.T) {
	mockSend := NewMockSendForStorage()
	mockSend.packetLUT["move_to_storage"] = "00F3"

	sm := NewStorageManager(mockSend)

	// Test moving an item to storage
	index := uint16(1)
	amount := uint16(10)
	err := sm.MoveToStorage(index, amount)
	if err != nil {
		t.Fatalf("MoveToStorage() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00F3"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["index"] != index {
		t.Errorf("args[\"index\"] = %v, want %v", args["index"], index)
	}

	if args["amount"] != amount {
		t.Errorf("args[\"amount\"] = %v, want %v", args["amount"], amount)
	}
}

// TestMoveFromStorage tests the MoveFromStorage method
func TestMoveFromStorage(t *testing.T) {
	mockSend := NewMockSendForStorage()
	mockSend.packetLUT["move_from_storage"] = "00F5"

	sm := NewStorageManager(mockSend)

	// Test moving an item from storage
	index := uint16(1)
	amount := uint16(10)
	err := sm.MoveFromStorage(index, amount)
	if err != nil {
		t.Fatalf("MoveFromStorage() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["00F5"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["index"] != index {
		t.Errorf("args[\"index\"] = %v, want %v", args["index"], index)
	}

	if args["amount"] != amount {
		t.Errorf("args[\"amount\"] = %v, want %v", args["amount"], amount)
	}
}

// TestArrangeStorage tests the ArrangeStorage method
func TestArrangeStorage(t *testing.T) {
	mockSend := NewMockSendForStorage()
	mockSend.packetLUT["arrange_storage"] = "00F8"

	sm := NewStorageManager(mockSend)

	// Test arranging storage
	err := sm.ArrangeStorage()
	if err != nil {
		t.Fatalf("ArrangeStorage() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	_, exists := mockSend.reconstructArgs["00F8"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}
}
