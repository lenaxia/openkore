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
