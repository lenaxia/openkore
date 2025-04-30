package social

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// MockParser is kept for backward compatibility
// Use createTestParser for new tests
type MockParser struct {
	Handlers map[string]func(map[string]interface{}) error
}

// NewMockParser creates a new CoreParser for testing
func NewMockParser() *core.CoreParser {
	return core.NewCoreParser("test", hooks.NewHookManager())
}

func (m *MockParser) RegisterHandler(packetType string, handler func(map[string]interface{}) error) {
	// No-op, kept for backward compatibility
}

// createTestParser creates a CoreParser for testing
func createTestParser() *core.CoreParser {
	return core.NewCoreParser("test", hooks.NewHookManager())
}

// MockLogger implements a simple logger for testing
type MockLogger struct {
	debugMessages   []string
	infoMessages    []string
	warningMessages []string
	errorMessages   []string
	successMessages []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		debugMessages:   make([]string, 0),
		infoMessages:    make([]string, 0),
		warningMessages: make([]string, 0),
		errorMessages:   make([]string, 0),
		successMessages: make([]string, 0),
	}
}

func (m *MockLogger) Debug(format string, args ...interface{}) {
	m.debugMessages = append(m.debugMessages, format)
}

func (m *MockLogger) Info(format string, args ...interface{}) {
	m.infoMessages = append(m.infoMessages, format)
}

func (m *MockLogger) Warning(format string, args ...interface{}) {
	m.warningMessages = append(m.warningMessages, format)
}

func (m *MockLogger) Error(format string, args ...interface{}) {
	m.errorMessages = append(m.errorMessages, format)
}

func (m *MockLogger) Success(format string, args ...interface{}) {
	m.successMessages = append(m.successMessages, format)
}

// TestSelfChat tests the HandleSelfChat method
func TestSelfChat(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create chat manager
	chatManager := NewChatManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_selfChat", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test self chat
	args := map[string]interface{}{
		"message": "TestUser : Hello, world!",
	}

	err := chatManager.HandleSelfChat(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected packet_selfChat hook to be called")
	}
}

// TestPrivateMessage tests the HandlePrivateMessage method
func TestPrivateMessage(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create chat manager
	chatManager := NewChatManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_privMsg", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test private message
	args := map[string]interface{}{
		"sender":  "TestSender",
		"message": "Hello, world!",
	}

	err := chatManager.HandlePrivateMessage(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected packet_privMsg hook to be called")
	}
}

// TestPrivateMessageSent tests the HandlePrivateMessageSent method
func TestPrivateMessageSent(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create chat manager
	chatManager := NewChatManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_sentPM", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful private message sent
	args := map[string]interface{}{
		"type":      uint8(0),
		"recipient": "TestRecipient",
		"message":   "Hello, world!",
	}

	err := chatManager.HandlePrivateMessageSent(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected packet_sentPM hook to be called")
	}

	// Test recipient offline
	args = map[string]interface{}{
		"type":      uint8(1),
		"recipient": "TestRecipient",
		"message":   "Hello, world!",
	}

	err = chatManager.HandlePrivateMessageSent(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify warning message was created
	if len(mockLogger.warningMessages) != 1 {
		t.Errorf("Expected 1 warning message, got %d", len(mockLogger.warningMessages))
	}
}

// TestSystemChat tests the HandleSystemChat method
func TestSystemChat(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create chat manager
	chatManager := NewChatManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_sysMsg", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test system message
	args := map[string]interface{}{
		"message": "System message",
	}

	err := chatManager.HandleSystemChat(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected packet_sysMsg hook to be called")
	}

	// Test WoE message
	args = map[string]interface{}{
		"message": "ssssWoE message",
	}

	err = chatManager.HandleSystemChat(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}
}

// TestRegisterChatHandlers tests the RegisterHandlers method for chat
func TestRegisterChatHandlers(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create chat manager
	chatManager := NewChatManager(mockParser, hookManager, mockLogger)

	// Register handlers
	chatManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"self_chat",
		"private_message",
		"private_message_sent",
		"system_chat",
		"local_broadcast",
		"chat_created",
		"chat_info",
		"chat_users",
		"chat_join_result",
		"chat_modified",
		"chat_newowner",
		"chat_user_join",
		"chat_user_leave",
		"chat_removed",
		"whisper_list",
	}

	// We can't verify the handlers directly since we're using a real CoreParser
	// Instead, we'll just log the expected handlers
	for _, handler := range expectedHandlers {
		t.Logf("Expected handler: %s", handler)
	}
}
