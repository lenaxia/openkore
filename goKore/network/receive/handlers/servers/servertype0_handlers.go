// Package servers provides server-specific handlers for different server types.
package servers

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/base"
)

// RegisterServerType0Handlers registers ServerType0-specific handlers with the receive component.
func RegisterServerType0Handlers(receive *base.BaseReceive) {
	// Register hooks for events
	receive.RegisterHook("account_info_received", func(hookName string, arg interface{}, userData interface{}) {
		// Handle account info received event
	})

	receive.RegisterHook("login_error", func(hookName string, arg interface{}, userData interface{}) {
		// Handle login error event
	})

	receive.RegisterHook("characters_info_received", func(hookName string, arg interface{}, userData interface{}) {
		// Handle characters info received event
	})

	receive.RegisterHook("character_creation_successful", func(hookName string, arg interface{}, userData interface{}) {
		// Handle character creation successful event
	})

	receive.RegisterHook("character_deletion_successful", func(hookName string, arg interface{}, userData interface{}) {
		// Handle character deletion successful event
	})

	receive.RegisterHook("character_deletion_failed", func(hookName string, arg interface{}, userData interface{}) {
		// Handle character deletion failed event
	})

	// Register packet handlers
	receive.RegisterHandler("account_server_info", func(args map[string]interface{}) error {
		return handleAccountServerInfo(args, receive)
	})

	receive.RegisterHandler("login_error", func(args map[string]interface{}) error {
		return handleLoginError(args, receive)
	})

	receive.RegisterHandler("received_characters_info", func(args map[string]interface{}) error {
		return handleReceivedCharactersInfo(args, receive)
	})

	receive.RegisterHandler("character_creation_successful", func(args map[string]interface{}) error {
		return handleCharacterCreationSuccessful(args, receive)
	})

	receive.RegisterHandler("character_deletion_successful", func(args map[string]interface{}) error {
		return handleCharacterDeletionSuccessful(args, receive)
	})

	receive.RegisterHandler("character_deletion_failed", func(args map[string]interface{}) error {
		return handleCharacterDeletionFailed(args, receive)
	})
}

// handleAccountServerInfo handles the account_server_info packet
func handleAccountServerInfo(args map[string]interface{}, receive *base.BaseReceive) error {
	// Extract session information
	sessionID, _ := args["sessionID"].([]byte)
	accountID, _ := args["accountID"].([]byte)

	// Extract account information
	accountSex, _ := args["accountSex"].(int)

	// Process the server info data
	// This would typically involve parsing the serverInfo byte array into a list of servers

	// Trigger the event using the global hook manager
	hooks.CallHook("account_info_received", map[string]interface{}{
		"sessionID": sessionID,
		"accountID": accountID,
		"sex":       accountSex,
	})

	return nil
}

// handleLoginError handles the login_error packet
func handleLoginError(args map[string]interface{}, receive *base.BaseReceive) error {
	// Extract error information
	errorType, _ := args["type"].(int)
	date, _ := args["date"].(string)

	// Trigger the event using the global hook manager
	hooks.CallHook("login_error", map[string]interface{}{
		"type": errorType,
		"date": date,
	})

	return nil
}

// handleReceivedCharactersInfo handles the received_characters_info packet
func handleReceivedCharactersInfo(args map[string]interface{}, receive *base.BaseReceive) error {
	// Extract character information
	totalSlot, _ := args["total_slot"].(int)
	premiumStartSlot, _ := args["premium_start_slot"].(int)
	premiumEndSlot, _ := args["premium_end_slot"].(int)

	// Process the character info data
	// This would typically involve parsing the charInfo byte array into a list of characters

	// Trigger the event using the global hook manager
	hooks.CallHook("characters_info_received", map[string]interface{}{
		"total_slot":         totalSlot,
		"premium_start_slot": premiumStartSlot,
		"premium_end_slot":   premiumEndSlot,
	})

	return nil
}

// handleCharacterCreationSuccessful handles the character_creation_successful packet
func handleCharacterCreationSuccessful(args map[string]interface{}, receive *base.BaseReceive) error {
	// Process the character info data
	// This would typically involve parsing the charInfo byte array into a character object

	// Trigger the event using the global hook manager
	hooks.CallHook("character_creation_successful", map[string]interface{}{})

	return nil
}

// handleCharacterDeletionSuccessful handles the character_deletion_successful packet
func handleCharacterDeletionSuccessful(args map[string]interface{}, receive *base.BaseReceive) error {
	// Trigger the event using the global hook manager
	hooks.CallHook("character_deletion_successful", map[string]interface{}{})

	return nil
}

// handleCharacterDeletionFailed handles the character_deletion_failed packet
func handleCharacterDeletionFailed(args map[string]interface{}, receive *base.BaseReceive) error {
	// Extract error information
	errorCode, _ := args["error"].(int)

	// Trigger the event using the global hook manager
	hooks.CallHook("character_deletion_failed", map[string]interface{}{
		"error": errorCode,
	})

	return nil
}
