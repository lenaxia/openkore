// Package main provides an example of how to use the receive component.
package main

import (
	"fmt"
	"os"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/factory"
	"github.com/lenaxia/goKore/network/receive/handlers/game"
	"github.com/lenaxia/goKore/network/receive/handlers/login"
)

func main() {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a receive factory
	receiveFactory := factory.NewReceiveFactory()

	// Register default server types
	receiveFactory.RegisterDefaultServerTypes()

	// Create a receive component for ServerType0
	receive, err := receiveFactory.CreateReceive("ServerType0", hookManager)
	if err != nil {
		fmt.Printf("Error creating receive component: %v\n", err)
		os.Exit(1)
	}

	// Register domain-specific handlers
	login.RegisterHandlers(receive)
	game.RegisterActorHandlers(receive)

	// Example packet data (account_server_info packet)
	// This is a simplified example - in a real implementation, this would come from the network
	packet := []byte{
		0x69, 0x00, // Packet ID (0069)
		0x4C, 0x00, // Length (76 bytes)
		0x01, 0x02, 0x03, 0x04, // Session ID
		0x05, 0x06, 0x07, 0x08, // Account ID
		0x09, 0x0A, 0x0B, 0x0C, // Session ID 2
		0x0D, 0x0E, 0x0F, 0x10, // Last Login IP
		// Last Login Time (26 bytes)
		0x32, 0x30, 0x32, 0x35, 0x2D, 0x30, 0x34, 0x2D, 0x32, 0x37, 0x20, 0x32, 0x30, 0x3A, 0x31, 0x34, 0x3A, 0x30, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, // Account Sex (1 = male)
		// Server Info (remaining bytes)
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
		0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14,
		0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E,
	}

	// Process the packet
	fmt.Println("Processing account_server_info packet...")
	err = receive.Process(packet)
	if err != nil {
		fmt.Printf("Error processing packet: %v\n", err)
		os.Exit(1)
	}

	// Example packet data (login_error packet)
	packet = []byte{
		0x6A, 0x00, // Packet ID (006A)
		0x01, // Error type
		// Error message (20 bytes)
		0x49, 0x6E, 0x76, 0x61, 0x6C, 0x69, 0x64, 0x20, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6F, 0x72, 0x64, 0x00, 0x00, 0x00, 0x00,
	}

	// Process the packet
	fmt.Println("\nProcessing login_error packet...")
	err = receive.Process(packet)
	if err != nil {
		fmt.Printf("Error processing packet: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nAll packets processed successfully!")
}
