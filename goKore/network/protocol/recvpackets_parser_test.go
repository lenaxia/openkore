package protocol

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseRecvPackets(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Create test recvpackets.txt file
	testFilePath := filepath.Join(tempDir, "recvpackets.txt")
	testContent := `# Packet lengths for ServerType0
# This file contains the packet lengths for the server type.
# The format is: packet ID, packet length
# Packet length is the fixed length of the packet or 0 if the packet has a variable length.

0064 55
0065 17
0066 3
0067 37
0068 46
0069 -1
006A 23
006B -1
006C 3
006D 149
006E 3
006F 2
0070 3
0071 28
0072 22 4 8 # Variable length packet with additional fields
0073 11
0074 3
0075 -1 # Variable length packet
0076 9
0077 5
0078 55
0079 53
007A 58
007B 60
007C 44
007D 2
007E 105
007F 6
0080 7
0081 3
0082 2
0083 2
0084 2
0085 5
0086 16
0087 12
0088 10
0089 11
008A 29
008B 23
008C 14
008D -1
008E -1
0090 7
0091 22
0092 28
0093 2
0094 19
0095 30
0096 -1
0097 -1
0098 3
0099 -1
009A -1
009B 34
009C 9
009D 17
009E 17
009F 20
00A0 23
00A1 6
00A2 6
00A3 -1
00A4 -1
00A5 -1
00A6 -1
00A7 8
00A8 7
00A9 6
00AA 9
00AB 4
00AC 7
00AE -1
00AF 6
00B0 8
00B1 8
00B2 3
00B3 3
00B4 -1
00B5 6
00B6 6
00B7 -1
00B8 7
00B9 6
00BA 2
00BB 5
00BC 6
00BD 44
00BE 5
00BF 3
00C0 7
00C1 2
00C2 6
00C3 8
00C4 6
00C5 7
00C6 -1
00C7 -1
00C8 -1
00C9 -1
00CA 3
00CB 3
00CC 6
00CD 3
00CE 2
00CF 27
00D0 3
00D1 4
00D2 4
00D3 2
00D4 -1
00D5 -1
00D6 3
00D7 -1
00D8 6
00D9 14
00DA 3
00DB -1
00DC 28
00DD 29
00DE -1
00DF -1
00E0 30
00E1 30
00E2 26
00E3 2
00E4 6
00E5 26
00E6 3
00E7 3
00E8 8
00E9 19
00EA 5
00EB 2
00EC 3
00ED 2
00EE 2
00EF 2
00F0 3
00F1 2
00F2 6
00F3 8
00F4 21
00F5 8
00F6 8
00F7 22
00F8 2
00F9 26
00FA 3
00FB -1
00FC 6
00FD 27
00FE 30
00FF 10
0100 2
0101 6
0102 6
0103 30
0104 79
0105 31
0106 10
0107 10
0108 -1
0109 -1
010A 4
010B 6
010C 6
010D 2
010E 11
010F -1
0110 10
0111 39
0112 4
0113 25
0114 31
0115 35
0116 17
0117 18
0118 2
0119 13
011A 15
011B 20
011C 68
011D 2
011E 3
011F 16
0120 6
0121 14
0122 -1
0123 -1
0124 21
0125 8
0126 8
0127 8
0128 8
0129 8
012A 2
012B 2
012C 3
012D 4
012E 2
012F -1
0130 6
0131 86
0132 6
0133 -1
0134 -1
0135 7
0136 -1
0137 6
0138 3
0139 16
013A 4
013B 4
013C 4
013D 6
013E 24
013F 26
0140 22
0141 14
0142 6
0143 10
0144 23
0145 19
0146 6
0147 39
0148 8
0149 9
014A 6
014B 27
014C -1
014D 2
014E 6
014F 6
0150 110
0151 6
0152 -1
0153 -1
0154 -1
0155 -1
0156 -1
0157 6
0158 -1
0159 54
015A 66
015B 54
015C 90
015D 42
015E 6
015F 42
0160 -1
0161 -1
0162 -1
0163 -1
0164 -1
0165 30
0166 -1
0167 3
0168 14
0169 3
016A 30
016B 10
016C 43
016D 14
016E 186
016F 182
0170 14
0171 30
0172 10
0173 3
0174 -1
0175 6
0176 106
0177 -1
0178 4
0179 5
017A 4
017B -1
017C 6
017D 7
017E -1
017F -1
0180 6
0181 3
0182 106
0183 10
0184 10
0185 34
0187 6
0188 8
0189 4
018A 4
018B 4
018C 29
018D -1
018E 10
018F 6
0190 23
0191 86
0192 24
0193 2
0194 30
0195 102
0196 9
0197 4
0198 8
0199 4
019A 14
019B 10
`
	err := os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test parsing recvpackets.txt
	t.Run("Parse recvpackets.txt", func(t *testing.T) {
		packetDefs, err := ParseRecvPackets(testFilePath)
		if err != nil {
			t.Errorf("ParseRecvPackets() error = %v", err)
			return
		}

		// Check a few specific packet definitions
		expectedDefs := map[string]PacketLengthDef{
			"0064": {Length: 55, HasLength: false},
			"0069": {Length: -1, HasLength: true},
			"0072": {Length: 22, HasLength: false},
			"0075": {Length: -1, HasLength: true},
			"00B4": {Length: -1, HasLength: true},
			"0152": {Length: -1, HasLength: true},
			"0196": {Length: 9, HasLength: false},
		}

		for id, expected := range expectedDefs {
			actual, exists := packetDefs[id]
			if !exists {
				t.Errorf("ParseRecvPackets() missing packet definition for %s", id)
				continue
			}

			if actual.Length != expected.Length {
				t.Errorf("ParseRecvPackets() packet %s length = %d, want %d", id, actual.Length, expected.Length)
			}

			if actual.HasLength != expected.HasLength {
				t.Errorf("ParseRecvPackets() packet %s hasLength = %v, want %v", id, actual.HasLength, expected.HasLength)
			}
		}

		// Check total number of packet definitions
		expectedCount := 309 // Number of packet definitions in the test file
		if len(packetDefs) != expectedCount {
			t.Errorf("ParseRecvPackets() got %d packet definitions, want %d", len(packetDefs), expectedCount)
		}
	})

	// Test parsing invalid recvpackets.txt
	t.Run("Parse invalid recvpackets.txt", func(t *testing.T) {
		invalidFilePath := filepath.Join(tempDir, "invalid.txt")
		invalidContent := `# Invalid recvpackets.txt
0064 invalid
0065 17
0066 3
`
		err := os.WriteFile(invalidFilePath, []byte(invalidContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		_, err = ParseRecvPackets(invalidFilePath)
		if err == nil {
			t.Errorf("ParseRecvPackets() expected error for invalid file, got nil")
		}
	})

	// Test parsing non-existent file
	t.Run("Parse non-existent file", func(t *testing.T) {
		_, err := ParseRecvPackets(filepath.Join(tempDir, "nonexistent.txt"))
		if err == nil {
			t.Errorf("ParseRecvPackets() expected error for non-existent file, got nil")
		}
	})
}

func TestConvertToTokenizerPacketDefs(t *testing.T) {
	// Create test packet definitions
	recvPackets := map[string]struct {
		Length int
	}{
		"0064": {Length: 55},
		"0065": {Length: 17},
		"0066": {Length: 3},
		"0069": {Length: -1},
		"006A": {Length: 23},
		"006B": {Length: -1},
		"006C": {Length: 3},
	}

	// Convert to tokenizer packet definitions
	tokenizerDefs := ConvertToTokenizerPacketDefs(recvPackets)

	// Check conversion
	expected := map[string]PacketLengthDef{
		"0064": {Length: 55, HasLength: false},
		"0065": {Length: 17, HasLength: false},
		"0066": {Length: 3, HasLength: false},
		"0069": {Length: -1, HasLength: true},
		"006A": {Length: 23, HasLength: false},
		"006B": {Length: -1, HasLength: true},
		"006C": {Length: 3, HasLength: false},
	}

	if !reflect.DeepEqual(tokenizerDefs, expected) {
		t.Errorf("ConvertToTokenizerPacketDefs() = %v, want %v", tokenizerDefs, expected)
	}
}

func TestLoadRecvPackets(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Create test folders
	folders := []string{
		filepath.Join(tempDir, "kRO/RagexeRE_2020_04_01b"),
		filepath.Join(tempDir, "iRO/official"),
		filepath.Join(tempDir, "iRO"),
	}

	for _, folder := range folders {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	// Create test recvpackets.txt files
	testFiles := map[string]string{
		filepath.Join(tempDir, "kRO/RagexeRE_2020_04_01b/recvpackets.txt"): "0064 55\n0065 17\n0066 3\n",
		filepath.Join(tempDir, "iRO/official/recvpackets.txt"):             "0064 60\n0065 20\n0066 5\n",
		filepath.Join(tempDir, "iRO/recvpackets.txt"):                      "0064 65\n0065 25\n0066 8\n",
		filepath.Join(tempDir, "recvpackets.txt"):                          "0064 70\n0065 30\n0066 10\n",
	}

	for path, content := range testFiles {
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	tests := []struct {
		name         string
		tableFolders []string
		wantDefs     map[string]PacketLengthDef
		wantErr      bool
	}{
		{
			name:         "Load from first folder",
			tableFolders: []string{"kRO/RagexeRE_2020_04_01b", "iRO/official", "iRO"},
			wantDefs: map[string]PacketLengthDef{
				"0064": {Length: 55, HasLength: false},
				"0065": {Length: 17, HasLength: false},
				"0066": {Length: 3, HasLength: false},
			},
			wantErr: false,
		},
		{
			name:         "Load from second folder",
			tableFolders: []string{"nonexistent", "iRO/official", "iRO"},
			wantDefs: map[string]PacketLengthDef{
				"0064": {Length: 60, HasLength: false},
				"0065": {Length: 20, HasLength: false},
				"0066": {Length: 5, HasLength: false},
			},
			wantErr: false,
		},
		{
			name:         "Load from third folder",
			tableFolders: []string{"nonexistent1", "nonexistent2", "iRO"},
			wantDefs: map[string]PacketLengthDef{
				"0064": {Length: 65, HasLength: false},
				"0065": {Length: 25, HasLength: false},
				"0066": {Length: 8, HasLength: false},
			},
			wantErr: false,
		},
		{
			name:         "Load from base folder",
			tableFolders: []string{"nonexistent1", "nonexistent2", "nonexistent3"},
			wantDefs: map[string]PacketLengthDef{
				"0064": {Length: 70, HasLength: false},
				"0065": {Length: 30, HasLength: false},
				"0066": {Length: 10, HasLength: false},
			},
			wantErr: false,
		},
		{
			name:         "File not found",
			tableFolders: []string{"nonexistent1", "nonexistent2", "nonexistent3"},
			wantDefs:     nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// In the last test case, rename recvpackets.txt to force an error
			if tt.wantErr {
				err := os.Rename(filepath.Join(tempDir, "recvpackets.txt"), filepath.Join(tempDir, "recvpackets.txt.bak"))
				if err != nil {
					t.Fatalf("Failed to rename test file: %v", err)
				}
				defer os.Rename(filepath.Join(tempDir, "recvpackets.txt.bak"), filepath.Join(tempDir, "recvpackets.txt"))
			}

			gotDefs, err := LoadRecvPackets(tempDir, tt.tableFolders)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadRecvPackets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(gotDefs, tt.wantDefs) {
				t.Errorf("LoadRecvPackets() = %v, want %v", gotDefs, tt.wantDefs)
			}
		})
	}
}
