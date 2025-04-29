// Package game provides game-related packet handlers.
package game

// RegisterAuctionHandlers registers all auction-related packet handlers with the factory.
func RegisterAuctionHandlers(send HandlerRegistrar) {
	// Auction handlers are registered through the auction manager
	// No need to register individual handlers here as they're handled by the AuctionManager

	// Note: In a real implementation, we might register specific handlers here
	// or create an auction manager and use it to register handlers
}
