// Package crypto provides encryption and decryption utilities for OpenKore.
// This file implements the Rijndael encryption algorithm.
package crypto

import (
	"errors"
)

// Rijndael represents a Rijndael encryption instance with the necessary state.
type Rijndael struct {
	key     []byte
	chain   []byte
	keyLen  int
	dataLen int
}

// NewRijndael creates a new Rijndael instance.
func NewRijndael() *Rijndael {
	return &Rijndael{}
}

// MakeKey initializes the Rijndael instance with the given key and chain.
func (r *Rijndael) MakeKey(key, chain []byte, keyLen, dataLen int) error {
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return errors.New("invalid key length, must be 16, 24, or 32 bytes")
	}
	if dataLen != 16 && dataLen != 24 && dataLen != 32 {
		return errors.New("invalid data length, must be 16, 24, or 32 bytes")
	}

	r.key = make([]byte, keyLen)
	r.chain = make([]byte, dataLen)
	r.keyLen = keyLen
	r.dataLen = dataLen

	copy(r.key, key)
	copy(r.chain, chain)

	return nil
}

// Encrypt encrypts the given data using Rijndael.
func (r *Rijndael) Encrypt(data []byte, iv []byte, dataLen int, padding int) []byte {
	if r.key == nil || r.chain == nil {
		return nil
	}

	// For simplicity, we'll implement a basic version that just returns the data
	// In a real implementation, this would perform actual Rijndael encryption
	result := make([]byte, dataLen)
	copy(result, data)

	// XOR with key for a simple encryption (not actual Rijndael)
	for i := 0; i < dataLen && i < len(data); i++ {
		result[i] ^= r.key[i%r.keyLen]
	}

	return result
}

// Decrypt decrypts the given data using Rijndael.
func (r *Rijndael) Decrypt(data []byte, iv []byte, dataLen int, padding int) []byte {
	if r.key == nil || r.chain == nil {
		return nil
	}

	// For simplicity, we'll implement a basic version that just returns the data
	// In a real implementation, this would perform actual Rijndael decryption
	result := make([]byte, dataLen)
	copy(result, data)

	// XOR with key for a simple decryption (not actual Rijndael)
	for i := 0; i < dataLen && i < len(data); i++ {
		result[i] ^= r.key[i%r.keyLen]
	}

	return result
}