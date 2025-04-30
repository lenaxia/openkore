package send

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// Happy path test - normal creation with valid hook manager
func TestNewBaseSend(t *testing.T) {
	// Create a new hook manager
	hookManager := hooks.NewHookManager()

	// Create a new BaseSend instance
	baseSend := NewBaseSend(hookManager)

	// Check that the BaseSend instance is not nil
	if baseSend == nil {
		t.Fatal("NewBaseSend() returned nil")
	}
}

// Edge case - creating with nil hook manager
func TestNewBaseSendWithNilHookManager(t *testing.T) {
	// Create a new BaseSend instance with nil hook manager
	baseSend := NewBaseSend(nil)

	// Check that the BaseSend instance is not nil
	// This test will fail if the implementation doesn't handle nil hook managers
	if baseSend == nil {
		t.Fatal("NewBaseSend() with nil hook manager returned nil")
	}
}

// Test type re-exports
func TestTypeReExports(t *testing.T) {
	// Create instances of re-exported types to ensure they work correctly
	var pc PacketConstruction
	var sh SendHandler
	var s Send
	var bs BaseSend

	// Just verify that the types exist and can be instantiated
	// This is a compile-time check, but we include it for completeness
	_ = pc
	_ = sh
	_ = s
	_ = bs
}
