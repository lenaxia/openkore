package security

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestHandleLoginErrorGameLoginServer(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)

	// Create a hook manager to capture the hook call
	hookManager := hooks.NewHookManager()
	var hookCalled bool
	var hookCode int
	var hookMessage string

	hookManager.AddHook("security/login_error", func(hookName string, arg interface{}, userData interface{}) {
		hookCalled = true
		if data, ok := arg.(map[string]interface{}); ok {
			if code, ok := data["code"].(int); ok {
				hookCode = code
			}
			if message, ok := data["message"].(string); ok {
				hookMessage = message
			}
		}
	}, nil)

	manager := NewLoginManager(parser, hookManager)

	// Set initial state
	manager.state = LoginStateLoggedIn

	// Test with error type = 0 (REFUSE_INVALID_ID)
	args := map[string]interface{}{
		"type": byte(0),
	}

	err := manager.handleLoginErrorGameLoginServer(args)
	if err != nil {
		t.Fatalf("handleLoginErrorGameLoginServer() returned error: %v", err)
	}

	// Check that login error was set
	if manager.loginError == nil {
		t.Fatal("manager.loginError is nil")
	}

	if manager.loginError.Code != 0 {
		t.Errorf("manager.loginError.Code = %d, want 0", manager.loginError.Code)
	}

	expectedMessage := "Account name doesn't exist"
	if manager.loginError.Message != expectedMessage {
		t.Errorf("manager.loginError.Message = %s, want %s", manager.loginError.Message, expectedMessage)
	}

	// Check that state was updated to disconnected
	if manager.state != LoginStateDisconnected {
		t.Errorf("manager.state = %v, want %v", manager.state, LoginStateDisconnected)
	}

	// Check that hook was called with correct arguments
	if !hookCalled {
		t.Error("Hook was not called")
	}

	if hookCode != 0 {
		t.Errorf("Hook code = %d, want 0", hookCode)
	}

	if hookMessage != expectedMessage {
		t.Errorf("Hook message = %s, want %s", hookMessage, expectedMessage)
	}
}
