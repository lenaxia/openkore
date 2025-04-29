package packets

import (
	"testing"
)

func TestPacketDefinition(t *testing.T) {
	// Test creating a packet definition
	def := NewPacketDefinition("0069", "account_server_info", "v a4 a4 a4 a4 a26 C a*", []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"})

	if def.ID != "0069" {
		t.Errorf("Expected ID to be 0069, got %s", def.ID)
	}

	if def.Name != "account_server_info" {
		t.Errorf("Expected Name to be account_server_info, got %s", def.Name)
	}

	if def.Format != "v a4 a4 a4 a4 a26 C a*" {
		t.Errorf("Expected Format to be 'v a4 a4 a4 a4 a26 C a*', got %s", def.Format)
	}

	if len(def.ParamNames) != 8 {
		t.Errorf("Expected 8 parameter names, got %d", len(def.ParamNames))
	}

	expectedParams := []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"}
	for i, param := range def.ParamNames {
		if param != expectedParams[i] {
			t.Errorf("Expected parameter %d to be %s, got %s", i, expectedParams[i], param)
		}
	}
}

func TestPacketDatabase(t *testing.T) {
	// Test creating a packet database
	db := NewPacketDatabase()

	// Test adding a packet definition
	def := NewPacketDefinition("0069", "account_server_info", "v a4 a4 a4 a4 a26 C a*", []string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"})
	db.AddPacketDefinition(def)

	// Test getting a packet definition by ID
	retrievedDef, exists := db.GetPacketByID("0069")
	if !exists {
		t.Error("Expected packet definition to exist")
	}

	if retrievedDef.Name != "account_server_info" {
		t.Errorf("Expected Name to be account_server_info, got %s", retrievedDef.Name)
	}

	// Test getting a packet definition by name
	retrievedDef, exists = db.GetPacketByName("account_server_info")
	if !exists {
		t.Error("Expected packet definition to exist")
	}

	if retrievedDef.ID != "0069" {
		t.Errorf("Expected ID to be 0069, got %s", retrievedDef.ID)
	}

	// Test getting a non-existent packet
	_, exists = db.GetPacketByID("9999")
	if exists {
		t.Error("Expected packet definition to not exist")
	}

	_, exists = db.GetPacketByName("non_existent_packet")
	if exists {
		t.Error("Expected packet definition to not exist")
	}
}

func TestDefaultPacketDatabase(t *testing.T) {
	// Test that the default packet database is initialized with some packet definitions
	db := NewDefaultPacketDatabase()

	// Check for a few common packet definitions
	packets := []struct {
		id   string
		name string
	}{
		{"0069", "account_server_info"},
		{"006A", "login_error"},
		{"0071", "received_character_ID_and_Map"},
		{"0073", "map_loaded"},
		{"0078", "actor_exists"},
		{"0079", "actor_connected"},
		{"007B", "actor_moved"},
		{"008D", "public_chat"},
		{"008E", "self_chat"},
		{"0091", "map_change"},
		{"0092", "map_changed"},
	}

	for _, p := range packets {
		def, exists := db.GetPacketByID(p.id)
		if !exists {
			t.Errorf("Expected packet definition %s to exist", p.id)
			continue
		}

		if def.Name != p.name {
			t.Errorf("Expected packet %s to have name %s, got %s", p.id, p.name, def.Name)
		}

		defByName, exists := db.GetPacketByName(p.name)
		if !exists {
			t.Errorf("Expected packet definition %s to exist", p.name)
			continue
		}

		if defByName.ID != p.id {
			t.Errorf("Expected packet %s to have ID %s, got %s", p.name, p.id, defByName.ID)
		}
	}
}

func TestPacketLength(t *testing.T) {
	// Test packet length calculation
	db := NewDefaultPacketDatabase()

	// Test fixed length packet
	def, exists := db.GetPacketByID("0073") // map_loaded
	if !exists {
		t.Fatal("Expected packet definition 0073 to exist")
	}

	length := def.GetLength()
	if length != 11 {
		t.Errorf("Expected packet 0073 to have length 11, got %d", length)
	}

	// Test variable length packet
	def, exists = db.GetPacketByID("0069") // account_server_info
	if !exists {
		t.Fatal("Expected packet definition 0069 to exist")
	}

	length = def.GetLength()
	if length != -1 {
		t.Errorf("Expected packet 0069 to have variable length (-1), got %d", length)
	}
}
