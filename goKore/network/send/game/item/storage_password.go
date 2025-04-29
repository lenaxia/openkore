// Package item provides item-related packet sending functionality.
package item

import (
	"fmt"
)

// ReconstructStoragePassword reconstructs the storage password packet.
// This is equivalent to the reconstruct_storage_password function in Send.pm.
func ReconstructStoragePassword(args map[string]interface{}) error {
	// Get the type and password
	type_, ok := args["type"].(int)
	if !ok {
		return fmt.Errorf("type is not an int")
	}

	pass, ok := args["pass"].([]byte)
	if !ok {
		return fmt.Errorf("pass is not a byte slice")
	}

	// Create the auxiliary data
	aux := []byte{
		0xEC, 0x62, 0xE5, 0x39, 0xBB, 0x6B, 0xBC, 0x81,
		0x1A, 0x60, 0xC0, 0x6F, 0xAC, 0xCB, 0x7E, 0xC8,
	}

	// Create the data based on the type
	if type_ == 3 {
		// Check password
		args["data"] = append(pass, aux...)
	} else if type_ == 2 {
		// Change password
		args["data"] = append(aux, pass...)
	} else {
		return fmt.Errorf("invalid type: %d", type_)
	}

	return nil
}

// SendStoragePassword sends a storage password packet.
// This is equivalent to the sendStoragePassword function in Send.pm.
func (sm *StorageManager) SendStoragePassword(pass []byte, type_ int) error {
	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID("storage_password")
	if !exists {
		return fmt.Errorf("storage_password packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type": type_,
		"pass": pass,
	}

	// Reconstruct the storage password
	if err := ReconstructStoragePassword(args); err != nil {
		return err
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SendStorageGetToCart sends a request to move an item from storage to cart.
// This is equivalent to the sendStorageGetToCart function in Send.pm.
func (sm *StorageManager) SendStorageGetToCart(index, amount uint16, isGuildStorage bool) error {
	var packetName string
	if isGuildStorage {
		packetName = "guild_storage_to_cart"
	} else {
		packetName = "storage_to_cart"
	}

	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID(packetName)
	if !exists {
		return fmt.Errorf("%s packet ID not found", packetName)
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":     index,
		"amount": amount,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}

// SendStorageAddFromCart sends a request to move an item from cart to storage.
// This is equivalent to the sendStorageAddFromCart function in Send.pm.
func (sm *StorageManager) SendStorageAddFromCart(index, amount uint16, isGuildStorage bool) error {
	var packetName string
	if isGuildStorage {
		packetName = "cart_to_guild_storage"
	} else {
		packetName = "cart_to_storage"
	}

	// Get the packet ID
	packetID, exists := sm.baseSend.GetPacketID(packetName)
	if !exists {
		return fmt.Errorf("%s packet ID not found", packetName)
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":     index,
		"amount": amount,
	}

	// Construct and send the packet
	packet, err := sm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return sm.baseSend.SendToServer(packet)
}
