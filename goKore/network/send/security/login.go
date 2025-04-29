// Package security provides security-related packet sending functionality.
package security

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
	"github.com/lenaxia/goKore/utils/crypto"
)

var (
	// ErrLoginFailed is returned when login fails.
	ErrLoginFailed = errors.New("login failed")

	// ErrInvalidCredentials is returned when credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// LoginManager handles login-related packet sending.
type LoginManager struct {
	// Base send implementation
	baseSend core.Send

	// Credentials
	username string
	password string

	// Server version
	version int

	// Master version
	masterVersion int

	// Game code
	gameCode string

	// Flag
	flag string

	// Account ID
	accountID []byte

	// Session ID
	sessionID []byte

	// Session ID2
	sessionID2 []byte

	// Account sex
	accountSex int

	// MAC address
	mac string

	// IP address
	ip string

	// Client hash
	clientHash string
}

// NewLoginManager creates a new login manager.
func NewLoginManager(baseSend core.Send) *LoginManager {
	return &LoginManager{
		baseSend:      baseSend,
		version:       23,
		masterVersion: 1,
		gameCode:      "0011", // kRO Ragnarok game code
		flag:          "G000", // Maybe this say that we are connecting from client
		mac:           "111111111111",
		ip:            "192.168.0.2", // Default IP
	}
}

// SetCredentials sets the credentials for login.
func (lm *LoginManager) SetCredentials(username, password string) {
	lm.username = username
	lm.password = password
}

// SetVersion sets the client version.
func (lm *LoginManager) SetVersion(version int) {
	lm.version = version
}

// SetMasterVersion sets the master server version.
func (lm *LoginManager) SetMasterVersion(masterVersion int) {
	lm.masterVersion = masterVersion
}

// SendClientVersion sends the client version to the server.
func (lm *LoginManager) SendClientVersion(version int) error {
	// Get the packet ID
	packetID, exists := lm.baseSend.GetPacketID("client_version")
	if !exists {
		return fmt.Errorf("client_version packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"clientVersion": version,
	}

	// Construct and send the packet
	packet, err := lm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return lm.baseSend.SendToServer(packet)
}

// SetGameCode sets the game code.
func (lm *LoginManager) SetGameCode(gameCode string) {
	lm.gameCode = gameCode
}

// SetFlag sets the flag.
func (lm *LoginManager) SetFlag(flag string) {
	lm.flag = flag
}

// SetMAC sets the MAC address.
func (lm *LoginManager) SetMAC(mac string) {
	lm.mac = mac
}

// SetIP sets the IP address.
func (lm *LoginManager) SetIP(ip string) {
	lm.ip = ip
}

// SetAccountInfo sets the account information.
func (lm *LoginManager) SetAccountInfo(accountID, sessionID, sessionID2 []byte, accountSex int) {
	lm.accountID = accountID
	lm.sessionID = sessionID
	lm.sessionID2 = sessionID2
	lm.accountSex = accountSex
}

// SendMasterLogin sends a master login packet.
func (lm *LoginManager) SendMasterLogin() error {
	// Get the packet ID
	packetID, exists := lm.baseSend.GetPacketID("master_login")
	if !exists {
		return fmt.Errorf("master_login packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"version":        lm.version,
		"master_version": lm.masterVersion,
		"username":       lm.username,
		"password":       lm.password,
		"game_code":      lm.gameCode,
		"flag":           lm.flag,
	}

	// Construct and send the packet
	packet, err := lm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return lm.baseSend.SendToServer(packet)
}

// SendMasterSecureLogin sends a secure master login packet.
func (lm *LoginManager) SendMasterSecureLogin(salt []byte, loginType int, account int) error {
	// Get the packet ID
	packetID, exists := lm.baseSend.GetPacketID("master_login")
	if !exists {
		return fmt.Errorf("master_login packet ID not found")
	}

	// Create the salted MD5 hash
	saltedMD5 := lm.secureLoginHash(lm.password, salt, loginType)

	// Create the arguments
	args := map[string]interface{}{
		"version":             lm.version,
		"master_version":      lm.masterVersion,
		"username":            lm.username,
		"password_salted_md5": saltedMD5,
		"clientInfo": func() int {
			if account > 0 {
				return account - 1
			} else {
				return 0
			}
		}(),
	}

	// Construct and send the packet
	packet, err := lm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return lm.baseSend.SendToServer(packet)
}

// secureLoginHash creates a secure login hash.
func (lm *LoginManager) secureLoginHash(password string, salt []byte, loginType int) []byte {
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

// SendGameLogin sends a game login packet.
func (lm *LoginManager) SendGameLogin() error {
	// Get the packet ID
	packetID, exists := lm.baseSend.GetPacketID("game_login")
	if !exists {
		return fmt.Errorf("game_login packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"accountID":  lm.accountID,
		"sessionID":  lm.sessionID,
		"sessionID2": lm.sessionID2,
		"mac":        lm.mac,
		"accountSex": lm.accountSex,
	}

	// Construct and send the packet
	packet, err := lm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return lm.baseSend.SendToServer(packet)
}

// SendCharLogin sends a character login packet.
func (lm *LoginManager) SendCharLogin(slot int) error {
	// Get the packet ID
	packetID, exists := lm.baseSend.GetPacketID("char_login")
	if !exists {
		return fmt.Errorf("char_login packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"slot": slot,
	}

	// Construct and send the packet
	packet, err := lm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return lm.baseSend.SendToServer(packet)
}

// SendMapLogin sends a map login packet.
func (lm *LoginManager) SendMapLogin() error {
	// Get the packet ID
	packetID, exists := lm.baseSend.GetPacketID("map_login")
	if !exists {
		return fmt.Errorf("map_login packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"accountID": lm.accountID,
		"charID":    lm.accountID, // In a real implementation, this would be the character ID
		"sessionID": lm.sessionID,
		"tick":      lm.baseSend.GetTime(),
		"sex":       lm.accountSex,
	}

	// Construct and send the packet
	packet, err := lm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return lm.baseSend.SendToServer(packet)
}

// SendMapLoaded sends a map loaded packet.
func (lm *LoginManager) SendMapLoaded() error {
	// Get the packet ID
	packetID, exists := lm.baseSend.GetPacketID("map_loaded")
	if !exists {
		return fmt.Errorf("map_loaded packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{}

	// Construct and send the packet
	packet, err := lm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return lm.baseSend.SendToServer(packet)
}

// SendClientHash sends a client hash packet.
func (lm *LoginManager) SendClientHash(hash []byte) error {
	// Get the packet ID
	packetID, exists := lm.baseSend.GetPacketID("client_hash")
	if !exists {
		return fmt.Errorf("client_hash packet ID not found")
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

// SendTokenLogin sends a token login packet.
func (lm *LoginManager) SendTokenLogin(token []byte) error {
	// Get the packet ID
	packetID, exists := lm.baseSend.GetPacketID("token_login")
	if !exists {
		return fmt.Errorf("token_login packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"version":        lm.version,
		"master_version": lm.masterVersion,
		"username":       lm.username,
		"mac":            lm.mac,
		"ip":             lm.ip,
		"token":          token,
	}

	// Construct and send the packet
	packet, err := lm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return lm.baseSend.SendToServer(packet)
}

// SendRijndaelLogin sends a login packet with Rijndael encryption.
func (lm *LoginManager) SendRijndaelLogin() error {
	// Get the packet ID
	packetID, exists := lm.baseSend.GetPacketID("master_login")
	if !exists {
		return fmt.Errorf("master_login packet ID not found")
	}

	// Create the Rijndael key and chain
	key := []byte{6, 169, 33, 64, 54, 184, 161, 91, 81, 46, 3, 213, 52, 18, 0, 6, 61, 175, 186, 66, 157, 158, 180, 48}
	chain := []byte{61, 175, 186, 66, 157, 158, 180, 48, 180, 34, 218, 128, 44, 159, 172, 65, 1, 2, 4, 8, 16, 32, 128}

	// Create the Rijndael encryptor
	rijndael := crypto.NewRijndael()

	// Make the key
	rijndael.MakeKey(key, chain, 24, 24)

	// Encrypt the password
	passwordBytes := []byte(lm.password)
	if len(passwordBytes) < 24 {
		// Pad the password to 24 bytes
		paddedPassword := make([]byte, 24)
		copy(paddedPassword, passwordBytes)
		passwordBytes = paddedPassword
	} else if len(passwordBytes) > 24 {
		// Truncate the password to 24 bytes
		passwordBytes = passwordBytes[:24]
	}

	// Encrypt the password
	encryptedPassword := rijndael.Encrypt(passwordBytes, nil, 24, 0)

	// Create the arguments
	args := map[string]interface{}{
		"version":           lm.version,
		"username":          lm.username,
		"password_rijndael": encryptedPassword,
		"master_version":    lm.masterVersion,
	}

	// Construct and send the packet
	packet, err := lm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return lm.baseSend.SendToServer(packet)
}

// SendMD5Login sends a login packet with MD5 encryption.
func (lm *LoginManager) SendMD5Login() error {
	// Get the packet ID
	packetID, exists := lm.baseSend.GetPacketID("master_login")
	if !exists {
		return fmt.Errorf("master_login packet ID not found")
	}

	// Create the MD5 hash
	hasher := md5.New()
	hasher.Write([]byte(lm.password))
	passwordMD5 := hasher.Sum(nil)
	passwordMD5Hex := hex.EncodeToString(passwordMD5)

	// Create the arguments
	args := map[string]interface{}{
		"version":          lm.version,
		"username":         lm.username,
		"password_md5_hex": passwordMD5Hex,
		"master_version":   lm.masterVersion,
	}

	// Construct and send the packet
	packet, err := lm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return lm.baseSend.SendToServer(packet)
}
