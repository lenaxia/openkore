// Package card provides card-related functionality for the network stack.
package card

import (
	"fmt"
	"sync"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// CardManager manages card-related functionality
type CardManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	mutex       sync.RWMutex

	// Card merge state
	cardMergeItemsID []int
}

// NewCardManager creates a new card manager
func NewCardManager(parser *core.CoreParser, hookManager *hooks.HookManager) *CardManager {
	return &CardManager{
		parser:           parser,
		hookManager:      hookManager,
		cardMergeItemsID: make([]int, 0),
	}
}

// RegisterHandlers registers card-related packet handlers
func (m *CardManager) RegisterHandlers() {
	// Register handlers for card-related packets
	m.parser.RegisterHandlerFunc("0A10", "card_merge_list", "v a*",
		[]string{"len", "item_list"},
		m.handleCardMergeList)

	m.parser.RegisterHandlerFunc("0A11", "card_merge_status", "B B B",
		[]string{"fail", "item_index", "card_index"},
		m.handleCardMergeStatus)
}

// ClearCardMergeItems clears the card merge items list
func (m *CardManager) ClearCardMergeItems() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.cardMergeItemsID = make([]int, 0)
}

// GetCardMergeItems returns a copy of the card merge items list
func (m *CardManager) GetCardMergeItems() []int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy to avoid race conditions
	result := make([]int, len(m.cardMergeItemsID))
	copy(result, m.cardMergeItemsID)

	return result
}

// handleCardMergeList handles the card_merge_list packet
func (m *CardManager) handleCardMergeList(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Clear previous list
	m.cardMergeItemsID = make([]int, 0)

	// Extract item list
	var itemList []byte
	if itemListVal, ok := args["item_list"].([]byte); ok {
		itemList = itemListVal
	}

	// Process item list
	for i := 0; i < len(itemList); i += 2 {
		if i+1 >= len(itemList) {
			break // Avoid index out of range
		}

		// Extract item ID
		itemID := int(itemList[i]) | (int(itemList[i+1]) << 8)

		// Add to list
		m.cardMergeItemsID = append(m.cardMergeItemsID, itemID)
	}

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("game/card/merge_list", map[string]interface{}{
			"items": m.GetCardMergeItems(),
		})
	}

	return nil
}

// handleCardMergeStatus handles the card_merge_status packet
func (m *CardManager) handleCardMergeStatus(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract parameters
	var fail byte
	var itemIndex byte
	var cardIndex byte

	if failVal, ok := args["fail"].(byte); ok {
		fail = failVal
	}

	if itemIndexVal, ok := args["item_index"].(byte); ok {
		itemIndex = itemIndexVal
	}

	if cardIndexVal, ok := args["card_index"].(byte); ok {
		cardIndex = cardIndexVal
	}

	// Process result
	var message string
	var success bool

	if fail != 0 {
		message = "Card merging failed"
		success = false
	} else {
		message = fmt.Sprintf("Card has been successfully merged into item")
		success = true
	}

	// Clear card merge items list
	m.cardMergeItemsID = make([]int, 0)

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("game/card/merge_status", map[string]interface{}{
			"success":    success,
			"message":    message,
			"item_index": int(itemIndex),
			"card_index": int(cardIndex),
		})
	}

	return nil
}
