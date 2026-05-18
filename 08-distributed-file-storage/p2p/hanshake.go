package p2p

// Handshake func
type HandShakeFunc func(Peer) error

func NOPShakeHands(peer Peer) error {
	return nil
}