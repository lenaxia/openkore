package factory

import (
	"testing"

	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
)

// TestNewReceiveFactory tests the NewReceiveFactory function
func TestNewReceiveFactory(t *testing.T) {
	// Create a receive factory
	factory := NewReceiveFactory()

	// Check that the factory was created
	if factory == nil {
		t.Fatal("NewReceiveFactory() returned nil")
	}

	// Check that the packet construction providers map was initialized
	if factory.packetConstructionProviders == nil {
		t.Error("factory.packetConstructionProviders was not initialized")
	}
}

// TestRegisterServerType tests the RegisterServerType method
func TestRegisterServerType(t *testing.T) {
	// Create a receive factory
	factory := NewReceiveFactory()

	// Create a packet construction provider
	provider := func() map[string]common.PacketConstruction {
		return map[string]common.PacketConstruction{
			"0064": {
				ID:         "0064",
				Name:       "login_response",
				Format:     "v V C x2 a4 a4 a4 V C2 a13 a3 x2",
				FieldNames: []string{"length", "account_id", "login_id1", "login_id2", "server_name", "server_ip", "server_port", "sex"},
			},
		}
	}

	// Register the server type
	factory.RegisterServerType("ServerType0", provider)

	// Check that the provider was registered
	if len(factory.packetConstructionProviders) != 1 {
		t.Errorf("len(factory.packetConstructionProviders) = %v, want %v", len(factory.packetConstructionProviders), 1)
	}

	// Check that the provider can be retrieved
	registeredProvider, exists := factory.packetConstructionProviders["ServerType0"]
	if !exists {
		t.Error("ServerType0 provider was not registered")
	} else {
		// Call the provider and check the result
		defs := registeredProvider()
		if len(defs) != 1 {
			t.Errorf("len(defs) = %v, want %v", len(defs), 1)
		}

		// Check that the packet definition was returned
		def, exists := defs["0064"]
		if !exists {
			t.Error("0064 packet definition was not returned")
		} else {
			// Check the packet definition fields
			if def.Name != "login_response" {
				t.Errorf("def.Name = %v, want %v", def.Name, "login_response")
			}
			if def.Format != "v V C x2 a4 a4 a4 V C2 a13 a3 x2" {
				t.Errorf("def.Format = %v, want %v", def.Format, "v V C x2 a4 a4 a4 V C2 a13 a3 x2")
			}
			if len(def.FieldNames) != 8 {
				t.Errorf("len(def.FieldNames) = %v, want %v", len(def.FieldNames), 8)
			}
		}
	}
}

// TestCreateReceive tests the CreateReceive method
func TestCreateReceive(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a receive factory
	factory := NewReceiveFactory()

	// Create a packet construction provider
	provider := func() map[string]common.PacketConstruction {
		return map[string]common.PacketConstruction{
			"0064": {
				ID:         "0064",
				Name:       "login_response",
				Format:     "v V C x2 a4 a4 a4 V C2 a13 a3 x2",
				FieldNames: []string{"length", "account_id", "login_id1", "login_id2", "server_name", "server_ip", "server_port", "sex"},
			},
		}
	}

	// Register the server type
	factory.RegisterServerType("ServerType0", provider)

	// Create a receive implementation
	receive, err := factory.CreateReceive("ServerType0", hookManager)
	if err != nil {
		t.Fatalf("CreateReceive() returned error: %v", err)
	}

	// Check that the receive implementation was created
	if receive == nil {
		t.Fatal("CreateReceive() returned nil")
	}

	// Test creating a receive implementation for a non-existent server type
	_, err = factory.CreateReceive("NonExistentServerType", hookManager)
	if err == nil {
		t.Error("CreateReceive() did not return an error for a non-existent server type")
	}
}

// TestRegisterDefaultServerTypes tests the RegisterDefaultServerTypes method
func TestRegisterDefaultServerTypes(t *testing.T) {
	// Create a receive factory
	factory := NewReceiveFactory()

	// Register the default server types
	factory.RegisterDefaultServerTypes()

	// Check that the server types were registered
	serverTypes := []string{"ServerType0", "ServerTypeSakray"}
	for _, serverType := range serverTypes {
		_, exists := factory.packetConstructionProviders[serverType]
		if !exists {
			t.Errorf("%s provider was not registered", serverType)
		}
	}
}

// TestIntegration tests the integration between the factory and the base receive
func TestIntegration(t *testing.T) {
	// Skip this test for now as it requires more complex setup
	t.Skip("Skipping integration test")
}
