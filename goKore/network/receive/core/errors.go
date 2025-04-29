package core

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
)

// Error type constants
const (
	ErrorServerShutdown       = 0
	ErrorServerClosed         = 1
	ErrorDualLogin            = 2
	ErrorOutOfSync            = 3
	ErrorServerJammed         = 4
	ErrorUnderaged            = 5
	ErrorMustPay              = 6
	ErrorLastConnection       = 8
	ErrorIPCapacityFull       = 9
	ErrorOutOfTime            = 10
	ErrorForcedDisconnectByGM = 15
	ErrorAccountSuspended     = 101
	ErrorTooManyConnections   = 102
)

// ErrorsManager manages error-related functionality
type ErrorsManager struct {
	hookManager *hooks.HookManager
	netState    int
	config      map[string]interface{}
	appName     string
}

// NewErrorsManager creates a new errors manager
func NewErrorsManager(hookManager *hooks.HookManager, netState int, config map[string]interface{}, appName string) *ErrorsManager {
	return &ErrorsManager{
		hookManager: hookManager,
		netState:    netState,
		config:      config,
		appName:     appName,
	}
}

// HandleErrors handles the errors packet
// Packet format: 0081 <type>.B
func (m *ErrorsManager) HandleErrors(args map[string]interface{}) error {
	// Extract error type with safety check
	var errorType byte
	if errorTypeVal, ok := args["type"].(byte); ok {
		errorType = errorTypeVal
	}

	// Call disconnected hook if in game
	if m.netState == NetworkStateInGame && m.hookManager != nil {
		m.hookManager.CallHook("disconnected", nil)
	}

	// Check if we should auto disconnect
	autoDisconnect := false
	if m.netState == NetworkStateInGame {
		dcOnDisconnect, _ := m.getConfigInt("dcOnDisconnect")
		if dcOnDisconnect > 1 || (dcOnDisconnect > 0 && errorType != ErrorOutOfSync && errorType != ErrorOutOfTime) {
			m.logError("Auto disconnecting on Disconnect!")
			m.logChat("*** You disconnected, auto disconnect! ***")
			autoDisconnect = true
		}
	}

	// Set network state to disconnected
	m.setNetState(NetworkStateDisconnected)

	// Handle specific error types
	switch errorType {
	case ErrorServerShutdown:
		dcOnServerShutDown, _ := m.getConfigInt("dcOnServerShutDown")
		if dcOnServerShutDown == 1 {
			m.logError("Auto disconnecting on ServerShutDown!")
			m.logChat("*** Server shutting down, auto disconnect! ***")
			autoDisconnect = true
		} else {
			m.logError("Server shutting down")
		}

	case ErrorServerClosed:
		dcOnServerClose, _ := m.getConfigInt("dcOnServerClose")
		if dcOnServerClose == 1 {
			m.logError("Auto disconnecting on ServerClose!")
			m.logChat("*** Server is closed, auto disconnect! ***")
			autoDisconnect = true
		} else {
			m.logError("Error: Server is closed")
		}

	case ErrorDualLogin:
		dcOnDualLogin, _ := m.getConfigInt("dcOnDualLogin")
		if dcOnDualLogin == 1 {
			m.logError(fmt.Sprintf("Critical Error: Dual login prohibited - Someone trying to login!\n\n%s will now immediately disconnect.", m.appName))
			m.logChat("*** DualLogin, auto disconnect! ***")
			autoDisconnect = true
		} else if dcOnDualLogin >= 2 {
			m.logError("Critical Error: Dual login prohibited - Someone trying to login!")
			m.logMessage(fmt.Sprintf("Reconnecting, wait %d seconds...", dcOnDualLogin))
			// Set reconnect timeout
			// In a real implementation, this would set a timeout for reconnection
		} else {
			m.logError("Critical Error: Dual login prohibited - Someone trying to login!")
		}

	case ErrorOutOfSync:
		m.logError("Error: Out of sync with server")

	case ErrorServerJammed:
		m.logError("Error: Server is jammed due to over-population.")

	case ErrorUnderaged:
		m.logError("Error: You are underaged and cannot join this server.")

	case ErrorMustPay:
		m.logError("Critical Error: You must pay to play this account!")
		if m.netState != 1 { // Not sure what this check is for in the original code
			autoDisconnect = true
		}

	case ErrorLastConnection:
		m.logError("Error: The server still recognizes your last connection")

	case ErrorIPCapacityFull:
		m.logError("Error: IP capacity of this Internet Cafe is full. Would you like to pay the personal base?")

	case ErrorOutOfTime:
		m.logError("Error: You are out of available time paid for")

	case ErrorForcedDisconnectByGM:
		m.logError("Error: You have been forced to disconnect by a GM")

	case ErrorAccountSuspended:
		m.logError("Error: Your account has been suspended until the next maintenance period for possible use of 3rd party programs")

	case ErrorTooManyConnections:
		m.logError("Error: For an hour, more than 10 connections having same IP address, have made. Please check this matter.")

	default:
		m.logError(fmt.Sprintf("Unknown error %d", errorType))
	}

	// Disconnect from server if not type 0
	if errorType != 0 {
		m.serverDisconnect()
	}

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("core.errors", map[string]interface{}{
			"type":           errorType,
			"autoDisconnect": autoDisconnect,
		})
	}

	return nil
}

// Helper methods (these would be implemented in a real system)
func (m *ErrorsManager) getConfigInt(key string) (int, bool) {
	if val, ok := m.config[key]; ok {
		if intVal, ok := val.(int); ok {
			return intVal, true
		}
	}
	return 0, false
}

func (m *ErrorsManager) setNetState(state int) {
	m.netState = state
}

func (m *ErrorsManager) serverDisconnect() {
	// In a real implementation, this would disconnect from the server
}

func (m *ErrorsManager) logError(message string) {
	// In a real implementation, this would log an error message
	// logger.Error(message)
}

func (m *ErrorsManager) logChat(message string) {
	// In a real implementation, this would log a chat message
	// chatLogger.Log(message)
}

func (m *ErrorsManager) logMessage(message string) {
	// In a real implementation, this would log a message
	// logger.Info(message)
}

// Network state constants
const (
	NetworkStateDisconnected = 1
	NetworkStateInGame       = 4
)

// RegisterHandlers registers error-related packet handlers with the given parser
func (m *ErrorsManager) RegisterHandlers(parser *CoreParser) {
	parser.RegisterHandlerFunc("0081", "errors", "B",
		[]string{"type"},
		m.HandleErrors)
}
