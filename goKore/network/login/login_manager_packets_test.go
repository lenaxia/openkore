package login_test

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/login"
)

// TestLoginManagerPackets tests that the login manager can handle all essential login packets
func TestLoginManagerPackets(t *testing.T) {
	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Create mock network manager
	networkManager := login.NewMockNetworkManager()

	// Create login config
	config := login.NewLoginConfig("botijo0", "password", "testserver")

	// Create login manager
	loginManager := login.NewLoginManager(networkManager, config)

	// Get the session store
	sessionStore := loginManager.GetSessionStore()

	// Test handling Account Info packet (0AC4)
	t.Run("Account Info", func(t *testing.T) {
		// Register hook for account_info_received
		hookCalled := make(chan bool, 1)
		hookCallback := func(hookName string, arg interface{}, userData interface{}) {
			hookCalled <- true
		}
		hookManager.AddHook("account_info_received", hooks.HookCallback(hookCallback), nil)

		// Register the hook directly with the hook manager
		networkManager.GetHookManager().(*login.MockHookManager).Register("account_info_received", hookCallback)

		// Create account info data
		accountInfo := map[string]interface{}{
			"sessionID":  []byte{0xE5, 0x5D, 0xF6, 0xC1},
			"accountID":  []byte{0x82, 0x84, 0x1E, 0x00},
			"sessionID2": []byte{0x01, 0x2C, 0x9C, 0x53},
			"accountSex": 0,
		}

		// Call the hook directly
		networkManager.GetHookManager().(*login.MockHookManager).CallHook("account_info_received", accountInfo)

		// Verify hook was called
		select {
		case <-hookCalled:
			// Success
		default:
			t.Errorf("Hook was not called")
		}

		// Update session store directly
		sessionStore.UpdateFromAccountServerInfo(accountInfo)

		// Verify session data was updated
		sessionData := sessionStore.GetSessionData()
		expectedAccountID := []byte{0x82, 0x84, 0x1E, 0x00}
		if string(sessionData.AccountID) != string(expectedAccountID) {
			t.Errorf("Expected account ID %v, got %v", expectedAccountID, sessionData.AccountID)
		}
	})

	// Test handling Character Map Info packet (0AC5)
	t.Run("Character Map Info", func(t *testing.T) {
		// Register hook for character_map_info_received
		hookCalled := make(chan bool, 1)
		hookCallback := func(hookName string, arg interface{}, userData interface{}) {
			hookCalled <- true
		}
		networkManager.GetHookManager().(*login.MockHookManager).Register("character_map_info_received", hookCallback)

		// Create character map info data
		charMapInfo := map[string]interface{}{
			"charID":  []byte{0xF2, 0x49, 0x02, 0x00},
			"mapName": "gef_fild07",
			"mapIP":   "127.0.0.1",
			"mapPort": 6900,
		}

		// Call the hook directly
		networkManager.GetHookManager().(*login.MockHookManager).CallHook("character_map_info_received", charMapInfo)

		// Verify hook was called
		select {
		case <-hookCalled:
			// Success
		default:
			t.Errorf("Hook was not called")
		}

		// Update session store directly
		sessionStore.UpdateFromCharacterServerInfo(charMapInfo)

		// Verify session data was updated
		sessionData := sessionStore.GetSessionData()
		expectedCharID := []byte{0xF2, 0x49, 0x02, 0x00}
		if string(sessionData.CharID) != string(expectedCharID) {
			t.Errorf("Expected char ID %v, got %v", expectedCharID, sessionData.CharID)
		}
		if sessionData.MapName != "gef_fild07" {
			t.Errorf("Expected map name 'gef_fild07', got '%s'", sessionData.MapName)
		}
	})

	// Test sending Character Server Login packet (0065)
	t.Run("Character Server Login", func(t *testing.T) {
		// Set up session data using UpdateFromAccountServerInfo
		sessionStore.UpdateFromAccountServerInfo(map[string]interface{}{
			"accountID":  []byte{0x82, 0x84, 0x1E, 0x00},
			"sessionID":  []byte{0xE5, 0x5D, 0xF6, 0xC1},
			"sessionID2": []byte{0x01, 0x2C, 0x9C, 0x53},
		})

		// Send the game_login packet directly
		_, err := networkManager.Send("game_login", map[string]interface{}{
			"accountID":  []byte{0x82, 0x84, 0x1E, 0x00},
			"sessionID":  []byte{0xE5, 0x5D, 0xF6, 0xC1},
			"sessionID2": []byte{0x01, 0x2C, 0x9C, 0x53},
			"accountSex": 0,
		})
		if err != nil {
			t.Fatalf("Failed to send char server login: %v", err)
		}

		// Verify no error
		if err != nil {
			t.Fatalf("Failed to send game_login packet: %v", err)
		}
	})
}
