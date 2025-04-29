package security

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/factory"
	"github.com/lenaxia/goKore/network/receive/types"
)

// RegisterAllHandlers registers all handlers in the security package with the given receive component
func RegisterAllHandlers(receive types.Receive) {
	// Register the captcha_preview handler
	receive.RegisterHandler("captcha_preview", func(args map[string]interface{}) error {
		// Create a captcha manager for this specific call
		manager := NewCaptchaManager(nil, nil)
		return manager.handleCaptchaPreview(args)
	})

	// Register the captcha_preview_image handler
	receive.RegisterHandler("captcha_preview_image", func(args map[string]interface{}) error {
		// Create a captcha manager for this specific call
		manager := NewCaptchaManager(nil, nil)
		return manager.handleCaptchaPreviewImage(args)
	})

	// Register PIN handlers with default configuration
	RegisterPINHandlers(receive)

	// Register login error handlers
	RegisterLoginErrorHandlers(receive)

	// Register login token handlers
	RegisterLoginTokenHandlers(receive)
}

// RegisterWithFactory registers all handlers in the security package with the given factory
func RegisterWithFactory(receiveFactory *factory.ReceiveFactory) {
	// Currently, the factory doesn't have a direct method to register packet handlers
	// This would typically be done through server-specific packet definitions
}

// RegisterWithCoreParser registers all handlers in the security package with the given core parser
func RegisterWithCoreParser(parser *core.CoreParser, hookManager *hooks.HookManager) {
	// Create the security managers
	captchaManager := NewCaptchaManager(parser, hookManager)
	pinManager := NewPINManager(parser, hookManager)
	loginTokenManager := NewLoginTokenManager(hookManager, 0) // Default to network version 0

	// Register handlers
	captchaManager.RegisterHandlers()
	pinManager.RegisterHandlers()

	// Register login token handler
	parser.RegisterHandlerFunc("0AE3", "received_login_token", "a32 W a32 W",
		[]string{"login_token", "len", "OTP_ip", "OTP_port"},
		loginTokenManager.HandleReceivedLoginToken)
}

// GetPacketDefinitions returns the packet definitions for the security package
func GetPacketDefinitions() map[string]common.PacketDef {
	return map[string]common.PacketDef{
		"0A6A": {
			ID:         "0A6A",
			Name:       "captcha_preview",
			Format:     "B V v",
			FieldNames: []string{"status", "image_size", "captcha_key"},
		},
		"0A6B": {
			ID:         "0A6B",
			Name:       "captcha_preview_image",
			Format:     "a*",
			FieldNames: []string{"captcha_image"},
		},
		"0AE3": {
			ID:         "0AE3",
			Name:       "received_login_token",
			Format:     "a32 W a32 W",
			FieldNames: []string{"login_token", "len", "OTP_ip", "OTP_port"},
		},
		// Add other security-related packet definitions here
	}
}

// RegisterPINHandlers registers all PIN-related handlers with the given receive component
func RegisterPINHandlers(receive types.Receive) {
	// Register the login_pin_code_request handler
	receive.RegisterHandler("login_pin_code_request", func(args map[string]interface{}) error {
		// Create a PIN manager for this specific call
		manager := NewPINManager(nil, nil)
		return manager.handleLoginPinCodeRequest(args)
	})

	// Register the login_pin_new_code_result handler
	receive.RegisterHandler("login_pin_new_code_result", func(args map[string]interface{}) error {
		// Create a PIN manager for this specific call
		manager := NewPINManager(nil, nil)
		return manager.handleLoginPinNewCodeResult(args)
	})
}

// RegisterLoginErrorHandlers registers all login error handlers with the given receive component
func RegisterLoginErrorHandlers(receive types.Receive) {
	// Register login error handlers here
	// This is a placeholder for future implementation
}

// RegisterLoginTokenHandlers registers all login token handlers with the given receive component
func RegisterLoginTokenHandlers(receive types.Receive) {
	// Register the received_login_token handler
	receive.RegisterHandler("received_login_token", func(args map[string]interface{}) error {
		// Create a login token manager for this specific call
		manager := NewLoginTokenManager(nil, 0) // Default to network version 0
		return manager.HandleReceivedLoginToken(args)
	})
}
