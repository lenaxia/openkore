// Package servers provides server implementations for different Ragnarok Online server types.
package servers

import (
	"github.com/lenaxia/goKore/network/common"
	receiveServers "github.com/lenaxia/goKore/network/receive/servers"
	sendServers "github.com/lenaxia/goKore/network/send/servers"
)

// ServerConfig contains both send and receive packet configurations
type ServerConfig struct {
	SendConfig    map[string]common.PacketConstruction
	ReceiveConfig map[string]common.PacketConstruction
}

// GetServerType0Config returns the send and receive packet configurations for ServerType0
func GetServerType0Config() *ServerConfig {
	return &ServerConfig{
		SendConfig:    GetServerType0SendConfig(),
		ReceiveConfig: GetServerType0ReceiveConfig(),
	}
}

// GetServerType0SendConfig returns the send packet configurations for ServerType0
func GetServerType0SendConfig() map[string]common.PacketConstruction {
	return sendServers.ServerType0PacketConstructions()
}

// GetServerType0ReceiveConfig returns the receive packet configurations for ServerType0
func GetServerType0ReceiveConfig() map[string]common.PacketConstruction {
	return receiveServers.ServerType0PacketConstructions()
}

// GetSakrayConfig returns the send and receive packet configurations for Sakray servers
func GetSakrayConfig() *ServerConfig {
	return &ServerConfig{
		SendConfig:    GetSakraySendConfig(),
		ReceiveConfig: GetSakrayReceiveConfig(),
	}
}

// GetSakraySendConfig returns the send packet configurations for Sakray servers
func GetSakraySendConfig() map[string]common.PacketConstruction {
	return sendServers.SakrayPacketConstructions()
}

// GetSakrayReceiveConfig returns the receive packet configurations for Sakray servers
func GetSakrayReceiveConfig() map[string]common.PacketConstruction {
	// For now, use ServerType0 packet constructions since Sakray doesn't have its own
	return receiveServers.ServerType0PacketConstructions()
}

// GetServerConfigByType returns the server configuration for the specified server type
func GetServerConfigByType(serverType ServerType) *ServerConfig {
	switch serverType {
	case ServerTypeSakray:
		return GetSakrayConfig()
	case ServerTypeOfficial:
		return GetServerType0Config()
	default:
		// Default to ServerType0
		return GetServerType0Config()
	}
}
