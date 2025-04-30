package buyingstore

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// BuyingStoreItem represents an item in a buying store
type BuyingStoreItem struct {
	ID     uint32
	NameID uint16
	Type   uint8
	Amount uint16
	Price  uint32
	Name   string
}

// BuyingStoreInfo represents a buying store
type BuyingStoreInfo struct {
	ID    uint32
	Title string
	Items []*BuyingStoreItem
}

// BuyingStoreManager handles buying store-related packet handling
type BuyingStoreManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for buying store interactions
	selfBuyerItemList []*BuyingStoreItem
	buyerLists        map[uint32]*BuyingStoreInfo
	buyerListsID      []uint32
	buyerPriceLimit   uint32
	buyerID           uint32
	buyingStoreID     uint32
	buyerShopStarted  bool
}

// NewBuyingStoreManager creates a new buying store manager
func NewBuyingStoreManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *BuyingStoreManager {
	return &BuyingStoreManager{
		parser:            parser,
		hookManager:       hookManager,
		logger:            logger,
		selfBuyerItemList: make([]*BuyingStoreItem, 0),
		buyerLists:        make(map[uint32]*BuyingStoreInfo),
		buyerListsID:      make([]uint32, 0),
		buyerPriceLimit:   0,
		buyerID:           0,
		buyingStoreID:     0,
		buyerShopStarted:  false,
	}
}

// RegisterHandlers registers all buying store-related packet handlers
func (bsm *BuyingStoreManager) RegisterHandlers() {
	// Register open buying store handler
	bsm.parser.RegisterHandlerFunc("0810", "open_buying_store", "V",
		[]string{"amount"}, bsm.HandleOpenBuyingStore)

	// Register buyer items handler
	bsm.parser.RegisterHandlerFunc("0814", "buyer_items", "v V a*",
		[]string{"len", "venderID", "msg"}, bsm.HandleBuyerItems)

	// Register open buying store item list handler
	bsm.parser.RegisterHandlerFunc("0813", "open_buying_store_item_list", "v a*",
		[]string{"len", "RAW_MSG"}, bsm.HandleOpenBuyingStoreItemList)

	// Register buying store found handler
	bsm.parser.RegisterHandlerFunc("0816", "buying_store_found", "V Z80",
		[]string{"ID", "title"}, bsm.HandleBuyingStoreFound)

	// Register buying store lost handler
	bsm.parser.RegisterHandlerFunc("0817", "buying_store_lost", "V",
		[]string{"ID"}, bsm.HandleBuyingStoreLost)

	// Register buying store items list handler
	bsm.parser.RegisterHandlerFunc("0818", "buying_store_items_list", "V V V a*",
		[]string{"buyerID", "buyingStoreID", "zeny", "itemList"}, bsm.HandleBuyingStoreItemsList)

	// Register buying store item delete handler
	bsm.parser.RegisterHandlerFunc("081A", "buying_store_item_delete", "V v V",
		[]string{"ID", "amount", "zeny"}, bsm.HandleBuyingStoreItemDelete)

	// Register buying store fail handler
	bsm.parser.RegisterHandlerFunc("081B", "buying_store_fail", "W",
		[]string{"result"}, bsm.HandleBuyingStoreFail)

	// Register buying store update handler
	bsm.parser.RegisterHandlerFunc("081C", "buying_store_update", "W v",
		[]string{"itemID", "count"}, bsm.HandleBuyingStoreUpdate)

	// Register buyer found handler
	bsm.parser.RegisterHandlerFunc("083A", "buyer_found", "V Z80",
		[]string{"ID", "title"}, bsm.HandleBuyerFound)

	// Register buyer lost handler
	bsm.parser.RegisterHandlerFunc("083B", "buyer_lost", "V",
		[]string{"ID"}, bsm.HandleBuyerLost)

	// Register buying buy fail handler
	bsm.parser.RegisterHandlerFunc("083E", "buying_buy_fail", "W",
		[]string{"result"}, bsm.HandleBuyingBuyFail)

	// Register open buying store fail handler
	bsm.parser.RegisterHandlerFunc("0812", "open_buying_store_fail", "W",
		[]string{"result"}, bsm.HandleOpenBuyingStoreFail)

	// Register search store open handler
	bsm.parser.RegisterHandlerFunc("0835", "search_store_open", "C V",
		[]string{"type", "amount"}, bsm.HandleSearchStoreOpen)

	// Register search store fail handler
	bsm.parser.RegisterHandlerFunc("0837", "search_store_fail", "W",
		[]string{"reason"}, bsm.HandleSearchStoreFail)

	// Register search store result handler
	bsm.parser.RegisterHandlerFunc("0836", "search_store_result", "C C a*",
		[]string{"first_page", "has_next", "storeInfo"}, bsm.HandleSearchStoreResult)

	// Register search store pos handler
	bsm.parser.RegisterHandlerFunc("0838", "search_store_pos", "W W",
		[]string{"x", "y"}, bsm.HandleSearchStorePos)
}

// GetItemName gets the name of an item by its nameID
func (bsm *BuyingStoreManager) GetItemName(nameID uint16) string {
	// In a real implementation, this would look up the item name from a database
	// For now, we'll just return a placeholder
	return fmt.Sprintf("Item#%d", nameID)
}

// HandleOpenBuyingStore handles the open_buying_store packet
func (bsm *BuyingStoreManager) HandleOpenBuyingStore(args map[string]interface{}) error {
	// Extract packet data
	amount, ok := args["amount"].(uint32)
	if !ok {
		return fmt.Errorf("invalid amount in open_buying_store packet")
	}

	// Log message
	bsm.logger.Info("Your buying store can buy %d items", amount)

	// Call hook
	bsm.hookManager.CallHook("open_buying_store", map[string]interface{}{
		"amount": amount,
	})

	return nil
}

// HandleBuyerItems handles the buyer_items packet
func (bsm *BuyingStoreManager) HandleBuyerItems(args map[string]interface{}) error {
	// Extract packet data
	venderID, ok := args["venderID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid venderID in buyer_items packet")
	}

	msg, ok := args["msg"].([]byte)
	if !ok {
		return fmt.Errorf("invalid msg in buyer_items packet")
	}

	// Process buyer items
	headerLen := 12
	total := uint32(0)
	if len(msg) >= headerLen+4 {
		total = uint32(msg[headerLen]) | uint32(msg[headerLen+1])<<8 | uint32(msg[headerLen+2])<<16 | uint32(msg[headerLen+3])<<24
	}
	headerLen += 4

	items := make([]*BuyingStoreItem, 0)
	for i := headerLen; i < len(msg); i += 9 {
		if i+9 > len(msg) {
			break
		}

		price := uint32(msg[i]) | uint32(msg[i+1])<<8 | uint32(msg[i+2])<<16 | uint32(msg[i+3])<<24
		amount := uint16(msg[i+4]) | uint16(msg[i+5])<<8
		// Skip one byte (unused)
		nameID := uint16(msg[i+7]) | uint16(msg[i+8])<<8

		item := &BuyingStoreItem{
			Price:  price,
			Amount: amount,
			NameID: nameID,
			Name:   bsm.GetItemName(nameID),
		}

		items = append(items, item)
	}

	// Call hook
	bsm.hookManager.CallHook("buyer_items", map[string]interface{}{
		"venderID": venderID,
		"total":    total,
		"items":    items,
	})

	return nil
}

// HandleOpenBuyingStoreItemList handles the open_buying_store_item_list packet
func (bsm *BuyingStoreManager) HandleOpenBuyingStoreItemList(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleBuyingStoreFound handles the buying_store_found packet
func (bsm *BuyingStoreManager) HandleBuyingStoreFound(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleBuyingStoreLost handles the buying_store_lost packet
func (bsm *BuyingStoreManager) HandleBuyingStoreLost(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleBuyingStoreItemsList handles the buying_store_items_list packet
func (bsm *BuyingStoreManager) HandleBuyingStoreItemsList(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleBuyingStoreItemDelete handles the buying_store_item_delete packet
func (bsm *BuyingStoreManager) HandleBuyingStoreItemDelete(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleBuyingStoreFail handles the buying_store_fail packet
func (bsm *BuyingStoreManager) HandleBuyingStoreFail(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleBuyingStoreUpdate handles the buying_store_update packet
func (bsm *BuyingStoreManager) HandleBuyingStoreUpdate(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleBuyerFound handles the buyer_found packet
func (bsm *BuyingStoreManager) HandleBuyerFound(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleBuyerLost handles the buyer_lost packet
func (bsm *BuyingStoreManager) HandleBuyerLost(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleBuyingBuyFail handles the buying_buy_fail packet
func (bsm *BuyingStoreManager) HandleBuyingBuyFail(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleOpenBuyingStoreFail handles the open_buying_store_fail packet
func (bsm *BuyingStoreManager) HandleOpenBuyingStoreFail(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleSearchStoreOpen handles the search_store_open packet
func (bsm *BuyingStoreManager) HandleSearchStoreOpen(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleSearchStoreFail handles the search_store_fail packet
func (bsm *BuyingStoreManager) HandleSearchStoreFail(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleSearchStoreResult handles the search_store_result packet
func (bsm *BuyingStoreManager) HandleSearchStoreResult(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}

// HandleSearchStorePos handles the search_store_pos packet
func (bsm *BuyingStoreManager) HandleSearchStorePos(args map[string]interface{}) error {
	// Implementation will be added in the next part
	return nil
}
