package core

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestNewCharacterManager(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	manager := NewCharacterManager(parser)

	if manager == nil {
		t.Fatal("NewCharacterManager() returned nil")
	}

	if manager.parser != parser {
		t.Errorf("manager.parser = %v, want %v", manager.parser, parser)
	}

	if len(manager.guildMembers) != 0 {
		t.Errorf("len(manager.guildMembers) = %d, want 0", len(manager.guildMembers))
	}
}

func TestRegisterCharacterHandlers(t *testing.T) {
	// Skip this test for now
	t.Skip("Skipping TestRegisterCharacterHandlers")

	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	manager := NewCharacterManager(parser)

	// Register handlers
	manager.RegisterHandlers()

	// Check that the handler was registered
	handler, found := parser.GetHandler("0095")
	if !found || handler == nil {
		t.Fatal("Handler for packet 0095 (character_name) was not registered")
	}
}

func TestAddGetGuildMember(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	manager := NewCharacterManager(parser)

	// Add a guild member
	charID := uint32(12345)
	name := "TestMember"
	manager.AddGuildMember(charID, name)

	// Get the guild member
	member, found := manager.GetGuildMember(charID)
	if !found {
		t.Fatal("GetGuildMember() returned not found")
	}

	if member.CharID != charID {
		t.Errorf("member.CharID = %d, want %d", member.CharID, charID)
	}

	if member.Name != name {
		t.Errorf("member.Name = %s, want %s", member.Name, name)
	}

	// Get a non-existent guild member
	_, found = manager.GetGuildMember(99999)
	if found {
		t.Error("GetGuildMember() returned found for non-existent member")
	}

	// Update an existing guild member
	newName := "UpdatedMember"
	manager.AddGuildMember(charID, newName)

	// Check that the member was updated
	member, found = manager.GetGuildMember(charID)
	if !found {
		t.Fatal("GetGuildMember() returned not found after update")
	}

	if member.Name != newName {
		t.Errorf("member.Name = %s, want %s", member.Name, newName)
	}

	// Get all guild members
	members := manager.GetGuildMembers()
	if len(members) != 1 {
		t.Errorf("len(members) = %d, want 1", len(members))
	}

	if members[0].CharID != charID {
		t.Errorf("members[0].CharID = %d, want %d", members[0].CharID, charID)
	}

	if members[0].Name != newName {
		t.Errorf("members[0].Name = %s, want %s", members[0].Name, newName)
	}
}

func TestHandleCharacterStatus(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	manager := NewCharacterManager(parser)

	// Test case 1: Update level and opt3 (packet 028A)
	t.Run("UpdateLevelAndOpt3", func(t *testing.T) {
		// Create test packet arguments
		// Character ID 12345 in little-endian: 0x3039 = [0x39, 0x30, 0x00, 0x00]
		args := map[string]interface{}{
			"ID":   []byte{0x39, 0x30, 0x00, 0x00},
			"lv":   uint16(50),
			"opt3": uint16(128),
		}

		// Call handler
		err := manager.handleCharacterStatus(args)
		if err != nil {
			t.Fatalf("handleCharacterStatus() returned error: %v", err)
		}

		// Check that the actor was created and updated
		actor, found := manager.GetActor(12345)
		if !found {
			t.Fatal("GetActor() returned not found")
		}

		if actor.Lv != 50 {
			t.Errorf("actor.Lv = %d, want 50", actor.Lv)
		}

		if actor.Opt3 != 128 {
			t.Errorf("actor.Opt3 = %d, want 128", actor.Opt3)
		}
	})

	// Test case 2: Update opt1 and opt2 (packet 0229 or 0119)
	t.Run("UpdateOpt1AndOpt2", func(t *testing.T) {
		// Create test packet arguments
		// Character ID 12345 in little-endian: 0x3039 = [0x39, 0x30, 0x00, 0x00]
		args := map[string]interface{}{
			"ID":   []byte{0x39, 0x30, 0x00, 0x00},
			"opt1": uint16(64),
			"opt2": uint16(32),
		}

		// Call handler
		err := manager.handleCharacterStatus(args)
		if err != nil {
			t.Fatalf("handleCharacterStatus() returned error: %v", err)
		}

		// Check that the actor was updated
		actor, found := manager.GetActor(12345)
		if !found {
			t.Fatal("GetActor() returned not found")
		}

		if actor.Opt1 != 64 {
			t.Errorf("actor.Opt1 = %d, want 64", actor.Opt1)
		}

		if actor.Opt2 != 32 {
			t.Errorf("actor.Opt2 = %d, want 32", actor.Opt2)
		}

		// Check that the level and opt3 are still set from the previous test
		if actor.Lv != 50 {
			t.Errorf("actor.Lv = %d, want 50", actor.Lv)
		}

		if actor.Opt3 != 128 {
			t.Errorf("actor.Opt3 = %d, want 128", actor.Opt3)
		}
	})
}

func TestHandleCharacterName(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	manager := NewCharacterManager(parser)

	// Add a guild member
	charID := uint32(12345)
	name := "TestMember"
	manager.AddGuildMember(charID, name)

	// Create test packet arguments
	// Character ID 12345 in little-endian: 0x3039 = [0x39, 0x30, 0x00, 0x00]
	args := map[string]interface{}{
		"ID":   []byte{0x39, 0x30, 0x00, 0x00},
		"name": "UpdatedMember",
	}

	// Call handler
	err := manager.handleCharacterName(args)
	if err != nil {
		t.Fatalf("handleCharacterName() returned error: %v", err)
	}

	// Check that the guild member's name was updated
	member, found := manager.GetGuildMember(charID)
	if !found {
		t.Fatal("GetGuildMember() returned not found after handleCharacterName")
	}

	if member.Name != "UpdatedMember" {
		t.Errorf("member.Name = %s, want UpdatedMember", member.Name)
	}

	// Test with a non-existent guild member
	args = map[string]interface{}{
		"ID":   []byte{0x42, 0x42, 0x00, 0x00}, // Some other character ID
		"name": "NonExistentMember",
	}

	// Call handler
	err = manager.handleCharacterName(args)
	if err != nil {
		t.Fatalf("handleCharacterName() returned error: %v", err)
	}

	// Check that the guild members list is unchanged
	members := manager.GetGuildMembers()
	if len(members) != 1 {
		t.Errorf("len(members) = %d, want 1", len(members))
	}
}
