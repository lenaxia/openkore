// Package security provides security-related functionality for the network stack.
package security

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// Macro Reporter Status constants
const (
	MCRMonitoring = iota
	MCRNoData
	MCRInProgress
)

// Macro Detector Status constants
const (
	MCDTimeout = iota
	MCDIncorrect
	MCDGood
)

// MacroDetectManager manages macro detection-related functionality
type MacroDetectManager struct {
	parser           *core.CoreParser
	hookManager      *hooks.HookManager
	mutex            sync.RWMutex
	logFolder        string
	captchaSize      uint32
	captchaKey       string
	captchaImage     []byte
	playersList      interface{} // This would be a reference to the players list in a real implementation
	messageSender    interface{} // This would be a reference to the message sender in a real implementation
	networkInterface interface{} // This would be a reference to the network interface in a real implementation
}

// NewMacroDetectManager creates a new macro detection manager
func NewMacroDetectManager(parser *core.CoreParser, hookManager *hooks.HookManager) *MacroDetectManager {
	return &MacroDetectManager{
		parser:      parser,
		hookManager: hookManager,
		logFolder:   "logs", // Default log folder
	}
}

// RegisterHandlers registers macro detection-related packet handlers
func (m *MacroDetectManager) RegisterHandlers() {
	// Register handlers for macro detection-related packets
	m.parser.RegisterHandlerFunc("0A57", "macro_reporter_status", "B",
		[]string{"status"},
		m.handleMacroReporterStatus)

	m.parser.RegisterHandlerFunc("0A58", "macro_detector", "V v",
		[]string{"image_size", "captcha_key"},
		m.handleMacroDetector)

	m.parser.RegisterHandlerFunc("0A59", "macro_detector_image", "a*",
		[]string{"captcha_image"},
		m.handleMacroDetectorImage)

	m.parser.RegisterHandlerFunc("0A5B", "macro_detector_show", "B L",
		[]string{"remaining_chances", "remaining_time"},
		m.handleMacroDetectorShow)

	m.parser.RegisterHandlerFunc("0A5D", "macro_detector_status", "B",
		[]string{"status"},
		m.handleMacroDetectorStatus)

	m.parser.RegisterHandlerFunc("0A68", "macro_reporter_select", "a*",
		[]string{"account_list"},
		m.handleMacroReporterSelect)
}

// SetLogFolder sets the log folder for saving captcha images
func (m *MacroDetectManager) SetLogFolder(folder string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.logFolder = folder
}

// SetPlayersList sets the players list reference
func (m *MacroDetectManager) SetPlayersList(playersList interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.playersList = playersList
}

// SetMessageSender sets the message sender reference
func (m *MacroDetectManager) SetMessageSender(messageSender interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.messageSender = messageSender
}

// SetNetworkInterface sets the network interface reference
func (m *MacroDetectManager) SetNetworkInterface(networkInterface interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.networkInterface = networkInterface
}

// handleMacroReporterStatus handles the macro_reporter_status packet
// Packet format: 0A57 <status>.B
func (m *MacroDetectManager) handleMacroReporterStatus(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract status
	var status byte
	if statusVal, ok := args["status"].(byte); ok {
		status = statusVal
	}

	// Determine status text
	var statusText string
	switch status {
	case MCRMonitoring:
		statusText = "Monitoring"
	case MCRNoData:
		statusText = "No Data"
	case MCRInProgress:
		statusText = "In Progress"
	default:
		statusText = "Unknown"
	}

	// Format message
	message := fmt.Sprintf("Macro Reporter - Status: %s", statusText)

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/macro_reporter_status", map[string]interface{}{
			"status":      status,
			"status_text": statusText,
			"message":     message,
		})
	}

	return nil
}

// handleMacroDetector handles the macro_detector packet
// Packet format: 0A58 <image_size>.V <captcha_key>.v
func (m *MacroDetectManager) handleMacroDetector(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract image size and captcha key
	if imageSizeVal, ok := args["image_size"].(uint32); ok {
		m.captchaSize = imageSizeVal
	}

	if captchaKeyVal, ok := args["captcha_key"].(string); ok {
		m.captchaKey = captchaKeyVal
	}

	// Format debug message
	message := fmt.Sprintf("Macro Detector - image_size: %d bytes - captcha_key: %s", m.captchaSize, m.captchaKey)

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/macro_detector", map[string]interface{}{
			"image_size":  m.captchaSize,
			"captcha_key": m.captchaKey,
			"message":     message,
		})
	}

	return nil
}

// handleMacroDetectorImage handles the macro_detector_image packet
// Packet format: 0A59 <captcha_image>.a*
func (m *MacroDetectManager) handleMacroDetectorImage(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract captcha image data
	var captchaImageData []byte
	if captchaImageVal, ok := args["captcha_image"].([]byte); ok {
		captchaImageData = captchaImageVal
	}

	// Append to existing image data
	m.captchaImage = append(m.captchaImage, captchaImageData...)

	// Check if we have received the complete image
	if uint32(len(m.captchaImage)) >= m.captchaSize {
		// Decompress the image
		var uncompressedImage []byte
		var err error

		reader, err := zlib.NewReader(bytes.NewReader(m.captchaImage))
		if err != nil {
			if m.hookManager != nil {
				m.hookManager.CallHook("security/macro_detector_error", map[string]interface{}{
					"message": fmt.Sprintf("Failed to decompress macro detector image: %v", err),
				})
			}
			return nil
		}
		defer reader.Close()

		uncompressedImage, err = io.ReadAll(reader)
		if err != nil {
			if m.hookManager != nil {
				m.hookManager.CallHook("security/macro_detector_error", map[string]interface{}{
					"message": fmt.Sprintf("Failed to read decompressed macro detector image: %v", err),
				})
			}
			return nil
		}

		// Process the image (similar to the Perl implementation)
		imageHex := hex.EncodeToString(uncompressedImage)
		for i := 102; i < 3564 && i+6 <= len(imageHex); i += 6 {
			byte1, err1 := hex.DecodeString(imageHex[i : i+2])
			byte2 := imageHex[i+2 : i+4]
			byte3, err3 := hex.DecodeString(imageHex[i+4 : i+6])

			if err1 == nil && err3 == nil && len(byte1) > 0 && len(byte3) > 0 {
				if byte1[0] > 250 && byte2 == "00" && byte3[0] > 250 {
					// Replace the middle byte with "FF"
					imageHex = imageHex[:i+2] + "FF" + imageHex[i+4:]
				}
			}
		}

		// Convert back to binary
		finalImage, err := hex.DecodeString(imageHex)
		if err != nil {
			if m.hookManager != nil {
				m.hookManager.CallHook("security/macro_detector_error", map[string]interface{}{
					"message": fmt.Sprintf("Failed to decode processed image: %v", err),
				})
			}
			return nil
		}

		// Create filename
		filename := fmt.Sprintf("captcha_%s.bmp", m.captchaKey)

		// Ensure log folder exists
		if err := os.MkdirAll(m.logFolder, 0755); err != nil {
			if m.hookManager != nil {
				m.hookManager.CallHook("security/macro_detector_error", map[string]interface{}{
					"message": fmt.Sprintf("Failed to create log folder: %v", err),
				})
			}
			return nil
		}

		// Save to file
		filePath := filepath.Join(m.logFolder, filename)
		if err := os.WriteFile(filePath, finalImage, 0644); err != nil {
			if m.hookManager != nil {
				m.hookManager.CallHook("security/macro_detector_error", map[string]interface{}{
					"message": fmt.Sprintf("Failed to save macro detector image: %v", err),
				})
			}
			return nil
		}

		// Call hook for captcha image
		var hookReturn bool
		if m.hookManager != nil {
			// Create a channel to capture hook return values
			returnChan := make(chan bool, 1)

			// Add a temporary hook to capture the return value
			tempHook := m.hookManager.AddHook("captcha_image", func(hookName string, arg interface{}, userData interface{}) {
				if hookArg, ok := arg.(map[string]interface{}); ok {
					if returnVal, ok := hookArg["return"].(bool); ok && returnVal {
						returnChan <- true
					}
				}
			}, nil)

			// Call the hook
			m.hookManager.CallHook("captcha_image", map[string]interface{}{
				"captcha_image": finalImage,
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
		if !hookReturn {
			// Call hook for captcha file
			if m.hookManager != nil {
				// Create a channel to capture hook return values
				returnChan := make(chan bool, 1)

				// Add a temporary hook to capture the return value
				tempHook := m.hookManager.AddHook("captcha_file", func(hookName string, arg interface{}, userData interface{}) {
					if hookArg, ok := arg.(map[string]interface{}); ok {
						if returnVal, ok := hookArg["return"].(bool); ok && returnVal {
							returnChan <- true
						}
					}
				}, nil)

				// Call the hook
				m.hookManager.CallHook("captcha_file", map[string]interface{}{
					"file": filePath,
				})

				// Remove the temporary hook
				m.hookManager.DelHook(tempHook)

				// Check if a plugin handled the captcha file
				select {
				case <-returnChan:
					// A plugin handled the captcha file
					hookReturn = true
				default:
					// No plugin handled the captcha file
				}
			}
		}

		// If no plugin handled the captcha, show a warning
		if !hookReturn {
			if m.hookManager != nil {
				m.hookManager.CallHook("security/macro_detector_warning", map[string]interface{}{
					"message": fmt.Sprintf("Macro Detector - captcha has been saved in: %s, open it, solve it and use the command: captcha <text>", filePath),
				})
			}
		}

		// Reset state
		m.captchaImage = nil
		m.captchaSize = 0
		m.captchaKey = ""

		// Send download confirmation if we have a message sender
		// This would be implemented in a real system
	}

	return nil
}

// handleMacroDetectorShow handles the macro_detector_show packet
// Packet format: 0A5B <remaining_chances>.B <remaining_time>.L
func (m *MacroDetectManager) handleMacroDetectorShow(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract remaining chances and time
	var remainingChances byte
	var remainingTime uint32

	if remainingChancesVal, ok := args["remaining_chances"].(byte); ok {
		remainingChances = remainingChancesVal
	}

	if remainingTimeVal, ok := args["remaining_time"].(uint32); ok {
		remainingTime = remainingTimeVal
	}

	// Format messages
	message1 := "Macro Detector"
	message2 := fmt.Sprintf("Remaining Chances: %d - Remaining Time: %d seconds", remainingChances, remainingTime/1000)

	// Call hooks
	if m.hookManager != nil {
		m.hookManager.CallHook("security/macro_detector_show", map[string]interface{}{
			"remaining_chances": remainingChances,
			"remaining_time":    remainingTime,
			"message1":          message1,
			"message2":          message2,
		})
	}

	return nil
}

// handleMacroDetectorStatus handles the macro_detector_status packet
// Packet format: 0A5D <status>.B
func (m *MacroDetectManager) handleMacroDetectorStatus(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract status
	var status byte
	if statusVal, ok := args["status"].(byte); ok {
		status = statusVal
	}

	// Determine status text
	var statusText string
	switch status {
	case MCDTimeout:
		statusText = "Timeout"
	case MCDIncorrect:
		statusText = "Incorrect"
	case MCDGood:
		statusText = "Correct"
	default:
		statusText = "Unknown"
	}

	// Format message
	message := fmt.Sprintf("Macro Detector Status: %s", statusText)

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/macro_detector_status", map[string]interface{}{
			"status":      status,
			"status_text": statusText,
			"message":     message,
		})
	}

	return nil
}

// handleMacroReporterSelect handles the macro_reporter_select packet
// Packet format: 0A68 <account_list>.a*
func (m *MacroDetectManager) handleMacroReporterSelect(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract account list
	var accountList []byte
	if accountListVal, ok := args["account_list"].([]byte); ok {
		accountList = accountListVal
	}

	// Format message
	message := "Macro Reporter - Account List:"

	// Process account list
	accountIDs := make([]string, 0)
	for i := 0; i < len(accountList); i += 4 {
		if i+4 <= len(accountList) {
			accountID := accountList[i : i+4]
			accountIDs = append(accountIDs, fmt.Sprintf("%X", accountID))
		}
	}

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/macro_reporter_select", map[string]interface{}{
			"account_list": accountIDs,
			"message":      message,
		})
	}

	return nil
}
