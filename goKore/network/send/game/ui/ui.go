// Package ui provides UI-related packet sending functionality.
package ui

import (
	"encoding/binary"
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// UIManager handles UI-related packet sending.
type UIManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewUIManager creates a new UI manager.
func NewUIManager(baseSend core.Send) *UIManager {
	return &UIManager{
		baseSend: baseSend,
	}
}

// SendMiscConfigSet sends a request to set a miscellaneous configuration.
// This is equivalent to the sendMiscConfigSet function in Send.pm.
// configType:
//
//	0 = show equip windows to other players
//	1 = being summoned by skills: Urgent Call, Romantic Rendezvous, Come to me, honey~ & Let's Go, Family!
//	2 = pet autofeeding
//	3 = homunculus autofeeding
//
// flag:
//
//	0 = disabled
//	1 = enabled
func (um *UIManager) SendMiscConfigSet(configType, flag uint32) error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("misc_config_set")
	if !exists {
		return fmt.Errorf("misc_config_set packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type": configType,
		"flag": flag,
	}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendProgress sends a notification that the progress bar is complete.
// This is equivalent to the sendProgress function in Send.pm.
func (um *UIManager) SendProgress() error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("notify_progress_bar_complete")
	if !exists {
		return fmt.Errorf("notify_progress_bar_complete packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendShowEquipPlayer sends a request to view another player's equipment.
// This is equivalent to the sendShowEquipPlayer function in Send.pm.
func (um *UIManager) SendShowEquipPlayer(playerID uint32) error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("view_player_equip_request")
	if !exists {
		return fmt.Errorf("view_player_equip_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": playerID,
	}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendRefineUISelect sends a request to select an item for refining.
// This is equivalent to the sendRefineUISelect function in Send.pm.
func (um *UIManager) SendRefineUISelect(itemIndex uint16) error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("refineui_select")
	if !exists {
		return fmt.Errorf("refineui_select packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index": itemIndex,
	}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendRefineUIRefine sends a request to refine an item.
// This is equivalent to the sendRefineUIRefine function in Send.pm.
func (um *UIManager) SendRefineUIRefine(itemIndex, materialNameID uint16, useCatalyst uint8) error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("refineui_refine")
	if !exists {
		return fmt.Errorf("refineui_refine packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"index":    itemIndex,
		"catalyst": materialNameID,
		"bless":    useCatalyst,
	}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendRefineUIClose sends a request to close the refine UI.
// This is equivalent to the sendRefineUIClose function in Send.pm.
func (um *UIManager) SendRefineUIClose() error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("refineui_close")
	if !exists {
		return fmt.Errorf("refineui_close packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendItemListWindowSelected sends a request to select items from a list window.
// This is equivalent to the sendItemListWindowSelected function in Send.pm.
// itemType:
//
//	0 = Change Material
//	1 = Elemental Analysis (Level 1: Pure to Rough)
//	2 = Elemental Analysis (Level 1: Rough to Pure)
//
// act:
//
//	0 = Cancel
//	1 = Process
//
// items: List of items [itemIndex, amount, itemName]
func (um *UIManager) SendItemListWindowSelected(itemType, act uint8, items []map[string]interface{}) error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("item_list_window_selected")
	if !exists {
		return fmt.Errorf("item_list_window_selected packet ID not found")
	}

	// Calculate the length
	length := (len(items) * 4) + 12

	// Create the arguments
	args := map[string]interface{}{
		"len":   length,
		"type":  itemType,
		"act":   act,
		"items": items,
	}

	// Reconstruct the itemInfo field
	if err := um.ReconstructItemListWindowSelected(args); err != nil {
		return err
	}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// ReconstructItemListWindowSelected reconstructs the itemInfo field for the item_list_window_selected packet.
// This is equivalent to the reconstruct_item_list_window_selected function in Send.pm.
func (um *UIManager) ReconstructItemListWindowSelected(args map[string]interface{}) error {
	items, ok := args["items"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("items not found or not a slice of maps")
	}

	// Each item is 4 bytes: 2 bytes for itemIndex, 2 bytes for amount
	itemInfo := make([]byte, len(items)*4)

	for i, item := range items {
		itemIndex, ok := item["itemIndex"].(uint16)
		if !ok {
			return fmt.Errorf("itemIndex not found or not a uint16 for item %d", i)
		}

		amount, ok := item["amount"].(uint16)
		if !ok {
			return fmt.Errorf("amount not found or not a uint16 for item %d", i)
		}

		binary.LittleEndian.PutUint16(itemInfo[i*4:i*4+2], itemIndex)
		binary.LittleEndian.PutUint16(itemInfo[i*4+2:i*4+4], amount)
	}

	args["itemInfo"] = itemInfo
	return nil
}

// SendMemo sends a request to create a memo at the current location.
// This is equivalent to the sendMemo function in Send.pm.
func (um *UIManager) SendMemo() error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("memo_request")
	if !exists {
		return fmt.Errorf("memo_request packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendStylistChange sends a request to change the character's appearance.
// This is equivalent to the sendStylistChange function in Send.pm.
func (um *UIManager) SendStylistChange(hairColor, hairStyle, clothColor, headTop, headMid, headBottom uint16) error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("stylist_change")
	if !exists {
		return fmt.Errorf("stylist_change packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"hair_color":  hairColor,
		"hair_style":  hairStyle,
		"cloth_color": clothColor,
		"head_top":    headTop,
		"head_mid":    headMid,
		"head_bottom": headBottom,
	}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendOpenUIRequest sends a request to open a UI window.
// This is equivalent to the sendOpenUIRequest function in Send.pm.
func (um *UIManager) SendOpenUIRequest(uiType uint8) error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("open_ui_request")
	if !exists {
		return fmt.Errorf("open_ui_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"UIType": uiType,
	}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendAttendanceRewardRequest sends a request to retrieve today's attendance reward.
// This is equivalent to the sendAttendanceRewardRequest function in Send.pm.
func (um *UIManager) SendAttendanceRewardRequest() error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("attendance_reward_request")
	if !exists {
		return fmt.Errorf("attendance_reward_request packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendRouletteWindowOpen sends a request to open the roulette window.
// This is equivalent to the sendRouletteWindowOpen function in Send.pm.
func (um *UIManager) SendRouletteWindowOpen() error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("roulette_window_open")
	if !exists {
		return fmt.Errorf("roulette_window_open packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendRouletteInfoRequest sends a request for roulette information.
// This is equivalent to the sendRouletteInfoRequest function in Send.pm.
func (um *UIManager) SendRouletteInfoRequest() error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("roulette_info_request")
	if !exists {
		return fmt.Errorf("roulette_info_request packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendRouletteClose sends a request to close the roulette window.
// This is equivalent to the sendRouletteClose function in Send.pm.
func (um *UIManager) SendRouletteClose() error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("roulette_close")
	if !exists {
		return fmt.Errorf("roulette_close packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendRouletteStart sends a request to start the roulette.
// This is equivalent to the sendRouletteStart function in Send.pm.
func (um *UIManager) SendRouletteStart() error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("roulette_start")
	if !exists {
		return fmt.Errorf("roulette_start packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendRouletteClaimPrize sends a request to claim a roulette prize.
// This is equivalent to the sendRouletteClaimPrize function in Send.pm.
func (um *UIManager) SendRouletteClaimPrize() error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("roulette_claim_prize")
	if !exists {
		return fmt.Errorf("roulette_claim_prize packet ID not found")
	}

	// No arguments needed for this packet
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}

// SendQuestState sends a request to change a quest state.
// This is equivalent to the sendQuestState function in Send.pm.
// state:
//
//	0 = active
//	1 = inactive
func (um *UIManager) SendQuestState(questID uint32, state uint8) error {
	// Get the packet ID
	packetID, exists := um.baseSend.GetPacketID("send_quest_state")
	if !exists {
		return fmt.Errorf("send_quest_state packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"questID": questID,
		"state":   state,
	}

	// Construct and send the packet
	packet, err := um.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return um.baseSend.SendToServer(packet)
}
