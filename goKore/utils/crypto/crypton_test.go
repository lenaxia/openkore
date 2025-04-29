package crypto

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// TestCryptonGenerateTable tests that the S-box and S-tab tables are generated correctly
func TestCryptonGenerateTable(t *testing.T) {
	// Reset the tableGenerated flag to force table generation
	tableGenerated = false

	// Generate the tables
	generateTable()

	// Check that the tables have been populated
	// Instead of checking for non-zero values (which might be valid),
	// we'll check that the inverse relationship holds for the S-box
	for i := 0; i < 256; i++ {
		if sBox[1][sBox[0][i]] != byte(i) {
			t.Errorf("S-box inverse relationship failed for index %d: sBox[1][sBox[0][%d]] = %d, expected %d",
				i, i, sBox[1][sBox[0][i]], i)
		}
	}

	// Check that the S-tab tables have been populated
	// We'll check that at least some entries are non-zero
	nonZeroCount := 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 256; j++ {
			if sTab[i][j] != 0 {
				nonZeroCount++
			}
		}
	}

	if nonZeroCount < 100 {
		t.Errorf("Expected at least 100 non-zero entries in S-tab tables, got %d", nonZeroCount)
	}
}

// TestCryptonNew tests the creation of a new Crypton instance
func TestCryptonNew(t *testing.T) {
	// Create a test key
	key := []byte("0123456789ABCDEF0123456789ABCDEF")
	keyLen := len(key) * 8 // Key length in bits

	// Create a new Crypton instance
	c := New(key, keyLen)

	// Check that the instance was created
	if c == nil {
		t.Error("New returned nil")
	}

	// Check that the key was set up
	if len(c.eKey) != 52 {
		t.Errorf("eKey length should be 52, got %d", len(c.eKey))
	}

	// Check that at least some key values are non-zero
	nonZeroCount := 0
	for _, v := range c.eKey {
		if v != 0 {
			nonZeroCount++
		}
	}

	if nonZeroCount == 0 {
		t.Error("All eKey values are 0")
	}
}

// TestCryptonEncrypt tests the encryption of a block
func TestCryptonEncrypt(t *testing.T) {
	// Test vectors
	testCases := []struct {
		key       string
		keyLen    int
		plaintext string
		expected  string
	}{
		{
			// Test case with a 128-bit key
			key:       "000102030405060708090A0B0C0D0E0F",
			keyLen:    128,
			plaintext: "00112233445566778899AABBCCDDEEFF",
			// This expected value would need to be generated from the Perl implementation
			// or from a known-good implementation of Crypton
			expected: "", // To be filled in with actual expected value
		},
	}

	for i, tc := range testCases {
		// Decode hex strings
		key, err := hex.DecodeString(tc.key)
		if err != nil {
			t.Errorf("Case %d: Failed to decode key: %v", i, err)
			continue
		}

		plaintext, err := hex.DecodeString(tc.plaintext)
		if err != nil {
			t.Errorf("Case %d: Failed to decode plaintext: %v", i, err)
			continue
		}

		// Create a new Crypton instance
		c := New(key, tc.keyLen)

		// Encrypt the plaintext
		ciphertext := c.Encrypt(plaintext)

		// If we have an expected value, check it
		if tc.expected != "" {
			expected, err := hex.DecodeString(tc.expected)
			if err != nil {
				t.Errorf("Case %d: Failed to decode expected: %v", i, err)
				continue
			}

			if !bytes.Equal(ciphertext, expected) {
				t.Errorf("Case %d: Expected %x, got %x", i, expected, ciphertext)
			}
		} else {
			// If we don't have an expected value, just print the result
			t.Logf("Case %d: Ciphertext: %x", i, ciphertext)
		}
	}
}

// TestHelperFunctions tests the helper functions used by Crypton
func TestHelperFunctions(t *testing.T) {
	// Test byteAt
	data := uint32(0x12345678)
	if byteAt(data, 0) != 0x78 {
		t.Errorf("byteAt(0x12345678, 0) should be 0x78, got %x", byteAt(data, 0))
	}
	if byteAt(data, 1) != 0x56 {
		t.Errorf("byteAt(0x12345678, 1) should be 0x56, got %x", byteAt(data, 1))
	}
	if byteAt(data, 2) != 0x34 {
		t.Errorf("byteAt(0x12345678, 2) should be 0x34, got %x", byteAt(data, 2))
	}
	if byteAt(data, 3) != 0x12 {
		t.Errorf("byteAt(0x12345678, 3) should be 0x12, got %x", byteAt(data, 3))
	}

	// Test rotl
	if rotl(0x12345678, 8) != 0x34567812 {
		t.Errorf("rotl(0x12345678, 8) should be 0x34567812, got %x", rotl(0x12345678, 8))
	}
	if rotl(0x12345678, 16) != 0x56781234 {
		t.Errorf("rotl(0x12345678, 16) should be 0x56781234, got %x", rotl(0x12345678, 16))
	}
	if rotl(0x12345678, 24) != 0x78123456 {
		t.Errorf("rotl(0x12345678, 24) should be 0x78123456, got %x", rotl(0x12345678, 24))
	}

	// Test pi
	// Initialize the tables first
	if !tableGenerated {
		generateTable()
		tableGenerated = true
	}

	// Create a test array
	testArray := []uint32{0x12345678, 0x9ABCDEF0, 0x13579BDF, 0x2468ACE0}

	// Test pi with different parameters
	result := pi(testArray, 0, 0, 1, 2, 3)
	if result == 0 {
		t.Error("pi returned 0")
	}

	// Test phiN
	result = phiN(0x12345678, 0, 1, 2, 3)
	if result == 0 {
		t.Error("phiN returned 0")
	}
}
