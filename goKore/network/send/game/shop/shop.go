// Package shop provides shop-related packet sending functionality.
package shop

import (
	"encoding/binary"
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// ShopManager handles shop-related packet sending.
type ShopManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewShopManager creates a new shop manager.
func NewShopManager(baseSend core.Send) *ShopManager {
	return &ShopManager{
		baseSend: baseSend,
	}
}

// ParseBuyBulkVender parses the buy bulk vender packet.
// This is equivalent to the parse_buy_bulk_vender function in Send.pm.
func ParseBuyBulkVender(args map[string]interface{}) {
	itemInfo, ok := args["itemInfo"].([]byte)
	if !ok {
		return
	}

	items := []map[string]interface{}{}
	for i := 0; i < len(itemInfo); i += 4 {
		if i+4 > len(itemInfo) {
			break
		}

		amount := binary.LittleEndian.Uint16(itemInfo[i : i+2])
		itemIndex := binary.LittleEndian.Uint16(itemInfo[i+2 : i+4])

		items = append(items, map[string]interface{}{
			"amount":    amount,
			"itemIndex": itemIndex,
		})
	}

	args["items"] = items
}

// ReconstructBuyBulkVender reconstructs the buy bulk vender packet.
// This is equivalent to the reconstruct_buy_bulk_vender function in Send.pm.
func ReconstructBuyBulkVender(args map[string]interface{}) {
	items, ok := args["items"].([]map[string]interface{})
	if !ok {
		return
	}

	itemInfo := []byte{}
	for _, item := range items {
		amount, ok1 := item["amount"].(uint16)
		itemIndex, ok2 := item["itemIndex"].(uint16)
		if !ok1 || !ok2 {
			continue
		}

		amountBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(amountBytes, amount)

		itemIndexBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(itemIndexBytes, itemIndex)

		itemInfo = append(itemInfo, amountBytes...)
		itemInfo = append(itemInfo, itemIndexBytes...)
	}

	args["itemInfo"] = itemInfo
}

// SendBuyBulkVender sends a request to buy items from a vender.
// This is equivalent to the sendBuyBulkVender function in Send.pm.
func (sm *ShopManager) SendBuyBulkVender(venderID uint32, items []map[string]interface{}, venderCID uint32) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("buy_bulk_vender")
	if !exists {
		return fmt.Errorf("buy_bulk_vender packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"venderID":  venderID,
		"venderCID": venderCID,
		"items":     items,
	}

	// Reconstruct the item info
	ReconstructBuyBulkVender(args)

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// ReconstructBuyBulkBuyer reconstructs the buy bulk buyer packet.
// This is equivalent to the reconstruct_buy_bulk_buyer function in Send.pm.
func ReconstructBuyBulkBuyer(args map[string]interface{}) {
	items, ok := args["items"].([]map[string]interface{})
	if !ok {
		return
	}

	itemInfo := []byte{}
	for _, item := range items {
		ID, ok1 := item["ID"].([]byte)
		itemID, ok2 := item["itemID"].(uint16)
		amount, ok3 := item["amount"].(uint16)
		if !ok1 || !ok2 || !ok3 {
			continue
		}

		itemIDBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(itemIDBytes, itemID)

		amountBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(amountBytes, amount)

		itemInfo = append(itemInfo, ID...)
		itemInfo = append(itemInfo, itemIDBytes...)
		itemInfo = append(itemInfo, amountBytes...)
	}

	args["itemInfo"] = itemInfo
}

// SendBuyBulkBuyer sends a request to buy items from a buyer.
// This is equivalent to the sendBuyBulkBuyer function in Send.pm.
func (sm *ShopManager) SendBuyBulkBuyer(buyerID uint32, items []map[string]interface{}, buyingStoreID uint32) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("buy_bulk_buyer")
	if !exists {
		return fmt.Errorf("buy_bulk_buyer packet ID not found")
	}

	// Calculate the length
	len := 12 + (len(items) * 8)

	// Create the arguments
	args := map[string]interface{}{
		"len":           len,
		"buyerID":       buyerID,
		"buyingStoreID": buyingStoreID,
		"items":         items,
	}

	// Reconstruct the item info
	ReconstructBuyBulkBuyer(args)

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SendEnteringBuyer sends a request to enter a buyer's shop.
// This is equivalent to the sendEnteringBuyer function in Send.pm.
func (sm *ShopManager) SendEnteringBuyer(ID uint32) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("buy_bulk_request")
	if !exists {
		return fmt.Errorf("buy_bulk_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": ID,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// ReconstructBuyBulkOpenShop reconstructs the buy bulk open shop packet.
// This is equivalent to the reconstruct_buy_bulk_openShop function in Send.pm.
func ReconstructBuyBulkOpenShop(args map[string]interface{}) {
	items, ok := args["items"].([]map[string]interface{})
	if !ok {
		return
	}

	itemInfo := []byte{}
	for _, item := range items {
		nameID, ok1 := item["nameID"].(uint16)
		amount, ok2 := item["amount"].(uint16)
		price, ok3 := item["price"].(uint32)
		if !ok1 || !ok2 || !ok3 {
			continue
		}

		nameIDBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(nameIDBytes, nameID)

		amountBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(amountBytes, amount)

		priceBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(priceBytes, price)

		itemInfo = append(itemInfo, nameIDBytes...)
		itemInfo = append(itemInfo, amountBytes...)
		itemInfo = append(itemInfo, priceBytes...)
	}

	args["itemInfo"] = itemInfo
}

// SendBuyBulkOpenShop sends a request to open a buying shop.
// This is equivalent to the sendBuyBulkOpenShop function in Send.pm.
func (sm *ShopManager) SendBuyBulkOpenShop(limitZeny uint32, result uint8, storeName string, items []map[string]interface{}) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("buy_bulk_openShop")
	if !exists {
		return fmt.Errorf("buy_bulk_openShop packet ID not found")
	}

	// Calculate the length
	len := 89 + (len(items) * 8)

	// Create the arguments
	args := map[string]interface{}{
		"len":       len,
		"limitZeny": limitZeny,
		"result":    result,
		"storeName": storeName,
		"items":     items,
	}

	// Reconstruct the item info
	ReconstructBuyBulkOpenShop(args)

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SendCloseShop sends a request to close a vending shop.
// This is equivalent to the sendCloseShop function in Send.pm.
func (sm *ShopManager) SendCloseShop() error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("shop_close")
	if !exists {
		return fmt.Errorf("shop_close packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SendCloseBuyShop sends a request to close a buying shop.
// This is equivalent to the sendCloseBuyShop function in Send.pm.
func (sm *ShopManager) SendCloseBuyShop() error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("buy_bulk_closeShop")
	if !exists {
		return fmt.Errorf("buy_bulk_closeShop packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SendEnteringVender sends a request to enter a vender's shop.
// This is equivalent to the sendEnteringVender function in Send.pm.
func (sm *ShopManager) SendEnteringVender(accountID uint32) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("send_entering_vending")
	if !exists {
		return fmt.Errorf("send_entering_vending packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"accountID": accountID,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// ReconstructShopOpen reconstructs the shop open packet.
// This is equivalent to the reconstruct_shop_open function in Send.pm.
func ReconstructShopOpen(args map[string]interface{}) {
	items, ok := args["items"].([]map[string]interface{})
	if !ok {
		return
	}

	vendingInfo := []byte{}
	for _, item := range items {
		ID, ok1 := item["ID"].([]byte)
		amount, ok2 := item["amount"].(uint16)
		price, ok3 := item["price"].(uint32)
		if !ok1 || !ok2 || !ok3 || len(ID) != 2 {
			continue
		}

		amountBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(amountBytes, amount)

		priceBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(priceBytes, price)

		vendingInfo = append(vendingInfo, ID...)
		vendingInfo = append(vendingInfo, amountBytes...)
		vendingInfo = append(vendingInfo, priceBytes...)
	}

	args["vendingInfo"] = vendingInfo
}

// SendOpenShop sends a request to open a vending shop.
// This is equivalent to the sendOpenShop function in Send.pm.
func (sm *ShopManager) SendOpenShop(title string, items []map[string]interface{}) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("shop_open")
	if !exists {
		return fmt.Errorf("shop_open packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"title":  []byte(title), // stringToBytes in the original
		"result": 1,
		"items":  items,
	}

	// Reconstruct the vending info
	ReconstructShopOpen(args)

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}
