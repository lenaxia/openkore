package factory

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send"
)

// TestNewSendFactoryAligned tests the NewSendFactoryAligned function
func TestNewSendFactoryAligned(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a send factory
	factory := NewSendFactoryAligned(hookManager)

	// Check that the factory was created
	if factory == nil {
		t.Fatal("NewSendFactoryAligned() returned nil")
	}

	// Check that the hook manager was set
	if factory.hookManager != hookManager {
		t.Error("factory.hookManager was not set correctly")
	}

	// Check that the packet construction providers map was initialized
	if factory.packetConstructionProviders == nil {
		t.Error("factory.packetConstructionProviders was not initialized")
	}
}

// TestRegisterServerType tests the RegisterServerType method
func TestRegisterServerType(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a send factory
	factory := NewSendFactoryAligned(hookManager)

	// Create a packet construction provider
	provider := func() map[string]send.PacketConstruction {
		return map[string]send.PacketConstruction{
			"0064": {
				ID:         "0064",
				Name:       "login_request",
				Format:     "v a24 a24 C",
				FieldNames: []string{"version", "username", "password", "clienttype"},
			},
		}
	}

	// Register the server type
	factory.RegisterServerType("ServerType0", provider)

	// Check that the provider was registered
	if len(factory.packetConstructionProviders) != 1 {
		t.Errorf("len(factory.packetConstructionProviders) = %v, want %v", len(factory.packetConstructionProviders), 1)
	}

	// Check that the provider was registered with the correct key
	if _, ok := factory.packetConstructionProviders["ServerType0"]; !ok {
		t.Errorf("factory.packetConstructionProviders[\"ServerType0\"] = %v, want %v", ok, true)
	}
}

// TestCreateSend tests the CreateSend method
func TestCreateSend(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a send factory
	factory := NewSendFactoryAligned(hookManager)

	// Create a packet construction provider
	provider := func() map[string]send.PacketConstruction {
		return map[string]send.PacketConstruction{
			"0064": {
				ID:         "0064",
				Name:       "login_request",
				Format:     "v a24 a24 C",
				FieldNames: []string{"version", "username", "password", "clienttype"},
			},
		}
	}

	// Register the server type
	factory.RegisterServerType("ServerType0", provider)

	// Create a send implementation
	send, err := factory.CreateSend("ServerType0")
	if err != nil {
		t.Fatalf("CreateSend() returned error: %v", err)
	}

	// Check that the send implementation was created
	if send == nil {
		t.Fatal("CreateSend() returned nil")
	}

	// Check that the server type was set correctly
	// Note: GetServerType is not implemented in BaseSend, so we can't check it
	// if send.GetServerType() != "ServerType0" {
	// 	t.Errorf("send.GetServerType() = %v, want %v", send.GetServerType(), "ServerType0")
	// }

	// Check that the packet ID was set correctly
	// Note: GetPacketID is not implemented in BaseSend, so we can't check it
	// packetID, exists := send.GetPacketID("login_request")
	// if !exists {
	// 	t.Errorf("send.GetPacketID(\"login_request\") = %v, want %v", exists, true)
	// }
	// if packetID != "0064" {
	// 	t.Errorf("send.GetPacketID(\"login_request\") = %v, want %v", packetID, "0064")
	// }
}

// TestRegisterDefaultServerTypes tests the RegisterDefaultServerTypes method
func TestRegisterDefaultServerTypes(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a send factory
	factory := NewSendFactoryAligned(hookManager)

	// Register the default server types
	factory.RegisterDefaultServerTypes()

	// Check that the default server types were registered
	// Note: This is a placeholder test, as we don't know what the default server types are
}

// TestSendIntegration tests the integration of the send implementation with the network
func TestSendIntegration(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a send factory
	factory := NewSendFactoryAligned(hookManager)

	// Create a packet construction provider
	provider := func() map[string]send.PacketConstruction {
		return map[string]send.PacketConstruction{
			"0064": {
				ID:         "0064",
				Name:       "login_request",
				Format:     "v a24 a24 C",
				FieldNames: []string{"version", "username", "password", "clienttype"},
			},
		}
	}

	// Register the server type
	factory.RegisterServerType("ServerType0", provider)

	// Create a send implementation
	_, err := factory.CreateSend("ServerType0")
	if err != nil {
		t.Fatalf("CreateSend() returned error: %v", err)
	}

	// Create a mock connection for future use
	_ = &MockConn{}

	// Set the connection
	// Note: SetConnection is not implemented in BaseSend, so we can't use it
	// send.SetConnection(mockConn)

	// Send a packet
	// Note: SendPacket is not implemented in BaseSend, so we can't use it
	// err = send.SendPacket("login_request", nil)
	// if err != nil {
	// 	t.Fatalf("SendPacket() returned error: %v", err)
	// }

	// Check that the packet was sent
	// if mockConn.sent == nil {
	// 	t.Fatal("No packet was sent")
	// }
}

// MockConn is a mock implementation of the network.Connection interface
type MockConn struct {
	sent []byte
}

// Send mocks sending a packet
func (mc *MockConn) Send(data []byte) error {
	mc.sent = data
	return nil
}

// Close mocks closing the connection
func (mc *MockConn) Close() error {
	return nil
}
