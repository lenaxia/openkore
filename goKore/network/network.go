// Package network provides the core networking functionality for OpenKore.
// It defines connection states, interfaces, and error types used throughout
// the network stack.
package network

import (
	"errors"
)

// Connection states
const (
	// NotConnected represents the state when not connected to any server
	NotConnected = 1

	// ConnectedToMasterServer represents the state when connected to the master server
	ConnectedToMasterServer = 2

	// ConnectedToLoginServer represents the state when connected to the login server
	ConnectedToLoginServer = 3

	// ConnectedToCharServer represents the state when connected to the character server
	ConnectedToCharServer = 4

	// InGame represents the state when connected to the map server and fully functional
	InGame = 5

	// InGameButUninitialized represents the state when in game but without enough information
	// This can only happen in XKore or XKoreProxy mode when the RO client is already
	// logged in before OpenKore, and OpenKore doesn't have enough information
	// (such as character name) to work properly
	InGameButUninitialized = -1
)

// Common errors
var (
	// ErrNotConnected is returned when an operation requires a connection but none exists
	ErrNotConnected = errors.New("not connected to server")

	// ErrInvalidState is returned when an operation is attempted in an invalid connection state
	ErrInvalidState = errors.New("invalid connection state")

	// ErrTimeout is returned when a network operation times out
	ErrTimeout = errors.New("connection timed out")

	// ErrConnectionClosed is returned when attempting to use a closed connection
	ErrConnectionClosed = errors.New("connection closed")

	// ErrPacketTooLarge is returned when a packet exceeds the maximum allowed size
	ErrPacketTooLarge = errors.New("packet too large")

	// ErrInvalidPacket is returned when a packet has an invalid structure
	ErrInvalidPacket = errors.New("invalid packet structure")
)

// NetworkInterface defines the common interface for all network implementations
type NetworkInterface interface {
	// Connect establishes a connection to the server
	Connect() error

	// Disconnect terminates the connection to the server
	Disconnect() error

	// IsConnected checks if there is an active connection
	IsConnected() bool

	// GetState returns the current connection state
	GetState() int

	// SetState changes the connection state
	SetState(state int)

	// Send transmits data to the server
	Send(data []byte) error

	// Receive retrieves data from the server
	Receive() ([]byte, error)
}

// StateChangeCallback is called when connection state changes
type StateChangeCallback func(oldState, newState int)

// PacketHandler defines the interface for packet handling
type PacketHandler interface {
	// Handle processes a packet and returns an error if processing fails
	Handle(packet []byte) error
}

// PacketSender defines the interface for packet sending
type PacketSender interface {
	// Send constructs and sends a packet with the given name and fields
	Send(packetName string, fields map[string]interface{}) ([]byte, error)

	// GetCashShopManager returns a manager for cash shop-related operations
	GetCashShopManager() interface{}

	// GetMiscManager returns a manager for miscellaneous operations
	GetMiscManager() interface{}

	// GetInfoChatManager returns a manager for chat info operations
	GetInfoChatManager() interface{}
}

// NetworkManager manages the network connections and packet handling
type NetworkManager struct {
	// Current connection state
	state int

	// Callback for state changes
	stateChangeCallback StateChangeCallback

	// Network interface
	networkInterface NetworkInterface

	// Packet sender
	packetSender PacketSender

	// Packet handler
	packetHandler PacketHandler
}

// NewNetworkManager creates a new network manager
func NewNetworkManager(networkInterface NetworkInterface, packetSender PacketSender, packetHandler PacketHandler) *NetworkManager {
	// Check for nil parameters
	if networkInterface == nil {
		panic("networkInterface cannot be nil")
	}
	if packetSender == nil {
		panic("packetSender cannot be nil")
	}
	if packetHandler == nil {
		panic("packetHandler cannot be nil")
	}

	return &NetworkManager{
		state:            NotConnected,
		networkInterface: networkInterface,
		packetSender:     packetSender,
		packetHandler:    packetHandler,
	}
}

// SetStateChangeCallback sets the callback for state changes
func (nm *NetworkManager) SetStateChangeCallback(callback StateChangeCallback) {
	nm.stateChangeCallback = callback
}

// Connect establishes a connection to the server
func (nm *NetworkManager) Connect() error {
	return nm.networkInterface.Connect()
}

// Disconnect terminates the connection to the server
func (nm *NetworkManager) Disconnect() error {
	return nm.networkInterface.Disconnect()
}

// IsConnected checks if there is an active connection
func (nm *NetworkManager) IsConnected() bool {
	return nm.networkInterface.IsConnected()
}

// GetState returns the current connection state
func (nm *NetworkManager) GetState() int {
	return nm.state
}

// SetState changes the connection state
func (nm *NetworkManager) SetState(state int) {
	oldState := nm.state
	nm.state = state

	if nm.stateChangeCallback != nil {
		nm.stateChangeCallback(oldState, state)
	}

	nm.networkInterface.SetState(state)
}

// Send constructs and sends a packet with the given name and fields
func (nm *NetworkManager) Send(packetName string, fields map[string]interface{}) ([]byte, error) {
	return nm.packetSender.Send(packetName, fields)
}

// GetCashShopManager returns a manager for cash shop-related operations
func (nm *NetworkManager) GetCashShopManager() interface{} {
	return nm.packetSender.GetCashShopManager()
}

// GetMiscManager returns a manager for miscellaneous operations
func (nm *NetworkManager) GetMiscManager() interface{} {
	return nm.packetSender.GetMiscManager()
}

// GetInfoChatManager returns a manager for chat info operations
func (nm *NetworkManager) GetInfoChatManager() interface{} {
	return nm.packetSender.GetInfoChatManager()
}

// HandlePacket processes a packet and returns an error if processing fails
func (nm *NetworkManager) HandlePacket(packet []byte) error {
	return nm.packetHandler.Handle(packet)
}
