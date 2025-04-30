package factory

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lenaxia/goKore/network/config"
	"github.com/lenaxia/goKore/network/protocol"
)

func TestServerTypeFactory_LoadPacketDefinitions(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Create test folder structure
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

	// Create test recvpackets.txt files with different content
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
		serverType   config.ServerType
		tableFolders []string
		wantDefs     map[string]protocol.PacketLengthDef
		wantErr      bool
	}{
		{
			name:         "ServerType0 with kRO folders",
			serverType:   config.ServerType0,
			tableFolders: []string{"kRO/RagexeRE_2020_04_01b", "iRO/official", "iRO"},
			wantDefs: map[string]protocol.PacketLengthDef{
				"0064": {Length: 55, HasLength: false},
				"0065": {Length: 17, HasLength: false},
				"0066": {Length: 3, HasLength: false},
			},
			wantErr: false,
		},
		{
			name:         "ServerTypeSakray with iRO/official folders",
			serverType:   config.ServerTypeSakray,
			tableFolders: []string{"iRO/official", "iRO"},
			wantDefs: map[string]protocol.PacketLengthDef{
				"0064": {Length: 60, HasLength: false},
				"0065": {Length: 20, HasLength: false},
				"0066": {Length: 5, HasLength: false},
			},
			wantErr: false,
		},
		{
			name:         "ServerTypeIRO with iRO folders",
			serverType:   config.ServerTypeIRO,
			tableFolders: []string{"iRO"},
			wantDefs: map[string]protocol.PacketLengthDef{
				"0064": {Length: 65, HasLength: false},
				"0065": {Length: 25, HasLength: false},
				"0066": {Length: 8, HasLength: false},
			},
			wantErr: false,
		},
		{
			name:         "Unknown server type with base folder",
			serverType:   config.ServerTypeUnknown,
			tableFolders: []string{},
			wantDefs: map[string]protocol.PacketLengthDef{
				"0064": {Length: 70, HasLength: false},
				"0065": {Length: 30, HasLength: false},
				"0066": {Length: 10, HasLength: false},
			},
			wantErr: false,
		},
		{
			name:         "Non-existent folders",
			serverType:   config.ServerTypeUnknown,
			tableFolders: []string{"nonexistent1", "nonexistent2"},
			wantDefs:     nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create server config
			serverConfig := &config.ServerConfig{
				Type:         tt.serverType,
				TableFolders: tt.tableFolders,
			}

			// Create factory
			factory := NewServerTypeFactory()

			// Load packet definitions
			packetDefs, err := factory.LoadPacketDefinitions(tempDir, serverConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPacketDefinitions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Check packet definitions
			for id, wantDef := range tt.wantDefs {
				gotDef, exists := packetDefs[id]
				if !exists {
					t.Errorf("LoadPacketDefinitions() missing packet definition for %s", id)
					continue
				}

				if gotDef.Length != wantDef.Length {
					t.Errorf("LoadPacketDefinitions() packet %s length = %d, want %d", id, gotDef.Length, wantDef.Length)
				}

				if gotDef.HasLength != wantDef.HasLength {
					t.Errorf("LoadPacketDefinitions() packet %s hasLength = %v, want %v", id, gotDef.HasLength, wantDef.HasLength)
				}
			}
		})
	}
}

func TestServerTypeFactory_CreateTokenizer(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Create test recvpackets.txt file
	testFilePath := filepath.Join(tempDir, "recvpackets.txt")
	testContent := "0064 55\n0065 17\n0066 3\n"
	err := os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create server config
	serverConfig := &config.ServerConfig{
		Type:         config.ServerType0,
		TableFolders: []string{},
	}

	// Create factory
	factory := NewServerTypeFactory()

	// Create tokenizer
	tokenizer, err := factory.CreateTokenizer(tempDir, serverConfig)
	if err != nil {
		t.Fatalf("CreateTokenizer() error = %v", err)
	}

	// Test tokenizer with a packet
	packet := make([]byte, 55)
	packet[0] = 0x64 // Little-endian byte order
	packet[1] = 0x00

	tokenizer.Add(packet)

	result, msgType, err := tokenizer.ReadNext()
	if err != nil {
		t.Errorf("ReadNext() error = %v", err)
	}

	if msgType != protocol.KnownMessage {
		t.Errorf("ReadNext() msgType = %v, want %v", msgType, protocol.KnownMessage)
	}

	if len(result) != 55 {
		t.Errorf("ReadNext() result length = %d, want %d", len(result), 55)
	}
}
