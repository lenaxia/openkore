package gm

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
			"gm_summon_player":        "0094",
			"gm_kick":                 "00CC",
			"gm_kick_all":             "00CD",
			"gm_item_mob_create":      "013F",
			"gm_move_to_map":          "01BD", // Changed from 01B9
			"gm_reset_state_skill":    "0196", // Changed from 0197
			"gm_change_cell_type":     "019C", // Changed from 0198
			"gm_change_effect_state":  "019B",
			"gm_remove":               "01BE", // Changed from 01BA
			"gm_shift":                "01BB",
			"gm_recall":               "01BC",
			"manner_by_name":          "0212",
			"gm_request_status":       "0213",
			"gm_request_account_name": "01DF",
			"ban_check":               "0090",
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

// TestNewGMManager tests the NewGMManager function
func TestNewGMManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	if gmManager == nil {
		t.Fatal("NewGMManager() returned nil")
	}

	if gmManager.baseSend == nil {
		t.Error("gmManager.baseSend was not set correctly")
	}
}

// TestSendGMSummon tests the SendGMSummon function
func TestSendGMSummon(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test sending a GM summon request
	playerName := "TestPlayer"
	err := gmManager.SendGMSummon(playerName)
	if err != nil {
		t.Fatalf("SendGMSummon() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "gm_summon_player" {
		t.Errorf("Expected packet ID 'gm_summon_player', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if nameVal, ok := mockSend.LastArgs()["playerName"].([]byte); !ok {
		t.Errorf("Expected playerName to be []byte, got %T", mockSend.LastArgs()["playerName"])
	} else {
		nameStr := string(nameVal)
		if nameStr != playerName {
			t.Errorf("Expected playerName=%s, got %s", playerName, nameStr)
		}
	}
}

// TestSendGMKick tests the SendGMKick function
func TestSendGMKick(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test sending a GM kick request
	accountID := uint32(12345)
	err := gmManager.SendGMKick(accountID)
	if err != nil {
		t.Fatalf("SendGMKick() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "gm_kick" {
		t.Errorf("Expected packet ID 'gm_kick', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["targetAccountID"].(uint32); !ok || idVal != accountID {
		t.Errorf("Expected targetAccountID=%d, got %v", accountID, mockSend.LastArgs()["targetAccountID"])
	}
}

// TestSendGMKickAll tests the SendGMKickAll function
func TestSendGMKickAll(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test sending a GM kick all request
	err := gmManager.SendGMKickAll()
	if err != nil {
		t.Fatalf("SendGMKickAll() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "gm_kick_all" {
		t.Errorf("Expected packet ID 'gm_kick_all', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendGMMonsterItem tests the SendGMMonsterItem function
func TestSendGMMonsterItem(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test sending a GM monster/item creation request
	name := "Poring"
	err := gmManager.SendGMMonsterItem(name)
	if err != nil {
		t.Fatalf("SendGMMonsterItem() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "gm_item_mob_create" {
		t.Errorf("Expected packet ID 'gm_item_mob_create', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if nameVal, ok := mockSend.LastArgs()["name"].([]byte); !ok {
		t.Errorf("Expected name to be []byte, got %T", mockSend.LastArgs()["name"])
	} else {
		nameStr := string(nameVal)
		if nameStr != name {
			t.Errorf("Expected name=%s, got %s", name, nameStr)
		}
	}
}

// TestSendGMMapMove tests the SendGMMapMove function
func TestSendGMMapMove(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test sending a GM map move request
	mapName := "prontera"
	x := uint16(150)
	y := uint16(150)
	err := gmManager.SendGMMapMove(mapName, x, y)
	if err != nil {
		t.Fatalf("SendGMMapMove() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "gm_move_to_map" {
		t.Errorf("Expected packet ID 'gm_move_to_map', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if nameVal, ok := mockSend.LastArgs()["mapName"].([]byte); !ok {
		t.Errorf("Expected mapName to be []byte, got %T", mockSend.LastArgs()["mapName"])
	} else {
		nameStr := string(nameVal)
		if nameStr != mapName {
			t.Errorf("Expected mapName=%s, got %s", mapName, nameStr)
		}
	}

	if xVal, ok := mockSend.LastArgs()["x"].(uint16); !ok || xVal != x {
		t.Errorf("Expected x=%d, got %v", x, mockSend.LastArgs()["x"])
	}

	if yVal, ok := mockSend.LastArgs()["y"].(uint16); !ok || yVal != y {
		t.Errorf("Expected y=%d, got %v", y, mockSend.LastArgs()["y"])
	}
}

// TestSendGMResetStateSkill tests the SendGMResetStateSkill function
func TestSendGMResetStateSkill(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test cases for different types
	testCases := []struct {
		resetType uint8
		name      string
	}{
		{0, "status"},
		{1, "skills"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := gmManager.SendGMResetStateSkill(tc.resetType)
			if err != nil {
				t.Fatalf("SendGMResetStateSkill(%d) returned error: %v", tc.resetType, err)
			}

			if mockSend.sent == nil {
				t.Fatal("No packet was sent")
			}

			// Check that the correct packet ID was used
			if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "gm_reset_state_skill" {
				t.Errorf("Expected packet ID 'gm_reset_state_skill', got '%s'", packetID)
			}

			// Check that the correct arguments were used
			if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || typeVal != tc.resetType {
				t.Errorf("Expected type=%d, got %v", tc.resetType, mockSend.LastArgs()["type"])
			}
		})
	}

	// Test with invalid type
	err := gmManager.SendGMResetStateSkill(2)
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

// TestSendGMChangeMapType tests the SendGMChangeMapType function
func TestSendGMChangeMapType(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test cases for different types
	testCases := []struct {
		cellType uint8
		name     string
	}{
		{0, "not walkable"},
		{1, "walkable"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x := uint16(150)
			y := uint16(150)
			err := gmManager.SendGMChangeMapType(x, y, tc.cellType)
			if err != nil {
				t.Fatalf("SendGMChangeMapType(%d, %d, %d) returned error: %v", x, y, tc.cellType, err)
			}

			if mockSend.sent == nil {
				t.Fatal("No packet was sent")
			}

			// Check that the correct packet ID was used
			if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "gm_change_cell_type" {
				t.Errorf("Expected packet ID 'gm_change_cell_type', got '%s'", packetID)
			}

			// Check that the correct arguments were used
			if xVal, ok := mockSend.LastArgs()["x"].(uint16); !ok || xVal != x {
				t.Errorf("Expected x=%d, got %v", x, mockSend.LastArgs()["x"])
			}

			if yVal, ok := mockSend.LastArgs()["y"].(uint16); !ok || yVal != y {
				t.Errorf("Expected y=%d, got %v", y, mockSend.LastArgs()["y"])
			}

			if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || typeVal != tc.cellType {
				t.Errorf("Expected type=%d, got %v", tc.cellType, mockSend.LastArgs()["type"])
			}
		})
	}

	// Test with invalid type
	err := gmManager.SendGMChangeMapType(150, 150, 2)
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

// TestSendGMChangeEffectState tests the SendGMChangeEffectState function
func TestSendGMChangeEffectState(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test sending a GM change effect state request
	effectState := uint32(12345)
	err := gmManager.SendGMChangeEffectState(effectState)
	if err != nil {
		t.Fatalf("SendGMChangeEffectState() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "gm_change_effect_state" {
		t.Errorf("Expected packet ID 'gm_change_effect_state', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if stateVal, ok := mockSend.LastArgs()["effect_state"].(uint32); !ok || stateVal != effectState {
		t.Errorf("Expected effect_state=%d, got %v", effectState, mockSend.LastArgs()["effect_state"])
	}
}

// TestSendGMRemove tests the SendGMRemove function
func TestSendGMRemove(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test sending a GM remove request
	playerName := "TestPlayer"
	err := gmManager.SendGMRemove(playerName)
	if err != nil {
		t.Fatalf("SendGMRemove() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "gm_remove" {
		t.Errorf("Expected packet ID 'gm_remove', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if nameVal, ok := mockSend.LastArgs()["playerName"].([]byte); !ok {
		t.Errorf("Expected playerName to be []byte, got %T", mockSend.LastArgs()["playerName"])
	} else {
		nameStr := string(nameVal)
		if nameStr != playerName {
			t.Errorf("Expected playerName=%s, got %s", playerName, nameStr)
		}
	}
}

// TestSendGMShift tests the SendGMShift function
func TestSendGMShift(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test sending a GM shift request
	playerName := "TestPlayer"
	err := gmManager.SendGMShift(playerName)
	if err != nil {
		t.Fatalf("SendGMShift() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "gm_shift" {
		t.Errorf("Expected packet ID 'gm_shift', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if nameVal, ok := mockSend.LastArgs()["playerName"].([]byte); !ok {
		t.Errorf("Expected playerName to be []byte, got %T", mockSend.LastArgs()["playerName"])
	} else {
		nameStr := string(nameVal)
		if nameStr != playerName {
			t.Errorf("Expected playerName=%s, got %s", playerName, nameStr)
		}
	}
}

// TestSendGMRecall tests the SendGMRecall function
func TestSendGMRecall(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test sending a GM recall request
	playerName := "TestPlayer"
	err := gmManager.SendGMRecall(playerName)
	if err != nil {
		t.Fatalf("SendGMRecall() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "gm_recall" {
		t.Errorf("Expected packet ID 'gm_recall', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if nameVal, ok := mockSend.LastArgs()["playerName"].([]byte); !ok {
		t.Errorf("Expected playerName to be []byte, got %T", mockSend.LastArgs()["playerName"])
	} else {
		nameStr := string(nameVal)
		if nameStr != playerName {
			t.Errorf("Expected playerName=%s, got %s", playerName, nameStr)
		}
	}
}

// TestSendGMGiveMannerByName tests the SendGMGiveMannerByName function
func TestSendGMGiveMannerByName(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test sending a GM give manner by name request
	playerName := "TestPlayer"
	err := gmManager.SendGMGiveMannerByName(playerName)
	if err != nil {
		t.Fatalf("SendGMGiveMannerByName() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "manner_by_name" {
		t.Errorf("Expected packet ID 'manner_by_name', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if nameVal, ok := mockSend.LastArgs()["playerName"].([]byte); !ok {
		t.Errorf("Expected playerName to be []byte, got %T", mockSend.LastArgs()["playerName"])
	} else {
		nameStr := string(nameVal)
		if nameStr != playerName {
			t.Errorf("Expected playerName=%s, got %s", playerName, nameStr)
		}
	}
}

// TestSendGMRequestStatus tests the SendGMRequestStatus function
func TestSendGMRequestStatus(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test sending a GM request status request
	playerName := "TestPlayer"
	err := gmManager.SendGMRequestStatus(playerName)
	if err != nil {
		t.Fatalf("SendGMRequestStatus() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "gm_request_status" {
		t.Errorf("Expected packet ID 'gm_request_status', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if nameVal, ok := mockSend.LastArgs()["playerName"].([]byte); !ok {
		t.Errorf("Expected playerName to be []byte, got %T", mockSend.LastArgs()["playerName"])
	} else {
		nameStr := string(nameVal)
		if nameStr != playerName {
			t.Errorf("Expected playerName=%s, got %s", playerName, nameStr)
		}
	}
}

// TestSendGMReqAccName tests the SendGMReqAccName function
func TestSendGMReqAccName(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test sending a GM request account name request
	targetID := uint32(12345)
	err := gmManager.SendGMReqAccName(targetID)
	if err != nil {
		t.Fatalf("SendGMReqAccName() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "gm_request_account_name" {
		t.Errorf("Expected packet ID 'gm_request_account_name', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["targetID"].(uint32); !ok || idVal != targetID {
		t.Errorf("Expected targetID=%d, got %v", targetID, mockSend.LastArgs()["targetID"])
	}
}

// TestSendBanCheck tests the SendBanCheck function
func TestSendBanCheck(t *testing.T) {
	mockSend := NewMockSend()
	gmManager := NewGMManager(mockSend)

	// Test sending a ban check request
	accountID := uint32(12345)
	err := gmManager.SendBanCheck(accountID)
	if err != nil {
		t.Fatalf("SendBanCheck() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "ban_check" {
		t.Errorf("Expected packet ID 'ban_check', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["accountID"].(uint32); !ok || idVal != accountID {
		t.Errorf("Expected accountID=%d, got %v", accountID, mockSend.LastArgs()["accountID"])
	}
}
