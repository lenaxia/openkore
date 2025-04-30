package refining

import (
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// RefiningManager handles refining-related packet handling
type RefiningManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for refining interactions
	refineUI    map[string]interface{}
	refineList  []uint16
	upgradeList []uint16
}

// NewRefiningManager creates a new refining manager
func NewRefiningManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *RefiningManager {
	return &RefiningManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
		refineUI:    make(map[string]interface{}),
		refineList:  make([]uint16, 0),
		upgradeList: make([]uint16, 0),
	}
}

// RegisterHandlers registers all refining-related packet handlers
func (rm *RefiningManager) RegisterHandlers() {
	// Register refineui opened handler
	rm.parser.RegisterHandlerFunc("0AA1", "refineui_opened", "",
		[]string{}, rm.HandleRefineUIOpened)

	// Register refineui info handler
	rm.parser.RegisterHandlerFunc("0AA2", "refineui_info", "v v C a*",
		[]string{"index", "bless", "len", "materials"}, rm.HandleRefineUIInfo)

	// Register refine result handler
	rm.parser.RegisterHandlerFunc("0188", "refine_result", "C v",
		[]string{"fail", "nameID"}, rm.HandleRefineResult)

	// Register refine status handler
	rm.parser.RegisterHandlerFunc("0AFB", "refine_status", "C v Z24",
		[]string{"status", "refine_level", "name"}, rm.HandleRefineStatus)

	// Register upgrade list handler
	rm.parser.RegisterHandlerFunc("0146", "upgrade_list", "v a*",
		[]string{"len", "item_list"}, rm.HandleUpgradeList)

	// Register upgrade message handler
	rm.parser.RegisterHandlerFunc("0223", "upgrade_message", "C v",
		[]string{"type", "itemID"}, rm.HandleUpgradeMessage)
}

// HandleRefineUIOpened handles the refineui_opened packet (lines 7885-7889)
func (rm *RefiningManager) HandleRefineUIOpened(args map[string]interface{}) error {
	// Set refineUI open flag
	rm.refineUI["open"] = true

	rm.logger.Info("RefineUI is opened. Type 'i' to check equipment and its index. To continue: refineui select [ItemIdx]")

	return nil
}

// HandleRefineUIInfo handles the refineui_info packet (lines 7895-7936)
func (rm *RefiningManager) HandleRefineUIInfo(args map[string]interface{}) error {
	// Extract packet data
	index, ok := args["index"].(uint16)
	if !ok {
		return fmt.Errorf("invalid index in refineui_info packet")
	}

	bless, ok := args["bless"].(uint8)
	if !ok {
		return fmt.Errorf("invalid bless in refineui_info packet")
	}

	materials, ok := args["materials"].([]byte)
	if !ok {
		return fmt.Errorf("invalid materials in refineui_info packet")
	}

	// Check if we have valid data
	if len(materials) > 0 {
		// Store refine info
		rm.refineUI["itemIndex"] = index
		rm.refineUI["bless"] = bless

		// Get item name (simplified)
		itemName := fmt.Sprintf("Item-%d", index)

		// Display refine info
		rm.logger.Info("========= RefineUI Info =========")
		rm.logger.Info("Target Equip:\n- Index: %d\n- Name: %s", index, itemName)
		rm.logger.Info("Blacksmith Blessing:\n- Needed: %d\n- Owned: %d", bless, 0) // In a real implementation, this would check inventory

		// Process materials
		msg := centerString(" Possible Materials ", 53, '-') + "\n" +
			"Mat_ID      %           Zeny        Material                        \n"

		// Parse materials (simplified)
		// In a real implementation, this would use proper unpacking
		for i := 0; i < len(materials); i += 7 {
			if i+7 <= len(materials) {
				nameID := uint16(materials[i]) | uint16(materials[i+1])<<8
				chance := uint8(materials[i+2])
				zeny := uint32(materials[i+3]) | uint32(materials[i+4])<<8 | uint32(materials[i+5])<<16 | uint32(materials[i+6])<<24

				// Format material info
				materialName := fmt.Sprintf("Item-%d", nameID)
				myMatCount := fmt.Sprintf("0 ea %s", materialName)

				msg += fmt.Sprintf("%8d %5d %14d   %s\n", nameID, chance, zeny, myMatCount)
			}
		}

		msg += strings.Repeat("-", 53) + "\n"
		rm.logger.Info(msg)
		rm.logger.Info("Continue: refineui refine %d [Mat_ID] [catalyst_toggle] to continue.", index)
	} else {
		rm.logger.Error("Equip cannot be refined, try different equipment. Type 'i' to check equipment and its index.")
	}

	return nil
}

// HandleRefineResult handles the refine_result packet (lines 9161-9176)
func (rm *RefiningManager) HandleRefineResult(args map[string]interface{}) error {
	// Extract packet data
	fail, ok := args["fail"].(uint8)
	if !ok {
		return fmt.Errorf("invalid fail in refine_result packet")
	}

	nameID, ok := args["nameID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid nameID in refine_result packet")
	}

	// Process based on fail code
	switch fail {
	case 0:
		rm.logger.Info("You successfully refined a weapon (ID %d)!", nameID)
	case 1:
		rm.logger.Info("You failed to refine a weapon (ID %d)!", nameID)
	case 2:
		rm.logger.Info("You successfully made a potion (ID %d)!", nameID)
	case 3:
		rm.logger.Info("You failed to make a potion (ID %d)!", nameID)
	case 6:
		rm.logger.Info("You successfully cook a item (ID %d)!", nameID)
	default:
		rm.logger.Info("You tried to refine a weapon (ID %d); result: unknown %d", nameID, fail)
	}

	return nil
}

// HandleRefineStatus handles the refine_status packet (lines 7938-7943)
func (rm *RefiningManager) HandleRefineStatus(args map[string]interface{}) error {
	// Extract packet data
	status, ok := args["status"].(uint8)
	if !ok {
		return fmt.Errorf("invalid status in refine_status packet")
	}

	refineLevel, ok := args["refine_level"].(uint16)
	if !ok {
		return fmt.Errorf("invalid refine_level in refine_status packet")
	}

	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in refine_status packet")
	}

	// Get item name (simplified)
	itemName := fmt.Sprintf("Item-%d", refineLevel) // In a real implementation, this would use the actual item name

	// Process based on status
	if status == 0 {
		rm.logger.Warning("Weapon upgraded: %s +%d %s", name, refineLevel, itemName)
	} else {
		rm.logger.Warning("Weapon not upgraded: %s +%d %s", name, refineLevel, itemName)
	}

	return nil
}

// HandleUpgradeList handles the upgrade_list packet (lines 9114-9134)
func (rm *RefiningManager) HandleUpgradeList(args map[string]interface{}) error {
	// Extract packet data
	itemList, ok := args["item_list"].([]byte)
	if !ok {
		return fmt.Errorf("invalid item_list in upgrade_list packet")
	}

	// Clear upgrade list
	rm.upgradeList = make([]uint16, 0)

	// Process upgrade items
	msg := centerString(" Upgrade List ", 79, '-') + "\n"

	for i := 0; i < len(itemList); i += 13 {
		if i+13 <= len(itemList) {
			// Extract item data (simplified)
			// In a real implementation, this would use proper unpacking
			index := uint16(itemList[i]) | uint16(itemList[i+1])<<8
			nameID := uint8(itemList[i+8])

			// Add item to upgrade list
			rm.upgradeList = append(rm.upgradeList, index)

			// Format item info for display
			itemName := fmt.Sprintf("Item-%d", nameID) // In a real implementation, this would use a function to get the item name
			msg += fmt.Sprintf("%2d - %-50s (%3d)\n", len(rm.upgradeList)-1, itemName, index)
		}
	}

	msg += strings.Repeat("-", 79) + "\n"
	rm.logger.Info(msg)
	rm.logger.Info("You can now use the 'refine' command.")

	return nil
}

// HandleUpgradeMessage handles the upgrade_message packet (lines 9179-9192)
func (rm *RefiningManager) HandleUpgradeMessage(args map[string]interface{}) error {
	// Extract packet data
	upgradeType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in upgrade_message packet")
	}

	itemID, ok := args["itemID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid itemID in upgrade_message packet")
	}

	// Get item name (simplified)
	itemName := fmt.Sprintf("Item-%d", itemID) // In a real implementation, this would use a function to get the item name

	// Process based on upgrade type
	switch upgradeType {
	case 0: // Success
		rm.logger.Info("Weapon upgraded: %s", itemName)
	case 1: // Fail
		rm.logger.Info("Weapon not upgraded: %s", itemName)
	case 2: // Fail Lvl
		rm.logger.Error("Cannot upgrade %s until you level up the upgrade weapon skill.", itemName)
	case 3: // Fail Item
		rm.logger.Info("You lack item %s to upgrade the weapon.", itemName)
	default:
		rm.logger.Info("Unknown upgrade result for %s: %d", itemName, upgradeType)
	}

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
