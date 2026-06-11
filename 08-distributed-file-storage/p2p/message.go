package p2p

const (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

// RPC holds any arbitrary data that is being
// sent over tcp between two nodes within a network.
type RPC struct {
	Payload []byte
	From    string
	Stream  bool
}
