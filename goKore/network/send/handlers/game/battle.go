// Package game provides game-related packet handlers.
package game

// RegisterBattleHandlers registers all battle-related packet handlers with the factory.
func RegisterBattleHandlers(send HandlerRegistrar) {
	// Battle handlers are registered through the battle manager
	// No need to register individual handlers here as they're handled by the BattleManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a battle manager and use it to register handlers
}
