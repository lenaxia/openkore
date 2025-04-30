package battle

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// BattleManager handles battle-related packet handling
type BattleManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger
}

// NewBattleManager creates a new battle manager
func NewBattleManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *BattleManager {
	return &BattleManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
	}
}

// RegisterHandlers registers all battle-related packet handlers
func (bm *BattleManager) RegisterHandlers() {
	// Register battleground message handler
	bm.parser.RegisterHandlerFunc("02DC", "battleground_message", "Z*",
		[]string{"message"}, bm.HandleBattlegroundMessage)

	// Register battleground emblem handler
	bm.parser.RegisterHandlerFunc("02DD", "battleground_emblem", "V",
		[]string{"emblemID"}, bm.HandleBattlegroundEmblem)

	// Register instance window start handler
	bm.parser.RegisterHandlerFunc("02CB", "instance_window_start", "V Z16 V V",
		[]string{"instanceID", "name", "time", "progress"}, bm.HandleInstanceWindowStart)

	// Register instance window queue handler
	bm.parser.RegisterHandlerFunc("02CC", "instance_window_queue", "V",
		[]string{"instanceID"}, bm.HandleInstanceWindowQueue)

	// Register instance window join handler
	bm.parser.RegisterHandlerFunc("02CD", "instance_window_join", "V V",
		[]string{"instanceID", "result"}, bm.HandleInstanceWindowJoin)

	// Register instance window leave handler
	bm.parser.RegisterHandlerFunc("02CE", "instance_window_leave", "V C",
		[]string{"instanceID", "flag"}, bm.HandleInstanceWindowLeave)
}

// HandleBattlegroundMessage handles the battleground_message packet (lines 10356-10359)
func (bm *BattleManager) HandleBattlegroundMessage(args map[string]interface{}) error {
	// Extract packet data
	message, ok := args["message"].(string)
	if !ok {
		return fmt.Errorf("invalid message in battleground_message packet")
	}

	// Log message
	bm.logger.Debug("Battleground message: %s", message)

	// Call hook
	bm.hookManager.CallHook("battleground_message", map[string]interface{}{
		"message": message,
	})

	return nil
}

// HandleBattlegroundEmblem handles the battleground_emblem packet (lines 10363-10366)
func (bm *BattleManager) HandleBattlegroundEmblem(args map[string]interface{}) error {
	// Extract packet data
	emblemID, ok := args["emblemID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid emblemID in battleground_emblem packet")
	}

	// Log emblem
	bm.logger.Debug("Battleground emblem: %d", emblemID)

	// Call hook
	bm.hookManager.CallHook("battleground_emblem", map[string]interface{}{
		"emblemID": emblemID,
	})

	return nil
}

// HandleInstanceWindowStart handles the instance_window_start packet (lines 10476-10479)
func (bm *BattleManager) HandleInstanceWindowStart(args map[string]interface{}) error {
	// Extract packet data
	instanceID, ok := args["instanceID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid instanceID in instance_window_start packet")
	}

	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in instance_window_start packet")
	}

	time, ok := args["time"].(uint32)
	if !ok {
		return fmt.Errorf("invalid time in instance_window_start packet")
	}

	progress, ok := args["progress"].(uint32)
	if !ok {
		return fmt.Errorf("invalid progress in instance_window_start packet")
	}

	// Log instance window start
	bm.logger.Debug("Instance window start: ID=%d, Name=%s, Time=%d, Progress=%d", instanceID, name, time, progress)

	// Call hook
	bm.hookManager.CallHook("instance_window_start", map[string]interface{}{
		"instanceID": instanceID,
		"name":       name,
		"time":       time,
		"progress":   progress,
	})

	return nil
}

// HandleInstanceWindowQueue handles the instance_window_queue packet (lines 10484-10487)
func (bm *BattleManager) HandleInstanceWindowQueue(args map[string]interface{}) error {
	// Extract packet data
	instanceID, ok := args["instanceID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid instanceID in instance_window_queue packet")
	}

	// Log instance window queue
	bm.logger.Debug("Instance window queue: ID=%d", instanceID)

	// Call hook
	bm.hookManager.CallHook("instance_window_queue", map[string]interface{}{
		"instanceID": instanceID,
	})

	return nil
}

// HandleInstanceWindowJoin handles the instance_window_join packet (lines 10491-10496)
func (bm *BattleManager) HandleInstanceWindowJoin(args map[string]interface{}) error {
	// Extract packet data
	instanceID, ok := args["instanceID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid instanceID in instance_window_join packet")
	}

	result, ok := args["result"].(uint32)
	if !ok {
		return fmt.Errorf("invalid result in instance_window_join packet")
	}

	// Log instance window join
	bm.logger.Debug("Instance window join: ID=%d, Result=%d", instanceID, result)

	// Call hook
	bm.hookManager.CallHook("instance_ready", map[string]interface{}{
		"instanceID": instanceID,
		"result":     result,
	})

	return nil
}

// HandleInstanceWindowLeave handles the instance_window_leave packet (lines 10507-10523)
func (bm *BattleManager) HandleInstanceWindowLeave(args map[string]interface{}) error {
	// Extract packet data
	instanceID, ok := args["instanceID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid instanceID in instance_window_leave packet")
	}

	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in instance_window_leave packet")
	}

	// Handle based on flag
	switch flag {
	case 0: // TYPE_NOTIFY = 0x0
		bm.logger.Debug("Received Memory Dungeon reservation update")
	case 1: // TYPE_DESTROY_LIVE_TIMEOUT = 0x1
		bm.logger.Info("The Memorial Dungeon expired it has been destroyed.")
	case 2: // TYPE_DESTROY_ENTER_TIMEOUT = 0x2
		bm.logger.Info("The Memorial Dungeon's entry time limit expired it has been destroyed.")
	case 3: // TYPE_DESTROY_USER_REQUEST = 0x3
		bm.logger.Info("The Memorial Dungeon has been removed.")
	case 4: // TYPE_CREATE_FAIL = 0x4
		bm.logger.Info("The instance windows has been removed, possibly due to party/guild leave.")
	default:
		bm.logger.Warning("Unknown results in instance_window_leave (flag: %d)", flag)
	}

	// Call hook
	bm.hookManager.CallHook("instance_window_leave", map[string]interface{}{
		"instanceID": instanceID,
		"flag":       flag,
	})

	return nil
}
