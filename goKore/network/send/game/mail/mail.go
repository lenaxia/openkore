// Package mail provides mail system-related packet sending functionality.
package mail

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// MailManager handles mail system-related packet sending.
type MailManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewMailManager creates a new mail system manager.
func NewMailManager(baseSend core.Send) *MailManager {
	return &MailManager{
		baseSend: baseSend,
	}
}

// SendMailboxOpen sends a request to open the mailbox.
// This is equivalent to the sendMailboxOpen function in Send.pm.
func (mm *MailManager) SendMailboxOpen() error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("mailbox_open")
	if !exists {
		return fmt.Errorf("mailbox_open packet ID not found")
	}

	// Create the arguments (empty for this packet)
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendMailRead sends a request to read a mail.
// This is equivalent to the sendMailRead function in Send.pm.
func (mm *MailManager) SendMailRead(mailID uint32) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("mail_read")
	if !exists {
		return fmt.Errorf("mail_read packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"mailID": mailID,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendMailDelete sends a request to delete a mail.
// This is equivalent to the sendMailDelete function in Send.pm.
func (mm *MailManager) SendMailDelete(mailID uint32) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("mail_delete")
	if !exists {
		return fmt.Errorf("mail_delete packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"mailID": mailID,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendMailGetAttach sends a request to get a mail attachment.
// This is equivalent to the sendMailGetAttach function in Send.pm.
func (mm *MailManager) SendMailGetAttach(mailID uint32) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("mail_attachment_get")
	if !exists {
		return fmt.Errorf("mail_attachment_get packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"mailID": mailID,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendMailOperateWindow sends a request to operate the mail window.
// This is equivalent to the sendMailOperateWindow function in Send.pm.
func (mm *MailManager) SendMailOperateWindow(flag uint8) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("mail_remove")
	if !exists {
		return fmt.Errorf("mail_remove packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"flag": flag,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendMailSetAttach sends a request to set a mail attachment.
// This is equivalent to the sendMailSetAttach function in Send.pm.
func (mm *MailManager) SendMailSetAttach(amount uint32, ID uint16) error {
	// Before setting an attachment, we must remove any zeny/item that was attached but the mail wasn't sent
	// Otherwise the attachment will be lost
	if ID != 0 {
		err := mm.SendMailOperateWindow(1)
		if err != nil {
			return err
		}
	} else {
		err := mm.SendMailOperateWindow(2)
		if err != nil {
			return err
		}
	}

	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("mail_attachment_set")
	if !exists {
		return fmt.Errorf("mail_attachment_set packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID":     ID,
		"amount": amount,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// ReconstructMailSend reconstructs the mail send packet.
// This is equivalent to the reconstruct_mail_send function in Send.pm.
func ReconstructMailSend(args map[string]interface{}) {
	body, ok := args["body"].(string)
	if !ok {
		return
	}

	bodyLen, ok := args["body_len"].(int)
	if !ok {
		return
	}

	// Convert body to bytes with null terminator
	bodyBytes := make([]byte, bodyLen+1)
	copy(bodyBytes, []byte(body))
	bodyBytes[bodyLen] = 0

	args["body"] = bodyBytes
}

// SendMailSend sends a request to send a mail.
// This is equivalent to the sendMailSend function in Send.pm.
func (mm *MailManager) SendMailSend(receiver, title, message string) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("mail_send")
	if !exists {
		return fmt.Errorf("mail_send packet ID not found")
	}

	// Limit message length to 255 characters
	messageLen := len(message)
	if messageLen > 255 {
		messageLen = 255
		message = message[:255]
	}

	// Create the arguments
	args := map[string]interface{}{
		"recipient": []byte(receiver), // stringToBytes in the original
		"title":     []byte(title),    // stringToBytes in the original
		"body_len":  messageLen,
		"body":      message,
	}

	// Reconstruct the body
	ReconstructMailSend(args)

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendMailReturn sends a request to return a mail.
// This is equivalent to the sendMailReturn function in Send.pm.
func (mm *MailManager) SendMailReturn(mailID, sender uint32) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("mail_return")
	if !exists {
		return fmt.Errorf("mail_return packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"mailID": mailID,
		"sender": sender,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}
