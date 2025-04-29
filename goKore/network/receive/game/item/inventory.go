package item

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// Parser defines the interface for packet parsing and handling
type Parser interface {
	RegisterHandler(packetID, name, format string, paramNames []string, handler core.PacketHandler)
}

// InventoryManager handles inventory-related packet handling
type InventoryManager struct {
	baseParse Parser
	hooks     *hooks.HookManager
	logger    core.Logger
}

// NewInventoryManager creates a new inventory manager
func NewInventoryManager(baseParse Parser, hooks *hooks.HookManager, logger core.Logger) *InventoryManager {
	return &InventoryManager{
		baseParse: baseParse,
		hooks:     hooks,
		logger:    logger,
	}
}

// RegisterHandlers registers all inventory-related packet handlers
func (im *InventoryManager) RegisterHandlers() {
	// Register inventory item added handler (0x07FA)
	im.baseParse.RegisterHandler("07FA", "inventory_item_added", "v V v C3 a8 C V2 a*",
		[]string{"index", "amount", "nameID", "identified", "damaged", "refine", "card", "type", "location", "wear_state", "bindOnEquipType", "options"},
		im.HandleInventoryItemAdded)

	// Register inventory item removed handler (0x07FB)
	im.baseParse.RegisterHandler("07FB", "inventory_item_removed", "v V",
		[]string{"index", "amount"},
		im.HandleInventoryItemRemoved)

	// Register inventory items stackable handler (0x0991)
	im.baseParse.RegisterHandler("0991", "inventory_items_stackable", "v a*",
		[]string{"len", "itemInfo"},
		im.HandleInventoryItemsStackable)

	// Register inventory items non-stackable handler (0x0992)
	im.baseParse.RegisterHandler("0992", "inventory_items_nonstackable", "v a*",
		[]string{"len", "itemInfo"},
		im.HandleInventoryItemsNonstackable)

	// Register inventory item favorite handler (0x0997)
	im.baseParse.RegisterHandler("0997", "inventory_item_favorite", "v C",
		[]string{"index", "flag"},
		im.HandleInventoryItemFavorite)

	// Register inventory expansion result handler (0x0A37)
	im.baseParse.RegisterHandler("0A37", "inventory_expansion_result", "C",
		[]string{"result"},
		im.HandleInventoryExpansionResult)

	// Register item preview handler (0x0A38)
	im.baseParse.RegisterHandler("0A38", "item_preview", "v C C a8 a*",
		[]string{"index", "broken", "upgrade", "cards", "options"},
		im.HandleItemPreview)

	// Register special item obtain handler (0x0A0E)
	im.baseParse.RegisterHandler("0A0E", "special_item_obtain", "C v Z24 a*",
		[]string{"type", "nameID", "holder", "etc"},
		im.HandleSpecialItemObtain)

	// Register item list start handler (0x09A0)
	im.baseParse.RegisterHandler("09A0", "item_list_start", "v V",
		[]string{"type", "name"},
		im.HandleItemListStart)

	// Register item list stackable handler (0x09A3)
	im.baseParse.RegisterHandler("09A3", "item_list_stackable", "v v a*",
		[]string{"type", "itemCount", "itemInfo"},
		im.HandleItemListStackable)

	// Register item list non-stackable handler (0x09A4)
	im.baseParse.RegisterHandler("09A4", "item_list_nonstackable", "v v a*",
		[]string{"type", "itemCount", "itemInfo"},
		im.HandleItemListNonstackable)

	// Register item list end handler (0x09A5)
	im.baseParse.RegisterHandler("09A5", "item_list_end", "v",
		[]string{"type"},
		im.HandleItemListEnd)

	// Register ground item handlers

	// Register item appeared handler (0x009E)
	im.baseParse.RegisterHandler("009E", "item_appeared", "L W B W W B B W",
		[]string{"ID", "nameID", "identified", "x", "y", "subX", "subY", "amount"},
		im.HandleItemAppeared)

	// Register item appeared handler with type (0x084B)
	im.baseParse.RegisterHandler("084B", "item_appeared_with_type", "L W W B W W B B W",
		[]string{"ID", "nameID", "type", "identified", "x", "y", "subX", "subY", "amount"},
		im.HandleItemAppeared)

	// Register item appeared handler with effects (0x0ADD)
	im.baseParse.RegisterHandler("0ADD", "item_appeared_with_effects", "L W W B W W B B W B W",
		[]string{"ID", "nameID", "type", "identified", "x", "y", "subX", "subY", "amount", "show_effect", "effect_type"},
		im.HandleItemAppeared)

	// Register item exists handler (0x00A0)
	im.baseParse.RegisterHandler("00A0", "item_exists", "L W B W W B B W",
		[]string{"ID", "nameID", "identified", "x", "y", "subX", "subY", "amount"},
		im.HandleItemExists)

	// Register item exists handler with type (0x0ADD)
	im.baseParse.RegisterHandler("0ADD", "item_exists_with_effects", "L W W B W W B B W B W",
		[]string{"ID", "nameID", "type", "identified", "x", "y", "subX", "subY", "amount", "show_effect", "effect_type"},
		im.HandleItemExists)

	// Register item disappeared handler (0x00A1)
	im.baseParse.RegisterHandler("00A1", "item_disappeared", "L",
		[]string{"ID"},
		im.HandleItemDisappeared)
}

// HandleInventoryItemAdded handles the inventory_item_added packet (lines 7068-7110)
func (im *InventoryManager) HandleInventoryItemAdded(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(string)
	if !ok {
		return fmt.Errorf("invalid ID in inventory_item_added packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in inventory_item_added packet")
	}

	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in inventory_item_added packet")
	}

	// Process based on fail code
	if fail != 0 {
		switch fail {
		case 1:
			im.logger.Info("Failed to add item to inventory: inventory weight limit exceeded")
		case 2:
			im.logger.Info("Failed to add item to inventory: inventory is full")
		case 3:
			im.logger.Info("Failed to add item to inventory: item does not exist")
		case 4:
			im.logger.Info("Failed to add item to inventory: stack overflow")
		case 5:
			im.logger.Info("Failed to add item to inventory: stack underflow")
		default:
			im.logger.Info("Failed to add item to inventory: unknown reason (%d)", fail)
		}
		return nil
	}

	// Extract additional item data for successful addition
	nameID, ok := args["nameID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid nameID in inventory_item_added packet")
	}

	itemType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in inventory_item_added packet")
	}

	identified, ok := args["identified"].(uint8)
	if !ok {
		return fmt.Errorf("invalid identified in inventory_item_added packet")
	}

	broken, ok := args["broken"].(uint8)
	if !ok {
		return fmt.Errorf("invalid broken in inventory_item_added packet")
	}

	upgrade, ok := args["upgrade"].(uint8)
	if !ok {
		return fmt.Errorf("invalid upgrade in inventory_item_added packet")
	}

	cards, ok := args["cards"].(string)
	if !ok {
		return fmt.Errorf("invalid cards in inventory_item_added packet")
	}

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

	im.logger.Info("Item added to inventory: %s%s%s[%d] x%d",
		upgradeStr, identifiedStr, brokenStr, nameID, amount)

	// Call hooks
	im.hooks.CallHook("item_gathered", map[string]interface{}{
		"ID":         ID,
		"nameID":     nameID,
		"amount":     amount,
		"type":       itemType,
		"identified": identified == 1,
		"broken":     broken == 1,
		"upgrade":    upgrade,
		"cards":      cards,
	})

	return nil
}

// HandleInventoryItemRemoved handles the inventory_item_removed packet (lines 7111-7142)
func (im *InventoryManager) HandleInventoryItemRemoved(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(string)
	if !ok {
		return fmt.Errorf("invalid ID in inventory_item_removed packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in inventory_item_removed packet")
	}

	reason, ok := args["reason"].(uint8)
	if !ok {
		return fmt.Errorf("invalid reason in inventory_item_removed packet")
	}

	// Log the item removal with reason
	var reasonStr string
	switch reason {
	case 0:
		reasonStr = "Normal"
	case 1:
		reasonStr = "Used to cast a skill"
	case 2:
		reasonStr = "Dropped"
	case 3:
		reasonStr = "Traded"
	case 4:
		reasonStr = "Sold"
	case 5:
		reasonStr = "Consumed"
	default:
		reasonStr = fmt.Sprintf("Unknown (%d)", reason)
	}

	im.logger.Debug("Item removed from inventory: [%s] x%d (Reason: %s)", ID, amount, reasonStr)

	// Call hooks
	im.hooks.CallHook("packet_item_removed", map[string]interface{}{
		"ID":     ID,
		"amount": amount,
		"reason": reason,
	})

	return nil
}

// HandleInventoryItemsStackable handles the inventory_items_stackable packet (lines 7143-7179)
func (im *InventoryManager) HandleInventoryItemsStackable(args map[string]interface{}) error {
	// Extract packet data
	itemInfo, ok := args["itemInfo"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid itemInfo in inventory_items_stackable packet")
	}

	// Log the stackable items
	im.logger.Debug("Received %d stackable inventory items", len(itemInfo))

	// Process each item
	for _, item := range itemInfo {
		ID, _ := item["ID"].(string)
		nameID, _ := item["nameID"].(uint16)
		amount, _ := item["amount"].(uint16)
		itemType, _ := item["type"].(uint8)
		identified, _ := item["identified"].(uint8)

		// Call hooks for each item
		im.hooks.CallHook("packet_inventory", map[string]interface{}{
			"ID":         ID,
			"nameID":     nameID,
			"amount":     amount,
			"type":       itemType,
			"identified": identified == 1,
		})
	}

	return nil
}

// HandleInventoryItemsNonstackable handles the inventory_items_nonstackable packet (lines 7180-7229)
func (im *InventoryManager) HandleInventoryItemsNonstackable(args map[string]interface{}) error {
	// Extract packet data
	itemInfo, ok := args["itemInfo"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid itemInfo in inventory_items_nonstackable packet")
	}

	// Log the non-stackable items
	im.logger.Debug("Received %d non-stackable inventory items", len(itemInfo))

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
		im.hooks.CallHook("packet_inventory", map[string]interface{}{
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

// HandleInventoryItemFavorite handles the inventory_item_favorite packet (lines 7230-7246)
func (im *InventoryManager) HandleInventoryItemFavorite(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(string)
	if !ok {
		return fmt.Errorf("invalid ID in inventory_item_favorite packet")
	}

	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in inventory_item_favorite packet")
	}

	// Log the favorite status change
	if flag == 1 {
		im.logger.Info("Item [%s] marked as favorite", ID)
	} else {
		im.logger.Info("Item [%s] unmarked as favorite", ID)
	}

	// Call hooks
	im.hooks.CallHook("inventory_item_favorite", map[string]interface{}{
		"ID":       ID,
		"favorite": flag == 1,
	})

	return nil
}

// HandleItemPreview handles the item_preview packet (lines 12182-12193)
func (im *InventoryManager) HandleItemPreview(args map[string]interface{}) error {
	// Extract packet data
	index, ok := args["index"].(uint16)
	if !ok {
		return fmt.Errorf("invalid index in item_preview packet")
	}

	// Log the item preview
	im.logger.Debug("Item preview for index: %d", index)

	// Extract optional parameters
	var broken uint8
	if brokenVal, ok := args["broken"].(uint8); ok {
		broken = brokenVal
	}

	var upgrade uint8
	if upgradeVal, ok := args["upgrade"].(uint8); ok {
		upgrade = upgradeVal
	}

	var cards string
	if cardsVal, ok := args["cards"].(string); ok {
		cards = cardsVal
	}

	var options string
	if optionsVal, ok := args["options"].(string); ok {
		options = optionsVal
	}

	// Call hooks
	im.hooks.CallHook("item_preview", map[string]interface{}{
		"index":   index,
		"broken":  broken,
		"upgrade": upgrade,
		"cards":   cards,
		"options": options,
	})

	return nil
}

// Constants for special item obtain types
const (
	TYPE_BOXITEM      = 0
	TYPE_MONSTER_ITEM = 1
)

// HandleSpecialItemObtain handles the special_item_obtain packet (lines 9904-9953)
func (im *InventoryManager) HandleSpecialItemObtain(args map[string]interface{}) error {
	// Extract packet data
	obtainType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in special_item_obtain packet")
	}

	nameID, ok := args["nameID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid nameID in special_item_obtain packet")
	}

	holder, ok := args["holder"].(string)
	if !ok {
		return fmt.Errorf("invalid holder in special_item_obtain packet")
	}

	etc, ok := args["etc"].([]byte)
	if !ok {
		return fmt.Errorf("invalid etc in special_item_obtain packet")
	}

	// Get item name (simplified for now)
	itemName := fmt.Sprintf("Item#%d", nameID)

	var sourceItemID uint16
	var sourceName string
	var msg string

	// Process based on obtain type
	switch obtainType {
	case TYPE_BOXITEM:
		// Parse box item ID
		var boxNameID uint16
		if len(etc) > 0 && etc[0] == 2 {
			// c/v format (2 bytes)
			if len(etc) >= 3 {
				boxNameID = uint16(etc[2]) | (uint16(etc[1]) << 8)
			}
		} else {
			// c/V format (4 bytes)
			if len(etc) >= 5 {
				boxNameID = uint16(etc[4]) | (uint16(etc[3]) << 8) | (uint16(etc[2]) << 16) | (uint16(etc[1]) << 24)
			}
		}

		// Get box item name (simplified for now)
		boxItemName := fmt.Sprintf("Item#%d", boxNameID)
		sourceName = boxItemName
		sourceItemID = boxNameID

		// Log the message
		msg = fmt.Sprintf("%s has got %s from %s.", holder, itemName, boxItemName)
		im.logger.Info(msg)

	case TYPE_MONSTER_ITEM:
		// Parse monster name
		if len(etc) > 0 {
			monsterNameLen := int(etc[0])
			if len(etc) > monsterNameLen {
				monsterName := string(etc[1 : monsterNameLen+1])
				sourceName = monsterName

				// Log the message
				msg = fmt.Sprintf("%s has got %s from %s.", holder, itemName, monsterName)
				im.logger.Info(msg)
			}
		}

	default:
		// Unknown type
		msg = fmt.Sprintf("%s has got %s (from Unknown type %d).", holder, itemName, obtainType)
		im.logger.Warning(msg)
	}

	// Call hooks
	im.hooks.CallHook("packet_special_item_obtain", map[string]interface{}{
		"ObtainType":   obtainType,
		"ItemName":     itemName,
		"ItemID":       nameID,
		"Holder":       holder,
		"SourceItemID": sourceItemID,
		"SourceName":   sourceName,
		"Msg":          msg,
	})

	return nil
}

// HandleInventoryExpansionResult handles the inventory_expansion_result packet (lines 7247-7263)
func (im *InventoryManager) HandleInventoryExpansionResult(args map[string]interface{}) error {
	// Extract packet data
	result, ok := args["result"].(uint8)
	if !ok {
		return fmt.Errorf("invalid result in inventory_expansion_result packet")
	}

	// Process based on result code
	switch result {
	case 0: // EXPAND_INVENTORY_RESULT_SUCCESS
		im.logger.Info("Inventory expansion successful")
	case 1: // EXPAND_INVENTORY_RESULT_FAILED
		im.logger.Info("Inventory expansion failed")
	case 2: // EXPAND_INVENTORY_RESULT_MAXIMUM_SIZE
		im.logger.Info("Inventory expansion failed: already at maximum size")
	case 3: // EXPAND_INVENTORY_RESULT_MISSING_ITEM
		im.logger.Info("Inventory expansion failed: missing required item")
	default:
		im.logger.Info("Inventory expansion failed: unknown reason (%d)", result)
	}

	// Call hooks
	im.hooks.CallHook("inventory_expansion_result", map[string]interface{}{
		"result": result,
	})

	return nil
}
