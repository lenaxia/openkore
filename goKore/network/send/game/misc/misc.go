// Package misc provides miscellaneous packet sending functionality.
package misc

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// MiscManager handles miscellaneous packet sending.
type MiscManager struct {
	// Base send implementation
	baseSend core.Send
}

// NewMiscManager creates a new miscellaneous manager.
func NewMiscManager(baseSend core.Send) *MiscManager {
	return &MiscManager{
		baseSend: baseSend,
	}
}

// GetManagerName returns the name of the manager.
// This implements the ManagerProvider interface.
func (mm *MiscManager) GetManagerName() string {
	return "MiscManager"
}

// SendReqRemainTime sends a request for remaining time.
// This is equivalent to the sendReqRemainTime function in Send.pm.
func (mm *MiscManager) SendReqRemainTime() error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("request_remain_time")
	if !exists {
		return fmt.Errorf("request_remain_time packet ID not found")
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, nil)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendClientVersion sends the client version to the server.
// This is equivalent to the sendClientVersion function in Send.pm.
func (mm *MiscManager) SendClientVersion(version int) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("client_version")
	if !exists {
		return fmt.Errorf("client_version packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"version": version,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendIgnoreAll sends an ignore all request.
// This is equivalent to the sendIgnoreAll function in Send.pm.
func (mm *MiscManager) SendIgnoreAll(flag int) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("ignore_all")
	if !exists {
		return fmt.Errorf("ignore_all packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"flag": flag,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendIgnorePlayer sends an ignore player request.
// This is equivalent to the sendIgnorePlayer function in Send.pm.
func (mm *MiscManager) SendIgnorePlayer(name string, flag int) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("ignore_player")
	if !exists {
		return fmt.Errorf("ignore_player packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"name": name,
		"flag": flag,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendIgnoreList sends a request for the ignore list.
// This is equivalent to the sendIgnoreList function in Send.pm.
func (mm *MiscManager) SendIgnoreList() error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("ignore_list")
	if !exists {
		return fmt.Errorf("ignore_list packet ID not found")
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, nil)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendTokenToServer sends a token to the server.
// This is equivalent to the sendTokenToServer function in Send.pm.
func (mm *MiscManager) SendTokenToServer(username string, password string, masterVersion uint32, version uint32, token string, length uint16, otpIP string, otpPort uint16) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("token_login")
	if !exists {
		return fmt.Errorf("token_to_server packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"username":       username,
		"password":       password,
		"master_version": masterVersion,
		"version":        version,
		"token":          token,
		"length":         length,
		"otp_ip":         otpIP,
		"otp_port":       otpPort,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// EncryptPassword encrypts a password.
// This is equivalent to the encryptPassword function in Send.pm.
func (mm *MiscManager) EncryptPassword(password string) string {
	// This is a placeholder implementation
	// In a real implementation, this would encrypt the password
	return password
}

// SendBlockingPlayerCancel sends a blocking player cancel request.
// This is equivalent to the sendBlockingPlayerCancel function in Send.pm.
func (mm *MiscManager) SendBlockingPlayerCancel() error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("blocking_play_cancel")
	if !exists {
		return fmt.Errorf("blocking_player_cancel packet ID not found")
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, nil)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendRecallSso sends a recall SSO request.
// This is equivalent to the sendRecallSso function in Send.pm.
func (mm *MiscManager) SendRecallSso(accountID uint32) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("recall_sso")
	if !exists {
		return fmt.Errorf("recall_sso packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": accountID,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendRemoveAidSso sends a remove AID SSO request.
// This is equivalent to the sendRemoveAidSso function in Send.pm.
func (mm *MiscManager) SendRemoveAidSso(accountID uint32) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("remove_aid_sso")
	if !exists {
		return fmt.Errorf("remove_aid_sso packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"ID": accountID,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendFeelSaveOk sends a feel save ok request.
// This is equivalent to the sendFeelSaveOk function in Send.pm.
func (mm *MiscManager) SendFeelSaveOk(flag uint8) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("starplace_agree")
	if !exists {
		return fmt.Errorf("feel_save_ok packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"flag": flag,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}

// SendReplySyncRequestEx sends a reply sync request ex.
// This is equivalent to the sendReplySyncRequestEx function in Send.pm.
func (mm *MiscManager) SendReplySyncRequestEx(syncID uint16) error {
	// Get the packet ID
	packetID, exists := mm.baseSend.GetPacketID("sync_request_ex")
	if !exists {
		return fmt.Errorf("sync_request_ex packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"syncID": syncID,
	}

	// Construct and send the packet
	packet, err := mm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return mm.baseSend.SendToServer(packet)
}
