package quest

import (
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// QuestManager handles quest-related packet handling
type QuestManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for quest interactions
	questList           map[uint32]map[string]interface{}
	questGeneration     int
	lastQuestGeneration int
}

// NewQuestManager creates a new quest manager
func NewQuestManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *QuestManager {
	return &QuestManager{
		parser:              parser,
		hookManager:         hookManager,
		logger:              logger,
		questList:           make(map[uint32]map[string]interface{}),
		questGeneration:     0,
		lastQuestGeneration: 0,
	}
}

// RegisterHandlers registers all quest-related packet handlers
func (qm *QuestManager) RegisterHandlers() {
	// Register quest_all_list handler
	qm.parser.RegisterHandlerFunc("02B1", "quest_all_list", "v V a*",
		[]string{"quest_amount", "len", "message"}, qm.HandleQuestAllList)
	qm.parser.RegisterHandlerFunc("097A", "quest_all_list", "v V a*",
		[]string{"quest_amount", "len", "message"}, qm.HandleQuestAllList)
	qm.parser.RegisterHandlerFunc("09F8", "quest_all_list", "v V a*",
		[]string{"quest_amount", "len", "message"}, qm.HandleQuestAllList)
	qm.parser.RegisterHandlerFunc("0AFF", "quest_all_list", "v V a*",
		[]string{"quest_amount", "len", "message"}, qm.HandleQuestAllList)

	// Register quest_all_mission handler
	qm.parser.RegisterHandlerFunc("02B2", "quest_all_mission", "v V a*",
		[]string{"mission_amount", "len", "message"}, qm.HandleQuestAllMission)

	// Register quest_add handler
	qm.parser.RegisterHandlerFunc("02B3", "quest_add", "V C V2 v a*",
		[]string{"questID", "active", "time_start", "time_expire", "mission_amount", "message"}, qm.HandleQuestAdd)
	qm.parser.RegisterHandlerFunc("09F9", "quest_add", "V C V2 v a*",
		[]string{"questID", "active", "time_start", "time_expire", "mission_amount", "message"}, qm.HandleQuestAdd)
	qm.parser.RegisterHandlerFunc("0B0C", "quest_add", "V C V2 v a*",
		[]string{"questID", "active", "time_start", "time_expire", "mission_amount", "message"}, qm.HandleQuestAdd)

	// Register quest_update_mission_hunt handler
	qm.parser.RegisterHandlerFunc("02B5", "quest_update_mission_hunt", "v V a*",
		[]string{"mission_amount", "len", "message"}, qm.HandleQuestUpdateMissionHunt)
	qm.parser.RegisterHandlerFunc("09FA", "quest_update_mission_hunt", "v V a*",
		[]string{"mission_amount", "len", "message"}, qm.HandleQuestUpdateMissionHunt)
	qm.parser.RegisterHandlerFunc("0AFE", "quest_update_mission_hunt", "v V a*",
		[]string{"mission_amount", "len", "message"}, qm.HandleQuestUpdateMissionHunt)
	qm.parser.RegisterHandlerFunc("08FE", "quest_update_mission_hunt", "V a*",
		[]string{"len", "message"}, qm.HandleQuestUpdateMissionHunt)

	// Register quest_delete handler
	qm.parser.RegisterHandlerFunc("02B4", "quest_delete", "V",
		[]string{"questID"}, qm.HandleQuestDelete)

	// Register quest_active handler
	qm.parser.RegisterHandlerFunc("02B7", "quest_active", "V C",
		[]string{"questID", "active"}, qm.HandleQuestActive)
}

// HandleQuestAllList handles the quest_all_list packet (lines 4599-4694)
func (qm *QuestManager) HandleQuestAllList(args map[string]interface{}) error {
	// Extract packet data
	questAmount, ok := args["quest_amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid quest_amount in quest_all_list packet")
	}

	message, ok := args["message"].([]byte)
	if !ok {
		return fmt.Errorf("invalid message in quest_all_list packet")
	}

	// Determine packet format based on switch
	var questInfo map[string]interface{}
	switch args["switch"].(string) {
	case "02B1": // DEFAULT PACKET
		questInfo = map[string]interface{}{
			"quest_pack":   "V C",
			"quest_keys":   []string{"quest_id", "active"},
			"quest_len":    5,
			"mission_pack": "",
			"mission_keys": []string{},
			"mission_len":  0,
		}
	case "097A": // SERVERTYPE >= 20141022
		questInfo = map[string]interface{}{
			"quest_pack":   "V C V2 v",
			"quest_keys":   []string{"quest_id", "active", "time_expire", "time_start", "mission_amount"},
			"quest_len":    15,
			"mission_pack": "V v2 Z24",
			"mission_keys": []string{"mob_id", "mob_count", "mob_goal", "mob_name_original"},
			"mission_len":  32,
		}
	case "09F8": // SERVERTYPE >= 20150513
		questInfo = map[string]interface{}{
			"quest_pack":   "V C V2 v",
			"quest_keys":   []string{"quest_id", "active", "time_expire", "time_start", "mission_amount"},
			"quest_len":    15,
			"mission_pack": "V3 v4 Z24",
			"mission_keys": []string{"hunt_id", "mob_type", "mob_id", "min_level", "max_level", "mob_count", "mob_goal", "mob_name_original"},
			"mission_len":  44,
		}
	case "0AFF": // SERVERTYPE >= 20181010
		questInfo = map[string]interface{}{
			"quest_pack":   "V C V2 v",
			"quest_keys":   []string{"quest_id", "active", "time_expire", "time_start", "mission_amount"},
			"quest_len":    15,
			"mission_pack": "V4 v4 Z24",
			"mission_keys": []string{"hunt_id", "hunt_id_cont", "mob_type", "mob_id", "min_level", "max_level", "mob_count", "mob_goal", "mob_name_original"},
			"mission_len":  48,
		}
	default:
		return fmt.Errorf("unknown packet switch: %s", args["switch"].(string))
	}

	// Long quest lists are split up over multiple packets. Only reset the quest list if we've switched maps.
	if qm.lastQuestGeneration != qm.questGeneration {
		qm.lastQuestGeneration = qm.questGeneration
		qm.questList = make(map[uint32]map[string]interface{})
	}

	offset := 0
	questLen := questInfo["quest_len"].(int)
	missionLen := questInfo["mission_len"].(int)

	for i := 0; i < int(questAmount); i++ {
		// Extract quest data
		quest := make(map[string]interface{})

		// In a real implementation, this would use proper unpacking based on questInfo["quest_pack"]
		// For simplicity, we'll just extract the quest_id and active status
		if offset+4 <= len(message) {
			questID := uint32(message[offset]) | uint32(message[offset+1])<<8 | uint32(message[offset+2])<<16 | uint32(message[offset+3])<<24
			quest["quest_id"] = questID

			if offset+5 <= len(message) {
				quest["active"] = uint8(message[offset+4])
			}

			// Add quest to questList
			qm.questList[questID] = quest

			qm.logger.Debug("Quest ID: %d - active: %d", questID, quest["active"])
		}

		offset += questLen

		// Skip if no mission data
		if missionLen == 0 || !strings.Contains(args["switch"].(string), "09") {
			continue
		}

		// Extract mission amount if available
		missionAmount := uint16(0)
		if questLen >= 15 && offset-questLen+14 < len(message) {
			missionAmount = uint16(message[offset-questLen+14])
		}

		quest["mission_amount"] = missionAmount
		quest["missions"] = make(map[uint32]map[string]interface{})

		qm.logger.Debug("- Mission amount: %d", missionAmount)

		// Process missions
		for j := 0; j < int(missionAmount); j++ {
			if offset+missionLen <= len(message) {
				mission := make(map[string]interface{})

				// In a real implementation, this would use proper unpacking based on questInfo["mission_pack"]
				// For simplicity, we'll just extract the mob_id and mob_count
				mobID := uint32(message[offset]) | uint32(message[offset+1])<<8 | uint32(message[offset+2])<<16 | uint32(message[offset+3])<<24
				mission["mob_id"] = mobID

				if offset+6 <= len(message) {
					mission["mob_count"] = uint16(message[offset+4]) | uint16(message[offset+5])<<8
				}

				// Extract mob name (simplified)
				mobName := "Unknown"
				mission["mob_name"] = mobName
				mission["mission_index"] = j

				// Add mission to quest
				if questID, ok := quest["quest_id"].(uint32); ok {
					qm.questList[questID]["missions"].(map[uint32]map[string]interface{})[mobID] = mission
				}

				qm.logger.Debug("- MobID: %d - Name: %s - Count: %d", mobID, mobName, mission["mob_count"])

				// Call hook for mission added
				qm.hookManager.CallHook("quest_mission_added", map[string]interface{}{
					"questID":    quest["quest_id"],
					"mission_id": mobID,
				})
			}

			offset += missionLen
		}
	}

	// Call hook for quest list end
	qm.hookManager.CallHook("quest_all_list_end", nil)

	return nil
}

// HandleQuestAllMission handles the quest_all_mission packet (lines 4698-4754)
func (qm *QuestManager) HandleQuestAllMission(args map[string]interface{}) error {
	// Extract packet data
	missionAmount, ok := args["mission_amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid mission_amount in quest_all_mission packet")
	}

	message, ok := args["message"].([]byte)
	if !ok {
		return fmt.Errorf("invalid message in quest_all_mission packet")
	}

	// Define quest info
	questInfo := map[string]interface{}{
		"quest_pack":   "V3 v",
		"quest_keys":   []string{"quest_id", "time_start", "time_expire", "mission_amount"},
		"quest_len":    14,
		"mission_pack": "V v Z24",
		"mission_keys": []string{"mob_id", "mob_count", "mob_name_original"},
		"mission_len":  30,
	}

	offset := 0
	questLen := questInfo["quest_len"].(int)
	missionLen := questInfo["mission_len"].(int)

	for i := 0; i < int(missionAmount); i++ {
		// Extract quest data
		quest := make(map[string]interface{})

		// In a real implementation, this would use proper unpacking based on questInfo["quest_pack"]
		// For simplicity, we'll just extract the quest_id
		if offset+4 <= len(message) {
			questID := uint32(message[offset]) | uint32(message[offset+1])<<8 | uint32(message[offset+2])<<16 | uint32(message[offset+3])<<24
			quest["quest_id"] = questID

			// Get or create quest in questList
			if _, exists := qm.questList[questID]; !exists {
				qm.questList[questID] = make(map[string]interface{})
			}

			// Copy quest data to questList
			for key, value := range quest {
				qm.questList[questID][key] = value
			}

			qm.logger.Debug("Quest ID: %d - active: %v", questID, qm.questList[questID]["active"])
		}

		offset += questLen

		// Process missions (always 3 in this packet)
		for j := 0; j < 3; j++ {
			// Skip if beyond mission amount
			if questID, ok := quest["quest_id"].(uint32); ok {
				if charQuest, exists := qm.questList[questID]; exists {
					if missionAmount, ok := charQuest["mission_amount"].(uint16); ok && j >= int(missionAmount) {
						offset += missionLen
						continue
					}
				}
			}

			if offset+missionLen <= len(message) {
				mission := make(map[string]interface{})

				// In a real implementation, this would use proper unpacking based on questInfo["mission_pack"]
				// For simplicity, we'll just extract the mob_id and mob_count
				mobID := uint32(message[offset]) | uint32(message[offset+1])<<8 | uint32(message[offset+2])<<16 | uint32(message[offset+3])<<24
				mission["mob_id"] = mobID

				if offset+6 <= len(message) {
					mission["mob_count"] = uint16(message[offset+4]) | uint16(message[offset+5])<<8
				}

				// Extract mob name (simplified)
				mobName := "Unknown"
				mission["mob_name"] = mobName
				mission["mission_index"] = j

				// Add mission to quest
				if questID, ok := quest["quest_id"].(uint32); ok {
					if _, exists := qm.questList[questID]["missions"]; !exists {
						qm.questList[questID]["missions"] = make(map[uint32]map[string]interface{})
					}

					qm.questList[questID]["missions"].(map[uint32]map[string]interface{})[mobID] = mission
				}

				qm.logger.Debug("- MobID: %d - Name: %s - Count: %d", mobID, mobName, mission["mob_count"])

				// Call hook for mission added
				qm.hookManager.CallHook("quest_mission_added", map[string]interface{}{
					"questID":    quest["quest_id"],
					"mission_id": mobID,
				})
			}

			offset += missionLen
		}
	}

	// Call hook for quest all mission end
	qm.hookManager.CallHook("quest_all_mission_end", nil)

	return nil
}

// HandleQuestAdd handles the quest_add packet (lines 4759-4823)
func (qm *QuestManager) HandleQuestAdd(args map[string]interface{}) error {
	// Extract packet data
	questID, ok := args["questID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid questID in quest_add packet")
	}

	active, ok := args["active"].(uint8)
	if !ok {
		return fmt.Errorf("invalid active in quest_add packet")
	}

	timeStart, ok := args["time_start"].(uint32)
	if !ok {
		return fmt.Errorf("invalid time_start in quest_add packet")
	}

	timeExpire, ok := args["time_expire"].(uint32)
	if !ok {
		return fmt.Errorf("invalid time_expire in quest_add packet")
	}

	missionAmount, ok := args["mission_amount"].(uint16)
	if !ok {
		return fmt.Errorf("invalid mission_amount in quest_add packet")
	}

	message, ok := args["message"].([]byte)
	if !ok {
		return fmt.Errorf("invalid message in quest_add packet")
	}

	// Determine mission info based on switch
	var missionInfo map[string]interface{}
	switch args["switch"].(string) {
	case "09F9": // SERVERTYPE >= 20150513
		missionInfo = map[string]interface{}{
			"mission_pack": "V3 v3 Z24",
			"mission_keys": []string{"hunt_id", "mob_type", "mob_id", "min_level", "max_level", "mob_count", "mob_name_original"},
			"mission_len":  42,
		}
	case "0B0C": // SERVERTYPE >= 20150513
		missionInfo = map[string]interface{}{
			"mission_pack": "V4 v3 Z24",
			"mission_keys": []string{"hunt_id", "hunt_id_cont", "mob_type", "mob_id", "min_level", "max_level", "mob_count", "mob_name_original"},
			"mission_len":  46,
		}
	default: // DEFAULT PACKET - 02B3
		missionInfo = map[string]interface{}{
			"mission_pack": "V v Z24",
			"mission_keys": []string{"mob_id", "mob_count", "mob_name_original"},
			"mission_len":  30,
		}
	}

	// Create or get quest in questList
	if _, exists := qm.questList[questID]; !exists {
		qm.questList[questID] = make(map[string]interface{})
	}

	// Update quest data
	qm.questList[questID]["quest_id"] = questID
	qm.questList[questID]["active"] = active
	qm.questList[questID]["time_start"] = timeStart
	qm.questList[questID]["time_expire"] = timeExpire
	qm.questList[questID]["mission_amount"] = missionAmount

	// Log quest added
	questTitle := fmt.Sprintf("%d", questID) // In a real implementation, this would use quest title from a lookup table
	qm.logger.Info("Quest: %s has been added.", questTitle)

	offset := 0
	missionLen := missionInfo["mission_len"].(int)

	// Initialize missions map if not exists
	if _, exists := qm.questList[questID]["missions"]; !exists {
		qm.questList[questID]["missions"] = make(map[uint32]map[string]interface{})
	}

	// Process missions (always 3 in this packet)
	for j := 0; j < 3; j++ {
		// Skip if beyond mission amount
		if j >= int(missionAmount) {
			offset += missionLen
			continue
		}

		if offset+missionLen <= len(message) {
			mission := make(map[string]interface{})

			// In a real implementation, this would use proper unpacking based on missionInfo["mission_pack"]
			// For simplicity, we'll just extract the mob_id and mob_count
			mobID := uint32(message[offset]) | uint32(message[offset+1])<<8 | uint32(message[offset+2])<<16 | uint32(message[offset+3])<<24
			mission["mob_id"] = mobID

			if offset+6 <= len(message) {
				mission["mob_count"] = uint16(message[offset+4]) | uint16(message[offset+5])<<8
			}

			// Extract mob name (simplified)
			mobName := "Unknown"
			mission["mob_name"] = mobName
			mission["mission_index"] = j

			// Add mission to quest
			qm.questList[questID]["missions"].(map[uint32]map[string]interface{})[mobID] = mission

			qm.logger.Debug("- MobID: %d - Name: %s - Count: %d", mobID, mobName, mission["mob_count"])

			// Call hook for mission added
			qm.hookManager.CallHook("quest_mission_added", map[string]interface{}{
				"questID":    questID,
				"mission_id": mobID,
			})
		}

		offset += missionLen
	}

	// Call hook for quest added
	qm.hookManager.CallHook("quest_added", map[string]interface{}{
		"questID": questID,
	})

	return nil
}

// HandleQuestUpdateMissionHunt handles the quest_update_mission_hunt packet (lines 4827-4925)
func (qm *QuestManager) HandleQuestUpdateMissionHunt(args map[string]interface{}) error {
	// Extract packet data
	var missionAmount uint16
	var message []byte
	var ok bool

	// Handle 08FE which doesn't have mission_amount
	if args["switch"].(string) == "08FE" {
		message, ok = args["message"].([]byte)
		if !ok {
			return fmt.Errorf("invalid message in quest_update_mission_hunt packet")
		}

		// Calculate mission amount based on message length
		missionInfo := map[string]interface{}{
			"mission_pack": "V2 v2",
			"mission_keys": []string{"questID", "mob_id", "mob_goal", "mob_count"},
			"mission_len":  12,
		}

		missionAmount = uint16(len(message) / missionInfo["mission_len"].(int))
	} else {
		missionAmount, ok = args["mission_amount"].(uint16)
		if !ok {
			return fmt.Errorf("invalid mission_amount in quest_update_mission_hunt packet")
		}

		message, ok = args["message"].([]byte)
		if !ok {
			return fmt.Errorf("invalid message in quest_update_mission_hunt packet")
		}
	}

	// Determine mission info based on switch
	var missionInfo map[string]interface{}
	switch args["switch"].(string) {
	case "09FA":
		missionInfo = map[string]interface{}{
			"mission_pack": "V2 v2",
			"mission_keys": []string{"questID", "hunt_id", "mob_goal", "mob_count"},
			"mission_len":  12,
		}
	case "0AFE":
		missionInfo = map[string]interface{}{
			"mission_pack": "V3 v2",
			"mission_keys": []string{"questID", "hunt_id", "hunt_id_cont", "mob_goal", "mob_count"},
			"mission_len":  16,
		}
	default: // 02B5 and 08FE
		missionInfo = map[string]interface{}{
			"mission_pack": "V2 v2",
			"mission_keys": []string{"questID", "mob_id", "mob_goal", "mob_count"},
			"mission_len":  12,
		}
	}

	offset := 0
	missionLen := missionInfo["mission_len"].(int)

	for i := 0; i < int(missionAmount); i++ {
		if offset+missionLen <= len(message) {
			// Extract mission data
			questID := uint32(message[offset]) | uint32(message[offset+1])<<8 | uint32(message[offset+2])<<16 | uint32(message[offset+3])<<24

			var missionID uint32
			if args["switch"].(string) == "09FA" || args["switch"].(string) == "0AFE" {
				// Hunt ID
				missionID = uint32(message[offset+4]) | uint32(message[offset+5])<<8 | uint32(message[offset+6])<<16 | uint32(message[offset+7])<<24
			} else {
				// Mob ID
				missionID = uint32(message[offset+4]) | uint32(message[offset+5])<<8 | uint32(message[offset+6])<<16 | uint32(message[offset+7])<<24
			}

			// Extract mob goal and count
			var mobGoal, mobCount uint16
			if args["switch"].(string) == "0AFE" {
				mobGoal = uint16(message[offset+12]) | uint16(message[offset+13])<<8
				mobCount = uint16(message[offset+14]) | uint16(message[offset+15])<<8
			} else {
				mobGoal = uint16(message[offset+8]) | uint16(message[offset+9])<<8
				mobCount = uint16(message[offset+10]) | uint16(message[offset+11])<<8
			}

			// Get quest and mission
			if _, exists := qm.questList[questID]; !exists {
				qm.questList[questID] = make(map[string]interface{})
				qm.questList[questID]["missions"] = make(map[uint32]map[string]interface{})
			}

			// Find the correct mission ID
			var actualMissionID uint32

			// Check if mission exists with the given ID
			if missions, ok := qm.questList[questID]["missions"].(map[uint32]map[string]interface{}); ok {
				if _, exists := missions[missionID]; exists {
					actualMissionID = missionID
				} else {
					// Search for mission with matching mob_id or hunt_id
					for id, mission := range missions {
						if (args["switch"].(string) == "09FA" || args["switch"].(string) == "0AFE") &&
							mission["hunt_id"] != nil && mission["hunt_id"].(uint32) == missionID {
							actualMissionID = id
							break
						} else if mission["mob_id"] != nil && mission["mob_id"].(uint32) == missionID {
							actualMissionID = id
							break
						}
					}
				}
			}

			// If mission doesn't exist, create it
			if actualMissionID == 0 {
				actualMissionID = missionID
				qm.questList[questID]["missions"].(map[uint32]map[string]interface{})[actualMissionID] = make(map[string]interface{})
			}

			// Update mission data
			mission := qm.questList[questID]["missions"].(map[uint32]map[string]interface{})[actualMissionID]
			mission["mob_count"] = mobCount
			mission["mob_goal"] = mobGoal

			// Log mission update
			mobName := "Unknown"
			if mission["mob_name"] != nil {
				mobName = mission["mob_name"].(string)
			}

			qm.logger.Debug("- MobID: %d - Name: %s - Count: %d - Goal: %d", missionID, mobName, mobCount, mobGoal)

			// Display quest progress based on config
			// In a real implementation, this would check config settings
			qm.logger.Info("%s [%d/%d]", mobName, mobCount, mobGoal)

			// Call hook for mission updated
			qm.hookManager.CallHook("quest_mission_updated", map[string]interface{}{
				"questID":    questID,
				"mission_id": actualMissionID,
				"mobID":      missionID,
				"count":      mobCount,
				"goal":       mobGoal,
			})
		}

		offset += missionLen
	}

	// Call hook for quest update mission hunt end
	qm.hookManager.CallHook("quest_update_mission_hunt_end", nil)

	return nil
}

// HandleQuestDelete handles the quest_delete packet (lines 4928-4934)
func (qm *QuestManager) HandleQuestDelete(args map[string]interface{}) error {
	// Extract packet data
	questID, ok := args["questID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid questID in quest_delete packet")
	}

	// Log quest deleted
	questTitle := fmt.Sprintf("%d", questID) // In a real implementation, this would use quest title from a lookup table
	qm.logger.Info("Quest: %s has been deleted.", questTitle)

	// Delete quest from questList
	delete(qm.questList, questID)

	// Call hook for quest delete
	qm.hookManager.CallHook("quest_delete", nil)

	return nil
}

// HandleQuestActive handles the quest_active packet (lines 4937-4948)
func (qm *QuestManager) HandleQuestActive(args map[string]interface{}) error {
	// Extract packet data
	questID, ok := args["questID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid questID in quest_active packet")
	}

	active, ok := args["active"].(uint8)
	if !ok {
		return fmt.Errorf("invalid active in quest_active packet")
	}

	// Log quest active status
	questTitle := fmt.Sprintf("%d", questID) // In a real implementation, this would use quest title from a lookup table
	if active != 0 {
		qm.logger.Info("Quest %s is now active.", questTitle)
	} else {
		qm.logger.Info("Quest %s is now inactive.", questTitle)
	}

	// Update quest active status
	if _, exists := qm.questList[questID]; !exists {
		qm.questList[questID] = make(map[string]interface{})
	}

	qm.questList[questID]["active"] = active

	// Call hook for quest active
	qm.hookManager.CallHook("quest_active", nil)

	return nil
}

// IncrementQuestGeneration increments the quest generation counter
// This should be called when the character changes maps
func (qm *QuestManager) IncrementQuestGeneration() {
	qm.questGeneration++
}
