package main

import (
	"bytes"
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
		EncKey: newEncryptionKey(),
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
	s3 := makeServer(":5000", ":4000", ":3000");

	go s1.Start()
	time.Sleep(1 * time.Second)

	go s2.Start()
	time.Sleep(1 * time.Second)

	go s3.Start()
	time.Sleep(1 * time.Second)

	key := "coolPicture.png"
	data := bytes.NewReader([]byte("This is a cool picture in png."))
	s3.Store("coolPicture.png", data)

	 if err := s3.store.Delete(key); err != nil { // delete file locally (just to test if we could fetch over the network)
		log.Fatal(err)
	}
	r, err := s3.Get(key);
	if err != nil {
		log.Fatal(err)
		
	}

	b, err := io.ReadAll(r);
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Read these bytes: %s\n", string(b))

}