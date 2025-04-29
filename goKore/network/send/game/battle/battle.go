// Package battle provides battle-related packet sending functionality.
package battle

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// BattleManager handles battle-related packet sending.
type BattleManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewBattleManager creates a new battle manager.
func NewBattleManager(baseSend core.Send) *BattleManager {
	return &BattleManager{
		baseSend: baseSend,
	}
}

// SendShowEquipPlayer sends a request to view another player's equipment.
// This is equivalent to the sendShowEquipPlayer function in Send.pm.
func (bm *BattleManager) SendShowEquipPlayer(playerID uint32) error {
	// Get the packet ID
	packetID, exists := bm.baseSend.GetPacketID("view_player_equip_request")
	if !exists {
		return fmt.Errorf("view_player_equip_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": playerID,
	}

	// Construct and send the packet
	packet, err := bm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bm.baseSend.SendToServer(packet)
}

// SendEmotion sends an emotion.
// This is equivalent to the sendEmotion function in Send.pm.
func (bm *BattleManager) SendEmotion(emotionID uint8) error {
	// Get the packet ID
	packetID, exists := bm.baseSend.GetPacketID("send_emotion")
	if !exists {
		return fmt.Errorf("send_emotion packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": emotionID,
	}

	// Construct and send the packet
	packet, err := bm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bm.baseSend.SendToServer(packet)
}

// SendNoviceDoriDori sends a novice dori dori skill request.
// This is equivalent to the sendNoviceDoriDori function in Send.pm.
func (bm *BattleManager) SendNoviceDoriDori() error {
	// Get the packet ID
	packetID, exists := bm.baseSend.GetPacketID("novice_dori_dori")
	if !exists {
		return fmt.Errorf("novice_dori_dori packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := bm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bm.baseSend.SendToServer(packet)
}

// SendNoviceExplosionSpirits sends a novice explosion spirits skill request.
// This is equivalent to the sendNoviceExplosionSpirits function in Send.pm.
func (bm *BattleManager) SendNoviceExplosionSpirits() error {
	// Get the packet ID
	packetID, exists := bm.baseSend.GetPacketID("novice_explosion_spirits")
	if !exists {
		return fmt.Errorf("novice_explosion_spirits packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := bm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bm.baseSend.SendToServer(packet)
}

// SendMemorialDungeonCommand sends a memorial dungeon command.
// This is equivalent to the sendMemorialDungeonCommand function in Send.pm.
func (bm *BattleManager) SendMemorialDungeonCommand(command uint32) error {
	// Get the packet ID
	packetID, exists := bm.baseSend.GetPacketID("memorial_dungeon_command")
	if !exists {
		return fmt.Errorf("memorial_dungeon_command packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"command": command,
	}

	// Construct and send the packet
	packet, err := bm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return bm.baseSend.SendToServer(packet)
}
