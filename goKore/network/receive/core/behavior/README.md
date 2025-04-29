# Character Behavior Package

This package provides handlers for character behavior-related packets in the OpenKore network stack. It follows the domain-driven organization pattern described in the network package README.md.

## Overview

The package implements handlers for the following packets:

- `manner_message` (0149): Handles manner point notifications
- `hack_shield_alarm` (08B3): Handles hack shield alarm notifications

## Architecture

The package uses a domain-driven approach where related packet handlers are grouped together in a domain-specific manager:

```
CharacterBehaviorManager
├── handleMannerMessage
└── handleHackShieldAlarm
```

Each handler:
1. Processes the packet data
2. Returns structured results
3. Notifies other components through the hooks system

## Usage

There are two ways to use this package:

### 1. Direct Registration with CoreParser

```go
// Create a hook manager
hookManager := hooks.NewHookManager()

// Create a core parser
parser := core.NewCoreParser("ServerType0", hookManager)

// Register behavior handlers
behavior.RegisterWithParser(parser, hookManager)
```

### 2. Registration with BaseReceive

```go
// Create a hook manager
hookManager := hooks.NewHookManager()

// Create a base receive
baseReceive := core.NewBaseReceive(hookManager)

// Configure the base receive
// ...

// Register behavior handlers
behavior.RegisterWithBaseReceive(baseReceive)
```

## Hook Events

The package publishes the following hook events:

- `character.manner_message`: Published when a manner message is received
  - Data: `map[string]interface{}{"flag": uint8, "message": string}`

- `character.hack_shield_alarm`: Published when a hack shield alarm is received
  - Data: `map[string]interface{}{"message": string}`

## Example

See the `example/main.go` file for a complete example of how to use this package.