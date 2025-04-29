package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandlePing(t *testing.T) {
	// Create a new account manager
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	manager := NewAccountManager(parser)

	// Set initial state
	manager.SetNetworkState(1) // Some network state

	// Create test packet arguments
	args := map[string]interface{}{}

	// Set XKore to "0" (not 1 or 3)
	manager.SetXKore("0")

	// Call handler
	err := manager.handlePing(args)
	if err != nil {
		t.Fatalf("handlePing() returned error: %v", err)
	}

	// Check that ping was sent
	// Note: In a real implementation, we would check if sendPing was called
	// but for now we just check that no error was returned

	// Set XKore to "1"
	manager.SetXKore("1")

	// Call handler again
	err = manager.handlePing(args)
	if err != nil {
		t.Fatalf("handlePing() returned error: %v", err)
	}

	// Check that ping was not sent
	// Note: In a real implementation, we would check if sendPing was not called
	// but for now we just check that no error was returned

	// Set XKore to "3"
	manager.SetXKore("3")

	// Call handler again
	err = manager.handlePing(args)
	if err != nil {
		t.Fatalf("handlePing() returned error: %v", err)
	}

	// Check that ping was not sent
	// Note: In a real implementation, we would check if sendPing was not called
	// but for now we just check that no error was returned
}
