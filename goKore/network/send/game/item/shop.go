// Package item provides item-related packet sending functionality.
package item

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// VendingItem represents an item to be sold in a vending shop.
type VendingItem struct {
	Index  uint16
	Amount uint16
	Price  uint32
}

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

// OpenNpcShop sends a request to open an NPC shop.
func (sm *ShopManager) OpenNpcShop(npcID uint32) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("open_npc_shop")
	if !exists {
		return fmt.Errorf("open_npc_shop packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"npc_id": npcID,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// BuyItem sends a request to buy an item from an NPC shop.
func (sm *ShopManager) BuyItem(itemID, amount uint16) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("buy_item")
	if !exists {
		return fmt.Errorf("buy_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"item_id": itemID,
		"amount":  amount,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SellItem sends a request to sell an item to an NPC shop.
func (sm *ShopManager) SellItem(index, amount uint16) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("sell_item")
	if !exists {
		return fmt.Errorf("sell_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index":  index,
		"amount": amount,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// CloseShop sends a request to close an NPC shop.
func (sm *ShopManager) CloseShop() error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("close_shop")
	if !exists {
		return fmt.Errorf("close_shop packet ID not found")
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

// OpenVendingShop sends a request to open a vending shop.
func (sm *ShopManager) OpenVendingShop(title string, items []VendingItem) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("open_vending_shop")
	if !exists {
		return fmt.Errorf("open_vending_shop packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"title": title,
		"items": items,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// BuyVendingItem sends a request to buy an item from a vending shop.
func (sm *ShopManager) BuyVendingItem(vendorID uint32, index, amount uint16) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("buy_vending_item")
	if !exists {
		return fmt.Errorf("buy_vending_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"vendor_id": vendorID,
		"index":     index,
		"amount":    amount,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// CloseVendingShop sends a request to close a vending shop.
func (sm *ShopManager) CloseVendingShop() error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("close_vending_shop")
	if !exists {
		return fmt.Errorf("close_vending_shop packet ID not found")
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

// SearchVendingShops sends a request to search for vending shops.
func (sm *ShopManager) SearchVendingShops(itemID uint16) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("search_vending_shops")
	if !exists {
		return fmt.Errorf("search_vending_shops packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"item_id": itemID,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// OpenBuyingStore sends a request to open a buying store.
func (sm *ShopManager) OpenBuyingStore(title string, zeny uint32, items []VendingItem) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("open_buying_store")
	if !exists {
		return fmt.Errorf("open_buying_store packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"title": title,
		"zeny":  zeny,
		"items": items,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SellToBuyingStore sends a request to sell an item to a buying store.
func (sm *ShopManager) SellToBuyingStore(buyerID uint32, index, amount uint16) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("sell_to_buying_store")
	if !exists {
		return fmt.Errorf("sell_to_buying_store packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"buyer_id": buyerID,
		"index":    index,
		"amount":   amount,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}
