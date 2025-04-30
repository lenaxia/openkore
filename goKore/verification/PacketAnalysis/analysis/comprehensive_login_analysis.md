

---


# Ragnarok Online Login Sequence Analysis

## Account Server Login Sequence
→ 0064: Account Server Login
← 0AC4: Account Info With Server Info

## Character Server Login Sequence
→ 0065: Character Server Login
→ 09A1: Character Received Sync
→ 09A1: Character Received Sync
→ 09A1: Character Received Sync
→ 0066: Char Login
← 082D: Received characters from Game Login Server
← 006B: Received characters from Game Login Server
← 09A0: Received characters from Game Login Server
← 020D: Char Ban List
← 08B9: PinCode Request
← 099D: Received characters Info from Game Login Server
← 099D: Received characters Info from Game Login Server
← 099D: Received characters Info from Game Login Server
← 0AC5: Received character ID and Map IP from Game Login Server

## Map Server Login Sequence
→ 0436: Map Login
→ 007D: Map Loaded
→ 0360: Sync
→ 014F: Guild Query Page
→ 0447: Request Blocking Player Cancel
→ 007D: Map Loaded
→ 0447: Request Blocking Player Cancel
← 0283: Account ID
← 0B18: Inventory expansion result
← 02EB: Enter Map
← 0201: Friend List
← 008E: Public message by yourself
← 008E: Public message by yourself
← 0091: Map Changed
← 00B0: Your Status Info
← 0141: Your Status Changed (Str, Agi, Vit, Int, Dex, Luk, Bonus)
← 0141: Your Status Changed (Str, Agi, Vit, Int, Dex, Luk, Bonus)
← 0141: Your Status Changed (Str, Agi, Vit, Int, Dex, Luk, Bonus)
← 0141: Your Status Changed (Str, Agi, Vit, Int, Dex, Luk, Bonus)
← 0141: Your Status Changed (Str, Agi, Vit, Int, Dex, Luk, Bonus)
← 0141: Your Status Changed (Str, Agi, Vit, Int, Dex, Luk, Bonus)
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 013A: Attack Range
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 00B0: Your Status Info
← 09E7: Rodex Unread Mails
← 0A24: Achievement Update
← 0A23: Achievement List
← 0ADE: Overweight Percent
← 01D7: Weapon / Shield Display
← 01D7: Weapon / Shield Display
← 0B08: Item List Start
← 0B09: Item List Stackable
← 0B0A: Item List Non-Stackable (Equips)
← 0B0B: Item List End
← 0A9B: Switch Equip System Log
← 00B0: Your Status Info
← 00B0: Your Status Info
← 099B: Map Property
← 09FF: actor_display (actor exists)
← 09FF: actor_display (actor exists)
← 010F: Skills List
← 0B20: Hotkey List
← 0B20: Hotkey List
← 0ACB: Your Status Info
← 0ACB: Your Status Info
← 0ACB: Your Status Info
← 0ACB: Your Status Info
← 00B0: Your Status Info
← 00BD: Your Status Info (Calculated)
← 0141: Your Status Changed (Str, Agi, Vit, Int, Dex, Luk, Bonus)
← 0141: Your Status Changed (Str, Agi, Vit, Int, Dex, Luk, Bonus)
← 0141: Your Status Changed (Str, Agi, Vit, Int, Dex, Luk, Bonus)
← 0141: Your Status Changed (Str, Agi, Vit, Int, Dex, Luk, Bonus)
← 0141: Your Status Changed (Str, Agi, Vit, Int, Dex, Luk, Bonus)
← 0141: Your Status Changed (Str, Agi, Vit, Int, Dex, Luk, Bonus)
← 013A: Attack Range
← 00B0: Your Status Info
← 02C9: Party Allow Invite
← 02DA: Show Equipment Window Flag (You)
← 02D9: Show Equipment Window Flag (other)
← 02D9: Show Equipment Window Flag (other)
← 02D9: Show Equipment Window Flag (other)
← 0B1B: Load confirm


---


# Ragnarok Online Packet Structure Analysis

## Key Login Sequence Packets

### 0064 - Account Server Login
**Direction:** → Sent to server
**Packet ID:** 0x0064

**Structure:**

| Offset | Length | Description | Constant | Values |
|--------|--------|-------------|----------|--------|
| 0 | 12 |  | Yes | 0x64 00 1C 00 00 00 62 6F 74 69 6A 6F |
| 12 | 3 |  | No | Variable |
| 15 | 1 |  | Yes | 0x00 |

**Raw Examples:**

Example 1:
```
0>  64 00 1C 00 00 00 62 6F    74 69 6A 6F 30 00 00 00    d.....botijo0...
```

Example 2:
```
0>  64 00 1C 00 00 00 62 6F    74 69 6A 6F 31 32 33 00    d.....botijo123.
```

Example 3:
```
0>  64 00 1C 00 00 00 62 6F    74 69 6A 6F 30 00 00 00    d.....botijo0...
```

### 0AC4 - Account Info With Server Info
**Direction:** ← Received from server
**Packet ID:** 0x0AC4

**Structure:**

| Offset | Length | Description | Constant | Values |
|--------|--------|-------------|----------|--------|
| 0 | 4 |  | Yes | 0xC4 0A E0 00 |
| 4 | 6 |  | No | Variable |
| 10 | 2 |  | Yes | 0x1E 00 |
| 12 | 4 |  | No | Variable |

**Raw Examples:**

Example 1:
```
0>  C4 0A E0 00 CB 3D 2A 6B    82 84 1E 00 FE 74 9A C8    .....=*k.....t..
```

Example 2:
```
0>  C4 0A E0 00 B2 59 B5 1A    FD 84 1E 00 9E 0D 9A E5    .....Y..........
```

Example 3:
```
0>  C4 0A E0 00 CD 1F 6F BC    82 84 1E 00 9B 94 E2 3E    ......o........>
```

### 0065 - Character Server Login
**Direction:** → Sent to server
**Packet ID:** 0x0065

**Structure:**

| Offset | Length | Description | Constant | Values |
|--------|--------|-------------|----------|--------|
| 0 | 2 | Packet ID | Yes | 0x65 00 |
| 2 | 2 | Packet Length | No | Variable |
| 4 | 2 |  | Yes | 0x1E 00 |
| 6 | 8 |  | No | Variable |
| 14 | 2 |  | Yes | 0x00 00 |

**Raw Examples:**

Example 1:
```
0>  65 00 82 84 1E 00 CB 3D    2A 6B FE 74 9A C8 00 00    e......=*k.t....
```

Example 2:
```
0>  65 00 FD 84 1E 00 B2 59    B5 1A 9E 0D 9A E5 00 00    e......Y........
```

Example 3:
```
0>  65 00 82 84 1E 00 CD 1F    6F BC 9B 94 E2 3E 00 00    e.......o....>..
```

### 082D - Received characters from Game Login Server
**Direction:** ← Received from server
**Packet ID:** 0x082D

**Structure:**

| Offset | Length | Description | Constant | Values |
|--------|--------|-------------|----------|--------|
| 0 | 16 |  | Yes | 0x2D 08 1D 00 0F 00 00 09 0F 00 00 00 00 00 00 00 |

**Raw Examples:**

Example 1:
```
0>  2D 08 1D 00 0F 00 00 09    0F 00 00 00 00 00 00 00    -...............
```

Example 2:
```
0>  2D 08 1D 00 0F 00 00 09    0F 00 00 00 00 00 00 00    -...............
```

Example 3:
```
0>  2D 08 1D 00 0F 00 00 09    0F 00 00 00 00 00 00 00    -...............
```

### 006B - Received characters from Game Login Server
**Direction:** ← Received from server
**Packet ID:** 0x006B

**Structure:**

| Offset | Length | Description | Constant | Values |
|--------|--------|-------------|----------|--------|
| 0 | 16 |  | Yes | 0x6B 00 B6 00 0F 0F 0F 00 00 00 00 00 00 00 00 00 |

**Raw Examples:**

Example 1:
```
0>  6B 00 B6 00 0F 0F 0F 00    00 00 00 00 00 00 00 00    k...............
```

Example 2:
```
0>  6B 00 B6 00 0F 0F 0F 00    00 00 00 00 00 00 00 00    k...............
```

Example 3:
```
0>  6B 00 B6 00 0F 0F 0F 00    00 00 00 00 00 00 00 00    k...............
```

### 08B9 - PinCode Request
**Direction:** ← Received from server
**Packet ID:** 0x08B9

**Structure:**

| Offset | Length | Description | Constant | Values |
|--------|--------|-------------|----------|--------|
| 0 | 2 | Packet ID | Yes | 0xB9 08 |
| 2 | 2 | Packet Length | No | Variable |
| 4 | 2 |  | Yes | 0x00 00 |
| 6 | 2 |  | No | Variable |
| 8 | 4 |  | Yes | 0x1E 00 00 00 |

**Raw Examples:**

Example 1:
```
0>  B9 08 22 B0 00 00 82 84    1E 00 00 00                ..".........
```

Example 2:
```
0>  B9 08 01 9E 00 00 FD 84    1E 00 00 00                ............
```

Example 3:
```
0>  B9 08 71 ED 00 00 82 84    1E 00 00 00                ..q.........
```

### 0AC5 - Received character ID and Map IP from Game Login Server
**Direction:** ← Received from server
**Packet ID:** 0x0AC5

**Structure:**

| Offset | Length | Description | Constant | Values |
|--------|--------|-------------|----------|--------|
| 0 | 2 | Packet ID | Yes | 0xC5 0A |
| 2 | 2 | Packet Length | No | Variable |
| 4 | 12 |  | Yes | 0x02 00 67 65 66 5F 66 69 6C 64 30 37 |

**Raw Examples:**

Example 1:
```
0>  C5 0A F2 49 02 00 67 65    66 5F 66 69 6C 64 30 37    ...I..gef_fild07
```

Example 2:
```
0>  C5 0A 6D 4A 02 00 67 65    66 5F 66 69 6C 64 30 37    ..mJ..gef_fild07
```

Example 3:
```
0>  C5 0A F2 49 02 00 67 65    66 5F 66 69 6C 64 30 37    ...I..gef_fild07
```

### 0436 - Map Login
**Direction:** → Sent to server
**Packet ID:** 0x0436

**Structure:**

| Offset | Length | Description | Constant | Values |
|--------|--------|-------------|----------|--------|
| 0 | 2 | Packet ID | Yes | 0x36 04 |
| 2 | 2 | Packet Length | No | Variable |
| 4 | 2 |  | Yes | 0x1E 00 |
| 6 | 2 |  | No | Variable |
| 8 | 2 |  | Yes | 0x02 00 |
| 10 | 6 |  | No | Variable |

**Raw Examples:**

Example 1:
```
0>  36 04 82 84 1E 00 F2 49    02 00 CB 3D 2A 6B 51 88    6......I...=*kQ.
```

Example 2:
```
0>  36 04 FD 84 1E 00 6D 4A    02 00 B2 59 B5 1A 4D 60    6.....mJ...Y..M`
```

Example 3:
```
0>  36 04 82 84 1E 00 F2 49    02 00 CD 1F 6F BC 96 F3    6......I....o...
```

### 02EB - Enter Map
**Direction:** ← Received from server
**Packet ID:** 0x02EB

**Structure:**

| Offset | Length | Description | Constant | Values |
|--------|--------|-------------|----------|--------|
| 0 | 2 | Packet ID | Yes | 0xEB 02 |
| 2 | 7 |  | No | Variable |
| 9 | 4 |  | Yes | 0x05 05 00 00 |

**Raw Examples:**

Example 1:
```
0>  EB 02 44 6C 84 02 3E 4B    C0 05 05 00 00             ..Dl..>K.....
```

Example 2:
```
0>  EB 02 45 44 A3 02 50 CB    D0 05 05 00 00             ..ED..P......
```

Example 3:
```
0>  EB 02 8A D7 12 01 41 8C    80 05 05 00 00             ......A......
```



---


## Sequence Diagram


A PlantUML sequence diagram has been generated in `login_sequence_diagram.puml`.

To visualize this diagram, use PlantUML or an online tool like https://www.planttext.com/



---


# Packet Sequence Comparison Report

## Overview
Number of dumps compared: 7
- dump7_packets.json: 135 packets
- dump3_packets.json: 120 packets
- dump4_packets.json: 107 packets
- dump6_packets.json: 138 packets
- dump1_packets.json: 112 packets
- dump2_packets.json: 106 packets
- dump5_packets.json: 137 packets

## Common Packets
Number of packets common to all dumps: 46

### Character Server
- → 0066
- ← 099D
- ← 08B9
- → 09A1
- ← 020D
- ← 09A0
- ← 082D
- ← 0AC5
- → 0065
- ← 006B

### Map Server
- ← 0B08
- ← 0141
- ← 0201
- → 0360
- ← 099B
- ← 02DA
- → 0436
- ← 0ADE
- ← 0ACB
- ← 01D7
- ← 0B20
- → 007D
- → 0447
- ← 0091
- ← 0A23
- ← 0B0B
- ← 0283
- ← 0A24
- ← 09FF
- ← 00BD
- ← 02C9
- ← 02D9
- ← 0B18
- ← 09E7
- ← 0A9B
- ← 0B0A
- ← 010F
- ← 008E
- ← 0B1B
- ← 013A
- ← 02EB
- → 014F
- ← 00B0
- ← 0B09

### Account Server
- ← 0AC4
- → 0064

## Unique Packets
These packets appear in some dumps but not others:
- ← 07FB (Map Server): Found in dump3_packets.json, dump4_packets.json, dump6_packets.json, dump1_packets.json, dump2_packets.json, dump5_packets.json
- → 0437 (Map Server): Found in dump5_packets.json
- ← 007F (Map Server): Found in dump3_packets.json, dump4_packets.json, dump6_packets.json
- → 035F (Map Server): Found in dump7_packets.json, dump3_packets.json, dump6_packets.json, dump1_packets.json, dump5_packets.json
- ← 0080 (Map Server): Found in dump7_packets.json
- → 0368 (Map Server): Found in dump7_packets.json, dump3_packets.json, dump6_packets.json, dump1_packets.json, dump5_packets.json
- ← 00C0 (Map Server): Found in dump3_packets.json, dump6_packets.json, dump1_packets.json, dump2_packets.json, dump5_packets.json
- → 0361 (Map Server): Found in dump6_packets.json
- ← 099A (Map Server): Found in dump7_packets.json, dump6_packets.json
- ← 09FD (Map Server): Found in dump7_packets.json, dump3_packets.json, dump6_packets.json, dump2_packets.json
- ← 0087 (Map Server): Found in dump7_packets.json, dump3_packets.json, dump6_packets.json, dump1_packets.json, dump5_packets.json
- ← 009D (Map Server): Found in dump6_packets.json
- ← 0088 (Map Server): Found in dump5_packets.json
- ← 0ADF (Map Server): Found in dump7_packets.json, dump3_packets.json, dump6_packets.json, dump1_packets.json, dump5_packets.json

## Sequence Differences
These packets appear in a different order in some dumps:

### ← 0141 (Map Server)
- In dump7_packets.json:
  - Expected after: ← 0B08 (Map Server)
  - Actual position: 32
  - Expected position: 82
- In dump3_packets.json:
  - Expected after: ← 0B08 (Map Server)
  - Actual position: 31
  - Expected position: 65
- In dump4_packets.json:
  - Expected after: ← 0B08 (Map Server)
  - Actual position: 31
  - Expected position: 65
- In dump6_packets.json:
  - Expected after: ← 0B08 (Map Server)
  - Actual position: 33
  - Expected position: 86
- In dump1_packets.json:
  - Expected after: ← 0B08 (Map Server)
  - Actual position: 31
  - Expected position: 65
- In dump2_packets.json:
  - Expected after: ← 0B08 (Map Server)
  - Actual position: 31
  - Expected position: 65
- In dump5_packets.json:
  - Expected after: ← 0B08 (Map Server)
  - Actual position: 33
  - Expected position: 79

### ← 0201 (Map Server)
- In dump7_packets.json:
  - Expected after: ← 0141 (Map Server)
  - Actual position: 24
  - Expected position: 33
- In dump3_packets.json:
  - Expected after: ← 0141 (Map Server)
  - Actual position: 24
  - Expected position: 32
- In dump4_packets.json:
  - Expected after: ← 0141 (Map Server)
  - Actual position: 24
  - Expected position: 32
- In dump6_packets.json:
  - Expected after: ← 0141 (Map Server)
  - Actual position: 24
  - Expected position: 34
- In dump1_packets.json:
  - Expected after: ← 0141 (Map Server)
  - Actual position: 24
  - Expected position: 32
- In dump2_packets.json:
  - Expected after: ← 0141 (Map Server)
  - Actual position: 24
  - Expected position: 32
- In dump5_packets.json:
  - Expected after: ← 0141 (Map Server)
  - Actual position: 24
  - Expected position: 34

### → 0360 (Map Server)
- In dump7_packets.json:
  - Expected after: ← 0201 (Map Server)
  - Actual position: 21
  - Expected position: 25
- In dump3_packets.json:
  - Expected after: ← 0201 (Map Server)
  - Actual position: 21
  - Expected position: 25
- In dump4_packets.json:
  - Expected after: ← 0201 (Map Server)
  - Actual position: 21
  - Expected position: 25
- In dump6_packets.json:
  - Expected after: ← 0201 (Map Server)
  - Actual position: 21
  - Expected position: 25
- In dump1_packets.json:
  - Expected after: ← 0201 (Map Server)
  - Actual position: 21
  - Expected position: 25
- In dump2_packets.json:
  - Expected after: ← 0201 (Map Server)
  - Actual position: 21
  - Expected position: 25
- In dump5_packets.json:
  - Expected after: ← 0201 (Map Server)
  - Actual position: 21
  - Expected position: 25

### ← 099D (Character Server)
- In dump7_packets.json:
  - Expected after: ← 02DA (Map Server)
  - Actual position: 11
  - Expected position: 112
- In dump3_packets.json:
  - Expected after: ← 02DA (Map Server)
  - Actual position: 11
  - Expected position: 98
- In dump4_packets.json:
  - Expected after: ← 02DA (Map Server)
  - Actual position: 11
  - Expected position: 100
- In dump6_packets.json:
  - Expected after: ← 02DA (Map Server)
  - Actual position: 11
  - Expected position: 118
- In dump1_packets.json:
  - Expected after: ← 02DA (Map Server)
  - Actual position: 11
  - Expected position: 96
- In dump2_packets.json:
  - Expected after: ← 02DA (Map Server)
  - Actual position: 11
  - Expected position: 97
- In dump5_packets.json:
  - Expected after: ← 02DA (Map Server)
  - Actual position: 11
  - Expected position: 109

### ← 08B9 (Character Server)
- In dump7_packets.json:
  - Expected after: → 0436 (Map Server)
  - Actual position: 10
  - Expected position: 17
- In dump3_packets.json:
  - Expected after: → 0436 (Map Server)
  - Actual position: 10
  - Expected position: 17
- In dump4_packets.json:
  - Expected after: → 0436 (Map Server)
  - Actual position: 10
  - Expected position: 17
- In dump6_packets.json:
  - Expected after: → 0436 (Map Server)
  - Actual position: 10
  - Expected position: 17
- In dump1_packets.json:
  - Expected after: → 0436 (Map Server)
  - Actual position: 10
  - Expected position: 17
- In dump2_packets.json:
  - Expected after: → 0436 (Map Server)
  - Actual position: 10
  - Expected position: 17
- In dump5_packets.json:
  - Expected after: → 0436 (Map Server)
  - Actual position: 10
  - Expected position: 17

### → 09A1 (Character Server)
- In dump7_packets.json:
  - Expected after: ← 0ACB (Map Server)
  - Actual position: 6
  - Expected position: 97
- In dump3_packets.json:
  - Expected after: ← 0ACB (Map Server)
  - Actual position: 6
  - Expected position: 83
- In dump4_packets.json:
  - Expected after: ← 0ACB (Map Server)
  - Actual position: 6
  - Expected position: 85
- In dump6_packets.json:
  - Expected after: ← 0ACB (Map Server)
  - Actual position: 6
  - Expected position: 103
- In dump1_packets.json:
  - Expected after: ← 0ACB (Map Server)
  - Actual position: 6
  - Expected position: 81
- In dump2_packets.json:
  - Expected after: ← 0ACB (Map Server)
  - Actual position: 6
  - Expected position: 82
- In dump5_packets.json:
  - Expected after: ← 0ACB (Map Server)
  - Actual position: 6
  - Expected position: 94

### → 007D (Map Server)
- In dump7_packets.json:
  - Expected after: ← 0B20 (Map Server)
  - Actual position: 20
  - Expected position: 95
- In dump3_packets.json:
  - Expected after: ← 0B20 (Map Server)
  - Actual position: 20
  - Expected position: 81
- In dump4_packets.json:
  - Expected after: ← 0B20 (Map Server)
  - Actual position: 20
  - Expected position: 83
- In dump6_packets.json:
  - Expected after: ← 0B20 (Map Server)
  - Actual position: 20
  - Expected position: 101
- In dump1_packets.json:
  - Expected after: ← 0B20 (Map Server)
  - Actual position: 20
  - Expected position: 79
- In dump2_packets.json:
  - Expected after: ← 0B20 (Map Server)
  - Actual position: 20
  - Expected position: 80
- In dump5_packets.json:
  - Expected after: ← 0B20 (Map Server)
  - Actual position: 20
  - Expected position: 92

### ← 0AC4 (Account Server)
- In dump7_packets.json:
  - Expected after: → 007D (Map Server)
  - Actual position: 1
  - Expected position: 21
- In dump3_packets.json:
  - Expected after: → 007D (Map Server)
  - Actual position: 1
  - Expected position: 21
- In dump4_packets.json:
  - Expected after: → 007D (Map Server)
  - Actual position: 1
  - Expected position: 21
- In dump6_packets.json:
  - Expected after: → 007D (Map Server)
  - Actual position: 1
  - Expected position: 21
- In dump1_packets.json:
  - Expected after: → 007D (Map Server)
  - Actual position: 1
  - Expected position: 21
- In dump2_packets.json:
  - Expected after: → 007D (Map Server)
  - Actual position: 1
  - Expected position: 21
- In dump5_packets.json:
  - Expected after: → 007D (Map Server)
  - Actual position: 1
  - Expected position: 21

### ← 0283 (Map Server)
- In dump7_packets.json:
  - Expected after: ← 0B0B (Map Server)
  - Actual position: 17
  - Expected position: 85
- In dump3_packets.json:
  - Expected after: ← 0B0B (Map Server)
  - Actual position: 17
  - Expected position: 68
- In dump4_packets.json:
  - Expected after: ← 0B0B (Map Server)
  - Actual position: 17
  - Expected position: 68
- In dump6_packets.json:
  - Expected after: ← 0B0B (Map Server)
  - Actual position: 17
  - Expected position: 89
- In dump1_packets.json:
  - Expected after: ← 0B0B (Map Server)
  - Actual position: 17
  - Expected position: 68
- In dump2_packets.json:
  - Expected after: ← 0B0B (Map Server)
  - Actual position: 17
  - Expected position: 68
- In dump5_packets.json:
  - Expected after: ← 0B0B (Map Server)
  - Actual position: 17
  - Expected position: 82

### ← 09A0 (Character Server)
- In dump7_packets.json:
  - Expected after: ← 0A24 (Map Server)
  - Actual position: 5
  - Expected position: 60
- In dump3_packets.json:
  - Expected after: ← 0A24 (Map Server)
  - Actual position: 5
  - Expected position: 60
- In dump4_packets.json:
  - Expected after: ← 0A24 (Map Server)
  - Actual position: 5
  - Expected position: 60
- In dump6_packets.json:
  - Expected after: ← 0A24 (Map Server)
  - Actual position: 5
  - Expected position: 61
- In dump1_packets.json:
  - Expected after: ← 0A24 (Map Server)
  - Actual position: 5
  - Expected position: 60
- In dump2_packets.json:
  - Expected after: ← 0A24 (Map Server)
  - Actual position: 5
  - Expected position: 60
- In dump5_packets.json:
  - Expected after: ← 0A24 (Map Server)
  - Actual position: 5
  - Expected position: 62

### ← 082D (Character Server)
- In dump7_packets.json:
  - Expected after: ← 00BD (Map Server)
  - Actual position: 3
  - Expected position: 102
- In dump3_packets.json:
  - Expected after: ← 00BD (Map Server)
  - Actual position: 3
  - Expected position: 88
- In dump4_packets.json:
  - Expected after: ← 00BD (Map Server)
  - Actual position: 3
  - Expected position: 90
- In dump6_packets.json:
  - Expected after: ← 00BD (Map Server)
  - Actual position: 3
  - Expected position: 108
- In dump1_packets.json:
  - Expected after: ← 00BD (Map Server)
  - Actual position: 3
  - Expected position: 86
- In dump2_packets.json:
  - Expected after: ← 00BD (Map Server)
  - Actual position: 3
  - Expected position: 87
- In dump5_packets.json:
  - Expected after: ← 00BD (Map Server)
  - Actual position: 3
  - Expected position: 99

### ← 0AC5 (Character Server)
- In dump7_packets.json:
  - Expected after: ← 02D9 (Map Server)
  - Actual position: 15
  - Expected position: 113
- In dump3_packets.json:
  - Expected after: ← 02D9 (Map Server)
  - Actual position: 15
  - Expected position: 99
- In dump4_packets.json:
  - Expected after: ← 02D9 (Map Server)
  - Actual position: 15
  - Expected position: 101
- In dump6_packets.json:
  - Expected after: ← 02D9 (Map Server)
  - Actual position: 15
  - Expected position: 119
- In dump1_packets.json:
  - Expected after: ← 02D9 (Map Server)
  - Actual position: 15
  - Expected position: 97
- In dump2_packets.json:
  - Expected after: ← 02D9 (Map Server)
  - Actual position: 15
  - Expected position: 98
- In dump5_packets.json:
  - Expected after: ← 02D9 (Map Server)
  - Actual position: 15
  - Expected position: 110

### → 0064 (Account Server)
- In dump7_packets.json:
  - Expected after: ← 0B18 (Map Server)
  - Actual position: 0
  - Expected position: 19
- In dump3_packets.json:
  - Expected after: ← 0B18 (Map Server)
  - Actual position: 0
  - Expected position: 19
- In dump4_packets.json:
  - Expected after: ← 0B18 (Map Server)
  - Actual position: 0
  - Expected position: 19
- In dump6_packets.json:
  - Expected after: ← 0B18 (Map Server)
  - Actual position: 0
  - Expected position: 19
- In dump1_packets.json:
  - Expected after: ← 0B18 (Map Server)
  - Actual position: 0
  - Expected position: 19
- In dump2_packets.json:
  - Expected after: ← 0B18 (Map Server)
  - Actual position: 0
  - Expected position: 19
- In dump5_packets.json:
  - Expected after: ← 0B18 (Map Server)
  - Actual position: 0
  - Expected position: 19

### ← 0B0A (Map Server)
- In dump7_packets.json:
  - Expected after: ← 0A9B (Map Server)
  - Actual position: 83
  - Expected position: 88
- In dump3_packets.json:
  - Expected after: ← 0A9B (Map Server)
  - Actual position: 66
  - Expected position: 69
- In dump4_packets.json:
  - Expected after: ← 0A9B (Map Server)
  - Actual position: 66
  - Expected position: 69
- In dump6_packets.json:
  - Expected after: ← 0A9B (Map Server)
  - Actual position: 87
  - Expected position: 92
- In dump1_packets.json:
  - Expected after: ← 0A9B (Map Server)
  - Actual position: 66
  - Expected position: 69
- In dump2_packets.json:
  - Expected after: ← 0A9B (Map Server)
  - Actual position: 66
  - Expected position: 69
- In dump5_packets.json:
  - Expected after: ← 0A9B (Map Server)
  - Actual position: 80
  - Expected position: 83

### → 0065 (Character Server)
- In dump7_packets.json:
  - Expected after: ← 010F (Map Server)
  - Actual position: 2
  - Expected position: 94
- In dump3_packets.json:
  - Expected after: ← 010F (Map Server)
  - Actual position: 2
  - Expected position: 80
- In dump4_packets.json:
  - Expected after: ← 010F (Map Server)
  - Actual position: 2
  - Expected position: 82
- In dump6_packets.json:
  - Expected after: ← 010F (Map Server)
  - Actual position: 2
  - Expected position: 100
- In dump1_packets.json:
  - Expected after: ← 010F (Map Server)
  - Actual position: 2
  - Expected position: 78
- In dump2_packets.json:
  - Expected after: ← 010F (Map Server)
  - Actual position: 2
  - Expected position: 79
- In dump5_packets.json:
  - Expected after: ← 010F (Map Server)
  - Actual position: 2
  - Expected position: 91

### ← 013A (Map Server)
- In dump7_packets.json:
  - Expected after: ← 0B1B (Map Server)
  - Actual position: 52
  - Expected position: 116
- In dump3_packets.json:
  - Expected after: ← 0B1B (Map Server)
  - Actual position: 52
  - Expected position: 104
- In dump4_packets.json:
  - Expected after: ← 0B1B (Map Server)
  - Actual position: 52
  - Expected position: 105
- In dump6_packets.json:
  - Expected after: ← 0B1B (Map Server)
  - Actual position: 53
  - Expected position: 123
- In dump1_packets.json:
  - Expected after: ← 0B1B (Map Server)
  - Actual position: 52
  - Expected position: 100
- In dump2_packets.json:
  - Expected after: ← 0B1B (Map Server)
  - Actual position: 52
  - Expected position: 101
- In dump5_packets.json:
  - Expected after: ← 0B1B (Map Server)
  - Actual position: 54
  - Expected position: 113

### ← 02EB (Map Server)
- In dump7_packets.json:
  - Expected after: ← 013A (Map Server)
  - Actual position: 19
  - Expected position: 53
- In dump3_packets.json:
  - Expected after: ← 013A (Map Server)
  - Actual position: 19
  - Expected position: 53
- In dump4_packets.json:
  - Expected after: ← 013A (Map Server)
  - Actual position: 19
  - Expected position: 53
- In dump6_packets.json:
  - Expected after: ← 013A (Map Server)
  - Actual position: 19
  - Expected position: 54
- In dump1_packets.json:
  - Expected after: ← 013A (Map Server)
  - Actual position: 19
  - Expected position: 53
- In dump2_packets.json:
  - Expected after: ← 013A (Map Server)
  - Actual position: 19
  - Expected position: 53
- In dump5_packets.json:
  - Expected after: ← 013A (Map Server)
  - Actual position: 19
  - Expected position: 55

### ← 006B (Character Server)
- In dump7_packets.json:
  - Expected after: ← 00B0 (Map Server)
  - Actual position: 4
  - Expected position: 31
- In dump3_packets.json:
  - Expected after: ← 00B0 (Map Server)
  - Actual position: 4
  - Expected position: 31
- In dump4_packets.json:
  - Expected after: ← 00B0 (Map Server)
  - Actual position: 4
  - Expected position: 31
- In dump6_packets.json:
  - Expected after: ← 00B0 (Map Server)
  - Actual position: 4
  - Expected position: 31
- In dump1_packets.json:
  - Expected after: ← 00B0 (Map Server)
  - Actual position: 4
  - Expected position: 31
- In dump2_packets.json:
  - Expected after: ← 00B0 (Map Server)
  - Actual position: 4
  - Expected position: 31
- In dump5_packets.json:
  - Expected after: ← 00B0 (Map Server)
  - Actual position: 4
  - Expected position: 31

## Server Transitions
Server transitions in each dump:

### dump7_packets.json
- Account Server → Character Server
- Character Server → Map Server

### dump3_packets.json
- Account Server → Character Server
- Character Server → Map Server

### dump4_packets.json
- Account Server → Character Server
- Character Server → Map Server

### dump6_packets.json
- Account Server → Character Server
- Character Server → Map Server

### dump1_packets.json
- Account Server → Character Server
- Character Server → Map Server

### dump2_packets.json
- Account Server → Character Server
- Character Server → Map Server

### dump5_packets.json
- Account Server → Character Server
- Character Server → Map Server