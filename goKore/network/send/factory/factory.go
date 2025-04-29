// Package factory provides a factory for creating send managers.
package factory

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send"
	"github.com/lenaxia/goKore/network/send/core"
	"github.com/lenaxia/goKore/network/send/game/cashshop"
	"github.com/lenaxia/goKore/network/send/game/chat"
	"github.com/lenaxia/goKore/network/send/game/misc"
)

// SendFactory is a factory for creating send managers
type SendFactory struct {
	hookManager                 *hooks.HookManager
	packetConstructionProviders map[string]send.PacketConstructionProvider
}

// NewSendFactoryAligned creates a new send factory with the given hook manager
func NewSendFactoryAligned(hookManager *hooks.HookManager) *SendFactory {
	return &SendFactory{
		hookManager:                 hookManager,
		packetConstructionProviders: make(map[string]send.PacketConstructionProvider),
	}
}

// CreateSend creates a new Send implementation for the given server type
func (sf *SendFactory) CreateSend(serverType string) (*send.BaseSend, error) {
	// Get the packet construction provider for the server type
	provider, exists := sf.packetConstructionProviders[serverType]
	if !exists {
		return nil, fmt.Errorf("no packet construction provider for server type: %s", serverType)
	}

	// Create a new Send implementation
	baseSend := send.NewBaseSend(sf.hookManager)

	// Configure the Send implementation with the packet constructions
	err := baseSend.Configure(serverType, provider())
	if err != nil {
		return nil, fmt.Errorf("failed to configure send: %w", err)
	}

	return baseSend, nil
}

// RegisterPacketConstructionProvider registers a packet construction provider for a server type
func (sf *SendFactory) RegisterPacketConstructionProvider(serverType string, provider send.PacketConstructionProvider) {
	sf.packetConstructionProviders[serverType] = provider
}

// RegisterServerType is an alias for RegisterPacketConstructionProvider
// It registers a packet construction provider for a server type
func (sf *SendFactory) RegisterServerType(serverType string, provider send.PacketConstructionProvider) {
	sf.RegisterPacketConstructionProvider(serverType, provider)
}

// RegisterDefaultServerTypes registers the default server types
func (sf *SendFactory) RegisterDefaultServerTypes() {
	// This is a placeholder for registering default server types
	// In a real implementation, this would register all the default server types
}

// CoreSendAdapter adapts an interface{} to core.Send
type CoreSendAdapter struct {
	core.Send
}

// NewCoreSendAdapter creates a new CoreSendAdapter
func NewCoreSendAdapter(send interface{}) *CoreSendAdapter {
	if coreSend, ok := send.(core.Send); ok {
		return &CoreSendAdapter{Send: coreSend}
	}
	// Return a default implementation if the conversion fails
	return &CoreSendAdapter{Send: &core.BaseSend{}}
}

// CreateCashShopManager creates a new cash shop manager
func (sf *SendFactory) CreateCashShopManager(baseSend interface{}) *cashshop.CashShopManager {
	adapter := NewCoreSendAdapter(baseSend)
	return cashshop.NewCashShopManager(adapter)
}

// CreateMiscManager creates a new miscellaneous manager
func (sf *SendFactory) CreateMiscManager(baseSend interface{}) *misc.MiscManager {
	adapter := NewCoreSendAdapter(baseSend)
	return misc.NewMiscManager(adapter)
}

// CreateInfoChatManager creates a new info chat manager
func (sf *SendFactory) CreateInfoChatManager(baseSend interface{}) *chat.InfoChatManager {
	adapter := NewCoreSendAdapter(baseSend)
	return chat.NewInfoChatManager(adapter)
}
