// Package login provides login-related packet constructions and handlers
package login

import (
	"github.com/lenaxia/goKore/network/send/types"
)

// RegisterHandlers registers all login-related handlers
func RegisterHandlers(send types.Send) {
	// Register login packet handlers
	send.RegisterHandler("login_request", constructLoginRequest)
	send.RegisterHandler("login_otp", constructLoginOTP)
	send.RegisterHandler("login_restart", constructLoginRestart)
	// More login handlers...
}

// constructLoginRequest constructs a login request packet
func constructLoginRequest(args map[string]interface{}) ([]byte, error) {
	// Implementation for login_request
	// This is a placeholder - real implementation would use the args to construct the packet
	return []byte{0x64, 0x00, 0x01, 0x02, 0x03, 0x04}, nil
}

// constructLoginOTP constructs a login OTP packet
func constructLoginOTP(args map[string]interface{}) ([]byte, error) {
	// Implementation for login_otp
	// This is a placeholder - real implementation would use the args to construct the packet
	return []byte{0x66, 0x01, 0x01, 0x02, 0x03, 0x04}, nil
}

// constructLoginRestart constructs a login restart packet
func constructLoginRestart(args map[string]interface{}) ([]byte, error) {
	// Implementation for login_restart
	// This is a placeholder - real implementation would use the args to construct the packet
	return []byte{0xB2, 0x00, 0x01}, nil
}

// GetPacketConstructions returns login-related packet constructions
func GetPacketConstructions() map[string]types.PacketConstruction {
	return map[string]types.PacketConstruction{
		"0064": {
			ID:         "0064",
			Name:       "login_request",
			Format:     "v a24 a24 C",
			FieldNames: []string{"version", "username", "password", "clienttype"},
		},
		"0066": {
			ID:         "0066",
			Name:       "login_otp",
			Format:     "v a24 a24 C C a6",
			FieldNames: []string{"version", "username", "password", "clienttype", "state", "otp"},
		},
		"00B2": {
			ID:         "00B2",
			Name:       "login_restart",
			Format:     "C",
			FieldNames: []string{"type"},
		},
		// More packet constructions...
	}
}
