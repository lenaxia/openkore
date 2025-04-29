// Package game provides handlers for game-related packets.
package game

import (
	"github.com/lenaxia/goKore/network/send/core"
)

// HandlerRegistrar is an interface for registering packet handlers
type HandlerRegistrar interface {
	RegisterHandler(packetName string, handler core.SendHandler)
}

// RegisterCharacterHandlers registers all character-related handlers with the send component.
func RegisterCharacterHandlers(send HandlerRegistrar) {
	// Register restart handler
	send.RegisterHandler("restart", func(args map[string]interface{}) ([]byte, error) {
		// Extract the type from the args
		typeVal, ok := args["type"].(uint8)
		if !ok {
			typeVal = 0 // Default to respawn
		}

		// Construct the packet
		// Format: 0x00b2,3,restart,2
		// 3-byte packet with a 2-byte field for the type
		packet := []byte{0xb2, 0x00, typeVal}
		return packet, nil
	})
}
