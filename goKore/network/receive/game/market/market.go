package market

import (
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// Market buy result constants
const (
	MarketBuyResultSuccess      = 0
	MarketBuyResultNoZeny       = 1
	MarketBuyResultOverWeight   = 2
	MarketBuyResultOutOfSpace   = 3
	MarketBuyResultAmountTooBig = 9
)

// MarketBuyResultError is defined separately because it doesn't fit in uint8
const MarketBuyResultError = 0xffff // -1

// MarketManager handles market-related packet handling
type MarketManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for market interactions
	storeList    map[string]interface{}
	storeItems   []map[string]interface{}
	npcTalkState map[string]interface{}
	inMarket     bool
}

// NewMarketManager creates a new market manager
func NewMarketManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *MarketManager {
	return &MarketManager{
		parser:       parser,
		hookManager:  hookManager,
		logger:       logger,
		storeList:    make(map[string]interface{}),
		storeItems:   make([]map[string]interface{}, 0),
		npcTalkState: make(map[string]interface{}),
		inMarket:     false,
	}
}

// RegisterHandlers registers all market-related packet handlers
func (mm *MarketManager) RegisterHandlers() {
	// Register NPC market info handler
	mm.parser.RegisterHandlerFunc("09D5", "npc_market_info", "v Z*",
		[]string{"len", "itemList"}, mm.HandleNpcMarketInfo)

	// Register NPC market purchase result handler
	mm.parser.RegisterHandlerFunc("09D7", "npc_market_purchase_result", "v C Z*",
		[]string{"len", "result", "itemList"}, mm.HandleNpcMarketPurchaseResult)
}

// HandleNpcMarketInfo handles the npc_market_info packet (lines 7710-7750)
func (mm *MarketManager) HandleNpcMarketInfo(args map[string]interface{}) error {
	// Extract packet data
	itemList, ok := args["itemList"].([]byte)
	if !ok {
		return fmt.Errorf("invalid itemList in npc_market_info packet")
	}

	// Define packet format (simplified)
	// In a real implementation, this would depend on the server type
	packLen := 11 // Length of the packed data (v C V2 v)

	// Clear store list and talk state
	mm.storeItems = make([]map[string]interface{}, 0)
	mm.storeList = make(map[string]interface{})
	mm.npcTalkState = make(map[string]interface{})

	// Process store items
	for i := 0; i < len(itemList); i += packLen {
		if i+packLen <= len(itemList) {
			// Extract item data (simplified)
			// In a real implementation, this would use proper unpacking
			nameID := uint16(itemList[i]) | uint16(itemList[i+1])<<8
			itemType := uint8(itemList[i+2])
			price := uint32(itemList[i+3]) | uint32(itemList[i+4])<<8 | uint32(itemList[i+5])<<16 | uint32(itemList[i+6])<<24
			amount := uint32(itemList[i+7]) | uint32(itemList[i+8])<<8 | uint32(itemList[i+9])<<16 | uint32(itemList[i+10])<<24
			weight := uint16(0) // Not included in the packet, but needed for the item structure

			// Skip items with amount 0
			if amount == 0 {
				continue
			}

			// Create item entry
			item := map[string]interface{}{
				"nameID": nameID,
				"type":   itemType,
				"price":  price,
				"amount": amount,
				"weight": weight,
				"ID":     len(mm.storeItems),
				"name":   fmt.Sprintf("Item-%d", nameID), // In a real implementation, this would use a function to get the item name
			}

			// Add item to store items
			mm.storeItems = append(mm.storeItems, item)

			mm.logger.Debug("Item added to Market: %s - %dz", item["name"], price)
		}
	}

	// Check if there are any items
	if len(mm.storeItems) == 0 {
		return nil
	}

	// Update NPC talk state
	mm.npcTalkState["talk"] = "store"
	mm.npcTalkState["time"] = getCurrentTime()

	// Set in_market flag
	mm.inMarket = true

	mm.logger.Info("Market information received with %d items", len(mm.storeItems))

	return nil
}

// HandleNpcMarketPurchaseResult handles the npc_market_purchase_result packet (lines 7763-7826)
func (mm *MarketManager) HandleNpcMarketPurchaseResult(args map[string]interface{}) error {
	// Extract packet data
	result, ok := args["result"].(uint8)
	if !ok {
		return fmt.Errorf("invalid result in npc_market_purchase_result packet")
	}

	itemList, ok := args["itemList"].([]byte)
	if !ok {
		return fmt.Errorf("invalid itemList in npc_market_purchase_result packet")
	}

	// Log result
	mm.logger.Debug("Npc market purchase result: %d", result)

	// Process based on result code
	if result == MarketBuyResultSuccess {
		mm.logger.Info("Item bought Successfully")
	} else if result == MarketBuyResultNoZeny {
		mm.logger.Error("Error Market Store (You don't have the necessary zeny)")
	} else if result == MarketBuyResultOverWeight {
		mm.logger.Error("Error Market Store (You are Overweight)")
	} else if result == MarketBuyResultOutOfSpace {
		mm.logger.Error("Error Market Store (You don't have space in inventory)")
	} else if result == MarketBuyResultAmountTooBig {
		mm.logger.Error("Error Market Store (You tried to buy an amount higher than NPC is selling)")
	} else if uint16(result) == MarketBuyResultError {
		mm.logger.Error("Error while trying to buy in a Market Store")
	} else {
		mm.logger.Error("Error while trying to buy in a Market Store (Unknown). (%d)", result)
	}

	// Call hook
	mm.hookManager.CallHook("market_buy_result", map[string]interface{}{
		"result": result,
	})

	// Define packet format (simplified)
	// In a real implementation, this would depend on the server type
	packLen := 11 // Length of the packed data (v C V2 v)

	// Clear store list and talk state
	mm.storeItems = make([]map[string]interface{}, 0)
	mm.storeList = make(map[string]interface{})
	mm.npcTalkState = make(map[string]interface{})

	// Process store items
	for i := 0; i < len(itemList); i += packLen {
		if i+packLen <= len(itemList) {
			// Extract item data (simplified)
			// In a real implementation, this would use proper unpacking
			nameID := uint16(itemList[i]) | uint16(itemList[i+1])<<8
			itemType := uint8(itemList[i+2])
			price := uint32(itemList[i+3]) | uint32(itemList[i+4])<<8 | uint32(itemList[i+5])<<16 | uint32(itemList[i+6])<<24
			amount := uint32(itemList[i+7]) | uint32(itemList[i+8])<<8 | uint32(itemList[i+9])<<16 | uint32(itemList[i+10])<<24
			weight := uint16(0) // Not included in the packet, but needed for the item structure

			// Skip items with amount 0
			if amount == 0 {
				continue
			}

			// Create item entry
			item := map[string]interface{}{
				"nameID": nameID,
				"type":   itemType,
				"price":  price,
				"amount": amount,
				"weight": weight,
				"ID":     len(mm.storeItems),
				"name":   fmt.Sprintf("Item-%d", nameID), // In a real implementation, this would use a function to get the item name
			}

			// Add item to store items
			mm.storeItems = append(mm.storeItems, item)

			mm.logger.Debug("Item added to Market: %s - %dz", item["name"], price)
		}
	}

	// Check if there are any items
	if len(mm.storeItems) == 0 {
		return nil
	}

	// Update NPC talk state
	mm.npcTalkState["talk"] = "store"
	mm.npcTalkState["time"] = getCurrentTime()

	// Set in_market flag
	mm.inMarket = true

	return nil
}

// Helper function to get current time
func getCurrentTime() int64 {
	return 0 // In a real implementation, this would return the current time
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
