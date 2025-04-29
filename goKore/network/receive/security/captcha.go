// Package security provides security-related functionality for the network stack.
package security

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// CaptchaState represents the state of the captcha
type CaptchaState int

// Captcha states
const (
	CaptchaStateUnknown CaptchaState = iota
	CaptchaStateReady
	CaptchaStateError
	CaptchaStateComplete
	CaptchaStateAnswerGood
	CaptchaStateAnswerBad
)

// String returns the string representation of the captcha state
func (s CaptchaState) String() string {
	switch s {
	case CaptchaStateUnknown:
		return "Unknown"
	case CaptchaStateReady:
		return "Ready"
	case CaptchaStateError:
		return "Error"
	case CaptchaStateComplete:
		return "Complete"
	case CaptchaStateAnswerGood:
		return "AnswerGood"
	case CaptchaStateAnswerBad:
		return "AnswerBad"
	default:
		return "Invalid"
	}
}

// CaptchaManager manages captcha-related functionality
type CaptchaManager struct {
	parser       *core.CoreParser
	hookManager  *hooks.HookManager
	state        CaptchaState
	mutex        sync.RWMutex
	imageSize    uint32
	captchaKey   string
	captchaImage []byte
	logFolder    string
	charName     string
	sessionID    string
}

// NewCaptchaManager creates a new captcha manager
func NewCaptchaManager(parser *core.CoreParser, hookManager *hooks.HookManager) *CaptchaManager {
	return &CaptchaManager{
		parser:      parser,
		hookManager: hookManager,
		state:       CaptchaStateUnknown,
		logFolder:   "logs", // Default log folder
	}
}

// RegisterHandlers registers captcha-related packet handlers
func (m *CaptchaManager) RegisterHandlers() {
	// Register handlers for captcha-related packets
	m.parser.RegisterHandlerFunc("0A6A", "captcha_preview", "B V v",
		[]string{"status", "image_size", "captcha_key"},
		m.handleCaptchaPreview)

	m.parser.RegisterHandlerFunc("0A6B", "captcha_preview_image", "a*",
		[]string{"captcha_image"},
		m.handleCaptchaPreviewImage)

	m.parser.RegisterHandlerFunc("07E7", "captcha_session_ID", "a24",
		[]string{"session_id"},
		m.handleCaptchaSessionID)

	m.parser.RegisterHandlerFunc("07E8", "captcha_image", "a*",
		[]string{"image"},
		m.handleCaptchaImage)

	m.parser.RegisterHandlerFunc("07E9", "captcha_answer", "B",
		[]string{"flag"},
		m.handleCaptchaAnswer)

	m.parser.RegisterHandlerFunc("0A54", "captcha_upload_request", "B",
		[]string{"status"},
		m.handleCaptchaUploadRequest)

	m.parser.RegisterHandlerFunc("0A55", "captcha_upload_request_status", "",
		[]string{},
		m.handleCaptchaUploadRequestStatus)
}

// SetLogFolder sets the log folder for saving captcha images
func (m *CaptchaManager) SetLogFolder(folder string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.logFolder = folder
}

// SetCharName sets the character name for captcha image filenames
func (m *CaptchaManager) SetCharName(name string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.charName = name
}

// GetState returns the current captcha state
func (m *CaptchaManager) GetState() CaptchaState {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.state
}

// GetImageSize returns the captcha image size
func (m *CaptchaManager) GetImageSize() uint32 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.imageSize
}

// GetCaptchaKey returns the captcha key
func (m *CaptchaManager) GetCaptchaKey() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.captchaKey
}

// GetSessionID returns the captcha session ID
func (m *CaptchaManager) GetSessionID() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.sessionID
}

// handleCaptchaPreview handles the captcha_preview packet
func (m *CaptchaManager) handleCaptchaPreview(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract status, image size, and captcha key
	var status byte
	if statusVal, ok := args["status"].(byte); ok {
		status = statusVal
	}

	if imageSizeVal, ok := args["image_size"].(uint32); ok {
		m.imageSize = imageSizeVal
	}

	if captchaKeyVal, ok := args["captcha_key"].(string); ok {
		m.captchaKey = captchaKeyVal
	}

	// Update state based on status
	switch status {
	case 0:
		// Captcha preview is ready
		m.state = CaptchaStateReady
		if m.hookManager != nil {
			m.hookManager.CallHook("security/captcha_preview_ready", map[string]interface{}{
				"message": "Captcha Preview - Now you can download the image",
			})
		}
	case 1:
		// Captcha preview failed
		m.state = CaptchaStateError
		if m.hookManager != nil {
			m.hookManager.CallHook("security/captcha_preview_error", map[string]interface{}{
				"message": "Captcha Preview - Failed to Request Captcha (ID is out of range)",
			})
		}
	default:
		// Unknown status
		m.state = CaptchaStateError
		if m.hookManager != nil {
			m.hookManager.CallHook("security/captcha_preview_error", map[string]interface{}{
				"message": fmt.Sprintf("Captcha Preview - Unknown status: %d", status),
			})
		}
	}

	// Call hook with all information
	if m.hookManager != nil {
		m.hookManager.CallHook("security/captcha_preview", map[string]interface{}{
			"status":      status,
			"image_size":  m.imageSize,
			"captcha_key": m.captchaKey,
			"state":       m.state,
		})
	}

	return nil
}

// handleCaptchaPreviewImage handles the captcha_preview_image packet
func (m *CaptchaManager) handleCaptchaPreviewImage(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if we're in the right state
	if m.state != CaptchaStateReady {
		if m.hookManager != nil {
			m.hookManager.CallHook("security/captcha_error", map[string]interface{}{
				"message": "Received captcha image data without preview request",
			})
		}
		return nil
	}

	// Extract captcha image data
	var captchaImageData []byte
	if captchaImageVal, ok := args["captcha_image"].([]byte); ok {
		captchaImageData = captchaImageVal
	}

	// Append to existing image data
	m.captchaImage = append(m.captchaImage, captchaImageData...)

	// Check if we have received the complete image
	if uint32(len(m.captchaImage)) >= m.imageSize {
		// Decompress the image
		var uncompressedImage []byte
		var err error

		reader, err := zlib.NewReader(bytes.NewReader(m.captchaImage))
		if err != nil {
			if m.hookManager != nil {
				m.hookManager.CallHook("security/captcha_error", map[string]interface{}{
					"message": fmt.Sprintf("Failed to decompress captcha image: %v", err),
				})
			}
			return nil
		}
		defer reader.Close()

		uncompressedImage, err = io.ReadAll(reader)
		if err != nil {
			if m.hookManager != nil {
				m.hookManager.CallHook("security/captcha_error", map[string]interface{}{
					"message": fmt.Sprintf("Failed to read decompressed captcha image: %v", err),
				})
			}
			return nil
		}

		// Create filename
		var filename string
		if m.charName != "" {
			filename = fmt.Sprintf("captcha_preview_%s_%s.bmp", m.charName, m.captchaKey)
		} else {
			filename = fmt.Sprintf("captcha_preview_%s.bmp", m.captchaKey)
		}

		// Ensure log folder exists
		if err := os.MkdirAll(m.logFolder, 0755); err != nil {
			if m.hookManager != nil {
				m.hookManager.CallHook("security/captcha_error", map[string]interface{}{
					"message": fmt.Sprintf("Failed to create log folder: %v", err),
				})
			}
			return nil
		}

		// Save to file
		filePath := filepath.Join(m.logFolder, filename)
		if err := os.WriteFile(filePath, uncompressedImage, 0644); err != nil {
			if m.hookManager != nil {
				m.hookManager.CallHook("security/captcha_error", map[string]interface{}{
					"message": fmt.Sprintf("Failed to save captcha image: %v", err),
				})
			}
			return nil
		}

		// Reset state
		m.state = CaptchaStateComplete
		m.captchaImage = nil
		m.imageSize = 0
		m.captchaKey = ""

		// Call hook
		if m.hookManager != nil {
			m.hookManager.CallHook("security/captcha_preview_image_complete", map[string]interface{}{
				"message":  fmt.Sprintf("Captcha Preview - captcha has been saved in: %s", filePath),
				"filePath": filePath,
			})
		}
	} else {
		// Call hook for partial image
		if m.hookManager != nil {
			m.hookManager.CallHook("security/captcha_preview_image_partial", map[string]interface{}{
				"received": len(m.captchaImage),
				"total":    m.imageSize,
			})
		}
	}

	return nil
}

// handleCaptchaSessionID handles the captcha_session_ID packet
func (m *CaptchaManager) handleCaptchaSessionID(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract session ID
	if sessionIDVal, ok := args["session_id"].([]byte); ok {
		m.sessionID = string(sessionIDVal)
	}

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/captcha_session_id", map[string]interface{}{
			"session_id": m.sessionID,
		})
	}

	return nil
}

// handleCaptchaImage handles the captcha_image packet
func (m *CaptchaManager) handleCaptchaImage(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract image data
	var imageData []byte
	if imageVal, ok := args["image"].([]byte); ok {
		imageData = imageVal
	}

	// Call hook for plugins
	var hookReturn bool
	if m.hookManager != nil {
		// Create a channel to capture hook return values
		returnChan := make(chan bool, 1)

		// Add a temporary hook to capture the return value
		tempHook := m.hookManager.AddHook("security/captcha_image", func(hookName string, arg interface{}, userData interface{}) {
			if hookArg, ok := arg.(map[string]interface{}); ok {
				if returnVal, ok := hookArg["return"].(bool); ok && returnVal {
					returnChan <- true
				}
			}
		}, nil)

		// Call the hook
		m.hookManager.CallHook("security/captcha_image", map[string]interface{}{
			"image": imageData,
		})

		// Remove the temporary hook
		m.hookManager.DelHook(tempHook)

		// Check if a plugin handled the captcha
		select {
		case hookReturn = <-returnChan:
			// A plugin handled the captcha
		default:
			// No plugin handled the captcha
		}
	}

	// If a plugin handled the captcha, we're done
	if hookReturn {
		return nil
	}

	// Create filename
	var filename string
	if m.charName != "" {
		filename = fmt.Sprintf("captcha_%s.bmp", m.charName)
	} else {
		filename = "captcha.bmp"
	}

	// Ensure log folder exists
	if err := os.MkdirAll(m.logFolder, 0755); err != nil {
		if m.hookManager != nil {
			m.hookManager.CallHook("security/captcha_error", map[string]interface{}{
				"message": fmt.Sprintf("Failed to create log folder: %v", err),
			})
		}
		return nil
	}

	// Save to file
	filePath := filepath.Join(m.logFolder, filename)
	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		if m.hookManager != nil {
			m.hookManager.CallHook("security/captcha_error", map[string]interface{}{
				"message": fmt.Sprintf("Failed to save captcha image: %v", err),
			})
		}
		return nil
	}

	// Call hook for file
	if m.hookManager != nil {
		// Create a channel to capture hook return values
		returnChan := make(chan bool, 1)

		// Add a temporary hook to capture the return value
		tempHook := m.hookManager.AddHook("security/captcha_file", func(hookName string, arg interface{}, userData interface{}) {
			if hookArg, ok := arg.(map[string]interface{}); ok {
				if returnVal, ok := hookArg["return"].(bool); ok && returnVal {
					returnChan <- true
				}
			}
		}, nil)

		// Call the hook
		m.hookManager.CallHook("security/captcha_file", map[string]interface{}{
			"file": filePath,
		})

		// Remove the temporary hook
		m.hookManager.DelHook(tempHook)

		// Check if a plugin handled the captcha file
		select {
		case <-returnChan:
			// A plugin handled the captcha file
			return nil
		default:
			// No plugin handled the captcha file
		}
	}

	// Log a warning message
	if m.hookManager != nil {
		m.hookManager.CallHook("security/captcha_warning", map[string]interface{}{
			"message": fmt.Sprintf("captcha.bmp has been saved to: %s, open it, solve it and use the command: captcha <text>", m.logFolder),
		})
	}

	return nil
}

// handleCaptchaUploadRequest handles the captcha_upload_request packet
// Packet format: 0A54 <status>.B
func (m *CaptchaManager) handleCaptchaUploadRequest(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract status
	var status byte
	if statusVal, ok := args["status"].(byte); ok {
		status = statusVal
	}

	var message string
	switch status {
	case 0:
		message = "Captcha Register - Now you can upload the image"
	case 1:
		message = "Captcha Register - Failed to upload the image"
	default:
		message = fmt.Sprintf("Captcha Register - Unknown status: %d", status)
	}

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/captcha_upload_request", map[string]interface{}{
			"status":  status,
			"message": message,
		})
	}

	return nil
}

// handleCaptchaUploadRequestStatus handles the captcha_upload_request_status packet
// Packet format: 0A55
func (m *CaptchaManager) handleCaptchaUploadRequestStatus(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	message := "Captcha Register - Image uploaded successfully"

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/captcha_upload_request_status", map[string]interface{}{
			"message": message,
		})
	}

	return nil
}

// handleCaptchaAnswer handles the captcha_answer packet
func (m *CaptchaManager) handleCaptchaAnswer(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract flag
	var flag byte
	if flagVal, ok := args["flag"].(byte); ok {
		flag = flagVal
	}

	// Update state based on flag
	if flag != 0 {
		m.state = CaptchaStateAnswerGood
	} else {
		m.state = CaptchaStateAnswerBad
	}

	// Call hook
	if m.hookManager != nil {
		message := "Captcha answer was incorrect"
		if flag != 0 {
			message = "Captcha answer was correct"
		}

		m.hookManager.CallHook("security/captcha_answer", map[string]interface{}{
			"flag":    flag,
			"success": flag != 0,
			"message": message,
		})
	}

	return nil
}
