// Package actor provides actor-related packet sending functionality.
package actor

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// Stat type constants
const (
	StatStr = 13 // Strength
	StatAgi = 14 // Agility
	StatVit = 15 // Vitality
	StatInt = 16 // Intelligence
	StatDex = 17 // Dexterity
	StatLuk = 18 // Luck
)

// AppearanceManager handles actor appearance-related packet sending.
type AppearanceManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewAppearanceManager creates a new appearance manager.
func NewAppearanceManager(baseSend core.Send) *AppearanceManager {
	return &AppearanceManager{
		baseSend: baseSend,
	}
}

// SendChangeClothes sends a request to change the character's equipment appearance.
func (am *AppearanceManager) SendChangeClothes(headTop, headMid, headBottom int) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("change_clothes")
	if !exists {
		return fmt.Errorf("change_clothes packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"head_top":    headTop,
		"head_mid":    headMid,
		"head_bottom": headBottom,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendChangeHair sends a request to change the character's hair style and color.
func (am *AppearanceManager) SendChangeHair(hairStyle, hairColor int) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("change_hair")
	if !exists {
		return fmt.Errorf("change_hair packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"hair_style": hairStyle,
		"hair_color": hairColor,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendChangeStat sends a request to change a character stat.
func (am *AppearanceManager) SendChangeStat(statType, amount int) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("change_stat")
	if !exists {
		return fmt.Errorf("change_stat packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"stat_type": statType,
		"amount":    amount,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}
