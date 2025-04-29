// Package buyingstore provides buying store-related packet sending functionality.
package buyingstore

import (
	"encoding/binary"
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// BuyingStoreManager handles buying store-related packet sending.
type BuyingStoreManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewBuyingStoreManager creates a new buying store manager.
func NewBuyingStoreManager(baseSend core.Send) *BuyingStoreManager {
	return &BuyingStoreManager{
		baseSend: baseSend,
	}
}

// SendBuyBulk sends a request to buy multiple items at once.
// This is equivalent to the sendBuyBulk function in Send.pm.
func (bsm *BuyingStoreManager) SendBuyBulk(items []map[string]interface{}) error {
	// Get the packet ID
	packetID, exists := bsm.baseSend.GetPacketID("buy_bulk")
	if !exists {
		return fmt.Errorf("buy_bulk packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"items": items,
	}

	// Reconstruct the buyInfo field
	if err := bsm.ReconstructBuyBulk(args); err != nil {
		return err
	}

	// Construct and send the packet
	packet, err := bsm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bsm.baseSend.SendToServer(packet)
}

// ReconstructBuyBulk reconstructs the buyInfo field for the buy_bulk packet.
// This is equivalent to the reconstruct_buy_bulk function in Send.pm.
func (bsm *BuyingStoreManager) ReconstructBuyBulk(args map[string]interface{}) error {
	items, ok := args["items"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("items not found or not a slice of maps")
	}

	// Each item is 4 bytes: 2 bytes for amount, 2 bytes for itemID
	buyInfo := make([]byte, len(items)*4)

	for i, item := range items {
		amount, ok := item["amount"].(uint16)
		if !ok {
			return fmt.Errorf("amount not found or not a uint16 for item %d", i)
		}

		itemID, ok := item["itemID"].(uint16)
		if !ok {
			return fmt.Errorf("itemID not found or not a uint16 for item %d", i)
		}

		binary.LittleEndian.PutUint16(buyInfo[i*4:i*4+2], amount)
		binary.LittleEndian.PutUint16(buyInfo[i*4+2:i*4+4], itemID)
	}

	args["buyInfo"] = buyInfo
	return nil
}

// SendSellBulk sends a request to sell multiple items at once.
// This is equivalent to the sendSellBulk function in Send.pm.
func (bsm *BuyingStoreManager) SendSellBulk(items []map[string]interface{}) error {
	// Get the packet ID
	packetID, exists := bsm.baseSend.GetPacketID("sell_bulk")
	if !exists {
		return fmt.Errorf("sell_bulk packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"items": items,
	}

	// Reconstruct the sellInfo field
	if err := bsm.ReconstructSellBulk(args); err != nil {
		return err
	}

	// Construct and send the packet
	packet, err := bsm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bsm.baseSend.SendToServer(packet)
}

// ReconstructSellBulk reconstructs the sellInfo field for the sell_bulk packet.
// This is equivalent to the reconstruct_sell_bulk function in Send.pm.
func (bsm *BuyingStoreManager) ReconstructSellBulk(args map[string]interface{}) error {
	items, ok := args["items"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("items not found or not a slice of maps")
	}

	// Each item is 4 bytes: 2 bytes for ID, 2 bytes for amount
	sellInfo := make([]byte, len(items)*4)

	for i, item := range items {
		id, ok := item["ID"].([]byte)
		if !ok {
			return fmt.Errorf("ID not found or not a byte slice for item %d", i)
		}
		if len(id) != 2 {
			return fmt.Errorf("ID must be 2 bytes for item %d", i)
		}

		amount, ok := item["amount"].(uint16)
		if !ok {
			return fmt.Errorf("amount not found or not a uint16 for item %d", i)
		}

		copy(sellInfo[i*4:i*4+2], id)
		binary.LittleEndian.PutUint16(sellInfo[i*4+2:i*4+4], amount)
	}

	args["sellInfo"] = sellInfo
	return nil
}

// SendSearchStoreClose sends a request to close the search store window.
// This is equivalent to the sendSearchStoreClose function in Send.pm.
func (bsm *BuyingStoreManager) SendSearchStoreClose() error {
	// Get the packet ID
	packetID, exists := bsm.baseSend.GetPacketID("search_store_close")
	if !exists {
		return fmt.Errorf("search_store_close packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := bsm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bsm.baseSend.SendToServer(packet)
}

// SendSearchStoreSearch sends a request to search for items in stores.
// This is equivalent to the sendSearchStoreSearch function in Send.pm.
func (bsm *BuyingStoreManager) SendSearchStoreSearch(searchType uint8, maxPrice, minPrice uint32, itemList, cardList []uint16) error {
	// Get the packet ID
	packetID, exists := bsm.baseSend.GetPacketID("search_store_info")
	if !exists {
		return fmt.Errorf("search_store_info packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type":      searchType,
		"max_price": maxPrice,
		"min_price": minPrice,
		"item_list": itemList,
		"card_list": cardList,
	}

	// Reconstruct the item_card_list field
	if err := bsm.ReconstructSearchStoreInfo(args); err != nil {
		return err
	}

	// Construct and send the packet
	packet, err := bsm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bsm.baseSend.SendToServer(packet)
}

// ReconstructSearchStoreInfo reconstructs the item_card_list field for the search_store_info packet.
// This is equivalent to the reconstruct_search_store_info function in Send.pm.
func (bsm *BuyingStoreManager) ReconstructSearchStoreInfo(args map[string]interface{}) error {
	itemList, ok := args["item_list"].([]uint16)
	if !ok {
		return fmt.Errorf("item_list not found or not a slice of uint16")
	}

	cardList, ok := args["card_list"].([]uint16)
	if !ok {
		return fmt.Errorf("card_list not found or not a slice of uint16")
	}

	// Set the item and card counts
	args["item_count"] = len(itemList)
	args["card_count"] = len(cardList)

	// Combine the item and card lists
	idList := append(itemList, cardList...)

	// Each ID is 2 bytes
	itemCardList := make([]byte, len(idList)*2)

	for i, id := range idList {
		binary.LittleEndian.PutUint16(itemCardList[i*2:i*2+2], id)
	}

	args["item_card_list"] = itemCardList
	return nil
}

// SendSearchStoreRequestNextPage sends a request to get the next page of search results.
// This is equivalent to the sendSearchStoreRequestNextPage function in Send.pm.
func (bsm *BuyingStoreManager) SendSearchStoreRequestNextPage() error {
	// Get the packet ID
	packetID, exists := bsm.baseSend.GetPacketID("search_store_request_next_page")
	if !exists {
		return fmt.Errorf("search_store_request_next_page packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := bsm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bsm.baseSend.SendToServer(packet)
}

// SendSearchStoreSelect sends a request to select an item from the search results.
// This is equivalent to the sendSearchStoreSelect function in Send.pm.
func (bsm *BuyingStoreManager) SendSearchStoreSelect(accountID, storeID uint32, nameID uint16) error {
	// Get the packet ID
	packetID, exists := bsm.baseSend.GetPacketID("search_store_select")
	if !exists {
		return fmt.Errorf("search_store_select packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"accountID": accountID,
		"storeID":   storeID,
		"nameID":    nameID,
	}

	// Construct and send the packet
	packet, err := bsm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bsm.baseSend.SendToServer(packet)
}
