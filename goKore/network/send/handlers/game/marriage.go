// Package game provides game-related packet handlers.
package game

// RegisterMarriageHandlers registers all marriage-related packet handlers with the factory.
func RegisterMarriageHandlers(send HandlerRegistrar) {
	// Marriage handlers are registered through the marriage manager
	// No need to register individual handlers here as they're handled by the MarriageManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a marriage manager and use it to register handlers
}
