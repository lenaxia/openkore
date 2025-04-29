// Package ranking provides handlers for ranking-related packets.
package ranking

import (
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// RankingManager manages ranking-related packet handlers
type RankingManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager

	// Internal state for PVP ranking
	pvpRank    uint16
	pvpNum     uint16
	pvpEnabled bool
}

// NewRankingManager creates a new ranking manager
func NewRankingManager(parser *core.CoreParser, hookManager *hooks.HookManager) *RankingManager {
	return &RankingManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterRankingHandlers registers all handlers related to ranking
func (m *RankingManager) RegisterRankingHandlers() {
	// Register pvp_rank handler
	if m.parser != nil {
		m.parser.RegisterHandlerFunc("09A1", "pvp_rank", "V v2",
			[]string{"ID", "rank", "num"},
			m.handlePvpRank)

		// Register taekwon_rank handler
		m.parser.RegisterHandlerFunc("0224", "taekwon_rank", "V",
			[]string{"rank"},
			m.handleTaekwonRank)

		// Register taekwon_packets handler
		m.parser.RegisterHandlerFunc("0226", "taekwon_packets", "C C Z24",
			[]string{"flag", "value", "name"},
			m.handleTaekwonPackets)

		// Register top10_taekwon_rank handler
		m.parser.RegisterHandlerFunc("0223", "top10_taekwon_rank", "a*",
			[]string{"RAW_MSG"},
			m.handleTop10TaekwonRank)

		// Register top10_pk_rank handler
		m.parser.RegisterHandlerFunc("0238", "top10_pk_rank", "a*",
			[]string{"RAW_MSG"},
			m.handleTop10PkRank)

		// Register top10_blacksmith_rank handler
		m.parser.RegisterHandlerFunc("0219", "top10_blacksmith_rank", "a*",
			[]string{"RAW_MSG"},
			m.handleTop10BlacksmithRank)

		// Register top10_alchemist_rank handler
		m.parser.RegisterHandlerFunc("021A", "top10_alchemist_rank", "a*",
			[]string{"RAW_MSG"},
			m.handleTop10AlchemistRank)

		// Register top10 handler
		m.parser.RegisterHandlerFunc("097D", "top10", "W a*",
			[]string{"type", "RAW_MSG"},
			m.handleTop10)
	}
}

// RegisterAllHandlers registers all ranking-related handlers
func (m *RankingManager) RegisterAllHandlers() {
	// Register ranking handlers
	m.RegisterRankingHandlers()
}

// handlePvpRank handles the pvp_rank packet
// Packet format: 09A1 <ID>.L <rank>.W <num>.W
func (m *RankingManager) handlePvpRank(args map[string]interface{}) error {
	// Process the packet
	result := m.processPvpRank(args)

	// Only notify through hooks system if there's a status message
	if m.hookManager != nil && result["status"] != "" {
		m.hookManager.CallHook("ranking.pvp_rank", result)
	}

	return nil
}

// processPvpRank processes the pvp_rank packet and returns a structured result
func (m *RankingManager) processPvpRank(args map[string]interface{}) map[string]interface{} {
	var id uint32
	var rank, num uint16
	var status string

	// Extract ID from args
	if idVal, ok := args["ID"].(uint32); ok {
		id = idVal
	}

	// Extract rank from args
	if rankVal, ok := args["rank"].(uint16); ok {
		rank = rankVal
	}

	// Extract num from args
	if numVal, ok := args["num"].(uint16); ok {
		num = numVal
	}

	// Check if rank or num has changed
	if rank != m.pvpRank || num != m.pvpNum {
		// Update internal state
		m.pvpRank = rank
		m.pvpNum = num

		// If PVP is enabled, create a status message
		if m.pvpEnabled {
			status = fmt.Sprintf("Your PvP rank is: %d/%d", rank, num)
		}
	}

	// Create the result
	result := map[string]interface{}{
		"ID":     id,
		"rank":   rank,
		"num":    num,
		"status": status,
	}

	return result
}

// handleTaekwonRank handles the taekwon_rank packet
// Packet format: 0224 <rank>.L
func (m *RankingManager) handleTaekwonRank(args map[string]interface{}) error {
	// Process the packet
	result := m.processTaekwonRank(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("ranking.taekwon_rank", result)
	}

	return nil
}

// processTaekwonRank processes the taekwon_rank packet and returns a structured result
func (m *RankingManager) processTaekwonRank(args map[string]interface{}) map[string]interface{} {
	var rank uint32

	// Extract rank from args
	if rankVal, ok := args["rank"].(uint32); ok {
		rank = rankVal
	}

	// Create status message
	status := fmt.Sprintf("TaeKwon Mission Rank : %d", rank)

	// Create the result
	result := map[string]interface{}{
		"rank":   rank,
		"status": status,
	}

	return result
}

// handleTaekwonPackets handles the taekwon_packets packet
// Packet format: 0226 <flag>.B <value>.B <name>.Z24
func (m *RankingManager) handleTaekwonPackets(args map[string]interface{}) error {
	// Process the packet
	result := m.processTaekwonPackets(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("ranking.taekwon_packets", result)
	}

	return nil
}

// processTaekwonPackets processes the taekwon_packets packet and returns a structured result
func (m *RankingManager) processTaekwonPackets(args map[string]interface{}) map[string]interface{} {
	var flag, value byte
	var name string
	var status string

	// Extract flag from args
	if flagVal, ok := args["flag"].(byte); ok {
		flag = flagVal
	}

	// Extract value from args
	if valueVal, ok := args["value"].(byte); ok {
		value = valueVal
	}

	// Extract name from args
	if nameBytes, ok := args["name"].([]byte); ok {
		// Convert bytes to string and trim null bytes
		name = bytesToString(nameBytes)
	}

	// Determine the string based on value
	var valueString string
	switch value {
	case 1:
		valueString = "Sun"
	case 2:
		valueString = "Moon"
	case 3:
		valueString = "Stars"
	default:
		valueString = fmt.Sprintf("Unknown (%d)", value)
	}

	// Process based on flag value
	switch flag {
	case 0: // Info about Star Gladiator save map: Map registered
		status = fmt.Sprintf("You have now marked: %s as Place of the %s.", name, valueString)
	case 1: // Info about Star Gladiator save map: Information
		status = fmt.Sprintf("%s is marked as Place of the %s.", name, valueString)
	case 10: // Info about Star Gladiator hate mob: Register mob
		status = fmt.Sprintf("You have now marked %s as Target of the %s.", name, valueString)
	case 11: // Info about Star Gladiator hate mob: Information
		status = fmt.Sprintf("%s is marked as Target of the %s.", name, valueString)
	case 20: // Info about TaeKwon Do TK_MISSION mob
		status = fmt.Sprintf("[TaeKwon Mission] Target Monster : %s (%d%%)", name, value)
	case 30: // Feel/Hate reset
		status = "Your Hate and Feel targets have been resetted."
	default:
		status = fmt.Sprintf("Unknown results in taekwon_packets (flag: %d)", flag)
	}

	// Create the result
	result := map[string]interface{}{
		"flag":   flag,
		"value":  value,
		"name":   name,
		"status": status,
	}

	return result
}

// handleTop10TaekwonRank handles the top10_taekwon_rank packet
// Packet format: 0223 { <name>.24B }*10 { <point>.L }*10
func (m *RankingManager) handleTop10TaekwonRank(args map[string]interface{}) error {
	// Process the packet
	result := m.processTop10TaekwonRank(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("ranking.top10_taekwon_rank", result)
	}

	return nil
}

// processTop10TaekwonRank processes the top10_taekwon_rank packet and returns a structured result
func (m *RankingManager) processTop10TaekwonRank(args map[string]interface{}) map[string]interface{} {
	var rawMsg []byte
	var rankings []map[string]interface{}

	// Extract raw message from args
	if rawMsgVal, ok := args["RAW_MSG"].([]byte); ok {
		rawMsg = rawMsgVal
	}

	// Extract names and points
	names := make([]string, 10)
	points := make([]uint32, 10)

	// Extract 10 names (24 bytes each)
	for i := 0; i < 10; i++ {
		if 2+(24*i)+24 <= len(rawMsg) {
			nameBytes := rawMsg[2+(24*i) : 2+(24*i)+24]
			names[i] = bytesToString(nameBytes)
		}
	}

	// Extract 10 points (4 bytes each)
	for i := 0; i < 10; i++ {
		if 2+(24*10)+(i*4)+4 <= len(rawMsg) {
			pointBytes := rawMsg[2+(24*10)+(i*4) : 2+(24*10)+(i*4)+4]
			points[i] = uint32(pointBytes[0]) |
				(uint32(pointBytes[1]) << 8) |
				(uint32(pointBytes[2]) << 16) |
				(uint32(pointBytes[3]) << 24)
		}
	}

	// Create rankings list
	rankings = make([]map[string]interface{}, 10)
	for i := 0; i < 10; i++ {
		rankings[i] = map[string]interface{}{
			"rank":   i + 1,
			"name":   names[i],
			"points": points[i],
		}
	}

	// Create formatted text list
	var textList strings.Builder
	textList.WriteString("=============== TAEKWON RANK ================\n")
	textList.WriteString("#    Name                             Points\n")

	for i := 0; i < 10; i++ {
		line := fmt.Sprintf("[%d] %-30s %d\n", i+1, names[i], points[i])
		textList.WriteString(line)
	}

	textList.WriteString("=============================================")

	// Create the result
	result := map[string]interface{}{
		"rankings": rankings,
		"status":   textList.String(),
	}

	return result
}

// handleTop10PkRank handles the top10_pk_rank packet
// Packet format: 0238 { <name>.24B }*10 { <point>.L }*10
func (m *RankingManager) handleTop10PkRank(args map[string]interface{}) error {
	// Process the packet
	result := m.processTop10PkRank(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("ranking.top10_pk_rank", result)
	}

	return nil
}

// processTop10PkRank processes the top10_pk_rank packet and returns a structured result
func (m *RankingManager) processTop10PkRank(args map[string]interface{}) map[string]interface{} {
	var rawMsg []byte
	var rankings []map[string]interface{}

	// Extract raw message from args
	if rawMsgVal, ok := args["RAW_MSG"].([]byte); ok {
		rawMsg = rawMsgVal
	}

	// Extract names and points
	names := make([]string, 10)
	points := make([]uint32, 10)

	// Extract 10 names (24 bytes each)
	for i := 0; i < 10; i++ {
		if 2+(24*i)+24 <= len(rawMsg) {
			nameBytes := rawMsg[2+(24*i) : 2+(24*i)+24]
			names[i] = bytesToString(nameBytes)
		}
	}

	// Extract 10 points (4 bytes each)
	for i := 0; i < 10; i++ {
		if 2+(24*10)+(i*4)+4 <= len(rawMsg) {
			pointBytes := rawMsg[2+(24*10)+(i*4) : 2+(24*10)+(i*4)+4]
			points[i] = uint32(pointBytes[0]) |
				(uint32(pointBytes[1]) << 8) |
				(uint32(pointBytes[2]) << 16) |
				(uint32(pointBytes[3]) << 24)
		}
	}

	// Create rankings list
	rankings = make([]map[string]interface{}, 10)
	for i := 0; i < 10; i++ {
		rankings[i] = map[string]interface{}{
			"rank":   i + 1,
			"name":   names[i],
			"points": points[i],
		}
	}

	// Create formatted text list
	var textList strings.Builder
	textList.WriteString("================ PVP RANK ===================\n")
	textList.WriteString("#    Name                             Points\n")

	for i := 0; i < 10; i++ {
		line := fmt.Sprintf("[%d] %-30s %d\n", i+1, names[i], points[i])
		textList.WriteString(line)
	}

	textList.WriteString("=============================================")

	// Create the result
	result := map[string]interface{}{
		"rankings": rankings,
		"status":   textList.String(),
	}

	return result
}

// handleTop10BlacksmithRank handles the top10_blacksmith_rank packet
// Packet format: 0219 { <name>.24B }*10 { <point>.L }*10
func (m *RankingManager) handleTop10BlacksmithRank(args map[string]interface{}) error {
	// Process the packet
	result := m.processTop10BlacksmithRank(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("ranking.top10_blacksmith_rank", result)
	}

	return nil
}

// processTop10BlacksmithRank processes the top10_blacksmith_rank packet and returns a structured result
func (m *RankingManager) processTop10BlacksmithRank(args map[string]interface{}) map[string]interface{} {
	var rawMsg []byte
	var rankings []map[string]interface{}

	// Extract raw message from args
	if rawMsgVal, ok := args["RAW_MSG"].([]byte); ok {
		rawMsg = rawMsgVal
	}

	// Extract names and points
	names := make([]string, 10)
	points := make([]uint32, 10)

	// Extract 10 names (24 bytes each)
	for i := 0; i < 10; i++ {
		if 2+(24*i)+24 <= len(rawMsg) {
			nameBytes := rawMsg[2+(24*i) : 2+(24*i)+24]
			names[i] = bytesToString(nameBytes)
		}
	}

	// Extract 10 points (4 bytes each)
	for i := 0; i < 10; i++ {
		if 2+(24*10)+(i*4)+4 <= len(rawMsg) {
			pointBytes := rawMsg[2+(24*10)+(i*4) : 2+(24*10)+(i*4)+4]
			points[i] = uint32(pointBytes[0]) |
				(uint32(pointBytes[1]) << 8) |
				(uint32(pointBytes[2]) << 16) |
				(uint32(pointBytes[3]) << 24)
		}
	}

	// Create rankings list
	rankings = make([]map[string]interface{}, 10)
	for i := 0; i < 10; i++ {
		rankings[i] = map[string]interface{}{
			"rank":   i + 1,
			"name":   names[i],
			"points": points[i],
		}
	}

	// Create formatted text list
	var textList strings.Builder
	textList.WriteString("============= BLACKSMITH RANK ===============\n")
	textList.WriteString("#    Name                             Points\n")

	for i := 0; i < 10; i++ {
		line := fmt.Sprintf("[%d] %-30s %d\n", i+1, names[i], points[i])
		textList.WriteString(line)
	}

	textList.WriteString("=============================================")

	// Create the result
	result := map[string]interface{}{
		"rankings": rankings,
		"status":   textList.String(),
	}

	return result
}

// handleTop10AlchemistRank handles the top10_alchemist_rank packet
// Packet format: 021A { <name>.24B }*10 { <point>.L }*10
func (m *RankingManager) handleTop10AlchemistRank(args map[string]interface{}) error {
	// Process the packet
	result := m.processTop10AlchemistRank(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("ranking.top10_alchemist_rank", result)
	}

	return nil
}

// processTop10AlchemistRank processes the top10_alchemist_rank packet and returns a structured result
func (m *RankingManager) processTop10AlchemistRank(args map[string]interface{}) map[string]interface{} {
	var rawMsg []byte
	var rankings []map[string]interface{}

	// Extract raw message from args
	if rawMsgVal, ok := args["RAW_MSG"].([]byte); ok {
		rawMsg = rawMsgVal
	}

	// Extract names and points
	names := make([]string, 10)
	points := make([]uint32, 10)

	// Extract 10 names (24 bytes each)
	for i := 0; i < 10; i++ {
		if 2+(24*i)+24 <= len(rawMsg) {
			nameBytes := rawMsg[2+(24*i) : 2+(24*i)+24]
			names[i] = bytesToString(nameBytes)
		}
	}

	// Extract 10 points (4 bytes each)
	for i := 0; i < 10; i++ {
		if 2+(24*10)+(i*4)+4 <= len(rawMsg) {
			pointBytes := rawMsg[2+(24*10)+(i*4) : 2+(24*10)+(i*4)+4]
			points[i] = uint32(pointBytes[0]) |
				(uint32(pointBytes[1]) << 8) |
				(uint32(pointBytes[2]) << 16) |
				(uint32(pointBytes[3]) << 24)
		}
	}

	// Create rankings list
	rankings = make([]map[string]interface{}, 10)
	for i := 0; i < 10; i++ {
		rankings[i] = map[string]interface{}{
			"rank":   i + 1,
			"name":   names[i],
			"points": points[i],
		}
	}

	// Create formatted text list
	var textList strings.Builder
	textList.WriteString("============= ALCHEMIST RANK ================\n")
	textList.WriteString("#    Name                             Points\n")

	for i := 0; i < 10; i++ {
		line := fmt.Sprintf("[%d] %-30s %d\n", i+1, names[i], points[i])
		textList.WriteString(line)
	}

	textList.WriteString("=============================================")

	// Create the result
	result := map[string]interface{}{
		"rankings": rankings,
		"status":   textList.String(),
	}

	return result
}

// handleTop10 handles the top10 packet
// Packet format: 097D <RankingType>.W {<CharName>.24B <point>L}*10 <mypoint>L
func (m *RankingManager) handleTop10(args map[string]interface{}) error {
	// Process the packet
	var rankType byte
	var rawMsg []byte

	// Extract type from args
	if typeVal, ok := args["type"].(byte); ok {
		rankType = typeVal
	} else if typeVal, ok := args["type"].(uint16); ok {
		rankType = byte(typeVal)
	}

	// Extract raw message from args
	if rawMsgVal, ok := args["RAW_MSG"].([]byte); ok {
		rawMsg = rawMsgVal
	}

	// Create a new args map with just the RAW_MSG
	// No need to skip the first 2 bytes as they are already handled by the parser
	newArgs := map[string]interface{}{
		"RAW_MSG": rawMsg,
	}

	// Dispatch to the appropriate handler based on type
	switch rankType {
	case 0: // Blacksmith
		return m.handleTop10BlacksmithRank(newArgs)
	case 1: // Alchemist
		return m.handleTop10AlchemistRank(newArgs)
	case 2: // Taekwon
		return m.handleTop10TaekwonRank(newArgs)
	case 3: // PK
		return m.handleTop10PkRank(newArgs)
	default:
		// Unknown type
		if m.hookManager != nil {
			m.hookManager.CallHook("ranking.top10_unknown", map[string]interface{}{
				"type":    rankType,
				"status":  fmt.Sprintf("Unknown top10 type %d.", rankType),
				"RAW_MSG": rawMsg,
			})
		}
	}

	return nil
}

// bytesToString converts a byte slice to a string, trimming null bytes
func bytesToString(b []byte) string {
	n := 0
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			break
		}
		n++
	}
	return string(b[:n])
}
