package security

import (
	"github.com/lenaxia/goKore/network/hooks"
)

// LoginTokenManager manages login token-related functionality
type LoginTokenManager struct {
	hookManager *hooks.HookManager
	netVersion  int // Network version
}

// NewLoginTokenManager creates a new login token manager
func NewLoginTokenManager(hookManager *hooks.HookManager, netVersion int) *LoginTokenManager {
	return &LoginTokenManager{
		hookManager: hookManager,
		netVersion:  netVersion,
	}
}

// HandleReceivedLoginToken handles the received_login_token packet
// Packet format: 0AE3 <login_token>.32B <len>.W <OTP_ip>.32B <OTP_port>.W
func (m *LoginTokenManager) HandleReceivedLoginToken(args map[string]interface{}) error {
	// XKore mode 1 / 3
	if m.netVersion == 1 {
		return nil
	}

	// Extract login token with safety check
	loginToken, ok := args["login_token"].([]byte)
	if !ok {
		return nil // Silently ignore if login_token is missing or invalid
	}

	// Extract length with safety check
	length, ok := args["len"].(uint16)
	if !ok {
		return nil // Silently ignore if len is missing or invalid
	}

	// Extract OTP IP with safety check
	otpIP, ok := args["OTP_ip"].([]byte)
	if !ok {
		return nil // Silently ignore if OTP_ip is missing or invalid
	}

	// Extract OTP port with safety check
	otpPort, ok := args["OTP_port"].(uint16)
	if !ok {
		return nil // Silently ignore if OTP_port is missing or invalid
	}

	// In the original implementation, this would call sendTokenToServer
	// We'll use the hook system to handle this
	if m.hookManager != nil {
		m.hookManager.CallHook("security.received_login_token", map[string]interface{}{
			"login_token": loginToken,
			"len":         length,
			"OTP_ip":      otpIP,
			"OTP_port":    otpPort,
		})
	}

	return nil
}

// SetNetVersion sets the network version
func (m *LoginTokenManager) SetNetVersion(version int) {
	m.netVersion = version
}
