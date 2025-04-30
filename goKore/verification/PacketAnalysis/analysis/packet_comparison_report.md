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