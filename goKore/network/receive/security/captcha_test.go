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

func TestNewCaptchaManager(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewCaptchaManager(parser, hookManager)

	if manager == nil {
		t.Fatal("NewCaptchaManager() returned nil")
	}

	if manager.parser != parser {
		t.Error("manager.parser was not set correctly")
	}

	if manager.hookManager != hookManager {
		t.Error("manager.hookManager was not set correctly")
	}

	if manager.state != CaptchaStateUnknown {
		t.Errorf("manager.state = %v, want %v", manager.state, CaptchaStateUnknown)
	}

	if manager.logFolder != "logs" {
		t.Errorf("manager.logFolder = %s, want logs", manager.logFolder)
	}
}

func TestCaptchaRegisterHandlers(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewCaptchaManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Verify handlers were registered
	handlerNames := []string{
		"captcha_preview",
		"captcha_preview_image",
		"captcha_session_ID",
		"captcha_image",
		"captcha_answer",
		"captcha_upload_request",
		"captcha_upload_request_status",
	}

	for _, name := range handlerNames {
		if _, exists := parser.GetHandler(name); !exists {
			t.Errorf("Handler %s was not registered", name)
		}
	}
}

func TestSetGetCaptchaState(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewCaptchaManager(parser, nil)

	// Set state
	manager.state = CaptchaStateReady

	// Get state
	state := manager.GetState()
	if state != CaptchaStateReady {
		t.Errorf("GetState() = %v, want %v", state, CaptchaStateReady)
	}
}

func TestSetGetImageSize(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewCaptchaManager(parser, nil)

	// Set image size
	manager.imageSize = 1000

	// Get image size
	size := manager.GetImageSize()
	if size != 1000 {
		t.Errorf("GetImageSize() = %d, want 1000", size)
	}
}

func TestSetGetCaptchaKey(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewCaptchaManager(parser, nil)

	// Set captcha key
	manager.captchaKey = "test_key"

	// Get captcha key
	key := manager.GetCaptchaKey()
	if key != "test_key" {
		t.Errorf("GetCaptchaKey() = %s, want test_key", key)
	}
}

func TestSetGetSessionID(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewCaptchaManager(parser, nil)

	// Set session ID
	manager.sessionID = "test_session_id"

	// Get session ID
	sessionID := manager.GetSessionID()
	if sessionID != "test_session_id" {
		t.Errorf("GetSessionID() = %s, want test_session_id", sessionID)
	}
}

func TestSetLogFolder(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewCaptchaManager(parser, nil)

	// Set log folder
	manager.SetLogFolder("test_logs")

	// Check log folder
	if manager.logFolder != "test_logs" {
		t.Errorf("manager.logFolder = %s, want test_logs", manager.logFolder)
	}
}

func TestSetCharName(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewCaptchaManager(parser, nil)

	// Set character name
	manager.SetCharName("TestChar")

	// Check character name
	if manager.charName != "TestChar" {
		t.Errorf("manager.charName = %s, want TestChar", manager.charName)
	}
}

func TestHandleCaptchaPreview(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewCaptchaManager(parser, hookManager)

	// Create a channel to receive hook calls
	hookCalled := make(chan bool, 1)
	hookManager.RegisterHook("security/captcha_preview", func(args interface{}) {
		hookCalled <- true
	})

	// Test with status = 0 (success)
	args := map[string]interface{}{
		"status":      byte(0),
		"image_size":  uint32(1000),
		"captcha_key": "test_key",
	}

	err := manager.handleCaptchaPreview(args)
	if err != nil {
		t.Fatalf("handleCaptchaPreview() returned error: %v", err)
	}

	// Check if hook was called
	select {
	case <-hookCalled:
		// Hook was called, which is expected
	default:
		t.Error("Hook security/captcha_preview was not called")
	}

	// Check state
	if manager.state != CaptchaStateReady {
		t.Errorf("manager.state = %v, want %v", manager.state, CaptchaStateReady)
	}

	// Check image size
	if manager.imageSize != 1000 {
		t.Errorf("manager.imageSize = %d, want 1000", manager.imageSize)
	}

	// Check captcha key
	if manager.captchaKey != "test_key" {
		t.Errorf("manager.captchaKey = %s, want test_key", manager.captchaKey)
	}

	// Test with status = 1 (error)
	args = map[string]interface{}{
		"status":      byte(1),
		"image_size":  uint32(1000),
		"captcha_key": "test_key",
	}

	err = manager.handleCaptchaPreview(args)
	if err != nil {
		t.Fatalf("handleCaptchaPreview() returned error: %v", err)
	}

	// Check state
	if manager.state != CaptchaStateError {
		t.Errorf("manager.state = %v, want %v", manager.state, CaptchaStateError)
	}
}

func TestHandleCaptchaPreviewImage(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewCaptchaManager(parser, hookManager)

	// Create a channel to receive hook calls
	partialHookCalled := make(chan bool, 1)
	completeHookCalled := make(chan bool, 1)

	hookManager.RegisterHook("security/captcha_preview_image_partial", func(args interface{}) {
		partialHookCalled <- true
	})

	hookManager.RegisterHook("security/captcha_preview_image_complete", func(args interface{}) {
		completeHookCalled <- true
	})

	// Set up initial state
	manager.captchaKey = "test_key"
	manager.state = CaptchaStateReady

	// Create test image data (compressed)
	var compressedData bytes.Buffer
	w := zlib.NewWriter(&compressedData)
	w.Write([]byte("test image data"))
	w.Close()

	// Get the compressed data
	compressedBytes := compressedData.Bytes()

	// Calculate half of the data for partial sending
	halfSize := len(compressedBytes) / 2

	// Set the image size to match the total size of our compressed data
	totalSize := len(compressedBytes)
	manager.imageSize = uint32(totalSize)

	// Test with partial image
	args := map[string]interface{}{
		"captcha_image": compressedBytes[:halfSize], // Send first half
	}

	err := manager.handleCaptchaPreviewImage(args)
	if err != nil {
		t.Fatalf("handleCaptchaPreviewImage() returned error: %v", err)
	}

	// Check if hook was called
	select {
	case <-partialHookCalled:
		// Hook was called, which is expected
	default:
		t.Error("Hook security/captcha_preview_image_partial was not called")
	}

	// Check image data was appended
	if len(manager.captchaImage) != halfSize {
		t.Errorf("len(manager.captchaImage) = %d, want %d", len(manager.captchaImage), halfSize)
	}

	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "captcha_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set log folder
	manager.logFolder = tempDir

	// Test with complete image
	args = map[string]interface{}{
		"captcha_image": compressedBytes[halfSize:], // Send second half
	}

	err = manager.handleCaptchaPreviewImage(args)
	if err != nil {
		t.Fatalf("handleCaptchaPreviewImage() returned error: %v", err)
	}

	// Check if hook was called
	select {
	case <-completeHookCalled:
		// Hook was called, which is expected
	default:
		t.Error("Hook security/captcha_preview_image_complete was not called")
	}

	// Check file was created
	expectedFilePath := filepath.Join(tempDir, "captcha_preview_test_key.bmp")
	_, err = os.Stat(expectedFilePath)
	if err != nil {
		t.Errorf("File %s was not created: %v", expectedFilePath, err)
	}

	// Check file contents
	fileData, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(fileData) != "test image data" {
		t.Errorf("File contents = %s, want 'test image data'", string(fileData))
	}

	// Check state was reset
	if len(manager.captchaImage) != 0 {
		t.Errorf("manager.captchaImage was not reset")
	}

	if manager.imageSize != 0 {
		t.Errorf("manager.imageSize = %d, want 0", manager.imageSize)
	}

	if manager.captchaKey != "" {
		t.Errorf("manager.captchaKey = %s, want empty string", manager.captchaKey)
	}
}

func TestHandleCaptchaPreviewImageErrors(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewCaptchaManager(parser, hookManager)

	// Create a channel to receive hook calls
	errorHookCalled := make(chan bool, 1)
	hookManager.RegisterHook("security/captcha_error", func(args interface{}) {
		errorHookCalled <- true
	})

	// Test with wrong state
	manager.state = CaptchaStateUnknown
	args := map[string]interface{}{
		"captcha_image": []byte("test image data"),
	}

	err := manager.handleCaptchaPreviewImage(args)
	if err != nil {
		t.Fatalf("handleCaptchaPreviewImage() returned error: %v", err)
	}

	// Check if hook was called
	select {
	case <-errorHookCalled:
		// Hook was called, which is expected
	default:
		t.Error("Hook security/captcha_error was not called")
	}
}

func TestHandleCaptchaSessionID(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewCaptchaManager(parser, hookManager)

	// Create a channel to receive hook calls
	hookCalled := make(chan bool, 1)
	hookManager.RegisterHook("security/captcha_session_id", func(args interface{}) {
		hookCalled <- true
	})

	// Test with valid session ID
	args := map[string]interface{}{
		"session_id": []byte("test_session_id"),
	}

	err := manager.handleCaptchaSessionID(args)
	if err != nil {
		t.Fatalf("handleCaptchaSessionID() returned error: %v", err)
	}

	// Check if hook was called
	select {
	case <-hookCalled:
		// Hook was called, which is expected
	default:
		t.Error("Hook security/captcha_session_id was not called")
	}

	// Check session ID
	if manager.sessionID != "test_session_id" {
		t.Errorf("manager.sessionID = %s, want test_session_id", manager.sessionID)
	}
}

func TestHandleCaptchaImage(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewCaptchaManager(parser, hookManager)

	// Create channels to receive hook calls
	imageHookCalled := make(chan bool, 1)
	warningHookCalled := make(chan bool, 1)

	hookManager.RegisterHook("security/captcha_image", func(args interface{}) {
		imageHookCalled <- true
	})

	hookManager.RegisterHook("security/captcha_warning", func(args interface{}) {
		warningHookCalled <- true
	})

	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "captcha_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set log folder
	manager.logFolder = tempDir

	// Test with valid image
	args := map[string]interface{}{
		"image": []byte("test captcha image"),
	}

	err = manager.handleCaptchaImage(args)
	if err != nil {
		t.Fatalf("handleCaptchaImage() returned error: %v", err)
	}

	// Check if hooks were called
	select {
	case <-imageHookCalled:
		// Hook was called, which is expected
	default:
		t.Error("Hook security/captcha_image was not called")
	}

	select {
	case <-warningHookCalled:
		// Hook was called, which is expected
	default:
		t.Error("Hook security/captcha_warning was not called")
	}

	// Check file was created
	expectedFilePath := filepath.Join(tempDir, "captcha.bmp")
	_, err = os.Stat(expectedFilePath)
	if err != nil {
		t.Errorf("File %s was not created: %v", expectedFilePath, err)
	}

	// Check file contents
	fileData, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(fileData) != "test captcha image" {
		t.Errorf("File contents = %s, want 'test captcha image'", string(fileData))
	}
}

func TestHandleCaptchaAnswer(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewCaptchaManager(parser, hookManager)

	// Create a channel to receive hook calls
	hookCalled := make(chan bool, 1)
	hookManager.RegisterHook("security/captcha_answer", func(args interface{}) {
		hookCalled <- true
	})

	// Test cases
	testCases := []struct {
		name          string
		flag          byte
		expectedState CaptchaState
	}{
		{
			name:          "Correct Answer",
			flag:          1,
			expectedState: CaptchaStateAnswerGood,
		},
		{
			name:          "Incorrect Answer",
			flag:          0,
			expectedState: CaptchaStateAnswerBad,
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
				"flag": tc.flag,
			}

			err := manager.handleCaptchaAnswer(args)
			if err != nil {
				t.Fatalf("handleCaptchaAnswer() returned error: %v", err)
			}

			// Check if hook was called
			select {
			case <-hookCalled:
				// Hook was called, which is expected
			default:
				t.Error("Hook security/captcha_answer was not called")
			}

			// Check state
			if manager.state != tc.expectedState {
				t.Errorf("manager.state = %v, want %v", manager.state, tc.expectedState)
			}
		})
	}
}

func TestCaptchaStateString(t *testing.T) {
	tests := []struct {
		state CaptchaState
		want  string
	}{
		{CaptchaStateUnknown, "Unknown"},
		{CaptchaStateReady, "Ready"},
		{CaptchaStateError, "Error"},
		{CaptchaStateComplete, "Complete"},
		{CaptchaStateAnswerGood, "AnswerGood"},
		{CaptchaStateAnswerBad, "AnswerBad"},
		{CaptchaState(99), "Invalid"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Errorf("CaptchaState(%d).String() = %v, want %v", tt.state, got, tt.want)
		}
	}
}

func TestHandleCaptchaUploadRequest(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewCaptchaManager(parser, hookManager)

	// Create a channel to receive hook calls
	hookCalled := make(chan bool, 1)
	hookManager.RegisterHook("security/captcha_upload_request", func(args interface{}) {
		hookCalled <- true
	})

	// Test cases
	testCases := []struct {
		name    string
		status  byte
		message string
	}{
		{
			name:    "Upload Ready",
			status:  0,
			message: "Captcha Register - Now you can upload the image",
		},
		{
			name:    "Upload Failed",
			status:  1,
			message: "Captcha Register - Failed to upload the image",
		},
		{
			name:    "Unknown Status",
			status:  2,
			message: "Captcha Register - Unknown status: 2",
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

			err := manager.handleCaptchaUploadRequest(args)
			if err != nil {
				t.Fatalf("handleCaptchaUploadRequest() returned error: %v", err)
			}

			// Check if hook was called
			select {
			case <-hookCalled:
				// Hook was called, which is expected
			default:
				t.Error("Hook security/captcha_upload_request was not called")
			}
		})
	}
}

func TestHandleCaptchaUploadRequestStatus(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewCaptchaManager(parser, hookManager)

	// Create a channel to receive hook calls
	hookCalled := make(chan bool, 1)
	hookManager.RegisterHook("security/captcha_upload_request_status", func(args interface{}) {
		hookCalled <- true
	})

	// Call handler
	err := manager.handleCaptchaUploadRequestStatus(map[string]interface{}{})
	if err != nil {
		t.Fatalf("handleCaptchaUploadRequestStatus() returned error: %v", err)
	}

	// Check if hook was called
	select {
	case <-hookCalled:
		// Hook was called, which is expected
	default:
		t.Error("Hook security/captcha_upload_request_status was not called")
	}
}
