// Package rodex provides RODEX mail system-related packet sending functionality.
package rodex

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// RodexManager handles RODEX mail system-related packet sending.
type RodexManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewRodexManager creates a new RODEX mail system manager.
func NewRodexManager(baseSend core.Send) *RodexManager {
	return &RodexManager{
		baseSend: baseSend,
	}
}

// RodexDeleteMail sends a request to delete a mail.
// This is equivalent to the rodex_delete_mail function in Send.pm.
func (rm *RodexManager) RodexDeleteMail(type_ uint8, mailID1, mailID2 uint32) error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_delete_mail")
	if !exists {
		return fmt.Errorf("rodex_delete_mail packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type":    type_,
		"mailID1": mailID1,
		"mailID2": mailID2,
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// RodexRequestZeny sends a request to receive zeny from a mail.
// This is equivalent to the rodex_request_zeny function in Send.pm.
func (rm *RodexManager) RodexRequestZeny(mailID1, mailID2 uint32, type_ uint8) error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_request_zeny")
	if !exists {
		return fmt.Errorf("rodex_request_zeny packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"mailID1": mailID1,
		"mailID2": mailID2,
		"type":    type_,
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// RodexRequestItems sends a request to receive items from a mail.
// This is equivalent to the rodex_request_items function in Send.pm.
func (rm *RodexManager) RodexRequestItems(mailID1, mailID2 uint32, type_ uint8) error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_request_items")
	if !exists {
		return fmt.Errorf("rodex_request_items packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"mailID1": mailID1,
		"mailID2": mailID2,
		"type":    type_,
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// RodexCancelWriteMail sends a request to cancel writing a mail.
// This is equivalent to the rodex_cancel_write_mail function in Send.pm.
func (rm *RodexManager) RodexCancelWriteMail() error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_cancel_write_mail")
	if !exists {
		return fmt.Errorf("rodex_cancel_write_mail packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// RodexAddItem sends a request to add an item to a mail.
// This is equivalent to the rodex_add_item function in Send.pm.
func (rm *RodexManager) RodexAddItem(ID, amount uint16) error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_add_item")
	if !exists {
		return fmt.Errorf("rodex_add_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":     ID,
		"amount": amount,
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// RodexRemoveItem sends a request to remove an item from a mail.
// This is equivalent to the rodex_remove_item function in Send.pm.
func (rm *RodexManager) RodexRemoveItem(ID, amount uint16) error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_remove_item")
	if !exists {
		return fmt.Errorf("rodex_remove_item packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":     ID,
		"amount": amount,
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// RodexOpenWriteMail sends a request to open a mail to write.
// This is equivalent to the rodex_open_write_mail function in Send.pm.
func (rm *RodexManager) RodexOpenWriteMail(name string) error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_open_write_mail")
	if !exists {
		return fmt.Errorf("rodex_open_write_mail packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"name": []byte(name), // stringToBytes in the original
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// RodexCheckname sends a request to check a name.
// This is equivalent to the rodex_checkname function in Send.pm.
func (rm *RodexManager) RodexCheckname(name string) error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_checkname")
	if !exists {
		return fmt.Errorf("rodex_checkname packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"name": []byte(name), // stringToBytes in the original
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// RodexSendMail sends a request to send a mail.
// This is equivalent to the rodex_send_mail function in Send.pm.
func (rm *RodexManager) RodexSendMail(receiver, sender string, zeny uint32, title, body string, charID uint32) error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_send_mail")
	if !exists {
		return fmt.Errorf("rodex_send_mail packet ID not found")
	}

	// Add null terminator to title and body
	titleBytes := append([]byte(title), 0)
	bodyBytes := append([]byte(body), 0)

	// Create the arguments
	args := map[string]interface{}{
		"receiver":  receiver,
		"sender":    []byte(sender), // stringToBytes in the original
		"zeny1":     zeny,
		"zeny2":     uint32(0),
		"title_len": uint16(len(titleBytes)),
		"body_len":  uint16(len(bodyBytes)),
		"char_id":   charID,
		"title":     titleBytes,
		"body":      bodyBytes,
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// RodexRefreshMaillist sends a request to refresh the mail list.
// This is equivalent to the rodex_refresh_maillist function in Send.pm.
func (rm *RodexManager) RodexRefreshMaillist(type_ uint8, mailID1, mailID2 uint32) error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_refresh_maillist")
	if !exists {
		return fmt.Errorf("rodex_refresh_maillist packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type":           type_,
		"mailID1":        mailID1,
		"mailID2":        mailID2,
		"mailReturnID1":  uint32(0),
		"mailReturnID2":  uint32(0),
		"mailAccountID1": uint32(0),
		"mailAccountID2": uint32(0),
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// RodexReadMail sends a request to read a mail.
// This is equivalent to the rodex_read_mail function in Send.pm.
func (rm *RodexManager) RodexReadMail(type_ uint8, mailID1, mailID2 uint32) error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_read_mail")
	if !exists {
		return fmt.Errorf("rodex_read_mail packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type":    type_,
		"mailID1": mailID1,
		"mailID2": mailID2,
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// RodexNextMaillist sends a request to get the next mail list.
// This is equivalent to the rodex_next_maillist function in Send.pm.
func (rm *RodexManager) RodexNextMaillist(type_ uint8, mailID1, mailID2 uint32) error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_next_maillist")
	if !exists {
		return fmt.Errorf("rodex_next_maillist packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type":    type_,
		"mailID1": mailID1,
		"mailID2": mailID2,
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// RodexOpenMailbox sends a request to open the mailbox.
// This is equivalent to the rodex_open_mailbox function in Send.pm.
func (rm *RodexManager) RodexOpenMailbox(type_ uint8, mailID1, mailID2 uint32) error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_open_mailbox")
	if !exists {
		return fmt.Errorf("rodex_open_mailbox packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type":           type_,
		"mailID1":        mailID1,
		"mailID2":        mailID2,
		"mailReturnID1":  uint32(0),
		"mailReturnID2":  uint32(0),
		"mailAccountID1": uint32(0),
		"mailAccountID2": uint32(0),
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// RodexCloseMailbox sends a request to close the mailbox.
// This is equivalent to the rodex_close_mailbox function in Send.pm.
func (rm *RodexManager) RodexCloseMailbox() error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rodex_close_mailbox")
	if !exists {
		return fmt.Errorf("rodex_close_mailbox packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}
