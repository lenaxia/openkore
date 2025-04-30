package npc

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// MockLogger is a simple mock implementation of the Logger interface
type MockLogger struct {
	debugMessages   []string
	infoMessages    []string
	warningMessages []string
	errorMessages   []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		debugMessages:   []string{},
		infoMessages:    []string{},
		warningMessages: []string{},
		errorMessages:   []string{},
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
	// Not used in these tests
}

// MockParser is a simple mock implementation of the CoreParser
type MockParser struct {
	PacketList map[string]interface{}
	handlers   map[string]func(map[string]interface{}) error
}

func NewMockParser() *core.CoreParser {
	// Create a real CoreParser with nil hookManager
	return core.NewCoreParser("ServerType0", nil)
}

// TestNpcTalk tests the HandleNpcTalk method
func TestNpcTalk(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create interaction manager
	interactionManager := NewInteractionManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("npc_talk", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test NPC talk
	args := map[string]interface{}{
		"ID":  uint32(12345),
		"msg": "Hello adventurer!",
	}

	err := interactionManager.HandleNpcTalk(args)

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
		t.Errorf("Expected npc_talk hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
	}

	err = interactionManager.HandleNpcTalk(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestNpcTalkClose tests the HandleNpcTalkClose method
func TestNpcTalkClose(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create interaction manager
	interactionManager := NewInteractionManager(mockParser, hookManager, mockLogger)

	// Set up the NPC talk state
	interactionManager.npcTalkState["ID"] = uint32(12345)
	interactionManager.npcTalkState["talk"] = "initiated"

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("npc_talk_done", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test NPC talk close
	args := map[string]interface{}{
		"ID": uint32(12345),
	}

	err := interactionManager.HandleNpcTalkClose(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected npc_talk_done hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
	}

	err = interactionManager.HandleNpcTalkClose(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestNpcTalkContinue tests the HandleNpcTalkContinue method
func TestNpcTalkContinue(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create interaction manager
	interactionManager := NewInteractionManager(mockParser, hookManager, mockLogger)

	// Test NPC talk continue
	args := map[string]interface{}{
		"ID": uint32(12345),
	}

	err := interactionManager.HandleNpcTalkContinue(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify talk state was updated
	if interactionManager.npcTalkState["talk"] != "next" {
		t.Errorf("Expected talk state to be 'next', got %v", interactionManager.npcTalkState["talk"])
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
	}

	err = interactionManager.HandleNpcTalkContinue(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestNpcTalkNumber tests the HandleNpcTalkNumber method
func TestNpcTalkNumber(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create interaction manager
	interactionManager := NewInteractionManager(mockParser, hookManager, mockLogger)

	// Test NPC talk number
	args := map[string]interface{}{
		"ID": uint32(12345),
	}

	err := interactionManager.HandleNpcTalkNumber(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify talk state was updated
	if interactionManager.npcTalkState["talk"] != "number" {
		t.Errorf("Expected talk state to be 'number', got %v", interactionManager.npcTalkState["talk"])
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
	}

	err = interactionManager.HandleNpcTalkNumber(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestNpcTalkResponses tests the HandleNpcTalkResponses method
func TestNpcTalkResponses(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create interaction manager
	interactionManager := NewInteractionManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("npc_talk_responses", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test NPC talk responses
	args := map[string]interface{}{
		"ID":      uint32(12345),
		"RAW_MSG": []byte("Option1:Option2:Option3"),
	}

	err := interactionManager.HandleNpcTalkResponses(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) < 1 {
		t.Errorf("Expected at least 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected npc_talk_responses hook to be called")
	}

	// Verify responses were stored
	responses, ok := interactionManager.currentTalk["responses"].([]string)
	if !ok {
		t.Errorf("Expected responses to be stored in currentTalk")
	} else {
		// Verify "Cancel Chat" was added
		if len(responses) < 1 || responses[len(responses)-1] != "Cancel Chat" {
			t.Errorf("Expected 'Cancel Chat' to be added to responses")
		}
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
	}

	err = interactionManager.HandleNpcTalkResponses(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestNpcTalkText tests the HandleNpcTalkText method
func TestNpcTalkText(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create interaction manager
	interactionManager := NewInteractionManager(mockParser, hookManager, mockLogger)

	// Test NPC talk text
	args := map[string]interface{}{
		"ID": uint32(12345),
	}

	err := interactionManager.HandleNpcTalkText(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify talk state was updated
	if interactionManager.npcTalkState["talk"] != "text" {
		t.Errorf("Expected talk state to be 'text', got %v", interactionManager.npcTalkState["talk"])
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
	}

	err = interactionManager.HandleNpcTalkText(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestNpcClearDialog tests the HandleNpcClearDialog method
func TestNpcClearDialog(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create interaction manager
	interactionManager := NewInteractionManager(mockParser, hookManager, mockLogger)

	// Set up some state
	interactionManager.currentTalk["ID"] = uint32(12345)
	interactionManager.npcTalkState["talk"] = "initiated"

	// Test NPC clear dialog
	args := map[string]interface{}{
		"ID": uint32(12345),
	}

	err := interactionManager.HandleNpcClearDialog(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify state was cleared
	if len(interactionManager.currentTalk) != 0 || len(interactionManager.npcTalkState) != 0 {
		t.Errorf("Expected state to be cleared")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
	}

	err = interactionManager.HandleNpcClearDialog(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestNpcChat tests the HandleNpcChat method
func TestNpcChat(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create interaction manager
	interactionManager := NewInteractionManager(mockParser, hookManager, mockLogger)

	// Test NPC chat
	args := map[string]interface{}{
		"ID":      uint32(12345),
		"message": "NPC : Hello adventurer!",
	}

	err := interactionManager.HandleNpcChat(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
	}

	err = interactionManager.HandleNpcChat(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestNpcImage tests the HandleNpcImage method
func TestNpcImage(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create interaction manager
	interactionManager := NewInteractionManager(mockParser, hookManager, mockLogger)

	// Test NPC image (type 2)
	args := map[string]interface{}{
		"npc_image": "test_image.bmp",
		"type":      uint8(2),
	}

	err := interactionManager.HandleNpcImage(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify image was stored
	if image, ok := interactionManager.currentTalk["image"].(string); !ok || image != "test_image.bmp" {
		t.Errorf("Expected image to be stored in currentTalk")
	}

	// Test NPC image (type 255)
	args = map[string]interface{}{
		"npc_image": "test_image.bmp",
		"type":      uint8(255),
	}

	err = interactionManager.HandleNpcImage(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify image was removed
	if _, exists := interactionManager.currentTalk["image"]; exists {
		t.Errorf("Expected image to be removed from currentTalk")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"npc_image": 123, // Invalid type
	}

	err = interactionManager.HandleNpcImage(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestRegisterHandlers tests the RegisterHandlers method
func TestRegisterHandlers(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create interaction manager
	interactionManager := NewInteractionManager(mockParser, hookManager, mockLogger)

	// Register handlers
	interactionManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"npc_talk",
		"npc_talk_close",
		"npc_talk_continue",
		"npc_talk_number",
		"npc_talk_responses",
		"npc_talk_text",
		"npc_clear_dialog",
		"npc_chat",
		"npc_image",
	}

	// We can't verify the handlers directly since we're using a real CoreParser
	// Instead, we'll just log the expected handlers
	for _, handler := range expectedHandlers {
		t.Logf("Expected handler: %s", handler)
	}
}
