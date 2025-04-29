package core

import (
	"bytes"
	"testing"
)

// TestNewPacketBuilder tests the NewPacketBuilder function
func TestNewPacketBuilder(t *testing.T) {
	pb := NewPacketBuilder()
	if pb == nil {
		t.Fatal("NewPacketBuilder() returned nil")
	}

	if pb.packetDefinitions == nil {
		t.Error("pb.packetDefinitions is nil")
	}

	if pb.packetLUT == nil {
		t.Error("pb.packetLUT is nil")
	}
}

// TestRegisterPacket tests the RegisterPacket method
func TestRegisterPacket(t *testing.T) {
	pb := NewPacketBuilder()

	// Register a packet
	pb.RegisterPacket("0064", "master_login", "V Z24 Z24 C", []string{"version", "username", "password", "master_version"})

	// Check that the packet was registered
	if _, exists := pb.packetDefinitions["0064"]; !exists {
		t.Error("Packet was not registered in packetDefinitions")
	}

	// Check that the packet ID was registered in the lookup table
	id, exists := pb.GetPacketID("master_login")
	if !exists {
		t.Error("Packet ID was not registered in packetLUT")
	}
	if id != "0064" {
		t.Errorf("GetPacketID() = %v, want 0064", id)
	}
}

// TestGetPacketID tests the GetPacketID method
func TestGetPacketID(t *testing.T) {
	pb := NewPacketBuilder()

	// Register a packet
	pb.RegisterPacket("0064", "master_login", "V Z24 Z24 C", []string{"version", "username", "password", "master_version"})

	// Test getting a registered packet ID
	id, exists := pb.GetPacketID("master_login")
	if !exists {
		t.Error("GetPacketID() returned exists = false, want true")
	}
	if id != "0064" {
		t.Errorf("GetPacketID() = %v, want 0064", id)
	}

	// Test getting an unregistered packet ID
	id, exists = pb.GetPacketID("unknown_packet")
	if exists {
		t.Error("GetPacketID() returned exists = true, want false")
	}
	if id != "" {
		t.Errorf("GetPacketID() = %v, want \"\"", id)
	}
}

// TestBuildPacket tests the BuildPacket method
func TestBuildPacket(t *testing.T) {
	pb := NewPacketBuilder()

	// Register a packet
	pb.RegisterPacket("0064", "master_login", "V Z24 Z24 C", []string{"version", "username", "password", "master_version"})

	// Test building a packet
	args := map[string]interface{}{
		"version":        23,
		"username":       "testuser",
		"password":       "testpass",
		"master_version": 1,
	}

	packet, err := pb.BuildPacket("0064", args)
	if err != nil {
		t.Fatalf("BuildPacket() returned error: %v", err)
	}

	// Check that the packet has the correct ID
	if len(packet) < 2 {
		t.Fatalf("BuildPacket() returned packet of length %d, want at least 2", len(packet))
	}
	if packet[0] != 0x64 || packet[1] != 0x00 {
		t.Errorf("BuildPacket() returned packet with ID %02X%02X, want 0064", packet[1], packet[0])
	}

	// Test building an unregistered packet
	_, err = pb.BuildPacket("0065", args)
	if err == nil {
		t.Error("BuildPacket() did not return error for unregistered packet")
	}

	// Test building a packet with invalid ID
	_, err = pb.BuildPacket("006", args)
	if err == nil {
		t.Error("BuildPacket() did not return error for invalid packet ID")
	}

	// Test building a packet with missing arguments
	args = map[string]interface{}{
		"version":  23,
		"username": "testuser",
		// Missing password and master_version
	}
	_, err = pb.BuildPacket("0064", args)
	if err == nil {
		t.Error("BuildPacket() did not return error for missing arguments")
	}
}

// TestParseFormat tests the parseFormat function
func TestParseFormat(t *testing.T) {
	testCases := []struct {
		format string
		want   []string
	}{
		{"V Z24 Z24 C", []string{"V", "Z24", "Z24", "C"}},
		{"v2 a4 a4 V C", []string{"v2", "a4", "a4", "V", "C"}},
		{"", []string{}},
		{"C", []string{"C"}},
	}

	for _, tc := range testCases {
		got := parseFormat(tc.format)
		if len(got) != len(tc.want) {
			t.Errorf("parseFormat(%q) returned %d parts, want %d", tc.format, len(got), len(tc.want))
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("parseFormat(%q)[%d] = %q, want %q", tc.format, i, got[i], tc.want[i])
			}
		}
	}
}

// TestParseFormatPart tests the parseFormatPart function
func TestParseFormatPart(t *testing.T) {
	testCases := []struct {
		part      string
		wantType  string
		wantCount int
	}{
		{"C", "C", 1},
		{"v2", "v", 2},
		{"Z24", "Z", 24},
		{"a4", "a", 4},
		{"V", "V", 1},
		{"", "", 0},
	}

	for _, tc := range testCases {
		gotType, gotCount := parseFormatPart(tc.part)
		if gotType != tc.wantType {
			t.Errorf("parseFormatPart(%q) returned type %q, want %q", tc.part, gotType, tc.wantType)
		}
		if gotCount != tc.wantCount {
			t.Errorf("parseFormatPart(%q) returned count %d, want %d", tc.part, gotCount, tc.wantCount)
		}
	}
}

// TestGetIntValue tests the getIntValue function
func TestGetIntValue(t *testing.T) {
	testCases := []struct {
		value   interface{}
		index   int
		want    int64
		wantErr bool
	}{
		{123, 0, 123, false},
		{int8(123), 0, 123, false},
		{int16(123), 0, 123, false},
		{int32(123), 0, 123, false},
		{int64(123), 0, 123, false},
		{uint(123), 0, 123, false},
		{uint8(123), 0, 123, false},
		{uint16(123), 0, 123, false},
		{uint32(123), 0, 123, false},
		{uint64(123), 0, 123, false},
		{float32(123.45), 0, 123, false},
		{float64(123.45), 0, 123, false},
		{"123", 0, 123, false},
		{"0x7B", 0, 123, false},
		{[]int{123, 456}, 0, 123, false},
		{[]int{123, 456}, 1, 456, false},
		{[]int{123, 456}, 2, 0, true},
		{"abc", 0, 0, true},
		{struct{}{}, 0, 0, true},
	}

	for _, tc := range testCases {
		got, err := getIntValue(tc.value, tc.index)
		if (err != nil) != tc.wantErr {
			t.Errorf("getIntValue(%v, %d) error = %v, wantErr %v", tc.value, tc.index, err, tc.wantErr)
			continue
		}
		if !tc.wantErr && got != tc.want {
			t.Errorf("getIntValue(%v, %d) = %d, want %d", tc.value, tc.index, got, tc.want)
		}
	}
}

// TestGetFloatValue tests the getFloatValue function
func TestGetFloatValue(t *testing.T) {
	testCases := []struct {
		value   interface{}
		index   int
		want    float64
		wantErr bool
	}{
		{123, 0, 123.0, false},
		{int8(123), 0, 123.0, false},
		{int16(123), 0, 123.0, false},
		{int32(123), 0, 123.0, false},
		{int64(123), 0, 123.0, false},
		{uint(123), 0, 123.0, false},
		{uint8(123), 0, 123.0, false},
		{uint16(123), 0, 123.0, false},
		{uint32(123), 0, 123.0, false},
		{uint64(123), 0, 123.0, false},
		{float32(123.45), 0, 123.45, false},
		{float64(123.45), 0, 123.45, false},
		{"123.45", 0, 123.45, false},
		{[]float64{123.45, 456.78}, 0, 123.45, false},
		{[]float64{123.45, 456.78}, 1, 456.78, false},
		{[]float64{123.45, 456.78}, 2, 0, true},
		{"abc", 0, 0, true},
		{struct{}{}, 0, 0, true},
	}

	for _, tc := range testCases {
		got, err := getFloatValue(tc.value, tc.index)
		if (err != nil) != tc.wantErr {
			t.Errorf("getFloatValue(%v, %d) error = %v, wantErr %v", tc.value, tc.index, err, tc.wantErr)
			continue
		}
		if !tc.wantErr {
			// Use a small epsilon for floating point comparison
			epsilon := 0.0001
			if got < tc.want-epsilon || got > tc.want+epsilon {
				t.Errorf("getFloatValue(%v, %d) = %f, want %f", tc.value, tc.index, got, tc.want)
			}
		}
	}
}

// TestGetBytesValue tests the getBytesValue function
func TestGetBytesValue(t *testing.T) {
	testCases := []struct {
		value   interface{}
		want    []byte
		wantErr bool
	}{
		{[]byte{1, 2, 3}, []byte{1, 2, 3}, false},
		{"abc", []byte("abc"), false},
		{[]rune("abc"), []byte("abc"), false},
		{123, nil, true},
		{struct{}{}, nil, true},
	}

	for _, tc := range testCases {
		got, err := getBytesValue(tc.value)
		if (err != nil) != tc.wantErr {
			t.Errorf("getBytesValue(%v) error = %v, wantErr %v", tc.value, err, tc.wantErr)
			continue
		}
		if !tc.wantErr && !bytes.Equal(got, tc.want) {
			t.Errorf("getBytesValue(%v) = %v, want %v", tc.value, got, tc.want)
		}
	}
}

// TestGetStringValue tests the getStringValue function
func TestGetStringValue(t *testing.T) {
	testCases := []struct {
		value   interface{}
		want    string
		wantErr bool
	}{
		{"abc", "abc", false},
		{[]byte("abc"), "abc", false},
		{[]rune("abc"), "abc", false},
		{123, "123", false},
		{int8(123), "123", false},
		{int16(123), "123", false},
		{int32(123), "123", false},
		{int64(123), "123", false},
		{uint(123), "123", false},
		{uint8(123), "123", false},
		{uint16(123), "123", false},
		{uint32(123), "123", false},
		{uint64(123), "123", false},
		{float32(123.45), "123.45", false},
		{float64(123.45), "123.45", false},
		{struct{}{}, "", true},
	}

	for _, tc := range testCases {
		got, err := getStringValue(tc.value)
		if (err != nil) != tc.wantErr {
			t.Errorf("getStringValue(%v) error = %v, wantErr %v", tc.value, err, tc.wantErr)
			continue
		}
		if !tc.wantErr && got != tc.want {
			t.Errorf("getStringValue(%v) = %q, want %q", tc.value, got, tc.want)
		}
	}
}
