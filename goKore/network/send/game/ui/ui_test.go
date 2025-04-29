package ui

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
			"misc_config_set":              "02D8",
			"notify_progress_bar_complete": "02F1",
			"view_player_equip_request":    "02D6",
			"refineui_select":              "0AA1",
			"refineui_refine":              "0AA3",
			"refineui_close":               "0AA4",
			"item_list_window_selected":    "09A4",
			"memo_request":                 "011D",
			"stylist_change":               "0A18",
			"open_ui_request":              "0A68",
			"attendance_reward_request":    "0AEF",
			"roulette_window_open":         "0A19",
			"roulette_info_request":        "0A1B",
			"roulette_close":               "0A1D",
			"roulette_start":               "0A1F",
			"roulette_claim_prize":         "0A21",
			"send_quest_state":             "02B6",
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

// TestNewUIManager tests the NewUIManager function
func TestNewUIManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	if uiManager == nil {
		t.Fatal("NewUIManager() returned nil")
	}

	if uiManager.baseSend == nil {
		t.Error("uiManager.baseSend was not set correctly")
	}
}

// TestSendMiscConfigSet tests the SendMiscConfigSet function
func TestSendMiscConfigSet(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a misc config set request
	configType := uint32(0) // show equip windows to other players
	flag := uint32(1)       // enabled
	err := uiManager.SendMiscConfigSet(configType, flag)
	if err != nil {
		t.Fatalf("SendMiscConfigSet() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "misc_config_set" {
		t.Errorf("Expected packet ID 'misc_config_set', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if typeVal, ok := mockSend.LastArgs()["type"].(uint32); !ok || typeVal != configType {
		t.Errorf("Expected type=%d, got %v", configType, mockSend.LastArgs()["type"])
	}

	if flagVal, ok := mockSend.LastArgs()["flag"].(uint32); !ok || flagVal != flag {
		t.Errorf("Expected flag=%d, got %v", flag, mockSend.LastArgs()["flag"])
	}
}

// TestSendProgress tests the SendProgress function
func TestSendProgress(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a progress bar complete notification
	err := uiManager.SendProgress()
	if err != nil {
		t.Fatalf("SendProgress() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "notify_progress_bar_complete" {
		t.Errorf("Expected packet ID 'notify_progress_bar_complete', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendShowEquipPlayer tests the SendShowEquipPlayer function
func TestSendShowEquipPlayer(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a show equip player request
	playerID := uint32(12345)
	err := uiManager.SendShowEquipPlayer(playerID)
	if err != nil {
		t.Fatalf("SendShowEquipPlayer() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "view_player_equip_request" {
		t.Errorf("Expected packet ID 'view_player_equip_request', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["ID"].(uint32); !ok || idVal != playerID {
		t.Errorf("Expected ID=%d, got %v", playerID, mockSend.LastArgs()["ID"])
	}
}

// TestSendRefineUISelect tests the SendRefineUISelect function
func TestSendRefineUISelect(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a refine UI select request
	itemIndex := uint16(123)
	err := uiManager.SendRefineUISelect(itemIndex)
	if err != nil {
		t.Fatalf("SendRefineUISelect() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "refineui_select" {
		t.Errorf("Expected packet ID 'refineui_select', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if indexVal, ok := mockSend.LastArgs()["index"].(uint16); !ok || indexVal != itemIndex {
		t.Errorf("Expected index=%d, got %v", itemIndex, mockSend.LastArgs()["index"])
	}
}

// TestSendRefineUIRefine tests the SendRefineUIRefine function
func TestSendRefineUIRefine(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a refine UI refine request
	itemIndex := uint16(123)
	materialNameID := uint16(456)
	useCatalyst := uint8(1)
	err := uiManager.SendRefineUIRefine(itemIndex, materialNameID, useCatalyst)
	if err != nil {
		t.Fatalf("SendRefineUIRefine() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "refineui_refine" {
		t.Errorf("Expected packet ID 'refineui_refine', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if indexVal, ok := mockSend.LastArgs()["index"].(uint16); !ok || indexVal != itemIndex {
		t.Errorf("Expected index=%d, got %v", itemIndex, mockSend.LastArgs()["index"])
	}

	if catalystVal, ok := mockSend.LastArgs()["catalyst"].(uint16); !ok || catalystVal != materialNameID {
		t.Errorf("Expected catalyst=%d, got %v", materialNameID, mockSend.LastArgs()["catalyst"])
	}

	if blessVal, ok := mockSend.LastArgs()["bless"].(uint8); !ok || blessVal != useCatalyst {
		t.Errorf("Expected bless=%d, got %v", useCatalyst, mockSend.LastArgs()["bless"])
	}
}

// TestSendRefineUIClose tests the SendRefineUIClose function
func TestSendRefineUIClose(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a refine UI close request
	err := uiManager.SendRefineUIClose()
	if err != nil {
		t.Fatalf("SendRefineUIClose() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "refineui_close" {
		t.Errorf("Expected packet ID 'refineui_close', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendItemListWindowSelected tests the SendItemListWindowSelected function
func TestSendItemListWindowSelected(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending an item list window selected request
	itemType := uint8(0) // Change Material
	act := uint8(1)      // Process
	items := []map[string]interface{}{
		{
			"itemIndex": uint16(123),
			"amount":    uint16(1),
			"itemName":  "Test Item 1",
		},
		{
			"itemIndex": uint16(456),
			"amount":    uint16(2),
			"itemName":  "Test Item 2",
		},
	}
	err := uiManager.SendItemListWindowSelected(itemType, act, items)
	if err != nil {
		t.Fatalf("SendItemListWindowSelected() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "item_list_window_selected" {
		t.Errorf("Expected packet ID 'item_list_window_selected', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || typeVal != itemType {
		t.Errorf("Expected type=%d, got %v", itemType, mockSend.LastArgs()["type"])
	}

	if actVal, ok := mockSend.LastArgs()["act"].(uint8); !ok || actVal != act {
		t.Errorf("Expected act=%d, got %v", act, mockSend.LastArgs()["act"])
	}

	if itemsVal, ok := mockSend.LastArgs()["items"].([]map[string]interface{}); !ok {
		t.Errorf("Expected items to be []map[string]interface{}, got %T", mockSend.LastArgs()["items"])
	} else {
		if len(itemsVal) != 2 {
			t.Errorf("Expected 2 items, got %d", len(itemsVal))
		}
	}

	// Check that the len field was calculated correctly
	if lenVal, ok := mockSend.LastArgs()["len"].(int); !ok || lenVal != 20 { // (2 * 4) + 12 = 20
		t.Errorf("Expected len=20, got %v", mockSend.LastArgs()["len"])
	}
}

// TestReconstructItemListWindowSelected tests the ReconstructItemListWindowSelected function
func TestReconstructItemListWindowSelected(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test reconstructing an item list window selected packet
	items := []map[string]interface{}{
		{
			"itemIndex": uint16(123),
			"amount":    uint16(1),
			"itemName":  "Test Item 1",
		},
		{
			"itemIndex": uint16(456),
			"amount":    uint16(2),
			"itemName":  "Test Item 2",
		},
	}
	args := map[string]interface{}{
		"items": items,
	}
	err := uiManager.ReconstructItemListWindowSelected(args)
	if err != nil {
		t.Fatalf("ReconstructItemListWindowSelected() returned error: %v", err)
	}

	// Check that the itemInfo field was added
	if _, ok := args["itemInfo"]; !ok {
		t.Error("itemInfo field was not added to args")
	}
}

// TestSendMemo tests the SendMemo function
func TestSendMemo(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a memo request
	err := uiManager.SendMemo()
	if err != nil {
		t.Fatalf("SendMemo() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "memo_request" {
		t.Errorf("Expected packet ID 'memo_request', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendStylistChange tests the SendStylistChange function
func TestSendStylistChange(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a stylist change request
	hairColor := uint16(1)
	hairStyle := uint16(2)
	clothColor := uint16(3)
	headTop := uint16(4)
	headMid := uint16(5)
	headBottom := uint16(6)
	err := uiManager.SendStylistChange(hairColor, hairStyle, clothColor, headTop, headMid, headBottom)
	if err != nil {
		t.Fatalf("SendStylistChange() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "stylist_change" {
		t.Errorf("Expected packet ID 'stylist_change', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if colorVal, ok := mockSend.LastArgs()["hair_color"].(uint16); !ok || colorVal != hairColor {
		t.Errorf("Expected hair_color=%d, got %v", hairColor, mockSend.LastArgs()["hair_color"])
	}

	if styleVal, ok := mockSend.LastArgs()["hair_style"].(uint16); !ok || styleVal != hairStyle {
		t.Errorf("Expected hair_style=%d, got %v", hairStyle, mockSend.LastArgs()["hair_style"])
	}

	if colorVal, ok := mockSend.LastArgs()["cloth_color"].(uint16); !ok || colorVal != clothColor {
		t.Errorf("Expected cloth_color=%d, got %v", clothColor, mockSend.LastArgs()["cloth_color"])
	}

	if topVal, ok := mockSend.LastArgs()["head_top"].(uint16); !ok || topVal != headTop {
		t.Errorf("Expected head_top=%d, got %v", headTop, mockSend.LastArgs()["head_top"])
	}

	if midVal, ok := mockSend.LastArgs()["head_mid"].(uint16); !ok || midVal != headMid {
		t.Errorf("Expected head_mid=%d, got %v", headMid, mockSend.LastArgs()["head_mid"])
	}

	if bottomVal, ok := mockSend.LastArgs()["head_bottom"].(uint16); !ok || bottomVal != headBottom {
		t.Errorf("Expected head_bottom=%d, got %v", headBottom, mockSend.LastArgs()["head_bottom"])
	}
}

// TestSendOpenUIRequest tests the SendOpenUIRequest function
func TestSendOpenUIRequest(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending an open UI request
	uiType := uint8(1)
	err := uiManager.SendOpenUIRequest(uiType)
	if err != nil {
		t.Fatalf("SendOpenUIRequest() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "open_ui_request" {
		t.Errorf("Expected packet ID 'open_ui_request', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if typeVal, ok := mockSend.LastArgs()["UIType"].(uint8); !ok || typeVal != uiType {
		t.Errorf("Expected UIType=%d, got %v", uiType, mockSend.LastArgs()["UIType"])
	}
}

// TestSendAttendanceRewardRequest tests the SendAttendanceRewardRequest function
func TestSendAttendanceRewardRequest(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending an attendance reward request
	err := uiManager.SendAttendanceRewardRequest()
	if err != nil {
		t.Fatalf("SendAttendanceRewardRequest() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "attendance_reward_request" {
		t.Errorf("Expected packet ID 'attendance_reward_request', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendRouletteWindowOpen tests the SendRouletteWindowOpen function
func TestSendRouletteWindowOpen(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a roulette window open request
	err := uiManager.SendRouletteWindowOpen()
	if err != nil {
		t.Fatalf("SendRouletteWindowOpen() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "roulette_window_open" {
		t.Errorf("Expected packet ID 'roulette_window_open', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendRouletteInfoRequest tests the SendRouletteInfoRequest function
func TestSendRouletteInfoRequest(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a roulette info request
	err := uiManager.SendRouletteInfoRequest()
	if err != nil {
		t.Fatalf("SendRouletteInfoRequest() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "roulette_info_request" {
		t.Errorf("Expected packet ID 'roulette_info_request', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendRouletteClose tests the SendRouletteClose function
func TestSendRouletteClose(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a roulette close request
	err := uiManager.SendRouletteClose()
	if err != nil {
		t.Fatalf("SendRouletteClose() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "roulette_close" {
		t.Errorf("Expected packet ID 'roulette_close', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendRouletteStart tests the SendRouletteStart function
func TestSendRouletteStart(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a roulette start request
	err := uiManager.SendRouletteStart()
	if err != nil {
		t.Fatalf("SendRouletteStart() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "roulette_start" {
		t.Errorf("Expected packet ID 'roulette_start', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendRouletteClaimPrize tests the SendRouletteClaimPrize function
func TestSendRouletteClaimPrize(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a roulette claim prize request
	err := uiManager.SendRouletteClaimPrize()
	if err != nil {
		t.Fatalf("SendRouletteClaimPrize() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "roulette_claim_prize" {
		t.Errorf("Expected packet ID 'roulette_claim_prize', got '%s'", packetID)
	}

	// Check that no arguments were used
	if len(mockSend.LastArgs()) != 0 {
		t.Errorf("Expected no arguments, got %v", mockSend.LastArgs())
	}
}

// TestSendQuestState tests the SendQuestState function
func TestSendQuestState(t *testing.T) {
	mockSend := NewMockSend()
	uiManager := NewUIManager(mockSend)

	// Test sending a quest state request
	questID := uint32(123)
	state := uint8(0) // active
	err := uiManager.SendQuestState(questID, state)
	if err != nil {
		t.Fatalf("SendQuestState() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "send_quest_state" {
		t.Errorf("Expected packet ID 'send_quest_state', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["questID"].(uint32); !ok || idVal != questID {
		t.Errorf("Expected questID=%d, got %v", questID, mockSend.LastArgs()["questID"])
	}

	if stateVal, ok := mockSend.LastArgs()["state"].(uint8); !ok || stateVal != state {
		t.Errorf("Expected state=%d, got %v", state, mockSend.LastArgs()["state"])
	}
}
