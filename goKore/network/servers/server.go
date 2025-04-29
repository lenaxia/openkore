// Package servers provides server implementations for different Ragnarok Online server types.
package servers

import (
	"errors"
	"fmt"
	"time"

	"github.com/lenaxia/goKore/network/config"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/protocol"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/security"
)

// Errors
var (
	ErrServerNotConnected     = errors.New("server not connected")
	ErrServerAlreadyConnected = errors.New("server already connected")
	ErrInvalidServerType      = errors.New("invalid server type")
	ErrInvalidServerConfig    = errors.New("invalid server configuration")
	ErrConnectionFailed       = errors.New("connection failed")
	ErrLoginFailed            = errors.New("login failed")
	ErrAuthenticationFailed   = errors.New("authentication failed")
	ErrServerMaintenance      = errors.New("server under maintenance")
	ErrServerFull             = errors.New("server full")
)

// ServerType represents the type of server
type ServerType int

// Server types
const (
	ServerTypeUnknown ServerType = iota
	ServerTypeOfficial
	ServerTypeSakray
	ServerTypeRenewal
	ServerTypeClassic
	ServerTypePreRenewal
	ServerTypePrivate
	ServerTypeThor
	ServerTypeZero
	ServerTypeRagexe
	ServerTypeRagexeRE
	ServerTypeCustom
)

// String returns the string representation of the server type
func (t ServerType) String() string {
	switch t {
	case ServerTypeUnknown:
		return "Unknown"
	case ServerTypeOfficial:
		return "Official"
	case ServerTypeSakray:
		return "Sakray"
	case ServerTypeRenewal:
		return "Renewal"
	case ServerTypeClassic:
		return "Classic"
	case ServerTypePreRenewal:
		return "PreRenewal"
	case ServerTypePrivate:
		return "Private"
	case ServerTypeThor:
		return "Thor"
	case ServerTypeZero:
		return "Zero"
	case ServerTypeRagexe:
		return "Ragexe"
	case ServerTypeRagexeRE:
		return "RagexeRE"
	case ServerTypeCustom:
		return "Custom"
	default:
		return "Invalid"
	}
}

// ServerState represents the state of the server connection
type ServerState int

// Server states
const (
	ServerStateDisconnected ServerState = iota
	ServerStateConnecting
	ServerStateHandshaking
	ServerStateLoggingIn
	ServerStateLoggedIn
	ServerStateInGame
	ServerStateDisconnecting
)

// String returns the string representation of the server state
func (s ServerState) String() string {
	switch s {
	case ServerStateDisconnected:
		return "Disconnected"
	case ServerStateConnecting:
		return "Connecting"
	case ServerStateHandshaking:
		return "Handshaking"
	case ServerStateLoggingIn:
		return "LoggingIn"
	case ServerStateLoggedIn:
		return "LoggedIn"
	case ServerStateInGame:
		return "InGame"
	case ServerStateDisconnecting:
		return "Disconnecting"
	default:
		return "Invalid"
	}
}

// ServerInfo contains information about a server
type ServerInfo struct {
	Name        string
	IP          string
	Port        int
	Type        ServerType
	Maintenance bool
	Users       int
	MaxUsers    int
	New         bool
	PvP         bool
}

// ServerCredentials contains credentials for connecting to a server
type ServerCredentials struct {
	Username string
	Password string
	PINCode  string
}

// ServerOptions contains options for connecting to a server
type ServerOptions struct {
	Timeout        time.Duration
	ReconnectDelay time.Duration
	MaxRetries     int
	UseProxy       bool
	ProxyAddress   string
	ProxyPort      int
	ProxyUsername  string
	ProxyPassword  string
	UseCompression bool
	UseEncryption  bool
	UseAntiCheat   bool
	AntiCheatType  security.AntiCheatType
	PacketVersion  int
	ClientVersion  string
	ServerVersion  string
	CharacterSlot  int
	CharacterName  string
	CustomOptions  map[string]interface{}
}

// Server is the interface for a Ragnarok Online server
type Server interface {
	// Connect connects to the server
	Connect() error

	// Disconnect disconnects from the server
	Disconnect() error

	// IsConnected returns whether the server is connected
	IsConnected() bool

	// GetState returns the current server state
	GetState() ServerState

	// GetServerInfo returns information about the server
	GetServerInfo() ServerInfo

	// GetServerType returns the server type
	GetServerType() ServerType

	// SetCredentials sets the credentials for connecting to the server
	SetCredentials(credentials ServerCredentials)

	// SetOptions sets the options for connecting to the server
	SetOptions(options ServerOptions)

	// Login logs in to the server
	Login() error

	// SelectCharacter selects a character
	SelectCharacter(slot int) error

	// CreateCharacter creates a new character
	CreateCharacter(name string, slot int, hair, hairColor, job, gender, str, agi, vit, int, dex, luk int) error

	// DeleteCharacter deletes a character
	DeleteCharacter(slot int) error

	// SendPacket sends a packet to the server
	SendPacket(packetID string, args map[string]interface{}) error

	// SendRawPacket sends a raw packet to the server
	SendRawPacket(data []byte) error

	// ProcessPackets processes incoming packets
	ProcessPackets() error

	// RegisterPacketHandler registers a handler for a specific packet
	RegisterPacketHandler(packetID, name, format string, paramNames []string, handler core.PacketHandler)

	// RegisterHook registers a hook for a specific event
	RegisterHook(hookName string, callback hooks.HookCallback, userData interface{}) *hooks.HookHandle

	// GetParser returns the packet parser
	GetParser() *core.CoreParser

	// GetTokenizer returns the message tokenizer
	GetTokenizer() *protocol.Tokenizer

	// GetConfig returns the server configuration
	GetConfig() *config.ServerConfig

	// GetNetworkConfig returns the network configuration
	GetNetworkConfig() *config.NetworkConfig

	// GetLoginManager returns the login manager
	GetLoginManager() *security.LoginManager

	// GetPINManager returns the PIN manager
	GetPINManager() *security.PINManager

	// GetAntiCheatManager returns the anti-cheat manager
	GetAntiCheatManager() *security.AntiCheatManager
}

// BaseServerConfig contains the base configuration for a server
type BaseServerConfig struct {
	// Server information
	Name        string
	IP          string
	Port        int
	Type        ServerType
	Description string

	// Network configuration
	NetworkConfig *config.NetworkConfig

	// Server configuration
	ServerConfig *config.ServerConfig

	// Credentials
	Credentials ServerCredentials

	// Options
	Options ServerOptions

	// Custom configuration
	CustomConfig map[string]interface{}
}

// NewBaseServerConfig creates a new base server configuration
func NewBaseServerConfig() *BaseServerConfig {
	return &BaseServerConfig{
		NetworkConfig: config.NewNetworkConfigManager().CreateDefaultNetworkConfig("default"),
		ServerConfig:  config.NewServerConfigManager().CreateDefaultServerConfig("default"),
		Options: ServerOptions{
			Timeout:        30 * time.Second,
			ReconnectDelay: 5 * time.Second,
			MaxRetries:     3,
			UseCompression: true,
			UseEncryption:  true,
			UseAntiCheat:   false,
			AntiCheatType:  security.AntiCheatNone,
			CustomOptions:  make(map[string]interface{}),
		},
		CustomConfig: make(map[string]interface{}),
	}
}

// DetectServerType detects the server type from the server configuration
func DetectServerType(serverConfig *config.ServerConfig) ServerType {
	// Check if the server type is explicitly set
	if serverType, ok := serverConfig.CustomFields["serverType"].(string); ok {
		switch serverType {
		case "official":
			return ServerTypeOfficial
		case "sakray":
			return ServerTypeSakray
		case "renewal":
			return ServerTypeRenewal
		case "classic":
			return ServerTypeClassic
		case "prerenewal":
			return ServerTypePreRenewal
		case "private":
			return ServerTypePrivate
		case "thor":
			return ServerTypeThor
		case "zero":
			return ServerTypeZero
		case "ragexe":
			return ServerTypeRagexe
		case "ragexere":
			return ServerTypeRagexeRE
		case "custom":
			return ServerTypeCustom
		}
	}

	// Try to detect the server type from the server version
	if serverVersion, ok := serverConfig.CustomFields["serverVersion"].(string); ok {
		if len(serverVersion) >= 6 {
			// Check for Sakray
			if serverVersion[:6] == "sakray" {
				return ServerTypeSakray
			}

			// Check for Zero
			if serverVersion[:4] == "zero" {
				return ServerTypeZero
			}

			// Check for Thor
			if serverVersion[:4] == "thor" {
				return ServerTypeThor
			}

			// Check for Ragexe
			if serverVersion[:6] == "ragexe" {
				if len(serverVersion) >= 8 && serverVersion[6:8] == "re" {
					return ServerTypeRagexeRE
				}
				return ServerTypeRagexe
			}
		}
	}

	// Try to detect the server type from the packet version
	if packetVersion, ok := serverConfig.CustomFields["packetVersion"].(int); ok {
		if packetVersion >= 20170000 {
			return ServerTypeRenewal
		} else if packetVersion >= 20150000 {
			return ServerTypeRenewal
		} else if packetVersion >= 20120000 {
			return ServerTypeRenewal
		} else if packetVersion >= 20100000 {
			return ServerTypeRenewal
		} else if packetVersion >= 20090000 {
			return ServerTypeRenewal
		} else {
			return ServerTypePreRenewal
		}
	}

	// Default to unknown
	return ServerTypeUnknown
}

// ValidateServerConfig validates a server configuration
func ValidateServerConfig(serverConfig *config.ServerConfig) error {
	// Check if the server IP is set
	if serverConfig.IP == "" {
		return errors.New("server IP is not set")
	}

	// Check if the server port is set
	if serverConfig.Port == 0 {
		return errors.New("server port is not set")
	}

	// Check if the server name is set
	if serverConfig.Name == "" {
		return errors.New("server name is not set")
	}

	return nil
}

// ValidateNetworkConfig validates a network configuration
func ValidateNetworkConfig(networkConfig *config.NetworkConfig) error {
	// Check if the connection timeout is set
	if networkConfig.Timeouts.Connect <= 0 {
		return errors.New("connection timeout must be greater than 0")
	}

	// Check if the reconnect delay is set
	if networkConfig.ReconnectPolicy.InitialInterval <= 0 {
		return errors.New("reconnect delay must be greater than 0")
	}

	// Check if the max retries is set
	if networkConfig.ReconnectPolicy.MaxAttempts < 0 {
		return errors.New("max retries must be greater than or equal to 0")
	}

	return nil
}

// ValidateCredentials validates server credentials
func ValidateCredentials(credentials ServerCredentials) error {
	// Check if the username is set
	if credentials.Username == "" {
		return errors.New("username is not set")
	}

	// Check if the password is set
	if credentials.Password == "" {
		return errors.New("password is not set")
	}

	return nil
}

// CreateServerFromConfig creates a server from a configuration
func CreateServerFromConfig(baseConfig *BaseServerConfig) (Server, error) {
	// Validate the server configuration
	if err := ValidateServerConfig(baseConfig.ServerConfig); err != nil {
		return nil, err
	}

	// Validate the network configuration
	if err := ValidateNetworkConfig(baseConfig.NetworkConfig); err != nil {
		return nil, err
	}

	// Detect the server type if not set
	if baseConfig.Type == ServerTypeUnknown {
		baseConfig.Type = DetectServerType(baseConfig.ServerConfig)
	}

	// Create the server factory
	factory := NewServerFactory()

	// Create the server
	server, err := factory.CreateServer(baseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	return server, nil
}
