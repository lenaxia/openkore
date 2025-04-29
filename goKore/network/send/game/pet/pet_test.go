package pet

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
)

// MockSend is a mock implementation of the core.Send interface for testing
type MockSend struct {
	packetIDs      map[string]string
	reconstructed  []byte
	sent           []byte
	time           uint32
	lastPacketName string
	lastArgs       map[string]interface{}
}

// NewMockSend creates a new MockSend instance with default values
func NewMockSend() *MockSend {
	return &MockSend{
		packetIDs: map[string]string{
			"pet_capture":   "019F",
			"pet_menu":      "01A1",
			"pet_hatch":     "01A7",
			"pet_name":      "01A5",
			"pet_emotion":   "01A9",
			"pet_evolution": "09FB",
		},
		time:     12345,
		lastArgs: make(map[string]interface{}),
	}
}

// SendToServer mocks sending a packet to the server
func (ms *MockSend) SendToServer(msg []byte) error {
	ms.sent = msg
	return nil
}

// EncryptMessageID mocks encrypting a message ID
func (ms *MockSend) EncryptMessageID(msg *[]byte) error {
	return nil
}

// CryptKeys mocks setting encryption keys
func (ms *MockSend) CryptKeys(key1, key2, key3 uint32) {}

// PinEncode mocks encoding a PIN
func (ms *MockSend) PinEncode(seed, pin int) string {
	return ""
}

// InjectMessage mocks injecting a message
func (ms *MockSend) InjectMessage(message string) error {
	return nil
}

// InjectAdminMessage mocks injecting an admin message
func (ms *MockSend) InjectAdminMessage(message string) error {
	return nil
}

// SendRaw mocks sending a raw packet
func (ms *MockSend) SendRaw(raw string) error {
	return nil
}

// Reconstruct mocks reconstructing a packet
func (ms *MockSend) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the last packet name and arguments for testing
	for name, id := range ms.packetIDs {
		if id == packetID {
			ms.lastPacketName = name
			break
		}
	}

	// Store the arguments for testing
	ms.lastArgs = args

	// Simple mock implementation that just returns the packet ID as bytes
	ms.reconstructed = []byte{0x00, 0x00}
	return ms.reconstructed, nil
}

// GetPacketID mocks getting a packet ID by name
func (ms *MockSend) GetPacketID(name string) (string, bool) {
	id, ok := ms.packetIDs[name]
	if ok {
		ms.lastPacketName = name
	}
	return id, ok
}

// RegisterPacketHandler mocks registering a packet handler
func (ms *MockSend) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
}

// RegisterHook mocks registering a hook
func (ms *MockSend) RegisterHook(hookName string, callback hooks.HookCallback) {}

// SetConnection mocks setting a connection
func (ms *MockSend) SetConnection(conn interface{}) {}

// GetConnection mocks getting a connection
func (ms *MockSend) GetConnection() interface{} {
	return nil
}

// GetTime mocks getting the current time
func (ms *MockSend) GetTime() uint32 {
	return ms.time
}

// LastPacketID returns the name of the last packet that was requested
func (ms *MockSend) LastPacketID() (string, bool) {
	if ms.lastPacketName == "" {
		return "", false
	}
	return ms.lastPacketName, true
}

// LastArgs returns the arguments of the last packet that was reconstructed
func (ms *MockSend) LastArgs() map[string]interface{} {
	return ms.lastArgs
}

// TestNewPetManager tests the NewPetManager function
func TestNewPetManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	petManager := NewPetManager(mockSend)

	if petManager == nil {
		t.Fatal("NewPetManager() returned nil")
	}

	if petManager.baseSend == nil {
		t.Error("petManager.baseSend was not set correctly")
	}
}

// TestSendPetCapture tests the SendPetCapture function
func TestSendPetCapture(t *testing.T) {
	mockSend := NewMockSend()
	petManager := NewPetManager(mockSend)

	// Test sending a pet capture command
	monsterID := uint32(12345)
	err := petManager.SendPetCapture(monsterID)
	if err != nil {
		t.Fatalf("SendPetCapture() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "pet_capture" {
		t.Errorf("Expected packet ID 'pet_capture', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != monsterID {
		t.Errorf("Expected ID=%d, got %v", monsterID, mockSend.LastArgs()["ID"])
	}
}

// TestSendPetMenu tests the SendPetMenu function
func TestSendPetMenu(t *testing.T) {
	mockSend := NewMockSend()
	petManager := NewPetManager(mockSend)

	// Test cases for different menu actions
	testCases := []struct {
		action int
		name   string
	}{
		{0, "info"},
		{1, "feed"},
		{2, "performance"},
		{3, "return to egg"},
		{4, "unequip accessory"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := petManager.SendPetMenu(tc.action)
			if err != nil {
				t.Fatalf("SendPetMenu(%d) returned error: %v", tc.action, err)
			}

			if mockSend.sent == nil {
				t.Fatal("No packet was sent")
			}

			// Check that the correct packet ID was used
			if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "pet_menu" {
				t.Errorf("Expected packet ID 'pet_menu', got '%s'", packetID)
			}

			// Check that the correct arguments were used
			if actionVal, ok := mockSend.LastArgs()["action"].(uint8); !ok || int(actionVal) != tc.action {
				t.Errorf("Expected action=%d, got %v", tc.action, mockSend.LastArgs()["action"])
			}
		})
	}

	// Test with invalid action
	err := petManager.SendPetMenu(5)
	if err == nil {
		t.Error("Expected error for invalid action, got nil")
	}
}

// TestSendPetHatch tests the SendPetHatch function
func TestSendPetHatch(t *testing.T) {
	mockSend := NewMockSend()
	petManager := NewPetManager(mockSend)

	// Test sending a pet hatch command
	itemID := uint32(7778) // Example pet egg item ID
	err := petManager.SendPetHatch(itemID)
	if err != nil {
		t.Fatalf("SendPetHatch() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "pet_hatch" {
		t.Errorf("Expected packet ID 'pet_hatch', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != itemID {
		t.Errorf("Expected ID=%d, got %v", itemID, mockSend.LastArgs()["ID"])
	}
}

// TestSendPetName tests the SendPetName function
func TestSendPetName(t *testing.T) {
	mockSend := NewMockSend()
	petManager := NewPetManager(mockSend)

	// Test sending a pet name command
	petName := "Fluffy"
	err := petManager.SendPetName(petName)
	if err != nil {
		t.Fatalf("SendPetName() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "pet_name" {
		t.Errorf("Expected packet ID 'pet_name', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if nameVal, ok := mockSend.LastArgs()["name"].(string); !ok || nameVal != petName {
		t.Errorf("Expected name=%s, got %v", petName, mockSend.LastArgs()["name"])
	}

	// Test with empty name
	err = petManager.SendPetName("")
	if err == nil {
		t.Error("Expected error for empty name, got nil")
	}

	// Test with too long name
	longName := "ThisNameIsTooLongForAPet"
	err = petManager.SendPetName(longName)
	if err == nil {
		t.Error("Expected error for too long name, got nil")
	}
}

// TestSendPetEmotion tests the SendPetEmotion function
func TestSendPetEmotion(t *testing.T) {
	mockSend := NewMockSend()
	petManager := NewPetManager(mockSend)

	// Test sending a pet emotion command
	emotionID := uint8(4) // Example emotion ID
	err := petManager.SendPetEmotion(emotionID)
	if err != nil {
		t.Fatalf("SendPetEmotion() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "pet_emotion" {
		t.Errorf("Expected packet ID 'pet_emotion', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint8); !ok || idVal != emotionID {
		t.Errorf("Expected ID=%d, got %v", emotionID, mockSend.LastArgs()["ID"])
	}
}

// TestParsePetEvolution tests the ParsePetEvolution function
func TestParsePetEvolution(t *testing.T) {
	mockSend := NewMockSend()
	petManager := NewPetManager(mockSend)

	// Create test data
	itemInfo := []byte{
		0x01, 0x00, 0x05, 0x00, // Item 1: index 1, amount 5
		0x02, 0x00, 0x03, 0x00, // Item 2: index 2, amount 3
	}

	args := map[string]interface{}{
		"itemInfo": itemInfo,
	}

	// Parse the pet evolution data
	err := petManager.ParsePetEvolution(args)
	if err != nil {
		t.Fatalf("ParsePetEvolution() returned error: %v", err)
	}

	// Check that the items were parsed correctly
	items, ok := args["items"].([]map[string]uint16)
	if !ok {
		t.Fatal("items not parsed correctly")
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	// Check first item
	if items[0]["itemIndex"] != 1 || items[0]["amount"] != 5 {
		t.Errorf("First item not parsed correctly, got %v", items[0])
	}

	// Check second item
	if items[1]["itemIndex"] != 2 || items[1]["amount"] != 3 {
		t.Errorf("Second item not parsed correctly, got %v", items[1])
	}
}

// TestReconstructPetEvolution tests the ReconstructPetEvolution function
func TestReconstructPetEvolution(t *testing.T) {
	mockSend := NewMockSend()
	petManager := NewPetManager(mockSend)

	// Create test data
	items := []map[string]uint16{
		{"itemIndex": 1, "amount": 5},
		{"itemIndex": 2, "amount": 3},
	}

	args := map[string]interface{}{
		"items": items,
	}

	// Reconstruct the pet evolution data
	err := petManager.ReconstructPetEvolution(args)
	if err != nil {
		t.Fatalf("ReconstructPetEvolution() returned error: %v", err)
	}

	// Check that the itemInfo was reconstructed correctly
	itemInfo, ok := args["itemInfo"].([]byte)
	if !ok {
		t.Fatal("itemInfo not reconstructed correctly")
	}

	// Expected binary data
	expected := []byte{
		0x01, 0x00, 0x05, 0x00, // Item 1: index 1, amount 5
		0x02, 0x00, 0x03, 0x00, // Item 2: index 2, amount 3
	}

	if len(itemInfo) != len(expected) {
		t.Fatalf("Expected itemInfo length %d, got %d", len(expected), len(itemInfo))
	}

	for i := 0; i < len(expected); i++ {
		if itemInfo[i] != expected[i] {
			t.Errorf("itemInfo[%d] = %d, expected %d", i, itemInfo[i], expected[i])
		}
	}
}

// TestSendPetEvolution tests the SendPetEvolution function
func TestSendPetEvolution(t *testing.T) {
	mockSend := NewMockSend()
	petManager := NewPetManager(mockSend)

	// Test sending a pet evolution command
	petEggID := uint32(9001)
	items := []map[string]uint16{
		{"itemIndex": 1, "amount": 5},
		{"itemIndex": 2, "amount": 3},
	}

	err := petManager.SendPetEvolution(petEggID, items)
	if err != nil {
		t.Fatalf("SendPetEvolution() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "pet_evolution" {
		t.Errorf("Expected packet ID 'pet_evolution', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != petEggID {
		t.Errorf("Expected ID=%d, got %v", petEggID, mockSend.LastArgs()["ID"])
	}

	// Check that the items were passed correctly
	if itemsVal, ok := mockSend.LastArgs()["items"].([]map[string]uint16); !ok {
		t.Errorf("Expected items to be []map[string]uint16, got %T", mockSend.LastArgs()["items"])
	} else {
		if len(itemsVal) != len(items) {
			t.Errorf("Expected %d items, got %d", len(items), len(itemsVal))
		}
	}

	// Test with empty items
	err = petManager.SendPetEvolution(petEggID, []map[string]uint16{})
	if err == nil {
		t.Error("Expected error for empty items, got nil")
	}
}
