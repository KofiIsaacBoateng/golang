package p2p

// interface that represents a remote network node
type Peer interface {
	Close() error
}


// communation layer between nodes
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}


