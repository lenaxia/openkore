// Package gm provides GM-related packet sending functionality.
package gm

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// GMManager handles GM-related packet sending.
type GMManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewGMManager creates a new GM manager.
func NewGMManager(baseSend core.Send) *GMManager {
	return &GMManager{
		baseSend: baseSend,
	}
}

// SendGMSummon sends a request to summon a player to the GM's location.
// This is equivalent to the sendGMSummon function in Send.pm.
func (gm *GMManager) SendGMSummon(playerName string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("gm_summon_player")
	if !exists {
		return fmt.Errorf("gm_summon_player packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"playerName": []byte(playerName),
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGMKick sends a request to kick a player from the server.
// This is equivalent to the sendGMKick function in Send.pm.
func (gm *GMManager) SendGMKick(accountID uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("gm_kick")
	if !exists {
		return fmt.Errorf("gm_kick packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"targetAccountID": accountID,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGMKickAll sends a request to kick all players from the server.
// This is equivalent to the sendGMKickAll function in Send.pm.
func (gm *GMManager) SendGMKickAll() error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("gm_kick_all")
	if !exists {
		return fmt.Errorf("gm_kick_all packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGMMonsterItem sends a request to create a monster or item.
// This is equivalent to the sendGMMonsterItem function in Send.pm.
func (gm *GMManager) SendGMMonsterItem(name string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("gm_item_mob_create")
	if !exists {
		return fmt.Errorf("gm_item_mob_create packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"name": []byte(name),
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGMMapMove sends a request to move to a specific map location.
// This is equivalent to the sendGMMapMove function in Send.pm.
func (gm *GMManager) SendGMMapMove(mapName string, x, y uint16) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("gm_move_to_map")
	if !exists {
		return fmt.Errorf("gm_move_to_map packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"mapName": []byte(mapName),
		"x":       x,
		"y":       y,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGMResetStateSkill sends a request to reset a player's stats or skills.
// This is equivalent to the sendGMResetStateSkill function in Send.pm.
// resetType:
//
//	0 => status
//	1 => skills
func (gm *GMManager) SendGMResetStateSkill(resetType uint8) error {
	// Validate reset type
	if resetType > 1 {
		return fmt.Errorf("invalid reset type: %d (must be 0 or 1)", resetType)
	}

	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("gm_reset_state_skill")
	if !exists {
		return fmt.Errorf("gm_reset_state_skill packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type": resetType,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGMChangeMapType sends a request to change a map cell's type.
// This is equivalent to the sendGMChangeMapType function in Send.pm.
// cellType:
//
//	0 => not walkable
//	1 => walkable
func (gm *GMManager) SendGMChangeMapType(x, y uint16, cellType uint8) error {
	// Validate cell type
	if cellType > 1 {
		return fmt.Errorf("invalid cell type: %d (must be 0 or 1)", cellType)
	}

	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("gm_change_cell_type")
	if !exists {
		return fmt.Errorf("gm_change_cell_type packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"x":    x,
		"y":    y,
		"type": cellType,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGMChangeEffectState sends a request to change a player's effect state.
// This is equivalent to the sendGMChangeEffectState function in Send.pm.
func (gm *GMManager) SendGMChangeEffectState(effectState uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("gm_change_effect_state")
	if !exists {
		return fmt.Errorf("gm_change_effect_state packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"effect_state": effectState,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGMRemove sends a request to remove a player from the server.
// This is equivalent to the sendGMRemove function in Send.pm.
func (gm *GMManager) SendGMRemove(playerName string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("gm_remove")
	if !exists {
		return fmt.Errorf("gm_remove packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"playerName": []byte(playerName),
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGMShift sends a request to teleport to a player's location.
// This is equivalent to the sendGMShift function in Send.pm.
func (gm *GMManager) SendGMShift(playerName string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("gm_shift")
	if !exists {
		return fmt.Errorf("gm_shift packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"playerName": []byte(playerName),
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGMRecall sends a request to recall a player to the GM's location.
// This is equivalent to the sendGMRecall function in Send.pm.
func (gm *GMManager) SendGMRecall(playerName string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("gm_recall")
	if !exists {
		return fmt.Errorf("gm_recall packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"playerName": []byte(playerName),
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGMGiveMannerByName sends a request to give manner points to a player by name.
// This is equivalent to the sendGMGiveMannerByName function in Send.pm.
func (gm *GMManager) SendGMGiveMannerByName(playerName string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("manner_by_name")
	if !exists {
		return fmt.Errorf("manner_by_name packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"playerName": []byte(playerName),
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGMRequestStatus sends a request to get a player's status.
// This is equivalent to the sendGMRequestStatus function in Send.pm.
func (gm *GMManager) SendGMRequestStatus(playerName string) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("gm_request_status")
	if !exists {
		return fmt.Errorf("gm_request_status packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"playerName": []byte(playerName),
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendGMReqAccName sends a request to get a player's account name.
// This is equivalent to the sendGMReqAccName function in Send.pm.
func (gm *GMManager) SendGMReqAccName(targetID uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("gm_request_account_name")
	if !exists {
		return fmt.Errorf("gm_request_account_name packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"targetID": targetID,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}

// SendBanCheck sends a request to check if an account is banned.
// This is equivalent to the sendBanCheck function in Send.pm.
func (gm *GMManager) SendBanCheck(accountID uint32) error {
	// Get the packet ID
	packetID, exists := gm.baseSend.GetPacketID("ban_check")
	if !exists {
		return fmt.Errorf("ban_check packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"accountID": accountID,
	}

	// Construct and send the packet
	packet, err := gm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return gm.baseSend.SendToServer(packet)
}
