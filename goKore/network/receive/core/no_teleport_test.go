package core

import (
	"testing"

	"github.com/lenaxia/goKore/network"
	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleNoTeleport(t *testing.T) {
	// Test case 1: Unavailable Area To Teleport (fail code 0)
	t.Run("UnavailableAreaToTeleport", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		parser := NewCoreParser("ServerType0", hookManager)
		logger := NewMockLogger()
		manager := NewAccountManager(parser, hookManager, logger)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create a channel to capture hook calls
		teleportErrorCalled := false
		teleportErrorMessage := ""
		hookManager.RegisterHook("teleport.error", func(args interface{}) {
			teleportErrorCalled = true
			if hookArgs, ok := args.(map[string]interface{}); ok {
				if msg, ok := hookArgs["message"].(string); ok {
					teleportErrorMessage = msg
				}
			}
		})

		// Create a channel to capture AI clear calls
		aiClearCalled := false
		hookManager.RegisterHook("ai.clear", func(args interface{}) {
			aiClearCalled = true
		})

		// Create test packet arguments
		args := map[string]interface{}{
			"fail": byte(0), // Unavailable Area To Teleport
		}

		// Call handler
		err := manager.handleNoTeleport(args)
		if err != nil {
			t.Fatalf("handleNoTeleport() returned error: %v", err)
		}

		// Check that hooks were called correctly
		if !teleportErrorCalled {
			t.Error("teleport.error hook was not called")
		}
		if teleportErrorMessage != "Unavailable Area To Teleport" {
			t.Errorf("teleport.error message = %q, want %q", teleportErrorMessage, "Unavailable Area To Teleport")
		}
		if !aiClearCalled {
			t.Error("ai.clear hook was not called")
		}
	})

	// Test case 2: Unavailable Area To Memo (fail code 1)
	t.Run("UnavailableAreaToMemo", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		parser := NewCoreParser("ServerType0", hookManager)
		logger := NewMockLogger()
		manager := NewAccountManager(parser, hookManager, logger)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create a channel to capture hook calls
		teleportErrorCalled := false
		teleportErrorMessage := ""
		hookManager.RegisterHook("teleport.error", func(args interface{}) {
			teleportErrorCalled = true
			if hookArgs, ok := args.(map[string]interface{}); ok {
				if msg, ok := hookArgs["message"].(string); ok {
					teleportErrorMessage = msg
				}
			}
		})

		// Create a channel to capture AI clear calls
		aiClearCalled := false
		hookManager.RegisterHook("ai.clear", func(args interface{}) {
			aiClearCalled = true
		})

		// Create test packet arguments
		args := map[string]interface{}{
			"fail": byte(1), // Unavailable Area To Memo
		}

		// Call handler
		err := manager.handleNoTeleport(args)
		if err != nil {
			t.Fatalf("handleNoTeleport() returned error: %v", err)
		}

		// Check that hooks were called correctly
		if !teleportErrorCalled {
			t.Error("teleport.error hook was not called")
		}
		if teleportErrorMessage != "Unavailable Area To Memo" {
			t.Errorf("teleport.error message = %q, want %q", teleportErrorMessage, "Unavailable Area To Memo")
		}
		// AI clear should not be called for memo errors
		if aiClearCalled {
			t.Error("ai.clear hook was called but should not have been")
		}
	})

	// Test case 3: Unknown fail code
	t.Run("UnknownFailCode", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		parser := NewCoreParser("ServerType0", hookManager)
		logger := NewMockLogger()
		manager := NewAccountManager(parser, hookManager, logger)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create a channel to capture hook calls
		teleportErrorCalled := false
		teleportErrorMessage := ""
		hookManager.RegisterHook("teleport.error", func(args interface{}) {
			teleportErrorCalled = true
			if hookArgs, ok := args.(map[string]interface{}); ok {
				if msg, ok := hookArgs["message"].(string); ok {
					teleportErrorMessage = msg
				}
			}
		})

		// Create test packet arguments
		args := map[string]interface{}{
			"fail": byte(2), // Unknown fail code
		}

		// Call handler
		err := manager.handleNoTeleport(args)
		if err != nil {
			t.Fatalf("handleNoTeleport() returned error: %v", err)
		}

		// Check that hooks were called correctly
		if !teleportErrorCalled {
			t.Error("teleport.error hook was not called")
		}
		if teleportErrorMessage != "Unavailable Area To Teleport (fail code 2)" {
			t.Errorf("teleport.error message = %q, want %q", teleportErrorMessage, "Unavailable Area To Teleport (fail code 2)")
		}
	})

	// Test case 4: Missing fail code
	t.Run("MissingFailCode", func(t *testing.T) {
		hookManager := hooks.NewHookManager()
		parser := NewCoreParser("ServerType0", hookManager)
		logger := NewMockLogger()
		manager := NewAccountManager(parser, hookManager, logger)

		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create test packet arguments with missing fail code
		args := map[string]interface{}{}

		// Call handler
		err := manager.handleNoTeleport(args)
		if err == nil {
			t.Fatalf("handleNoTeleport() should return error for missing fail code")
		}
	})
}
