package core

import (
	"testing"
	"time"

	"github.com/lenaxia/goKore/network"
	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleMoveInterrupt(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	manager := NewAccountManager(parser)

	// Set initial state
	manager.SetNetworkState(network.InGame)
	manager.SetState(AccountStateInGame)

	// Create test packet arguments - this packet doesn't have any specific arguments
	args := map[string]interface{}{}

	// Call handler
	err := manager.handleMoveInterrupt(args)
	if err != nil {
		t.Fatalf("handleMoveInterrupt() returned error: %v", err)
	}

	// Check that session state remains unchanged
	if manager.GetNetworkState() != network.InGame {
		t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.InGame)
	}

	session := manager.GetSession()
	if session.State != AccountStateInGame {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateInGame)
	}

	// Check that last packet time was updated
	if time.Since(session.LastPacketTime) > time.Second {
		t.Errorf("LastPacketTime was not updated")
	}
}

// Test with invalid arguments
func TestHandleMoveInterruptInvalidArgs(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	manager := NewAccountManager(parser)

	// Set initial state
	manager.SetNetworkState(network.InGame)
	manager.SetState(AccountStateInGame)
	// Store the initial time for comparison if needed

	// Create test packet with invalid arguments
	args := map[string]interface{}{
		"invalidArg": "value",
	}

	// Call handler - should still work with invalid args since this handler doesn't use any args
	err := manager.handleMoveInterrupt(args)
	if err != nil {
		t.Fatalf("handleMoveInterrupt() returned error with invalid args: %v", err)
	}

	// Check that session state remains unchanged
	if manager.GetNetworkState() != network.InGame {
		t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.InGame)
	}

	session := manager.GetSession()
	if session.State != AccountStateInGame {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateInGame)
	}

	// Check that last packet time was updated
	if time.Since(session.LastPacketTime) > time.Second {
		t.Errorf("LastPacketTime was not updated")
	}
}
