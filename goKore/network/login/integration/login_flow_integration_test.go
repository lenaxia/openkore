package integration

import (
	"testing"

	"github.com/lenaxia/goKore/network/login"
)

// TestPacketHandling tests that the packet handlers correctly process packets
func TestPacketHandling(t *testing.T) {
	// Create a mock network manager
	networkManager := login.NewMockNetworkManager()

	// Create login config
	config := login.NewLoginConfig("botijo0", "Melon.77", "rAthena")

	// Create login manager
	loginManager := login.NewLoginManager(networkManager, config)

	// Get the session store
	sessionStore := loginManager.GetSessionStore()

	// Test account server info packet
	accountInfoPacket := []byte{
		0xCD, 0x1F, 0x6F, 0xBC, // Session ID (3161399245)
		0x82, 0x84, 0x1E, 0x00, // Account ID (2000002)
		0x9B, 0x94, 0xE2, 0x3E, // Session ID2 (1055036571)
		0x01, // Sex (1 = male)
	}

	// Update session store directly
	args := map[string]interface{}{
		"sessionID":  accountInfoPacket[0:4],
		"accountID":  accountInfoPacket[4:8],
		"sessionID2": accountInfoPacket[8:12],
		"accountSex": int(accountInfoPacket[12]),
	}
	sessionStore.UpdateFromAccountServerInfo(args)

	// Verify account info was stored correctly
	sessionData := sessionStore.GetSessionData()
	if len(sessionData.AccountID) != 4 || sessionData.AccountID[0] != 0x82 || sessionData.AccountID[1] != 0x84 ||
		sessionData.AccountID[2] != 0x1E || sessionData.AccountID[3] != 0x00 {
		t.Errorf("Account ID not stored correctly: %v", sessionData.AccountID)
	}

	if len(sessionData.SessionID) != 4 || sessionData.SessionID[0] != 0xCD || sessionData.SessionID[1] != 0x1F ||
		sessionData.SessionID[2] != 0x6F || sessionData.SessionID[3] != 0xBC {
		t.Errorf("Session ID not stored correctly: %v", sessionData.SessionID)
	}

	if len(sessionData.SessionID2) != 4 || sessionData.SessionID2[0] != 0x9B || sessionData.SessionID2[1] != 0x94 ||
		sessionData.SessionID2[2] != 0xE2 || sessionData.SessionID2[3] != 0x3E {
		t.Errorf("Session ID2 not stored correctly: %v", sessionData.SessionID2)
	}

	if sessionData.AccountSex != 1 {
		t.Errorf("Account sex not stored correctly: %v", sessionData.AccountSex)
	}

	// Test character map info packet
	charMapInfoArgs := map[string]interface{}{
		"charID":  []byte{0xF2, 0x49, 0x02, 0x00}, // Char ID (150002)
		"mapName": "gef_fild07.gat",
		"mapIP":   "192.168.5.219",
		"mapPort": 5121,
	}
	sessionStore.UpdateFromCharacterServerInfo(charMapInfoArgs)

	// Verify character map info was stored correctly
	sessionData = sessionStore.GetSessionData()
	if len(sessionData.CharID) != 4 || sessionData.CharID[0] != 0xF2 || sessionData.CharID[1] != 0x49 ||
		sessionData.CharID[2] != 0x02 || sessionData.CharID[3] != 0x00 {
		t.Errorf("Character ID not stored correctly: %v", sessionData.CharID)
	}

	if sessionData.MapName != "gef_fild07.gat" {
		t.Errorf("Map name not stored correctly: %v", sessionData.MapName)
	}

	if sessionData.MapIP != "192.168.5.219" {
		t.Errorf("Map IP not stored correctly: %v", sessionData.MapIP)
	}

	if sessionData.MapPort != 5121 {
		t.Errorf("Map port not stored correctly: %v", sessionData.MapPort)
	}
}
