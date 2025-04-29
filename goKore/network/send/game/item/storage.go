// Package item provides item-related packet sending functionality.
package item

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// StorageManager handles storage-related packet sending.
type StorageManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewStorageManager creates a new storage manager.
func NewStorageManager(baseSend core.Send) *StorageManager {
	return &StorageManager{
		baseSend: baseSend,
	}
}

// OpenStorage sends a request to open the storage.
func (sm *StorageManager) OpenStorage() error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("open_storage")
	if !exists {
		return fmt.Errorf("open_storage packet ID not found")
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

// CloseStorage sends a request to close the storage.
func (sm *StorageManager) CloseStorage() error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("close_storage")
	if !exists {
		return fmt.Errorf("close_storage packet ID not found")
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

// MoveToStorage sends a request to move an item from the inventory to storage.
func (sm *StorageManager) MoveToStorage(index, amount uint16) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("move_to_storage")
	if !exists {
		return fmt.Errorf("move_to_storage packet ID not found")
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

// MoveFromStorage sends a request to move an item from storage to the inventory.
func (sm *StorageManager) MoveFromStorage(index, amount uint16) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("move_from_storage")
	if !exists {
		return fmt.Errorf("move_from_storage packet ID not found")
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

// ArrangeStorage sends a request to arrange the storage.
func (sm *StorageManager) ArrangeStorage() error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("arrange_storage")
	if !exists {
		return fmt.Errorf("arrange_storage packet ID not found")
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

// OpenGuildStorage sends a request to open the guild storage.
func (sm *StorageManager) OpenGuildStorage() error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("open_guild_storage")
	if !exists {
		return fmt.Errorf("open_guild_storage packet ID not found")
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

// CloseGuildStorage sends a request to close the guild storage.
func (sm *StorageManager) CloseGuildStorage() error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("close_guild_storage")
	if !exists {
		return fmt.Errorf("close_guild_storage packet ID not found")
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

// MoveToGuildStorage sends a request to move an item from the inventory to guild storage.
func (sm *StorageManager) MoveToGuildStorage(index, amount uint16) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("move_to_guild_storage")
	if !exists {
		return fmt.Errorf("move_to_guild_storage packet ID not found")
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

// MoveFromGuildStorage sends a request to move an item from guild storage to the inventory.
func (sm *StorageManager) MoveFromGuildStorage(index, amount uint16) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("move_from_guild_storage")
	if !exists {
		return fmt.Errorf("move_from_guild_storage packet ID not found")
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
