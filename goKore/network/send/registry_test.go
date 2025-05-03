package send

import (
	"testing"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
)

// MockLogger is a mock implementation of the core.Logger interface for testing
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

func (ml *MockLogger) Debug(format string, args ...interface{}) {
	ml.debugMessages = append(ml.debugMessages, format)
}

func (ml *MockLogger) Info(format string, args ...interface{}) {
	ml.infoMessages = append(ml.infoMessages, format)
}

func (ml *MockLogger) Warning(format string, args ...interface{}) {
	ml.warningMessages = append(ml.warningMessages, format)
}

func (ml *MockLogger) Error(format string, args ...interface{}) {
	ml.errorMessages = append(ml.errorMessages, format)
}

func (ml *MockLogger) Success(format string, args ...interface{}) {
	ml.successMessages = append(ml.successMessages, format)
}

// NewMockBaseSend creates a new mock BaseSend for testing
func NewMockBaseSend() *core.BaseSend {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a base send
	baseSend := core.NewBaseSend(hookManager)

	// Set the server type
	packetConstructions := make(map[string]common.PacketConstruction)
	baseSend.Configure("ServerType0", packetConstructions)

	return baseSend
}

// TestNewHandlerRegistry tests the NewHandlerRegistry function
func TestNewHandlerRegistry(t *testing.T) {
	// Create dependencies
	baseSend := NewMockBaseSend()
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()

	// Create the registry
	registry := NewHandlerRegistry(baseSend, hookManager, logger)

	// Verify the registry was created correctly
	if registry == nil {
		t.Fatal("NewHandlerRegistry returned nil")
	}

	if registry.baseSend != baseSend {
		t.Error("registry.baseSend was not set correctly")
	}

	if registry.hookManager != hookManager {
		t.Error("registry.hookManager was not set correctly")
	}

	if registry.logger != logger {
		t.Error("registry.logger was not set correctly")
	}

	if registry.managers == nil {
		t.Error("registry.managers was not initialized")
	}
}

// TestRegisterAllHandlers tests the RegisterAllHandlers function
func TestRegisterAllHandlers(t *testing.T) {
	// Create dependencies
	baseSend := NewMockBaseSend()
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()

	// Create the registry
	registry := NewHandlerRegistry(baseSend, hookManager, logger)

	// Register all handlers
	registry.RegisterAllHandlers()

	// Verify that the info message was logged
	if len(logger.infoMessages) == 0 || logger.infoMessages[0] != "Registered all send handlers" {
		t.Error("Expected 'Registered all send handlers' info message")
	}
}

// TestGetManager tests the GetManager function
func TestGetManager(t *testing.T) {
	// Create dependencies
	baseSend := NewMockBaseSend()
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()

	// Create the registry
	registry := NewHandlerRegistry(baseSend, hookManager, logger)

	// Add a test manager
	testManager := "test manager"
	registry.managers["test"] = testManager

	// Test getting an existing manager
	manager, exists := registry.GetManager("test")
	if !exists {
		t.Error("GetManager returned exists=false for an existing manager")
	}

	if manager != testManager {
		t.Error("GetManager returned the wrong manager")
	}

	// Test getting a non-existent manager
	_, exists = registry.GetManager("non-existent")
	if exists {
		t.Error("GetManager returned exists=true for a non-existent manager")
	}
}

// TestConfigureServerType tests the ConfigureServerType function
func TestConfigureServerType(t *testing.T) {
	// Create dependencies
	baseSend := NewMockBaseSend()
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()

	// Create the registry
	registry := NewHandlerRegistry(baseSend, hookManager, logger)

	// Configure the server type
	packetDefs := make(map[string]common.PacketConstruction)
	err := registry.ConfigureServerType("TestServerType", packetDefs)

	// Verify that no error was returned
	if err != nil {
		t.Errorf("ConfigureServerType returned an error: %v", err)
	}

	// Verify that the server type was set correctly
	if baseSend.GetServerType() != "TestServerType" {
		t.Errorf("Expected server type to be 'TestServerType', got '%s'", baseSend.GetServerType())
	}
}

// TestGetPacketDefinitions tests the GetPacketDefinitions function
func TestGetPacketDefinitions(t *testing.T) {
	// Create dependencies
	baseSend := NewMockBaseSend()
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()

	// Create the registry
	registry := NewHandlerRegistry(baseSend, hookManager, logger)

	// Get packet definitions
	packetDefs := registry.GetPacketDefinitions()

	// Verify that a map was returned
	if packetDefs == nil {
		t.Error("GetPacketDefinitions returned nil")
	}
}
