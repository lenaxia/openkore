package packets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPacketDefinitionsFromFile(t *testing.T) {
	// Create a temporary file with packet definitions
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "packet_db.txt")

	// Write test packet definitions to the file
	content := `
# Packet definitions for testing
# Format: PacketID, PacketName, PacketFormat, ParamName1, ParamName2, ...

# Login packets
0069 account_server_info v a4 a4 a4 a4 a26 C a* len sessionID accountID sessionID2 lastLoginIP lastLoginTime accountSex serverInfo
006A login_error C Z20 type date

# Character packets
0071 received_character_ID_and_Map a4 Z16 a4 v charID mapName mapIP mapPort
0072 received_characters v a* len charInfo
`
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	// Load packet definitions from the file
	db, err := LoadPacketDefinitionsFromFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to load packet definitions: %v", err)
	}

	// Verify that the packet definitions were loaded correctly
	testCases := []struct {
		id   string
		name string
	}{
		{"0069", "account_server_info"},
		{"006A", "login_error"},
		{"0071", "received_character_ID_and_Map"},
		{"0072", "received_characters"},
	}

	for _, tc := range testCases {
		def, exists := db.GetPacketByID(tc.id)
		if !exists {
			t.Errorf("Expected packet definition %s to exist", tc.id)
			continue
		}

		if def.Name != tc.name {
			t.Errorf("Expected packet %s to have name %s, got %s", tc.id, tc.name, def.Name)
		}

		defByName, exists := db.GetPacketByName(tc.name)
		if !exists {
			t.Errorf("Expected packet definition %s to exist", tc.name)
			continue
		}

		if defByName.ID != tc.id {
			t.Errorf("Expected packet %s to have ID %s, got %s", tc.name, tc.id, defByName.ID)
		}
	}
}

func TestLoadPacketDefinitionsFromInvalidFile(t *testing.T) {
	// Test loading from a non-existent file
	_, err := LoadPacketDefinitionsFromFile("non_existent_file.txt")
	if err == nil {
		t.Error("Expected error when loading from non-existent file, got nil")
	}

	// Create a temporary file with invalid packet definitions
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "invalid_packet_db.txt")

	// Write invalid test packet definitions to the file
	content := `
# Invalid packet definitions for testing
# Missing fields
0069 account_server_info
# Invalid format
006A login_error INVALID_FORMAT type date
`
	err = os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	// Load packet definitions from the file
	db, err := LoadPacketDefinitionsFromFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to load packet definitions: %v", err)
	}

	// Verify that the invalid packet definitions were skipped
	_, exists := db.GetPacketByID("0069")
	if exists {
		t.Error("Expected invalid packet definition 0069 to be skipped")
	}

	_, exists = db.GetPacketByID("006A")
	if exists {
		t.Error("Expected invalid packet definition 006A to be skipped")
	}
}

func TestMergePacketDatabases(t *testing.T) {
	// Create two packet databases
	db1 := NewPacketDatabase()
	db2 := NewPacketDatabase()

	// Add packet definitions to the first database
	db1.AddPacketDefinition(NewPacketDefinition("0069", "account_server_info", "v a4 a4 a4 a4 a26 C a*", []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"}))
	db1.AddPacketDefinition(NewPacketDefinition("006A", "login_error", "C Z20", []string{"type", "date"}))

	// Add packet definitions to the second database
	db2.AddPacketDefinition(NewPacketDefinition("0071", "received_character_ID_and_Map", "a4 Z16 a4 v", []string{"charID", "mapName", "mapIP", "mapPort"}))
	db2.AddPacketDefinition(NewPacketDefinition("0072", "received_characters", "v a*", []string{"len", "charInfo"}))

	// Add a packet definition with the same ID but different name to the second database
	db2.AddPacketDefinition(NewPacketDefinition("0069", "different_name", "v a4", []string{"len", "data"}))

	// Merge the databases
	merged := MergePacketDatabases(db1, db2)

	// Verify that all packet definitions were merged correctly
	testCases := []struct {
		id   string
		name string
	}{
		{"0069", "different_name"}, // Second database should override the first
		{"006A", "login_error"},
		{"0071", "received_character_ID_and_Map"},
		{"0072", "received_characters"},
	}

	for _, tc := range testCases {
		def, exists := merged.GetPacketByID(tc.id)
		if !exists {
			t.Errorf("Expected packet definition %s to exist", tc.id)
			continue
		}

		if def.Name != tc.name {
			t.Errorf("Expected packet %s to have name %s, got %s", tc.id, tc.name, def.Name)
		}
	}
}

func TestLoadServerSpecificPacketDefinitions(t *testing.T) {
	// Create a temporary directory structure for server-specific packet definitions
	tempDir := t.TempDir()
	baseDir := filepath.Join(tempDir, "packets")
	err := os.MkdirAll(baseDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	// Create base packet definitions file
	baseFile := filepath.Join(baseDir, "base.txt")
	baseContent := `
# Base packet definitions
0069 account_server_info v a4 a4 a4 a4 a26 C a* len sessionID accountID sessionID2 lastLoginIP lastLoginTime accountSex serverInfo
006A login_error C Z20 type date
`
	err = os.WriteFile(baseFile, []byte(baseContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create base file: %v", err)
	}

	// Create server-specific packet definitions file
	serverDir := filepath.Join(baseDir, "serverType1")
	err = os.MkdirAll(serverDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create server directory: %v", err)
	}

	serverFile := filepath.Join(serverDir, "packets.txt")
	serverContent := `
# Server-specific packet definitions
# Override base packet
0069 account_server_info_v2 v a4 a4 a4 a4 a26 C a* v len sessionID accountID sessionID2 lastLoginIP lastLoginTime accountSex serverInfo version
# Add new packet
0071 received_character_ID_and_Map a4 Z16 a4 v charID mapName mapIP mapPort
`
	err = os.WriteFile(serverFile, []byte(serverContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create server file: %v", err)
	}

	// Load packet definitions
	db, err := LoadServerSpecificPacketDefinitions(baseFile, "serverType1", baseDir)
	if err != nil {
		t.Fatalf("Failed to load server-specific packet definitions: %v", err)
	}

	// Verify that the packet definitions were loaded correctly
	testCases := []struct {
		id   string
		name string
	}{
		{"0069", "account_server_info_v2"},        // Server-specific definition should override base
		{"006A", "login_error"},                   // Base definition should be preserved
		{"0071", "received_character_ID_and_Map"}, // Server-specific definition should be added
	}

	for _, tc := range testCases {
		def, exists := db.GetPacketByID(tc.id)
		if !exists {
			t.Errorf("Expected packet definition %s to exist", tc.id)
			continue
		}

		if def.Name != tc.name {
			t.Errorf("Expected packet %s to have name %s, got %s", tc.id, tc.name, def.Name)
		}
	}
}
