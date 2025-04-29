package behavior

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestMannerMessage(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the behavior manager
	manager := NewCharacterBehaviorManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases
	testCases := []struct {
		name           string
		flag           uint8
		expectedResult string
	}{
		{
			name:           "Manner point aligned",
			flag:           0,
			expectedResult: "A manner point has been successfully aligned.",
		},
		{
			name:           "Chat block by GM",
			flag:           3,
			expectedResult: "Chat Block has been applied by GM due to your ill-mannerous action.",
		},
		{
			name:           "Chat block by anti-spam",
			flag:           4,
			expectedResult: "Automated Chat Block has been applied due to Anti-Spam System.",
		},
		{
			name:           "Good point received",
			flag:           5,
			expectedResult: "You got a good point.",
		},
		{
			name:           "Unknown flag",
			flag:           99,
			expectedResult: "Unknown manner message result (flag: 99)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook("character.manner_message", func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Create packet data
			args := map[string]interface{}{
				"flag": tc.flag,
			}

			// Call the handler
			err := manager.handleMannerMessage(args)
			if err != nil {
				t.Errorf("handleMannerMessage() returned error: %v", err)
			}

			// Check that the hook was called
			if !hookCalled {
				t.Error("Hook was not called")
			}

			// Check the hook result
			if hookResult == nil {
				t.Fatal("Hook result is nil")
			}

			// Check the flag
			if flag, ok := hookResult["flag"].(uint8); !ok || flag != tc.flag {
				t.Errorf("Expected flag %d, got %v", tc.flag, hookResult["flag"])
			}

			// Check the message
			if message, ok := hookResult["message"].(string); !ok || message != tc.expectedResult {
				t.Errorf("Expected message %q, got %q", tc.expectedResult, message)
			}
		})
	}
}

func TestHackShieldAlarm(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the behavior manager
	manager := NewCharacterBehaviorManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Create a channel to receive hook events
	hookCalled := false
	var hookResult map[string]interface{}

	// Register a hook to capture the event
	hookManager.AddHook("character.hack_shield_alarm", func(hookName string, arg interface{}, userData interface{}) {
		hookCalled = true
		if result, ok := arg.(map[string]interface{}); ok {
			hookResult = result
		}
	}, nil)

	// Call the handler
	err := manager.handleHackShieldAlarm(map[string]interface{}{})
	if err != nil {
		t.Errorf("handleHackShieldAlarm() returned error: %v", err)
	}

	// Check that the hook was called
	if !hookCalled {
		t.Error("Hook was not called")
	}

	// Check the hook result
	if hookResult == nil {
		t.Fatal("Hook result is nil")
	}

	// Check the message
	expectedMessage := "Error: You have been forced to disconnect by a Hack Shield. Please check Poseidon."
	if message, ok := hookResult["message"].(string); !ok || message != expectedMessage {
		t.Errorf("Expected message %q, got %q", expectedMessage, message)
	}
}

func TestRegisterHandlers(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the behavior manager
	manager := NewCharacterBehaviorManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Check that the handlers were registered
	mannerHandler, found := parser.GetHandler("manner_message")
	if !found || mannerHandler == nil {
		t.Error("manner_message handler was not registered")
	}

	hackShieldHandler, found := parser.GetHandler("hack_shield_alarm")
	if !found || hackShieldHandler == nil {
		t.Error("hack_shield_alarm handler was not registered")
	}
}
