package security

import (
	"testing"
	"time"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TestSecurityIntegration tests the integration between login, pin, and anticheat components
func TestSecurityIntegration(t *testing.T) {
	// Create a parser and hook manager
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()

	// Create the security components
	loginManager := NewLoginManager(parser, hookManager)
	pinManager := NewPINManager(parser, hookManager)
	antiCheatManager := NewAntiCheatManager(parser, hookManager)

	// Register handlers for all components
	loginManager.RegisterHandlers()
	pinManager.RegisterHandlers()
	antiCheatManager.RegisterHandlers()

	// Test login flow
	testLoginFlow(t, loginManager, pinManager, antiCheatManager, hookManager)

	// Test PIN flow
	testPINFlow(t, loginManager, pinManager, antiCheatManager, hookManager)

	// Test anti-cheat flow
	testAntiCheatFlow(t, loginManager, pinManager, antiCheatManager, hookManager)
}

// testLoginFlow tests the login flow
func testLoginFlow(t *testing.T, loginManager *LoginManager, pinManager *PINManager, antiCheatManager *AntiCheatManager, hookManager *hooks.HookManager) {
	// Set up test credentials
	loginManager.SetCredentials("testuser", "testpass")

	// Set up hook to capture login success
	var loginSuccessful bool
	hookManager.AddHook("security/login_success", func(hookName string, arg interface{}, userData interface{}) {
		loginSuccessful = true
	}, nil)

	// Simulate receiving account server info packet
	err := loginManager.handleAccountServerInfo(map[string]interface{}{
		"sessionID":     []byte{1, 2, 3, 4},
		"accountID":     []byte{5, 6, 7, 8},
		"sessionID2":    []byte{9, 10, 11, 12},
		"lastLoginIP":   []byte{127, 0, 0, 1},
		"lastLoginTime": "2023-01-01",
		"accountSex":    byte(1),
		"serverInfo":    []byte{},
	})

	// Check that there was no error
	if err != nil {
		t.Errorf("handleAccountServerInfo returned error: %v", err)
	}

	// Check that login was successful
	if !loginSuccessful {
		t.Error("Login success hook was not called")
	}

	// Check that login state is correct
	if loginManager.GetState() != LoginStateLoggedIn {
		t.Errorf("Login state = %v, want %v", loginManager.GetState(), LoginStateLoggedIn)
	}

	// Check that session IDs were set correctly
	sessionID, sessionID2, accountID := loginManager.GetSessionIDs()
	if len(sessionID) != 4 || sessionID[0] != 1 || sessionID[1] != 2 || sessionID[2] != 3 || sessionID[3] != 4 {
		t.Errorf("sessionID = %v, want [1 2 3 4]", sessionID)
	}
	if len(sessionID2) != 4 || sessionID2[0] != 9 || sessionID2[1] != 10 || sessionID2[2] != 11 || sessionID2[3] != 12 {
		t.Errorf("sessionID2 = %v, want [9 10 11 12]", sessionID2)
	}
	if accountID != 0x08070605 {
		t.Errorf("accountID = %v, want %v", accountID, 0x08070605)
	}
}

// testPINFlow tests the PIN flow
func testPINFlow(t *testing.T, loginManager *LoginManager, pinManager *PINManager, antiCheatManager *AntiCheatManager, hookManager *hooks.HookManager) {
	// Set up test PIN
	err := pinManager.SetPIN("1234")
	if err != nil {
		t.Fatalf("SetPIN returned error: %v", err)
	}

	// Set up hook to capture PIN verification
	var pinVerified bool
	hookManager.AddHook("security/pin_verified", func(hookName string, arg interface{}, userData interface{}) {
		pinVerified = true
	}, nil)

	// Simulate receiving PIN code request packet
	err = pinManager.handleLoginPinCodeRequest(map[string]interface{}{
		"seed":      uint32(12345),
		"accountID": []byte{5, 6, 7, 8},
		"flag":      uint16(1), // PIN is required
	})

	// Check that there was no error
	if err != nil {
		t.Errorf("handleLoginPinCodeRequest returned error: %v", err)
	}

	// Check that PIN state is correct
	if pinManager.GetState() != PINStateRequested {
		t.Errorf("PIN state = %v, want %v", pinManager.GetState(), PINStateRequested)
	}

	// Verify PIN
	err = pinManager.VerifyPIN("1234")
	if err != nil {
		t.Errorf("VerifyPIN returned error: %v", err)
	}

	// Check that PIN verification was successful
	if !pinVerified {
		t.Error("PIN verified hook was not called")
	}

	// Check that PIN state is correct
	if pinManager.GetState() != PINStateVerified {
		t.Errorf("PIN state = %v, want %v", pinManager.GetState(), PINStateVerified)
	}
}

// testAntiCheatFlow tests the anti-cheat flow
func testAntiCheatFlow(t *testing.T, loginManager *LoginManager, pinManager *PINManager, antiCheatManager *AntiCheatManager, hookManager *hooks.HookManager) {
	// Set up hook to capture anti-cheat verification
	var antiCheatVerified bool
	hookManager.AddHook("security/anticheat_verified", func(hookName string, arg interface{}, userData interface{}) {
		antiCheatVerified = true
	}, nil)

	// Simulate receiving GameGuard request packet
	err := antiCheatManager.handleGameGuardRequest(map[string]interface{}{})
	if err != nil {
		t.Errorf("handleGameGuardRequest returned error: %v", err)
	}

	// Check that anti-cheat state is correct
	if antiCheatManager.GetState() != AntiCheatStateWaitingForChallenge {
		t.Errorf("Anti-cheat state = %v, want %v", antiCheatManager.GetState(), AntiCheatStateWaitingForChallenge)
	}

	// Generate challenge
	challenge := antiCheatManager.GenerateChallenge()
	if len(challenge) != 20 {
		t.Errorf("len(challenge) = %d, want 20", len(challenge))
	}

	// Check that anti-cheat state is correct
	if antiCheatManager.GetState() != AntiCheatStateWaitingForResponse {
		t.Errorf("Anti-cheat state = %v, want %v", antiCheatManager.GetState(), AntiCheatStateWaitingForResponse)
	}

	// Generate response
	response := antiCheatManager.GenerateGameGuardResponse(challenge)

	// Verify response
	err = antiCheatManager.VerifyResponse(response)
	if err != nil {
		t.Errorf("VerifyResponse returned error: %v", err)
	}

	// Check that anti-cheat verification was successful
	if !antiCheatVerified {
		t.Error("Anti-cheat verified hook was not called")
	}

	// Check that anti-cheat state is correct
	if antiCheatManager.GetState() != AntiCheatStateVerified {
		t.Errorf("Anti-cheat state = %v, want %v", antiCheatManager.GetState(), AntiCheatStateVerified)
	}
}

// TestSecurityTimeout tests timeout handling in the security components
func TestSecurityTimeout(t *testing.T) {
	// Create a parser and hook manager
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()

	// Create the security components
	loginManager := NewLoginManager(parser, hookManager)
	pinManager := NewPINManager(parser, hookManager)
	antiCheatManager := NewAntiCheatManager(parser, hookManager)

	// Test login timeout
	testLoginTimeout(t, loginManager)

	// Test PIN timeout
	testPINTimeout(t, pinManager)

	// Test anti-cheat timeout
	testAntiCheatTimeout(t, antiCheatManager)
}

// testLoginTimeout tests login timeout handling
func testLoginTimeout(t *testing.T, loginManager *LoginManager) {
	// Set login state to logging in
	loginManager.SetState(LoginStateLoggingIn)

	// Set last activity to a long time ago
	loginManager.lastActivity = time.Now().Add(-60 * time.Second)

	// Check if session is expired with a 30-second timeout
	if !loginManager.IsSessionExpired(30 * time.Second) {
		t.Error("IsSessionExpired() = false, want true for expired session")
	}

	// Update activity
	loginManager.UpdateActivity()

	// Check if session is no longer expired
	if loginManager.IsSessionExpired(30 * time.Second) {
		t.Error("IsSessionExpired() = true, want false after updating activity")
	}
}

// testPINTimeout tests PIN timeout handling
func testPINTimeout(t *testing.T, pinManager *PINManager) {
	// Set PIN and state
	pinManager.pin = "1234"
	pinManager.state = PINStateRequested

	// Set last attempt to a long time ago
	pinManager.lastAttempt = time.Now().Add(-60 * time.Second)

	// Try to verify PIN with incorrect value
	for i := 0; i < pinManager.maxAttempts; i++ {
		err := pinManager.VerifyPIN("5678")
		if i < pinManager.maxAttempts-1 {
			if err != ErrInvalidPIN {
				t.Errorf("VerifyPIN() attempt %d returned error: %v, want %v", i, err, ErrInvalidPIN)
			}
		} else {
			// Last attempt should lock the PIN
			if err != ErrPINLocked {
				t.Errorf("VerifyPIN() last attempt returned error: %v, want %v", err, ErrPINLocked)
			}
		}
	}

	// Check that PIN is locked
	if pinManager.state != PINStateLocked {
		t.Errorf("PIN state = %v, want %v after max attempts", pinManager.state, PINStateLocked)
	}

	// Unlock PIN
	pinManager.UnlockPIN()

	// Check that PIN is unlocked
	if pinManager.state != PINStateSet {
		t.Errorf("PIN state = %v, want %v after unlock", pinManager.state, PINStateSet)
	}
}

// testAntiCheatTimeout tests anti-cheat timeout handling
func testAntiCheatTimeout(t *testing.T, antiCheatManager *AntiCheatManager) {
	// Enable anti-cheat
	antiCheatManager.Enable(AntiCheatGameGuard)

	// Generate challenge
	antiCheatManager.GenerateChallenge()

	// Set last challenge to a long time ago
	antiCheatManager.lastChallenge = time.Now().Add(-60 * time.Second)

	// Check if timed out
	if !antiCheatManager.IsTimedOut() {
		t.Error("IsTimedOut() = false, want true for timed out challenge")
	}

	// Try to verify response
	err := antiCheatManager.VerifyResponse([]byte{1, 2, 3, 4})
	if err != ErrAntiCheatTimeout {
		t.Errorf("VerifyResponse() returned error: %v, want %v", err, ErrAntiCheatTimeout)
	}

	// Check that anti-cheat state is rejected
	if antiCheatManager.state != AntiCheatStateRejected {
		t.Errorf("Anti-cheat state = %v, want %v after timeout", antiCheatManager.state, AntiCheatStateRejected)
	}
}

// TestSecurityErrorHandling tests error handling in the security components
func TestSecurityErrorHandling(t *testing.T) {
	// Create a parser and hook manager
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()

	// Create the security components
	loginManager := NewLoginManager(parser, hookManager)
	pinManager := NewPINManager(parser, hookManager)
	antiCheatManager := NewAntiCheatManager(parser, hookManager)

	// Test login error handling
	testLoginErrorHandling(t, loginManager, hookManager)

	// Test PIN error handling
	testPINErrorHandling(t, pinManager)

	// Test anti-cheat error handling
	testAntiCheatErrorHandling(t, antiCheatManager)
}

// testLoginErrorHandling tests login error handling
func testLoginErrorHandling(t *testing.T, loginManager *LoginManager, hookManager *hooks.HookManager) {
	// Set up hook to capture login error
	var loginError bool
	hookManager.AddHook("security/login_error", func(hookName string, arg interface{}, userData interface{}) {
		loginError = true
	}, nil)

	// Simulate receiving login error packet
	err := loginManager.handleLoginError(map[string]interface{}{
		"type": byte(1),
		"date": "2023-01-01",
	})

	// Check that there was no error
	if err != nil {
		t.Errorf("handleLoginError returned error: %v", err)
	}

	// Check that login error hook was called
	if !loginError {
		t.Error("Login error hook was not called")
	}

	// Check that login state is disconnected
	if loginManager.GetState() != LoginStateDisconnected {
		t.Errorf("Login state = %v, want %v after error", loginManager.GetState(), LoginStateDisconnected)
	}

	// Check that login error was set
	loginErrorObj := loginManager.GetLoginError()
	if loginErrorObj == nil {
		t.Fatal("GetLoginError() returned nil")
	}

	if loginErrorObj.Code != 1 {
		t.Errorf("Login error code = %d, want 1", loginErrorObj.Code)
	}

	if loginErrorObj.Date != "2023-01-01" {
		t.Errorf("Login error date = %s, want 2023-01-01", loginErrorObj.Date)
	}
}

// testPINErrorHandling tests PIN error handling
func testPINErrorHandling(t *testing.T, pinManager *PINManager) {
	// Test invalid PIN format
	err := pinManager.SetPIN("123")
	if err != ErrPINWrongLength {
		t.Errorf("SetPIN() with wrong length returned error: %v, want %v", err, ErrPINWrongLength)
	}

	err = pinManager.SetPIN("123a")
	if err != ErrPINInvalidFormat {
		t.Errorf("SetPIN() with non-digits returned error: %v, want %v", err, ErrPINInvalidFormat)
	}

	// Set valid PIN
	err = pinManager.SetPIN("1234")
	if err != nil {
		t.Fatalf("SetPIN() returned error: %v", err)
	}

	// Test PIN verification with PIN not required
	pinManager.state = PINStateSet
	err = pinManager.VerifyPIN("1234")
	if err != ErrPINNotSet {
		t.Errorf("VerifyPIN() with PIN not required returned error: %v, want %v", err, ErrPINNotSet)
	}

	// Test PIN verification with locked PIN
	pinManager.state = PINStateLocked
	err = pinManager.VerifyPIN("1234")
	if err != ErrPINLocked {
		t.Errorf("VerifyPIN() with locked PIN returned error: %v, want %v", err, ErrPINLocked)
	}
}

// testAntiCheatErrorHandling tests anti-cheat error handling
func testAntiCheatErrorHandling(t *testing.T, antiCheatManager *AntiCheatManager) {
	// Test verify response with anti-cheat disabled
	antiCheatManager.enabled = false
	err := antiCheatManager.VerifyResponse([]byte{1, 2, 3, 4})
	if err != ErrAntiCheatDisabled {
		t.Errorf("VerifyResponse() with anti-cheat disabled returned error: %v, want %v", err, ErrAntiCheatDisabled)
	}

	// Enable anti-cheat
	antiCheatManager.Enable(AntiCheatGameGuard)

	// Test verify response with invalid state
	antiCheatManager.state = AntiCheatStateInitializing
	err = antiCheatManager.VerifyResponse([]byte{1, 2, 3, 4})
	if err != ErrInvalidResponse {
		t.Errorf("VerifyResponse() with invalid state returned error: %v, want %v", err, ErrInvalidResponse)
	}

	// Test verify response with invalid response
	antiCheatManager.state = AntiCheatStateWaitingForResponse
	antiCheatManager.lastChallenge = time.Now()
	antiCheatManager.challenge = []byte{0xFF, 0xFF, 0xFF, 0xFF}
	err = antiCheatManager.VerifyResponse([]byte{0x01, 0x01, 0x01, 0x01})
	if err != ErrInvalidResponse {
		t.Errorf("VerifyResponse() with invalid response returned error: %v, want %v", err, ErrInvalidResponse)
	}

	// Check that anti-cheat state is rejected
	if antiCheatManager.state != AntiCheatStateRejected {
		t.Errorf("Anti-cheat state = %v, want %v after invalid response", antiCheatManager.state, AntiCheatStateRejected)
	}
}
