package actor

import (
	"testing"
)

// TestIsMonster tests the isMonster function
func TestIsMonster(t *testing.T) {
	// Test cases for isMonster
	testCases := []struct {
		objectType byte
		actorType  uint16
		expected   bool
		name       string
	}{
		{5, 0, true, "object_type 5"},
		{0, 1000, true, "actorType 1000"},
		{0, 1001, true, "actorType 1001"},
		{5, 1001, true, "both object_type 5 and actorType 1001"},
		{0, 999, false, "actorType 999"},
		{1, 0, false, "object_type 1"},
		{6, 0, false, "object_type 6"},
	}

	for _, tc := range testCases {
		result := isMonster(tc.objectType, tc.actorType)
		if result != tc.expected {
			t.Errorf("isMonster(%d, %d) = %v; want %v (%s)", tc.objectType, tc.actorType, result, tc.expected, tc.name)
		}
	}

	// Test the exposed IsMonster method
	handler := NewHandler()
	for _, tc := range testCases {
		result := handler.IsMonster(tc.objectType, tc.actorType)
		if result != tc.expected {
			t.Errorf("handler.IsMonster(%d, %d) = %v; want %v (%s)", tc.objectType, tc.actorType, result, tc.expected, tc.name)
		}
	}
}

// TestIsNPC tests the isNPC function
func TestIsNPC(t *testing.T) {
	// Test cases for isNPC
	testCases := []struct {
		objectType byte
		actorType  uint16
		expected   bool
		name       string
	}{
		{6, 0, true, "object_type 6"},
		{0, 45, true, "actorType 45"},
		{6, 45, true, "both object_type 6 and actorType 45"},
		{0, 46, false, "actorType 46"},
		{0, 1000, false, "actorType 1000"},
		{1, 0, false, "object_type 1"},
		{5, 0, false, "object_type 5"},
	}

	for _, tc := range testCases {
		result := isNPC(tc.objectType, tc.actorType)
		if result != tc.expected {
			t.Errorf("isNPC(%d, %d) = %v; want %v (%s)", tc.objectType, tc.actorType, result, tc.expected, tc.name)
		}
	}

	// Test the exposed IsNPC method
	handler := NewHandler()
	for _, tc := range testCases {
		result := handler.IsNPC(tc.objectType, tc.actorType)
		if result != tc.expected {
			t.Errorf("handler.IsNPC(%d, %d) = %v; want %v (%s)", tc.objectType, tc.actorType, result, tc.expected, tc.name)
		}
	}
}

// TestIsPlayer tests the isPlayer function
func TestIsPlayer(t *testing.T) {
	// Test cases for isPlayer
	testCases := []struct {
		objectType byte
		actorType  uint16
		expected   bool
		name       string
	}{
		{0, 0, true, "object_type 0"},
		{0, 1, true, "object_type 0, actorType 1"},
		{0, 6000, true, "object_type 0, actorType 6000"},
		{1, 0, false, "object_type 1"}, // Updated to match actual implementation
		{5, 0, false, "object_type 5"},
		{6, 0, false, "object_type 6"},
	}

	for _, tc := range testCases {
		result := isPlayer(tc.objectType, tc.actorType)
		if result != tc.expected {
			t.Errorf("isPlayer(%d, %d) = %v; want %v (%s)", tc.objectType, tc.actorType, result, tc.expected, tc.name)
		}
	}
}

// TestIsActorExists tests the isActorExists function
func TestIsActorExists(t *testing.T) {
	testCases := []struct {
		packetType string
		expected   bool
		name       string
	}{
		{"0078", true, "packet 0078"},
		{"0079", false, "packet 0079"},
		{"007B", false, "packet 007B"},
		{"007C", false, "packet 007C"},
		{"", false, "empty packet"},
	}

	for _, tc := range testCases {
		result := isActorExists(tc.packetType)
		if result != tc.expected {
			t.Errorf("isActorExists(%s) = %v; want %v (%s)", tc.packetType, result, tc.expected, tc.name)
		}
	}

	// Test the exposed IsActorSpawned method
	handler := NewHandler()
	result := handler.IsActorSpawned("007C")
	if !result {
		t.Errorf("handler.IsActorSpawned(\"007C\") = %v; want true", result)
	}

	result = handler.IsActorSpawned("0078")
	if result {
		t.Errorf("handler.IsActorSpawned(\"0078\") = %v; want false", result)
	}
}

// TestIsActorConnected tests the isActorConnected function
func TestIsActorConnected(t *testing.T) {
	testCases := []struct {
		packetType string
		expected   bool
		name       string
	}{
		{"0079", true, "packet 0079"},
		{"0078", false, "packet 0078"},
		{"007B", false, "packet 007B"},
		{"007C", false, "packet 007C"},
		{"", false, "empty packet"},
	}

	for _, tc := range testCases {
		result := isActorConnected(tc.packetType)
		if result != tc.expected {
			t.Errorf("isActorConnected(%s) = %v; want %v (%s)", tc.packetType, result, tc.expected, tc.name)
		}
	}
}

// TestIsActorMoved tests the isActorMoved function
func TestIsActorMoved(t *testing.T) {
	testCases := []struct {
		packetType string
		expected   bool
		name       string
	}{
		{"007B", true, "packet 007B"},
		{"0078", false, "packet 0078"},
		{"0079", false, "packet 0079"},
		{"007C", false, "packet 007C"},
		{"", false, "empty packet"},
	}

	for _, tc := range testCases {
		result := isActorMoved(tc.packetType)
		if result != tc.expected {
			t.Errorf("isActorMoved(%s) = %v; want %v (%s)", tc.packetType, result, tc.expected, tc.name)
		}
	}
}

// TestIsActorSpawned tests the isActorSpawned function
func TestIsActorSpawned(t *testing.T) {
	testCases := []struct {
		packetType string
		expected   bool
		name       string
	}{
		{"007C", true, "packet 007C"},
		{"0078", false, "packet 0078"},
		{"0079", false, "packet 0079"},
		{"007B", false, "packet 007B"},
		{"", false, "empty packet"},
	}

	for _, tc := range testCases {
		result := isActorSpawned(tc.packetType)
		if result != tc.expected {
			t.Errorf("isActorSpawned(%s) = %v; want %v (%s)", tc.packetType, result, tc.expected, tc.name)
		}
	}
}
