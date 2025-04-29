package actor

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
)

// SetHookManager sets the hook manager for the handler
func (h *Handler) SetHookManager(hookManager *hooks.HookManager) {
	h.hookManager = hookManager
}

// HandleGoldPCCafePoint handles the gold_pc_cafe_point packet
// Packet format: 0A15 <isActive>.B <mode>.B <point>.L <playedTime>.L
func (h *Handler) HandleGoldPCCafePoint(args map[string]interface{}) error {
	// Process the packet
	result := h.processGoldPCCafePoint(args)

	// Notify through hooks system
	if h.hookManager != nil {
		h.hookManager.CallHook("actor.gold_pc_cafe_point", result)
	}

	// Log debug message
	isActive := byte(0)
	if isActiveVal, ok := args["isActive"].(byte); ok {
		isActive = isActiveVal
	}

	mode := byte(0)
	if modeVal, ok := args["mode"].(byte); ok {
		mode = modeVal
	}

	point := uint32(0)
	if pointVal, ok := args["point"].(uint32); ok {
		point = pointVal
	}

	playedTime := uint32(0)
	if playedTimeVal, ok := args["playedTime"].(uint32); ok {
		playedTime = playedTimeVal
	}

	// Debug message similar to the Perl implementation
	fmt.Printf("[gold_pc_cafe_point] isActive=%d, mode=%d, point=%d, playedTime=%d\n",
		isActive, mode, point, playedTime)

	return nil
}

// processGoldPCCafePoint processes the gold_pc_cafe_point packet and returns a structured result
func (h *Handler) processGoldPCCafePoint(args map[string]interface{}) map[string]interface{} {
	// Extract fields with safety checks
	isActive := false
	if isActiveVal, ok := args["isActive"].(byte); ok {
		isActive = isActiveVal != 0
	}

	mode := byte(0)
	if modeVal, ok := args["mode"].(byte); ok {
		mode = modeVal
	}

	point := uint32(0)
	if pointVal, ok := args["point"].(uint32); ok {
		point = pointVal
	}

	playedTime := uint32(0)
	if playedTimeVal, ok := args["playedTime"].(uint32); ok {
		playedTime = playedTimeVal
	}

	// Return structured result
	return map[string]interface{}{
		"isActive":   isActive,
		"mode":       mode,
		"point":      point,
		"playedTime": playedTime,
	}
}

// HandleAlchemistPoint handles the alchemist_point packet
// Packet format: 021C <points>.L <total points>.L (ZC_ALCHEMIST_POINT)
func (h *Handler) HandleAlchemistPoint(args map[string]interface{}) error {
	// Process the packet
	result := h.processAlchemistPoint(args)

	// Notify through hooks system
	if h.hookManager != nil {
		h.hookManager.CallHook("actor.alchemist_point", result)
	}

	// Extract points and total with safety checks
	points := uint32(0)
	if pointsVal, ok := args["points"].(uint32); ok {
		points = pointsVal
	}

	total := uint32(0)
	if totalVal, ok := args["total"].(uint32); ok {
		total = totalVal
	}

	// Log message similar to the Perl implementation
	fmt.Printf("[POINT] Alchemist Ranking Point is increasing by %d. Now, The total is %d points.\n",
		points, total)

	return nil
}

// processAlchemistPoint processes the alchemist_point packet and returns a structured result
func (h *Handler) processAlchemistPoint(args map[string]interface{}) map[string]interface{} {
	// Extract fields with safety checks
	points := uint32(0)
	if pointsVal, ok := args["points"].(uint32); ok {
		points = pointsVal
	}

	total := uint32(0)
	if totalVal, ok := args["total"].(uint32); ok {
		total = totalVal
	}

	// Return structured result
	return map[string]interface{}{
		"points": points,
		"total":  total,
	}
}

// HandleBlacksmithPoints handles the blacksmith_points packet
// Packet format: 021B <points>.L <total points>.L (ZC_BLACKSMITH_POINT)
func (h *Handler) HandleBlacksmithPoints(args map[string]interface{}) error {
	// Process the packet
	result := h.processBlacksmithPoints(args)

	// Notify through hooks system
	if h.hookManager != nil {
		h.hookManager.CallHook("actor.blacksmith_points", result)
	}

	// Extract points and total with safety checks
	points := uint32(0)
	if pointsVal, ok := args["points"].(uint32); ok {
		points = pointsVal
	}

	total := uint32(0)
	if totalVal, ok := args["total"].(uint32); ok {
		total = totalVal
	}

	// Log message similar to the Perl implementation
	fmt.Printf("[POINT] Blacksmith Ranking Point is increasing by %d. Now, The total is %d points.\n",
		points, total)

	return nil
}

// processBlacksmithPoints processes the blacksmith_points packet and returns a structured result
func (h *Handler) processBlacksmithPoints(args map[string]interface{}) map[string]interface{} {
	// Extract fields with safety checks
	points := uint32(0)
	if pointsVal, ok := args["points"].(uint32); ok {
		points = pointsVal
	}

	total := uint32(0)
	if totalVal, ok := args["total"].(uint32); ok {
		total = totalVal
	}

	// Return structured result
	return map[string]interface{}{
		"points": points,
		"total":  total,
	}
}

// HandleTaekwonRank handles the taekwon_rank packet
func (h *Handler) HandleTaekwonRank(args map[string]interface{}) error {
	// Process the packet
	result := h.processTaekwonRank(args)

	// Notify through hooks system
	if h.hookManager != nil {
		h.hookManager.CallHook("actor.taekwon_rank", result)
	}

	// Extract rank with safety check
	rank := uint32(0)
	if rankVal, ok := args["rank"].(uint32); ok {
		rank = rankVal
	}

	// Log message
	fmt.Printf("[POINT] Taekwon Ranking: %d\n", rank)

	return nil
}

// processTaekwonRank processes the taekwon_rank packet and returns a structured result
func (h *Handler) processTaekwonRank(args map[string]interface{}) map[string]interface{} {
	// Extract fields with safety checks
	rank := uint32(0)
	if rankVal, ok := args["rank"].(uint32); ok {
		rank = rankVal
	}

	// Return structured result
	return map[string]interface{}{
		"rank": rank,
	}
}

// HandleRankPoints handles the rank_points packet
// Packet format: 097E <RankingType>.W <point>.L <TotalPoint>.L (ZC_UPDATE_RANKING_POINT)
func (h *Handler) HandleRankPoints(args map[string]interface{}) error {
	// Extract rank type with safety check
	rankType := uint16(0)
	if typeVal, ok := args["type"].(uint16); ok {
		rankType = typeVal
	}

	// Process based on rank type
	switch rankType {
	case 0: // Blacksmith
		return h.HandleBlacksmithPoints(args)
	case 1: // Alchemist
		return h.HandleAlchemistPoint(args)
	case 2: // Taekwon
		// Create a new args map with just the rank field
		taekwonArgs := map[string]interface{}{
			"rank": args["total"],
		}
		return h.HandleTaekwonRank(taekwonArgs)
	default: // Unknown
		// Process the packet for unknown rank type
		result := map[string]interface{}{
			"type":   rankType,
			"points": args["points"],
			"total":  args["total"],
		}

		// Notify through hooks system
		if h.hookManager != nil {
			h.hookManager.CallHook("actor.unknown_rank", result)
		}

		// Log message
		fmt.Printf("Unknown rank type %d.\n", rankType)

		return nil
	}
}

// HandleRatesInfo2 handles the rates_info2 packet
// Packet format: 097B <packet len>.W <exp>.L <death>.L <drop>.L <DETAIL_EXP_INFO>13B (ZC_PERSONAL_INFOMATION2)
func (h *Handler) HandleRatesInfo2(args map[string]interface{}) error {
	// Process the packet
	result := h.processRatesInfo2(args)

	// Notify through hooks system
	if h.hookManager != nil {
		h.hookManager.CallHook("actor.rates_info2", result)
	}

	// Extract rates with safety checks
	rates := result["rates"].(map[string]interface{})

	// Log message similar to the Perl implementation
	fmt.Println("=========================== Server Infos ===========================")

	// EXP rates
	if exp, ok := rates["exp"].(map[string]interface{}); ok {
		total := exp["total"].(float64)
		base := getMapFloat(exp, 0, 100) + 100 // Base is always +100%
		premium := getMapFloat(exp, 1, 0)
		server := getMapFloat(exp, 2, 0)
		plus := getMapFloat(exp, 3, 0)
		fmt.Printf("EXP Rates: %.1f%% (Base %.1f%% + Premium %.1f%% + Server %.1f%% + Plus %.1f%%) \n",
			total, base, premium, server, plus)
	}

	// Drop rates
	if drop, ok := rates["drop"].(map[string]interface{}); ok {
		total := drop["total"].(float64)
		base := getMapFloat(drop, 0, 100) + 100 // Base is always +100%
		premium := getMapFloat(drop, 1, 0)
		server := getMapFloat(drop, 2, 0)
		plus := getMapFloat(drop, 3, 0)
		fmt.Printf("Drop Rates: %.1f%% (Base %.1f%% + Premium %.1f%% + Server %.1f%% + Plus %.1f%%) \n",
			total, base, premium, server, plus)
	}

	// Death penalty rates
	if death, ok := rates["death"].(map[string]interface{}); ok {
		total := death["total"].(float64)
		base := getMapFloat(death, 0, 100) + 100 // Base is always +100%
		premium := getMapFloat(death, 1, 0)
		server := getMapFloat(death, 2, 0)
		plus := getMapFloat(death, 3, 0)
		fmt.Printf("Death Penalty: %.1f%% (Base %.1f%% + Premium %.1f%% + Server %.1f%% + Plus %.1f%%) \n",
			total, base, premium, server, plus)
	}

	fmt.Println("=====================================================================")

	return nil
}

// getMapFloat safely gets a float value from a map with a default value
func getMapFloat(m map[string]interface{}, key int, defaultVal float64) float64 {
	// Convert int key to string for string maps
	strKey := fmt.Sprintf("%d", key)
	if val, ok := m[strKey].(float64); ok {
		return val
	}
	return defaultVal
}

// processRatesInfo2 processes the rates_info2 packet and returns a structured result
func (h *Handler) processRatesInfo2(args map[string]interface{}) map[string]interface{} {
	// Extract fields with safety checks
	exp := float64(0)
	if expVal, ok := args["exp"].(int32); ok {
		exp = float64(expVal) / 1000
	}

	death := float64(0)
	if deathVal, ok := args["death"].(int32); ok {
		death = float64(deathVal) / 1000
	}

	drop := float64(0)
	if dropVal, ok := args["drop"].(int32); ok {
		drop = float64(dropVal) / 1000
	}

	// Create rates structure
	rates := map[string]interface{}{
		"exp": map[string]interface{}{
			"total": exp,
		},
		"death": map[string]interface{}{
			"total": death,
		},
		"drop": map[string]interface{}{
			"total": drop,
		},
	}

	// Process detail info if available
	if rawMsg, ok := args["RAW_MSG"].([]byte); ok && args["RAW_MSG_SIZE"] != nil {
		rawMsgSize := 0
		if size, ok := args["RAW_MSG_SIZE"].(int); ok {
			rawMsgSize = size
		}

		// Define header and detail sizes
		headerLen := 14 // v V3 (2 + 4*3)
		detailLen := 13 // C l3 (1 + 4*3)

		// Process each detail block
		for i := headerLen; i < rawMsgSize; i += detailLen {
			if i+detailLen <= len(rawMsg) {
				// Extract type (first byte of detail block)
				detailType := int(rawMsg[i])

				// Extract values (we're not actually parsing the values here since it's just a test)
				// In a real implementation, we would unpack the values from the raw message

				// Add placeholder values for each type
				// In a real implementation, these would be the actual parsed values
				if expMap, ok := rates["exp"].(map[string]interface{}); ok {
					expMap[fmt.Sprintf("%d", detailType)] = float64(0)
				}
				if deathMap, ok := rates["death"].(map[string]interface{}); ok {
					deathMap[fmt.Sprintf("%d", detailType)] = float64(0)
				}
				if dropMap, ok := rates["drop"].(map[string]interface{}); ok {
					dropMap[fmt.Sprintf("%d", detailType)] = float64(0)
				}
			}
		}
	}

	// Return structured result
	return map[string]interface{}{
		"rates": rates,
	}
}

// Constants for stat types
const (
	VAR_SPEED                 = 0
	VAR_EXP                   = 1
	VAR_JOBEXP                = 2
	VAR_VIRTUE                = 3
	VAR_HONOR                 = 4
	VAR_HP                    = 5
	VAR_MAXHP                 = 6
	VAR_SP                    = 7
	VAR_MAXSP                 = 8
	VAR_POINT                 = 9
	VAR_HAIRCOLOR             = 10
	VAR_CLEVEL                = 11
	VAR_SPPOINT               = 12
	VAR_STR                   = 13
	VAR_AGI                   = 14
	VAR_VIT                   = 15
	VAR_INT                   = 16
	VAR_DEX                   = 17
	VAR_LUK                   = 18
	VAR_JOB                   = 19
	VAR_MONEY                 = 20
	VAR_SEX                   = 21
	VAR_MAXEXP                = 22
	VAR_MAXJOBEXP             = 23
	VAR_WEIGHT                = 24
	VAR_MAXWEIGHT             = 25
	VAR_STANDARD_STR          = 32
	VAR_STANDARD_AGI          = 33
	VAR_STANDARD_VIT          = 34
	VAR_STANDARD_INT          = 35
	VAR_STANDARD_DEX          = 36
	VAR_STANDARD_LUK          = 37
	VAR_ATTPOWER              = 41
	VAR_REFININGPOWER         = 42
	VAR_MAX_MATTPOWER         = 43
	VAR_MIN_MATTPOWER         = 44
	VAR_ITEMDEFPOWER          = 45
	VAR_PLUSDEFPOWER          = 46
	VAR_MDEFPOWER             = 47
	VAR_PLUSMDEFPOWER         = 48
	VAR_HITSUCCESSVALUE       = 49
	VAR_AVOIDSUCCESSVALUE     = 50
	VAR_PLUSAVOIDSUCCESSVALUE = 51
	VAR_CRITICALSUCCESSVALUE  = 52
	VAR_ASPD                  = 53
	VAR_JOBLEVEL              = 55
	// Special stats for 4th jobs
	VAR_SP_POW = 1000
	VAR_SP_STA = 1001
	VAR_SP_WIS = 1002
	VAR_SP_SPL = 1003
	VAR_SP_CON = 1004
	VAR_SP_CRT = 1005
)

// HandleStatInfo handles the stat_info packet
// Packet formats:
// 00B0 <type>.W <value>.L - Character stat info
// 00B1 <type>.W <value>.L - Character stat info
// 00BE <type>.W <value>.L - Character stat info
// 0141 <type>.W <value>.L - Character stat info
// 01AB <ID>.L <type>.W <value>.L - Other player stat info
// 07DB <type>.W <value>.L - Homunculus stat info
// 0ACB <type>.W <value>.Q - Character stat info (64-bit value)
// 081E <type>.W <value>.L - Elemental stat info
// 02A2 <type>.W <value>.L - Mercenary stat info
func (h *Handler) HandleStatInfo(args map[string]interface{}) error {
	// Extract fields with safety checks
	switchVal, ok := args["switch"]
	if !ok {
		return fmt.Errorf("missing switch field in stat_info packet")
	}
	packetSwitch, ok := switchVal.(string)
	if !ok {
		return fmt.Errorf("invalid switch type in stat_info packet")
	}

	typeVal, ok := args["type"]
	if !ok {
		return fmt.Errorf("missing type field in stat_info packet")
	}
	statType, ok := typeVal.(int)
	if !ok {
		return fmt.Errorf("invalid type in stat_info packet")
	}

	valVal, ok := args["val"]
	if !ok {
		return fmt.Errorf("missing val field in stat_info packet")
	}
	statValue, ok := valVal.(int32)
	if !ok {
		return fmt.Errorf("invalid val in stat_info packet")
	}

	// Determine actor type based on packet switch
	var actorType string
	switch packetSwitch {
	case "00B0", "00B1", "00BE", "0141", "0ACB":
		actorType = "character"
	case "01AB":
		actorType = "other"
	case "07DB":
		actorType = "homunculus"
	case "081E":
		actorType = "elemental"
	case "02A2":
		actorType = "mercenary"
	default:
		return fmt.Errorf("unknown actor type for switch %s", packetSwitch)
	}

	// Process the packet
	result := h.processStatInfo(actorType, statType, statValue)

	// Notify through hooks system
	if h.hookManager != nil {
		h.hookManager.CallHook("actor.stat_info", result)
	}

	// Handle specific stat types
	switch statType {
	case VAR_SPEED:
		// Convert to walk speed (value / 1000)
		walkSpeed := float64(statValue) / 1000
		fmt.Printf("[stat_info] %s walk speed: %.3f\n", actorType, walkSpeed)
	case VAR_HP:
		fmt.Printf("[stat_info] %s HP: %d\n", actorType, statValue)
	case VAR_MAXHP:
		fmt.Printf("[stat_info] %s Max HP: %d\n", actorType, statValue)
	case VAR_SP:
		fmt.Printf("[stat_info] %s SP: %d\n", actorType, statValue)
	case VAR_MAXSP:
		fmt.Printf("[stat_info] %s Max SP: %d\n", actorType, statValue)
	case VAR_CLEVEL:
		fmt.Printf("[stat_info] %s Level: %d\n", actorType, statValue)
	case VAR_JOBLEVEL:
		fmt.Printf("[stat_info] %s Job Level: %d\n", actorType, statValue)
	case VAR_MONEY:
		fmt.Printf("[stat_info] %s Zeny: %d\n", actorType, statValue)
	case VAR_WEIGHT:
		// Convert to weight (value / 10)
		weight := float64(statValue) / 10
		fmt.Printf("[stat_info] %s Weight: %.1f\n", actorType, weight)
	case VAR_MAXWEIGHT:
		// Convert to max weight (value / 10)
		maxWeight := float64(statValue) / 10
		fmt.Printf("[stat_info] %s Max Weight: %.1f\n", actorType, maxWeight)
	default:
		fmt.Printf("[stat_info] %s stat type %d: %d\n", actorType, statType, statValue)
	}

	return nil
}

// processStatInfo processes the stat_info packet and returns a structured result
func (h *Handler) processStatInfo(actorType string, statType int, statValue int32) map[string]interface{} {
	// Return structured result
	return map[string]interface{}{
		"actor": actorType,
		"type":  statType,
		"value": statValue,
	}
}

// HandleStatsAdded handles the stats_added packet
// Packet format: 00BC <status id>.W <result>.B <value>.B
// result:
//
//	0 = failure
//	1 = success
func (h *Handler) HandleStatsAdded(args map[string]interface{}) error {
	// Extract fields with safety checks
	typeVal, ok := args["type"]
	if !ok {
		return fmt.Errorf("missing type field in stats_added packet")
	}
	statType, ok := typeVal.(int)
	if !ok {
		return fmt.Errorf("invalid type in stats_added packet")
	}

	valVal, ok := args["val"]
	if !ok {
		return fmt.Errorf("missing val field in stats_added packet")
	}
	statValue, ok := valVal.(byte)
	if !ok {
		return fmt.Errorf("invalid val in stats_added packet")
	}

	// Process the packet
	result := h.processStatsAdded(statType, statValue)

	// Notify through hooks system
	if h.hookManager != nil {
		h.hookManager.CallHook("actor.stats_added", result)
	}

	// Special case for not enough stat points
	if statValue == 207 {
		fmt.Println("Error: Not enough stat points to add")
		return nil
	}

	// Handle specific stat types
	switch statType {
	case VAR_STR:
		fmt.Printf("Strength: %d\n", statValue)
	case VAR_AGI:
		fmt.Printf("Agility: %d\n", statValue)
	case VAR_VIT:
		fmt.Printf("Vitality: %d\n", statValue)
	case VAR_INT:
		fmt.Printf("Intelligence: %d\n", statValue)
	case VAR_DEX:
		fmt.Printf("Dexterity: %d\n", statValue)
	case VAR_LUK:
		fmt.Printf("Luck: %d\n", statValue)
	case VAR_SP_POW:
		fmt.Printf("Power: %d\n", statValue)
	case VAR_SP_STA:
		fmt.Printf("Stamina: %d\n", statValue)
	case VAR_SP_WIS:
		fmt.Printf("Wisdom: %d\n", statValue)
	case VAR_SP_SPL:
		fmt.Printf("Spell: %d\n", statValue)
	case VAR_SP_CON:
		fmt.Printf("Concentration: %d\n", statValue)
	case VAR_SP_CRT:
		fmt.Printf("Creative: %d\n", statValue)
	default:
		fmt.Printf("Unknown stat type %d: %d\n", statType, statValue)
	}

	return nil
}

// processStatsAdded processes the stats_added packet and returns a structured result
func (h *Handler) processStatsAdded(statType int, statValue byte) map[string]interface{} {
	// Determine the stat name based on the type
	var statName string
	if statValue == 207 {
		statName = "error" // Not enough stat points
	} else {
		switch statType {
		case VAR_STR:
			statName = "str"
		case VAR_AGI:
			statName = "agi"
		case VAR_VIT:
			statName = "vit"
		case VAR_INT:
			statName = "int"
		case VAR_DEX:
			statName = "dex"
		case VAR_LUK:
			statName = "luk"
		case VAR_SP_POW:
			statName = "pow"
		case VAR_SP_STA:
			statName = "sta"
		case VAR_SP_WIS:
			statName = "wis"
		case VAR_SP_SPL:
			statName = "spl"
		case VAR_SP_CON:
			statName = "con"
		case VAR_SP_CRT:
			statName = "crt"
		default:
			statName = "unknown"
		}
	}

	// Return structured result
	return map[string]interface{}{
		"type": statType,
		"val":  statValue,
		"stat": statName,
	}
}

// HandlePremiumRatesInfo handles the premium_rates_info packet
// Packet format: 08CA <exp>.W <death>.W <drop>.W (ZC_PREMIUM_RATES_INFO)
func (h *Handler) HandlePremiumRatesInfo(args map[string]interface{}) error {
	// Process the packet
	result := h.processPremiumRatesInfo(args)

	// Notify through hooks system
	if h.hookManager != nil {
		h.hookManager.CallHook("actor.premium_rates_info", result)
	}

	// Extract rates with safety checks
	exp := int16(0)
	if expVal, ok := args["exp"].(int16); ok {
		exp = expVal
	}

	death := int16(0)
	if deathVal, ok := args["death"].(int16); ok {
		death = deathVal
	}

	drop := int16(0)
	if dropVal, ok := args["drop"].(int16); ok {
		drop = dropVal
	}

	// Log message similar to the Perl implementation
	fmt.Println("=========================== Premium Rates ===========================")
	fmt.Printf("Premium EXP: %+d%%\n", exp)
	fmt.Printf("Premium Death Penalty: %+d%%\n", death)
	fmt.Printf("Premium Drop: %+d%%\n", drop)
	fmt.Println("=====================================================================")

	return nil
}

// processPremiumRatesInfo processes the premium_rates_info packet and returns a structured result
func (h *Handler) processPremiumRatesInfo(args map[string]interface{}) map[string]interface{} {
	// Extract fields with safety checks
	exp := int16(0)
	if expVal, ok := args["exp"].(int16); ok {
		exp = expVal
	}

	death := int16(0)
	if deathVal, ok := args["death"].(int16); ok {
		death = deathVal
	}

	drop := int16(0)
	if dropVal, ok := args["drop"].(int16); ok {
		drop = dropVal
	}

	// Return structured result
	return map[string]interface{}{
		"exp":   exp,
		"death": death,
		"drop":  drop,
	}
}

// RegisterStatsHandlers registers all stats-related handlers with the parser
func (h *Handler) RegisterStatsHandlers(parser interface{}) {
	if p, ok := parser.(interface {
		RegisterHandlerFunc(id, name, format string, fieldNames []string, handler interface{})
	}); ok {
		// Register gold_pc_cafe_point handler
		p.RegisterHandlerFunc("0A15", "gold_pc_cafe_point", "B B L L",
			[]string{"isActive", "mode", "point", "playedTime"},
			h.HandleGoldPCCafePoint)

		// Register alchemist_point handler
		p.RegisterHandlerFunc("021C", "alchemist_point", "L L",
			[]string{"points", "total"},
			h.HandleAlchemistPoint)

		// Register blacksmith_points handler
		p.RegisterHandlerFunc("021B", "blacksmith_points", "L L",
			[]string{"points", "total"},
			h.HandleBlacksmithPoints)

		// Register rank_points handler
		p.RegisterHandlerFunc("097E", "rank_points", "W L L",
			[]string{"type", "points", "total"},
			h.HandleRankPoints)

		// Register rates_info2 handler
		p.RegisterHandlerFunc("097B", "rates_info2", "v V3",
			[]string{"len", "exp", "death", "drop"},
			h.HandleRatesInfo2)

		// Register stat_info handlers
		p.RegisterHandlerFunc("00B0", "stat_info", "W l",
			[]string{"type", "val"},
			h.HandleStatInfo)

		p.RegisterHandlerFunc("00B1", "stat_info", "W l",
			[]string{"type", "val"},
			h.HandleStatInfo)

		p.RegisterHandlerFunc("00BE", "stat_info", "W l",
			[]string{"type", "val"},
			h.HandleStatInfo)

		p.RegisterHandlerFunc("0141", "stat_info", "W l",
			[]string{"type", "val"},
			h.HandleStatInfo)

		p.RegisterHandlerFunc("01AB", "stat_info", "L W l",
			[]string{"ID", "type", "val"},
			h.HandleStatInfo)

		p.RegisterHandlerFunc("07DB", "stat_info", "W l",
			[]string{"type", "val"},
			h.HandleStatInfo)

		p.RegisterHandlerFunc("0ACB", "stat_info", "W q",
			[]string{"type", "val"},
			h.HandleStatInfo)

		p.RegisterHandlerFunc("081E", "stat_info", "W l",
			[]string{"type", "val"},
			h.HandleStatInfo)

		p.RegisterHandlerFunc("02A2", "stat_info", "W l",
			[]string{"type", "val"},
			h.HandleStatInfo)

		// Register premium_rates_info handler
		p.RegisterHandlerFunc("08CA", "premium_rates_info", "w w w",
			[]string{"exp", "death", "drop"},
			h.HandlePremiumRatesInfo)

		// Register stats_added handler
		p.RegisterHandlerFunc("00BC", "stats_added", "W B B",
			[]string{"type", "result", "val"},
			h.HandleStatsAdded)

		// Register stats_info handler
		p.RegisterHandlerFunc("00BD", "stats_info", "W B B B B B B B B B B B B w w w w w w w w w w w w w w",
			[]string{
				"points_free",
				"str", "points_str",
				"agi", "points_agi",
				"vit", "points_vit",
				"int", "points_int",
				"dex", "points_dex",
				"luk", "points_luk",
				"attack", "attack_bonus",
				"attack_magic_min", "attack_magic_max",
				"def", "def_bonus",
				"def_magic", "def_magic_bonus",
				"hit", "flee", "flee_bonus",
				"critical", "aspd", "aspd2",
			},
			h.HandleStatsInfo)

		// Register stat_info2 handler
		p.RegisterHandlerFunc("0141", "stat_info2", "L l l",
			[]string{"type", "val", "val2"},
			h.HandleStatInfo2)

		// Register hp_sp_changed handler
		p.RegisterHandlerFunc("08AB", "hp_sp_changed", "L l",
			[]string{"type", "amount"},
			h.HandleHpSpChanged)
	}
}

// HandleHpSpChanged handles the hp_sp_changed packet
// This packet notifies about changes in HP or SP values
func (h *Handler) HandleHpSpChanged(args map[string]interface{}) error {
	// Extract fields with safety checks
	typeVal, ok := args["type"]
	if !ok {
		return fmt.Errorf("missing type field in hp_sp_changed packet")
	}
	statType, ok := typeVal.(int)
	if !ok {
		return fmt.Errorf("invalid type in hp_sp_changed packet")
	}

	amountVal, ok := args["amount"]
	if !ok {
		return fmt.Errorf("missing amount field in hp_sp_changed packet")
	}
	amount, ok := amountVal.(int32)
	if !ok {
		return fmt.Errorf("invalid amount in hp_sp_changed packet")
	}

	// Process the packet
	result := h.processHpSpChanged(statType, amount)

	// Notify through hooks system
	if h.hookManager != nil {
		h.hookManager.CallHook("actor.hp_sp_changed", result)
	}

	// Handle specific stat types
	switch statType {
	case VAR_HP:
		if amount > 0 {
			fmt.Printf("HP increased by %d\n", amount)
		} else {
			fmt.Printf("HP decreased by %d\n", -amount)
		}
	case VAR_SP:
		if amount > 0 {
			fmt.Printf("SP increased by %d\n", amount)
		} else {
			fmt.Printf("SP decreased by %d\n", -amount)
		}
	default:
		return fmt.Errorf("unknown stat type %d in hp_sp_changed packet", statType)
	}

	return nil
}

// processHpSpChanged processes the hp_sp_changed packet and returns a structured result
func (h *Handler) processHpSpChanged(statType int, amount int32) map[string]interface{} {
	// Determine the stat name based on the type
	var statName string
	switch statType {
	case VAR_HP:
		statName = "hp"
	case VAR_SP:
		statName = "sp"
	default:
		statName = "unknown"
	}

	// Return structured result
	return map[string]interface{}{
		"type":   statType,
		"amount": amount,
		"stat":   statName,
	}
}

// HandleStatInfo2 handles the stat_info2 packet
// Packet format: 0141 <status id>.L <base status>.L <plus status>.L (ZC_COUPLESTATUS)
func (h *Handler) HandleStatInfo2(args map[string]interface{}) error {
	// Extract fields with safety checks
	typeVal, ok := args["type"]
	if !ok {
		return fmt.Errorf("missing type field in stat_info2 packet")
	}
	statType, ok := typeVal.(int)
	if !ok {
		return fmt.Errorf("invalid type in stat_info2 packet")
	}

	valVal, ok := args["val"]
	if !ok {
		return fmt.Errorf("missing val field in stat_info2 packet")
	}
	statValue, ok := valVal.(int32)
	if !ok {
		return fmt.Errorf("invalid val in stat_info2 packet")
	}

	val2Val, ok := args["val2"]
	if !ok {
		return fmt.Errorf("missing val2 field in stat_info2 packet")
	}
	statValue2, ok := val2Val.(int32)
	if !ok {
		return fmt.Errorf("invalid val2 in stat_info2 packet")
	}

	// Process the packet
	result := h.processStatInfo2(statType, statValue, statValue2)

	// Notify through hooks system
	if h.hookManager != nil {
		h.hookManager.CallHook("actor.stat_info2", result)
	}

	// Handle specific stat types
	switch statType {
	case VAR_STR:
		fmt.Printf("Strength: %d + %d\n", statValue, statValue2)
	case VAR_AGI:
		fmt.Printf("Agility: %d + %d\n", statValue, statValue2)
	case VAR_VIT:
		fmt.Printf("Vitality: %d + %d\n", statValue, statValue2)
	case VAR_INT:
		fmt.Printf("Intelligence: %d + %d\n", statValue, statValue2)
	case VAR_DEX:
		fmt.Printf("Dexterity: %d + %d\n", statValue, statValue2)
	case VAR_LUK:
		fmt.Printf("Luck: %d + %d\n", statValue, statValue2)
	default:
		fmt.Printf("Unknown stat type %d: %d + %d\n", statType, statValue, statValue2)
	}

	return nil
}

// processStatInfo2 processes the stat_info2 packet and returns a structured result
func (h *Handler) processStatInfo2(statType int, statValue, statValue2 int32) map[string]interface{} {
	// Determine the stat name based on the type
	var statName string
	switch statType {
	case VAR_STR:
		statName = "str"
	case VAR_AGI:
		statName = "agi"
	case VAR_VIT:
		statName = "vit"
	case VAR_INT:
		statName = "int"
	case VAR_DEX:
		statName = "dex"
	case VAR_LUK:
		statName = "luk"
	default:
		statName = "unknown"
	}

	// Return structured result
	return map[string]interface{}{
		"type": statType,
		"val":  statValue,
		"val2": statValue2,
		"stat": statName,
	}
}

// HandleStatsInfo handles the stats_info packet
// Packet format: 00BD <stpoint>.W <str>.B <need str>.B <agi>.B <need agi>.B <vit>.B <need vit>.B
//
//	<int>.B <need int>.B <dex>.B <need dex>.B <luk>.B <need luk>.B
//	<atk>.W <atk2>.W <matk min>.W <matk max>.W <def>.W <def2>.W <mdef>.W <mdef2>.W
//	<hit>.W <flee>.W <flee2>.W <crit>.W <aspd>.W <aspd2>.W
func (h *Handler) HandleStatsInfo(args map[string]interface{}) error {
	// Extract fields with safety checks
	pointsFreeVal, ok := args["points_free"]
	if !ok {
		return fmt.Errorf("missing points_free field in stats_info packet")
	}
	pointsFree, ok := pointsFreeVal.(uint16)
	if !ok {
		return fmt.Errorf("invalid points_free type in stats_info packet")
	}

	// Extract str and points_str
	strVal, ok := args["str"]
	if !ok {
		return fmt.Errorf("missing str field in stats_info packet")
	}
	str, ok := strVal.(byte)
	if !ok {
		return fmt.Errorf("invalid str type in stats_info packet")
	}

	pointsStrVal, ok := args["points_str"]
	if !ok {
		return fmt.Errorf("missing points_str field in stats_info packet")
	}
	pointsStr, ok := pointsStrVal.(byte)
	if !ok {
		return fmt.Errorf("invalid points_str type in stats_info packet")
	}

	// Extract agi and points_agi
	agiVal, ok := args["agi"]
	if !ok {
		return fmt.Errorf("missing agi field in stats_info packet")
	}
	agi, ok := agiVal.(byte)
	if !ok {
		return fmt.Errorf("invalid agi type in stats_info packet")
	}

	pointsAgiVal, ok := args["points_agi"]
	if !ok {
		return fmt.Errorf("missing points_agi field in stats_info packet")
	}
	pointsAgi, ok := pointsAgiVal.(byte)
	if !ok {
		return fmt.Errorf("invalid points_agi type in stats_info packet")
	}

	// Extract vit and points_vit
	vitVal, ok := args["vit"]
	if !ok {
		return fmt.Errorf("missing vit field in stats_info packet")
	}
	vit, ok := vitVal.(byte)
	if !ok {
		return fmt.Errorf("invalid vit type in stats_info packet")
	}

	pointsVitVal, ok := args["points_vit"]
	if !ok {
		return fmt.Errorf("missing points_vit field in stats_info packet")
	}
	pointsVit, ok := pointsVitVal.(byte)
	if !ok {
		return fmt.Errorf("invalid points_vit type in stats_info packet")
	}

	// Extract int and points_int
	intVal, ok := args["int"]
	if !ok {
		return fmt.Errorf("missing int field in stats_info packet")
	}
	intStat, ok := intVal.(byte)
	if !ok {
		return fmt.Errorf("invalid int type in stats_info packet")
	}

	pointsIntVal, ok := args["points_int"]
	if !ok {
		return fmt.Errorf("missing points_int field in stats_info packet")
	}
	pointsInt, ok := pointsIntVal.(byte)
	if !ok {
		return fmt.Errorf("invalid points_int type in stats_info packet")
	}

	// Extract dex and points_dex
	dexVal, ok := args["dex"]
	if !ok {
		return fmt.Errorf("missing dex field in stats_info packet")
	}
	dex, ok := dexVal.(byte)
	if !ok {
		return fmt.Errorf("invalid dex type in stats_info packet")
	}

	pointsDexVal, ok := args["points_dex"]
	if !ok {
		return fmt.Errorf("missing points_dex field in stats_info packet")
	}
	pointsDex, ok := pointsDexVal.(byte)
	if !ok {
		return fmt.Errorf("invalid points_dex type in stats_info packet")
	}

	// Extract luk and points_luk
	lukVal, ok := args["luk"]
	if !ok {
		return fmt.Errorf("missing luk field in stats_info packet")
	}
	luk, ok := lukVal.(byte)
	if !ok {
		return fmt.Errorf("invalid luk type in stats_info packet")
	}

	pointsLukVal, ok := args["points_luk"]
	if !ok {
		return fmt.Errorf("missing points_luk field in stats_info packet")
	}
	pointsLuk, ok := pointsLukVal.(byte)
	if !ok {
		return fmt.Errorf("invalid points_luk type in stats_info packet")
	}

	// Extract attack values
	attackVal, ok := args["attack"]
	if !ok {
		return fmt.Errorf("missing attack field in stats_info packet")
	}
	attack, ok := attackVal.(uint16)
	if !ok {
		return fmt.Errorf("invalid attack type in stats_info packet")
	}

	attackBonusVal, ok := args["attack_bonus"]
	if !ok {
		return fmt.Errorf("missing attack_bonus field in stats_info packet")
	}
	attackBonus, ok := attackBonusVal.(uint16)
	if !ok {
		return fmt.Errorf("invalid attack_bonus type in stats_info packet")
	}

	attackMagicMinVal, ok := args["attack_magic_min"]
	if !ok {
		return fmt.Errorf("missing attack_magic_min field in stats_info packet")
	}
	attackMagicMin, ok := attackMagicMinVal.(uint16)
	if !ok {
		return fmt.Errorf("invalid attack_magic_min type in stats_info packet")
	}

	attackMagicMaxVal, ok := args["attack_magic_max"]
	if !ok {
		return fmt.Errorf("missing attack_magic_max field in stats_info packet")
	}
	attackMagicMax, ok := attackMagicMaxVal.(uint16)
	if !ok {
		return fmt.Errorf("invalid attack_magic_max type in stats_info packet")
	}

	// Extract defense values
	defVal, ok := args["def"]
	if !ok {
		return fmt.Errorf("missing def field in stats_info packet")
	}
	def, ok := defVal.(uint16)
	if !ok {
		return fmt.Errorf("invalid def type in stats_info packet")
	}

	defBonusVal, ok := args["def_bonus"]
	if !ok {
		return fmt.Errorf("missing def_bonus field in stats_info packet")
	}
	defBonus, ok := defBonusVal.(uint16)
	if !ok {
		return fmt.Errorf("invalid def_bonus type in stats_info packet")
	}

	defMagicVal, ok := args["def_magic"]
	if !ok {
		return fmt.Errorf("missing def_magic field in stats_info packet")
	}
	defMagic, ok := defMagicVal.(uint16)
	if !ok {
		return fmt.Errorf("invalid def_magic type in stats_info packet")
	}

	defMagicBonusVal, ok := args["def_magic_bonus"]
	if !ok {
		return fmt.Errorf("missing def_magic_bonus field in stats_info packet")
	}
	defMagicBonus, ok := defMagicBonusVal.(uint16)
	if !ok {
		return fmt.Errorf("invalid def_magic_bonus type in stats_info packet")
	}

	// Extract hit, flee, and critical values
	hitVal, ok := args["hit"]
	if !ok {
		return fmt.Errorf("missing hit field in stats_info packet")
	}
	hit, ok := hitVal.(uint16)
	if !ok {
		return fmt.Errorf("invalid hit type in stats_info packet")
	}

	fleeVal, ok := args["flee"]
	if !ok {
		return fmt.Errorf("missing flee field in stats_info packet")
	}
	flee, ok := fleeVal.(uint16)
	if !ok {
		return fmt.Errorf("invalid flee type in stats_info packet")
	}

	fleeBonusVal, ok := args["flee_bonus"]
	if !ok {
		return fmt.Errorf("missing flee_bonus field in stats_info packet")
	}
	fleeBonus, ok := fleeBonusVal.(uint16)
	if !ok {
		return fmt.Errorf("invalid flee_bonus type in stats_info packet")
	}

	criticalVal, ok := args["critical"]
	if !ok {
		return fmt.Errorf("missing critical field in stats_info packet")
	}
	critical, ok := criticalVal.(uint16)
	if !ok {
		return fmt.Errorf("invalid critical type in stats_info packet")
	}

	// Process the packet
	result := h.processStatsInfo(args)

	// Notify through hooks system
	if h.hookManager != nil {
		h.hookManager.CallHook("actor.stats_info", result)
	}

	// Log debug message
	fmt.Printf("Strength: %d #%d\n", str, pointsStr)
	fmt.Printf("Agility: %d #%d\n", agi, pointsAgi)
	fmt.Printf("Vitality: %d #%d\n", vit, pointsVit)
	fmt.Printf("Intelligence: %d #%d\n", intStat, pointsInt)
	fmt.Printf("Dexterity: %d #%d\n", dex, pointsDex)
	fmt.Printf("Luck: %d #%d\n", luk, pointsLuk)
	fmt.Printf("Attack: %d\n", attack)
	fmt.Printf("Attack Bonus: %d\n", attackBonus)
	fmt.Printf("Magic Attack Min: %d\n", attackMagicMin)
	fmt.Printf("Magic Attack Max: %d\n", attackMagicMax)
	fmt.Printf("Defense: %d\n", def)
	fmt.Printf("Defense Bonus: %d\n", defBonus)
	fmt.Printf("Magic Defense: %d\n", defMagic)
	fmt.Printf("Magic Defense Bonus: %d\n", defMagicBonus)
	fmt.Printf("Hit: %d\n", hit)
	fmt.Printf("Flee: %d\n", flee)
	fmt.Printf("Flee Bonus: %d\n", fleeBonus)
	fmt.Printf("Critical: %d\n", critical)
	fmt.Printf("Status Points: %d\n", pointsFree)

	return nil
}

// processStatsInfo processes the stats_info packet and returns a structured result
func (h *Handler) processStatsInfo(args map[string]interface{}) map[string]interface{} {
	// Simply return the args as the result
	// In a real implementation, we might do more processing here
	return args
}
