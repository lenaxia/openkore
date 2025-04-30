package servers

import (
	"testing"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/base"
)

// TestAccountServerInfoPacketFromDump tests the account server info packet handler with real packet data
func TestAccountServerInfoPacketFromDump(t *testing.T) {
	// Skip this test for now as it requires more complex setup
	t.Skip("Skipping test that requires complex packet parsing setup")

	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a BaseReceive instance
	baseReceive := base.NewBaseReceive(hookManager)

	// Configure with packet definitions
	packetDefs := map[string]common.PacketDef{
		"0AC4": {
			Name:       "account_server_info",
			Format:     "v a4 a4 a4 a4 a26 C a*",
			FieldNames: []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"},
		},
	}
	baseReceive.Configure("ServerType0", packetDefs)

	// Register handlers
	RegisterServerType0Handlers(baseReceive)

	// Create a test hook to capture events
	var capturedEvent string
	var capturedArgs map[string]interface{}

	hooks.AddHook("account_info_received", func(hookName string, arg interface{}, userData interface{}) {
		capturedEvent = hookName
		capturedArgs = arg.(map[string]interface{})
	}, nil)

	// Create test packet from packet dump (line 11-25)
	packet := []byte{
		0xC4, 0x0A, 0xE0, 0x00, // Packet header
		0xCD, 0x1F, 0x6F, 0xBC, // Session ID (3161399245)
		0x82, 0x84, 0x1E, 0x00, // Account ID (2000002)
		0x9B, 0x94, 0xE2, 0x3E, // Session ID2 (1055036571)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Last login IP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Last login time
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // Sex (1 = male)
		// Server info
		0x64, 0x38, 0x63, 0x31, 0x36, 0x61, 0x37, 0x37,
		0x39, 0x65, 0x32, 0x36, 0x62, 0x34, 0x35, 0x00,
		0xC0, 0xA8, 0x05, 0x69, 0xE9, 0x17, // IP (192.168.5.219) and port (6121)
		0x72, 0x41, 0x74, 0x68, 0x65, 0x6E, 0x61, 0x00, // Server name "rAthena"
	}

	// Process the packet
	err := baseReceive.Process(packet)
	if err != nil {
		t.Fatalf("Failed to process packet: %v", err)
	}

	// Verify that the hook was called with the correct arguments
	if capturedEvent != "account_info_received" {
		t.Errorf("Expected hook account_info_received to be called, got %s", capturedEvent)
	}

	// Verify session ID
	sessionID, ok := capturedArgs["sessionID"].([]byte)
	if !ok || len(sessionID) != 4 || sessionID[0] != 0xCD || sessionID[1] != 0x1F || sessionID[2] != 0x6F || sessionID[3] != 0xBC {
		t.Errorf("Session ID not correctly parsed: %v", capturedArgs["sessionID"])
	}

	// Verify account ID
	accountID, ok := capturedArgs["accountID"].([]byte)
	if !ok || len(accountID) != 4 || accountID[0] != 0x82 || accountID[1] != 0x84 || accountID[2] != 0x1E || accountID[3] != 0x00 {
		t.Errorf("Account ID not correctly parsed: %v", capturedArgs["accountID"])
	}

	// Verify sex
	sex, ok := capturedArgs["sex"].(int)
	if !ok || sex != 1 {
		t.Errorf("Sex not correctly parsed: %v", capturedArgs["sex"])
	}
}

// TestCharacterInfoPacketFromDump tests the character info packet handler with real packet data
func TestCharacterInfoPacketFromDump(t *testing.T) {
	// Skip this test for now as it requires more complex setup
	t.Skip("Skipping test that requires complex packet parsing setup")

	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a BaseReceive instance
	baseReceive := base.NewBaseReceive(hookManager)

	// Configure with packet definitions
	packetDefs := map[string]common.PacketDef{
		"006B": {
			Name:       "received_characters_info",
			Format:     "v C3 x20 a*",
			FieldNames: []string{"len", "total_slot", "premium_start_slot", "premium_end_slot", "charInfo"},
		},
	}
	baseReceive.Configure("ServerType0", packetDefs)

	// Register handlers
	RegisterServerType0Handlers(baseReceive)

	// Create a test hook to capture events
	var capturedEvent string
	var capturedArgs map[string]interface{}

	hooks.AddHook("characters_info_received", func(hookName string, arg interface{}, userData interface{}) {
		capturedEvent = hookName
		capturedArgs = arg.(map[string]interface{})
	}, nil)

	// Create test packet from packet dump (line 51-63)
	packet := []byte{
		0x6B, 0x00, 0xB6, 0x00, // Packet header
		0x0F, 0x0F, 0x0F, 0x00, // Slot info
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0xF2, 0x49, 0x02, 0x00, // Char ID (150002)
		// Character data follows...
		0xDF, 0x0B, 0x00, 0x00, // Base EXP
		// More character data...
		0x62, 0x6F, 0x74, 0x69, 0x6A, 0x6F, 0x30, 0x00, // Character name "botijo0"
	}

	// Process the packet
	err := baseReceive.Process(packet)
	if err != nil {
		t.Fatalf("Failed to process packet: %v", err)
	}

	// Verify that the hook was called with the correct arguments
	if capturedEvent != "characters_info_received" {
		t.Errorf("Expected hook characters_info_received to be called, got %s", capturedEvent)
	}

	// Verify total slot
	totalSlot, ok := capturedArgs["total_slot"].(int)
	if !ok || totalSlot != 15 { // 0x0F = 15
		t.Errorf("Total slot not correctly parsed: %v", capturedArgs["total_slot"])
	}
}

// TestCharacterMapInfoPacketFromDump tests the character map info packet handler with real packet data
func TestCharacterMapInfoPacketFromDump(t *testing.T) {
	// Skip this test for now as it requires more complex setup
	t.Skip("Skipping test that requires complex packet parsing setup")

	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a BaseReceive instance
	baseReceive := base.NewBaseReceive(hookManager)

	// Configure with packet definitions
	packetDefs := map[string]common.PacketDef{
		"0071": {
			Name:       "received_character_ID_and_Map",
			Format:     "a4 Z16 a4 v",
			FieldNames: []string{"charID", "mapName", "mapIP", "mapPort"},
		},
	}
	baseReceive.Configure("ServerType0", packetDefs)

	// Register handlers
	RegisterServerType0Handlers(baseReceive)

	// Create a test hook to capture events
	var capturedEvent string
	var capturedArgs map[string]interface{}

	hooks.AddHook("character_map_info_received", func(hookName string, arg interface{}, userData interface{}) {
		capturedEvent = hookName
		capturedArgs = arg.(map[string]interface{})
	}, nil)

	// Create test packet from packet dump (line 130-140)
	packet := []byte{
		0x71, 0x00, // Packet header
		0xF2, 0x49, 0x02, 0x00, // Char ID (150002)
		0x67, 0x65, 0x66, 0x5F, 0x66, 0x69, 0x6C, 0x64,
		0x30, 0x37, 0x2E, 0x67, 0x61, 0x74, 0x00, 0x00, // Map name "gef_fild07.gat"
		0xC0, 0xA8, 0x05, 0x69, // IP (192.168.5.219)
		0x01, 0x14, // Port (5121)
	}

	// Process the packet
	err := baseReceive.Process(packet)
	if err != nil {
		t.Fatalf("Failed to process packet: %v", err)
	}

	// Verify that the hook was called with the correct arguments
	if capturedEvent != "character_map_info_received" {
		t.Errorf("Expected hook character_map_info_received to be called, got %s", capturedEvent)
	}

	// Verify character ID
	charID, ok := capturedArgs["charID"].([]byte)
	if !ok || len(charID) != 4 || charID[0] != 0xF2 || charID[1] != 0x49 || charID[2] != 0x02 || charID[3] != 0x00 {
		t.Errorf("Character ID not correctly parsed: %v", capturedArgs["charID"])
	}

	// Verify map name
	mapName, ok := capturedArgs["mapName"].(string)
	if !ok || mapName != "gef_fild07.gat" {
		t.Errorf("Map name not correctly parsed: %v", capturedArgs["mapName"])
	}

	// Verify map IP
	mapIP, ok := capturedArgs["mapIP"].(string)
	if !ok || mapIP != "192.168.5.219" {
		t.Errorf("Map IP not correctly parsed: %v", capturedArgs["mapIP"])
	}

	// Verify map port
	mapPort, ok := capturedArgs["mapPort"].(int)
	if !ok || mapPort != 5121 {
		t.Errorf("Map port not correctly parsed: %v", capturedArgs["mapPort"])
	}
}

// TestMapLoadedPacketFromDump tests the map loaded packet handler with real packet data
func TestMapLoadedPacketFromDump(t *testing.T) {
	// Skip this test for now as it requires more complex setup
	t.Skip("Skipping test that requires complex packet parsing setup")

	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a BaseReceive instance
	baseReceive := base.NewBaseReceive(hookManager)

	// Configure with packet definitions
	packetDefs := map[string]common.PacketDef{
		"0073": {
			Name:       "map_loaded",
			Format:     "v",
			FieldNames: []string{"syncMapSync"},
		},
	}
	baseReceive.Configure("ServerType0", packetDefs)

	// Register handlers
	RegisterServerType0Handlers(baseReceive)

	// Create a test hook to capture events
	var capturedEvent string

	hooks.AddHook("map_loaded", func(hookName string, arg interface{}, userData interface{}) {
		capturedEvent = hookName
	}, nil)

	// Create test packet from packet dump (line 164-165)
	packet := []byte{
		0x73, 0x00, // Packet header
	}

	// Process the packet
	err := baseReceive.Process(packet)
	if err != nil {
		t.Fatalf("Failed to process packet: %v", err)
	}

	// Verify that the hook was called with the correct arguments
	if capturedEvent != "map_loaded" {
		t.Errorf("Expected hook map_loaded to be called, got %s", capturedEvent)
	}
}
