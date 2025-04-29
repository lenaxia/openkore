package misc

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestDynamicNPCCreateResult(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("misc.dynamicnpc_create_result", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewMiscManager(nil, hookManager)

	// Test cases
	testCases := []struct {
		name           string
		args           map[string]interface{}
		expectedStatus string
	}{
		{
			name: "Success Result",
			args: map[string]interface{}{
				"result": uint8(0), // DYNAMICNPC_RESULT_SUCCESS
			},
			expectedStatus: "Success",
		},
		{
			name: "Unknown Result",
			args: map[string]interface{}{
				"result": uint8(1), // DYNAMICNPC_RESULT_UNKNOWN
			},
			expectedStatus: "Unknown",
		},
		{
			name: "Unknown NPC Result",
			args: map[string]interface{}{
				"result": uint8(2), // DYNAMICNPC_RESULT_UNKNOWNNPC
			},
			expectedStatus: "Unknown NPC",
		},
		{
			name: "Duplicate Result",
			args: map[string]interface{}{
				"result": uint8(3), // DYNAMICNPC_RESULT_DUPLICATE
			},
			expectedStatus: "Duplicate",
		},
		{
			name: "Out of Time Result",
			args: map[string]interface{}{
				"result": uint8(4), // DYNAMICNPC_RESULT_OUTOFTIME
			},
			expectedStatus: "Out of time",
		},
		{
			name: "Invalid Result",
			args: map[string]interface{}{
				"result": uint8(99), // Invalid result
			},
			expectedStatus: "Unknown Result: 99",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := manager.handleDynamicNPCCreateResult(tc.args)
			if err != nil {
				t.Errorf("handleDynamicNPCCreateResult returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-resultChan

			// Verify the result
			if result["status"] != tc.expectedStatus {
				t.Errorf("Expected status %q, got %q", tc.expectedStatus, result["status"])
			}
			if result["result"] != tc.args["result"] {
				t.Errorf("Expected result %v, got %v", tc.args["result"], result["result"])
			}
		})
	}
}

func TestRegisterWithParser(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a parser for testing
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Register the handlers
	// This is a bit tricky to test directly, so we'll just check that the registration doesn't panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("RegisterWithParser panicked: %v", r)
			}
		}()
		RegisterWithCoreParser(parser, hookManager)
	}()
}

func TestRegisterWithBaseReceive(t *testing.T) {
	// Create a base receive for testing
	baseReceive := core.NewBaseReceive(nil)

	// Register the handlers
	// This is a bit tricky to test directly, so we'll just check that the registration doesn't panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("RegisterWithBaseReceive panicked: %v", r)
			}
		}()
		RegisterWithBaseReceive(baseReceive)
	}()
}
