# Essential Login Sequence Packets

## Account Server

→ 0064: Account Server Login [55 bytes] 2025.04.30 16:24:58
Raw data:
0>  64 00 1C 00 00 00 62 6F    74 69 6A 6F 30 00 00 00    d.....botijo0...

← 0AC4: Account Info With Server Info [224 bytes] 2025.04.30 16:24:58
Raw data:
0>  C4 0A E0 00 E5 5D F6 C1    82 84 1E 00 01 2C 9C 53    .....].......,.S

## Character Server

→ 0065: Character Server Login [17 bytes] 2025.04.30 16:24:58
Raw data:
0>  65 00 82 84 1E 00 E5 5D    F6 C1 01 2C 9C 53 00 00    e......]...,.S..

→ 0066: Char Login [3 bytes] 2025.04.30 16:25:00
Raw data:
0>  C5 0A F2 49 02 00 67 65    66 5F 66 69 6C 64 30 37    ...I..gef_fild07

← 082D: Received characters from Game Login Server [29 bytes] 2025.04.30 16:24:58
Raw data:
0>  2D 08 1D 00 0F 00 00 09    0F 00 00 00 00 00 00 00    -...............

← 006B: Received characters from Game Login Server [182 bytes] 2025.04.30 16:24:58
Raw data:
0>  6B 00 B6 00 0F 0F 0F 00    00 00 00 00 00 00 00 00    k...............

← 08B9: PinCode Request [12 bytes] 2025.04.30 16:24:58
Raw data:
0>  B9 08 37 37 00 00 82 84    1E 00 00 00                ..77........

← 0AC5: Received character ID and Map IP from Game Login Server [156 bytes] 2025.04.30 16:25:00
Raw data:
0>  C5 0A F2 49 02 00 67 65    66 5F 66 69 6C 64 30 37    ...I..gef_fild07

## Map Server

→ 0436: Map Login [19 bytes] 2025.04.30 16:25:00
Raw data:
0>  36 04 82 84 1E 00 F2 49    02 00 E5 5D F6 C1 D6 5A    6......I...]...Z

→ 007D: Map Loaded [2 bytes] 2025.04.30 16:25:00

→ 007D: Map Loaded [2 bytes] 2025.04.30 16:25:00
Raw data:
0>  41 01 0D 00 00 00 5A 00    00 00 0A 00 00 00          A.....Z.......

← 0283: Account ID [6 bytes] 2025.04.30 16:25:00
Raw data:
0>  EB 02 C9 3E 82 02 3D 8B    F0 05 05 00 00             ...>..=......

← 02EB: Enter Map [13 bytes] 2025.04.30 16:25:00
Raw data:
0>  EB 02 C9 3E 82 02 3D 8B    F0 05 05 00 00             ...>..=......
