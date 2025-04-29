// Package servers provides server implementations for different Ragnarok Online server types.
package servers

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/lenaxia/goKore/network/config"
	"github.com/lenaxia/goKore/network/connection"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/packets"
	"github.com/lenaxia/goKore/network/protocol"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/security"
)

// BaseServer implements the base functionality for all server types
type BaseServer struct {
	// Server information
	serverInfo ServerInfo
	serverType ServerType
	state      ServerState
	stateMutex sync.RWMutex

	// Configuration
	config        *config.ServerConfig
	networkConfig *config.NetworkConfig
	credentials   ServerCredentials
	options       ServerOptions

	// Connection
	conn         connection.Connection
	connManager  *connection.ConnectionManager
	lastActivity time.Time

	// Protocol
	tokenizer *protocol.Tokenizer
	parser    *core.CoreParser

	// Security
	loginManager     *security.LoginManager
	pinManager       *security.PINManager
	antiCheatManager *security.AntiCheatManager

	// Hooks
	hookManager *hooks.HookManager

	// Context for cancellation
	ctx        context.Context
	cancelFunc context.CancelFunc

	// Packet database
	packetDB *packets.PacketDatabase

	// Character selection
	selectedCharacter int
}

// NewBaseServer creates a new base server
func NewBaseServer(baseConfig *BaseServerConfig) (*BaseServer, error) {
	// Validate the server configuration
	if err := ValidateServerConfig(baseConfig.ServerConfig); err != nil {
		return nil, err
	}

	// Validate the network configuration
	if err := ValidateNetworkConfig(baseConfig.NetworkConfig); err != nil {
		return nil, err
	}

	// Create a new hook manager
	hookManager := hooks.NewHookManager()

	// Create a new tokenizer with empty packet definitions (will be populated later)
	tokenizer := protocol.NewTokenizer(make(map[string]protocol.PacketDef))

	// Create a new parser
	parser := core.NewCoreParser(baseConfig.Type.String(), hookManager)

	// Create security managers
	loginManager := security.NewLoginManager(parser, hookManager)
	pinManager := security.NewPINManager(parser, hookManager)
	antiCheatManager := security.NewAntiCheatManager(parser, hookManager)

	// Create a context with cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())

	// Create the server info
	serverInfo := ServerInfo{
		Name:        baseConfig.Name,
		IP:          baseConfig.IP,
		Port:        baseConfig.Port,
		Type:        baseConfig.Type,
		Maintenance: false,
		Users:       0,
		MaxUsers:    0,
		New:         false,
		PvP:         false,
	}

	// Create the base server
	server := &BaseServer{
		serverInfo:        serverInfo,
		serverType:        baseConfig.Type,
		state:             ServerStateDisconnected,
		config:            baseConfig.ServerConfig,
		networkConfig:     baseConfig.NetworkConfig,
		credentials:       baseConfig.Credentials,
		options:           baseConfig.Options,
		tokenizer:         tokenizer,
		parser:            parser,
		loginManager:      loginManager,
		pinManager:        pinManager,
		antiCheatManager:  antiCheatManager,
		hookManager:       hookManager,
		ctx:               ctx,
		cancelFunc:        cancelFunc,
		packetDB:          packets.NewPacketDatabase(),
		selectedCharacter: -1,
	}

	// Register default packet handlers
	server.registerDefaultHandlers()

	return server, nil
}

// Connect connects to the server
func (s *BaseServer) Connect() error {
	s.stateMutex.Lock()
	defer s.stateMutex.Unlock()

	if s.state != ServerStateDisconnected {
		return ErrServerAlreadyConnected
	}

	// Create a new connection based on the server configuration
	var conn connection.Connection
	var err error

	// Create connection configuration
	connConfig := &connection.ConnectionConfig{
		Host:        s.serverInfo.IP,
		Port:        s.serverInfo.Port,
		Timeout:     s.options.Timeout,
		RecvTimeout: s.options.Timeout,
		SendTimeout: s.options.Timeout,
	}

	// Create the appropriate connection type
	if s.options.UseEncryption {
		// Create a TLS connection
		tlsConn := connection.NewTLSConnection(connConfig)
		conn = tlsConn
	} else {
		// Create a direct connection
		directConn := connection.NewDirectConnection(connConfig)
		conn = directConn
	}

	// Connect to the server
	if err = conn.Connect(); err != nil {
		return fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	// Create a connection manager
	s.connManager = connection.NewConnectionManager(conn)

	// Configure the connection manager with network config settings
	if s.networkConfig != nil {
		if s.networkConfig.Timeouts.Connect > 0 {
			s.connManager.SetConnectTimeout(s.networkConfig.Timeouts.Connect)
		}
		if s.networkConfig.Timeouts.Read > 0 {
			s.connManager.SetReceiveTimeout(s.networkConfig.Timeouts.Read)
		}
		if s.networkConfig.Timeouts.Write > 0 {
			s.connManager.SetSendTimeout(s.networkConfig.Timeouts.Write)
		}
		if s.networkConfig.ReconnectPolicy.InitialInterval > 0 {
			s.connManager.SetReconnectDelay(s.networkConfig.ReconnectPolicy.InitialInterval)
		}
		if s.networkConfig.ReconnectPolicy.MaxAttempts >= 0 {
			s.connManager.SetMaxReconnectAttempts(s.networkConfig.ReconnectPolicy.MaxAttempts)
		}
	}

	// Set the connection
	s.conn = conn

	// Update the state
	s.state = ServerStateConnecting

	// Update the last activity time
	s.lastActivity = time.Now()

	// Call the connection hook
	s.hookManager.CallHook("server/connect", nil)

	return nil
}

// Disconnect disconnects from the server
func (s *BaseServer) Disconnect() error {
	s.stateMutex.Lock()
	defer s.stateMutex.Unlock()

	if s.state == ServerStateDisconnected {
		return nil
	}

	// Cancel the context
	s.cancelFunc()

	// Disconnect the connection
	if s.conn != nil {
		if err := s.conn.Disconnect(); err != nil {
			return err
		}
	}

	// Update the state
	s.state = ServerStateDisconnected

	// Call the disconnection hook
	s.hookManager.CallHook("server/disconnect", nil)

	return nil
}

// IsConnected returns whether the server is connected
func (s *BaseServer) IsConnected() bool {
	s.stateMutex.RLock()
	defer s.stateMutex.RUnlock()

	return s.state != ServerStateDisconnected
}

// GetState returns the current server state
func (s *BaseServer) GetState() ServerState {
	s.stateMutex.RLock()
	defer s.stateMutex.RUnlock()

	return s.state
}

// GetServerInfo returns information about the server
func (s *BaseServer) GetServerInfo() ServerInfo {
	return s.serverInfo
}

// GetServerType returns the server type
func (s *BaseServer) GetServerType() ServerType {
	return s.serverType
}

// SetCredentials sets the credentials for connecting to the server
func (s *BaseServer) SetCredentials(credentials ServerCredentials) {
	s.credentials = credentials
	s.loginManager.SetCredentials(credentials.Username, credentials.Password)
	s.pinManager.SetPIN(credentials.PINCode)
}

// SetOptions sets the options for connecting to the server
func (s *BaseServer) SetOptions(options ServerOptions) {
	s.options = options
}

// Login logs in to the server
func (s *BaseServer) Login() error {
	s.stateMutex.Lock()
	defer s.stateMutex.Unlock()

	if s.state != ServerStateConnecting {
		return errors.New("server not in connecting state")
	}

	// Set the credentials
	s.loginManager.SetCredentials(s.credentials.Username, s.credentials.Password)

	// Update the state
	s.state = ServerStateLoggingIn

	// Send the login packet
	// Since there's no SendLogin method in LoginManager, we need to construct and send the login packet manually
	// This is a placeholder - in a real implementation, you would construct the appropriate login packet
	loginPacketID, exists := s.parser.GetPacketID("login")
	if !exists {
		return fmt.Errorf("login packet ID not found")
	}

	// Construct login packet
	loginArgs := map[string]interface{}{
		"switch":   loginPacketID,
		"username": s.credentials.Username,
		"password": s.credentials.Password,
	}

	// Send the login packet
	if err := s.SendPacket(loginPacketID, loginArgs); err != nil {
		return fmt.Errorf("%w: %v", ErrLoginFailed, err)
	}

	// Update the last activity time
	s.lastActivity = time.Now()

	return nil
}

// SelectCharacter selects a character
func (s *BaseServer) SelectCharacter(slot int) error {
	s.stateMutex.Lock()
	defer s.stateMutex.Unlock()

	if s.state != ServerStateLoggedIn {
		return errors.New("server not in logged in state")
	}

	// Set the selected character
	s.selectedCharacter = slot

	// Send the character selection packet
	// Since there's no SelectCharacter method in LoginManager, we need to construct and send the character selection packet manually
	// This is a placeholder - in a real implementation, you would construct the appropriate character selection packet
	selectCharPacketID, exists := s.parser.GetPacketID("char_select")
	if !exists {
		return fmt.Errorf("character selection packet ID not found")
	}

	// Construct character selection packet
	selectCharArgs := map[string]interface{}{
		"switch": selectCharPacketID,
		"slot":   slot,
	}

	// Send the character selection packet
	if err := s.SendPacket(selectCharPacketID, selectCharArgs); err != nil {
		return err
	}

	// Update the last activity time
	s.lastActivity = time.Now()

	return nil
}

// CreateCharacter creates a new character
func (s *BaseServer) CreateCharacter(name string, slot int, hair, hairColor, job, gender, str, agi, vit, int, dex, luk int) error {
	s.stateMutex.Lock()
	defer s.stateMutex.Unlock()

	if s.state != ServerStateLoggedIn {
		return errors.New("server not in logged in state")
	}

	// Send the character creation packet
	// Since there's no CreateCharacter method in LoginManager, we need to construct and send the character creation packet manually
	// This is a placeholder - in a real implementation, you would construct the appropriate character creation packet
	createCharPacketID, exists := s.parser.GetPacketID("char_create")
	if !exists {
		return fmt.Errorf("character creation packet ID not found")
	}

	// Construct character creation packet
	createCharArgs := map[string]interface{}{
		"switch":    createCharPacketID,
		"name":      name,
		"slot":      slot,
		"hair":      hair,
		"hairColor": hairColor,
		"job":       job,
		"gender":    gender,
		"str":       str,
		"agi":       agi,
		"vit":       vit,
		"int":       int,
		"dex":       dex,
		"luk":       luk,
	}

	// Send the character creation packet
	if err := s.SendPacket(createCharPacketID, createCharArgs); err != nil {
		return err
	}

	// Update the last activity time
	s.lastActivity = time.Now()

	return nil
}

// DeleteCharacter deletes a character
func (s *BaseServer) DeleteCharacter(slot int) error {
	s.stateMutex.Lock()
	defer s.stateMutex.Unlock()

	if s.state != ServerStateLoggedIn {
		return errors.New("server not in logged in state")
	}

	// Send the character deletion packet
	// Since there's no DeleteCharacter method in LoginManager, we need to construct and send the character deletion packet manually
	// This is a placeholder - in a real implementation, you would construct the appropriate character deletion packet
	deleteCharPacketID, exists := s.parser.GetPacketID("char_delete")
	if !exists {
		return fmt.Errorf("character deletion packet ID not found")
	}

	// Construct character deletion packet
	deleteCharArgs := map[string]interface{}{
		"switch": deleteCharPacketID,
		"slot":   slot,
	}

	// Send the character deletion packet
	if err := s.SendPacket(deleteCharPacketID, deleteCharArgs); err != nil {
		return err
	}

	// Update the last activity time
	s.lastActivity = time.Now()

	return nil
}

// SendPacket sends a packet to the server
func (s *BaseServer) SendPacket(packetID string, args map[string]interface{}) error {
	if !s.IsConnected() {
		return ErrServerNotConnected
	}

	// Construct the packet
	// Use the protocol parser to construct the packet
	// Directly use the protocol parser to construct the packet
	// Since CoreParser doesn't have a direct method to construct packets,
	// we'll need to implement this functionality

	// This is a placeholder implementation
	packetData := []byte{0x00, 0x00} // Empty packet
	var err error

	// In a real implementation, you would construct the packet based on the packetID and args
	// For example:
	// 1. Look up the packet format from a packet database
	// 2. Encode the arguments according to the format
	// 3. Return the constructed packet
	if err != nil {
		return err
	}

	// Send the packet
	return s.SendRawPacket(packetData)
}

// SendRawPacket sends a raw packet to the server
func (s *BaseServer) SendRawPacket(data []byte) error {
	if !s.IsConnected() {
		return ErrServerNotConnected
	}

	// Send the packet
	if err := s.conn.Send(data); err != nil {
		return err
	}

	// Update the last activity time
	s.lastActivity = time.Now()

	return nil
}

// ProcessPackets processes incoming packets
func (s *BaseServer) ProcessPackets() error {
	if !s.IsConnected() {
		return ErrServerNotConnected
	}

	// Receive data from the server
	data, err := s.conn.Receive()
	if err != nil {
		return err
	}

	// Add the data to the tokenizer
	s.tokenizer.Add(data)

	// Process all available messages
	for {
		// Get the next message
		msg, msgType, err := s.tokenizer.ReadNext()
		if err != nil {
			if errors.Is(err, protocol.ErrIncompletePacket) {
				break
			}
			return err
		}

		// Process the message
		// Process the message based on its type
		switch msgType {
		case protocol.KnownMessage:
			if err := s.parser.Process(msg); err != nil {
				return err
			}
		case protocol.AccountID:
			// Handle account ID message
			if s.hookManager != nil {
				s.hookManager.CallHook("receive/account_id", map[string]interface{}{
					"accountID": msg,
				})
			}
		case protocol.UnknownMessage:
			// Handle unknown message
			if s.hookManager != nil {
				s.hookManager.CallHook("receive/unknown_packet", map[string]interface{}{
					"packet": msg,
				})
			}
		default:
			return fmt.Errorf("unknown message type: %d", msgType)
		}
	}

	// Update the last activity time
	s.lastActivity = time.Now()

	return nil
}

// RegisterPacketHandler registers a handler for a specific packet
func (s *BaseServer) RegisterPacketHandler(packetID, name, format string, paramNames []string, handler core.PacketHandler) {
	s.parser.RegisterHandler(packetID, name, format, paramNames, handler)
}

// RegisterHook registers a hook for a specific event
func (s *BaseServer) RegisterHook(hookName string, callback hooks.HookCallback, userData interface{}) *hooks.HookHandle {
	return s.hookManager.AddHook(hookName, callback, userData)
}

// GetParser returns the packet parser
func (s *BaseServer) GetParser() *core.CoreParser {
	return s.parser
}

// GetTokenizer returns the message tokenizer
func (s *BaseServer) GetTokenizer() *protocol.Tokenizer {
	return s.tokenizer
}

// GetConfig returns the server configuration
func (s *BaseServer) GetConfig() *config.ServerConfig {
	return s.config
}

// GetNetworkConfig returns the network configuration
func (s *BaseServer) GetNetworkConfig() *config.NetworkConfig {
	return s.networkConfig
}

// GetLoginManager returns the login manager
func (s *BaseServer) GetLoginManager() *security.LoginManager {
	return s.loginManager
}

// GetPINManager returns the PIN manager
func (s *BaseServer) GetPINManager() *security.PINManager {
	return s.pinManager
}

// GetAntiCheatManager returns the anti-cheat manager
func (s *BaseServer) GetAntiCheatManager() *security.AntiCheatManager {
	return s.antiCheatManager
}

// registerDefaultHandlers registers the default packet handlers
func (s *BaseServer) registerDefaultHandlers() {
	// Register login-related handlers
	s.loginManager.RegisterHandlers()

	// Register PIN-related handlers
	s.pinManager.RegisterHandlers()

	// Register anti-cheat-related handlers
	s.antiCheatManager.RegisterHandlers()

	// Register server state change handlers
	s.registerStateChangeHandlers()

	// Register common packet handlers
	s.registerCommonHandlers()
}

// registerStateChangeHandlers registers handlers for state changes
func (s *BaseServer) registerStateChangeHandlers() {
	// Handle login success
	s.RegisterHook("security/login_success", func(hookName string, arg interface{}, userData interface{}) {
		s.stateMutex.Lock()
		s.state = ServerStateLoggedIn
		s.stateMutex.Unlock()
	}, nil)

	// Handle login failure
	s.RegisterHook("security/login_failure", func(hookName string, arg interface{}, userData interface{}) {
		s.stateMutex.Lock()
		s.state = ServerStateConnecting
		s.stateMutex.Unlock()
	}, nil)

	// Handle character selection success
	s.RegisterHook("security/character_selected", func(hookName string, arg interface{}, userData interface{}) {
		s.stateMutex.Lock()
		s.state = ServerStateInGame
		s.stateMutex.Unlock()
	}, nil)

	// Handle disconnection
	s.RegisterHook("connection/disconnected", func(hookName string, arg interface{}, userData interface{}) {
		s.stateMutex.Lock()
		s.state = ServerStateDisconnected
		s.stateMutex.Unlock()
	}, nil)
}

// registerCommonHandlers registers common packet handlers
func (s *BaseServer) registerCommonHandlers() {
	// Register handler for server ping (0x007F)
	s.RegisterPacketHandler("007F", "received_sync", "V", []string{"time"}, func(args map[string]interface{}) error {
		// Update the last activity time
		s.lastActivity = time.Now()

		// Call the ping hook
		s.hookManager.CallHook("server/ping", args)

		return nil
	})

	// Register handler for map loaded (0x0073)
	s.RegisterPacketHandler("0073", "map_loaded", "V a3 C2", []string{"syncMapSync", "coords", "xSize", "ySize"}, func(args map[string]interface{}) error {
		// Update the state
		s.stateMutex.Lock()
		s.state = ServerStateInGame
		s.stateMutex.Unlock()

		// Call the map loaded hook
		s.hookManager.CallHook("server/map_loaded", args)

		return nil
	})

	// Register handler for map change (0x0091)
	s.RegisterPacketHandler("0091", "map_change", "Z16 v2", []string{"map", "x", "y"}, func(args map[string]interface{}) error {
		// Call the map change hook
		s.hookManager.CallHook("server/map_change", args)

		return nil
	})

	// Register handler for map changed (0x0092)
	s.RegisterPacketHandler("0092", "map_changed", "Z16 v2 a4 v", []string{"map", "x", "y", "IP", "port"}, func(args map[string]interface{}) error {
		// Call the map changed hook
		s.hookManager.CallHook("server/map_changed", args)

		return nil
	})

	// Register handler for system chat (0x009A)
	s.RegisterPacketHandler("009A", "system_chat", "v a*", []string{"len", "message"}, func(args map[string]interface{}) error {
		// Call the system chat hook
		s.hookManager.CallHook("server/system_chat", args)

		return nil
	})

	// Register handler for errors (0x0081)
	s.RegisterPacketHandler("0081", "errors", "C", []string{"type"}, func(args map[string]interface{}) error {
		// Call the error hook
		s.hookManager.CallHook("server/error", args)

		return nil
	})
}
