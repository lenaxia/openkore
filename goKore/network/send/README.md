# Send Package

This package provides functionality for sending packets to the server in the goKore network stack.

## Overview

The send package consists of:

1. A central registry for managing all handler registrations
2. Core functionality for sending packets
3. Handlers for different types of packets (login, game, server-specific)
4. A unified registration pattern for standardizing handler registration

## HandlerRegistry

The `HandlerRegistry` is the central component that manages all handler registrations. It provides methods for registering different types of handlers:

- `RegisterAllHandlers`: Registers all handlers
- `registerLoginHandlers`: Registers login-related handlers
- `registerGameHandlers`: Registers game-related handlers
- `registerServerHandlers`: Registers server-specific handlers

## Current Registration Status

Currently, most handlers are registered through the `game.RegisterHandlers` function in the `handlers/game` package. This function registers handlers for:

- Character-related packets
- Pet-related packets
- Mercenary-related packets
- Battle-related packets
- Marriage-related packets
- Auction-related packets
- Buying store-related packets
- UI-related packets
- Deal-related packets
- Ranking-related packets
- GM-related packets
- Macro-related packets
- Captcha-related packets
- Card-related packets
- Cash shop-related packets
- Miscellaneous packets

However, many of these registrations are just placeholders and don't actually register any handlers. The actual handlers are implemented in the corresponding packages in the `game/` directory, but they don't have registration functions.

Only the `card` package has been migrated to the unified registration pattern, which provides a more standardized way of registering handlers.

## Usage

### Creating a HandlerRegistry

```go
// Create a BaseSend instance
baseSend := core.NewBaseSend(hookManager)

// Create a logger
logger := // your logger implementation

// Create a HandlerRegistry
registry := NewHandlerRegistry(baseSend, hookManager, logger)
```

### Registering Handlers

```go
// Register all handlers
registry.RegisterAllHandlers()
```

### Configuring Server Type

```go
// Get packet definitions
packetDefs := registry.GetPacketDefinitions()

// Configure server type
err := registry.ConfigureServerType("ServerType0", packetDefs)
if err != nil {
    // Handle error
}
```

## Unified Registration Pattern

The send package is transitioning to a unified registration pattern for standardizing how handlers are registered across different packages. This pattern consists of:

1. A `register.go` file that exposes standard registration functions:
   - `RegisterWithSend`: Registers all handlers with the send component
   - `GetPacketDefinitions`: Returns packet definitions for this package

2. A package-specific implementation file that contains the actual handler logic:
   - `Manager`: Handles packet sending
   - `RegisterHandlers`: Registers all handlers
   - Handler functions for each packet type

### Migrated Packages

The following packages have been migrated to the unified pattern:

- `card`: Card-related functionality

### Legacy Packages

Most packages still use the legacy registration pattern or don't have registration functions at all.

See [LEGACY_PACKAGES.md](LEGACY_PACKAGES.md) for a list of packages that need to be migrated to the unified pattern.

## Migration Guide

See [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) for instructions on how to migrate packages to the unified pattern.

## Template Package

A template package is available in `game/template` to help with creating new packages that follow the unified pattern. See [game/template/README.md](game/template/README.md) for more information.

## Example

An example of how to use the registry is available in the `examples/registry_usage` directory.