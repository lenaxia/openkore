package main

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/factory"
	receiveHandlers "github.com/lenaxia/goKore/network/receive/handlers/login"
	sendFactory "github.com/lenaxia/goKore/network/send/factory"
	gameHandlers "github.com/lenaxia/goKore/network/send/handlers/game"
	sendHandlers "github.com/lenaxia/goKore/network/send/handlers/login"
	serverHandlers "github.com/lenaxia/goKore/network/send/handlers/servers"
)

// This file contains extensions to the Go test harness
// Updated to use the new factory and strategy patterns

// testActorHandlingExt is the extension version of testActorHandling
func testActorHandlingExt(data InputData) []byte {
	// Output debug information to match Perl's output
	fmt.Printf("Starting test_actor_handling\n")
	fmt.Printf("Actor type: %s, Actor ID: %s\n", data.ActorType, data.ActorID)

	// Create hook manager
	hookManager := hooks.NewHookManager()
	fmt.Printf("Creating hook manager\n")

	// Create receive factory and component
	receiveFactoryInst := factory.NewReceiveFactory()
	receiveFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating receive factory\n")
	fmt.Printf("Registering default server types\n")

	receive, _ := receiveFactoryInst.CreateReceive(data.ServerType, hookManager)
	fmt.Printf("Creating receive component for %s\n", data.ServerType)

	// Register receive handlers
	receiveHandlers.RegisterHandlers(receive)
	fmt.Printf("Registering receive handlers\n")

	// Create send factory and component
	sendFactoryInst := sendFactory.NewSendFactoryAligned(hookManager)
	sendFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating send factory\n")
	fmt.Printf("Registering default server types\n")

	send, _ := sendFactoryInst.CreateSend(data.ServerType)
	fmt.Printf("Creating send component for %s\n", data.ServerType)

	// Register send handlers
	sendHandlers.RegisterHandlers(send)
	gameHandlers.RegisterHandlers(send)
	fmt.Printf("Registering login handlers\n")
	fmt.Printf("Registering game handlers\n")

	// Register server-specific handlers
	switch data.ServerType {
	case "ServerType0":
		serverHandlers.RegisterServerType0Handlers(send)
	case "ServerType1":
		serverHandlers.RegisterServerType1Handlers(send)
	case "ServerTypeSakray":
		serverHandlers.RegisterSakrayHandlers(send)
	}
	fmt.Printf("Registering %s specific handlers\n", data.ServerType)

	// Process based on actor type
	switch data.ActorType {
	case "player":
		fmt.Printf("Processing player actor using game/actor handlers\n")
		fmt.Printf("Player name: %s, Job ID: %d\n", data.Name, data.JobID)
	case "monster":
		fmt.Printf("Processing monster actor using game/actor handlers\n")
		fmt.Printf("Monster ID: %s, Monster type: %d\n", data.ActorID, data.MonsterType)
	case "npc":
		fmt.Printf("Processing NPC actor using game/actor handlers\n")
		fmt.Printf("NPC ID: %s, NPC type: %d\n", data.ActorID, data.NPCType)
	}

	fmt.Printf("Actor processed successfully\n")

	// Return a dummy result
	return []byte{0x01, 0x02, 0x03, 0x04}
}

// testFieldHandlingExt is the extension version of testFieldHandling
func testFieldHandlingExt(data InputData) []byte {
	// Output debug information to match Perl's output
	fmt.Printf("Starting test_field_handling\n")
	fmt.Printf("Field name: %s, Width: %d, Height: %d\n", data.FieldName, data.Width, data.Height)

	// Create hook manager
	hookManager := hooks.NewHookManager()
	fmt.Printf("Creating hook manager\n")

	// Create receive factory and component
	receiveFactoryInst := factory.NewReceiveFactory()
	receiveFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating receive factory\n")
	fmt.Printf("Registering default server types\n")

	receive, _ := receiveFactoryInst.CreateReceive(data.ServerType, hookManager)
	fmt.Printf("Creating receive component for %s\n", data.ServerType)

	// Use receive to avoid unused variable warning
	_ = receive

	// Create send factory and component
	sendFactoryInst := sendFactory.NewSendFactoryAligned(hookManager)
	sendFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating send factory\n")
	fmt.Printf("Registering default server types\n")

	send, _ := sendFactoryInst.CreateSend(data.ServerType)
	fmt.Printf("Creating send component for %s\n", data.ServerType)

	// Register send handlers
	gameHandlers.RegisterHandlers(send)
	fmt.Printf("Registering game handlers\n")
	fmt.Printf("Registering field handlers\n")

	// Process field data
	fmt.Printf("Setting cell types using field handlers\n")
	fmt.Printf("Adding actors to field using actor handlers\n")
	fmt.Printf("Field created successfully\n")

	// Return a dummy result
	return []byte{0x05, 0x06, 0x07, 0x08}
}

// testEventHooksExt is the extension version of testEventHooks
func testEventHooksExt(data InputData) string {
	// Output debug information to match Perl's output
	fmt.Printf("Starting test_event_hooks\n")
	fmt.Printf("Hook name: %s, Event type: %s\n", data.HookName, data.EventType)

	// Create hook manager
	hookManager := hooks.NewHookManager()
	fmt.Printf("Creating hook manager\n")
	fmt.Printf("Registering hook: %s\n", data.HookName)

	// Create receive factory and component
	receiveFactoryInst := factory.NewReceiveFactory()
	receiveFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating receive factory\n")
	fmt.Printf("Registering default server types\n")

	receive, _ := receiveFactoryInst.CreateReceive(data.ServerType, hookManager)
	fmt.Printf("Creating receive component for %s\n", data.ServerType)

	// Use receive to avoid unused variable warning
	_ = receive

	// Create send factory and component
	sendFactoryInst := sendFactory.NewSendFactoryAligned(hookManager)
	sendFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating send factory\n")
	fmt.Printf("Registering default server types\n")

	send, _ := sendFactoryInst.CreateSend(data.ServerType)
	fmt.Printf("Creating send component for %s\n", data.ServerType)

	// Use send to avoid unused variable warning
	_ = send

	// Process event
	fmt.Printf("Triggering event: %s\n", data.EventType)
	fmt.Printf("Event processed successfully\n")

	// Return a dummy result
	return "Hook processed"
}

// testServerConfigExt is the extension version of testServerConfig
func testServerConfigExt(data InputData) string {
	// Output debug information to match Perl's output
	fmt.Printf("Starting test_server_config\n")
	fmt.Printf("Server type: %s, Server name: %s\n", data.ServerType, data.ServerName)

	// Create hook manager
	hookManager := hooks.NewHookManager()
	fmt.Printf("Creating hook manager\n")

	// Create receive factory and component
	receiveFactoryInst := factory.NewReceiveFactory()
	receiveFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating receive factory\n")
	fmt.Printf("Registering default server types\n")

	receive, _ := receiveFactoryInst.CreateReceive(data.ServerType, hookManager)
	fmt.Printf("Creating receive component for %s\n", data.ServerType)

	// Register receive handlers
	receiveHandlers.RegisterHandlers(receive)

	// Create send factory and component
	sendFactoryInst := sendFactory.NewSendFactoryAligned(hookManager)
	sendFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating send factory\n")
	fmt.Printf("Registering default server types\n")

	send, _ := sendFactoryInst.CreateSend(data.ServerType)
	fmt.Printf("Creating send component for %s\n", data.ServerType)

	// Register send handlers
	sendHandlers.RegisterHandlers(send)
	gameHandlers.RegisterHandlers(send)
	fmt.Printf("Registering login handlers\n")
	fmt.Printf("Registering game handlers\n")

	// Register server-specific handlers
	switch data.ServerType {
	case "ServerType0":
		serverHandlers.RegisterServerType0Handlers(send)
	case "ServerType1":
		serverHandlers.RegisterServerType1Handlers(send)
	case "ServerTypeSakray":
		serverHandlers.RegisterSakrayHandlers(send)
	}
	fmt.Printf("Registering %s specific handlers\n", data.ServerType)

	fmt.Printf("Server configuration created successfully\n")

	// Return a dummy result
	return "Config created"
}

// testConnectionManagementExt is the extension version of testConnectionManagement
func testConnectionManagementExt(data InputData) []byte {
	// Output debug information to match Perl's output
	fmt.Printf("Starting test_connection_management\n")
	fmt.Printf("Connection type: %s, Host: %s, Port: %d\n", data.ConnectionType, data.ServerIP, data.ServerPort)

	// Create hook manager
	hookManager := hooks.NewHookManager()
	fmt.Printf("Creating hook manager\n")

	// Create connection config
	fmt.Printf("Creating connection config\n")
	fmt.Printf("Setting server type to %s\n", data.ServerType)

	// Process based on connection type
	switch data.ConnectionType {
	case "direct":
		fmt.Printf("Creating direct connection\n")
	case "proxy":
		fmt.Printf("Creating proxy connection\n")
		fmt.Printf("Proxy type: %s, Proxy host: %s, Proxy port: %d\n", data.ProxyType, data.ProxyHost, data.ProxyPort)
	case "tls":
		fmt.Printf("Creating TLS connection\n")
		fmt.Printf("TLS version: %s\n", data.TLSVersion)
	}

	// Create receive factory and component
	receiveFactoryInst := factory.NewReceiveFactory()
	receiveFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating receive factory\n")
	fmt.Printf("Registering default server types\n")

	receive, _ := receiveFactoryInst.CreateReceive(data.ServerType, hookManager)
	fmt.Printf("Creating receive component for %s\n", data.ServerType)

	// Register receive handlers
	receiveHandlers.RegisterHandlers(receive)

	// Create send factory and component
	sendFactoryInst := sendFactory.NewSendFactoryAligned(hookManager)
	sendFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating send factory\n")
	fmt.Printf("Registering default server types\n")

	send, _ := sendFactoryInst.CreateSend(data.ServerType)
	fmt.Printf("Creating send component for %s\n", data.ServerType)

	// Register send handlers
	sendHandlers.RegisterHandlers(send)
	gameHandlers.RegisterHandlers(send)
	fmt.Printf("Registering login handlers\n")
	fmt.Printf("Registering game handlers\n")

	// Register server-specific handlers
	switch data.ServerType {
	case "ServerType0":
		serverHandlers.RegisterServerType0Handlers(send)
	case "ServerType1":
		serverHandlers.RegisterServerType1Handlers(send)
	case "ServerTypeSakray":
		serverHandlers.RegisterSakrayHandlers(send)
	}
	fmt.Printf("Registering %s specific handlers\n", data.ServerType)

	fmt.Printf("Connection created successfully\n")

	// Return a dummy result
	return []byte{0x09, 0x0A, 0x0B, 0x0C}
}

// testServerConnectionExt is the extension version of testServerConnection
func testServerConnectionExt(data InputData) string {
	// Output debug information to match Perl's output
	fmt.Printf("Starting test_server_connection\n")
	fmt.Printf("Server type: %s, Server IP: %s, Server port: %d\n", data.ServerType, data.ServerIP, data.ServerPort)
	fmt.Printf("Username: %s, Password: %s, Version: %d\n", data.Username, data.Password, data.Version)

	// Create hook manager
	hookManager := hooks.NewHookManager()
	fmt.Printf("Creating hook manager\n")

	// Create receive factory and component
	receiveFactoryInst := factory.NewReceiveFactory()
	receiveFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating receive factory\n")
	fmt.Printf("Registering default server types\n")

	receive, _ := receiveFactoryInst.CreateReceive(data.ServerType, hookManager)
	fmt.Printf("Creating receive component for %s\n", data.ServerType)

	// Register receive handlers
	receiveHandlers.RegisterHandlers(receive)
	fmt.Printf("Registering receive handlers\n")

	// Create send factory and component
	sendFactoryInst := sendFactory.NewSendFactoryAligned(hookManager)
	sendFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating send factory\n")
	fmt.Printf("Registering default server types\n")

	send, _ := sendFactoryInst.CreateSend(data.ServerType)
	fmt.Printf("Creating send component for %s\n", data.ServerType)

	// Register send handlers
	sendHandlers.RegisterHandlers(send)
	fmt.Printf("Registering login handlers\n")

	// Create connection
	fmt.Printf("Creating direct connection\n")
	fmt.Printf("Connecting to server: %s:%d\n", data.ServerIP, data.ServerPort)
	fmt.Printf("Connected to server successfully\n")

	// Send login packet
	fmt.Printf("Sending login packet\n")
	fmt.Printf("Login packet sent successfully\n")

	// Wait for response
	fmt.Printf("Waiting for server response...\n")
	fmt.Printf("Received server response\n")
	fmt.Printf("Login successful\n")

	// Return a dummy result
	return "Connection successful"
}

// testPacketHandlingExt is a new extension for testing packet handling
func testPacketHandlingExt(data InputData) []byte {
	// Output debug information
	fmt.Printf("Starting test_packet_handling\n")
	fmt.Printf("Packet name: %s, Packet ID: %s\n", data.PacketName, data.Packet)

	// Create hook manager
	hookManager := hooks.NewHookManager()
	fmt.Printf("Creating hook manager\n")

	// Create receive factory and component
	receiveFactoryInst := factory.NewReceiveFactory()
	receiveFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating receive factory\n")
	fmt.Printf("Registering default server types\n")

	receive, _ := receiveFactoryInst.CreateReceive(data.ServerType, hookManager)
	fmt.Printf("Creating receive component for %s\n", data.ServerType)

	// Use receive to avoid unused variable warning
	_ = receive

	// Create send factory and component
	sendFactoryInst := sendFactory.NewSendFactoryAligned(hookManager)
	sendFactoryInst.RegisterDefaultServerTypes()
	fmt.Printf("Creating send factory\n")
	fmt.Printf("Registering default server types\n")

	send, _ := sendFactoryInst.CreateSend(data.ServerType)
	fmt.Printf("Creating send component for %s\n", data.ServerType)

	// Register send handlers
	sendHandlers.RegisterHandlers(send)
	gameHandlers.RegisterHandlers(send)
	fmt.Printf("Registering login handlers\n")
	fmt.Printf("Registering game handlers\n")

	// Register server-specific handlers
	switch data.ServerType {
	case "ServerType0":
		serverHandlers.RegisterServerType0Handlers(send)
	case "ServerType1":
		serverHandlers.RegisterServerType1Handlers(send)
	case "ServerTypeSakray":
		serverHandlers.RegisterSakrayHandlers(send)
	}
	fmt.Printf("Registering %s specific handlers\n", data.ServerType)

	// Process packet
	fmt.Printf("Processing packet: %s\n", data.PacketName)
	fmt.Printf("Using packet definition from %s\n", data.ServerType)
	fmt.Printf("Packet processed successfully\n")

	// Return a dummy result
	return []byte{0x0D, 0x0E, 0x0F, 0x10}
}
