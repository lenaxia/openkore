package mercenary

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// MercInfo represents mercenary information
type MercInfo struct {
	ID           uint32
	Name         string
	Level        uint16
	HP           uint32
	MaxHP        uint32
	SP           uint32
	MaxSP        uint32
	Attack       uint16
	MagicAttack  uint16
	Hit          uint16
	Critical     uint16
	Defense      uint16
	MagicDefense uint16
	Flee         uint16
	AspeedDelay  uint16
	AttackRange  uint16
	ExpireTime   uint32
	Faith        uint16
	Calls        uint32
	Kills        uint32
	Map          string
}

// MercManager handles mercenary-related packet handling
type MercManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for mercenary interactions
	mercInfo *MercInfo
	charID   uint32
}

// NewMercManager creates a new mercenary manager
func NewMercManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *MercManager {
	return &MercManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
		mercInfo:    nil,
	}
}

// RegisterHandlers registers all mercenary-related packet handlers
func (mm *MercManager) RegisterHandlers() {
	// Register mercenary init handler
	mm.parser.RegisterHandlerFunc("029B", "mercenary_init", "L W W W W W W W W Z24 W L L L L L W L L W",
		[]string{"ID", "atk", "matk", "hit", "crit", "def", "mdef", "flee", "aspd", "name", "level", "hp", "maxhp", "sp", "maxsp", "expire_time", "faith", "calls", "kills", "atk_range"}, mm.HandleMercenaryInit)

	// Register mercenary off handler
	mm.parser.RegisterHandlerFunc("02A5", "mercenary_off", "",
		[]string{}, mm.HandleMercenaryOff)
}

// SetCharID sets the character's ID
func (mm *MercManager) SetCharID(id uint32) {
	mm.charID = id
}

// SlaveCalcPropertyHandler calculates slave properties
func SlaveCalcPropertyHandler(slave map[string]interface{}, args map[string]interface{}) {
	// Calculate attack speed
	aspd, ok := args["aspd"].(uint16)
	if !ok {
		return
	}

	// Attack speed calculation
	// attack_speed = int(200 - ((aspd < 10) ? 10 : (aspd / 10)))
	var attackSpeed int
	if aspd < 10 {
		attackSpeed = 200 - 10
	} else {
		attackSpeed = 200 - int(aspd/10)
	}

	slave["attack_speed"] = attackSpeed
}

// HandleMercenaryInit handles the mercenary_init packet (lines 11093-11128)
func (mm *MercManager) HandleMercenaryInit(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in mercenary_init packet")
	}

	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in mercenary_init packet")
	}

	// Create mercenary info if it doesn't exist or ID is different
	if mm.mercInfo == nil || mm.mercInfo.ID != id {
		mm.mercInfo = &MercInfo{
			ID:  id,
			Map: "unknown", // This would be set from the field name in the original code
		}
	}

	// Update mercenary info
	mm.mercInfo.Name = name

	// Extract and update other properties
	if level, ok := args["level"].(uint16); ok {
		mm.mercInfo.Level = level
	}
	if hp, ok := args["hp"].(uint32); ok {
		mm.mercInfo.HP = hp
	}
	if maxhp, ok := args["maxhp"].(uint32); ok {
		mm.mercInfo.MaxHP = maxhp
	}
	if sp, ok := args["sp"].(uint32); ok {
		mm.mercInfo.SP = sp
	}
	if maxsp, ok := args["maxsp"].(uint32); ok {
		mm.mercInfo.MaxSP = maxsp
	}
	if atk, ok := args["atk"].(uint16); ok {
		mm.mercInfo.Attack = atk
	}
	if matk, ok := args["matk"].(uint16); ok {
		mm.mercInfo.MagicAttack = matk
	}
	if hit, ok := args["hit"].(uint16); ok {
		mm.mercInfo.Hit = hit
	}
	if crit, ok := args["crit"].(uint16); ok {
		mm.mercInfo.Critical = crit
	}
	if def, ok := args["def"].(uint16); ok {
		mm.mercInfo.Defense = def
	}
	if mdef, ok := args["mdef"].(uint16); ok {
		mm.mercInfo.MagicDefense = mdef
	}
	if flee, ok := args["flee"].(uint16); ok {
		mm.mercInfo.Flee = flee
	}
	if aspd, ok := args["aspd"].(uint16); ok {
		mm.mercInfo.AspeedDelay = aspd
	}
	if atkRange, ok := args["atk_range"].(uint16); ok {
		mm.mercInfo.AttackRange = atkRange
	}
	if expireTime, ok := args["expire_time"].(uint32); ok {
		mm.mercInfo.ExpireTime = expireTime
	}
	if faith, ok := args["faith"].(uint16); ok {
		mm.mercInfo.Faith = faith
	}
	if calls, ok := args["calls"].(uint32); ok {
		mm.mercInfo.Calls = calls
	}
	if kills, ok := args["kills"].(uint32); ok {
		mm.mercInfo.Kills = kills
	}

	// Calculate slave properties
	slaveProps := make(map[string]interface{})
	SlaveCalcPropertyHandler(slaveProps, args)

	// Log mercenary info
	mm.logger.Info("Mercenary initialized: %s (ID: %d, Level: %d)", mm.mercInfo.Name, mm.mercInfo.ID, mm.mercInfo.Level)
	mm.logger.Info("HP: %d/%d, SP: %d/%d", mm.mercInfo.HP, mm.mercInfo.MaxHP, mm.mercInfo.SP, mm.mercInfo.MaxSP)
	mm.logger.Info("ATK: %d, MATK: %d, DEF: %d, MDEF: %d", mm.mercInfo.Attack, mm.mercInfo.MagicAttack, mm.mercInfo.Defense, mm.mercInfo.MagicDefense)
	mm.logger.Info("HIT: %d, CRIT: %d, FLEE: %d, ASPD: %d", mm.mercInfo.Hit, mm.mercInfo.Critical, mm.mercInfo.Flee, mm.mercInfo.AspeedDelay)
	mm.logger.Info("Attack Range: %d, Expire Time: %d", mm.mercInfo.AttackRange, mm.mercInfo.ExpireTime)
	mm.logger.Info("Faith: %d, Calls: %d, Kills: %d", mm.mercInfo.Faith, mm.mercInfo.Calls, mm.mercInfo.Kills)

	// Call hook
	mm.hookManager.CallHook("mercenary_init", map[string]interface{}{
		"mercInfo": mm.mercInfo,
	})

	return nil
}

// HandleMercenaryOff handles the mercenary_off packet (lines 11131-11137)
func (mm *MercManager) HandleMercenaryOff(args map[string]interface{}) error {
	// Check if mercenary exists
	if mm.mercInfo == nil {
		mm.logger.Warning("Received mercenary_off packet but no mercenary exists")
		return nil
	}

	// Log mercenary removal
	mm.logger.Info("Mercenary removed: %s (ID: %d)", mm.mercInfo.Name, mm.mercInfo.ID)

	// Call hook
	mm.hookManager.CallHook("mercenary_off", map[string]interface{}{
		"ID": mm.mercInfo.ID,
	})

	// Remove mercenary
	mm.mercInfo = nil

	return nil
}
