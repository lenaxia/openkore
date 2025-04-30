package servers

import (
	"testing"
)

// Happy path - all expected packet constructions are returned
func TestServerType0PacketConstructions(t *testing.T) {
	// Get the packet constructions
	constructions := ServerType0PacketConstructions()

	// Check that the packet constructions were returned
	expectedPackets := []string{
		"0064", // login_request
		"0085", // move_to
		"00B2", // restart
	}

	for _, id := range expectedPackets {
		if _, exists := constructions[id]; !exists {
			t.Errorf("Packet construction %s was not returned", id)
		}
	}

	// Check that the packet constructions have the correct format
	if constructions["0064"].Format != "v a24 a24 C" {
		t.Errorf("Packet construction 0064 has incorrect format: %s", constructions["0064"].Format)
	}

	if constructions["0085"].Format != "v3" {
		t.Errorf("Packet construction 0085 has incorrect format: %s", constructions["0085"].Format)
	}

	if constructions["00B2"].Format != "C" {
		t.Errorf("Packet construction 00B2 has incorrect format: %s", constructions["00B2"].Format)
	}

	// Check field names for a packet construction
	expectedFieldNames := []string{"version", "username", "password", "clienttype"}
	for i, fieldName := range constructions["0064"].FieldNames {
		if i >= len(expectedFieldNames) {
			t.Errorf("Unexpected field name at index %d: %s", i, fieldName)
			continue
		}
		if fieldName != expectedFieldNames[i] {
			t.Errorf("Packet construction 0064 has incorrect field name at index %d: expected %s, got %s",
				i, expectedFieldNames[i], fieldName)
		}
	}

	// Check that the packet IDs match the names
	if constructions["0064"].Name != "login_request" {
		t.Errorf("Packet construction 0064 has incorrect name: %s", constructions["0064"].Name)
	}

	if constructions["0085"].Name != "move_to" {
		t.Errorf("Packet construction 0085 has incorrect name: %s", constructions["0085"].Name)
	}

	if constructions["00B2"].Name != "restart" {
		t.Errorf("Packet construction 00B2 has incorrect name: %s", constructions["00B2"].Name)
	}
}

// Edge case - check for duplicate packet names
func TestServerType0PacketConstructions_NoDuplicateNames(t *testing.T) {
	constructions := ServerType0PacketConstructions()

	// Create a map to track packet names
	nameCount := make(map[string]int)

	// Count occurrences of each name
	for _, pc := range constructions {
		nameCount[pc.Name]++
	}

	// Check for duplicates
	for name, count := range nameCount {
		if count > 1 {
			t.Errorf("Duplicate packet name found: %s appears %d times", name, count)
		}
	}
}

// Happy path - all expected packet constructions are returned
func TestServerType1PacketConstructions(t *testing.T) {
	// Get the packet constructions
	constructions := ServerType1PacketConstructions()

	// Check that the packet constructions were returned and have the expected format
	if len(constructions) == 0 {
		t.Error("No packet constructions were returned")
	}

	// Check that packet constructions have the required fields
	for id, pc := range constructions {
		if pc.ID == "" {
			t.Errorf("Packet construction %s has empty ID", id)
		}

		if pc.Name == "" {
			t.Errorf("Packet construction %s has empty name", id)
		}

		// Format can be empty for some packets, so we don't check it

		// Some packets might not have field names, so we don't check for empty slices
	}
}

// Happy path - all expected packet constructions are returned
func TestSakrayPacketConstructions(t *testing.T) {
	// Get the packet constructions
	constructions := SakrayPacketConstructions()

	// Check that the packet constructions were returned and have the expected format
	if len(constructions) == 0 {
		t.Error("No packet constructions were returned")
	}

	// Check that packet constructions have the required fields
	for id, pc := range constructions {
		if pc.ID == "" {
			t.Errorf("Packet construction %s has empty ID", id)
		}

		if pc.Name == "" {
			t.Errorf("Packet construction %s has empty name", id)
		}
	}
}

// Test consistency between server types
func TestPacketConstructionConsistency(t *testing.T) {
	// Get packet constructions from all server types
	type0 := ServerType0PacketConstructions()
	type1 := ServerType1PacketConstructions()
	sakray := SakrayPacketConstructions()

	// Find common packet IDs manually
	commonIDs := []string{}

	// Check each ID in type0
	for id := range type0 {
		// If it exists in all three maps, add it to commonIDs
		if _, exists1 := type1[id]; exists1 {
			if _, exists2 := sakray[id]; exists2 {
				commonIDs = append(commonIDs, id)
			}
		}
	}

	for _, id := range commonIDs {
		// Check that the packet name is consistent across server types
		if type0[id].Name != type1[id].Name || type0[id].Name != sakray[id].Name {
			t.Errorf("Inconsistent packet name for ID %s: type0=%s, type1=%s, sakray=%s",
				id, type0[id].Name, type1[id].Name, sakray[id].Name)
		}

		// Check that the field names are consistent across server types
		// (only if they have the same format)
		if type0[id].Format == type1[id].Format && type0[id].Format == sakray[id].Format {
			if !compareStringSlices(type0[id].FieldNames, type1[id].FieldNames) ||
				!compareStringSlices(type0[id].FieldNames, sakray[id].FieldNames) {
				t.Errorf("Inconsistent field names for ID %s with same format", id)
			}
		}
	}
}

// Helper function to compare string slices (ignoring order)
func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := make(map[string]int)
	for _, s := range a {
		aMap[s]++
	}

	bMap := make(map[string]int)
	for _, s := range b {
		bMap[s]++
	}

	for s, count := range aMap {
		if bMap[s] != count {
			return false
		}
	}

	return true
}
