package item

import (
	"fmt"
)

// HandleItemAppeared handles the item_appeared packet (lines 7015-7052)
// This packet makes an item appear on the ground
// 009E <id>.L <name id>.W <identified>.B <x>.W <y>.W <subX>.B <subY>.B <amount>.W (ZC_ITEM_FALL_ENTRY)
// 084B <id>.L <name id>.W <type>.W <identified>.B <x>.W <y>.W <subX>.B <subY>.B <amount>.W (ZC_ITEM_FALL_ENTRY4)
// 0ADD <id>.L <name id>.W <type>.W <identified>.B <x>.W <y>.W <subX>.B <subY>.B <amount>.W <show drop effect>.B <drop effect mode>.W (ZC_ITEM_FALL_ENTRY5)
func (im *InventoryManager) HandleItemAppeared(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(string)
	if !ok {
		return fmt.Errorf("invalid ID in item_appeared packet")
	}

	nameID, ok := args["nameID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid nameID in item_appeared packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in item_appeared packet")
	}

	identified, ok := args["identified"].(uint8)
	if !ok {
		return fmt.Errorf("invalid identified in item_appeared packet")
	}

	x, ok := args["x"].(uint16)
	if !ok {
		return fmt.Errorf("invalid x in item_appeared packet")
	}

	y, ok := args["y"].(uint16)
	if !ok {
		return fmt.Errorf("invalid y in item_appeared packet")
	}

	// Get item type if available
	var itemType uint16
	if typeVal, ok := args["type"].(uint16); ok {
		itemType = typeVal
	}

	// Log the item appearance
	im.logger.Info("Item Appeared: [%s] x%d at (%d, %d)", ID, amount, x, y)

	// Call hooks
	im.hooks.CallHook("item_appeared", map[string]interface{}{
		"ID":         ID,
		"nameID":     nameID,
		"amount":     amount,
		"identified": identified == 1,
		"x":          x,
		"y":          y,
		"type":       itemType,
	})

	return nil
}

// HandleItemExists handles the item_exists packet (lines 7054-7084)
// This packet indicates an item exists on the ground
func (im *InventoryManager) HandleItemExists(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(string)
	if !ok {
		return fmt.Errorf("invalid ID in item_exists packet")
	}

	nameID, ok := args["nameID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid nameID in item_exists packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in item_exists packet")
	}

	identified, ok := args["identified"].(uint8)
	if !ok {
		return fmt.Errorf("invalid identified in item_exists packet")
	}

	x, ok := args["x"].(uint16)
	if !ok {
		return fmt.Errorf("invalid x in item_exists packet")
	}

	y, ok := args["y"].(uint16)
	if !ok {
		return fmt.Errorf("invalid y in item_exists packet")
	}

	// Get optional parameters
	var itemType uint16
	if typeVal, ok := args["type"].(uint16); ok {
		itemType = typeVal
	}

	var showEffect uint8
	if effectVal, ok := args["show_effect"].(uint8); ok {
		showEffect = effectVal
	}

	var effectType uint16
	if effectTypeVal, ok := args["effect_type"].(uint16); ok {
		effectType = effectTypeVal
	}

	// Log the item existence
	im.logger.Info("Item Exists: [%s] x%d at (%d, %d)", ID, amount, x, y)

	// Call hooks
	im.hooks.CallHook("item_exists", map[string]interface{}{
		"ID":          ID,
		"nameID":      nameID,
		"amount":      amount,
		"identified":  identified == 1,
		"x":           x,
		"y":           y,
		"type":        itemType,
		"show_effect": showEffect,
		"effect_type": effectType,
	})

	return nil
}

// HandleItemDisappeared handles the item_disappeared packet (lines 7088-7119)
// This packet makes an item disappear from the ground
// 00A1 <id>.L (ZC_ITEM_DISAPPEAR)
func (im *InventoryManager) HandleItemDisappeared(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(string)
	if !ok {
		return fmt.Errorf("invalid ID in item_disappeared packet")
	}

	// Log the item disappearance
	im.logger.Debug("Item Disappeared: [%s]", ID)

	// Call hooks
	im.hooks.CallHook("item_disappeared", map[string]interface{}{
		"ID": ID,
	})

	return nil
}
