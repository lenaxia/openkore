// Package game provides game-related packet handlers.
package game

// RegisterCashShopHandlers registers all cash shop-related packet handlers with the factory.
func RegisterCashShopHandlers(send HandlerRegistrar) {
	// Cash shop handlers are registered through the cash shop manager
	// No need to register individual handlers here as they're handled by the CashShopManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a cash shop manager and use it to register handlers
}
