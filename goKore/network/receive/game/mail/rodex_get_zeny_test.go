package mail

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TestRodexGetZeny tests the HandleRodexGetZeny method
func TestRodexGetZeny(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create rodex manager
	rodexManager := NewRodexManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("rodex_get_zeny", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Initialize rodexList with a mail
	rodexManager.rodexList = make(map[string]interface{})
	rodexManager.rodexList["mails"] = make(map[uint32]map[string]interface{})
	rodexManager.rodexList["mails"].(map[uint32]map[string]interface{})[1] = map[string]interface{}{
		"mailID1": uint32(1),
		"sender":  "TestSender",
		"title":   "TestTitle",
		"zeny1":   uint32(1000),
		"zeny2":   uint16(0),
	}

	// Test rodex get zeny
	args := map[string]interface{}{
		"mailID1": uint32(1),
		"mailID2": uint32(0),
		"zeny":    uint32(1000),
	}

	err := rodexManager.HandleRodexGetZeny(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify mail was updated
	if mail, ok := rodexManager.rodexList["mails"].(map[uint32]map[string]interface{})[1]; ok {
		// Check if zeny was reset
		if zeny1, ok := mail["zeny1"].(uint32); !ok || zeny1 != 0 {
			t.Errorf("Expected zeny1 to be 0, got %v", mail["zeny1"])
		}

		if zeny2, ok := mail["zeny2"].(uint16); !ok || zeny2 != 0 {
			t.Errorf("Expected zeny2 to be 0, got %v", mail["zeny2"])
		}
	} else {
		t.Errorf("Mail with ID 1 not found in rodexList")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected rodex_get_zeny hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"mailID1": "invalid", // Invalid type
		"mailID2": uint32(0),
		"zeny":    uint32(1000),
	}

	err = rodexManager.HandleRodexGetZeny(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
