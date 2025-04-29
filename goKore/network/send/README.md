# Send Component Alignment

This directory contains the implementation of the Send component alignment according to the hybrid architecture approach. The Send component is responsible for constructing and sending packets to the server.

## Architecture Overview

The Send component has been aligned with the Receive component to provide a consistent architecture across both sides of the network implementation. The key components of this architecture are:

1. **Packet Definition System**: A shared system for defining packet structures, formats, and field names.
2. **Configurable Send Component**: A configurable Send interface that can be adapted to different server types.
3. **Factory for Send Components**: A factory that creates and configures Send implementations based on server type.
4. **Domain-Organized Packet Constructions**: Packet constructions organized by domain (login, movement, chat, etc.).

## Directory Structure

```
send/
├── core/                  # Core Send implementation
│   ├── base_send.go       # Base Send implementation
│   ├── base_send_new.go   # New aligned Send implementation
│   └── send.go            # Send interface definition
├── factory/               # Factory for creating Send implementations
│   ├── factory.go         # Original factory implementation
│   └── factory_new.go     # New aligned factory implementation
├── game/                  # Game-related packet constructions
│   ├── login/             # Login-related packet constructions
│   ├── movement/          # Movement-related packet constructions
│   ├── chat/              # Chat-related packet constructions
│   └── item/              # Item-related packet constructions
├── servers/               # Server-specific packet constructions
│   ├── sakray/            # Sakray server-specific packet constructions
│   └── ...                # Other server types
├── types/                 # Type definitions
│   └── send.go            # Send interface definition
└── integration.go         # Integration example
```

## Key Components

### Packet Definition

Packet definitions are shared between Send and Receive components through the `common` package:

```go
// PacketDef defines the structure of a packet for receiving
type PacketDef struct {
    ID         string
    Name       string
    Format     string
    FieldNames []string
}

// PacketConstruction defines how to construct a packet for sending
type PacketConstruction struct {
    ID         string
    Name       string
    Format     string
    FieldNames []string
}
```

### Send Interface

The Send interface defines the contract for Send implementations:

```go
// Send defines the interface for packet construction and sending
type Send interface {
    // RegisterHandler registers a handler for a specific packet
    RegisterHandler(packetName string, handler SendHandler)

    // ConstructPacket constructs a packet from a packet name and arguments
    ConstructPacket(packetName string, args map[string]interface{}) ([]byte, error)

    // SendPacket constructs and sends a packet
    SendPacket(packetName string, args map[string]interface{}) error

    // SendToServer sends a raw packet to the server
    SendToServer(packet []byte) error

    // Configure configures the send component with server-specific packet constructions
    Configure(serverType string, packetConstructions map[string]PacketConstruction) error
}
```

### BaseSend Implementation

The BaseSend implementation provides the core functionality for sending packets:

```go
// BaseSend implements the Send interface
type BaseSend struct {
    conn               interface{}
    hookManager        *hooks.HookManager
    serverType         string
    handlers           map[string]SendHandler
    packetConstructions map[string]PacketConstruction
    packetLUT          map[string]string
    // ...
}
```

### SendFactory

The SendFactory creates and configures Send implementations based on server type:

```go
// SendFactory creates and configures Send implementations
type SendFactory struct {
    packetConstructionProviders map[string]PacketConstructionProvider
    hookManager                 *hooks.HookManager
}
```

### Domain-Specific Packet Constructions

Packet constructions are organized by domain to improve maintainability:

```go
// Login-related packet constructions
func GetLoginPacketConstructions() map[string]PacketConstruction {
    return map[string]PacketConstruction{
        "0064": {
            ID:         "0064",
            Name:       "login_request",
            Format:     "v a24 a24 C",
            FieldNames: []string{"version", "username", "password", "clienttype"},
        },
        // ...
    }
}
```

## Usage Example

```go
// Create a hook manager
hookManager := hooks.NewHookManager()

// Create a send factory
factory := NewSendFactory(hookManager)

// Register server types
factory.RegisterDefaultServerTypes()

// Create a send implementation for a specific server type
send, err := factory.CreateSend("ServerType0")
if err != nil {
    // Handle error
}

// Use the send implementation
err = send.SendPacket("login_request", map[string]interface{}{
    "version":    1,
    "username":   "user",
    "password":   "pass",
    "clienttype": 0,
})
if err != nil {
    // Handle error
}
```

## Alignment with Receive Side

The Send component has been aligned with the Receive component to provide a consistent architecture:

1. **Shared Packet Definitions**: Both components use the same packet definition structure.
2. **Configurable Components**: Both components can be configured for different server types.
3. **Factory Pattern**: Both components use a factory to create and configure implementations.
4. **Domain Organization**: Both components organize handlers by domain.

This alignment improves maintainability, reduces duplication, and provides a consistent architecture across both sides of the network implementation.