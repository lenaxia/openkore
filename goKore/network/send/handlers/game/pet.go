// Package game provides game-related packet handlers.
package game

// RegisterPetHandlers registers all pet-related packet handlers with the factory.
func RegisterPetHandlers(send HandlerRegistrar) {
	// Pet handlers are registered through the pet manager
	// No need to register individual handlers here as they're handled by the PetManager

	// Note: In a real implementation, we might register specific handlers here
	// or create a pet manager and use it to register handlers
}
