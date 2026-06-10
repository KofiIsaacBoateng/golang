package p2p

// RPC holds any arbitrary data that is being
// sent over tcp between two nodes within a network.
type RPC struct {
	Payload []byte
	From string
}
