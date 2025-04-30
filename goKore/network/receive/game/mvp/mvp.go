package mvp

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// MVPManager handles MVP-related packet handling
type MVPManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger
}

// NewMVPManager creates a new MVP manager
func NewMVPManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *MVPManager {
	return &MVPManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
	}
}

// RegisterHandlers registers all MVP-related packet handlers
func (mm *MVPManager) RegisterHandlers() {
	// Register MVP item handler
	mm.parser.RegisterHandlerFunc("010A", "mvp_item", "w",
		[]string{"itemID"}, mm.HandleMVPItem)

	// Register MVP other handler
	mm.parser.RegisterHandlerFunc("010C", "mvp_other", "L",
		[]string{"ID"}, mm.HandleMVPOther)

	// Register MVP you handler
	mm.parser.RegisterHandlerFunc("010B", "mvp_you", "L",
		[]string{"expAmount"}, mm.HandleMVPYou)
}

// GetItemName gets the name of an item by its ID
func (mm *MVPManager) GetItemName(itemID uint16) string {
	// In a real implementation, this would look up the item name from a database
	// For now, we'll just return a placeholder
	return fmt.Sprintf("Item#%d", itemID)
}

// GetActorName gets the name of an actor by its ID
func (mm *MVPManager) GetActorName(actorID uint32) string {
	// In a real implementation, this would look up the actor name from the actor list
	// For now, we'll just return a placeholder
	return fmt.Sprintf("Actor#%d", actorID)
}

// HandleMVPItem handles the mvp_item packet (lines 11165-11170)
func (mm *MVPManager) HandleMVPItem(args map[string]interface{}) error {
	// Extract packet data
	itemID, ok := args["itemID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid itemID in mvp_item packet")
	}

	// Get item name
	itemName := mm.GetItemName(itemID)

	// Log message
	message := fmt.Sprintf("Get MVP item %s", itemName)
	mm.logger.Info(message)

	// Call hook
	mm.hookManager.CallHook("mvp_item", map[string]interface{}{
		"itemID":   itemID,
		"itemName": itemName,
	})

	return nil
}

// HandleMVPOther handles the mvp_other packet (lines 11172-11177)
func (mm *MVPManager) HandleMVPOther(args map[string]interface{}) error {
	// Extract packet data
	actorID, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in mvp_other packet")
	}

	// Get actor name
	actorName := mm.GetActorName(actorID)

	// Log message
	message := fmt.Sprintf("%s become MVP!", actorName)
	mm.logger.Info(message)

	// Call hook
	mm.hookManager.CallHook("mvp_other", map[string]interface{}{
		"ID":        actorID,
		"actorName": actorName,
	})

	return nil
}

// HandleMVPYou handles the mvp_you packet (lines 11179-11184)
func (mm *MVPManager) HandleMVPYou(args map[string]interface{}) error {
	// Extract packet data
	expAmount, ok := args["expAmount"].(uint32)
	if !ok {
		return fmt.Errorf("invalid expAmount in mvp_you packet")
	}

	// Log message
	message := fmt.Sprintf("Congratulations, you are the MVP! Your reward is %d exp!", expAmount)
	mm.logger.Info(message)

	// Call hook
	mm.hookManager.CallHook("mvp_you", map[string]interface{}{
		"expAmount": expAmount,
	})

	return nil
}
