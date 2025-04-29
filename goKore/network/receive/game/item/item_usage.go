package item

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// ItemUsageManager handles item usage-related packet handling
type ItemUsageManager struct {
	baseParse Parser
	hooks     *hooks.HookManager
	logger    core.Logger
	accountID string // Player's account ID
}

// NewItemUsageManager creates a new item usage manager
func NewItemUsageManager(baseParse Parser, hooks *hooks.HookManager, logger core.Logger) *ItemUsageManager {
	return &ItemUsageManager{
		baseParse: baseParse,
		hooks:     hooks,
		logger:    logger,
		accountID: "", // Will be set when the player logs in
	}
}

// RegisterHandlers registers all item usage-related packet handlers
func (im *ItemUsageManager) RegisterHandlers() {
	// Register item_used handler (0x0A0E)
	im.baseParse.RegisterHandler("0A0E", "item_used", "v v S v C",
		[]string{"ID", "itemID", "actorID", "remaining", "success"},
		im.HandleItemUsed)

	// Register use_item handler (0x01C8)
	im.baseParse.RegisterHandler("01C8", "use_item", "v v",
		[]string{"ID", "amount"},
		im.HandleUseItem)
}

// SetAccountID sets the player's account ID
func (im *ItemUsageManager) SetAccountID(accountID string) {
	im.accountID = accountID
}

// HandleItemUsed handles the item_used packet (lines 6957-7002)
func (im *ItemUsageManager) HandleItemUsed(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid ID in item_used packet")
	}

	itemID, ok := args["itemID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid itemID in item_used packet")
	}

	actorID, ok := args["actorID"].(string)
	if !ok {
		return fmt.Errorf("invalid actorID in item_used packet")
	}

	remaining, ok := args["remaining"].(uint16)
	if !ok {
		return fmt.Errorf("invalid remaining in item_used packet")
	}

	success, ok := args["success"].(uint8)
	if !ok {
		return fmt.Errorf("invalid success in item_used packet")
	}

	// Process based on who used the item
	if actorID == im.accountID {
		// The player used the item
		if success == 1 {
			// Successfully used the item
			im.logger.Info("You used Item: #%d (ID: %d) - %d remaining", itemID, ID, remaining)
		} else {
			// Failed to use the item
			im.logger.Info("You failed to use item: #%d (ID: %d) - %d remaining", itemID, ID, remaining)
		}
	} else {
		// Another actor used the item
		im.logger.Info("Actor %s used Item: #%d - %d remaining", actorID, itemID, remaining)
	}

	// Call hooks
	im.hooks.CallHook("packet_useitem", map[string]interface{}{
		"serverIndex": ID,
		"itemID":      itemID,
		"userID":      actorID,
		"remaining":   remaining,
		"success":     success == 1,
	})

	return nil
}

// HandleUseItem handles the use_item packet (lines 11673-11681)
func (im *ItemUsageManager) HandleUseItem(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid ID in use_item packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in use_item packet")
	}

	// Log the item usage
	im.logger.Info("You used Item: ID %d x %d", ID, amount)

	// Call hooks
	im.hooks.CallHook("item_used", map[string]interface{}{
		"ID":     ID,
		"amount": amount,
	})

	return nil
}
