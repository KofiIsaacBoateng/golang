package main

import (
	"distributed-fs/p2p"
	"fmt"
)

func OnPeer (peer p2p.Peer) error {
	return nil;
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr: ":5000",
		ShakeHands: p2p.NOPShakeHands,
		Decoder: p2p.DefaultDecoder{},
		OnPeer: OnPeer,
	}
	tcpTransport := p2p.NewTCPTransport(tcpOpts);

	if err := tcpTransport.ListenAndAccept(); err != nil {
		fmt.Println(err)
	}


	go func () {
		for{
			msg := <-tcpTransport.Consume();
			fmt.Printf("Message received from peer: %+v\n", msg)
		}
	}()

	select {}
}