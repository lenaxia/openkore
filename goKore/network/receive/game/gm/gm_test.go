package gm

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestGMSilence(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("gm.silence", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewGMManager(nil, hookManager)

	// Test case 1: Muted
	t.Run("Muted", func(t *testing.T) {
		args := map[string]interface{}{
			"flag": byte(1),
			"name": []byte("GameMaster\x00"),
		}

		// Call the handler
		err := manager.HandleGMSilence(args)
		if err != nil {
			t.Errorf("HandleGMSilence returned an error: %v", err)
		}

		// Get the result from the channel
		result := <-resultChan

		// Verify the result
		if result["flag"] != byte(1) {
			t.Errorf("Expected flag 1, got %d", result["flag"])
		}
		if result["name"] != "GameMaster" {
			t.Errorf("Expected name 'GameMaster', got '%s'", result["name"])
		}
		if result["status"] != "You have been: muted by GameMaster." {
			t.Errorf("Expected status 'You have been: muted by GameMaster.', got '%s'", result["status"])
		}
	})

	// Test case 2: Unmuted
	t.Run("Unmuted", func(t *testing.T) {
		args := map[string]interface{}{
			"flag": byte(0),
			"name": []byte("GameMaster\x00"),
		}

		// Call the handler
		err := manager.HandleGMSilence(args)
		if err != nil {
			t.Errorf("HandleGMSilence returned an error: %v", err)
		}

		// Get the result from the channel
		result := <-resultChan

		// Verify the result
		if result["flag"] != byte(0) {
			t.Errorf("Expected flag 0, got %d", result["flag"])
		}
		if result["name"] != "GameMaster" {
			t.Errorf("Expected name 'GameMaster', got '%s'", result["name"])
		}
		if result["status"] != "You have been: unmuted by GameMaster." {
			t.Errorf("Expected status 'You have been: unmuted by GameMaster.', got '%s'", result["status"])
		}
	})
}

func TestGMReqAccName(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("gm.req_acc_name", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewGMManager(nil, hookManager)

	// Test case
	args := map[string]interface{}{
		"targetID":    uint32(12345),
		"accountName": "TestAccount",
	}

	// Call the handler
	err := manager.HandleGMReqAccName(args)
	if err != nil {
		t.Errorf("HandleGMReqAccName returned an error: %v", err)
	}

	// Get the result from the channel
	result := <-resultChan

	// Verify the result
	if result["targetID"] != uint32(12345) {
		t.Errorf("Expected targetID 12345, got %d", result["targetID"])
	}
	if result["accountName"] != "TestAccount" {
		t.Errorf("Expected accountName 'TestAccount', got '%s'", result["accountName"])
	}
	if result["status"] != "The accountName for ID 12345 is TestAccount." {
		t.Errorf("Expected status 'The accountName for ID 12345 is TestAccount.', got '%s'", result["status"])
	}
}
