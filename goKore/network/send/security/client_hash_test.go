package security

import (
	"bytes"
	"testing"
)

// TestReconstructClientHash tests the ReconstructClientHash function
func TestReconstructClientHash(t *testing.T) {
	// Test with code
	hash, err := ReconstructClientHash("02 04 7B 8A A8 90 2F D8 E8 30 F8 A5 25 7A 0D 3B CE 52")
	if err != nil {
		t.Fatalf("ReconstructClientHash() returned error: %v", err)
	}

	expected := []byte{0x7B, 0x8A, 0xA8, 0x90, 0x2F, 0xD8, 0xE8, 0x30, 0xF8, 0xA5, 0x25, 0x7A, 0x0D, 0x3B, 0xCE, 0x52}
	if !bytes.Equal(hash, expected) {
		t.Errorf("ReconstructClientHash() = %v, want %v", hash, expected)
	}

	// Test with type 1
	hash, err = ReconstructClientHashByType(1)
	if err != nil {
		t.Fatalf("ReconstructClientHashByType() returned error: %v", err)
	}

	expected = []byte{0x7B, 0x8A, 0xA8, 0x90, 0x2F, 0xD8, 0xE8, 0x30, 0xF8, 0xA5, 0x25, 0x7A, 0x0D, 0x3B, 0xCE, 0x52}
	if !bytes.Equal(hash, expected) {
		t.Errorf("ReconstructClientHashByType() = %v, want %v", hash, expected)
	}

	// Test with type 2
	hash, err = ReconstructClientHashByType(2)
	if err != nil {
		t.Fatalf("ReconstructClientHashByType() returned error: %v", err)
	}

	expected = []byte{0x27, 0x6A, 0x2C, 0xCE, 0xAF, 0x88, 0x01, 0x87, 0xCB, 0xB1, 0xFC, 0xD5, 0x90, 0xC4, 0xED, 0xD2}
	if !bytes.Equal(hash, expected) {
		t.Errorf("ReconstructClientHashByType() = %v, want %v", hash, expected)
	}

	// Test with type 3
	hash, err = ReconstructClientHashByType(3)
	if err != nil {
		t.Fatalf("ReconstructClientHashByType() returned error: %v", err)
	}

	expected = []byte{0x42, 0x00, 0xB0, 0xCA, 0x10, 0x49, 0x3D, 0x89, 0x49, 0x42, 0x82, 0x57, 0xB1, 0x68, 0x5B, 0x85}
	if !bytes.Equal(hash, expected) {
		t.Errorf("ReconstructClientHashByType() = %v, want %v", hash, expected)
	}

	// Test with type 4
	hash, err = ReconstructClientHashByType(4)
	if err != nil {
		t.Fatalf("ReconstructClientHashByType() returned error: %v", err)
	}

	expected = []byte{0x22, 0x37, 0xD7, 0xFC, 0x8E, 0x9B, 0x05, 0x79, 0x60, 0xAE, 0x02, 0x33, 0x6D, 0x0D, 0x82, 0xC6}
	if !bytes.Equal(hash, expected) {
		t.Errorf("ReconstructClientHashByType() = %v, want %v", hash, expected)
	}

	// Test with type 5
	hash, err = ReconstructClientHashByType(5)
	if err != nil {
		t.Fatalf("ReconstructClientHashByType() returned error: %v", err)
	}

	expected = []byte{0xC7, 0x0A, 0x94, 0xC2, 0x7A, 0xCC, 0x38, 0x9A, 0x47, 0xF5, 0x54, 0x39, 0x7C, 0xA4, 0xD0, 0x39}
	if !bytes.Equal(hash, expected) {
		t.Errorf("ReconstructClientHashByType() = %v, want %v", hash, expected)
	}

	// Test with invalid type
	_, err = ReconstructClientHashByType(6)
	if err == nil {
		t.Error("ReconstructClientHashByType() with invalid type should return error")
	}
}

// TestSendClientMD5Hash tests the SendClientMD5Hash method
func TestSendClientMD5Hash(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.RegisterPacketHandler("0204", "client_hash", "", nil, nil)

	lm := NewLoginManager(mockSend)
	lm.SetClientHash("82d12c914f5ad48fd96fcf7ef4cc492d")

	err := lm.SendClientMD5Hash()
	if err != nil {
		t.Fatalf("SendClientMD5Hash() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0204"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	// Check that the hash was correctly converted from hex to bytes
	expectedHash := []byte{
		0x82, 0xd1, 0x2c, 0x91, 0x4f, 0x5a, 0xd4, 0x8f,
		0xd9, 0x6f, 0xcf, 0x7e, 0xf4, 0xcc, 0x49, 0x2d,
	}

	if !bytes.Equal(args["hash"].([]byte), expectedHash) {
		t.Errorf("args[\"hash\"] = %v, want %v", args["hash"], expectedHash)
	}
}
