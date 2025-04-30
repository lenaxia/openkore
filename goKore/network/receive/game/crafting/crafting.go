package crafting

import (
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// CraftingManager handles crafting-related packet handling
type CraftingManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for crafting interactions
	makableList   []uint16
	arrowCraftIDs []uint16
	cookingList   []uint16
	cookingType   uint8
	repairList    []map[string]interface{}
	identifyIDs   []uint16
	selectedCraft uint16
}

// NewCraftingManager creates a new crafting manager
func NewCraftingManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *CraftingManager {
	return &CraftingManager{
		parser:        parser,
		hookManager:   hookManager,
		logger:        logger,
		makableList:   make([]uint16, 0),
		arrowCraftIDs: make([]uint16, 0),
		cookingList:   make([]uint16, 0),
		cookingType:   0,
		repairList:    make([]map[string]interface{}, 0),
		identifyIDs:   make([]uint16, 0),
		selectedCraft: 0,
	}
}

// RegisterHandlers registers all crafting-related packet handlers
func (cm *CraftingManager) RegisterHandlers() {
	// Register makable item list handler
	cm.parser.RegisterHandlerFunc("018D", "makable_item_list", "v Z*",
		[]string{"len", "item_list"}, cm.HandleMakableItemList)

	// Register arrowcraft list handler
	cm.parser.RegisterHandlerFunc("01AD", "arrowcraft_list", "v Z*",
		[]string{"len", "RAW_MSG"}, cm.HandleArrowcraftList)

	// Register cooking list handler
	cm.parser.RegisterHandlerFunc("025A", "cooking_list", "C v Z*",
		[]string{"type", "len", "item_list"}, cm.HandleCookingList)

	// Register repair list handler
	cm.parser.RegisterHandlerFunc("01FC", "repair_list", "v Z*",
		[]string{"len", "RAW_MSG"}, cm.HandleRepairList)

	// Register repair result handler
	cm.parser.RegisterHandlerFunc("01FE", "repair_result", "v C",
		[]string{"index", "flag"}, cm.HandleRepairResult)

	// Register identify list handler
	cm.parser.RegisterHandlerFunc("0177", "identify_list", "v Z*",
		[]string{"len", "RAW_MSG"}, cm.HandleIdentifyList)

	// Register identify result handler
	cm.parser.RegisterHandlerFunc("0179", "identify", "v C",
		[]string{"ID", "flag"}, cm.HandleIdentify)
}

// HandleMakableItemList handles the makable_item_list packet (lines 5000-5019)
func (cm *CraftingManager) HandleMakableItemList(args map[string]interface{}) error {
	// Extract packet data
	itemList, ok := args["item_list"].([]byte)
	if !ok {
		return fmt.Errorf("invalid item_list in makable_item_list packet")
	}

	// Clear makable list
	cm.makableList = make([]uint16, 0)

	// Define packet format (simplified)
	// In a real implementation, this would depend on the server type
	packLen := 8 // Length of the packed data (v4)

	// Process makable items
	msg := centerString(" Create Item List ", 79, '-') + "\n"

	for i := 0; i < len(itemList); i += packLen {
		if i+packLen <= len(itemList) {
			// Extract item data (simplified)
			// In a real implementation, this would use proper unpacking
			nameID := uint16(itemList[i]) | uint16(itemList[i+1])<<8
			// material_1, material_2, material_3 are also in the packet but not used here

			// Add item to makable list
			cm.makableList = append(cm.makableList, nameID)

			// Format item info for display
			itemName := fmt.Sprintf("Item-%d", nameID) // In a real implementation, this would use a function to get the item name
			msg += fmt.Sprintf("%2d %-50s (%6d)\n", len(cm.makableList)-1, itemName, nameID)
		}
	}

	msg += strings.Repeat("-", 79) + "\n"
	cm.logger.Info(msg)
	cm.logger.Info("You can now use the 'create' command.")

	// Call hook
	cm.hookManager.CallHook("makable_item_list", map[string]interface{}{
		"item_list": cm.makableList,
	})

	return nil
}

// HandleArrowcraftList handles the arrowcraft_list packet (lines 10184-10215)
func (cm *CraftingManager) HandleArrowcraftList(args map[string]interface{}) error {
	// Extract packet data
	msg, ok := args["RAW_MSG"].([]byte)
	if !ok {
		return fmt.Errorf("invalid RAW_MSG in arrowcraft_list packet")
	}

	msgSize, ok := args["len"].(uint16)
	if !ok {
		return fmt.Errorf("invalid len in arrowcraft_list packet")
	}

	// Clear arrowcraft IDs
	cm.arrowCraftIDs = make([]uint16, 0)
	cm.selectedCraft = 0

	// Process arrowcraft items
	for i := 4; i < int(msgSize); i += 2 {
		if i+2 <= len(msg) {
			// Extract item data
			id := uint16(msg[i]) | uint16(msg[i+1])<<8

			// Add item to arrowcraft IDs
			cm.arrowCraftIDs = append(cm.arrowCraftIDs, id)
		}
	}

	// Check if this is for poisoning weapon
	if cm.selectedCraft == 2027 { // GC_POISONINGWEAPON
		cm.logger.Info("Received Possible Poison List - type 'poison'")
	} else {
		cm.logger.Info("Received Possible Item List - type 'arrowcraft' or 'poison'")
	}

	// Call hook
	cm.hookManager.CallHook("arrowcraft_list", map[string]interface{}{
		"item_list": cm.arrowCraftIDs,
	})

	return nil
}

// HandleCookingList handles the cooking_list packet (lines 9137-9159)
func (cm *CraftingManager) HandleCookingList(args map[string]interface{}) error {
	// Extract packet data
	itemList, ok := args["item_list"].([]byte)
	if !ok {
		return fmt.Errorf("invalid item_list in cooking_list packet")
	}

	cookingType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in cooking_list packet")
	}

	// Clear cooking list
	cm.cookingList = make([]uint16, 0)
	cm.cookingType = cookingType

	// Process cooking items
	msg := centerString(" Cooking List ", 79, '-') + "\n"

	for i := 0; i < len(itemList); i += 2 {
		if i+2 <= len(itemList) {
			// Extract item data
			nameID := uint16(itemList[i]) | uint16(itemList[i+1])<<8

			// Add item to cooking list
			cm.cookingList = append(cm.cookingList, nameID)

			// Format item info for display
			itemName := fmt.Sprintf("Item-%d", nameID) // In a real implementation, this would use a function to get the item name
			msg += fmt.Sprintf("%2d %-50s\n", len(cm.cookingList)-1, itemName)
		}
	}

	msg += strings.Repeat("-", 79) + "\n"
	cm.logger.Info(msg)
	cm.logger.Info("You can now use the 'cook' command.")

	// Call hook
	cm.hookManager.CallHook("cooking_list", map[string]interface{}{
		"cooking_list": cm.cookingList,
	})

	return nil
}

// HandleRepairList handles the repair_list packet (lines 11265-11318)
func (cm *CraftingManager) HandleRepairList(args map[string]interface{}) error {
	// Extract packet data
	msg, ok := args["RAW_MSG"].([]byte)
	if !ok {
		return fmt.Errorf("invalid RAW_MSG in repair_list packet")
	}

	msgSize, ok := args["len"].(uint16)
	if !ok {
		return fmt.Errorf("invalid len in repair_list packet")
	}

	// Clear repair list
	cm.repairList = make([]map[string]interface{}, 0)

	// Process repair items
	msg1 := centerString(" Repair List ", 80, '-') + "\n" +
		"   # Short name                     Full name\n"
	var msg2 string

	for i := 4; i < int(msgSize); i += 13 {
		if i+13 <= len(msg) {
			// Extract item data
			index := uint16(msg[i]) | uint16(msg[i+1])<<8
			nameID := uint16(msg[i+2]) | uint16(msg[i+3])<<8
			upgrade := uint8(msg[i+4])
			// cards is also in the packet but not used here

			// Create repair item
			repairItem := map[string]interface{}{
				"index":   index,
				"nameID":  nameID,
				"upgrade": upgrade,
				"name":    fmt.Sprintf("Item-%d", nameID), // In a real implementation, this would use a function to get the item name
			}

			// Add item to repair list
			for len(cm.repairList) <= int(index) {
				cm.repairList = append(cm.repairList, nil)
			}
			cm.repairList[index] = repairItem

			// Format item info for display
			shortName := fmt.Sprintf("Item-%d", nameID) // In a real implementation, this would use a function to get the item name
			msg2 += fmt.Sprintf("%4d %-30s %s\n", index, shortName, repairItem["name"])
		}
	}

	msg2 += strings.Repeat("-", 80) + "\n"
	cm.logger.Info(msg1 + msg2)

	return nil
}

// HandleRepairResult handles the repair_result packet (lines 11327-11339)
func (cm *CraftingManager) HandleRepairResult(args map[string]interface{}) error {
	// Extract packet data
	index, ok := args["index"].(uint16)
	if !ok {
		return fmt.Errorf("invalid index in repair_result packet")
	}

	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in repair_result packet")
	}

	// Adjust index
	index -= 2

	// Check if repair item exists
	if int(index) >= len(cm.repairList) || cm.repairList[index] == nil {
		return fmt.Errorf("invalid repair item index: %d", index)
	}

	// Process repair result
	if flag != 0 {
		cm.logger.Error("Repair of %s failed.", cm.repairList[index]["name"])
	} else {
		cm.logger.Info("Successfully repaired '%s'.", cm.repairList[index]["name"])
	}

	// Clear repair list
	cm.repairList = make([]map[string]interface{}, 0)

	return nil
}

// HandleIdentifyList handles the identify_list packet (lines 6900-6915)
func (cm *CraftingManager) HandleIdentifyList(args map[string]interface{}) error {
	// Extract packet data
	msg, ok := args["RAW_MSG"].([]byte)
	if !ok {
		return fmt.Errorf("invalid RAW_MSG in identify_list packet")
	}

	msgSize, ok := args["len"].(uint16)
	if !ok {
		return fmt.Errorf("invalid len in identify_list packet")
	}

	// Clear identify IDs
	cm.identifyIDs = make([]uint16, 0)

	// Process identify items
	for i := 4; i < int(msgSize); i += 2 {
		if i+2 <= len(msg) {
			// Extract item data
			index := uint16(msg[i]) | uint16(msg[i+1])<<8

			// Add item to identify IDs
			cm.identifyIDs = append(cm.identifyIDs, index)
		}
	}

	cm.logger.Info("Received Possible Identify List (%d item(s)) - type 'identify'", len(cm.identifyIDs))

	return nil
}

// HandleIdentify handles the identify packet (lines 6919-6930)
func (cm *CraftingManager) HandleIdentify(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid ID in identify packet")
	}

	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in identify packet")
	}

	// Process identify result
	if flag == 0 {
		cm.logger.Info("Item Identified: Item-%d (%d)", id, id)
	} else {
		cm.logger.Error("Item Appraisal has failed.")
	}

	// Clear identify IDs
	cm.identifyIDs = make([]uint16, 0)

	return nil
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
