// Package security provides security-related packet sending functionality.
package security

import (
	"crypto/md5"
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// TokenManager handles token-based authentication packet sending.
type TokenManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewTokenManager creates a new token manager.
func NewTokenManager(baseSend core.Send) *TokenManager {
	return &TokenManager{
		baseSend: baseSend,
	}
}

// SendTokenLogin sends a token-based login packet.
func (tm *TokenManager) SendTokenLogin(username string, token []byte, mac, ip string, version, masterVersion int) error {
	// Get the packet ID
	packetID, exists := tm.baseSend.GetPacketID("token_login")
	if !exists {
		return fmt.Errorf("token_login packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"version":        version,
		"master_version": masterVersion,
		"username":       username,
		"mac":            mac,
		"ip":             ip,
		"token":          token,
	}

	// Construct and send the packet
	packet, err := tm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return tm.baseSend.SendToServer(packet)
}

// SendSecureLogin sends a secure login packet with salted password.
func (tm *TokenManager) SendSecureLogin(username, password string, salt []byte, loginType int) error {
	// Get the packet ID
	packetID, exists := tm.baseSend.GetPacketID("secure_login")
	if !exists {
		return fmt.Errorf("secure_login packet ID not found")
	}

	// Create the salted MD5 hash
	passwordHash := tm.secureLoginHash(password, salt, loginType)

	// Create the arguments
	args := map[string]interface{}{
		"username":      username,
		"password_hash": passwordHash,
		"login_type":    loginType,
	}

	// Construct and send the packet
	packet, err := tm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return tm.baseSend.SendToServer(packet)
}

// secureLoginHash creates a secure login hash.
func (tm *TokenManager) secureLoginHash(password string, salt []byte, loginType int) []byte {
	// Create the MD5 hasher
	hasher := md5.New()

	// Add the salt and password based on the login type
	if loginType%2 == 1 {
		hasher.Write(salt)
		hasher.Write([]byte(password))
	} else {
		hasher.Write([]byte(password))
		hasher.Write(salt)
	}

	// Return the digest
	return hasher.Sum(nil)
}
