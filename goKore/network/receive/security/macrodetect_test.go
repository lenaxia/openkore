package security

import (
	"bytes"
	"compress/zlib"
	"os"
	"path/filepath"
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestNewMacroDetectManager(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewMacroDetectManager(parser, hookManager)

	if manager == nil {
		t.Fatal("NewMacroDetectManager() returned nil")
	}

	if manager.parser != parser {
		t.Error("manager.parser was not set correctly")
	}

	if manager.hookManager != hookManager {
		t.Error("manager.hookManager was not set correctly")
	}

	if manager.logFolder != "logs" {
		t.Errorf("manager.logFolder = %s, want logs", manager.logFolder)
	}
}

func TestMacroDetectRegisterHandlers(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewMacroDetectManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Verify handlers were registered
	handlerNames := []string{
		"macro_reporter_status",
		"macro_detector",
		"macro_detector_image",
		"macro_detector_show",
		"macro_detector_status",
		"macro_reporter_select",
	}

	for _, name := range handlerNames {
		if _, exists := parser.GetHandler(name); !exists {
			t.Errorf("Handler %s was not registered", name)
		}
	}
}

func TestMacroDetectSetLogFolder(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewMacroDetectManager(parser, nil)

	// Set log folder
	manager.SetLogFolder("test_logs")

	// Check log folder
	if manager.logFolder != "test_logs" {
		t.Errorf("manager.logFolder = %s, want test_logs", manager.logFolder)
	}
}

func TestHandleMacroReporterStatus(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewMacroDetectManager(parser, hookManager)

	// Create a channel to receive hook calls
	hookCalled := make(chan bool, 1)
	hookManager.RegisterHook("security/macro_reporter_status", func(args interface{}) {
		hookCalled <- true
	})

	// Test cases
	testCases := []struct {
		name       string
		status     byte
		statusText string
	}{
		{
			name:       "Monitoring",
			status:     MCRMonitoring,
			statusText: "Monitoring",
		},
		{
			name:       "No Data",
			status:     MCRNoData,
			statusText: "No Data",
		},
		{
			name:       "In Progress",
			status:     MCRInProgress,
			statusText: "In Progress",
		},
		{
			name:       "Unknown",
			status:     99,
			statusText: "Unknown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset hook channel
			select {
			case <-hookCalled:
				// Drain channel
			default:
				// Channel is already empty
			}

			// Call handler
			args := map[string]interface{}{
				"status": tc.status,
			}

			err := manager.handleMacroReporterStatus(args)
			if err != nil {
				t.Fatalf("handleMacroReporterStatus() returned error: %v", err)
			}

			// Check if hook was called
			select {
			case <-hookCalled:
				// Hook was called, which is expected
			default:
				t.Error("Hook security/macro_reporter_status was not called")
			}
		})
	}
}

func TestHandleMacroDetector(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewMacroDetectManager(parser, hookManager)

	// Create a channel to receive hook calls
	hookCalled := make(chan bool, 1)
	hookManager.RegisterHook("security/macro_detector", func(args interface{}) {
		hookCalled <- true
	})

	// Call handler
	args := map[string]interface{}{
		"image_size":  uint32(1000),
		"captcha_key": "test_key",
	}

	err := manager.handleMacroDetector(args)
	if err != nil {
		t.Fatalf("handleMacroDetector() returned error: %v", err)
	}

	// Check if hook was called
	select {
	case <-hookCalled:
		// Hook was called, which is expected
	default:
		t.Error("Hook security/macro_detector was not called")
	}

	// Check values
	if manager.captchaSize != 1000 {
		t.Errorf("manager.captchaSize = %d, want 1000", manager.captchaSize)
	}

	if manager.captchaKey != "test_key" {
		t.Errorf("manager.captchaKey = %s, want test_key", manager.captchaKey)
	}
}

func TestHandleMacroDetectorShow(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewMacroDetectManager(parser, hookManager)

	// Create a channel to receive hook calls
	hookCalled := make(chan bool, 1)
	hookManager.RegisterHook("security/macro_detector_show", func(args interface{}) {
		hookCalled <- true
	})

	// Call handler
	args := map[string]interface{}{
		"remaining_chances": byte(3),
		"remaining_time":    uint32(30000),
	}

	err := manager.handleMacroDetectorShow(args)
	if err != nil {
		t.Fatalf("handleMacroDetectorShow() returned error: %v", err)
	}

	// Check if hook was called
	select {
	case <-hookCalled:
		// Hook was called, which is expected
	default:
		t.Error("Hook security/macro_detector_show was not called")
	}
}

func TestHandleMacroDetectorStatus(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewMacroDetectManager(parser, hookManager)

	// Create a channel to receive hook calls
	hookCalled := make(chan bool, 1)
	hookManager.RegisterHook("security/macro_detector_status", func(args interface{}) {
		hookCalled <- true
	})

	// Test cases
	testCases := []struct {
		name       string
		status     byte
		statusText string
	}{
		{
			name:       "Timeout",
			status:     MCDTimeout,
			statusText: "Timeout",
		},
		{
			name:       "Incorrect",
			status:     MCDIncorrect,
			statusText: "Incorrect",
		},
		{
			name:       "Correct",
			status:     MCDGood,
			statusText: "Correct",
		},
		{
			name:       "Unknown",
			status:     99,
			statusText: "Unknown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset hook channel
			select {
			case <-hookCalled:
				// Drain channel
			default:
				// Channel is already empty
			}

			// Call handler
			args := map[string]interface{}{
				"status": tc.status,
			}

			err := manager.handleMacroDetectorStatus(args)
			if err != nil {
				t.Fatalf("handleMacroDetectorStatus() returned error: %v", err)
			}

			// Check if hook was called
			select {
			case <-hookCalled:
				// Hook was called, which is expected
			default:
				t.Error("Hook security/macro_detector_status was not called")
			}
		})
	}
}

func TestHandleMacroReporterSelect(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewMacroDetectManager(parser, hookManager)

	// Create a channel to receive hook calls
	hookCalled := make(chan bool, 1)
	hookManager.RegisterHook("security/macro_reporter_select", func(args interface{}) {
		hookCalled <- true
	})

	// Call handler
	args := map[string]interface{}{
		"account_list": []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	err := manager.handleMacroReporterSelect(args)
	if err != nil {
		t.Fatalf("handleMacroReporterSelect() returned error: %v", err)
	}

	// Check if hook was called
	select {
	case <-hookCalled:
		// Hook was called, which is expected
	default:
		t.Error("Hook security/macro_reporter_select was not called")
	}
}

func TestHandleMacroDetectorImage(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewMacroDetectManager(parser, hookManager)

	// Create a channel to receive hook calls
	warningCalled := make(chan bool, 1)
	hookManager.RegisterHook("security/macro_detector_warning", func(args interface{}) {
		warningCalled <- true
	})

	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "macro_detector_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set log folder
	manager.logFolder = tempDir

	// Set up initial state
	manager.captchaKey = "test_key"
	manager.captchaSize = 10

	// Create test image data (compressed)
	var compressedData bytes.Buffer
	w := zlib.NewWriter(&compressedData)
	w.Write([]byte("test image"))
	w.Close()

	// Get the compressed data
	compressedBytes := compressedData.Bytes()

	// Set the image size to match the total size of our compressed data
	totalSize := len(compressedBytes)
	manager.captchaSize = uint32(totalSize)

	// Call handler
	args := map[string]interface{}{
		"captcha_image": compressedBytes,
	}

	err = manager.handleMacroDetectorImage(args)
	if err != nil {
		t.Fatalf("handleMacroDetectorImage() returned error: %v", err)
	}

	// Check if warning hook was called
	select {
	case <-warningCalled:
		// Hook was called, which is expected
	default:
		t.Error("Hook security/macro_detector_warning was not called")
	}

	// Check file was created
	expectedFilePath := filepath.Join(tempDir, "captcha_test_key.bmp")
	_, err = os.Stat(expectedFilePath)
	if err != nil {
		t.Errorf("File %s was not created: %v", expectedFilePath, err)
	}

	// Check state was reset
	if len(manager.captchaImage) != 0 {
		t.Errorf("manager.captchaImage was not reset")
	}

	if manager.captchaSize != 0 {
		t.Errorf("manager.captchaSize = %d, want 0", manager.captchaSize)
	}

	if manager.captchaKey != "" {
		t.Errorf("manager.captchaKey = %s, want empty string", manager.captchaKey)
	}
}
