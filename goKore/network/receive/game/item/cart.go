package item

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// CartManager handles cart-related packet handling
type CartManager struct {
	baseParse Parser
	hooks     *hooks.HookManager
	logger    core.Logger
}

// NewCartManager creates a new cart manager
func NewCartManager(baseParse Parser, hooks *hooks.HookManager, logger core.Logger) *CartManager {
	return &CartManager{
		baseParse: baseParse,
		hooks:     hooks,
		logger:    logger,
	}
}

// RegisterHandlers registers all cart-related packet handlers
func (cm *CartManager) RegisterHandlers() {
	// Register cart off handler (0x012B)
	cm.baseParse.RegisterHandler("012B", "cart_off", "",
		[]string{},
		cm.HandleCartOff)

	// Register cart items stackable handler (0x0993)
	cm.baseParse.RegisterHandler("0993", "cart_items_stackable", "v a*",
		[]string{"len", "itemInfo"},
		cm.HandleCartItemsStackable)

	// Register cart items non-stackable handler (0x0994)
	cm.baseParse.RegisterHandler("0994", "cart_items_nonstackable", "v a*",
		[]string{"len", "itemInfo"},
		cm.HandleCartItemsNonstackable)

	// Register cart item added handler (0x0A0C)
	cm.baseParse.RegisterHandler("0A0C", "cart_item_added", "v V v C3 a8 C V2 a*",
		[]string{"index", "amount", "nameID", "identified", "damaged", "refine", "card", "type", "location", "wear_state", "bindOnEquipType", "options"},
		cm.HandleCartItemAdded)

	// Register cart item removed handler (0x0A0D)
	cm.baseParse.RegisterHandler("0A0D", "cart_item_removed", "v V",
		[]string{"index", "amount"},
		cm.HandleCartItemRemoved)

	// Register cart info handler (0x0121)
	cm.baseParse.RegisterHandler("0121", "cart_info", "v V V V",
		[]string{"items", "items_max", "weight", "weight_max"},
		cm.HandleCartInfo)

	// Register cart add failed handler (0x00B0)
	cm.baseParse.RegisterHandler("00B0", "cart_add_failed", "C",
		[]string{"fail"},
		cm.HandleCartAddFailed)
}

// HandleCartOff handles the cart_off packet (lines 3804-3807)
func (cm *CartManager) HandleCartOff(args map[string]interface{}) error {
	// Log the cart release
	cm.logger.Success("Cart released")

	// Call hooks
	cm.hooks.CallHook("cart_off", nil)

	return nil
}

// HandleCartItemsStackable handles the cart_items_stackable packet (lines 5120-5131)
func (cm *CartManager) HandleCartItemsStackable(args map[string]interface{}) error {
	// Extract packet data
	itemInfo, ok := args["itemInfo"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid itemInfo in cart_items_stackable packet")
	}

	// Log the stackable items
	cm.logger.Debug("Received %d stackable cart items", len(itemInfo))

	// Process each item
	for _, item := range itemInfo {
		ID, _ := item["ID"].(string)
		nameID, _ := item["nameID"].(uint16)
		amount, _ := item["amount"].(uint16)
		itemType, _ := item["type"].(uint8)
		identified, _ := item["identified"].(uint8)

		// Call hooks for each item
		cm.hooks.CallHook("packet_cart", map[string]interface{}{
			"ID":         ID,
			"nameID":     nameID,
			"amount":     amount,
			"type":       itemType,
			"identified": identified == 1,
		})
	}

	return nil
}

// HandleCartItemsNonstackable handles the cart_items_nonstackable packet (lines 5133-5144)
func (cm *CartManager) HandleCartItemsNonstackable(args map[string]interface{}) error {
	// Extract packet data
	itemInfo, ok := args["itemInfo"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid itemInfo in cart_items_nonstackable packet")
	}

	// Log the non-stackable items
	cm.logger.Debug("Received %d non-stackable cart items", len(itemInfo))

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
		cm.hooks.CallHook("packet_cart", map[string]interface{}{
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

// HandleCartItemAdded handles the cart_item_added packet (lines 5146-5174)
func (cm *CartManager) HandleCartItemAdded(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(string)
	if !ok {
		return fmt.Errorf("invalid ID in cart_item_added packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in cart_item_added packet")
	}

	nameID, ok := args["nameID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid nameID in cart_item_added packet")
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

	cm.logger.Info("Item added to cart: %s%s%s[%d] x%d",
		upgradeStr, identifiedStr, brokenStr, nameID, amount)

	// Call hooks
	cm.hooks.CallHook("packet_cart_add", map[string]interface{}{
		"ID":         ID,
		"nameID":     nameID,
		"amount":     amount,
		"identified": identified == 1,
		"broken":     broken == 1,
		"upgrade":    upgrade,
	})

	return nil
}

// HandleCartItemRemoved handles the cart_item_removed packet (lines 5176-5186)
func (cm *CartManager) HandleCartItemRemoved(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(string)
	if !ok {
		return fmt.Errorf("invalid ID in cart_item_removed packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in cart_item_removed packet")
	}

	// Call hooks
	cm.hooks.CallHook("packet_cart_remove", map[string]interface{}{
		"ID":     ID,
		"amount": amount,
	})

	return nil
}

// HandleCartInfo handles the cart_info packet (lines 5190-5194)
func (cm *CartManager) HandleCartInfo(args map[string]interface{}) error {
	// Extract packet data
	items, _ := args["items"].(uint16)
	itemsMax, _ := args["items_max"].(uint16)
	weight, _ := args["weight"].(uint32)
	weightMax, _ := args["weight_max"].(uint32)

	// Log the cart info
	cm.logger.Debug("Cart info: %d/%d items, %d/%d weight", items, itemsMax, weight, weightMax)

	// Call hooks
	cm.hooks.CallHook("cart_info", map[string]interface{}{
		"items":      items,
		"items_max":  itemsMax,
		"weight":     weight,
		"weight_max": weightMax,
	})

	return nil
}

// HandleCartAddFailed handles the cart_add_failed packet (lines 5196-5208)
func (cm *CartManager) HandleCartAddFailed(args map[string]interface{}) error {
	// Extract packet data
	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in cart_add_failed packet")
	}

	// Process based on fail code
	var reason string
	switch fail {
	case 0:
		reason = "overweight"
	case 1:
		reason = "too many items"
	default:
		reason = fmt.Sprintf("Unknown code %d", fail)
	}

	// Log the failure
	cm.logger.Error("Can't Add Cart Item (%s)", reason)

	// Call hooks
	cm.hooks.CallHook("cart_add_failed", map[string]interface{}{
		"fail":   fail,
		"reason": reason,
	})

	return nil
}
