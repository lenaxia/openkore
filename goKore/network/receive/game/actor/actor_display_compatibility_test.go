package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleActorDisplayCompatibility(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create channels to track hook calls
	preHookCalled := make(chan bool, 1)
	postHookCalled := make(chan bool, 1)
	skipHandling := false

	// Register hooks
	hookManager.AddHook("packet_pre/actor_display", func(hookName string, arg interface{}, userData interface{}) {
		// Modify args if skipHandling is true
		if skipHandling {
			args := arg.(map[string]interface{})
			args["return"] = true
		}
		preHookCalled <- true
	}, nil)

	hookManager.AddHook("packet/actor_display", func(hookName string, arg interface{}, userData interface{}) {
		postHookCalled <- true
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Create test packet arguments
	args := map[string]interface{}{
		"ID":          []byte{1, 2, 3, 4},
		"object_type": byte(0),
		"type":        uint16(0),
		"name":        "TestPlayer",
		"coords":      []byte{10, 20, 30, 40, 50, 60},
		"switch":      "0086", // actor_moved
	}

	// Test case 1: Normal handling
	t.Run("Normal handling", func(t *testing.T) {
		skipHandling = false

		// Call the handler
		err := handler.HandleActorDisplayCompatibility(args)
		if err != nil {
			t.Errorf("HandleActorDisplayCompatibility returned an error: %v", err)
		}

		// Verify pre-hook was called
		select {
		case <-preHookCalled:
			// Success
		default:
			t.Error("Pre-hook was not called")
		}

		// Verify post-hook was called
		select {
		case <-postHookCalled:
			// Success
		default:
			t.Error("Post-hook was not called")
		}
	})

	// Test case 2: Skip handling
	t.Run("Skip handling", func(t *testing.T) {
		skipHandling = true

		// Call the handler
		err := handler.HandleActorDisplayCompatibility(args)
		if err != nil {
			t.Errorf("HandleActorDisplayCompatibility returned an error: %v", err)
		}

		// Verify pre-hook was called
		select {
		case <-preHookCalled:
			// Success
		default:
			t.Error("Pre-hook was not called")
		}

		// Verify post-hook was not called (since we skipped handling)
		select {
		case <-postHookCalled:
			t.Error("Post-hook was called when it should have been skipped")
		default:
			// Success
		}
	})
}
