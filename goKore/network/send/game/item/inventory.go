// Package item provides item-related packet sending functionality.
package item

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// InventoryManager handles inventory-related packet sending.
type InventoryManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewInventoryManager creates a new inventory manager.
func NewInventoryManager(baseSend core.Send) *InventoryManager {
	return &InventoryManager{
		baseSend: baseSend,
	}
}

// UseItem sends a request to use an item.
// This is equivalent to the sendItemUse function in Send.pm.
func (im *InventoryManager) UseItem(index uint16, targetID uint32) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("use_item")
	if !exists {
		return fmt.Errorf("use_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index":     index,
		"target_id": targetID,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendTake sends a request to pick up an item from the ground.
// This is equivalent to the sendTake function in Send.pm.
func (im *InventoryManager) SendTake(objectID uint32) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("take")
	if !exists {
		return fmt.Errorf("take packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"object_id": objectID,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// DropItem sends a request to drop an item from the inventory.
// This is equivalent to the sendDrop function in Send.pm.
func (im *InventoryManager) DropItem(index, amount uint16) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("drop_item")
	if !exists {
		return fmt.Errorf("drop_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index":  index,
		"amount": amount,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// MoveItem sends a request to move an item from one inventory slot to another.
// This is equivalent to the sendMoveItem function in Send.pm.
func (im *InventoryManager) MoveItem(fromIndex, toIndex, amount uint16) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("move_item")
	if !exists {
		return fmt.Errorf("move_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"from_index": fromIndex,
		"to_index":   toIndex,
		"amount":     amount,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SplitItem sends a request to split a stack of items.
// This is equivalent to the sendSplitItem function in Send.pm.
func (im *InventoryManager) SplitItem(index, amount uint16) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("split_item")
	if !exists {
		return fmt.Errorf("split_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index":  index,
		"amount": amount,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendEquip sends a request to equip an item.
// This is equivalent to the sendEquip function in Send.pm.
func (im *InventoryManager) SendEquip(inventoryIndex, equipType uint32) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("equip")
	if !exists {
		return fmt.Errorf("equip packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"inventory_index": inventoryIndex,
		"equip_type":      equipType,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendUnequip sends a request to unequip an item.
// This is equivalent to the sendUnequip function in Send.pm.
func (im *InventoryManager) SendUnequip(equipIndex uint32) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("unequip")
	if !exists {
		return fmt.Errorf("unequip packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"equip_index": equipIndex,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// IdentifyItem sends a request to identify an item.
// This is equivalent to the sendIdentify function in Send.pm.
func (im *InventoryManager) IdentifyItem(index uint16) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("identify_item")
	if !exists {
		return fmt.Errorf("identify_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index": index,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendCardMergeRequest sends a request to merge a card with an item.
// This is equivalent to the sendCardMergeRequest function in Send.pm.
func (im *InventoryManager) SendCardMergeRequest(cardIndex uint32) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("card_merge_request")
	if !exists {
		return fmt.Errorf("card_merge_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"card_index": cardIndex,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendCardMerge sends a request to merge a card with an item.
// This is equivalent to the sendCardMerge function in Send.pm.
func (im *InventoryManager) SendCardMerge(cardIndex, itemIndex uint32) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("card_merge")
	if !exists {
		return fmt.Errorf("card_merge packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"card_index": cardIndex,
		"item_index": itemIndex,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendEquipSwitchAdd sends a request to add an item to the equip switch window.
// This is equivalent to the sendEquipSwitchAdd function in Send.pm.
func (im *InventoryManager) SendEquipSwitchAdd(inventoryIndex, position uint32) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("equip_switch_add")
	if !exists {
		return fmt.Errorf("equip_switch_add packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":       inventoryIndex,
		"position": position,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendEquipSwitchRemove sends a request to remove an item from the equip switch window.
// This is equivalent to the sendEquipSwitchRemove function in Send.pm.
func (im *InventoryManager) SendEquipSwitchRemove(inventoryIndex, position uint32) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("equip_switch_remove")
	if !exists {
		return fmt.Errorf("equip_switch_remove packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":       inventoryIndex,
		"position": position,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendEquipSwitchRun sends a request to execute a full equip switch.
// This is equivalent to the sendEquipSwitchRun function in Send.pm.
func (im *InventoryManager) SendEquipSwitchRun() error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("equip_switch_run")
	if !exists {
		return fmt.Errorf("equip_switch_run packet ID not found")
	}

	// Create the arguments (no arguments for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendEquipSwitchSingle sends a request to execute a single equip switch.
// This is equivalent to the sendEquipSwitchSingle function in Send.pm.
func (im *InventoryManager) SendEquipSwitchSingle(inventoryIndex uint32) error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("equip_switch_single")
	if !exists {
		return fmt.Errorf("equip_switch_single packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": inventoryIndex,
	}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendChangeDress sends a request to change dress.
// This is equivalent to the sendChangeDress function in Send.pm.
func (im *InventoryManager) SendChangeDress() error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("change_dress")
	if !exists {
		return fmt.Errorf("change_dress packet ID not found")
	}

	// Create the arguments (no arguments for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendInventoryExpansionRequest sends a request to expand the inventory.
// This is equivalent to the sendInventoryExpansionRequest function in Send.pm.
func (im *InventoryManager) SendInventoryExpansionRequest() error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("inventory_expansion_request")
	if !exists {
		return fmt.Errorf("inventory_expansion_request packet ID not found")
	}

	// Create the arguments (no arguments for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}

// SendInventoryExpansionRejected sends a request to reject inventory expansion.
// This is equivalent to the sendInventoryExpansionRejected function in Send.pm.
func (im *InventoryManager) SendInventoryExpansionRejected() error {
	// Get the packet ID
	packetID, exists := im.baseSend.GetPacketID("inventory_expansion_rejected")
	if !exists {
		return fmt.Errorf("inventory_expansion_rejected packet ID not found")
	}

	// Create the arguments (no arguments for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := im.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return im.baseSend.SendToServer(packet)
}
