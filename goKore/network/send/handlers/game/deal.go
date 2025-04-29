// Package game provides game-related packet handlers.
package game

// RegisterDealHandlers registers all deal-related packet handlers with the factory.
func RegisterDealHandlers(send HandlerRegistrar) {
	// Deal handlers are registered through the deal manager
	// No need to register individual handlers here as they're handled by the DealManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a deal manager and use it to register handlers
}
