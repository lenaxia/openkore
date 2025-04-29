// Package party provides party-related packet sending functionality.
package party

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// PartyManager handles party-related packet sending.
type PartyManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewPartyManager creates a new party manager.
func NewPartyManager(baseSend core.Send) *PartyManager {
	return &PartyManager{
		baseSend: baseSend,
	}
}

// SendPartyOption sends a request to change party options.
// This is equivalent to the sendPartyOption function in Send.pm.
func (pm *PartyManager) SendPartyOption(exp, itemPickup, itemDivision int) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("party_setting")
	if !exists {
		return fmt.Errorf("party_setting packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"exp":          exp,
		"itemPickup":   itemPickup,
		"itemDivision": itemDivision,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendPartyLeader sends a request to change the party leader.
// This is equivalent to the sendPartyLeader function in Send.pm.
func (pm *PartyManager) SendPartyLeader(accountID uint32) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("party_leader")
	if !exists {
		return fmt.Errorf("party_leader packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"accountID": accountID,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendPartyJoinRequest sends a request to join a party.
// This is equivalent to the sendPartyJoinRequest function in Send.pm.
func (pm *PartyManager) SendPartyJoinRequest(ID uint32) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("party_join_request")
	if !exists {
		return fmt.Errorf("party_join_request packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": ID,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendPartyJoin sends a response to a party join request.
// This is equivalent to the sendPartyJoin function in Send.pm.
func (pm *PartyManager) SendPartyJoin(ID uint32, flag int) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("party_join")
	if !exists {
		return fmt.Errorf("party_join packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":   ID,
		"flag": flag,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendPartyKick sends a request to kick a member from the party.
// This is equivalent to the sendPartyKick function in Send.pm.
func (pm *PartyManager) SendPartyKick(ID uint32, name string) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("party_kick")
	if !exists {
		return fmt.Errorf("party_kick packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":   ID,
		"name": []byte(name), // stringToBytes in the original
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendPartyJoinRequestByName sends a request to join a party by name.
// This is equivalent to the sendPartyJoinRequestByName function in Send.pm.
func (pm *PartyManager) SendPartyJoinRequestByName(name string) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("party_join_request_by_name")
	if !exists {
		return fmt.Errorf("party_join_request_by_name packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"partyName": []byte(name), // stringToBytes in the original
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendPartyJoinRequestByNameReply sends a reply to a party join request by name.
// This is equivalent to the sendPartyJoinRequestByNameReply function in Send.pm.
func (pm *PartyManager) SendPartyJoinRequestByNameReply(accountID uint32, flag int) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("party_join_request_by_name_reply")
	if !exists {
		return fmt.Errorf("party_join_request_by_name_reply packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"accountID": accountID,
		"flag":      flag,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendPartyBookingRegister sends a request to register a party booking.
// This is equivalent to the sendPartyBookingRegister function in Send.pm.
func (pm *PartyManager) SendPartyBookingRegister(level, mapID int, jobList []int) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("booking_register")
	if !exists {
		return fmt.Errorf("booking_register packet ID not found")
	}

	// Ensure jobList has 6 elements
	for len(jobList) < 6 {
		jobList = append(jobList, 0)
	}

	// Create the arguments
	args := map[string]interface{}{
		"level": level,
		"MapID": mapID,
		"job0":  jobList[0],
		"job1":  jobList[1],
		"job2":  jobList[2],
		"job3":  jobList[3],
		"job4":  jobList[4],
		"job5":  jobList[5],
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendPartyBookingReqSearch sends a request to search for party bookings.
// This is equivalent to the sendPartyBookingReqSearch function in Send.pm.
func (pm *PartyManager) SendPartyBookingReqSearch(level, mapID, job, lastIndex, resultCount int) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("booking_search")
	if !exists {
		return fmt.Errorf("booking_search packet ID not found")
	}

	// Apply defaults
	if job == 0 {
		job = 65535 // job null = 65535
	}
	if resultCount == 0 {
		resultCount = 10 // ResultCount default = 10
	}

	// Create the arguments
	args := map[string]interface{}{
		"level":       level,
		"MapID":       mapID,
		"job":         job,
		"LastIndex":   lastIndex,
		"ResultCount": resultCount,
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendPartyBookingDelete sends a request to delete a party booking.
// This is equivalent to the sendPartyBookingDelete function in Send.pm.
func (pm *PartyManager) SendPartyBookingDelete() error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("booking_delete")
	if !exists {
		return fmt.Errorf("booking_delete packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendPartyBookingUpdate sends a request to update a party booking.
// This is equivalent to the sendPartyBookingUpdate function in Send.pm.
func (pm *PartyManager) SendPartyBookingUpdate(jobList []int) error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("booking_update")
	if !exists {
		return fmt.Errorf("booking_update packet ID not found")
	}

	// Ensure jobList has 6 elements
	for len(jobList) < 6 {
		jobList = append(jobList, 0)
	}

	// Create the arguments
	args := map[string]interface{}{
		"job0": jobList[0],
		"job1": jobList[1],
		"job2": jobList[2],
		"job3": jobList[3],
		"job4": jobList[4],
		"job5": jobList[5],
	}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}

// SendPartyLeave sends a request to leave the party.
// This is equivalent to the sendPartyLeave function in Send.pm.
func (pm *PartyManager) SendPartyLeave() error {
	// Get the packet ID
	packetID, exists := pm.baseSend.GetPacketID("party_leave")
	if !exists {
		return fmt.Errorf("party_leave packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := pm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return pm.baseSend.SendToServer(packet)
}
