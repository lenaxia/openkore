package mail

import (
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// MailManager handles regular mail-related packet handling
type MailManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for mail interactions
	mailList map[uint32]map[string]interface{}
}

// NewMailManager creates a new mail manager
func NewMailManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *MailManager {
	return &MailManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
		mailList:    make(map[uint32]map[string]interface{}),
	}
}

// HandleMailDelete handles the mail_delete packet (lines 10729-10737)
func (mm *MailManager) HandleMailDelete(args map[string]interface{}) error {
	// Extract packet data
	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in mail_delete packet")
	}

	mailID, ok := args["mailID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid mailID in mail_delete packet")
	}

	if fail != 0 {
		mm.logger.Info("Failed to delete mail with ID: %d.", mailID)
	} else {
		mm.logger.Info("Succeeded to delete mail with ID: %d.", mailID)

		// Remove mail from mailList if it exists
		delete(mm.mailList, mailID)
	}

	// Call hook
	mm.hookManager.CallHook("mail_delete", map[string]interface{}{
		"mailID": mailID,
		"fail":   fail != 0,
	})

	return nil
}

// HandleMailWindow handles the mail_window packet (lines 10739-10747)
func (mm *MailManager) HandleMailWindow(args map[string]interface{}) error {
	// Extract packet data
	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in mail_window packet")
	}

	if flag != 0 {
		mm.logger.Info("Mail window is now closed.")
	} else {
		mm.logger.Info("Mail window is now opened.")
	}

	// Call hook
	mm.hookManager.CallHook("mail_window", map[string]interface{}{
		"flag": flag != 0,
	})

	return nil
}

// HandleMailReturn handles the mail_return packet (lines 10749-10754)
func (mm *MailManager) HandleMailReturn(args map[string]interface{}) error {
	// Extract packet data
	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in mail_return packet")
	}

	mailID, ok := args["mailID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid mailID in mail_return packet")
	}

	if fail != 0 {
		mm.logger.Error("The mail with ID: %d does not exist.", mailID)
	} else {
		mm.logger.Info("The mail with ID: %d is returned to the sender.", mailID)
	}

	// Call hook
	mm.hookManager.CallHook("mail_return", map[string]interface{}{
		"mailID": mailID,
		"fail":   fail != 0,
	})

	return nil
}

// HandleMailRead handles the mail_read packet (lines 10756-10778)
func (mm *MailManager) HandleMailRead(args map[string]interface{}) error {
	// Extract packet data
	mailID, ok := args["mailID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid mailID in mail_read packet")
	}

	nameID, ok := args["nameID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid nameID in mail_read packet")
	}

	title, ok := args["title"].(string)
	if !ok {
		return fmt.Errorf("invalid title in mail_read packet")
	}

	sender, ok := args["sender"].(string)
	if !ok {
		return fmt.Errorf("invalid sender in mail_read packet")
	}

	message, ok := args["message"].(string)
	if !ok {
		return fmt.Errorf("invalid message in mail_read packet")
	}

	amount, _ := args["amount"].(uint16)  // Optional
	zeny, _ := args["zeny"].(uint32)      // Optional
	upgrade, _ := args["upgrade"].(uint8) // Optional
	cards, _ := args["cards"].([]uint16)  // Optional
	broken, _ := args["broken"].(uint8)   // Optional

	// Create item info
	item := map[string]interface{}{
		"nameID":  nameID,
		"upgrade": upgrade,
		"cards":   cards,
		"broken":  broken,
		"amount":  amount,
	}

	// Store mail data
	mailData := map[string]interface{}{
		"mailID":  mailID,
		"title":   title,
		"sender":  sender,
		"message": message,
		"item":    item,
		"zeny":    zeny,
	}

	// Add to mailList
	mm.mailList[mailID] = mailData

	// Log mail content
	mm.logger.Info("Mail from %s: %s", sender, title)
	mm.logger.Info("Message: %s", message)
	if nameID > 0 {
		mm.logger.Info("Item: %d x %d", nameID, amount)
	}
	if zeny > 0 {
		mm.logger.Info("Zeny: %d", zeny)
	}

	// Call hook
	mm.hookManager.CallHook("mail_read", mailData)

	return nil
}

// HandleMailRefreshInbox handles the mail_refreshinbox packet (lines 10780-10822)
func (mm *MailManager) HandleMailRefreshInbox(args map[string]interface{}) error {
	// Extract packet data
	count, ok := args["count"].(uint32)
	if !ok {
		return fmt.Errorf("invalid count in mail_refreshinbox packet")
	}

	if count == 0 {
		mm.logger.Info("There is no mail in your inbox.")
		return nil
	}

	// Clear existing mail list
	mm.mailList = make(map[uint32]map[string]interface{})

	// Process mail list from raw message
	rawMsg, ok := args["RAW_MSG"].([]byte)
	if !ok {
		return fmt.Errorf("invalid RAW_MSG in mail_refreshinbox packet")
	}

	mm.logger.Info("You've got %d mail in your Mailbox.", count)

	// Parse mail entries
	for i := 0; i < int(count); i++ {
		offset := 8 + i*73 // Starting offset + mail entry size

		// Ensure we have enough data
		if offset+73 > len(rawMsg) {
			return fmt.Errorf("mail_refreshinbox packet data too short")
		}

		// Extract mail data
		mailID := uint32(rawMsg[offset]) | uint32(rawMsg[offset+1])<<8 | uint32(rawMsg[offset+2])<<16 | uint32(rawMsg[offset+3])<<24
		title := string(rawMsg[offset+4 : offset+44])
		read := rawMsg[offset+44] != 0
		sender := string(rawMsg[offset+45 : offset+69])
		timestamp := uint32(rawMsg[offset+69]) | uint32(rawMsg[offset+70])<<8 | uint32(rawMsg[offset+71])<<16 | uint32(rawMsg[offset+72])<<24

		// Clean up strings
		title = strings.TrimRight(title, "\x00")
		sender = strings.TrimRight(sender, "\x00")

		// Store mail data
		mm.mailList[mailID] = map[string]interface{}{
			"mailID":    mailID,
			"title":     title,
			"read":      read,
			"sender":    sender,
			"timestamp": timestamp,
		}

		mm.logger.Info("Mail #%d: %s from %s", i+1, title, sender)
	}

	// Call hook
	mm.hookManager.CallHook("mail_refreshinbox", map[string]interface{}{
		"count":    count,
		"mailList": mm.mailList,
	})

	return nil
}

// HandleMailGetAttachment handles the mail_getattachment packet (lines 10824-10833)
func (mm *MailManager) HandleMailGetAttachment(args map[string]interface{}) error {
	// Extract packet data
	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in mail_getattachment packet")
	}

	if fail == 0 {
		mm.logger.Info("Successfully added attachment to inventory.")
	} else if fail == 2 {
		mm.logger.Error("Failed to get the attachment to inventory due to your weight.")
	} else {
		mm.logger.Error("Failed to get the attachment to inventory.")
	}

	// Call hook
	mm.hookManager.CallHook("mail_getattachment", map[string]interface{}{
		"fail": fail,
	})

	return nil
}

// HandleMailSetAttachment handles the mail_setattachment packet (lines 10835-10862)
func (mm *MailManager) HandleMailSetAttachment(args map[string]interface{}) error {
	// Extract packet data
	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in mail_setattachment packet")
	}

	id, _ := args["ID"].(uint32)         // Optional, could be item ID or 0 for zeny
	amount, _ := args["amount"].(uint32) // Optional

	if fail != 0 {
		if id > 0 {
			mm.logger.Info("Failed to attach item: %d.", id)
		} else {
			mm.logger.Info("Failed to attach zeny.")
		}
	} else {
		if id > 0 {
			mm.logger.Info("Succeeded to attach item: %d.", id)
		} else {
			mm.logger.Info("Succeeded to attach zeny: %d.", amount)
		}
	}

	// Call hook
	mm.hookManager.CallHook("mail_setattachment", map[string]interface{}{
		"fail":   fail != 0,
		"ID":     id,
		"amount": amount,
	})

	return nil
}

// HandleMailSend handles the mail_send packet (lines 10864-10869)
func (mm *MailManager) HandleMailSend(args map[string]interface{}) error {
	// Extract packet data
	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in mail_send packet")
	}

	if fail != 0 {
		mm.logger.Error("Failed to send mail, the recipient does not exist.")
	} else {
		mm.logger.Info("Mail sent successfully.")
	}

	// Call hook
	mm.hookManager.CallHook("mail_send", map[string]interface{}{
		"fail": fail != 0,
	})

	return nil
}

// HandleMailNew handles the mail_new packet (lines 10871-10874)
func (mm *MailManager) HandleMailNew(args map[string]interface{}) error {
	// Extract packet data
	sender, ok := args["sender"].(string)
	if !ok {
		return fmt.Errorf("invalid sender in mail_new packet")
	}

	title, ok := args["title"].(string)
	if !ok {
		return fmt.Errorf("invalid title in mail_new packet")
	}

	mm.logger.Info("New mail from sender: %s titled: %s.", sender, title)

	// Call hook
	mm.hookManager.CallHook("mail_new", map[string]interface{}{
		"sender": sender,
		"title":  title,
	})

	return nil
}
