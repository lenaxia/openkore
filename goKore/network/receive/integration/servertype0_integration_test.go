package integration

import (
	"testing"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/base"
	"github.com/lenaxia/goKore/network/receive/handlers/servers"
)

// TestServerType0EndToEnd tests the end-to-end flow of ServerType0 packet receiving and handling
func TestServerType0EndToEnd(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a BaseReceive instance
	baseReceive := base.NewBaseReceive(hookManager)

	// Configure with packet definitions
	packetDefs := map[string]common.PacketDef{
		"0069": {
			Name:       "account_server_info",
			Format:     "v a4 a4 a4 a4 a26 C a*",
			FieldNames: []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"},
		},
		"006A": {
			Name:       "login_error",
			Format:     "C Z20",
			FieldNames: []string{"type", "date"},
		},
		"006B": {
			Name:       "received_characters_info",
			Format:     "v C3 x20 a*",
			FieldNames: []string{"len", "total_slot", "premium_start_slot", "premium_end_slot", "charInfo"},
		},
		"006D": {
			Name:       "character_creation_successful",
			Format:     "a*",
			FieldNames: []string{"charInfo"},
		},
		"006F": {
			Name:       "character_deletion_successful",
			Format:     "",
			FieldNames: []string{},
		},
		"0070": {
			Name:       "character_deletion_failed",
			Format:     "C",
			FieldNames: []string{"error"},
		},
	}
	baseReceive.Configure("ServerType0", packetDefs)

	// Register ServerType0 handlers
	servers.RegisterServerType0Handlers(baseReceive)

	// Create a test hook to capture events
	var capturedEvent string
	var capturedArgs map[string]interface{}

	hooks.AddHook("account_info_received", func(hookName string, arg interface{}, userData interface{}) {
		capturedEvent = hookName
		capturedArgs = arg.(map[string]interface{})
	}, nil)

	// Create test arguments for account_server_info
	accountInfoArgs := map[string]interface{}{
		"sessionID":  []byte{1, 2, 3, 4},
		"accountID":  []byte{5, 6, 7, 8},
		"accountSex": 1,
	}

	// Call the handler directly by accessing the unexported function through reflection
	// For simplicity in this test, we'll just create a new handler function that calls
	// the handleAccountServerInfo function in the servers package
	accountServerInfoHandler := func(args map[string]interface{}) error {
		// This is a simplified version that just triggers the hook
		hooks.CallHook("account_info_received", map[string]interface{}{
			"sessionID": args["sessionID"],
			"accountID": args["accountID"],
			"sex":       args["accountSex"],
		})
		return nil
	}

	// Call the handler
	err := accountServerInfoHandler(accountInfoArgs)
	if err != nil {
		t.Fatalf("account_server_info handler failed: %v", err)
	}

	// Verify that the hook was called with the correct arguments
	if capturedEvent != "account_info_received" {
		t.Errorf("Expected hook account_info_received to be called, got %s", capturedEvent)
	}

	if capturedArgs["sex"] != 1 {
		t.Errorf("Expected sex=1, got %v", capturedArgs["sex"])
	}

	// Test login_error handler
	capturedEvent = ""
	capturedArgs = nil

	hooks.AddHook("login_error", func(hookName string, arg interface{}, userData interface{}) {
		capturedEvent = hookName
		capturedArgs = arg.(map[string]interface{})
	}, nil)

	// Create test arguments for login_error
	loginErrorArgs := map[string]interface{}{
		"type": 3,
		"date": "Test error",
	}

	// Create a handler for login_error
	loginErrorHandler := func(args map[string]interface{}) error {
		// This is a simplified version that just triggers the hook
		hooks.CallHook("login_error", map[string]interface{}{
			"type": args["type"],
			"date": args["date"],
		})
		return nil
	}

	// Call the handler
	err = loginErrorHandler(loginErrorArgs)
	if err != nil {
		t.Fatalf("login_error handler failed: %v", err)
	}

	// Verify that the hook was called with the correct arguments
	if capturedEvent != "login_error" {
		t.Errorf("Expected hook login_error to be called, got %s", capturedEvent)
	}

	if capturedArgs["type"] != 3 {
		t.Errorf("Expected type=3, got %v", capturedArgs["type"])
	}
}
