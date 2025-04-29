// Package game provides game-related packet handlers.
package game

// RegisterMercenaryHandlers registers all mercenary-related packet handlers with the factory.
func RegisterMercenaryHandlers(send HandlerRegistrar) {
	// Mercenary handlers are registered through the mercenary manager
	// No need to register individual handlers here as they're handled by the MercenaryManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a mercenary manager and use it to register handlers
}
