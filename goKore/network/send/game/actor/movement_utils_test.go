package actor

import (
	"testing"
)

// TestParseActorMove tests the ParseActorMove function
func TestParseActorMove(t *testing.T) {
	// Create a test case
	args := map[string]interface{}{
		"coords": []byte{0x12, 0x34, 0x56, 0x78}, // Example coordinate bytes
	}

	// Call the function
	ParseActorMove(args)

	// Check that the coordinates were parsed correctly
	if args["x"] == nil || args["y"] == nil {
		t.Fatal("ParseActorMove did not set x and y coordinates")
	}

	// The actual values depend on the implementation of makeCoordsDir
	// We're just checking that they exist for now
}

// TestReconstructActorMove tests the ReconstructActorMove function
func TestReconstructActorMove(t *testing.T) {
	// Test with no_padding = false
	args := map[string]interface{}{
		"x": 123,
		"y": 456,
	}

	// Call the function
	ReconstructActorMove(args)

	// Check that the coords field was set
	if args["coords"] == nil {
		t.Fatal("ReconstructActorMove did not set coords")
	}

	// Test with no_padding = true
	args = map[string]interface{}{
		"x":          123,
		"y":          456,
		"no_padding": true,
	}

	// Call the function
	ReconstructActorMove(args)

	// Check that the coords field was set
	if args["coords"] == nil {
		t.Fatal("ReconstructActorMove did not set coords")
	}
}

// TestGetCoordString tests the GetCoordString function
func TestGetCoordString(t *testing.T) {
	// Test with no padding
	coords := GetCoordString(123, 456, true)
	if len(coords) != 3 {
		t.Errorf("GetCoordString with no padding returned %d bytes, expected 3", len(coords))
	}

	// Test with padding
	coords = GetCoordString(123, 456, false)
	if len(coords) != 3 {
		t.Errorf("GetCoordString with padding returned %d bytes, expected 3", len(coords))
	}
}

// TestMakeCoordsDir tests the MakeCoordsDir function
func TestMakeCoordsDir(t *testing.T) {
	// Create a test case
	args := map[string]interface{}{
		"coords": []byte{0x12, 0x34, 0x56}, // Example coordinate bytes
	}

	// Call the function
	MakeCoordsDir(args)

	// Check that the coordinates were parsed correctly
	if args["x"] == nil || args["y"] == nil {
		t.Fatal("MakeCoordsDir did not set x and y coordinates")
	}

	// The actual values depend on the implementation
	// We're just checking that they exist for now
}
