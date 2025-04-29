// Package login provides handlers for login-related packets.
package login

import (
	"fmt"

	"github.com/lenaxia/goKore/network/receive/types"
)

// RegisterHandlers registers all login-related handlers with the receive component.
func RegisterHandlers(receive types.Receive) {
	receive.RegisterHandler("account_server_info", handleAccountServerInfo)
	receive.RegisterHandler("login_error", handleLoginError)
	// More login handlers would be registered here
}

// handleAccountServerInfo handles the account_server_info packet.
func handleAccountServerInfo(args map[string]interface{}) error {
	// Extract fields from args
	sessionID, ok := args["sessionID"].([]byte)
	if !ok {
		return fmt.Errorf("invalid sessionID")
	}

	accountID, ok := args["accountID"].([]byte)
	if !ok {
		return fmt.Errorf("invalid accountID")
	}

	// Process the account server info
	fmt.Printf("Received account server info: sessionID=%v, accountID=%v\n", sessionID, accountID)

	// In a real implementation, this would update the game state
	// and trigger appropriate actions

	return nil
}

// handleLoginError handles the login_error packet.
func handleLoginError(args map[string]interface{}) error {
	// Extract fields from args
	errorType, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid error type")
	}

	date, ok := args["date"].(string)
	if !ok {
		return fmt.Errorf("invalid date")
	}

	// Process the login error
	fmt.Printf("Login error: type=%d, date=%s\n", errorType, date)

	// In a real implementation, this would update the game state
	// and trigger appropriate actions based on the error type

	return nil
}
