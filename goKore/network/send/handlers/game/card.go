// Package game provides game-related packet handlers.
package game

// RegisterCardHandlers registers all card-related packet handlers with the factory.
func RegisterCardHandlers(send HandlerRegistrar) {
	// Card handlers are registered through the card manager
	// No need to register individual handlers here as they're handled by the CardManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a card manager and use it to register handlers
}
