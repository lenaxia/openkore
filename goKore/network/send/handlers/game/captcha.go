// Package game provides game-related packet handlers.
package game

// RegisterCaptchaHandlers registers all captcha-related packet handlers with the factory.
func RegisterCaptchaHandlers(send HandlerRegistrar) {
	// Captcha handlers are registered through the captcha manager
	// No need to register individual handlers here as they're handled by the CaptchaManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a captcha manager and use it to register handlers
}
