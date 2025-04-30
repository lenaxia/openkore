package effects

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
	successMessages []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		debugMessages:   []string{},
		infoMessages:    []string{},
		warningMessages: []string{},
		errorMessages:   []string{},
		successMessages: []string{},
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

// TestMiscEffect tests the HandleMiscEffect method
func TestMiscEffect(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create effects manager
	effectsManager := NewEffectsManager(mockParser, hookManager, mockLogger)

	// Set effect name
	effectsManager.SetEffectName(1, "Test Effect")

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("misc_effect", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test misc effect with known effect ID
	args := map[string]interface{}{
		"ID":     uint32(12345),
		"effect": uint32(1),
	}

	err := effectsManager.HandleMiscEffect(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected misc_effect hook to be called")
	}

	// Test misc effect with unknown effect ID
	hookCalled = false
	mockLogger = NewMockLogger()
	effectsManager.logger = mockLogger

	args = map[string]interface{}{
		"ID":     uint32(12345),
		"effect": uint32(999), // Unknown effect ID
	}

	err = effectsManager.HandleMiscEffect(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected misc_effect hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID":     "invalid", // Invalid type
		"effect": uint32(1),
	}

	err = effectsManager.HandleMiscEffect(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestSoundEffect tests the HandleSoundEffect method
func TestSoundEffect(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create effects manager
	effectsManager := NewEffectsManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("sound_effect", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test sound effect (play once)
	args := map[string]interface{}{
		"name": "test.wav",
		"type": uint8(SOUND_PLAY_ONCE),
		"term": uint32(0),
		"ID":   uint32(12345),
	}

	err := effectsManager.HandleSoundEffect(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected sound_effect hook to be called")
	}

	// Test sound effect (no actor)
	hookCalled = false
	mockLogger = NewMockLogger()
	effectsManager.logger = mockLogger

	args = map[string]interface{}{
		"name": "test.wav",
		"type": uint8(SOUND_PLAY_ONCE),
		"term": uint32(0),
	}

	err = effectsManager.HandleSoundEffect(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected sound_effect hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"name": 123, // Invalid type
		"type": uint8(SOUND_PLAY_ONCE),
		"term": uint32(0),
	}

	err = effectsManager.HandleSoundEffect(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestHatEffect tests the HandleHatEffect method
func TestHatEffect(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create effects manager
	effectsManager := NewEffectsManager(mockParser, hookManager, mockLogger)

	// Set hat effect handle and name
	effectsManager.SetHatEffectHandle(1, "HAT_EFFECT_1")
	effectsManager.SetHatEffectName("HAT_EFFECT_1", "Test Hat Effect")

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("hat_effect", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test hat effect (flag = 1)
	args := map[string]interface{}{
		"ID":     uint32(12345),
		"flag":   uint8(1),
		"effect": []byte{0x01, 0x00}, // Hat effect ID 1
	}

	err := effectsManager.HandleHatEffect(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected hat_effect hook to be called")
	}

	// Test hat effect (flag = 0)
	hookCalled = false
	mockLogger = NewMockLogger()
	effectsManager.logger = mockLogger

	args = map[string]interface{}{
		"ID":     uint32(12345),
		"flag":   uint8(0),
		"effect": []byte{0x01, 0x00}, // Hat effect ID 1
	}

	err = effectsManager.HandleHatEffect(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected hat_effect hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID":     "invalid", // Invalid type
		"flag":   uint8(1),
		"effect": []byte{0x01, 0x00},
	}

	err = effectsManager.HandleHatEffect(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestEmoticon tests the HandleEmoticon method
func TestEmoticon(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create effects manager
	effectsManager := NewEffectsManager(mockParser, hookManager, mockLogger)

	// Set emotion lookup table
	effectsManager.SetEmotionLut(1, "Smile")

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_emotion", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test emoticon with known emotion type
	args := map[string]interface{}{
		"ID":   uint32(12345),
		"type": uint8(1),
	}

	err := effectsManager.HandleEmoticon(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected packet_emotion hook to be called")
	}

	// Test emoticon with unknown emotion type
	hookCalled = false
	mockLogger = NewMockLogger()
	effectsManager.logger = mockLogger

	args = map[string]interface{}{
		"ID":   uint32(12345),
		"type": uint8(255), // Unknown emotion type
	}

	err = effectsManager.HandleEmoticon(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected packet_emotion hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID":   "invalid", // Invalid type
		"type": uint8(1),
	}

	err = effectsManager.HandleEmoticon(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestParseHatEffect tests the ParseHatEffect method
func TestParseHatEffect(t *testing.T) {
	// Create effects manager
	effectsManager := NewEffectsManager(nil, nil, nil)

	// Set hat effect handle and name
	effectsManager.SetHatEffectHandle(1, "HAT_EFFECT_1")
	effectsManager.SetHatEffectName("HAT_EFFECT_1", "Test Hat Effect")

	// Test parsing hat effect
	effect := []byte{0x01, 0x00, 0x02, 0x00} // Hat effect IDs 1 and 2

	hatEffects := effectsManager.ParseHatEffect(effect)

	// Verify hat effects were parsed correctly
	if len(hatEffects) != 2 {
		t.Errorf("Expected 2 hat effects, got %d", len(hatEffects))
	}

	// Verify first hat effect
	if hatEffects[0].HatEFID != 1 {
		t.Errorf("Expected hat effect ID 1, got %d", hatEffects[0].HatEFID)
	}
	if hatEffects[0].Handle != "HAT_EFFECT_1" {
		t.Errorf("Expected hat effect handle HAT_EFFECT_1, got %s", hatEffects[0].Handle)
	}
	if hatEffects[0].Name != "Test Hat Effect" {
		t.Errorf("Expected hat effect name 'Test Hat Effect', got %s", hatEffects[0].Name)
	}

	// Verify second hat effect (unknown)
	if hatEffects[1].HatEFID != 2 {
		t.Errorf("Expected hat effect ID 2, got %d", hatEffects[1].HatEFID)
	}
	if hatEffects[1].Handle != "" {
		t.Errorf("Expected empty hat effect handle, got %s", hatEffects[1].Handle)
	}
	if hatEffects[1].Name != "Unknown #2" {
		t.Errorf("Expected hat effect name 'Unknown #2', got %s", hatEffects[1].Name)
	}
}
