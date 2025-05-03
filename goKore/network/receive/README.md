# Receive Component

This package implements the receive side of the network architecture for goKore, following the domain-organized approach outlined in the main network architecture. It is designed to align with the Send component, sharing common patterns and structures.

## Architecture Overview

The receive component is designed with a modular, domain-organized architecture that allows for easy extension and customization. The key components are:

1. **Packet Definition System**: Defines packet structures with metadata (ID, name, format, field names)
2. **Configurable Receive Component**: Implements the Receive interface with server-specific configurations
3. **Factory for Receive Components**: Creates and configures Receive implementations for different server types
4. **Domain-Specific Managers**: Organizes packet handlers by domain (character, inventory, field, etc.)
5. **Hook Notification System**: Provides a consistent way for components to react to packet events

This architecture mirrors the Send component's structure, with both sides sharing common types and patterns.

## Package Structure

- `types/`: Contains the Receive interface and common types
- `common/`: Shared types between send and receive components
- `core/`: Core functionality for parsing and processing packets
- `base/`: Base implementation of the Receive interface
- `factory/`: Factory for creating and configuring Receive implementations
- `game/`: Domain-specific packet managers
  - `actor/`: Character and actor-related packet handlers
  - `field/`: Field and map-related packet handlers
  - `card/`: Card system packet handlers
- `security/`: Security-related packet handlers
- `misc/`: Miscellaneous packet handlers

## Domain-Organized Packet Handling

The receive component uses a domain-organized approach that leverages the hooks system to provide better organization while maintaining Go's preference for simplicity:

```go
// Domain-specific manager for field-related packets
type FieldManager struct {
    parser      *CoreParser
    hookManager *hooks.HookManager
}

// Create a new field manager
func NewFieldManager(parser *CoreParser, hookManager *hooks.HookManager) *FieldManager {
    return &FieldManager{
        parser:      parser,
        hookManager: hookManager,
    }
}

// Register all handlers related to field operations
func (m *FieldManager) RegisterHandlers() {
    // Register warp_portal handler
    m.parser.RegisterHandlerFunc("0091", "warp_portal", "Z24",
        []string{"map_name"},
        m.handleWarpPortal)
        
    // Register boss_map_info handler
    m.parser.RegisterHandlerFunc("09EB", "boss_map_info", "V2 V2 V4",
        []string{"mapID", "x", "y"},
        m.handleBossMapInfo)
}

// Handle warp_portal packet
func (m *FieldManager) handleWarpPortal(args map[string]interface{}) error {
    // Process the packet
    result := m.processWarpPortal(args)
    
    // Notify through hooks system
    if m.hookManager != nil {
        m.hookManager.CallHook("field.warp_portal", result)
    }
    
    return nil
}

// Process warp_portal packet and return structured result
func (m *FieldManager) processWarpPortal(args map[string]interface{}) map[string]interface{} {
    // Process packet data and return structured result
    mapName, _ := args["map_name"].(string)
    
    return map[string]interface{}{
        "map_name": mapName,
        "processed_at": time.Now(),
    }
}
```

## Usage

### Creating a Receive Component

```go
// Create a hook manager
hookManager := hooks.NewHookManager()

// Create a receive factory
receiveFactory := factory.NewReceiveFactory()

// Register default server types
receiveFactory.RegisterDefaultServerTypes()

// Create a receive component for ServerType0
receive, err := receiveFactory.CreateReceive("ServerType0", hookManager)
if err != nil {
    // Handle error
}
```

### Registering Domain Managers

```go
// Create core parser
parser := core.NewCoreParser(receive)

// Create and register domain-specific managers
fieldManager := field.NewFieldManager(parser, hookManager)
fieldManager.RegisterHandlers()

actorManager := actor.NewActorManager(parser, hookManager)
actorManager.RegisterHandlers()

securityManager := security.NewSecurityManager(parser, hookManager)
securityManager.RegisterHandlers()
```

### Subscribing to Events

```go
// Subscribe to field.warp_portal events
hookManager.AddHook("field.warp_portal", func(args map[string]interface{}) {
    mapName := args["map_name"].(string)
    fmt.Printf("Warping to map: %s\n", mapName)
})
```

### Processing Packets

```go
// Process a packet
err := receive.Process(packet)
if err != nil {
    // Handle error
}
```

## Extending the Component

### Adding a New Server Type

1. Create a new packet definition provider function:

```go
func NewServerTypePacketDefs() map[string]common.PacketConstruction {
    // Start with base definitions
    defs := ServerType0PacketDefs()

    // Override or add specific packet definitions
    defs["0069"] = common.PacketConstruction{
        ID:         "0069",
        Name:       "account_server_info",
        Format:     "v a4 a4 a4 a4 a26 C a* v v",
        FieldNames: []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo", "additionalField1", "additionalField2"},
    }

    return defs
}
```

2. Register the new server type with the factory:

```go
receiveFactory.RegisterServerType("NewServerType", NewServerTypePacketDefs)
```

### Adding a New Domain Manager

1. Create a new package for the domain manager:

```go
// Package inventory provides handlers for inventory-related packets.
package inventory

import (
    "github.com/lenaxia/goKore/network/receive/core"
    "github.com/lenaxia/goKore/network/hooks"
)

// InventoryManager handles inventory-related packets
type InventoryManager struct {
    parser      *core.CoreParser
    hookManager *hooks.HookManager
}

// NewInventoryManager creates a new inventory manager
func NewInventoryManager(parser *core.CoreParser, hookManager *hooks.HookManager) *InventoryManager {
    return &InventoryManager{
        parser:      parser,
        hookManager: hookManager,
    }
}

// RegisterHandlers registers all inventory-related handlers
func (m *InventoryManager) RegisterHandlers() {
    m.parser.RegisterHandlerFunc("00A0", "inventory_item_add", "v3 C v C2 a8 l",
        []string{"index", "amount", "itemID", "identified", "broken", "upgrade", "cards", "expire_time"},
        m.handleInventoryItemAdd)
        
    // More inventory handlers...
}

// Handler implementations...
```

2. Create and register the domain manager:

```go
inventoryManager := inventory.NewInventoryManager(parser, hookManager)
inventoryManager.RegisterHandlers()
```

## Benefits of This Architecture

1. **Domain-Driven Organization**: Group related packets by domain (character, inventory, field, etc.)
2. **Improved Testability**: Domain-specific managers can be tested in isolation
3. **Consistent Hook Notifications**: Standard pattern for notifying other components
4. **Reduced Duplication**: Common patterns extracted to domain managers
5. **Simplicity**: Leverages existing hooks system without introducing new abstractions
6. **Consistent Architecture**: Both send and receive sides follow the same patterns
7. **Centralized Configuration**: Server-specific configurations in one place
8. **Better Scalability**: Easy to add new server types