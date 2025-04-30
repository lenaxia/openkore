package shop

import (
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// ShopManager handles shop-related packet handling
type ShopManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for shop interactions
	storeList     map[string]interface{}
	storeItems    []map[string]interface{}
	npcTalkState  map[string]interface{}
	venderLists   map[uint32]map[string]interface{}
	venderListsID []uint32
}

// NewShopManager creates a new shop manager
func NewShopManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *ShopManager {
	return &ShopManager{
		parser:        parser,
		hookManager:   hookManager,
		logger:        logger,
		storeList:     make(map[string]interface{}),
		storeItems:    make([]map[string]interface{}, 0),
		npcTalkState:  make(map[string]interface{}),
		venderLists:   make(map[uint32]map[string]interface{}),
		venderListsID: make([]uint32, 0),
	}
}

// RegisterHandlers registers all shop-related packet handlers
func (sm *ShopManager) RegisterHandlers() {
	// Register NPC store begin handler
	sm.parser.RegisterHandlerFunc("00C4", "npc_store_begin", "V",
		[]string{"ID"}, sm.HandleNpcStoreBegin)

	// Register NPC store info handler
	sm.parser.RegisterHandlerFunc("00C6", "npc_store_info", "v Z*",
		[]string{"len", "RAW_MSG"}, sm.HandleNpcStoreInfo)

	// Register NPC sell list handler
	sm.parser.RegisterHandlerFunc("00C7", "npc_sell_list", "v Z*",
		[]string{"len", "itemsdata"}, sm.HandleNpcSellList)

	// Register buy result handler
	sm.parser.RegisterHandlerFunc("00CA", "buy_result", "C",
		[]string{"fail"}, sm.HandleBuyResult)

	// Register sell result handler
	sm.parser.RegisterHandlerFunc("00CB", "sell_result", "C",
		[]string{"fail"}, sm.HandleSellResult)

	// Register vending start handler
	sm.parser.RegisterHandlerFunc("0136", "vending_start", "v Z*",
		[]string{"len", "itemList"}, sm.HandleVendingStart)

	// Register vender items list handler
	sm.parser.RegisterHandlerFunc("0133", "vender_items_list", "V V v Z*",
		[]string{"venderID", "venderCID", "len", "itemList"}, sm.HandleVenderItemsList)

	// Register vender found handler
	sm.parser.RegisterHandlerFunc("0131", "vender_found", "V Z80",
		[]string{"ID", "title"}, sm.HandleVenderFound)

	// Register vender lost handler
	sm.parser.RegisterHandlerFunc("0132", "vender_lost", "V",
		[]string{"ID"}, sm.HandleVenderLost)

	// Register vender buy fail handler
	sm.parser.RegisterHandlerFunc("00D5", "vender_buy_fail", "C v V",
		[]string{"fail", "amount", "ID"}, sm.HandleVenderBuyFail)

	// Register shop skill handler
	sm.parser.RegisterHandlerFunc("012D", "shop_skill", "v",
		[]string{"number"}, sm.HandleShopSkill)

	// Register shop sold handler
	sm.parser.RegisterHandlerFunc("0137", "shop_sold", "v v",
		[]string{"number", "amount"}, sm.HandleShopSold)

	// Register shop sold long handler
	sm.parser.RegisterHandlerFunc("0136", "shop_sold_long", "v v V V",
		[]string{"number", "amount", "zeny", "charID"}, sm.HandleShopSoldLong)

	// Register open store status handler
	sm.parser.RegisterHandlerFunc("012E", "open_store_status", "C",
		[]string{"flag"}, sm.HandleOpenStoreStatus)
}

// HandleNpcStoreBegin handles the npc_store_begin packet (lines 7573-7581)
func (sm *ShopManager) HandleNpcStoreBegin(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in npc_store_begin packet")
	}

	// Clear talk state
	sm.npcTalkState = make(map[string]interface{})

	// Set NPC ID
	sm.npcTalkState["ID"] = id

	// Update NPC talk state
	sm.npcTalkState["talk"] = "buy_or_sell"
	sm.npcTalkState["time"] = getCurrentTime()

	// Get NPC name (in a real implementation, this would use a function to get the name from the ID)
	npcName := fmt.Sprintf("NPC-%d", id)

	// Store NPC name
	sm.storeList["npcName"] = npcName

	sm.logger.Info("Entered shop dialog with %s", npcName)

	return nil
}

// HandleNpcStoreInfo handles the npc_store_info packet (lines 7588-7633)
func (sm *ShopManager) HandleNpcStoreInfo(args map[string]interface{}) error {
	// Extract packet data
	rawMsg, ok := args["RAW_MSG"].([]byte)
	if !ok {
		return fmt.Errorf("invalid RAW_MSG in npc_store_info packet")
	}

	msgSize, ok := args["len"].(uint16)
	if !ok {
		return fmt.Errorf("invalid len in npc_store_info packet")
	}

	// Clear store list and talk state
	sm.storeList = make(map[string]interface{})
	sm.storeItems = make([]map[string]interface{}, 0)
	sm.npcTalkState = make(map[string]interface{})

	// Define packet format (simplified)
	// In a real implementation, this would depend on the server type
	packLen := 11 // Length of the packed data

	// Process store items
	for i := 4; i < int(msgSize); i += packLen {
		if i+packLen <= len(rawMsg) {
			// Extract item data (simplified)
			// In a real implementation, this would use proper unpacking
			price := uint32(rawMsg[i]) | uint32(rawMsg[i+1])<<8 | uint32(rawMsg[i+2])<<16 | uint32(rawMsg[i+3])<<24
			itemType := uint8(rawMsg[i+8])
			nameID := uint16(rawMsg[i+9]) | uint16(rawMsg[i+10])<<8

			// Create item entry
			item := map[string]interface{}{
				"price":  price,
				"type":   itemType,
				"nameID": nameID,
				"ID":     len(sm.storeItems),
				"name":   fmt.Sprintf("Item-%d", nameID), // In a real implementation, this would use a function to get the item name
			}

			// Add item to store items
			sm.storeItems = append(sm.storeItems, item)

			sm.logger.Debug("Item added to Store: %s - %dz", item["name"], price)
		}
	}

	// Update NPC talk state
	sm.npcTalkState["talk"] = "store"
	sm.npcTalkState["time"] = getCurrentTime()

	sm.logger.Info("Store information received with %d items", len(sm.storeItems))

	return nil
}

// HandleNpcSellList handles the npc_sell_list packet (lines 7637-7663)
func (sm *ShopManager) HandleNpcSellList(args map[string]interface{}) error {
	// Extract packet data
	itemsData, ok := args["itemsdata"].([]byte)
	if !ok {
		return fmt.Errorf("invalid itemsdata in npc_sell_list packet")
	}

	// Clear talk state
	sm.npcTalkState = make(map[string]interface{})

	sm.logger.Info("You can sell:")

	// Process sellable items
	for i := 0; i < len(itemsData); i += 10 {
		if i+10 <= len(itemsData) {
			// Extract item data (simplified)
			// In a real implementation, this would use proper unpacking
			index := uint16(itemsData[i]) | uint16(itemsData[i+1])<<8
			priceOvercharge := uint32(itemsData[i+6]) | uint32(itemsData[i+7])<<8 | uint32(itemsData[i+8])<<16 | uint32(itemsData[i+9])<<24

			// In a real implementation, we would get the item from the inventory
			itemName := fmt.Sprintf("Item-%d", index)
			itemAmount := 1

			sm.logger.Info("%d x %s for %dz each", itemAmount, itemName, priceOvercharge)
		}
	}

	// Update NPC talk state
	sm.npcTalkState["talk"] = "sell"
	sm.npcTalkState["time"] = getCurrentTime()

	sm.logger.Info("Ready to start selling items")

	return nil
}

// HandleBuyResult handles the buy_result packet (lines 7681-7704)
func (sm *ShopManager) HandleBuyResult(args map[string]interface{}) error {
	// Extract packet data
	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in buy_result packet")
	}

	// Process based on fail code
	switch fail {
	case 0:
		sm.logger.Info("Buy completed")
	case 1:
		sm.logger.Error("Buy failed (insufficient zeny)")
	case 2:
		sm.logger.Error("Buy failed (insufficient weight capacity)")
	case 3:
		sm.logger.Error("Buy failed (too many different inventory items)")
	case 4:
		sm.logger.Error("Buy failed (item does not exist in store)")
	case 5:
		sm.logger.Error("Buy failed (item cannot be exchanged)")
	case 6:
		sm.logger.Error("Buy failed (invalid store)")
	default:
		sm.logger.Error("Buy failed (failure code %d)", fail)
	}

	// Call hooks
	sm.hookManager.CallHook("buy_result", map[string]interface{}{
		"fail": fail,
	})

	return nil
}

// HandleSellResult handles the sell_result packet (lines 9519-9531)
func (sm *ShopManager) HandleSellResult(args map[string]interface{}) error {
	// Extract packet data
	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in sell_result packet")
	}

	// Process based on fail code
	if fail != 0 {
		sm.logger.Error("Sell failed")
	} else {
		sm.logger.Info("Sold items")
		sm.logger.Info("Sell completed")
	}

	// Call hooks
	sm.hookManager.CallHook("sell_result", map[string]interface{}{
		"fail": fail,
	})

	return nil
}

// HandleVendingStart handles the vending_start packet (lines 3908-3940)
func (sm *ShopManager) HandleVendingStart(args map[string]interface{}) error {
	// Extract packet data
	itemList, ok := args["itemList"].([]byte)
	if !ok {
		return fmt.Errorf("invalid itemList in vending_start packet")
	}

	listLen, ok := args["len"].(uint16)
	if !ok {
		return fmt.Errorf("invalid len in vending_start packet")
	}

	// In a real implementation, we would get the shop title from somewhere
	shopTitle := "My Shop"

	// Log shop opening
	sm.logger.Info("Shop '%s' opened!", shopTitle)

	// Clear articles
	sm.storeItems = make([]map[string]interface{}, 0)

	// Define item pack format (simplified)
	itemLen := 20 // Length of each item entry

	// Process item list
	sm.logger.Info("%s", centerString(" "+shopTitle+" ", 83, '-'))
	sm.logger.Info("#  Name                                       Type                     Price Amount")

	for i := 0; i < int(listLen); i += itemLen {
		if i+itemLen <= len(itemList) {
			// Extract item data (simplified)
			// In a real implementation, this would use proper unpacking
			price := uint32(itemList[i]) | uint32(itemList[i+1])<<8 | uint32(itemList[i+2])<<16 | uint32(itemList[i+3])<<24
			number := uint16(itemList[i+4]) | uint16(itemList[i+5])<<8
			quantity := uint16(itemList[i+6]) | uint16(itemList[i+7])<<8
			itemType := uint8(itemList[i+8])
			nameID := uint16(itemList[i+9]) | uint16(itemList[i+10])<<8

			// Create item entry
			item := map[string]interface{}{
				"price":    price,
				"number":   number,
				"quantity": quantity,
				"type":     itemType,
				"nameID":   nameID,
				"name":     fmt.Sprintf("Item-%d", nameID), // In a real implementation, this would use a function to get the item name
			}

			// Add item to store items
			sm.storeItems = append(sm.storeItems, item)

			// Log item
			sm.logger.Debug("Item added to Vender Store: %s - %d z", item["name"], price)

			// Format item info for display
			itemTypeName := fmt.Sprintf("Type-%d", itemType) // In a real implementation, this would use a function to get the item type name
			sm.logger.Info("%d %s %s %dz %d", len(sm.storeItems), item["name"], itemTypeName, price, quantity)
		}
	}

	sm.logger.Info("%s", strings.Repeat("-", 83))

	return nil
}

// Helper function to get current time
func getCurrentTime() int64 {
	return 0 // In a real implementation, this would return the current time
}

// HandleVenderItemsList handles the vender_items_list packet (lines 3942-3990)
func (sm *ShopManager) HandleVenderItemsList(args map[string]interface{}) error {
	// Extract packet data
	venderID, ok := args["venderID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid venderID in vender_items_list packet")
	}

	venderCID, ok := args["venderCID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid venderCID in vender_items_list packet")
	}

	itemList, ok := args["itemList"].([]byte)
	if !ok {
		return fmt.Errorf("invalid itemList in vender_items_list packet")
	}

	listLen, ok := args["len"].(uint16)
	if !ok {
		return fmt.Errorf("invalid len in vender_items_list packet")
	}

	// Get player name (in a real implementation, this would use a function to get the player name from the ID)
	playerName := fmt.Sprintf("Player-%d", venderID)

	// Clear vender item list
	sm.storeItems = make([]map[string]interface{}, 0)

	// Define item pack format (simplified)
	itemLen := 20 // Length of each item entry

	// Log vender items list
	sm.logger.Info("%s", centerString(" Vender: "+playerName+" ", 88, '-'))
	sm.logger.Info("#  Name                                      Type                           Price Amount")

	// Process item list
	for i := 0; i < int(listLen); i += itemLen {
		if i+itemLen <= len(itemList) {
			// Extract item data (simplified)
			// In a real implementation, this would use proper unpacking
			price := uint32(itemList[i]) | uint32(itemList[i+1])<<8 | uint32(itemList[i+2])<<16 | uint32(itemList[i+3])<<24
			amount := uint16(itemList[i+4]) | uint16(itemList[i+5])<<8
			id := uint16(itemList[i+6]) | uint16(itemList[i+7])<<8
			itemType := uint8(itemList[i+8])
			nameID := uint16(itemList[i+9]) | uint16(itemList[i+10])<<8

			// Create item entry
			item := map[string]interface{}{
				"price":  price,
				"amount": amount,
				"ID":     id,
				"type":   itemType,
				"nameID": nameID,
				"name":   fmt.Sprintf("Item-%d", nameID), // In a real implementation, this would use a function to get the item name
				"binID":  len(sm.storeItems),
			}

			// Add item to store items
			sm.storeItems = append(sm.storeItems, item)

			// Log item
			sm.logger.Debug("Item added to Vender Store: %s - %d z", item["name"], price)

			// Call hook for each item
			sm.hookManager.CallHook("packet_vender_store", map[string]interface{}{
				"item": item,
			})

			// Format item info for display
			itemTypeName := fmt.Sprintf("Type-%d", itemType) // In a real implementation, this would use a function to get the item type name
			sm.logger.Info("%d %s %s %dz %d", item["binID"], item["name"], itemTypeName, price, amount)
		}
	}

	sm.logger.Info("%s", strings.Repeat("-", 88))

	// Check for expire date
	expireDate := uint32(0)
	if expireDateVal, ok := args["expireDate"].(uint32); ok && expireDateVal > 0 {
		expireDate = expireDateVal
		// In a real implementation, this would format the date properly
		sm.logger.Info("Expire Date: %d", expireDate)
	}

	// Call hook for the entire list
	sm.hookManager.CallHook("packet_vender_store2", map[string]interface{}{
		"venderID":   venderID,
		"venderCID":  venderCID,
		"itemList":   sm.storeItems,
		"expireDate": expireDate,
	})

	return nil
}

// HandleVenderFound handles the vender_found packet (lines 11689-11702)
func (sm *ShopManager) HandleVenderFound(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in vender_found packet")
	}

	title, ok := args["title"].(string)
	if !ok {
		return fmt.Errorf("invalid title in vender_found packet")
	}

	// Check if vender already exists
	if _, exists := sm.venderLists[id]; !exists {
		// Add vender ID to list
		sm.venderListsID = append(sm.venderListsID, id)

		// Call hook
		sm.hookManager.CallHook("packet_vender", map[string]interface{}{
			"ID":    id,
			"title": title,
		})
	}

	// Store vender information
	sm.venderLists[id] = map[string]interface{}{
		"title": title,
		"id":    id,
	}

	sm.logger.Debug("Vender found: %s (ID: %d)", title, id)

	return nil
}

// HandleVenderLost handles the vender_lost packet (lines 11704-11710)
func (sm *ShopManager) HandleVenderLost(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in vender_lost packet")
	}

	// Remove vender ID from list
	for i, venderID := range sm.venderListsID {
		if venderID == id {
			sm.venderListsID = append(sm.venderListsID[:i], sm.venderListsID[i+1:]...)
			break
		}
	}

	// Delete vender information
	delete(sm.venderLists, id)

	sm.logger.Debug("Vender lost: %d", id)

	return nil
}

// HandleVenderBuyFail handles the vender_buy_fail packet (lines 9986-10002)
func (sm *ShopManager) HandleVenderBuyFail(args map[string]interface{}) error {
	// Extract packet data
	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in vender_buy_fail packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in vender_buy_fail packet")
	}

	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in vender_buy_fail packet")
	}

	// Process based on fail code
	switch fail {
	case 1:
		sm.logger.Error("Failed to buy %d of item #%d from vender (insufficient zeny) (error code %d)", amount, id, fail)
	case 2:
		sm.logger.Error("Failed to buy %d of item #%d from vender (overweight) (error code %d)", amount, id, fail)
	case 4:
		sm.logger.Error("Failed to buy %d of item #%d from vender (requested to purchase more than vender had in stock) (error code %d)", amount, id, fail)
	case 6:
		sm.logger.Error("Failed to buy %d of item #%d from vender (vender refreshed shop before purchase request) (error code %d)", amount, id, fail)
	case 8:
		sm.logger.Error("Failed to buy %d of item #%d from vender (vender would go over max zeny with the purchase) (error code %d)", amount, id, fail)
	default:
		sm.logger.Error("Failed to buy %d of item #%d from vender (unknown error code %d)", amount, id, fail)
	}

	return nil
}

// HandleOpenStoreStatus handles the open_store_status packet (lines 11839-11851)
func (sm *ShopManager) HandleOpenStoreStatus(args map[string]interface{}) error {
	// Extract packet data
	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in open_store_status packet")
	}

	// Process based on flag
	if flag == 0 {
		sm.logger.Info("Store set up successfully")

		// Call hook
		sm.hookManager.CallHook("open_store_success", nil)
	} else {
		sm.logger.Error("Failed setting up shop with error code %d", flag)

		// Call hook
		sm.hookManager.CallHook("open_store_fail", map[string]interface{}{
			"flag": flag,
		})
	}

	return nil
}

// HandleShopSkill handles the shop_skill packet (lines 3810-3816)
func (sm *ShopManager) HandleShopSkill(args map[string]interface{}) error {
	// Extract packet data
	number, ok := args["number"].(uint16)
	if !ok {
		return fmt.Errorf("invalid number in shop_skill packet")
	}

	// Log message
	sm.logger.Info("You can sell %d items!", number)

	return nil
}

// HandleShopSold handles the shop_sold packet (lines 3820-3859)
func (sm *ShopManager) HandleShopSold(args map[string]interface{}) error {
	// Extract packet data
	number, ok := args["number"].(uint16)
	if !ok {
		return fmt.Errorf("invalid number in shop_sold packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in shop_sold packet")
	}

	// Check if the article exists
	if int(number) >= len(sm.storeItems) {
		return fmt.Errorf("invalid article number %d", number)
	}

	// Get article
	article := sm.storeItems[number]

	// Update article
	if sold, ok := article["sold"].(uint16); ok {
		article["sold"] = sold + amount
	} else {
		article["sold"] = amount
	}

	// Calculate earned zeny
	price, _ := article["price"].(uint32)
	earned := uint32(amount) * price

	// Update shop earned
	shopEarned := uint32(0)
	if val, ok := sm.storeList["earned"].(uint32); ok {
		shopEarned = val
	}
	sm.storeList["earned"] = shopEarned + earned

	// Update article quantity
	if quantity, ok := article["quantity"].(uint16); ok {
		article["quantity"] = quantity - amount
	}

	// Log message
	itemName := "Unknown"
	if name, ok := article["name"].(string); ok {
		itemName = name
	}
	sm.logger.Info("Sold: %s x %d - %dz", itemName, amount, earned)

	// Call hook
	sm.hookManager.CallHook("vending_item_sold", map[string]interface{}{
		"vendShopIndex": number,
		"amount":        amount,
		"vendArticle":   article,
		"zenyEarned":    earned,
		"packetType":    "short",
	})

	// Check if item sold out
	if quantity, ok := article["quantity"].(uint16); ok && quantity < 1 {
		sm.logger.Info("Sold out: %s", itemName)

		// Call hook
		sm.hookManager.CallHook("vending_item_sold_out", map[string]interface{}{
			"vendShopIndex": number,
			"vendArticle":   article,
		})

		// Check if all items sold out
		allSoldOut := true
		for _, item := range sm.storeItems {
			if quantity, ok := item["quantity"].(uint16); ok && quantity > 0 {
				allSoldOut = false
				break
			}
		}

		if allSoldOut {
			sm.logger.Info("Items have been sold out.")
			// In a real implementation, this would close the shop
		}
	}

	return nil
}

// HandleShopSoldLong handles the shop_sold_long packet (lines 3861-3905)
func (sm *ShopManager) HandleShopSoldLong(args map[string]interface{}) error {
	// Extract packet data
	number, ok := args["number"].(uint16)
	if !ok {
		return fmt.Errorf("invalid number in shop_sold_long packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in shop_sold_long packet")
	}

	zeny, ok := args["zeny"].(uint32)
	if !ok {
		return fmt.Errorf("invalid zeny in shop_sold_long packet")
	}

	charID, ok := args["charID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid charID in shop_sold_long packet")
	}

	// Check if the article exists
	if int(number) >= len(sm.storeItems) {
		return fmt.Errorf("invalid article number %d", number)
	}

	// Get article
	article := sm.storeItems[number]

	// Update article
	if sold, ok := article["sold"].(uint16); ok {
		article["sold"] = sold + amount
	} else {
		article["sold"] = amount
	}

	// Update shop earned
	shopEarned := uint32(0)
	if val, ok := sm.storeList["earned"].(uint32); ok {
		shopEarned = val
	}
	sm.storeList["earned"] = shopEarned + zeny

	// Update article quantity
	if quantity, ok := article["quantity"].(uint16); ok {
		article["quantity"] = quantity - amount
	}

	// Log message
	itemName := "Unknown"
	if name, ok := article["name"].(string); ok {
		itemName = name
	}
	sm.logger.Info("Sold: %s x %d - %dz (Buyer charID: %d)", itemName, amount, zeny, charID)

	// Call hook
	sm.hookManager.CallHook("vending_item_sold", map[string]interface{}{
		"vendShopIndex": number,
		"amount":        amount,
		"vendArticle":   article,
		"buyerCharID":   charID,
		"zenyEarned":    zeny,
		"time":          getCurrentTime(),
		"packetType":    "long",
	})

	// Check if item sold out
	if quantity, ok := article["quantity"].(uint16); ok && quantity < 1 {
		sm.logger.Info("Sold out: %s", itemName)

		// Call hook
		sm.hookManager.CallHook("vending_item_sold_out", map[string]interface{}{
			"vendShopIndex": number,
			"vendArticle":   article,
		})

		// Check if all items sold out
		allSoldOut := true
		for _, item := range sm.storeItems {
			if quantity, ok := item["quantity"].(uint16); ok && quantity > 0 {
				allSoldOut = false
				break
			}
		}

		if allSoldOut {
			sm.logger.Info("Items have been sold out.")
			// In a real implementation, this would close the shop
		}
	}

	return nil
}

// Helper function to center a string
func centerString(s string, width int, fill byte) string {
	if len(s) >= width {
		return s
	}

	leftPad := (width - len(s)) / 2
	rightPad := width - len(s) - leftPad

	return strings.Repeat(string(fill), leftPad) + s + strings.Repeat(string(fill), rightPad)
}
