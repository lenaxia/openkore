package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleFlag(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a character manager
	parser := NewCoreParser("ServerType0", hookManager)
	charManager := NewCharacterManager(parser)

	// Register a hook to capture the flag event
	var flagCalled bool
	var capturedArgs map[string]interface{}
	hookManager.AddHook("character.flag", func(hookName string, arg interface{}, userData interface{}) {
		flagCalled = true
		capturedArgs = arg.(map[string]interface{})
	}, nil)

	// Test case: Basic flag handling
	t.Run("BasicFlagHandling", func(t *testing.T) {
		// Reset flag called
		flagCalled = false
		capturedArgs = nil

		// Create test args
		testArgs := map[string]interface{}{
			"test_key": "test_value",
		}

		// Call the handler
		err := charManager.handleFlag(testArgs)

		// Verify results
		if err != nil {
			t.Fatalf("handleFlag() returned error: %v", err)
		}

		// Check that the hook was called
		if !flagCalled {
			t.Error("character.flag hook was not called")
		}

		// Check that the args were passed correctly
		if capturedArgs == nil {
			t.Fatal("capturedArgs is nil")
		}

		if val, ok := capturedArgs["test_key"]; !ok || val != "test_value" {
			t.Errorf("capturedArgs[\"test_key\"] = %v, want \"test_value\"", val)
		}
	})
}
