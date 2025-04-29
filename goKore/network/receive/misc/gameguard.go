package misc

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// GameGuardConfig holds configuration for GameGuard
type GameGuardConfig struct {
	// GameGuard setting (0 = disabled, 2 = enabled)
	GameGuard int
	// Network version
	NetVersion int
}

// DefaultGameGuardConfig returns a default GameGuard configuration
func DefaultGameGuardConfig() *GameGuardConfig {
	return &GameGuardConfig{
		GameGuard:  0, // Disabled by default
		NetVersion: 1,
	}
}

// GameGuardManager manages GameGuard-related packet handlers
type GameGuardManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	config      *GameGuardConfig
}

// NewGameGuardManager creates a new GameGuard manager
func NewGameGuardManager(parser *core.CoreParser, hookManager *hooks.HookManager, config *GameGuardConfig) *GameGuardManager {
	if config == nil {
		config = DefaultGameGuardConfig()
	}

	return &GameGuardManager{
		parser:      parser,
		hookManager: hookManager,
		config:      config,
	}
}

// RegisterHandlers registers all handlers related to GameGuard
func (m *GameGuardManager) RegisterHandlers() {
	// Register gameguard_request handler
	m.parser.RegisterHandlerFunc("0277", "gameguard_request", "",
		[]string{},
		m.handleGameGuardRequest)

	// Register gameguard_grant handler
	m.parser.RegisterHandlerFunc("02DC", "gameguard_grant", "C",
		[]string{"server"},
		m.handleGameGuardGrant)
}

// handleGameGuardRequest handles the gameguard_request packet
// Packet format: 0277
func (m *GameGuardManager) handleGameGuardRequest(args map[string]interface{}) error {
	// Check if GameGuard is enabled
	if (m.config.NetVersion == 1 && m.config.GameGuard != 2) || (m.config.GameGuard == 0) {
		// GameGuard is disabled, so we don't need to do anything
		return nil
	}

	// Process the packet
	result := m.processGameGuardRequest(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("misc.gameguard_request", result)
	}

	return nil
}

// processGameGuardRequest processes the gameguard_request packet and returns a structured result
func (m *GameGuardManager) processGameGuardRequest(args map[string]interface{}) map[string]interface{} {
	var rawMsg []byte
	var rawMsgSize int

	// Extract RAW_MSG from args
	if rawMsgVal, ok := args["RAW_MSG"].([]byte); ok {
		rawMsg = rawMsgVal
	}

	// Extract RAW_MSG_SIZE from args
	if rawMsgSizeVal, ok := args["RAW_MSG_SIZE"].(int); ok {
		rawMsgSize = rawMsgSizeVal
	}

	// In the original implementation, this would query Poseidon
	// Poseidon::Client::getInstance()->query(substr($args->{RAW_MSG}, 0, $args->{RAW_MSG_SIZE}));
	// We'll need to implement this functionality elsewhere

	// Return structured result
	return map[string]interface{}{
		"raw_msg":      rawMsg[:rawMsgSize],
		"raw_msg_size": rawMsgSize,
		"message":      "Querying Poseidon",
	}
}

// handleGameGuardGrant handles the gameguard_grant packet
// Packet format: 02DC <server>.B
func (m *GameGuardManager) handleGameGuardGrant(args map[string]interface{}) error {
	// Process the packet
	result := m.processGameGuardGrant(args)

	// Notify through hooks system
	if m.hookManager != nil {
		m.hookManager.CallHook("misc.gameguard_grant", result)
	}

	return nil
}

// processGameGuardGrant processes the gameguard_grant packet and returns a structured result
func (m *GameGuardManager) processGameGuardGrant(args map[string]interface{}) map[string]interface{} {
	var server uint8
	var state string
	var message string

	// Extract server from args
	if serverVal, ok := args["server"].(uint8); ok {
		server = serverVal
	}

	// Process based on server value
	switch server {
	case 0:
		state = "denied"
		message = "The server Denied the login because GameGuard packets where not replied " +
			"correctly or too many time has been spent to send the response. " +
			"Please verify the version of your poseidon server and try again"
	case 1:
		state = "account_server"
		message = "Server granted login request to account server"
	default:
		state = "char_map_server"
		message = "Server granted login request to char/map server"

		// In the original implementation, this would call change_to_constate25()
		// if ($masterServer->{'gameGuard'} eq "2")
		// We'll need to implement this functionality elsewhere
	}

	// In the original implementation, this would set the network state
	// $net->setState(1.3) if ($net->getState() == 1.2);
	// We'll need to implement this functionality elsewhere

	// Return structured result
	return map[string]interface{}{
		"server":  server,
		"state":   state,
		"message": message,
	}
}

// UpdateConfig updates the GameGuard configuration
func (m *GameGuardManager) UpdateConfig(config *GameGuardConfig) {
	if config != nil {
		m.config = config
	}
}

// GetConfig returns the current GameGuard configuration
func (m *GameGuardManager) GetConfig() *GameGuardConfig {
	return m.config
}
