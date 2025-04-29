// Package gm provides handlers for GM-related packets.
package gm

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// GMManager manages GM-related packet handlers
type GMManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
}

// NewGMManager creates a new GM manager
func NewGMManager(parser *core.CoreParser, hookManager *hooks.HookManager) *GMManager {
	return &GMManager{
		parser:      parser,
		hookManager: hookManager,
	}
}

// RegisterGMHandlers registers all handlers related to GM commands
func (m *GMManager) RegisterGMHandlers() {
	// Register GM handlers
	if m.parser != nil {
		// Register GM_silence handler
		m.parser.RegisterHandlerFunc("0149", "GM_silence", "C Z24",
			[]string{"flag", "name"},
			m.HandleGMSilence)

		// Register GM_req_acc_name handler
		m.parser.RegisterHandlerFunc("01B3", "GM_req_acc_name", "V Z24",
			[]string{"targetID", "accountName"},
			m.HandleGMReqAccName)
	}
}

// RegisterAllHandlers registers all GM-related handlers
func (m *GMManager) RegisterAllHandlers() {
	// Register GM handlers
	m.RegisterGMHandlers()
}

// HandleGMSilence handles the GM_silence packet
// Packet format: 0149 <flag>.B <name>.24B
func (m *GMManager) HandleGMSilence(args map[string]interface{}) error {
	// Process the packet
	result := m.processGMSilence(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("gm.silence", result)
	}

	return nil
}

// processGMSilence processes the GM_silence packet and returns a structured result
func (m *GMManager) processGMSilence(args map[string]interface{}) map[string]interface{} {
	var flag byte
	var name string
	var status string

	// Extract flag from args
	if flagVal, ok := args["flag"].(byte); ok {
		flag = flagVal
	}

	// Extract name from args
	if nameBytes, ok := args["name"].([]byte); ok {
		// Convert bytes to string and trim null bytes
		name = bytesToString(nameBytes)
	}

	// Create status message based on flag
	if flag != 0 {
		status = fmt.Sprintf("You have been: muted by %s.", name)
	} else {
		status = fmt.Sprintf("You have been: unmuted by %s.", name)
	}

	// Create the result
	result := map[string]interface{}{
		"flag":   flag,
		"name":   name,
		"status": status,
	}

	return result
}

// HandleGMReqAccName handles the GM_req_acc_name packet
// Packet format: 01B3 <ID>.L <account name>.24B
func (m *GMManager) HandleGMReqAccName(args map[string]interface{}) error {
	// Process the packet
	result := m.processGMReqAccName(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("gm.req_acc_name", result)
	}

	return nil
}

// processGMReqAccName processes the GM_req_acc_name packet and returns a structured result
func (m *GMManager) processGMReqAccName(args map[string]interface{}) map[string]interface{} {
	var targetID uint32
	var accountName string

	// Extract targetID from args
	if targetIDVal, ok := args["targetID"].(uint32); ok {
		targetID = targetIDVal
	}

	// Extract accountName from args
	if accountNameVal, ok := args["accountName"].(string); ok {
		accountName = accountNameVal
	} else if accountNameBytes, ok := args["accountName"].([]byte); ok {
		// Convert bytes to string and trim null bytes
		accountName = bytesToString(accountNameBytes)
	}

	// Create status message
	status := fmt.Sprintf("The accountName for ID %d is %s.", targetID, accountName)

	// Create the result
	result := map[string]interface{}{
		"targetID":    targetID,
		"accountName": accountName,
		"status":      status,
	}

	return result
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
