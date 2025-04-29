package item

import (
	"fmt"
	"strconv"
	"strings"
)

// Constants for inventory types
const (
	INVTYPE_INVENTORY     uint8 = 0
	INVTYPE_CART          uint8 = 1
	INVTYPE_STORAGE       uint8 = 2
	INVTYPE_GUILD_STORAGE uint8 = 3
)

// HandleItemListStart handles the item_list_start packet (lines 5232-5247)
// This packet marks the beginning of an item list packet sequence
func (im *InventoryManager) HandleItemListStart(args map[string]interface{}) error {
	// Extract packet data
	invType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in item_list_start packet")
	}

	// Get container name if provided
	containerName := ""
	if name, ok := args["name"].(string); ok {
		containerName = name
	}

	// Log the item list start
	if containerName != "" {
		im.logger.Debug("Starting Item List. Type: %d (%s)", invType, containerName)
	} else {
		im.logger.Debug("Starting Item List. Type: %d", invType)
	}

	// Get inventory type name for hook
	var invTypeName string
	switch invType {
	case INVTYPE_INVENTORY:
		invTypeName = "inventory"
	case INVTYPE_CART:
		invTypeName = "cart"
	case INVTYPE_STORAGE:
		invTypeName = "storage"
	case INVTYPE_GUILD_STORAGE:
		invTypeName = "guild storage"
	default:
		invTypeName = "unknown"
		im.logger.Warning("Unsupported item_list_start type (%d)", invType)
	}

	// Call hooks
	im.hooks.CallHook("item_list_start", map[string]interface{}{
		"type":      invType,
		"typeName":  invTypeName,
		"container": containerName,
	})

	return nil
}

// HandleItemListStackable handles the item_list_stackable packet (lines 5249-5289)
// This packet contains a list of stackable items for a specific container type
func (im *InventoryManager) HandleItemListStackable(args map[string]interface{}) error {
	// Extract packet data
	invType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in item_list_stackable packet")
	}

	// Extract item info
	itemInfo, ok := args["itemInfo"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid itemInfo in item_list_stackable packet")
	}

	// Log the stackable items
	im.logger.Debug("Received %d stackable items for container type %d", len(itemInfo), invType)

	// Determine hook name based on container type
	var hookName string
	switch invType {
	case INVTYPE_INVENTORY:
		hookName = "packet_inventory"
	case INVTYPE_CART:
		hookName = "packet_cart"
	case INVTYPE_STORAGE, INVTYPE_GUILD_STORAGE:
		hookName = "packet_storage"
	default:
		im.logger.Warning("Unsupported item_list_stackable type (%d)", invType)
		return nil
	}

	// Process each item
	for _, item := range itemInfo {
		// Check for arrow equipment
		ID, _ := item["ID"].(string)
		nameID, _ := item["nameID"].(uint16)
		amount, _ := item["amount"].(uint16)
		itemType, _ := item["type"].(uint8)
		identified, _ := item["identified"].(uint8)

		// Call hooks for each item
		im.hooks.CallHook(hookName, map[string]interface{}{
			"ID":         ID,
			"nameID":     nameID,
			"amount":     amount,
			"type":       itemType,
			"identified": identified == 1,
		})
	}

	return nil
}

// HandleItemListNonstackable handles the item_list_nonstackable packet (lines 5291-5339)
// This packet contains a list of non-stackable items for a specific container type
func (im *InventoryManager) HandleItemListNonstackable(args map[string]interface{}) error {
	// Extract packet data
	invType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in item_list_nonstackable packet")
	}

	// Extract item info
	itemInfo, ok := args["itemInfo"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid itemInfo in item_list_nonstackable packet")
	}

	// Log the non-stackable items
	im.logger.Debug("Received %d non-stackable items for container type %d", len(itemInfo), invType)

	// Determine hook name based on container type
	var hookName string
	switch invType {
	case INVTYPE_INVENTORY:
		hookName = "packet_inventory"
	case INVTYPE_CART:
		hookName = "packet_cart"
	case INVTYPE_STORAGE, INVTYPE_GUILD_STORAGE:
		hookName = "packet_storage"
	default:
		im.logger.Warning("Unsupported item_list_nonstackable type (%d)", invType)
		return nil
	}

	// Process each item
	for _, item := range itemInfo {
		// Extract item data
		ID, _ := item["ID"].(string)
		nameID, _ := item["nameID"].(uint16)
		amount, _ := item["amount"].(uint16)
		itemType, _ := item["type"].(uint8)
		identified, _ := item["identified"].(uint8)
		broken, _ := item["broken"].(uint8)
		upgrade, _ := item["upgrade"].(uint8)
		cards, _ := item["cards"].(string)
		expire, _ := item["expire"].(uint32)
		equipped, _ := item["equipped"].(uint32)

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
		im.hooks.CallHook(hookName, map[string]interface{}{
			"ID":         ID,
			"nameID":     nameID,
			"amount":     amount,
			"type":       itemType,
			"identified": identified == 1,
			"broken":     broken == 1,
			"upgrade":    upgrade,
			"cards":      cardIDs,
			"expire":     expire,
			"equipped":   equipped,
		})

		// Handle equipped items (in the original code this is done in the callback)
		if equipped > 0 && hookName == "packet_inventory" {
			// In the original code, this would update character equipment
			// For now, we'll just log it
			im.logger.Debug("Item [%s] is equipped (flag: %d)", ID, equipped)
		}
	}

	return nil
}

// HandleItemListEnd handles the item_list_end packet (lines 5341-5354)
// This packet marks the end of an item list packet sequence
func (im *InventoryManager) HandleItemListEnd(args map[string]interface{}) error {
	// Extract packet data
	invType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in item_list_end packet")
	}

	// Get inventory type name for logging
	var invTypeName string
	switch invType {
	case INVTYPE_INVENTORY:
		invTypeName = "inventory"
	case INVTYPE_CART:
		invTypeName = "cart"
	case INVTYPE_STORAGE:
		invTypeName = "storage"
	case INVTYPE_GUILD_STORAGE:
		invTypeName = "guild storage"
	default:
		invTypeName = "unknown"
		im.logger.Warning("Unsupported item_list_end type (%d)", invType)
	}

	// Log the item list end
	im.logger.Info("Ending Item List. Type: %d (%s)", invType, invTypeName)

	// Call hooks
	im.hooks.CallHook("item_list_end", map[string]interface{}{
		"type":     invType,
		"typeName": invTypeName,
	})

	return nil
}
