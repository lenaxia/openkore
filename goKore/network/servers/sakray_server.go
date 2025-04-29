// Package servers provides server implementations for different Ragnarok Online server types.
package servers

import (
	"fmt"

	"github.com/lenaxia/goKore/network/packets"
)

// SakrayServer implements the Sakray server type
type SakrayServer struct {
	*BaseServer
}

// NewSakrayServer creates a new Sakray server
func NewSakrayServer(baseConfig *BaseServerConfig) (*SakrayServer, error) {
	// Create the base server
	baseServer, err := NewBaseServer(baseConfig)
	if err != nil {
		return nil, err
	}

	// Create the Sakray server
	server := &SakrayServer{
		BaseServer: baseServer,
	}

	// Register Sakray-specific packet handlers
	server.registerSakrayHandlers()

	return server, nil
}

// registerSakrayHandlers registers Sakray-specific packet handlers
func (s *SakrayServer) registerSakrayHandlers() {
	// Override or add Sakray-specific packet handlers
	// These would be different from the base server handlers for packets that
	// have different structures or behaviors in Sakray servers

	// Example: Override a packet handler for a Sakray-specific packet format
	s.RegisterPacketHandler("0069", "account_server_info", "v a4 a4 a4 a4 a26 C a*",
		[]string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"},
		func(args map[string]interface{}) error {
			// Call the hook for this packet
			if s.hookManager != nil {
				s.hookManager.CallHook("sakray/account_server_info", args)
			}

			// Process the packet (Sakray-specific logic would go here)
			// For now, just log that we received a Sakray-specific packet
			fmt.Println("Received Sakray-specific account_server_info packet")

			return nil
		})

	// Add more Sakray-specific packet handlers as needed
}

// LoadSakrayPacketDefinitions loads Sakray-specific packet definitions
func (s *SakrayServer) LoadSakrayPacketDefinitions() error {
	// Load Sakray-specific packet definitions
	// This would typically load from a configuration file or embedded data

	// Example: Create a new packet database with Sakray-specific packet definitions
	packetDB := packets.NewPacketDatabase()

	// In a real implementation, you would load packet definitions from a file or embedded data
	// For now, just add a few example packet definitions

	// Example packet definition
	packetDef := packets.NewPacketDefinition(
		"0069",
		"account_server_info",
		"v a4 a4 a4 a4 a26 C a*",
		[]string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"},
	)
	packetDB.AddPacketDefinition(packetDef)

	// Merge with existing packet definitions
	// This is a placeholder - in a real implementation, you would merge the
	// loaded packet definitions with the existing ones

	return nil
}

// GetServerType returns the server type
func (s *SakrayServer) GetServerType() ServerType {
	return ServerTypeSakray
}

// Connect overrides the base Connect method to add Sakray-specific connection logic
func (s *SakrayServer) Connect() error {
	// Call the base Connect method
	if err := s.BaseServer.Connect(); err != nil {
		return err
	}

	// Add Sakray-specific connection logic
	// For example, Sakray servers might require a different handshake or initialization

	// Log that we're connecting to a Sakray server
	fmt.Println("Connecting to Sakray server:", s.serverInfo.Name)

	return nil
}

// Login overrides the base Login method to add Sakray-specific login logic
func (s *SakrayServer) Login() error {
	// Call the base Login method
	if err := s.BaseServer.Login(); err != nil {
		return err
	}

	// Add Sakray-specific login logic
	// For example, Sakray servers might require additional authentication steps

	// Log that we're logging in to a Sakray server
	fmt.Println("Logging in to Sakray server:", s.serverInfo.Name)

	return nil
}

// ProcessPackets overrides the base ProcessPackets method to add Sakray-specific packet processing
func (s *SakrayServer) ProcessPackets() error {
	// Call the base ProcessPackets method
	if err := s.BaseServer.ProcessPackets(); err != nil {
		return err
	}

	// Add Sakray-specific packet processing
	// For example, Sakray servers might have additional packet types or different packet formats

	return nil
}
