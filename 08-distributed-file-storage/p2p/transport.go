package p2p

import "net"

// interface that represents a remote network node
type Peer interface {
	Send([]byte) error 
	RemoteAddr() net.Addr
	Close() error
}


// communation layer between nodes
type Transport interface {
	Dial(addr string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}


