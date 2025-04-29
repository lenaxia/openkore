// Package game provides game-related packet handlers.
package game

// RegisterMiscHandlers registers all miscellaneous packet handlers with the factory.
func RegisterMiscHandlers(send HandlerRegistrar) {
	// Miscellaneous handlers are registered through the misc manager
	// No need to register individual handlers here as they're handled by the MiscManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a misc manager and use it to register handlers
}
