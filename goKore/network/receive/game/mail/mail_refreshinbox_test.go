package mail

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TestMailRefreshInbox tests the HandleMailRefreshInbox method
func TestMailRefreshInbox(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create mail manager
	mailManager := NewMailManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("mail_refreshinbox", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test mail refresh inbox with no mails
	args := map[string]interface{}{
		"count": uint32(0),
	}

	err := mailManager.HandleMailRefreshInbox(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify correct message
	expectedMsg := "There is no mail in your inbox."
	if mockLogger.infoMessages[0] != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, mockLogger.infoMessages[0])
	}

	// Test mail refresh inbox with mails
	hookCalled = false
	mockLogger = NewMockLogger()
	mailManager.logger = mockLogger

	// Create a mock raw message with 2 mail entries
	// Each mail entry is 73 bytes:
	// - mailID (4 bytes)
	// - title (40 bytes)
	// - read flag (1 byte)
	// - sender (24 bytes)
	// - timestamp (4 bytes)
	rawMsg := make([]byte, 8+2*73) // 8 bytes header + 2 mail entries

	// Mail 1
	// mailID = 1
	rawMsg[8] = 1
	// title = "Test Mail 1"
	copy(rawMsg[8+4:], []byte("Test Mail 1"))
	// read = 0 (unread)
	rawMsg[8+44] = 0
	// sender = "Sender 1"
	copy(rawMsg[8+45:], []byte("Sender 1"))
	// timestamp = 1234567890
	rawMsg[8+69] = 210
	rawMsg[8+70] = 2
	rawMsg[8+71] = 150
	rawMsg[8+72] = 73

	// Mail 2
	// mailID = 2
	rawMsg[8+73] = 2
	// title = "Test Mail 2"
	copy(rawMsg[8+73+4:], []byte("Test Mail 2"))
	// read = 1 (read)
	rawMsg[8+73+44] = 1
	// sender = "Sender 2"
	copy(rawMsg[8+73+45:], []byte("Sender 2"))
	// timestamp = 1234567891
	rawMsg[8+73+69] = 211
	rawMsg[8+73+70] = 2
	rawMsg[8+73+71] = 150
	rawMsg[8+73+72] = 73

	args = map[string]interface{}{
		"count":   uint32(2),
		"RAW_MSG": rawMsg,
	}

	err = mailManager.HandleMailRefreshInbox(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) != 3 { // 1 for count + 2 for mails
		t.Errorf("Expected 3 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify mailList was populated
	if len(mailManager.mailList) != 2 {
		t.Errorf("Expected 2 mails in mailList, got %d", len(mailManager.mailList))
	}

	// Verify mail 1 data
	if mail, exists := mailManager.mailList[1]; !exists {
		t.Errorf("Expected mail with ID 1 to be in mailList")
	} else {
		if title, ok := mail["title"].(string); !ok || title != "Test Mail 1" {
			t.Errorf("Expected title 'Test Mail 1', got '%v'", mail["title"])
		}
		if read, ok := mail["read"].(bool); !ok || read != false {
			t.Errorf("Expected read to be false, got %v", mail["read"])
		}
		if sender, ok := mail["sender"].(string); !ok || sender != "Sender 1" {
			t.Errorf("Expected sender 'Sender 1', got '%v'", mail["sender"])
		}
	}

	// Verify mail 2 data
	if mail, exists := mailManager.mailList[2]; !exists {
		t.Errorf("Expected mail with ID 2 to be in mailList")
	} else {
		if title, ok := mail["title"].(string); !ok || title != "Test Mail 2" {
			t.Errorf("Expected title 'Test Mail 2', got '%v'", mail["title"])
		}
		if read, ok := mail["read"].(bool); !ok || read != true {
			t.Errorf("Expected read to be true, got %v", mail["read"])
		}
		if sender, ok := mail["sender"].(string); !ok || sender != "Sender 2" {
			t.Errorf("Expected sender 'Sender 2', got '%v'", mail["sender"])
		}
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_refreshinbox hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"count": "invalid", // Invalid type
	}

	err = mailManager.HandleMailRefreshInbox(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
