package servers

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/base"
)

func TestAccountServerInfoPacket_DirectCall(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a BaseReceive instance
	baseReceive := base.NewBaseReceive(hookManager)

	// Create a test hook to capture events
	var capturedEvent string
	var capturedArgs map[string]interface{}

	hooks.AddHook("account_info_received", func(hookName string, arg interface{}, userData interface{}) {
		capturedEvent = hookName
		capturedArgs = arg.(map[string]interface{})
	}, nil)

	// Create test arguments for account_server_info
	args := map[string]interface{}{
		"sessionID":  []byte{0xCD, 0x1F, 0x6F, 0xBC}, // Session ID (3161399245)
		"accountID":  []byte{0x82, 0x84, 0x1E, 0x00}, // Account ID (2000002)
		"sessionID2": []byte{0x9B, 0x94, 0xE2, 0x3E}, // Session ID2 (1055036571)
		"accountSex": 1,                              // Sex (1 = male)
	}

	// Call the handler directly
	err := handleAccountServerInfo(args, baseReceive)
	if err != nil {
		t.Fatalf("handleAccountServerInfo failed: %v", err)
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

func TestLoginErrorPacket_DirectCall(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a BaseReceive instance
	baseReceive := base.NewBaseReceive(hookManager)

	// Create a test hook to capture events
	var capturedEvent string
	var capturedArgs map[string]interface{}

	hooks.AddHook("login_error", func(hookName string, arg interface{}, userData interface{}) {
		capturedEvent = hookName
		capturedArgs = arg.(map[string]interface{})
	}, nil)

	// Create test arguments for login_error
	args := map[string]interface{}{
		"type": 1,
		"date": "Error message",
	}

	// Call the handler directly
	err := handleLoginError(args, baseReceive)
	if err != nil {
		t.Fatalf("handleLoginError failed: %v", err)
	}

	// Verify that the hook was called with the correct arguments
	if capturedEvent != "login_error" {
		t.Errorf("Expected hook login_error to be called, got %s", capturedEvent)
	}

	// Verify error type
	errorType, ok := capturedArgs["type"].(int)
	if !ok || errorType != 1 {
		t.Errorf("Error type not correctly parsed: %v", capturedArgs["type"])
	}
}
