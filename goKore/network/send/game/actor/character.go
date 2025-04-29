// Package actor provides actor-related packet sending functionality.
package actor

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// CharacterManager handles character-related packet sending.
type CharacterManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewCharacterManager creates a new character manager.
func NewCharacterManager(baseSend core.Send) *CharacterManager {
	return &CharacterManager{
		baseSend: baseSend,
	}
}

// SendLook sends a look command to change the character's direction and head direction.
// This is equivalent to the sendLook function in Send.pm.
func (cm *CharacterManager) SendLook(bodyDirection, headDirection int) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("look")
	if !exists {
		return fmt.Errorf("look packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"body_direction": uint8(bodyDirection),
		"head_direction": uint8(headDirection),
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendAction sends an action command for the character.
// This is equivalent to the sendAction function in Send.pm.
// flag: 0 attack (once), 7 attack (continuous), 2 sit, 3 stand
func (cm *CharacterManager) SendAction(targetID uint32, action int) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("action")
	if !exists {
		return fmt.Errorf("action packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"target_id": targetID,
		"action":    uint8(action),
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendEmotion sends an emotion command for the character.
// This is equivalent to the sendEmotion function in Send.pm.
func (cm *CharacterManager) SendEmotion(emotion int) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("emotion")
	if !exists {
		return fmt.Errorf("emotion packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"emotion": uint8(emotion),
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendRespawn sends a respawn command for the character.
// This is equivalent to the sendRespawn function in Send.pm.
func (cm *CharacterManager) SendRespawn() error {
	return cm.SendRestart(0)
}

// SendQuit sends a command to quit the game.
// This is equivalent to the sendQuit function in Send.pm.
func (cm *CharacterManager) SendQuit() error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("quit_request")
	if !exists {
		return fmt.Errorf("quit_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type": 0,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendQuitToCharSelect sends a command to quit to the character selection screen.
// This is equivalent to the sendQuitToCharSelect function in Send.pm.
func (cm *CharacterManager) SendQuitToCharSelect() error {
	return cm.SendRestart(1)
}

// SendRestart sends a restart command for the character.
// This is equivalent to the sendRestart function in Send.pm.
func (cm *CharacterManager) SendRestart(type_ int) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("restart")
	if !exists {
		return fmt.Errorf("restart packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type": uint8(type_),
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendCharCreate sends a character creation request.
// This is equivalent to the sendCharCreate function in Send.pm.
// Different packet IDs require different parameters:
// - 0067: All stats are required
// - 0970: Only slot, name, hair_style, and hair_color are required
// - 0A39: slot, name, hair_style, hair_color, job_id, and sex are required
func (cm *CharacterManager) SendCharCreate(slot int, name string, str, agi, vit, int_, dex, luk, hairStyle, hairColor, jobID, sex int) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("char_create")
	if !exists {
		return fmt.Errorf("char_create packet ID not found")
	}

	// Create the arguments based on the packet ID
	args := map[string]interface{}{
		"slot":       slot,
		"name":       name,
		"hair_style": hairStyle,
		"hair_color": hairColor,
	}

	// Add additional arguments based on the packet ID
	if packetID == "0067" {
		args["str"] = str
		args["agi"] = agi
		args["vit"] = vit
		args["int"] = int_
		args["dex"] = dex
		args["luk"] = luk
	} else if packetID == "0A39" {
		args["job_id"] = jobID
		args["sex"] = sex
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendCharDelete sends a character deletion request.
// This is equivalent to the sendCharDelete function in Send.pm.
func (cm *CharacterManager) SendCharDelete(charID []byte, email string) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("char_delete")
	if !exists {
		return fmt.Errorf("char_delete packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"charID": charID,
		"email":  email,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendCharDelete2 sends a character deletion request (version 2).
// This is equivalent to the sendCharDelete2 function in Send.pm.
func (cm *CharacterManager) SendCharDelete2(charID []byte) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("char_delete2")
	if !exists {
		return fmt.Errorf("char_delete2 packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"charID": charID,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendCharDelete2Accept sends a character deletion acceptance request.
// This is equivalent to the sendCharDelete2Accept function in Send.pm.
func (cm *CharacterManager) SendCharDelete2Accept(charID []byte, code string) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("char_delete2_accept")
	if !exists {
		return fmt.Errorf("char_delete2_accept packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"charID": charID,
		"code":   code,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendCharDelete2Cancel sends a character deletion cancellation request.
// This is equivalent to the sendCharDelete2Cancel function in Send.pm.
func (cm *CharacterManager) SendCharDelete2Cancel(charID []byte) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("char_delete2_cancel")
	if !exists {
		return fmt.Errorf("char_delete2_cancel packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"charID": charID,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendAddStatusPoint sends a request to add a status point.
// This is equivalent to the sendAddStatusPoint function in Send.pm.
func (cm *CharacterManager) SendAddStatusPoint(statusID, amount int) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("send_add_status_point")
	if !exists {
		return fmt.Errorf("send_add_status_point packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"statusID": statusID,
		"Amount":   amount,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendAddSkillPoint sends a request to add a skill point.
// This is equivalent to the sendAddSkillPoint function in Send.pm.
func (cm *CharacterManager) SendAddSkillPoint(skillID int) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("send_add_skill_point")
	if !exists {
		return fmt.Errorf("send_add_skill_point packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"skillID": skillID,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendHotKeyChange sends a request to change a hotkey.
// This is equivalent to the sendHotKeyChange function in Send.pm.
func (cm *CharacterManager) SendHotKeyChange(idx, type_, id, lvl int) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("hotkey_change")
	if !exists {
		return fmt.Errorf("hotkey_change packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"idx":  idx,
		"type": type_,
		"id":   id,
		"lvl":  lvl,
	}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendChangeTitle sends a request to change the character's title.
// This is equivalent to the sendchangetitle function in Send.pm.
func (cm *CharacterManager) SendChangeTitle(titleID int) error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("send_change_title")
	if !exists {
		return fmt.Errorf("send_change_title packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": titleID,
	}

	// Construct and send the packet

	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}

// SendAutoRevive sends a request to auto revive.
// This is equivalent to the sendAutoRevive function in Send.pm.
func (cm *CharacterManager) SendAutoRevive() error {
	// Get the packet ID
	packetID, exists := cm.baseSend.GetPacketID("auto_revive")
	if !exists {
		return fmt.Errorf("auto_revive packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := cm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return cm.baseSend.SendToServer(packet)
}
