package main

import (
	// "bytes"
	"distributed-fs/p2p"
	"io"
	"log"
	"time"
)

func makeServer (listenAddr string, nodes ...string) *FileServer {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr: listenAddr,
		ShakeHands: p2p.NOPShakeHands,
		Decoder: p2p.DefaultDecoder{},
		// OnPeer: OnPeer,
	}
	tcpTransport := p2p.NewTCPTransport(tcpOpts);

	fileServerOpts := FileServerOpts{
		StorageRoot: listenAddr[1:] + "_store",
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
	// data := bytes.NewReader([]byte("This is a test file data in store."))
	// s2.Store("keytodatainstore", data)

	key1 := "keytodatainstore"
	// key2 := "doesn't exist"
	r, err := s2.Get(key1);
	if err != nil {
		log.Fatal(err)
		
	}

	b, err := io.ReadAll(r);
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Read these bytes: %s\n", string(b))

	select{}
}