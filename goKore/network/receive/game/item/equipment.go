package item

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// EquipmentManager handles equipment-related packet handling
type EquipmentManager struct {
	baseParse Parser
	hooks     *hooks.HookManager
	logger    core.Logger
}

// NewEquipmentManager creates a new equipment manager
func NewEquipmentManager(baseParse Parser, hooks *hooks.HookManager, logger core.Logger) *EquipmentManager {
	return &EquipmentManager{
		baseParse: baseParse,
		hooks:     hooks,
		logger:    logger,
	}
}

// RegisterHandlers registers all equipment-related packet handlers
func (em *EquipmentManager) RegisterHandlers() {
	// Register equip item handler (0x0999)
	em.baseParse.RegisterHandler("0999", "equip_item", "v V C",
		[]string{"index", "type", "success"},
		em.HandleEquipItem)

	// Register unequip item handler (0x099A)
	em.baseParse.RegisterHandler("099A", "unequip_item", "v V C",
		[]string{"index", "type", "success"},
		em.HandleUnequipItem)

	// Register equip item switch handler (0x0A9B)
	em.baseParse.RegisterHandler("0A9B", "equip_item_switch", "v V C",
		[]string{"index", "position", "success"},
		em.HandleEquipItemSwitch)

	// Register unequip item switch handler (0x0A9C)
	em.baseParse.RegisterHandler("0A9C", "unequip_item_switch", "v V C",
		[]string{"ID", "type", "success"},
		em.HandleUnequipItemSwitch)

	// Register equip switch run result handler (0x0A9D)
	em.baseParse.RegisterHandler("0A9D", "equip_switch_run_res", "C",
		[]string{"success"},
		em.HandleEquipSwitchRunRes)

	// Register equip switch log handler (0x0A9E)
	em.baseParse.RegisterHandler("0A9E", "equip_switch_log", "v a*",
		[]string{"len", "items"},
		em.HandleEquipSwitchLog)

	// Register arrow equipped handler (0x0A3B)
	em.baseParse.RegisterHandler("0A3B", "arrow_equipped", "v",
		[]string{"index"},
		em.HandleArrowEquipped)

	// Register arrow none handler (0x013B)
	em.baseParse.RegisterHandler("013B", "arrow_none", "C",
		[]string{"type"},
		em.HandleArrowNone)
}

// HandleEquipItem handles the equip_item packet (lines 6045-6071)
func (em *EquipmentManager) HandleEquipItem(args map[string]interface{}) error {
	// Extract packet data
	index, ok := args["index"].(uint16)
	if !ok {
		return fmt.Errorf("invalid index in equip_item packet")
	}

	equipType, ok := args["type"].(uint16)
	if !ok {
		return fmt.Errorf("invalid type in equip_item packet")
	}

	success, ok := args["success"].(uint8)
	if !ok {
		return fmt.Errorf("invalid success in equip_item packet")
	}

	// Get equipment type name
	equipTypeName := getEquipTypeName(equipType)

	// Process based on success code
	if success == 1 {
		em.logger.Info("Item [%d] equipped as %s", index, equipTypeName)

		// Call hooks
		em.hooks.CallHook("item_equipped", map[string]interface{}{
			"index":     index,
			"type":      equipType,
			"type_name": equipTypeName,
		})
	} else {
		em.logger.Error("Failed to equip item [%d] as %s", index, equipTypeName)
	}

	return nil
}

// HandleUnequipItem handles the unequip_item packet (lines 11292-11318)
func (em *EquipmentManager) HandleUnequipItem(args map[string]interface{}) error {
	// Extract packet data
	index, ok := args["index"].(uint16)
	if !ok {
		return fmt.Errorf("invalid index in unequip_item packet")
	}

	equipType, ok := args["type"].(uint16)
	if !ok {
		return fmt.Errorf("invalid type in unequip_item packet")
	}

	success, ok := args["success"].(uint8)
	if !ok {
		return fmt.Errorf("invalid success in unequip_item packet")
	}

	// Get equipment type name
	equipTypeName := getEquipTypeName(equipType)

	// Process based on success code
	if success == 1 {
		em.logger.Info("Item [%d] unequipped from %s", index, equipTypeName)

		// Call hooks
		em.hooks.CallHook("item_unequipped", map[string]interface{}{
			"index":     index,
			"type":      equipType,
			"type_name": equipTypeName,
		})
	} else {
		em.logger.Error("Failed to unequip item [%d] from %s", index, equipTypeName)
	}

	return nil
}

// HandleEquipItemSwitch handles the equip_item_switch packet (lines 6072-6098)
func (em *EquipmentManager) HandleEquipItemSwitch(args map[string]interface{}) error {
	// Extract packet data
	index, ok := args["index"].(uint16)
	if !ok {
		return fmt.Errorf("invalid index in equip_item_switch packet")
	}

	position, ok := args["position"].(uint32)
	if !ok {
		return fmt.Errorf("invalid position in equip_item_switch packet")
	}

	success, ok := args["success"].(uint8)
	if !ok {
		return fmt.Errorf("invalid success in equip_item_switch packet")
	}

	// Process based on success code
	if success == 1 {
		em.logger.Info("Item [%d] equipped in equipment switch at position %d", index, position)

		// Call hooks
		em.hooks.CallHook("item_equipped_switch", map[string]interface{}{
			"index":    index,
			"position": position,
		})
	} else {
		em.logger.Error("Failed to equip item [%d] in equipment switch at position %d", index, position)
	}

	return nil
}

// HandleEquipSwitchRunRes handles the equip_switch_run_res packet (lines 6099-6115)
func (em *EquipmentManager) HandleEquipSwitchRunRes(args map[string]interface{}) error {
	// Extract packet data
	success, ok := args["success"].(uint8)
	if !ok {
		return fmt.Errorf("invalid success in equip_switch_run_res packet")
	}

	// Process based on success code
	if success == 1 {
		em.logger.Success("Equipment set switched successfully")

		// Call hooks
		em.hooks.CallHook("equipment_set_switched", nil)
	} else {
		em.logger.Error("Failed to switch equipment set")
	}

	return nil
}

// HandleEquipSwitchLog handles the equip_switch_log packet (lines 6116-6138)
func (em *EquipmentManager) HandleEquipSwitchLog(args map[string]interface{}) error {
	// Extract packet data
	items, ok := args["items"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid items in equip_switch_log packet")
	}

	// Log the equipment switch items
	em.logger.Info("Equipment switch log:")

	// Process each item
	for _, item := range items {
		nameID, _ := item["nameID"].(uint16)
		itemType, _ := item["type"].(uint8)

		em.logger.Info("- Item [%d] (Type: %d)", nameID, itemType)
	}

	// Call hooks
	em.hooks.CallHook("equipment_switch_log", map[string]interface{}{
		"items": items,
	})

	return nil
}

// HandleArrowEquipped handles the arrow_equipped packet (lines 3520-3536)
func (em *EquipmentManager) HandleArrowEquipped(args map[string]interface{}) error {
	// Extract packet data
	index, ok := args["index"].(uint16)
	if !ok {
		return fmt.Errorf("invalid index in arrow_equipped packet")
	}

	// Log the arrow equipment
	em.logger.Info("Arrow equipped: [%d]", index)

	// Call hooks
	em.hooks.CallHook("arrow_equipped", map[string]interface{}{
		"index": index,
	})

	return nil
}

// HandleUnequipItemSwitch handles the unequip_item_switch packet (lines 11643-11670)
func (em *EquipmentManager) HandleUnequipItemSwitch(args map[string]interface{}) error {
	// Extract packet data
	ID, ok := args["ID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid ID in unequip_item_switch packet")
	}

	equipType, ok := args["type"].(uint16)
	if !ok {
		return fmt.Errorf("invalid type in unequip_item_switch packet")
	}

	success, ok := args["success"].(uint8)
	if !ok {
		return fmt.Errorf("invalid success in unequip_item_switch packet")
	}

	// Process based on success code
	if success == 1 {
		// Handle special case for arrows
		if equipType == 10 || equipType == 32768 {
			em.logger.Info("Arrow unequipped from equipment switch: [%d]", ID)
		} else {
			// Get equipment type name
			equipTypeName := getEquipTypeName(equipType)
			em.logger.Info("Item [%d] unequipped from equipment switch: %s", ID, equipTypeName)
		}

		// Call hooks
		em.hooks.CallHook("item_unequipped_switch", map[string]interface{}{
			"ID":   ID,
			"type": equipType,
		})
	} else {
		em.logger.Error("Failed to unequip item [%d] from equipment switch", ID)
	}

	return nil
}

// HandleArrowNone handles the arrow_none packet (lines 10161-10182)
func (em *EquipmentManager) HandleArrowNone(args map[string]interface{}) error {
	// Extract packet data
	arrowType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in arrow_none packet")
	}

	// Process based on arrow type
	switch arrowType {
	case 0:
		// No arrows
		em.logger.Error("Please equip arrow first")

		// In the original code, there's logic to disconnect if dcOnEmptyArrow is set
		// We'll just log this for now, as disconnection logic would be handled elsewhere
	case 1:
		// Weight limit exceeded - can't attack or use skills
		em.logger.Debug("You can't Attack or use Skills because your Weight Limit has been exceeded")
	case 2:
		// Weight limit exceeded - can't use skills
		em.logger.Debug("You can't use Skills because Weight Limit has been exceeded")
	case 3:
		// Arrow equipped
		em.logger.Debug("Arrow equipped")
	default:
		em.logger.Debug("Unknown arrow status: %d", arrowType)
	}

	// Call hooks
	em.hooks.CallHook("arrow_none", map[string]interface{}{
		"type": arrowType,
	})

	return nil
}

// Helper function to get equipment type name
func getEquipTypeName(equipType uint16) string {
	switch equipType {
	case 0:
		return "Lowhead"
	case 1:
		return "Righthand"
	case 2:
		return "Lefthand"
	case 3:
		return "Armor"
	case 4:
		return "Midhead"
	case 5:
		return "Tophead"
	case 6:
		return "Accessory1"
	case 7:
		return "Accessory2"
	case 8:
		return "Shoes"
	case 9:
		return "Garment"
	case 10:
		return "Arrow"
	default:
		return fmt.Sprintf("Unknown(%d)", equipType)
	}
}
