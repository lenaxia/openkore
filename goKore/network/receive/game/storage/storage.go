// Package storage provides handlers for storage-related packets.
package storage

import (
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// StorageManager manages storage-related packet handlers
type StorageManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewStorageManager creates a new storage manager
func NewStorageManager(parser *core.CoreParser, hookManager *hooks.HookManager) *StorageManager {
	return &StorageManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterStorageHandlers registers all handlers related to storage
func (m *StorageManager) RegisterStorageHandlers() {
	// Register storage handlers
	if m.parser != nil {
		// Register guild_storage_log handler
		m.parser.RegisterHandlerFunc("09A6", "guild_storage_log", "v a*",
			[]string{"result", "log"},
			m.HandleGuildStorageLog)
	}
}

// RegisterAllHandlers registers all storage-related handlers
func (m *StorageManager) RegisterAllHandlers() {
	// Register storage handlers
	m.RegisterStorageHandlers()
}

// HandleGuildStorageLog handles the guild_storage_log packet
// Packet format: 09A6 <result>.W <storage log>.B
func (m *StorageManager) HandleGuildStorageLog(args map[string]interface{}) error {
	// Process the packet
	result := m.processGuildStorageLog(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("storage.guild_storage_log", result)
	}

	return nil
}

// processGuildStorageLog processes the guild_storage_log packet and returns a structured result
func (m *StorageManager) processGuildStorageLog(args map[string]interface{}) map[string]interface{} {
	var resultCode byte
	var log []byte
	var status string
	var items []map[string]interface{}

	// Extract result from args
	if resultVal, ok := args["result"].(byte); ok {
		resultCode = resultVal
	} else if resultVal, ok := args["result"].(uint16); ok {
		resultCode = byte(resultVal)
	}

	// Extract log from args
	if logVal, ok := args["log"].([]byte); ok {
		log = logVal
	}

	// Process based on result code
	if resultCode == 0 || resultCode == 1 {
		// Define action strings
		actions := map[byte]string{
			0: "Get",
			1: "Put",
		}

		// Create header for status message
		var message strings.Builder
		message.WriteString(centerString("[ Guild Storage LOG ]", 80, '-') + "\n")
		message.WriteString("#  Name                     Item-Name                                         Amount  Action          Time\n")

		// Process each item in the log
		items = make([]map[string]interface{}, 0)
		for i := 0; i < len(log); i += 83 {
			if i+83 > len(log) {
				break
			}

			item := make(map[string]interface{})

			// Extract ID (4 bytes)
			item["ID"] = uint32(log[i]) |
				(uint32(log[i+1]) << 8) |
				(uint32(log[i+2]) << 16) |
				(uint32(log[i+3]) << 24)

			// Extract nameID (2 bytes)
			item["nameID"] = uint16(log[i+4]) |
				(uint16(log[i+5]) << 8)

			// Extract amount (4 bytes)
			item["amount"] = uint32(log[i+6]) |
				(uint32(log[i+7]) << 8) |
				(uint32(log[i+8]) << 16) |
				(uint32(log[i+9]) << 24)

			// Extract action (1 byte)
			item["action"] = log[i+10]

			// Extract upgrade (4 bytes)
			item["upgrade"] = uint32(log[i+11]) |
				(uint32(log[i+12]) << 8) |
				(uint32(log[i+13]) << 16) |
				(uint32(log[i+14]) << 24)

			// Extract uniqueID (8 bytes)
			item["uniqueID"] = uint64(log[i+15]) |
				(uint64(log[i+16]) << 8) |
				(uint64(log[i+17]) << 16) |
				(uint64(log[i+18]) << 24) |
				(uint64(log[i+19]) << 32) |
				(uint64(log[i+20]) << 40) |
				(uint64(log[i+21]) << 48) |
				(uint64(log[i+22]) << 56)

			// Extract identified (1 byte)
			item["identified"] = log[i+23]

			// Extract type_equip (2 bytes)
			item["type_equip"] = uint16(log[i+24]) |
				(uint16(log[i+25]) << 8)

			// Extract cards (8 bytes)
			cards := make([]uint16, 4)
			for j := 0; j < 4; j++ {
				cards[j] = uint16(log[i+26+(j*2)]) |
					(uint16(log[i+27+(j*2)]) << 8)
			}
			item["cards"] = cards

			// Extract charName (24 bytes)
			charNameEnd := i + 34
			for j := i + 34; j < i+58; j++ {
				if log[j] == 0 {
					charNameEnd = j
					break
				}
			}
			item["charName"] = string(log[i+34 : charNameEnd])

			// Extract time (24 bytes)
			timeEnd := i + 58
			for j := i + 58; j < i+82; j++ {
				if log[j] == 0 {
					timeEnd = j
					break
				}
			}
			item["time"] = string(log[i+58 : timeEnd])

			// Extract attribute (1 byte)
			item["attribute"] = log[i+82]

			// Add to items list
			items = append(items, item)

			// Add to message
			index := len(items) - 1
			actionStr := actions[item["action"].(byte)]
			// In a real implementation, we would use itemName() to get the item name
			// For now, we'll just use the nameID as a placeholder
			itemName := fmt.Sprintf("Item #%d", item["nameID"])

			message.WriteString(fmt.Sprintf("%2d %-24s %-48s %6d %-7s %20s\n",
				index, item["charName"], itemName, item["amount"], actionStr, item["time"]))
		}

		// Add footer to message
		message.WriteString(strings.Repeat("-", 80))

		// Set status message
		status = message.String()

	} else if resultCode == 2 {
		status = "Guild Storage empty."
	} else if resultCode == 3 {
		status = "You are not currently using Guild Storage. Please try later."
	}

	// Create the result
	result := map[string]interface{}{
		"result": resultCode,
		"status": status,
	}

	// Add items if available
	if items != nil {
		result["items"] = items
	}

	return result
}

// centerString centers a string within a field of a given width
func centerString(s string, width int, fill byte) string {
	if len(s) >= width {
		return s
	}

	leftPad := (width - len(s)) / 2
	rightPad := width - len(s) - leftPad

	return strings.Repeat(string(fill), leftPad) + s + strings.Repeat(string(fill), rightPad)
}
