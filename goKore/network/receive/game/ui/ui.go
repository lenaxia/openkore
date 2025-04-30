package ui

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// UI Types
const (
	BANK_UI       = 0
	STYLIST_UI    = 1
	CAPTCHA_UI    = 2
	MACRO_UI      = 3
	UI_UNUSED     = 4
	TIPBOX_UI     = 5
	RENEWQUEST_UI = 6
	ATTENDANCE_UI = 7
)

// Config Types
const (
	CONFIG_OPEN_EQUIPMENT_WINDOW = 0
	CONFIG_CALL                  = 1
	CONFIG_PET_AUTOFEED          = 2
	CONFIG_HOMUNCULUS_AUTOFEED   = 3
)

// Inventory Expansion Result
const (
	EXPAND_INVENTORY_RESULT_SUCCESS      = 0
	EXPAND_INVENTORY_RESULT_FAILED       = 1
	EXPAND_INVENTORY_RESULT_OTHER_WORK   = 2
	EXPAND_INVENTORY_RESULT_MISSING_ITEM = 3
	EXPAND_INVENTORY_RESULT_MAX_SIZE     = 4
)

// EquipInfo represents equipment information
type EquipInfo struct {
	ID              uint32
	NameID          uint16
	Type            uint8
	Identified      uint8
	TypeEquip       uint16
	Equipped        uint16
	Broken          uint8
	Upgrade         uint8
	Cards           []uint32
	Expire          uint32
	BindOnEquipType uint16
	SpriteID        uint16
	NumOptions      uint8
	Options         []uint32
}

// RouletteInfo represents roulette information
type RouletteInfo struct {
	Serial         uint32
	Result         uint8
	Stage          uint8
	Price          uint8
	AdditionalItem uint16
	Gold           uint32
	Silver         uint32
	Bronze         uint32
}

// AttendanceInfo represents attendance information
type AttendanceInfo struct {
	Data             uint32
	AlreadyRequested bool
	AttendanceCount  int
}

// UIManager handles UI-related packet handling
type UIManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for UI interactions
	progressBarActive bool
	attendanceRewards map[string]interface{}
	roulette          map[string]interface{}
}

// NewUIManager creates a new UI manager
func NewUIManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *UIManager {
	return &UIManager{
		parser:            parser,
		hookManager:       hookManager,
		logger:            logger,
		progressBarActive: false,
		attendanceRewards: make(map[string]interface{}),
		roulette:          make(map[string]interface{}),
	}
}

// RegisterHandlers registers all UI-related packet handlers
func (um *UIManager) RegisterHandlers() {
	// Register show equipment handler
	um.parser.RegisterHandlerFunc("02D7", "show_eq", "v Z24 a*",
		[]string{"type", "name", "equips_info"}, um.HandleShowEq)

	// Register misc config handler
	um.parser.RegisterHandlerFunc("0A95", "misc_config", "B B",
		[]string{"show_eq_flag", "call_flag"}, um.HandleMiscConfig)

	// Register misc config reply handler
	um.parser.RegisterHandlerFunc("02D9", "misc_config_reply", "L L",
		[]string{"type", "flag"}, um.HandleMiscConfigReply)

	// Register show eq msg self handler
	um.parser.RegisterHandlerFunc("02DA", "show_eq_msg_self", "B",
		[]string{"type"}, um.HandleShowEqMsgSelf)

	// Register show script handler
	um.parser.RegisterHandlerFunc("08B3", "show_script", "L Z*",
		[]string{"ID", "message"}, um.HandleShowScript)

	// Register progress bar handler
	um.parser.RegisterHandlerFunc("02F0", "progress_bar", "L",
		[]string{"time"}, um.HandleProgressBar)

	// Register progress bar stop handler
	um.parser.RegisterHandlerFunc("02F2", "progress_bar_stop", "",
		[]string{}, um.HandleProgressBarStop)

	// Register progress bar unit handler
	um.parser.RegisterHandlerFunc("09D1", "progress_bar_unit", "L L",
		[]string{"GID", "time"}, um.HandleProgressBarUnit)

	// Register book read handler
	um.parser.RegisterHandlerFunc("0294", "book_read", "L L",
		[]string{"bookID", "page"}, um.HandleBookRead)

	// Register font handler
	um.parser.RegisterHandlerFunc("0AEB", "font", "L L",
		[]string{"ID", "fontID"}, um.HandleFont)

	// Register open UI handler
	um.parser.RegisterHandlerFunc("0AEF", "open_ui", "B",
		[]string{"type"}, um.HandleOpenUI)

	// Register action UI handler
	um.parser.RegisterHandlerFunc("0AF0", "action_ui", "L L",
		[]string{"type", "data"}, um.HandleActionUI)

	// Register attendance UI handler
	um.parser.RegisterHandlerFunc("0AF1", "attendance_ui", "L",
		[]string{"data"}, um.HandleAttendanceUI)

	// Register load confirm handler
	um.parser.RegisterHandlerFunc("0B02", "load_confirm", "",
		[]string{}, um.HandleLoadConfirm)

	// Register inventory expansion result handler
	um.parser.RegisterHandlerFunc("0B18", "inventory_expansion_result", "W",
		[]string{"result"}, um.HandleInventoryExpansionResult)

	// Register roulette window handler
	um.parser.RegisterHandlerFunc("0A1A", "roulette_window", "L L B B B W L L L",
		[]string{"serial", "result", "stage", "price", "additional_item", "gold", "silver", "bronze"}, um.HandleRouletteWindow)
}

// GetItemName gets the name of an item by its nameID
func (um *UIManager) GetItemName(nameID uint16) string {
	// In a real implementation, this would look up the item name from a database
	// For now, we'll just return a placeholder
	return fmt.Sprintf("Item#%d", nameID)
}

// GetNpcName gets the name of an NPC by its ID
func (um *UIManager) GetNpcName(npcID uint32) string {
	// In a real implementation, this would look up the NPC name from the NPC list
	// For now, we'll just return a placeholder
	return fmt.Sprintf("NPC#%d", npcID)
}

// HandleShowEq handles the show_eq packet (lines 3216-3281)
func (um *UIManager) HandleShowEq(args map[string]interface{}) error {
	// Extract packet data
	equipType, ok := args["type"].(uint16)
	if !ok {
		return fmt.Errorf("invalid type in show_eq packet")
	}

	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in show_eq packet")
	}

	equipsInfo, ok := args["equips_info"].([]byte)
	if !ok {
		return fmt.Errorf("invalid equips_info in show_eq packet")
	}

	// Log message
	um.logger.Info("Showing equipment info for %s", name)

	// Parse equipment info based on packet type
	var itemInfo struct {
		len   int
		types string
		keys  []string
	}

	switch equipType {
	case 0x02D7: // PACKETVER DEFAULT
		itemInfo.len = 26
		itemInfo.types = "a2 v C2 v2 C2 a8 l v"
		itemInfo.keys = []string{"ID", "nameID", "type", "identified", "type_equip", "equipped", "broken", "upgrade", "cards", "expire", "bindOnEquipType"}
	case 0x0906: // PACKETVER >= ?? NOT IMPLEMENTED ON EATHENA BASED EMULATOR
		itemInfo.len = 27
		itemInfo.types = "v2 C v2 C a8 l v2 C"
		itemInfo.keys = []string{"ID", "nameID", "type", "type_equip", "equipped", "upgrade", "cards", "expire", "bindOnEquipType", "sprite_id", "identified"}
	case 0x0859: // PACKETVER >= 20101124
		itemInfo.len = 28
		itemInfo.types = "a2 v C2 v2 C2 a8 l v2"
		itemInfo.keys = []string{"ID", "nameID", "type", "identified", "type_equip", "equipped", "broken", "upgrade", "cards", "expire", "bindOnEquipType", "sprite_id"}
	case 0x0997: // PACKETVER >= 20120925
		itemInfo.len = 31
		itemInfo.types = "a2 v C V2 C a8 l v2 C"
		itemInfo.keys = []string{"ID", "nameID", "type", "type_equip", "equipped", "upgrade", "cards", "expire", "bindOnEquipType", "sprite_id", "identified"}
	case 0x0A2D: // PACKETVER >= 20150226
		itemInfo.len = 57
		itemInfo.types = "a2 v C V2 C a8 l v2 C a25 C"
		itemInfo.keys = []string{"ID", "nameID", "type", "type_equip", "equipped", "upgrade", "cards", "expire", "bindOnEquipType", "sprite_id", "num_options", "options", "identified"}
	case 0x0B03: // PACKETVER >= 20150226
		itemInfo.len = 67
		itemInfo.types = "a2 V C V2 C a16 l v2 C a25 C"
		itemInfo.keys = []string{"ID", "nameID", "type", "type_equip", "equipped", "upgrade", "cards", "expire", "bindOnEquipType", "sprite_id", "num_options", "options", "identified"}
	default:
		return fmt.Errorf("unknown equip type: %d", equipType)
	}

	// Parse equipment info
	equipItems := make([]EquipInfo, 0)
	for i := 0; i < len(equipsInfo); i += itemInfo.len {
		if i+itemInfo.len > len(equipsInfo) {
			break
		}

		// In a real implementation, this would parse the equipment info based on the itemInfo
		// For now, we'll just create a placeholder
		equipItem := EquipInfo{
			ID:              uint32(i),
			NameID:          uint16(i),
			Type:            uint8(i % 256),
			Identified:      1,
			TypeEquip:       uint16(i % 32768),
			Equipped:        uint16(i % 32768),
			Broken:          0,
			Upgrade:         uint8(i % 256),
			Cards:           []uint32{0, 0, 0, 0},
			Expire:          0,
			BindOnEquipType: 0,
			SpriteID:        0,
			NumOptions:      0,
			Options:         []uint32{},
		}

		equipItems = append(equipItems, equipItem)
	}

	// Call hook
	um.hookManager.CallHook("show_eq", map[string]interface{}{
		"type":        equipType,
		"name":        name,
		"equips_info": equipItems,
	})

	return nil
}

// HandleMiscConfig handles the misc_config packet (lines 3287-3321)
func (um *UIManager) HandleMiscConfig(args map[string]interface{}) error {
	// Extract packet data
	showEqFlag, ok := args["show_eq_flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid show_eq_flag in misc_config packet")
	}

	callFlag, ok := args["call_flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid call_flag in misc_config packet")
	}

	// Handle show equipment flag
	if showEqFlag == 1 {
		um.logger.Info("Your Equipment information is now open to the public.")
	} else {
		um.logger.Info("Your Equipment information is now not open to the public.")
	}

	// Handle call flag
	if callFlag == 1 {
		um.logger.Info("Allowed being summoned by skills: Urgent Call, Marriage Skills, etc.")
	} else {
		um.logger.Info("Not Allowed being summoned by skills: Urgent Call, Marriage Skills, etc.")
	}

	// Handle pet autofeed flag if present
	if petAutofeedFlag, ok := args["pet_autofeed_flag"].(uint8); ok {
		if petAutofeedFlag == 1 {
			um.logger.Info("Pet automatic feeding is ON. (Ragexe Client Feature)")
		} else {
			um.logger.Info("Pet automatic feeding is OFF. (Ragexe Client Feature)")
		}
	}

	// Handle homunculus autofeed flag if present
	if homunculusAutofeedFlag, ok := args["homunculus_autofeed_flag"].(uint8); ok {
		if homunculusAutofeedFlag == 1 {
			um.logger.Info("Homunculus automatic feeding is ON. (Ragexe Client Feature)")
		} else {
			um.logger.Info("Homunculus automatic feeding is OFF. (Ragexe Client Feature)")
		}
	}

	// Call hook
	um.hookManager.CallHook("misc_config", map[string]interface{}{
		"show_eq_flag": showEqFlag,
		"call_flag":    callFlag,
	})

	return nil
}

// HandleMiscConfigReply handles the misc_config_reply packet (lines 3333-3363)
func (um *UIManager) HandleMiscConfigReply(args map[string]interface{}) error {
	// Extract packet data
	configType, ok := args["type"].(uint32)
	if !ok {
		return fmt.Errorf("invalid type in misc_config_reply packet")
	}

	flag, ok := args["flag"].(uint32)
	if !ok {
		return fmt.Errorf("invalid flag in misc_config_reply packet")
	}

	// Handle based on config type
	switch configType {
	case CONFIG_OPEN_EQUIPMENT_WINDOW:
		if flag != 0 {
			um.logger.Info("Your Equipment information is now open to the public.")
		} else {
			um.logger.Info("Your Equipment information is now not open to the public.")
		}
	case CONFIG_CALL:
		if flag != 0 {
			um.logger.Info("Allowed being summoned by skills: Urgent Call, Marriage Skills, etc.")
		} else {
			um.logger.Info("Not Allowed being summoned by skills: Urgent Call, Marriage Skills, etc.")
		}
	case CONFIG_PET_AUTOFEED:
		if flag != 0 {
			um.logger.Info("Pet automatic feeding is ON. (Ragexe Client Feature)")
		} else {
			um.logger.Info("Pet automatic feeding is OFF. (Ragexe Client Feature)")
		}
	case CONFIG_HOMUNCULUS_AUTOFEED:
		if flag != 0 {
			um.logger.Info("Homunculus automatic feeding is ON. (Ragexe Client Feature)")
		} else {
			um.logger.Info("Homunculus automatic feeding is OFF. (Ragexe Client Feature)")
		}
	default:
		um.logger.Warning("Unknown Config Type: %d, Flag: %d", configType, flag)
	}

	// Call hook
	um.hookManager.CallHook("misc_config_reply", map[string]interface{}{
		"type": configType,
		"flag": flag,
	})

	return nil
}

// HandleShowEqMsgSelf handles the show_eq_msg_self packet (lines 3365-3372)
func (um *UIManager) HandleShowEqMsgSelf(args map[string]interface{}) error {
	// Extract packet data
	showType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in show_eq_msg_self packet")
	}

	// Handle based on type
	if showType != 0 {
		um.logger.Info("Your Equipment information is now open to the public.")
	} else {
		um.logger.Info("Your Equipment information is now not open to the public.")
	}

	// Call hook
	um.hookManager.CallHook("show_eq_msg_self", map[string]interface{}{
		"type": showType,
	})

	return nil
}

// HandleShowScript handles the show_script packet (lines 3375-3387)
func (um *UIManager) HandleShowScript(args map[string]interface{}) error {
	// Extract packet data
	npcID, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in show_script packet")
	}

	message, ok := args["message"].(string)
	if !ok {
		return fmt.Errorf("invalid message in show_script packet")
	}

	// Get NPC name
	npcName := um.GetNpcName(npcID)

	// Log message
	um.logger.Debug("%s (%d): %s", npcName, npcID, message)

	// Call hook
	um.hookManager.CallHook("show_script", map[string]interface{}{
		"ID":      npcID,
		"message": message,
	})

	return nil
}

// HandleProgressBar handles the progress_bar packet (lines 4576-4588)
func (um *UIManager) HandleProgressBar(args map[string]interface{}) error {
	// Extract packet data
	time, ok := args["time"].(uint32)
	if !ok {
		return fmt.Errorf("invalid time in progress_bar packet")
	}

	// Log message
	um.logger.Info("Progress bar loading (time: %d).", time)

	// Set progress bar active
	um.progressBarActive = true

	// Call hook
	um.hookManager.CallHook("progress_bar", map[string]interface{}{
		"time": time,
	})

	return nil
}

// HandleProgressBarStop handles the progress_bar_stop packet (lines 4590-4593)
func (um *UIManager) HandleProgressBarStop(args map[string]interface{}) error {
	// Log message
	um.logger.Info("Progress bar finished.")

	// Set progress bar inactive
	um.progressBarActive = false

	// Call hook
	um.hookManager.CallHook("progress_bar_stop", map[string]interface{}{})

	return nil
}

// HandleProgressBarUnit handles the progress_bar_unit packet (lines 11241-11244)
func (um *UIManager) HandleProgressBarUnit(args map[string]interface{}) error {
	// Extract packet data
	gid, ok := args["GID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid GID in progress_bar_unit packet")
	}

	time, ok := args["time"].(uint32)
	if !ok {
		return fmt.Errorf("invalid time in progress_bar_unit packet")
	}

	// Log message
	um.logger.Debug("Displays progress bar (GID: %d time: %d)", gid, time)

	// Call hook
	um.hookManager.CallHook("progress_bar_unit", map[string]interface{}{
		"GID":  gid,
		"time": time,
	})

	return nil
}

// HandleBookRead handles the book_read packet (lines 10601-10604)
func (um *UIManager) HandleBookRead(args map[string]interface{}) error {
	// Extract packet data
	bookID, ok := args["bookID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid bookID in book_read packet")
	}

	page, ok := args["page"].(uint32)
	if !ok {
		return fmt.Errorf("invalid page in book_read packet")
	}

	// Log message
	um.logger.Debug("Reading book: %d page: %d", bookID, page)

	// Call hook
	um.hookManager.CallHook("book_read", map[string]interface{}{
		"bookID": bookID,
		"page":   page,
	})

	return nil
}

// HandleFont handles the font packet (lines 10707-10710)
func (um *UIManager) HandleFont(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in font packet")
	}

	fontID, ok := args["fontID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid fontID in font packet")
	}

	// Log message
	um.logger.Debug("Account: %d is using fontID: %d", id, fontID)

	// Call hook
	um.hookManager.CallHook("font", map[string]interface{}{
		"ID":     id,
		"fontID": fontID,
	})

	return nil
}

// HandleOpenUI handles the open_ui packet (lines 11878-11903)
func (um *UIManager) HandleOpenUI(args map[string]interface{}) error {
	// Extract packet data
	uiType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in open_ui packet")
	}

	// Log message
	um.logger.Debug("Received request from server to open UI: %d", uiType)

	// Handle based on UI type
	switch uiType {
	case BANK_UI:
		um.logger.Info("Server requested to open Bank UI.")
	case STYLIST_UI:
		um.logger.Info("Server requested to open Stylist UI.")
	case CAPTCHA_UI:
		um.logger.Info("Server requested to open Captcha UI.")
	case MACRO_UI:
		um.logger.Info("Server requested to open Macro Recorder UI.")
	case UI_UNUSED:
		um.logger.Info("Server requested to open Unused UI.")
	case TIPBOX_UI:
		um.logger.Info("Server requested to open Tip Box UI.")
	case RENEWQUEST_UI:
		um.logger.Info("Server requested to open Quest UI.")
	case ATTENDANCE_UI:
		um.logger.Info("Server requested to open Attendance UI.")
		// Handle attendance UI
		if data, ok := args["data"].(uint32); ok {
			um.HandleAttendanceUI(map[string]interface{}{"data": data})
		}
	default:
		um.logger.Error("Received request from server to open unknown UI: %d", uiType)
	}

	// Call hook
	um.hookManager.CallHook("open_ui", map[string]interface{}{
		"type": uiType,
	})

	return nil
}

// HandleActionUI handles the action_ui packet (lines 11909-11913)
func (um *UIManager) HandleActionUI(args map[string]interface{}) error {
	// Extract packet data
	actionType, ok := args["type"].(uint32)
	if !ok {
		return fmt.Errorf("invalid type in action_ui packet")
	}

	// Log message
	um.logger.Debug("Received request from server to close UI: %d", actionType)

	// Call hook
	um.hookManager.CallHook("action_ui", map[string]interface{}{
		"type": actionType,
	})

	return nil
}

// HandleAttendanceUI handles the attendance_ui packet (lines 11922-11959)
func (um *UIManager) HandleAttendanceUI(args map[string]interface{}) error {
	// Extract packet data
	data, ok := args["data"].(uint32)
	if !ok {
		return fmt.Errorf("invalid data in attendance_ui packet")
	}

	// Parse attendance data
	alreadyRequested := data % 10
	attendanceCount := int(data/10) + 1 - int(alreadyRequested)

	// Log message
	um.logger.Info("Attendance Day: %d", attendanceCount)

	// Call hook
	um.hookManager.CallHook("attendance_ui", map[string]interface{}{
		"data":             data,
		"alreadyRequested": alreadyRequested != 0,
		"attendanceCount":  attendanceCount,
	})

	return nil
}

// HandleLoadConfirm handles the load_confirm packet (lines 12150-12153)
func (um *UIManager) HandleLoadConfirm(args map[string]interface{}) error {
	// Log message
	um.logger.Debug("You are allowed to use Keyboard")

	// Call hook
	um.hookManager.CallHook("load_confirm", map[string]interface{}{})

	return nil
}

// HandleInventoryExpansionResult handles the inventory_expansion_result packet (lines 12163-12170)
func (um *UIManager) HandleInventoryExpansionResult(args map[string]interface{}) error {
	// Extract packet data
	result, ok := args["result"].(uint16)
	if !ok {
		return fmt.Errorf("invalid result in inventory_expansion_result packet")
	}

	// Handle based on result
	switch result {
	case EXPAND_INVENTORY_RESULT_SUCCESS:
		um.logger.Info("You have successfully expanded the possession limit.")
	case EXPAND_INVENTORY_RESULT_FAILED:
		um.logger.Info("Failed to expand the maximum possession limit.")
	case EXPAND_INVENTORY_RESULT_OTHER_WORK:
		um.logger.Info("You cannot expand your inventory at this time.")
	case EXPAND_INVENTORY_RESULT_MISSING_ITEM:
		um.logger.Info("You do not have the required item to expand your inventory.")
	case EXPAND_INVENTORY_RESULT_MAX_SIZE:
		um.logger.Info("You have reached the maximum inventory size.")
	default:
		um.logger.Warning("Unknown inventory expansion result: %d", result)
	}

	// Call hook
	um.hookManager.CallHook("inventory_expansion_result", map[string]interface{}{
		"result": result,
	})

	return nil
}

// HandleRouletteWindow handles the roulette_window packet (lines 12061-12081)
func (um *UIManager) HandleRouletteWindow(args map[string]interface{}) error {
	// Extract packet data
	serial, ok := args["serial"].(uint32)
	if !ok {
		return fmt.Errorf("invalid serial in roulette_window packet")
	}

	result, ok := args["result"].(uint32)
	if !ok {
		return fmt.Errorf("invalid result in roulette_window packet")
	}

	stage, ok := args["stage"].(uint8)
	if !ok {
		return fmt.Errorf("invalid stage in roulette_window packet")
	}

	price, ok := args["price"].(uint8)
	if !ok {
		return fmt.Errorf("invalid price in roulette_window packet")
	}

	additionalItem, ok := args["additional_item"].(uint16)
	if !ok {
		return fmt.Errorf("invalid additional_item in roulette_window packet")
	}

	gold, ok := args["gold"].(uint32)
	if !ok {
		return fmt.Errorf("invalid gold in roulette_window packet")
	}

	silver, ok := args["silver"].(uint32)
	if !ok {
		return fmt.Errorf("invalid silver in roulette_window packet")
	}

	bronze, ok := args["bronze"].(uint32)
	if !ok {
		return fmt.Errorf("invalid bronze in roulette_window packet")
	}

	// Store roulette info
	um.roulette = map[string]interface{}{
		"serial":          serial,
		"result":          result,
		"stage":           stage,
		"price":           price,
		"additional_item": additionalItem,
		"gold":            gold,
		"silver":          silver,
		"bronze":          bronze,
	}

	// Handle based on result
	if result == 1 {
		um.logger.Warning("Roulette: Something went wrong")
		return nil
	} else if result == 2 {
		um.logger.Warning("Roulette: No enough Point (coin) to roll")
		return nil
	}

	// Get result string
	resultLut := []string{"Success", "Failed", "No_Enought_Point", "Losing"}
	resultStr := "Unknown"
	if int(result) < len(resultLut) {
		resultStr = resultLut[result]
	}

	// Log message
	um.logger.Info("[Roulette] - %d", serial)
	um.logger.Info("Result: %s  Row: %d  Column: %d  Bonus Item: %s", resultStr, stage, price, um.GetItemName(additionalItem))
	um.logger.Info("Coins:")
	um.logger.Info("Gold: %d  Silver: %d  Bronze: %d", gold, silver, bronze)

	// Call hook
	um.hookManager.CallHook("roulette_window", map[string]interface{}{
		"serial":          serial,
		"result":          result,
		"stage":           stage,
		"price":           price,
		"additional_item": additionalItem,
		"gold":            gold,
		"silver":          silver,
		"bronze":          bronze,
	})

	return nil
}
