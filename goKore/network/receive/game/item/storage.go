package item

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// StorageManager handles storage-related packet handling
type StorageManager struct {
	baseParse Parser
	hooks     *hooks.HookManager
	logger    core.Logger
}

// NewStorageManager creates a new storage manager
func NewStorageManager(baseParse Parser, hooks *hooks.HookManager, logger core.Logger) *StorageManager {
	return &StorageManager{
		baseParse: baseParse,
		hooks:     hooks,
		logger:    logger,
	}
}

// RegisterHandlers registers all storage-related packet handlers
func (sm *StorageManager) RegisterHandlers() {
	// Register storage opened handler (0x09A9)
	sm.baseParse.RegisterHandler("09A9", "storage_opened", "v V",
		[]string{"items", "items_max"},
		sm.HandleStorageOpened)

	// Register storage closed handler (0x09AA)
	sm.baseParse.RegisterHandler("09AA", "storage_closed", "",
		[]string{},
		sm.HandleStorageClosed)

	// Register storage items stackable handler (0x0995)
	sm.baseParse.RegisterHandler("0995", "storage_items_stackable", "v a*",
		[]string{"len", "itemInfo"},
		sm.HandleStorageItemsStackable)

	// Register storage items non-stackable handler (0x0996)
	sm.baseParse.RegisterHandler("0996", "storage_items_nonstackable", "v a*",
		[]string{"len", "itemInfo"},
		sm.HandleStorageItemsNonstackable)

	// Register storage item added handler (0x0A0A)
	sm.baseParse.RegisterHandler("0A0A", "storage_item_added", "v V v C3 a8 C V2 a*",
		[]string{"index", "amount", "nameID", "identified", "damaged", "refine", "card", "type", "location", "wear_state", "bindOnEquipType", "options"},
		sm.HandleStorageItemAdded)

	// Register storage item removed handler (0x0A0B)
	sm.baseParse.RegisterHandler("0A0B", "storage_item_removed", "v V",
		[]string{"index", "amount"},
		sm.HandleStorageItemRemoved)

	// Register storage password request handler (0x023A)
	sm.baseParse.RegisterHandler("023A", "storage_password_request", "C",
		[]string{"flag"},
		sm.HandleStoragePasswordRequest)

	// Register storage password result handler (0x023E)
	sm.baseParse.RegisterHandler("023E", "storage_password_result", "C",
		[]string{"type"},
		sm.HandleStoragePasswordResult)

	// Register guild storage log handler (0x09D0)
	sm.baseParse.RegisterHandler("09D0", "guild_storage_log", "C a*",
		[]string{"result", "log"},
		sm.HandleGuildStorageLog)
}

// HandleStorageOpened handles the storage_opened packet (lines 10923-10935)
func (sm *StorageManager) HandleStorageOpened(args map[string]interface{}) error {
	// Extract packet data if available
	items, ok := args["items"].(uint16)
	if ok {
		sm.logger.Info("Storage opened (%d/%d items)", items, args["items_max"])
	} else {
		sm.logger.Info("Storage opened")
	}

	// Call hooks
	sm.hooks.CallHook("storage_opened", nil)

	return nil
}

// HandleStorageClosed handles the storage_closed packet (lines 10936-10947)
func (sm *StorageManager) HandleStorageClosed(args map[string]interface{}) error {
	// Log the storage closure
	sm.logger.Info("Storage closed")

	// Call hooks
	sm.hooks.CallHook("storage_closed", nil)

	return nil
}

// HandleStorageItemAdded handles the storage_item_added packet (lines 10948-10979)
func (sm *StorageManager) HandleStorageItemAdded(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(string)
	if !ok {
		return fmt.Errorf("invalid ID in storage_item_added packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in storage_item_added packet")
	}

	nameID, ok := args["nameID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid nameID in storage_item_added packet")
	}

	// Extract additional item data if available
	identified, _ := args["identified"].(uint8)
	broken, _ := args["broken"].(uint8)
	upgrade, _ := args["upgrade"].(uint8)

	// Log the item addition
	identifiedStr := ""
	if identified == 0 {
		identifiedStr = "unidentified "
	}

	upgradeStr := ""
	if upgrade > 0 {
		upgradeStr = fmt.Sprintf("+%d ", upgrade)
	}

	brokenStr := ""
	if broken == 1 {
		brokenStr = "broken "
	}

	sm.logger.Info("Item added to storage: %s%s%s[%d] x%d",
		upgradeStr, identifiedStr, brokenStr, nameID, amount)

	// Call hooks
	sm.hooks.CallHook("packet_storage_add", map[string]interface{}{
		"ID":         ID,
		"nameID":     nameID,
		"amount":     amount,
		"identified": identified == 1,
		"broken":     broken == 1,
		"upgrade":    upgrade,
	})

	return nil
}

// HandleStorageItemRemoved handles the storage_item_removed packet (lines 10980-11000)
func (sm *StorageManager) HandleStorageItemRemoved(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(string)
	if !ok {
		return fmt.Errorf("invalid ID in storage_item_removed packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in storage_item_removed packet")
	}

	// Call hooks
	sm.hooks.CallHook("packet_storage_remove", map[string]interface{}{
		"ID":     ID,
		"amount": amount,
	})

	return nil
}

// HandleStorageItemsStackable handles the storage_items_stackable packet (lines 11001-11037)
func (sm *StorageManager) HandleStorageItemsStackable(args map[string]interface{}) error {
	// Extract packet data
	itemInfo, ok := args["itemInfo"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid itemInfo in storage_items_stackable packet")
	}

	// Log the stackable items
	sm.logger.Debug("Received %d stackable storage items", len(itemInfo))

	// Process each item
	for _, item := range itemInfo {
		ID, _ := item["ID"].(string)
		nameID, _ := item["nameID"].(uint16)
		amount, _ := item["amount"].(uint16)
		itemType, _ := item["type"].(uint8)
		identified, _ := item["identified"].(uint8)

		// Call hooks for each item
		sm.hooks.CallHook("packet_storage", map[string]interface{}{
			"ID":         ID,
			"nameID":     nameID,
			"amount":     amount,
			"type":       itemType,
			"identified": identified == 1,
		})
	}

	return nil
}

// HandleStorageItemsNonstackable handles the storage_items_nonstackable packet (lines 11038-11087)
func (sm *StorageManager) HandleStorageItemsNonstackable(args map[string]interface{}) error {
	// Extract packet data
	itemInfo, ok := args["itemInfo"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid itemInfo in storage_items_nonstackable packet")
	}

	// Log the non-stackable items
	sm.logger.Debug("Received %d non-stackable storage items", len(itemInfo))

	// Process each item
	for _, item := range itemInfo {
		ID, _ := item["ID"].(string)
		nameID, _ := item["nameID"].(uint16)
		amount, _ := item["amount"].(uint16)
		itemType, _ := item["type"].(uint8)
		identified, _ := item["identified"].(uint8)
		broken, _ := item["broken"].(uint8)
		upgrade, _ := item["upgrade"].(uint8)
		cards, _ := item["cards"].(string)
		expire, _ := item["expire"].(uint32)

		// Parse cards
		cardIDs := []uint16{}
		if cards != "" {
			cardParts := strings.Split(cards, ",")
			for _, cardPart := range cardParts {
				cardID, err := strconv.ParseUint(cardPart, 10, 16)
				if err == nil {
					cardIDs = append(cardIDs, uint16(cardID))
				}
			}
		}

		// Call hooks for each item
		sm.hooks.CallHook("packet_storage", map[string]interface{}{
			"ID":         ID,
			"nameID":     nameID,
			"amount":     amount,
			"type":       itemType,
			"identified": identified == 1,
			"broken":     broken == 1,
			"upgrade":    upgrade,
			"cards":      cardIDs,
			"expire":     expire,
		})
	}

	return nil
}

// HandleStoragePasswordRequest handles the storage_password_request packet (lines 11088-11115)
func (sm *StorageManager) HandleStoragePasswordRequest(args map[string]interface{}) error {
	// Extract packet data
	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in storage_password_request packet")
	}

	// Process based on flag
	switch flag {
	case 0: // STORE_PASSWORD_NONE (new password)
		sm.logger.Info("Please enter a new storage password")
	case 1: // STORE_PASSWORD_CURRENT (verify current password)
		sm.logger.Info("Please enter your storage password")
	case 2: // STORE_PASSWORD_CHANGE (change password)
		sm.logger.Info("Please enter your current storage password")
	case 3: // STORE_PASSWORD_CHECK (checking password)
		sm.logger.Info("Checking storage password...")
	case 8: // STORE_PASSWORD_ERROR (too many incorrect attempts)
		sm.logger.Error("Too many incorrect storage password attempts")
	default:
		sm.logger.Warning("Unknown storage password request: %d", flag)
	}

	// Call hooks
	sm.hooks.CallHook("storage_password_request", map[string]interface{}{
		"flag": flag,
	})

	return nil
}

// HandleStoragePasswordResult handles the storage_password_result packet (lines 11116-11151)
func (sm *StorageManager) HandleStoragePasswordResult(args map[string]interface{}) error {
	// Extract packet data
	type_, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in storage_password_result packet")
	}

	// Process based on type code
	switch type_ {
	case 0: // STORE_PASSWORD_EMPTY
		sm.logger.Info("Storage password is not set")
	case 1: // STORE_PASSWORD_SET
		sm.logger.Success("Storage password set")
	case 2: // STORE_PASSWORD_CLEARED
		sm.logger.Success("Storage password cleared")
	case 3: // STORE_PASSWORD_CHANGE_READY
		sm.logger.Info("Enter new storage password")
	case 4: // STORE_PASSWORD_CHANGE_OK
		sm.logger.Success("Storage password changed")
	case 5: // STORE_PASSWORD_CHANGE_NG
		sm.logger.Error("Failed to change storage password")
	case 6: // STORE_PASSWORD_CHECK_OK
		sm.logger.Success("Storage password verified")
	case 7: // STORE_PASSWORD_CHECK_NG
		sm.logger.Error("Incorrect storage password")
	default:
		sm.logger.Warning("Unknown storage password result: %d", type_)
	}

	// Call hooks
	sm.hooks.CallHook("storage_password_result", map[string]interface{}{
		"type": type_,
	})

	return nil
}

// HandleGuildStorageLog handles the guild_storage_log packet (lines 9577-9613)
func (sm *StorageManager) HandleGuildStorageLog(args map[string]interface{}) error {
	// Extract packet data
	result, ok := args["result"].(uint8)
	if !ok {
		return fmt.Errorf("invalid result in guild_storage_log packet")
	}

	// Process based on result code
	switch result {
	case 0, 1: // Get or Put actions
		// In the original code, this parses a complex log structure
		// For simplicity, we'll just log that we received guild storage log data
		sm.logger.Info("Guild Storage Log received (action: %d)", result)

		// In a full implementation, we would parse the log data and display it
		// But for now, we'll just acknowledge receipt of the data
		if log, ok := args["log"].([]byte); ok {
			sm.logger.Debug("Guild Storage Log contains %d bytes of data", len(log))
		}

	case 2: // Empty storage
		sm.logger.Info("Guild Storage empty")

	case 3: // Not using guild storage
		sm.logger.Info("You are not currently using Guild Storage. Please try later")

	default:
		sm.logger.Warning("Unknown guild storage log result: %d", result)
	}

	// Call hooks
	sm.hooks.CallHook("guild_storage_log", map[string]interface{}{
		"result": result,
	})

	return nil
}
