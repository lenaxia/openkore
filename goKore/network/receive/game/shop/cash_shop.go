package shop

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// CashShopManager handles cash shop-related packet handling
type CashShopManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for cash shop interactions
	cashShopList  map[string]interface{}
	cashShopItems []map[string]interface{}
	cashPoints    map[string]uint32
	mergeItemList map[uint16]map[string]interface{}
}

// NewCashShopManager creates a new cash shop manager
func NewCashShopManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *CashShopManager {
	return &CashShopManager{
		parser:        parser,
		hookManager:   hookManager,
		logger:        logger,
		cashShopList:  make(map[string]interface{}),
		cashShopItems: make([]map[string]interface{}, 0),
		cashPoints:    make(map[string]uint32),
		mergeItemList: make(map[uint16]map[string]interface{}),
	}
}

// RegisterHandlers registers all cash shop-related packet handlers
func (csm *CashShopManager) RegisterHandlers() {
	// Register cash shop list handler
	csm.parser.RegisterHandlerFunc("0287", "cash_shop_list", "v C Z*",
		[]string{"len", "tabcode", "itemInfo"}, csm.HandleCashShopList)

	// Register cash shop open result handler
	csm.parser.RegisterHandlerFunc("0289", "cash_shop_open_result", "V V",
		[]string{"cash_points", "kafra_points"}, csm.HandleCashShopOpenResult)

	// Register cash shop buy result handler
	csm.parser.RegisterHandlerFunc("0289", "cash_shop_buy_result", "C V V",
		[]string{"result", "item_id", "updated_points"}, csm.HandleCashShopBuyResult)

	// Register cash dealer handler
	csm.parser.RegisterHandlerFunc("0287", "cash_dealer", "v V a*",
		[]string{"len", "cash_points", "item_list"}, csm.HandleCashDealer)

	// Register cash buy fail handler
	csm.parser.RegisterHandlerFunc("0289", "cash_buy_fail", "V V C",
		[]string{"cash_points", "kafra_points", "fail"}, csm.HandleCashBuyFail)

	// Register merge item open handler
	csm.parser.RegisterHandlerFunc("096D", "merge_item_open", "v Z*",
		[]string{"len", "itemList"}, csm.HandleMergeItemOpen)

	// Register merge item result handler
	csm.parser.RegisterHandlerFunc("096F", "merge_item_result", "W W C",
		[]string{"itemIndex", "total", "result"}, csm.HandleMergeItemResult)
}

// HandleCashShopList handles the cash_shop_list packet (lines 4443-4481)
func (csm *CashShopManager) HandleCashShopList(args map[string]interface{}) error {
	// Extract packet data
	tabcode, ok := args["tabcode"].(uint8)
	if !ok {
		return fmt.Errorf("invalid tabcode in cash_shop_list packet")
	}

	itemInfo, ok := args["itemInfo"].([]byte)
	if !ok {
		return fmt.Errorf("invalid itemInfo in cash_shop_list packet")
	}

	// Define cash shop tab names
	cashItemTabs := map[uint8]string{
		0: "New",
		1: "Popular",
		2: "Limited",
		3: "Rental",
		4: "Perpetuity",
		5: "Buff",
		6: "Recovery",
		7: "Etc",
	}

	// Get tab name
	tabName, ok := cashItemTabs[tabcode]
	if !ok {
		tabName = fmt.Sprintf("Unknown(%d)", tabcode)
	}

	// Log tab name
	csm.logger.Info("Cash Shop Tab: %s", tabName)

	// Initialize tab in cashShopList if it doesn't exist
	if _, ok := csm.cashShopList[fmt.Sprintf("%d", tabcode)]; !ok {
		csm.cashShopList[fmt.Sprintf("%d", tabcode)] = make([]map[string]interface{}, 0)
	}

	// Parse item info
	// The item_pack format is typically 'v V' (item ID and price)
	// Each item entry is 6 bytes (2 bytes for ID, 4 bytes for price)
	itemLen := 6
	for i := 0; i < len(itemInfo); i += itemLen {
		if i+itemLen > len(itemInfo) {
			break // Prevent out of bounds access
		}

		// Extract item ID and price
		itemID := uint16(itemInfo[i]) | uint16(itemInfo[i+1])<<8
		price := uint32(itemInfo[i+2]) | uint32(itemInfo[i+3])<<8 | uint32(itemInfo[i+4])<<16 | uint32(itemInfo[i+5])<<24

		// Create item data
		item := map[string]interface{}{
			"item_id": itemID,
			"price":   price,
		}

		// Add to cash shop list for this tab
		tabItems := csm.cashShopList[fmt.Sprintf("%d", tabcode)].([]map[string]interface{})
		csm.cashShopList[fmt.Sprintf("%d", tabcode)] = append(tabItems, item)

		// Log item info
		csm.logger.Info("Cash Shop Item: ID %d, Price %d", itemID, price)
	}

	// Call hook
	csm.hookManager.CallHook("cash_shop_list", map[string]interface{}{
		"tabcode":  tabcode,
		"tabName":  tabName,
		"itemInfo": itemInfo,
	})

	return nil
}

// HandleCashShopOpenResult handles the cash_shop_open_result packet (lines 4483-4491)
func (csm *CashShopManager) HandleCashShopOpenResult(args map[string]interface{}) error {
	// Extract packet data
	cashPoints, ok := args["cash_points"].(uint32)
	if !ok {
		return fmt.Errorf("invalid cash_points in cash_shop_open_result packet")
	}

	kafraPoints, ok := args["kafra_points"].(uint32)
	if !ok {
		return fmt.Errorf("invalid kafra_points in cash_shop_open_result packet")
	}

	// Store points
	csm.cashPoints["cash"] = cashPoints
	csm.cashPoints["kafra"] = kafraPoints

	// Log points
	csm.logger.Info("Cash Points: %d - Kafra Points: %d", cashPoints, kafraPoints)

	// Call hook
	csm.hookManager.CallHook("cash_shop_open_result", map[string]interface{}{
		"cash_points":  cashPoints,
		"kafra_points": kafraPoints,
	})

	return nil
}

// HandleCashShopBuyResult handles the cash_shop_buy_result packet (lines 4493-4526)
func (csm *CashShopManager) HandleCashShopBuyResult(args map[string]interface{}) error {
	// Extract packet data
	result, ok := args["result"].(uint8)
	if !ok {
		return fmt.Errorf("invalid result in cash_shop_buy_result packet")
	}

	itemID, ok := args["item_id"].(uint32)
	if !ok {
		return fmt.Errorf("invalid item_id in cash_shop_buy_result packet")
	}

	updatedPoints, ok := args["updated_points"].(uint32)
	if !ok {
		return fmt.Errorf("invalid updated_points in cash_shop_buy_result packet")
	}

	// Define result messages
	resultMessages := map[uint8]string{
		0:  "Success",
		1:  "Wrong Tab",
		2:  "Shortage cash",
		3:  "Unknown item",
		4:  "Inventory weight",
		5:  "Inventory item count",
		9:  "Rune overcount",
		10: "Eachitem overcount",
		11: "Unknown",
		12: "Busy",
	}

	// Get result message
	resultMsg, ok := resultMessages[result]
	if !ok {
		resultMsg = fmt.Sprintf("Unknown(%d)", result)
	}

	// Handle result
	if result > 0 {
		// Error occurred
		csm.logger.Error("Error while buying item ID %d from cash shop. Error code: %d (%s)", itemID, result, resultMsg)
	} else {
		// Success
		csm.logger.Info("Bought item ID %d from cash shop. Current CASH: %d", itemID, updatedPoints)

		// Update cash points
		csm.cashPoints["cash"] = updatedPoints
	}

	// Debug log
	csm.logger.Debug("Got result ID [%d] while buying item ID %d from CASH Shop. Current CASH: %d", result, itemID, updatedPoints)

	// Call hook
	csm.hookManager.CallHook("cash_shop_buy_result", map[string]interface{}{
		"result":         result,
		"resultMessage":  resultMsg,
		"item_id":        itemID,
		"updated_points": updatedPoints,
	})

	return nil
}

// HandleCashDealer handles the cash_dealer packet (lines 10008-10040)
func (csm *CashShopManager) HandleCashDealer(args map[string]interface{}) error {
	// Extract packet data
	cashPoints, ok := args["cash_points"].(uint32)
	if !ok {
		return fmt.Errorf("invalid cash_points in cash_dealer packet")
	}

	itemList, ok := args["item_list"].([]byte)
	if !ok {
		return fmt.Errorf("invalid item_list in cash_dealer packet")
	}

	// Extract kafra points if available
	var kafraPoints uint32
	if kafraPointsArg, ok := args["kafra_points"].(uint32); ok {
		kafraPoints = kafraPointsArg
	}

	// Clear cash list
	csm.cashShopItems = make([]map[string]interface{}, 0)

	// Log points
	csm.logger.Info("Cash Point: %d, Kafra Points: %d", cashPoints, kafraPoints)

	// Parse item list
	// Each item entry is 11 bytes:
	// - price (4 bytes)
	// - price_discount (4 bytes)
	// - type (1 byte)
	// - nameid (2 bytes)
	itemLen := 11
	for i := 0; i < len(itemList); i += itemLen {
		if i+itemLen > len(itemList) {
			break // Prevent out of bounds access
		}

		// Extract item data
		price := uint32(itemList[i]) | uint32(itemList[i+1])<<8 | uint32(itemList[i+2])<<16 | uint32(itemList[i+3])<<24
		priceDiscount := uint32(itemList[i+4]) | uint32(itemList[i+5])<<8 | uint32(itemList[i+6])<<16 | uint32(itemList[i+7])<<24
		itemType := uint8(itemList[i+8])
		nameID := uint16(itemList[i+9]) | uint16(itemList[i+10])<<8

		// Create item data
		item := map[string]interface{}{
			"price":          price,
			"price_discount": priceDiscount,
			"type":           itemType,
			"nameID":         nameID,
			"ID":             len(csm.cashShopItems),
		}

		// Add to cash shop items
		csm.cashShopItems = append(csm.cashShopItems, item)

		// Log item info
		csm.logger.Info("Cash Shop Item: ID %d, Type %d, Price %d, Discount %d", nameID, itemType, price, priceDiscount)
	}

	// Call hook
	csm.hookManager.CallHook("cash_dealer", map[string]interface{}{
		"cash_points":  cashPoints,
		"kafra_points": kafraPoints,
		"items":        csm.cashShopItems,
	})

	return nil
}

// HandleCashBuyFail handles the cash_buy_fail packet (lines 10614-10617)
func (csm *CashShopManager) HandleCashBuyFail(args map[string]interface{}) error {
	// Extract packet data
	cashPoints, ok := args["cash_points"].(uint32)
	if !ok {
		return fmt.Errorf("invalid cash_points in cash_buy_fail packet")
	}

	kafraPoints, ok := args["kafra_points"].(uint32)
	if !ok {
		return fmt.Errorf("invalid kafra_points in cash_buy_fail packet")
	}

	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in cash_buy_fail packet")
	}

	// Log failure
	csm.logger.Debug("Cash buy failed: cash_points=%d kafra_points=%d fail=%d", cashPoints, kafraPoints, fail)

	// Call hook
	csm.hookManager.CallHook("cash_buy_fail", map[string]interface{}{
		"cash_points":  cashPoints,
		"kafra_points": kafraPoints,
		"fail":         fail,
	})

	return nil
}

// HandleMergeItemOpen handles the merge_item_open packet (lines 10046-10061)
func (csm *CashShopManager) HandleMergeItemOpen(args map[string]interface{}) error {
	// Extract packet data
	itemList, ok := args["itemList"].([]byte)
	if !ok {
		return fmt.Errorf("invalid itemList in merge_item_open packet")
	}

	// Clear merge item list
	csm.mergeItemList = make(map[uint16]map[string]interface{})

	// Parse item IDs from itemList
	// Each item ID is 2 bytes
	var itemIDs []uint16
	for i := 0; i < len(itemList); i += 2 {
		if i+2 > len(itemList) {
			break // Prevent out of bounds access
		}

		itemID := uint16(itemList[i]) | uint16(itemList[i+1])<<8
		itemIDs = append(itemIDs, itemID)
	}

	// Log number of items
	csm.logger.Debug("Enable to merge %d items", len(itemIDs))

	// Call hook with parsed item IDs
	csm.hookManager.CallHook("merge_item_open", map[string]interface{}{
		"itemIDs": itemIDs,
	})

	// Log message to user
	csm.logger.Info("Received %d items that can be merged. Use 'merge' to continue", len(itemIDs))

	return nil
}

// HandleMergeItemResult handles the merge_item_result packet (lines 10072-10093)
func (csm *CashShopManager) HandleMergeItemResult(args map[string]interface{}) error {
	// Extract packet data
	itemIndex, ok := args["itemIndex"].(uint16)
	if !ok {
		return fmt.Errorf("invalid itemIndex in merge_item_result packet")
	}

	total, ok := args["total"].(uint16)
	if !ok {
		return fmt.Errorf("invalid total in merge_item_result packet")
	}

	result, ok := args["result"].(uint8)
	if !ok {
		return fmt.Errorf("invalid result in merge_item_result packet")
	}

	// Handle result
	if result == 0 {
		// Success
		csm.logger.Info("Items were merged successfully!")
		csm.logger.Info("Updated amount of item index %d: new amount: %d", itemIndex, total)
	} else if result == 1 {
		// Cannot merge
		csm.logger.Error("Items cannot be merged.")
	} else if result == 2 {
		// Exceed stack limit
		csm.logger.Error("The amount of merged item will exceed stack limit.")
	} else {
		// Unknown error
		csm.logger.Error("An error occurred while merging items. Error: %d", result)
	}

	// Debug log
	csm.logger.Debug("Merge item result: itemIndex:%d total:%d result:%d", itemIndex, total, result)

	// Call hook
	csm.hookManager.CallHook("merge_item_result", map[string]interface{}{
		"itemIndex": itemIndex,
		"total":     total,
		"result":    result,
	})

	return nil
}
