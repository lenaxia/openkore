# Ragnarok Online Packet Structures

This document provides detailed information about the packet structures used in the Ragnarok Online login sequence. It complements the `LOGIN_PACKET_SEQUENCE.md` file by focusing on the specific structure and fields of each packet.

## Secure Login Phase

### Client → Server: Secure Login Request

#### Packet 0x01DB (Secure Login)
```
struct {
    uint16 packetID;    // 0x01DB
    // No additional fields
}
```

#### Packet 0x0204 (Secure Login - Alternative)
```
struct {
    uint16 packetID;    // 0x0204
    // No additional fields
}
```

### Server → Client: Secure Login Key Response

#### Packet 0x01DC (Secure Login Key)
```
struct {
    uint16 packetID;    // 0x01DC
    uint16 packetLength;    // 0x14
    uint8 secure_key[17];   // Encryption key
}
```

## Master Login Phase

### Client → Server: Login Packets

#### Packet 0x0064 (Traditional Login)
```
struct {
    uint16 packetID;    // 0x0064
    uint32 version;     // Client version
    char username[24];  // Account username
    char password[24];  // Account password
    uint8 master_version;   // Master version
}
```

#### Packet 0x01DD (Secure Login with MD5)
```
struct {
    uint16 packetID;    // 0x01DD
    uint32 version;     // Client version
    char username[24];  // Account username
    uint8 password_salted_md5[16];   // MD5 hashed password
    uint8 master_version;   // Master version
}
```

#### Packet 0x01FA (Secure Master Login)
```
struct {
    uint16 packetID;    // 0x01FA
    uint32 version;     // Client version
    char username[24];  // Account username
    uint8 password_salted_md5[16];   // MD5 hashed password
    uint8 master_version;   // Master version
    uint8 clientInfo;   // Client information
}
```

#### Packet 0x0825 (Token Login)
```
struct {
    uint16 packetID;    // 0x0825
    uint16 packetLength;    // Packet length
    uint16 version;     // Client version
    uint16 master_version;  // Master version
    char username[24];  // Account username
    uint8 password_rijndael[27];   // Rijndael encrypted password
    char mac[17];       // MAC address
    char ip[15];        // IP address
    uint8 token[];      // Token data
}
```

#### Packet 0x0987 (Login with MD5 Hex)
```
struct {
    uint16 packetID;    // 0x0987
    uint32 version;     // Client version
    char username[24];  // Account username
    uint8 password_md5_hex[32];   // MD5 hex password
    uint8 master_version;   // Master version
}
```

### Server → Client: Account Server Info

#### Packet 0x0069 (Standard Account Server Info)
```
struct {
    uint16 packetID;    // 0x0069
    uint16 packetLength;    // 0x4F
    uint32 sessionID;   // Session ID
    uint32 accountID;   // Account ID
    uint32 sessionID2;  // Secondary session ID
    uint8 padding[30];  // Padding
    uint8 sex;          // Gender (1=male, 0=female)
    uint8 serverIP[4];  // Server IP address
    uint16 serverPort;  // Server port
    char serverName[20];    // Server name
    uint32 serverUsers; // Number of users on server
    uint8 padding2[2];  // Padding
}
```

#### Packet 0x0AC9 (Newer Account Server Info)
```
struct {
    uint16 packetID;    // 0x0AC9
    uint16 packetLength;    // 0xCF
    uint32 sessionID;   // Session ID
    uint32 accountID;   // Account ID
    uint32 sessionID2;  // Secondary session ID
    uint8 lastLoginIP[4];   // Last login IP
    char lastLoginTime[26]; // Last login time
    uint8 sex;          // Gender
    uint8 padding[6];   // Padding
    char serverName[20];    // Server name
    uint32 serverUsers; // Number of users
    uint8 property[2];  // Server property
    char serverInfo[];  // Server info (IP:port)
    uint8 padding2[114];    // Padding
}
```

## Server Selection Phase

### Client → Server: Server Choice

#### Packet 0x0065 (Server Choice)
```
struct {
    uint16 packetID;    // 0x0065
    uint8 serverIndex;  // Index of selected server
}
```

### Server → Client: Character List

#### Packet 0x006B (Standard Character List)
```
struct {
    uint16 packetID;    // 0x006B
    uint16 packetLength;    // 7 + (character block size * number of characters)
    uint32 accountID;   // Account ID
    uint16 totalSlots;  // Total character slots
    uint8 premium_start_slot; // Premium start slot
    uint8 premium_end_slot;   // Premium end slot
    // Followed by character blocks
}
```

#### Character Block (varies by server type)
```
struct {
    uint32 charID;      // Character ID
    uint32 baseExp;     // Base experience
    uint32 zeny;        // Zeny
    uint32 jobExp;      // Job experience
    uint32 jobLevel;    // Job level
    uint32 opt1;        // Option 1
    uint32 opt2;        // Option 2
    uint32 option;      // Character options
    uint32 stance;      // Stance
    uint32 manner;      // Manner
    uint16 statPoint;   // Status points
    uint32 hp;          // HP
    uint32 maxHp;       // Max HP
    uint16 sp;          // SP
    uint16 maxSp;       // Max SP
    uint16 walkSpeed;   // Walk speed
    uint16 jobId;       // Job ID
    uint16 hairStyle;   // Hair style
    uint16 weapon;      // Weapon
    uint16 level;       // Base level
    uint16 skillPoint;  // Skill points
    uint16 headLow;     // Lower headgear
    uint16 shield;      // Shield
    uint16 headTop;     // Upper headgear
    uint16 headMid;     // Middle headgear
    uint16 hairColor;   // Hair color
    uint16 clothesColor;    // Clothes color
    char name[24];      // Character name
    uint8 str;          // Strength
    uint8 agi;          // Agility
    uint8 vit;          // Vitality
    uint8 int_;         // Intelligence
    uint8 dex;          // Dexterity
    uint8 luk;          // Luck
    uint8 characterSlot;    // Character slot
    uint16 rename;      // Rename
    // Additional fields may be present depending on server type
}
```

## Character Selection Phase

### Client → Server: Character Choice

#### Packet 0x0066 (Character Choice)
```
struct {
    uint16 packetID;    // 0x0066
    uint8 characterSlot;    // Selected character slot
}
```

### Server → Client: Character ID and Map

#### Packet 0x0071 (Standard Character ID and Map)
```
struct {
    uint16 packetID;    // 0x0071
    uint32 charID;      // Character ID
    char mapName[16];   // Map name
    uint8 mapIP[4];     // Map server IP
    uint16 mapPort;     // Map server port
}
```

#### Packet 0x0AC5 (Extended Character ID and Map)
```
struct {
    uint16 packetID;    // 0x0AC5
    uint32 charID;      // Character ID
    char mapName[16];   // Map name
    uint8 mapIP[4];     // Map server IP
    uint16 mapPort;     // Map server port
    char mapURL[128];   // Map server URL
}
```

## Map Login Phase

### Client → Server: Map Login

#### Packet 0x0072 (Map Login)
```
struct {
    uint16 packetID;    // 0x0072
    uint32 accountID;   // Account ID
    uint32 charID;      // Character ID
    uint32 sessionID;   // Session ID
    uint32 tick;        // Client tick
    uint8 sex;          // Gender
}
```

### Server → Client: Map Login Response

#### Packet 0x0073 (Map Loaded)
```
struct {
    uint16 packetID;    // 0x0073
    uint32 tick;        // Server tick
    uint8 pos[3];       // Position (x, y, direction)
    uint8 padding[2];   // Padding
}
```

#### Packet 0x0283 (Account ID Confirmation)
```
struct {
    uint16 packetID;    // 0x0283
    uint32 accountID;   // Account ID
}
```

#### Packet 0x00BD (Stats Info)
```
struct {
    uint16 packetID;    // 0x00BD
    uint16 points_free; // Free status points
    uint8 str;          // Strength
    uint8 points_str;   // Strength points
    uint8 agi;          // Agility
    uint8 points_agi;   // Agility points
    uint8 vit;          // Vitality
    uint8 points_vit;   // Vitality points
    uint8 int_;         // Intelligence
    uint8 points_int;   // Intelligence points
    uint8 dex;          // Dexterity
    uint8 points_dex;   // Dexterity points
    uint8 luk;          // Luck
    uint8 points_luk;   // Luck points
    uint16 attack;      // Attack
    uint16 attack_bonus;    // Attack bonus
    uint16 attack_magic_min;    // Minimum magic attack
    uint16 attack_magic_max;    // Maximum magic attack
    uint16 def;         // Defense
    uint16 def_bonus;   // Defense bonus
    uint16 def_magic;   // Magic defense
    uint16 def_magic_bonus;  // Magic defense bonus
    uint16 hit;         // Hit
    uint16 flee;        // Flee
    uint16 flee_bonus;  // Flee bonus
    uint16 critical;    // Critical
    uint16 aspd;        // Attack speed
    uint16 aspd_bonus;  // Attack speed bonus
}
```

## Notes on Packet Structures

1. **Packet Variations**: The exact structure of packets may vary depending on the server type and version.

2. **Encryption**: Some packets may be encrypted, requiring decryption before processing.

3. **Field Sizes**: The sizes of fields may vary between server types. For example, some servers use different lengths for username and password fields.

4. **Padding**: Many packets include padding fields to align data or reserve space for future use.

5. **Extensions**: Newer server versions often extend existing packet structures with additional fields.

## Implementation Considerations

When implementing packet handling in Go:

1. Use appropriate data structures to represent packets.
2. Implement proper serialization and deserialization functions.
3. Handle endianness correctly (Ragnarok Online uses little-endian).
4. Validate packet lengths and field values.
5. Handle encryption and decryption as needed.

Example Go structure for a login packet:

```go
type LoginPacket struct {
    PacketID    uint16
    Version     uint32
    Username    [24]byte
    Password    [24]byte
    MasterVer   uint8
}

func (p *LoginPacket) Serialize() []byte {
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.LittleEndian, p.PacketID)
    binary.Write(buf, binary.LittleEndian, p.Version)
    buf.Write(p.Username[:])
    buf.Write(p.Password[:])
    buf.WriteByte(p.MasterVer)
    return buf.Bytes()
}