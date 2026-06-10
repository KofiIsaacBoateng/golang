package main

import (
	"bytes"
	"distributed-fs/p2p"
	"log"
	"time"
)

func makeServer (listeningAddr string, nodes ...string) *FileServer {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr: listeningAddr,
		ShakeHands: p2p.NOPShakeHands,
		Decoder: p2p.DefaultDecoder{},
		// OnPeer: OnPeer,
	}
	tcpTransport := p2p.NewTCPTransport(tcpOpts);

	fileServerOpts := FileServerOpts{
		StorageRoot: listeningAddr + "_store",
		PathTransformerFunc: CASPathTransformFunc,
		Transport: tcpTransport,
		BootStrapNodes: nodes,
	}

	s := NewFileServer(fileServerOpts);

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":3000");
	s2 := makeServer(":4000", ":3000")

	go func(){
		log.Fatal(s1.Start())
	}()
	time.Sleep(1 * time.Second)


	go s2.Start()
	time.Sleep(1 * time.Second)
	data := bytes.NewReader([]byte("This is a test file data in store."))
	s2.StoreData("keytodatainstore", data)

	select{}
}