// Package connection provides interfaces and implementations for network connections
// to Ragnarok Online servers. It handles connection establishment, data transfer,
// and connection state management.
package connection

import (
	"context"
	"fmt"
	"net"
	"time"
)

// ConnectionState represents the current state of the network connection
type ConnectionState int

const (
	// NOT_CONNECTED indicates that the client is not connected to any server
	NOT_CONNECTED ConnectionState = 1

	// CONNECTED_TO_MASTER_SERVER indicates that the client is connected to the master server
	CONNECTED_TO_MASTER_SERVER ConnectionState = 2

	// CONNECTED_TO_LOGIN_SERVER indicates that the client is connected to the login server
	CONNECTED_TO_LOGIN_SERVER ConnectionState = 3

	// CONNECTED_TO_CHAR_SERVER indicates that the client is connected to the character server
	CONNECTED_TO_CHAR_SERVER ConnectionState = 4

	// IN_GAME indicates that the client is connected to the map server and is ready to play
	IN_GAME ConnectionState = 5

	// IN_GAME_BUT_UNINITIALIZED indicates that the client is in game but without enough information
	IN_GAME_BUT_UNINITIALIZED ConnectionState = -1
)

// String returns the string representation of the ConnectionState
func (s ConnectionState) String() string {
	switch s {
	case NOT_CONNECTED:
		return "NOT_CONNECTED"
	case CONNECTED_TO_MASTER_SERVER:
		return "CONNECTED_TO_MASTER_SERVER"
	case CONNECTED_TO_LOGIN_SERVER:
		return "CONNECTED_TO_LOGIN_SERVER"
	case CONNECTED_TO_CHAR_SERVER:
		return "CONNECTED_TO_CHAR_SERVER"
	case IN_GAME:
		return "IN_GAME"
	case IN_GAME_BUT_UNINITIALIZED:
		return "IN_GAME_BUT_UNINITIALIZED"
	default:
		return fmt.Sprintf("UNKNOWN_STATE(%d)", s)
	}
}

// ConnectionEvent represents events that can occur during a connection's lifecycle
type ConnectionEvent int

const (
	// EventConnecting is triggered when a connection attempt starts
	EventConnecting ConnectionEvent = iota
	// EventConnected is triggered when a connection is successfully established
	EventConnected
	// EventDisconnecting is triggered when a disconnection attempt starts
	EventDisconnecting
	// EventDisconnected is triggered when a connection is fully closed
	EventDisconnected
	// EventDataSent is triggered when data is sent to the server
	EventDataSent
	// EventDataReceived is triggered when data is received from the server
	EventDataReceived
	// EventError is triggered when an error occurs
	EventError
	// EventStateChange is triggered when the connection state changes
	EventStateChange
)

// String returns the string representation of the ConnectionEvent
func (e ConnectionEvent) String() string {
	switch e {
	case EventConnecting:
		return "EventConnecting"
	case EventConnected:
		return "EventConnected"
	case EventDisconnecting:
		return "EventDisconnecting"
	case EventDisconnected:
		return "EventDisconnected"
	case EventDataSent:
		return "EventDataSent"
	case EventDataReceived:
		return "EventDataReceived"
	case EventError:
		return "EventError"
	case EventStateChange:
		return "EventStateChange"
	default:
		return "UnknownEvent"
	}
}

// EventCallback is a function that is called when a connection event occurs
type EventCallback func(event ConnectionEvent, data interface{})

// ConnectionConfig contains configuration options for a connection
type ConnectionConfig struct {
	// Host is the server hostname or IP address
	Host string
	// Port is the server port
	Port int
	// Timeout is the connection timeout
	Timeout time.Duration
	// RecvTimeout is the receive operation timeout
	RecvTimeout time.Duration
	// SendTimeout is the send operation timeout
	SendTimeout time.Duration
	// ServerType is the type of server (login, char, map, etc.)
	ServerType string
	// UseTLS indicates whether to use TLS for the connection
	UseTLS bool
	// TLSVerify indicates whether to verify TLS certificates
	TLSVerify bool
	// ProxyHost is the proxy hostname or IP address (if any)
	ProxyHost string
	// ProxyPort is the proxy port (if any)
	ProxyPort int
	// ProxyType is the type of proxy (SOCKS4, SOCKS5, HTTP)
	ProxyType string
	// ProxyUsername is the username for proxy authentication (if any)
	ProxyUsername string
	// ProxyPassword is the password for proxy authentication (if any)
	ProxyPassword string
}

// Connection extends the base network.Connection interface with additional
// functionality specific to Ragnarok Online connections.
type Connection interface {
	// Connect establishes a connection to the server
	Connect() error

	// Disconnect closes the connection to the server
	Disconnect() error

	// Send sends data to the server
	Send(data []byte) error

	// Receive receives data from the server
	Receive() ([]byte, error)

	// GetState returns the current connection state
	GetState() ConnectionState

	// SetState sets the connection state
	SetState(state ConnectionState)

	// IsConnected returns true if the connection is established
	IsConnected() bool

	// ConnectWithContext establishes a connection to the server with context for timeout/cancellation
	ConnectWithContext(ctx context.Context) error

	// RegisterCallback registers a callback function for a specific event
	RegisterCallback(event ConnectionEvent, callback EventCallback)

	// UnregisterCallback removes a callback function for a specific event
	UnregisterCallback(event ConnectionEvent, callback EventCallback)

	// GetConfig returns the connection configuration
	GetConfig() *ConnectionConfig

	// SetConfig updates the connection configuration
	SetConfig(config *ConnectionConfig)

	// GetRemoteAddress returns the remote server address
	GetRemoteAddress() net.Addr

	// GetLocalAddress returns the local client address
	GetLocalAddress() net.Addr

	// GetLastError returns the last error that occurred
	GetLastError() error

	// GetConnectedTime returns the time when the connection was established
	GetConnectedTime() time.Time

	// GetLastActivityTime returns the time of the last send or receive activity
	GetLastActivityTime() time.Time

	// IsIdle returns true if no activity has occurred for the specified duration
	IsIdle(duration time.Duration) bool

	// SendWithContext sends data to the server with context for timeout/cancellation
	SendWithContext(ctx context.Context, data []byte) error

	// ReceiveWithContext receives data from the server with context for timeout/cancellation
	ReceiveWithContext(ctx context.Context) ([]byte, error)
}

// BaseConnection provides a basic implementation of the Connection interface
// that can be embedded in specific connection implementations.
type BaseConnection struct {
	state            ConnectionState
	config           *ConnectionConfig
	callbacks        map[ConnectionEvent][]EventCallback
	lastError        error
	connectedTime    time.Time
	lastActivityTime time.Time
}

// NewBaseConnection creates a new BaseConnection with the given configuration
func NewBaseConnection(config *ConnectionConfig) *BaseConnection {
	return &BaseConnection{
		state:     NOT_CONNECTED,
		config:    config,
		callbacks: make(map[ConnectionEvent][]EventCallback),
	}
}

// GetState returns the current connection state
func (c *BaseConnection) GetState() ConnectionState {
	return c.state
}

// SetState sets the connection state and triggers the state change event
func (c *BaseConnection) SetState(state ConnectionState) {
	oldState := c.state
	c.state = state
	c.triggerEvent(EventStateChange, map[string]interface{}{
		"oldState": oldState,
		"newState": state,
	})
}

// GetConfig returns the connection configuration
func (c *BaseConnection) GetConfig() *ConnectionConfig {
	return c.config
}

// SetConfig updates the connection configuration
func (c *BaseConnection) SetConfig(config *ConnectionConfig) {
	c.config = config
}

// RegisterCallback registers a callback function for a specific event
func (c *BaseConnection) RegisterCallback(event ConnectionEvent, callback EventCallback) {
	if c.callbacks[event] == nil {
		c.callbacks[event] = make([]EventCallback, 0)
	}
	c.callbacks[event] = append(c.callbacks[event], callback)
}

// UnregisterCallback removes a callback function for a specific event
func (c *BaseConnection) UnregisterCallback(event ConnectionEvent, callback EventCallback) {
	if c.callbacks[event] == nil {
		return
	}

	// Find and remove the callback
	for i, cb := range c.callbacks[event] {
		if &cb == &callback {
			c.callbacks[event] = append(c.callbacks[event][:i], c.callbacks[event][i+1:]...)
			break
		}
	}
}

// triggerEvent calls all registered callbacks for the given event
func (c *BaseConnection) triggerEvent(event ConnectionEvent, data interface{}) {
	if c.callbacks[event] == nil {
		return
	}

	for _, callback := range c.callbacks[event] {
		callback(event, data)
	}
}

// GetLastError returns the last error that occurred
func (c *BaseConnection) GetLastError() error {
	return c.lastError
}

// setLastError sets the last error and triggers the error event
func (c *BaseConnection) setLastError(err error) {
	c.lastError = err
	if err != nil {
		c.triggerEvent(EventError, err)
	}
}

// GetConnectedTime returns the time when the connection was established
func (c *BaseConnection) GetConnectedTime() time.Time {
	return c.connectedTime
}

// GetLastActivityTime returns the time of the last send or receive activity
func (c *BaseConnection) GetLastActivityTime() time.Time {
	return c.lastActivityTime
}

// IsIdle returns true if no activity has occurred for the specified duration
func (c *BaseConnection) IsIdle(duration time.Duration) bool {
	return time.Since(c.lastActivityTime) > duration
}

// updateLastActivityTime updates the last activity time to the current time
func (c *BaseConnection) updateLastActivityTime() {
	c.lastActivityTime = time.Now()
}
