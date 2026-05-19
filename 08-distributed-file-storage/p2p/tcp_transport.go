package p2p

import (
	"errors"
	"fmt"
	"net"
)

type TCPPeer struct {
	net.Conn
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
		Conn: conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Send(bytes []byte) error {
	_, err := p.Conn.Write(bytes)
	return err
}


func NewTCPTransport (opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		tcpch: make(chan RPC),
	}
}


// Consume implements the transport interface
func (t *TCPTransport) Consume() <-chan RPC {
	return t.tcpch
}

// Close implements the transport interface
func (t *TCPTransport) Close() error {
	return t.listener.Close();
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

	fmt.Printf("Server is listening on PORT: %s\n", t.ListenAddr);

	return nil

}


func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept();
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			fmt.Printf("TCP accept error: %v\n", err)
		}

		go t.handleConn(conn, true)
	}
}

func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr);
	if(err != nil) {
		return err
	}

	go t.handleConn(conn, false)

	return nil
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("Dropping peer connection: %s\n", err)
		conn.Close()
	}()


	peer := NewTCPPeer(conn, outbound);

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

		if errors.Is(err, net.ErrClosed){
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