package mail

import (
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// RodexManager handles RODEX mail-related packet handling
type RodexManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for RODEX mail interactions
	rodexList   map[string]interface{}
	rodexWrite  map[string]interface{}
	currentType uint8
}

// NewRodexManager creates a new RODEX mail manager
func NewRodexManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *RodexManager {
	return &RodexManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
		rodexList:   make(map[string]interface{}),
		rodexWrite:  make(map[string]interface{}),
		currentType: 0,
	}
}

// RegisterHandlers registers all RODEX mail-related packet handlers
func (rm *RodexManager) RegisterHandlers() {
	// Register rodex_mail_list handler
	rm.parser.RegisterHandlerFunc("09F0", "rodex_mail_list", "C V2 C V a*",
		[]string{"type", "len", "amount", "isEnd", "isEnd2", "mailList"}, rm.HandleRodexMailList)
	rm.parser.RegisterHandlerFunc("0A7D", "rodex_mail_list", "C V2 C V a*",
		[]string{"type", "len", "amount", "isEnd", "isEnd2", "mailList"}, rm.HandleRodexMailList)
	rm.parser.RegisterHandlerFunc("0AC2", "rodex_mail_list", "C V2 C V a*",
		[]string{"openType", "len", "amount", "isEnd", "isEnd2", "mailList"}, rm.HandleRodexMailList)
	rm.parser.RegisterHandlerFunc("0B5F", "rodex_mail_list", "C V2 C V a*",
		[]string{"openType", "len", "amount", "isEnd", "isEnd2", "mailList"}, rm.HandleRodexMailList)
}

// Helper function to convert bytes to string
func bytesToString(data []byte) string {
	// Find the first null byte
	nullIndex := -1
	for i, b := range data {
		if b == 0 {
			nullIndex = i
			break
		}
	}

	// If no null byte found, use the entire slice
	if nullIndex == -1 {
		nullIndex = len(data)
	}

	// Convert to string
	return string(data[:nullIndex])
}

// Helper function to center a string
func centerString(s string, width int, fill byte) string {
	if len(s) >= width {
		return s
	}

	leftPad := (width - len(s)) / 2
	rightPad := width - len(s) - leftPad

	return strings.Repeat(string(fill), leftPad) + s + strings.Repeat(string(fill), rightPad)
}

// HandleRodexMailList handles the rodex_mail_list packet (lines 8490-8578)
func (rm *RodexManager) HandleRodexMailList(args map[string]interface{}) error {
	// Extract packet data
	mailList, ok := args["mailList"].([]byte)
	if !ok {
		return fmt.Errorf("invalid mailList in rodex_mail_list packet")
	}

	isEnd, ok := args["isEnd"].(uint8)
	if !ok {
		return fmt.Errorf("invalid isEnd in rodex_mail_list packet")
	}

	amount, ok := args["amount"].(uint32)
	if !ok {
		return fmt.Errorf("invalid amount in rodex_mail_list packet")
	}

	// Determine mail info based on switch
	var mailInfo map[string]interface{}
	switch args["switch"].(string) {
	case "0B5F":
		mailInfo = map[string]interface{}{
			"len":   45,
			"types": "C V2 C2 Z24 V v x4",
			"keys":  []string{"openType", "mailID1", "mailID2", "isRead", "attach", "sender", "expireSecconds", "Titlelength"},
		}
	case "0AC2":
		mailInfo = map[string]interface{}{
			"len":   41,
			"types": "C V2 C2 Z24 V v",
			"keys":  []string{"openType", "mailID1", "mailID2", "isRead", "attach", "sender", "expireSecconds", "Titlelength"},
		}
	default: // 09F0, 0A7D
		mailInfo = map[string]interface{}{
			"len":   44,
			"types": "V2 C2 Z24 V2 v",
			"keys":  []string{"mailID1", "mailID2", "isRead", "attach", "sender", "regDateTime", "expireSecconds", "Titlelength"},
		}
	}

	// Update current type if needed
	if args["switch"].(string) == "09F0" || args["switch"].(string) == "0A7D" {
		if attach, ok := args["type"].(uint8); ok {
			rm.currentType = attach
		}
	}

	// Initialize or reset rodexList
	if args["switch"].(string) == "0A7D" || args["switch"].(string) == "0AC2" || args["switch"].(string) == "0B5F" {
		rm.rodexList = make(map[string]interface{})
		rm.rodexList["current_page"] = 0
		rm.rodexList["mails"] = make(map[uint32]map[string]interface{})
	} else {
		if _, ok := rm.rodexList["current_page"]; !ok {
			rm.rodexList["current_page"] = 0
		} else {
			rm.rodexList["current_page"] = rm.rodexList["current_page"].(int) + 1
		}

		// Ensure mails map exists
		if _, ok := rm.rodexList["mails"]; !ok {
			rm.rodexList["mails"] = make(map[uint32]map[string]interface{})
		}
	}

	// Update last page if needed
	if isEnd == 1 {
		rm.rodexList["last_page"] = rm.rodexList["current_page"]
	} else {
		rm.rodexList["mails_per_page"] = amount
	}

	// Format header for display
	msg := centerString(" Rodex Mail Page "+fmt.Sprintf("%d", rm.rodexList["current_page"]), 119, '-') + "\n" +
		" #  ID       From                    Att  New  Expire    Title\n"

	// Process mail list
	mailLen := mailInfo["len"].(int)
	index := 0

	for i := 0; i < len(mailList); {
		if i+mailLen > len(mailList) {
			break
		}

		// Extract mail data (simplified)
		mail := make(map[string]interface{})

		// In a real implementation, this would use proper unpacking based on mailInfo["types"]
		// For simplicity, we'll just extract some basic fields
		if i+4 <= len(mailList) {
			mailID1 := uint32(mailList[i]) | uint32(mailList[i+1])<<8 | uint32(mailList[i+2])<<16 | uint32(mailList[i+3])<<24
			mail["mailID1"] = mailID1

			if i+5 <= len(mailList) {
				mail["isRead"] = uint8(mailList[i+8])
				mail["attach"] = uint8(mailList[i+9])

				// Extract sender name (simplified)
				senderBytes := mailList[i+10 : i+34]
				mail["sender"] = bytesToString(senderBytes)

				// Extract expire seconds
				if i+38 <= len(mailList) {
					expireSeconds := uint32(mailList[i+34]) | uint32(mailList[i+35])<<8 | uint32(mailList[i+36])<<16 | uint32(mailList[i+37])<<24
					mail["expireSecconds"] = expireSeconds
					mail["expireDay"] = int(expireSeconds / 60 / 60 / 24)

					// Extract title length
					if i+40 <= len(mailList) {
						titleLength := uint16(mailList[i+38]) | uint16(mailList[i+39])<<8
						mail["Titlelength"] = titleLength

						// Extract title
						if i+mailLen+int(titleLength) <= len(mailList) {
							titleBytes := mailList[i+mailLen : i+mailLen+int(titleLength)]
							mail["title"] = bytesToString(titleBytes)

							// Debug title extraction
							rm.logger.Debug("Extracted title: '%s' from bytes at offset %d with length %d",
								mail["title"], i+mailLen, titleLength)
						} else {
							mail["title"] = "Unknown" // Default title if we can't extract it
						}
					}
				}
			}

			// Add mail to rodexList
			if mails, ok := rm.rodexList["mails"].(map[uint32]map[string]interface{}); ok {
				mail["page"] = rm.rodexList["current_page"]
				mail["page_index"] = index
				mails[mailID1] = mail
				rm.rodexList["current_page_last_mailID"] = mailID1
			}

			// Format attachment type
			attachType := "-"
			switch mail["attach"].(uint8) {
			case 2:
				attachType = "z" // only zeny
			case 4:
				attachType = "i" // only item
			case 6:
				attachType = "z+i" // zeny + item
			case 12:
				attachType = "gift" // a gift from the admin
			}
			mail["attach"] = attachType

			// Format mail info for display
			isReadStr := "Yes"
			if mail["isRead"].(uint8) != 0 {
				isReadStr = "No"
			}

			msg += fmt.Sprintf("%2d  %7d %-24s %-3s %-3s  %7d  %-50s\n",
				index, mail["mailID1"], mail["sender"], mail["attach"], isReadStr, mail["expireDay"], mail["title"])

			// Increment index and position
			index++
			i += mailLen + int(mail["Titlelength"].(uint16))
		} else {
			break
		}
	}

	msg += strings.Repeat("-", 119) + "\n"
	rm.logger.Info(msg)

	// Call hook
	rm.hookManager.CallHook("rodex_mail_list", map[string]interface{}{
		"mails":        rm.rodexList["mails"],
		"current_page": rm.rodexList["current_page"],
		"last_mailID":  rm.rodexList["current_page_last_mailID"],
		"isEnd":        isEnd,
	})

	return nil
}

// HandleRodexReadMail handles the rodex_read_mail packet (lines 8580-8664)
func (rm *RodexManager) HandleRodexReadMail(args map[string]interface{}) error {
	// Extract packet data
	mailID1, ok := args["mailID1"].(uint32)
	if !ok {
		return fmt.Errorf("invalid mailID1 in rodex_read_mail packet")
	}

	message, ok := args["message"].([]byte)
	if !ok {
		return fmt.Errorf("invalid message in rodex_read_mail packet")
	}

	textLen, ok := args["text_len"].(uint16)
	if !ok {
		return fmt.Errorf("invalid text_len in rodex_read_mail packet")
	}

	zeny1, ok := args["zeny1"].(uint32)
	if !ok {
		return fmt.Errorf("invalid zeny1 in rodex_read_mail packet")
	}

	zeny2, ok := args["zeny2"].(uint16)
	if !ok {
		return fmt.Errorf("invalid zeny2 in rodex_read_mail packet")
	}

	mailType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in rodex_read_mail packet")
	}

	itemCount, ok := args["itemCount"].(uint32)
	if !ok {
		return fmt.Errorf("invalid itemCount in rodex_read_mail packet")
	}

	// Create mail object
	mail := make(map[string]interface{})

	// Extract header length
	headerLen := 24 // Simplified, in a real implementation this would be calculated based on the packet structure

	// Extract mail body
	if int(textLen) <= len(message) {
		mail["body"] = bytesToString(message[:textLen])
	} else {
		mail["body"] = ""
	}

	// Set mail data
	mail["zeny1"] = zeny1
	mail["zeny2"] = zeny2
	mail["type"] = mailType

	// Map mail type to string
	mailTypeStr := "Mail from players"
	switch mailType {
	case 0:
		mailTypeStr = "Mail from players"
	case 1:
		mailTypeStr = "Account mail"
	case 2:
		mailTypeStr = "Return"
	case 3:
		mailTypeStr = "Unset"
	}

	// Initialize items array
	mail["items"] = []map[string]interface{}{}

	// Get mail from rodexList
	var mailSender string
	var mailTitle string
	if mails, ok := rm.rodexList["mails"].(map[uint32]map[string]interface{}); ok {
		if mailData, ok := mails[mailID1]; ok {
			mailSender = mailData["sender"].(string)
			if title, ok := mailData["title"].(string); ok {
				mailTitle = title
			} else {
				mailTitle = "No Title"
			}
		}
	}

	// Format mail header for display
	printMsg := centerString(" Mail "+fmt.Sprintf("%d", mailID1)+" from "+mailSender+" ", 119, '-') + "\n"
	printMsg += fmt.Sprintf("%-12s %s\n", "Mail type:", mailTypeStr)
	printMsg += fmt.Sprintf("%-12s %s\n", "Title:", mailTitle)
	printMsg += "Message:     " + mail["body"].(string) + "\n"
	rm.logger.Info(printMsg)

	// Format mail items for display
	printMsg = fmt.Sprintf("%-12s %d\n", "Item count:", itemCount)
	printMsg += fmt.Sprintf("%-12s %d\n", "Zeny:", zeny1)

	// Process items
	// In a real implementation, this would be based on the server type
	// itemPack := "v2 C3 a8 a4 C a4 a25"
	itemLen := 50 // Simplified, in a real implementation this would be calculated based on the item pack

	index := 0
	for i := headerLen + int(textLen); i < len(message); i += itemLen {
		if i+itemLen > len(message) {
			break
		}

		// Extract item data (simplified)
		item := make(map[string]interface{})

		// In a real implementation, this would use proper unpacking based on itemPack
		// For simplicity, we'll just extract some basic fields
		if i+2 <= len(message) {
			item["amount"] = uint16(message[i]) | uint16(message[i+1])<<8

			if i+4 <= len(message) {
				item["nameID"] = uint16(message[i+2]) | uint16(message[i+3])<<8

				if i+7 <= len(message) {
					item["identified"] = uint8(message[i+4])
					item["broken"] = uint8(message[i+5])
					item["upgrade"] = uint8(message[i+6])

					// Extract cards and other data would go here in a real implementation

					// Set item name (simplified)
					item["name"] = fmt.Sprintf("Item-%d", item["nameID"])

					// Format item display
					display := item["name"].(string)
					if item["amount"].(uint16) > 1 {
						display += fmt.Sprintf(" x %d", item["amount"])
					}

					printMsg += fmt.Sprintf("%3d %s\n", index, display)

					// Add item to mail
					mail["items"] = append(mail["items"].([]map[string]interface{}), item)

					index++
				}
			}
		}
	}

	printMsg += strings.Repeat("-", 119) + "\n"
	rm.logger.Info(printMsg)

	// Update mail in rodexList
	if mails, ok := rm.rodexList["mails"].(map[uint32]map[string]interface{}); ok {
		if mailData, ok := mails[mailID1]; ok {
			mailData["body"] = mail["body"]
			mailData["items"] = mail["items"]
			mailData["zeny1"] = mail["zeny1"]
			mailData["zeny2"] = mail["zeny2"]
			mailData["isRead"] = 1
		}
	}

	// Set current read mail
	rm.rodexList["current_read"] = mailID1

	// Call hook
	rm.hookManager.CallHook("rodex_mail", map[string]interface{}{
		"mailID":    mailID1,
		"from":      mailSender,
		"title":     mailTitle,
		"content":   mail["body"],
		"zeny":      zeny1,
		"itemCount": itemCount,
		"items":     mail["items"],
	})

	return nil
}

// HandleRodexRemoveItem handles the rodex_remove_item packet (lines 8672-8690)
func (rm *RodexManager) HandleRodexRemoveItem(args map[string]interface{}) error {
	// Extract packet data
	result, ok := args["result"].(uint16)
	if !ok {
		return fmt.Errorf("invalid result in rodex_remove_item packet")
	}

	if result == 0 {
		rm.logger.Error("You failed to remove an item from rodex mail.")
		return nil
	}

	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in rodex_remove_item packet")
	}

	amount, ok := args["amount"].(uint8)
	if !ok {
		return fmt.Errorf("invalid amount in rodex_remove_item packet")
	}

	// Get item from rodexWrite
	var rodexItem map[string]interface{}
	if _, ok := rm.rodexWrite["items"]; !ok {
		return fmt.Errorf("rodexWrite items not initialized")
	}

	items, ok := rm.rodexWrite["items"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid rodexWrite items format")
	}

	for _, item := range items {
		if itemID, ok := item["ID"].(uint32); ok && itemID == id {
			rodexItem = item
			break
		}
	}

	if rodexItem == nil {
		return fmt.Errorf("item with ID %d not found in rodexWrite", id)
	}

	// Log item removal
	binID := uint32(0)
	if val, ok := rodexItem["binID"].(uint32); ok {
		binID = val
	}

	itemType := "Unknown"
	if val, ok := rodexItem["type"].(uint8); ok {
		itemType = fmt.Sprintf("%d", val)
	}

	rm.logger.Info("Item removed from rodex mail message: %s (%d) x %d - %s",
		rodexItem["name"], binID, amount, itemType)

	// Update item amount
	if itemAmount, ok := rodexItem["amount"].(uint16); ok {
		rodexItem["amount"] = itemAmount - uint16(amount)
	} else {
		rodexItem["amount"] = uint16(0)
	}

	// Remove item if amount is 0 or less
	if itemAmount, ok := rodexItem["amount"].(uint16); ok && itemAmount <= 0 {
		newItems := make([]map[string]interface{}, 0)
		for _, item := range items {
			if itemID, ok := item["ID"].(uint32); ok && itemID != id {
				newItems = append(newItems, item)
			}
		}
		rm.rodexWrite["items"] = newItems
	}

	return nil
}

// HandleRodexAddItem handles the rodex_add_item packet (lines 8692-8732)
func (rm *RodexManager) HandleRodexAddItem(args map[string]interface{}) error {
	// Extract packet data
	result, ok := args["result"].(uint16)
	if !ok {
		return fmt.Errorf("invalid result in rodex_add_item packet")
	}

	if result == 0 {
		rm.logger.Error("You failed to add an item to rodex mail.")
		return nil
	}

	index, ok := args["index"].(uint16)
	if !ok {
		return fmt.Errorf("invalid index in rodex_add_item packet")
	}

	amount, ok := args["amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid amount in rodex_add_item packet")
	}

	itemID, ok := args["itemID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid itemID in rodex_add_item packet")
	}

	identified, ok := args["identified"].(uint8)
	if !ok {
		return fmt.Errorf("invalid identified in rodex_add_item packet")
	}

	broken, ok := args["broken"].(uint8)
	if !ok {
		return fmt.Errorf("invalid broken in rodex_add_item packet")
	}

	upgrade, ok := args["upgrade"].(uint8)
	if !ok {
		return fmt.Errorf("invalid upgrade in rodex_add_item packet")
	}

	itemType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in rodex_add_item packet")
	}

	// Extract cards (simplified)
	cards := []uint16{0, 0, 0, 0}
	if cardsData, ok := args["cards"].([]uint16); ok && len(cardsData) == 4 {
		cards = cardsData
	}

	// Extract option data (simplified)
	options := []map[string]interface{}{}
	if optionsData, ok := args["options"].([]map[string]interface{}); ok {
		options = optionsData
	}

	// Create item
	item := map[string]interface{}{
		"ID":         uint32(index),
		"nameID":     itemID,
		"amount":     amount,
		"identified": identified,
		"broken":     broken,
		"upgrade":    upgrade,
		"cards":      cards,
		"options":    options,
		"type":       itemType,
		"name":       fmt.Sprintf("Item-%d", itemID), // Simplified, in a real implementation this would be looked up
	}

	// Initialize rodexWrite if needed
	if rm.rodexWrite == nil {
		rm.rodexWrite = make(map[string]interface{})
	}

	// Initialize items array if needed
	if _, ok := rm.rodexWrite["items"]; !ok {
		rm.rodexWrite["items"] = make([]map[string]interface{}, 0)
	}

	// Add item to rodexWrite
	items, ok := rm.rodexWrite["items"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid rodexWrite items format")
	}

	rm.rodexWrite["items"] = append(items, item)

	// Log item addition
	rm.logger.Info("Item added to rodex mail message: %s (%d) x %d - %s",
		item["name"], index, amount, fmt.Sprintf("%d", itemType))

	return nil
}

// HandleRodexCheckPlayer handles the rodex_check_player packet (lines 8748-8777)
func (rm *RodexManager) HandleRodexCheckPlayer(args map[string]interface{}) error {
	// Extract packet data
	char_id, ok := args["char_id"].(uint32)
	if !ok {
		return fmt.Errorf("invalid char_id in rodex_check_player packet")
	}

	result, ok := args["result"].(uint8)
	if !ok {
		return fmt.Errorf("invalid result in rodex_check_player packet")
	}

	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in rodex_check_player packet")
	}

	// Process result
	if result == 0 {
		rm.logger.Error("Player '%s' (char_id: %d) doesn't exist.", name, char_id)

		// Update rodexWrite
		if rm.rodexWrite != nil {
			rm.rodexWrite["receiver"] = ""
			rm.rodexWrite["target_char_id"] = uint32(0)
		}
	} else {
		rm.logger.Info("Player '%s' (char_id: %d) exists.", name, char_id)

		// Update rodexWrite
		if rm.rodexWrite != nil {
			rm.rodexWrite["receiver"] = name
			rm.rodexWrite["target_char_id"] = char_id
		}
	}

	// Call hook
	rm.hookManager.CallHook("rodex_check_player", map[string]interface{}{
		"char_id": char_id,
		"name":    name,
		"exists":  result == 1,
	})

	return nil
}

// HandleRodexWriteResult handles the rodex_write_result packet (lines 8779-8789)
func (rm *RodexManager) HandleRodexWriteResult(args map[string]interface{}) error {
	// Extract packet data
	result, ok := args["result"].(uint8)
	if !ok {
		return fmt.Errorf("invalid result in rodex_write_result packet")
	}

	// Process result
	if result == 0 {
		rm.logger.Error("Failed to send rodex mail.")
	} else {
		rm.logger.Success("Rodex mail sent successfully.")

		// Reset rodexWrite
		rm.rodexWrite = make(map[string]interface{})
		rm.rodexWrite["items"] = make([]map[string]interface{}, 0)
		rm.rodexWrite["receiver"] = ""
		rm.rodexWrite["title"] = ""
		rm.rodexWrite["body"] = ""
		rm.rodexWrite["zeny"] = 0
	}

	// Call hook
	rm.hookManager.CallHook("rodex_write_result", map[string]interface{}{
		"success": result == 1,
	})

	return nil
}

// HandleRodexGetZeny handles the rodex_get_zeny packet (lines 8791-8803)
func (rm *RodexManager) HandleRodexGetZeny(args map[string]interface{}) error {
	// Extract packet data
	mailID1, ok := args["mailID1"].(uint32)
	if !ok {
		return fmt.Errorf("invalid mailID1 in rodex_get_zeny packet")
	}

	mailID2, ok := args["mailID2"].(uint32)
	if !ok {
		return fmt.Errorf("invalid mailID2 in rodex_get_zeny packet")
	}

	zeny, ok := args["zeny"].(uint32)
	if !ok {
		return fmt.Errorf("invalid zeny in rodex_get_zeny packet")
	}

	// Get mail from rodexList
	var mailTitle string
	var mailSender string
	if mails, ok := rm.rodexList["mails"].(map[uint32]map[string]interface{}); ok {
		if mailData, ok := mails[mailID1]; ok {
			if title, ok := mailData["title"].(string); ok {
				mailTitle = title
			} else {
				mailTitle = "No Title"
			}
			if sender, ok := mailData["sender"].(string); ok {
				mailSender = sender
			} else {
				mailSender = "Unknown"
			}

			// Update mail data
			mailData["zeny1"] = uint32(0)
			mailData["zeny2"] = uint16(0)
		}
	}

	// Log zeny retrieval
	rm.logger.Info("Retrieved %d zeny from mail %d from %s (Title: %s)", zeny, mailID1, mailSender, mailTitle)

	// Call hook
	rm.hookManager.CallHook("rodex_get_zeny", map[string]interface{}{
		"mailID1": mailID1,
		"mailID2": mailID2,
		"zeny":    zeny,
	})

	return nil
}

// HandleRodexGetItem handles the rodex_get_item packet (lines 8805-8817)
func (rm *RodexManager) HandleRodexGetItem(args map[string]interface{}) error {
	// Extract packet data
	mailID1, ok := args["mailID1"].(uint32)
	if !ok {
		return fmt.Errorf("invalid mailID1 in rodex_get_item packet")
	}

	mailID2, ok := args["mailID2"].(uint32)
	if !ok {
		return fmt.Errorf("invalid mailID2 in rodex_get_item packet")
	}

	itemCount, ok := args["itemCount"].(uint8)
	if !ok {
		return fmt.Errorf("invalid itemCount in rodex_get_item packet")
	}

	// Get mail from rodexList
	var mailTitle string
	var mailSender string
	if mails, ok := rm.rodexList["mails"].(map[uint32]map[string]interface{}); ok {
		if mailData, ok := mails[mailID1]; ok {
			if title, ok := mailData["title"].(string); ok {
				mailTitle = title
			} else {
				mailTitle = "No Title"
			}
			if sender, ok := mailData["sender"].(string); ok {
				mailSender = sender
			} else {
				mailSender = "Unknown"
			}

			// Update mail data - remove items
			mailData["items"] = make([]map[string]interface{}, 0)
		}
	}

	// Log item retrieval
	rm.logger.Info("Retrieved %d items from mail %d from %s (Title: %s)", itemCount, mailID1, mailSender, mailTitle)

	// Call hook
	rm.hookManager.CallHook("rodex_get_item", map[string]interface{}{
		"mailID1":   mailID1,
		"mailID2":   mailID2,
		"itemCount": itemCount,
	})

	return nil
}

// HandleRodexDelete handles the rodex_delete packet (lines 8819-8831)
func (rm *RodexManager) HandleRodexDelete(args map[string]interface{}) error {
	// Extract packet data
	mailID1, ok := args["mailID1"].(uint32)
	if !ok {
		return fmt.Errorf("invalid mailID1 in rodex_delete packet")
	}

	mailID2, ok := args["mailID2"].(uint32)
	if !ok {
		return fmt.Errorf("invalid mailID2 in rodex_delete packet")
	}

	// Get mail from rodexList
	var mailTitle string
	var mailSender string
	if mails, ok := rm.rodexList["mails"].(map[uint32]map[string]interface{}); ok {
		if mailData, ok := mails[mailID1]; ok {
			if title, ok := mailData["title"].(string); ok {
				mailTitle = title
			} else {
				mailTitle = "No Title"
			}
			if sender, ok := mailData["sender"].(string); ok {
				mailSender = sender
			} else {
				mailSender = "Unknown"
			}

			// Delete mail
			delete(mails, mailID1)
		}
	}

	// Log mail deletion
	rm.logger.Info("Mail %d from %s (Title: %s) has been deleted.", mailID1, mailSender, mailTitle)

	// Call hook
	rm.hookManager.CallHook("rodex_delete", map[string]interface{}{
		"mailID1": mailID1,
		"mailID2": mailID2,
	})

	return nil
}

// HandleRodexOpenWrite handles the rodex_open_write packet (lines 8734-8746)
func (rm *RodexManager) HandleRodexOpenWrite(args map[string]interface{}) error {
	// Extract packet data
	result, ok := args["result"].(uint8)
	if !ok {
		return fmt.Errorf("invalid result in rodex_open_write packet")
	}

	if result != 1 {
		rm.logger.Error("You failed to open rodex write window.")
		return nil
	}

	// Initialize rodexWrite
	rm.rodexWrite = make(map[string]interface{})
	rm.rodexWrite["items"] = make([]map[string]interface{}, 0)
	rm.rodexWrite["receiver"] = ""
	rm.rodexWrite["title"] = ""
	rm.rodexWrite["body"] = ""
	rm.rodexWrite["zeny"] = 0

	rm.logger.Info("Rodex write window opened.")

	// Call hook
	rm.hookManager.CallHook("rodex_open_write", nil)

	return nil
}

// HandleUnreadRodex handles the unread_rodex packet (lines 8666-8670)
func (rm *RodexManager) HandleUnreadRodex(args map[string]interface{}) error {
	rm.logger.Info("You have new unread rodex mails.")

	// Call hook
	rm.hookManager.CallHook("rodex_unread_mail", nil)

	return nil
}
