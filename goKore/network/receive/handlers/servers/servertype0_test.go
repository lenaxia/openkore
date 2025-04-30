package servers

import (
	"testing"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/base"
)

// setupTestBaseReceive creates a configured BaseReceive for testing
func setupTestBaseReceive() *base.BaseReceive {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a BaseReceive
	baseReceive := base.NewBaseReceive(hookManager)

	// Configure with test packet definitions
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

	// Configure the BaseReceive
	baseReceive.Configure("ServerType0", packetDefs)

	return baseReceive
}

// TestRegisterServerType0Handlers tests that ServerType0 handlers are properly registered
func TestRegisterServerType0Handlers(t *testing.T) {
	// Create a configured BaseReceive
	mockReceive := setupTestBaseReceive()

	// Register handlers
	RegisterServerType0Handlers(mockReceive)

	// In a real implementation with actual handlers, we would verify that they were registered
	// For now, we just make sure the function doesn't panic
}

// TestServerType0Integration tests the integration between ServerType0 handlers and the BaseReceive
func TestServerType0Integration(t *testing.T) {
	// Create a configured BaseReceive
	baseReceive := setupTestBaseReceive()

	// Register handlers
	RegisterServerType0Handlers(baseReceive)

	// Create a test hook to capture events
	var capturedEvent string
	var capturedArgs map[string]interface{}

	hooks.AddHook("account_info_received", func(hookName string, arg interface{}, userData interface{}) {
		capturedEvent = hookName
		capturedArgs = arg.(map[string]interface{})
	}, nil)

	// Create a test packet for account_server_info
	// This is a simplified packet - in a real implementation, we would create a proper packet
	// with the correct format and data
	args := map[string]interface{}{
		"sessionID":  []byte{1, 2, 3, 4},
		"accountID":  []byte{5, 6, 7, 8},
		"accountSex": 1,
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

	if capturedArgs["sex"] != 1 {
		t.Errorf("Expected sex=1, got %v", capturedArgs["sex"])
	}
}
