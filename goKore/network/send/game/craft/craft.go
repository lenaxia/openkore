// Package craft provides crafting-related packet sending functionality.
package craft

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// CraftManager handles crafting-related packet sending.
type CraftManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewCraftManager creates a new craft manager.
func NewCraftManager(baseSend core.Send) *CraftManager {
	return &CraftManager{
		baseSend: baseSend,
	}
}

// SendIdentify sends a request to identify an item.
// This is equivalent to the sendIdentify function in Send.pm.
func (cm *CraftManager) SendIdentify(ID uint16) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("identify")
	if !exists {
		return fmt.Errorf("identify packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": ID,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendWeaponRefine sends a request to refine a weapon.
// This is equivalent to the sendWeaponRefine function in Send.pm.
func (cm *CraftManager) SendWeaponRefine(ID uint16) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("refine_item")
	if !exists {
		return fmt.Errorf("refine_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": ID,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendCooking sends a request to cook an item.
// This is equivalent to the sendCooking function in Send.pm.
func (cm *CraftManager) SendCooking(type_, nameID uint16) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("cook_request")
	if !exists {
		return fmt.Errorf("cook_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type":   type_,
		"nameID": nameID,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendRepairItem sends a request to repair an item.
// This is equivalent to the sendRepairItem function in Send.pm.
func (cm *CraftManager) SendRepairItem(index uint16, nameID uint16, upgrade uint8, cards []uint32) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("repair_item")
	if !exists {
		return fmt.Errorf("repair_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index":   index,
		"nameID":  nameID,
		"upgrade": upgrade,
		"cards":   cards,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendArrowCraft sends a request to craft arrows.
// This is equivalent to the sendArrowCraft function in Send.pm.
func (cm *CraftManager) SendArrowCraft(nameID uint16) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("make_arrow")
	if !exists {
		return fmt.Errorf("make_arrow packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"nameID": nameID,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendCardMergeRequest sends a request to merge a card.
// This is equivalent to the sendCardMergeRequest function in Send.pm.
func (cm *CraftManager) SendCardMergeRequest(cardID uint16) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("card_merge_request")
	if !exists {
		return fmt.Errorf("card_merge_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"cardID": cardID,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendCardMerge sends a request to merge a card with an item.
// This is equivalent to the sendCardMerge function in Send.pm.
func (cm *CraftManager) SendCardMerge(cardID, itemID uint16) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("card_merge")
	if !exists {
		return fmt.Errorf("card_merge packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"cardID": cardID,
		"itemID": itemID,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendMakeItemRequest sends a request to make an item.
// This is equivalent to the sendMakeItemRequest function in Send.pm.
func (cm *CraftManager) SendMakeItemRequest(nameID, material1, material2, material3 uint16) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("make_item_request")
	if !exists {
		return fmt.Errorf("make_item_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"nameID":           nameID,
		"material_nameID1": material1,
		"material_nameID2": material2,
		"material_nameID3": material3,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}
