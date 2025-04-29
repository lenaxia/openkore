// Package item provides item-related packet sending functionality.
package item

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// Equipment position constants
const (
	EquipPosLowerHead  = 1
	EquipPosRightHand  = 2
	EquipPosGarment    = 4
	EquipPosLeftHand   = 8
	EquipPosArmor      = 16
	EquipPosMiddleHead = 32
	EquipPosUpperHead  = 64
	EquipPosRightAcc   = 128
	EquipPosLeftAcc    = 256
	EquipPosShoes      = 512
	EquipPosMouth      = 1024
	EquipPosAmmo       = 2048
	EquipPosShadow     = 4096
)

// EquipmentManager handles equipment-related packet sending.
type EquipmentManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewEquipmentManager creates a new equipment manager.
func NewEquipmentManager(baseSend core.Send) *EquipmentManager {
	return &EquipmentManager{
		baseSend: baseSend,
	}
}

// EquipItem sends a request to equip an item.
func (em *EquipmentManager) EquipItem(index, position uint16) error {
	// Get the packet ID
	packetID, exists := em.baseSend.GetPacketID("equip_item")
	if !exists {
		return fmt.Errorf("equip_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index":    index,
		"position": position,
	}

	// Construct and send the packet
	packet, err := em.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return em.baseSend.SendToServer(packet)
}

// UnequipItem sends a request to unequip an item.
func (em *EquipmentManager) UnequipItem(index uint16) error {
	// Get the packet ID
	packetID, exists := em.baseSend.GetPacketID("unequip_item")
	if !exists {
		return fmt.Errorf("unequip_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index": index,
	}

	// Construct and send the packet
	packet, err := em.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return em.baseSend.SendToServer(packet)
}

// UpgradeItem sends a request to upgrade an item.
func (em *EquipmentManager) UpgradeItem(index uint16) error {
	// Get the packet ID
	packetID, exists := em.baseSend.GetPacketID("upgrade_item")
	if !exists {
		return fmt.Errorf("upgrade_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index": index,
	}

	// Construct and send the packet
	packet, err := em.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return em.baseSend.SendToServer(packet)
}

// RefineItem sends a request to refine an item.
func (em *EquipmentManager) RefineItem(index uint16) error {
	// Get the packet ID
	packetID, exists := em.baseSend.GetPacketID("refine_item")
	if !exists {
		return fmt.Errorf("refine_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index": index,
	}

	// Construct and send the packet
	packet, err := em.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return em.baseSend.SendToServer(packet)
}

// RepairItem sends a request to repair an item.
func (em *EquipmentManager) RepairItem(index, itemID uint16) error {
	// Get the packet ID
	packetID, exists := em.baseSend.GetPacketID("repair_item")
	if !exists {
		return fmt.Errorf("repair_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index":   index,
		"item_id": itemID,
	}

	// Construct and send the packet
	packet, err := em.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return em.baseSend.SendToServer(packet)
}

// EnchantItem sends a request to enchant an item.
func (em *EquipmentManager) EnchantItem(index uint16, cardID uint32) error {
	// Get the packet ID
	packetID, exists := em.baseSend.GetPacketID("enchant_item")
	if !exists {
		return fmt.Errorf("enchant_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index":   index,
		"card_id": cardID,
	}

	// Construct and send the packet
	packet, err := em.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return em.baseSend.SendToServer(packet)
}

// ArrangeEquipment sends a request to arrange the equipment window.
func (em *EquipmentManager) ArrangeEquipment() error {
	// Get the packet ID
	packetID, exists := em.baseSend.GetPacketID("arrange_equipment")
	if !exists {
		return fmt.Errorf("arrange_equipment packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := em.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return em.baseSend.SendToServer(packet)
}
