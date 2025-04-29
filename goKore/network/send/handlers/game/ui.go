// Package game provides game-related packet handlers.
package game

// RegisterUIHandlers registers all UI-related packet handlers with the factory.
func RegisterUIHandlers(send HandlerRegistrar) {
	// UI handlers are registered through the UI manager
	// No need to register individual handlers here as they're handled by the UIManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a UI manager and use it to register handlers
}
