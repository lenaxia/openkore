package rental

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// RentalManager handles rental-related packet handling
type RentalManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger
	inventory   map[uint32]map[string]interface{} // Simple inventory representation
}

// NewRentalManager creates a new rental manager
func NewRentalManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *RentalManager {
	return &RentalManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
		inventory:   make(map[uint32]map[string]interface{}),
	}
}

// RegisterHandlers registers all rental-related packet handlers
func (rm *RentalManager) RegisterHandlers() {
	// Register rental expired handler
	rm.parser.RegisterHandlerFunc("0299", "rental_expired", "L W",
		[]string{"ID", "nameID"}, rm.HandleRentalExpired)

	// Register rental time handler
	rm.parser.RegisterHandlerFunc("01DF", "rental_time", "W L",
		[]string{"nameID", "seconds"}, rm.HandleRentalTime)
}

// SetInventory sets the inventory for the rental manager
func (rm *RentalManager) SetInventory(inventory map[uint32]map[string]interface{}) {
	rm.inventory = inventory
}

// GetItemName gets the name of an item by its nameID
func (rm *RentalManager) GetItemName(nameID uint16) string {
	// In a real implementation, this would look up the item name from a database
	// For now, we'll just return a placeholder
	return fmt.Sprintf("Item#%d", nameID)
}

// HandleRentalExpired handles the rental_expired packet (lines 3789-3801)
func (rm *RentalManager) HandleRentalExpired(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in rental_expired packet")
	}

	nameID, ok := args["nameID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid nameID in rental_expired packet")
	}

	// Get item name
	itemName := rm.GetItemName(nameID)

	// Log rental expiration
	rm.logger.Info("Rental item '%s' has expired!", itemName)

	// Check if item exists in inventory
	item, exists := rm.inventory[id]
	if exists {
		// Get binID
		binID, ok := item["binID"].(uint16)
		if !ok {
			return fmt.Errorf("invalid binID in inventory item")
		}

		// Call hook
		rm.hookManager.CallHook("rental_expired", map[string]interface{}{
			"index":  binID,
			"nameID": nameID,
		})

		// Remove item from inventory
		delete(rm.inventory, id)
	}

	return nil
}

// HandleRentalTime handles the rental_time packet (lines 10607-10610)
func (rm *RentalManager) HandleRentalTime(args map[string]interface{}) error {
	// Extract packet data
	nameID, ok := args["nameID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid nameID in rental_time packet")
	}

	seconds, ok := args["seconds"].(uint32)
	if !ok {
		return fmt.Errorf("invalid seconds in rental_time packet")
	}

	// Get item name
	itemName := rm.GetItemName(nameID)

	// Calculate minutes
	minutes := seconds / 60

	// Log rental time
	rm.logger.Info("The '%s' item will disappear in %d minutes.", itemName, minutes)

	// Call hook
	rm.hookManager.CallHook("rental_time", map[string]interface{}{
		"nameID":  nameID,
		"seconds": seconds,
		"minutes": minutes,
	})

	return nil
}
