// Package captcha provides captcha-related packet sending functionality.
package captcha

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// CaptchaManager handles captcha-related packet sending.
type CaptchaManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewCaptchaManager creates a new captcha manager.
func NewCaptchaManager(baseSend core.Send) *CaptchaManager {
	return &CaptchaManager{
		baseSend: baseSend,
	}
}

// SendCaptchaAnswer sends an answer to a captcha.
// This is equivalent to the sendCaptchaAnswer function in Send.pm.
func (cm *CaptchaManager) SendCaptchaAnswer(answer string) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("captcha_answer")
	if !exists {
		return fmt.Errorf("captcha_answer packet ID not found")
	}

	// In the real implementation, we would get the account ID from the session
	// For now, we'll use a hardcoded value
	accountID := uint32(12345)

	// Create the arguments
	args := map[string]interface{}{
		"accountID": accountID,
		"answer":    answer,
		"len":       32, // Fixed length as mentioned in the original code
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendCaptchaPreviewRequest sends a request to preview a captcha.
// This is equivalent to the sendCaptchaPreviewRequest function in Send.pm.
func (cm *CaptchaManager) SendCaptchaPreviewRequest(captchaKey uint32) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("captcha_preview_request")
	if !exists {
		return fmt.Errorf("captcha_preview_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"captcha_key": captchaKey,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}
