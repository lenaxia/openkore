# Ragnarok Online Login Packet Sequence

This document describes the complete login packet sequence for Ragnarok Online, mapping it to the connection states defined in `network.go`. Understanding this sequence is crucial for implementing and debugging the login process in OpenKore.

## Connection States (from network.go)

```go
// Connection states
const (
    // NotConnected represents the state when not connected to any server
    NotConnected = 1

    // ConnectedToMasterServer represents the state when connected to the master server
    ConnectedToMasterServer = 2

    // ConnectedToLoginServer represents the state when connected to the login server
    ConnectedToLoginServer = 3

    // ConnectedToCharServer represents the state when connected to the character server
    ConnectedToCharServer = 4

    // InGame represents the state when connected to the map server and fully functional
    InGame = 5

    // InGameButUninitialized represents the state when in game but without enough information
    // This can only happen in XKore or XKoreProxy mode when the RO client is already
    // logged in before OpenKore, and OpenKore doesn't have enough information
    // (such as character name) to work properly
    InGameButUninitialized = -1
)
```

## Login Packet Sequence

### 1. Secure Login Phase (NotConnected → ConnectedToMasterServer)

The sequence begins with secure login handshaking:

- **Client → Server**: `0x01DB` or `0x0204` (Secure Login request)
- **Server → Client**: `0x01DC` (Secure Login Key response)

This establishes the secure connection and encryption parameters.

### 2. Master Login Phase (ConnectedToMasterServer)

The client authenticates with username/password:

- **Client → Server**: Various packets depending on client version:
  - `0x0064`: Traditional login
  - `0x01DD`: Secure login with salted MD5 password
  - `0x01FA`: Secure Master Login with additional client info
  - `0x0277`, `0x027C`, `0x02B0`: Various login packet versions
  - `0x0825`: Token-based login (kRO Zero)
  - `0x0987`: Login with MD5 hex password
  - `0x0A76`: Login with Rijndael password (tRO)
  - `0x0AAC`: Login with hex password
  - `0x0B04`: Login with access tokens

- **Server → Client**: Various responses depending on server type:
  - `0x0069`: Standard account server info
  - `0x0276`: tRO account server info
  - `0x0AC4`: kRO Zero account server info
  - `0x0AC9`: Newer account server info
  - `0x0B07`: Latest account server info
  - `0x0B60`: twRO account server info

These packets contain session IDs, account ID, and server information.

### 3. Server Selection Phase (ConnectedToMasterServer → ConnectedToLoginServer)

After authentication, the client selects a game server:

- **Client → Server**: `0x0065` or `0x0275` (Server choice packet)
- **Server → Client**: Character list packets:
  - `0x006b`: Standard character list
  - `0x082d`: Newer character list format
  - `0x099D` or `0x0B72`: Latest character list formats

The character list packet contains detailed information about available characters.

### 4. Character Selection Phase (ConnectedToLoginServer → ConnectedToCharServer)

The client selects a character to play:

- **Client → Server**: `0x0066` (Character choice packet)
- **Server → Client**: Character and map information:
  - `0x0071`: Standard character ID and map info
  - `0x0AC5`: Extended character ID and map info with URL

This provides the client with the character ID and map server connection details.

### 5. Map Login Phase (ConnectedToCharServer → InGame)

Finally, the client connects to the map server:

- **Client → Server**: Map login packet (varies by server type)
  - Common packet IDs include `0x0072` and others defined in server type configuration
- **Server → Client**: Multiple packets to complete login:
  - `0x0283`: Account ID confirmation
  - Map loaded notification:
    - `0x0073`: Standard map loaded
    - `0x02EB`: Extended map loaded
    - `0x0A18`: Latest map loaded format
  - `0x013A`: Attack range information
  - `0x00BD`: Character stats information
  - `0x0B1B`: Load confirmation (unlocks keyboard)

At this point, the client is fully logged in and can interact with the game world.

## Packet Encryption

The code also handles packet encryption for certain server types:
- When a character is selected, encryption keys may be initialized
- The `DecryptMessageID` function handles decryption of incoming packet IDs
- Different server types use different encryption methods

## Implementation Notes

1. **State Transitions**: The `NetworkManager.SetState()` method in `network.go` should be called at appropriate points during the login sequence to update the connection state.

2. **Error Handling**: Network errors during any phase should result in appropriate error handling and potentially reverting to the `NotConnected` state.

3. **Server Type Variations**: The packet sequence may vary slightly depending on the server type (kRO, iRO, bRO, etc.). The implementation should be flexible enough to handle these variations.

4. **Packet Handlers**: Implement appropriate packet handlers for each packet type in the sequence to process server responses and trigger state transitions.

5. **Timeout Handling**: Implement timeout handling for each phase of the login sequence to prevent the client from hanging if a server response is not received.

## References

- The packet structures and sequences are derived from analysis of the Poseidon module in OpenKore, particularly `src/Poseidon/RagnarokServer.pm`.
- Additional packet information can be found in the server type configuration files.