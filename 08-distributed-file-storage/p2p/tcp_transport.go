package p2p

import (
	"fmt"
	"net"
)

type TCPPeer struct {
	conn net.Conn
	outbound bool // connection coming from outside this server

}

type TCPTransportOpts struct {
	ListenAddr string
	ShakeHands HandShakeFunc
	Decoder Decoder
	OnPeer func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener

	//consume channel
	tcpch chan RPC
}


func NewTCPPeer (conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn,
		outbound,
	}
}

// close implements the Peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}


func NewTCPTransport (opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
			TCPTransportOpts: opts,
			tcpch: make(chan RPC),
	}
}


func (t *TCPTransport) Consume() <-chan RPC {
	return t.tcpch
}



func (t *TCPTransport) ListenAndAccept() error {
	var err error

	// initiate listener over a tcp connection
	t.listener, err = net.Listen("tcp", t.ListenAddr);
	if err != nil {
		return err
	}

	// start accept loop (accept any tcp connection coming over the specified addr)
	go t.startAcceptLoop()

	return nil

}


func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept();
		if err != nil {
			fmt.Printf("TCP accept error: %v\n", err)
		}

		fmt.Printf("New incoming connection from %+v\n", conn.RemoteAddr())
		go t.handleConn(conn)
	}
}


func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error

	defer func() {
		fmt.Printf("Dropping peer connection: %s\n", err)
		conn.Close()
	}()


	peer := NewTCPPeer(conn, true);

	// shakehands
	if err = t.ShakeHands(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// read loop
	var rpc RPC
	for {
		err = t.Decoder.Decode(conn, &rpc); 

		if err == net.ErrClosed {
			return
		}

		if err != nil {
			fmt.Printf("TCP Error: %v\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()

		t.tcpch <- rpc
	}
}