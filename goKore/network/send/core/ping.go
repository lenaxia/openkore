package core

// SendPing sends a ping packet to the server.
// This corresponds to packet 0B1C (PACKET_CZ_PING) in the original Perl implementation.
func (bs *BaseSend) SendPing() error {
	// Send the ping packet with an empty map of arguments
	return bs.SendPacket("ping", map[string]interface{}{})
}
