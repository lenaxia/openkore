# Send Alignment Implementation

This document describes the implementation of the Send alignment according to the hybrid architecture approach.

## Overview

The Send component has been aligned with the Receive component to provide a consistent architecture across both sides of the network implementation. The key components of this architecture are:

1. **Packet Definition System**: A shared system for defining packet structures, formats, and field names.
2. **Configurable Send Component**: A configurable Send interface that can be adapted to different server types.
3. **Factory for Send Components**: A factory that creates and configures Send implementations based on server type.
4. **Domain-Organized Packet Constructions**: Packet constructions organized by domain (login, movement, chat, etc.).

## Implementation Details

### 1. Shared Packet Definition System

We've created a shared package `common` that contains packet definition structures used by both Send and Receive components:

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

This ensures consistency in packet definitions across both components.

### 2. Updated Send Interface

The Send interface has been updated to align with the Receive interface:

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

This interface provides a clear contract for Send implementations and aligns with the Receive interface.

### 3. BaseSend Implementation

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

Key methods include:

- `Configure`: Configures the send component with server-specific packet constructions.
- `RegisterHandler`: Registers a handler for a specific packet.
- `ConstructPacket`: Constructs a packet from a packet name and arguments.
- `SendPacket`: Constructs and sends a packet.
- `SendToServer`: Sends a raw packet to the server.

### 4. SendFactory

The SendFactory creates and configures Send implementations based on server type:

```go
// SendFactory creates and configures Send implementations
type SendFactory struct {
    packetConstructionProviders map[string]PacketConstructionProvider
    hookManager                 *hooks.HookManager
}
```

Key methods include:

- `RegisterServerType`: Registers packet constructions for a server type.
- `CreateSend`: Creates and configures a Send implementation for a server type.
- `RegisterDefaultServerTypes`: Registers the default server types.

### 5. Domain-Organized Packet Constructions

Packet constructions are organized by domain to improve maintainability:

- `login`: Login-related packet constructions and handlers.
- `movement`: Movement-related packet constructions and handlers.
- `chat`: Chat-related packet constructions and handlers.
- `item`: Item-related packet constructions and handlers.

Each domain package provides:

- `RegisterHandlers`: Registers all domain-specific handlers.
- `GetPacketConstructions`: Returns domain-specific packet constructions.

### 6. Server-Specific Packet Constructions

Server-specific packet constructions are organized by server type:

- `sakray`: Sakray server-specific packet constructions and handlers.
- Other server types...

Each server package provides:

- `RegisterHandlers`: Registers all server-specific handlers.
- `GetPacketConstructions`: Returns server-specific packet constructions that override or extend base constructions.

## Benefits of the Implementation

1. **Consistent Architecture**: Both send and receive sides follow the same patterns.
2. **Centralized Configuration**: Server-specific configurations in one place.
3. **Improved Maintainability**: Clear separation of concerns.
4. **Better Scalability**: Easy to add new server types.
5. **Reduced Duplication**: Shared code and patterns between send and receive.

## Next Steps

1. **Migration**: Migrate existing code to the new architecture.
2. **Testing**: Create comprehensive tests for the new implementation.
3. **Documentation**: Update documentation to reflect the new architecture.
4. **Integration**: Integrate the new Send component with the rest of the system.

## Conclusion

The Send alignment implementation provides a consistent architecture across both sides of the network implementation. It improves maintainability, reduces duplication, and provides a clear path for future development.