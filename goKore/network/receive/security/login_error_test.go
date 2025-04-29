package security

import (
	"testing"

	"github.com/lenaxia/goKore/network/receive/core"
)

func TestGetLoginErrorMessageExtended(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	testCases := []struct {
		code int
		want string
	}{
		// Basic error codes
		{0, "Account name doesn't exist"},
		{1, "Password Error"},
		{3, "The server has denied your connection"},
		{4, "Critical Error: Your account has been blocked"},
		{5, "Connect failed, something is wrong with the login settings"},
		{6, "The server is temporarily blocking your connection"},
		{11, "Critical Error: Your account has been blocked"},

		// Extended error codes
		{0x69, "Please dial to activate the login procedure"},
		{0x6a, "Mobile Authentication: Max number of simultaneous IP addresses reached"},
		{0xa, "Account email address not confirmed"},
		{0x1452, "The server is blocking connection from this user"},
		{0x1453, "The server is blocking connections from your country"},
		{0x1456, "The server is blocking your connection due to billing issues"},
		{0x1458, "The server is blocking your connection due to billing issues"},
		{0x1459, "The server demands a password change for this account"},
		{0x14B5, "Account doesn't have access to Premium Server"},
		{0x16, "Your connection is currently delayed. You can connect again later"},
		{0xf3, "Your connection was refused due to expired Token"},

		// Unknown error code
		{999, "The server has denied your connection for unknown reason (999)"},
	}

	for _, tc := range testCases {
		got := manager.getLoginErrorMessage(tc.code)
		if got != tc.want {
			t.Errorf("getLoginErrorMessage(%d) = %s, want %s", tc.code, got, tc.want)
		}
	}
}
