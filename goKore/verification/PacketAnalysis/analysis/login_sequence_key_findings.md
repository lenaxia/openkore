# Ragnarok Online Login Sequence - Key Findings

This document summarizes the key findings from analyzing the packet exchange sequence required for a successful login to a Ragnarok Online server.

## Login Process Overview

The login process involves three distinct phases, each communicating with a different server:

1. **Account Server Authentication**
2. **Character Server Selection**
3. **Map Server Connection**

## Critical Packet Exchanges

### 1. Account Server Authentication

```
Client → Server: 0064 (Account Server Login)
Server → Client: 0AC4 (Account Info With Server Info)
```

**Key Components:**
- Packet 0064 contains the username and password credentials
- Packet 0AC4 contains:
  - Account ID (constant across sessions: 0x82841E00 = 2000002)
  - Session ID (variable, changes each login)
  - SessionID2 (variable, changes each login)
  - Server list information

### 2. Character Server Selection

```
Client → Server: 0065 (Character Server Login)
Server → Client: 082D (Received characters info)
Server → Client: 006B (Received characters info)
Server → Client: 08B9 (PinCode Request)
Client → Server: 0066 (Char Login - character selection)
Server → Client: 0AC5 (Received character ID and Map IP)
```

**Key Components:**
- Packet 0065 contains:
  - Account ID from Account Server
  - Session IDs from Account Server
- Packet 0AC5 contains:
  - Character ID (0xF2490200 = 150002)
  - Map name ("gef_fild07.gat")
  - Map server IP and port

### 3. Map Server Connection

```
Client → Server: 0436 (Map Login)
Server → Client: 0283 (Account ID)
Server → Client: 02EB (Enter Map)
Client → Server: 007D (Map Loaded)
```

**Key Components:**
- Packet 0436 contains:
  - Account ID
  - Character ID
  - Session ID from Account Server
- Packet 02EB contains map coordinates and other initialization data
- Packet 007D is a simple acknowledgment that the map has loaded

## Session Data Flow

A critical aspect of the login sequence is the proper handling of session data:

1. Account Server provides:
   - Account ID (0x82841E00)
   - Session ID (variable)
   - SessionID2 (variable)

2. These values must be passed to the Character Server in packet 0065

3. Character Server provides:
   - Character ID (0xF2490200)
   - Map server information

4. Account ID, Character ID, and Session ID must be passed to the Map Server in packet 0436

## Packet Structure Analysis

### Account Server Login (0064)
- Fixed structure with username and password
- Username: "botijo0" (visible in hex dump)
- Password: "Melon.77" (visible in hex dump)

### Character Server Login (0065)
- First 6 bytes: Fixed header with Account ID
- Next 8 bytes: Variable session IDs from Account Server
- Last 2 bytes: Fixed (0x0000)

### Map Server Login (0436)
- First 10 bytes: Fixed header with Account ID and Character ID
- Next 6 bytes: Variable session data

## Implementation Considerations

1. **Session Handling**: The session IDs must be preserved between server transitions
2. **PIN Code**: The server may require PIN code verification (packet 08B9)
3. **Acknowledgments**: The client must properly acknowledge certain server packets (e.g., 007D after map load)
4. **Character Selection**: The character selection packet (0066) is simple but critical

## Conclusion

The login sequence follows a clear pattern across all analyzed login attempts. The key to a successful implementation is properly handling the session data flow between the three server types and correctly formatting each packet according to the server's expectations.