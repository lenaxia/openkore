// Package security provides security-related packet sending functionality.
package security

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidHashType is returned when an invalid hash type is provided.
var ErrInvalidHashType = errors.New("invalid hash type")

// ReconstructClientHash reconstructs a client hash from a code string.
// The code string is expected to be a space-separated list of hex bytes.
// If the code starts with "02 04", it will be removed.
func ReconstructClientHash(code string) ([]byte, error) {
	// Remove "02 04" prefix if present
	code = strings.TrimPrefix(code, "02 04 ")

	// Split the code into individual hex bytes
	hexBytes := strings.Split(code, " ")

	// Convert each hex byte to a byte
	result := make([]byte, len(hexBytes))
	for i, hexByte := range hexBytes {
		b, err := hex.DecodeString(hexByte)
		if err != nil {
			return nil, fmt.Errorf("failed to decode hex byte %s: %w", hexByte, err)
		}
		result[i] = b[0]
	}

	return result, nil
}

// ReconstructClientHashByType reconstructs a client hash based on a type.
// The type can be 1, 2, 3, 4, or 5, each corresponding to a different hash.
func ReconstructClientHashByType(hashType int) ([]byte, error) {
	switch hashType {
	case 1:
		return []byte{0x7B, 0x8A, 0xA8, 0x90, 0x2F, 0xD8, 0xE8, 0x30, 0xF8, 0xA5, 0x25, 0x7A, 0x0D, 0x3B, 0xCE, 0x52}, nil
	case 2:
		return []byte{0x27, 0x6A, 0x2C, 0xCE, 0xAF, 0x88, 0x01, 0x87, 0xCB, 0xB1, 0xFC, 0xD5, 0x90, 0xC4, 0xED, 0xD2}, nil
	case 3:
		return []byte{0x42, 0x00, 0xB0, 0xCA, 0x10, 0x49, 0x3D, 0x89, 0x49, 0x42, 0x82, 0x57, 0xB1, 0x68, 0x5B, 0x85}, nil
	case 4:
		return []byte{0x22, 0x37, 0xD7, 0xFC, 0x8E, 0x9B, 0x05, 0x79, 0x60, 0xAE, 0x02, 0x33, 0x6D, 0x0D, 0x82, 0xC6}, nil
	case 5:
		return []byte{0xC7, 0x0A, 0x94, 0xC2, 0x7A, 0xCC, 0x38, 0x9A, 0x47, 0xF5, 0x54, 0x39, 0x7C, 0xA4, 0xD0, 0x39}, nil
	default:
		return nil, ErrInvalidHashType
	}
}

// SetClientHash sets the client hash for login.
func (lm *LoginManager) SetClientHash(clientHash string) {
	lm.clientHash = clientHash
}

// SendClientMD5Hash sends the client MD5 hash to the server.
func (lm *LoginManager) SendClientMD5Hash() error {
	// Get the packet ID
	packetID, exists := lm.baseSend.GetPacketID("client_hash")
	if !exists {
		return fmt.Errorf("client_hash packet ID not found")
	}

	// Convert the client hash from hex to bytes
	hash, err := hex.DecodeString(lm.clientHash)
	if err != nil {
		return fmt.Errorf("failed to decode client hash: %w", err)
	}

	// Create the arguments
	args := map[string]interface{}{
		"hash": hash,
	}

	// Construct and send the packet
	packet, err := lm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return lm.baseSend.SendToServer(packet)
}
