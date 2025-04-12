# Error Handling

This file details all error handling mechanisms used in Receive.pm.

## Core Error Types

### Network Errors
- Connection loss
- Packet corruption 
- Invalid packet sequence
- Server timeout
- Authentication failure

### Game State Errors  
- Invalid state transitions
- Desynchronization
- Invalid actor references
- Map load failures
- Inventory corruption

### Resource Errors
- Memory allocation failures
- File I/O errors
- Database access errors
- Configuration errors

### Logic Errors
- Invalid calculations
- Infinite loops
- Stack overflow
- Race conditions

## Error Handling Patterns

### Try-Catch Pattern
```perl
eval {
  # Risky operation
};
if (my $e = caught('FileNotFoundException', 'IOException')) {
  error TF("Cannot load field %s: %s\n", $map_noinstance, $e);
} elsif ($@) {
  die $@; 
}
```

### Status Check Pattern
```perl
if (!defined $actor) {
  error TF("Actor %d not found\n", $ID);
  return;
}
```

### Validation Pattern
```perl
unless (changeToInGameState()) {
  error "Not in game state\n";
  return;
}
```

## Error Recovery Strategies

### Immediate Recovery
- Retry failed operation
- Use cached/default values
- Reconnect to server
- Reload resources

### Graceful Degradation  
- Continue with reduced functionality
- Fall back to safe mode
- Log and continue

### Fatal Errors
- Save state and exit
- Emergency shutdown
- Alert administrator

## Error Reporting

### Log Levels
- debug: Detailed debugging info
- info: General operational info  
- warning: Potential issues
- error: Error conditions
- critical: System failure

### Error Context
- Stack trace
- System state
- Configuration
- Resource usage

## Server-Specific Error Handling

### bRO
- Custom packet validation
- Encryption error handling
- Rate limiting errors

### iRO
- Payment errors
- Regional restrictions
- Version mismatch

### Private Servers  
- Custom packet formats
- Non-standard behaviors
- Missing features

## Plugin Error Integration

### Hook Points
- pre_error
- on_error  
- post_error
- error_recovery

### Error Data
```perl
{
  type => 'network',
  code => 5,
  message => 'Connection lost',
  context => {
    state => $state,
    lastPacket => $lastPacket
  }
}
```

## Testing Error Conditions

### Unit Tests
- Test each error handler
- Verify recovery
- Check logging

### Integration Tests  
- End-to-end scenarios
- Server simulation
- Load testing

## Best Practices

### Error Prevention
- Input validation
- State verification
- Resource checking
- Timeouts

### Error Handling
- Specific error types
- Meaningful messages
- Proper logging
- Recovery plans

### Error Reporting
- Clear messages
- Relevant context
- Action items
- Contact info

## Implementation Details

### Error Types
```perl
package Network::Receive::Error;

use constant {
  NETWORK => 1,
  STATE => 2,
  RESOURCE => 3,
  LOGIC => 4
};
```

### Error Handlers
```perl
sub handleError {
  my ($self, $type, $code, $msg) = @_;
  
  if ($type == NETWORK) {
    $self->handleNetworkError($code, $msg);
  } elsif ($type == STATE) {
    $self->handleStateError($code, $msg);
  }
  # etc
}
```

### Recovery Functions
```perl
sub recoverFromError {
  my ($self, $error) = @_;
  
  given ($error->type) {
    when (NETWORK) { $self->reconnect() }
    when (STATE) { $self->resetState() }
    default { $self->emergencyShutdown() }
  }
}
```

## General Error Handling Structure

1. Error Detection
2. Logging
3. Error-specific Actions
4. Plugin Notification
5. Recovery Attempts (if applicable)

## Common Error Scenarios

### Invalid Packet Structure
- Detection: Packet length or content doesn't match expected format
- Logging: `$logger->error("Invalid packet structure for packet $packetID")`
- Action: Skip packet processing
- Plugin Hook: `packet_error`
- Recovery: Request resend of packet if critical

### Disconnection
- Detection: Socket closed or timeout
- Logging: `$logger->error("Disconnected from server")`
- Action: Set `$conState` to 'disconnected', attempt reconnection
- Plugin Hook: `disconnected`
- Recovery: Initiate reconnection process

### Invalid Game State
- Detection: Received packet inconsistent with current game state
- Logging: `$logger->warning("Unexpected packet $packetID in state $conState")`
- Action: Attempt to recover game state
- Plugin Hook: `invalid_game_state`
- Recovery: Request state sync from server

### Inventory Desync
- Detection: Server and client inventory states don't match
- Logging: `$logger->error("Inventory desync detected")`
- Action: Request full inventory sync from server
- Plugin Hook: `inventory_desync`
- Recovery: Perform full inventory update

### Map Load Failure
- Detection: Unable to load map data
- Logging: `$logger->error("Failed to load map $map_name")`
- Action: Attempt to rerequest map data, if fails, disconnect
- Plugin Hook: `map_load_error`
- Recovery: Reattempt map load or switch to last known good map

## Packet-Specific Error Handling

### 0x0097 - Private Message
- Error: Message too long
- Handling: Truncate message, log warning
- Logging: `$logger->warning("Private message truncated")`
- Recovery: Split long messages if supported by server

### 0x00A0 - Inventory Item Added
- Error: Item doesn't exist in game database
- Handling: Create temporary item entry, log error
- Logging: `$logger->error("Unknown item ID: $item_id")`
- Recovery: Request item data from server if possible

### 0x0078 - Character Move
- Error: Invalid coordinates
- Handling: Ignore movement, request position sync
- Logging: `$logger->warning("Invalid move coordinates: $x, $y")`
- Recovery: Use last known valid position

### 0x00B0 - Stats Info
- Error: Invalid stat value
- Handling: Ignore update, log error
- Logging: `$logger->error("Invalid stat value for $stat_type: $value")`
- Recovery: Request full stats update

### 0x00B4 - NPC Talk
- Error: NPC ID not found
- Handling: Ignore dialog, log error
- Logging: `$logger->error("NPC not found for dialog: $npc_id")`
- Recovery: Attempt to end NPC interaction

### 0x00BE - Update Status
- Error: Unknown status effect
- Handling: Ignore status, log warning
- Logging: `$logger->warning("Unknown status effect: $status_id")`
- Recovery: Request status effect list if supported

### 0x0080 - Remove Entity
- Error: Entity not found in lists
- Handling: Ignore removal, log warning
- Logging: `$logger->warning("Attempted to remove non-existent entity: $entity_id")`
- Recovery: Perform full entity list update

### 0x00B6 - NPC Talk Close
- Error: No active NPC conversation
- Handling: Ignore close request, log warning
- Logging: `$logger->warning("Attempted to close non-existent NPC dialog")`
- Recovery: Reset NPC interaction state

## Server-Specific Error Handling

### bRO (Brazilian Ragnarok Online)
- Packet 0x0078: Additional checksum validation
- Error: Checksum mismatch
- Handling: Log error, request resend
- Logging: `$logger->error("bRO move packet checksum mismatch")`
- Recovery: Use last known valid position

### Custom Private Servers
- Error: Unexpected packet structure
- Handling: Attempt to parse with fallback method, log error
- Logging: `$logger->error("Unexpected packet structure for server type: $server_type")`
- Recovery: Switch to more permissive parsing mode

## Critical Error Handling

### Severe Desync
- Detection: Multiple consecutive invalid packets or states
- Action: Force disconnect and full reconnection
- Logging: `$logger->critical("Severe desync detected, forcing reconnection")`
- Recovery: Perform full login sequence

### Memory Corruption
- Detection: Unexpected data types or values in critical structures
- Action: Emergency data dump, client shutdown
- Logging: `$logger->emergency("Critical memory corruption detected, shutting down")`
- Recovery: Require manual restart and log analysis

Note: This error handling overview covers the main scenarios in Receive.pm. The actual implementation includes more specific error cases and may vary based on server type and OpenKore configuration.
