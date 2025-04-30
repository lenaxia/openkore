package mail

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TestMailRead tests the HandleMailRead method
func TestMailRead(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create mail manager
	mailManager := NewMailManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("mail_read", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test mail read with item and zeny
	args := map[string]interface{}{
		"mailID":  uint32(1),
		"nameID":  uint16(501),
		"title":   "Test Mail",
		"sender":  "TestSender",
		"message": "This is a test message",
		"amount":  uint16(5),
		"zeny":    uint32(1000),
		"upgrade": uint8(7),
		"cards":   []uint16{4001, 4002},
		"broken":  uint8(0),
	}

	err := mailManager.HandleMailRead(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) != 4 {
		t.Errorf("Expected 4 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify mail was added to mailList
	if mail, exists := mailManager.mailList[1]; !exists {
		t.Errorf("Expected mail with ID 1 to be added to mailList")
	} else {
		// Check mail data
		if mail["title"] != "Test Mail" {
			t.Errorf("Expected title 'Test Mail', got '%s'", mail["title"])
		}
		if mail["sender"] != "TestSender" {
			t.Errorf("Expected sender 'TestSender', got '%s'", mail["sender"])
		}
		if mail["message"] != "This is a test message" {
			t.Errorf("Expected message 'This is a test message', got '%s'", mail["message"])
		}
		if mail["zeny"] != uint32(1000) {
			t.Errorf("Expected zeny 1000, got %d", mail["zeny"])
		}

		// Check item data
		if item, ok := mail["item"].(map[string]interface{}); !ok {
			t.Errorf("Expected item to be a map")
		} else {
			if item["nameID"] != uint16(501) {
				t.Errorf("Expected nameID 501, got %d", item["nameID"])
			}
			if item["amount"] != uint16(5) {
				t.Errorf("Expected amount 5, got %d", item["amount"])
			}
			if item["upgrade"] != uint8(7) {
				t.Errorf("Expected upgrade 7, got %d", item["upgrade"])
			}
		}
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_read hook to be called")
	}

	// Test mail read with no item
	hookCalled = false
	mockLogger = NewMockLogger()
	mailManager.logger = mockLogger

	args = map[string]interface{}{
		"mailID":  uint32(2),
		"nameID":  uint16(0),
		"title":   "Test Mail 2",
		"sender":  "TestSender2",
		"message": "This is another test message",
		"zeny":    uint32(2000),
	}

	err = mailManager.HandleMailRead(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created (no item message)
	if len(mockLogger.infoMessages) != 3 {
		t.Errorf("Expected 3 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify mail was added to mailList
	if _, exists := mailManager.mailList[2]; !exists {
		t.Errorf("Expected mail with ID 2 to be added to mailList")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_read hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"mailID": "invalid", // Invalid type
		"title":  "Test Mail",
	}

	err = mailManager.HandleMailRead(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
