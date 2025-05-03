// Package registry provides adapters to integrate send and receive registries with the network stack.
package registry

import (
	"github.com/lenaxia/goKore/network"
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive"
	"github.com/lenaxia/goKore/network/receive/core"
	receivetypes "github.com/lenaxia/goKore/network/receive/types"
	"github.com/lenaxia/goKore/network/send"
	sendcore "github.com/lenaxia/goKore/network/send/core"
)

// Logger defines a common interface that satisfies both sendcore.Logger and core.Logger
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
	Success(format string, args ...interface{})
}

// SendRegistryAdapter adapts the send.HandlerRegistry to implement the network.PacketSender interface
type SendRegistryAdapter struct {
	registry    *send.HandlerRegistry
	baseSend    *sendcore.BaseSend
	hookManager *hooks.HookManager
	logger      sendcore.Logger
}

// NewSendRegistryAdapter creates a new adapter for the send registry
func NewSendRegistryAdapter(hookManager *hooks.HookManager, logger sendcore.Logger) *SendRegistryAdapter {
	baseSend := sendcore.NewBaseSend(hookManager)
	registry := send.NewHandlerRegistry(baseSend, hookManager, logger)

	return &SendRegistryAdapter{
		registry:    registry,
		baseSend:    baseSend,
		hookManager: hookManager,
		logger:      logger,
	}
}

// Send constructs and sends a packet with the given name and fields
func (a *SendRegistryAdapter) Send(packetName string, fields map[string]interface{}) ([]byte, error) {
	// SendPacket only returns error, so we need to construct the packet first to return the bytes
	packet, err := a.baseSend.ConstructPacket(packetName, fields)
	if err != nil {
		return nil, err
	}

	// Send the packet
	err = a.baseSend.SendPacket(packetName, fields)
	return packet, err
}

// managers stores references to all available managers
var managerTypes = map[string]string{
	"CashShop": "cashShopManager",
	"Misc":     "miscManager",
	"InfoChat": "infoChatManager",
}

// GetCashShopManager returns a manager for cash shop-related operations
func (a *SendRegistryAdapter) GetCashShopManager() interface{} {
	return a.getManager("CashShop")
}

// GetMiscManager returns a manager for miscellaneous operations
func (a *SendRegistryAdapter) GetMiscManager() interface{} {
	return a.getManager("Misc")
}

// GetInfoChatManager returns a manager for chat info operations
func (a *SendRegistryAdapter) GetInfoChatManager() interface{} {
	return a.getManager("InfoChat")
}

// getManager is a helper method to get a manager by name
func (a *SendRegistryAdapter) getManager(managerType string) interface{} {
	// This would use reflection to get the manager from the registry
	// For now, we'll return nil as a placeholder
	return nil
}

// Initialize initializes the send registry
func (a *SendRegistryAdapter) Initialize() {
	a.registry.RegisterAllHandlers()
}

// SetConnection sets the connection to use for sending packets
func (a *SendRegistryAdapter) SetConnection(conn interface{}) {
	a.baseSend.SetConnection(conn)
}

// ReceiveRegistryAdapter adapts the receive.HandlerRegistry to implement the network.PacketHandler interface
type ReceiveRegistryAdapter struct {
	registry    *receive.HandlerRegistry
	coreParser  *core.CoreParser
	baseReceive *core.BaseReceive
	hookManager *hooks.HookManager
	logger      core.Logger
}

// NewReceiveRegistryAdapter creates a new adapter for the receive registry
func NewReceiveRegistryAdapter(hookManager *hooks.HookManager, logger core.Logger) *ReceiveRegistryAdapter {
	coreParser := core.NewCoreParser("", hookManager)
	baseReceive := core.NewBaseReceive(hookManager)

	// Create a receive interface implementation
	receiveImpl := &receiveImpl{
		baseReceive: baseReceive,
		logger:      logger,
	}

	registry := receive.NewHandlerRegistry(coreParser, hookManager, receiveImpl, logger)

	return &ReceiveRegistryAdapter{
		registry:    registry,
		coreParser:  coreParser,
		baseReceive: baseReceive,
		hookManager: hookManager,
		logger:      logger,
	}
}

// Handle processes a packet and returns an error if processing fails
func (a *ReceiveRegistryAdapter) Handle(packet []byte) error {
	return a.baseReceive.Process(packet)
}

// Initialize initializes the receive registry
func (a *ReceiveRegistryAdapter) Initialize() {
	a.registry.RegisterAllHandlers()
}

// receiveImpl implements the receivetypes.Receive interface
type receiveImpl struct {
	baseReceive *core.BaseReceive
	logger      core.Logger
}

// Implement the necessary methods from types.Receive interface
// This is a minimal implementation to satisfy the interface requirements
// In a real implementation, these methods would delegate to the appropriate components

func (r *receiveImpl) GetServerType() string {
	return r.baseReceive.GetServerType()
}

// Configure configures the receive component with server-specific packet definitions
func (r *receiveImpl) Configure(serverType string, packetDefs map[string]common.PacketConstruction) error {
	return r.baseReceive.Configure(serverType, packetDefs)
}

// RegisterHandler registers a handler for a specific packet
func (r *receiveImpl) RegisterHandler(packetName string, handler receivetypes.ReceiveHandler) {
	// Convert receivetypes.ReceiveHandler to core.PacketHandler
	// Since they have the same function signature, we can use a simple wrapper
	coreHandler := func(args map[string]interface{}) error {
		return handler(args)
	}
	r.baseReceive.RegisterHandler(packetName, coreHandler)
}

// Process processes a packet, calling the appropriate handler and hooks
func (r *receiveImpl) Process(packet []byte) error {
	return r.baseReceive.Process(packet)
}

// GetPacketID returns the packet ID for a given packet name
func (r *receiveImpl) GetPacketID(name string) (string, bool) {
	// This would typically delegate to the core parser
	// For now, return empty values as a placeholder
	return "", false
}

// RegisterHook registers a hook for a specific event
func (r *receiveImpl) RegisterHook(hookName string, callback hooks.HookCallback) {
	// This would typically delegate to the hook manager
	// For now, it's a no-op
}

// SetDebugMode sets the debug mode
func (r *receiveImpl) SetDebugMode(debug bool) {
	// This would typically set debug mode on the core parser
	// For now, it's a no-op
}

// ParsePacket parses a packet and returns the parsed arguments
func (r *receiveImpl) ParsePacket(packet []byte) (map[string]interface{}, error) {
	// This would typically delegate to the core parser
	// For now, return empty values as a placeholder
	return nil, nil
}

// NetworkRegistryIntegrator integrates the send and receive registries with the network stack
type NetworkRegistryIntegrator struct {
	sendAdapter    *SendRegistryAdapter
	receiveAdapter *ReceiveRegistryAdapter
	hookManager    *hooks.HookManager
	logger         Logger
}

// NewNetworkRegistryIntegrator creates a new integrator for the network registries
func NewNetworkRegistryIntegrator(logger Logger) *NetworkRegistryIntegrator {
	hookManager := hooks.NewHookManager()

	// Create adapters for the specific logger types needed
	sendLogger := &loggerAdapter{logger: logger}
	receiveLogger := &loggerAdapter{logger: logger}

	return &NetworkRegistryIntegrator{
		sendAdapter:    NewSendRegistryAdapter(hookManager, sendLogger),
		receiveAdapter: NewReceiveRegistryAdapter(hookManager, receiveLogger),
		hookManager:    hookManager,
		logger:         logger,
	}
}

// CreateNetworkManager creates a new NetworkManager with the integrated registries
func (i *NetworkRegistryIntegrator) CreateNetworkManager(networkInterface network.NetworkInterface) *network.NetworkManager {
	// Initialize the registries
	i.sendAdapter.Initialize()
	i.receiveAdapter.Initialize()

	// Set the connection on the send adapter
	i.sendAdapter.SetConnection(networkInterface)

	// Create the network manager with the adapters
	return network.NewNetworkManager(networkInterface, i.sendAdapter, i.receiveAdapter)
}

// GetHookManager returns the hook manager used by the integrator
func (i *NetworkRegistryIntegrator) GetHookManager() *hooks.HookManager {
	return i.hookManager
}

// loggerAdapter adapts our Logger interface to the specific logger interfaces required by packages
type loggerAdapter struct {
	logger Logger
}

func (l *loggerAdapter) Debug(format string, args ...interface{}) {
	l.logger.Debug(format, args...)
}

func (l *loggerAdapter) Info(format string, args ...interface{}) {
	l.logger.Info(format, args...)
}

func (l *loggerAdapter) Warning(format string, args ...interface{}) {
	l.logger.Warning(format, args...)
}

func (l *loggerAdapter) Error(format string, args ...interface{}) {
	l.logger.Error(format, args...)
}

func (l *loggerAdapter) Success(format string, args ...interface{}) {
	l.logger.Success(format, args...)
}
