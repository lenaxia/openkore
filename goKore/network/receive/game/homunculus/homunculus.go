package homunculus

import (
	"fmt"
	"time"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// Homunculus state constants
const (
	HO_PRE_INIT             = 0
	HO_RELATIONSHIP_CHANGED = 1
	HO_FULLNESS_CHANGED     = 2
	HO_ACCESSORY_CHANGED    = 3
	HO_HEADTYPE_CHANGED     = 4
)

// HomState represents the homunculus state bits
type HomState struct {
	Named     bool // bit 0
	Vaporized bool // bit 1
	Dead      bool // bit 2
}

// HomInfo represents homunculus information
type HomInfo struct {
	ID            uint32
	Name          string
	Level         uint16
	Hunger        uint16
	Intimacy      uint16
	EquipID       uint16
	Attack        uint16
	MagicAttack   uint16
	Hit           uint16
	Critical      uint16
	Defense       uint16
	MagicDefense  uint16
	Flee          uint16
	AspeedDelay   uint16
	HP            uint32
	MaxHP         uint32
	SP            uint16
	MaxSP         uint16
	Exp           uint32
	MaxExp        uint32
	SkillPoints   uint16
	AttackRange   uint16
	State         HomState
	RenameFlagSet bool
	Vaporized     bool
	Dead          bool
	VaporizeTime  int64
}

// HomManager handles homunculus-related packet handling
type HomManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for homunculus interactions
	homInfo *HomInfo
	charID  uint32
}

// NewHomManager creates a new homunculus manager
func NewHomManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *HomManager {
	return &HomManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
		homInfo:     nil,
	}
}

// RegisterHandlers registers all homunculus-related packet handlers
func (hm *HomManager) RegisterHandlers() {
	// Register homunculus property handler
	hm.parser.RegisterHandlerFunc("022E", "homunculus_property", "Z24 B W W W W W W W W W W W W W W W W L L W W",
		[]string{"name", "modified", "level", "hunger", "intimacy", "equip_id", "atk", "matk", "hit", "crit", "def", "mdef", "flee", "aspd", "hp", "max_hp", "sp", "max_sp", "exp", "max_exp", "skill_points", "atk_range"}, hm.HandleHomProperty)

	// Register homunculus state handler
	hm.parser.RegisterHandlerFunc("0230", "homunculus_info", "B B L L",
		[]string{"type", "state", "ID", "val"}, hm.HandleHomInfo)

	// Register homunculus food handler
	hm.parser.RegisterHandlerFunc("022F", "homunculus_food", "B W",
		[]string{"success", "foodID"}, hm.HandleHomFood)

	// Register egg list handler
	hm.parser.RegisterHandlerFunc("01A6", "egg_list", "v a*",
		[]string{"len", "RAW_MSG"}, hm.HandleEggList)
}

// SetCharID sets the character's ID
func (hm *HomManager) SetCharID(id uint32) {
	hm.charID = id
}

// EnforceHomState ensures homunculus state exists
func (hm *HomManager) EnforceHomState() bool {
	if hm.homInfo == nil {
		hm.logger.Debug("[Homunculus] Received homunculus property without the homunculus object existing, creating a temporary one.")
		hm.homInfo = &HomInfo{}
		return true
	}
	return true
}

// HandleHomProperty handles the homunculus_property packet (lines 2897-2925)
func (hm *HomManager) HandleHomProperty(args map[string]interface{}) error {
	// Check if homunculus state exists
	if !hm.EnforceHomState() {
		return fmt.Errorf("failed to enforce homunculus state")
	}

	// Extract packet data
	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in homunculus_property packet")
	}

	// Update homunculus info
	hm.homInfo.Name = name

	// Extract and update other properties
	if level, ok := args["level"].(uint16); ok {
		hm.homInfo.Level = level
	}
	if hunger, ok := args["hunger"].(uint16); ok {
		hm.homInfo.Hunger = hunger
	}
	if intimacy, ok := args["intimacy"].(uint16); ok {
		hm.homInfo.Intimacy = intimacy
	}
	if equipID, ok := args["equip_id"].(uint16); ok {
		hm.homInfo.EquipID = equipID
	}
	if atk, ok := args["atk"].(uint16); ok {
		hm.homInfo.Attack = atk
	}
	if matk, ok := args["matk"].(uint16); ok {
		hm.homInfo.MagicAttack = matk
	}
	if hit, ok := args["hit"].(uint16); ok {
		hm.homInfo.Hit = hit
	}
	if crit, ok := args["crit"].(uint16); ok {
		hm.homInfo.Critical = crit
	}
	if def, ok := args["def"].(uint16); ok {
		hm.homInfo.Defense = def
	}
	if mdef, ok := args["mdef"].(uint16); ok {
		hm.homInfo.MagicDefense = mdef
	}
	if flee, ok := args["flee"].(uint16); ok {
		hm.homInfo.Flee = flee
	}
	if aspd, ok := args["aspd"].(uint16); ok {
		hm.homInfo.AspeedDelay = aspd
	}
	if hp, ok := args["hp"].(uint16); ok {
		hm.homInfo.HP = uint32(hp)
	}
	if maxHP, ok := args["max_hp"].(uint16); ok {
		hm.homInfo.MaxHP = uint32(maxHP)
	}
	if sp, ok := args["sp"].(uint16); ok {
		hm.homInfo.SP = sp
	}
	if maxSP, ok := args["max_sp"].(uint16); ok {
		hm.homInfo.MaxSP = maxSP
	}
	if exp, ok := args["exp"].(uint32); ok {
		hm.homInfo.Exp = exp
	}
	if maxExp, ok := args["max_exp"].(uint32); ok {
		hm.homInfo.MaxExp = maxExp
	}
	if skillPoints, ok := args["skill_points"].(uint16); ok {
		hm.homInfo.SkillPoints = skillPoints
	}
	if atkRange, ok := args["atk_range"].(uint16); ok {
		hm.homInfo.AttackRange = atkRange
	}

	// Process homunculus state
	if modified, ok := args["modified"].(uint8); ok {
		hm.ProcessHomState(uint8(modified))
	}

	// Call hook
	hm.hookManager.CallHook("homunculus_property", map[string]interface{}{
		"homInfo": hm.homInfo,
	})

	return nil
}

// ProcessHomState processes the homunculus state bits
func (hm *HomManager) ProcessHomState(state uint8) {
	// Process state bits
	// Bit 0: Named
	// Bit 1: Vaporized
	// Bit 2: Dead

	// Create new state
	newState := HomState{
		Named:     (state & 1) != 0,
		Vaporized: (state & 2) != 0,
		Dead:      (state & 4) == 0, // Note: The bit is inverted in the original code
	}

	// If this is the first time we're setting the state
	if hm.homInfo.State.Named != newState.Named {
		if newState.Named {
			hm.logger.Info("Your Homunculus has already been renamed")
			hm.homInfo.RenameFlagSet = true
		} else {
			hm.logger.Info("Your Homunculus has not been renamed")
			hm.homInfo.RenameFlagSet = false
		}
	}

	if hm.homInfo.State.Vaporized != newState.Vaporized {
		if newState.Vaporized {
			hm.logger.Info("Your Homunculus is vaporized")
			hm.homInfo.Vaporized = true
		} else {
			hm.logger.Info("Your Homunculus is not vaporized")
			hm.homInfo.Vaporized = false
		}
	}

	if hm.homInfo.State.Dead != newState.Dead {
		if newState.Dead {
			hm.logger.Info("Your Homunculus is dead")
			hm.homInfo.Dead = true
		} else {
			hm.logger.Info("Your Homunculus is not dead")
			hm.homInfo.Dead = false
		}
	}

	// Update state
	hm.homInfo.State = newState
}

// HandleHomInfo handles the homunculus_info packet (lines 3022-3055)
func (hm *HomManager) HandleHomInfo(args map[string]interface{}) error {
	// Extract packet data
	state, ok := args["state"].(uint8)
	if !ok {
		return fmt.Errorf("invalid state in homunculus_info packet")
	}

	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in homunculus_info packet")
	}

	val, ok := args["val"].(uint32)
	if !ok {
		return fmt.Errorf("invalid val in homunculus_info packet")
	}

	// Debug log
	hm.logger.Debug("homunculus_info type: %d", args["type"])

	// Handle based on state
	switch state {
	case HO_PRE_INIT:
		// Some servers won't send 'homunculus_property' after a teleport, so we don't delete homInfo object
		if hm.homInfo != nil && hm.homInfo.Dead {
			hm.logger.Debug("[Homunculus] We received a homunculus_info packet while our homunculus is dead, assume it was resurrected.")
			hm.homInfo.Dead = false
		}

		// Update homunculus ID if needed
		if hm.homInfo == nil || hm.homInfo.ID != id {
			if hm.homInfo == nil {
				hm.homInfo = &HomInfo{}
			}
			hm.homInfo.ID = id
		}

		// Call hook
		hm.hookManager.CallHook("homunculus_appear", map[string]interface{}{
			"ID": id,
		})

	case HO_RELATIONSHIP_CHANGED:
		if hm.homInfo != nil {
			hm.homInfo.Intimacy = uint16(val)
		}

	case HO_FULLNESS_CHANGED:
		if hm.homInfo != nil {
			hm.homInfo.Hunger = uint16(val)
		}

	case HO_ACCESSORY_CHANGED:
		if hm.homInfo != nil {
			hm.homInfo.EquipID = uint16(val)
		}

	case HO_HEADTYPE_CHANGED:
		// Not implemented in original code
	}

	// Call hook
	hm.hookManager.CallHook("homunculus_info", map[string]interface{}{
		"state": state,
		"ID":    id,
		"val":   val,
	})

	return nil
}

// HandleHomFood handles the homunculus_food packet (lines 6265-6278)
func (hm *HomManager) HandleHomFood(args map[string]interface{}) error {
	// Extract packet data
	success, ok := args["success"].(uint8)
	if !ok {
		return fmt.Errorf("invalid success in homunculus_food packet")
	}

	foodID, ok := args["foodID"].(uint16)
	if !ok {
		return fmt.Errorf("invalid foodID in homunculus_food packet")
	}

	// Handle based on success
	if success != 0 {
		hm.logger.Info("Fed homunculus with item ID %d", foodID)
	} else {
		hm.logger.Error("Failed to feed homunculus with item ID %d: no food in inventory.", foodID)

		// Auto-vaporize if hunger is critical
		if hm.homInfo != nil && hm.homInfo.Hunger <= 11 && (hm.homInfo.VaporizeTime == 0 || (currentTime()-hm.homInfo.VaporizeTime) > 5) {
			// In the original code, this would send a skill use packet
			// $messageSender->sendSkillUse(244, 1, $accountID);
			hm.homInfo.VaporizeTime = currentTime()
			hm.logger.Error("Critical hunger level reached. Homunculus is put to rest.")
		}
	}

	// Call hook
	hm.hookManager.CallHook("homunculus_food", map[string]interface{}{
		"success": success != 0,
		"foodID":  foodID,
	})

	return nil
}

// HandleEggList handles the egg_list packet (lines 5948-5959)
func (hm *HomManager) HandleEggList(args map[string]interface{}) error {
	// Extract packet data
	rawMsg, ok := args["RAW_MSG"].([]byte)
	if !ok {
		return fmt.Errorf("invalid RAW_MSG in egg_list packet")
	}

	// Parse egg list
	var eggList []map[string]interface{}
	for i := 4; i < len(rawMsg); i += 2 {
		if i+2 > len(rawMsg) {
			break // Prevent out of bounds access
		}

		// Extract item index
		index := uint16(rawMsg[i]) | uint16(rawMsg[i+1])<<8

		// Add to egg list
		eggList = append(eggList, map[string]interface{}{
			"index": index,
		})
	}

	// Log egg list
	hm.logger.Info("===== Egg Hatch Candidates =====")
	for _, egg := range eggList {
		hm.logger.Info("Egg Index: %d", egg["index"])
	}
	hm.logger.Info("Ready to use command 'pet [hatch|h] #'")

	// Call hook
	hm.hookManager.CallHook("egg_list", map[string]interface{}{
		"eggs": eggList,
	})

	return nil
}

// Variable for current time function to allow mocking in tests
var currentTime = func() int64 {
	return time.Now().Unix()
}
