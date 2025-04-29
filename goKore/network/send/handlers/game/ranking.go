// Package game provides game-related packet handlers.
package game

// RegisterRankingHandlers registers all ranking-related packet handlers with the factory.
func RegisterRankingHandlers(send HandlerRegistrar) {
	// Ranking handlers are registered through the ranking manager
	// No need to register individual handlers here as they're handled by the RankingManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a ranking manager and use it to register handlers
}
