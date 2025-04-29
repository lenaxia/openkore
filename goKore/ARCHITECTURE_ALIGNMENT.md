# Aligning Send and Receive Architectures

Based on our analysis of the original Perl implementation and the current Go code, here's what would be required to align both sides of the network architecture using the hybrid approach:

## Aligning the Receive Side

1. **Create a Packet Definition System**
   - Define a `PacketDef` structure to hold packet metadata (ID, name, format, field names)
   - Implement server-specific packet definition providers
   - Create inheritance mechanism for packet definitions (base definitions with overrides)

2. **Implement a Configurable Receive Component**
   ```go
   // Receive interface
   type Receive interface {
       RegisterHandler(packetName string, handler ReceiveHandler)
       Process(packet []byte) error
       Configure(serverType string, packetDefs map[string]PacketDef) error
   }
   
   // BaseReceive implementation
   type BaseReceive struct {
       serverType  string
       packetDefs  map[string]PacketDef
       handlers    map[string]ReceiveHandler
       hookManager *hooks.HookManager
   }
   ```

3. **Create a Factory for Receive Components**
   ```go
   // ReceiveFactory creates and configures Receive implementations
   type ReceiveFactory struct {
       packetDefProviders map[string]PacketDefProvider
   }
   
   // CreateReceive creates a configured Receive component
   func (rf *ReceiveFactory) CreateReceive(serverType string, hookManager *hooks.HookManager) (Receive, error) {
       // Get packet definitions for server type
       // Create and configure a BaseReceive
       // Return the configured component
   }
   ```

4. **Organize Handlers by Domain**
   ```go
   // Register all login-related handlers
   func RegisterLoginHandlers(receive Receive) {
       receive.RegisterHandler("account_server_info", handleAccountServerInfo)
       receive.RegisterHandler("login_error", handleLoginError)
       // More login handlers...
   }
   ```

## Aligning the Send Side

1. **Refine the Packet Construction System**
   - Define a `PacketConstruction` structure similar to `PacketDef`
   - Implement server-specific packet construction providers
   - Create inheritance mechanism for packet constructions

2. **Update the Send Component**
   ```go
   // Send interface (already exists but may need refinement)
   type Send interface {
       ConstructPacket(packetName string, args map[string]interface{}) ([]byte, error)
       SendPacket(packetName string, args map[string]interface{}) error
       Configure(serverType string, packetConstructions map[string]PacketConstruction) error
   }
   
   // BaseSend implementation (update existing)
   type BaseSend struct {
       serverType          string
       packetConstructions map[string]PacketConstruction
       connection          network.Connection
       hookManager         *hooks.HookManager
   }
   ```

3. **Refine the Factory for Send Components**
   ```go
   // SendFactory (update existing)
   type SendFactory struct {
       packetConstructionProviders map[string]PacketConstructionProvider
   }
   
   // CreateSend creates a configured Send component
   func (sf *SendFactory) CreateSend(serverType string, hookManager *hooks.HookManager) (Send, error) {
       // Get packet constructions for server type
       // Create and configure a BaseSend
       // Return the configured component
   }
   ```

4. **Organize Packet Construction by Domain**
   ```go
   // Register all login-related packet constructions
   func RegisterLoginPacketConstructors(send Send) {
       // Register login packet constructors
   }
   ```

## Implementation Steps

### Phase 1: Core Infrastructure (2-3 weeks)
1. Define common interfaces and structures
2. Implement the packet definition system
3. Create the base receive component
4. Update the base send component

### Phase 2: Factory Implementation (1-2 weeks)
1. Implement the receive factory
2. Update the send factory
3. Create server type providers

### Phase 3: Handler Organization (2-3 weeks)
1. Group handlers by domain
2. Create registration functions
3. Implement domain-specific modules

### Phase 4: Migration and Testing (2-3 weeks)
1. Migrate existing code to the new architecture
2. Create comprehensive tests
3. Validate with real server connections

## Key Benefits of Alignment

1. **Consistent Architecture**: Both send and receive sides follow the same patterns
2. **Centralized Configuration**: Server-specific configurations in one place
3. **Improved Maintainability**: Clear separation of concerns
4. **Better Scalability**: Easy to add new server types
5. **Reduced Duplication**: Shared code and patterns between send and receive

## Potential Challenges

1. **Migration Complexity**: Carefully migrating 400+ packet handlers
2. **Performance Considerations**: Ensuring the new architecture doesn't impact performance
3. **Backward Compatibility**: Maintaining compatibility with existing code
4. **Testing Coverage**: Ensuring comprehensive test coverage for all server types

This hybrid approach combines the strengths of both factory and strategy patterns, providing a consistent architecture across both send and receive sides while maintaining the flexibility needed to handle the large number of packet types.