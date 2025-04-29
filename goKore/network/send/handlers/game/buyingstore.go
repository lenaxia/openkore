// Package game provides game-related packet handlers.
package game

// RegisterBuyingStoreHandlers registers all buying store-related packet handlers with the factory.
func RegisterBuyingStoreHandlers(send HandlerRegistrar) {
	// Buying store handlers are registered through the buying store manager
	// No need to register individual handlers here as they're handled by the BuyingStoreManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a buying store manager and use it to register handlers
}
