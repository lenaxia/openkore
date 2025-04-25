# Comprehensive Implementation Plan for OpenKore Network Stack in Go

## Phase 1: Core Infrastructure (3 weeks)

1. **network/network.go**
   - Source: `Network.pm`
   - Dependencies: None
   - Purpose: Define core network interfaces and constants
   - Tasks:
     - Define connection states
     - Create base interfaces
     - Implement error types

2. **utils/crypto/crypton.go**
   - Source: `Utils/Crypton.pm`
   - Dependencies: None
   - Purpose: Implement encryption algorithms
   - Tasks:
     - Port Crypton algorithm to Go
     - Implement table generation
     - Create encryption/decryption methods

3. **network/protocol/tokenizer.go**
   - Source: `Network/MessageTokenizer.pm`
   - Dependencies: network/network.go
   - Purpose: Extract discrete packets from byte stream
   - Tasks:
     - Implement buffer management
     - Create packet boundary detection
     - Handle message types

4. **network/protocol/parser.go**
   - Source: `Network/PacketParser.pm`
   - Dependencies: network/protocol/tokenizer.go
   - Purpose: Parse and construct packet structures
   - Tasks:
     - Create packet parsing framework
     - Implement handler registration
     - Build packet construction utilities

5. **network/packets/definitions.go**
   - Source: `rpackets.txt`, packet structures from `Network/Receive/ServerType0.pm`
   - Dependencies: None
   - Purpose: Define packet structures and field layouts
   - Tasks:
     - Define packet structure types
     - Create field definition types
     - Implement packet lookup utilities

6. **network/packets/database.go**
   - Source: Packet loading code from various files
   - Dependencies: network/packets/definitions.go
   - Purpose: Load and manage packet definitions
   - Tasks:
     - Implement config file parsing
     - Create packet database structure
     - Build dynamic packet registration

## Phase 2: Connection Management (3 weeks)

1. **network/connection/connection.go**
   - Source: Connection interfaces from `Network/DirectConnection.pm`
   - Dependencies: network/network.go
   - Purpose: Define connection interfaces
   - Tasks:
     - Create connection interface
     - Define connection events
     - Implement connection states

2. **network/connection/direct.go**
   - Source: `Network/DirectConnection.pm`
   - Dependencies: network/connection/connection.go
   - Purpose: Implement direct TCP connection to server
   - Tasks:
     - Implement TCP connection handling
     - Create send/receive methods
     - Handle connection errors

3. **network/connection/tls.go**
   - Source: SSL/TLS concepts (not explicitly in provided files)
   - Dependencies: network/connection/connection.go
   - Purpose: Implement SSL/TLS connections
   - Tasks:
     - Implement TLS connection handling
     - Create certificate verification
     - Handle secure connections

4. **network/proxy/proxy.go**
   - Source: Proxy concepts from `Network/DirectConnection.pm`
   - Dependencies: network/connection/connection.go
   - Purpose: Define proxy interfaces
   - Tasks:
     - Create proxy interface
     - Define proxy configuration
     - Implement proxy selection

5. **network/proxy/socks.go**
   - Source: SOCKS proxy concepts (not explicitly in provided files)
   - Dependencies: network/proxy/proxy.go
   - Purpose: Implement SOCKS proxy support
   - Tasks:
     - Implement SOCKS4/5 protocol
     - Create authentication methods
     - Handle proxy negotiation

6. **network/proxy/http.go**
   - Source: HTTP proxy concepts (not explicitly in provided files)
   - Dependencies: network/proxy/proxy.go
   - Purpose: Implement HTTP proxy support
   - Tasks:
     - Implement HTTP CONNECT method
     - Handle proxy authentication
     - Create tunneling support

7. **network/connection/manager.go**
   - Source: Connection management logic from `Network/DirectConnection.pm`
   - Dependencies: network/connection/*, network/proxy/*
   - Purpose: Manage connection lifecycle and reconnection
   - Tasks:
     - Implement connection state machine
     - Create reconnection strategies
     - Handle connection events

## Phase 3: Protocol Implementation (3 weeks)

1. **network/protocol/padding.go**
   - Source: `Network/PaddedPackets.pm`
   - Dependencies: utils/crypto/crypton.go
   - Purpose: Handle special packet formatting for security
   - Tasks:
     - Implement packet padding algorithms
     - Create hash data management
     - Build packet generation methods

2. **network/packets/constructor.go**
   - Source: Packet construction methods from `Network/PacketParser.pm`
   - Dependencies: network/packets/definitions.go
   - Purpose: Construct outgoing packets
   - Tasks:
     - Implement packet field encoding
     - Create packet assembly methods
     - Build validation utilities

3. **network/hooks/event_hooks.go**
   - Source: Plugin hook system in `Network/PacketParser.pm`
   - Dependencies: network/packets/definitions.go
   - Purpose: Provide hook points for plugins
   - Tasks:
     - Create hook registration system
     - Implement hook calling mechanism
     - Build hook priority handling

4. **network/state/game_state.go**
   - Source: State management in `Network.pm` and `Network/DirectConnection.pm`
   - Dependencies: network/network.go
   - Purpose: Track and manage game connection state
   - Tasks:
     - Implement state machine
     - Create state transition events
     - Build state persistence

5. **network/state/session.go**
   - Source: Session management concepts from various files
   - Dependencies: network/state/game_state.go
   - Purpose: Manage game session data
   - Tasks:
     - Implement session creation
     - Create session data storage
     - Build session recovery

## Phase 4: Server Implementation (4 weeks)

1. **network/config/server_config.go**
   - Source: Server configuration handling in various files
   - Dependencies: None
   - Purpose: Load and manage server-specific configurations
   - Tasks:
     - Implement config loading
     - Create server type detection
     - Build configuration validation

2. **network/config/network_config.go**
   - Source: Network configuration from various files
   - Dependencies: network/config/server_config.go
   - Purpose: Configure network behavior
   - Tasks:
     - Implement timeout settings
     - Create reconnection policies
     - Build proxy configuration

3. **network/receive/core/parse.go**
   - Source: Core parsing functionality from `Network/PacketParser.pm`
   - Dependencies: network/protocol/parser.go
   - Purpose: Implement core packet parsing functionality
   - Tasks:
     - Create parsing framework
     - Implement packet validation
     - Build packet processing pipeline

4. **network/receive/core/account.go**
   - Source: Account-related core functionality from `Network/Receive/ServerType0.pm`
   - Dependencies: network/receive/core/parse.go
   - Purpose: Handle account-related core functionality
   - Tasks:
     - Implement account ID handling
     - Create session management
     - Build account state tracking

5. **network/receive/security/login.go**
   - Source: Login and authentication from `Network/Receive/ServerType0.pm`
   - Dependencies: network/receive/core/account.go
   - Purpose: Handle login and authentication
   - Tasks:
     - Implement login packet handling
     - Create authentication flow
     - Build session establishment

6. **network/receive/security/pin.go**
   - Source: PIN code handling from `Network/Receive/ServerType0.pm`
   - Dependencies: network/receive/security/login.go
   - Purpose: Handle PIN code authentication
   - Tasks:
     - Implement PIN code request handling
     - Create PIN code validation
     - Build PIN code change functionality

7. **network/receive/security/anticheat.go**
   - Source: Anti-cheat systems from `Network/Receive/ServerType0.pm`
   - Dependencies: network/receive/core/parse.go
   - Purpose: Handle anti-cheat systems
   - Tasks:
     - Implement GameGuard handling
     - Create captcha handling
     - Build other anti-cheat mechanisms

8. **network/receive/game/actor/player.go**
   - Source: Player-related functionality from `Network/Receive/ServerType0.pm`
   - Dependencies: network/receive/core/parse.go
   - Purpose: Handle player-related packets
   - Tasks:
     - Implement player appearance handling
     - Create player movement handling
     - Build player status handling

9. **network/receive/game/actor/monster.go**
   - Source: Monster-related functionality from `Network/Receive/ServerType0.pm`
   - Dependencies: network/receive/core/parse.go
   - Purpose: Handle monster-related packets
   - Tasks:
     - Implement monster appearance handling
     - Create monster movement handling
     - Build monster status handling

10. **network/receive/game/actor/npc.go**
    - Source: NPC-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle NPC-related packets
    - Tasks:
      - Implement NPC appearance handling
      - Create NPC interaction handling
      - Build NPC dialog handling

11. **network/receive/game/item/inventory.go**
    - Source: Inventory-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle inventory-related packets
    - Tasks:
      - Implement inventory item handling
      - Create item use handling
      - Build item movement handling

12. **network/receive/game/item/storage.go**
    - Source: Storage-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle storage-related packets
    - Tasks:
      - Implement storage item handling
      - Create storage open/close handling
      - Build item transfer handling

13. **network/receive/game/item/equipment.go**
    - Source: Equipment-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle equipment-related packets
    - Tasks:
      - Implement equipment change handling
      - Create equipment status handling
      - Build equipment upgrade handling

14. **network/receive/game/social/chat.go**
    - Source: Chat-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle chat-related packets
    - Tasks:
      - Implement public chat handling
      - Create private chat handling
      - Build party/guild chat handling

15. **network/receive/game/social/guild.go**
    - Source: Guild-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle guild-related packets
    - Tasks:
      - Implement guild info handling
      - Create guild member handling
      - Build guild action handling

16. **network/receive/game/social/party.go**
    - Source: Party-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle party-related packets
    - Tasks:
      - Implement party info handling
      - Create party member handling
      - Build party action handling

17. **network/receive/game/social/friends.go**
    - Source: Friend-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle friend-related packets
    - Tasks:
      - Implement friend list handling
      - Create friend request handling
      - Build friend status handling

18. **network/receive/game/economy/shop.go**
    - Source: Shop-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle shop-related packets
    - Tasks:
      - Implement NPC shop handling
      - Create buy/sell handling
      - Build shop item handling

19. **network/receive/game/economy/vending.go**
    - Source: Vending-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle vending-related packets
    - Tasks:
      - Implement vending creation handling
      - Create vending transaction handling
      - Build vending search handling

20. **network/receive/game/economy/auction.go**
    - Source: Auction-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle auction-related packets
    - Tasks:
      - Implement auction listing handling
      - Create auction bid handling
      - Build auction result handling

21. **network/receive/game/economy/banking.go**
    - Source: Banking-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle banking-related packets
    - Tasks:
      - Implement bank open/close handling
      - Create deposit/withdraw handling
      - Build transaction result handling

22. **network/receive/game/world/map.go**
    - Source: Map-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle map-related packets
    - Tasks:
      - Implement map change handling
      - Create map property handling
      - Build map effect handling

23. **network/receive/game/world/movement.go**
    - Source: Movement-related functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle movement-related packets
    - Tasks:
      - Implement character movement handling
      - Create teleport handling
      - Build movement restriction handling

24. **network/receive/game/world/npc.go**
    - Source: NPC interaction functionality from `Network/Receive/ServerType0.pm`
    - Dependencies: network/receive/core/parse.go
    - Purpose: Handle NPC interaction packets
    - Tasks:
      - Implement NPC dialog handling
      - Create NPC menu handling
      - Build NPC script handling

25. **network/receive/types/constants.go**
    - Source: Constants from various files
    - Dependencies: None
    - Purpose: Define shared constants
    - Tasks:
      - Define packet type constants
      - Create status effect constants
      - Build other shared constants

26. **network/receive/types/packets.go**
    - Source: Packet structures from various files
    - Dependencies: None
    - Purpose: Define shared packet structures
    - Tasks:
      - Define common packet structures
      - Create packet field types
      - Build packet utility functions

27. **network/servers/server.go**
   - Source: Server interface concepts from various files
   - Dependencies: network/receive/*, network/config/server_config.go
   - Purpose: Define server interface
   - Tasks:
     - Create server interface
     - Define server capabilities
     - Implement server events

28. **network/servers/base_server.go**
   - Source: `Network/Receive/ServerType0.pm`
   - Dependencies: network/servers/server.go, network/receive/*
   - Purpose: Implement base server functionality
   - Tasks:
     - Create base server implementation
     - Implement handler registration
     - Build packet processing pipeline

29. **network/servers/sakray_server.go**
   - Source: `Network/Receive/Sakray.pm`
   - Dependencies: network/servers/base_server.go
   - Purpose: Implement Sakray-specific functionality
   - Tasks:
     - Override base server methods
     - Implement Sakray-specific handlers
     - Create Sakray packet mappings

30. **network/servers/factory.go**
    - Source: Server creation logic from `Network/Receive.pm`
    - Dependencies: network/servers/*, network/config/server_config.go
    - Purpose: Create appropriate server implementation based on configuration
    - Tasks:
      - Implement server type detection
      - Create server instantiation
      - Build server configuration

## Phase 5: Packet Sending Implementation (3 weeks)

1. **network/send/send.go**
   - Source: Send interface concepts from various files
   - Dependencies: network/packets/constructor.go
   - Purpose: Define packet sending interface
   - Tasks:
     - Create send interface
     - Define sending capabilities
     - Implement sending events

2. **network/send/base_send.go**
   - Source: `Network/Send/ServerType0.pm` (inferred)
   - Dependencies: network/send/send.go, network/packets/constructor.go
   - Purpose: Implement base packet sending functionality
   - Tasks:
     - Create base send implementation
     - Implement common packet construction
     - Build sending methods for all packet types

3. **network/send/sakray_send.go**
   - Source: `Network/Send/Sakray.pm` (inferred)
   - Dependencies: network/send/base_send.go
   - Purpose: Implement Sakray-specific packet sending
   - Tasks:
     - Override base send methods
     - Implement Sakray-specific packet construction
     - Create Sakray packet mappings

4. **network/send/factory.go**
   - Source: Send creation logic (inferred)
   - Dependencies: network/send/*, network/config/server_config.go
   - Purpose: Create appropriate send implementation based on configuration
   - Tasks:
     - Implement send type detection
     - Create send instantiation
     - Build send configuration

## Phase 6: Integration and Testing (4 weeks)

1. **network/client.go**
   - Source: Integration of all components
   - Dependencies: All network/* packages
   - Purpose: Provide a unified client interface
   - Tasks:
     - Create client interface
     - Implement client lifecycle
     - Build event system
     - Integrate all components

2. **Unit Tests**
   - Source: All Perl files
   - Dependencies: All Go files
   - Purpose: Test individual components
   - Tasks:
     - Write tests for each component
     - Create mock objects
     - Implement test utilities

3. **Integration Tests**
   - Source: All Perl files
   - Dependencies: All Go files
   - Purpose: Test components working together
   - Tasks:
     - Create end-to-end tests
     - Implement test servers
     - Build test scenarios

4. **Performance Testing**
   - Source: All Perl files
   - Dependencies: All Go files
   - Purpose: Test performance and optimize
   - Tasks:
     - Create benchmarks
     - Identify bottlenecks
     - Implement optimizations

5. **Documentation**
   - Source: All Perl files
   - Dependencies: All Go files
   - Purpose: Document the Go implementation
   - Tasks:
     - Write package documentation
     - Create usage examples
     - Build API reference

## Timeline Summary

1. **Phase 1: Core Infrastructure** - 3 weeks
2. **Phase 2: Connection Management** - 3 weeks
3. **Phase 3: Protocol Implementation** - 3 weeks
4. **Phase 4: Server Implementation** - 4 weeks
5. **Phase 5: Packet Sending Implementation** - 3 weeks
6. **Phase 6: Integration and Testing** - 4 weeks

**Total Estimated Time: 20 weeks (5 months)**

This comprehensive implementation plan covers the entire network stack for OpenKore (excluding XKore), including all the additional components needed for a complete implementation. The plan is structured to build from the core infrastructure up to the complete client interface, with each phase building on the previous ones.


# Comprehensive OpenKore Network Stack Implementation in Go

## Complete Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                                                 │
│                                              OpenKore Core                                                      │
│                                                                                                                 │
└───────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────┘
                                                        │
                                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                             network/client.go                                                   │
│                                       (Unified Network Client Interface)                                        │
└───────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────┘
                                                        │
                  ┌────────────────────────────────────┬┴┬────────────────────────────────────────┐
                  │                                    │ │                                        │
                  ▼                                    │ │                                        ▼
┌─────────────────────────────────────┐               │ │               ┌─────────────────────────────────────────┐
│        network/state/*.go           │               │ │               │         network/config/*.go             │
│  (Game State & Session Management)  │◄──────────────┘ └──────────────►│    (Server & Network Configuration)     │
└─────────────────┬───────────────────┘                                 └─────────────────────┬───────────────────┘
                  │                                                                           │
                  │                                                                           │
                  ▼                                                                           ▼
┌─────────────────────────────────────┐               ┌─────────────────────────────────────────────────────────┐
│      network/servers/*.go           │◄──────────────┤                network/send/*.go                        │
│  (Server-specific Implementations)  │               │         (Packet Sending Implementation)                 │
└─────────────────┬───────────────────┘               └─────────────────────┬───────────────────────────────────┘
                  │                                                         │
                  ▼                                                         ▼
┌─────────────────────────────────────┐               ┌─────────────────────────────────────────────────────────┐
│      network/receive/*.go           │               │            network/packets/*.go                         │
│     (Packet Handler Implementation) │               │      (Packet Definitions & Construction)                │
└─────────────────┬───────────────────┘               └─────────────────────┬───────────────────────────────────┘
                  │                                                         │
                  └─────────────────────┐   ┌─────────────────────────────┬─┘
                                        │   │                             │
                                        ▼   ▼                             ▼
                  ┌─────────────────────────────────────┐     ┌─────────────────────────────────────┐
                  │      network/protocol/*.go          │     │       network/hooks/event_hooks.go   │
                  │  (Protocol Parsing & Processing)    │     │         (Plugin Hook System)         │
                  └─────────────────────┬───────────────┘     └─────────────────────────────────────┘
                                        │
                                        ▼
                  ┌─────────────────────────────────────┐
                  │      network/connection/*.go        │
                  │    (Connection Management)          │
                  └─────────────────────┬───────────────┘
                                        │
                                        ▼
                  ┌─────────────────────────────────────┐
                  │        network/proxy/*.go           │
                  │      (Proxy Support)                │
                  └─────────────────────┬───────────────┘
                                        │
                                        ▼
                  ┌─────────────────────────────────────┐
                  │       utils/crypto/*.go             │
                  │    (Encryption Support)             │
                  └─────────────────────┬───────────────┘
                                        │
                                        ▼
                  ┌─────────────────────────────────────┐
                  │             RO Server               │
                  └─────────────────────────────────────┘
```

## Updated Architecture Diagram with Detailed Receive Structure

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                              OpenKore Core                                                      │
└───────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────┘
                                                        │
                                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                             network/client.go                                                   │
└───────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────┘
                                                        │
                  ┌────────────────────────────────────┬┴┬────────────────────────────────────────┐
                  │                                    │ │                                        │
                  ▼                                    │ │                                        ▼
┌─────────────────────────────────────┐               │ │               ┌─────────────────────────────────────────┐
│        network/state/*.go           │               │ │               │         network/config/*.go             │
└─────────────────┬───────────────────┘               │ │               └─────────────────────┬───────────────────┘
                  │                                    │ │                                    │
                  │                                    │ │                                    │
                  ▼                                    │ │                                    ▼
┌─────────────────────────────────────┐               │ │               ┌─────────────────────────────────────────┐
│      network/servers/*.go           │◄──────────────┘ └──────────────►│         network/send/*.go               │
└─────────────────┬───────────────────┘                                 └─────────────────────┬───────────────────┘
                  │                                                                           │
                  ▼                                                                           ▼
┌─────────────────────────────────────────────────────────────────────┐ ┌─────────────────────────────────────────┐
│                     network/receive/                                │ │      network/packets/*.go               │
│  ┌─────────────────────────────────────────────────────────────┐    │ └─────────────────────┬───────────────────┘
│  │                        core/                                │    │                       │
│  │  ┌─────────────┐  ┌─────────────┐                          │    │                       │
│  │  │  parse.go   │  │ account.go  │                          │    │                       │
│  │  └─────────────┘  └─────────────┘                          │    │                       │
│  └─────────────────────────────────────────────────────────────┘    │                       │
│                                                                     │                       │
│  ┌─────────────────────────────────────────────────────────────┐    │                       │
│  │                      security/                              │    │                       │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐      │    │                       │
│  │  │  login.go   │  │   pin.go    │  │  anticheat.go   │      │    │                       │
│  │  └─────────────┘  └─────────────┘  └─────────────────┘      │    │                       │
│  └─────────────────────────────────────────────────────────────┘    │                       │
│                                                                     │                       │
│  ┌─────────────────────────────────────────────────────────────┐    │                       │
│  │                        game/                                │    │                       │
│  │  ┌─────────────────────────────────────────────────────┐    │    │                       │
│  │  │                      actor/                         │    │    │                       │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │    │    │                       │
│  │  │  │  player.go  │  │ monster.go  │  │   npc.go    │  │    │    │                       │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘  │    │    │                       │
│  │  └─────────────────────────────────────────────────────┘    │    │                       │
│  │                                                             │    │                       │
│  │  ┌─────────────────────────────────────────────────────┐    │    │                       │
│  │  │                      item/                          │    │    │                       │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │    │    │                       │
│  │  │  │inventory.go │  │ storage.go  │  │equipment.go │  │    │    │                       │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘  │    │    │                       │
│  │  └─────────────────────────────────────────────────────┘    │    │                       │
│  │                                                             │    │                       │
│  │  ┌─────────────────────────────────────────────────────┐    │    │                       │
│  │  │                     social/                         │    │    │                       │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │    │    │                       │
│  │  │  │   chat.go   │  │  guild.go   │  │  party.go   │  │    │    │                       │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘  │    │    │                       │
│  │  │  ┌─────────────┐                                    │    │    │                       │
│  │  │  │ friends.go  │                                    │    │    │                       │
│  │  │  └─────────────┘                                    │    │    │                       │
│  │  └─────────────────────────────────────────────────────┘    │    │                       │
│  │                                                             │    │                       │
│  │  ┌─────────────────────────────────────────────────────┐    │    │                       │
│  │  │                    economy/                         │    │    │                       │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │    │    │                       │
│  │  │  │   shop.go   │  │ vending.go  │  │ auction.go  │  │    │    │                       │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘  │    │    │                       │
│  │  │  ┌─────────────┐                                    │    │    │                       │
│  │  │  │ banking.go  │                                    │    │    │                       │
│  │  │  └─────────────┘                                    │    │    │                       │
│  │  └─────────────────────────────────────────────────────┘    │    │                       │
│  │                                                             │    │                       │
│  │  ┌─────────────────────────────────────────────────────┐    │    │                       │
│  │  │                     world/                          │    │    │                       │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │    │    │                       │
│  │  │  │   map.go    │  │ movement.go │  │   npc.go    │  │    │    │                       │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘  │    │    │                       │
│  │  └─────────────────────────────────────────────────────┘    │    │                       │
│  └─────────────────────────────────────────────────────────────┘    │                       │
│                                                                     │                       │
│  ┌─────────────────────────────────────────────────────────────┐    │                       │
│  │                       types/                                │    │                       │
│  │  ┌─────────────────┐  ┌─────────────────┐                   │    │                       │
│  │  │  constants.go   │  │   packets.go    │                   │    │                       │
│  │  └─────────────────┘  └─────────────────┘                   │    │                       │
│  └─────────────────────────────────────────────────────────────┘    │                       │
└─────────────────────────────────────────────────────────────────────┘                       │
                                                                                              │
                  ┌─────────────────────────────────────┐                                     │
                  │      network/protocol/*.go          │◄────────────────────────────────────┘
                  └─────────────────────┬───────────────┘
                                        │
                                        ▼
                  ┌─────────────────────────────────────┐
                  │      network/connection/*.go        │
                  └─────────────────────┬───────────────┘
                                        │
                                        ▼
                  ┌─────────────────────────────────────┐
                  │        network/proxy/*.go           │
                  └─────────────────────┬───────────────┘
                                        │
                                        ▼
                  ┌─────────────────────────────────────┐
                  │       utils/crypto/*.go             │
                  └─────────────────────┬───────────────┘
                                        │
                                        ▼
                  ┌─────────────────────────────────────┐
                  │             RO Server               │
                  └─────────────────────────────────────┘
```

## Updated Code Snippets

### 1. network/network.go

```go
package network

import "errors"

// Connection states
const (
    NotConnected           = 1
    ConnectedToMasterServer = 2
    ConnectedToLoginServer  = 3
    ConnectedToCharServer   = 4
    InGame                  = 5
    InGameButUninitialized  = -1
)

// Common errors
var (
    ErrNotConnected = errors.New("not connected to server")
    ErrInvalidState = errors.New("invalid connection state")
    ErrTimeout      = errors.New("connection timed out")
)

// NetworkInterface defines the common interface for all network implementations
type NetworkInterface interface {
    Connect() error
    Disconnect() error
    IsConnected() bool
    GetState() int
    SetState(state int)
    Send(data []byte) error
    Receive() ([]byte, error)
}

// StateChangeCallback is called when connection state changes
type StateChangeCallback func(oldState, newState int)
```

### 2. network/connection/connection.go

```go
package connection

import (
    "context"
    "net"
    "time"

    "github.com/yourusername/openkore-go/network"
)

// Connection defines the interface for all connection types
type Connection interface {
    Connect(ctx context.Context, host string, port int) error
    Disconnect() error
    IsConnected() bool
    Send(data []byte) error
    Receive() ([]byte, error)
    SetTimeout(timeout time.Duration)
    GetRemoteAddress() string
}

// ConnectionEventType defines types of connection events
type ConnectionEventType int

const (
    EventConnected ConnectionEventType = iota
    EventDisconnected
    EventDataSent
    EventDataReceived
    EventError
)

// ConnectionEvent represents a connection-related event
type ConnectionEvent struct {
    Type    ConnectionEventType
    Data    []byte
    Error   error
    Address string
}

// ConnectionEventHandler handles connection events
type ConnectionEventHandler func(event ConnectionEvent)
```

### 3. network/connection/direct.go

```go
package connection

import (
    "context"
    "net"
    "time"

    "github.com/yourusername/openkore-go/network"
)

// DirectConnection implements a direct TCP connection to the server
type DirectConnection struct {
    socket      net.Conn
    state       int
    serverHost  string
    serverPort  int
    timeout     time.Duration
    lastReceive time.Time
    handlers    []ConnectionEventHandler
}

// NewDirectConnection creates a new direct connection
func NewDirectConnection() *DirectConnection {
    return &DirectConnection{
        state:   network.NotConnected,
        timeout: 30 * time.Second,
        handlers: make([]ConnectionEventHandler, 0),
    }
}

// Connect establishes a connection to the server
func (c *DirectConnection) Connect(ctx context.Context, host string, port int) error {
    if c.IsConnected() {
        c.Disconnect()
    }

    c.serverHost = host
    c.serverPort = port

    dialer := net.Dialer{Timeout: c.timeout}
    var err error
    c.socket, err = dialer.DialContext(ctx, "tcp", net.JoinHostPort(host, fmt.Sprintf("%d", port)))

    if err != nil {
        c.emitEvent(EventError, nil, err, "")
        return err
    }

    c.state = network.ConnectedToMasterServer
    c.emitEvent(EventConnected, nil, nil, c.GetRemoteAddress())
    return nil
}

// Send sends data to the server
func (c *DirectConnection) Send(data []byte) error {
    if !c.IsConnected() {
        return network.ErrNotConnected
    }

    if err := c.socket.SetWriteDeadline(time.Now().Add(c.timeout)); err != nil {
        return err
    }

    _, err := c.socket.Write(data)
    if err != nil {
        c.emitEvent(EventError, data, err, "")
        return err
    }

    c.emitEvent(EventDataSent, data, nil, "")
    return nil
}

// Receive receives data from the server
func (c *DirectConnection) Receive() ([]byte, error) {
    if !c.IsConnected() {
        return nil, network.ErrNotConnected
    }

    if err := c.socket.SetReadDeadline(time.Now().Add(c.timeout)); err != nil {
        return nil, err
    }

    buffer := make([]byte, 4096)
    n, err := c.socket.Read(buffer)
    if err != nil {
        if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
            return nil, network.ErrTimeout
        }
        c.emitEvent(EventError, nil, err, "")
        return nil, err
    }

    data := buffer[:n]
    c.lastReceive = time.Now()
    c.emitEvent(EventDataReceived, data, nil, "")
    return data, nil
}

// Disconnect closes the connection
func (c *DirectConnection) Disconnect() error {
    if c.socket != nil {
        err := c.socket.Close()
        c.socket = nil
        c.state = network.NotConnected
        c.emitEvent(EventDisconnected, nil, err, "")
        return err
    }
    return nil
}

// IsConnected checks if the connection is active
func (c *DirectConnection) IsConnected() bool {
    return c.socket != nil && c.state != network.NotConnected
}

// GetRemoteAddress returns the remote server address
func (c *DirectConnection) GetRemoteAddress() string {
    if c.socket != nil {
        return c.socket.RemoteAddr().String()
    }
    return ""
}

// SetTimeout sets the connection timeout
func (c *DirectConnection) SetTimeout(timeout time.Duration) {
    c.timeout = timeout
}

// AddEventHandler adds a handler for connection events
func (c *DirectConnection) AddEventHandler(handler ConnectionEventHandler) {
    c.handlers = append(c.handlers, handler)
}

// emitEvent notifies all handlers of an event
func (c *DirectConnection) emitEvent(eventType ConnectionEventType, data []byte, err error, address string) {
    event := ConnectionEvent{
        Type:    eventType,
        Data:    data,
        Error:   err,
        Address: address,
    }

    for _, handler := range c.handlers {
        handler(event)
    }
}
```

### 4. network/protocol/tokenizer.go

```go
package protocol

import (
    "encoding/binary"
    "errors"
    "fmt"
)

// MessageType represents the type of message received
type MessageType int

const (
    KnownMessage MessageType = iota
    UnknownMessage
    AccountID
)

// Errors
var (
    ErrIncompletePacket = errors.New("incomplete packet")
    ErrInvalidPacket    = errors.New("invalid packet")
)

// PacketDef contains information about packet structure
type PacketDef struct {
    Length    int  // Fixed length or -1 for variable length
    HasLength bool // Whether packet has length field
}

// Tokenizer handles breaking byte streams into discrete packets
type Tokenizer struct {
    buffer     []byte
    packetDefs map[string]PacketDef
    nextMightBeAccountID bool
}

// NewTokenizer creates a new message tokenizer
func NewTokenizer(packetDefs map[string]PacketDef) *Tokenizer {
    return &Tokenizer{
        buffer:     make([]byte, 0),
        packetDefs: packetDefs,
    }
}

// Add appends data to the buffer
func (t *Tokenizer) Add(data []byte) {
    t.buffer = append(t.buffer, data...)
}

// Clear removes data from the buffer
func (t *Tokenizer) Clear(size int) {
    if size <= 0 || size > len(t.buffer) {
        t.buffer = make([]byte, 0)
    } else {
        t.buffer = t.buffer[size:]
    }
}

// GetBuffer returns the current buffer contents
func (t *Tokenizer) GetBuffer() []byte {
    return t.buffer
}

// NextMessageMightBeAccountID marks that the next message might be an account ID
func (t *Tokenizer) NextMessageMightBeAccountID() {
    t.nextMightBeAccountID = true
}

// GetMessageID extracts the message ID from a packet
func GetMessageID(packet []byte) string {
    if len(packet) < 2 {
        return ""
    }
    return fmt.Sprintf("%02X%02X", packet[1], packet[0])
}

// ReadNext extracts the next complete packet from the buffer
func (t *Tokenizer) ReadNext() ([]byte, MessageType, error) {
    if len(t.buffer) < 2 {
        return nil, UnknownMessage, ErrIncompletePacket
    }

    // Check if this might be an account ID
    if t.nextMightBeAccountID && len(t.buffer) >= 4 {
        t.nextMightBeAccountID = false
        // Logic to check if this is an account ID
        // For now, just a placeholder
        // In real implementation, compare with global accountID
        return t.buffer[:4], AccountID, nil
    }

    // Get packet ID and look up definition
    packetID := GetMessageID(t.buffer)
    packetDef, exists := t.packetDefs[packetID]

    if !exists {
        // Unknown packet type
        result := make([]byte, len(t.buffer))
        copy(result, t.buffer)
        t.buffer = make([]byte, 0)
        return result, UnknownMessage, nil
    }

    if packetDef.Length > 0 {
        // Fixed length packet
        if len(t.buffer) < packetDef.Length {
            return nil, UnknownMessage, ErrIncompletePacket
        }

        result := make([]byte, packetDef.Length)
        copy(result, t.buffer[:packetDef.Length])
        t.buffer = t.buffer[packetDef.Length:]
        return result, KnownMessage, nil
    } else if packetDef.HasLength {
        // Variable length packet
        if len(t.buffer) < 4 {
            return nil, UnknownMessage, ErrIncompletePacket
        }

        length := int(binary.LittleEndian.Uint16(t.buffer[2:4]))
        if len(t.buffer) < length {
            return nil, UnknownMessage, ErrIncompletePacket
        }

        result := make([]byte, length)
        copy(result, t.buffer[:length])
        t.buffer = t.buffer[length:]
        return result, KnownMessage, nil
    }

    // Should never reach here if packet definitions are correct
    return nil, UnknownMessage, ErrInvalidPacket
}
```

### 5. network/protocol/parser.go

```go
package protocol

import (
    "encoding/binary"
    "errors"
    "fmt"
    "reflect"

    "github.com/yourusername/openkore-go/network/packets"
)

// Errors
var (
    ErrUnknownPacket = errors.New("unknown packet type")
    ErrParsingFailed = errors.New("packet parsing failed")
)

// Parser handles parsing and constructing network packets
type Parser struct {
    packetDefs    map[string]packets.PacketDef
    packetHandlers map[string]PacketHandler
}

// PacketHandler defines the interface for packet handlers
type PacketHandler interface {
    Handle(packet []byte) error
}

// NewParser creates a new packet parser
func NewParser(packetDefs map[string]packets.PacketDef) *Parser {
    return &Parser{
        packetDefs:    packetDefs,
        packetHandlers: make(map[string]PacketHandler),
    }
}

// RegisterHandler registers a handler for a specific packet type
func (p *Parser) RegisterHandler(packetID string, handler PacketHandler) {
    p.packetHandlers[packetID] = handler
}

// Parse parses a raw packet and calls the appropriate handler
func (p *Parser) Parse(data []byte) error {
    if len(data) < 2 {
        return ErrParsingFailed
    }

    packetID := GetMessageID(data)
    handler, exists := p.packetHandlers[packetID]
    if !exists {
        return fmt.Errorf("%w: %s", ErrUnknownPacket, packetID)
    }

    return handler.Handle(data)
}

// Construct builds a packet from structured data
func (p *Parser) Construct(packetID string, fields map[string]interface{}) ([]byte, error) {
    def, exists := p.packetDefs[packetID]
    if !exists {
        return nil, fmt.Errorf("%w: %s", ErrUnknownPacket, packetID)
    }

    // Determine packet size
    size := def.Length
    if size <= 0 {
        // Variable length packet, start with 4 bytes (ID + length)
        size = 4
        for _, field := range def.Fields {
            if field.Type == "string" {
                // For strings, get actual length from the field value
                if strVal, ok := fields[field.Name].(string); ok {
                    size += len(strVal)
                } else {
                    size += field.Length // Default length
                }
            } else {
                size += field.Length
            }
        }
    }

    // Create packet buffer
    packet := make([]byte, size)

    // Set packet ID
    idBytes, err := hex.DecodeString(packetID)
    if err != nil {
        return nil, err
    }
    packet[0] = idBytes[1] // Little endian
    packet[1] = idBytes[0]

    // Set packet length for variable length packets
    if def.Length <= 0 {
        binary.LittleEndian.PutUint16(packet[2:4], uint16(size))
    }

    // Fill in fields
    offset := 2
    if def.Length <= 0 {
        offset = 4 // Skip length field for variable length packets
    }

    for _, field := range def.Fields {
        value, exists := fields[field.Name]
        if !exists {
            continue // Skip fields not provided
        }

        switch field.Type {
        case "byte":
            if v, ok := value.(byte); ok {
                packet[offset] = v
            }
            offset += 1
        case "short":
            if v, ok := value.(uint16); ok {
                binary.LittleEndian.PutUint16(packet[offset:offset+2], v)
            }
            offset += 2
        case "int":
            if v, ok := value.(uint32); ok {
                binary.LittleEndian.PutUint32(packet[offset:offset+4], v)
            }
            offset += 4
        case "string":
            if v, ok := value.(string); ok {
                copy(packet[offset:offset+len(v)], v)
            }
            offset += field.Length
        // Add more types as needed
        }
    }

    return packet, nil
}
```

### 6. network/client.go

```go
package network

import (
    "context"
    "sync"
    "time"

    "github.com/yourusername/openkore-go/network/connection"
    "github.com/yourusername/openkore-go/network/config"
    "github.com/yourusername/openkore-go/network/protocol"
    "github.com/yourusername/openkore-go/network/servers"
    "github.com/yourusername/openkore-go/network/send"
    "github.com/yourusername/openkore-go/network/state"
)

// Client represents the complete network client
type Client struct {
    conn          connection.Connection
    tokenizer     *protocol.Tokenizer
    parser        *protocol.Parser
    server        servers.Server
    sender        send.Sender
    config        *config.NetworkConfig
    gameState     *state.GameState
    session       *state.Session

    ctx           context.Context
    cancel        context.CancelFunc
    mutex         sync.Mutex

    stateCallbacks []StateChangeCallback
}

// NewClient creates a new network client
func NewClient(config *config.NetworkConfig) *Client {
    ctx, cancel := context.WithCancel(context.Background())

    client := &Client{
        config:        config,
        ctx:           ctx,
        cancel:        cancel,
        stateCallbacks: make([]StateChangeCallback, 0),
    }

    // Create connection based on config
    if config.UseTLS {
        client.conn = connection.NewTLSConnection()
    } else {
        client.conn = connection.NewDirectConnection()
    }

    // Set up proxy if configured
    if config.UseProxy {
        // Configure proxy
    }

    // Create game state and session
    client.gameState = state.NewGameState()
    client.session = state.NewSession()

    // Create server and sender based on server type
    serverFactory := servers.NewServerFactory()
    client.server = serverFactory.CreateServer(config.ServerType)

    senderFactory := send.NewSenderFactory()
    client.sender = senderFactory.CreateSender(config.ServerType)

    // Set up packet definitions
    packetDefs := loadPacketDefinitions(config.ServerType)
    client.tokenizer = protocol.NewTokenizer(packetDefs)
    client.parser = protocol.NewParser(packetDefs)

    // Set up connection event handlers
    client.conn.AddEventHandler(client.handleConnectionEvent)

    return client
}

// Connect connects to the server
func (c *Client) Connect() error {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    return c.conn.Connect(c.ctx, c.config.ServerHost, c.config.ServerPort)
}

// Disconnect disconnects from the server
func (c *Client) Disconnect() error {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    return c.conn.Disconnect()
}

// IsConnected checks if connected to the server
func (c *Client) IsConnected() bool {
    return c.conn.IsConnected()
}

// GetState returns the current connection state
func (c *Client) GetState() int {
    return c.gameState.GetState()
}

// SetState sets the connection state
func (c *Client) SetState(state int) {
    oldState := c.gameState.GetState()
    c.gameState.SetState(state)

    // Notify callbacks
    for _, callback := range c.stateCallbacks {
        callback(oldState, state)
    }
}

// AddStateChangeCallback adds a callback for state changes
func (c *Client) AddStateChangeCallback(callback StateChangeCallback) {
    c.stateCallbacks = append(c.stateCallbacks, callback)
}

// Start starts the client processing
func (c *Client) Start() {
    go c.processPackets()
}

// Stop stops the client processing
func (c *Client) Stop() {
    c.cancel()
    c.Disconnect()
}

// Send sends a packet to the server
func (c *Client) Send(packetName string, fields map[string]interface{}) error {
    packet, err := c.sender.Send(packetName, fields)
    if err != nil {
        return err
    }

    return c.conn.Send(packet)
}

// processPackets processes incoming packets
func (c *Client) processPackets() {
    buffer := make([]byte, 4096)

    for {
        select {
        case <-c.ctx.Done():
            return
        default:
            // Receive data
            n, err := c.conn.Receive()
            if err != nil {
                // Handle error
                continue
            }

            // Add to tokenizer
            c.tokenizer.Add(n)

            // Process all complete packets
            for {
                packet, msgType, err := c.tokenizer.ReadNext()
                if err != nil {
                    break // No more complete packets
                }

                switch msgType {
                case protocol.KnownMessage:
                    c.parser.Parse(packet)
                case protocol.AccountID:
                    // Handle account ID
                case protocol.UnknownMessage:
                    // Handle unknown message
                }
            }
        }
    }
}

// handleConnectionEvent handles connection events
func (c *Client) handleConnectionEvent(event connection.ConnectionEvent) {
    switch event.Type {
    case connection.EventConnected:
        // Handle connection
    case connection.EventDisconnected:
        // Handle disconnection
    case connection.EventError:
        // Handle error
    }
}

// loadPacketDefinitions loads packet definitions for the given server type
func loadPacketDefinitions(serverType string) map[string]protocol.PacketDef {
    // Load packet definitions from configuration
    return nil // Placeholder
}
```

These code snippets provide a comprehensive foundation for implementing the OpenKore network stack in Go. They demonstrate the key components and their interactions, following Go's idioms and best practices while maintaining compatibility with the original Perl implementation.

