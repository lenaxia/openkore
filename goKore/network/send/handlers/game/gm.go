// Package game provides game-related packet handlers.
package game

// RegisterGMHandlers registers all GM-related packet handlers with the factory.
func RegisterGMHandlers(send HandlerRegistrar) {
	// GM handlers are registered through the GM manager
	// No need to register individual handlers here as they're handled by the GMManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a GM manager and use it to register handlers
}
