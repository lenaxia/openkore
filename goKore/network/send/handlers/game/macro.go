// Package game provides game-related packet handlers.
package game

// RegisterMacroHandlers registers all macro-related packet handlers with the factory.
func RegisterMacroHandlers(send HandlerRegistrar) {
	// Macro handlers are registered through the macro manager
	// No need to register individual handlers here as they're handled by the MacroManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a macro manager and use it to register handlers
}
