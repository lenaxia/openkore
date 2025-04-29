// Package actor provides actor-related packet sending functionality.
package actor

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// Action constants
const (
	ActionAttack = 1
	ActionSit    = 2
	ActionStand  = 3
	ActionSkill  = 4
	ActionTalk   = 5
)

// ActionManager handles actor action-related packet sending.
type ActionManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewActionManager creates a new action manager.
func NewActionManager(baseSend core.Send) *ActionManager {
	return &ActionManager{
		baseSend: baseSend,
	}
}

// SendAction sends an action command for the player character.
func (am *ActionManager) SendAction(targetID uint32, action int) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("send_action")
	if !exists {
		return fmt.Errorf("send_action packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"target_id": targetID,
		"action":    action,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendLook sends a look command to change the character's direction.
func (am *ActionManager) SendLook(body, head int) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("send_look")
	if !exists {
		return fmt.Errorf("send_look packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"body": body,
		"head": head,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendEmotion sends an emotion command.
func (am *ActionManager) SendEmotion(emotion int) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("send_emotion")
	if !exists {
		return fmt.Errorf("send_emotion packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"emotion": emotion,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendSit sends a sit command.
func (am *ActionManager) SendSit() error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("send_sit")
	if !exists {
		return fmt.Errorf("send_sit packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"action": ActionSit,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendStand sends a stand command.
func (am *ActionManager) SendStand() error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("send_stand")
	if !exists {
		return fmt.Errorf("send_stand packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"action": ActionStand,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendSlaveAttack sends an attack command for a slave (homunculus, mercenary, etc.).
func (am *ActionManager) SendSlaveAttack(slaveID, targetID uint32, flag int) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("slave_attack")
	if !exists {
		return fmt.Errorf("slave_attack packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"slaveID":  slaveID,
		"targetID": targetID,
		"flag":     flag,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendSlaveStandBy sends a command for a slave to return to its master.
func (am *ActionManager) SendSlaveStandBy(slaveID uint32) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("slave_move_to_master")
	if !exists {
		return fmt.Errorf("slave_move_to_master packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"slaveID": slaveID,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}

// SendAlignment sends an alignment command.
func (am *ActionManager) SendAlignment(targetID uint32, alignment, point int) error {
	// Get the packet ID
	packetID, exists := am.baseSend.GetPacketID("alignment")
	if !exists {
		return fmt.Errorf("alignment packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"targetID": targetID,
		"type":     alignment,
		"point":    point,
	}

	// Construct and send the packet
	packet, err := am.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return am.baseSend.SendToServer(packet)
}
